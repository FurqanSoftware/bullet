# Commands

## setup

Prepare servers for application deployment. Installs Docker and creates the application directory structure.

```sh
bullet -H 192.168.0.3 setup
bullet -H 192.168.0.3 setup --environ env.production
```

| Flag        | Description                                        |
|-------------|----------------------------------------------------|
| `--environ` | Path to an environment file to push during setup.  |

What it does:

1. Installs Docker (skipped if already installed).
2. Creates `/opt/<identifier>/releases/` directory.
3. Creates an empty `/opt/<identifier>/env` file.
4. Optionally uploads the environment file.

## deploy

Package and deploy the application to servers.

```sh
bullet -H 192.168.0.3 deploy app.tar.gz
```

The argument is a path to a tarball containing your application files.

What it does:

1. Computes the SHA256 hash of the tarball.
2. For each selected node:
   - Checks if this release is already deployed (by comparing hashes). Skips if so.
   - Uploads the tarball to `/tmp/` on the server.
   - Extracts it to `/opt/<identifier>/releases/<timestamp>-<hash>/`.
   - Removes the temporary tarball.
   - Updates `/opt/<identifier>/current` (via symlink or copy, based on `deploy.current`).
   - Writes the hash to `/opt/<identifier>/current.hash`.
   - Builds Docker images for programs that have a `dockerfile` defined.
   - Reloads running containers using the configured [reload method](bulletspec.md#reload).
   - Prunes old releases, keeping the 5 most recent.

## status

Print the status of all application containers across nodes.

```sh
bullet -H 192.168.0.3 status
```

Outputs a table showing how many instances of each program are running and healthy on each node, with totals.

## restart

Restart all application containers on selected nodes.

```sh
bullet -H 192.168.0.3 restart
```

Stops and recreates each running container.

## run

Run a program as a one-off interactive container.

```sh
bullet -H 192.168.0.3 run web
```

The argument is a program key from the Bulletspec. The container is removed after it exits. This is useful for running management commands, debugging, or one-off tasks.

## scale

Scale program instances on selected nodes.

```sh
# Scale to specific counts
bullet -H 192.168.0.3 scale web=4 worker=2

# Use default scaling rules from Bulletspec
bullet -H 192.168.0.3 scale
```

When called with arguments, each argument is in `program=count` format.

When called without arguments, Bullet evaluates the [scaling rules](bulletspec.md#scaling-rules) defined in the Bulletspec against each node's properties (tags, hardware) to determine the desired instance counts.

Bullet reports how many containers were created or removed:

```
Scaled program web
∟ Desired: 4
∟ Ready: 4 (+2)
```

## cron:enable

Enable one or more cron jobs on selected nodes.

```sh
bullet -H 192.168.0.3 cron:enable cleanup
bullet -H 192.168.0.3 cron:enable cleanup report
```

Creates systemd timer and service units for each job and enables them. The timer triggers a Docker container that runs the job's command.

If the job has a `healthcheck.url`, the service pings `{url}/start` before the job and `{url}/{exit_status}` after.

If the job has a `jitter`, the timer uses `RandomizedDelaySec` with `FixedRandomDelay` to spread execution.

## cron:disable

Disable one or more cron jobs on selected nodes.

```sh
bullet -H 192.168.0.3 cron:disable cleanup
```

Stops the systemd timer, removes the timer and service unit files, and reloads the systemd daemon.

## cron:status

Print the status of all cron jobs.

```sh
bullet -H 192.168.0.3 cron:status
```

Shows whether each job's timer is active or disabled, along with the next trigger time.

## environ:push

Upload an environment file to selected nodes.

```sh
bullet -H 192.168.0.3 environ:push env.production
```

The file is uploaded to `/opt/<identifier>/env` on the server. All containers and cron jobs use this file for environment variables via Docker's `--env-file` flag.

## log

Tail container logs for a program.

```sh
# Tail instance 1 (default)
bullet -H 192.168.0.3 log web

# Tail instance 3
bullet -H 192.168.0.3 log web:3
```

Shows the last 10 lines, then streams new log output until interrupted.

## forward

Forward a port from a remote server to your local machine.

```sh
# Forward local port 8080 to remote port 8080
bullet -H 192.168.0.3 forward 8080

# Forward local port 3000 to remote port 8080
bullet -H 192.168.0.3 forward 3000:8080
```

The forwarding stays active until interrupted. Multiple connections are supported concurrently.

## prune

Remove old releases to free disk space.

```sh
bullet -H 192.168.0.3 prune
```

Keeps the 5 most recent releases in `/opt/<identifier>/releases/` and deletes the rest. This is also run automatically at the end of each deploy.

## host:shell

Open an interactive shell on a remote server.

```sh
bullet -H 192.168.0.3 host:shell
```

Starts a bash session over SSH with full PTY support.

## host:df

Show disk space usage on a remote server.

```sh
bullet -H 192.168.0.3 host:df
bullet -H 192.168.0.3 host:df --watch
bullet -H 192.168.0.3 host:df --arguments "-h"
```

| Flag          | Short | Description                             |
|---------------|-------|-----------------------------------------|
| `--watch`     | `-w`  | Continuously update output with `watch`.|
| `--arguments` | `-a`  | Additional arguments to pass to `df`.   |

## host:top

Show running processes on a remote server.

```sh
bullet -H 192.168.0.3 host:top
```

Runs `top` interactively with full PTY support.
