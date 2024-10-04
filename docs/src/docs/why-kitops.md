# Why Use KitOps?

## The Problem

There is no standard and versioned packaging system for AI/ML projects. Today each part of the project is kept somewhere different:
* Code in Jupyter notebooks or (if you're lucky) git repositories 
* Datasets in DvC or storage buckets like S3
* Configuration in Jupyter notebooks, feature stores, MLOps tools, or ...
* Pipeline definitions in proprietary tools

Jupyter notebooks are great, but extracting the model, datasets, and metadata from one is tricky. Similarly, ML-specific experimentation tools like MLFlow or Weights & Biases are excellent at training, but they save everything in proprietary formats that are confusing for software engineers and SREs.

When the only people using AI were data scientists this was annoying but workable. Now there are application teams trying to integrate model versions with their application, testing teams trying to validate models, and DevOps teams trying to deploy and maintain models in production.

Without unified packaging teams take on risk and give up speed:
* Which dataset version was used to train and validate this model version?
* When did the dataset change? Did that effect my test run?
* Where are the best configuration parameters for the model we're running in production?
* Where did the model come from? Can we trust the source?
* What changes did we make to the model?

...and if you have to rollback a model deployment in production...good luck. With leaders demanding teams "add AI/ML" to their portfolios, many have fallen into a "throw it over the wall and hope it works" process that adds risk, delay, and frustration to self-hosting models.

This problem is only getting worse and the stakes are rising each day as more and more teams start deploying models to production without proper operational safeguards.

## The Solution

> [!NOTE]
> The goal of KitOps is to be a library of versioned packages for your AI project, stored in an enterprise registry you already use.

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

Use `kit pack` to package up your Jupyter notebook, serialized model, and datasets (based on a [Kitfile](./kitfile/structure.md)).

Then `kit push` it to any OCI-compliant registry, even a private one.

Most people won't need everything, so just `kit unpack` from the remote registry to get just the model, only the datasets, or just the notebook. Or, if you need everything then a `kit pull` will grab everything.

Check out our [getting started doc](./get-started.md), see the power and flexibility of our [CLI commands](./cli/cli-reference.md), or learn more about packaging your AI/ML project with [ModelKits](./modelkit/intro.md).
