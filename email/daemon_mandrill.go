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

package email

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/mail"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/utils/log"
	"github.com/mattbaird/gochimp"
)

// @todo Add AmazonSES, Mailgun, MailJet, SendGrid, PostMark, etc
// https://github.com/aws/aws-sdk-go/blob/master/service%2Fses%2Fexamples_test.go

// MandrillOptions can be used as an argument to SetMandrill() function.
type MandrillOptions func(*gochimp.MandrillAPI)

// SetMandrillRoundTripper sets a round tripper to the MandrillAPI mainly for testing.
func SetMandrillRoundTripper(transport http.RoundTripper) MandrillOptions {
	return func(m *gochimp.MandrillAPI) {
		m.Transport = transport
	}
}

// SetMandrill sets the Mandrill API for sending emails. This function is not
// recursive and returns nil. @todo
func SetMandrill(opts ...MandrillOptions) DaemonOption {
	return func(da *Daemon) DaemonOption {
		// this whole func is just a quick write down. no idea if it's working
		// and refactor ... 8-)
		apiKey := da.Config.GetString(config.ScopeStore(da.Scope), config.Path(PathSmtpMandrillAPIKey))

		if apiKey == "" {
			da.lastErrs = append(da.lastErrs, errors.New("Mandrill API Key is empty."))
			return nil
		}

		md, err := gochimp.NewMandrill(apiKey)
		if err != nil {
			da.lastErrs = append(da.lastErrs, err)
			return nil
		}
		for _, o := range opts {
			o(md)
		}

		da.sendFunc = func(from string, to []string, msg io.WriterTo) error {

			// @todo figure out if "to" contains To, CC and BCC addresses.

			addr, err := mail.ParseAddress(from)
			if err != nil {
				return log.Error("mail.daemon.Mandrill.ParseAddress", "err", err, "from", from, "to", to)
			}

			r := gochimp.Recipient{
				Name:  addr.Name,
				Email: addr.Address,
			}

			var buf bytes.Buffer
			if _, err := msg.WriteTo(&buf); err != nil {
				return log.Error("mail.daemon.Mandrill.MessageWriteTo", "err", err, "from", from, "to", to, "msg", buf.String())
			}

			resp, err := md.MessageSendRaw(buf.String(), to, r, false)
			if err != nil {
				return log.Error("mail.daemon.Mandrill.MessageSendRaw", "err", err, "from", from, "to", to, "msg", buf.String())
			}
			if log.IsDebug() {
				log.Debug("mail.daemon.Mandrill.MessageSendRaw", "resp", resp, "from", from, "to", to, "msg", buf.String())
			}
			// The last arg in MessageSendRaw means async in the Mandrill API:
			// Async: enable a background sending mode that is optimized for bulk sending.
			// In async mode, messages/send will immediately return a status of "queued"
			// for every recipient. To handle rejections when sending in async mode, set
			// up a webhook for the 'reject' event. Defaults to false for messages with
			// no more than 10 recipients; messages with more than 10 recipients are
			// always sent asynchronously, regardless of the value of async.
			return nil
		}
		da.dialer = nil

		return nil
	}
}
