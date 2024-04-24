# KitOps Workflow

KitOps is the market's only open source, standards-based packaging and versioning system designed for AI/ML projects. You can read more about why KitOps is a powerful collaboration solution for AI/ML projects in our [Why KitOps](./why-kitops.md) page.

## How KitOps Gets Used 🛠️

Organizations have spent the last 15+ years building up a suite of developer and DevOps tools to make developing software applications and services faster, safer, and easier. As they start working with AI/ML, they are realizing that those tools are not well suited to a world of massive models and even bigger datasets. ModelKits were created by a team who lived that nightmare and wanted others to avoid it.

Here’s how we used ModelKits along with the tools we already used in our data science, software development, and infrastructure management teams.

Here's how we implemented those tags and changes in our repo throughout the AI project development lifecycle.

### Data Access 💾

Each time our team has built AI/ML models for use in our products they needed to be trained on customer data that we couldn't let out of our private compute environment. We also had to be careful about who could even access that data at all. To help with this we would create two ModelKit repos in our registry for each project:

* One for the data science team with all the datasets they needed ("project-X-DS")
* One for the software development and infrastructure teams with only "safe" datasets that could be used for testing and in-production validation ("project-X")

When the data science team had finished training the model, they would create a new ModelKit that included only safe data and push that to the non-data science repo.

This guaranteed that anyone could quickly see exactly what datasets were used for each model and project for auditing, security, or compliance reasons. This is often challenging without a mechanism like ModelKits because at most organizations, datasets aren’t housed in a versioned repository and can easily (and silently) change during the course of model development.

### Development 🛒

Our data scientists developed models in Jupyter notebooks, but saved their work as a ModelKit at each development milestone to facilitate collaboration across our teams. This was simple by using the `kit pack` and `kit push` commands in their notebook.

For example, if the data scientist had completed fine-tuning of the model and our enterprise registry was GitLab (`registry.gitlab.com`), our repository was “chatbot” and the model was called “legalchat” then the commands would have been:

``` shell
kit pack . -t registry.gitlab.com/chatbot/legalchat:tuned
kit push registry.gitlab.com/chatbot/legalchat:tuned
```

Now everyone knew that the Legal Chat model had completed fine-tuning and they could quickly pull, run, or inspect the updated model, datasets, or code.

### Testing 🧑‍🍳

Once the ModelKit was pushed to the shared repo, a QA team would validate the model. This team would develop their own datasets, often built from what had been observed in production. The goal was to create a minimal dataset that would efficiently test the core and outer edges of how the model would be used by real customers.

This team generally didn't need the codebases in the ModelKit so they would just unpack the model and test it with their own dataset.

`kit unpack registry.gitlab.com/chatbot/legalchat:tuned --model`

After running their tests, they would add their dataset and results to the ModelKit so anyone could look over their results and understand any differences in accuracy from what the data science team had observed. They would tag this ModelKit as `tested` to alert the app team that they should do their integration tests next.

```shell
kit pack . -t registry.gitlab.com/chatbot/legalchat:tested
kit push registry.gitlab.com/chatbot/legalchat:tested
```

By doing this the project team maintained a ModelKit that included all the relevant assets that they might need: the model and codebases, any safe datasets used by the data science team, the testing dataset used by QA, and the results of both the data science team's and the QA team's tests.

This context was critical so that each team member could understand what work had already been done, look over the results, and potentially spot items that may have been missed.

### Integration 🍽️

Once a new ModelKit completed testing, our app team would pull the updated ModelKit so they could build the integrations with other in-production services and applications:

```shell
kit pull registry.gitlab.com/chatbot/legalchat:tested
```

The team would test their service integration with the updated model and could use any of the datasets included in the ModelKit for their own testing. This made it faster for them to execute integration work and kept them focused on the performance characteristics rather than worrying about dataset integrity or access.

Once integration was complete, they would add to the ModelKit the codebases for any applications that needed updates to work with the new model. Finally they would tag the ModelKit `integrated` to indicate that the model was ready for production deployment.

