#!/usr/bin/env bash
# strip-tests.sh: turn the working tree into the tree Gomega releases.
#
# Gomega uses Ginkgo for its own tests, and the go toolchain pulls a dependency's *test*
# dependencies into consuming projects as indirect requirements - so a project that uses only
# Gomega would end up with all of Ginkgo in its go.mod.  Releases therefore ship a tree with every
# _test.go file deleted and go.mod re-tidied, which drops Ginkgo (see CHANGELOG.md 1.40.0).
# scripts/release.sh runs this to build the master-lite commit that carries the vX.Y.Z tag, and the
# test workflow runs it on every push so a non-test file importing Ginkgo is caught long before a
# release.
#
# This rewrites the working tree.  Run it on a throwaway checkout, not on your working copy.
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

echo "deleting $(find . -path ./.git -prune -o -name '*_test.go' -print | wc -l | tr -d ' ') _test.go files"
find . -path ./.git -prune -o -name '*_test.go' -delete
go mod tidy
# The whole point: the released module must build, and must not require Ginkgo.
go build ./...
if grep -q 'onsi/ginkgo' go.mod; then
	echo >&2 "strip-tests.sh: go.mod still requires Ginkgo after stripping the tests:"
	grep >&2 'onsi/ginkgo' go.mod
	echo >&2 "Something outside a _test.go file imports Ginkgo.  Gomega's released module cannot depend on it."
	exit 1
fi
