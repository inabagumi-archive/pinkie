# Changelog

## [4.3.6](https://github.com/inabagumi/pinkie/compare/v4.3.5...v4.3.6) (2022-06-11)


### Bug Fixes

* **terraform:** add missing role for gha ([#268](https://github.com/inabagumi/pinkie/issues/268)) ([0994ee6](https://github.com/inabagumi/pinkie/commit/0994ee65cd6f4d6785fc9dee92aa485221366fe7))
* **terraform:** grant reader role to all users ([#269](https://github.com/inabagumi/pinkie/issues/269)) ([61e9504](https://github.com/inabagumi/pinkie/commit/61e9504f64517b4ff3fe5da19fbc89562183dd51))
* **terraform:** use POST method ([#266](https://github.com/inabagumi/pinkie/issues/266)) ([2899d6f](https://github.com/inabagumi/pinkie/commit/2899d6fe7152598331b2c565423f5576302aefcd))

## [4.3.5](https://github.com/inabagumi/pinkie/compare/v4.3.4...v4.3.5) (2022-06-10)


### Bug Fixes

* **github-actions:** ignore `gha-creds-*.json` ([#262](https://github.com/inabagumi/pinkie/issues/262)) ([cdd8335](https://github.com/inabagumi/pinkie/commit/cdd8335d9cd4838a4f72819b0e229bc0bce1495f))

## [4.3.4](https://github.com/inabagumi/pinkie/compare/v4.3.3...v4.3.4) (2022-06-10)


### Bug Fixes

* **github-actions:** add `fetch-depth: 0` ([#260](https://github.com/inabagumi/pinkie/issues/260)) ([43dff01](https://github.com/inabagumi/pinkie/commit/43dff01dc0e023373a6d32a27bc66843a9becae7))

## [4.3.3](https://github.com/inabagumi/pinkie/compare/v4.3.2...v4.3.3) (2022-06-10)


### Bug Fixes

* **goreleaser:** escape quote ([#259](https://github.com/inabagumi/pinkie/issues/259)) ([2d47d7d](https://github.com/inabagumi/pinkie/commit/2d47d7d88790389a3ed7f41e52b15d90c3ab7a4d))
* **terraform:** binding roles to service account for terraform ([#257](https://github.com/inabagumi/pinkie/issues/257)) ([40702b5](https://github.com/inabagumi/pinkie/commit/40702b5cc1f52a6038a18ca0a89fa2eaf5a321ec))

## [4.3.2](https://github.com/inabagumi/pinkie/compare/v4.3.1...v4.3.2) (2022-06-09)


### Bug Fixes

* **terraform:** fix invalid workload identity provider ([#255](https://github.com/inabagumi/pinkie/issues/255)) ([2cdfd97](https://github.com/inabagumi/pinkie/commit/2cdfd97e22630908e0dddc7c7067e76d99755230))

## [4.3.1](https://github.com/inabagumi/pinkie/compare/v4.3.0...v4.3.1) (2022-06-09)


### Bug Fixes

* **github-actions:** add missing args ([#253](https://github.com/inabagumi/pinkie/issues/253)) ([897076a](https://github.com/inabagumi/pinkie/commit/897076a2591a0a8a3b1d73d5142dc73e21a6db9d))

## [4.3.0](https://github.com/inabagumi/pinkie/compare/v4.2.5...v4.3.0) (2022-06-09)


### Features

* release to artifact registry ([#251](https://github.com/inabagumi/pinkie/issues/251)) ([c0fc07d](https://github.com/inabagumi/pinkie/commit/c0fc07d1567429e2987dd9bd310c36a6c3a6dd3d))

### [4.2.5](https://github.com/inabagumi/pinkie/compare/v4.2.4...v4.2.5) (2022-02-19)


### Bug Fixes

* **goreleaser:** use goreleaser v1.15 ([#219](https://github.com/inabagumi/pinkie/issues/219)) ([d830702](https://github.com/inabagumi/pinkie/commit/d8307029ffa79cb2b70f6cea511a63ef3e4988ae))

### [4.2.4](https://github.com/inabagumi/pinkie/compare/v4.2.3...v4.2.4) (2022-02-19)


### Bug Fixes

* **release-please:** re-add actions token ([#217](https://github.com/inabagumi/pinkie/issues/217)) ([f5f609c](https://github.com/inabagumi/pinkie/commit/f5f609cc5d860d1734fc9104414aa1c84827550e))

### [4.2.3](https://github.com/inabagumi/pinkie/compare/v4.2.2...v4.2.3) (2022-02-19)


### Bug Fixes

* **github-actions:** use github token ([#213](https://github.com/inabagumi/pinkie/issues/213)) ([b00630f](https://github.com/inabagumi/pinkie/commit/b00630feefa4ae1687fa6b91fb2aee1a09fbc468))
* **release-please:** add permission ([#216](https://github.com/inabagumi/pinkie/issues/216)) ([bad1da8](https://github.com/inabagumi/pinkie/commit/bad1da8c8f94d3749eced04da7aff5154e9f2f13))
* **release-please:** use github token ([#214](https://github.com/inabagumi/pinkie/issues/214)) ([23563ce](https://github.com/inabagumi/pinkie/commit/23563ce51af151ba27ceae067832af3a11716d23))

### [4.2.2](https://github.com/inabagumi/pinkie/compare/v4.2.1...v4.2.2) (2022-02-19)


### Bug Fixes

* **github-actions:** fix password ([#211](https://github.com/inabagumi/pinkie/issues/211)) ([56cf750](https://github.com/inabagumi/pinkie/commit/56cf75077a5dbe75f5e1d4a667b46d5486b8d679))

### [4.2.1](https://github.com/inabagumi/pinkie/compare/v4.2.0...v4.2.1) (2022-02-19)


### Bug Fixes

* **github-actions:** use `docker/login-action` ([#209](https://github.com/inabagumi/pinkie/issues/209)) ([e95beea](https://github.com/inabagumi/pinkie/commit/e95beea6d4c6a53ff6fc60e6db0a6464739be0e7))

## [4.2.0](https://github.com/inabagumi/pinkie/compare/v4.1.2...v4.2.0) (2022-02-19)


### Features

* **deps:** update go to v1.17 ([#207](https://github.com/inabagumi/pinkie/issues/207)) ([9463044](https://github.com/inabagumi/pinkie/commit/946304444d0fd07375ae8b38f4da8e47c637b537))
