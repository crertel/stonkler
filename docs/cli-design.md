# CLI Design

`stonk` is a domain-first command-line tool for financial data. The CLI should
keep precise provider-backed commands available as composable "pipes", while
also offering a small porcelain layer for common workflows.

## Command Model

Top-level commands are organized by domain first:

```text
stonk stocks ...
stonk funds ...
stonk crypto ...
stonk forex ...
stonk commodities ...
stonk indexes ...
stonk search ...
stonk get ...
stonk config ...
```

Domain commands should be explicit and predictable. They should avoid broad
symbol inference, hidden cross-domain fallback, or silently changing the asset
class. Porcelain commands can infer, search, disambiguate, and choose useful
defaults.

## Domain Commands

### Stocks

```text
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
stonk stocks watch AAPL MSFT NVDA
```

### Funds

```text
stonk funds info SPY
stonk funds holdings SPY
stonk funds exposure AAPL
stonk funds sector-weightings SPY
stonk funds country-weightings VXUS
stonk funds watch SPY VTI VXUS
```

### Crypto, Forex, Commodities, And Indexes

```text
stonk crypto quote BTCUSD
stonk crypto history BTCUSD --from 2024-01-01
stonk crypto watch BTCUSD ETHUSD

stonk forex quote EURUSD
stonk forex history EURUSD --from 2024-01-01
stonk forex watch EURUSD USDJPY

stonk commodities quote GCUSD
stonk commodities history GCUSD --from 2024-01-01
stonk commodities watch GCUSD CLUSD

stonk indexes quote GSPC
stonk indexes history GSPC --from 2024-01-01
stonk indexes watch GSPC DJI IXIC
```

## Search

`search` is the discovery surface. It can expose generic lookup first, then
domain-specific lookup as the CLI grows.

```text
stonk search apple
stonk search stocks apple
stonk search funds spy
stonk search cik 320193
stonk search isin US0378331005
stonk search screener --sector Technology --country US --market-cap-min 100B
```

## Get

`get` is the porcelain layer. It should compose domain commands, perform
reasonable inference, and prefer helpful defaults over exhaustive flags.

Examples:

```text
stonk get quote AAPL
stonk get quote apple
stonk get company AAPL
stonk get history AAPL --from 2024-01-01
stonk get holdings SPY
stonk get transcript AAPL --latest
```

Unlike domain commands, `get` may:

- Resolve names to symbols.
- Infer asset class from a symbol or search result.
- Ask the user to disambiguate when interactive.
- Return a clear ambiguity error when non-interactive.
- Choose common defaults such as the latest transcript or annual statements.

## Watch

Each market-like domain should eventually expose a `watch` subcommand. `watch`
is an interactive terminal monitor, closer to `mtr` than a one-shot quote fetch.

Initial behavior:

```text
stonk stocks watch AAPL MSFT NVDA
stonk funds watch SPY VTI
stonk crypto watch BTCUSD ETHUSD
stonk forex watch EURUSD USDJPY
stonk stocks watch AAPL MSFT --stream
```

Expected features:

- Refresh on a configurable interval.
- Render a stable terminal table without scrolling.
- Highlight price, absolute change, percent change, volume, and update time.
- Preserve row order from the command line.
- Show stale data or request failures per symbol without tearing down the UI.
- Support `--jsonl` for streaming machine-readable updates.
- Support `--interval`, `--sort`, and `--fields` once the basic UI works.
- Support `stocks watch --stream` for FMP's real-time stock websocket feed.

`watch` should use the same provider interfaces as one-shot commands. The watch
loop belongs in CLI/application code, not inside a backend implementation.

## Output

Human-readable table output should be the default. JSON and CSV should be
available on commands that fetch data:

```text
stonk stocks quote AAPL --json
stonk stocks quotes AAPL MSFT --json
stonk stocks quotes AAPL MSFT --csv
stonk stocks watch AAPL MSFT --jsonl
```

Structured output should avoid decoration, progress messages, and shell prompts.
Any diagnostic output should go to stderr.

CSV should be intended for spreadsheet and batch export workflows. It should
include a header row by default, use stable column ordering, and avoid nested
data. Commands that naturally return nested records should either flatten common
fields or reject CSV with a clear error until a useful shape exists.

## Configuration

Configuration starts with environment variables:

```text
FMP_API_KEY
```

The `config` command should expose diagnostics before persistent secret storage:

```text
stonk config show
stonk config doctor
stonk config providers
```

Provider API keys must not be written to git, command logs, test fixtures, or
example output.
