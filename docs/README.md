# KitOps Documentation

This is the documentation for [KitsOps](https://kitops.ml). You can read the docs at https://docs.kitops.ml.

## Table of Contents

- [Introduction](#introduction)
- [Development](#development)
  * [Prerequisites](#prerequisites)
  * [Installation](#installation)
- [Contributing Guidelines](#contributing-guidelines)
  + [Reporting Bugs and Issues](#reporting-bugs-and-issues)
  + [Submitting Pull Requests](#submitting-pull-requests)

## Introduction

This documentation has been built using [VitePress](https://vitepress.dev/). VitePress is a fast static site generator for building modern documentations using Vue.js. For more information please refer to the [VitePress Website](https://vitepress.dev/).

## Development

To get started with our, follow these steps:

### Prerequisites

Ensure you have the following installed on your system:
* Node.js (v18.x or higher)
* npm or pnpm package manager

For local development, [pnpm](https://pnpm.io/) is preferred as package manager.

### Running locally

1. Clone the repository:
   ```sh
   git clone https://github.com/jozu-ai/kitops.git
   ```
2. Navigate to the docs directory:
   ```sh
   cd docs
   ```
3. Install dependencies:
   ```sh
   pnpm i
   ```
4. Start the development server:
   ```sh
   pnpm docs:dev
   ```

The documentation should now be available at `http://localhost:5173`.

## Contributing Guidelines

We welcome contributions from the community to help improve our project and documentation. Please follow these guidelines when contributing:

### Reporting Bugs and Issues

If you encounter any bugs or issues with the documentation, please report them in our [GitHub issue tracker](https://github.com/jozu-ai/kitops/issues). Be sure to provide clear details about the problem, including steps to reproduce if possible.

### Submitting Pull Requests

If you'd like to submit a pull request for updates or improvements to the documentation, please follow these steps:

1. Fork the repository on GitHub.
2. Create a new branch for your changes: `git checkout -b feature/your-feature`
3. Make the necessary changes and additions to the documentation.
4. Commit your changes: `git commit -m 'Add documentation for XYZ'`
5. Push the branch to your forked repository: `git push origin feature/your-feature`
6. Submit a pull request on GitHub, detailing the changes you made.
