# Welcome to KitOps 🚀

## Streamlined Collaboration for AI/ML Projects

KitOps is your toolkit for transforming how you package, share, and deploy AI/ML models. Say goodbye to compatibility concerns and hello to smooth AI/ML collaboration.

KitOps simplifies the handoffs between data scientists, application developers, and SREs working on self-hosted AI/ML models (including LLMs). KitOps' ModelKits create a unified package for models, their dependencies, configurations, and environments. The ModelKit is portable and uses open standards for compatibility with the tools you already use.

https://github.com/jozu-ai/kitops/assets/10517533/8f7539e1-b0d2-43c4-abfe-31841d4c68bd

### What is in the box?

**ModelKit:** At the heart of KitOps is the ModelKit, an OCI-compliant packaging format that enables the seamless sharing of all necessary artifacts involved in the AI/ML model lifecycle. This includes datasets, code, configurations, and the models themselves. By standardizing the way these components are packaged, ModelKit facilitates a more streamlined and collaborative development process that is compatible with nearly any tool.

**Kitfile:** Complementing the ModelKit is the Kitfile, your AI/ML project's blueprint. It's a YAML-based configuration file that simplifies the sharing of model, dataset, and code configurations. Kitfiles are designed with both ease of use and security in mind, ensuring that configurations can be efficiently packaged and shared without compromising on safety or governance.

**Kit CLI:** Your magic wand for AI/ML collaboration. The Kit CLI not only enables users to create, manage, run, and deploy ModelKits... it lets you pull only the pieces you need. Just need the serialized model for deployment? Use `unpack --model` or maybe you just want the training datasets? `unpack --datasets`. So, whether you are packaging a new model for development or deploying an existing model into production, the Kit CLI provides the flexibility and power to streamline your workflow.

## Quick Start with Kit

Dive into the world of KitOps with ease! Whether you're looking to streamline your AI/ML workflows or explore the power of ModelKits, getting started with Kit is straightforward.

### Running Kit with Pre-built Binaries

First, download the Kit CLI. Choose the `latest` [tagged version](https://github.com/jozu-ai/kitops/tags) for the most stable release, or explore the `next` tag for our development builds.

For installation instructions and selecting the right binary for your platform, please refer to our [Installation Guide](./docs/src/docs/cli/installation.md).

To launch Kit, simply open a terminal and type:

```shell
kit
```
This command will display a list of available actions to supercharge your AI/ML projects.

### Building Kit from Source

For those who prefer to build from the source, follow these steps to get the latest version directly from our repository:

1. Clone the Repository: Clone the KitOps source code to your local machine:

```shell
git clone https://github.com/jozu-ai/kitops.git
cd kitops
```

2. Build the Kit CLI: Compile the source code into an executable named kit:
```shell
go build -o kit
```

3. Run Your Build: Execute the built CLI to see all available commands:

```shell
./kit
```

Or, for direct execution during development:

```shell
go run .
```

## Using Kit

The easiest way to get introduced to Kit is with our [quick start guide](./docs/src/docs/getting-started.md)

## Your Voice Matters

### Reporting Issues and Suggesting Features

Your insights help Kit evolve as an open standard for AI/ML. We *deeply value* the issues and feature requests we get from users in our community :sparkling_heart:. To contribute your thoughts,navigate to the **Issues** tab and hitting the **New Issue** green button. Our templates guide you in providing essential details to address your request effectively.

### Joining the KitOps Contributors

We ❤️ our Kit community and contributors. To learn more about the many ways you can contribute (you don't need to be a coder) and how to get started see our [Contributor's Guide](./CONTRIBUTING.md). Please read our [Governance](./GOVERNANCE.md) and our [Code of Conduct](./CODE-OF-CONDUCT.md) before contributing.

### Heed Help?

If you need help there are several ways to reach our community and [Maintainers](./MAINTAINERS.md) outlined in our [support doc](./SUPPORT.md)

### A Community Built on Respect

At KitOps, inclusivity, empathy, and responsibility are at our core. Please read our [Code of Conduct](./CODE-OF-CONDUCT.md) to understand the values guiding our community.

### Join KitOps community

For support, release updates, and general KitOps discussion, please join the [KitOps Discord](https://discord.gg/qCvhJgf2).
