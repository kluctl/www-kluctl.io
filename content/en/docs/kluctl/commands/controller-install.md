---
description: controller command
github_branch: main
github_repo: https://github.com/kluctl/kluctl
lastmod: "2023-08-26T09:38:51+02:00"
linkTitle: controller install
path_base_for_github_subdir:
    from: .*
    to: docs/kluctl/commands/controller-install.md
title: controller install
weight: 10
---



## Command
<!-- BEGIN SECTION "controller install" "Usage" false -->
Usage: kluctl controller install [flags]

Install the Kluctl controller
This command will install the kluctl-controller to the current Kubernetes clusters.

<!-- END SECTION -->

## Arguments
The following sets of arguments are available:
1. [command results arguments](./common-arguments.md#command-results-arguments)

In addition, the following arguments are available:
<!-- BEGIN SECTION "controller install" "Misc arguments" true -->
```
Misc arguments:
  Command specific arguments.

      --context string          Override the context to use.
      --dry-run                 Performs all kubernetes API calls in dry-run mode.
      --kluctl-version string   Specify the controller version to install.
  -y, --yes                     Suppresses 'Are you sure?' questions and proceeds as if you would answer 'yes'.

```
<!-- END SECTION -->
