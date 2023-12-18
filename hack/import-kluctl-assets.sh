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

    go run ./sync-docs --repo kluctl/kluctl --subdir docs/kluctl --dest content/en/docs/kluctl --min-version=v2.21.0
    go run ./sync-docs --repo kluctl/kluctl --subdir docs/gitops --dest content/en/docs/gitops --min-version=v2.21.0
    go run ./sync-docs --repo kluctl/kluctl --subdir docs/webui --dest content/en/docs/webui --min-version=v2.21.0
    go run ./sync-docs --repo kluctl/template-controller --subdir docs --dest content/en/docs/template-controller --with-root-readme --min-version=v0.7.1
}

{
  # provide Kluctl install script
  curl -f -s -# -Lf https://raw.githubusercontent.com/kluctl/kluctl/main/install/kluctl.sh -o static/install.sh
}
