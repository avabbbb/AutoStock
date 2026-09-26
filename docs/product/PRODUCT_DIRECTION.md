# AutoStock product direction

## Product statement

AutoStock is an **agent-first local trading workbench**. It is not another chat box embedded in a stock terminal.

External reasoning agents such as Codex, OMP, WorkBuddy and other MCP-capable agents are the reasoning plane. AutoStock owns deterministic market capabilities, state, risk controls, approval and execution.

> LLM proposes. Risk engine validates. Human authorizes. Broker executes.

## Core owner workflow

AutoStock is optimized first for a compact owner workflow:

```text
Market Pulse
  ↓
Screen
  ↓
Agent Research
  ↓
Recommendation
  ↓
Track
  ↓
Review
```

1. **Market Pulse** — understand index direction, breadth, market emotion, sector strength and major events at a glance.
2. **Screen** — run indicator/factor/natural-language screens and get a compact candidate set.
3. **Agent Research** — let the local/external agent use AutoStock's existing data tools plus configured external APIs to research the candidates.
4. **Recommendation** — save selected ideas into a first-class recommendation record with rating, thesis, entry zone, take-profit, stop-loss and provenance.
5. **Track** — continuously show current price, recent trend and whether entry / take-profit / stop-loss conditions have been reached.
6. **Review** — open any recommendation for chart, thesis, risk, model/prompt provenance and backtest/outcome review.

This path is the primary UI and architecture priority. Features that do not strengthen this loop should not compete for permanent navigation.

## Primary agent experience

The user should also be able to start from a coding agent:

1. install or discover the AutoStock Skill;
2. let the Skill discover/start the local AutoStock runtime;
3. ask natural-language research or portfolio questions;
4. let the agent call typed AutoStock operations through MCP/CLI/local API;
5. open the Workbench only when visual context, state inspection or human approval is useful.

The GUI is therefore a **Workbench**, not the product's only entry point.

## Product surfaces

### Agent surface

- Skill: bootstrap, doctor, configuration, capability discovery.
- MCP: structured tool protocol for agent runtimes.
- CLI: stable JSON interface for scripts and coding agents.
- Local API: shared transport for desktop and future integrations.

### Workbench surface

The desktop app should focus on five jobs:

- **Today** — market state, portfolio state, alerts and current agent activity.
- **Research** — quotes, charts, news, fundamentals, reports and research artifacts.
- **Portfolio** — positions, cash, PnL, trading journal and review.
- **Automations** — agents, schedules, alerts, skills, MCP servers and audit history.
- **Settings** — data providers, models, brokers, permissions and safety limits.

Anything outside these jobs should be treated as optional/library content rather than permanent primary navigation.

## Explicit non-goals

AutoStock should not:

- keep two competing assistant UIs;
- expose every upstream feature as a permanent sidebar entry;
- preserve sponsor/VIP/community-plaza product logic from the source snapshot;
- make the LLM the authority for financial arithmetic, risk checks or order validity;
- let an agent silently mutate broker state;
- make GUI state the only way to access capabilities.

## Trading autonomy ladder

- L0 Research — read-only market/research operations.
- L1 Workspace — local workspace mutations such as watchlists, alerts and journal.
- L2 Paper — automatic paper orders.
- L3 Live Confirm — live proposals require explicit human approval.
- L4 Bounded Auto — automation only inside user-defined limits.

The first live-trading target is **L3**, not L4.

## Product decisions

1. **Agent-first, not embedded-chat-first.**
2. **One capability contract, many transports.**
3. **Workbench is an observable state surface.**
4. **Mutation is explicit and auditable.**
5. **Upstream data capabilities are reusable assets, not the product architecture.**
6. **Dense information is allowed; visual chrome is not.**
7. **Low-frequency features belong in Library/Command Palette, not the main sidebar.**
