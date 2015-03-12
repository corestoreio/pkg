# CoreStore FrameWork

This repository contains the main framework.

Please see the README files in the sub-packages.

Magento is a trademark of [MAGENTO, INC.](http://www.magentocommerce.com/license/).

## Usage

To properly use the CoreStore framework some environment variables must be set before running `go generate`.

### Required settings

- `CS_DSN` the environment variable for the MySQL connection.
- ...

### Optional settings

- `CS_EAV_MAP` the environment variable which points to a JSON file for mapping EAV entities.
- ...

```
$ go get github.com/corestoreio/csfw
$ cd $GOPATH/src/github.com/corestoreio/csfw
$ go generate
```

## Contributing

Please have a look at the [contribution guidelines](https://github.com/corestoreio/corestore/blob/master/CONTRIBUTING.md).

## Licensing

CoreStore is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/corestoreio/corestore/blob/master/LICENSE) for the full license text.

## Copyright

[Cyrill Schumacher](http://cyrillschumacher.com) - [PGP Key](https://keybase.io/cyrill)
