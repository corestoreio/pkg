# CoreStore FrameWork

[![Join the chat at https://gitter.im/corestoreio/csfw](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/corestoreio/csfw?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

This repository contains the main framework.

Please see [godoc.org](https://godoc.org/github.com/corestoreio/csfw) which is more up-to-date than this README.md file.

Magento is a trademark of [MAGENTO, INC.](http://www.magentocommerce.com/license/).

## Badges

[goreportcard](http://goreportcard.com/report/Corestoreio/csfw) [![GoDoc](https://godoc.org/github.com/corestoreio/csfw?status.svg)](https://godoc.org/github.com/corestoreio/csfw)

@todo add travis

## Usage

To properly use the CoreStore framework some environment variables must be set before running `go generate`.

### Required settings

`CS_DSN` the environment variable for the MySQL connection.

```shell
$ export CS_DSN='magento1:magento1@tcp(localhost:3306)/magento1'
$ export CS_DSN='magento2:magento2@tcp(localhost:3306)/magento2'
```

```
$ go get github.com/corestoreio/csfw
$ export CS_DSN_TEST='see next section'
$ cd $GOPATH/src/github.com/corestoreio/csfw
$ go generate ./...
```

## Testing

Setup two databases. One for Magento 1 and one for Magento 2 and fill them with the provided [test data](https://github.com/corestoreio/csfw/tree/master/testData).

Create a DSN env var `CS_DSN_TEST` and point it to Magento 1 database. Run the tests.
Change the env var to let it point to Magento 2 database. Rerun the tests.

```shell
$ export CS_DSN_TEST='magento1:magento1@tcp(localhost:3306)/magento1'
$ export CS_DSN_TEST='magento2:magento2@tcp(localhost:3306)/magento2'
```

## IDE

Currently using the IntelliJ IDEA Community Edition with the [go-lang-idea-plugin](https://github.com/go-lang-plugin-org/go-lang-idea-plugin).
At the moment Q2/2015: There are no official jar files for downloading so the go lang plugin will be 
compiled on a daily basis. The plugin works very well!

IDEA has been configured with goimports, gofmt, golint, govet and ... with the file watcher plugin.

Why am I not using vim? Because I would only generate passwords ;-|.

## Contributing

Please have a look at the [contribution guidelines](https://github.com/corestoreio/corestore/blob/master/CONTRIBUTING.md).

## Licensing

CoreStore is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/corestoreio/corestore/blob/master/LICENSE) for the full license text.

## Copyright

[Cyrill Schumacher](http://cyrillschumacher.com) - [PGP Key](https://keybase.io/cyrill)
