<script setup>
import vGaTrack from '@theme/directives/ga'
</script>

# Installing Kit

This page includes instructions for:

* Installing on [MacOS](#🍎-macos-package-download)
* Installing on [Windows](#🪟-windows-package-download)
* Installing on [Linux](#🐧-linux-package-download)
* Building from the [source code](#🛠️-install-from-source)

## 🍎 MacOS Install

There are two generations of Mac hardware, if you aren't sure which you have [check here](https://www.sweetwater.com/sweetcare/articles/intel-based-mac-or-mac-with-apple-silicon/#:~:text=Choose%20About%20This%20Mac.,a%20Mac%20with%20Apple%20silicon.&text=As%20of%20this%20writing%2C%20Apple,have%20an%20Intel%2Dbased%20Mac.).

1. <a href="https://github.com/jozu-ai/kitops/releases/latest/download/kitops-darwin-arm64.zip"
  v-ga-track="{
    category: 'link',
    label: 'MacOS (Apple Silicon)',
    location: 'docs/installation'
  }">
  Apple Silicon / ARM64
</a>

2. <a href="https://github.com/jozu-ai/kitops/releases/latest/download/kitops-darwin-x86_64.zip"
  v-ga-track="{
    category: 'link',
    label: 'MacOS (Intel)',
    location: 'docs/installation'
  }">
  Intel / x86_64
</a>

Once the file is downloaded, open your Finder and double-click the `.zip` file to unpack it. Then select the unpacked folder and move it to `/usr/local/bin`.

You can verify that `kit` is correctly installed by opening a new terminal or command prompt and typing:

```shell
kit version
```

This command should display the version number of `kit` you have installed, indicating that the installation was successful.

### Follow the Quick Start

Now that everything is set up you can follow our [Quick Start](https://kitops.ml/docs/quick-start.html) to learn how to pack and share your first ModelKit.

That's it!

## 🪟 Windows Install

Make sure you get the correct download for your hardware.

1. <a href="https://github.com/jozu-ai/kitops/releases/latest/download/kitops-windows-x86_64.zip"
  v-ga-track="{
    category: 'link',
    label: 'Windows (AMD64)',
    location: 'docs/installation'
  }">
  AMD 64-bit
</a>

2. <a href="https://github.com/jozu-ai/kitops/releases/latest/download/kitops-windows-arm64.zip"
  v-ga-track="{
    category: 'link',
    label: 'Windows (ARM64)',
    location: 'docs/installation'
  }">
  ARM 64-bit
</a>

3. <a href="https://github.com/jozu-ai/kitops/releases/latest/download/kitops-windows-i386.zip"
  v-ga-track="{
    category: 'link',
    label: 'Windows (x86_32)',
    location: 'docs/installation'
  }">
  Intel / AMD, 32-bit
</a>

Once downloaded, right-click on the `.zip` file and select "Extract All..." to unzip the archive. Now, Move the extracted `kit.exe` to a directory that is included in your system's PATH variable. This will allow you to run `kit` from any command prompt or terminal window.

You can verify that `kit` is correctly installed by opening a new terminal or command prompt and typing:

```shell
kit version
```

This command should display the version number of `kit` you have installed, indicating that the installation was successful.

### Follow the Quick Start

Now that everything is set up you can follow our [Quick Start](https://kitops.ml/docs/quick-start.html) to learn how to pack and share your first ModelKit.

That's it!

## 🐧 Linux Install

Make sure you get the correct download for your hardware.

1. <a href="https://github.com/jozu-ai/kitops/releases/latest/download/kitops-linux-arm64.tar.gz"
  v-ga-track="{
    category: 'link',
    label: 'Linux (ARM64)',
    location: 'docs/installation'
  }">
  ARM 64-bit
</a>

2. <a href="https://github.com/jozu-ai/kitops/releases/latest/download/kitops-linux-x86_64.tar.gz"
  v-ga-track="{
    category: 'link',
    label: 'Linux (AMD64)',
    location: 'docs/installation'
  }">
  AMD 64-bit
</a>

3. <a href="https://github.com/jozu-ai/kitops/releases/latest/download/kitops-linux-i386.tar.gz"
  v-ga-track="{
    category: 'link',
    label: 'Linux (x86_32)',
    location: 'docs/installation'
  }">
  Intel / AMD, 32-bit
</a>

Once downloaded, open a terminal window and use the `tar` command to extract the downloaded file. For example, if you downloaded the `kitops-linux-x86_64.tar.gz` file, you would use the following command:

```shell
tar -xzvf kitops-linux-x86_64.tar.gz
```

Move the extracted `kit` executable to a location in your system's PATH. A common choice is `/usr/local/bin`. You can do this with the `mv` command:

```
sudo mv kit /usr/local/bin/
```

This step may require administrator privileges.

After installation, you can verify that `kit` is correctly installed by opening a new terminal or command prompt and typing:

```shell
kit version
```

This command should display the version number of `kit` you have installed, indicating that the installation was successful.

### Follow the Quick Start

Now that everything is set up you can follow our [Quick Start](https://kitops.ml/docs/quick-start.html) to learn how to pack and share your first ModelKit.

## 🛠️ Build from Source Code

For those who prefer or require building `kit` from the source code, this section will guide you through the necessary steps.

Before you begin, make sure you have the following installed on your system:

- Git
- A recent version of Go

You can check if you have Go installed by running `go version` in your terminal. If you need to install Go, visit the [official Go download page](https://golang.org/dl/) for instructions.

### Clone the Repository

First, clone the `kitops` GitHub repository to your local machine. Open a terminal and run:


```shell
git clone https://github.com/jozu-ai/kitops.git
cd kitops
```

This command clones the repository and changes your current directory to the cloned repository's root.

### Build Sources

Once inside the `kitops` directory, you can build the `kit` tool using the Go compiler. Run:

```shell
go build -o kit
```

This command compiles the source code into an executable named `kit`. If you are on Windows, you might want to name the executable `kit.exe`.

### Install the Executable

After the build process completes, you need to move the `kit` executable to a location in your system's PATH to make it accessible from anywhere in the terminal:

#### For MacOS and Linux:

```shell
sudo mv kit /usr/local/bin/
```

#### For Windows:

Move `kit.exe` to a directory that's included in your system's PATH variable. This step may vary based on your specific Windows setup.

### Verify the Installation

To verify that `kit` was installed successfully, open a new terminal window and type:

```shell
kit version
```

### Follow the Quick Start

Now that everything is set up you can follow our [Quick Start](https://kitops.ml/docs/quick-start.html) to learn how to pack and share your first ModelKit.