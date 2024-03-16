# lotusdocs

We use git subtrees to include the lotusdocs theme. We do this so that we can add changes without copying whole
files into our own layouts.

# Updating

To pull the latest release of lotusdocs, run the following command (from the root of the project):

```shell
git subtree pull --prefix themes/lotusdocs https://github.com/colinwilson/lotusdocs.git release --squash
```
