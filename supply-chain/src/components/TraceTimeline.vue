<template>
  <div class="trace-timeline">
    <div v-for="(event, index) in events" :key="index" class="timeline-item">
      <div class="timeline-marker" :class="event.type">
        <i :class="getIcon(event.type)"></i>
      </div>
      <div class="timeline-content">
        <div class="timeline-header">
          <span class="timeline-time">{{ formatTime(event.time) }}</span>
          <Tag :value="event.action" :severity="getSeverity(event.type)" size="small" />
        </div>
        <div class="timeline-title">{{ event.title }}</div>
        <div class="timeline-desc">{{ event.description }}</div>
        <div v-if="event.refNo" class="timeline-ref">
          <span>关联单号: {{ event.refNo }}</span>
          <span v-if="event.operator">操作人: {{ event.operator }}</span>
        </div>
      </div>
    </div>
    <div v-if="events.length === 0" class="empty-timeline">
      暂无流转记录
    </div>
  </div>
</template>

<script setup lang="ts">
import type { TimelineEvent } from '@/types'

defineProps<{
  events: TimelineEvent[]
}>()

const formatTime = (time: string) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

const getIcon = (type: string) => {
  const map: Record<string, string> = {
    procurement: 'ri-shopping-cart-line',
    inventory: 'ri-exchange-line',
    sales: 'ri-file-list-3-line',
    logistics: 'ri-truck-line',
    return: 'ri-arrow-go-back-line'
  }
  return map[type] || 'ri-information-line'
}

const getSeverity = (type: string) => {
  const map: Record<string, string> = {
    procurement: 'info',
    inventory: 'warning',
    sales: 'success',
    logistics: 'secondary',
    return: 'danger'
  }
  return map[type] || 'info'
}
</script>

<style scoped>
.trace-timeline {
  position: relative;
  padding-left: 40px;
  background: var(--surface-card);
  padding: 24px 24px 24px 64px;
  border-radius: var(--radius-lg);
  margin-bottom: 24px;
}

.trace-timeline::before {
  content: '';
  position: absolute;
  left: 39px;
  top: 24px;
  bottom: 24px;
  width: 2px;
  background: var(--surface-border);
}

.timeline-item {
  position: relative;
  padding-bottom: 20px;
}

.timeline-item:last-child {
  padding-bottom: 0;
}

.timeline-marker {
  position: absolute;
  left: -40px;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
  background: var(--surface-card);
  border: 2px solid var(--surface-border);
}

.timeline-marker.procurement { background: var(--blue-500); border-color: var(--blue-500); }
.timeline-marker.inventory { background: var(--orange-500); border-color: var(--orange-500); }
.timeline-marker.sales { background: var(--green-500); border-color: var(--green-500); }
.timeline-marker.logistics { background: var(--purple-500); border-color: var(--purple-500); }
.timeline-marker.return { background: var(--red-500); border-color: var(--red-500); }

.timeline-content {
  background: var(--surface-ground);
  padding: 16px;
  border-radius: 8px;
}

.timeline-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.timeline-time {
  font-size: 12px;
  color: var(--text-muted);
}

.timeline-title {
  font-weight: 600;
  margin-bottom: 4px;
}

.timeline-desc {
  font-size: 14px;
  color: var(--text-secondary);
}

.timeline-ref {
  margin-top: 8px;
  font-size: 12px;
  color: var(--text-muted);
  display: flex;
  gap: 16px;
}

.empty-timeline {
  text-align: center;
  padding: 40px;
  color: var(--text-muted);
}
</style>