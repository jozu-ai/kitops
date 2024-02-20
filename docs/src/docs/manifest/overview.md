# Manifest

## Atamafiles

Atama uses what we call an `atamafile` as its manifest. An `atamafile` is like a little sandbox for your code, and it comes in the form of a `yaml` or `json` file, like `atama.yaml` or `atama.json`, it can even be prefixed with your own name like `myfile.atama.yaml`. This also means that an atamafile doesn't know anything about your computer or files. It runs on the given environment. Atamafiles have everything your code needs to run, even a basic operating system.

In this walkthrough, you'll view and explore an actual atamafile.

Before you start, get the latest version of [atama cli](/docs/cli/installation). Atama adds new features regularly and some parts of this guide may work only with the latest version of Atama Desktop.

The Atamafile supports the following instructions:

| Instruction |	Description |
| ---- | --- |
| `ADD` |	Add local or remote files and directories. |
| `ARG` |	Use build-time variables. |
| `CMD` |	Specify default commands. |
| `COPY` |	Copy files and directories. |
| `ENTRYPOINT` | Specify default executable. |
| `ENV` |	Set environment variables. |
| `EXPOSE` |	Describe which ports your application is listening on. |
| `FROM` |	Create a new build stage from a base image. |
| `HEALTHCHECK` |	Check a container's health on startup. |
| `LABEL` |	Add metadata to an image. |
| `MAINTAINER` |	Specify the author of an image. |
| `ONBUILD` |	Specify instructions for when the image is used in a build. |
| `RUN` |	Execute build commands. |
| `SHELL` |	Set the default shell of an image. |
| `STOPSIGNAL` | Specify the system call signal for exiting a container. |
| `USER` | Set user and group ID. |
| `VOLUME` | Create volume mounts. |
| `WORKDIR` |	Change working directory. |

## Format
The format of the atamafile is:

```yaml
# Comment
INSTRUCTION arguments
```

The instruction is not case-sensitive. However, convention is for them to be UPPERCASE to distinguish them from arguments more easily.

Atama runs instructions in an `atamafile` in order. An atamafile must begin with a FROM instruction. This may be after parser directives, comments, and globally scoped ARGs. The FROM instruction specifies the parent image from which you are building. FROM may only be preceded by one or more ARG instructions, which declare arguments that are used in FROM lines in the Atamafile.

:::info Note on whitespace
For backward compatibility, leading whitespace before comments (`#`) and instructions (such as `RUN`) are ignored, but discouraged. Leading whitespace is not preserved in these cases, and the following examples are therefore equivalent:


```sh
        # this is a comment-line
    RUN echo hello
RUN echo world
```

&NewLine;

```sh
# this is a comment-line
RUN echo hello
RUN echo world
````

Whitespace in instruction arguments, however, isn't ignored. The following example prints hello world with leading whitespace as specified:

```sh
RUN echo "\
     hello\
     world"
```
:::

## Step 1: Set up the walkthrough

The first thing you need is a running Atamafile. Use the following instructions to run an atamafile.

## Step 2: Setting up environment variables

Occaecati repudiandae aliquid nostrum et dolores repellat. Vero corporis ducimus placeat. Deserunt animi alias recusandae in eos aut facere sed praesentium. In eos illo ab tenetur et. Cumque dolor dolorem aut et ea nostrum eos ratione.

Vel enim impedit in hic ut est sit aut. Suscipit recusandae et delectus quasi aut sapiente ratione. Totam deserunt sunt nemo. Aperiam iste inventore velit. Officiis quas saepe hic quisquam optio rerum non voluptatibus.

Blanditiis non numquam mollitia. Placeat impedit earum. Non optio quibusdam autem veritatis rerum omnis. Labore cum autem tempora.

## Step 3: Manage Dependencies

Eaque et soluta. Eveniet dignissimos modi quasi adipisci nesciunt et iure in eligendi. Perferendis inventore quod placeat ut nostrum occaecati. Quaerat voluptatem ipsum quo et nobis error laboriosam.

Eveniet ut sunt eius ea atque saepe. Est omnis qui nihil quam dolor. Illum ut excepturi dolorum possimus sed quis asperiores officiis. Qui molestias rerum et quasi eveniet modi. Id et facere quia. Illo possimus officiis reprehenderit adipisci odit dolor praesentium debitis.

Tenetur eum quos. Reprehenderit est nulla quod autem et officia quasi. Facere dolorem minus nihil. Commodi qui est dolor quia soluta soluta esse. Assumenda et hic modi sequi numquam maiores. Quis aut exercitationem earum occaecati sed placeat.

## Step 4: Best practices and tips

Dolorem officiis laborum voluptatem doloribus id suscipit. Non officiis modi dolorem eum quis aut numquam distinctio. Est vero blanditiis vel aliquam dolor. Rerum pariatur hic libero. Esse fuga est dolores dolor provident nobis quis quae eius. Ex consectetur et blanditiis et eveniet aliquid autem.

Debitis quia quidem omnis quam nihil est recusandae. Odio vitae et deserunt ullam rerum et. Quae quisquam omnis aliquam.

Saepe consectetur sunt quidem dolorem est et vero. Saepe neque sed quo. Facere commodi explicabo. Perferendis et mollitia consectetur ipsum sint voluptatem accusantium dignissimos. Laborum facilis ut consequatur. Laborum consequatur omnis dolorem ipsam officia omnis explicabo ab.
