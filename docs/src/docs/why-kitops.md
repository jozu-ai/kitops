# Why Use KitOps?

Because hand-offs between AI/ML and app/SRE teams are painful.
Because getting an LLM or other AI/ML model deployed safely to production is hard.
Because they can be easy, and painless.

Jupyter notebooks are great, but extracting the model, datasets, and metadata from one is tricky. Similarly, ML-specific experimentation tools like MLFlow or Weights & Biases are excellent at training, but they save everything in proprietary formats that are confusing for software engineers and SREs.

Worse yet, none of these AI/ML tools are compatible with the toolchains organizations have successfully used for years with their applications.

With leaders demanding teams "add AI/ML" to their protfolios, many have fallen into a "throw it over the wall and hope it works" process that adds risk, delay, and frustration to self-hosting models.

> **The goal of KitOps is to simplify the sharing of AI/ML models, datasets, code, and configuration so that they can be run anywhere.**

We're not trying to replace your other tools, we just want to improve the experience of packaging and sharing models.

We dreamt of a better solution we call a "ModelKit." ModelKits:
* Combine models, datasets, code and all the context teams need to integrate, test, or deploy:
  * Training code
  * Model code
  * Serialized model
  * Training, validation, and other datasets
  * Metadata
* Let teams reuse their existing container registries by packaging everything as an OCI-compliant artifact
* Support unpacking only a piece of the model package to your local machine (saving time and space)
* Remove tampering risks by using an immutable package
* Reduces risks by including the provenance of the model and datasets

Using the kit CLI, you no longer have to remember the repo with the code in it, the registry with the model in it, the storage URI with the datasets, etc...

The better way:
Use `kit pack` to package up your Jupyter notebook, serialized model, and datasets (based on a [Kitfile](./kitfile/structure.md)).

Then `kit push` it to any OCI-compliant registry, even a private one.

Most people won't need everything, so just `kit unpack` from the remote registry to get just the model, only the datasets, or just the notebook. Or, if you need everything then a `kit pull` will grab everything.

Check out our [quick start](./quick-start.md), see the power and flexibility of our [CLI commands](./cli/cli-reference.md), or learn more about packaging your AI/ML project with [ModelKits](./modelkit/intro.md).
