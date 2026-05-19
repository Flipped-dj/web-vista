import React, { useState, useEffect } from 'react'
import { Layout, Menu, Table, Card, Button, Modal, Form, Input, InputNumber, Select, message, Tag, Space, Statistic, Row, Col, Upload, Progress } from 'antd'
import {
    DashboardOutlined,
    BookOutlined,
    UserOutlined,
    AppstoreOutlined,
    PlusOutlined,
    LogoutOutlined,
    UploadOutlined,
    VideoCameraOutlined,
    PictureOutlined,
    InboxOutlined,
    ReloadOutlined,
    ShoppingCartOutlined,
    DollarOutlined,
    HomeOutlined,
    EyeOutlined
} from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import request from '../utils/request'
import './Admin.css'

const { Header, Sider, Content } = Layout
const { Dragger } = Upload

function Admin() {
    const navigate = useNavigate()
    const [collapsed, setCollapsed] = useState(false)
    const [activeMenu, setActiveMenu] = useState('dashboard')
    const [userInfo, setUserInfo] = useState(null)

    useEffect(() => {
        const info = localStorage.getItem('userInfo')
        if (info) {
            const parsed = JSON.parse(info)
            if (parsed.role !== 2) {
                message.error('无管理员权限')
                navigate('/')
                return
            }
            setUserInfo(parsed)
        } else {
            navigate('/login')
        }
    }, [navigate])

    const handleLogout = () => {
        localStorage.removeItem('token')
        localStorage.removeItem('userInfo')
        navigate('/login')
    }

    const menuItems = [
        { key: 'dashboard', icon: <DashboardOutlined />, label: '控制台' },
        { key: 'courses', icon: <BookOutlined />, label: '课程管理' },
        { key: 'users', icon: <UserOutlined />, label: '用户管理' },
        { key: 'categories', icon: <AppstoreOutlined />, label: '分类管理' },
        { type: 'divider' },
        { key: 'home', icon: <HomeOutlined />, label: '返回首页' },
        { key: 'browse', icon: <EyeOutlined />, label: '浏览课程' },
    ]

    const handleMenuClick = ({ key }) => {
        if (key === 'home') {
            navigate('/')
        } else if (key === 'browse') {
            navigate('/courses')
        } else {
            setActiveMenu(key)
        }
    }

    const renderContent = () => {
        switch (activeMenu) {
            case 'dashboard':
                return <Dashboard />
            case 'courses':
                return <CourseManagement />
            case 'users':
                return <UserManagement />
            case 'categories':
                return <CategoryManagement />
            default:
                return <Dashboard />
        }
    }

    return (
        <Layout style={{ minHeight: '100vh' }}>
            <Sider collapsible collapsed={collapsed} onCollapse={setCollapsed}>
                <div className="admin-logo">
                    {collapsed ? 'GE' : 'GeekEdu 管理'}
                </div>
                <Menu
                    theme="dark"
                    selectedKeys={[activeMenu]}
                    items={menuItems}
                    onClick={handleMenuClick}
                />
            </Sider>
            <Layout>
                <Header className="admin-header">
                    <span>欢迎，{userInfo?.nickname || '管理员'}</span>
                    <Space>
                        <Button type="link" onClick={() => navigate('/')}>
                            <HomeOutlined /> 首页
                        </Button>
                        <Button type="link" danger icon={<LogoutOutlined />} onClick={handleLogout}>
                            退出登录
                        </Button>
                    </Space>
                </Header>
                <Content className="admin-content">
                    {renderContent()}
                </Content>
            </Layout>
        </Layout>
    )
}

