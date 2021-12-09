// Copyright 2021 PingCAP, Inc.
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

package checker

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pingcap/errors"
	"github.com/pingcap/tidb-tools/pkg/dbutil"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/charset"
	"github.com/pingcap/tidb/parser/mysql"
)

// AutoIncrementKeyChecking is an identification for auto increment key checking.
const AutoIncrementKeyChecking = "auto-increment key checking"

// hold information of incompatibility option.
type incompatibilityOption struct {
	state       State
	instruction string
	errMessage  string
}

// String returns raw text of this incompatibility option.
func (o *incompatibilityOption) String() string {
	var text bytes.Buffer

	if len(o.errMessage) > 0 {
		fmt.Fprintf(&text, "information: %s\n", o.errMessage)
	}

	if len(o.instruction) > 0 {
		fmt.Fprintf(&text, "instruction: %s\n", o.instruction)
	}

	return text.String()
}

// TablesChecker checks compatibility of table structures, there are differences between MySQL and TiDB.
// In generally we need to check definitions of columns, constraints and table options.
// Because of the early TiDB engineering design, we did not have a complete list of check items, which are all based on experience now.
type TablesChecker struct {
	db     *sql.DB
	dbinfo *dbutil.DBConfig
	tables map[string][]string // schema => []table; if []table is empty, query tables from db
}

// NewTablesChecker returns a RealChecker.
func NewTablesChecker(db *sql.DB, dbinfo *dbutil.DBConfig, tables map[string][]string) RealChecker {
	return &TablesChecker{
		db:     db,
		dbinfo: dbinfo,
		tables: tables,
	}
}

// Check implements RealChecker interface.
func (c *TablesChecker) Check(ctx context.Context) *Result {
	r := &Result{
		Name:  c.Name(),
		Desc:  "check compatibility of table structure",
		State: StateSuccess,
		Extra: fmt.Sprintf("address of db instance - %s:%d", c.dbinfo.Host, c.dbinfo.Port),
	}

	var (
		err     error
		options = make(map[string][]*incompatibilityOption)
	)
	for schema, tables := range c.tables {
		if len(tables) == 0 {
			tables, err = dbutil.GetTables(ctx, c.db, schema)
			if err != nil {
				markCheckError(r, err)
				return r
			}
		}

		for _, table := range tables {
			tableName := dbutil.TableName(schema, table)
			statement, err := dbutil.GetCreateTableSQL(ctx, c.db, schema, table)
			if err != nil {
				// continue if table was deleted when checking
				if isMySQLError(err, mysql.ErrNoSuchTable) {
					continue
				}
				markCheckError(r, err)
				return r
			}

			opts := c.checkCreateSQL(ctx, statement)
			if len(opts) > 0 {
				options[tableName] = opts
			}
		}
	}

	for name, opts := range options {
		if len(opts) == 0 {
			continue
		}
		tableMsg := "table " + name + " "

		for _, option := range opts {
			switch option.state {
			case StateWarning:
				r.State = StateWarning
				e := NewError(tableMsg + option.errMessage)
				e.Severity = StateWarning
				e.Instruction = option.instruction
				r.Errors = append(r.Errors, e)
			case StateFailure:
				r.State = StateFailure
				e := NewError(tableMsg + option.errMessage)
				e.Instruction = option.instruction
				r.Errors = append(r.Errors, e)
			}
		}
	}

	return r
}

// Name implements RealChecker interface.
func (c *TablesChecker) Name() string {
	return "table structure compatibility check"
}

func (c *TablesChecker) checkCreateSQL(ctx context.Context, statement string) []*incompatibilityOption {
	parser2, err := dbutil.GetParserForDB(ctx, c.db)
	if err != nil {
		return []*incompatibilityOption{
			{
				state:      StateFailure,
				errMessage: err.Error(),
			},
		}
	}

	stmt, err := parser2.ParseOneStmt(statement, "", "")
	if err != nil {
		return []*incompatibilityOption{
			{
				state:      StateFailure,
				errMessage: err.Error(),
			},
		}
	}
	// Analyze ast
	return c.checkAST(stmt)
}

func (c *TablesChecker) checkAST(stmt ast.StmtNode) []*incompatibilityOption {
	st, ok := stmt.(*ast.CreateTableStmt)
	if !ok {
		return []*incompatibilityOption{
			{
				state:      StateFailure,
				errMessage: fmt.Sprintf("Expect CreateTableStmt but got %T", stmt),
			},
		}
	}

	var options []*incompatibilityOption
	// check colum def
	for _, def := range st.Cols {
		option := checkColumnDef(def)
		if option != nil {
			options = append(options, option)
		}
	}
	// check constrains
	for _, cst := range st.Constraints {
		option := checkConstraint(cst)
		if option != nil {
			options = append(options, option)
		}
	}
	// check primary/unique key
	hasUnique := false
	for _, cst := range st.Constraints {
		if checkUnique(cst) {
			hasUnique = true
			break
		}
	}
	if !hasUnique {
		options = append(options, &incompatibilityOption{
			state:       StateWarning,
			instruction: "please set primary/unique key for the table",
			errMessage:  "primary/unique key does not exist",
		})
	}

	// check options
	for _, opt := range st.Options {
		option := checkTableOption(opt)
		if option != nil {
			options = append(options, option)
		}
	}
	return options
}

