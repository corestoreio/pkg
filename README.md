# CoreStore library WIP

WIP = Work in Progress

[![Build Status](https://travis-ci.org/corestoreio/csfw.svg?branch=master)](https://travis-ci.org/corestoreio/csfw) [![wercker status](https://app.wercker.com/status/d7d0bdda415d2228b6fb5bb01681b5c4/s/master "wercker status")](https://app.wercker.com/project/bykey/d7d0bdda415d2228b6fb5bb01681b5c4) [![Appveyor Build status](https://ci.appveyor.com/api/projects/status/lrlnbpcjdy585mg1/branch/master?svg=true)](https://ci.appveyor.com/project/SchumacherFM/csfw/branch/master) [![GoDoc](http://godoc.org/github.com/corestoreio/csfw?status.svg)](http://godoc.org/github.com/corestoreio/csfw) [![Join the chat at https://gitter.im/corestoreio/csfw](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/corestoreio/csfw?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) [goreportcard](http://goreportcard.com/report/Corestoreio/csfw)

eCommerce library which is compatible to Magento 1 and 2 database schema.

Magento is a trademark of [MAGENTO, INC.](http://www.magentocommerce.com/license/).

## Usage

To properly use the CoreStore library some environment variables must be set before running `go generate`. (TODO)

### Required settings

`CS_DSN` the environment variable for the MySQL connection.

```shell
$ export CS_DSN='magento1:magento1@tcp(localhost:3306)/magento1'
$ export CS_DSN='magento2:magento2@tcp(localhost:3306)/magento2'
```

```
$ go get github.com/corestoreio/csfw
$ export CS_DSN='see previous'
$ cd $GOPATH/src/github.com/corestoreio/csfw
$ go run codegen/tableToStruct/*.go
```

## Testing

Setup two databases. One for Magento 1 and one for Magento 2 and fill them with the provided [test data](https://github.com/corestoreio/csfw/tree/master/testData).

Create a DSN env var `CS_DSN` and point it to Magento 1 database. Run the tests.
Change the env var to let it point to Magento 2 database. Rerun the tests.

### Finding Allocations

Side note: There is a new testing approach: TBDD = Test n' Benchmark Driven Development.

On the first run we got this result:

```
$ go test -run=😶 -bench=Benchmark_WithInitStoreByToken .
PASS
Benchmark_WithInitStoreByToken-4	  100000	     17297 ns/op	    9112 B/op	     203 allocs/op
ok  	github.com/corestoreio/csfw/store	2.569s

```

Quite shocking to use 203 allocs for just figuring out the current store view within a request.

Now compile your tests into an executable binary:

```
$ go test -c
```

This compilation reduces the noise in the below output trace log.

```
$ GODEBUG=allocfreetrace=1 ./store -test.run=😶 -test.bench=Benchmark_WithInitStoreByToken -test.benchtime=10ms 2>trace.log
```

Now open the trace.log file (around 26MB) and investigate all the allocations and refactor your code. Once finished you can achieve results like:

```
$ go test -run=😶  -bench=Benchmark_WithInitStoreByToken .
PASS
Benchmark_WithInitStoreByToken-4	 2000000	       826 ns/op	     128 B/op	       5 allocs/op
ok  	github.com/corestoreio/csfw/store	2.569s
```

#### Other development helpers

- [go get github.com/maruel/panicparse/cmd/pp](https://github.com/maruel/panicparse)
- [go get github.com/alecthomas/gometalinter](https://github.com/alecthomas/gometalinter)

## TODO

- Create Magento 1+2 modules to setup test database and test Magento system.

## Contributing

Please have a look at the [contribution guidelines](https://github.com/corestoreio/corestore/blob/master/CONTRIBUTING.md).

## Acknowledgements

| Name | Package | License |
| -------|----------|-------|
| Steve Francia | [cast](http://github.com/corestoreio/csfw/tree/master/utils/cast) | MIT Copyright (c) 2014 |
| Steve Francia | [bufferpool](http://github.com/corestoreio/csfw/tree/master/utils/bufferpool) | [Simple Public License, V2.0](http://opensource.org/licenses/Simple-2.0) |
| Jonathan Novak, Tyler Smith, Michal Bohuslávek | [dbr](http://github.com/corestoreio/csfw/tree/master/storage/dbr) | The MIT License (MIT) 2014 |
| Martin Angers and Contributors. | [ctxthrottled](http://github.com/corestoreio/csfw/tree/master/net/ctxthrottled) | The MIT License (MIT) 2014 |
| Canonical Ltd. | [errgo](https://github.com/juju/errgo) | [LGPLv3](http://www.gnu.org/licenses/lgpl-3.0.en.html) |
| Jad Dittmar | [finance](https://github.com/Confunctionist/finance) aka. [money](http://github.com/corestoreio/csfw/tree/master/storage/money) | Copyright (c) 2011 |
| Wenbin Xiao | [sqlparser](https://github.com/xwb1989/sqlparser) | Copyright 2015 BSD Style |
| Google Inc | [youtube/vitess\sqlparser](https://github.com/youtube/vitess) | Copyright 2012 BSD Style |
| Julien Schmidt | [ctxrouter](https://github.com/corestoreio/csfw/tree/master/net/ctxrouter) aka. [httprouter](https://github.com/julienschmidt/httprouter) | Copyright (c) 2013 MIT |
| Olivier Poitrey| [ctxmw.WithAccessLog](https://github.com/corestoreio/csfw/tree/master/net/ctxmw) & CORS | Copyright (c) 2014-2015  MIT License |
| Dave Grijalva| [csjwt](https://github.com/corestoreio/csfw/tree/master/util/csjwt) | Copyright (c) 2012 MIT License |

## Licensing

CoreStore is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/corestoreio/corestore/blob/master/LICENSE) for the full license text.

## Copyright

[Cyrill Schumacher](http://cyrillschumacher.com) - [PGP Key](https://keybase.io/cyrill)