// 控制台
function Dashboard() {
    const [stats, setStats] = useState({ users: 0, courses: 0, orders: 0, income: 0 })
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        const fetchStats = async () => {
            try {
                setLoading(true)
                // 获取用户总数
                const usersData = await request.get('/admin/users', { params: { page: 1, pageSize: 1 } })
                // 获取课程总数
                const coursesData = await request.get('/courses', { params: { pageNum: 1, pageSize: 1 } })
                // 获取订单总数
                const ordersData = await request.get('/admin/orders', { params: { pageNum: 1, pageSize: 1 } })

                setStats({
                    users: usersData.total || 0,
                    courses: coursesData.total || 0,
                    orders: ordersData.total || 0,
                    income: ordersData.totalIncome || 0
                })
            } catch (error) {
                console.error('获取统计数据失败:', error)
            } finally {
                setLoading(false)
            }
        }
        fetchStats()
    }, [])

    const handleRefresh = () => {
        const fetchStats = async () => {
            try {
                setLoading(true)
                const usersData = await request.get('/admin/users', { params: { page: 1, pageSize: 1 } })
                const coursesData = await request.get('/courses', { params: { pageNum: 1, pageSize: 1 } })
                const ordersData = await request.get('/admin/orders', { params: { pageNum: 1, pageSize: 1 } })

                setStats({
                    users: usersData.total || 0,
                    courses: coursesData.total || 0,
                    orders: ordersData.total || 0,
                    income: ordersData.totalIncome || 0
                })
                message.success('数据已刷新')
            } catch (error) {
                console.error('刷新失败:', error)
            } finally {
                setLoading(false)
            }
        }
        fetchStats()
    }

    const [recentOrders, setRecentOrders] = useState([])

    useEffect(() => {
        const fetchRecentOrders = async () => {
            try {
                const data = await request.get('/admin/orders', { params: { pageNum: 1, pageSize: 5 } })
                setRecentOrders(data.list || [])
            } catch (error) {
                console.error('获取订单失败:', error)
            }
        }
        fetchRecentOrders()
    }, [])

    const orderColumns = [
        { title: '订单号', dataIndex: 'orderNo', width: 180 },
        { title: '课程', dataIndex: 'courseTitle', ellipsis: true },
        { title: '用户', dataIndex: 'username' },
        { title: '金额', dataIndex: 'amount', render: v => `¥${v?.toFixed(2) || '0.00'}` },
        { title: '状态', dataIndex: 'status', render: s => s === 1 ? <Tag color="green">已支付</Tag> : <Tag color="orange">待支付</Tag> },
        { title: '创建时间', dataIndex: 'createTime', width: 170 }
    ]

    return (
        <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
                <h2 style={{ margin: 0 }}>📊 控制台</h2>
                <Button icon={<ReloadOutlined />} onClick={handleRefresh} loading={loading}>
                    刷新
                </Button>
            </div>
            <Row gutter={[24, 24]}>
                <Col xs={24} sm={12} lg={6}>
                    <Card loading={loading} className="stat-card stat-card-users">
                        <Statistic
                            title="用户总数"
                            value={stats.users}
                            prefix={<UserOutlined />}
                            valueStyle={{ color: '#6366f1', fontSize: 32, fontWeight: 700 }}
                        />
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card loading={loading} className="stat-card stat-card-courses">
                        <Statistic
                            title="课程总数"
                            value={stats.courses}
                            prefix={<BookOutlined />}
                            valueStyle={{ color: '#10b981', fontSize: 32, fontWeight: 700 }}
                        />
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card loading={loading} className="stat-card stat-card-orders">
                        <Statistic
                            title="订单总数"
                            value={stats.orders}
                            prefix={<ShoppingCartOutlined />}
                            valueStyle={{ color: '#f59e0b', fontSize: 32, fontWeight: 700 }}
                        />
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card loading={loading} className="stat-card stat-card-income">
                        <Statistic
                            title="总收入"
                            value={stats.income}
                            prefix={<DollarOutlined />}
                            precision={2}
                            suffix="元"
                            valueStyle={{ color: '#ef4444', fontSize: 32, fontWeight: 700 }}
                        />
                    </Card>
                </Col>
            </Row>

            <Card style={{ marginTop: 24 }} title={<span>📋 最近订单</span>}>
                <Table
                    columns={orderColumns}
                    dataSource={recentOrders}
                    rowKey="id"
                    pagination={false}
                    locale={{ emptyText: '暂无订单数据' }}
                />
            </Card>
        </div>
    )
}

