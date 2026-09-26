# Legacy cleanup baseline

This PR removes obvious source-product residue while preserving AutoStock's useful market, screening, agent, recommendation and backtest capabilities.

## Removed now

- inactive legacy `FloatingAiAssistant` surface (the newer Agent assistant remains);
- backup-only `agent-chat_bk.vue`;
- VIP/sponsor gates that blocked core owner workflows:
  - K-line analysis entry;
  - indicator-screen K-line inspection;
  - recommendation tracking signals and recommendation detail;
- duplicated binary copies of manuals when Markdown source is already present;
- inherited promotional/payment/screenshots that are not runtime assets;
- obsolete root demo images;
- old go-stock package/product metadata on the Wails/frontend package surface.

## Explicitly preserved

- market overview / market emotion;
- natural-language indicator screening;
- market/data tool handlers;
- agent runtime;
- MCP and Skill infrastructure;
- recommendation persistence and automatic recommendation saving;
- recommendation backtest;
- K-line/charting;
- alerts and trading journal;
- data providers.

## Deferred cleanup

Some inherited product areas are coupled across frontend/backend and should be removed in dedicated follow-up PRs after dependency checks:

- sponsor backend/data models and remaining About-page sponsor presentation;
- prompt/community plaza services;
- duplicate/legacy route hierarchy;
- generated Wails bindings;
- the separate `ai-assistant-web` surface;
- updater logic tied to the upstream release model.

The rule for subsequent cleanup is: **remove legacy presentation/coupling, not valuable capability.**
