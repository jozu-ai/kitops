# KitOps Roadmap

As a public community project the KitOps roadmap is always open to feedback. If you think a feature lower in the list is more important, or feel that we've missed an important feature please:

* [Open an issue](https://github.com/jozu-ai/kitops/issues) in the repo
* Tell us in the [KitOps Discord](https://discord.gg/Tapeh8agYy) channel
* Send us a message on [X / Twitter](https://twitter.com/Kit_Ops)

## Roadmap
🛳️ = Shipped

🚗 = In-progress

📅 = Planned

💡 = Idea (please provide feedback)

### Dev Mode Improvements

* Change `kit dev` implementation to use [Llamafile](https://github.com/Mozilla-Ocho/llamafile?tab=readme-ov-file) 📅
* Export [Llamafile](https://github.com/Mozilla-Ocho/llamafile?tab=readme-ov-file) from ModelKit 📅
* Generate example code as the parameters and prompts are entered 💡
* A way to see JSON communicated between the server and responses 💡
* Hide the parameters that are not frequently changed 💡
* Add `kit dev` command to Kit CLI on *Windows* 💡
* Add `kit dev` command to Kit CLI on *Linux* 💡

### Tutorials

#### Demos

* Creating a ModelKit from a Jupyter Notebook 📅

### Signing

Users can optionally sign their ModelKit using something like Cosign in order to add an extra layer of security to their packaging.

* Add docs on signing and verifying to docs 🛳️
* Decide on a signing utility 💡
* Add signing to Kit CLI 💡

### Deployment 💡

Users want to be able to deploy their models through existing CI/CD/CT pipelines. Since KitOps doesn't know the details of a user's deployment pipeline or process we will `unpack` a ModelKit into an appropriate directory structure for one of several deployment targets.

### CLI Distribution

Users want to be able to get the Kit CLI from locations like Brew, Choco, and Conda.

* Add support for Brew 🚗
* Create a Kit Python library 💡
* Add support for Choco 💡
* Add support for Conda 💡

### Attestation

This feature will come in two parts: build attestation and self-attestation. Build attestations will be done by KitOps itself, adding a SLSA attestation about how the ModelKit was built. ModelKit creators can add an optional attestation for specific assets in the ModelKit, or the whole ModelKit. Users can include any 3rd party attestation Verification Summary Attestation (VSA). Additional attestations could be added as predicates. Attestation would be included as a separate layer in the ModelKit.

* Add provenance for KitOps ModelKit builds 💡
* Add CLI warnings if attestation was expected and not found 💡
