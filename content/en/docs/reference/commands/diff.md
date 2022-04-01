---
title: "diff"
linkTitle: "diff"
weight: 10
description: >
    diff command
---

## Command
<!-- BEGIN SECTION "diff" "Usage" false -->
Usage: kluctl diff

<!-- END SECTION -->

## Arguments
The following sets of arguments are available:
1. [project arguments]({{< ref "./common-arguments#project-arguments" >}})
1. [image arguments]({{< ref "./common-arguments#image-arguments" >}})
1. [inclusion/exclusion arguments]({{< ref "./common-arguments#inclusionexclusion-arguments" >}})

In addition, the following arguments are available:
<!-- BEGIN SECTION "diff" "Misc arguments" true -->
```
Misc arguments:
  Command specific arguments.

      --force-apply                        Force conflict resolution when applying. See documentation for details
      --replace-on-error                   When patching an object fails, try to replace it. See documentation for more
                                           details.
      --force-replace-on-error             Same as --replace-on-error, but also try to delete and re-create objects. See
                                           documentation for more details.
      --ignore-tags                        Ignores changes in tags when diffing
      --ignore-labels                      Ignores changes in labels when diffing
      --ignore-annotations                 Ignores changes in annotations when diffing
  -o, --output-format=OUTPUT-FORMAT,...    Specify output format and target file, in the format 'format=path'. Format
                                           can either be 'text' or 'yaml'. Can be specified multiple times. The actual
                                           format for yaml is currently not documented and subject to change.
      --render-output-dir=STRING           Specifies the target directory to render the project into. If omitted, a
                                           temporary directory is used.

```
<!-- END SECTION -->

`--force-apply` and `--replace-on-error` have the same meaning as in [deploy](#deploy).
