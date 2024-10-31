# Kit Dev: Run an LLM Locally

:::info
This is a beta feature only available on MacOS today. To provide feedback (we love that) you can [file an issue](https://github.com/jozu-ai/kitops/issues) in our GitHub repo, or [join our Discord](https://discord.gg/Tapeh8agYy) channel.
:::

If you're using Kit with LLMs you can quickly run the model locally to speed integration, testing, or experimentation.

To run the ModelKit locally, first create a new directory for your LLM:

```sh
mkdir kitdev
cd kitdev
```

Now unpack an LLM ModelKit - there are several on [Jozu Hub](https://jozu.ml/discover), but here we're using Phi3 Mini because of its size:


```sh
kit unpack jozu.ml/jozu/phi3:3.8b-mini-instruct-4k-q4_K_M
```

Now start your LLM dev server locally using the [kit dev start command](./cli/cli-reference.md#kit-dev-start):

```sh
kit dev start .
```

In the command output you'll see a URL you can use to interact with the LLM (there's a command flag to always use the same port). You can control parameters of the model, change the prompt, or chat with the LLM.

If you need to get logs use the [dev logs command](./cli/cli-reference.md#kit-dev-logs):

```sh
kit dev logs
```

When you're done don't forget to stop the Kit dev server:

```sh
kit dev stop
```

**Questions or suggestions?** Drop an [issue in our GitHub repository](https://github.com/jozu-ai/kitops/issues) or join [our Discord server](https://discord.gg/Tapeh8agYy) to get support or share your feedback.