#!/usr/bin/env bash

set -e
set -o pipefail

case "$(uname -s)" in
      Linux*)     os=linux;arch=amd64;;
      Darwin*)    os=darwin;arch=amd64;;
      MINGW*)     os=windows;arch=amd64;exe=.exe;;
      MSYS*)      os=windows;arch=amd64;exe=.exe;;
      *)          echo "unknown os"; exit 1;
esac

./hack/import-kluctl-assets.sh

if [ "$BASE_URL" != "" ]; then
  BASE_URL_ARG="-b $BASE_URL"
elif [ "$DEPLOY_PRIME_URL" != "" ]; then
  BASE_URL_ARG="-b $DEPLOY_PRIME_URL/"
fi

if [ "$CONTEXT" != "" ]; then
  export HUGO_ENV=$CONTEXT
fi

hugo $BASE_URL_ARG
