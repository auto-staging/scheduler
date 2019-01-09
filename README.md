# Auto Staging Scheduler

[![Maintainability](https://api.codeclimate.com/v1/badges/4081e8c1d9f05200133d/maintainability)](https://codeclimate.com/github/auto-staging/scheduler/maintainability)
[![GoDoc](https://godoc.org/github.com/auto-staging/scheduler?status.svg)](https://godoc.org/github.com/auto-staging/scheduler)
[![Go Report Card](https://goreportcard.com/badge/github.com/auto-staging/scheduler)](https://goreportcard.com/report/github.com/auto-staging/scheduler)

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

## License and Author

Author: Jan Ritter

License: MIT