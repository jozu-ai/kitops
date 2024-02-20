# Flags

`atama` supports global and local flags.

## Global flags

These flags must be run with no command.

| Flag |	Description |
| ---- | ---- |
| `--config <string>` | config file (default is `$HOME/.jozu`) |
| `-h`, `--help` | help for atama |

### Example

<PlatformSnippet>
  <template #windows>

  ```sh
  ./atama --config c:\Applications\caches\.jozucache
  ```

  </template>

  <template #linux>

  ```sh
  ./atama --config /var/usr/.jozucache
  ```

  </template>

  <template #mac>

  ```sh
  ./atama --config ~/Library/.jozucache
  ```

  </template>
</PlatformSnippet>

## Local flags

These flags run next to a command.

| Flag |	Description |
| ---- | ---- |
| `-h`, `--help` | help for atama |

### Example

`./atama models --help`

or

`./atama models -h`
