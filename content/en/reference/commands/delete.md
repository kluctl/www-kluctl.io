---
title: "delete"
linkTitle: "delete"
weight: 10
description: >
    delete command
---

## Command
<!-- BEGIN SECTION "delete" "Usage" false -->
Usage: kluctl delete

Delete a target (or parts of it) from the corresponding cluster

Objects are located based on 'commonLabels”, configured in 'deployment.yml'

WARNING: This command will also delete objects which are not part of your deployment project (anymore). It really only
decides based on the 'deleteByLabel' labels and does NOT take the local target/state into account!

<!-- END SECTION -->

## Arguments
The following sets of arguments are available:
1. [project arguments]({{< ref "./common-arguments#project-arguments" >}})
1. [image arguments]({{< ref "./common-arguments#image-arguments" >}})
1. [inclusion/exclusion arguments]({{< ref "./common-arguments#inclusionexclusion-arguments" >}})

In addition, the following arguments are available:
<!-- BEGIN SECTION "delete" "Misc arguments" true -->
```
Misc arguments:
  Command specific arguments.

  -y, --yes                                    Suppresses 'Are you sure?' questions and proceeds as if you would answer
                                               'yes' ($KLUCTL_YES).
      --dry-run                                Performs all kubernetes API calls in dry-run mode ($KLUCTL_DRY_RUN).
  -o, --output-format=OUTPUT-FORMAT,...        Specify output format and target file, in the format 'format=path'.
                                               Format can either be 'text' or 'yaml'. Can be specified multiple times.
                                               The actual format for yaml is currently not documented and subject to
                                               change ($KLUCTL_OUTPUT_FORMAT).
  -l, --delete-by-label=DELETE-BY-LABEL,...    Override the labels used to find objects for deletion
                                               ($KLUCTL_DELETE_BY_LABEL).

```
<!-- END SECTION -->

They have the same meaning as described in [deploy](#deploy).
