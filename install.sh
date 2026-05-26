#!/usr/bin/env sh
set -eu

prefix=${PREFIX:-"$HOME/.local"}
bindir=${BINDIR:-"$prefix/bin"}
tmpdir=${TMPDIR:-/tmp}/regi-install-$$

cleanup() {
	rm -rf "$tmpdir"
}

trap cleanup EXIT INT TERM

if ! command -v go >/dev/null 2>&1; then
	echo "go is required to install regi" >&2
	exit 1
fi

mkdir -p "$bindir" "$tmpdir"
go build -trimpath -ldflags="-s -w" -o "$tmpdir/regi" .
install -m 755 "$tmpdir/regi" "$bindir/regi"

echo "installed regi to $bindir/regi"
echo "homebrew install: brew tap hamb1y/tap && brew install regi"

case ":$PATH:" in
	*:"$bindir":*) ;;
	*) echo "add $bindir to PATH to run regi from anywhere" ;;
esac
