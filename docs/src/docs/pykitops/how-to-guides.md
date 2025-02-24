# How-to Guides

This part of the project documentation focuses on a **problem-oriented** approach. You'll tackle common tasks that you might have, with the help of the code provided in this project.

## How To Create A Kitfile Object?

Whether you're working with an existing ModelKit's Kitfile,
or starting from nothing, the `kitops` package can help you
get this done.

Install the `kitops` package from PYPI into your project's environment
with the following command

```sh
pip install kitops
```

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
my_kitfile.manifestVersion = "3.0"
my_kitfile.package = {
    "name": "Another-Package",
    "version": "3.0.0",
    "description": "Another description",
    "authors": ["Someone"]
}
print(my_kitfile.to_yaml())

# OUTPUT:
#   manifestVersion: '3.0'
#   package:
#       name: Another-Package
#       version: 3.0.0
#       description: Another description
#       authors:
#       - Someone
```
