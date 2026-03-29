#!/bin/sh
set -e

IMAGE_DIR="${IMAGE_DIR:-/var/www/html/images}"

if [ -d "$IMAGE_DIR" ]; then
    chown -R 1001:1001 "$IMAGE_DIR"
else
    mkdir -p "$IMAGE_DIR"
    chown -R 1001:1001 "$IMAGE_DIR"
fi

exec su-exec appuser "$@"
