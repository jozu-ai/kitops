## kit login

Log in to an OCI registry

### Synopsis

Log in to an OCI registry

```
kit login <registry> [flags]
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
      --config string   Config file (default $HOME/.kitops)
  -v, --verbose         Include additional information in output (default false)
```

### SEE ALSO

* [kit](kit.md)	 - Streamline the lifecycle of AI/ML models

