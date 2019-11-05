# Pinkie

[![GitHub Actions](https://github.com/inabagumi/pinkie/workflows/Go/badge.svg)](https://github.com/inabagumi/pinkie/actions?query=workflow%3AGo) [![Codecov](https://codecov.io/gh/inabagumi/pinkie/branch/master/graph/badge.svg)](https://codecov.io/gh/inabagumi/pinkie)

## Installation

```console
$ export PINKIE_VERSION=3.3.0
$ cd $(mktemp -d)
$ curl -sSL https://github.com/inabagumi/pinkie/releases/download/v${PINKIE_VERSION}/ytc_${PINKIE_VERSION}_Linux_x86_64.tar.gz | tar xzf -
$ sudo install pinkie /usr/local/bin
```

## Usage

```console
$ export GOOGLE_API_KEY=xxxxx
$ export ALGOLIA_APPLICATION_ID=xxxxx
$ export ALGOLIA_API_KEY=xxxxx
$ export ALGOLIA_INDEX_NAME=xxxxx
$ pinkie --channel UC0Owc36U9lOyi9Gx9Ic-4qg
```

## LICENSE

[MIT](LICENSE)
