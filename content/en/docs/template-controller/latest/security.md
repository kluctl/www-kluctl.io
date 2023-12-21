---
description: Security documentation.
github_repo: https://github.com/kluctl/template-controller
lastmod: "2022-12-14T16:47:29+01:00"
path_base_for_github_subdir:
    from: .*
    to: main/docs/security.md
title: Security
weight: 20
---





The Template Controller is a powerful controller that is able to create/apply arbitrary objects from templates and an
input matrix. This has some security implications as it requires you to make sure that you don't open potential
security vulnerabilities inside your cluster.

This means, you must make sure that your `ObjectTemplate` objects are either not dependent on external inputs (which
might contain malicious input) or tha the used [service account](./spec/v1alpha1/objecttemplate.md#serviceaccountname)
is restricted enough to not allow malicious modifications to the cluster.

## cluster-admin role

Especially watch out when using the cluster-admin (or comparable) role. It can easily lead to privilege escalation if
templates and inputs are too dynamic. 
