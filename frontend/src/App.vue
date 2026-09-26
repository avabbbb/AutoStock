<script setup>
import {
  EventsEmit,
  EventsOff,
  EventsOn,
  Quit,Hide ,
  WindowFullscreen,
  WindowUnfullscreen,
  WindowSetTitle
} from '../wailsjs/runtime'
import {h, onBeforeMount, onBeforeUnmount, onMounted, ref, watch} from "vue";
import {RouterLink, useRoute, useRouter} from 'vue-router'
import {createDiscreteApi,darkTheme,lightTheme , NIcon, NText,NButton,NProgress,dateZhCN,zhCN} from 'naive-ui'
import {
  AlarmOutline,
  AnalyticsOutline,
  BarChartSharp, Bonfire, BonfireOutline, BookOutline, CalendarOutline, DiamondOutline, DocumentTextOutline, EaselSharp,
  ExpandOutline, Flag,
  Flame, FlameSharp, FlaskOutline, GlobeOutline, HomeOutline, InformationOutline,
  LogoGithub,
  ChatbubblesOutline,
  NewspaperOutline,
  NewspaperSharp, Notifications,
  PowerOutline, Pulse,
  ReorderTwoOutline,
  SettingsOutline, ServerOutline, ScaleOutline, Skull, SkullOutline, SkullSharp,
  SparklesOutline, FlashOutline, Star,
  StarOutline,
  StatsChartOutline,
  Wallet, WarningOutline, TimeOutline, SearchOutline, BookmarkOutline,
} from '@vicons/ionicons5'
import {AnalyzeSentiment, GetConfig, GetGroupList, GetVersionInfo, IsTradingTime, IsHKTradingTime, IsUSTradingTime} from "../wailsjs/go/main/App";
import FloatingAgentAssistant from "./components/FloatingAgentAssistant.vue";
import SignalMonitorPanel from "./components/SignalMonitorPanel.vue";
import {Dragon, Fire, FirefoxBrowser, Gripfire, Robot} from "@vicons/fa";
import {Prompt, ReportAnalytics, ReportMoney, ReportSearch, TrendingUp} from "@vicons/tabler";
import {LocalFireDepartmentRound} from "@vicons/material";
import {AppsList20Regular, BoxSearch20Regular,SlideHide24Filled, CommentNote20Filled} from "@vicons/fluent";
import {FireFilled, MoneyCollectOutlined, NotificationFilled, StockOutlined} from "@vicons/antd";




const router = useRouter()
const loading = ref(true)
const loadingMsg = ref("加载数据中...")
const enableNews = ref(false)
const contentStyle = ref("")
const enableFund = ref(false)
const enableAgent = ref(false)
const enableDarkTheme = ref(darkTheme)
const content = ref('未经授权,禁止商业目的!\n\n数据来源于网络,仅供参考;投资有风险,入市需谨慎')
const isFullscreen = ref(false)
const activeKey = ref('home')
const route = useRoute()
// 路由变化时同步菜单高亮（如重定向、前进/后退）
const primaryRouteKey = {
  home: 'marketPulse',
  marketPulse: 'marketPulse',
  screener: 'screener',
  agent: 'agentResearch',
  agentResearch: 'agentResearch',
  picks: 'picks',
  settings: 'settings',
}

watch(() => route.name, (name) => {
  if (name && typeof name === 'string') {
    activeKey.value = primaryRouteKey[name] || name
  }
})
const containerRef = ref({})
const realtimeProfit = ref(0)
const telegraph = ref([])
const groupList = ref([])
const officialStatement= ref("")
const marketStatus = ref('')
let marketStatusTimer = null

const downloadState = ref({
  active: false, percentage: 0, speed: 0, avgSpeed: 0,
  downloaded: 0, total: 0, version: '', proxy: '',
  proxySpeed: 0, message: '', retrying: false,
})
let downloadNotification = null

