import request from '../utils/request'

// 用户注册
export const register = (data) => {
    return request.post('/register', data)
}

// 用户登录
export const login = (data) => {
    return request.post('/auth/login', data)
}

// 获取用户信息
export const getUserInfo = () => {
    return request.get('/user/info')
}

// 更新用户信息
export const updateProfile = (data) => {
    return request.put('/user/profile', data)
}
