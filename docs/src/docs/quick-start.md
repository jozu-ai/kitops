# KitOps Quick Start

In this guide, we'll use ModelKits and the kit CLI to easily:
* Package up a model, notebook, and datasets into a single ModelKit you can share through your existing tools
* Push the ModelKit package to a public or private registry
* Grab only the assets you need from the ModelKit for testing, integration, local running, or deployment

## Before we start...

1. Make sure you've got the [Kit CLI setup](./cli/installation.md) on your machine. If you'd like to learn more about ModelKits [check out the overview](./modelkit/intro.md). If you are wondering about specific flags you can check out everything in the [CLI command reference](./cli/cli-reference.md)

2. Create and navigate to a new folder on your filesystem. We'd suggest calling it `KitStart` but any name works.

One more note...with Kit, packing and pushing are separate steps (same with unpacking and pulling). Packing builds the ModelKit using the content outlined in the `Kitfile` manifest. Pushing takes a built ModelKit from your local registry and sends it to any remote registry. We strongly recommend tagging your ModelKits with a version number and any other tags that will help your team (e.g., challenger, champion, v1.3, dev, production, etc...)

## Learning to use the CLI

### 1/ Check your CLI Version

Check that the Kit CLI is properly installed.

```sh
kit version
```

You'll see information about the version of Kit you're running. If you get an error check to make sure you have [Kit installed](./cli/installation.md) and in your path.

### 2/ Login to Your Registry

You can log into any container registry, but below are a few popular ones.

**Docker Hub**

```sh
kit login docker.io
```

**GitHub Registry**

```sh
kit login ghcr.io
```

You'll see `Log in successful`.

<!--
1. Pull a sample ModelKit. You'll use this to learn the CLI before creating your own ModelKit.

```sh
kit pull ghcr.io/jozu-ai/modelkit-examples/scikitlearn-tabular:latest
```

You'll see a message `Pulled [SHA256 digest ID]`.

1. Check your repository.

```sh
kit list
```

You'll see `ghcr.io/jozu-ai/modelkit-examples/scikitlearn-tabular   latest   Jozu`

-->

### 3/ Get a Sample ModelKit

Let's unpack a sample ModelKit to our machine that we can play with. In this case we'll unpack the whole thing, but one of the great things about Kit is that you can also selectively unpack only the thigs you need: just the model, the model and dataset, the code, the configuration...whatever you want. Check out the `unpack` [command reference](./cli/cli-reference.md) for details.

