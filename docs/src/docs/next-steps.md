# Next Steps with Kit

In this guide you'll learn how to:
* Make your own Kitfile
* The power of `unpack`
* Read the Kitfile or manifest from a ModelKit
* Tag ModelKits and keep your registry tidy

## Making your own Kitfile

A Kitfile is the configuration document for your ModelKit. It's written in YAML so it's easy to read. There are four main parts to a Kitfile:
1. ModelKit metadata in the `package` section, including the author, description, and license
1. Path to the Jupyter notebook folder in the `code` section
1. Path to the datasets in the `datasets` section (you can have multiple datasets in the same page)
1. Path to the serialized model in the `model` section

A Kitfile only needs the `package` section, plus one of the other sections.

The `datasets` and `code` sections can include multiple assets, but the `model` section can only contain a single model (you can chain models by using multiple ModelKits).

Here's an example Kitfile:

```yaml
manifestVersion: v1.0.0

package:
  authors:
  - Jozu
  description: Small language model based on Mistral-7B fine tuned for answering film photography questions.
  license: Apache-2.0
  name: FilmSLM

model:
  description: Film photography Q&A model using Mistral-7B
  framework: Mistral-7B
  license: Apache-2.0
  name: FilmSLM
  path: ./models/film_slm:champion
  version: 1.2.6

datasets:
- description: Forum postings from sites like rangefinderforum, DPreview, PhotographyTalk, and r/AnalogCommunity
  name: training data
  path: ./data/forum-to-2023-train.csv
- description: validation data
  name: validation data
  path: ./data/test.csv

code:
- description: Jupyter notebook with model training code in Python
  path: ./notebooks
```

A minimal ModelKit for distributing a pair of datasets might look like this:
```yaml
manifestVersion: v1.0.0

datasets:
- name: training data
  path: ./data/train.csv
- description: validation data
  name: validation data
  path: ./data/validate.csv
```

Right now you can only build ModelKits from files on your local system...but don't worry we're already working towards allowing you to [reference remote files](https://github.com/jozu-ai/kitops/issues/85). For example, building a ModelKit from a local notebook and model, but a dataset hosted on DvC, S3, or anywhere else.

Once you've authored a Kitfile for your AI/ML project you can pack it up and store it in your local or remote repository.

```sh
kit pack . -t film-slm:champion
kit push docker.io/jozubrad/film-slm:champion
```

# The power of unpack

Models and their datasets can be very large and take a long time to push or pull, so Kit includes the `unpack` command that allows you to pull only pieces of the ModelKit you need, saving time and storage space:

`unpack` can take arguments for partial pulling of a ModelKit:
* `--model` to pull only the model to the destination file system
* `--datasets` to pull only the datasets to the destination file system
* `--code` to pull only the code bases to the destination file system
* `--config` to pull only the `Kitfile` to the destination file system

For example:

```sh
kit unpack mymodel:challenger --model -d ./model
```

Will extract the model from the `mymodel:challenger` ModelKit and place it in a local directory called `/model`.

The `unpack` command is part of the typical push and pull commands:
* `pack` will pack up a set of assets into a ModelKit package.
* `push` will push the whole ModelKit to a registry.
* `pull` will pull the whole ModelKit from a registry.
* `unpack` will extract all the assets from the ModelKit package. 

## Read the Kitfile or manifest from a ModelKit

For any ModelKit in your local or remote registry you can easily read the Kitfile without pulling or unpacking it. This is a great way to understand what's in a ModelKit you might be interested in without needing to execute the more time-consuming unpack/pull comamnds.


```sh
kit info mymodel:challenger
```

Will print the following to your terminal:

```sh
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

If you need more details, like the size or file format of the contents, you can print the manifest to the terminal using the inspect command:

```sh
kit inspect mymodel:challenger
```

```sh
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

## Tag ModelKits and keep your registry tidy

Tagging is a great way to version your ModelKits as they move through the development lifecycle. For example, during development the model I'm working on currently may always be tagged as the `latest` so my team knows which is most current. At the same time the model that's operating in production for my customers may be tagged the `champion`.

However, after testing my latest model, I find that its scores are much higher than the current champion model. At that point I may want to tag it `challenger` so everyone knows that this is likely to be the next model we deploy to production, to replace our current champion model.

To do that I can simply add a new tag based on the existing latest tag.

```sh
kit tag mymodel:latest mymodel:challenger
```

If you run `kit list` you'll now see that you have two models in your local registry:

```sh
kit list

REPOSITORY  TAG         MAINTAINER   NAME             SIZE       DIGEST
mymodel     latest      Rajat        Finetuning_SLM   13.1 MiB   sha256:f268a74ff85a00f2a68400dfc17b045bc7c1638da7f096c7ae400ad5bdfd520c

mymodel     challenger  Rajat        Finetuning_SLM   13.1 MiB   sha256:f268a74ff85a00f2a68400dfc17b045bc7c1638da7f096c7ae400ad5bdfd520c
```

You'll notice that the digest of the two ModelKits is the same so while you have two tags, the ModelKit itself is only stored once for efficiency. This makes tags an efficient way to mark models that have hit specific milestones that are meaningful to your organization's development lifecycle.

Now let's imagine that we deploy our challenger model. At that point we'd tag it as `champion`:

```sh
kit tag mymodel:challenger mymodel:champion
```

If you run `kit list` you'll now see that you have two models in your local registry:

```sh
kit list

REPOSITORY  TAG         MAINTAINER   NAME             SIZE       DIGEST
mymodel     latest      Rajat        Finetuning_SLM   13.1 MiB   sha256:f268a74ff85a00f2a68400dfc17b045bc7c1638da7f096c7ae400ad5bdfd520c

mymodel     challenger  Rajat        Finetuning_SLM   13.1 MiB   sha256:f268a74ff85a00f2a68400dfc17b045bc7c1638da7f096c7ae400ad5bdfd520c

mymodel     champion    Rajat        Finetuning_SLM   13.1 MiB   sha256:f268a74ff85a00f2a68400dfc17b045bc7c1638da7f096c7ae400ad5bdfd520c
```

Now we no longer want this ModelKit to be tagged as `challenger` because it's the `champion` so we can remove it from our registry:

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

Thanks for taking some time to play with Kit. We'd love to hear what you think. Feel free to drop us an [issue in our GitHub repository](https://github.com/jozu-ai/kitops/issues) or join [our Discord server](https://discord.gg/YyAfWnEg).
