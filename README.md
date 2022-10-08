![Homettp](https://user-images.githubusercontent.com/1419087/194723881-3ae9add0-f5a0-404c-aab0-2025dd5572a8.png)

## Badges

[![Build Status](https://github.com/homettp/homettp/workflows/tests/badge.svg)](https://github.com/homettp/homettp/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](LICENSE.md)

## Getting Started

Before you start, you need to install the prerequisites.

### Prerequisites

- Redis: `Version >= 5.0` for data storage
- GO: `Version >= 1.19` for building
- Node.js: `Version >= 14.0` for building
- Yarn or NPM: for building

### Install with Docker

Image can be found at package page on [GitHub](https://github.com/homettp/homettp/pkgs/container/homettp).

### Install from binary

Downloads can be found at releases page on [GitHub](https://github.com/homettp/homettp/releases).

### Install from source

1. Clone the repository:

```
git clone git@github.com:homettp/homettp.git
```

2. Open the folder:

```
cd homettp
```

3. Install the UI dependencies

```
yarn install
```

4. Build the UI

```
yarn prod
```

5. Build the Homettp:

```
go build
```

6. Copy the example configuration:

```
cp .env.example .env
```

## Configuration

The configruation is stored in the `.env` file.

#### Application Address:

```
APP_ADDR=:4000
```

#### Application URL:

```
APP_URL=http://127.0.0.1:4000
```

---

#### Redis URL:

```
REDIS_URL=redis://127.0.0.1:6379/0
```

#### Redis Key Prefix:

```
REDIS_KEY_PREFIX=homettp:
```

---

#### Command Timeout (in seconds):

```
COMMAND_TIMEOUT=60
```

#### Command Worker Count:

```
COMMAND_WORKER_COUNT=2
```

#### Command Histroy Limit

```
COMMAND_HISTORY_LIMIT=100
```

## Usage

Run the app using the following command:

```
./homettp web serve
```

## Reporting Issues

If you are facing a problem with this package or found any bug, please open an issue on [GitHub](https://github.com/homettp/homettp/issues).

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
