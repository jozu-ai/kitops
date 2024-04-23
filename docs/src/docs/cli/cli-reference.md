---
outline: 2
---
# Kit CLI Reference

## kit info

Show the configuration for a modelkit

### Synopsis

Print the contents of a modelkit config to the screen.

By default, kit will check local storage for the specified modelkit. To see
the configuration for a modelkit stored on a remote registry, use the
--remote flag.

```
kit info [flags] MODELKIT
```

### Examples

```
# See configuration for a local modelkit:
kit info mymodel:mytag

# See configuration for a local modelkit by digest:
kit info mymodel@sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a

# See configuration for a remote modelkit:
kit info --remote registry.example.com/my-model:1.0.0
```

### Options

```
  -h, --help         help for info
      --plain-http   Use plain HTTP when connecting to remote registries
  -r, --remote       Check remote registry instead of local storage
      --tls-verify   Require TLS and verify certificates when connecting to remote registries (default true)
```

### Options inherited from parent commands

```
      --config string     Alternate path to root storage directory for CLI
      --progress string   Configure progress bars for longer operations (options: none, plain, fancy) (default "plain")
  -v, --verbose           Include additional information in output (default false)
```

## kit inspect

Inspect a modelkit's manifest

### Synopsis

Print the contents of a modelkit manifest to the screen.

By default, kit will check local storage for the specified modelkit. To
inspect a modelkit stored on a remote registry, use the --remote flag.

```
kit inspect [flags] MODELKIT
```

### Examples

```
# Inspect a local modelkit:
kit inspect mymodel:mytag

# Inspect a local modelkit by digest:
kit inspect mymodel@sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a

# Inspect a remote modelkit:
kit inspect --remote registry.example.com/my-model:1.0.0
```

### Options

```
  -h, --help         help for inspect
      --plain-http   Use plain HTTP when connecting to remote registries
  -r, --remote       Check remote registry instead of local storage
      --tls-verify   Require TLS and verify certificates when connecting to remote registries (default true)
```

### Options inherited from parent commands

```
      --config string     Alternate path to root storage directory for CLI
      --progress string   Configure progress bars for longer operations (options: none, plain, fancy) (default "plain")
  -v, --verbose           Include additional information in output (default false)
```

## kit list

List modelkits in a repository

### Synopsis

Displays a list of modelkits available in a repository.

This command provides an overview of modelkits stored either in the local
repository or a specified remote repository. It displays each modelkit along
with its associated tags and the cumulative size of its contents. Modelkits
comprise multiple artifacts, including models, datasets, code, and
configuration, designed to enhance reusability and modularity. However, this
command focuses on the aggregate rather than listing individual artifacts.

Each modelkit entry includes its DIGEST, a unique identifier that ensures
distinct versions of a modelkit are easily recognizable, even if they share
the same name or tags. Modelkits with multiple tags or repository names will
appear multiple times in the list, distinguished by their DIGEST.

The SIZE displayed for each modelkit represents the total storage space
occupied by all its components.

```
kit list [flags] [REPOSITORY]
```

### Examples

```
# List local modelkits
kit list

# List modelkits from a remote repository
kit list registry.example.com/my-namespace/my-model
```

### Options

```
  -h, --help         help for list
      --plain-http   Use plain HTTP when connecting to remote registries
      --tls-verify   Require TLS and verify certificates when connecting to remote registries (default true)
```

### Options inherited from parent commands

```
      --config string     Alternate path to root storage directory for CLI
      --progress string   Configure progress bars for longer operations (options: none, plain, fancy) (default "plain")
  -v, --verbose           Include additional information in output (default false)
```

## kit login

Log in to an OCI registry

### Synopsis

Log in to a specified OCI-compatible registry. Credentials are saved and used
automatically for future CLI operations

```
kit login [flags] [REGISTRY]
```

### Examples

```
# Login to ghcr.io
kit login ghcr.io -u github_user -p personal_token
```

### Options

```
  -h, --help              help for login
  -p, --password string   registry password or token
      --password-stdin    read password from stdin
      --plain-http        Use plain HTTP when connecting to remote registries
      --tls-verify        Require TLS and verify certificates when connecting to remote registries (default true)
  -u, --username string   registry username
```

### Options inherited from parent commands

```
      --config string     Alternate path to root storage directory for CLI
      --progress string   Configure progress bars for longer operations (options: none, plain, fancy) (default "plain")
  -v, --verbose           Include additional information in output (default false)
```

## kit logout

Log out from an OCI registry

