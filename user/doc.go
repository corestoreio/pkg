// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

// Package user takes care of backend/API users.
// Enables admin users to manage and assign roles to administrators and other non-customer users,
// reset user passwords, and invalidate access tokens.
// Different roles can be assigned to different users to define their permissions.

//
// TODO: use punycode to transform email addresses when storing in DB.
// One Quick Note: Though not strictly required, using punycode conversion from
// John@Gıthub.com to xn--john@gthub-2ub.com would have helped prevent this
// issue. It's doubtful any web apps do this as part of the user registration
// process. https://eng.getwisdom.io/hacking-github-with-unicode-dotless-i/ John
// discovered a flaw in the way email addresses were being normalized to
// standard character sets when used to look up accounts during the password
// recovery flow. Password reset tokens are associated with email addresses and
// initiating a password reset with an email address that normalizes to another
// email address would result in the reset token for one user being delivered to
// the email address of another account. The attack only works if an email
// provider allows Unicode in the “local part” of the email address and an
// attacker can claim an email address containing Unicode that would improperly
// normalize to the email address of another account (e.g. mike@example.org vs
// mıke@example.org). Unicode in the “domain part” is not allowed by GitHub's
// outgoing mail server and therefore cannot be used as part of a broader attack
// on common domains (e.g. gmail.com vs gmaıl.com). GitHub addressed the
// vulnerability by making sure the email address in the database matches the
// email address that initiated the reset flow. This ensures that the email
// address used to generate the token matches the email address to which the
// reset token gets delivered. GitHub Security Team

package user
