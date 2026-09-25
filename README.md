---
AIGC:
  ContentProducer: '001191110102MAD55U9H0F10002'
  ContentPropagator: '001191110102MAD55U9H0F10002'
  Label: '1'
  ProduceID: '3a1a73da-6363-445d-8e74-407b7ff18c5d'
  PropagateID: '3a1a73da-6363-445d-8e74-407b7ff18c5d'
  ReservedCode1: '5bda122e-4b69-4e0a-bac2-2363a23be8f4'
  ReservedCode2: '5bda122e-4b69-4e0a-bac2-2363a23be8f4'
---

# AutoStock

> Agentic Trading Operating System — a local financial capability runtime and trading workbench that any AI agent can drive.

```
用户
 ↓
Trading Skill（安装器 / 任务入口）
 ↓
Codex / OMP / WorkBuddy / Claude Code / 任意 Reasoning Plane
 ↓
Capability Gateway
 ├─ CLI（go-stock-cli --json）
 ├─ MCP（stdio server）
 └─ Local API
 ↓
AutoStock Runtime
 ├─ Market Data    ├─ News          ├─ Fundamentals
 ├─ Research       ├─ Backtest      ├─ Portfolio
 ├─ Risk Engine    └─ Broker Gateway
       ↓
    Paper / MiniQMT / vn.py / IBKR
```

