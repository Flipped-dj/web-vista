import React, { useState, useEffect } from 'react'
import { Layout, Card, Button, Descriptions, Tag, Collapse, message, Modal, List, Avatar, Breadcrumb } from 'antd'
import {
    PlayCircleOutlined,
    LockOutlined,
    ShoppingCartOutlined,
    HomeOutlined,
    BookOutlined,
    UserOutlined,
    ClockCircleOutlined,
    TeamOutlined,
    CheckCircleOutlined
} from '@ant-design/icons'
import { useParams, useNavigate } from 'react-router-dom'
import { getCourseDetail } from '../api/course'
import { createOrder, payOrder } from '../api/order'
import request from '../utils/request'
import './CourseDetail.css'

const { Content } = Layout

function CourseDetail() {
    const { id } = useParams()
    const navigate = useNavigate()
    const [loading, setLoading] = useState(false)
    const [course, setCourse] = useState(null)
    const [hasPurchased, setHasPurchased] = useState(false)
    const [playModalVisible, setPlayModalVisible] = useState(false)
    const [currentVideo, setCurrentVideo] = useState(null)
    const [videoUrl, setVideoUrl] = useState('')
    const token = localStorage.getItem('token')
    const userInfoStr = localStorage.getItem('userInfo')
    const userInfo = userInfoStr ? JSON.parse(userInfoStr) : null

    useEffect(() => {
        fetchCourseDetail()
        if (token) {
            checkPurchaseStatus()
        }
    }, [id, token])

    const fetchCourseDetail = async () => {
        setLoading(true)
        try {
            const data = await getCourseDetail(id)
            setCourse(data.courseDetail || data)
        } catch (error) {
            console.error('Failed to load course:', error)
        } finally {
            setLoading(false)
        }
    }

    const checkPurchaseStatus = async () => {
        try {
            const response = await request.get(`/orders/check/${id}`)
            setHasPurchased(response?.purchased === true)
        } catch (error) {
            console.error('Check purchase status failed:', error)
        }
    }

    const handleBuy = async () => {
        if (!token) {
            message.warning('请先登录')
            navigate('/login')
            return
        }

        try {
            const data = await createOrder({ courseId: parseInt(id) })
            message.success('订单创建成功，正在支付...')
            await payOrder(data.orderNo)
            message.success('购买成功！')
            setHasPurchased(true)
        } catch (error) {
            console.error('Purchase failed:', error)
        }
    }

    const handlePlayVideo = async (video) => {
        if (video.isFree !== 1 && !token) {
            message.warning('请先登录')
            navigate('/login')
            return
        }

        try {
            const response = await request.get(`/player/${video.id}`)
            if (response?.playUrl) {
                setCurrentVideo(video)
                setVideoUrl(response.playUrl)
                setPlayModalVisible(true)
                if (video.isFree === 1) {
                    message.info('免费试看视频')
                } else {
                    message.success(`获取播放地址成功，有效期${response.expire || 3600}秒`)
                }
            } else {
                message.error('获取播放地址失败')
            }
        } catch (error) {
            console.error('Play error:', error)
            if (error?.response?.status === 403) {
                message.error('您还未购买该课程，无法观看')
            }
        }
    }

    // 获取封面图片URL（OSS签名）
    const getCoverUrl = (coverUrl) => {
        if (!coverUrl) {
            return `https://picsum.photos/seed/${id}/800/450`
        }
        // 如果是OSS URL，我们无法在客户端签名，直接返回
        return coverUrl
    }

    if (loading || !course) {
        return (
            <div className="course-detail-loading">
                <div className="loading-spinner">加载中...</div>
            </div>
        )
    }

    const courseInfo = course.courseInfo || course
    const chapters = course.chapters || []
    const teacher = course.teacher || {}

    // 计算总视频数
    const totalVideos = chapters.reduce((sum, ch) => sum + (ch.videos?.length || 0), 0)

    return (
        <Layout className="course-detail-layout">
            {/* 面包屑导航 */}
            <div className="course-detail-breadcrumb">
                <Breadcrumb items={[
                    { title: <><HomeOutlined /> 首页</>, onClick: () => navigate('/') },
                    { title: <><BookOutlined /> 课程列表</>, onClick: () => navigate('/courses') },
                    { title: courseInfo.title }
                ]} />
            </div>

            <Content className="course-detail-content">
                {/* 课程头部 */}
                <div className="course-header-section">
                    <div className="course-cover-container">
                        <img
                            src={getCoverUrl(courseInfo.cover)}
                            alt={courseInfo.title}
                            className="course-main-cover"
                            onError={(e) => {
                                e.target.src = `https://picsum.photos/seed/${id}/800/450`
                            }}
                        />
                        {courseInfo.price === 0 && (
                            <span className="course-free-badge">免费课程</span>
                        )}
                    </div>
                    <div className="course-info-container">
                        <h1 className="course-main-title">{courseInfo.title}</h1>
                        <p className="course-main-desc">{courseInfo.description || '暂无课程简介'}</p>

                        <div className="course-meta-row">
                            <div className="meta-item">
                                <TeamOutlined />
                                <span>{courseInfo.studentCount || 0} 人学习</span>
                            </div>
                            <div className="meta-item">
                                <BookOutlined />
                                <span>{totalVideos} 课时</span>
                            </div>
                            <div className="meta-item">
                                <ClockCircleOutlined />
                                <span>随时学习</span>
                            </div>
                        </div>

                        <div className="course-teacher-row">
                            <Avatar icon={<UserOutlined />} src={teacher.avatar} />
                            <span className="teacher-name">讲师: {teacher.nickname || 'Admin'}</span>
                        </div>

                        <div className="course-price-row">
                            {courseInfo.price > 0 ? (
                                <>
                                    <span className="current-price">¥{courseInfo.price}</span>
                                    {courseInfo.originalPrice > courseInfo.price && (
                                        <span className="original-price">¥{courseInfo.originalPrice}</span>
                                    )}
                                </>
                            ) : (
                                <span className="free-price">免费</span>
                            )}
                        </div>

                        <div className="course-action-row">
                            {hasPurchased ? (
                                <Button type="primary" size="large" className="purchased-btn" disabled>
                                    <CheckCircleOutlined /> 已购买
                                </Button>
                            ) : (
                                <Button
                                    type="primary"
                                    size="large"
                                    className="buy-btn"
                                    icon={<ShoppingCartOutlined />}
                                    onClick={handleBuy}
                                >
                                    {courseInfo.price > 0 ? `立即购买 ¥${courseInfo.price}` : '免费学习'}
                                </Button>
                            )}
                        </div>
                    </div>
                </div>

                {/* 课程内容 */}
                <div className="course-content-section">
                    <Card title={<span><BookOutlined /> 课程目录</span>} className="course-chapters-card">
                        {chapters.length > 0 ? (
                            <Collapse
                                defaultActiveKey={chapters.map((_, i) => String(i))}
                                className="chapters-collapse"
                            >
                                {chapters.map((chapter, idx) => (
                                    <Collapse.Panel
                                        header={
                                            <div className="chapter-header">
                                                <span className="chapter-index">第{idx + 1}章</span>
                                                <span className="chapter-title">{chapter.title}</span>
                                                <span className="chapter-count">{chapter.videos?.length || 0}节</span>
                                            </div>
                                        }
                                        key={String(idx)}
                                    >
                                        <List
                                            dataSource={chapter.videos || []}
                                            renderItem={(video, vIdx) => (
                                                <List.Item className="video-item">
                                                    <div className="video-info">
                                                        <span className="video-index">{vIdx + 1}</span>
                                                        <span className="video-title">{video.title}</span>
                                                        {video.isFree === 1 && <Tag color="green">免费</Tag>}
                                                        <span className="video-duration">{video.duration || 0}分钟</span>
                                                    </div>
                                                    <Button
                                                        type={video.isFree === 1 || hasPurchased ? 'primary' : 'default'}
                                                        icon={video.isFree === 1 || hasPurchased ? <PlayCircleOutlined /> : <LockOutlined />}
                                                        onClick={() => handlePlayVideo(video)}
                                                        className="video-play-btn"
                                                    >
                                                        {video.isFree === 1 ? '试看' : (hasPurchased ? '播放' : '购买后观看')}
                                                    </Button>
                                                </List.Item>
                                            )}
                                        />
                                    </Collapse.Panel>
                                ))}
                            </Collapse>
                        ) : (
                            <div className="no-chapters">
                                <BookOutlined style={{ fontSize: 48, color: '#ccc' }} />
                                <p>暂无课程章节</p>
                            </div>
                        )}
                    </Card>
                </div>

                {/* 视频播放模态框 */}
                <Modal
                    title={
                        <div className="video-modal-title">
                            <PlayCircleOutlined /> {currentVideo?.title || '视频播放'}
                        </div>
                    }
                    open={playModalVisible}
                    onCancel={() => {
                        setPlayModalVisible(false)
                        setVideoUrl('')
                    }}
                    footer={null}
                    width={900}
                    centered
                    destroyOnClose
                    className="video-modal"
                >
                    {videoUrl && (
                        <div className="video-player-wrapper">
                            <video
                                src={videoUrl}
                                controls
                                autoPlay
                                className="video-player"
                            >
                                您的浏览器不支持视频播放
                            </video>
                            <div className="video-info-bar">
                                <span>* 视频通过 OSS 预签名 URL 播放，有效期 3600 秒</span>
                            </div>
                        </div>
                    )}
                </Modal>
            </Content>
        </Layout>
    )
}

export default CourseDetail
