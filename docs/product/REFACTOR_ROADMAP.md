# Refactor roadmap

## Phase 0 — establish boundaries

- freeze new feature sprawl in the legacy navigation;
- inventory routes, components, services and generated artifacts;
- classify each item: keep / wrap / migrate / archive / delete;
- introduce product and UI contracts before visual rewrites.

## Phase 1 — clean baseline

Remove or archive source-product residue that is not part of AutoStock:

- old branding and package metadata;
- sponsor/VIP gates;
- community/plaza promotion surfaces;
- backup components and generated artifacts that are committed redundantly;
- duplicate binary manuals and historical marketing screenshots;
- redundant assistant surfaces.

Do not delete market-data, charting, research or backtest capability merely because the current UI is being retired.

## Phase 2 — shell first

Build the compact shell before redesigning feature pages:

- semantic design tokens;
- flat sidebar;
- command/search entry;
- page header/breadcrumb;
- consistent local tabs;
- route compatibility layer.

The shell PR should not simultaneously rewrite every feature page.

## Phase 3 — capability gateway

Introduce the Operation Registry and wrap the highest-value read-only operations first:

1. quote / market status;
2. K-line;
3. news;
4. fundamentals;
5. flows;
6. reports;
7. positions/journal read paths.

Expose the same operations through MCP, JSON CLI and local API.

## Phase 4 — external-agent loop

- thin Trading Skill;
- runtime discovery/start/doctor;
- capability discovery;
- visible agent run/tool-call records;
- proposal artifacts;
- audit log.

## Phase 5 — paper trading

- PaperBroker;
- risk validation;
- preview/place/cancel lifecycle;
- fills/positions/PnL;
- deterministic replay tests.

## Phase 6 — live confirm

- broker adapter;
- read-only account sync first;
- order preview;
- signed/hashed proposal approval;
- live place/cancel only after explicit confirmation.

## Pull-request discipline

Keep refactor PRs single-purpose:

- docs/architecture;
- legacy cleanup;
- UI foundation;
- shell migration;
- capability slice;
- transport adapter;
- broker/risk changes.

A PR that mixes architecture, visual redesign, provider rewrites and trading execution should be split.
