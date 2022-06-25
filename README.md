# tdc-bumpversion-action

This action will create a new tag version (semantic versioning) in a repository based on the latest created version.

For example, if the latest tag version on a repository is `v2.2.0` and you run this action with `incrementType` as `minor` the new tag will be `v2.3.0`.

## Pre-requirements

Make sure the repository that you are using this action with has atleast one tag with in the semantic versioning format, for example: `v0.0.1`.

## Using this action

```yaml
jobs:
  increment_version:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: the-developer-cafe/tdc-bumpversion-action@v2.3.0 # check for latest version
        with:
          gitEmail: youremail@email.com #required
          gitName: Your Name #required
          incrementType: minor #required (can be 'major', 'minor' or 'patch')
```
