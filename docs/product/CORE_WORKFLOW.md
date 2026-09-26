# Core workflow contract

## The product loop

```text
MARKET → SCREEN → RESEARCH → RECOMMEND → TRACK → REVIEW
```

This loop should be reachable with minimal navigation and should map to stable capability operations.

## 1. Market Pulse

Purpose: answer "what is the market doing now, and what is driving it?"

Preserve and consolidate existing capabilities such as:

- major A-share indexes and intraday trend;
- market breadth / up-down distribution;
- limit-up / limit-down statistics;
- market emotion;
- sector/concept strength and fund flow;
- important market/news events.

Target UI: one compact page with a small number of dense sections, not separate first-level menu items for every dataset.

Candidate operations:

```text
market.snapshot
market.indices
market.breadth
market.emotion
market.sector_flow
market.events
```

## 2. Screen

Purpose: quickly narrow the universe into candidates.

Preserve:

- natural-language indicator screening;
- hot/built-in screening strategies;
- saved custom screens;
- sortable result table;
- quick chart inspection.

Target UX:

- screen expression at the top;
- saved screens nearby;
- results immediately below;
- one-click "Research with Agent";
- one-click add to watchlist/research set.

Candidate operations:

```text
screen.run
screen.list_presets
screen.save
screen.delete
```

## 3. Agent Research

Purpose: use the agent as the reasoning plane while AutoStock remains the data/capability plane.

The agent should be able to combine:

- AutoStock market/data tools;
- company fundamentals;
- news/reports;
- flows and technical data;
- user-configured external APIs/MCP servers.

Research output must preserve provenance:

- model/provider;
- system/user prompt or task context;
- tool calls / source references;
- timestamp;
- candidate universe.

The UI needs one canonical agent activity surface, not multiple assistant windows.

## 4. Recommendation

A recommendation is a first-class artifact, not merely a paragraph in chat.

Minimum fields:

```text
symbol
name
rating
thesis
entry_zone
take_profit
stop_loss
risk_notes
strategy
recommended_at
reference_price
model
prompt/task provenance
source research run
```

AutoStock already contains an AI recommendation record and automatic-save fallback. Refactoring should strengthen this path instead of replacing it.

Required UX:

- "Save as recommendation" from screen/research result;
- automatic save when an agent emits the structured recommendation contract;
- easy manual correction before/after saving;
- de-duplication / recommendation history per symbol.

## 5. Track

The recommendation list is a core product page.

Each row/card should make the state obvious without opening details:

- current price and change;
- recommendation/reference price;
- entry zone;
- take-profit;
- stop-loss;
- recent mini trend;
- current state: waiting / entry reached / active / TP reached / SL reached / expired;
- age since recommendation;
- alert state.

Avoid VIP gating or community-product residue on this owner workflow.

## 6. Recommendation detail

Clicking a recommendation should show a dense detail workspace:

- price/K-line chart with entry, stop and take-profit overlays;
- recommendation thesis and risk notes;
- recent market/stock context;
- generated strategy;
- model and prompt/task provenance;
- alert history;
- backtest/outcome data;
- optional follow-up Agent action.

The user should not have to jump between several legacy pages to reconstruct the recommendation context.

## Existing implementation assets to preserve

Current snapshot already includes useful foundations:

- `AnalyzeMartket.vue` — broad market and emotion views;
- `SelectStock.vue` — natural-language indicator screening;
- `backend/data/tool_market_data.go` and other registered data tools;
- `backend/agent/auto_recommend_saver.go` — structured recommendation auto-save fallback;
- `backend/data/ai_recommend_stocks_api.go` — recommendation persistence and aggregation;
- `aiRecommendStocksList.vue` — current/recommended prices, entry/TP/SL tracking, mini trend and detail;
- recommendation backtest infrastructure.

These are migration inputs. The refactor should first wrap/recompose them, then replace legacy presentation and coupling incrementally.

## Product success criterion

A normal session should be possible with this sequence:

1. open AutoStock and understand market state within seconds;
2. run a screen;
3. send 3–10 candidates to the Agent;
4. receive structured research/recommendations;
5. persist selected recommendations in one action;
6. track all active recommendations from one page;
7. open a recommendation and understand exactly what changed and why.

If a feature does not help this loop, it should default to Library/More, settings, or removal.
