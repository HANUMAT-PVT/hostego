import React, { useState, useEffect } from 'react'
import axiosClient from '@/app/utils/axiosClient'

const NotificationManager = () => {
    const [notifications, setNotifications] = useState([])
    const [loading, setLoading] = useState(false)
    const [formData, setFormData] = useState({
        title: '',
        body: '',
        link: '',
        notification_image_url: ''
    })
    const [isFormOpen, setIsFormOpen] = useState(false)

    const fetchNotifications = async () => {
        const response = await axiosClient.get('/api/notifications')
        setNotifications(response?.data?.data)
    }

    // Mock data for demonstration - replace with actual API calls
    useEffect(() => {
        // Simulate fetching notifications
        fetchNotifications()
    }, [])

    const handleInputChange = (e) => {
        const { name, value } = e.target
        setFormData(prev => ({
            ...prev,
            [name]: value
        }))
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        setLoading(true)

        try {
            const response = await axiosClient.post('/api/notifications', formData)
            // Validate form data

            if (!formData.title.trim() || !formData.body.trim()) {
                // toast.error('Title and Body are required')
                return
            }
            setFormData({
                title: '',
                body: '',
                link: '',
                notification_image_url: ''
            })
            fetchNotifications()
            setIsFormOpen(false)
            // toast.success('Notification sent successfully!')
        } catch (error) {

            // toast.error('Failed to send notification')
            console.error('Error sending notification:', error)
        } finally {
            setLoading(false)
        }
    }



    const formatDate = (dateString) => {
        return new Date(dateString).toLocaleString()
    }

    return (
        <div className='flex flex-col gap-6 p-6 bg-gray-50 min-h-screen'>
            {/* Header */}
            <div className='flex justify-between items-center'>
                <div>
                    <h1 className='text-3xl font-bold text-gray-800'>Notification Manager</h1>
                    <p className='text-gray-600 mt-1'>Manage and send notifications to users</p>
                </div>
                <button
                    onClick={() => setIsFormOpen(true)}
                    className='bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-semibold transition-colors duration-200 flex items-center gap-2'
                >
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                    </svg>
                    Send Notification
                </button>
            </div>

            {/* Send Notification Form */}
            {isFormOpen && (
                <div className='bg-white rounded-xl shadow-lg p-6 border border-gray-200'>
                    <div className='flex justify-between items-center mb-6'>
                        <h2 className='text-2xl font-bold text-gray-800'>Send New Notification</h2>
                        <button
                            onClick={() => setIsFormOpen(false)}
                            className='text-gray-500 hover:text-gray-700 transition-colors'
                        >
                            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>

                    <form onSubmit={handleSubmit} className='space-y-6'>
                        <div className='grid grid-cols-1 md:grid-cols-2 gap-6'>
                            <div>
                                <label className='block text-sm font-medium text-gray-700 mb-2'>
                                    Title *
                                </label>
                                <input
                                    type="text"
                                    name="title"
                                    value={formData.title}
                                    onChange={handleInputChange}
                                    placeholder="Enter notification title"
                                    className='w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200'
                                    maxLength={255}
                                    required
                                />
                            </div>

                            <div>
                                <label className='block text-sm font-medium text-gray-700 mb-2'>
                                    Link (Optional)
                                </label>
                                <input
                                    type="url"
                                    name="link"
                                    value={formData.link}
                                    onChange={handleInputChange}
                                    placeholder="https://example.com"
                                    className='w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200'
                                    maxLength={255}
                                />
                            </div>
                        </div>

                        <div>
                            <label className='block text-sm font-medium text-gray-700 mb-2'>
                                Body *
                            </label>
                            <textarea
                                name="body"
                                value={formData.body}
                                onChange={handleInputChange}
                                placeholder="Enter notification message"
                                rows={4}
                                className='w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 resize-none'
                                maxLength={255}
                                required
                            />
                            <p className='text-sm text-gray-500 mt-1'>
                                {formData.body.length}/255 characters
                            </p>
                        </div>

                        <div>
                            <label className='block text-sm font-medium text-gray-700 mb-2'>
                                Image URL (Optional)
                            </label>
                            <input
                                type="url"
                                name="notification_image_url"
                                value={formData.notification_image_url}
                                onChange={handleInputChange}
                                placeholder="https://example.com/image.jpg"
                                className='w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200'
                                maxLength={255}
                            />
                        </div>

                        <div className='flex gap-4 pt-4'>
                            <button
                                type="submit"
                                disabled={loading}
                                className='bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white px-8 py-3 rounded-lg font-semibold transition-colors duration-200 flex items-center gap-2'
                            >
                                {loading ? (
                                    <>
                                        <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                                        Sending...
                                    </>
                                ) : (
                                    <>
                                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                                        </svg>
                                        Send Notification
                                    </>
                                )}
                            </button>
                            <button
                                type="button"
                                onClick={() => setIsFormOpen(false)}
                                className='bg-gray-200 hover:bg-gray-300 text-gray-700 px-8 py-3 rounded-lg font-semibold transition-colors duration-200'
                            >
                                Cancel
                            </button>
                        </div>
                    </form>
                </div>
            )}

            {/* Notifications List */}
            <div className='bg-white rounded-xl shadow-lg border border-gray-200'>
                <div className='p-6 border-b border-gray-200'>
                    <h2 className='text-2xl font-bold text-gray-800'>Notification History</h2>
                    <p className='text-gray-600 mt-1'>Recent notifications sent to users</p>
                </div>

                <div className='divide-y divide-gray-200'>
                    {notifications?.length === 0 ? (
                        <div className='p-8 text-center'>
                            <div className='w-16 h-16 mx-auto mb-4 bg-gray-100 rounded-full flex items-center justify-center'>
                                <svg className="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-5 5v-5zM4 19h6a2 2 0 002-2V7a2 2 0 00-2-2H4a2 2 0 00-2 2v10a2 2 0 002 2z" />
                                </svg>
                            </div>
                            <h3 className='text-lg font-medium text-gray-900 mb-2'>No notifications yet</h3>
                            <p className='text-gray-500'>Start by sending your first notification to users.</p>
                        </div>
                    ) : (
                        notifications?.map((notification) => (
                            <div key={notification.id} className='p-6 hover:bg-gray-50 transition-colors duration-200'>
                                <div className='flex items-start justify-between'>
                                    <div className='flex-1'>
                                        <div className='flex items-start gap-4'>
                                            {notification.notification_image_url && (
                                                <img
                                                    src={notification.notification_image_url}
                                                    alt="Notification"
                                                    className='w-12 h-12 rounded-lg object-cover flex-shrink-0'
                                                    onError={(e) => {
                                                        e.target.style.display = 'none'
                                                    }}
                                                />
                                            )}
                                            <div className='flex-1 min-w-0'>
                                                <h3 className='text-lg font-semibold text-gray-900 mb-1'>
                                                    {notification.title}
                                                </h3>
                                                <p className='text-gray-600 mb-2 line-clamp-2'>
                                                    {notification.body}
                                                </p>
                                                {notification.link && (
                                                    <a
                                                        href={notification.link}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                        className='text-blue-600 hover:text-blue-800 text-sm font-medium inline-flex items-center gap-1'
                                                    >
                                                        View Link
                                                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                                                        </svg>
                                                    </a>
                                                )}
                                                <div className='flex items-center gap-4 mt-3 text-sm text-gray-500'>
                                                    <span className='flex items-center gap-1'>
                                                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                                                        </svg>
                                                        Created: {formatDate(notification.created_at)}
                                                    </span>
                                                    <span className='flex items-center gap-1'>
                                                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                                                        </svg>
                                                        Updated: {formatDate(notification.updated_at)}
                                                    </span>
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                </div>
                            </div>
                        ))
                    )}
                </div>
            </div>
        </div>
    )
}

export default NotificationManager