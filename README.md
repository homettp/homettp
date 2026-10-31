![Homettp](https://user-images.githubusercontent.com/1419087/194723881-3ae9add0-f5a0-404c-aab0-2025dd5572a8.png)

## Badges

[![Build Status](https://github.com/homettp/homettp/workflows/tests/badge.svg)](https://github.com/homettp/homettp/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](LICENSE.md)

## Features

- **HTTP-based command runner** — token access, payload parameters, multiline support, configurable timeouts
- **Background worker processing** — async queue with real-time status tracking
- **Call history** — command output with auto-refresh
- **Command management** — regenerable tokens, custom icons, shareable URLs
- **User authentication** — remember me, Redis-backed sessions
- **User management** — CLI and web interface, Gravatar integration
- **Modern UI** — Vue 3, Tailwind CSS, dark/light/system theme, responsive
- **Cross-platform** — Linux, macOS, Windows
- **Docker support** — multi-architecture (amd64, arm64)
- **Redis-backed storage** — configurable key prefix for multi-instance use

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

- Go: `Version >= 1.27`
- Node.js: `Version >= 24.0`
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

- Must be exactly 32 characters, generate with `openssl rand -hex 16` command

```
APP_KEY=
```

#### Application Debug Mode

- Serves the assets from the `./static` directory, use it only for local development

```
APP_DEBUG=false
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

### 3. Call a command

Make a `GET` or `POST` request to the command URL shown on the command page. The optional `payload` query parameter is available in the command as `%p`:

```bash
curl "http://127.0.0.1:4000/call?commandID=1&token=<token>&payload=hello"
```

The payload is never inserted into the command text. It is passed in the `HOMETTP_PAYLOAD` environment variable, and `%p` expands to a quoted reference to it (`${HOMETTP_PAYLOAD:+"$HOMETTP_PAYLOAD"}` with `/bin/sh`), so it is a single argument, or no argument when the payload is empty. Don't wrap `%p` in quotes yourself, and don't start the command with it.

On Windows `%p` expands to `"%HOMETTP_PAYLOAD%"` (or nothing when the payload is empty), but `cmd.exe` expands variables before parsing the command line, so it is not shell-safe there. Payloads containing `" % ! ^ & | < > ( )` or line breaks, or ending in a backslash, are rejected with `422 Unprocessable Entity`.

## Reporting Issues

If you are facing a problem with this package or found any bug, please open an issue on [GitHub](https://github.com/homettp/homettp/issues).

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
