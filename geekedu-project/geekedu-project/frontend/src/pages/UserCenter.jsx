import React, { useState, useEffect } from 'react'
import { Layout, Card, Row, Col, Tabs, Descriptions, Button, Table, Tag, message, Form, Input, Modal, Avatar, Upload, Statistic, Space, Empty } from 'antd'
import { useNavigate } from 'react-router-dom'
import {
    UserOutlined,
    LogoutOutlined,
    EditOutlined,
    HomeOutlined,
    RocketOutlined,
    LockOutlined,
    ShoppingCartOutlined,
    BookOutlined,
    SettingOutlined,
    CameraOutlined
} from '@ant-design/icons'
import { getUserInfo } from '../api/user'
import { getOrderList } from '../api/order'
import request from '../utils/request'
import './UserCenter.css'

const { Header, Content } = Layout

function UserCenter() {
    const navigate = useNavigate()
    const [userInfo, setUserInfo] = useState(null)
    const [orders, setOrders] = useState([])
    const [loading, setLoading] = useState(false)
    const [editModalVisible, setEditModalVisible] = useState(false)
    const [passwordModalVisible, setPasswordModalVisible] = useState(false)
    const [avatarUrl, setAvatarUrl] = useState('')
    const [uploading, setUploading] = useState(false)
    const [form] = Form.useForm()
    const [passwordForm] = Form.useForm()

    useEffect(() => {
        fetchUserInfo()
        fetchOrders()
    }, [])

    const fetchUserInfo = async () => {
        try {
            const data = await getUserInfo()
            setUserInfo(data)
            setAvatarUrl(data.avatar || '')
        } catch (error) {
            console.error('获取用户信息失败:', error)
        }
    }

    const fetchOrders = async () => {
        setLoading(true)
        try {
            const data = await getOrderList()
            setOrders(data.list || [])
        } catch (error) {
            console.error('获取订单列表失败:', error)
        } finally {
            setLoading(false)
        }
    }

    const handleLogout = () => {
        localStorage.removeItem('token')
        localStorage.removeItem('userInfo')
        message.success('退出成功')
        navigate('/login')
    }

    const handleEditProfile = () => {
        form.setFieldsValue({
            nickname: userInfo?.nickname || ''
        })
        setAvatarUrl(userInfo?.avatar || '')
        setEditModalVisible(true)
    }

    const handleAvatarUpload = async (info) => {
        const file = info.file
        const formData = new FormData()
        formData.append('file', file)

        setUploading(true)
        try {
            const token = localStorage.getItem('token')
            const response = await fetch('/api/v1/upload', {
                method: 'POST',
                headers: { 'Authorization': `Bearer ${token}` },
                body: formData
            })
            const result = await response.json()
            if (result.code === 0) {
                const ossUrl = result.data.url
                // 直接使用原始URL保存，不需要签名
                setAvatarUrl(ossUrl)
                message.success('头像上传成功')
            } else {
                message.error(result.message || '上传失败')
            }
        } catch (error) {
            console.error('上传错误:', error)
            message.error('上传失败，请重试')
        } finally {
            setUploading(false)
        }
    }

    const handleSaveProfile = async (values) => {
        try {
            await request.put('/user/profile', {
                nickname: values.nickname,
                avatar: avatarUrl || ''
            })
            message.success('保存成功')
            const newUserInfo = { ...userInfo, nickname: values.nickname, avatar: avatarUrl }
            localStorage.setItem('userInfo', JSON.stringify(newUserInfo))
            setUserInfo(newUserInfo)
            setEditModalVisible(false)
            fetchUserInfo()
        } catch (error) {
            console.error('保存失败:', error)
        }
    }

    const handleChangePassword = async (values) => {
        try {
            await request.put('/user/password', {
                oldPassword: values.oldPassword,
                newPassword: values.newPassword
            })
            message.success('密码修改成功，请重新登录')
            setPasswordModalVisible(false)
            handleLogout()
        } catch (error) {
            console.error('修改密码失败:', error)
        }
    }

    const paidOrders = orders.filter(o => o.status === 1)
    const totalSpent = paidOrders.reduce((sum, o) => sum + (o.amount || 0), 0)

    const columns = [
        {
            title: '订单号',
            dataIndex: 'orderNo',
            key: 'orderNo',
            width: 200,
            ellipsis: true
        },
        {
            title: '课程ID',
            dataIndex: 'courseId',
            key: 'courseId',
            width: 100
        },
        {
            title: '金额',
            dataIndex: 'amount',
            key: 'amount',
            width: 120,
            render: (amount) => <span style={{ color: '#f5222d', fontWeight: 600 }}>¥{amount?.toFixed(2)}</span>
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status) => {
                const statusMap = {
                    0: { text: '待支付', color: 'warning' },
                    1: { text: '已支付', color: 'success' },
                }
                const { text, color } = statusMap[status] || { text: '未知', color: 'default' }
                return <Tag color={color}>{text}</Tag>
            }
        },
        {
            title: '创建时间',
            dataIndex: 'createTime',
            key: 'createTime',
            width: 180
        }
    ]

    const getRoleName = (role) => {
        const roles = { 0: '学员', 1: '讲师', 2: '管理员' }
        return roles[role] || '未知'
    }

    const getRoleColor = (role) => {
        const colors = { 0: '#52c41a', 1: '#1890ff', 2: '#f5222d' }
        return colors[role] || '#999'
    }

    return (
        <Layout className="user-center-layout">
            <Header className="user-center-header">
                <div className="header-container">
                    <div className="logo" onClick={() => navigate('/')}>
                        <RocketOutlined className="logo-icon" />
                        <span className="logo-text">GeekEdu</span>
                    </div>
                    <div className="header-actions">
                        <Button type="text" icon={<HomeOutlined />} onClick={() => navigate('/')}>首页</Button>
                        {userInfo?.role === 2 && (
                            <Button type="text" icon={<SettingOutlined />} onClick={() => navigate('/admin')}>管理后台</Button>
                        )}
                        <Button type="text" danger icon={<LogoutOutlined />} onClick={handleLogout}>退出登录</Button>
                    </div>
                </div>
            </Header>

            <Content className="user-center-content">
                {/* 用户信息卡片 */}
                <Card className="user-profile-card">
                    <Row gutter={32} align="middle">
                        <Col flex="none">
                            <div className="avatar-container">
                                <Avatar
                                    size={120}
                                    src={userInfo?.avatar}
                                    icon={<UserOutlined />}
                                    className="user-avatar"
                                />
                            </div>
                        </Col>
                        <Col flex="auto">
                            <div className="user-info-main">
                                <h2 className="username">{userInfo?.nickname || userInfo?.username || '用户'}</h2>
                                <p className="user-account">账号: {userInfo?.username}</p>
                                <Tag color={getRoleColor(userInfo?.role)} className="role-tag">
                                    {getRoleName(userInfo?.role)}
                                </Tag>
                            </div>
                        </Col>
                        <Col flex="none">
                            <Space size="middle">
                                <Button type="primary" icon={<EditOutlined />} size="large" onClick={handleEditProfile}>
                                    编辑资料
                                </Button>
                                <Button icon={<LockOutlined />} size="large" onClick={() => setPasswordModalVisible(true)}>
                                    修改密码
                                </Button>
                            </Space>
                        </Col>
                    </Row>
                </Card>

                {/* 统计卡片 */}
                <Row gutter={24} className="stats-row">
                    <Col xs={24} sm={8}>
                        <Card className="stat-card">
                            <Statistic
                                title="我的订单"
                                value={orders.length}
                                prefix={<ShoppingCartOutlined />}
                                valueStyle={{ color: '#1890ff', fontSize: 32 }}
                            />
                        </Card>
                    </Col>
                    <Col xs={24} sm={8}>
                        <Card className="stat-card">
                            <Statistic
                                title="已购课程"
                                value={paidOrders.length}
                                prefix={<BookOutlined />}
                                valueStyle={{ color: '#52c41a', fontSize: 32 }}
                            />
                        </Card>
                    </Col>
                    <Col xs={24} sm={8}>
                        <Card className="stat-card">
                            <Statistic
                                title="累计消费"
                                value={totalSpent}
                                prefix="¥"
                                precision={2}
                                valueStyle={{ color: '#f5222d', fontSize: 32 }}
                            />
                        </Card>
                    </Col>
                </Row>

                {/* 订单列表 */}
                <Card className="orders-card" title={<span><ShoppingCartOutlined /> 我的订单</span>}>
                    {orders.length > 0 ? (
                        <Table
                            columns={columns}
                            dataSource={orders}
                            rowKey="id"
                            loading={loading}
                            pagination={{
                                pageSize: 10,
                                showTotal: (total) => `共 ${total} 条订单`
                            }}
                        />
                    ) : (
                        <Empty description="暂无订单" />
                    )}
                </Card>
            </Content>

            {/* 编辑资料弹窗 */}
            <Modal
                title="编辑个人资料"
                open={editModalVisible}
                onCancel={() => setEditModalVisible(false)}
                footer={null}
                width={500}
            >
                <Form form={form} onFinish={handleSaveProfile} layout="vertical" size="large">
                    <div className="avatar-upload-section">
                        <Upload
                            name="file"
                            showUploadList={false}
                            accept="image/*"
                            customRequest={({ file, onSuccess }) => {
                                handleAvatarUpload({ file })
                                onSuccess()
                            }}
                        >
                            <div className="avatar-upload-wrapper">
                                <Avatar
                                    size={100}
                                    src={avatarUrl}
                                    icon={<UserOutlined />}
                                    className={uploading ? 'uploading' : ''}
                                />
                                <div className="avatar-upload-mask">
                                    <CameraOutlined />
                                </div>
                            </div>
                        </Upload>
                        <p className="upload-tip">点击更换头像</p>
                    </div>
                    <Form.Item
                        name="nickname"
                        label="昵称"
                        rules={[{ required: true, message: '请输入昵称' }]}
                    >
                        <Input placeholder="请输入昵称" />
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit" block loading={uploading}>
                            保存修改
                        </Button>
                    </Form.Item>
                </Form>
            </Modal>

            {/* 修改密码弹窗 */}
            <Modal
                title="修改密码"
                open={passwordModalVisible}
                onCancel={() => setPasswordModalVisible(false)}
                footer={null}
                width={450}
            >
                <Form form={passwordForm} onFinish={handleChangePassword} layout="vertical" size="large">
                    <Form.Item
                        name="oldPassword"
                        label="原密码"
                        rules={[{ required: true, message: '请输入原密码' }]}
                    >
                        <Input.Password placeholder="请输入原密码" />
                    </Form.Item>
                    <Form.Item
                        name="newPassword"
                        label="新密码"
                        rules={[
                            { required: true, message: '请输入新密码' },
                            { min: 6, message: '密码至少6位' }
                        ]}
                    >
                        <Input.Password placeholder="请输入新密码" />
                    </Form.Item>
                    <Form.Item
                        name="confirmPassword"
                        label="确认密码"
                        dependencies={['newPassword']}
                        rules={[
                            { required: true, message: '请确认新密码' },
                            ({ getFieldValue }) => ({
                                validator(_, value) {
                                    if (!value || getFieldValue('newPassword') === value) {
                                        return Promise.resolve()
                                    }
                                    return Promise.reject(new Error('两次密码不一致'))
                                }
                            })
                        ]}
                    >
                        <Input.Password placeholder="请再次输入新密码" />
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit" block>确认修改</Button>
                    </Form.Item>
                </Form>
            </Modal>
        </Layout>
    )
}

export default UserCenter
