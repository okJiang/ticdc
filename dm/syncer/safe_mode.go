// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package syncer

import (
	"time"

	"github.com/pingcap/failpoint"
	"go.uber.org/zap"

	"github.com/pingcap/tiflow/dm/dm/unit"
	"github.com/pingcap/tiflow/dm/pkg/binlog"
	tcontext "github.com/pingcap/tiflow/dm/pkg/context"
	"github.com/pingcap/tiflow/dm/pkg/terror"
)

func (s *Syncer) enableSafeModeByTaskCliArgs(tctx *tcontext.Context) {
	//nolint:errcheck
	s.safeMode.Add(tctx, 1)
	s.tctx.L().Info("enable safe-mode because of task cli args")
}

func (s *Syncer) enableSafeModeInitializationPhase(tctx *tcontext.Context) {
	var err error
	defer func() {
		if err != nil {
			// send error to the fatal chan to interrupt the process
			s.runFatalChan <- unit.NewProcessError(err)
		}
		if !s.safeMode.Enable() {
			s.tctx.L().Warn("enable safe-mode failed")
		}
	}()

	s.safeMode.Reset(tctx) // in initialization phase, reset first

	// cliArgs has higher priority than config
	if s.cliArgs != nil && s.cliArgs.SafeModeDuration != "" {
		s.enableSafeModeByTaskCliArgs(tctx)
		return
	}

	if s.cfg.SafeMode {
		//nolint:errcheck
		s.safeMode.Add(tctx, 1) // add 1 but has no corresponding -1, so keeps enabled
		s.tctx.L().Info("enable safe-mode by config")
		return
	}
	exitPoint := s.checkpoint.SafeModeExitPoint()
	failpoint.Inject("SafeModeDurationIsZero", func() {
		exitPoint = nil
	})
	if exitPoint != nil {
		s.tctx.L().Info("compare exitPoint and beginLocation", zap.Stringer("exit point", *exitPoint), zap.Stringer("init executed location", s.beginLocation))
		if binlog.CompareLocation(*exitPoint, *s.beginLocation, s.cfg.EnableGTID) > 0 {
			//nolint:errcheck
			s.safeMode.Add(tctx, 1) // enable and will revert after pass SafeModeExitLoc
			s.tctx.L().Info("enable safe-mode for safe mode exit point, will exit at", zap.Stringer("location", *exitPoint))
		} else {
			s.tctx.L().Info("disable safe-mode because initExecutedLoc equal safeModeExitPoint")
		}
	} else {
		initPhaseSeconds := s.cfg.SafeModeDuration

		failpoint.Inject("SafeModeInitPhaseSeconds", func(val failpoint.Value) {
			seconds := val.(string)
			initPhaseSeconds = seconds
			s.tctx.L().Info("set initPhaseSeconds", zap.String("failpoint", "SafeModeInitPhaseSeconds"), zap.String("value", seconds))
		})
		var duration time.Duration
		if initPhaseSeconds == "" {
			duration = time.Second * time.Duration(2*s.cfg.CheckpointFlushInterval)
		} else {
			duration, err = time.ParseDuration(initPhaseSeconds)
			if err != nil {
				// send error to the fatal chan to interrupt the process
				s.runFatalChan <- unit.NewProcessError(err)
				s.tctx.L().Error("enable safe-mode failed due to duration parse failed", zap.String("duration", initPhaseSeconds))
				return
			}
		}
		s.tctx.L().Info("enable safe-mode because of task initialization", zap.Duration("duration", duration))

		if int64(duration) > 0 {
			//nolint:errcheck
			s.safeMode.Add(tctx, 1) // enable and will revert after 2 * CheckpointFlushInterval
			go func() {
				defer func() {
					err = s.safeMode.Add(tctx, -1)
					if err == nil && !s.safeMode.Enable() {
						s.tctx.L().Info("disable safe-mode after task initialization finished")
					}
				}()

				select {
				case <-tctx.Context().Done():
				case <-time.After(duration):
				}
			}()
		} else {
			fresh, err2 := s.IsFreshTask(tctx.Ctx)
			if err2 != nil {
				err = err2
			} else if !fresh && !s.safeMode.Enable() {
				err = terror.ErrSyncerReprocessWithSafeModeFail.Generate()
			}
		}
	}
}
