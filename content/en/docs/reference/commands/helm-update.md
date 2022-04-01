---
title: "helm-update"
linkTitle: "helm-update"
weight: 10
description: >
    helm-update command
---

## Command
<!-- BEGIN SECTION "helm-update" "Usage" false -->
Usage: kluctl helm-update

<!-- END SECTION -->

## Arguments
The following sets of arguments are available:
1. [project arguments]({{< ref "./common-arguments#project-arguments" >}}) (except `-a`)

In addition, the following arguments are available:
<!-- BEGIN SECTION "helm-update" "Misc arguments" true -->
```
Misc arguments:
  Command specific arguments.

  --upgrade    Write new versions into helm-chart.yml and perform helm-pull afterwards
  --commit     Create a git commit for every updated chart

```
<!-- END SECTION -->