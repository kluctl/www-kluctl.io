---
title: "Common Arguments"
linkTitle: "Common Arguments"
weight: 1
description: >
    Common arguments
---

A few sets of arguments are common between multiple commands. These arguments are still part of the command itself and
must be placed *after* the command name.

## Global arguments

These arguments are available for all commands.

<!-- BEGIN SECTION "deploy" "Global arguments" true -->
```
Global arguments:
      --debug             Enable debug logging
      --no-update-check   Disable update check on startup

```
<!-- END SECTION -->

## Project arguments

These arguments are available for all commands that are based on a Kluctl project.
They control where and how to load the kluctl project and deployment project.

<!-- BEGIN SECTION "deploy" "Project arguments" true -->
```
Project arguments:
  Define where and how to load the kluctl project and its components from.

  -a, --arg strings                          Template argument in the form name=value
      --cluster string                       Specify/Override cluster
      --from-archive existingpath            Load project (.kluctl.yml, cluster, ...) from archive. Given path can
                                             either be an archive file or a directory with the extracted contents.
      --from-archive-metadata existingfile   Specify where to load metadata (targets, ...) from. If not specified,
                                             metadata is assumed to be part of the archive.
      --git-cache-update-interval duration   Specify the time to wait between git cache updates. Defaults to not
                                             wait at all and always updating caches.
      --local-clusters existingdir           Local clusters directory. Overrides the project from .kluctl.yml
      --local-deployment existingdir         Local deployment directory. Overrides the project from .kluctl.yml
      --local-sealed-secrets existingdir     Local sealed-secrets directory. Overrides the project from .kluctl.yml
      --output-metadata path                 Specify the output path for the project metadata to be written to.
                                             When used with the 'archive' command, it will also cause the archive
                                             to not include the metadata.yaml file.
  -c, --project-config existingfile          Location of the .kluctl.yml config file. Defaults to $PROJECT/.kluctl.yml
  -b, --project-ref string                   Git ref of the kluctl project. Only used when --project-url was given.
  -p, --project-url string                   Git url of the kluctl project. If not specified, the current
                                             directory will be used instead of a remote Git project
  -t, --target string                        Target name to run command for. Target must exist in .kluctl.yml.
      --timeout duration                     Specify timeout for all operations, including loading of the project,
                                             all external api calls and waiting for readiness. (default 10m0s)

```
<!-- END SECTION -->

## Image arguments

These arguments are available on some target based commands.
They control image versions requested by `images.get_image(...)` [calls]({{< ref "docs/reference/deployments/images#imagesget_image" >}}).

<!-- BEGIN SECTION "deploy" "Image arguments" true -->
```
Image arguments:
  Control fixed images and update behaviour.

  -F, --fixed-image strings              Pin an image to a given version. Expects
                                         '--fixed-image=image<:namespace:deployment:container>=result'
      --fixed-images-file existingfile   Use .yml file to pin image versions. See output of list-images
                                         sub-command or read the documentation for details about the output format
  -u, --update-images                    This causes kluctl to prefer the latest image found in registries, based
                                         on the 'latest_image' filters provided to 'images.get_image(...)' calls.
                                         Use this flag if you want to update to the latest versions/tags of all
                                         images. '-u' takes precedence over '--fixed-image/--fixed-images-file',
                                         meaning that the latest images are used even if an older image is given
                                         via fixed images.

```
<!-- END SECTION -->

## Inclusion/Exclusion arguments

These arguments are available for some target based commands.
They control inclusion/exclusion based on tags and deployment item pathes.

<!-- BEGIN SECTION "deploy" "Inclusion/Exclusion arguments" true -->
```
Inclusion/Exclusion arguments:
  Control inclusion/exclusion.

      --exclude-deployment-dir strings   Exclude deployment dir. The path must be relative to the root deployment
                                         project. Exclusion has precedence over inclusion, same as in --exclude-tag
  -E, --exclude-tag strings              Exclude deployments with given tag. Exclusion has precedence over
                                         inclusion, meaning that explicitly excluded deployments will always be
                                         excluded even if an inclusion rule would match the same deployment.
      --include-deployment-dir strings   Include deployment dir. The path must be relative to the root deployment
                                         project.
  -I, --include-tag strings              Include deployments with given tag.

```
<!-- END SECTION -->
