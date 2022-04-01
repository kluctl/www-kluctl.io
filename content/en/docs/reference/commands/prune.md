---
title: "prune"
linkTitle: "prune"
weight: 10
description: >
    prune command
---

## Command
<!-- BEGIN SECTION "prune" "Usage" false -->
Usage: kluctl prune

<!-- END SECTION -->

## Arguments
The following sets of arguments are available:
1. [project arguments]({{< ref "./common-arguments#project-arguments" >}})
1. [image arguments]({{< ref "./common-arguments#image-arguments" >}})
1. [inclusion/exclusion arguments]({{< ref "./common-arguments#inclusionexclusion-arguments" >}})

In addition, the following arguments are available:
<!-- BEGIN SECTION "prune" "Misc arguments" true -->
```
Misc arguments:
  Command specific arguments.

  -y, --yes                                Suppresses 'Are you sure?' questions and proceeds as if you would answer
                                           'yes'.
      --dry-run                            Performs all kubernetes API calls in dry-run mode.
  -o, --output-format=OUTPUT-FORMAT,...    Specify output format and target file, in the format 'format=path'. Format
                                           can either be 'text' or 'yaml'. Can be specified multiple times. The actual
                                           format for yaml is currently not documented and subject to change.

```
<!-- END SECTION -->

They have the same meaning as described in [deploy](#deploy).