### Synopsis

Log out from a specified OCI-compatible registry. Any saved credentials are
removed from storage.

```
kit logout [flags] REGISTRY
```

### Examples

```
# Log out from ghcr.io
kit logout ghcr.io
```

### Options

```
  -h, --help   help for logout
```

### Options inherited from parent commands

```
      --config string     Alternate path to root storage directory for CLI
      --progress string   Configure progress bars for longer operations (options: none, plain, fancy) (default "plain")
  -v, --verbose           Include additional information in output (default false)
```

## kit pack

Pack a modelkit

### Synopsis

Pack a modelkit from a kitfile using the given context directory.

The packing process involves taking the configuration and resources defined in
your kitfile and using them to create a modelkit. This modelkit is then stored
in your local registry, making it readily available for further actions such
as pushing to a remote registry for collaboration.

Unless a different location is specified, this command looks for the kitfile
at the root of the provided context directory. Any relative paths defined
within the kitfile are interpreted as being relative to this context
directory.

```
kit pack [flags] DIRECTORY
```

### Examples

```
# Pack a modelkit using the kitfile in the current directory
kit pack .

# Pack a modelkit with a specific kitfile and tag
kit pack . -f /path/to/your/Kitfile -t registry/repository:modelv1
```

### Options

```
  -f, --file string   Specifies the path to the Kitfile explictly (use "-" to read from standard input)
  -h, --help          help for pack
  -t, --tag string    Assigns one or more tags to the built modelkit. Example: -t registry/repository:tag1,tag2
```

### Options inherited from parent commands

```
      --config string     Alternate path to root storage directory for CLI
      --progress string   Configure progress bars for longer operations (options: none, plain, fancy) (default "plain")
  -v, --verbose           Include additional information in output (default false)
```

## kit pull

Retrieve modelkits from a remote registry to your local environment.

### Synopsis

Downloads modelkits from a specified registry. The downloaded modelkits
are stored in the local registry.

```
kit pull [flags] registry/repository[:tag|@digest]
```

### Examples

```
# Pull the latest version of a modelkit from a remote registry
kit pull registry.example.com/my-model:latest
```

### Options

```
  -h, --help         help for pull
      --plain-http   Use plain HTTP when connecting to remote registries
      --tls-verify   Require TLS and verify certificates when connecting to remote registries (default true)
```

### Options inherited from parent commands

```
      --config string     Alternate path to root storage directory for CLI
      --progress string   Configure progress bars for longer operations (options: none, plain, fancy) (default "plain")
  -v, --verbose           Include additional information in output (default false)
```

## kit push

Upload a modelkit to a specified registry

### Synopsis

This command pushes modelkits to a remote registry.

The modelkits should be tagged with the target registry and repository before
they can be pushed

```
kit push [flags] registry/repository[:tag|@digest]
```

### Examples

```
# Push the latest modelkits to a remote registry
kit push registry.example.com/my-model:latest

# Push a specific version of a modelkits using a tag:
kit push registry.example.com/my-model:1.0.0
```

### Options

```
  -h, --help         help for push
      --plain-http   Use plain HTTP when connecting to remote registries
      --tls-verify   Require TLS and verify certificates when connecting to remote registries (default true)
```

### Options inherited from parent commands

```
      --config string     Alternate path to root storage directory for CLI
      --progress string   Configure progress bars for longer operations (options: none, plain, fancy) (default "plain")
  -v, --verbose           Include additional information in output (default false)
```

## kit remove

Remove a modelkit from local storage

### Synopsis

Removes a modelkit from storage on the local disk.

The model to be removed may be specifed either by a tag or by a digest. If
specified by digest, that modelkit will be removed along with any tags that
might refer to it. If specified by tag (and the --force flag is not used),
the modelkit will only be removed if no other tags refer to it; otherwise
it is only untagged.

```
kit remove [flags] registry/repository[:tag|@digest]
```

### Examples

```
# Remove modelkit by tag
kit remove my-registry.com/my-org/my-repo:my-tag

# Remove modelkit by digest
kit remove my-registry.com/my-org/my-repo@sha256:<digest>

# Remove multiple tags for a modelkit
kit remove my-registry.com/my-org/my-repo:tag1,tag2,tag3

# Remove all untagged modelkits
kit remove --all

# Remove all locally stored modelkits
kit remove --all --force
```

### Options

```
  -a, --all     remove all untagged modelkits
  -f, --force   remove modelkit and all other tags that refer to it
  -h, --help    help for remove
```

### Options inherited from parent commands

