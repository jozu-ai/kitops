# KitOps Roadmap

As a public community project the KitOps roadmap is always open to feedback. If you think a feature lower in the list is more important, or feel that we've missed an important feature please:

* [Open an issue](https://github.com/jozu-ai/kitops/issues) in the repo
* Tell us in the [KitOps Discord](https://discord.gg/Tapeh8agYy) channel
* Send us a message on [X / Twitter](https://twitter.com/Kit_Ops)

## Roadmap
🛳️ = Shipped

🏃‍➡️= In-progress

### ModelKit Chaining 🛳️

A ModelKit can refer to another ModelKit's model as its "base model." This is useful for situations where a user wants to take an existing model (e.g., Llama 3) but create their own model based on it. The new ModelKit would include a reference to the Llama 3 ModelKit and Kit would understand to pull both when the user asks for the new ModelKit.

* Add chaining to ModelKit schema 
* Add support for chaining to relevant Kit CLI commands 

### Dev Mode 🏃‍➡️

Users can run a ModelKit packaged LLM on localhost and interact with it via chat or prompt through a beautiful and user-friendly UI. This speeds up running, testing, and integrating LLMs with applications.

* Design UX for chat, prompt, and parameter interactions 🛳️
* Code UI for interactions focused on desktop form factor 🛳️
* Create ModelKits for popular LLMs 🛳️
* Add `dev` command to Kit CLI on *MacOS* 🛳️
* Generate example code as the parameters and prompts are entered 🏃‍➡️
* A way to see JSON communicated between the server and responses 🏃‍➡️
* Hide the parameters that are not frequently changed 🏃‍➡️
* Add `dev` command to Kit CLI on *Windows* (TBC)
* Add `dev` command to Kit CLI on *Linux* (TBC)

### Tutorials 🏃‍➡️

Users can learn how to use ModelKits as part of a RAG pipeline, or while fine-tuning their LLM. Helps people get the most value out of ModelKits during their LLM development work.

* Tutorial: Fine-tuning LLMs using KitOps 🚢
* Tutorial: Doing RAG with KitOps 🏃‍➡️

### Signing 🏃‍➡️

Users can optionally sign their ModelKit using something like Cosign in order to add an extra layer of security to their packaging.

* Decide on a signing utility 🏃‍➡️
* Add signing to Kit CLI 🏃‍➡️
* Add docs on signing and verifying to docs 🏃‍➡️

### Workflow Documentation 🏃‍➡️

To help users understand how ModelKits fit into their existing workflows we will create documentation explaining how CI/CD/CT, observability, and other tasks can be simplified with KitOps.

### Deployment

Users want to be able to deploy their models through existing CI/CD/CT pipelines. Since KitOps doesn't know the details of a user's deployment pipeline or process we will `unpack` a ModelKit into an appropriate directory structure for one of several deployment targets.

### CLI Distribution

Users want to be able to get the Kit CLI from locations like Brew, Choco, and Conda.

### Attestation

This feature will come in two parts: build attestation and self-attestation. Build attestations will be done by KitOps itself, adding a SLSA attestation about how the ModelKit was built. ModelKit creators can add an optional attestation for specific assets in the ModelKit, or the whole ModelKit. Users can include any 3rd party attestation Verification Summary Attestation (VSA). Additional attestations could be added as predicates. Attestation would be included as a separate layer in the ModelKit.

* Add attestation for KitOps ModelKit builds
* Add mechanism for adding attestations to ModelKits
* Update CLI to make use of attestations
* Add CLI warnings if attestation was expected and not found