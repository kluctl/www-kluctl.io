# Kluctl website

This repository contains the source code for the Kluctl website. It uses the [Lotusdocs](https://github.com/colinwilson/lotusdocs) theme. A development environment can be set up locally or by using GitHub Codespaces.

## Local Development

### Requirements

Ensure the following prerequisites are met.

- [Go](https://go.dev/dl/) (Golang) version: `1.21.3`
- [Hugo](https://gohugo.io/installation/) extended version: `0.120.4`

> **Note:**
The Hugo extended version is required for various purposes (like transforming SCSS to CSS, converting images to `webp` format). Without it, you may encounter the following error message:

```text
error: failed to transform resource: TOCSS: failed to transform "scss/main.scss" (text/x-scss): this feature is not available in your current Hugo version 
We release two set of binaries for technical reasons. The extended version is not what you get by default for some installation methods. On the release page, look for archives with extended in the name. To build hugo-extended, use go install --tags extended
```

### Setup

1. Clone the repository
   ```bash
   git clone git@github.com:kluctl/www-kluctl.io.git
   cd www-kluctl.io
   ```

2. Run Hugo server:
   ```bash
   hugo server
   ```

   Note: The first run might take a while as Hugo modules are fetched.

3. Open your web browser and go to `http://localhost:1313` to preview the Kluctl website.

## GitHub Codespaces

For GitHub Codespaces, no additional setup is needed. Follow these steps:

1. Start a codespace in the `main` branch using GitHub UI.

2. Start the Hugo server using the following command:
   ```bash
   npm run dev-cp
   ```

   This will start the Hugo server. The server URL will be shown as a notification in codespace. It can also be found ny navigating to the PORTS section of the terminal pane.

Feel free to explore and contribute to the Kluctl website. If you encounter any issues or have questions, please create an issue.