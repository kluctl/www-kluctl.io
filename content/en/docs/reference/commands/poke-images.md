---
title: "poke-images"
linkTitle: "poke-images"
weight: 10
description: >
    poke-images command
---

## Command
<!-- BEGIN SECTION "poke-images" "Usage" false -->
Usage: kluctl poke-images

<!-- END SECTION -->

## Arguments
The following sets of arguments are available:
1. [project arguments]({{< ref "./common-arguments#project-arguments" >}})
1. [image arguments]({{< ref "./common-arguments#image-arguments" >}})
1. [inclusion/exclusion arguments]({{< ref "./common-arguments#inclusionexclusion-arguments" >}})

In addition, the following arguments are available:
<!-- BEGIN SECTION "poke-images" "Misc arguments" true -->
```
Misc arguments:
  Command specific arguments.

  -y, --yes                         Suppresses 'Are you sure?' questions and proceeds as if you would answer 'yes'.
      --dry-run                     Performs all kubernetes API calls in dry-run mode.
  -o, --output=OUTPUT,...           Specify output target file. Can be specified multiple times
      --render-output-dir=STRING    Specifies the target directory to render the project into. If omitted, a temporary
                                    directory is used.

```
<!-- END SECTION -->