# Next Steps with Kit

In this guide you'll learn how to:
* Sign your ModelKit
* Make your own Kitfile
* The power of `unpack`
* Read the Kitfile or manifest from a ModelKit
* Tag ModelKits and keep your registry tidy

## Signing your ModelKit

Because ModelKits are OCI 1.1 artifacts, they can be signed like any other OCI artifact (you may already sign your containers, for example).

If you need a quick way to sign a ModelKit you can follow the same instructions as for a container, using a tool like [Cosign](https://docs.sigstore.dev/cosign/signing/signing_with_containers/).


## Making your own Kitfile

A Kitfile is the configuration document for your ModelKit. It's written in YAML so it's easy to read. There are five main parts to a Kitfile:

1. The `package` section: Metadata about the ModelKit, including the author, description, and license
1. The `model` section: Information about the serialized model
1. The `docs` section: Information about documentation for the ModelKit
1. The `code` section: Information about codebases related to the project, including Jupyter notebook folders
1. The `datasets` section: Information on included datasets

A Kitfile only needs the `package` section, plus one or more of the other sections.

The `model` section can contain a single model, or you can create model dependencies with `model parts` which is covered in the [KitFile format documentation](https://kitops.ml/docs/kitfile/format.html#model).

The `datasets`, `code`, and `docs` sections are lists, so each entry must start with a dash. The dash is required even if you are only packaging a single item of that type.

Here's a snippet of a KitFile that contains two datasets, notice that each starts with "-":

```yaml
datasets:
- description: Forum postings from photo sites
  name: training data
  path: ./data/forum-to-2023-train.csv

- name: validation data
  path: ./data/test.csv
```

...and here's an example of a single codebase, it still needs to start with a "-":

```yaml
code:
- description: Jupyter notebook with model training code in Python
  path: ./notebooks
```

Any relative paths defined within the Kitfile are interpreted as being relative to the context directory sent to the `kit pack` command.

### Kitfile Examples

Here's a complete Kitfile example, with a model, documentation, two datasets, and a codebase:

```yaml
manifestVersion: 1.0

package:
  authors:
    - Jozu
  description: Small language model based on Mistral-7B fine tuned for answering film photography questions.
  license: Apache-2.0
  name: FilmSLM

model:
  name: FilmSLM
  description: Film photography Q&A model using Mistral-7B
  framework: Mistral-7B
  license: Apache-2.0
  path: ./models/film_slm:champion
  version: 1.2.6

docs:
  - path: ./README.md
    description: Readme file for this ModelKit
  - path: ./USAGE.md
    description: Information on how to use this model for inference

datasets:
  - description: Forum postings from sites like rangefinderforum, PhotographyTalk, and r/AnalogCommunity
    name: training data
    path: ./data/forum-to-2023-train.csv
  - description: validation data
    name: validation data
    path: ./data/test.csv

code:
  - description: Jupyter notebook with model training code in Python
    path: ./notebooks
```

A minimal ModelKit for distributing a pair of datasets looks like this:

```yaml
manifestVersion: v1.0.0

package:
  authors:
    - Jozu

datasets:
  - name: training data
    path: ./data/train.csv
    license: Apache-2.0
  - description: validation data
    name: validation data
    path: ./data/validate.csv
```

More information on Kitfiles can be found in the [Overview](./kitfile/kf-overview.md) and [Format](./kitfile/format.md) documentation.

### Using the Kitfile to Pack a ModelKit

When you're done writing the Kitfile, name it `Kitfile` without an extension. Now you can use the [kit pack command](./cli/cli-reference.md#kit-pack) to build your ModelKit.

To pack a ModelKit with a Kitfile in the current directory, name it "film-slm", attach the "champion" tag, and store it in the local registry:

```sh
kit pack . -t film-slm:champion
```

To pack a ModelKit with the same settings, but using a Kitfile stored elsewhere:

```sh
kit pack . -f /path/to/your/Kitfile -t film-slm:champion
```

### Pushing to a Remote Registry

In each case this will pack a ModelKit and store it in your local registry. To push it to a remote registry for sharing with others, there are two steps:

1. Tagging the local copy with the remote registry's name
1. Pushing the remote-named copy from your local to the remote registry

Let's imagine we want to push our new ModelKit to Docker Hub:

First, you need to tag the image in your local registry with the remote registry's name:

```sh
kit pack . -t docker.io/jozubrad/film-slm:champion
```

Second, you will need to login, then [kit push](./cli/cli-reference.md#kit-push) your local image to the remote registry:

```sh
kit login docker.io
kit push docker.io/jozubrad/film-slm:champion
```

# The Power of Unpack

Models and their datasets can be very large and take a long time to push or pull, so Kit includes a unique and flexible [unpack command](./cli/cli-reference.md#kit-unpack) that allows you to pull only pieces of the ModelKit you need, saving time and storage space:

`unpack` can take arguments for partial unpacking of a ModelKit:
* `--model` to unpack only the model to the destination file system
* `--docs` to unpack only the documentation to the destination file system
* `--datasets` to unpack only the datasets to the destination file system
* `--code` to unpack only the code bases to the destination file system
* `--config` to unpack only the Kitfile to the destination file system

For example:

```sh
kit unpack mymodel:challenger --model -d ./model
```

Will extract only the model from the `mymodel:challenger` ModelKit, and place it in a local directory called `/model`.

The `unpack` command is part of the typical push and pull commands:
* `pack` will pack up a set of assets into a ModelKit package.
* `push` will push the whole ModelKit to a registry.
* `pull` will pull the whole ModelKit from a registry.
* `unpack` will extract all, or selected assets, from the ModelKit.

## Read the Kitfile or Manifest from a ModelKit

For any ModelKit in your local or remote registry you can use the [info command](./cli/cli-reference.md#kit-info) to easily read the Kitfile without pulling or unpacking it. This is a great way to understand what's in a ModelKit you might be interested in without needing to execute the more time-consuming unpack/pull commands.

```sh
kit info mymodel:challenger
```

Will print the following to your terminal:

```yaml
manifestVersion: v1.0.0
package:
  name: Finetuning_SLM
  description: This Kitfile contains all the necessary files required for finetuning SLM.
  license: Apache-2.0
  authors: [Rajat]
model:
  name: keras Model
  path: ./model_adapter
  framework: Scikit-learn
  version: 1.0.0
  description: Flight satisfaction and trait analysis model using Scikit-learn
  license: Apache-2.0
code:
  - path: ./SLM_Finetuning.ipynb
    description: Jupyter notebook with model training code in Python
datasets:
  - name: training data
    path: ./slm_tuning_dataset.csv
    description: UCF Video Dataset
```

If you need more details, like the size, file format, or SHA digest of the contents, you can use [kit inspect](./cli/cli-reference.md#kit-inspect) to print the manifest to the terminal using the inspect command:

```sh
kit inspect mymodel:challenger
```

```json
{
  "schemaVersion": 2,
  "config": {
    "mediaType": "application/vnd.kitops.modelkit.config.v1+json",
    "digest": "sha256:58444ef30d1cc7ee0fd2a24697e26252c38bf4317bfc791cd30c5af0d6f91f8f",
    "size": 619
  },
  "layers": [
    {
      "mediaType": "application/vnd.kitops.modelkit.code.v1.tar+gzip",
      "digest": "sha256:144ad8a89c3946ab794455450840f1b401511cdd8befc95884826376ef56d861",
      "size": 12657
    },
    {
      "mediaType": "application/vnd.kitops.modelkit.dataset.v1.tar+gzip",
      "digest": "sha256:7d0f8d60895bba5c56e032f509c92553cfba1b014eee58f01e46d7af923099e8",
      "size": 13787789
    },
    {
      "mediaType": "application/vnd.kitops.modelkit.model.v1.tar+gzip",
      "digest": "sha256:25646107e93d62ccd55b2cf14dfd55f7bbf426e5492a42d3d7f1bfd2cde30035",
      "size": 122
    }
  ]
}
```

`size` is shown in bytes. For more information on the manifest can be found in the [specification documentation](./modelkit/spec.md).

## Tag ModelKits and Keep Your Registry Tidy

### Tag Command

Tagging is a great way to version your ModelKits as they move through the development lifecycle. For example, during development the model I'm working on currently may always be tagged as the `latest` so my team knows which is most current. At the same time the model that's operating in production for my customers may be tagged the `champion`.

However, after testing my latest model, if I find that its scores are much higher than the current champion model I may tag it `challenger` so everyone knows that this is likely to be the next model we deploy to production, replacing our current champion model.

To do that I can create a new ModelKit and use the [tag command](./cli/cli-reference.md#kit-tag). For example to change from `latest` to `challenger`:

```sh
kit tag mymodel:latest mymodel:challenger
```

If you run [kit list](./cli/cli-reference.md#kit-list) you'll now see that you have two models in your local registry:

```sh
kit list

REPOSITORY  TAG         MAINTAINER   NAME             SIZE       DIGEST
mymodel     latest      Rajat        Finetuning_SLM   13.1 MiB   sha256:f268a74ff85a00f2a68400dfc17b045bc7c1638da7f096c7ae400ad5bdfd520c

mymodel     challenger  Rajat        Finetuning_SLM   13.1 MiB   sha256:f268a74ff85a00f2a68400dfc17b045bc7c1638da7f096c7ae400ad5bdfd520c
```

You'll notice that the digest of the two ModelKits is the same so while you have two tags, the ModelKit itself is only stored once for efficiency. This makes tags an efficient way to mark models that have hit specific milestones that are meaningful to your organization's development lifecycle. Of course, if the contents of the ModelKit have changed since I last packaged it, I should use the `pack -t` command with the new tag so I get an updated ModelKit.

Now let's imagine that we deploy our challenger model. At that point we'd tag it as `champion`:

```sh
kit tag mymodel:challenger mymodel:champion
```

If you run `kit list` you'll now see that you have three models in your local registry:

```sh
kit list

REPOSITORY  TAG         MAINTAINER   NAME             SIZE       DIGEST
mymodel     latest      Rajat        Finetuning_SLM   13.1 MiB   sha256:f268a74ff85a00f2a68400dfc17b045bc7c1638da7f096c7ae400ad5bdfd520c

mymodel     challenger  Rajat        Finetuning_SLM   13.1 MiB   sha256:f268a74ff85a00f2a68400dfc17b045bc7c1638da7f096c7ae400ad5bdfd520c

mymodel     champion    Rajat        Finetuning_SLM   13.1 MiB   sha256:f268a74ff85a00f2a68400dfc17b045bc7c1638da7f096c7ae400ad5bdfd520c
```

### Remove Command

Sometimes you want to remove a ModelKit that you've packed and stored in the repository. The `kit remove` command comes to the rescue.

In this case, we no longer want the `challenger` ModelKit since it's a duplicate of the `champion` now. [Removing it](./cli/cli-reference.md#kit-remove) from our registry will keep things clean and clearer for other users:

```sh
kit remove mymodel:challenger

Removed localhost/mymodel:challenger (digest sha256:f268a74ff85a00f2a68400dfc17b045bc7c1638da7f096c7ae400ad5bdfd520c)
```

Running kit list again will show our updated local registry contents:

```sh
kit list

REPOSITORY  TAG         MAINTAINER   NAME             SIZE       DIGEST
mymodel     latest      Rajat        Finetuning_SLM   13.1 MiB   sha256:f268a74ff85a00f2a68400dfc17b045bc7c1638da7f096c7ae400ad5bdfd520c

mymodel     champion    Rajat        Finetuning_SLM   13.1 MiB   sha256:f268a74ff85a00f2a68400dfc17b045bc7c1638da7f096c7ae400ad5bdfd520c
```

You can learn more about all the Kit CLI commands from our [command reference doc](./cli/cli-reference.md).

To learn about how to run an LLM locally using Kit, see our [Kit Dev](./dev-mode.md) documentation.

Thanks for taking some time to play with Kit. We'd love to hear what you think. Feel free to drop us an [issue in our GitHub repository](https://github.com/jozu-ai/kitops/issues) or join [our Discord server](https://discord.gg/Tapeh8agYy).
