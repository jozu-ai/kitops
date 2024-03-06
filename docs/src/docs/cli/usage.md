# Kit CLI tool

The `kit CLI` is a tool to easily and quickly manage models.

## Usage

```sh
$ ./kit [command]
```

You can always get help on a command by adding the `-h` flag.

Available Commands:

| Command |	Description |
| ---- | --- |
| `build` | Build a ModelKit |
| `completion` | Generate the autocompletion script for the specified shell |
| `dev` | Run the serialized model | <!-- starts a server on a given port and drops the model in there for inference -->
| `fetch` | Updating the local respository for a ModelKit from the remote |
| `help` | Help about any command |
| `list` | List ModelKits |
| `login` | Login to the remote repository |
| `logout` | Logout to the remote repository |
| `pull` | Pull one or more of the model, dataset, code, and Kitfile into a destination folder |
| `push` | Push ModelKit to respository |
| `remove` | Removed the ModelKit from the local repository |
| `tag` | Tags a ModelKit |
| `version` | Display the version information for Kit |

## A Few Examples

To list your available ModelKit:

```sh
$ ./kit list
```

To build a ModelKit for your model and tag it with `example-tag`:

```sh
$ ./kit build ../examples/onnx -t localhost:5050/example-repo:example-tag"
```

Then you can push it to your registry:

```sh
$ ./kit push localhost:5050/example-repo:example-tag --http
```

After you finish calling all your friends and telling them about Kit, they will want to fetch your ModelKit and run it. The `fetch` command is used to bring everything in the ModelKit to your local machine - the model, dataset(s), code, and the [Kitfile](../kitfile/kf-overview.md) manifest.

```sh
$ ./kit fetch localhost:5050/test-repo:test-tag --http
```

However, Kit is a *modular package* so if someone only needs the model they can `pull` only that part:

```sh
$ ./kit pull -filter model
```

Or just the dataset:

```sh
$ ./kit pull -filter dataset
```

You can also use `pull` to filer for the `code` or the Kitfile `manifest`. When you pull any filtered part of a ModelKit you always get the Kitfile as well.

The `dev` command will automatically generate a RESTful API for the model and then run the model and API locally so anyone can use the model:

```sh
$ ./kit dev
```

So with a few easy commands you've been able to package up a model, share it with others, and run it locally...and you never needed to learn Dockerfile syntax or how to deal with a Helm chart or other proprietary packaging method.
