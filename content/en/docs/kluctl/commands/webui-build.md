---
description: webui command
github_repo: https://github.com/kluctl/kluctl
lastmod: "2023-08-26T09:38:51+02:00"
linkTitle: webui build
path_base_for_github_subdir:
    from: .*
    to: main/docs/kluctl/commands/webui-build.md
title: webui build
weight: 10
---



## Command
<!-- BEGIN SECTION "webui build" "Usage" false -->
Usage: kluctl webui build [flags]

Build the static Kluctl Webui
This command will build the static Kluctl Webui.

<!-- END SECTION -->

## Arguments

The following arguments are available:
<!-- BEGIN SECTION "webui build" "Misc arguments" true -->
```
Misc arguments:
  Command specific arguments.

      --all-contexts          Use all Kubernetes contexts found in the kubeconfig.
      --context stringArray   List of kubernetes contexts to use. Defaults to the current context.
      --max-results int       Specify the maximum number of results per target. (default 1)
      --path string           Output path.

```
<!-- END SECTION -->

