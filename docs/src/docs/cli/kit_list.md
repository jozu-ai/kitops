## kit list

List modelkits in a repository

### Synopsis

Displays a list of modelkits available in a repository.

This command provides an overview of modelkits stored either in the local 
repository or a specified remote repository. It displays each modelkit along 
with its associated tags and the cumulative size of its contents. Modelkits 
comprise multiple artifacts, including models, datasets, code, and configuration,
designed to enhance reusability and modularity. However, this command focuses on
the aggregate rather than listing individual artifacts.

Each modelkit entry includes its DIGEST, a unique identifier that ensures
distinct versions of a modelkit are easily recognizable, even if they share the 
same name or tags. Modelkits with multiple tags or repository names will appear 
multiple times in the list, distinguished by their DIGEST.

The SIZE displayed for each modelkit represents the total storage space occupied
by all its components.

```
kit list [registry/repository] [flags]
```

### Examples

```
## List local modelkits
kit list

# List modelkits from a remote repository
kit list registry.example.com/my-model
```

### Options

```
  -h, --help         help for list
      --plain-http   Use plain HTTP when connecting to remote registries
      --tls-verify   Require TLS and verify certificates when connecting to remote registries (default true)
```

### Options inherited from parent commands

```
      --config string   Config file (default $HOME/.kitops)
  -v, --verbose         Include additional information in output (default false)
```

### SEE ALSO

* [kit](kit.md)	 - Streamline the lifecycle of AI/ML models

