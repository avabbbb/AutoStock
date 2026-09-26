<script setup>
import {onBeforeMount, onBeforeUnmount, onMounted, ref, computed} from 'vue'
import {GetConfig, GetTelegraphList, ReFleshTelegraphList} from "../../wailsjs/go/main/App";
import {EventsOff, EventsOn} from "../../wailsjs/runtime";
import {format} from 'date-fns';
import {zhCN} from 'date-fns/locale';
import AnalyzeMartket from "./AnalyzeMartket.vue";
import ConceptEventList from "./ConceptEventList.vue";
import RzrqRank from "./RzrqRank.vue";
import NewsList from "./newsList.vue";

const darkTheme = ref(false)
const telegraphList = ref([])
const time = ref(new Date())

const timeText = computed(() => format(time.value, 'yyyy-MM-dd HH:mm:ss EEEE QQQQ', {locale: zhCN}))

const updateTime = () => {
  time.value = new Date()
}

let timer = null

onBeforeMount(() => {
  GetConfig().then(res => {
    darkTheme.value = res.darkTheme
  }).catch(err => {
    console.error('Home GetConfig error:', err)
  })
})

onMounted(() => {
  timer = setInterval(updateTime, 1000)
  GetTelegraphList("财联社电报").then(res => {
    telegraphList.value = res || []
  }).catch(err => {
    console.error('GetTelegraphList error:', err)
  })
  EventsOn("newTelegraph", (data) => {
    if (data != null) {
      for (let i = 0; i < data.length; i++) {
        telegraphList.value.pop()
      }
      telegraphList.value.unshift(...data)
    }
  })
})

onBeforeUnmount(() => {
  if (timer) {
    clearInterval(timer)
  }
  EventsOff("newTelegraph")
})

function refreshTelegraph() {
  ReFleshTelegraphList("财联社电报").then(res => {
    telegraphList.value = res || []
  })
}
</script>

<template>
  <section class="market-pulse" style="--wails-draggable:no-drag">
    <header class="page-header">
      <div>
        <h1>Market</h1>
        <p>大盘走势、市场情绪与关键事件</p>
      </div>
      <n-text depth="3" class="market-clock">{{ timeText }}</n-text>
    </header>

    <n-flex vertical :size="10">
      <!-- 大盘分析：全球股指跑马灯 + 市场情绪 + 涨跌停/分时/融资融券走势 -->
      <AnalyzeMartket :dark-theme="darkTheme" :chart-height="280"/>

      <!-- 每日炒作题材 + 融资融券 + 财联社电报 三列等高 -->
      <n-flex :size="12" :wrap="false" style="align-items: stretch;">
        <div class="thin-scroll" style="flex: 1; min-width: 0; max-height: 600px; overflow-y: auto;">
          <ConceptEventList/>
        </div>
        <div class="thin-scroll" style="flex: 1; min-width: 0; max-height: 600px; overflow-y: auto;">
          <NewsList :newsList="telegraphList" :headerTitle="'财联社电报'" @update:message="refreshTelegraph"/>
        </div>
        <div class="thin-scroll" style="flex: 1; min-width: 0; max-height: 600px; overflow-y: auto;">
          <RzrqRank :dark-theme="darkTheme"/>
        </div>
      </n-flex>
    </n-flex>
  </section>
</template>

<style scoped>
.market-pulse {
  width: 100%;
}

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  padding: 2px 2px 12px;
  border-bottom: 1px solid rgba(128, 128, 128, 0.16);
  margin-bottom: 12px;
}

.page-header h1 {
  margin: 0;
  font-size: 18px;
  line-height: 24px;
  font-weight: 650;
}

.page-header p {
  margin: 2px 0 0;
  color: rgba(128, 128, 128, 0.9);
  font-size: 12px;
}

.market-clock {
  font-size: 11px;
  white-space: nowrap;
}

:deep(.thin-scroll) {
  scrollbar-width: none;
}
:deep(.thin-scroll::-webkit-scrollbar) {
  display: none;
}
</style>
