# How KitOps Is Used 🛠️

KitOps is the market's only open source, standards-based packaging and versioning system designed for AI/ML projects. Using the OCI standard allows KitOps to be painlessly adopted by any organization using containers and enterprise registries today (see a partial list of [compatible tools](./modelkit/compatibility.md)).

Today AI/ML development in enterprises relies on artifacts that are tightly coupled, but versioned and stored separately:
* Models in Jupyter notebooks or MLOps tools
* Datasets in data lakes, databases, or files systems
* Code in git repositories
* Metadata (hyperparameters, features, weights, etc...) in various locations based on their type

For this reason, organizations around the world are using KitOps as a "gate" between development and production. Those who are concerned about end-to-end auditing of their model development (like those in regulated industries, or under the jurisdiction of the [EU AI Act](https://artificialintelligenceact.eu/) extend KitOps usage to security and development use cases (see [Level 2](#level-2-adding-security-️) and [Level 3](#level-3-storage-for-all-ai-project-versions-) use cases below.)

## Level 1: Handoff From Development to Production 🤝

Organizations are having AI teams build a [ModelKit](./modelkit/intro.md) for each version of the AI project that is going to staging, user acceptance testing (UAT), or production. KitOps is ideally suited to CI/CD pipelines (e.g., using [KitOps in a GitHub Action](https://dev.to/kitops/introducing-the-new-github-action-for-using-kit-cli-on-mlops-pipelines-21ia)) either triggered manually by the model development team when they're ready to send the model to production, or automatically when a model or its artifacts are updated in their respective repositories.

This ensures that:
* __Operations teams have all the assets and information they need__ in order to determine how to test, deploy, audit, and manage these new workloads
* __AI versioned packages are held in the same enterprise registry__ as other production assets like containers making them easier to find, secure, and audit
* __Compliance teams have a catalogue of versioned models__ that can be used for [EU AI Act](https://artificialintelligenceact.eu/) or other regulatory reporting
* __Everyone has a library of immutable and signed ModelKits__ for intellectual property, progress tracking, or other requirements
* __Organizations are protected against vendor shifts__ in their MLOps and Serving Infrastructure domains (this also gives them negotiating leverage with vendors)

Teams working on model development continue to use their disparate repositories during the development cycle at this stage. This is where most organizations start their usage of KitOps, but once they start most continue on...

## Level 2: Adding Security 🛡️

Some organizations want to scan their models either before they enter the development phase (ideal), or before they are promoted beyond development. The open source [ModelScan project](https://github.com/protectai/modelscan) can help here.

### Creating a Curated Model Set 🧑‍🍳

Security-conscious organizations will often restrict the set of models that data science teams can use as a basis for their work. In these cases specific models can be pulled from public repositories like Hugging Face, then scanned with a tool like ModelScan, and finally packaged as a KitOps ModelKit stored in their enterprise registry.

This guarantees that models used by internal teams are safe and tamper-proof. By storing the models in the existing enterprise registry they're also easy for anyone to find or audit.

### Protecting Production 🚦

After model development has completed, the resulting ModelKit and its artifacts can again by scanned by something like ModelScan and only allowed to move forward if it passes. Using ModelKits here again ensures that a model that passes the scan is packaged and signed so that it cannot be tampered with on its way to production.

Even with this level of scrutiny, however, there remain some risks since the varying repositories and versioning of artifacts during development can invite accidental or malicious tampering. This leads us to Level 3 adoption...

## Level 3: Storage for all AI Project Versions 💾

Organizations that are more mature in their handling of AI projects, or are subject to extra scrutiny or regulations extend the use of ModelKits to the development phase as well. This solves the currently scattered storage of artifacts.

This ensures that:
* The organization uses a standards-based storage (OCI 1.1 artifacts) for all their AI/ML project work
* Data and AI teams who don't work closely together can share artifacts even if they're using different development tools
* Development artifacts can't be accidentally or maliciously tampered with

This can be a lightweight process automated with ModelKit packing via the Kit [CLI](./cli/cli-reference.md).

The beauty of KitOps' ModelKits is their flexibility and the fact that they fit into the already standard tools and processes that organizations have built around OCI artifacts and registries.

Below are details on how we have used KitOps internally from development to production.

### Using Tags 🏷️

We used ModelKit tag names to give team members a quick way to determine where an AI project was in its lifecycle and whether they needed to get involved. Since ModelKits are immutable but use content-addressable storage, two ModelKits with the same contents but different tags only result in one ModelKit in storage (with two tag "pointers").

* **_tuning_**: dataset is designed for model tuning
  * Used for ModelKits that contained only datasets. These were used almost exclusively by the data science team.
* **_validating_**: dataset is designed for model validation [1]
  * Used for ModelKits that contained only datasets. These were used almost exclusively by the data science team.
* **_trained_**: model has completed training phase [2]
  * ModelKit would include the model, plus training and validation datasets. Other assets were optional depending on the project.
* **_tuned_**: model has completed fine-tuning phase [2]
  * ModelKit would include the model, plus training and validation datasets. Other assets were optional depending on the project.
* **_challenger_**: model should be prepared to replace the current champion model
  * ModelKit would include any codebases and datasets that need to be deployed to production with the model. The idea was to create a package that could be deployed via our existing CI/CD pipelines (one for code, one for serialized models, one for datasets).
* **_champion_**: model is deployed in production
  * This is the same ModelKit as the `challenger` tagged one, but with an updated tag to show what's in production. This just creates a second reference to the same ModelKit, not 2x the ModelKit storage.
* **_rollback_**: the model was the previous "champion" model, before it was replaced by the "challenger"
  * This is the same ModelKit as the `champion` tagged one. We kept it so that we could quickly deploy it to production if there was a catastrophic failure with the current `champion` model.
* **_retired_**: the model is no longer appropriate for production usage
  * We kept old ModelKits so we could trace the history of development. When they were too old to be informative we used `kit remove` to get rid of them.

Here's how we implemented those tags and changes in our repo throughout the AI project development lifecycle.

### Development 🧪

Our data scientists developed models in Jupyter notebooks, but saved their work as a ModelKit at each development milestone to facilitate collaboration across our teams. This was simple by using the `kit pack` and `kit push` commands in their notebook:

For example, if the data scientist had completed fine-tuning of the model and our enterprise registry was GitLab (`registry.gitlab.com`), our repository was “chatbot” and the model was called “legalchat” then the commands would be:

``` shell
kit pack . -t registry.gitlab.com/chatbot/legalchat:tuned
kit push registry.gitlab.com/chatbot/legalchat:tuned
```

Now everyone knew that the Legal Chat model had completed fine-tuning and they could quickly pull, run, or inspect the updated model, datasets, or code.

### Integration 🧑‍🍳

Once a new ModelKit completed training or tuning, our app team would pull the updated ModelKit:

`kit pull registry.gitlab.com/chatbot/legalchat:tuned`

The team would test their service integration with the updated model, paying special attention to the performance characteristics and watching for features that may have changed between model versions.

### Testing 🍽️

Once the application team had confirmed the new model didn’t change the behavior of the application, an engineer from outside the data science team would validate the model. To do this they would run the new model with the validation dataset included in the ModelKit, and compare their results with the results from the data science team. This engineer generally didn't need the codebases in the ModelKit so they would just unpack the model and datasets.

`kit unpack registry.gitlab.com/chatbot/legalchat:tuned --model --datasets`

This guaranteed that they were using the same dataset as the data science team to replicate the results, which is otherwise difficult because at most organizations, datasets aren’t housed in a versioned repository and can easily (and silently) change during the course of model development.

If the model passed validation the QA team would retag it from `tuned` to `challenger` to alert the SRE team that it should be readied for deployment.

```shell
kit tag registry.gitlab.com/chatbot/legalchat:tuned registry.gitlab.com/chatbot/legalchat:challenger

kit push registry.gitlab.com/chatbot/legalchat:challenger
```

### Deployment 🚢

Adding the `challenger` tag to a ModelKit was a trigger in our team for the DevOps group to know that it was ready to deploy to production. They would take the serialized model from the ModelKit and ready it for deployment via our pipelines. This may mean putting the model into a container, but it may mean using an init container, a sidecar, entrypoint, or post-start hooks.

Once the model was deployed and validated in production the ModelKit for the model that was previously tagged `champion` would be retagged to `rollback`, and the ModelKit for the `challenger` would be retagged to `champion`. These changes were part of our deployment automation and ensured that we were always ready to quickly redeploy the `rollback` model in case something catastrophic happened in production withe current `champion`.

```shell
kit tag registry.gitlab.com/chatbot/legalchat:champion registry.gitlab.com/chatbot/legalchat:rollback

kit push registry.gitlab.com/chatbot/legalchat:rollback

kit tag registry.gitlab.com/chatbot/legalchat:challenger registry.gitlab.com/chatbot/legalchat:champion

kit push registry.gitlab.com/chatbot/legalchat:champion
```

If you use KitOps differently and want to share it go ahead - we'd love to hear about it!