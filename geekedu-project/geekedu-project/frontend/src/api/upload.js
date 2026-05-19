import request from '../utils/request'

// 上传文件
export const uploadFile = (formData) => {
    return request.post('/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
    })
}
