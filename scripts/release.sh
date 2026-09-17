#!/usr/bin/env bash
# release.sh: cut a Gomega release.  The Release workflow (.github/workflows/release.yml) runs this
# after the test workflow passes; see RELEASING.md.
#
# Usage: scripts/release.sh patch|minor
#
# One run: release ## Unreleased as vX.Y.Z on master, build the stripped-down release tree on
# master-lite (no tests, no Ginkgo in go.mod - see scripts/strip-tests.sh), tag *that* commit
# vX.Y.Z, push all three refs together, and create the GitHub release.  Re-running after a failure
# resumes: once the tag is on origin, the run builds from the tag and skips whatever is already
# released.
#
# GOMEGA_RELEASE_DRY_RUN=1 does everything locally, then stops short of pushing and creating the
# GitHub release.  A development aid only.
set -euo pipefail

bump=${1:-}
[[ "$bump" == patch || "$bump" == minor ]] || { echo "usage: $0 patch|minor" >&2; exit 2; }
dry_run=${GOMEGA_RELEASE_DRY_RUN:-}

cd "$(dirname "${BASH_SOURCE[0]}")/.."

say() { printf '\n==> %s\n' "$*"; }
fail() { printf 'release failed: %s\n' "$*" >&2; exit 1; }
tool() { .release/release-tool "$@"; }
# Commits made by the workflow, including the master-lite commit built with commit-tree.
as_bot() {
	git -c user.name="github-actions[bot]" -c user.email="41898282+github-actions[bot]@users.noreply.github.com" "$@"
}

[[ "$(git rev-parse --abbrev-ref HEAD)" == master ]] || fail "releases are cut from master"
[[ -z "$(git status --porcelain)" ]] || fail "the working tree is not clean"
scripts/check-full-tree.sh # a release is cut from the full tree, never from a stripped one

mkdir -p .release && go build -o .release/release-tool ./scripts/release
version=$(tool next "$bump")
tag="v$version"

# ls-remote exits 2 when the tag is absent; anything else nonzero means origin could not be asked.
tag_status=0
git ls-remote --exit-code --tags origin "refs/tags/$tag" >/dev/null || tag_status=$?
[[ $tag_status == 0 || $tag_status == 2 ]] || fail "could not check origin for $tag"

if [[ $tag_status == 0 ]]; then
	say "$tag is already on origin - resuming the release from it"
	git fetch --force origin "refs/tags/$tag:refs/tags/$tag"
	# The tag is on master-lite, so this checks out the stripped tree; the changelog and
	# gomega_dsl.go the release tool reads are both still in it.
	git checkout --quiet --detach "$tag"
	[[ "$(tool current)" == "$version" ]] || fail "$tag does not have GOMEGA_VERSION $version"
else
	say "Preparing $tag on master"
	tool prepare "$version"
	go build ./... # gomega_dsl.go is Go source; never ship a release that does not compile
	as_bot commit --quiet --all --message "$tag (full)"
	release_commit=$(git rev-parse HEAD)
	git show --stat HEAD

	# master-lite is the released branch: the same tree with the tests stripped out and go.mod
	# re-tidied, so projects that use Gomega don't inherit Ginkgo.  Every release adds one commit
	# to it, merging the release commit from master.
	git fetch --force origin refs/heads/master-lite:refs/remotes/origin/master-lite ||
		fail "could not fetch origin/master-lite - it is the parent of every released commit"
	lite_parent=$(git rev-parse origin/master-lite)

	say "Building the master-lite tree for $tag"
	git checkout --quiet --detach "$release_commit"
	scripts/strip-tests.sh
	as_bot commit --quiet --all --message "$tag (stripped)"
	lite_tree=$(git show --no-patch --format=%T HEAD)
	lite_commit=$(as_bot commit-tree "$lite_tree" -p "$lite_parent" -p "$release_commit" -m "$tag")
	git checkout --quiet master # back to the full tree; the stripped tree lives on in lite_commit

	# The tag is what the go toolchain resolves, so check the tree it points at before pushing it.
	[[ -z "$(git ls-tree -r --name-only "$lite_commit" | grep '_test\.go$' || true)" ]] ||
		fail "the master-lite tree still has _test.go files"
	! git show "$lite_commit:go.mod" | grep -q 'onsi/ginkgo' ||
		fail "the master-lite go.mod still requires Ginkgo"
	git show --stat "$lite_commit"

	git tag "$tag" "$lite_commit"
	if [[ -n "$dry_run" ]]; then
		say "dry run: not pushing master, master-lite, and $tag"
	else
		# Atomic: origin gets both branches and the tag together or not at all, so a tag on origin
		# is proof the commits are there too.  Fast-forward only - fails if either branch moved.
		git push --atomic origin \
			"$release_commit:refs/heads/master" \
			"$lite_commit:refs/heads/master-lite" \
			"refs/tags/$tag"
	fi
fi

tool notes "$version" >.release/notes.md
say "Release notes for $tag"
cat .release/notes.md

if [[ -n "$dry_run" ]]; then
	say "dry run: not creating the GitHub release $tag"
elif gh release view "$tag" >/dev/null 2>&1; then
	say "The GitHub release $tag already exists - refreshing its notes"
	gh release edit "$tag" --notes-file .release/notes.md
else
	say "Creating the GitHub release $tag"
	gh release create "$tag" --verify-tag --title "$tag" --notes-file .release/notes.md
fi

say "Released $tag${dry_run:+ (dry run)}"
