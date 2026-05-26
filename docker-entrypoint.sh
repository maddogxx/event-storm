#!/bin/sh
set -e

DATA_DIR="${NUXT_DATA_DIR:-/app/data}"
APP_UID="${APP_UID:-1001}"
APP_GID="${APP_GID:-1001}"

mkdir -p "$DATA_DIR"

if [ "$(id -u)" = "0" ]; then
  chown -R "$APP_UID:$APP_GID" "$DATA_DIR" 2>/dev/null || true
  exec gosu "$APP_UID:$APP_GID" "$@"
else
  exec "$@"
fi
