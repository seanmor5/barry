# Barry

Barry is a simple CLI written in Go for performing business banking and accounting tasks. Why use a dashboard or a dedicated finance tool when you can use your terminal?

## Installation

Nothing official right now.

## Usage

Once installed, you can use Barry to check balances, track spend, and track revenue.

```
$ barry
Barry is a CLI for performing common accounting and banking tasks

Usage:
  barry [command]

Available Commands:
  balances    View account balances
  help        Help about any command
  revenue     Track revenue across counterparties and periods
  spend       Track spend across counterparties and periods

Flags:
  -h, --help   help for barry

Use "barry [command] --help" for more information about a command.
```

## Account Access

Right now this uses the Mercury API but it will eventually switch to using [Teller](https://teller.io) when I can think of a good workflow for individuals connecting accounts to their own individual applications via the CLI.

## Why Barry?

Barry is my dad's name, and he is an accountant.