func checkColumnDef(def *ast.ColumnDef) *incompatibilityOption {
	return nil
}

func checkConstraint(cst *ast.Constraint) *incompatibilityOption {
	if cst.Tp == ast.ConstraintForeignKey {
		return &incompatibilityOption{
			state:       StateWarning,
			instruction: "please ref document: https://docs.pingcap.com/tidb/stable/mysql-compatibility#unsupported-features",
			errMessage:  fmt.Sprintf("Foreign Key %s is parsed but ignored by TiDB.", cst.Name),
		}
	}

	return nil
}

func checkUnique(cst *ast.Constraint) bool {
	switch cst.Tp {
	case ast.ConstraintPrimaryKey, ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex:
		return true
	}
	return false
}

func checkTableOption(opt *ast.TableOption) *incompatibilityOption {
	if opt.Tp == ast.TableOptionCharset {
		// Check charset
		cs := strings.ToLower(opt.StrValue)
		if cs != "binary" && !charset.ValidCharsetAndCollation(cs, "") {
			return &incompatibilityOption{
				state:       StateFailure,
				instruction: "https://docs.pingcap.com/tidb/stable/mysql-compatibility#unsupported-features",
				errMessage:  fmt.Sprintf("unsupport charset %s", opt.StrValue),
			}
		}
	}
	return nil
}

// ShardingTablesChecker checks consistency of table structures of one sharding group
// * check whether they have same column list
// * check whether they have auto_increment key.
type ShardingTablesChecker struct {
	name string

	dbs    map[string]*sql.DB
	tables map[string]map[string][]string // instance => {schema: [table1, table2, ...]}
}

// NewShardingTablesChecker returns a RealChecker.
func NewShardingTablesChecker(name string, dbs map[string]*sql.DB, tables map[string]map[string][]string) RealChecker {
	return &ShardingTablesChecker{
		name:   name,
		dbs:    dbs,
		tables: tables,
	}
}

// Check implements RealChecker interface.
func (c *ShardingTablesChecker) Check(ctx context.Context) *Result {
	r := &Result{
		Name:  c.Name(),
		Desc:  "check consistency of sharding table structures",
		State: StateSuccess,
		Extra: fmt.Sprintf("sharding %s", c.name),
	}

	var (
		stmtNode      *ast.CreateTableStmt
		firstTable    string
		firstInstance string
	)

	for instance, schemas := range c.tables {
		db, ok := c.dbs[instance]
		if !ok {
			markCheckError(r, errors.NotFoundf("client for instance %s", instance))
			return r
		}

		parser2, err := dbutil.GetParserForDB(ctx, db)
		if err != nil {
			markCheckError(r, err)
			r.Extra = fmt.Sprintf("fail to get parser for instance %s on sharding %s", instance, c.name)
			return r
		}

		for schema, tables := range schemas {
			for _, table := range tables {
				statement, err := dbutil.GetCreateTableSQL(ctx, db, schema, table)
				if err != nil {
					// continue if table was deleted when checking
					if isMySQLError(err, mysql.ErrNoSuchTable) {
						continue
					}
					markCheckError(r, err)
					r.Extra = fmt.Sprintf("instance %s on sharding %s", instance, c.name)
					return r
				}

				stmt, err := parser2.ParseOneStmt(statement, "", "")
				if err != nil {
					markCheckError(r, errors.Annotatef(err, "statement %s", statement))
					r.Extra = fmt.Sprintf("instance %s on sharding %s", instance, c.name)
					return r
				}

				ctStmt, ok := stmt.(*ast.CreateTableStmt)
				if !ok {
					markCheckError(r, errors.Errorf("Expect CreateTableStmt but got %T", stmt))
					r.Extra = fmt.Sprintf("instance %s on sharding %s", instance, c.name)
					return r
				}

				if stmtNode == nil {
					stmtNode = ctStmt
					firstTable = dbutil.TableName(schema, table)
					firstInstance = instance
					continue
				}

				checkErr := c.checkConsistency(stmtNode, ctStmt, firstTable, dbutil.TableName(schema, table), firstInstance, instance)
				if checkErr != nil {
					r.State = StateFailure
					r.Errors = append(r.Errors, checkErr)
					r.Extra = fmt.Sprintf("error on sharding %s", c.name)
					r.Instruction = "please set same table structure for sharding tables"
					return r
				}
			}
		}
	}

	return r
}

type briefColumnInfo struct {
	name         string
	tp           string
	isUniqueKey  bool
	isPrimaryKey bool
}

func (c *briefColumnInfo) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s %s", c.name, c.tp)
	if c.isPrimaryKey {
		fmt.Fprintln(&buf, " primary key")
	} else if c.isUniqueKey {
		fmt.Fprintln(&buf, " unique key")
	}

	return buf.String()
}

