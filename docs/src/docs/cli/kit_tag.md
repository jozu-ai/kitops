## kit tag

Create a tag that refers to a modelkit

### Synopsis

Create or update a tag `target-modelkit` that refers to `source-modelkit`

This command assigns a new tag to an existing modelkit (source-modelkit) or
updates an existing tag, effectively renaming or categorizing modelkits for
better organization and version control. Tags are identifiers linked to specific
modelkit versions within a repository.

A full modelkit reference has the following format:

[HOST[:PORT_NUMBER]/][NAMESPACE/]REPOSITORY[:TAG]

 * HOST: Optional. The registry hostname where the ModelKit is located. Defaults
   to localhost if unspecified. Must follow standard DNS rules
   (excluding underscores).

 * PORT_NUMBER: Optional. Specifies the registry's port number if a hostname is
   provided.

 * NAMESPACE: Represents a user or organization's namespace, consisting of
   slash-separated components that may include lowercase letters, digits, and
   specific separators (periods, underscores, hyphens).

 * REPOSITORY: The name of the repository, typically corresponding to the
   modelkit's name.

 * TAG: A human-readable identifier for the modelkit version or variant. Valid
   ASCII characters include lowercase and uppercase letters, digits, underscores,
   periods, and hyphens. It cannot start with a period or hyphen and is limited
   to 128 characters.

Tagging is a powerful way to manage different versions or configurations of your
modelkits, making it easier to organize, retrieve, and deploy specific
iterations. Ensure tags are meaningful and consistent across your team or
organization to maintain clarity and avoid confusion.

```
kit tag <source-modelkit>[:TAG] <target-modelkit>[:TAG] [flags]
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
      --config string   Config file (default $HOME/.kitops)
  -v, --verbose         Include additional information in output (default false)
```

### SEE ALSO

* [kit](kit.md)	 - Streamline the lifecycle of AI/ML models

