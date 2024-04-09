<script setup>
import vGaTrack from '@theme/directives/ga'
</script>

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

### 3/ Get a Sample ModelKit

Let's unpack a sample ModelKit to our machine that we can play with. In this case we'll unpack the whole thing, but one of the great things about Kit is that you can also selectively unpack only the thigs you need: just the model, the model and dataset, the code, the configuration...whatever you want. Check out the `unpack` [command reference](./cli/cli-reference.md) for details.

You can grab <a href="https://github.com/orgs/jozu-ai/packages"
  v-ga-track="{
    category: 'link',
    label: 'grab any of the ModelKits',
    location: 'docs/quick-start'
  }">any of the ModelKits</a>from our site, but we've chosen a small language model example below.

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

If you'd like to learn more about using Kit, try our [Next Steps with Kit](./next-steps.md) document that covers:
* Making your own Kitfile
* The power of `unpack`
* Tagging ModelKits
* Keeping your registry tidy

Thanks for taking some time to play with Kit. We'd love to hear what you think. Feel free to drop us an [issue in our GitHub repository](https://github.com/jozu-ai/kitops/issues) or join [our Discord server](https://discord.gg/3eDb4yAN).
