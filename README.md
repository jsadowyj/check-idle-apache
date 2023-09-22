[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/mnx-solutions/check-idle-apache)
![Go Test](https://github.com/mnx-solutions/check-idle-apache/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/mnx-solutions/check-idle-apache/workflows/goreleaser/badge.svg)

# check-idle-apache

## Table of Contents
- [Overview](#overview)
- [Files](#files)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

## Overview

The check-idle-apache is a [Sensu Check][6] that monitors idle apache worker count.

## Files

## Usage examples

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add mnx-solutions/check-idle-apache
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/mnx-solutions/check-idle-apache].

### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: check-idle-apache
  namespace: default
spec:
  command: check-idle-apache --url http://127.0.0.1/server-status?auto --critical 0 --warning 5
  subscriptions:
  - system
  runtime_assets:
  - mnx-solutions/check-idle-apache
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the check-idle-apache repository:

```
go build
```

## Additional notes

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[2]: https://github.com/sensu/sensu-plugin-sdk
[3]: https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md
[4]: https://github.com/mnx-solutions/check-idle-apache/blob/master/.github/workflows/release.yml
[5]: https://github.com/mnx-solutions/check-idle-apache/actions
[6]: https://docs.sensu.io/sensu-go/latest/reference/checks/
[7]: https://github.com/sensu/check-plugin-template/blob/master/main.go
[8]: https://bonsai.sensu.io/
[9]: https://github.com/sensu/sensu-plugin-tool
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
