// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

// Package xacml implements the OASIS/XACML standard for Policy-based
// authorization.
//
// TODO implement xacml
// - https://github.com/Hasdcorona/go-xacml
// - https://github.com/murphysean/xacml
// - https://github.com/enygma/xacmlphp
// - https://www.oasis-open.org/committees/download.php/2713/Brief_Introduction_to_XACML.html//
//
// The OASIS Standard
//
// The OASIS/XACML standard is a well-defined XML-based structure for evaluating
// attributes on Policies against attributes on Subjects to see if there's a
// match (based on Operation rules and combining Algorithms).
//
// Terminology
//
// PolicySet: Set of Policy objects
//
// Policy: Defines the policies to evaluate for authoriation. Policies contain
// sets of Rules that are evaluated and the results are combined according to
// the Policy's Algorithm for an overall Policy pass/fail status
//
// Rule: A Rule is made of of a set of Matches (inside a Target) that are used
// to evaluate authorization
//
// Match: An object that defines the property to look at (Designator) and the
// value to check against (Value) and the Operation to perform (like
// "StringEqual") for Permit/Deny result
//
// Attribute: Property on a Subject, Resource, Action or Environment
//
// Algorithm: Evaluation method for combining results of the object (like Policy
// or Rule). In the OASIS spec, these are called Functions.
//
// Effect: According to the spec, this can only be "PERMIT" or "DENY"
//
// Enforcer: Point of enforcement of the access, called the PEP (Policy
// Enforcement Point) in the OASIS spec.
//
// Decider: The object that handles the decision logic, tracing down from
// Policies to Matches. Called the PDP (Policy Decision Point) in the OASIS
// spec.
//
// Resource: An object representing a "something" the Subject is trying to
// access.
//
package xacml
