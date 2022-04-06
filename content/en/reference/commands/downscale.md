---
title: "downscale"
linkTitle: "downscale"
weight: 10
description: >
    downscale command
---

## Command
<!-- BEGIN SECTION "downscale" "Usage" false -->
Usage: kluctl downscale

Downscale all deployments

This command will downscale all Deployments, StatefulSets and CronJobs. It is also possible to influence the behaviour
with the help of annotations, as described in the documentation.

<!-- END SECTION -->

## Arguments
The following sets of arguments are available:
1. [project arguments]({{< ref "./common-arguments#project-arguments" >}})
1. [image arguments]({{< ref "./common-arguments#image-arguments" >}})
1. [inclusion/exclusion arguments]({{< ref "./common-arguments#inclusionexclusion-arguments" >}})

In addition, the following arguments are available:
<!-- BEGIN SECTION "downscale" "Misc arguments" true -->
```
Misc arguments:
  Command specific arguments.

  -y, --yes                                Suppresses 'Are you sure?' questions and proceeds as if you would answer
                                           'yes' ($KLUCTL_YES).
      --dry-run                            Performs all kubernetes API calls in dry-run mode ($KLUCTL_DRY_RUN).
  -o, --output-format=OUTPUT-FORMAT,...    Specify output format and target file, in the format 'format=path'. Format
                                           can either be 'text' or 'yaml'. Can be specified multiple times. The actual
                                           format for yaml is currently not documented and subject to change
                                           ($KLUCTL_OUTPUT_FORMAT).
      --render-output-dir=STRING           Specifies the target directory to render the project into. If omitted, a
                                           temporary directory is used ($KLUCTL_RENDER_OUTPUT_DIR).

```
<!-- END SECTION -->