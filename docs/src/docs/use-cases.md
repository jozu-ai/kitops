# KitOps Workflow

KitOps is the market's only open source, standards-based packaging and versioning system designed for AI/ML projects. You can read more about why KitOps is a powerful collaboration solution for AI/ML projects in our [Why KitOps](./why-kitops.md) page.

## How KitOps Gets Used 🛠️

Organizations have spent the last 15+ years building up a suite of developer and DevOps tools to make developing software applications and services faster, safer, and easier. As they start working with AI/ML, they are realizing that those tools are not well suited to a world of massive models and even bigger datasets. ModelKits were created by a team who lived that nightmare and wanted others to avoid it.

Here’s how we used ModelKits along with the tools we already used in our data science, software development, and infrastructure management teams.

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