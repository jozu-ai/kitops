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
| `login` | Login to the remote repository |
| `logout` | Logout to the remote repository |
| `list` | List ModelKits |
| `pull` | Pull one or more of the model, dataset, code, and Kitfile into a destination folder |
| `push` | Push ModelKit to respository |
| `remove` | Removed the ModelKit from the local repository |
| `tag` | Tags a ModelKit |
| `version` | Display the version information for Kit |

## Example

To list your available ModelKit:

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

After you finish calling all your friends and telling them about Kit they can fetch your ModelKit and run it:

```sh
$ ./kit fetch localhost:5050/test-repo:test-tag --http
$ ./kit dev
```

Maybe one of your friends only wants the dataset you used:

```sh
$ ./kit pull -filter dataset
```

Another friend only needs the model so they can integrate it with their application:

```sh
$ ./kit pull -filter model
```

To see the Kitfile associated with a ModelKit:

```sh
$ ./kit pull -filter kitfile
```
