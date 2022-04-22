---
title: "archive"
linkTitle: "archive"
weight: 10
description: >
    archive command
---

## Command
<!-- BEGIN SECTION "archive" "Usage" false -->
Usage: kluctl archive

Write project and all related components into single tgz

This archive can then be used with '--from-archive'

<!-- END SECTION -->

## Arguments
The following arguments are available:
<!-- BEGIN SECTION "archive" "Misc arguments" true -->
```
Misc arguments:
  Command specific arguments.

  --output-archive=STRING     Path to .tgz to write project to ($KLUCTL_OUTPUT_ARCHIVE).
  --output-metadata=STRING    Path to .yml to write metadata to. If not specified, metadata is written into the archive
                              ($KLUCTL_OUTPUT_METADATA).

```
<!-- END SECTION -->