```
      --config string     Alternate path to root storage directory for CLI
      --progress string   Configure progress bars for longer operations (options: none, plain, fancy) (default "plain")
  -v, --verbose           Include additional information in output (default false)
```

## kit tag

Create a tag that refers to a modelkit

### Synopsis

Create or update a tag {target-modelkit} that refers to {source-modelkit}

This command assigns a new tag to an existing modelkit (source-modelkit) or
updates an existing tag, effectively renaming or categorizing modelkits for
better organization and version control. Tags are identifiers linked to
specific modelkit versions within a repository.

A full modelkit reference has the following format:

[HOST[:PORT_NUMBER]/][NAMESPACE/]REPOSITORY[:TAG]

  * HOST: Optional. The registry hostname where the ModelKit is located.
    Defaults to localhost if unspecified. Must follow standard DNS rules
    (excluding underscores).

  * PORT_NUMBER: Optional. Specifies the registry's port number if a hostname
    is provided.

  * NAMESPACE: Represents a user or organization's namespace, consisting of
    slash-separated components that may include lowercase letters, digits, and
    specific separators (periods, underscores, hyphens).

  * REPOSITORY: The name of the repository, typically corresponding to the
    modelkit's name.

  * TAG: A human-readable identifier for the modelkit version or variant.
    Valid ASCII characters include lowercase and uppercase letters, digits,
    underscores, periods, and hyphens. It cannot start with a period or hyphen
    and is limited to 128 characters.

Tagging is a powerful way to manage different versions or configurations of
your modelkits, making it easier to organize, retrieve, and deploy specific
iterations. Ensure tags are meaningful and consistent across your team or
organization to maintain clarity and avoid confusion.

```
kit tag SOURCE_MODELKIT[:TAG] TARGET_MODELKIT[:TAG] [flags]
```

### Examples

```
kit tag myregistry.com/myrepo/mykit:latest myregistry.com/myrepo/mykit:v1.0.0
```

### Options

```
  -h, --help   help for tag
```

### Options inherited from parent commands

```
      --config string     Alternate path to root storage directory for CLI
      --progress string   Configure progress bars for longer operations (options: none, plain, fancy) (default "plain")
  -v, --verbose           Include additional information in output (default false)
```

## kit unpack

Produce the components from a modelkit on the local filesystem

### Synopsis

Produces all or selected components of a modelkit on the local filesystem.

This command unpacks a modelkit's components, including models, code,
datasets, and configuration files, to a specified directory on the local
filesystem. By default, it attempts to find the modelkit in local storage; if
not found, it searches the remote registry and retrieves it. This process
ensures that the necessary components are always available for unpacking,
optimizing for efficiency by fetching only specified components from the
remote registry when necessary

```
kit unpack [flags] [registry/]repository[:tag|@digest]
```

### Examples

```
# Unpack all components of a modelkit to the current directory
kit unpack myrepo/my-model:latest -d /path/to/unpacked

# Unpack only the model and datasets of a modelkit to a specified directory
kit unpack myrepo/my-model:latest --model --datasets -d /path/to/unpacked

# Unpack a modelkit from a remote registry with overwrite enabled
kit unpack registry.example.com/myrepo/my-model:latest -o -d /path/to/unpacked
```

### Options

```
      --code         Unpack only code
      --datasets     Unpack only datasets
  -d, --dir string   The target directory to unpack components into. This directory will be created if it does not exist
  -h, --help         help for unpack
      --kitfile      Unpack only Kitfile
      --model        Unpack only model
  -o, --overwrite    Overwrites existing files and directories in the target unpack directory without prompting
      --plain-http   Use plain HTTP when connecting to remote registries
      --tls-verify   Require TLS and verify certificates when connecting to remote registries (default true)
```

### Options inherited from parent commands

```
      --config string     Alternate path to root storage directory for CLI
      --progress string   Configure progress bars for longer operations (options: none, plain, fancy) (default "plain")
  -v, --verbose           Include additional information in output (default false)
```

## kit version

Display the version information for the CLI

### Synopsis

The version command prints detailed version information.

This information includes the current version of the tool, the Git commit that
the version was built from, the build time, and the version of Go it was
compiled with.

```
kit version [flags]
```

### Options

```
  -h, --help   help for version
```

### Options inherited from parent commands

```
      --config string     Alternate path to root storage directory for CLI
      --progress string   Configure progress bars for longer operations (options: none, plain, fancy) (default "plain")
  -v, --verbose           Include additional information in output (default false)
```

