# Kitfile Structure

The Kitfile defines the contents of your ModelKit. It is written in YAML and stored with the ModelKit. You can extract the Kitfile from any ModelKit with the Kit CLI:

```sh
kit unpack [registry/repo:tag] --config -d .
```

There are four main parts to a Kitfile:
1. ModelKit metadata in the `package` section
1. Path to the Jupyter notebook folder in the `code` section
1. Path to the serialized model in the `model` section
1. Path to the datasets in the `datasets` section (you can have multiple datasets in the same page)

Here's an example Kitfile:

```yaml
manifestVersion: v1.0.0

package:
  authors:
  - Jozu
  description: Updated model to analyze flight trait and passenger satisfaction data
  license: Apache-2.0
  name: FlightSatML

code:
- description: Jupyter notebook with model training code in Python
  path: ./notebooks

model:
  description: Flight satisfaction and trait analysis model using Scikit-learn
  framework: Scikit-learn
  license: Apache-2.0
  name: joblib Model
  path: ./models/scikit_class_model_v2.joblib
  version: 1.0.0

datasets:
- description: Flight traits and traveller satisfaction training data (tabular)
  name: training data
  path: ./data/train.csv
- description: validation data (tabular)
  name: validation data
  path: ./data/test.csv
```

The only mandatory parts of the Kitfile are:
* `manifestVersion`
* At least one of `code`, `model`, `or datasets` sections

A ModelKit can only contain one model, but multiple datasets or code bases are allowed. Also note that you can only use relative paths (no absolute paths) in your Kitfile. Right now you can only build ModelKits from files on your local system...but don't worry we're already working towards allowing you to reference remote files. For example, building a ModelKit from a local notebook and model, but a dataset hosted on DvC, S3, or anywhere else.

So a minimal ModelKit for distributing a pair of datasets might look like this:
```yaml
manifestVersion: v1.0.0

datasets:
- name: training data
  path: ./data/train.csv
- description: validation data (tabular)
  name: validation data
  path: ./data/test.csv
```

