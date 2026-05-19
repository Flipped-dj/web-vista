# AI-Vue (AI 心理情感咨询/知识共享系统)

## 📌 项目概览 (Project Overview)
本项目是一个基于 Vue 3 + Vite 的现代化前端应用，主要功能涵盖前台内容展示与 AI 互动（情感日记、智能心理咨询）以及后台管理系统（控制台、文章管理、互动管理等）。以 AI 为核心驱动，为用户提供智能情感倾诉、知识检索与心理辅导等服务。

## 🚀 技术栈 (Tech Stack)
- **前端框架**：[Vue 3](https://vuejs.org/) (Composition API, `<script setup>`)
- **构建工具**：[Vite](https://vitejs.dev/) - 极速的现代前端构建工具
- **路由管理**：[Vue Router 4](https://router.vuejs.org/)
- **状态管理**：[Pinia](https://pinia.vuejs.org/) - 新一代轻量级状态管理
- **UI 组件库**：[Element Plus](https://element-plus.org/)
- **网络请求**：[Axios](https://axios-http.com/) & SSE (`@microsoft/fetch-event-source` 用于 AI 流式对话)
- **图表展示**：[ECharts](https://echarts.apache.org/) - 用于后台数据统计和可视化分析
- **富文本编辑器**：[WangEditor 5](https://www.wangeditor.com/)
- **CSS 预处理器**：[Sass](https://sass-lang.com/)

## 📂 核心目录结构 (Project Structure)
```text
ai-vue/
├── public/                 # 静态资源 (构建时将直接复制到 dist)
├── src/                    # 源码目录
│   ├── api/                # API 接口统一管理 (admin接口与前端接口分离)
│   ├── assets/             # 全局静态资源 (图片、字体、公共样式等)
│   ├── components/         # 基础与全局复用组件 (Navbar, Sidebar, MarkdownRenderer 等)
│   ├── config/             # 项目全局配置
│   ├── router/             # 路由配置 (前台路由、后台管理权限路由)
│   ├── stores/             # 状态管理块 (Pinia stores)
│   ├── utils/              # 工具函数 (如 request.js 封装 axios 和拦截器)
│   ├── views/              # 页面视图层 (登录/注册/控制台/文章详情等)
│   ├── App.vue             # 根组件
│   └── main.js             # 项目全局入口文件
├── index.html              # HTML 模板入口
├── package.json            # 依赖项与 NPM 脚本配置
└── vite.config.js          # Vite 构建与开发服务器配置
```

## 🛠️ 安装与运行 (Getting Started)

### 1. 环境要求
- **Node.js**: 建议 `>= 16.x`
- **包管理器**: **npm** (推荐) 或 **yarn** / **pnpm**

### 2. 克隆与安装依赖
首先将项目克隆到本地，然后进入项目目录执行依赖安装：
```bash
# 以 npm 为例
npm install

# 或者使用 yarn/pnpm:
# yarn install
# pnpm install
```

### 3. 本地开发服务器启动
安装完成后，执行以下命令即可启动带有热更新（HMR）特性的本地服务器：
```bash
npm run dev
```
之后在浏览器中打开控制台输出的地址（通常为 `http://localhost:5173`）即可预览项目。

### 4. 项目构建与打包 (Production Build)
当需要部署到线上环境时，执行打包命令：
```bash
npm run build
```
执行完毕后，所有优化压缩过的静态资源将会输出在项目的 `dist/` 文件夹下。

### 5. 预览生产环境构建 (Preview)
此命令可让你在本地快速拉起一个服务器，专门用于测试刚打包出来的 `dist` 目录运行效果：
```bash
npm run preview
```

## ✨ 主要功能模块 (Key Features)

### 👤 前台用户端 (Frontend)
- **AI 情感咨询/心理疏导**：依托 LLM 大模型与 SSE 流式输出技术，带来无缝打字的沉浸式对话辅导体验（支持文章/知识关联推荐）。
- **情感日记墙与图表**：用户可记录每日情绪（`emotionDiary.vue`），平台可提供状态可视化追踪。
- **知识共享网络**：心理与情感类文章的浏览、检索，包含 Markdown 的实时渲染。

### ⚙️ 后台管理端 (Backend Admin)
- **仪表盘 (Dashboard)**：对应用内用户量、咨询频次、情绪走向等结构化数据通过 `ECharts` 进行精美呈现与统计分析。
- **内容发布与知识管理**：使用集成了王座编辑器 (`WangEditor`) 的后台录入组件进行文章、动态等富文本的增删改查操作记录控制。
- **业务管理体系**：采用 `admin.js` 集中管理所有中台业务状态，具备鉴权布局 (`BackendLayout.vue`) 与独立侧边栏操作空间。

## 💡 开发与扩展向导
- **请求配置：** 如果遇到跨域问题或需要对接实际的后端大模型服务接口，请在 `vite.config.js` 的 `server.proxy` 字段中配置对应的后端代理。同时可在 `src/utils/request.js` 中配置你需要的请求头、动态 Token 以及全局拦截层。
- **样式定制：** 本项目基于 `Element Plus` 构建核心 UI 系统。可随时覆盖全局 SCSS 变量（如有）调整主题色调。

