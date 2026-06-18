# stonkler

`stonk` is a domain-first command-line tool for financial market data. It keeps
precise, provider-backed commands available as composable building blocks, while
offering a small porcelain layer (`get`) for common workflows.

Data is currently sourced from [Financial Modeling Prep (FMP)](https://financialmodelingprep.com/),
but the CLI is structured around provider interfaces so additional backends can
be added without changing the command surface.

## Installation

### With Nix

This repository ships a flake. To run directly:

```sh
nix run github:crertel/stonkler -- stocks quote AAPL
```

Or build the binary:

```sh
nix build
./result/bin/stonk stocks quote AAPL
```

### With Go

```sh
go install github.com/crertel/stonkler/cmd/stonk@latest
```

Or build from a checkout:

```sh
go build -o stonk ./cmd/stonk
```

## Configuration

Configuration starts with environment variables. To use the FMP backend, set:

```sh
export FMP_API_KEY=your_api_key_here
```

Inspect configuration and provider readiness:

```sh
stonk config show       # print non-secret configuration
stonk config doctor     # check that a provider is configured (exit 1 if not)
stonk config providers  # list data providers and their state
```

API keys are read from the environment only and are never written to disk,
logs, or output.

## Usage

```text
stonk <command> [flags]
```

Top-level commands are organized by domain first:

| Command       | Description                                          |
| ------------- | ---------------------------------------------------- |
| `stocks`      | Stock quotes, history, fundamentals, and watch views |
| `funds`       | ETF and mutual fund data                             |
| `crypto`      | Cryptocurrency market data                           |
| `forex`       | Foreign exchange market data                         |
| `commodities` | Commodity market data                                |
| `indexes`     | Index market data                                    |
| `search`      | Discover symbols and securities                      |
| `get`         | Workflow-oriented shortcuts (porcelain)              |
| `config`      | Configuration and provider diagnostics               |
| `version`     | Print version information                            |

Use `stonk <command> --help` for command-specific help.

### Domain commands

Domain commands are explicit and predictable. They avoid broad symbol
inference, hidden cross-domain fallback, and silently changing the asset class.

```sh
# Stocks
stonk stocks quote AAPL
stonk stocks quotes AAPL MSFT NVDA
stonk stocks history AAPL --from 2024-01-01 --to 2024-12-31
stonk stocks profile AAPL
stonk stocks statements AAPL income --period annual --limit 5
stonk stocks ratios AAPL --ttm
stonk stocks metrics AAPL --ttm
stonk stocks peers AAPL
stonk stocks analyst AAPL
stonk stocks transcript AAPL --year 2025 --quarter 1
stonk stocks insiders AAPL
stonk stocks sec AAPL

# Funds
stonk funds info SPY
stonk funds holdings SPY
stonk funds exposure AAPL
stonk funds sector-weightings SPY
stonk funds country-weightings VXUS

# Crypto / Forex / Commodities / Indexes
stonk crypto quote BTCUSD
stonk forex quote EURUSD
stonk commodities quote GCUSD
stonk indexes quote GSPC
```

### Search

`search` is the discovery surface for resolving names to symbols and securities.

```sh
stonk search apple
stonk search stocks apple
stonk search funds spy
stonk search cik 320193
stonk search isin US0378331005
stonk search screener --sector Technology --country US --market-cap-min 100B
```

### Get (porcelain)

`get` composes domain commands, resolves names to symbols, infers asset class,
and prefers helpful defaults over exhaustive flags.

```sh
stonk get quote AAPL
stonk get quote apple
stonk get company AAPL
stonk get history AAPL --from 2024-01-01
stonk get holdings SPY
stonk get transcript AAPL --latest
```

When interactive, `get` can ask you to disambiguate; when non-interactive, it
returns a clear ambiguity error instead.

### Watch

Market-like domains expose an interactive `watch` monitor — closer to `mtr`
than a one-shot quote — that refreshes on an interval and renders a stable
terminal table.

```sh
stonk stocks watch AAPL MSFT NVDA
stonk funds watch SPY VTI VXUS
stonk crypto watch BTCUSD ETHUSD
stonk forex watch EURUSD USDJPY
```

Watch supports `--interval`, `--sort`, and `--fields`, plus `--jsonl` for
streaming machine-readable updates. Row order from the command line is
preserved, and a failed or stale symbol is shown per-row without tearing down
the UI.

## Output formats

Human-readable table output is the default. Commands that fetch data also
support structured output:

```sh
stonk stocks quote AAPL --json
stonk stocks quotes AAPL MSFT --csv
stonk stocks watch AAPL MSFT --jsonl
```

Structured output avoids decoration and progress messages; all diagnostics go
to stderr. CSV includes a header row, uses stable column ordering, and avoids
nested data.

## Development

A Nix dev shell provides the full toolchain (Go, gopls, delve, golangci-lint,
just, and more):

```sh
nix develop
```

With [direnv](https://direnv.net/), `direnv allow` loads the shell
automatically via `.envrc`.

Common tasks:

```sh
go build ./...
go test ./...
```

The codebase is organized as:

- `cmd/stonk` — CLI entry point.
- `internal/cli` — command parsing, output rendering, and the watch loop.
- `internal/fmp` — Financial Modeling Prep provider client.

See [`docs/cli-design.md`](docs/cli-design.md) for the full CLI design notes.
