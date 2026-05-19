import request from '../utils/request'

// 创建订单
export const createOrder = (data) => {
    return request.post('/orders', data)
}

// 获取订单列表
export const getOrderList = (params) => {
    return request.get('/orders', { params })
}

// 支付订单
export const payOrder = (orderNo, data) => {
    return request.post(`/orders/${orderNo}/pay`, data)
}