type briefColumnInfos []*briefColumnInfo

func (cs briefColumnInfos) String() string {
	colStrs := make([]string, 0, len(cs))
	for _, col := range cs {
		colStrs = append(colStrs, col.String())
	}

	return strings.Join(colStrs, "\n")
}

func (c *ShardingTablesChecker) checkConsistency(self, other *ast.CreateTableStmt, selfTable, otherTable, selfInstance, otherInstance string) *Error {
	selfColumnList := getBriefColumnList(self)
	otherColumnList := getBriefColumnList(other)

	if len(selfColumnList) != len(otherColumnList) {
		e := NewError("column length mismatch (self: %d vs other: %d)", len(selfColumnList), len(otherColumnList))
		getColumnNames := func(infos briefColumnInfos) []string {
			ret := make([]string, 0, len(infos))
			for _, info := range infos {
				ret = append(ret, info.name)
			}
			return ret
		}
		e.Self = fmt.Sprintf("instance %s table %s columns %v", selfInstance, selfTable, getColumnNames(selfColumnList))
		e.Other = fmt.Sprintf("instance %s table %s columns %v", otherInstance, otherTable, getColumnNames(otherColumnList))
		return e
	}

	for i := range selfColumnList {
		if *selfColumnList[i] != *otherColumnList[i] {
			e := NewError("different column definition")
			e.Self = fmt.Sprintf("instance %s table %s column %s", selfInstance, selfTable, selfColumnList[i])
			e.Other = fmt.Sprintf("instance %s table %s column %s", otherInstance, otherTable, otherColumnList[i])
			return e
		}
	}

	return nil
}

func getBriefColumnList(stmt *ast.CreateTableStmt) briefColumnInfos {
	columnList := make(briefColumnInfos, 0, len(stmt.Cols))

	for _, col := range stmt.Cols {
		bc := &briefColumnInfo{
			name: col.Name.Name.L,
			tp:   col.Tp.String(),
		}

		for _, opt := range col.Options {
			switch opt.Tp {
			case ast.ColumnOptionPrimaryKey:
				bc.isPrimaryKey = true
			case ast.ColumnOptionUniqKey:
				bc.isUniqueKey = true
			}
		}

		columnList = append(columnList, bc)
	}

	return columnList
}

// Name implements Checker interface.
func (c *ShardingTablesChecker) Name() string {
	return fmt.Sprintf("sharding table %s consistency checking", c.name)
}

type AutoIncrementChecker struct {
	db     *sql.DB
	dbinfo *dbutil.DBConfig
	tables map[string][]string // schema => []table; if []table is empty, query tables from db
}

// NewTablesChecker returns a RealChecker.
func NewAutoIncrementChecker(db *sql.DB, dbinfo *dbutil.DBConfig, tables map[string][]string) RealChecker {
	return &AutoIncrementChecker{
		db:     db,
		dbinfo: dbinfo,
		tables: tables,
	}
}

func (a *AutoIncrementChecker) Check(ctx context.Context) *Result {
	r := &Result{
		Name:  a.Name(),
		Desc:  "check compatibility of table structure",
		State: StateSuccess,
		Extra: fmt.Sprintf("address of db instance - %s:%d", a.dbinfo.Host, a.dbinfo.Port),
	}
	var err error
	for schema, tables := range a.tables {
		if len(tables) == 0 {
			tables, err = dbutil.GetTables(ctx, a.db, schema)
			if err != nil {
				markCheckError(r, err)
				return r
			}
		}

		for _, table := range tables {
			statement, err := dbutil.GetCreateTableSQL(ctx, a.db, schema, table)
			if err != nil {
				// continue if table was deleted when checking
				if isMySQLError(err, mysql.ErrNoSuchTable) {
					continue
				}
				markCheckError(r, err)
				return r
			}

			parser2, err := dbutil.GetParserForDB(ctx, a.db)
			if err != nil {
				markCheckError(r, err)
				return r
			}
			stmt, err := parser2.ParseOneStmt(statement, "", "")
			if err != nil {
				markCheckError(r, err)
				return r
			}
			if st, ok := stmt.(*ast.CreateTableStmt); !ok {
				r.State = StateFailure
				r.Instruction = fmt.Sprintf("Expect CreateTableStmt but got %T", stmt)
				return r
			} else {
				for _, def := range st.Cols {
					if hasAutoIncrement(def) {
						r.State = StateWarning
						r.Errors = append(r.Errors, &Error{
							Severity:    StateWarning,
							Instruction: "TODO: method resolved PK/UK conflicts",
						})
					}
				}
			}
		}
	}
	return r
}

func (a *AutoIncrementChecker) Name() string {
	return "check auto increment key."
}

func hasAutoIncrement(col *ast.ColumnDef) bool {
	for _, opt := range col.Options {
		if opt.Tp == ast.ColumnOptionAutoIncrement {
			return true
		}
	}
	return false
}
