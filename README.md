# Auto Staging Scheduler

[![Maintainability](https://api.codeclimate.com/v1/badges/4081e8c1d9f05200133d/maintainability)](https://codeclimate.com/github/auto-staging/scheduler/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/4081e8c1d9f05200133d/test_coverage)](https://codeclimate.com/github/auto-staging/scheduler/test_coverage)
[![GoDoc](https://godoc.org/github.com/auto-staging/scheduler?status.svg)](https://godoc.org/github.com/auto-staging/scheduler)
[![Go Report Card](https://goreportcard.com/badge/github.com/auto-staging/scheduler)](https://goreportcard.com/report/github.com/auto-staging/scheduler)
[![Build Status](https://travis-ci.com/auto-staging/scheduler.svg?branch=master)](https://travis-ci.com/auto-staging/scheduler)

> Scheduler gets invoked by CloudWatchEvents rules or the Tower Lambda function, it starts and stops EC2 Instances and RDS Clusters for the given
> repository and branch (Environment)

## CloudWatchEvents Bodys

### Start

```json
{
    "repository": "demo-app",
    "branch": "feat/branch",
    "action": "start"
}
```

### Stop

```json
{
    "repository": "demo-app",
    "branch": "feat/branch",
    "action": "stop"
}
```

## Usage

### Install dependencies

Go dep must be installed

```bash
make prepare
```

### Execute tests

```bash
make tests
```

### Build binary

```bash
make build
```

compiles to bin/auto-staging-scheduler

## License and Author

Author: Jan Ritter

License: MIT