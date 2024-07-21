---
description: Installing the Kluctl Controller
github_branch: main
github_repo: https://github.com/kluctl/kluctl
lastmod: "2024-07-20T07:51:52+02:00"
linkTitle: Installation
path_base_for_github_subdir:
    from: .*
    to: docs/gitops/installation.md
title: Installation
weight: 10
---





The controller can be installed via two available options.

## Using the "install" sub-command

The [`kluctl controller install`](../kluctl/commands/controller-install.md) command can be used to install the
controller. It will use an embedded version of the Controller Kluctl deployment project
found [here](https://github.com/kluctl/kluctl/tree/main/install/controller).

## Using a Kluctl deployment

To manage and install the controller via Kluctl, you can use a Git include in your own deployment:

```yaml
deployments:
  - git:
      url: https://github.com/kluctl/kluctl.git
      subDir: install/controller
      ref:
        tag: v2.25.1
```
