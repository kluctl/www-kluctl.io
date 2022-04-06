#!/usr/bin/env bash

set -e

case "$(uname -s)" in
      Linux*)     arch=linux-amd64;;
      Darwin*)    arch=darwin-amd64;;
      MINGW*)     arch=windows-amd64;exe=.exe;;
      MSYS*)      arch=windows-amd64;exe=.exe;;
      *)          echo "unknown os"; exit 1;
esac

echo "Determining latest_version"
latest_version=$(curl -s https://api.github.com/repos/kluctl/kluctl/releases/latest | jq '.tag_name' -r)
echo "latest_version=$latest_version"

mkdir -p /tmp/kluctl-for-docs
curl -L -o /tmp/kluctl-for-docs/kluctl$exe "https://github.com/kluctl/kluctl/releases/download/$latest_version/kluctl-$arch$exe"
chmod +x /tmp/kluctl-for-docs/kluctl$exe
ls -lah /tmp/kluctl-for-docs/

export PATH=/tmp/kluctl-for-docs:$PATH

go mod vendor
go run ./replace-commands-help --docs-dir ./content/en/docs/reference/commands

hugo
