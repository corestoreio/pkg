# eav library

Contains the logic for the Entity-Attribute-Value model based on the Magento database schema.

To use this library with additional columns in the EAV tables you must run from the
tools folder first `tableToStruct` and then build the program `eavToStruct` and run it.

@todo use build tags to separate between gen_tables.go file and custom created file.

## API Versioning

The folders `v0`, `v1`, `vX` contains later the official API endpoint. The functions in these
endpoints are wrappers for the underlaying code.

