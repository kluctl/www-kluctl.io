---
description: Annotations on hooks
github_repo: https://github.com/kluctl/kluctl
lastmod: "2023-09-15T15:13:22+02:00"
linkTitle: Hooks
path_base_for_github_subdir:
    from: .*
    to: main/docs/kluctl/deployments/annotations/hooks.md
title: Hooks
weight: 2
---





The following annotations control hook execution

See [hooks](../../deployments/hooks.md) for more details.

### kluctl.io/hook
Declares a resource to be a hook, which is deployed/executed as described in [hooks](../../deployments/hooks.md). The value of the
annotation determines when the hook is deployed/executed.

### kluctl.io/hook-weight
Specifies a weight for the hook, used to determine deployment/execution order.

### kluctl.io/hook-delete-policy
Defines when to delete the hook resource.

### kluctl.io/hook-wait
Defines whether kluctl should wait for hook-completion. It defaults to `true` and can be manually set to `false`.
