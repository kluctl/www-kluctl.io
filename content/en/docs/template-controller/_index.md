---
description: Template Controller documentation.
github_repo: https://github.com/kluctl/template-controller
lastmod: "2022-12-29T10:41:22+01:00"
linkTitle: Template Controller
path_base_for_github_subdir:
    from: .*
    to: main/README.md
title: Template Controller
weight: 40
---





The Template Controller is a controller originating from the [Kluctl](https://kluctl.io) project, but not limited to
Kluctl. It allows to define template objects which are rendered and applied into the cluster based on an input matrix.

In its easiest form, an `ObjectTemplate` takes one input object (e.g. a ConfigMap) and creates another object
(e.g. a Secret) which is then applied into the cluster.

The Template Controller also offers CRDs which allow to query external resources (e.g. GitHub Pull Requests) which can
then be used as inputs into `ObjectTemplates`.

## Use Cases

Template Controller has many use case. Some are for example:
1. [Dynamic environments for Pull Requests](./use-case-dynamic-environments.md)
2. [Transformation of Secrets/Objects](./use-case-transformation.md)

## Documentation

Reference documentation is available [here](./spec/v1alpha1).

The [announcement blog post](https://blog.kluctl.io/introducing-the-template-controller-and-building-gitops-preview-environments-2cce4041406a) also contains valuable explanations
and examples.

## Installation

Installation instructions can be found [here](./install.md)
