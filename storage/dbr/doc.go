// Package dbr has additions to Go's database/sql for super fast performance and
// type safety and convenience.
//
// Aim: Allow a developer to easily modify a SQL query without type assertion of
// parts of the query. This package gets extended during csfw development.
//
// Abbreviations
//
// DML (https://en.wikipedia.org/wiki/Data_manipulation_language) Select,
// Insert, Update and Delete.
//
// DDL (https://en.wikipedia.org/wiki/Data_definition_language) Create, Drop,
// Alter, and Rename.
//
// DCL (https://en.wikipedia.org/wiki/Data_control_language) Grant and Revoke.
//
// CRUD (https://en.wikipedia.org/wiki/Create,_read,_update_and_delete) Create,
// Read, Update and Delete.
//
// TODO(CyS) Add named parameter from GO1.8 to each query builder
package dbr
