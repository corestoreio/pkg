// Auto generated via github.com/corestoreio/pkg/sql/dmlgen
syntax = "proto3";
package {{.Package}};
import "github.com/gogo/protobuf/gogoproto/gogo.proto";
import "google/protobuf/timestamp.proto";
import "github.com/corestoreio/pkg/sql/dml/types_null.proto";
option go_package = "{{.Package}}";
{{range $opts := .GogoProtoOptions -}}
option {{$opts}};
{{end}}
