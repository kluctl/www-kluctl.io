---
title: "Templating"
linkTitle: "Templating"
weight: 2
description: >
    Templating Engine.
---

kluctl uses a Jinja2 Templating engine to pre-process/render every involved configuration file and resource before
actually interpreting it. Only files that are explicitly excluded via [templateExcludes]({{< ref "reference/deployments/deployment-yml#templateexcludes" >}})
are not rendered via Jinja2.

Generally, everything that is possible with Jinja2 is possible in kluctl configuration/resources. Please
read into the [Jinja2 documentation](https://jinja.palletsprojects.com/en/3.0.x/templates/) to understand what exactly
is possible and how to use it.

## Includes and imports
Standard Jinja2 [includes](https://jinja.palletsprojects.com/en/2.11.x/templates/#include) and
[imports](https://jinja.palletsprojects.com/en/2.11.x/templates/#import) can be used in all templates.

The path given to include/import is treated as relative to the template that is currently rendered.

## Macros

[Jinja2 macros](https://jinja.palletsprojects.com/en/2.11.x/templates/#macros) are fully supported. When writing
macros that produce yaml resources, you must use the `---` yaml separator in case you want to produce multiple resources
in one go.

## Why no Go Templating

kluctl started as a python project and was then migrated to be a Go project. In the python world, Jinja2 is the obvious
choice when it comes to templating. In the Go world, of course Go Templates would be the first choice.

When the migration to Go was performed, it was a conscious and opinionated decision to stick with Jinja2 templating.
The reason is that I (@codablock) believe that Go Templates are hard to read and write and at the same time quite limited
in their features (without extensive work). It never felt natural to write Go Templates.

This "feeling" was confirmed by multiple users of kluctl when it started and users described as "relieving" to not
be forced to use Go Templates.

The above is my personal experience and opinion. I'm still quite open for contributions in regard to Go Templating
support, as long as Jinja2 support is kept.