function formatBytes(bytes) {
  if (!bytes || bytes <= 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  let val = bytes, i = 0
  while (val >= 1024 && i < units.length - 1) { val /= 1024; i++ }
  return val.toFixed(2) + ' ' + units[i]
}
function formatSpeed(bps) {
  if (!bps || bps <= 0) return '0 B/s'
  return formatBytes(bps) + '/s'
}
function renderDownloadContent() {
  const children = []
  if (downloadState.value.message) {
    children.push(h('div', {
      style: { 'margin-bottom': '8px', 'color': '#999', 'font-size': '12px', 'white-space': 'pre-wrap', 'max-height': '120px', 'overflow': 'hidden' }
    }, { default: () => downloadState.value.message }))
  }
  children.push(h(NProgress, {
    type: 'line',
    status: downloadState.value.retrying ? 'warning' : 'success',
    percentage: Math.round(downloadState.value.percentage),
    showIndicator: false,
    height: 8,
    borderRadius: 4,
  }))
  const detail = downloadState.value.retrying
    ? '正在尝试其他下载源...'
    : `${formatBytes(downloadState.value.downloaded)} / ${formatBytes(downloadState.value.total)} · ${formatSpeed(downloadState.value.speed)}`
  children.push(h('div', {
    style: { 'margin-top': '6px', 'font-size': '12px', 'color': '#888' }
  }, { default: () => detail }))
  if (downloadState.value.proxy) {
    children.push(h('div', {
      style: { 'margin-top': '2px', 'font-size': '11px', 'color': '#aaa' }
    }, { default: () => `下载源: ${downloadState.value.proxy}` }))
  }
  return h('div', { style: { 'text-align': 'left', 'font-size': '14px', 'min-width': '280px' } }, { default: () => children })
}

const investmentMottos = [
  "投资有风险，入市需谨慎",
  "别人贪婪我恐惧，别人恐惧我贪婪",
  "股市有风险，投资需谨慎",
  "不要把所有鸡蛋放在一个篮子里",
  "时间是优秀企业的朋友",
  "买股票就是买公司",
  "市场短期是投票机，长期是称重机",
  "保住本金是投资的第一要务",
  "在别人恐慌时贪婪，在别人贪婪时恐慌",
  "风险来自于你不知道自己在做什么",
  "价格是你付出的，价值是你得到的",
  "投资最重要的品质是耐心",
  "机会总是留给有准备的人",
  "知行合一，方能致远",
  "顺势而为，逆势而思",
  "投资是一场马拉松，不是百米冲刺",
  "独立思考是投资成功的关键",
  "市场永远在波动，但价值终将回归",
  "控制风险比追求收益更重要",
  "学习是最好的投资",
]
const currentMotto = ref(investmentMottos[Math.floor(Math.random() * investmentMottos.length)])

function refreshMotto() {
  currentMotto.value = investmentMottos[Math.floor(Math.random() * investmentMottos.length)]
}

function updateMarketStatus() {
  Promise.all([
    IsTradingTime().catch(() => false),
    IsHKTradingTime().catch(() => false),
    IsUSTradingTime().catch(() => false)
  ]).then(([cn, hk, us]) => {
    const parts = []
    parts.push(cn ? 'A股交易中' : 'A股休市')
    parts.push(hk ? '港股交易中' : '港股休市')
    parts.push(us ? '美股交易中' : '美股休市')
    marketStatus.value = parts.join(' | ')
    WindowSetTitle("AutoStock " + marketStatus.value + " " + officialStatement.value + "  「" + currentMotto.value + "」  [数据来源于网络，仅供参考；投资有风险，入市需谨慎]")
  })
}

function handleKlineAnalysisClick() {
  activeKey.value = 'klineAnalysis'
  router.push({ name: 'klineAnalysis' })
}

const menuOptions = ref([
  {
    label: () => h(RouterLink, { to: { name: 'marketPulse' } }, { default: () => 'Market' }),
    key: 'marketPulse',
    icon: renderIcon(HomeOutline),
  },
  {
    label: () => h(RouterLink, { to: { name: 'screener' } }, { default: () => 'Screener' }),
    key: 'screener',
    icon: renderIcon(SearchOutline),
  },
  {
    label: () => h(RouterLink, { to: { name: 'agentResearch' } }, { default: () => 'Research' }),
    key: 'agentResearch',
    icon: renderIcon(SparklesOutline),
  },
  {
    label: () => h(RouterLink, { to: { name: 'picks' } }, { default: () => 'Picks' }),
    key: 'picks',
    icon: renderIcon(BookmarkOutline),
  },
  {
    label: 'More',
    key: 'more',
    icon: renderIcon(AppsList20Regular),
    children: [
      {
        label: () => h(RouterLink, { to: { name: 'stock', query: { groupName: '全部', groupId: 0 } } }, { default: () => 'Watchlist' }),
        key: 'stock',
      },
      {
        label: () => h(RouterLink, { to: { name: 'market' } }, { default: () => 'Market data' }),
        key: 'market',
      },
      {
        label: () => h(RouterLink, { to: { name: 'dailyReview' } }, { default: () => 'Daily review' }),
        key: 'dailyReview',
      },
      {
        label: () => h(RouterLink, { to: { name: 'morningStrategy' } }, { default: () => 'Morning strategy' }),
        key: 'morningStrategy',
      },
      {
        label: () => h(RouterLink, { to: { name: 'recommendBacktestStats' } }, { default: () => 'Backtests' }),
        key: 'recommendBacktestStats',
      },
      {
        label: () => h(RouterLink, { to: { name: 'cronTasks' } }, { default: () => 'Automations' }),
        key: 'cronTasks',
      },
      {
        label: () => h(RouterLink, { to: { name: 'mcpServers' } }, { default: () => 'MCP servers' }),
        key: 'mcpServers',
      },
      {
        label: () => h(RouterLink, { to: { name: 'research', query: { name: 'AI分析报告' } } }, { default: () => 'Legacy tools' }),
        key: 'research',
      },
    ],
  },
  {
    label: () => h(RouterLink, { to: { name: 'settings' } }, { default: () => 'Settings' }),
    key: 'settings',
    icon: renderIcon(SettingsOutline),
  },
])

function renderIcon(icon) {
  return () => h(NIcon, null, {default: () => h(icon)})
}

function toggleFullscreen(e) {
  activeKey.value = 'full'
  //console.log(e)
  if (isFullscreen.value) {
    WindowUnfullscreen()
    //e.target.innerHTML = '全屏'
  } else {
    WindowFullscreen()
    // e.target.innerHTML = '取消全屏'
  }
  isFullscreen.value = !isFullscreen.value
}

// const drag = ref(false)
// const lastPos= ref({x:0,y:0})
// function toggleStartMoveWindow(e) {
//   drag.value=!drag.value
//   lastPos.value={x:e.clientX,y:e.clientY}
// }
// function dragstart(e) {
//   if (drag.value) {
//     let x=e.clientX-lastPos.value.x
//     let y=e.clientY-lastPos.value.y
//     WindowGetPosition().then((pos) => {
//       WindowSetPosition(pos.x+x,pos.y+y)
//     })
//   }
// }
// window.addEventListener('mousemove', dragstart)

EventsOn("realtime_profit", (data) => {
  realtimeProfit.value = data
})
EventsOn("telegraph", (data) => {
  telegraph.value = data
})

EventsOn("loadingMsg", (data) => {
  if(data==="done"){
    loadingMsg.value = "加载完成..."
    EventsEmit("loadingDone", "app")
    loading.value  = false
  }else{
    loading.value  = true
    loadingMsg.value = data
  }
})

setTimeout(() => {
  if (loading.value) {
    loading.value = false
    loadingMsg.value = "加载完成..."
    EventsEmit("loadingDone", "app")
  }
}, 8000)

onBeforeUnmount(() => {
  if (marketStatusTimer) {
    clearInterval(marketStatusTimer)
    marketStatusTimer = null
  }
  EventsOff("realtime_profit")
  EventsOff("loadingMsg")
  EventsOff("telegraph")
  EventsOff("newsPush")
  EventsOff("dailyReviewGenerated")
  EventsOff("morningStrategyGenerated")
  EventsOff("groupListChanged")
  EventsOff("updateDownloadStart")
  EventsOff("downloadProgress")
  EventsOff("updateDownloadComplete")
  EventsOff("updateDownloadFailed")
})

window.onerror = function (msg, source, lineno, colno, error) {
  // 将错误信息发送给后端
  EventsEmit("frontendError", {
    page: "App.vue",
    message: msg,
    source: source,
    lineno: lineno,
    colno: colno,
    error: error ? error.stack : null,
  });
  return true;
};

onBeforeMount(() => {
  GetVersionInfo().then(result => {
    if(result.officialStatement){
      content.value = result.officialStatement+"\n\n"+content.value
    }
    officialStatement.value = result.officialStatement || ""
    if (result.customBuild) {
      // 定制版本编译（wails build -tags custom）：隐藏"关于"菜单
      menuOptions.value.forEach((item) => {
        if (item.key === 'about') {
          item.show = false
        }
      })
    }
    updateMarketStatus()
  }).catch(err => {
    console.error("GetVersionInfo error:", err)
  })

  refreshStockGroupMenu()
  // 监听分组变化（新增/改名/删除），实时刷新菜单栏
  EventsOn("groupListChanged", () => {
    refreshStockGroupMenu()
  })


  GetConfig().then((res) => {
    enableFund.value = res.enableFund
    enableAgent.value = res.enableAgent

    menuOptions.value.filter((item) => {
      if (item.key === 'fund') {
        item.show = res.enableFund
      }
      if (item.key === 'agent') {
        item.show = res.enableAgent
      }
    })

    if (res.darkTheme) {
      enableDarkTheme.value = darkTheme
    } else {
      enableDarkTheme.value = null
    }
  }).catch(err => {
    console.error("GetConfig error:", err)
  })
})

onMounted(() => {
  updateMarketStatus()
  marketStatusTimer = setInterval(() => {
    refreshMotto()
    updateMarketStatus()
  }, 60000)
  contentStyle.value = "max-height: calc(92vh);overflow: hidden"
  GetConfig().then((res) => {
    if (res.enableNews) {
      enableNews.value = true
    }
    enableFund.value = res.enableFund
    enableAgent.value = res.enableAgent
    const {notification } =createDiscreteApi(["notification"], {
      configProviderProps: {
        theme: enableDarkTheme.value ? darkTheme : lightTheme ,
        max: 3,
      },
    })
    EventsOn("newsPush", (data) => {
      //console.log(data)
      if(data.isRed){
        notification.create({
          //type:"error",
         // avatar: () => h(NIcon,{component:Notifications,color:"red"}),
          title: data.time,
          content: () => h('div',{type:"error",style:{
              "text-align":"left",
              "font-size":"14px",
              "color":"#f67979"
            }}, { default: () => data.content }),
          meta: () => h(NText,{type:"warning"}, { default: () => data.source}),
          duration:1000*40,
        })
      }else{
         notification.create({
          //type:"info",
          //avatar: () => h(NIcon,{component:Notifications}),
          title: data.time,
          content: () => h('div',{type:"info",style:{
            "text-align":"left",
              "font-size":"14px",
              "color": data.source==="go-stock"?"#F98C24":"#549EC8"
            }}, { default: () => data.content }),
          meta: () => h(NText,{type:"warning"}, { default: () => data.source}),
          duration:1000*30 ,
        })
      }
    })

    EventsOn("dailyReviewGenerated", (data) => {
      if (!data) return
      if (data.status === 'failed') {
        notification.create({
          title: `每日复盘生成失败（${data.date}）`,
          content: () => h('div', {style: {"text-align": "left", "font-size": "14px", "color": "#f67979"}},
              {default: () => data.errorMessage || '未知错误'}),
          duration: 1000 * 30,
        })
      } else if (data.status === 'success') {
        notification.create({
          title: `每日复盘已生成（${data.date}）`,
          content: () => h('div', {style: {"text-align": "left", "font-size": "14px"}},
              {default: () => data.summary || ''}),
          duration: 1000 * 30,
        })
      }
    })

    EventsOn("morningStrategyGenerated", (data) => {
      if (!data) return
      if (data.status === 'failed') {
        notification.create({
          title: `盘前策略生成失败（${data.date}）`,
          content: () => h('div', {style: {"text-align": "left", "font-size": "14px", "color": "#f67979"}},
              {default: () => data.errorMessage || '未知错误'}),
          duration: 1000 * 30,
        })
      } else if (data.status === 'success') {
        notification.create({
          title: `盘前策略已生成（${data.date}）`,
          content: () => h('div', {style: {"text-align": "left", "font-size": "14px"}},
              {default: () => data.summary || ''}),
          duration: 1000 * 30,
        })
      }
    })

    EventsOn("updateDownloadStart", (data) => {
      downloadState.value = {
        active: true, percentage: 0, speed: 0, avgSpeed: 0,
        downloaded: 0, total: data.total || 0, version: data.version || '',
        proxy: data.proxy || '(直连)', proxySpeed: data.proxySpeed || 0,
        message: data.message || '', retrying: false,
      }
      if (downloadNotification) { downloadNotification.destroy(); downloadNotification = null }
      downloadNotification = notification.create({
        title: () => '正在下载新版本 ' + downloadState.value.version,
        content: renderDownloadContent,
        meta: () => h(NText, { type: 'warning' }, { default: () => 'AutoStock' }),
        duration: 0,
      })
    })

    EventsOn("downloadProgress", (data) => {
      if (data.status === 'retrying') {
        downloadState.value.retrying = true
        downloadState.value.percentage = 0
        downloadState.value.downloaded = 0
        downloadState.value.speed = 0
        return
      }
      downloadState.value.retrying = false
      downloadState.value.percentage = data.percentage || 0
      downloadState.value.speed = data.speed || 0
      downloadState.value.avgSpeed = data.avgSpeed || 0
      downloadState.value.downloaded = data.downloaded || 0
      downloadState.value.total = data.total || 0
      if (data.proxy) {
        downloadState.value.proxy = data.proxy
      }
    })

    EventsOn("updateDownloadComplete", (data) => {
      downloadState.value.active = false
      downloadState.value.percentage = 100
      if (downloadNotification) { downloadNotification.destroy(); downloadNotification = null }
      notification.create({
        title: '版本下载完成',
        content: () => h('div', { style: { 'text-align': 'left', 'font-size': '14px', 'color': '#52c41a' } },
          { default: () => '新版本 ' + data.version + ' 下载完成，正在应用更新，下次重启生效...' }),
        meta: () => h(NText, { type: 'warning' }, { default: () => 'go-stock' }),
        duration: 5000,
      })
    })

    EventsOn("updateDownloadFailed", (data) => {
      downloadState.value.active = false
      if (downloadNotification) { downloadNotification.destroy(); downloadNotification = null }
      const items = [
        h('div', { style: { 'margin-bottom': '8px' } },
          { default: () => '新版本 ' + data.version + ' 自动下载失败: ' + data.error })
      ]
      if (data.manualLinks) {
        items.push(h('div', { style: { 'margin-bottom': '4px' } }, { default: () => '请手动下载后替换程序文件:' }))
        items.push(h('div', {
          style: { 'font-size': '12px', 'word-break': 'break-all', 'margin-bottom': '2px', 'color': '#549EC8' }
        }, { default: () => '加速镜像: ' + data.manualLinks.mirror }))
        items.push(h('div', {
          style: { 'font-size': '12px', 'word-break': 'break-all', 'color': '#549EC8' }
        }, { default: () => '原始地址: ' + data.manualLinks.original }))
      }
      notification.create({
        title: '版本下载失败',
        content: () => h('div', { style: { 'text-align': 'left', 'font-size': '14px', 'color': '#f67979' } },
          { default: () => items }),
        meta: () => h(NText, { type: 'warning' }, { default: () => 'go-stock' }),
        duration: 0,
      })
    })

  }).catch(err => {
    console.error("GetConfig(onMounted) error:", err)
  })
})
</script>
<template>
  <n-config-provider ref="containerRef" :theme="enableDarkTheme" :locale="zhCN" :date-locale="dateZhCN">
    <n-message-provider>
      <n-notification-provider>
        <n-modal-provider>
          <n-dialog-provider>
            <n-watermark
                :content="''"
                cross
                selectable
                :font-size="16"
                :line-height="16"
                :width="500"
                :height="400"
                :x-offset="50"
                :y-offset="150"
                :rotate="-15"
            >
              <FloatingAgentAssistant />
              <SignalMonitorPanel />
              <n-flex>
                <n-grid x-gap="12" :cols="1">
                  <n-gi>
                    <n-spin :show="loading">
                      <template #description>
                        {{ loadingMsg }}
                      </template>
                      <n-marquee :speed="100" style="position: relative;top:0;z-index: 19;width: 100%"
                                 v-if="(telegraph.length>0)&&(enableNews)">
                        <n-tag type="warning" v-for="item in telegraph" style="margin-right: 10px">
                          {{ item }}
                        </n-tag>
                      </n-marquee>
                      <n-scrollbar :style="contentStyle">
                        <n-skeleton v-if="loading" height="calc(100vh)" />
                        <RouterView/>
                      </n-scrollbar>
                    </n-spin>
                  </n-gi>
                  <n-gi style="position: fixed;bottom:0;z-index: 9;width: 100%;">
                    <n-card size="small" style="--wails-draggable:no-drag">
                      <n-menu style="font-size: 18px;"
                              v-model:value="activeKey"
                              mode="horizontal"
                              :options="menuOptions"
                              :dropdown-props="{ menuProps: () => ({ style: 'max-height: 60vh; overflow-y: auto;' }) }"
                              responsive
                      />
                    </n-card>
                  </n-gi>
                </n-grid>
              </n-flex>
            </n-watermark>
          </n-dialog-provider>
        </n-modal-provider>
      </n-notification-provider>
    </n-message-provider>
  </n-config-provider>
</template>
<style>
/* 菜单/下拉弹出层滚动条样式（naive-ui 弹出层渲染到 body，需全局作用域） */
.n-dropdown-menu,
.n-dropdown-menu .n-vm-list,
.n-base-select-menu,
.n-base-select-menu .n-vm-list {
  scrollbar-width: thin;
  scrollbar-color: rgba(128, 128, 128, 0.45) transparent;
}
.n-dropdown-menu::-webkit-scrollbar,
.n-dropdown-menu .n-vm-list::-webkit-scrollbar,
.n-dropdown-menu .n-scrollbar-container::-webkit-scrollbar,
.n-base-select-menu::-webkit-scrollbar,
.n-base-select-menu .n-vm-list::-webkit-scrollbar,
.n-base-select-menu .n-scrollbar-container::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}
.n-dropdown-menu::-webkit-scrollbar-thumb,
.n-dropdown-menu .n-vm-list::-webkit-scrollbar-thumb,
.n-dropdown-menu .n-scrollbar-container::-webkit-scrollbar-thumb,
.n-base-select-menu::-webkit-scrollbar-thumb,
.n-base-select-menu .n-vm-list::-webkit-scrollbar-thumb,
.n-base-select-menu .n-scrollbar-container::-webkit-scrollbar-thumb {
  background-color: rgba(128, 128, 128, 0.45);
  border-radius: 3px;
}
.n-dropdown-menu::-webkit-scrollbar-thumb:hover,
.n-dropdown-menu .n-vm-list::-webkit-scrollbar-thumb:hover,
.n-dropdown-menu .n-scrollbar-container::-webkit-scrollbar-thumb:hover,
.n-base-select-menu::-webkit-scrollbar-thumb:hover,
.n-base-select-menu .n-vm-list::-webkit-scrollbar-thumb:hover,
.n-base-select-menu .n-scrollbar-container::-webkit-scrollbar-thumb:hover {
  background-color: rgba(128, 128, 128, 0.7);
}
.n-dropdown-menu::-webkit-scrollbar-track,
.n-dropdown-menu .n-vm-list::-webkit-scrollbar-track,
.n-dropdown-menu .n-scrollbar-container::-webkit-scrollbar-track,
.n-base-select-menu::-webkit-scrollbar-track,
.n-base-select-menu .n-vm-list::-webkit-scrollbar-track,
.n-base-select-menu .n-scrollbar-container::-webkit-scrollbar-track {
  background: transparent;
}
</style>
