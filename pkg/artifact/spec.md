# Kit Model Specification v0.1

A *model package* is an ordered collection AI/ML artifacts and the corresponding execution and training parameters. This specification outlines the format of these artifacts and corresponding parameters and describes how to create and use them.

## Terminology

**Layer:** Models are composed of one or more layers. Each layer holds artifacts that are needed at model inference time. Layers do not carry metadata but these are properties of the model package.

**Model Manifest:**  Each model artifact has an associated JSON structure which describes some basic information about the model such as date created, author, format as well as inference time configuration.The JSON structure also references a cryptographic hash of each layer used by the model artifact, This JSON is considered to be immutable, because changing it would change the computed modelID. Changing it means creating a new derived model package, instead of changing the existing model image.

**ModelID:** Each modeld package's ID is given by the SHA256 hash of its manifest. It is represented as a hexadecimal encoding of 256 bits, e.g., sha256:a9561eb1b190625c9adb5a9513e72c4dedafc1cb2d4c5236c9a6957ec7dfd5a9. 

**Tag:** A tag serves to map a descriptive, user-given name to any single modelID. Tag values are limited to the set of characters [a-zA-Z0-9_.-], except they may not start with a . or - character. Tags are limited to 128 characters.

**Repository:** A collection of tags grouped under a common prefix (the name component before :). For example, in an model package tagged with the name myllm:3.1.4, myllm is the Repository component of the name. A repository name is made up of slash-separated name components, optionally prefixed by a DNS hostname. The hostname must comply with standard DNS rules, but may not contain _ characters. If a hostname is present, it may optionally be followed by a port number in the format :8080. Name components may contain lowercase characters, digits, and separators. A separator is defined as a period, one or two underscores, or one or more dashes. A name component may not start or end with a separator.


## Model Manifest Description


### Manifest fields

created `string`

ISO-8601 formatted combined date and time at which the image was created.

maintainer `string`

Gives the name and/or email address of the person or entity which created and is responsible for maintaining the image.

inputs `array of structs`

- **name:** Specifies the name of the input variable
- **type:** Indicates the data type of the input variable
- **dims:** Describes the dimensions of the input variable

outputs `array of structs`

- **name:** Specifies the name of the output variable
- **type:** Indicates the data type of the output variable
- **dims:** Describes the dimensions of the output variable