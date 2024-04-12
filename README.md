
<img width="1270" alt="KitOps" src="https://github.com/jozu-ai/kitops/assets/10517533/41295471-fe49-4011-adf6-a215f29890c2">


## The world needs a standard packaging / versioning system for AI/ML projects.

[![LICENSE](https://img.shields.io/badge/License-Apache%202.0-yellow.svg)](https://github.com/myscale/myscaledb/blob/main/LICENSE)
[![Language](https://img.shields.io/badge/Language-go-blue.svg)](https://go.dev/)
[![Discord](https://img.shields.io/discord/1098133460310294528?logo=Discord)](https://discord.gg/Tapeh8agYy)
[![Twitter](https://img.shields.io/twitter/url/http/shields.io.svg?style=social&label=Twitter)](https://twitter.com/kit_ops)
[![Hits](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fjozu-ai%2Fkitops&count_bg=%2379C83D&title_bg=%23555555&icon=&icon_color=%23E7E7E7&title=hits&edge_flat=false)](https://hits.seeyoufarm.com)

[![Official Website](<https://img.shields.io/badge/-Visit%20the%20Official%20Website%20%E2%86%92-rgb(255,175,82)?style=for-the-badge>)](https://kitops.ml/?utm_source=github&utm_medium=kitops-readme)
[![Use Cases](<https://img.shields.io/badge/-KitOps%20Quick%20Start%20%E2%86%92-rgb(122,140,225)?style=for-the-badge>)](https://kitops.ml/docs/quick-start.html/?utm_source=github&utm_medium=kitops-readme)

### What is KitOps?

KitOps is a packaging and versioning system for AI/ML projects that uses open standards so it works with the AI/ML, development, and DevOps tools you are already using.

KitOps simplifies the handoffs between data scientists, application developers, and SREs working with LLMs and other AI/ML models. KitOps' ModelKits are a standards-based package for models, their dependencies, configurations, and codebases. ModelKits are portable, reproducible, and work with the tools you already use.

### Features

* 🎁 **[Unified packaging](https://kitops.ml/docs/modelkit/intro.html):** A ModelKit package includes models, datasets, configurations, and code. Add as much or as little as your project needs.
* 🏭 **[Versioning](https://kitops.ml/docs/cli/cli-reference.html#kit-tag):** Each ModelKit is tagged so everyone knows which dataset and model work together.
* 🤖 **[Automation]](https://github.com/marketplace/actions/setup-kit-cli):** Pack or unpack a ModelKit locally or as part of your CI/CD workflow for testing, integration, or deployment.
* 🔒 **[Tamper-proofing](https://kitops.ml/docs/modelkit/spec.html):** Each ModelKit package includes a SHA digest for itself, and every artifact it holds.
* 🌈 **[Standards-based](https://kitops.ml/docs/modelkit/compatibility.html):** Store ModelKits in any container or artifact registry.
* 🥧 **[Simple syntax](https://kitops.ml/docs/kitfile/kf-overview.html):** Kitfiles are easy to write and read, using a familiar YAML syntax.
* 😻 **No GPU or internet:** Kit doesn't require GPUs, internet connectivity, your email, or favourite limb. It's a free tool you can use anywhere.
* 🤗 **Flexible:** ModelKits can be used with any AI, ML, or LLM project - even multi-modal models.
* 🧰 **Data science + DevOps:** Simplify asset management and versioning for training, experimentation, integration, deployment, and operations.
* 🏃‍♂️‍➡️ **Run locally:** Kit's Dev Mode lets your run an LLM locally, configure it, and prompt/chat with it instantly (coming soon).
* 🐳 **Deploy containers:** Generate a Docker container as part of your `kit unpack` (coming soon).
* 🚢 **Kubernetes-ready:** Generate a Kubernetes / KServe deployment config as part of your `kit unpack` (coming soon).
* 📝 **Signed packages:** ModelKits and their assets can be signed so you can be confident of their provenance.

### See KitOps in Action

https://github.com/jozu-ai/kitops/assets/4766570/05ae1362-afd3-4e78-bfce-e982c17f8df2

### What is in the box?

**ModelKit:** At the heart of KitOps is the ModelKit, an OCI-compliant packaging format that enables the seamless sharing of all necessary artifacts involved in the AI/ML model lifecycle. This includes datasets, code, configurations, and the models themselves. By standardizing the way these components are packaged, ModelKit facilitates a more streamlined and collaborative development process that is compatible with nearly any tool.

**Kitfile:** Complementing the ModelKit is the Kitfile, your AI/ML project's blueprint. It's a YAML-based configuration file that simplifies the sharing of model, dataset, and code configurations. Kitfiles are designed with both ease of use and security in mind, ensuring that configurations can be efficiently packaged and shared without compromising on safety or governance.

**Kit CLI:** Your magic wand for AI/ML collaboration. The Kit CLI not only enables users to create, manage, run, and deploy ModelKits... it lets you pull only the pieces you need. Just need the serialized model for deployment? Use `unpack --model` or maybe you just want the training datasets? `unpack --datasets`. So, whether you are packaging a new model for development or deploying an existing model into production, the Kit CLI provides the flexibility and power to streamline your workflow.

## Try Kit in under 15 Minutes

First, download the Kit CLI. Choose the `latest` [tagged version](https://github.com/jozu-ai/kitops/tags) for the most stable release, or explore the `next` tag for our development builds.

For installation instructions and selecting the right binary for your platform, please refer to our [Installation Guide](./docs/src/docs/cli/installation.md).

To launch Kit, simply open a terminal and type:

```shell
kit
```
This command will display a list of available actions to supercharge your AI/ML projects.

The [Kit Quick Start](https://kitops.ml/docs/quick-start.html) will guide you through the main features of kit in under 10 minutes.

## Building Kit from Source

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

## Your Voice Matters

### Reporting Issues and Suggesting Features

Your insights help Kit evolve as an open standard for AI/ML. We *deeply value* the issues and feature requests we get from users in our community :sparkling_heart:. To contribute your thoughts,navigate to the **Issues** tab and hitting the **New Issue** green button. Our templates guide you in providing essential details to address your request effectively.

### Joining the KitOps Contributors

We ❤️ our Kit community and contributors. To learn more about the many ways you can contribute (you don't need to be a coder) and how to get started see our [Contributor's Guide](./CONTRIBUTING.md). Please read our [Governance](./GOVERNANCE.md) and our [Code of Conduct](./CODE-OF-CONDUCT.md) before contributing.

### Need Help?

If you need help there are several ways to reach our community and [Maintainers](./MAINTAINERS.md) outlined in our [support doc](./SUPPORT.md)

### A Community Built on Respect

At KitOps, inclusivity, empathy, and responsibility are at our core. Please read our [Code of Conduct](./CODE-OF-CONDUCT.md) to understand the values guiding our community.

### Join KitOps community

For support, release updates, and general KitOps discussion, please join the [KitOps Discord](https://discord.gg/Tapeh8agYy). Follow [KitOps on X](https://twitter.com/Kit_Ops) for daily updates.

<a href="https://trackgit.com">Trackgit</a>
