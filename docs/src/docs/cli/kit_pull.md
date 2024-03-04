## kit pull

Retrieve modelkits from a remote registry to your local environment.

### Synopsis

Downloads modelkits from a specified registry. The downloaded 
modelkits are stored in the local registry.

```
kit pull registry/repository[:tag|@digest] [flags]
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
      --config string   Config file (default $HOME/.kitops)
  -v, --verbose         Include additional information in output (default false)
```

### SEE ALSO

* [kit](kit.md)	 - Streamline the lifecycle of AI/ML models

