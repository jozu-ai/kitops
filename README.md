# KitOps

## Description

The KitOps project provides tools for building, packaging, and deploying / running all kinds of ML models. The KitOps project enables data scientists and AI/ML developers to encapsulate their models along with all necessary dependencies, configurations, and environments into a standardized, portable format. This facilitates easy sharing and deployment across diverse computing environments, ensuring that models are readily usable without compatibility issues.

The project is composed of two main components:

* A `kitfile` manifest that lists a model, its dataset(s), code, and metadata - everything required to build, deploy, or run that model
* The `kit` command-line interface (CLI) tool designed to simplify the packaging, versioning, and distribution of AI/ML models as OCI (Open Container Initiative) artifacts

We refer to the final package that the `kit` CLI builds based on the `kitfile` as a **Model Kit**. Model Kits can be run in a Jupyter Notebook, or packaged as a Docker container that can be run anywhere. Kit supports a wide range of AI/ML frameworks, offers robust version control for model iterations, integrates security practices for safe distribution, and provides a seamless connection to model repositories for storing and retrieving models.

Kit's aim is to streamline the end-to-end lifecycle of AI/ML model management, making it as effortless as managing containerized applications.

## Running Kit with Binaries

You can download the Kit CLI using one of our [tagged versions](https://github.com/jozu-ai/kitops/tags). Make sure you get the right binary for your platform:

* MacOS: TODO
* Linux: TODO
* Windows: TODO

We suggest renaming the executable once it's downloaded to just `kit` (make sure it's in your path an executable).

Run Kit by opening a terminal and typing:

```shell
./kit
```

This will list all the commands you can use.

## Building and Running Kit from Source Code

You can get the Kit CLI sources from our [tagged versions](https://github.com/jozu-ai/kitops/tags).

```shell
go build -o kit
```

Then run the project:

```shell
./kit
```

Alternatively 

```shell
go run kit
```

## Filing Issues and Feature Requests

We want to see Kit become an open standard across the growing AI/ML industry so we *deeply value* the issues and feature requests we get from users in our community :sparkling_heart:. You can file an issue by selecting the **Issues** tab and hitting the **New Issue** green button. There are templates for Issues and Feature Requests, the more information you provide us the faster and better we can address your request. Thank you!

## Contributing

We ❤️ our Kit community and contributors. To learn more about the many ways you can contribute (you don't need to be a coder) and how to get started see our [Contributor's Guide](./CONTRIBUTING.md). Please read our [Governance](./GOVERNANCE.md) and our [Code of Conduct](./CODE-OF-CONDUCT.md) before contributing.

## Support

If you need help there are several ways to reach our community and [Maintainers](./MAINTAINERS.md) outlined in our [support doc](./SUPPORT.md)

## Code of Conduct

We run an inclusive, empathetic, and responsible community. Please read our [Code of Conduct](./CODE-OF-CONDUCT.md).
