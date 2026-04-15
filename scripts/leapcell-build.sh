#!/bin/sh
set -eu

echo "==> Go version"
go version

echo "==> Building MoodMap API"
go build -o app ./cmd/api
