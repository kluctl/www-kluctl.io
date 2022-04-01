---
title: "render"
linkTitle: "render"
weight: 10
description: >
    render command
---

## Command
<!-- BEGIN SECTION "render" "Usage" false -->
Usage: kluctl render

Renders all resources and configuration files

Renders all resources and configuration files and stores the result in either a temporary directory or a specified
directory.

<!-- END SECTION -->

## Arguments
The following sets of arguments are available:
1. [project arguments]({{< ref "./common-arguments#project-arguments" >}})
1. [image arguments]({{< ref "./common-arguments#image-arguments" >}})

In addition, the following arguments are available:
<!-- BEGIN SECTION "render" "Misc arguments" true -->
```
Misc arguments:
  Command specific arguments.

  --render-output-dir=STRING    Specifies the target directory to render the project into. If omitted, a temporary
                                directory is used ($KLUCTL_RENDER_OUTPUT_DIR).

```
<!-- END SECTION -->