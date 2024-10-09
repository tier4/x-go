# Changelog

## [0.16.0](https://www.github.com/tier4/x-go/compare/v0.15.0...v0.16.0) (2024-10-09)


### Features

* drop Go 1.19 to 1.21 supports and newly support Go 1.23 ([#73](https://www.github.com/tier4/x-go/issues/73)) ([db23e56](https://www.github.com/tier4/x-go/commit/db23e56a720ac290c9b204f9b12e101da3f264aa))

## [0.15.0](https://www.github.com/tier4/x-go/compare/v0.14.0...v0.15.0) (2024-09-27)


### Features

* **zstdx:** uncompress with custom size limit ([#70](https://www.github.com/tier4/x-go/issues/70)) ([fb2345e](https://www.github.com/tier4/x-go/commit/fb2345e14aa356ca6f7af9d47fba9b7131110e66))

## [0.14.0](https://www.github.com/tier4/x-go/compare/v0.13.0...v0.14.0) (2024-05-24)


### Features

* **dockertestx:** support SQS (ElasticMQ) ([#66](https://www.github.com/tier4/x-go/issues/66)) ([a2c292d](https://www.github.com/tier4/x-go/commit/a2c292dea4c48069ab18243bf5e78a8dc055e484))

## [0.13.0](https://www.github.com/tier4/x-go/compare/v0.12.0...v0.13.0) (2024-03-13)


### Features

* add Go 1.22 support ([#59](https://www.github.com/tier4/x-go/issues/59)) ([0037862](https://www.github.com/tier4/x-go/commit/0037862da7cd270c528c65a6f6d2f119c4e08946))
* introduce random package ([#58](https://www.github.com/tier4/x-go/issues/58)) ([61ff285](https://www.github.com/tier4/x-go/commit/61ff285ca303156e4d5b471e2e6855e6b69821ed))

## [0.12.0](https://www.github.com/tier4/x-go/compare/v0.11.0...v0.12.0) (2024-02-19)


### Features

* **dockertestx:** Support Prism ([#56](https://www.github.com/tier4/x-go/issues/56)) ([05ab7ef](https://www.github.com/tier4/x-go/commit/05ab7ef059b650921ddc6311bdeb70d62b585ba1))

## [0.11.0](https://www.github.com/tier4/x-go/compare/v0.10.0...v0.11.0) (2023-11-14)


### Features

* drop support Go v1.18  and add support Go v1.19 later ([#52](https://www.github.com/tier4/x-go/issues/52)) ([76f6b7c](https://www.github.com/tier4/x-go/commit/76f6b7c8308441a52ef89fe66e934d7b87ae5394))

## [0.10.0](https://www.github.com/tier4/x-go/compare/v0.9.0...v0.10.0) (2023-07-14)


### Features

* add zstdx package ([#49](https://www.github.com/tier4/x-go/issues/49)) ([001b422](https://www.github.com/tier4/x-go/commit/001b42293c4b9a8262876884751963d8ff7655f5))

## [0.9.0](https://www.github.com/tier4/x-go/compare/v0.8.0...v0.9.0) (2022-05-12)


### Features

* Utility function: Ref and Deref ([#42](https://www.github.com/tier4/x-go/issues/42)) ([cd34ba9](https://www.github.com/tier4/x-go/commit/cd34ba9722cd0f0aac1561230d59fdcfe37a9005))

## [0.8.0](https://www.github.com/tier4/x-go/compare/v0.7.0...v0.8.0) (2022-03-01)


### ⚠ BREAKING CHANGES

* Drop Go 1.16 support (#38)
* Remove popx (#36)

### Bug Fixes

* Drop Go 1.16 support ([#38](https://www.github.com/tier4/x-go/issues/38)) ([ef30f00](https://www.github.com/tier4/x-go/commit/ef30f00eca5acc25b85054f2ee3c2e85a8f4e797))
* Remove popx ([#36](https://www.github.com/tier4/x-go/issues/36)) ([3bffe78](https://www.github.com/tier4/x-go/commit/3bffe782c2eceee47c0539e7f7cc224f72d5fa6a))

## [0.7.0](https://www.github.com/tier4/x-go/compare/v0.6.0...v0.7.0) (2022-02-28)


### ⚠ BREAKING CHANGES

* Bump pop to v6 (#34)

### Bug Fixes

* Bump pop to v6 ([#34](https://www.github.com/tier4/x-go/issues/34)) ([a1fdcc2](https://www.github.com/tier4/x-go/commit/a1fdcc2d367f4a2f002f6cabd57d7bcf1637f321))

## [0.6.0](https://www.github.com/tier4/x-go/compare/v0.5.2...v0.6.0) (2022-01-15)


### Features

* New ORM wrapper; bunx ([#29](https://www.github.com/tier4/x-go/issues/29)) ([0c9d526](https://www.github.com/tier4/x-go/commit/0c9d5265883e1a0c94ed632699391586bc5c93fc))


### Bug Fixes

* Use stdlib errors instead of github.com/pkg/errors ([#32](https://www.github.com/tier4/x-go/issues/32)) ([7e62d2e](https://www.github.com/tier4/x-go/commit/7e62d2e5f854652e5435188afa807e91a99f9246))

### [0.5.2](https://www.github.com/tier4/x-go/compare/v0.5.1...v0.5.2) (2021-10-22)


### Bug Fixes

* CVE-2021-42576 github.com/microcosm-cc/bluemonday ([#27](https://www.github.com/tier4/x-go/issues/27)) ([62884e5](https://www.github.com/tier4/x-go/commit/62884e50964b8fb04e2a12f21561a0ccada4a2e1))

### [0.5.1](https://www.github.com/tier4/x-go/compare/v0.5.0...v0.5.1) (2021-08-05)


### Bug Fixes

* **dockertestx:** remove unnecessary log ([#22](https://www.github.com/tier4/x-go/issues/22)) ([d467ea2](https://www.github.com/tier4/x-go/commit/d467ea231ec8037f4b5c9bacbffa0290ba27eaa8))
* **popx:** Connection.Store compatible with sqlx.QueryerContext ([#25](https://www.github.com/tier4/x-go/issues/25)) ([0e499fd](https://www.github.com/tier4/x-go/commit/0e499fd4ecb4da0600bf53fd76706bc4c2824b06))

## [0.5.0](https://www.github.com/tier4/x-go/compare/v0.4.0...v0.5.0) (2021-05-31)


### Features

* **dockertestx:** S3 ([#18](https://www.github.com/tier4/x-go/issues/18)) ([5741787](https://www.github.com/tier4/x-go/commit/5741787f2e6e45a0d0cfe0bffef8f2cb4d935472))
* **popx:** Reset migration helper ([#20](https://www.github.com/tier4/x-go/issues/20)) ([7333ad4](https://www.github.com/tier4/x-go/commit/7333ad404fe511df33248890d6de468bf047997a))

## [0.4.0](https://www.github.com/tier4/x-go/compare/v0.3.2...v0.4.0) (2021-05-17)


### Features

* **dockertestx:** Save state to reuse container ([#13](https://www.github.com/tier4/x-go/issues/13)) ([dc78bca](https://www.github.com/tier4/x-go/commit/dc78bca42a409b92627a19163453bb7f5516d391))


### Bug Fixes

* **dockertestx:** Ignore not found container ([#14](https://www.github.com/tier4/x-go/issues/14)) ([d186536](https://www.github.com/tier4/x-go/commit/d186536fa943f8c75a3ae0775e81e21cdd05b5e6))
* **dockertestx:** Ignore stopped container ([#15](https://www.github.com/tier4/x-go/issues/15)) ([2784bc7](https://www.github.com/tier4/x-go/commit/2784bc7869428a1093fdb7498a1d7c510745136e))

### [0.3.2](https://www.github.com/tier4/x-go/compare/v0.3.1...v0.3.2) (2021-05-14)


### Bug Fixes

* **dockertestx:** Use temporary environment ([#11](https://www.github.com/tier4/x-go/issues/11)) ([baf72bf](https://www.github.com/tier4/x-go/commit/baf72bfc2d19ca0e4a5f6cf2c1aae96bb5d809c1))

### [0.3.1](https://www.github.com/tier4/x-go/compare/v0.3.0...v0.3.1) (2021-05-14)


### Bug Fixes

* **dockertestx:** Stub directory ([#9](https://www.github.com/tier4/x-go/issues/9)) ([5d2f7dd](https://www.github.com/tier4/x-go/commit/5d2f7ddda495beb90901a8e40af8fc13c5da9bc4))

## [0.3.0](https://www.github.com/tier4/x-go/compare/v0.2.1...v0.3.0) (2021-05-13)


### Features

* **dockertestx:** DynamoDB ([#7](https://www.github.com/tier4/x-go/issues/7)) ([c95c9e4](https://www.github.com/tier4/x-go/commit/c95c9e4afe1cb4e74ed309d1687cb2dcf3e8f2c0))

### [0.2.1](https://www.github.com/tier4/x-go/compare/v0.2.0...v0.2.1) (2021-05-13)


### Bug Fixes

* Rename package name ([#5](https://www.github.com/tier4/x-go/issues/5)) ([324c934](https://www.github.com/tier4/x-go/commit/324c934245074ae3d1fbfef42b2f9d00df15acc3))

## [0.2.0](https://www.github.com/tier4/x-go/compare/v0.1.0...v0.2.0) (2021-05-11)


### Features

* popx and dockertestx for database connection ([#3](https://www.github.com/tier4/x-go/issues/3)) ([63e2636](https://www.github.com/tier4/x-go/commit/63e2636d373d59aa9075d6759ab0741d64cc5bb6))

## 0.1.0 (2021-05-10)


### Bug Fixes

* CI token ([1a98085](https://www.github.com/tier4/x-go/commit/1a9808515b2592666acb0a6eb079ab983cfedcfb))
