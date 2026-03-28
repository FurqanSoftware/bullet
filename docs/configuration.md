# Configuration

Bullet can be configured through configuration files, environment variables, and command-line flags. These are applied in order, with later sources overriding earlier ones:

1. Configuration file (`Bulletcfg.<name>`)
2. Environment variables (`BULLET_*`)
3. Command-line flags

## Configuration File

Configuration files are YAML files named `Bulletcfg.<name>`. Select one with the `-c` flag:

```sh
bullet -c production deploy app.tar.gz
```

This loads `Bulletcfg.production` from the current directory.

```yaml
# Bulletcfg.production
hosts: 192.168.0.3,192.168.0.4
port: 22
identity: ~/.ssh/deploy_key

vars:
  PORT: "8080"
  ENV: production
```

| Field        | Type     | Default | Description                                |
|--------------|----------|---------|--------------------------------------------|
| `hosts`      | string   |         | Target hosts. See [Hosts](#hosts).         |
| `port`       | int      | `22`    | SSH port.                                  |
| `identity`   | string   |         | Path to SSH private key.                   |
| `vars`       | map      |         | Variables for Bulletspec expansion.        |
| `sshretries` | int      | `0`     | Number of SSH connection retries.          |
| `sshtimeout` | duration | `30s`   | SSH connection timeout.                    |

## Environment Variables

All configuration fields can be set via environment variables prefixed with `BULLET_`:

| Variable            | Field        |
|---------------------|--------------|
| `BULLET_HOSTS`      | `hosts`      |
| `BULLET_PORT`       | `port`       |
| `BULLET_IDENTITY`   | `identity`   |
| `BULLET_SSH_RETRIES`| `sshretries` |
| `BULLET_SSH_TIMEOUT` | `sshtimeout` |

Example:

```sh
export BULLET_HOSTS=192.168.0.3
export BULLET_SSH_RETRIES=3
bullet deploy app.tar.gz
```

## Command-Line Flags

| Flag               | Short | Description                              |
|--------------------|-------|------------------------------------------|
| `--config <name>`  | `-c`  | Configuration name (loads `Bulletcfg.<name>`). |
| `--hosts <hosts>`  | `-H`  | Target hosts.                            |
| `--port <port>`    | `-p`  | SSH port (default: 22).                  |
| `--identity <path>`| `-i`  | Path to SSH private key.                 |

## Hosts

The hosts field supports two formats.

### Comma-Separated Hosts

A simple list of hostnames or IPs:

```yaml
hosts: 192.168.0.3,192.168.0.4,192.168.0.5
```

Or via flag:

```sh
bullet -H 192.168.0.3,192.168.0.4 deploy app.tar.gz
```

### Node Manifest

For more complex setups, reference a YAML manifest file with `@`:

```yaml
hosts: "@nodes.yaml"
```

The manifest is a YAML array of nodes:

```yaml
# nodes.yaml
- name: web-1
  host: 192.168.0.3
  port: 22
  tags:
    - production
    - web
  hw:
    cores: 8
    memory: 16384

- name: web-2
  host: 192.168.0.4
  port: 22
  tags:
    - production
    - web
  hw:
    cores: 8
    memory: 16384

- name: worker-1
  host: 192.168.0.5
  port: 22
  tags:
    - production
    - worker
  hw:
    cores: 4
    memory: 8192
```

#### Node Fields

| Field    | Type     | Description                                          |
|----------|----------|------------------------------------------------------|
| `name`   | string   | Display name for the node.                           |
| `host`   | string   | Hostname or IP address.                              |
| `port`   | int      | SSH port.                                            |
| `tags`   | []string | Tags for filtering and scaling expressions.          |
| `hw`     | object   | Hardware info, used in scaling expressions.           |

#### Hardware Fields

| Field    | Type | Description          |
|----------|------|----------------------|
| `cores`  | int  | Number of CPU cores. |
| `memory` | int  | Memory in MB.        |

#### Filtering by Tags

Append tags after a colon to filter nodes from the manifest:

```yaml
# All nodes
hosts: "@nodes.yaml"

# Nodes tagged "production"
hosts: "@nodes.yaml:production"

# Nodes tagged both "production" AND "web"
hosts: "@nodes.yaml:production+web"

# Nodes tagged "production" OR "staging"
hosts: "@nodes.yaml:production,staging"
```

## Variables

Variables defined in the configuration file can be used in Bulletspec using Go template syntax:

```yaml
# Bulletcfg.production
vars:
  PORT: "8080"
  WORKERS: "4"
```

```yaml
# Bulletspec
application:
  programs:
    web:
      command: ./server --port {{.Vars.PORT}} --workers {{.Vars.WORKERS}}
```

## Node Selection

When multiple nodes are targeted, Bullet prompts for selection before running a command.

For commands that operate on a single node (`run`, `log`, `forward`, `host:shell`, `host:df`, `host:top`):

```
1. web-1
2. web-2
3. worker-1
? [1]
```

Enter a number to select a node. Press Enter to accept the default.

For commands that operate on multiple nodes (`deploy`, `setup`, `restart`, `scale`, `cron:*`, `environ:push`, `prune`, `status`):

```
1. web-1
2. web-2
3. worker-1
? [1-3]
```

Supported selection formats:

| Input   | Selects            |
|---------|--------------------|
| `2`     | Node 2             |
| `1-3`   | Nodes 1 through 3  |
| `1,3`   | Nodes 1 and 3      |
| `1-2,4` | Nodes 1, 2, and 4  |
