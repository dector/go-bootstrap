## Setup

- Replace all `mymyapp` with your mod name.

## Features

- Customizable bind host via `HOST` and port via `PORT` env variables.
- Automatic hot reload via `gust` on port `3000` (use `ror dev`).
- Centralized meta-info in `meta/meta.yml` file.

## Defaults

- Default host (when `HOST` is unset or empty): `localhost` (loopback only). Set `HOST=0.0.0.0` (or `HOST="*"`) to listen on all IPv4 interfaces.
- Default port: `8080`.
