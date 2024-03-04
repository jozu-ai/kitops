## kit export

Produce the components from a modelkit on the local filesystem

### Synopsis

Produces all or selected components of a modelkit on the local filesystem.
	
This command exports a modelkit's components, including models, code, datasets, 
and configuration files, to a specified directory on the local filesystem. 
By default, it attempts to find the modelkit in local storage; if not found, it 
searches the remote registry and retrieves it. This process ensures that the 
necessary components are always available for export, optimizing for efficiency 
by fetching only specified components from the remote registry when necessary

```
kit export [registry/]repository[:tag|@digest] [flags]
```

### Examples

```
# Export all components of a modelkit to the current directory
kit export myrepo/my-model:latest -d /path/to/export

# Export only the model and datasets of a modelkit to a specified directory
kit export myrepo/my-model:latest --model --datasets -d /path/to/export

# Export a modelkit from a remote registry with overwrite enabled
kit export registry.example.com/myrepo/my-model:latest -o -d /path/to/export
```

### Options

```
      --code         Export only code
      --config       Export only config file
      --datasets     Export only datasets
  -d, --dir string   The target directory to export components into. This directory will be created if it does not exist
  -h, --help         help for export
      --model        Export only model
  -o, --overwrite    Overwrites existing files and directories in the target export directory without prompting
      --plain-http   Use plain HTTP when connecting to remote registries
      --tls-verify   Require TLS and verify certificates when connecting to remote registries (default true)
```

### Options inherited from parent commands

```
  -v, --verbose   Include additional information in output (default false)
```

### SEE ALSO

* [kit](kit.md)	 - Streamline the lifecycle of AI/ML models

