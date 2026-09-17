#!/usr/bin/env bash
# check-full-tree.sh: assert this checkout is the full repository, not a released (lite) tree.
#
# Releases ship a stripped tree on master-lite: no tests, and a go.mod with no Ginkgo in it (see
# scripts/strip-tests.sh and RELEASING.md).  master-lite therefore always looks "ahead" of master,
# and merging or rebasing it back into master deletes every _test.go file and rewrites go.mod -
# cleanly, with no conflicts.  This is the guard against that: the test workflow runs it on every
# push, and scripts/release.sh runs it before cutting a release, so a stripped master is caught
# with a legible message instead of a mysteriously empty test run.
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

fail() {
	echo >&2 "check-full-tree.sh: $1"
	echo >&2 "This checkout looks like a released (lite) tree.  Has master-lite been merged or"
	echo >&2 "rebased into master?  master-lite is written only by the release workflow and must"
	echo >&2 "never be merged back - see RELEASING.md."
	exit 1
}

[[ -n "$(find . -path ./.git -prune -o -name '*_test.go' -print -quit)" ]] || fail "there are no _test.go files"
grep -q 'onsi/ginkgo' go.mod || fail "go.mod does not require Ginkgo"
