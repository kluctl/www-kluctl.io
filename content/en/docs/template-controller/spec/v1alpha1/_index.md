---
description: templates.kluctl.io/v1alpha1 documentation
github_branch: main
github_repo: https://github.com/kluctl/template-controller
lastmod: "2023-01-16T13:44:24+01:00"
linkTitle: v1alpha1 specs
path_base_for_github_subdir:
    from: .*
    to: docs/spec/v1alpha1/README.md
title: v1alpha1 specs
weight: 10
---



# templates.kluctl.io/v1alpha1

This is the v1alpha1 API specification for defining templating related resources.

## Specification

- [ObjectTemplate CRD](objecttemplate.md)
    + [Spec fields](objecttemplate.md#spec-fields)
- [TextTemplate CRD](texttemplate.md)
    + [Spec fields](objecttemplate.md#spec-fields)
- [GitProjector CRD](gitprojector.md)
    + [Spec fields](gitprojector.md#spec-fields)
- [ListGithubPullRequests CRD](listgithubpullrequests.md)
    + [Spec fields](listgithubpullrequests.md#spec-fields)
- [ListGitlabMergeRequests CRD](listgitlabmergerequests.md)
    + [Spec fields](listgitlabmergerequests.md#spec-fields)
- [GithubComment CRD](githubcomment.md)
    + [Spec fields](githubcomment.md#spec-fields)
- [GitlabComment CRD](gitlabcomment.md)
    + [Spec fields](gitlabcomment.md#spec-fields)

## Implementation

* [template-controller](https://github.com/kluctl/template-controller)