### Deployment 🚢

Adding the `integrated` tag to a ModelKit was a trigger for the DevOps team to get to work. They would take the serialized model from the ModelKit and ready it for deployment via our pipelines. Sometimes this meant putting the model into a container, but it often meant using an init container, a sidecar, entrypoint, or post-start hooks.

Once the project was packaged and ready for production they would add any necessary infrastructure-as-code (IAC) codebases to the ModelKit and retag it as `challenger` and begin deployment via our existing pipelines.

If the `challenger` model was replacing a model already in production (the `champion`), then once deployment and production validation were complete for the `challenger`, the DevOps team would retag the `champion` model to `rollback`:

```shell
kit tag registry.gitlab.com/chatbot/legalchat:champion registry.gitlab.com/chatbot/legalchat:rollback

kit push registry.gitlab.com/chatbot/legalchat:rollback
```

This way if there was a production incident we could always confidently roll back to a previous model revision without having to worry if the state of the model, or associated codebases had changed since it was last used successfully in production.

At this point the ModelKit for the `challenger` would be retagged to `champion`. 

```shell
kit tag registry.gitlab.com/chatbot/legalchat:challenger registry.gitlab.com/chatbot/legalchat:champion

kit push registry.gitlab.com/chatbot/legalchat:champion
```

These changes were part of our deployment automation and ensured that we always knew what was in production, what we could roll back to if needed, and what was coming up from development that would need production deployment in the near future.

### Using Tags 🏷️

As you can see, we used ModelKit tag names to give team members a quick way to determine where an AI project was in its lifecycle and whether they needed to get involved. Since ModelKits are immutable but use content-addressable storage, two ModelKits with the same contents but different tags only result in one ModelKit in storage (with two tag "pointers").

These were the tags we used, but each organization should tailor tags to suit their own lifecycle and processes.

#### Dataset-only Tags

* **_tuning_**: dataset is designed for model tuning
  * Used for ModelKits that contained only datasets. These were used almost exclusively by the data science team.
* **_validating_**: dataset is designed for model validation
  * Used for ModelKits that contained only datasets. These were used almost exclusively by the data science team.
* **_retired_**: the dataset is no longer appropriate for use, but needs to be retained for compliance or auditing reasons.
  * We kept old ModelKits so we could trace the history of development. When they were too old to be informative we used `kit remove` to get rid of them.

#### AI Project Tags

* **_trained_**: model has completed training phase
  * ModelKit would include the model, plus any training and validation datasets. Other assets were optional depending on the project.
* **_tuned_**: model has completed fine-tuning phase
  * ModelKit would include the model, plus tuning datasets. Other assets were optional depending on the project.
* **_tested_**: model has been tested successfully in QA
  * ModelKit would include all the assets passed into this phase, plus the validation dataset and any code that QA used.
* **_integrated_**: model has completed production application / service integration phase
  * ModelKit would include all the assets passed into this phase, plus codebases for any service that needed updating to work with the new model.
* **_challenger_**: model should be prepared to replace the current champion model
  * ModelKit would include any codebases and datasets that need to be deployed to production with the model. The idea was to create a package that could be deployed via our existing CI/CD pipelines (one for code, one for serialized models, one for datasets).
* **_champion_**: model is deployed in production
  * This is the same ModelKit as the `challenger` tagged one, but with an updated tag to show what's in production. This just creates a second reference to the same ModelKit, not 2x the ModelKit storage.
* **_rollback_**: the model was the previous "champion" model, before it was replaced by the "challenger"
  * This is the same ModelKit as the `champion` tagged one. We kept it so that we could quickly deploy it to production if there was a catastrophic failure with the current `champion` model.
* **_retired_**: the model is no longer appropriate for production usage
  * We kept old ModelKits so we could trace the history of development. When they were too old to be informative we used `kit remove` to get rid of them.

If you use KitOps differently and want to share it go ahead - we'd love to hear about it!