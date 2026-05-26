# regi

Tiny newline-delimited plaintext register CLI.

Registers are stored in:

```text
~/.config/regi/registers/[register].regi
```

## Install

With Homebrew:

```sh
brew tap hamb1y/tap
brew install regi
```

Or from source:

```sh
sh install.sh
```

The source installer writes `regi` to `~/.local/bin` by default. Override it with `PREFIX` or `BINDIR`:

```sh
PREFIX=/usr/local sh install.sh
BINDIR="$HOME/bin" sh install.sh
```

## Usage

```sh
regi
regi work
regi add "buy milk"
regi add work call Sam
regi del work "^done:"
```
