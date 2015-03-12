# eavToStruct Code generator

```shell
$ eavToStruct
Usage of eavToStruct:
  -o="": Output file name
  -p="": Package name in template
  -prefixName="": Table name prefix
  -prefixSearch="eav": Search Table Prefix. Used in where condition to list tables
  -run=false: If true program runs
```

This program depends on the generated table structs from `tableToStruct` program.

## Install

Please don't. This runs via go:generate.
