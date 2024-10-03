# Changelog

## [1.8.1](https://github.com/antonchaban/news-aggregator/compare/news-alligator@v1.8.0...news-alligator@v1.8.1) (2024-10-03)


### Bug Fixes

* Remove helm charts installation from cdk ([5bdf1ed](https://github.com/antonchaban/news-aggregator/commit/5bdf1ed1478c9851ebdb8b2c6d8240bbd78fcf22))

## [1.8.0](https://github.com/antonchaban/news-aggregator/compare/news-alligator@v1.7.1...news-alligator@v1.8.0) (2024-10-01)


### Features

* Add AWS cloud formation for cluster creation ([9ba22d8](https://github.com/antonchaban/news-aggregator/commit/9ba22d8cd7614ace2381559e86b8fe6956073309))

## [1.7.1](https://github.com/antonchaban/news-aggregator/compare/news-alligator@v1.7.0...news-alligator@v1.7.1) (2024-10-01)


### Bug Fixes

* Update ECR docker image name ([c78c557](https://github.com/antonchaban/news-aggregator/commit/c78c557cb16cc70e85c3842c7a3170dd42d8b0c0))

## [1.7.0](https://github.com/antonchaban/news-aggregator/compare/news-alligator@v1.6.0...news-alligator@v1.7.0) (2024-10-01)


### Features

* Create a k8s CronJob which creates and updates a k8s Secret with credentials to ECR; ([9befef3](https://github.com/antonchaban/news-aggregator/commit/9befef3b3affcbb3ada0bd6c0a3aaacf17e6f00f))
* create storages for all docker images and helm charts in private ECR ([dce631f](https://github.com/antonchaban/news-aggregator/commit/dce631f77aa4ccbbc23ad79cb6f69786b9f56f87))
* Update all k8s deployments to use the secret to pull images from ECR ([bff6f9b](https://github.com/antonchaban/news-aggregator/commit/bff6f9bec55fa7ce6a3f1d4475c74d3f5007941c))
* Update GithubActions to push image to ECR on it's release; ([cd6cf80](https://github.com/antonchaban/news-aggregator/commit/cd6cf800964c45d72b5e367acd52087e13b1b5aa))
* Update Taskfile to be able to push images/charts to ECR; ([e4c2088](https://github.com/antonchaban/news-aggregator/commit/e4c2088daf40aeb7fb2aa3002662fe25f5f25075))


### Bug Fixes

* Add task tool installation ([3fbea2d](https://github.com/antonchaban/news-aggregator/commit/3fbea2d8cdf9d116df4bd26db0c03add36bb2371))
* update release please to use version from release ([43501f5](https://github.com/antonchaban/news-aggregator/commit/43501f5be47d8d8a1e96805957e793d03214ec58))

## [1.6.0](https://github.com/antonchaban/news-aggregator/compare/news-alligator@v1.5.0...news-alligator@v1.6.0) (2024-10-01)


### Features

* Add scaling features for news aggregator application ([9093f69](https://github.com/antonchaban/news-aggregator/commit/9093f69944fb2584632705f5926828ec141565e5))

## [1.5.0](https://github.com/antonchaban/news-aggregator/compare/news-alligator@v1.4.0...news-alligator@v1.5.0) (2024-10-01)


### Features

* Facilitate the management of certificates with cert-manager ([3c79b2a](https://github.com/antonchaban/news-aggregator/commit/3c79b2acff8d716bd7a6ba70d40c35931ca39f3f))

## [1.4.0](https://github.com/antonchaban/news-aggregator/compare/news-alligator@v1.3.0...news-alligator@v1.4.0) (2024-10-01)


### Features

* Add endpoint to news aggregator to view all available sources ([df08aa0](https://github.com/antonchaban/news-aggregator/commit/df08aa02cf9f1e98e665344862fdbc73b904d44e))

## [1.3.0](https://github.com/antonchaban/news-aggregator/compare/news-alligator@v1.2.0...news-alligator@v1.3.0) (2024-10-01)


### Features

* Add helm for managing k8s manifests ([8f8352c](https://github.com/antonchaban/news-aggregator/commit/8f8352c3cb42ddb6f92b807d6783557506ff976f))

## [1.2.0](https://github.com/antonchaban/news-aggregator/compare/news-alligator@v1.1.0...news-alligator@v1.2.0) (2024-09-23)


### Features

* Add DB and separate fetcher service for news fetching as cron job ([12f0d27](https://github.com/antonchaban/news-aggregator/commit/12f0d27d2c8a38964b44cdb0407e8f61d5fc8eb5))

## [1.1.0](https://github.com/antonchaban/news-aggregator/compare/news-alligator@v1.0.1...news-alligator@v1.1.0) (2024-09-23)


### Features

* Add volumes to k8s cluster ([f77b924](https://github.com/antonchaban/news-aggregator/commit/f77b9245fc4e5eb45100addb096dfc1519948a2a))

## [1.0.1](https://github.com/antonchaban/news-aggregator/compare/news-alligator-v1.0.0...news-alligator@v1.0.1) (2024-08-01)


### Bug Fixes

* Add certificates ([7ba19a0](https://github.com/antonchaban/news-aggregator/commit/7ba19a04825e9feb7e2ffc0b2d035bf63b3a8649))
* add if statements for docker ([a77a5a8](https://github.com/antonchaban/news-aggregator/commit/a77a5a815f5f0894f8532492187d1e47f7691d04))
* mark tasks as internal: true ([c57578a](https://github.com/antonchaban/news-aggregator/commit/c57578a4eb90a5fb8a6a18eb3b51217a1d6adbb1))
* Modify workflow to upload image to dockerhub ([7607195](https://github.com/antonchaban/news-aggregator/commit/7607195f4752b1667f399a823d1976f9ef7374ec))
* Update branch name to master ([682c918](https://github.com/antonchaban/news-aggregator/commit/682c91847702627cae52e2e8a3bd5821c691e686))
* Update dockerfile ([d9b2b9d](https://github.com/antonchaban/news-aggregator/commit/d9b2b9d4010bffe8f050ec51b30e8c8c84426a80))
* Update dockerfile to copy only required dirs and fix wrong ENV syntax ([74adf2c](https://github.com/antonchaban/news-aggregator/commit/74adf2cc8e7288f4bfc2b4a9967b1ea619347347))
* update go.yml that build tasks needs lint ([f20a8c4](https://github.com/antonchaban/news-aggregator/commit/f20a8c4da62e08e85f908383ee0eeddfbd3e64fd))
* update lint task ([65e2522](https://github.com/antonchaban/news-aggregator/commit/65e252288a2994bda9d203b9ec5dd78a37255cc4))
* Update release please configuration ([931ba79](https://github.com/antonchaban/news-aggregator/commit/931ba793f40d3ae836ebd929d33c5690207f04f3))
* Update Taskfile.yml ([0749a7d](https://github.com/antonchaban/news-aggregator/commit/0749a7d9201812bae13040db6bce4fbacf06caa3))
