import axios from 'axios'
import { message } from 'antd'

// 创建axios实例
const request = axios.create({
    baseURL: '/api/v1',
    timeout: 10000
})

// 请求拦截器
request.interceptors.request.use(
    config => {
        // 从localStorage获取token
        const token = localStorage.getItem('token')
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    error => {
        console.error('请求错误:', error)
        return Promise.reject(error)
    }
)

// 响应拦截器
request.interceptors.response.use(
    response => {
        const { code, data, message: msg } = response.data

        if (code === 0) {
            return data
        } else {
            message.error(msg || '请求失败')
            return Promise.reject(new Error(msg || '请求失败'))
        }
    },
    error => {
        if (error.response) {
            const { status } = error.response

            if (status === 401) {
                message.error('未登录或登录已过期')
                localStorage.removeItem('token')
                window.location.href = '/login'
            } else if (status === 403) {
                message.error('没有权限访问')
            } else if (status === 404) {
                message.error('请求的资源不存在')
            } else if (status === 500) {
                message.error('服务器错误')
            } else {
                message.error(error.message || '请求失败')
            }
        } else {
            message.error('网络错误，请检查网络连接')
        }

        return Promise.reject(error)
    }
)

export default request
