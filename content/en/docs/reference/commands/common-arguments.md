---
title: "Common Arguments"
linkTitle: "Common Arguments"
weight: 1
description: >
    Common arguments
---

A few sets of arguments are common between multiple commands. These arguments are still part of the command itself and
must be placed *after* the command name.

The following sets of common arguments are available:

## project arguments
<!-- BEGIN SECTION "deploy" "Project arguments" true -->
```
Project arguments:
  Define where and how to load the kluctl project and its components from.

  -p, --project-url=STRING              Git url of the kluctl project. If not specified, the current directory will be
                                        used instead of a remote Git project
  -b, --project-ref=STRING              Git ref of the kluctl project. Only used when --project-url was given.
  -c, --project-config=STRING           Location of the .kluctl.yml config file. Defaults to $PROJECT/.kluctl.yml
      --local-clusters=STRING           Local clusters directory. Overrides the project from .kluctl.yml
      --local-deployment=STRING         Local deployment directory. Overrides the project from .kluctl.yml
      --local-sealed-secrets=STRING     Local sealed-secrets directory. Overrides the project from .kluctl.yml
      --from-archive=STRING             Load project (.kluctl.yml, cluster, ...) from archive. Given path can either be
                                        an archive file or a directory with the extracted contents.
      --from-archive-metadata=STRING    Specify where to load metadata (targets, ...) from. If not specified, metadata
                                        is assumed to be part of the archive.
      --cluster=STRING                  Specify/Override cluster
  -t, --target=STRING                   Target name to run command for. Target must exist in .kluctl.yml.
  -a, --arg=ARG,...                     Template argument in the form name=value

```
<!-- END SECTION -->

These arguments control where and how to load the kluctl project and deployment project.

## image arguments
<!-- BEGIN SECTION "deploy" "Image arguments" true -->
```
Image arguments:
  Control fixed images and update behaviour.

  -F, --fixed-image=FIXED-IMAGE,...    Pin an image to a given version. Expects
                                       '--fixed-image=image<:namespace:deployment:container>=result'
      --fixed-images-file=STRING       Use .yml file to pin image versions. See output of list-images sub-command or
                                       read the documentation for details about the output format
  -u, --update-images                  This causes kluctl to prefer the latest image found in registries, based on the
                                       'latest_image' filters provided to 'images.get_image(...)' calls. Use this flag
                                       if you want to update to the latest versions/tags of all images. '-u' takes
                                       precedence over '--fixed-image/--fixed-images-file', meaning that the latest
                                       images are used even if an older image is given via fixed images.

```
<!-- END SECTION -->

These arguments control image versions requested by `images.get_image(...)` [calls](./images.md#imagesget_image).

## Inclusion/Exclusion arguments
<!-- BEGIN SECTION "deploy" "Inclusion/Exclusion arguments" true -->
```
Inclusion/Exclusion arguments:
  Control inclusion/exclusion.

  -I, --include-tag=INCLUDE-TAG,...    Include deployments with given tag.
  -E, --exclude-tag=EXCLUDE-TAG,...    Exclude deployments with given tag. Exclusion has precedence over inclusion,
                                       meaning that explicitly excluded deployments will always be excluded even if an
                                       inclusion rule would match the same deployment.
      --include-deployment-dir=INCLUDE-DEPLOYMENT-DIR,...
                                       Include deployment dir. The path must be relative to the root deployment project.
      --exclude-deployment-dir=EXCLUDE-DEPLOYMENT-DIR,...
                                       Exclude deployment dir. The path must be relative to the root deployment project.
                                       Exclusion has precedence over inclusion, same as in --exclude-tag

```
<!-- END SECTION -->

These arguments control inclusion/exclusion based on tags and kustomize depoyment pathes.
