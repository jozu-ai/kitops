# How KitOps is Different

When people first come across KitOps they sometimes wonder, "how is this better than my favorite MLOps tool, git, Docker, etc...?"

## How AI Project Assets are Managed Today

Most teams working on AI projects store, track, and version their assets in one of two ways.

1. Using an MLOps tool
1. Using a combination of git, containers, cloud storage, and Jupyter notebooks

Neither solution is well suited to tracking and sharing AI project updates across data science, application development, infrastructure, and management teams... and neither is able to work seamlessly with the security, compliance, and efficiency processes organizations have spent decades perfecting.

Let's look at each option in a little more depth.

## KitOps vs. MLOps Tools

First off, it's important to understand that KitOps and its ModelKits don't completely replace the need for MLOps training and experimentation tools like Weights & Biases, MLFlow, or others.

However, [ModelKits](./modelkit/intro.md) are a better way to package, version, and share AI project assets outside of the data science team who use MLOps tools everyday.

Unlike MLOps tools, KitOps:

* Fits naturally (and without any changes) into organizations' existing deployment, security, and compliance processes
* Can already be used with *every* software, DevOps, and data science tool
* Uses existing, proven, and compliant registries organizations already depend on for their critical software assets
* Is simple enough for anyone to use, not just data science teams
* Leverages the same structure and syntax engineering teams are familiar with from containers and Kubernetes
* Is based on standards like OCI, that are vendor agnostic
* Is open source, and openly governed so it protects users and organizations from vendor lock-in
* Built by a community with decades of production operations and compliance experience

## KitOps vs. Jupyter, Containers, Dataset Storage, and Git

As with MLOps tools, Kit isn't designed to replace the other tools you already use and love. Jupyter notebooks, git, and containers all have their strengths and can be used for housing part of an AI project.

However, [ModelKits](./modelkit/intro.md) are a better way to package, version, and share *all* the assets for an AI project in one trackable place, for use by the data science team, software engineering, and infrastructure teams. This includes:

* Codebases
* Serialized models
* Datasets
* Feature lists
* Hyperparameters
* Prompts
* Deployment artifacts or IaC
* Etc...

The first and most important part of enabling fast and efficient inter-team collaboration on an AI project is to start with a central, accessible, understandable, and standards-based package for all these AI assets.

Let's look at some of the places parts of the project are likely housed today.

### Jupyter Notebooks

Jupyter notebooks are a fixture in the data science community for good reason - they're an excellent tool for mixing code, text, graphs, and images in a single file as is often required for research and model development.

Unfortunately, notebooks are not good at handling state or versioning. Although notebooks can be added to git, they're awkward to work with and resist the normal granular diff'ing that has made git so popular.

Additionally, while you can run a model in a notebook, the model isn't durably serialized making it painful to share models with colleagues who don't use notebooks.

#### ModelKits & Jupyter Notebooks

We suggest continuing to use notebooks, but include a [Kitfile](./kitfile/kf-overview.md) in each, and at the end of the notebook execute a `kit pack` command to save the [serialized model, dataset, and code](./cli/cli-reference.html#kit-pack) from the notebook into a ModelKit for versioning, centralized tracking, and sharing. This allows the data science team to continue to use Jupyter Notebooks, but allows software engineering, product management, and infrastructure teams to easily run, track, and use the models outside of notebooks.

### Containers

We love containers - they're great for running and deploying models. But they're not a natural way to distribute or version code or datasets. It's common for data scientists to have limited experience writing dockerfiles or building containers so most organizations have a platform engineering or SRE team build the container for the model.

#### ModelKits & Containers

We suggest having data science and production operations teams discuss how a model will be deployed (in a container, as a side-car, etc...) early in the development cycle. If a container is going to be used, you can [include a dockerfile or container in a ModelKit](./kitfile/kf-overview.html#kitfile-structure) so that a model can be easily run locally or pushed through a standard deployment pipeline when needed.

### Git

Git is excellent at managing software projects that consist of a large number of small files. Unfortunately, git was _never designed to manage large binary objects_ like the serialized models and datasets that are critical to AI/ML projects. Although git LFS can be used, it doesn't work smoothly with git's own versioning workflows. Plus, many data scientists don't know or like working with git, increasing the likelihood of repo errors and friction between teams.

#### ModelKits & Git

Code related to model development is often easier to store in ModelKits where it is always in-sync with the Jupyter Notebook, serialized model, and datasets used during development. Larger codebases and code related to application integrations is best kept in git, but is often helpful to also package into the ModelKit ([a codebase can be stored in a ModelKit](./kitfile/format.html#code) so that anyone can see the state of the code at the point that the project's ModelKit was versioned).

### Dataset Storage

Datasets are one of the things that are most difficult to version and store because there are so many types (SQL databases, CSVs, vector databases, images, audio/video files, etc...) and they're usually spread across many different storage locations (cloud storage, BI tools, databases, laptops, etc...).

It's easy to end up with near-duplicate datasets in different locations, making it extremely hard to know what dataset in what state was used to train a specific model, for example. Imagine if a dataset is discovered to include sensitive data - which models were trained with it? It's important to protect against these risks from the start.

#### ModelKits & Dataset Storage

Keeping datasets in versioned ModelKits ensures that it's always clear which data and state, was used with a specific version of the model. It avoids the risk of accidental data contamination and ensures you can always trace the [model/data lineage](./modelkit/spec.md#terminology-and-structure). A library of ModelKits for an AI project acts as a kind of audit record, allowing you to diff package contents to see when something changed and who made the change.