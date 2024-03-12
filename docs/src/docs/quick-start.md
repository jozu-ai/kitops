# KitOps Quick Start

In this guide, we'll use ModelKits and the kit CLI to easily:
* Package up a model, notebook, and datasets into a single ModelKit you can share through your existing mechanisms
* Push the ModelKit package to a public or private registry
* Grab just the things you need from the ModelKit for testing, integration, local running, or deployment

## Before we start...

Make sure you've got the Kit CLI setup on your machine. Our [installation instructions](./cli/installation.md) will help.

We recommend starting by pulling one of our [example ModelKits](https://github.com/orgs/jozu-ai/packages) to your machine and going through this getting started. From there you can try [writing a Kitfile](./kitfile/format.md) for your own AI/ML project.

## Preparing for packaging

The first step is to make a `Kitfile` - a YAML manifest for your ModelKit. There are four main parts to a Kitfile:
1. ModelKit metadata in the `package` section
1. Path to the Jupyter notebook folder in the `code` section
1. Path to the datasets in the `datasets` section (you can have multiple datasets in the same page)
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

If you have questions please join our [Discord server](https://discord.gg/eHXGmHds).