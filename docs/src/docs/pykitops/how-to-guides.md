---
description: Follow practical how-to guides for PyKitOps. Learn step-by-step methods for managing and deploying ModelKits in Python.
---
# How-to Guides

This part of the project documentation focuses on a **problem-oriented** approach. You'll tackle common tasks that you might have, with the help of the code provided in this project.

## How To Create A Kitfile Object?

Whether you're working with an existing ModelKit's Kitfile,
or starting from nothing, the `kitops` package can help you
get this done.

### Installation

Install the `kitops` package from PYPI into your project's environment
with the following command

```sh
pip install kitops
```

### Creating a Kitfile

There are two main ways to work with Kitfiles: creating from scratch or loading an existing one.

#### Loading an Existing Kitfile

Inside of your code you can now import the `Kitfile`
class from the `kitops.modelkit.kitfile` module:

```python
from kitops.modelkit.kitfile import Kitfile
```

After you've imported the class, you can use it
to create a Kitfile object from an existing ModelKit's Kitfile:

```python
from kitops.modelkit.kitfile import Kitfile

my_kitfile = Kitfile(path='/path/to/Kitfile')
print(my_kitfile.to_yaml())

# The output should match the contents of the Kitfile
# located at: /path/to/Kitfile
```

#### Creating a New Kitfile

You can also create an empty Kitfile from scratch:

```python
from kitops.modelkit.kitfile import Kitfile

my_kitfile = Kitfile()
print(my_kitfile.to_yaml())

# OUTPUT: {}
```

Regardless of how you created the Kitfile, you can update its contents
like you would do with any other python dictionary:

```python
from kitops.modelkit.kitfile import Kitfile

# Create new Kitfile
kitfile = Kitfile()

# Set basic metadata
kitfile.manifestVersion = "1.0"
kitfile.package = {
    "name": "sample-kitfile",
    "version": "1.0",
    "description": "Sample kitfile for PyKitOps demonstration"
}

# Configure model information
kitfile.model = {
    "name": "sample-model",
    "path": "model/model.pkl",
    "license": "Apache 2.0",
    "version": "1.0",
    "description": "Sample Model"
}

# Add code files
kitfile.code = [
    {
        "path": "demo.py",
        "description": "Sample model to demonstrate PyKitOps SDK",
        "license": "Apache 2.0"
    },
    {
        "path": "requirements.txt",
        "description": "Python dependencies"
    }
]

# Add datasets
kitfile.datasets = [
    {
        "name": "dataset",
        "path": "data/sample.csv",
        "description": "full dataset",
        "license": "Apache 2.0"
    }
]

# Add documentation
kitfile.docs = [
    {"path": "docs/README.md"},
    {"path": "docs/LICENSE"}
]

# OUTPUT:
# manifestVersion: '1.0'
# package:
#   name: sample-kitfile
#   version: '1.0'
#   description: Sample kitfile for PyKitOps demonstration
# code:
# - path: demo.py
#   description: Sample model to demonstrate PyKitOps SDK
#   license: Apache 2.0
# - path: requirements.txt
#   description: Python dependencies
# datasets:
# - name: dataset
#   path: data/sample.csv
#   description: full dataset
#   license: Apache 2.0
# docs:
# - path: docs/README.md
# - path: docs/LICENSE
# model:
#   name: sample-model
#   path: model/model.pkl
#   license: Apache 2.0
#   version: '1.0'
#   description: Sample Model
```

### Pushing to Jozu Hub

Once you've created your Kitfile, you can push it to Jozu Hub using the ModelKitManager. Here's how:

```python
from kitops.modelkit.manager import ModelKitManager

# Configure the ModelKit manager
modelkit_tag = "jozu.ml/yourname/reponame:latest"
manager = ModelKitManager(
    working_directory=".",
    modelkit_tag=modelkit_tag
)

# Assign your Kitfile
manager.kitfile = kitfile

# Pack and push to Jozu Hub
manager.pack_and_push_modelkit(save_kitfile=True)
```

### Complete Example

Here's a complete script that creates a Kitfile and pushes it to Jozu Hub:

```python
import os
from kitops.modelkit.kitfile import Kitfile
from kitops.modelkit.manager import ModelKitManager

if __name__ == "__main__":
    # Create the Kitfile
    kitfile = Kitfile()
    kitfile.manifestVersion = "1.0"
    kitfile.package = {
        "name": "sample-kitfile",
        "version": "1.0",
        "description": "Sample kitfile for PyKitOps demonstration"
    }
    
    kitfile.model = {
        "name": "sample-model",
        "path": "model/model.pkl",
        "license": "Apache 2.0",
        "version": "1.0",
        "description": "Sample Model"
    }
    
    kitfile.code = [
        {
            "path": "demo.py",
            "description": "Sample model to demonstrate PyKitOps SDK",
            "license": "Apache 2.0"
        },
        {
            "path": "requirements.txt",
            "description": "Python dependencies"
        }
    ]
    
    kitfile.datasets = [
        {
            "name": "dataset",
            "path": "data/sample.csv",
            "description": "full dataset",
            "license": "Apache 2.0"
        }
    ]
    
    kitfile.docs = [
        {"path": "docs/README.md"},
        {"path": "docs/LICENSE"}
    ]
    
    # Push to Jozu Hub
    modelkit_tag = "jozu.ml/yourname/reponame:latest"
    manager = ModelKitManager(
        working_directory=".",
        modelkit_tag=modelkit_tag
    )
    manager.kitfile = kitfile
    manager.pack_and_push_modelkit(save_kitfile=True)
    ```