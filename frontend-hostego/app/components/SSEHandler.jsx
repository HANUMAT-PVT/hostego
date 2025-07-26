"use client"

import { usePolling } from '../hooks/useSSE'
import { useEffect } from 'react'
import { subscribeToNotifications } from '../utils/webNotifications'


const SSEHandler = ({ userId }) => {
    // usePolling(userId, ({title,body,roles}) => {
    //     subscribeToNotifications(title,body)
    // })

    return null
}

export default SSEHandler
