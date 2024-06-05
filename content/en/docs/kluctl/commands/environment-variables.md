---
description: Controlling Kluctl via environment variables
github_branch: main
github_repo: https://github.com/kluctl/kluctl
lastmod: "2023-10-17T00:30:26+02:00"
linkTitle: Environment Variables
path_base_for_github_subdir:
    from: .*
    to: docs/kluctl/commands/environment-variables.md
title: Environment Variables
weight: 2
---



In addition to arguments, Kluctl can be controlled via a set of environment variables.

## Environment variables as arguments
All options/arguments accepted by kluctl can also be specified via environment variables. The name of the environment
variables always start with `KLUCTL_` and end with the option/argument in uppercase and dashes replaced with
underscores. As an example, `--dry-run` can also be specified with the environment variable
`KLUCTL_DRY_RUN=true`.

If an argument needs to be specified multiple times through environment variables, indexed can be appended to the
names of the environment variables, e.g. `KLUCTL_ARG_0=name1=value1` and `KLUCTL_ARG_1=name2=value2`.

## Additional environment variables
A few additional environment variables are supported which do not belong to an option/argument. These are:

1. `KLUCTL_REGISTRY_<idx>_HOST`, `KLUCTL_REGISTRY_<idx>_USERNAME`, and so on. See [OCI authentication](../deployments/oci.md#authentication) for details.
2. `KLUCTL_HELM_<idx>_HOST`, `KLUCTL_HELM_<idx>_USERNAME`, and so on. See [Helm private repositories](../deployments/helm.md#private-repositories) for details.
3. `KLUCTL_GIT_<idx>_HOST`, `KLUCTL_GIT_<idx>_USERNAME`, and so on.
4. `KLUCTL_SSH_DISABLE_STRICT_HOST_KEY_CHECKING`. Disable ssh host key checking when accessing git repositories.
