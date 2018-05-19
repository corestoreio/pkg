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

package config

import (
	"os"
	"unicode"

	"github.com/corestoreio/log"
)

// EnvNamePlaceHolder replaces in a file name or pattern argument applied to a
// WithFiles, WithFile or WithGlob function with the current environment name of
// *config.Service. You can load environment dependent configuration files.
const EnvNamePlaceHolder = `{CS_ENV}`

const DefaultOSEnvVariableName = `CS_ENV`

// Options applies configurations to the NewService function. Used mainly by external
// packages for providing different storage engines.
type Options struct {
	// Level1 defines a usually short lived cached like an LRU or TinyLFU or a
	// storage with un-resetble data (locked data). It can be set optionally.
	// Level1 only gets called for Set operation when a value is requested in
	// Level2 via Get.
	Level1       Storager
	Log          log.Logger
	EnablePubSub bool
	// OSEnvVariableName loads a string from an applied environment variable to
	// use it as a prefix for the Path type and when loading configuration files
	// as part of their filename or path (see cfgfile.EnvNamePlaceHolder). For
	// example different payment gateways and access credentials for a payment
	// method on STAGING or PRODUCTION systems.
	OSEnvVariableName string
	// EnvName defines the current name of the environment directly instead of
	// loading it from an operating system environment variable. It gets used as
	// a prefix for the Path type or as part of the directory or filename when
	// loading configuration files (see cfgfile.EnvNamePlaceHolder).
	EnvName string

	// EnableHotReload if the Service receives an OS signal, it triggers a hot
	// reload of the cached functions of type LoadDataFn. Errors during hot
	// reloading do not trigger an exit of the config.Service.
	EnableHotReload bool
	// HotReloadSignals specifies custom signals to listen to. Defaults to
	// syscall.SIGUSR2
	HotReloadSignals []os.Signal
}

// LoadDataFn allows other storage backends to pump their data into the
// config.Service during or after initialization via an OS signal and hot
// reloading.
// TODO LoadDataFn should support the target storage to load the data to allows to write locked values.
type LoadDataFn func(*Service) error

//type LoadData  struct {
//	Process func(*Service) error
//	TargetLevel int
//}

func isLetter(str string) bool {
	for _, r := range str {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
