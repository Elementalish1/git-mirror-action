[![build](https://github.com/mathisve/git-mirror-action/actions/workflows/go.yaml/badge.svg?branch=master)](https://github.com/mathisve/git-mirror-action/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/mathisve/git-mirror-action)](https://goreportcard.com/report/github.com/mathisve/git-mirror-action)

# git-mirror-action
Github Action to mirror repositories.

## Arguments
### - `originalURL` 
**required -**
URL of the repository you want to mirror. (Has to be `https` and have `.git` at the end. See examples) 
### - `originalBranch`
**optional -**
Branch in the original repository you want to mirror.
### - `mirrorURL`
**required -**
URL of the repository you want to mirror in. (Has to be `https` and have `.git` at the end. See examples)
### - `mirrorBranch`
**optional -**
Branch in the mirror repository you want the mirror to go in.
### - `pat`
**required -**
**Base64 encoded** Personal Access Token! Preferable a secret. Has to be Base64 encoded. Will not work otherwise! Create a token [here!](https://github.com/settings/tokens)
### - `force`
**optional -**
Whether you want to use force when pushing into the mirror. (Might be required if branch already exists) Use with caution!
### - `verbose`
**optional -**
Whether you want to more verbosely log what git commands it executes. Useful for debugging.
### - `tags`
**optional -**
Whether you want to more transfer tags to the mirror repository. Useful for debugging.

## Example Usage
Mirrors `torvalds/linux:master` into `mathisve/mirror:mirror` once per day.
```
name: Example
on: 
  schedule:
    - cron: '0 0 1-31 * *'

jobs:
  mirror:
    runs-on: ubuntu-latest
    steps:
      - name: mirror
        uses: mathisve/git-mirror-action@latest
        with:
          originalURL: https://github.com/torvalds/linux.git
          mirrorURL: https://github.com/mathisve/mirror.git
          pat: ${{ secrets.PAT }}
```
Mirrors `moby/moby:20.10` into `mathisve/mirror:master` every 7 days with `--force` and `verbose`.
```
name: Example
on: 
  schedule:
    - cron: '0 0 */7 1-12 *'

jobs:
  mirror:
    runs-on: ubuntu-latest
    steps:
      - name: mirror
        uses: mathisve/git-mirror-action@latest
        with:
          originalURL: https://github.com/moby/moby.git
          originalBranch: 20.10
          mirrorURL: https://github.com/mathisve/mirror.git
          mirrorBranch: master
          verbose: true
          force: true
          tags: true
          pat: ${{ secrets.PAT }}
```