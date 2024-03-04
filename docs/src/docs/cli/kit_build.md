## kit build

Builds a modelkit

### Synopsis

Build a modelkit from a kitfile using the given context directory. 

The build process involves taking the configuration and resources defined in 
your kitfile and using them to create a modelkit. This modelkit is then stored
in your local registry, making it readily available for further actions such 
as pushing to a remote registry for collaboration.

Unless a different location is specified, this command looks for the k	itfile 
at the root of the provided context directory. Any relative paths defined 
within the kitfile are interpreted as being relative to this context directory.

```
kit build DIRECTORY [flags]
```

### Examples

```
# Build a modelkit using the kitfile in the current directory
kit build .

# Build a modelkit with a specific kitfile and tag
kit build . -f /path/to/your/Kitfile -t registry/repository:modelv1
```

### Options

```
  -f, --file string   Specifies the path to the Kitfile if it's not located at the root of the context directory
  -h, --help          help for build
  -t, --tag string    Assigns one or more tags to the built modelkit. Example: -t registry/repository:tag1,tag2
```

### Options inherited from parent commands

```
      --config string   Config file (default $HOME/.kitops)
  -v, --verbose         Include additional information in output (default false)
```

### SEE ALSO

* [kit](kit.md)	 - Streamline the lifecycle of AI/ML models

