---
description: validate command
github_repo: https://github.com/kluctl/kluctl
lastmod: "2023-10-17T00:30:26+02:00"
linkTitle: validate
path_base_for_github_subdir:
    from: .*
    to: main/docs/kluctl/commands/validate.md
title: validate
weight: 10
---



## Command
<!-- BEGIN SECTION "validate" "Usage" false -->
Usage: kluctl validate [flags]

Validates the already deployed deployment
This means that all objects are retrieved from the cluster and checked for readiness.

TODO: This needs to be better documented!

<!-- END SECTION -->

## Arguments
The following sets of arguments are available:
1. [project arguments](./common-arguments.md#project-arguments)
1. [image arguments](./common-arguments.md#image-arguments)
1. [helm arguments](./common-arguments.md#helm-arguments)
1. [registry arguments](./common-arguments.md#registry-arguments)

In addition, the following arguments are available:
<!-- BEGIN SECTION "validate" "Misc arguments" true -->
```
Misc arguments:
  Command specific arguments.

  -o, --output stringArray         Specify output target file. Can be specified multiple times
      --render-output-dir string   Specifies the target directory to render the project into. If omitted, a
                                   temporary directory is used.
      --sleep duration             Sleep duration between validation attempts (default 5s)
      --wait duration              Wait for the given amount of time until the deployment validates
      --warnings-as-errors         Consider warnings as failures

```
<!-- END SECTION -->
