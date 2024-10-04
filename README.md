
<img width="1270" alt="KitOps" src="https://github.com/jozu-ai/kitops/assets/10517533/41295471-fe49-4011-adf6-a215f29890c2">


## Standards-based packaging and versioning system for AI/ML projects.

[![LICENSE](https://img.shields.io/badge/License-Apache%202.0-yellow.svg)](https://github.com/myscale/myscaledb/blob/main/LICENSE)
[![Language](https://img.shields.io/badge/Language-go-blue.svg)](https://go.dev/)
[![Discord](https://img.shields.io/discord/1098133460310294528?logo=Discord)](https://discord.gg/Tapeh8agYy)
[![Twitter](https://img.shields.io/twitter/url/http/shields.io.svg?style=social&label=Twitter)](https://twitter.com/kit_ops)
[![Hits](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fjozu-ai%2Fkitops&count_bg=%2379C83D&title_bg=%23555555&icon=&icon_color=%23E7E7E7&title=hits&edge_flat=false)](https://hits.seeyoufarm.com)

[![Official Website](<https://img.shields.io/badge/-Visit%20the%20Official%20Website%20%E2%86%92-rgb(255,175,82)?style=for-the-badge>)](https://kitops.ml?utm_source=github&utm_medium=kitops-readme)

[![Use Cases](<https://img.shields.io/badge/-KitOps%20Quick%20Start%20%E2%86%92-rgb(122,140,225)?style=for-the-badge>)](https://kitops.ml/docs/get-started.html?utm_source=github&utm_medium=kitops-readme)

### What is KitOps?

KitOps is a packaging and versioning system for AI/ML projects that uses open standards so it works with the AI/ML, development, and DevOps tools you are already using, and can be stored in your enterprise registry.

KitOps makes it easy for organizations to track, control, and audit access and changes to their AI project artifacts. It simplifies the handoffs between data scientists, application developers, and SREs working with LLMs and other AI/ML models. KitOps' ModelKits are an OCI-compliant package for models, their dependencies, configurations, codebases, features, hyperparmeters, and any other documentation. ModelKits are portable, reproducible, and work with the tools you already use.

Teams and enterprises use KitOps as a gate between development and production deployment. This ensures that projects are complete, versioned, and compatible with existing pipelines and tools. This makes deployments, rollbacks, and other production operations simpler and safer.

Use KitOps to speed up and de-risk all types of AI projects from small analysis models to large language models, including fine-tuning and RAG.

### 🇪🇺 EU AI Act Compliance 🔒
For our friends in the EU - ModelKits are the perfect way to create a library of model versions for EU AI Act compliance because they're tamper-proof, signable, and auditable.


### 🎉 New

Get the most out of KitOps' ModelKits by using them with the **[Jozu Hub](https://jozu.ml/)** repository. Or, continue using ModelKits with your existing OCI registry (even on-premises and air-gapped).


### Features

* 🎁 **[Unified packaging](https://kitops.ml/docs/modelkit/intro.html):** A ModelKit package includes models, datasets, configurations, and code. Add as much or as little as your project needs.
* 🏭 **[Versioning](https://kitops.ml/docs/cli/cli-reference.html#kit-tag):** Each ModelKit is tagged so everyone knows which dataset and model work together.
* 🤖 **[Automation](https://github.com/marketplace/actions/setup-kit-cli):** Pack or unpack a ModelKit locally or as part of your CI/CD workflow for testing, integration, or deployment.
* 🪛 **[LLM fine-tuning](https://dev.to/kitops/fine-tune-your-first-large-language-model-llm-with-lora-llamacpp-and-kitops-in-5-easy-steps-1g7f):** Use KitOps to fine-tune a large language model using LoRA.
* 🎯 **[RAG pipelines](https://www.codeproject.com/Articles/5384392/A-Step-by-Step-Guide-to-Building-and-Distributing):** Create a RAG pipeline for tailoring an LLM with KitOps.
* 🔒 **[Tamper-proofing](https://kitops.ml/docs/modelkit/spec.html):** Each ModelKit package includes a SHA digest for itself, and every artifact it holds.
* 📝 **[Artifact signing](https://kitops.ml/docs/next-steps.html):** ModelKits and their assets can be signed so you can be confident of their provenance.
* 🌈 **[Standards-based](https://kitops.ml/docs/modelkit/compatibility.html):** Store ModelKits in any OCI 1.1-compliant container or artifact registry.
* 🥧 **[Simple syntax](https://kitops.ml/docs/kitfile/kf-overview.html):** Kitfiles are easy to write and read, using a familiar YAML syntax.
* 🩰 **[Flexible](https://kitops.ml/docs/kitfile/format.html#model):** Reference base models using `model parts`, or store key-value pairs (or any YAML-compatible JSON data) in your Kitfile - use it to keep features, hyperparameters, links to MLOps tool experiments, or validation output.
* 🏃‍♂️‍➡️ **[Run locally](./docs/src/docs/dev-mode.md):** Kit's Dev Mode lets you run an LLM locally, configure it, and prompt/chat with it instantly.
* 🤗 **Universal:** ModelKits can be used with any AI, ML, or LLM project - even multi-modal models.
* 🐳 **Deploy containers:** Generate a Docker container as part of your `kit unpack` (coming soon).
* 🚢 **Kubernetes-ready:** Generate a Kubernetes / KServe deployment config as part of your `kit unpack` (coming soon).

### See KitOps in Action

There's a video of KitOps in action on the [KitOps site](https://kitops.ml/).

### What is in the box?

**[ModelKit](https://kitops.ml/docs/modelkit/intro.html):** At the heart of KitOps is the ModelKit, an OCI-compliant packaging format for sharing all AI project artifacts: datasets, code, configurations, and models. By standardizing the way these components are packaged, versioned, and shared, ModelKits facilitate a more streamlined and collaborative development process that is compatible with any MLOps or DevOps tool.

**[Kitfile](https://kitops.ml/docs/kitfile/kf-overview.html):** A ModelKit is defined by a Kitfile - your AI/ML project's blueprint. It uses YAML to describe where to find each of the artifacts that will be packaged into the ModelKit. The Kitfile outlines what each part of the project is.

**[Kit CLI](https://kitops.ml/docs/cli/cli-reference.html):** The Kit CLI not only enables users to create, manage, run, and deploy ModelKits -- it lets you pull only the pieces you need. Just need the serialized model for deployment? Use `unpack --model`, or maybe you just want the training datasets? `unpack --datasets`.

## 🚀 Try KitOps in under 15 Minutes

1. [Install the CLI](https://kitops.ml/docs/cli/installation.html) for your platform.
2. Follow the [Getting Started](https://kitops.ml/docs/get-started.html) docs to learn to pack, unpack, and share a ModelKit.
3. Test drive one of our [ModelKit Quick Starts](https://jozu.ml/organization/jozu-quickstarts) that include everything thing you need to run your model including a codebase, dataset, documentation, and of course the model.

- [Meta LLama 3.1](https://jozu.ml/repository/jozu/llama3.1-8b)
- [Google Gemma](https://jozu.ml/repository/jozu/gemma-7b)
- [Microsoft Phi3](https://jozu.ml/repository/jozu/phi3)
- [Fine-tuning](https://jozu.ml/repository/jozu/fine-tuning)
- [Rag pipeline](https://jozu.ml/repository/jozu/rag-pipeline)
- [Object detection](https://jozu.ml/repository/jozu/yolo-v10)

For those who prefer to build from the source, follow [these steps](https://kitops.ml/docs/cli/installation.html#🛠️-install-from-source) to get the latest version from our repository.

## ✨ What's New? 😍

We've been busy and shipping quickly!

* 📝 Add any arbitrary metadata to your Kitfile - hyperparameters, test results, links to resources, etc...
* 🔼 Automatic chunking of uploads to registries
* 🪵 New log-levels and per-request trace logging
* 🏎️ Reduced ModelKit packing time by >97%
* 🎯 Use KitOps to do LLM fine-tuning or package a RAG pipeline
* 🤓 Read Kitfile from `stdin`
* 👩‍💻 Use KitOps Dev mode to run an LLM instantly (no GPUs or internet required)

You can see all the gory details in our [release changelogs](https://github.com/jozu-ai/kitops/releases).

## Need Help?

### Join KitOps community

For support, release updates, and general KitOps discussion, please join the [KitOps Discord](https://discord.gg/Tapeh8agYy). Follow [KitOps on X](https://twitter.com/Kit_Ops) for daily updates.

If you need help there are several ways to reach our community and [Maintainers](./MAINTAINERS.md) outlined in our [support doc](./SUPPORT.md)

### Reporting Issues and Suggesting Features

Your insights help KitOps evolve as an open standard for AI/ML. We *deeply value* the issues and feature requests we get from users in our community :sparkling_heart:. To contribute your thoughts,navigate to the **Issues** tab and hitting the **New Issue** green button. Our templates guide you in providing essential details to address your request effectively.

### Joining the KitOps Contributors

We ❤️ our KitOps community and contributors. To learn more about the many ways you can contribute (you don't need to be a coder) and how to get started see our [Contributor's Guide](./CONTRIBUTING.md). Please read our [Governance](./GOVERNANCE.md) and our [Code of Conduct](./CODE-OF-CONDUCT.md) before contributing.

### A Community Built on Respect

At KitOps, inclusivity, empathy, and responsibility are at our core. Please read our [Code of Conduct](./CODE-OF-CONDUCT.md) to understand the values guiding our community.

## Roadmap

We [share our roadmap openly](./ROADMAP.md) so anyone in the community can provide feedback and ideas. Let us know what you'd like to see by pinging us on Discord or creating an issue.

---
---
---

<a href="https://trackgit.com">Trackgit</a>
