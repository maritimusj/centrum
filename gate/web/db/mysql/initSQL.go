package mysql

const initDBSQL = `
--
-- 由SQLiteStudio v3.2.1 产生的文件 周一 10月 28 14:01:33 2019
--
-- 文本编码：UTF-8
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- 表：alarms
CREATE TABLE alarms (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, status INTEGER NOT NULL DEFAULT 0, org_id INTEGER NOT NULL DEFAULT 0, device_id INTEGER NOT NULL DEFAULT 0, measure_id INTEGER NOT NULL DEFAULT 0, extra BLOB NOT NULL, created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL);

-- 表：api_resources
CREATE TABLE "api_resources" (

"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,

"name"  TEXT(64) NOT NULL,

"title"  TEXT(128) NOT NULL,

"desc"  TEXT(255) NOT NULL

);

-- 表：config
CREATE TABLE config (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, name VARCHAR (128) NOT NULL, extra BLOB NOT NULL, created_at DATETIME NOT NULL, update_at DATETIME NOT NULL);

-- 表：device_groups
CREATE TABLE device_groups (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, device_id INTEGER NOT NULL, group_id INTEGER NOT NULL, created_at DATETIME NOT NULL);

-- 表：devices
CREATE TABLE devices (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, org_id INTEGER NOT NULL DEFAULT 0, enable INTEGER NOT NULL DEFAULT 0, title TEXT (128) NOT NULL, options BLOB NOT NULL, created_at DATETIME NOT NULL);

-- 表：equipment_groups
CREATE TABLE equipment_groups (equipment_id INTEGER NOT NULL DEFAULT 0, group_id INTEGER NOT NULL DEFAULT 0, created_at DATETIME NOT NULL, PRIMARY KEY (equipment_id));

-- 表：equipments
CREATE TABLE equipments (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, org_id INTEGER NOT NULL, enable INTEGER NOT NULL DEFAULT 0, title TEXT (128) NOT NULL, "desc" TEXT (255) NOT NULL, created_at DATETIME NOT NULL);

-- 表：groups
CREATE TABLE groups (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, org_id INTEGER NOT NULL DEFAULT 0, parent_id INTEGER NOT NULL DEFAULT 0, title TEXT (128) NOT NULL, "desc" TEXT (255) NOT NULL, created_at DATETIME NOT NULL);

-- 表：logs
CREATE TABLE logs (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, org_id INTEGER NOT NULL DEFAULT 0, src TEXT (64) NOT NULL, level INTEGER NOT NULL DEFAULT 0, message TEXT (255) NOT NULL, extra BLOB NOT NULL, created_at DATETIME NOT NULL);

-- 表：measures
CREATE TABLE measures (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, enable INTEGER NOT NULL DEFAULT 0, device_id INTEGER NOT NULL DEFAULT 0, title TEXT (128) NOT NULL, tag TEXT (64) NOT NULL, kind INTEGER NOT NULL DEFAULT 0, created_at DATETIME NOT NULL);

-- 表：organizations
CREATE TABLE organizations (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, enable INTEGER NOT NULL DEFAULT 0, name TEXT (64) NOT NULL, title TEXT (128) NOT NULL, extra BLOB NOT NULL, created_at DATETIME NOT NULL);

-- 表：policies
CREATE TABLE policies (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, role_id INTEGER NOT NULL DEFAULT 0, resource_class INTEGER NOT NULL DEFAULT 0, resource_id INTEGER NOT NULL, "action" INTEGER NOT NULL DEFAULT 0, effect INTEGER NOT NULL DEFAULT 0, created_at DATETIME NOT NULL);

-- 表：roles
CREATE TABLE roles (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, org_id INTEGER NOT NULL DEFAULT 0, enable INTEGER NOT NULL DEFAULT 0, name TEXT (64) NOT NULL, title TEXT (128) NOT NULL, "desc" TEXT (512) NOT NULL, created_at DATETIME NOT NULL);

-- 表：states
CREATE TABLE states (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, enable INTEGER NOT NULL DEFAULT 0, title TEXT (64) NOT NULL, "desc" TEXT (512) NOT NULL, equipment_id INTEGER NOT NULL DEFAULT 0, measure_id INTEGER NOT NULL DEFAULT 0, script BLOB NOT NULL, created_at DATETIME NOT NULL);

-- 表：user_roles
CREATE TABLE user_roles (user_id INTEGER NOT NULL, role_id INTEGER NOT NULL, created_at DATETIME NOT NULL, PRIMARY KEY (user_id, role_id));

-- 表：users
CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, org_id INTEGER NOT NULL DEFAULT 0, enable INTEGER NOT NULL DEFAULT 0, name TEXT (64) NOT NULL, title TEXT (128) NOT NULL, password TEXT (64) NOT NULL, mobile TEXT (18) NOT NULL, email TEXT (64) NOT NULL, created_at DATETIME NOT NULL);

-- 索引：device
CREATE INDEX device ON measures ("device_id" ASC);

-- 索引：equipment
CREATE INDEX equipment ON equipment_groups ("equipment_id" ASC, "group_id" ASC);

-- 索引：group
CREATE INDEX "group" ON device_groups ("group_id" ASC, "device_id" ASC);

-- 索引：idx
CREATE UNIQUE INDEX idx ON policies ("role_id" ASC, "resource_class" ASC, "resource_id" ASC, "action" ASC);

-- 索引：name
CREATE UNIQUE INDEX name ON config (name ASC);

-- 索引：name_roles
CREATE UNIQUE INDEX name_roles ON roles ("name" ASC);

-- 索引：namex
CREATE UNIQUE INDEX "namex"

ON "api_resources" ("name" ASC);

-- 索引：orgx
CREATE INDEX orgx ON logs ("org_id" ASC, "src" ASC, "level" ASC);

-- 索引：user_namex
CREATE UNIQUE INDEX user_namex ON users ("name" ASC);

COMMIT TRANSACTION;
PRAGMA foreign_keys = on;

`