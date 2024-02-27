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
| `build` | Build a model |
| `completion` | Generate the autocompletion script for the specified shell |
| `dev` | ??? |
| `export` | Extract only the model, dataset, code, or Kitfile from a ModelKit |
| `help` | Help about any command |
| `login` | ??? |
| `list` | List model kits |
| `pull` | Pull model from registry |
| `push` | Push model to registry |
| `version` | Display the version information for Kit |

## Example

To list your available model kits:

```sh
$ ./kit list
```

To build a ModelKit for your model:

```sh
$ ./kit build ../examples/onnx -t localhost:5050/example-repo:example-tag"
```

Then you can push it to your registry:

```sh
$ ./kit push localhost:5050/example-repo:example-tag --http
```

After you finish calling all your friends and telling them about Kit they can pull your model and run it:

```sh
$ ./kit pull localhost:5050/test-repo:test-tag --http
$ ./kit dev
```

Maybe one of your friends only wants the dataset you used:

```sh
$ ./kit export dataset
```

Another friend only needs the model so they can integrate it with their application:

```sh
$ ./kit export model
```

To see the Kitfile associated with a ModelKit:

```sh
$ ./kit export kitfile
```
