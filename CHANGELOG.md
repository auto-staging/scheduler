# [1.3.0](https://github.com/auto-staging/scheduler/compare/1.2.0...1.3.0) (2019-07-22)


### Features

* added autoscaling group start and stop calls to main handler function ([b0cbcca](https://github.com/auto-staging/scheduler/commit/b0cbcca))
* added option to "start" and "stop" autoscaling groups by adapting the min instance number ([e2695ca](https://github.com/auto-staging/scheduler/commit/e2695ca))
* added option to start and stop autoscaling groups by lowering the min number to zero and setting it back to the previous min value on start ([14333a4](https://github.com/auto-staging/scheduler/commit/14333a4))

# [1.2.0](https://github.com/auto-staging/scheduler/compare/1.1.1...1.2.0) (2019-04-01)


### Features

* added endpoint for the tower to get the scheduler version information ([1212d44](https://github.com/auto-staging/scheduler/commit/1212d44))

## [1.1.1](https://github.com/auto-staging/scheduler/compare/1.1.0...1.1.1) (2019-03-28)


### Bug Fixes

* compile binary for linux ([c5ea20f](https://github.com/auto-staging/scheduler/commit/c5ea20f))

# [1.1.0](https://github.com/auto-staging/scheduler/compare/1.0.0...1.1.0) (2019-02-22)


### Bug Fixes

* fixed invalid clusterARN empty check ([0eb1d57](https://github.com/auto-staging/scheduler/commit/0eb1d57))


### Features

* added version info output ([dadd9b9](https://github.com/auto-staging/scheduler/commit/dadd9b9))

# 1.0.0 (2019-01-29)


### Features

* added start and stop for rds cluster ([a91516d](https://github.com/auto-staging/scheduler/commit/a91516d))
* improved error handling and success output ([c421fd3](https://github.com/auto-staging/scheduler/commit/c421fd3))
* improved logging output and parse cloudwatch event to struct ([b2a0406](https://github.com/auto-staging/scheduler/commit/b2a0406))
* project init ([8150da4](https://github.com/auto-staging/scheduler/commit/8150da4))
* start / stop ec2 instances based on repo and branch provided in the cloudwatch event ([11adb07](https://github.com/auto-staging/scheduler/commit/11adb07))
* update environment status after starting of stopping ([a442b35](https://github.com/auto-staging/scheduler/commit/a442b35))
