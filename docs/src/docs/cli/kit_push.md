## kit push

Uploads modelkits to a specified registry

### Synopsis

This command pushes modelkits to a remote registry.

The modelkits should be tagged with the target registry and repository before 
they can be pushed

```
kit push registry/repository[:tag|@digest] [flags]
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
      --config string   Config file (default $HOME/.kitops)
  -v, --verbose         Include additional information in output (default false)
```

### SEE ALSO

* [kit](kit.md)	 - Streamline the lifecycle of AI/ML models

