# Jozu AI/ML Packaging Manifest Format Reference

The Jozu manifest for AI/ML is a YAML file designed to encapsulate all the necessary information about the package, including code, datasets, models, and their metadata. This reference documentation outlines the structure and specifications of the manifest format.

## Overview

The manifest is structured into several key sections: `version`, `package`,`code`, `datasets` and `models`. Each section serves a specific purpose in describing the AI/ML package components and requirements.

### `version`

- **Description**: Specifies the manifest format version.
- **Type**: String
- **Example**: `1.0`

### `package`

This section provides general information about the AI/ML project.

#### `name`

- **Description**: The name of the AI/ML project.
- **Type**: String

#### `version`

- **Description**: The current version of the project.
- **Type**: String
- **Example**: `1.2.3`

#### `description`

- **Description**: A brief overview of the project's purpose and capabilities.
- **Type**: String

#### `authors`

- **Description**: A list of individuals or entities that have contributed to the project.
- **Type**: Array of Strings


#### `code`

- **Description**: Information about the source code.
- **Type**: Object Array
  - `path`: Location of the source cod files or directory relative to the context
  - `description`: Description of what the code does.
  - `license`: SPDX license identifier for the code.

#### `datasets`

- **Description**: Information about the datasets used.
- **Type**: Object Array
  - `name`: Name of the dataset.
  - `path`: Location of the dataset file or directory relative to the context.
  - `description`: Overview of the dataset.
  - `license`: SPDX license identifier for the dataset.
  - `preprocessing`: Reference to preprocessing steps.

#### `models`

- **Description**: Details of the trained models included in the package.
- **Type**: Object Array
  - `name`: Name of the model 
  - `path`: Location of the model file or directory relative to the context
  - `framework`: AI/ML framework
  - `version`: Version of the model
  - `description`: Overview of the model
  - `license`: SPDX license identifier for the dataset. 
  - `training`:
    - `dataset`: Name of the dataset
    - `parameters`: name value pairs
  - `validation`:
    - `dataset`: Name of the dataset
    - `metrics`: name value pairs


## Example

```yaml
version: 1.0
package:
  name: AIProjectName
  version: 1.2.3
  description: >-
    A brief description of the AI/ML project.
  authors: [Author Name, Contributor Name]
code:
  - path: src/
    description: Source code for the AI models.
    license: Apache-2.0
datasets:
  - name: DatasetName
    path: data/dataset.csv
    description: Description of the dataset.
    license: CC-BY-4.0
    preprocessing: Preprocessing steps.
models:
  - name: ModelName
    path: models/model.h5
    framework: TensorFlow
    version: 1.0
    description: Model description.
    license: Apache-2.0
    training:
      dataset: DatasetName
      parameters:
        learning_rate: 0.001
        epochs: 100
        batch_size: 32
    Validation:
      - dataset: DatasetName
        metrics:
          accuracy: 0.95
          f1_score: 0.94
```


## Future Considerations 

This section is for collecting future ideas.

### `dependencies`

**This is a possible future section that may be used for creating BOM.**

- **Description**: Lists the project's external dependencies.
- **Type**: Object Array
  - `name`: Name of the dependency.
  - `version`: Version of the dependency.
  - `license`: SPDX license identifier for the dependency.

##### Example for dependencies
```yaml
  dependencies:
  - name: numpy
    version: 1.19.2
    license: BSD-3-Clause
  - name: pandas
    version: 1.1.3
    license: BSD-3-Clause
```