You can grab [any of the ModelKits](https://github.com/orgs/jozu-ai/packages) from our site, but we've chosen a small language model example below.

```sh
kit unpack ghcr.io/jozu-ai/modelkit-examples/finetuning_slm:latest
```

You'll see a set of messages as Kit unpacked the configuration, code, datasets, and serialized model. Now list the directory contents:

```sh
ls
```

You'll see a single file (`Kitfile`) which is the manifest for our ModelKit, and a set of files or directories including adapters, a Jupyter notebook, and dataset.

### 4/ Check our Local Repository

Now let's check that we don't have anything in our local repository.

```sh
kit list
```

You'll see the column headings for an empty table with things like `REPOSITORY`, `TAG`, etc...

### 5/ Pack the ModelKit

Since our repository is empty we'll need to create our ModelKit. The ModelKit in your local registry will need to be named the same as your remote registry. So the command will look like:

`kit pack . -t [your registry address]/[your repository name]/mymodelkit:latest`

In my case I am pushing to my `jozubrad` repository on Docker Hub, so I use:

```sh
kit pack . -t docker.io/jozubrad/mymodelkit:latest
```

You'll see a set of `Saved ...` messages as each piece of the ModelKit is saved to the local repository.

Let's check our local registry now:

```sh
kit list
```

You'll see a new entry named the same as your pack command.

### 6/ Push the ModelKit to a Remote Repository

You're already logged in to your remote repository, so now you can just push. The naming of your ModelKit will need to be the same as what you see in your `kit list` command (REPOSITORY:TAG). You can even copy and paste it. In my case it looks like:

```sh
kit push docker.io/jozubrad/mymodelkit:latest
```

### Congratulations

You've learned how to unpack a ModelKit, pack one up, and push it. Anyone with access to your remote repository can now pull your new ModelKit and start playing with your model:

```sh
kit pull docker.io/jozubrad/mymodelkit:latest
```

Thanks for taking some time to play with Kit. We'd love to hear what you think. Feel free to drop us an [issue in our GitHub repository](https://github.com/jozu-ai/kitops/issues) or join [our Discord server](https://discord.gg/YyAfWnEg).


<!--


We recommend pulling one of our [example ModelKits](https://github.com/orgs/jozu-ai/packages) for this quick start. From there you can try [writing a Kitfile](./kitfile/format.md) for your own AI/ML project.

## Preparing for packaging

The first step is to make a `Kitfile` - a YAML manifest for your ModelKit. There are four main parts to a Kitfile:
1. ModelKit metadata in the `package` section
1. Path to the Jupyter notebook folder in the `code` section
1. Path to the datasets in the `datasets` section
1. Path to the serialized model in the `model` section

Here's an example Kitfile:

```yaml
manifestVersion: v1.0.0

package:
  authors:
  - Jozu
  description: Updated model to analyze flight trait and passenger satisfaction data
  license: Apache-2.0
  name: FlightSatML

code:
- description: Jupyter notebook with model training code in Python
  path: ./notebooks

model:
  description: Flight satisfaction and trait analysis model using Scikit-learn
  framework: Scikit-learn
  license: Apache-2.0
  name: joblib Model
  path: ./models/scikit_class_model_v2.joblib
  version: 1.0.0

datasets:
- description: Flight traits and traveller satisfaction training data (tabular)
  name: training data
  path: ./data/train.csv
- description: validation data (tabular)
  name: validation data
  path: ./data/test.csv
```

The only mandatory parts of the Kitfile are:
* `manifestVersion`
* At least one of `code`, `model`, `or datasets` sections

A ModelKit can only contain one model, but multiple datasets or code bases.

So a minimal ModelKit for distributing a pair of datasets might look like this:
```yaml
manifestVersion: v1.0.0

datasets:
- name: training data
  path: ./data/train.csv
- description: validation data (tabular)
  name: validation data
  path: ./data/test.csv
```

Right now you can only build ModelKits from files on your local system...but don't worry we're already working towards allowing you to reference remote files. For example, building a ModelKit from a local notebook and model, but a dataset hosted on DvC, S3, or anywhere else.

## Packing, tagging, and pushing

With Kit, packing and pushing are separate steps. Packing builds the ModelKit using the content outlined in the Kitfile manifest. Pushing takes a built ModelKit and sends it to an OCI-compliant registry. We strongly recommend tagging your ModelKits with a version number and any other tags that will help your team (e.g., challenger, champion, v1.3, dev, production, etc...)

For this example I want to tag my ModelKit with the "challenger" tag since I think this model can replace our current production model. Here I'm running the pack command from the directory where my Kitfile is located.

```sh
kit pack . -t mymodel:challenger
```

I can run `kit list` to see what's in my registry and confirm that `mymodel:challenger` is there. Then I can push it to my remote registry:

```sh
kit push ghcr.io/jozu-ai/modelkit-examples/mymodel:challenger
```

Once I've pushed the ModelKit to a shared registry I can let my team know.

## Unpacking and pulling

Like pack/push, unpack/pull are paired. Unpack allows you to pull only select parts of the ModelKit. This is important because ModelKits might be very large, especially if they contain large datasets. Unpacking gives you the flexibility to only pull what you need. The pull command will always pull down the whole ModelKit.

So if I only want to pull the serialized model to a `model` folder in my current directory:

```sh
kit unpack mymodel:challenger --model -d ./model
```

I can do the same with `--datasets` or `--code`. If I only want to get the Kitfile I can use `--config`.

We're also working on a `dev` command that will run a lightweight serving engine on a laptop and auto-generate a RESTful API to make it easier for software developers to test an in-development model or integrate it with their applications.

If you have questions please join our [Discord server](https://discord.gg/YyAfWnEg).
