// Auto generated via github.com/corestoreio/pkg/sql/dmlgen

namespace {{.Package}};

include "github.com/corestoreio/pkg/storage/null/null.fbs";

{{range $opts := .SerializerHeaderOptions -}}
{{$opts}};
{{end}}
