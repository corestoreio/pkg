// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jsonapi

// Errors implements the JSON API Errors document.
type Errors struct {
	Errors []Error `json:"errors"`
}

// NewErrors creates a new error object for serialization as JSON API Errors
// document.
func NewErrors(e Error, es ...Error) Errors {
	return Errors{
		Errors: append([]Error{e}, es...),
	}
}

// Error implements the JSON API Error object.
type Error struct {
	ID     string                 `json:"id,omitempty"`     // a unique identifier for this particular occurrence of the problem
	Links  map[string]string      `json:"links,omitempty"`  // a links object containing a field about that leads to further details about this particular occurrence of the problem
	Status string                 `json:"status,omitempty"` // the HTTP status code applicable to this problem, expressed as a string value
	Code   string                 `json:"code,omitempty"`   // an application-specific error code, expressed as a string value
	Title  string                 `json:"title,omitempty"`  // a short, human-readable summary of the problem that SHOULD NOT change from occurrence to occurrence of the problem, except for purposes of localization
	Detail string                 `json:"detail,omitempty"` // a human-readable explanation specific to this occurrence of the problem
	Source ErrorSource            `json:"source,omitempty"` // a Souce object containing references to the source of the error
	Meta   map[string]interface{} `json:"meta,omitempty"`   // a meta object containing non-standard meta-information about the error
}

// NewError returns a basic Error that is not related to a field or a URL parameter.
func NewError(about, status, code, title, detail string) Error {
	return Error{
		Links:  map[string]string{"about": about},
		Status: status,
		Code:   code,
		Title:  title,
		Detail: detail,
	}
}

// NewError returns a basic Error that is related to a field or a URL parameter.
func NewFieldError(status, code, title, detail, pointer string) Error {
	return Error{
		Status: status,
		Code:   code,
		Title:  title,
		Detail: detail,
		Source: ErrorSource{
			Pointer: pointer,
		},
	}
}

// ErrorSource contains references to the source of the error
type ErrorSource struct {
	Pointer   string `json:"pointer,omitempty"`   // a JSON Pointer [RFC6901] to the associated entity in the request document
	Parameter string `json:"parameter,omitempty"` // a string indicating which URI query parameter caused the error
}
