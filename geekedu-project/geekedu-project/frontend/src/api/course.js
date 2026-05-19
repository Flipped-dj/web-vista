import request from '../utils/request'

// 获取课程列表
export const getCourseList = (params) => {
    return request.get('/courses', { params })
}

// 获取课程详情
export const getCourseDetail = (id) => {
    return request.get(`/courses/${id}`)
}

// 获取分类树
export const getCategoryTree = () => {
    return request.get('/categories')
}

// 发布课程
export const createCourse = (formData) => {
    return request.post('/courses', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
    })
}

// 上传课程视频
export const uploadCourseVideo = (courseId, formData) => {
    return request.post(`/courses/${courseId}/videos`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
    })
}

// 获取视频播放地址
export const getVideoPlayUrl = (videoId) => {
    return request.get(`/player/${videoId}`)
}

// 获取课程视频列表
export const getCourseVideos = (courseId) => {
    return request.get(`/courses/${courseId}/videos`)
}
