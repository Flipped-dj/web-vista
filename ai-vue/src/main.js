// 导入 Vue 的 createApp 函数，用于创建 Vue 应用实例
import { createApp } from 'vue'
// 导入全局样式文件
import './style.css'
// 导入根组件 App
import App from './App.vue'
// 导入 Element Plus UI 库
import ElementPlus from 'element-plus'
// 导入 Element Plus 的样式文件
import 'element-plus/dist/index.css'
// 导入路由配置
import router from './router'
// 导入 Element Plus 图标库
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
// 导入 Pinia 状态管理库
import { createPinia } from 'pinia'

// 创建 Vue 应用实例，使用 App 作为根组件
const app = createApp(App)

// 创建 Pinia 实例
const pinia = createPinia()

// 遍历 ElementPlusIconsVue 对象，注册所有图标组件
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 应用使用 Element Plus、路由和 Pinia，然后挂载到 #app 元素
app.use(ElementPlus).use(router).use(pinia).mount('#app')
