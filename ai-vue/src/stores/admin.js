// 导入 Pinia 的 defineStore 函数
import { defineStore } from 'pinia'
// 导入 Vue 的 ref 函数
import { ref } from 'vue'

// 定义并导出 admin store
export const useAdminStore = defineStore('admin',()=> {
  // 侧边栏折叠状态（默认不折叠）
  const isCollapse = ref(false)
  
  // 切换侧边栏折叠状态的方法
  const toggleCollapse=()=> {
    isCollapse.value = !isCollapse.value
  }

  // 导出状态和方法
  return {
    isCollapse,     // 侧边栏折叠状态
    toggleCollapse  // 切换折叠状态的方法
  }
})
