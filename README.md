# Momentum

Some things are good for us but hard to start. Momentum keeps us moving in the
right direction. Build your first rhythm and start taking steps today.

Momentum is a phone-first web app for building a daily rhythm out of small
steps, alone or with friends.

## Documents

- [V0 spec](https://claude.ai/artifact/1rMtQBkgqDs4H3pLkvTuR3): workflows, vocabulary, momentum formula
- [Tech stack](https://claude.ai/artifact/35AokfuvYx7X6hwifiJURm): services, contracts, environments, compute
- [Implementation plan](https://claude.ai/artifact/TBwUCFPSJkNhQVU3Q2EiXJ): phases, subagents, test plan
- [Design canvas](https://claude.ai/artifact/PrGtTTRokJ9YKLxKtpxYQ2): every screen

## Layout

```text
proto/       API contracts (buf), from phase 1
gen/go/      generated Go code, from phase 1
pkg/         shared Go packages: clock, config, server
services/    one folder per service (web for now)
web/         React app, from phase 2
deploy/      Kustomize base and overlays (dev on OrbStack, prod on RKE2)
```

## Working on it

```sh
make check   # gofmt, go vet, staticcheck, unit tests
make itest   # integration tests (needs Docker)
make help    # every target
```

The whole repository is one Go module, `github.com/SpencerPergande/momentum`.
