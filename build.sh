#!/usr/bin/env bash

set -e

case "$(uname -s)" in
      Linux*)     arch=linux-amd64;;
      Darwin*)    arch=darwin-amd64;;
      MINGW*)     arch=windows-amd64;exe=.exe;;
      MSYS*)      arch=windows-amd64;exe=.exe;;
      *)          echo "unknown os"; exit 1;
esac

echo "Determining version"
version=$(cat config.toml | grep 'fullversion =' | sed 's/fullversion = "\(.*\)"/\1/')
echo "version=$version"

mkdir -p /tmp/kluctl-for-docs
curl -L -o /tmp/kluctl-for-docs/kluctl$exe "https://github.com/kluctl/kluctl/releases/download/$version/kluctl-$arch$exe"
chmod +x /tmp/kluctl-for-docs/kluctl$exe
ls -lah /tmp/kluctl-for-docs/

export PATH=/tmp/kluctl-for-docs:$PATH

go mod vendor
go run ./replace-commands-help --docs-dir ./content/en/docs/reference/commands

export HUGO_ENV=$CONTEXT

hugo
