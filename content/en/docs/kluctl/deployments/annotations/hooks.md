---
description: Annotations on hooks
github_branch: main
github_repo: https://github.com/kluctl/kluctl
lastmod: "2024-06-18T14:54:41+02:00"
linkTitle: Hooks
path_base_for_github_subdir:
    from: .*
    to: docs/kluctl/deployments/annotations/hooks.md
title: Hooks
weight: 2
---





The following annotations control hook execution

See [hooks](../../deployments/hooks.md) for more details.

### kluctl.io/hook
Declares a resource to be a hook, which is deployed/executed as described in [hooks](../../deployments/hooks.md). The value of the
annotation determines when the hook is deployed/executed.

### kluctl.io/hook-weight
Specifies a weight for the hook, used to determine deployment/execution order. For resources with the same `kluctl.io/hook` annotation, hooks are executed in ascending order based on hook-weight.

### kluctl.io/hook-delete-policy
Defines when to delete the hook resource.

### kluctl.io/hook-wait
Defines whether kluctl should wait for hook-completion. It defaults to `true` and can be manually set to `false`.
