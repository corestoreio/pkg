// Copyright 2015 CoreStore Authors
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

package i18n

//type (
//	Locale struct {
//		language.Tag
//	}
//
//	I18n struct {
//		DefaultLocale string
//		defaultLocale Locale
//	}
//
//	Option func(*I18n)
//)
//
//var ErrLocaleIncorrect = errors.New("Incorrect locale")
//var defaultI18n = NewI18n()
//
//func NewI18n() *I18n {
//	return &I18n{}
//}
//
//func Init(opts ...Option) error {
//	return defaultI18n.Init(opts...)
//}
//
//func (i *I18n) Init(opts ...Option) error {
//	for _, opt := range opts {
//		opt(i)
//	}
//
//	if i.DefaultLocale == "" {
//		return ErrLocaleIncorrect
//	}
//
//	i.defaultLocale = Locale{
//		Tag: language.Make(i.DefaultLocale),
//	}
//	return nil
//}
//
//func (i *I18n) __(translationID string, args ...interface{}) string {
//	return translationID
//}
//
//func __(translationID string, args ...interface{}) string {
//	return defaultI18n.__(translationID, args...)
//}
