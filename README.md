# kluctl.io website and documentation

[![Netlify Status](https://api.netlify.com/api/v1/badges/02a02d5e-b26e-495c-be3d-946de035115b/deploy-status)](https://app.netlify.com/sites/kluctl/deploys)

Built with [Docsy](https://github.com/google/docsy)

# Launch website locally

requirements:
- **[hugo](https://github.com/gohugoio/hugo/releases)** extended version
  - We need to process SCSS or SASS to CSS in our Hugo project, 
    you need the Hugo extended version, or else you may see this error message:
    ```
    error: failed to transform resource: TOCSS: failed to transform "scss/main.scss" (text/x-scss): this feature is not available in your current Hugo version 
    We release two set of binaries for technical reasons. The extended version is not what you get by default for some installation methods. On the release page, look for archives with extended in the name. To build hugo-extended, use go install --tags extended
    ```
- **[npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)**

```bash
git clone git@github.com:kluctl/www-kluctl.io.git
git submodule update --init --recursive
npm install
hugo server
```
