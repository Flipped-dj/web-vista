// 导入 axios 库
import axios from 'axios'
// 导入 Element Plus 的消息组件
import { ElMessage } from 'element-plus'

// 创建 axios 实例
const service = axios.create({
    baseURL: '/api', // 请求前缀，所有请求都会自动添加该前缀
    timeout: 5000 // 请求超时时间（5秒）
})

// 请求拦截器
service.interceptors.request.use(
    config => {
        // 在发送请求之前的处理，从本地存储获取 token
        const token = localStorage.getItem('token')
        // 如果存在 token，添加到请求头
        if (token) {
            config.headers['token'] = token
        }
        // 返回处理后的配置
        return config
    },
    error => {
        // 对请求错误的处理
        return Promise.reject(error)
    }
)

// 响应拦截器
service.interceptors.response.use(
    response => {
        // 对响应数据的处理
        const { data, config } = response

        // 处理业务状态码
        if (data.code === '200') {
            // 成功状态码，返回响应数据中的 data 部分
            return data.data
        } else {
            // 错误状态码
            if (data.code === '-1') {
                // token 过期或认证失败
                if (!config.url?.includes('/login')) {
                    // 如果不是登录接口，提示登录过期
                    ElMessage.error(data.msg || '登录过期，请重新登录')
                    // 清空本地存储的登录信息
                    localStorage.removeItem('token')
                    localStorage.removeItem('userInfo')
                    // 跳转到登录页面
                    window.location.href = '/auth/login'
                }
            } else {
                // 其他错误，显示错误消息
                ElMessage.error(data.msg || '登录过期，请重新登录')
                // 拒绝 promise，让调用方处理错误
                return Promise.reject('网络请求失败....')
            }
        }
        // 兜底返回响应对象
        return response
    },
    error => {
        // 对响应错误的处理
        console.log('响应错误:', error)
        // 拒绝 promise，让调用方处理错误
        return Promise.reject(error)
    }
)

// 导出 axios 实例
export default service