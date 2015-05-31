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

/*
Package config handles the scopes and the configuration via consul, etc or simple files.

Elements

The three elements Section, Group and Field represents front end configuration fields and more important
default values and their backend/source models (loading and saving).

Those three elements represents the PackageConfiguration variable which can be found in any package.

Your app which includes the csfw must merge all "PackageConfiguration"s into a single slice.
You should submit all default values (interface config.Sectioner) to the config.Manager.ApplyDefaults()
function. @todo merge core_config_data values in config.Manager
The models included in PackageConfiguration will be later used when handling the values
for each configuration field.

The JSON enconding of the three elements Section, Group and Field are intended to use
on the backend REST API and for debugging and testing. Only used in non performance critical parts.

Scope Values

To get a value from the configuration manager via any Get* method you have to set up the arguments.
At least a config.Path() is needed. If you need a config value from another scope (store or website)
you must also supply a Scope*() value. Without the scope the default value will be returned.

	val := config.Manager.GetString(config.Path("path/to/setting"))

Above code returns the default value for path/to/setting key.

Can also be rewritten without using slashes:

	val := config.Manager.GetString(config.Path("path", "to", "setting"))

Returning a website scope based value:

	w := store.Manager.Website()
	val := config.Manager.GetString(config.Path("path/to/setting"), config.Scope(config.IDScopeWebsite, w))

can be rewritten as:

	w := store.Manager.Website()
	val := config.Manager.GetString(config.Path("path/to/setting"), config.ScopeWebsite(w))

The code returns the value for a specific website scope. If the value has not been found then the
default value will be returned.

Returning a store scope based value:

	w := store.Manager.Website()
	val := config.Manager.GetString(config.Path("path/to/setting"), config.Scope(config.IDScopeStore, w))

can be rewritten as:

	w := store.Manager.Website()
	val := config.Manager.GetString(config.Path("path/to/setting"), config.ScopeStore(w))

The code returns the value for a specific store scope. If the value has not been found then the
default value will be returned.

Mixing Store and Website scope in calling of any Write/Get*() function will return that value which scope
will be added at last to the OptionFunc slice.

*/
package config