// 课程管理
function CourseManagement() {
    const [courses, setCourses] = useState([])
    const [categories, setCategories] = useState([])
    const [loading, setLoading] = useState(false)
    const [modalVisible, setModalVisible] = useState(false)
    const [editModalVisible, setEditModalVisible] = useState(false)
    const [videoModalVisible, setVideoModalVisible] = useState(false)
    const [chapterModalVisible, setChapterModalVisible] = useState(false)
    const [currentCourse, setCurrentCourse] = useState(null)
    const [coverUrl, setCoverUrl] = useState('')
    const [coverChanged, setCoverChanged] = useState(false)  // 是否重新上传了封面
    const [uploading, setUploading] = useState(false)
    const [chapters, setChapters] = useState([])
    const [form] = Form.useForm()
    const [editForm] = Form.useForm()
    const [videoForm] = Form.useForm()
    const [chapterForm] = Form.useForm()

    const fetchCourses = async () => {
        setLoading(true)
        try {
            const data = await request.get('/courses', { params: { pageNum: 1, pageSize: 100 } })
            setCourses(data.list || [])
        } catch (error) {
            console.error(error)
        } finally {
            setLoading(false)
        }
    }

    const fetchCategories = async () => {
        try {
            const data = await request.get('/categories')
            // 扁平化分类列表
            const flatCategories = []
            const flatten = (list) => {
                list.forEach(item => {
                    flatCategories.push({ id: item.id, name: item.name })
                    if (item.children && item.children.length > 0) {
                        flatten(item.children)
                    }
                })
            }
            flatten(data || [])
            setCategories(flatCategories)
        } catch (error) {
            console.error('获取分类失败:', error)
        }
    }

    useEffect(() => {
        fetchCourses()
        fetchCategories()
    }, [])

    // 上传封面图片到OSS
    const handleCoverUpload = async (info) => {
        const { file } = info
        if (file.status === 'uploading') {
            setUploading(true)
            return
        }

        const formData = new FormData()
        formData.append('file', file.originFileObj || file)

        try {
            setUploading(true)
            const token = localStorage.getItem('token')
            const response = await fetch('/api/v1/upload', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`
                },
                body: formData
            })
            const result = await response.json()
            if (result.code === 0) {
                setCoverUrl(result.data.url)
                setCoverChanged(true)  // 标记封面已更新
                message.success('封面上传成功！')
            } else {
                message.error(result.message || '上传失败')
            }
        } catch (error) {
            message.error('上传失败: ' + error.message)
        } finally {
            setUploading(false)
        }
    }

    // 创建课程
    const handleCreate = async (values) => {
        if (!coverUrl) {
            message.error('请先上传封面图片')
            return
        }
        try {
            await request.post('/courses', {
                ...values,
                coverUrl: coverUrl
            })
            message.success('课程创建成功')
            setModalVisible(false)
            form.resetFields()
            setCoverUrl('')
            setCoverChanged(false)
            fetchCourses()
        } catch (error) {
            console.error(error)
        }
    }

    // 编辑课程
    const handleEdit = (record) => {
        setCurrentCourse(record)
        setCoverUrl(record.cover || '')
        setCoverChanged(false)  // 重置封面更新标记
        editForm.setFieldsValue({
            title: record.title,
            description: record.description,
            price: record.price,
            categoryId: record.categoryId,
            status: record.status
        })
        setEditModalVisible(true)
    }

    const handleUpdate = async (values) => {
        try {
            // 只有在用户重新上传封面时才发送coverUrl
            const updateData = { ...values }
            if (coverChanged && coverUrl) {
                updateData.coverUrl = coverUrl
            }
            await request.put(`/courses/${currentCourse.id}`, updateData)
            message.success('课程更新成功')
            setEditModalVisible(false)
            editForm.resetFields()
            setCoverUrl('')
            setCoverChanged(false)
            setCurrentCourse(null)
            fetchCourses()
        } catch (error) {
            message.error('更新失败')
        }
    }

    // 删除课程
    const handleDelete = (record) => {
        Modal.confirm({
            title: '确认删除',
            content: `确定要删除课程"${record.title}"吗？此操作不可恢复。`,
            okText: '删除',
            okType: 'danger',
            cancelText: '取消',
            onOk: async () => {
                try {
                    await request.delete(`/courses/${record.id}`)
                    message.success('课程删除成功')
                    fetchCourses()
                } catch (error) {
                    message.error('删除失败')
                }
            }
        })
    }

    // 管理章节
    const handleManageChapters = async (course) => {
        setCurrentCourse(course)
        try {
            const data = await request.get(`/courses/${course.id}`)
            setChapters(data.chapters || [])
        } catch (error) {
            setChapters([])
        }
        setChapterModalVisible(true)
    }

    // 添加章节
    const handleAddChapter = async (values) => {
        try {
            await request.post(`/courses/${currentCourse.id}/chapters`, {
                title: values.chapterTitle,
                sort: chapters.length + 1
            })
            message.success('章节添加成功')
            chapterForm.resetFields()
            // 刷新章节列表
            const data = await request.get(`/courses/${currentCourse.id}`)
            setChapters(data.chapters || [])
        } catch (error) {
            message.error('添加章节失败')
        }
    }

    // 上传视频到章节
    const handleVideoUpload = async (course) => {
        setCurrentCourse(course)
        // 先加载章节列表
        try {
            const data = await request.get(`/courses/${course.id}`)
            setChapters(data.chapters || [])
        } catch (error) {
            setChapters([])
        }
        setVideoModalVisible(true)
    }

    const handleVideoSubmit = async (values) => {
        console.log('提交视频表单:', values)
        if (!values.videoUrl) {
            message.error('请先上传视频文件')
            return
        }
        if (!values.chapterId) {
            message.error('请选择视频所属章节')
            return
        }
        try {
            const payload = {
                title: values.title,
                videoUrl: values.videoUrl,
                duration: parseInt(values.duration) || 0,
                isFree: parseInt(values.isFree) || 0,
                chapterId: parseInt(values.chapterId)
            }
            console.log('发送请求:', payload)
            await request.post(`/courses/${currentCourse.id}/videos`, payload)
            message.success('视频添加成功！')
            setVideoModalVisible(false)
            videoForm.resetFields()
        } catch (error) {
            console.error('添加视频失败:', error)
            // 错误已在 request 拦截器中显示
        }
    }

    const columns = [
        { title: 'ID', dataIndex: 'id', width: 60 },
        {
            title: '封面',
            dataIndex: 'cover',
            width: 100,
            render: (url) => <img src={url} alt="封面" style={{ width: 80, height: 45, objectFit: 'cover', borderRadius: 4 }} />
        },
        { title: '课程名称', dataIndex: 'title', ellipsis: true },
        { title: '价格', dataIndex: 'price', width: 80, render: (v) => `¥${v}` },
        { title: '学员数', dataIndex: 'studentCount', width: 80 },
        {
            title: '状态',
            dataIndex: 'status',
            width: 80,
            render: (status) => {
                const statusMap = {
                    0: { color: 'orange', text: '待审核' },
                    1: { color: 'green', text: '已上架' },
                    2: { color: 'red', text: '已下架' }
                }
                const s = statusMap[status] || statusMap[0]
                return <Tag color={s.color}>{s.text}</Tag>
            }
        },
        {
            title: '操作',
            width: 280,
            render: (_, record) => (
                <Space size="small">
                    <Button size="small" type="link" onClick={() => handleEdit(record)}>编辑</Button>
                    <Button size="small" type="link" onClick={() => handleManageChapters(record)}>章节</Button>
                    <Button size="small" type="link" icon={<VideoCameraOutlined />} onClick={() => handleVideoUpload(record)}>视频</Button>
                    <Button size="small" type="link" danger onClick={() => handleDelete(record)}>删除</Button>
                </Space>
            )
        }
    ]

    return (
        <div>
            <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
                <h2>课程管理</h2>
                <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalVisible(true)}>
                    发布课程
                </Button>
            </div>
            <Table
                columns={columns}
                dataSource={courses}
                rowKey="id"
                loading={loading}
            />

            {/* 发布课程弹窗 */}
            <Modal
                title="发布新课程"
                open={modalVisible}
                onCancel={() => {
                    setModalVisible(false)
                    setCoverUrl('')
                    setCoverChanged(false)
                    form.resetFields()
                }}
                footer={null}
                width={600}
            >
                <Form form={form} onFinish={handleCreate} layout="vertical">
                    <Form.Item name="title" label="课程标题" rules={[{ required: true, message: '请输入课程标题' }]}>
                        <Input placeholder="请输入课程标题" />
                    </Form.Item>
                    <Form.Item name="description" label="课程简介" rules={[{ required: true, message: '请输入课程简介' }]}>
                        <Input.TextArea rows={3} placeholder="请输入课程简介" />
                    </Form.Item>

                    {/* 封面图片上传 */}
                    <Form.Item label="封面图片" required>
                        <Dragger
                            name="file"
                            multiple={false}
                            accept="image/*"
                            showUploadList={false}
                            customRequest={({ file, onSuccess }) => {
                                handleCoverUpload({ file })
                                onSuccess()
                            }}
                        >
                            {coverUrl ? (
                                <div>
                                    <img src={coverUrl} alt="封面" style={{ maxWidth: '100%', maxHeight: 200 }} />
                                    <p style={{ marginTop: 8, color: '#52c41a' }}>
                                        <PictureOutlined /> 封面已上传到阿里云OSS
                                    </p>
                                </div>
                            ) : (
                                <div>
                                    <p className="ant-upload-drag-icon">
                                        <InboxOutlined />
                                    </p>
                                    <p className="ant-upload-text">点击或拖拽上传封面图片</p>
                                    <p className="ant-upload-hint">支持 JPG、PNG 格式，图片将上传至阿里云OSS</p>
                                </div>
                            )}
                        </Dragger>
                        {uploading && <Progress percent={50} status="active" />}
                    </Form.Item>

                    <Form.Item name="price" label="价格" rules={[{ required: true, message: '请输入价格' }]}>
                        <InputNumber min={0} precision={2} style={{ width: '100%' }} placeholder="请输入价格" addonAfter="元" />
                    </Form.Item>
                    <Form.Item name="categoryId" label="分类" rules={[{ required: true, message: '请选择分类' }]}>
                        <Select placeholder="请选择分类">
                            {categories.map(cat => (
                                <Select.Option key={cat.id} value={cat.id}>{cat.name}</Select.Option>
                            ))}
                        </Select>
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit" block loading={uploading}>
                            发布课程
                        </Button>
                    </Form.Item>
                </Form>
            </Modal>

            {/* 上传视频弹窗 */}
            <Modal
                title={`上传视频 - ${currentCourse?.title || ''}`}
                open={videoModalVisible}
                onCancel={() => {
                    setVideoModalVisible(false)
                    videoForm.resetFields()
                    setChapters([])
                }}
                footer={null}
                width={600}
            >
                <VideoUploadForm
                    form={videoForm}
                    onFinish={handleVideoSubmit}
                    courseId={currentCourse?.id}
                    chapters={chapters}
                />
            </Modal>

            {/* 编辑课程弹窗 */}
            <Modal
                title="编辑课程"
                open={editModalVisible}
                onCancel={() => {
                    setEditModalVisible(false)
                    setCoverUrl('')
                    setCoverChanged(false)
                    editForm.resetFields()
                    setCurrentCourse(null)
                }}
                footer={null}
                width={600}
            >
                <Form form={editForm} onFinish={handleUpdate} layout="vertical">
                    <Form.Item name="title" label="课程标题" rules={[{ required: true, message: '请输入课程标题' }]}>
                        <Input placeholder="请输入课程标题" />
                    </Form.Item>
                    <Form.Item name="description" label="课程简介">
                        <Input.TextArea rows={3} placeholder="请输入课程简介" />
                    </Form.Item>
                    <Form.Item label="封面图片">
                        <Dragger
                            name="file"
                            multiple={false}
                            accept="image/*"
                            showUploadList={false}
                            customRequest={({ file, onSuccess }) => {
                                handleCoverUpload({ file })
                                onSuccess()
                            }}
                        >
                            {coverUrl ? (
                                <div>
                                    <img src={coverUrl} alt="封面" style={{ maxWidth: '100%', maxHeight: 150 }} />
                                    <p style={{ marginTop: 8, color: '#52c41a' }}>点击更换封面</p>
                                </div>
                            ) : (
                                <div>
                                    <p className="ant-upload-drag-icon"><InboxOutlined /></p>
                                    <p className="ant-upload-text">点击上传新封面</p>
                                </div>
                            )}
                        </Dragger>
                    </Form.Item>
                    <Row gutter={16}>
                        <Col span={12}>
                            <Form.Item name="price" label="价格">
                                <InputNumber min={0} precision={2} style={{ width: '100%' }} addonAfter="元" />
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item name="status" label="状态">
                                <Select>
                                    <Select.Option value={0}>待审核</Select.Option>
                                    <Select.Option value={1}>已上架</Select.Option>
                                    <Select.Option value={2}>已下架</Select.Option>
                                </Select>
                            </Form.Item>
                        </Col>
                    </Row>
                    <Form.Item name="categoryId" label="分类">
                        <Select placeholder="请选择分类">
                            {categories.map(cat => (
                                <Select.Option key={cat.id} value={cat.id}>{cat.name}</Select.Option>
                            ))}
                        </Select>
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit" block>保存修改</Button>
                    </Form.Item>
                </Form>
            </Modal>

            {/* 章节管理弹窗 */}
            <Modal
                title={`章节管理 - ${currentCourse?.title || ''}`}
                open={chapterModalVisible}
                onCancel={() => {
                    setChapterModalVisible(false)
                    setChapters([])
                    chapterForm.resetFields()
                }}
                footer={null}
                width={700}
            >
                <div style={{ marginBottom: 16 }}>
                    <Form form={chapterForm} layout="inline" onFinish={handleAddChapter}>
                        <Form.Item name="chapterTitle" rules={[{ required: true, message: '请输入章节标题' }]}>
                            <Input placeholder="输入章节标题" style={{ width: 300 }} />
                        </Form.Item>
                        <Form.Item>
                            <Button type="primary" htmlType="submit" icon={<PlusOutlined />}>添加章节</Button>
                        </Form.Item>
                    </Form>
                </div>
                <Table
                    dataSource={chapters}
                    rowKey="id"
                    size="small"
                    pagination={false}
                    columns={[
                        {
                            title: '序号',
                            width: 60,
                            render: (_, __, index) => index + 1
                        },
                        { title: '章节标题', dataIndex: 'title' },
                        {
                            title: '视频数',
                            dataIndex: 'videos',
                            width: 80,
                            render: (videos) => videos?.length || 0
                        }
                    ]}
                    locale={{ emptyText: '暂无章节，请添加' }}
                />
            </Modal>
        </div>
    )
}

// 视频上传表单组件
function VideoUploadForm({ form, onFinish, courseId, chapters = [] }) {
    const [videoUrl, setVideoUrl] = useState('')
    const [uploading, setUploading] = useState(false)
    const [uploadProgress, setUploadProgress] = useState(0)

    const handleVideoUpload = async (info) => {
        const { file } = info
        const formData = new FormData()
        formData.append('file', file.originFileObj || file)

        try {
            setUploading(true)
            setUploadProgress(0)

            const token = localStorage.getItem('token')

            // 使用XMLHttpRequest来获取上传进度
            const xhr = new XMLHttpRequest()
            xhr.upload.addEventListener('progress', (e) => {
                if (e.lengthComputable) {
                    const percent = Math.round((e.loaded / e.total) * 100)
                    setUploadProgress(percent)
                }
            })

            const response = await new Promise((resolve, reject) => {
                xhr.onload = () => {
                    if (xhr.status === 200) {
                        resolve(JSON.parse(xhr.responseText))
                    } else {
                        reject(new Error('上传失败'))
                    }
                }
                xhr.onerror = () => reject(new Error('网络错误'))
                xhr.open('POST', '/api/v1/upload')
                xhr.setRequestHeader('Authorization', `Bearer ${token}`)
                xhr.send(formData)
            })

            if (response.code === 0) {
                setVideoUrl(response.data.url)
                form.setFieldsValue({ videoUrl: response.data.url })
                message.success('视频上传成功！')
            } else {
                message.error(response.message || '上传失败')
            }
        } catch (error) {
            message.error('上传失败: ' + error.message)
        } finally {
            setUploading(false)
        }
    }

    const handleSubmit = (values) => {
        onFinish({ ...values, videoUrl })
    }

    return (
        <Form form={form} onFinish={handleSubmit} layout="vertical">
            <Form.Item name="title" label="视频标题" rules={[{ required: true, message: '请输入视频标题' }]}>
                <Input placeholder="请输入视频标题" />
            </Form.Item>

            <Form.Item
                name="chapterId"
                label="所属章节"
                rules={[{ required: true, message: '请选择章节' }]}
                extra={chapters.length === 0 ? <span style={{ color: '#ff4d4f' }}>请先在"章节"管理中添加章节</span> : null}
            >
                <Select placeholder="请选择章节" disabled={chapters.length === 0}>
                    {chapters.map(ch => (
                        <Select.Option key={ch.id} value={ch.id}>{ch.title}</Select.Option>
                    ))}
                </Select>
            </Form.Item>

            {/* 视频文件上传 */}
            <Form.Item label="视频文件" required>
                <Dragger
                    name="file"
                    multiple={false}
                    accept="video/*"
                    showUploadList={false}
                    customRequest={({ file, onSuccess }) => {
                        handleVideoUpload({ file })
                        onSuccess()
                    }}
                >
                    {videoUrl ? (
                        <div>
                            <p style={{ color: '#52c41a', fontSize: 16 }}>
                                <VideoCameraOutlined /> 视频已上传
                            </p>
                            <p style={{ color: '#999', fontSize: 12, wordBreak: 'break-all' }}>
                                {videoUrl.length > 60 ? videoUrl.slice(0, 60) + '...' : videoUrl}
                            </p>
                        </div>
                    ) : (
                        <div>
                            <p className="ant-upload-drag-icon">
                                <VideoCameraOutlined style={{ fontSize: 48, color: '#1890ff' }} />
                            </p>
                            <p className="ant-upload-text">点击或拖拽上传视频文件</p>
                            <p className="ant-upload-hint">支持 MP4、AVI、MOV 等格式</p>
                        </div>
                    )}
                </Dragger>
                {uploading && <Progress percent={uploadProgress} status="active" style={{ marginTop: 8 }} />}
            </Form.Item>

            <Form.Item name="videoUrl" hidden>
                <Input />
            </Form.Item>

            <Form.Item name="duration" label="视频时长（分钟）">
                <InputNumber min={0} style={{ width: '100%' }} placeholder="请输入视频时长" />
            </Form.Item>

            <Form.Item name="isFree" label="是否免费试看" initialValue={0}>
                <Select>
                    <Select.Option value={0}>付费观看</Select.Option>
                    <Select.Option value={1}>免费试看</Select.Option>
                </Select>
            </Form.Item>

            <Form.Item>
                <Button type="primary" htmlType="submit" block disabled={!videoUrl || uploading}>
                    添加视频
                </Button>
            </Form.Item>
        </Form>
    )
}

// 用户管理
function UserManagement() {
    const [users, setUsers] = useState([])
    const [loading, setLoading] = useState(false)
    const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 })

    const fetchUsers = async (page = 1, pageSize = 10) => {
        setLoading(true)
        try {
            const data = await request.get('/admin/users', { params: { page, pageSize } })
            setUsers(data.list || [])
            setPagination(prev => ({ ...prev, current: page, total: data.total || 0 }))
        } catch (error) {
            console.error('获取用户列表失败:', error)
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        fetchUsers()
    }, [])

    const handleTableChange = (pag) => {
        fetchUsers(pag.current, pag.pageSize)
    }

    const columns = [
        { title: 'ID', dataIndex: 'id', width: 60 },
        { title: '用户名', dataIndex: 'username' },
        { title: '昵称', dataIndex: 'nickname' },
        {
            title: '角色',
            dataIndex: 'role',
            render: (role) => {
                const roleMap = {
                    0: { color: 'blue', text: '学员' },
                    1: { color: 'green', text: '讲师' },
                    2: { color: 'red', text: '管理员' }
                }
                const r = roleMap[role] || roleMap[0]
                return <Tag color={r.color}>{r.text}</Tag>
            }
        },
        {
            title: '注册时间',
            dataIndex: 'createdAt',
            render: (v) => v ? v.replace('T', ' ').split('.')[0] : '-'
        }
    ]

    return (
        <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
                <h2 style={{ margin: 0 }}>用户管理</h2>
                <Button icon={<ReloadOutlined />} onClick={() => fetchUsers(pagination.current, pagination.pageSize)}>
                    刷新
                </Button>
            </div>
            <Table
                columns={columns}
                dataSource={users}
                rowKey="id"
                loading={loading}
                pagination={pagination}
                onChange={handleTableChange}
            />
        </div>
    )
}

// 分类管理
function CategoryManagement() {
    const [categories, setCategories] = useState([])
    const [loading, setLoading] = useState(false)
    const [modalVisible, setModalVisible] = useState(false)
    const [editingCategory, setEditingCategory] = useState(null)
    const [form] = Form.useForm()

    const fetchCategories = async () => {
        setLoading(true)
        try {
            const data = await request.get('/categories')
            // Flatten for table display with courseCount and sort
            const flatList = []
            const flatten = (list, parentName = '') => {
                list.forEach(item => {
                    flatList.push({
                        id: item.id,
                        name: item.name,
                        parentId: item.parentId || 0,
                        courseCount: item.courseCount || 0,
                        sort: item.sort || 0,
                        parentName: parentName
                    })
                    if (item.children && item.children.length > 0) {
                        flatten(item.children, item.name)
                    }
                })
            }
            flatten(data || [])
            // Sort by parentId then sort
            flatList.sort((a, b) => {
                if (a.parentId === b.parentId) return a.sort - b.sort
                return a.parentId - b.parentId
            })
            setCategories(flatList)
        } catch (error) {
            console.error('获取分类列表失败:', error)
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        fetchCategories()
    }, [])

    const handleAdd = () => {
        setEditingCategory(null)
        form.resetFields()
        form.setFieldsValue({ parentId: 0, sort: 0 })
        setModalVisible(true)
    }

    const handleEdit = (record) => {
        setEditingCategory(record)
        form.setFieldsValue({
            name: record.name,
            parentId: record.parentId,
            sort: record.sort
        })
        setModalVisible(true)
    }

    const handleSubmit = async (values) => {
        try {
            if (editingCategory) {
                await request.put(`/admin/categories/${editingCategory.id}`, values)
                message.success('修改成功')
            } else {
                await request.post('/admin/categories', values)
                message.success('添加成功')
            }
            setModalVisible(false)
            form.resetFields()
            setEditingCategory(null)
            fetchCategories()
        } catch (error) {
            message.error(editingCategory ? '修改失败' : '添加失败')
        }
    }

    const handleDelete = async (id) => {
        // Check if has children
        const hasChildren = categories.some(c => c.parentId === id)
        if (hasChildren) {
            message.warning('该分类下有子分类，请先删除子分类')
            return
        }
        try {
            await request.delete(`/admin/categories/${id}`)
            message.success('删除成功')
            fetchCategories()
        } catch (error) {
            message.error('删除失败')
        }
    }

    const topCategories = categories.filter(c => c.parentId === 0)

    const columns = [
        { title: 'ID', dataIndex: 'id', width: 80 },
        {
            title: '分类名称',
            dataIndex: 'name',
            render: (text, record) => (
                <span>
                    {record.parentId > 0 && <span style={{ color: '#999', marginRight: 8 }}>└─</span>}
                    <Tag color={record.parentId === 0 ? 'blue' : 'cyan'}>{text}</Tag>
                </span>
            )
        },
        {
            title: '父级分类',
            dataIndex: 'parentId',
            width: 150,
            render: (parentId) => {
                if (parentId === 0) return <Tag color="purple">顶级分类</Tag>
                const parent = categories.find(c => c.id === parentId)
                return parent?.name || '-'
            }
        },
        {
            title: '课程数',
            dataIndex: 'courseCount',
            width: 100,
            render: (v) => <Tag color={v > 0 ? 'green' : 'default'}>{v}</Tag>
        },
        {
            title: '排序',
            dataIndex: 'sort',
            width: 80
        },
        {
            title: '操作',
            width: 150,
            render: (_, record) => (
                <Space>
                    <Button size="small" type="link" onClick={() => handleEdit(record)}>编辑</Button>
                    <Button size="small" type="link" danger onClick={() => handleDelete(record.id)}>删除</Button>
                </Space>
            )
        }
    ]

    return (
        <div>
            <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <h2 style={{ margin: 0 }}>📂 分类管理</h2>
                <Space>
                    <Button icon={<ReloadOutlined />} onClick={fetchCategories}>刷新</Button>
                    <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
                        添加分类
                    </Button>
                </Space>
            </div>
            <Card>
                <Table
                    columns={columns}
                    dataSource={categories}
                    rowKey="id"
                    loading={loading}
                    pagination={false}
                    size="middle"
                />
            </Card>
            <Modal
                title={editingCategory ? '编辑分类' : '添加分类'}
                open={modalVisible}
                onCancel={() => {
                    setModalVisible(false)
                    setEditingCategory(null)
                    form.resetFields()
                }}
                footer={null}
            >
                <Form form={form} onFinish={handleSubmit} layout="vertical">
                    <Form.Item name="name" label="分类名称" rules={[{ required: true, message: '请输入分类名称' }]}>
                        <Input placeholder="请输入分类名称" />
                    </Form.Item>
                    <Form.Item name="parentId" label="父级分类">
                        <Select>
                            <Select.Option value={0}>无 (顶级分类)</Select.Option>
                            {topCategories.filter(c => !editingCategory || c.id !== editingCategory.id).map(c => (
                                <Select.Option key={c.id} value={c.id}>{c.name}</Select.Option>
                            ))}
                        </Select>
                    </Form.Item>
                    <Form.Item name="sort" label="排序 (数字越小越靠前)">
                        <InputNumber min={0} style={{ width: '100%' }} placeholder="0" />
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit" block>
                            {editingCategory ? '保存修改' : '添加分类'}
                        </Button>
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    )
}

export default Admin
