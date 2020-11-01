package cstrace

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
)

var ErrorKey = label.Key("error")

// Status codes for use with Span.SetStatus. These correspond to the status
// codes used by gRPC defined here: https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto

// Status if there is an error, it sets the error code "unknown" with the error
// string as span status otherwise status ok.
func Status(span otel.Span, err error, msg string) {
	if err == nil {
		span.SetStatus(codes.Ok, msg)
		return
	}
	span.SetStatus(codes.Error, msg)
	span.SetAttributes(ErrorKey.String(err.Error()))
}

// StatusErrorWithCode sets a custom code with an error.
// go.opencensus.io/trace/status_codes.go. These correspond to the status codes
// used by gRPC defined here:
// https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
