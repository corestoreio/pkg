# CoreStore FrameWork

This repository contains the main framework.

Please see the README files in the sub-packages.

Magento is a trademark of [MAGENTO, INC.](http://www.magentocommerce.com/license/).

## Badges

[goreportcard](http://goreportcard.com/report/Corestoreio/csfw)

@todo add travis

## Usage

To properly use the CoreStore framework some environment variables must be set before running `go generate`.

### Required settings

`CS_DSN` the environment variable for the MySQL connection.

```shell
$ export CS_DSN='magento1:magento1@tcp(localhost:3306)/magento1'
$ export CS_DSN='magento2:magento2@tcp(localhost:3306)/magento2'
```

### Optional settings

- `CS_EAV_MAP` the environment variable which points to a JSON file for mapping EAV entities.
- ...

```
$ go get github.com/corestoreio/csfw
$ export CS_DSN_TEST='see next section'
$ cd $GOPATH/src/github.com/corestoreio/csfw
$ go generate ./...
```

## Testing

Setup two databases. One for Magento 1 and one for Magento and fill them.

Create a DSN env var `CS_DSN_TEST` and point it to Magento 1 database. Run the tests.
Change the env var to let it point to Magento 2 database. Rerun the tests.

```shell
$ export CS_DSN_TEST='magento1:magento1@tcp(localhost:3306)/magento1'
$ export CS_DSN_TEST='magento2:magento2@tcp(localhost:3306)/magento2'
```

## Contributing

Please have a look at the [contribution guidelines](https://github.com/corestoreio/corestore/blob/master/CONTRIBUTING.md).

## Licensing

CoreStore is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/corestoreio/corestore/blob/master/LICENSE) for the full license text.

## Copyright

[Cyrill Schumacher](http://cyrillschumacher.com) - [PGP Key](https://keybase.io/cyrill)
