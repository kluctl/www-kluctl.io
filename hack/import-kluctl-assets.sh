#!/bin/bash

set -e

fatal() {
    echo '[ERROR] ' "$@" >&2
    exit 1
}

# Set os, fatal if operating system not supported
setup_verify_os() {
    if [ -z "${OS}" ]; then
        OS=$(uname)
    fi
    case ${OS} in
        Darwin)
            OS=darwin
            ;;
        Linux)
            OS=linux
            ;;
        *)
            fatal "Unsupported operating system ${OS}"
    esac
}

# Set arch, fatal if architecture not supported
setup_verify_arch() {
    if [ -z "${ARCH}" ]; then
        ARCH=$(uname -m)
    fi
    case ${ARCH} in
        arm|armv6l|armv7l)
            ARCH=arm
            ;;
        arm64|aarch64|armv8l)
            ARCH=arm64
            ;;
        amd64)
            ARCH=amd64
            ;;
        x86_64)
            ARCH=amd64
            ;;
        *)
            fatal "Unsupported architecture ${ARCH}"
    esac
}

{
    setup_verify_os
    setup_verify_arch

    go run ./sync-docs --repo kluctl/kluctl --subdir docs/kluctl --dest content/en/docs/kluctl
    go run ./sync-docs --repo kluctl/kluctl --subdir docs/gitops --dest content/en/docs/gitops
    go run ./sync-docs --repo kluctl/kluctl --subdir docs/webui --dest content/en/docs/webui
    go run ./sync-docs --repo kluctl/template-controller --subdir docs --dest content/en/docs/template-controller --with-root-readme
}

{
  # provide Kluctl install script
  if [ ! -d static ]; then
    mkdir static
  fi
  curl -f -s -# -Lf https://raw.githubusercontent.com/kluctl/kluctl/main/install/kluctl.sh -o static/install.sh
}
