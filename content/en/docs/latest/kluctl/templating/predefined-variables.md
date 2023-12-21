---
description: Available predefined variables.
github_repo: https://github.com/kluctl/kluctl
lastmod: "2023-08-26T09:38:51+02:00"
linkTitle: Predefined Variables
path_base_for_github_subdir:
    from: .*
    to: main/docs/kluctl/templating/predefined-variables.md
title: Predefined Variables
weight: 1
---





There are multiple variables available which are pre-defined by kluctl. These are:

### args
This is a dictionary of arguments given via command line. It contains every argument defined in
[deployment args](../deployments/deployment-yml.md#args).

### target
This is the target definition of the currently processed target. It contains all values found in the 
[target definition](../kluctl-project/targets), for example `target.name`.

### images
This global object provides the dynamic images features described in [images](../deployments/images.md).
