// 导入 Vue Router 核心函数
import { createRouter, createWebHistory } from 'vue-router'
// 导入后台布局组件
import backendLayout from '@/components/backendLayout.vue'
// 导入认证布局组件
import AuthLayout from '@/components/AuthLayout.vue'
// 导入前台布局组件
import FrontendLayout from '../components/FrontendLayout.vue'


// 后台路由配置
const backendRoutes = [
    {
        path: '/back', // 后台根路径
        redirect: '/back/dashboard', // 默认重定向到仪表盘
        component: backendLayout, // 使用后台布局组件
        children: [
            {
                path: 'dashboard', // 仪表盘路由
                component: () => import('@/views/dashboard.vue'), // 懒加载组件
                meta: {
                    title: '数据分析', // 页面标题
                    icon: 'PieChart' // 菜单图标
                }
            },
            {
                path: 'knowledge', // 知识文章管理
                component: () => import('@/views/knowledge.vue'),
                meta: {
                    title: '知识文章',
                    icon: 'ChatLineSquare'
                }
            },
            {
                path: 'consultations', // 咨询记录管理
                component: () => import('@/views/consultations.vue'),
                meta: {
                    title: '咨询记录',
                    icon: 'Message'
                }
            },
            {
                path: 'emotional', // 情绪日志管理
                component: () => import('@/views/emotional.vue'),
                meta: {
                    title: '情绪日志',
                    icon: 'User'
                }
            },
        ]
    },
    {
        path: '/auth', // 认证相关路由
        component: AuthLayout, // 使用认证布局组件
        children: [
            {
                path: 'login', // 登录页面
                component: () => import('@/views/login.vue'),
                meta: {
                    title: '登录',
                }
            },
            {
                path: 'register', // 注册页面
                component: () => import('@/views/register.vue'),
                meta: {
                    title: '注册',
                }
            }
        ]
    }
]

// 前台路由配置
const frontendRoutes = [
    {
        path: '/', // 前台根路径
        component: FrontendLayout, // 使用前台布局组件
        children: [
            {
                path: '', // 首页
                component: () => import('@/views/home.vue'),
            },
            {
                path: 'consultation', // 咨询页面
                component: () => import('@/views/consultation.vue'),
            },
            {
                path: 'emotion-diary', // 情绪日志页面
                component: () => import('@/views/emotionDiary.vue'),
            },
            {
                path: 'knowledge', // 知识文章列表
                component: () => import('@/views/frontendKnowledge.vue'),
            },
            {
                path: 'knowledge/article/:id', // 文章详情页（带参数）
                component: () => import('@/views/articleDetail.vue'),
                props: true // 启用 props 传递路由参数
            }
        ]
    }
]

// 创建路由实例
const router = createRouter({
    history: createWebHistory(), // 使用 history 模式（无 # 号）
    routes: [...backendRoutes, ...frontendRoutes] // 合并后台和前台路由
})

// 路由前置守卫（权限控制）
router.beforeEach((to, from, next) => {
    // 从本地存储获取 token
    const token = localStorage.getItem('token')
    
    // 判断用户是否登录
    if (token) {
        // 获取用户信息
        const userInfo = JSON.parse(localStorage.getItem('userInfo'))
        
        // 后台用户（userType == 2）
        if (userInfo.userType == 2) {
            // 后台用户只能访问后台路由
            if (to.path.startsWith('/back')) {
                next() // 允许访问
            } else {
                // 重定向到后台仪表盘
                next('/back/dashboard')
            }
        } else if (userInfo.userType == 1) {
            // 前台用户（userType == 1）
            // 前台用户不能访问后台和认证路由
            if (to.path.startsWith('/back') || to.path.startsWith('/auth')) {
                // 重定向到前台首页
                next('/')
            } else {
                next() // 允许访问前台路由
            }
        }
    } else {
        // 未登录用户
        // 访问后台页面需要登录
        if (to.path.startsWith('/back')) {
            // 重定向到登录页面
            next('/auth/login')
        } else {
            // 前台页面允许未登录访问
            next()
        }
    }
})

// 导出路由实例
export default router