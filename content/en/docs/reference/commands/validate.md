---
title: "validate"
linkTitle: "validate"
weight: 10
description: >
    validate command
---

## Command
<!-- BEGIN SECTION "validate" "Usage" false -->
Usage: kluctl validate

<!-- END SECTION -->

## Arguments
The following sets of arguments are available:
1. [project arguments]({{< ref "./common-arguments#project-arguments" >}})
1. [image arguments]({{< ref "./common-arguments#image-arguments" >}})

In addition, the following arguments are available:
<!-- BEGIN SECTION "validate" "Misc arguments" true -->
```
Misc arguments:
  Command specific arguments.

  -o, --output=OUTPUT,...           Specify output target file. Can be specified multiple times
      --render-output-dir=STRING    Specifies the target directory to render the project into. If omitted, a temporary
                                    directory is used.
      --wait=DURATION               Wait for the given amount of time until the deployment validates
      --sleep=5s                    Sleep duration between validation attempts
      --warnings-as-errors          Consider warnings as failures

```
<!-- END SECTION -->