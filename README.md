# Bump Version GitHub Action

This GitHub Action bumps the version based on commit messages. If a commit message contains either `#major`, `#minor` or `#patch` the action will increment the associated version number. If no commit message contains one of these keywords, the action will increment the patch version.

## Inputs

### `GITHUB_TOKEN`

**Required** The GitHub token to use for authentication. This should be a user generated Personal Access Token (PAT) if you want the new tags to trigger other GitHub action workflows.

## Example usage

```yaml
steps:
  - name: Checkout code
    uses: actions/checkout@v2

  - name: Bump version
    uses: gabeduke/bumpver-action@v1
    with:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```
