<!--
SPDX-FileCopyrightText: SAP SE or an SAP affiliate company
SPDX-License-Identifier: Apache-2.0
-->

# Release Guide

We use [GoReleaser][goreleaser] and GitHub workflows for automating the release
process. We follow [semantic versioning][semver].

This repository uses **PR-driven releases**: pushes to `master` trigger a
`release-pr` workflow that opens (or updates) a release PR with a changelog and
version bump. Merging that PR tags and publishes the release automatically.

## Cutting a release

1. Land your changes on `master` with entries under the `[Unreleased]`
   section of [`CHANGELOG.md`](./CHANGELOG.md).

2. Find the release PR opened by the `release-pr` workflow on branch
   `chore/release-next`. By default it prepares a **patch** bump.
   For a **minor** or **major** release, open the **Actions** tab, pick the
   **release-pr** workflow, click **Run workflow**, and select the desired
   bump from the `version` dropdown (`patch` / `minor` / `major`).

3. Review and merge the PR. The `goreleaser` workflow's `tag` job creates and
   pushes the `vX.Y.Z` tag, and the `release` job runs goreleaser to publish.

If you need to cut a release without going through the release PR (e.g. to
re-tag), you can fall back to the [manual flow](#manual-flow) below.

## Manual flow

1. Ensure local `master` branch is up to date with `origin/master`:

  ```sh
  git fetch --all --tags
  ```

2. Ensure all checks are passing:

  ```sh
  make check
  ```

3. Update the [`CHANGELOG`](./CHANGELOG.md).
  Make sure that the format is consistent especially the version heading.
  We follow [semantic versioning][semver] for our releases.

  You can check if the file format is correct by running [`release-info`][release-info] for the new version:

  ```sh
  go install github.com/sapcc/go-bits/tools/release-info@latest
  release-info CHANGELOG.md X.Y.Z
  ```

  where `X.Y.Z` is the version that you are planning to release.

4. Commit the updated changelog with message: `Release <version>`
5. Create and push a new Git tag:

  ```sh
  git tag vX.Y.Z
  git push
  git push --tags
  ```

  > [!IMPORTANT]
  > Tags are prefixed with `v` and the GitHub release workflow is triggered for tags that match the `v[0-9]+.[0-9]+.[0-9]+` [gh-pattern].

[release-info]: https://github.com/sapcc/go-bits/tree/master/tools/release-info
[semver]: https://semver.org/spec/v2.0.0.html
[gh-pattern]: https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#patterns-to-match-branches-and-tags
[goreleaser]: https://github.com/goreleaser/goreleaser
