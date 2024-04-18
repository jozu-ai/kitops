# How KitOps is Different

When people first come across KitOps they sometimes wonder, "how is this better than my favorite MLOps tool, git, Docker, etc...?"

## How AI Project Assets are Managed Today

Most teams working on AI projects store, track, and version their assets in one of two ways.

1. Using an MLOps tool
1. Using a combination of git, containers, block storage, and Jupyter notebooks

Neither solution is well suited to tracking and sharing AI project updates across data science, application development, infrastructure, and management teams... and neither is able to work seamlessly with the security, compliance, and efficiency processes organizations have spent decades perfecting.

Let's look at each option in a little more depth.

## KitOps vs. MLOps Tools

First off, it's important to understand that KitOps and its ModelKits don't completely replace the need for most MLOps tools like Weights & Biases, MLFlow, or others. However, ModelKits are a better way to package, version, and share AI project assets outside of the data science team who use MLOps tools everyday.

Unlike MLOps tools, KitOps:

* Fits naturally (and without any changes) into organizations' existing deployment, security, and compliance processes
* Can already be used with *every* software, DevOps, and data science tool
* Uses existing, proven, and compliant registries organizations already depend on for their critical software assets
* Is simple enough for anyone to use, not just data science teams
* Leverages the same structure and syntax engineering teams are used to from containers and Kubernetes
* Is based on standards like OCI, that are vendor agnostic
* Is open source, and openly governed so it protects users and organizations from vendor lock-in
* Built by a community with decades of production operations and compliance experience

When a model is "ready" from a data science perspective there is still a lot of work needed to make it ready for production usage. For example, even though many MLOps tools have a "deploy" button, no enterprise would allow someone to deploy directly from those tools to production for any critical applications because they bypass the years worth of security, compliance, and safety processes and tooling organizations rely on to users and their business safe.

Vendor lock-in should also be a concern for enterprises choosing their AI project tool chain. At this early stage there are thousands of MLOps tools, and many of those companies will fail or be acquired. Standards will evolve and some vendors will build a walled garden, trapping customers into increasingly costly contracts. KitOps is a free, open source project that is openly governed. We are working to create a vendor neutral packaging model that can work with any tool and vendor, so that if a vendor change is required at any point, you know you can easily move your projects without costly retooling or potential data loss.

## KitOps vs. Git, Containers, and Jupyter

As with MLOps tools, Kit isn't designed to replace the other tools you already use and love. Jupyter notebooks, git, and containers all have their strengths and can be used for housing part of an AI project. However, ModelKits are a better way to package, version, and share *all* the assets for an AI project in one trackable place, for use by the data science team, software engineering, and infrastructure teams. This includes:

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

Unfortunately, Jupyter is not good at handling state or versioning. Although the notebooks can be added to git, they're awkward to work with and resist the normal granular diff'ing that has made git so popular.

Additionally, although you can run a model in a notebook, the model isn't durably serialized making it painful to share models with colleagues who don't use notebooks.

We suggest continuing to use notebooks, but include a [Kitfile](./kitfile/kf-overview.md) in each, and at the end of the notebook execute a `kit pack` command to save a serialized model, the dataset, and code from the notebook into a ModelKit for versioning, centralized tracking, and sharing.

### Git

Git is excellent at managing software projects which typically consist of a large number of small files. Unfortunately, git was __never designed to manage large binary objects__ like serialized models that are often >10GB, or datasets that can often exceed 100GB. Although git LFS can be used, it is awkward and doesn't even work smoothly with git's own versioning workflows. Plus, many data scientists don't know or like working with git, increasing the likelihood of repo errors and friction between teams.

We suggest keeping code in git and using it as you do today. A codebase can be cloned into a ModelKit so that anyone can see the state of the code at the point that the project's ModelKit was versioned. The larger binary objects, however, should be kept in ModelKits where they can be more efficiently managed, versioned, and shared.

### Containers

We love containers - they're great for running and deploying models. But they're not a natural way to distribute or version code or datasets. You can include a dockerfile or container in a ModelKit so that a model can be easily pushed through a standard deployment pipeline, for instance.

### Data Storage

Datasets are one of the things that are most often problematic to version and store because there are so many types (SQL databases, CSVs, vector databases, images, audio/video files, etc...) and they're usually spread across many different places (cloud storage, BI tools, databases, laptops, etc...).

It's easy to end up with near-duplicate datasets in different locations, making it extremely hard to know what dataset in what state was used to train a specific model, for example. Imagine if a dataset is discovered to include sensitive data - which models were trained with it?

ModelKits simplify answering these and other questions by providing a clear lineage for every asset in the package and allowing you to diff package contents to see when things changed.

Regardless of where your datasets or other assets are housed, ModelKits make it simple to create a library of versions for each package that is vendor-neutral, shareable, verifiable, and safe.