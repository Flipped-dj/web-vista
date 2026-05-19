import React, { useState, useEffect } from 'react'
import { Layout, Card, Row, Col, Tag, Button, Pagination, Spin, Input, Avatar, Dropdown, Empty } from 'antd'
import { useNavigate } from 'react-router-dom'
import {
    SearchOutlined,
    UserOutlined,
    LogoutOutlined,
    SettingOutlined,
    TeamOutlined,
    RocketOutlined,
    PlayCircleOutlined
} from '@ant-design/icons'
import { getCourseList } from '../api/course'
import './CourseList.css'

const { Header, Content, Footer } = Layout

function CourseList() {
    const navigate = useNavigate()
    const [loading, setLoading] = useState(false)
    const [courses, setCourses] = useState([])
    const [total, setTotal] = useState(0)
    const [pageNum, setPageNum] = useState(1)
    const [pageSize] = useState(12)
    const [searchKey, setSearchKey] = useState('')

    const token = localStorage.getItem('token')
    const userInfo = token ? JSON.parse(localStorage.getItem('userInfo') || '{}') : null

    useEffect(() => {
        fetchCourses()
    }, [pageNum])

    const fetchCourses = async () => {
        setLoading(true)
        try {
            const data = await getCourseList({ pageNum, pageSize })
            setCourses(data.list || [])
            setTotal(data.total || 0)
        } catch (error) {
            console.error('Failed to load courses:', error)
        } finally {
            setLoading(false)
        }
    }

    const handlePageChange = (page) => {
        setPageNum(page)
    }

    const handleLogout = () => {
        localStorage.removeItem('token')
        localStorage.removeItem('userInfo')
        navigate('/login')
    }

    const userMenuItems = [
        { key: 'profile', icon: <UserOutlined />, label: '个人中心', onClick: () => navigate('/user') },
        userInfo?.role === 2 && { key: 'admin', icon: <SettingOutlined />, label: '管理后台', onClick: () => navigate('/admin') },
        { type: 'divider' },
        { key: 'logout', icon: <LogoutOutlined />, label: '退出登录', onClick: handleLogout }
    ].filter(Boolean)

    const filteredCourses = courses.filter(course =>
        !searchKey || course.title.toLowerCase().includes(searchKey.toLowerCase())
    )

    return (
        <Layout className="course-list-layout">
            {/* Header */}
            <Header className="course-list-header">
                <div className="header-container">
                    <div className="logo" onClick={() => navigate('/')}>
                        <RocketOutlined className="logo-icon" />
                        <span className="logo-text">GeekEdu</span>
                    </div>

                    <div className="header-search">
                        <Input
                            placeholder="搜索课程..."
                            prefix={<SearchOutlined />}
                            value={searchKey}
                            onChange={(e) => setSearchKey(e.target.value)}
                            className="search-input"
                        />
                    </div>

                    <div className="header-actions">
                        {token ? (
                            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
                                <div className="user-avatar-wrapper">
                                    <Avatar icon={<UserOutlined />} src={userInfo?.avatar} className="user-avatar" />
                                    <span className="user-name">{userInfo?.nickname || userInfo?.username}</span>
                                </div>
                            </Dropdown>
                        ) : (
                            <div className="auth-buttons">
                                <Button type="text" onClick={() => navigate('/login')}>登录</Button>
                                <Button type="primary" onClick={() => navigate('/register')}>注册</Button>
                            </div>
                        )}
                    </div>
                </div>
            </Header>

            <Content className="course-list-content">
                {/* Page Banner */}
                <div className="page-banner">
                    <h1>全部课程</h1>
                    <p>探索优质课程，开启学习之旅</p>
                </div>

                {/* Course Grid */}
                <div className="courses-container">
                    <Spin spinning={loading}>
                        {filteredCourses.length > 0 ? (
                            <>
                                <Row gutter={[24, 24]}>
                                    {filteredCourses.map(course => (
                                        <Col xs={24} sm={12} md={8} lg={6} key={course.id}>
                                            <Card
                                                hoverable
                                                className="course-card"
                                                onClick={() => navigate(`/courses/${course.id}`)}
                                                cover={
                                                    <div className="course-cover-wrapper">
                                                        <img
                                                            alt={course.title}
                                                            src={course.cover || `https://picsum.photos/seed/${course.id}/400/225`}
                                                            className="course-cover-img"
                                                            onError={(e) => {
                                                                e.target.src = `https://picsum.photos/seed/${course.id}/400/225`
                                                            }}
                                                        />
                                                        {course.price === 0 && (
                                                            <span className="free-badge">免费</span>
                                                        )}
                                                        <div className="course-overlay">
                                                            <PlayCircleOutlined className="play-icon" />
                                                        </div>
                                                    </div>
                                                }
                                            >
                                                <div className="course-content">
                                                    <h3 className="course-title">{course.title}</h3>
                                                    <p className="course-desc">{course.description || '暂无简介'}</p>
                                                    <div className="course-footer">
                                                        <div className="course-meta">
                                                            <TeamOutlined />
                                                            <span>{course.studentCount || 0}人学习</span>
                                                        </div>
                                                        <div className="course-price">
                                                            {course.price > 0 ? (
                                                                <span className="price">¥{course.price}</span>
                                                            ) : (
                                                                <span className="price free">免费</span>
                                                            )}
                                                        </div>
                                                    </div>
                                                </div>
                                            </Card>
                                        </Col>
                                    ))}
                                </Row>

                                {/* Pagination */}
                                <div className="pagination-wrapper">
                                    <Pagination
                                        current={pageNum}
                                        total={total}
                                        pageSize={pageSize}
                                        onChange={handlePageChange}
                                        showTotal={(t) => `共 ${t} 门课程`}
                                        showSizeChanger={false}
                                    />
                                </div>
                            </>
                        ) : (
                            <Empty description="暂无课程" className="empty-state" />
                        )}
                    </Spin>
                </div>
            </Content>

            {/* Footer */}
            <Footer className="course-list-footer">
                <div className="footer-content">
                    <div className="footer-brand">
                        <RocketOutlined /> GeekEdu
                    </div>
                    <div className="footer-copyright">
                        © 2026 GeekEdu. All rights reserved.
                    </div>
                </div>
            </Footer>
        </Layout>
    )
}

export default CourseList
