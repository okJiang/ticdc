drop database if exists `drop_column_with_index`;
create database `drop_column_with_index`;
use `drop_column_with_index`;
create table t1 (c1 int, c2 int, c3 int);
alter table t1 add index idx1(c1);
alter table t1 add index idx2(c2);
alter table t1 add index idx3(c3);
alter table t1 add index idx23(c2, c3);
alter table t1 add unique index uidx1(c1);
alter table t1 add unique index uidx2(c2);
