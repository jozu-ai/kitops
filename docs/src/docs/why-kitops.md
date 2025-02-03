# Why Use KitOps?

KitOps is the market's only open source, standards-based packaging and versioning system designed for AI/ML projects. Using the OCI standard allows KitOps to be painlessly adopted by any organization using containers and enterprise registries today (see a partial list of [compatible tools](./modelkit/compatibility.md)).

KitOps has been downloaded over 20,000 times in just the last three months. Users often use it as a:

* [Secure and immutable packaging and versioning standard](../modelkit/intro/) that is [compatible with their existing container registry](../modelkit/compatibility/#compliant-oci-registries)
* Point-of-control between development and production to [enforce consistency in packaging and documentation](../kitfile/kf-overview/)
* Catalogue of meaningful AI/ML project versions for regulatory compliance or change tracking
* Mechanism to simplify and unify the [creation of containers or Kubernetes deployment YAML](../deploy/)

> [!NOTE]
> The goal of KitOps is to be a library of versioned packages for your AI project, stored in an enterprise registry you already use.

## The Problem

There is no standard and versioned packaging system for AI/ML projects. Today each part of the project is kept somewhere different:
* Code in Jupyter notebooks or (if you're lucky) git repositories
* Datasets in DvC or storage buckets like S3
* Configuration in Jupyter notebooks, feature stores, MLOps tools, or ...
* Pipeline definitions in proprietary tools

This makes it difficult to track which versions of code, model, and datasets go together. It makes building containers harder and managing in-production AI/ML projects riskier.

Teams that use ModelKits report saving between 12 and 100 hours per AI/ML project iteration. While security and compliance teams appreciate that all AI/ML project assets are packaged together for each version and stored in an already secured and auditable enterprise container registry.

Suddenly tough questions like these are easy to answer:

* Where did the model come from? Can we trust the source?
* When did the dataset change? Which models were trained on it?
* Who build and signed off on the model?
* Which model is in production, which is coming, and which has been retired?

## The Solution

Kit's ModelKits are the better solution:
* Combine models, datasets, code and all the context teams need to integrate, test, or deploy:
  * Training code
  * Model code
  * Serialized model
  * Training, validation, and other datasets
  * Metadata
* Let teams reuse their existing container registries by packaging everything as an OCI-compliant artifact
* Support unpacking only a piece of the model package to your local machine (saving time and space)
* Remove tampering risks by using an immutable package
* Reduce risks by including the provenance of the model and datasets

Use `kit pack` to package up your Jupyter notebook, serialized model, and datasets (based on a [Kitfile](../kitfile/kf-overview/)).

Then `kit push` it to any OCI-compliant registry, even a private one.

Most people won't need everything, so just `kit unpack` only the layers you need (e.g., only model and datasets, or only code and docs) from the remote registry. Or, if you need everything then a `kit pull` will grab everything.

Finally [package it all up as a container or Kubernetes deployment](../deploy/).

Check out our [getting started doc](../get-started/), see the power and flexibility of our [CLI commands](../cli/cli-reference/), or learn more about packaging your AI/ML project with [ModelKits](../modelkit/intro/) and even making them [deployable](../deploy/).
