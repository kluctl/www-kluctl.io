---
title: "helm-update"
linkTitle: "helm-update"
weight: 10
description: >
    helm-update command
---

## Command
<!-- BEGIN SECTION "helm-update" "Usage" false -->
Usage: kluctl helm-update [flags]

Recursively searches for 'helm-chart.yml'' files and checks for new available versions
Optionally performs the actual upgrade and/or add a commit to version control.

<!-- END SECTION -->

## Arguments
The following sets of arguments are available:
1. [project arguments]({{< ref "./common-arguments#project-arguments" >}}) (except `-a`)

In addition, the following arguments are available:
<!-- BEGIN SECTION "helm-update" "Misc arguments" true -->
```
Misc arguments:
  Command specific arguments.

      --commit    Create a git commit for every updated chart
      --upgrade   Write new versions into helm-chart.yml and perform helm-pull afterwards

```
<!-- END SECTION -->