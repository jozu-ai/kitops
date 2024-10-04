<script setup>
import vGaTrack from '@theme/directives/ga'
</script>

# KitOps Getting Started

In this guide, we'll use ModelKits and the kit CLI to easily:
* Package up a model, notebook, and datasets into a single ModelKit you can share through your existing tools
* Push that versioned ModelKit package to a registry
* Grab only the assets you need from the ModelKit for testing, integration, local running, or deployment

## Before we start...

1. Make sure you've got the [Kit CLI setup](./cli/installation.md).
2. Create and navigate to a new folder on your filesystem - we suggest calling it `KitStart` but any name works.

## Learning to use the CLI

### 1/ Check your CLI Version

Check that the Kit CLI is properly installed by using the [version command](./cli/cli-reference.md#kit-version).

```sh
kit version
```

You'll see information about the version of Kit you're running. If you get an error check to make sure you have [Kit installed and in your path](./cli/installation.md).

### 2/ Login to Your Registry

You can use the [login command](./cli/cli-reference.md#kit-login) to authenticate with any OCI v1.1-compatible container registry - local or remote (you can see our [list of compliant registries](./modelkit/compatibility.md)). In this guide we'll use the [Jozu Hub](https://jozu.ml/) because it's free to sign-up and provides more detail on what's inside each ModelKit like whether it's signed or has provenance. You can substitute your own repository if preferred.

```sh
kit login jozu.ml
```

After entering your username and password, you'll see `Log in successful`. If you get an error it may be that you need an HTTP vs HTTPS (default) connection. Try the login command again but with `--plain-http`.

### 3/ Get a Sample ModelKit

Let's use the [unpack command](./cli/cli-reference.md#kit-unpack) to pull a [sample ModelKit from Jozu Hub](https://jozu.ml/browse) to our machine that we can play with. In this case we'll unpack the whole thing, but one of the great things about Kit is that you can also selectively unpack only the artifacts you need: just the model, the model and dataset, the code, the configuration...whatever you want. Check out the `unpack` [command reference](./cli/cli-reference.md#kit-unpack) for details.

You can grab <a href="https://jozu.ml/discover"
  v-ga-track="{
    category: 'link',
    label: 'grab any of the ModelKits',
    location: 'docs/get-started'
  }">any of the ModelKits</a> from the Jozu Hub, but we've chosen a fine-tuned model based on Llama3.

The unpack command will unpack the ModelKit contents to the current directory by default. If you want it unpacked to a specific directory use the `-d /path/to/unpacked`.

```sh
kit unpack jozu.ml/jozu/fine-tuning:latest
```

You'll see a set of messages as Kit unpacks the configuration, code, datasets, and serialized model. Now list the directory contents:

```sh
ls
```

You'll see:
* A Kitfile
* A README file
* A Llama3 model in GGUF format
* A LoRA adapter in GGUF format
* A training dataset

The [Kitfile](./kitfile/kf-overview.md) is the manifest for our ModelKit, the serialized model, and a set of files or directories including the adapter, dataset, and docs. Every ModelKit has a Kitfile and you can use the info and inspect commands to view them from the CLI (there's more on this in our [Next Steps](next-steps.md) doc).

### 4/ Check the Local Repository

Use the [list command](./cli/cli-reference.md#kit-list) to check what's in our local repository.

```sh
kit list
```

You'll see the column headings for an empty table with things like `REPOSITORY`, `TAG`, etc...

### 5/ Pack the ModelKit

Since our repository is empty we'll need use the [pack command](./cli/cli-reference.md#kit-pack) to create our ModelKit. The ModelKit in your local registry will need to be named the same as your remote registry. The command will look like: `kit pack . -t [your registry address]/[your registry user or organization name]/[your repository name]:[your tag name]`

In our case we are pushing a ModelKit tagged `latest` to:
* The [Jozu Hub](https://jozu.ml/) registry
* The `brad` user organization
* The `quick-start` repository

As a result, the command will look like:

```sh
kit pack . -t jozu.ml/brad/quick-start:latest
```

You may need to substitute your own registry, user, repository, or tag names.

Once complete, you'll see a set of `Saved ...` messages as each piece of the ModelKit is saved to the local repository.

Check your local registry again:

```sh
kit list
```

You should see an entry named based on whatever you used in your pack command.

### 6/ (Optional) Remove a ModelKit from a Local Repository

If you have a typo when packing a ModelKit you can easily remove it from your repository and try again. The [Next Steps guide includes information on how to remove ModelKits](./next-steps.md#remove-command).

Once you've removed the mistaken ModelKit from the repository, you can repeat the `kit pack` command in the previous step, being sure to provide the correct organization and repository name for your ModelKit.

### 7/ Push the ModelKit to a Remote Repository

The [push command](./cli/cli-reference.md#kit-push) will copy the newly built ModelKit from your local repository to the remote repository you logged into earlier. The naming of your ModelKit will need to be the same as what you see in your `kit list` command (REPOSITORY:TAG). You can even copy and paste it. In our case it looks like:

```sh
kit push jozu.ml/brad/quick-start:latest
```

Note that some registries, like Jozu Hub, don't automatically create a repository. If you receive an error from your `push` command, make sure you have created the repository in your target registry and that you have push rights to the repository.

### Congratulations

You've learned how to unpack a ModelKit, pack one up, and push it. Anyone with access to your remote repository can now pull your new ModelKit and start playing with your model using the `kit pull` or `kit unpack` commands.

If you'd like to learn more about using Kit, try our [Next Steps with Kit](./next-steps.md) document that covers:
* Signing your ModeKit
* Making your own Kitfile
* The power of `unpack`
* Tagging ModelKits
* Keeping your registry tidy

Or, if you want to run an LLM-based ModelKit locally try our [dev mode](./dev-mode.md)

Thanks for taking some time to play with Kit. We'd love to hear what you think. Feel free to drop us an [issue in our GitHub repository](https://github.com/jozu-ai/kitops/issues) or join [our Discord server](https://discord.gg/3eDb4yAN).
