#!/usr/bin/env bash

set -e

case "$(uname -s)" in
      Linux*)     os=linux;arch=amd64;;
      Darwin*)    os=darwin;arch=amd64;;
      MINGW*)     os=windows;arch=amd64;exe=.exe;;
      MSYS*)      os=windows;arch=amd64;exe=.exe;;
      *)          echo "unknown os"; exit 1;
esac

echo "Downloading kluctl install script"
curl -s -# -Lf https://raw.githubusercontent.com/kluctl/kluctl/main/install/kluctl.sh -o static/install.sh

echo "Determining version"
version=$(cat config.toml | grep 'fullversion =' | sed 's/fullversion = "\(.*\)"/\1/')
echo "version=$version"

mkdir -p /tmp/kluctl-for-docs
curl -L -o /tmp/kluctl-for-docs/kluctl.tar.gz "https://github.com/kluctl/kluctl/releases/download/$version/kluctl_${version}_${os}_${arch}.tar.gz"
(cd /tmp/kluctl-for-docs && tar xzf kluctl.tar.gz)
ls -lah /tmp/kluctl-for-docs/

export PATH=/tmp/kluctl-for-docs:$PATH

go mod vendor
go run ./replace-commands-help --docs-dir ./content/en/reference/commands

if [ "$BASE_URL" != "" ]; then
  BASE_URL_ARG="-b $BASE_URL"
elif [ "$DEPLOY_PRIME_URL" != "" ]; then
  BASE_URL_ARG="-b $DEPLOY_PRIME_URL/"
fi

if [ "$CONTEXT" != "" ]; then
  export HUGO_ENV=$CONTEXT
fi

hugo $BASE_URL_ARG