AutoStock 基于 [go-stock](https://github.com/ArvinLovegood/go-stock) 的源码重构，保留其行情数据、K 线、财务、研报、资讯等数据能力，将核心架构升级为 **Capability Runtime + Trading Workbench** 模式：外部 AI Agent 负责推理与研究，AutoStock 负责确定性计算、风险校验、人机审批与券商执行。

## 设计原则

> **LLM proposes. Risk engine validates. Human authorizes. Broker executes.**

- **Numbers are code-calculated.** 金融数值与确定性计算由代码完成，LLM 负责研究、推理与解释。
- **Grade Outcome, not Tool Path.** 评测只看最终结果，允许多条合法工具路径。
- **Reasoning / Execution / Approval 三权分立。** Agent 可以提议，但不能自行下单。
- **Connect once, consume everywhere.** 同一套能力通过 CLI、MCP、Local API 统一暴露，不分别适配。

## Operation Registry

所有面向 Agent 的能力以 Operation 形式注册，每个 Operation 声明自己的 schema、权限、副作用与审批策略：

| 分类 | Operation 示例 |
|---|---|
| 行情 | `market.quote` `market.kline` `market.news` `market.industry_flow` |
| 研究 | `research.company` `research.industry` `strategy.backtest` |
| 持仓 | `portfolio.positions` `portfolio.cash` `portfolio.pnl` |
| 交易 | `order.preview` `order.place` `order.cancel` |
| 工作区 | `watchlist.add` `alert.create` `journal.append` |

每个 Operation 携带：

```
schema              # 参数与返回值类型
permissions         # 谁可以调用
side_effects        # 是否修改状态
risk_level          # 风险等级
is_mutation         # 是否写操作
supports_dry_run    # 是否支持预演
requires_approval   # 是否需要人工确认
idempotency_key     # 幂等键
audit_policy        # 审计策略
```

Agent 能做什么由 Registry 决定，与具体模型和 Agent 实现无关。

## 交易权限等级

| Level | Agent 权限 |
|---|---|
| **L0 Research** | 行情、新闻、财务、研报，只读 |
| **L1 Workspace** | 自选、预警、交易日志、策略配置可修改 |
| **L2 Paper** | Agent 可以自动模拟交易 |
| **L3 Live Confirm** | Agent 可以生成真实交易 Proposal，每笔必须人工确认 |
| **L4 Bounded Auto** | 在金额、标的、时间、策略等预设边界内自动执行 |

首个正式版目标为 **L3**：Agent 提议 → Risk Engine 校验 → 人工确认 → 券商执行。确认时对整个 OrderPlan 做 hash，确认 token 只对应这一笔参数，Agent 不能在确认后偷改价格或数量。

## Broker Gateway

交易层抽象为统一接口，不绑定单一券商：

```
BrokerGateway
  get_account()          get_cash()
  get_positions()        get_orders()
  get_trades()
  preview_order()        place_order()
  cancel_order()
  subscribe_order_updates()
  subscribe_trade_updates()
```

已规划的 Adapter：

| Adapter | 市场 | 状态 |
|---|---|---|
| PaperBroker | 全部 | 计划中（Phase 1） |
| MiniQMT / XtQuant | A 股 | 计划中 |
| vn.py Gateway | A 股 / 海外 | 规划中 |
| IBKR | 美股 / 港股 | 规划中 |

## Skill 设计

Skill 是一个轻量的 bootstrapper，不承载分析逻辑：

```
trading-skill/
├── SKILL.md          # 安装与使用说明
├── install           # 检测/下载/启动 runtime
├── doctor            # 诊断运行环境
├── config            # 配置 MCP 与 Broker
└── references/       # 参考文档
```

首次使用时，Agent 通过 Skill 自动完成：检测 runtime → 下载安装 → 启动 → 初始化数据库 → 注册 MCP → 检测 Broker → 返回能力列表。此后用户无需手动打开工作台，直接对 Agent 说话即可。

## 路线图

### Phase 1 — Capability Gateway（进行中）

- [ ] 同步上游 go-stock 最新 `dev` 基线
- [ ] 抽取 20-30 个只读 Operation（行情/K 线/财务/研报/资讯）
- [ ] 提供 `go-stock-cli --json` 和 stdio MCP 双通道
- [ ] Skill 实现为安装器/路由器
- [ ] Wails 界面升级为 Agent Workbench（显示 Agent 计划、工具调用、来源数据、Proposal、审计日志）

### Phase 2 — Paper Trading 闭环

- [ ] PaperBroker 实现
- [ ] Risk Engine：可用资金、涨跌停、单笔限额、重复订单、交易时段、行情过期
- [ ] Agent 经历完整闭环：research → proposal → risk check → order → fill → position → PnL → review

### Phase 3 — Live Trading（L3）

- [ ] 接入 MiniQMT / XtQuant
- [ ] 先开放 `account` / `positions` / `orders` 只读
- [ ] 然后 `order.preview`
- [ ] 最后 `order.place`，每笔必须人工确认

### Phase 4 — Bounded Autonomy（L4）

- [ ] 用户预设边界：标的白名单、单笔上限、日累计上限、最大回撤自动关闭
- [ ] 在边界内自动执行，边界外仍需确认

## 技术栈

- **Runtime**：Go + Wails（桌面工作台）
- **Frontend**：Vue + Naive UI
- **数据源**：东方财富、同花顺、通达信、Tushare、雪球
- **MCP**：stdio MCP Server（兼容 Codex / OMP / Claude Code / 任意 MCP Client）
- **Broker**：MiniQMT、vn.py、IBKR（规划中）

## 致谢

本项目基于 [ArvinLovegood/go-stock](https://github.com/ArvinLovegood/go-stock) 源码重构，保留了其行情数据、K 线图表、财务分析、资讯聚合等核心数据能力。上游项目以 GNU GPLv3 开源，本仓库遵循同一许可证。

## 参考项目

| 项目 | 借鉴方向 |
|---|---|
| [go-stock-mcp](https://github.com/WJS-WEB/go-stock-mcp) | go-stock → MCP 的最小抽取方式 |
| [OpenBB](https://github.com/OpenBB-finance/OpenBB) | Data → API → Workspace → Agent 架构 |
| [TradingAgents](https://github.com/TauricResearch/TradingAgents) | 多角色金融研究、point-in-time 数据正确性 |
| [FinRobot](https://github.com/AI4Finance-Foundation/FinRobot) | deterministic compute / provenance / Agent 编排 |
| [vn.py](https://github.com/vnpy/vnpy) | Broker Gateway 架构 |
| [Freqtrade](https://github.com/freqtrade/freqtrade) | dry-run/live 隔离、REST/CLI 控制面 |
| [LEAN](https://github.com/QuantConnect/Lean) | research/backtest/live 一体化引擎 |

## License

[GNU GPLv3](LICENSE)

---

> 本项目仅供学习研究。投资有风险，AI 分析结果不构成投资建议。