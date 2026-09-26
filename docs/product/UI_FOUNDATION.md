# UI foundation — compact Notion-like workbench

## Goal

The UI should feel like a focused desktop workspace: flat, compact, quiet and information-dense.

"Notion-like" means interaction principles, not pixel cloning.

## Navigation model

Permanent sidebar:

```text
AutoStock
Search / Command

Today
Research
Portfolio
Automations

────────
Recent
Pinned
────────
Settings
```

Do **not** put every market dataset, backtest page, MCP page, skill page, plaza page and configuration page into the first navigation level.

Low-frequency tools should be discoverable through:

- contextual tabs;
- command palette;
- Library / More page;
- inline actions;
- settings subsections.

## Density rules

- Sidebar target width: 220–248 px, collapsible.
- Primary row height: 30–34 px.
- Toolbar controls: 28–32 px.
- Content max width only for prose/research reports; tables and charts can use the viewport.
- Prefer 8/12/16 px spacing rhythm.
- Prefer 6–8 px radius; avoid oversized 16–24 px cards.
- Use borders and background contrast before shadows.
- No decorative gradients in core workbench surfaces.
- No marketing banners inside operational screens.
- Icons are monochrome by default and should not carry navigation meaning alone.

## Page anatomy

```text
Breadcrumb / page title            contextual actions
optional 1-line status
────────────────────────────────────────────────────
local tabs / filters
────────────────────────────────────────────────────
primary content
```

A page should usually have one visual hierarchy, not nested cards inside cards.

## Information architecture

### Today

- market status;
- portfolio snapshot;
- alerts requiring attention;
- active/recent agent runs;
- pending proposals.

### Research

Tabs or contextual views for:

- watchlist;
- market;
- company;
- chart;
- news;
- fundamentals;
- reports;
- backtests.

### Portfolio

- positions;
- orders/fills;
- PnL;
- journal;
- daily review.

### Automations

- agent runs;
- schedules;
- alerts;
- skills;
- MCP;
- audit log.

### Settings

- data sources;
- model/providers;
- brokers;
- permissions;
- risk limits;
- appearance.

## Agent UX

There should be one canonical agent activity surface.

Avoid simultaneous:

- full-page agent chat;
- legacy floating AI assistant;
- floating agent assistant;
- multiple overlapping notification/progress panels.

Agent actions should produce first-class artifacts:

- ResearchRun
- ToolCall
- Finding
- Proposal
- Approval
- Execution
- AuditEvent

The user should be able to inspect them without reading raw chain-of-thought.

## Visual tokens

Start with semantic tokens rather than page-specific colors:

```text
--bg-canvas
--bg-sidebar
--bg-hover
--bg-selected
--border-subtle
--text-primary
--text-secondary
--text-muted
--accent
--positive
--negative
--warning
```

Market red/green are data semantics, not the application theme.

## Acceptance criteria for the UI refactor

- main sidebar has no more than five primary product destinations plus settings;
- no sponsor/VIP/community promotion in the primary shell;
- only one canonical agent entry/activity surface;
- every old route is either mapped into the new IA, moved to Library/More, or explicitly retired;
- light/dark mode both use the same semantic token system;
- keyboard navigation and command search are first-class;
- important tables/charts remain information-dense;
- no core task requires opening multiple floating panels.

## Reference rationale

Notion's 2026 redesign explicitly reduced sidebar crowding, split high-frequency areas into focused tabs, and moved less frequently accessed content into a Library. AutoStock adopts the same information-architecture principle for a financial workbench.
