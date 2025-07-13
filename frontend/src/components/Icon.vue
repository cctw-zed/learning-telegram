<template>
  <svg 
    :width="size" 
    :height="size" 
    viewBox="0 0 24 24" 
    fill="none" 
    xmlns="http://www.w3.org/2000/svg"
    :class="['icon', `icon-${name}`]"
  >
    <!-- 消息图标 -->
    <path v-if="name === 'message'" d="M20 2H4c-1.1 0-1.99.9-1.99 2L2 22l4-4h14c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zM6 9h12v2H6V9zm8 5H6v-2h8v2zm4-6H6V6h12v2z" :fill="color"/>
    
    <!-- 用户图标 -->
    <path v-if="name === 'user'" d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z" :fill="color"/>
    
    <!-- 群组图标 -->
    <path v-if="name === 'group'" d="M16 4c0-1.11.89-2 2-2s2 .89 2 2-.89 2-2 2-2-.89-2-2zm4 18v-6c0-1.11-.89-2-2-2s-2 .89-2 2v6h4zM12.5 11.5c.83 0 1.5-.67 1.5-1.5s-.67-1.5-1.5-1.5S11 9.17 11 10s.67 1.5 1.5 1.5zM17 22v-4c0-.83-.67-1.5-1.5-1.5S14 17.17 14 18v4h3zM6 4c0-1.11.89-2 2-2s2 .89 2 2-.89 2-2 2-2-.89-2-2zM4 22v-6c0-1.11.89-2 2-2s2 .89 2 2v6H4z" :fill="color"/>
    
    <!-- 发送图标 -->
    <path v-if="name === 'send'" d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z" :fill="color"/>
    
    <!-- 设置图标 -->
    <path v-if="name === 'settings'" d="M19.14,12.94c0.04-0.3,0.06-0.61,0.06-0.94c0-0.32-0.02-0.64-0.07-0.94l2.03-1.58c0.18-0.14,0.23-0.41,0.12-0.61 l-1.92-3.32c-0.12-0.22-0.37-0.29-0.59-0.22l-2.39,0.96c-0.5-0.38-1.03-0.7-1.62-0.94L14.4,2.81c-0.04-0.24-0.24-0.41-0.48-0.41 h-3.84c-0.24,0-0.43,0.17-0.47,0.41L9.25,5.35C8.66,5.59,8.12,5.92,7.63,6.29L5.24,5.33c-0.22-0.08-0.47,0-0.59,0.22L2.74,8.87 C2.62,9.08,2.66,9.34,2.86,9.48l2.03,1.58C4.84,11.36,4.8,11.69,4.8,12s0.02,0.64,0.07,0.94l-2.03,1.58 c-0.18,0.14-0.23,0.41-0.12,0.61l1.92,3.32c0.12,0.22,0.37,0.29,0.59,0.22l2.39-0.96c0.5,0.38,1.03,0.7,1.62,0.94l0.36,2.54 c0.05,0.24,0.24,0.41,0.48,0.41h3.84c0.24,0,0.44-0.17,0.47-0.41l0.36-2.54c0.59-0.24,1.13-0.56,1.62-0.94l2.39,0.96 c0.22,0.08,0.47,0,0.59-0.22l1.92-3.32c0.12-0.22,0.07-0.47-0.12-0.61L19.14,12.94z M12,15.6c-1.98,0-3.6-1.62-3.6-3.6 s1.62-3.6,3.6-3.6s3.6,1.62,3.6,3.6S13.98,15.6,12,15.6z" :fill="color"/>
    
    <!-- 退出图标 -->
    <path v-if="name === 'logout'" d="M17 7l-1.41 1.41L18.17 11H8v2h10.17l-2.58 2.58L17 17l5-5zM4 5h8V3H4c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h8v-2H4V5z" :fill="color"/>
    
    <!-- 搜索图标 -->
    <path v-if="name === 'search'" d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z" :fill="color"/>
    
    <!-- 在线状态图标 -->
    <circle v-if="name === 'online'" cx="12" cy="12" r="10" :fill="color"/>
    
    <!-- 离线状态图标 -->
    <circle v-if="name === 'offline'" cx="12" cy="12" r="10" fill="none" :stroke="color" stroke-width="2"/>
    
    <!-- 正在输入图标 -->
    <g v-if="name === 'typing'">
      <circle cx="7" cy="12" r="1" :fill="color">
        <animate attributeName="opacity" values="0;1;0" dur="1.4s" repeatCount="indefinite" begin="0s"/>
      </circle>
      <circle cx="12" cy="12" r="1" :fill="color">
        <animate attributeName="opacity" values="0;1;0" dur="1.4s" repeatCount="indefinite" begin="0.2s"/>
      </circle>
      <circle cx="17" cy="12" r="1" :fill="color">
        <animate attributeName="opacity" values="0;1;0" dur="1.4s" repeatCount="indefinite" begin="0.4s"/>
      </circle>
    </g>
    
    <!-- 检查图标 -->
    <path v-if="name === 'check'" d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z" :fill="color"/>
    
    <!-- 错误图标 -->
    <path v-if="name === 'error'" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" :fill="color"/>
    
    <!-- 加载图标 -->
    <g v-if="name === 'loading'">
      <circle cx="12" cy="12" r="10" fill="none" :stroke="color" stroke-width="2" opacity="0.3"/>
      <circle cx="12" cy="12" r="10" fill="none" :stroke="color" stroke-width="2" stroke-linecap="round" 
              stroke-dasharray="31.416" stroke-dashoffset="31.416">
        <animate attributeName="stroke-dasharray" dur="2s" values="0 31.416;15.708 15.708;0 31.416" repeatCount="indefinite"/>
        <animate attributeName="stroke-dashoffset" dur="2s" values="0;-15.708;-31.416" repeatCount="indefinite"/>
      </circle>
    </g>
  </svg>
</template>

<script setup lang="ts">
interface Props {
  name: string
  size?: number | string
  color?: string
}

withDefaults(defineProps<Props>(), {
  size: 24,
  color: 'currentColor'
})
</script>

<style scoped>
.icon {
  display: inline-block;
  vertical-align: middle;
  flex-shrink: 0;
}
</style> 