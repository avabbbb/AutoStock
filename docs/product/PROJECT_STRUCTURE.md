# Target project structure

This document defines the target boundaries for the AutoStock refactor. The current snapshot does not need to move all at once; new work should converge toward these boundaries.

```text
AutoStock/
├─ cmd/
│  ├─ autostock/             # desktop/runtime entry
│  └─ autostock-cli/         # machine-readable CLI
├─ internal/
│  ├─ capability/            # Operation Registry + schemas
│  ├─ market/                # quotes, K-line, flows, market calendars
│  ├─ research/              # news, fundamentals, reports, research artifacts
│  ├─ portfolio/             # positions, cash, PnL, journal
│  ├─ strategy/              # deterministic strategy/backtest services
│  ├─ risk/                  # validation and policy engine
│  ├─ proposal/              # human-reviewable mutations
│  ├─ broker/                # BrokerGateway + adapters
│  ├─ automation/            # schedules, alerts, agent jobs
│  ├─ audit/                 # immutable operation/audit records
│  └─ platform/              # config, storage, lifecycle
├─ transports/
│  ├─ mcp/                   # MCP adapter over capabilities
│  ├─ cli/                   # JSON CLI adapter
│  └─ localapi/              # local HTTP/RPC adapter
├─ skills/
│  └─ trading/               # thin bootstrap/doctor/config skill
├─ frontend/
│  └─ src/
│     ├─ app/                # shell, routes, navigation, command palette
│     ├─ features/
│     │  ├─ today/
│     │  ├─ research/
│     │  ├─ portfolio/
│     │  ├─ automations/
│     │  └─ settings/
│     ├─ entities/           # stock/order/position/proposal views
│     ├─ components/         # shared UI primitives
│     ├─ design/             # tokens/theme/density
│     └─ lib/
├─ docs/
│  ├─ product/
│  ├─ architecture/
│  └─ archive/               # historical upstream docs when still useful
└─ tests/
   ├─ contract/
   ├─ integration/
   └─ e2e/
```

## Core architecture

```text
Agent / Skill / User
        │
        ▼
Operation Registry
        │
        ├──── MCP
        ├──── CLI
        └──── Local API
        │
        ▼
Capability Runtime
 market · research · portfolio · strategy
        │
        ▼
Risk Engine
        │
        ▼
Proposal / Approval
        │
        ▼
Broker Gateway
```

## Operation contract

Every agent-facing operation must declare:

- stable name;
- input/output schema;
- read/write classification;
- side effects;
- risk level;
- required permission;
- dry-run support;
- approval requirement;
- idempotency behavior;
- audit behavior.

Examples:

```text
market.quote
market.kline
market.news
research.company
portfolio.positions
portfolio.pnl
watchlist.add
alert.create
order.preview
order.place
order.cancel
```

Transport-specific code must not contain business rules that belong to the operation implementation.

## Dependency rule

Dependencies flow inward:

```text
UI / MCP / CLI
      ↓
Capability contract
      ↓
Domain services
      ↓
Adapters / providers / storage
```

The frontend should not directly own provider logic; an MCP adapter should not reimplement market calculations; broker-specific SDK types should not leak into generic proposals.

## Migration rule

Prefer small strangler migrations:

1. wrap an existing upstream service behind a typed operation;
2. add contract tests;
3. point desktop/MCP/CLI consumers at the operation;
4. only then remove the old duplicate entry point.

Do not begin with a giant folder move.
