# Bulletspec Reference

The `Bulletspec` file is a YAML file placed in the root of your project. It defines your application, its programs, cron jobs, and how they should be built, deployed, and run.

## Application

The top-level key is `application`:

```yaml
application:
  name: My App
  identifier: myapp
```

| Field        | Type   | Description                                                                 |
|--------------|--------|-----------------------------------------------------------------------------|
| `name`       | string | Display name of the application.                                            |
| `identifier` | string | Unique identifier. Used in file paths, container names, and image names.    |
| `deploy`     | object | Deployment behavior. See [Deploy](#deploy).                                 |
| `programs`   | map    | Named programs to run as containers. See [Programs](#programs).             |
| `cron`       | object | Scheduled jobs. See [Cron Jobs](#cron-jobs).                                |

## Deploy

Controls how releases are made current on the server.

```yaml
application:
  deploy:
    current: symlink
```

| Field     | Type   | Default     | Description                                                        |
|-----------|--------|-------------|--------------------------------------------------------------------|
| `current` | string | `"symlink"` | `"symlink"` creates a symlink to the release. `"replace"` copies release contents into the current directory. |

## Programs

Each key in the `programs` map defines a long-running service. The key is used as the program identifier.

```yaml
application:
  programs:
    web:
      name: Web Server
      command: node server.js
      container:
        image: node:20-alpine
      ports:
        - 80:3000
    worker:
      name: Background Worker
      command: node worker.js
      container:
        dockerfile: Dockerfile.worker
```

### Program Fields

| Field         | Type     | Description                                                                 |
|---------------|----------|-----------------------------------------------------------------------------|
| `name`        | string   | Display name of the program.                                                |
| `command`     | string   | Command to run inside the container. Supports variable expansion (see [Variables](configuration.md#variables)). |
| `user`        | string   | User to run the container as.                                               |
| `container`   | object   | Container configuration. See [Container](#container).                       |
| `ports`       | []string | Port mappings in `HOST:CONTAINER` format. Host port is incremented per instance (e.g. instance 2 of `80:3000` maps to `81:3000`). |
| `volumes`     | []string | Volume mounts passed directly to Docker.                                    |
| `healthcheck` | object   | Container health check. See [Healthcheck](#healthcheck).                    |
| `scales`      | []object | Expression-based scaling rules. See [Scaling Rules](#scaling-rules).        |
| `reload`      | object   | How to reload on deploy. See [Reload](#reload).                             |
| `unsafe`      | object   | Unsafe flags. See [Unsafe](#unsafe).                                        |

### Container

Defines the Docker image or build configuration for a program.

```yaml
container:
  image: node:20-alpine
```

Or with a custom Dockerfile:

```yaml
container:
  dockerfile: Dockerfile.web
  entrypoint: "/entrypoint.sh"
  workingdir: /app
  applicationdir: /app
```

| Field            | Type   | Default            | Description                                                          |
|------------------|--------|--------------------|----------------------------------------------------------------------|
| `image`          | string |                    | Base Docker image (e.g. `node:20-alpine`).                           |
| `dockerfile`     | string |                    | Path to a Dockerfile. If set, the image is built from this file.     |
| `entrypoint`     | string |                    | Custom entrypoint for the container.                                 |
| `workingdir`     | string | `/<identifier>`    | Working directory inside the container.                              |
| `applicationdir` | string | `/<identifier>`    | Where the application directory is mounted inside the container.     |

When `dockerfile` is specified, Bullet tracks the Dockerfile's SHA256 hash. The image is only rebuilt when the Dockerfile changes.

### Healthcheck

Configures Docker's built-in health check for a program.

```yaml
healthcheck:
  command: curl -f http://localhost:3000/health
  interval: 30s
  timeout: 5s
  retries: 3
  startperiod: 10s
```

| Field         | Type     | Description                                      |
|---------------|----------|--------------------------------------------------|
| `command`     | string   | Health check command to run inside the container. |
| `interval`    | duration | How often to run the check.                      |
| `timeout`     | duration | Maximum time for a single check.                 |
| `retries`     | int      | Failures needed to mark unhealthy.               |
| `startperiod` | duration | Grace period after container start.              |

### Scaling Rules

Scaling rules let you define the default number of instances per program based on the target node's properties.

```yaml
scales:
  - if: 'hasTags("production")'
    n: 'hw.cores'
  - n: '2'
```

Rules are evaluated in order. The last matching rule wins.

| Field | Type   | Description                                                             |
|-------|--------|-------------------------------------------------------------------------|
| `if`  | string | Boolean expression. If omitted, the rule always applies.               |
| `n`   | string | Expression that evaluates to the desired instance count (must be int). |

Expressions have access to:

| Name                    | Type     | Description                                   |
|-------------------------|----------|-----------------------------------------------|
| `hasTags(tag1, tag2..)` | function | Returns true if the node has all listed tags.  |
| `hw.cores`              | int      | Number of CPU cores on the node.               |
| `hw.memory`             | int      | Memory in MB on the node.                      |

### Reload

Controls how containers are updated during a deploy.

```yaml
reload:
  method: signal
  signal: SIGHUP
  precommand: nginx -t
```

| Field        | Type   | Default     | Description                                                     |
|--------------|--------|-------------|-----------------------------------------------------------------|
| `method`     | string | `"restart"` | `"signal"`, `"command"`, or `"restart"`.                        |
| `signal`     | string |             | Signal to send when method is `"signal"` (e.g. `SIGHUP`).      |
| `command`    | string |             | Command to exec in container when method is `"command"`.        |
| `precommand` | string |             | Command to exec in container before the reload action.          |

If the Docker image was rebuilt during the deploy, the container is always restarted regardless of the reload method.

### Unsafe

```yaml
unsafe:
  networkhost: true
  ulimits:
    - nofile=65535:65535
    - memlock=-1:-1
```

| Field         | Type     | Description                                              |
|---------------|----------|----------------------------------------------------------|
| `networkhost` | bool     | Run container with `--network=host`. Default false.      |
| `ulimits`     | []string | Set container ulimits via `--ulimit` (e.g. `nofile=65535:65535`). |

## Cron Jobs

Cron jobs run as one-off Docker containers on a schedule, managed by systemd timers.

```yaml
application:
  cron:
    jobs:
      - key: cleanup
        command: node cleanup.js
        schedule: daily
        jitter: 30m
        healthcheck:
          url: https://hc-ping.com/abc123
```

### Job Fields

| Field         | Type   | Description                                                                      |
|---------------|--------|----------------------------------------------------------------------------------|
| `key`         | string | Unique identifier for the job.                                                   |
| `command`     | string | Command to run in the container.                                                 |
| `schedule`    | string | systemd OnCalendar expression (e.g. `daily`, `hourly`, `*-*-* 02:00:00`).       |
| `jitter`      | string | Random delay added to the schedule (e.g. `30m`). Uses systemd RandomizedDelaySec with FixedRandomDelay. |
| `healthcheck` | object | Optional healthcheck pings.                                                      |

### Job Healthcheck

| Field | Type   | Description                                                                                            |
|-------|--------|--------------------------------------------------------------------------------------------------------|
| `url` | string | Base URL for health pings. Bullet calls `{url}/start` before the job and `{url}/{exit_status}` after. |

## Full Example

```yaml
application:
  name: My Web App
  identifier: mywebapp
  deploy:
    current: symlink

  programs:
    web:
      name: Web Server
      command: ./server --port 3000
      container:
        dockerfile: Dockerfile
      ports:
        - 80:3000
      healthcheck:
        command: curl -f http://localhost:3000/health
        interval: 30s
        timeout: 5s
        retries: 3
        startperiod: 10s
      scales:
        - if: 'hasTags("production")'
          n: 'hw.cores'
        - n: '1'
      reload:
        method: signal
        signal: SIGHUP
        precommand: ./server --check-config

    worker:
      name: Background Worker
      command: ./worker
      container:
        image: mywebapp-base
      scales:
        - n: '2'

  cron:
    jobs:
      - key: cleanup
        command: ./cleanup
        schedule: daily
        jitter: 1h
        healthcheck:
          url: https://hc-ping.com/abc123
```
