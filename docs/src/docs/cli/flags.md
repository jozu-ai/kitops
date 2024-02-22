# Flags

`kit cli` supports global and local flags.

## Global flags

These flags must be run with no command.

| Flag |	Description |
| ---- | ---- |
| `--config <string>` | config file (default is `$HOME/.jozu`) |
| `-h`, `--help` | help for kit |

### Example

<PlatformSnippet>
  <template #windows>

  ```sh
  ./kit --config c:\Applications\caches\.jozucache
  ```

  </template>

  <template #linux>

  ```sh
  ./kit --config /var/usr/.jozucache
  ```

  </template>

  <template #mac>

  ```sh
  ./kit --config ~/Library/.jozucache
  ```

  </template>
</PlatformSnippet>

## Local flags

These flags run next to a command.

| Flag |	Description |
| ---- | ---- |
| `-h`, `--help` | help for kit |

### Example

`./kit list --help`

or

`./kit list -h`
