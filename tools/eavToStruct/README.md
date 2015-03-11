# eavToStruct Code generator

```shell
$ eavToStruct
Usage of eavToStruct:
  -dsn="test:test@tcp(localhost:3306)/test": MySQL DSN data source name. Can also be provided via ENV with key CS_DSN
  -file="table_struct.go": Output file name
  -package="eav": Package name in template
  -prefixName="": Table name prefix
  -prefixSearch="eav": Search Table Prefix. Used in where condition to list tables
  -run=false: If true program runs
```

This program depends on the generated table structs from `tableToStruct` program.

## Install

```shell
$ go install github.com:corestoreio/tools/eavToStruct
```
