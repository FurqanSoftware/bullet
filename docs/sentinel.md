# Sentinel

Bullet Sentinel is a small daemon that runs on each Bullet-managed host and
exposes a Prometheus-format `/metrics` endpoint. It self-discovers Bullet's
containers and cron timers via `docker` and `systemctl` — there's no spec
pushed to the host; the sentinel introspects what's actually running.

The sentinel is shipped as a separate binary embedded inside `bullet` itself
and pushed to hosts via `bullet sentinel:install`.

## Installing

```sh
bullet -H 192.168.0.3 sentinel:install
```

This pushes the embedded `bullet-sentinel` binary to
`/usr/local/bin/bullet-sentinel`, writes
`/etc/systemd/system/bullet-sentinel.service`, runs `systemctl daemon-reload`,
and enables and starts the unit.

The sentinel binary is selected to match the host's architecture (amd64 or
arm64). The `bullet` binary you run must have been built with the embed step
already done — see [Building from source](#building-from-source) below.

The configuration file at `/etc/bullet/sentinel.yaml` is **not** pushed by
`sentinel:install`. It's expected to be managed out-of-band (e.g. by your
provisioning tooling). If it's absent, the sentinel runs with built-in
defaults.

## Configuration

The sentinel reads its configuration in this order, with later sources
overriding earlier ones:

1. Built-in defaults
2. `/etc/bullet/sentinel.yaml` (or the path passed via `--config`)
3. `BULLET_SENTINEL_*` environment variables
4. CLI flags

### File

```yaml
# HTTP bind address for the /metrics endpoint. Default is localhost-only;
# set to 0.0.0.0:9479 to allow remote scraping (put auth in front of it).
addr: 127.0.0.1:9479

# Minimum interval between docker/systemctl introspections. Multiple scrapes
# within this window are served from a cached snapshot.
scrape_interval: 15s

# Path to the docker binary on the host.
docker_path: /usr/bin/docker
```

A working example ships with `bullet` at [`embed/sentinel.example.yaml`](../embed/sentinel.example.yaml).

### Environment variables

| Variable                          | Equivalent flag       |
|-----------------------------------|-----------------------|
| `BULLET_SENTINEL_ADDR`            | `--addr`              |
| `BULLET_SENTINEL_SCRAPE_INTERVAL` | `--scrape-interval`   |
| `BULLET_SENTINEL_DOCKER_PATH`     | `--docker-path`       |

### Flags

| Flag                | Default                      | Purpose                              |
|---------------------|------------------------------|--------------------------------------|
| `--config`          | `/etc/bullet/sentinel.yaml`  | Path to YAML config file             |
| `--addr`            | (from config)                | HTTP bind address                    |
| `--scrape-interval` | (from config)                | Min interval between introspections  |
| `--docker-path`     | (from config)                | Path to the docker binary            |

## Endpoints

| Path       | Purpose                                                |
|------------|--------------------------------------------------------|
| `/metrics` | Prometheus exposition format with the metrics below.   |
| `/healthz` | Liveness probe. Always returns `200 OK` with `ok\n`.   |

## Exported metrics

All metrics are gauges. Values are derived from a snapshot refreshed at most
once per `scrape_interval`.

### Sentinel self-state

| Metric                          | Labels    | Meaning                                      |
|---------------------------------|-----------|----------------------------------------------|
| `bullet_sentinel_build_info`    | `version` | Always `1`. Label carries the build version. |
| `bullet_sentinel_docker_up`     |           | `1` if the last `docker ps` call succeeded.  |
| `bullet_sentinel_systemd_up`    |           | `1` if the last `systemctl` call succeeded.  |

### Containers

The sentinel matches container names of the form `<app>_<program>_<instance>`
— Bullet's standard container naming.

| Metric                              | Labels                       | Meaning                                                                                |
|-------------------------------------|------------------------------|----------------------------------------------------------------------------------------|
| `bullet_container_up`               | `app`, `program`, `instance` | `1` if the container's status starts with `Up`.                                        |
| `bullet_container_healthy`          | `app`, `program`, `instance` | `1` healthy, `0` unhealthy or starting. Absent for containers without a healthcheck.   |
| `bullet_program_instances_running`  | `app`, `program`             | Count of running instances per program on this host.                                   |

### Cron jobs

The sentinel matches systemd timer units of the form `bullet_<app>_<job>.timer`
— the format `bullet cron:enable` writes.

| Metric                                          | Labels         | Meaning                                                       |
|-------------------------------------------------|----------------|---------------------------------------------------------------|
| `bullet_cron_timer_active`                      | `app`, `job`   | `1` if the timer's `ActiveState` is `active`.                 |
| `bullet_cron_last_trigger_timestamp_seconds`    | `app`, `job`   | Unix timestamp of the most recent trigger. `0` if never run.  |
| `bullet_cron_last_result`                       | `app`, `job`   | `1` if the last run reported a non-success `Result`.          |

## Building from source

The sentinel is cross-compiled and embedded into `bullet` via `//go:embed`.
Before building `bullet` for release, run:

```sh
task sentinel:embed
task build
```

`sentinel:embed` produces `embed/bullet-sentinel-linux-amd64` and
`embed/bullet-sentinel-linux-arm64`. Without this step, `bullet sentinel:install`
returns an error pointing at the missing embed.

## Limitations

- Containers managed by other tools that share the `<a>_<b>_<n>` naming pattern
  (e.g. docker-compose) will be picked up as Bullet containers.
- App identifiers and program/job keys containing underscores can confuse the
  parser, since the sentinel splits on the last underscore.
- The `/metrics` endpoint has no built-in authentication. If you bind to a
  non-localhost address, put a reverse proxy or firewall in front of it.
