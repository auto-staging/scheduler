# Auto Staging Scheduler

> Scheduler gets invoked by CloudWatchEvents rules or the Tower Lambda function, it starts and stops EC2 Instances and RDS Clusters for the given
> repository and branch (Environment)

## CloudWatch Event Bodys

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