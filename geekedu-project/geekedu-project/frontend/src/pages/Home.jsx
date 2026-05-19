import React, { useState, useEffect } from 'react'
import { Layout, Button, Card, Row, Col, Avatar, Dropdown, Spin, Empty, Input, Badge } from 'antd'
import { useNavigate } from 'react-router-dom'
import {
    SearchOutlined,
    UserOutlined,
    LogoutOutlined,
    SettingOutlined,
    FireOutlined,
    RocketOutlined,
    TeamOutlined,
    BookOutlined,
    PlayCircleOutlined,
    SafetyCertificateOutlined,
    CustomerServiceOutlined
} from '@ant-design/icons'
import { getCourseList } from '../api/course'
import './Home.css'

const { Header, Content, Footer } = Layout

function Home() {
    const navigate = useNavigate()
    const [courses, setCourses] = useState([])
    const [loading, setLoading] = useState(false)
    const [searchKeyword, setSearchKeyword] = useState('')
    const token = localStorage.getItem('token')
    const userInfo = token ? JSON.parse(localStorage.getItem('userInfo') || '{}') : null

    useEffect(() => {
        fetchCourses()
    }, [])

    const fetchCourses = async () => {
        setLoading(true)
        try {
            const data = await getCourseList({ pageNum: 1, pageSize: 4 })
            setCourses(data.list || [])
        } catch (error) {
            console.error('Failed to load courses:', error)
        } finally {
            setLoading(false)
        }
    }

    const handleLogout = () => {
        localStorage.removeItem('token')
        localStorage.removeItem('userInfo')
        window.location.reload()
    }

    const userMenuItems = [
        {
            key: 'profile',
            icon: <UserOutlined />,
            label: '个人中心',
            onClick: () => navigate('/user')
        },
        userInfo?.role === 2 && {
            key: 'admin',
            icon: <SettingOutlined />,
            label: '管理后台',
            onClick: () => navigate('/admin')
        },
        { type: 'divider' },
        {
            key: 'logout',
            icon: <LogoutOutlined />,
            label: '退出登录',
            onClick: handleLogout
        }
    ].filter(Boolean)

    const handleSearch = () => {
        if (searchKeyword.trim()) {
            navigate(`/courses?keyword=${encodeURIComponent(searchKeyword)}`)
        } else {
            navigate('/courses')
        }
    }

    return (
        <Layout className="home-layout">
            {/* Navigation Header */}
            <Header className="home-header">
                <div className="header-container">
                    <div className="logo" onClick={() => navigate('/')}>
                        <RocketOutlined className="logo-icon" />
                        <span className="logo-text">GeekEdu</span>
                    </div>

                    <div className="header-search">
                        <Input
                            placeholder="搜索课程..."
                            prefix={<SearchOutlined />}
                            value={searchKeyword}
                            onChange={e => setSearchKeyword(e.target.value)}
                            onPressEnter={handleSearch}
                            className="search-input"
                        />
                    </div>

                    <div className="header-actions">
                        {token ? (
                            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
                                <div className="user-avatar-wrapper">
                                    <Avatar
                                        icon={<UserOutlined />}
                                        src={userInfo?.avatar}
                                        className="user-avatar"
                                    />
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

            <Content>
                {/* Hero Banner */}
                <section className="hero-banner">
                    <div className="hero-bg-pattern"></div>
                    <div className="hero-content">
                        <div className="hero-badge">在线教育平台</div>
                        <h1 className="hero-title">
                            学习改变未来<br />
                            <span className="highlight">GeekEdu</span> 助你成长
                        </h1>
                        <p className="hero-desc">
                            专业的在线学习平台，汇聚优质课程资源，随时随地开启你的学习之旅
                        </p>
                        <div className="hero-actions">
                            <Button
                                type="primary"
                                size="large"
                                className="hero-btn-primary"
                                onClick={() => navigate('/courses')}
                            >
                                <PlayCircleOutlined /> 开始学习
                            </Button>
                            {!token && (
                                <Button
                                    size="large"
                                    className="hero-btn-secondary"
                                    onClick={() => navigate('/register')}
                                >
                                    免费注册
                                </Button>
                            )}
                        </div>
                        <div className="hero-stats">
                            <div className="stat-item">
                                <span className="stat-number">1000+</span>
                                <span className="stat-label">注册学员</span>
                            </div>
                            <div className="stat-divider"></div>
                            <div className="stat-item">
                                <span className="stat-number">100+</span>
                                <span className="stat-label">精品课程</span>
                            </div>
                            <div className="stat-divider"></div>
                            <div className="stat-item">
                                <span className="stat-number">50+</span>
                                <span className="stat-label">专业讲师</span>
                            </div>
                        </div>
                    </div>
                </section>

                {/* Features Section */}
                <section className="features-section">
                    <div className="section-container">
                        <Row gutter={[32, 32]}>
                            <Col xs={24} sm={12} md={6}>
                                <div className="feature-item">
                                    <div className="feature-icon icon-purple">
                                        <BookOutlined />
                                    </div>
                                    <h3>优质课程</h3>
                                    <p>精选优质内容，系统化学习路径</p>
                                </div>
                            </Col>
                            <Col xs={24} sm={12} md={6}>
                                <div className="feature-item">
                                    <div className="feature-icon icon-pink">
                                        <TeamOutlined />
                                    </div>
                                    <h3>名师授课</h3>
                                    <p>行业专家亲授，经验倾囊相传</p>
                                </div>
                            </Col>
                            <Col xs={24} sm={12} md={6}>
                                <div className="feature-item">
                                    <div className="feature-icon icon-blue">
                                        <SafetyCertificateOutlined />
                                    </div>
                                    <h3>品质保障</h3>
                                    <p>严格审核机制，确保内容质量</p>
                                </div>
                            </Col>
                            <Col xs={24} sm={12} md={6}>
                                <div className="feature-item">
                                    <div className="feature-icon icon-green">
                                        <CustomerServiceOutlined />
                                    </div>
                                    <h3>贴心服务</h3>
                                    <p>7x24小时在线，随时答疑解惑</p>
                                </div>
                            </Col>
                        </Row>
                    </div>
                </section>

                {/* Hot Courses Section */}
                <section className="courses-section">
                    <div className="section-container">
                        <div className="section-header">
                            <div className="section-title">
                                <FireOutlined className="title-icon" />
                                <h2>热门课程</h2>
                            </div>
                            <Button type="link" onClick={() => navigate('/courses')}>
                                查看全部 →
                            </Button>
                        </div>

                        <Spin spinning={loading}>
                            {courses.length > 0 ? (
                                <Row gutter={[24, 24]}>
                                    {courses.slice(0, 4).map(course => (
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
                            ) : (
                                <Empty description="暂无课程" />
                            )}
                        </Spin>
                    </div>
                </section>

                {/* CTA Section */}
                <section className="cta-section">
                    <div className="cta-content">
                        <h2>准备好开始学习了吗？</h2>
                        <p>加入我们，与数千名学员一起成长</p>
                        <Button
                            type="primary"
                            size="large"
                            className="cta-btn"
                            onClick={() => token ? navigate('/courses') : navigate('/register')}
                        >
                            {token ? '浏览课程' : '立即加入'}
                        </Button>
                    </div>
                </section>
            </Content>

            {/* Footer */}
            <Footer className="home-footer">
                <div className="footer-content">
                    <div className="footer-brand">
                        <RocketOutlined /> GeekEdu
                    </div>
                    <div className="footer-links">
                        <a onClick={() => navigate('/courses')}>课程中心</a>
                        <a onClick={() => navigate('/about')}>关于我们</a>
                        <a onClick={() => navigate('/help')}>帮助中心</a>
                    </div>
                    <div className="footer-copyright">
                        © 2026 GeekEdu. All rights reserved.
                    </div>
                </div>
            </Footer>
        </Layout>
    )
}

export default Home
