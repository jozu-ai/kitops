# Use Cases

KitOps removes friction from the handoffs between AI/ML and App/SRE teams:
* AI/ML teams can package their project, including models, code, datasets, and configuration, into a ModelKit that can be run by anyone, anywhere.
* Application teams can pull only what they need from a ModelKit to embed a model in their app, or generate a RESTful API for integration.
* SRE teams can pull the serialized model and validation dataset from a ModelKit and deploy it through their CI/CD pipelines.
* Governance teams can track and audit AI/ML assets with your existing processes and tools (we save everything as OCI-compatible artifacts).
* Security teams can integrate ModelKits into their SBOM processes.

To make the use cases a little more fun let's tell two stories:
1. Collaborating on an AI-enabled application for customers
2. Deploying a model for internal business usage

In both cases the context of the stories is the same.

Weyland-Yutani Corporation has been adding AI and ML models to its portfolio for both internal- and customer-facing deployment. There are three main groups that are involved in the development lifecycle for their AI/ML work:
* AI/ML Team, composed of data scientists, data engineers, and MLOps engineers. This specialized group builds, tunes, and validates the models that increasingly touch every aspect of the organization's work.
* Application Teams, composed of software and hardware engineers. This group builds the applications that the business and customers interact with - over time more and more of these applications are being integrated into the models the AI/ML team works on.
* SRE Team, composed of infrastructure, pipeline, and release specialists who make sure the models and applications are deployed quickly and safely. They also monitor all production applications to ensure security, correctness, availability, and performance.

There has been a heated argument within Weyland-Yutani about how best to manage the collaboration, testing, and deployment nuances of AI/ML-integrated applications. Today, this is a friction-filled process that relies on significant manual effort from every team. Each team is often having to ask for help from other teams to prepare models for the work each needs to do, slowing down progress and introducing needless human errors.

...but all that's about to change! (Hint: the change involves KitOps...surprise!)

## 1/ Collaborating on an AI-Enabled Application for Production

### AI/ML Team Collaboration

Rajat is working on tuning an open source foundational model for his company. After several days of work in his Jupyter notebook, he has a model that is outperforming the old model they've been using in production. But before he alerts the App and SRE teams Rajat wants to have another data scientist try the model and verify his findings.

Rajat adds two lines to the end of his Jupyter notebook to create a ModelKit for that will include the model, the datasets used for training and validation, the Jupyter notebook file with the code and context, and the Kitfile manifest and metadata:

```sh
$ ./kit pack -t corp-registry/app-model:challenger
$ ./kit push  --http corp-registry/app-model:challenger
```

Now that the ModelKit is published on Weyland-Yutani's private corporate registry, others can quickly get what they need from the project.

Rajat sends a note in Slack to his colleague Gorkem who also uses Juypter notebooks. He only needs the notebook file so he quickly runs:

```sh
$ ./kit pull -filter code corp-registry/app-model:challenger
```

Gorkem loads Rajat's notebook file and runs through Rajat's tests. He Slacks him back congratulating him on a great model and the painless way he was able to share his work.

Thanks KitOps!

### AI/ML and App Team Collaboration

Now that [Rajat and Gorkem have both agreed](#ai-ml-team-collaboration) that the newly trained model is a real challenger for production, they alert the Application Team who will need to do an integration test between their app and the model.

Rajat Slacks the application team a heads-up, and Nida volunteers to kick off integration testing. She needs the model and hopes that Rajat included a validation dataset. Nida uses the Kit CLI's flexible pull command to only extract the new model's Kitfile and check whether there's a validation dataset included:

```sh
$ ./kit pull -filter config corp-registry/app-model:challenger
```

After confirming that there _is_ a validation dataset, she pulls only the validation dataset from the ModelKit, saving her time and skipping the need to learn Rajat's repository and file structure:

```sh
$ ./kit pull -filter dataset:validation -filter model corp-registry/app-model:challenger
```

The model arrives as a `.tar` which she can drop directly into her integration test pipeline along with the validation dataset. After 13 minutes the tests complete and show that the application functions correctly with the new model except for one API call. Nida quickly corrects the API bug and rebuilds the app, running a second integration test with the new model. This time everything passes.

Nida is a bit of an introvert so the fact that she didn't need to ask anyone for help understanding where to find the assets she needed, and didn't need to repackage the serialized model contributed to her awesome afternoon (that and the fact that she brought her dog to work today).

Thanks KitOps!

### AI/ML, App, and SRE Team Collaboration

When the second [integration pipeline that Nida kicked off](#ai-ml-and-app-team-collaboration) completed successfully it automatically notified Annika in the SRE team. Annika quickly looks over the integration test results and sees that the model and application are ready for deployment. She issues a PR to Weyland-Yutani's GitOps repository to kick off the production deployment pipeline in GitLab.

Part of the GitOps automation calls the Kit CLI to extract only the model (saving a heap of time by not having to pull a +10GB dataset that isn't needed):

```sh
$ ./kit pull -filter model corp-registry/app-model:challenger
```

The deployment pipeline does an A/B deployment of both the model and application to their Kubernetes environment in one of their smaller regions (using Kserve). Annika is quickly able to confirm that the Rajat's Challenger model does work more efficiently and accurately than the in-production Champion model. Over the next hour she progressively shifts customer traffic from Champion to Challenger and, once all traffic is going to the Champion, she undeploys Champion. At that point she triggers parallel deployments to the other regions and goes through the same A/B analysis in each to make sure that the behaviour is consistent across the galaxy where Weyland-Yutani's customers are found.

Once the new Challenger model is taking all customer traffic in every galactic region, Annika retags the old Champion model with a version number, and tags the Challenger model as Champion:

```sh
$ ./kit tag 1.0 corp-registry/app-model:champion
$ ./kit tag champion corp-registry/app-model:challenger
```

Once the deployment is complete across all galactic regions, Annika sends an update to the AI/ML, app, and executive teams. Weyland-Yutani's leadership is ecstatic because the new model will give their spaceships a better way to scan for alien lifeforms when they receive random deep space distress signals. Everyone involved gets cake, champagne, and a high-five from Ripley herself...

Thanks KitOps!

## 2/ Collaborating on an Internal Model

### AI/ML and SRE Team Collaboration

Angel leads an internal data science team that builds models to help Weyland-Yutani's executives make faster and better decisions with the help of AI/ML models. Their team was asked in the most recent monthly business review to identify all the customers they have who might be a churn risk based on the set of common traits and behaviours seen in all churned customers over the last five years.

After a few days work, Angel's team is done. They build a ModelKit so their new model can easily be deployed:

```sh
$ ./kit pack -t corp-registry/churn-model:beta
$ ./kit push --plain-http corp-registry/churn-model:beta
```

Angel is experienced with Weyland-Yutani's GitOps process so they issue a PR with a reference to the ModelKit in the company's private repository. Annika reviews the PR and approves and merges it once all the automated tests pass.

Part of the GitOps automation that is triggered by the merged PR calls the Kit CLI to extract only the model:

```sh
$ ./kit pull -filter model corp-registry/churn-model:beta
```

After running for several hours against the company's internal sales, marketing, and product analytics data, the model spits out a list of at-risk customers ranked in order of most likely to churn. Angel's team, after sending the list to executive leadership, are given champagne and bonuses.

Thanks KitOps!



