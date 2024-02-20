# Markdown Extension Examples

This page demonstrates some of the built-in markdown extensions provided by VitePress.

## Syntax Highlighting

VitePress provides Syntax Highlighting powered by [Shiki](https://github.com/shikijs/shiki), with additional features like line-highlighting:

**Input**

````md
```js{4}
export default {
  data () {
    return {
      msg: 'Highlighted!'
    }
  }
}
```
````

**Output**

```js{4}
export default {
  data () {
    return {
      msg: 'Highlighted!'
    }
  }
}
```

## Dynamic platform code snippet

Dynamic snippets that depends on the selected platform needs to be wrapped into a `PlatformSnippet` tag, which is a custom-made Vue component. Then each language snippet must be inside a `<template #<language>>` tag.

### Examples:

````html
<PlatformSnippet>
  <template #windows>

  ```bat
  $ atama run --win atama.yaml
  ```

  </template>

  <template #mac>

  ```sh
  $ atama run --mac atama.yaml
  ```

  </template>

  <template #linux>

  ```sh
  $ atama run --linux atama.yaml
  ```

  </template>
</PlatformSnippet>
````

:::info Important to know
An extra line is required between `<template>` and the \`\`\` codeblock. And the code block must start at _4 spaces indentation_ at most or it will render as text instead of code.
:::



### Results:

<PlatformSnippet>
  <template #windows>

  ```bat
  $ atama run --win myfile.atamafile
  ```

  </template>

  <template #mac>

  ```sh
  $ atama run --mac myfile.atamafile
  ```

  </template>

  <template #linux>

  ```sh
  $ atama run --linux myfile.atamafile
  ```

  </template>
</PlatformSnippet>


## Custom Containers

**Input**

```md
::: info
This is an info box.
:::

::: tip
This is a tip.
:::

::: warning
This is a warning.
:::

::: danger
This is a dangerous warning.
:::

::: details
This is a details block.
:::
```

**Output**

::: info
This is an info box.
:::

::: tip
This is a tip.
:::

::: warning
This is a warning.
:::

::: danger
This is a dangerous warning.
:::

::: details
This is a details block.
:::

## More

Check out the documentation for the [full list of markdown extensions](https://vitepress.dev/guide/markdown).
