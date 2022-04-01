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

Validates the already deployed deployment

This means that all objects are retrieved from the cluster and checked for readiness.

TODO: This needs to be better documented!

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

  -o, --output=OUTPUT,...           Specify output target file. Can be specified multiple times ($KLUCTL_OUTPUT)
      --render-output-dir=STRING    Specifies the target directory to render the project into. If omitted, a temporary
                                    directory is used ($KLUCTL_RENDER_OUTPUT_DIR).
      --wait=DURATION               Wait for the given amount of time until the deployment validates ($KLUCTL_WAIT)
      --sleep=5s                    Sleep duration between validation attempts ($KLUCTL_SLEEP)
      --warnings-as-errors          Consider warnings as failures ($KLUCTL_WARNINGS_AS_ERRORS)

```
<!-- END SECTION -->
