"use client"

import { usePolling } from '../hooks/useSSE'
import { useEffect } from 'react'


const SSEHandler = ({userId}) => {
    usePolling(userId, (data) => {
        // data = { userId, type, msg }
        // alert("Order assingned")
    })

    return null
}

export default SSEHandler
