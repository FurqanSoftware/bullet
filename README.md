# Bullet

![](assets/bullet_128.png)

Bullet is a fast and flexible application deployment tool built by [Furqan Software](https://furqansoftware.com/). It automates the full lifecycle of deploying containerized applications to remote servers over SSH, without the overhead of a full orchestration platform.

## Features

- **Server setup**: Installs Docker and prepares remote servers for deployments.
- **Tarball-based deploys**: Packages application code into tarballs, uploads them over SCP, and builds Docker images on the server.
- **Program management**: Define multiple programs (e.g. web, worker) per application, each running in its own container.
- **Scaling**: Scale programs up or down with expression-based rules that can factor in host tags and hardware.
- **Cron jobs**: Manage scheduled tasks backed by systemd timers, with optional healthcheck pings.
- **Zero-downtime reloads**: Reload containers via signal, command, or restart, with support for pre-reload hooks.
- **Environment management**: Push environment files to servers.
- **Log tailing**: Tail container logs directly from the CLI.
- **Port forwarding**: Forward remote ports to your local machine over SSH.
- **Host access**: Open an interactive shell, view disk usage, or run `top` on remote servers.
- **Multi-node support**: Target multiple hosts per deployment with interactive node selection.
- **Release pruning**: Clean up old releases to free disk space.

## Getting Started

### Install

From source:

```sh
go install github.com/FurqanSoftware/bullet@latest
```

Or download a prebuilt binary from the [Releases](https://github.com/FurqanSoftware/bullet/releases) page.

### Define a Bulletspec

Create a `Bulletspec` file in your project root. This defines your application and its programs:

```yaml
application:
  name: Hello World
  identifier: hello

  programs:
    web:
      name: Hello World Web Server
      command: node index.js
      container:
        image: node:8.1-alpine
      ports:
        - 80:5000
```

### Configure Hosts

Pass hosts directly via flags:

```sh
bullet -H 192.168.0.3 <command>
```

Or create a `Bulletcfg.<name>` file:

```yaml
hosts: 192.168.0.3,192.168.0.4
port: 22
```

Then use it with:

```sh
bullet -c <name> <command>
```

You can also set hosts via environment variables:

```sh
export BULLET_HOSTS=192.168.0.3
```

### Set Up a Server

```sh
bullet -H 192.168.0.3 setup
```

This installs Docker and prepares the server for deployments.

### Deploy

Package your application as a tarball and deploy:

```sh
tar czf app.tar.gz <your files>
bullet -H 192.168.0.3 deploy app.tar.gz
```

### Scale Programs

```sh
bullet -H 192.168.0.3 scale web=2
```

### Other Commands

```sh
bullet status              # Show container status
bullet restart             # Restart containers
bullet run <program>       # Run a one-off container
bullet log <program>       # Tail container logs
bullet cron:enable <job>   # Enable a cron job
bullet cron:disable <job>  # Disable a cron job
bullet cron:status         # Show cron job status
bullet environ:push <file> # Push environment file to server
bullet forward <port>      # Forward a remote port locally
bullet prune               # Remove old releases
bullet host:shell          # SSH into a server
bullet host:df             # Show disk usage
bullet host:top            # Show running processes
```

## Shell Completion

Bullet supports autocompletion for bash, zsh, fish, and powershell. It completes command names, flag values (e.g. `-c` from `Bulletcfg.*` files), and arguments (e.g. program keys and cron job keys from `Bulletspec`).

To enable it, add the following to your shell configuration:

**Bash** (`~/.bashrc`):
```sh
eval "$(bullet completion bash)"
```

**Zsh** (`~/.zshrc`):
```sh
eval "$(bullet completion zsh)"
```

**Fish** (`~/.config/fish/config.fish`):
```sh
bullet completion fish | source
```

**PowerShell**:
```powershell
bullet completion powershell | Out-String | Invoke-Expression
```

## But, Kubernetes?

https://k8s.af/

## Acknowledgements

- [Nikita Golubev](http://www.flaticon.com/authors/nikita-golubev) - For the bullet icon
