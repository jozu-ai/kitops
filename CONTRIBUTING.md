# Contributing Guide

* [Ways to Contribute](#ways-to-contribute)
* [Development Environment Setup](#development-environment-setup)
* [Pull Request Lifecycle](#pull-request-lifecycle)
* [Sign Your Commits](#sign-your-commits)
* [Ask for Help](#ask-for-help)

Welcome! We are so excited that you want to contribute to our project! 💖

As you get started, you are in the best position to give us feedback on areas of our project that we need help with including:

* Problems found during setting up a new developer environment
* Gaps in our guides or documentation
* Bugs in our tools and automation scripts

If anything doesn't make sense, or doesn't work when you try it, please open a bug report and let us know!

## Ways to Contribute

We welcome many different types of contributions including:

* New features
* Builds and CI/CD
* Bug fixes
* Documentation
* Issue triage
* Answering questions on Discord, or the mailing list
* Web design
* Communications, social media, blog posts, or other marketing
* Release management

Not everything happens through a GitHub pull request. Please contact us in the [#general channel of our Discord server](https://discord.gg/XzSmtPn3) or during our [office hours meeting](./GOVERNANCE.md#meetings) and let's discuss how we can work together.

## Development Environment Setup

### Prerequisites
- Go (Golang): The latest version of Go, as the project is written in this language. Go's installation guide can be found at https://golang.org/doc/install.
- Git: Version control system for cloning the repository and managing code changes. Installation instructions are available at https://git-scm.com/book/en/v2/Getting-Started-Installing-Git.


### Setting up the project

1. Clone the Repository: Clone the KitOps source code to your local machine:

```shell
git clone https://github.com/jozu-ai/kitops.git
cd kitops
```

2. Install Go Dependencies: Inside the project directory, fetch and install the project's dependencies using the Go command:

```shell
go mod tidy
```

3. Build the Kit CLI: Compile the source code into an executable named kit:
```shell
go build -o kit
```

4. Run Your Build: Execute the built CLI to see all available commands:

```shell
go run .
```
5. Updating Dependencies: If you add or update dependencies, ensure to update the go.mod and go.sum files by running `go mod tidy` again and include these changes in your commits.


## Pull Request Lifecycle

Pull requests are often called a "PR". KitOps generally follows the standard [GitHub pull request process](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/about-pull-requests).

## Code Reviews

There are two aspects of code review: giving and receiving. To make it easier for your PR to receive reviews, consider that reviewers will need you to:

* Follow the project coding conventions
* Write good commit messages
* Break large changes into a logical series of smaller patches which individually make easily understandable changes, and in aggregate solve a broader issue

Reviewers, the people giving the review, are highly encouraged to revisit our [Code of Conduct](./CODE-OF-CONDUCT.md) and must go above and beyond to promote a collaborative, respectful community. When reviewing PRs from others [The Gentle Art of Patch Review](https://sage.thesharps.us/2014/09/01/the-gentle-art-of-patch-review/) suggests an iterative series of focuses which is designed to lead new contributors to positive collaboration without inundating them initially with a pile of suggestions:

1. Is the idea behind the contribution sound?
1. Is the contribution architected correctly?
1. Is the contribution polished?

If your PR isn't getting any attention after 3-4 days (remember Maintainers tend to be very busy) please ping one of the Maintainers in Discord.

## Sign Your Commits

Licensing is important to open source projects. It provides some assurances that the software will continue to be available based under the terms that the author(s) desired. We require that contributors sign off on commits submitted to our project's repositories. The [Developer Certificate of Origin (DCO)](https://probot.github.io/apps/dco/) is a way to certify that you wrote and have the right to contribute the code you are submitting to the project.

You sign-off by adding the following to your commit messages. Your sign-off must match the git user and email associated with the commit. Your commit message should be followed by:

    Signed-off-by: Your Name <your.name@example.com>

Git has a `-s` command line option to do this automatically:

    git commit -s -m 'This is my commit message'

If you forgot to do this and have not yet pushed your changes to the remote
repository, you can amend your commit with the sign-off by running

    git commit --amend -s

## Ask for Help

The best way to reach us with a question when contributing is to ask on:

* The original github issue
* Our Discord channel (TODO)
* At our [office hours meeting](./GOVERNANCE.md)
