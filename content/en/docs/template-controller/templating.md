---
description: Templating documentation.
github_branch: main
github_repo: https://github.com/kluctl/template-controller
lastmod: "2024-07-12T11:48:21+02:00"
path_base_for_github_subdir:
    from: .*
    to: docs/templating.md
title: Templating
weight: 30
---





The Template Controller reuses the Jinja2 templating engine of [Kluctl](https://kluctl.io).

Documentation is available [here](https://kluctl.io/docs/kluctl/templating/).

## Predefined variables

You can use multiple predefined variables in your templates. These are:

### objectTemplate

Available in templates inside [ObjectTemplate](./spec/v1alpha1/objecttemplate.md) and represents the whole
`ObjectTemplate` that was on your target BEFORE the reconciliation started.

### textTemplate

Available in templates inside [TextTemplate](./spec/v1alpha1/texttemplate.md) and represents the whole
`TextTemplate` that was on your target BEFORE the reconciliation started.
