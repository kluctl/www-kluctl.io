---
description: helm-pull command
github_repo: https://github.com/kluctl/kluctl
lastmod: "2023-10-17T00:30:26+02:00"
linkTitle: helm-pull
path_base_for_github_subdir:
    from: .*
    to: main/docs/kluctl/commands/helm-pull.md
title: helm-pull
weight: 10
---



## Command
<!-- BEGIN SECTION "helm-pull" "Usage" false -->
Usage: kluctl helm-pull [flags]

Recursively searches for 'helm-chart.yaml' files and pre-pulls the specified Helm charts
Kluctl requires Helm Charts to be pre-pulled by default, which is handled by this command. It will collect
all required Charts and versions and pre-pull them into .helm-charts. To disable pre-pulling for individual charts,
set 'skipPrePull: true' in helm-chart.yaml.

<!-- END SECTION -->

See [helm-integration](../deployments/helm.md) for more details.

## Arguments
The following sets of arguments are available:
1. [project arguments](./common-arguments.md#project-arguments) (except `-a`)
1. [helm arguments](./common-arguments.md#helm-arguments)
1. [registry arguments](./common-arguments.md#registry-arguments)
