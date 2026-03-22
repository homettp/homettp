![Homettp](https://user-images.githubusercontent.com/1419087/194723881-3ae9add0-f5a0-404c-aab0-2025dd5572a8.png)

## Badges

[![Build Status](https://github.com/homettp/homettp/workflows/tests/badge.svg)](https://github.com/homettp/homettp/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](LICENSE.md)

## Features

- **HTTP-based command runner** — token-based access with payload parameters and configurable timeouts
- **Background worker processing** — async command queue with real-time status tracking
- **Call history** — collapsible command output with auto-refresh and configurable limits
- **Command management** — re-generable tokens, custom icons, and shareable URLs
- **User authentication** — remember me option with secure Redis-backed sessions
- **User management** — CLI and web interface with Gravatar integration
- **Modern UI** — Vue 3, Tailwind CSS, dark/light/system theme, responsive design
- **Cross-platform** — Linux, macOS, and Windows support
- **Docker support** — multi-architecture image (amd64, arm64)
- **Redis-backed storage** — configurable key prefix for multi-instance deployments

## Getting Started

Follow the steps below to install and configure Homettp.

### Prerequisites

- Redis: `Version >= 5.0` for data storage
- `/bin/sh` on Linux/MacOS for command execution
- `cmd.exe` on Windows for command execution

### Run with Docker

Image can be found at the package page on [GitHub](https://github.com/homettp/homettp/pkgs/container/homettp).

```bash
docker run --rm \
-e APP_URL=http://127.0.0.1:4000 \
-e APP_KEY=$(openssl rand -hex 16) \
-e REDIS_URL=redis://192.168.0.200:6379/0 \
-p 4000:4000 \
ghcr.io/homettp/homettp
```

### Install from Binary

Download the latest release for your platform from the [GitHub Releases](https://github.com/homettp/homettp/releases) page.

---

### Install from Source

#### Prerequisites

- Go: `Version >= 1.26`
- Node.js: `Version >= 22.0`
- Yarn or NPM

#### Steps

1. Clone the repository:

```bash
git clone git@github.com:homettp/homettp.git
```

2. Install UI dependencies and build:

```bash
cd homettp
yarn install
yarn build
```

3. Build the binary:

```bash
go build
```

4. Copy and edit the configuration:

```bash
cp .env.example .env
```

## Configuration

All configuration is done through environment variables in the `.env` file.

### General

#### Application Key (encryption key)

- For example generate with `openssl rand -hex 16` command

```
APP_KEY=
```

#### Application Address

```
APP_ADDR=:4000
```

#### Application URL

```
APP_URL=http://127.0.0.1:4000
```

---

### Redis

#### Redis URL

```
REDIS_URL=redis://127.0.0.1:6379/0
```

#### Redis Key Prefix

```
REDIS_KEY_PREFIX=homettp:
```

---

### Command

#### Command Timeout (in seconds)

```
COMMAND_TIMEOUT=60
```

#### Command Worker Count

```
COMMAND_WORKER_COUNT=2
```

#### Command History Limit

```
COMMAND_HISTORY_LIMIT=100
```

## Usage

### 1. Make a user

```bash
./homettp make user <username> <email> <password>
```

### 2. Run the app

```bash
./homettp web serve
```

## Reporting Issues

If you are facing a problem with this package or found any bug, please open an issue on [GitHub](https://github.com/homettp/homettp/issues).

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
