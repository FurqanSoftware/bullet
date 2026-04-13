# Deploying Applications

This guide walks through the full deployment lifecycle with Bullet.

## Preparing Your Application

Bullet deploys applications as tarballs. Package your application files into a `.tar.gz` archive:

```sh
tar czf app.tar.gz src/ package.json server.js
```

A Makefile target is a common pattern:

```makefile
.PHONY: release
release:
	tar czf app.tar.gz *.js package.json
```

## Writing a Bulletspec

Create a `Bulletspec` file in your project root:

```yaml
application:
  name: My App
  identifier: myapp

  programs:
    web:
      name: Web Server
      command: node server.js
      container:
        image: node:20-alpine
      ports:
        - 80:3000
```

See the [Bulletspec Reference](bulletspec.md) for all available options.

## Setting Up Servers

Before deploying, prepare your servers:

```sh
bullet -H 192.168.0.3 setup
```

This installs Docker and creates the directory structure under `/opt/myapp/`.

To also push an environment file during setup:

```sh
bullet -H 192.168.0.3 setup --environ env.production
```

## Deploying

```sh
tar czf app.tar.gz <your files>
bullet -H 192.168.0.3 deploy app.tar.gz
```

Bullet uploads the tarball, extracts it, builds Docker images, and reloads running containers. If the release has already been deployed (same SHA256 hash), the node is skipped.

You can combine setup, environment push, and scaling into a single deploy:

```sh
bullet -H 192.168.0.3 deploy app.tar.gz --setup --environ env.production --scale
```

## Scaling

After the first deploy, start your containers:

```sh
bullet -H 192.168.0.3 scale web=2
```

Or use expression-based scaling rules defined in your Bulletspec:

```yaml
programs:
  web:
    scales:
      - if: 'hasTags("production")'
        n: 'hw.cores'
      - n: '1'
```

Then run `scale` without arguments to apply the rules:

```sh
bullet -H 192.168.0.3 scale
```

## Environment Variables

Push environment files to your servers:

```sh
bullet -H 192.168.0.3 environ:push env.production
```

The file is stored at `/opt/<identifier>/env` and is automatically loaded by all containers and cron jobs.

## Custom Docker Images

For applications that need a custom image, specify a Dockerfile:

```yaml
programs:
  web:
    container:
      dockerfile: Dockerfile
```

Bullet builds the image on the server during deploy. It tracks the Dockerfile's hash and only rebuilds when the Dockerfile changes.

## Reload Strategies

By default, containers are restarted on deploy. You can configure a lighter reload:

```yaml
programs:
  web:
    reload:
      method: signal
      signal: SIGHUP
```

Or run a command inside the container:

```yaml
programs:
  web:
    reload:
      method: command
      command: nginx -s reload
      pre_command: nginx -t
```

The `pre_command` runs before the reload action, useful for configuration validation. If the Docker image was rebuilt during the deploy, the container is always restarted regardless of the reload method.

## Directory Structure on the Server

After setup and deployment, the server has:

```
/opt/<identifier>/
  releases/
    1711612800-abc123/    # Release directories (timestamp-hash)
    1711699200-def456/
  current -> releases/1711699200-def456/   # Symlink to latest (or copy)
  current.hash            # SHA256 of current release
  env                     # Environment file
```

## Container Details

Containers are named `<identifier>_<program>_<instance>`, for example `myapp_web_1`, `myapp_web_2`.

Each container gets the following environment variables set by Bullet:

| Variable                    | Value                          |
|-----------------------------|--------------------------------|
| `BULLET_APPLICATION_NAME`   | Application name from spec     |
| `BULLET_APPLICATION_ID`     | Application identifier         |
| `BULLET_PROGRAM_KEY`        | Program key (e.g. `web`)       |
| `BULLET_PROGRAM_NAME`       | Program display name           |
| `BULLET_INSTANCE_ID`        | Container name                 |

Plus all variables from the environment file.

Containers use the `json-file` log driver with a 3 GB max size and are configured with `--restart always`.

### Port Mapping

Port mappings increment the host port per instance. If the spec says `80:3000`:

| Instance | Host Port | Container Port |
|----------|-----------|----------------|
| 1        | 80        | 3000           |
| 2        | 81        | 3000           |
| 3        | 82        | 3000           |

## Multi-Node Deployments

Target multiple hosts at once:

```sh
bullet -H 192.168.0.3,192.168.0.4,192.168.0.5 deploy app.tar.gz
```

Or use a node manifest with tag filtering for larger setups:

```sh
# Deploy to all production web servers
bullet -c production deploy app.tar.gz
```

With `Bulletcfg.production`:

```yaml
hosts: "@nodes.yaml:production+web"
```

See [Configuration](configuration.md) for details on node manifests and tag filtering.

## Cron Jobs

Define scheduled tasks in your Bulletspec:

```yaml
application:
  cron:
    jobs:
      - key: cleanup
        command: ./cleanup
        schedule: daily
        jitter: 1h
        healthcheck:
          url: https://hc-ping.com/abc123
```

Enable them on your servers:

```sh
bullet -H 192.168.0.3 cron:enable cleanup
```

Cron jobs run as one-off Docker containers managed by systemd timers. They use the same environment file and application directory as regular programs.

Check their status:

```sh
bullet -H 192.168.0.3 cron:status
```

Disable when no longer needed:

```sh
bullet -H 192.168.0.3 cron:disable cleanup
```

## Monitoring

Check container status across nodes:

```sh
bullet status
```

Tail logs:

```sh
bullet -H 192.168.0.3 log web
bullet -H 192.168.0.3 log web:2    # Instance 2
```

Forward a port for local debugging:

```sh
bullet -H 192.168.0.3 forward 3000:8080
```

Access the server directly:

```sh
bullet -H 192.168.0.3 host:shell
bullet -H 192.168.0.3 host:df
bullet -H 192.168.0.3 host:top
```

## Cleanup

Remove old releases to free disk space:

```sh
bullet -H 192.168.0.3 prune
```

This keeps the 5 most recent releases. Pruning also happens automatically after each deploy.
