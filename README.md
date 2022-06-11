# Pinkie

[![GitHub Actions](https://github.com/inabagumi/pinkie/workflows/Go/badge.svg)](https://github.com/inabagumi/pinkie/actions?query=workflow%3AGo) [![Codecov](https://codecov.io/gh/inabagumi/pinkie/branch/main/graph/badge.svg)](https://codecov.io/gh/inabagumi/pinkie)

## Usage

```console
$ docker pull ghcr.io/inabagumi/pinkie:latest
$ echo GOOGLE_API_KEY=xxxxx >> .env
$ echo ALGOLIA_APPLICATION_ID=xxxxx >> .env
$ echo ALGOLIA_API_KEY=xxxxx >> .env
$ echo ALGOLIA_INDEX_NAME=xxxxx >> .env
$ docker run --env-file .env --rm docker pull ghcr.io/inabagumi/pinkie:latest --channel UC0Owc36U9lOyi9Gx9Ic-4qg
```

## LICENSE

[MIT](LICENSE)
