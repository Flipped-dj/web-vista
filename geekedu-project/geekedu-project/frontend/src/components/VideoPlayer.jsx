import React, { useState, useEffect } from 'react'
import ReactPlayer from 'react-player'
import { message, Spin } from 'antd'
import { getVideoPlayUrl } from '../api/course'
import './VideoPlayer.css'

const VideoPlayer = ({ videoId }) => {
    const [playUrl, setPlayUrl] = useState(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)

    useEffect(() => {
        if (!videoId) {
            setLoading(false)
            return
        }

        fetchVideoUrl()
    }, [videoId])

    const fetchVideoUrl = async () => {
        try {
            setLoading(true)
            setError(null)

            const response = await getVideoPlayUrl(videoId)

            // 准点要求：只有当返回200和签名URL时才渲染播放器
            if (response && response.playUrl) {
                setPlayUrl(response.playUrl)
            } else {
                setError('无法获取播放地址')
                message.error('无法获取播放地址')
            }
        } catch (err) {
            console.error('获取视频播放地址失败:', err)
            setError(err.response?.data?.message || '视频加载失败，请检查是否已购买该课程')
            message.error(err.response?.data?.message || '视频加载失败')
        } finally {
            setLoading(false)
        }
    }

    if (loading) {
        return (
            <div className="video-player-loading">
                <Spin size="large" tip="加载中..." />
            </div>
        )
    }

    if (error) {
        return (
            <div className="video-player-error">
                <p>{error}</p>
            </div>
        )
    }

    // 核心：只有当playUrl存在（200状态 + 签名URL）时才渲染播放器
    if (!playUrl) {
        return (
            <div className="video-player-empty">
                <p>请选择要播放的视频</p>
            </div>
        )
    }

    return (
        <div className="video-player-container">
            <ReactPlayer
                url={playUrl}
                controls
                width="100%"
                height="100%"
                playing={false}
                config={{
                    file: {
                        attributes: {
                            controlsList: 'nodownload',
                            onContextMenu: (e) => e.preventDefault()
                        }
                    }
                }}
                onError={(error) => {
                    console.error('播放器错误:', error)
                    message.error('视频播放失败')
                }}
            />
            <div className="video-player-info">
                <p className="video-url-notice">
                    📺 视频地址有效期：1小时（3600秒）
                </p>
            </div>
        </div>
    )
}

export default VideoPlayer
