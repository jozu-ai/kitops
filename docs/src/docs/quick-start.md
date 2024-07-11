<script setup>
import vGaTrack from '@theme/directives/ga'
</script>

# KitOps Quick Start

In this guide, we'll use ModelKits and the kit CLI to easily:
* Package up a model, notebook, and datasets into a single ModelKit you can share through your existing tools
* Push the ModelKit package to a public or private registry
* Grab only the assets you need from the ModelKit for testing, integration, local running, or deployment
* Run an LLM locally to speed app integration, testing, or experimentation

## Before we start...

1. Make sure you've got the [Kit CLI setup](./cli/installation.md).
2. Create and navigate to a new folder on your filesystem - we suggest calling it `KitStart` but any name works.

## Learning to use the CLI

### 1/ Check your CLI Version

Check that the Kit CLI is properly installed by using the [version command](./cli/cli-reference.md#kit-version).

```sh
kit version
```

You'll see information about the version of Kit you're running. If you get an error check to make sure you have [Kit installed](./cli/installation.md) and in your path.

### 2/ Login to Your Registry

You can use the [login command](./cli/cli-reference.md#kit-login) to authenticate with any container registry, but we'll use **GitHub Registry** throughout this guide.

```sh
kit login ghcr.io
```

You'll see `Log in successful`. If you get an error it may be that you need an HTTP vs HTTPS (default) connection. Try the login command again but with `--plain-http`.

### 3/ Get a Sample ModelKit

Let's use the [unpack command](./cli/cli-reference.md#kit-unpack) to pull a sample ModelKit to our machine that we can play with. In this case we'll unpack the whole thing, but one of the great things about Kit is that you can also selectively unpack only the artifacts you need: just the model, the model and dataset, the code, the configuration...whatever you want. Check out the `unpack` [command reference](./cli/cli-reference.md#kit-unpack) for details.

You can grab <a href="https://github.com/orgs/jozu-ai/packages"
  v-ga-track="{
    category: 'link',
    label: 'grab any of the ModelKits',
    location: 'docs/quick-start'
  }">any of the ModelKits</a> from our site, but we've chosen a small language model example below.

```sh
kit unpack ghcr.io/jozu-ai/modelkit-examples/finetuning_slm:latest
```

You'll see a set of messages as Kit unpacked the configuration, code, datasets, and serialized model. Now list the directory contents:

```sh
ls
```

You'll see a single file (`Kitfile`) which is the manifest for our ModelKit, and a set of files or directories including adapters, a Jupyter notebook, and dataset.

### 4/ Check the Local Repository

Use the [list command](./cli/cli-reference.md#kit-list) to check what's in our local repository.

```sh
kit list
```

You'll see the column headings for an empty table with things like `REPOSITORY`, `TAG`, etc...

### 5/ Pack the ModelKit

Since our repository is empty we'll need use the [pack command](./cli/cli-reference.md#kit-pack) to create our ModelKit. The ModelKit in your local registry will need to be named the same as your remote registry. So the command will look like: `kit pack . -t [your registry address]/[your repository name]/mymodelkit:latest`

In my case I am pushing to the `jozubrad` repository:

```sh
kit pack . -t ghcr.io/jozubrad/mymodelkit:latest
```

You'll see a set of `Saved ...` messages as each piece of the ModelKit is saved to the local repository.

Checking your local registry again you should see an entry:

```sh
kit list
```

The new entry will be named based on whatever you used in your pack command.

### 6/ (Optional) Remove a ModelKit from a Local Repository

Let's pretend that the `pack` command we ran in the previous step contained a typo in the ModelKit's repository name causing the word 'model' to be entered as 'modle'.The output from the `kit list` command would display the ModelKit as:

```sh
ghcr.io/jozubrad/mymodlekit:latest
```

To correct this, we would `remove` the misspelled ModelKit from our local repository using the [remove command](./cli/cli-reference.md#kit-remove), being sure to provide reference the ModelKit using its mispelled name:

```sh
kit remove ghcr.io/jozubrad/mymodlekit:latest
```

Next, we would repeat the `kit pack` command in the previous step, being sure to provide the correct repository name for our ModelKit.

### 7/ Push the ModelKit to a Remote Repository

The [push command](./cli/cli-reference.md#kit-push) will copy the newly built ModelKit from your local repository to the remote repository you logged into earlier. The naming of your ModelKit will need to be the same as what you see in your `kit list` command (REPOSITORY:TAG). You can even copy and paste it. In my case it looks like:

```sh
kit push ghcr.io/jozubrad/mymodelkit:latest
```

### 8/ Run an LLM Locally

If you're using Kit with LLMs you can quickly run the model locally to speed integration, testing, or experimentation.

Create a new directory for your LLM:

```sh
mkdir devmode
cd devmode
```

Now unpack an LLM ModelKit - we have [several](https://github.com/orgs/jozu-ai/packages), but I've chosen Phi3:

```sh
kit unpack ghcr.io/jozu-ai/phi3:3.8b-mini-instruct-4k-q4_K_M
```

Now start your LLM dev server locally using the [dev start command](./cli/cli-reference.md#kit-dev-start):

```sh
kit dev start .
```

In the command output you'll see a URL you can use to interact with the LLM (there's a command flag if you want to always use the same port). You can control parameters of the model, change the prompt, or chat with the LLM.

If you need to get logs use the [dev logs command](./cli/cli-reference.md#kit-dev-logs):

```sh
kit dev logs
```

When you're done don't forget to stop the LLM dev server:

```sh
kit dev stop
```

### Congratulations

You've learned how to unpack a ModelKit, pack one up, push it, and run an LLM locally. Anyone with access to your remote repository can now pull your new ModelKit and start playing with your model:

```sh
kit pull ghcr.io/jozubrad/mymodelkit:latest
```

If you'd like to learn more about using Kit, try our [Next Steps with Kit](./next-steps.md) document that covers:
* Signing your ModeKit
* Making your own Kitfile
* The power of `unpack`
* Tagging ModelKits
* Keeping your registry tidy

Thanks for taking some time to play with Kit. We'd love to hear what you think. Feel free to drop us an [issue in our GitHub repository](https://github.com/jozu-ai/kitops/issues) or join [our Discord server](https://discord.gg/3eDb4yAN).
