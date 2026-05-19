import React, { useState } from 'react'
import { Form, Input, Button, Card, message, Tabs } from 'antd'
import { UserOutlined, LockOutlined, BookOutlined, SettingOutlined } from '@ant-design/icons'
import { useNavigate, Link } from 'react-router-dom'
import { login } from '../api/user'
import './Login.css'

function Login() {
    const navigate = useNavigate()
    const [loading, setLoading] = useState(false)
    const [activeTab, setActiveTab] = useState('student')

    const onFinish = async (values) => {
        setLoading(true)
        try {
            const data = await login(values)
            localStorage.setItem('token', data.token)
            localStorage.setItem('userInfo', JSON.stringify(data.userInfo))
            message.success('登录成功')

            // 根据角色跳转不同页面
            if (data.userInfo.role === 2) {
                navigate('/admin')
            } else {
                navigate('/')
            }
        } catch (error) {
            console.error('登录失败:', error)
        } finally {
            setLoading(false)
        }
    }

    const tabItems = [
        {
            key: 'student',
            label: (
                <span>
                    <BookOutlined />
                    学员登录
                </span>
            ),
            children: <LoginForm onFinish={onFinish} loading={loading} />
        },
        {
            key: 'admin',
            label: (
                <span>
                    <SettingOutlined />
                    管理员登录
                </span>
            ),
            children: <LoginForm onFinish={onFinish} loading={loading} />
        }
    ]

    return (
        <div className="login-container">
            <Card className="login-card">
                <h2 className="login-title">GeekEdu 在线教育平台</h2>
                <Tabs
                    activeKey={activeTab}
                    onChange={setActiveTab}
                    items={tabItems}
                    centered
                />
                <div className="login-footer">
                    还没有账号？<Link to="/register">立即注册</Link>
                </div>
            </Card>
        </div>
    )
}

// 登录表单组件
function LoginForm({ onFinish, loading }) {
    return (
        <Form
            name="login"
            onFinish={onFinish}
            autoComplete="off"
        >
            <Form.Item
                name="username"
                rules={[{ required: true, message: '请输入用户名' }]}
            >
                <Input
                    prefix={<UserOutlined />}
                    placeholder="用户名"
                    size="large"
                />
            </Form.Item>

            <Form.Item
                name="password"
                rules={[{ required: true, message: '请输入密码' }]}
            >
                <Input.Password
                    prefix={<LockOutlined />}
                    placeholder="密码"
                    size="large"
                />
            </Form.Item>

            <Form.Item>
                <Button
                    type="primary"
                    htmlType="submit"
                    size="large"
                    loading={loading}
                    block
                >
                    登录
                </Button>
            </Form.Item>
        </Form>
    )
}

export default Login
