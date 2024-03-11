# Installation Guide for `kit` CLI

## Installation from GitHub Releases
Welcome to the installation guide for the `kit`! This guide is designed to help you quickly and easily install the `kit` on your machine.

### Step 1: Downloading the `kit`

To begin, you will need to download the latest version of `kit`. You can find the most recent release on the official GitHub releases page:

[Download the latest `kit` release](https://github.com/jozu-ai/kitops/releases/latest)

Or, if you're feeling frisky try the cutting edge [`next` release](https://github.com/jozu-ai/kitops/releases/latest). Just remember that's not fully tested so YMMV.

#### Selecting the Correct Version for Your Platform

Depending on your operating system and its architecture, you will need to download a specific build of `kit`. Below is a table to help you identify the correct file to download for your platform:

| Platform                          | Release File Name               |
|-----------------------------------|---------------------------------|
| macOS (Apple Silicon, ARM64)      | `kitops-darwin-arm64.tar.gz`    |
| macOS (Intel, x86_64)             | `kitops-darwin-x86_64.tar.gz`   |
| Linux (ARM64)                     | `kitops-linux-arm64.tar.gz`     |
| Linux (AMD64/x86_64)              | `kitops-linux-x86_64.tar.gz`    |
| Linux (Intel/AMD, 32-bit)         | `kitops-linux-i386.tar.gz`      |
| Windows (AMD64/x86_64)            | `kitops-windows-x86_64.zip`     |
| Windows (ARM64)                   | `kitops-windows-arm64.zip`      |
| Windows (Intel/AMD, 32-bit)       | `kitops-windows-i386.zip`       |


### Optional: Verifying the Checksum

After downloading the `kit` and before proceeding with the installation, it's highly recommended to verify the checksum of the downloaded file. By verifying the checksum, you can ensure the authenticity and integrity of your downloaded `kit cli` file before installation.

Each release comes with a file that ends with `checksum.txt` that contains the SHA-256 hashes of the release files. Here's how to verify the checksum:

#### For macOS and Linux:

1. Open a terminal window.

2. Navigate to the directory where the downloaded file and the `checksum.txt` file are located.

3. Run the `sha256sum` command followed by the name of the downloaded file. For example, if you downloaded `kitops-linux-x86_64.tar.gz`, you would run:
```shell
shasum -a 256 kitops-darwin-arm64.tar.gz
```
3. Compare the output of this command with the corresponding checksum found in the `checksum.txt` file. If the checksums match, the file is verified and safe to use.

#### For Windows Users:

1. Open Command Prompt.
2. Use the `CertUtil` utility to generate the SHA-256 checksum of the downloaded file. For example, if you downloaded `kitops-windows-x86_64.zip`, run:
```shell
CertUtil -hashfile kitops-windows-x86_64.zip SHA256
```
3. Compare the output of this command with the corresponding checksum in the `checksum.txt` file. A matching checksum confirms the file's integrity.

### Step 2: Installing `kit`

Once you have downloaded the appropriate file for your system, follow these instructions to install the `kit`.

#### For macOS and Linux Users:

1. **Extract the Archive**: Open a terminal window and use the `tar` command to extract the downloaded file. For example, if you downloaded the `kitops-linux-x86_64.tar.gz` file, you would use the following command:

```shell
tar -xzvf kitops-linux-x86_64.tar.gz
```

2. **Move to Path**: Move the extracted `kit` executable to a location in your system's PATH. A common choice is `/usr/local/bin`. You can do this with the `mv` command:

```
sudo mv kit /usr/local/bin/
```
This step may require administrator privileges.

#### For Windows Users:

1. **Extract the Archive**: Right-click on the downloaded `.zip` file and select "Extract All..." to unzip the archive.

2. **Move to Path**: Move the extracted `kit.exe` to a directory that is included in your system's PATH variable. This will allow you to run `kit` from any command prompt or terminal window.


### Verifying the Installation

After installation, you can verify that `kit` is correctly installed by opening a new terminal or command prompt and typing:

```shell
kit version
```

This command should display the version number of `kit` you have installed, indicating that the installation was successful.


## Installation from Source

For those who prefer or require building `kit` from the source code, this section will guide you through the necessary steps.

### Prerequisites

Before you begin, make sure you have the following installed on your system:

- Git
- A recent version of Go

You can check if you have Go installed by running `go version` in your terminal. If you need to install Go, visit the [official Go download page](https://golang.org/dl/) for instructions.

### Step 1: Clone the Repository

First, clone the `kitops` GitHub repository to your local machine. Open a terminal and run:


```shell
git clone https://github.com/jozu-ai/kitops.git
cd kitops
```

This command clones the repository and changes your current directory to the cloned repository's root.

### Step 2: Build from Source

Once inside the `kitops` directory, you can build the `kit` tool using the Go compiler. Run:

```shell
go build -o kit
```

This command compiles the source code into an executable named `kit`. If you are on Windows, you might want to name the executable `kit.exe`.

### Step 3: Install the Executable

After the build process completes, you need to move the `kit` executable to a location in your system's PATH to make it accessible from anywhere in the terminal:

#### For macOS and Linux:

```shell
sudo mv kit /usr/local/bin/
```

#### For Windows:

Move `kit.exe` to a directory that's included in your system's PATH variable. This step may vary based on your specific Windows setup.

### Verifying the Installation

To verify that `kit` was installed successfully, open a new terminal window and type:

```shell
kit version
```
