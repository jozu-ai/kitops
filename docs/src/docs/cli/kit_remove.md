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
kit remove registry/repository[:tag|@digest] [flags]
```

### Examples

```
kit remove my-registry.com/my-org/my-repo:my-tag
kit remove my-registry.com/my-org/my-repo@sha256:<digest>
kit remove my-registry.com/my-org/my-repo:tag1,tag2,tag3
```

### Options

```
  -f, --force   remove manifest even if other tags refer to it
  -h, --help    help for remove
```

### Options inherited from parent commands

```
      --config string   Config file (default $HOME/.kitops)
  -v, --verbose         Include additional information in output (default false)
```

### SEE ALSO

* [kit](kit.md)	 - Streamline the lifecycle of AI/ML models

