## How to Run the Demos

Follow these steps to run the demo scripts successfully:

1. **Open Your Terminal:** Start by opening a terminal window on your computer.

2. **Navigate to the Demo Directory:** Change your current working directory to the demo's directory by using the `cd` command followed by the path to the directory. For example:

```shell 
cd  kitops/demo/
```


3. **Download the CLI Version:** Obtain the version of the Command Line Interface (CLI) tool that is required for the demo. You can get the [nightly](https://github.com/jozu-ai/kitops/releases/tag/nightly) builds. Once downloaded, move the CLI tool into the demo directory.Rename the CLI and remove any platform specifications from the name. Ensure that the CLI tool has the correct permissions to be executed. If necessary, you can change the permissions by running:

```shell
chmod +x ./kitops
```
4. **Run the Demo Script:** Execute the demo script by typing the name of the demo you want to execute to your terminal:

```shell
./my-demo.sh
```

Make sure that `my-demo.sh` is the correct name of the demo script you intend to run. This command assumes that `my-demo.sh` is located in the current directory and has execution permissions. If the script does not have execution permissions, grant them by running:

```shell
chmod +x my-demo.sh
```

Before running these steps, ensure that you have any dependencies or environmental requirements satisfied which are usually recorded as comments on the `.sh` file. 
