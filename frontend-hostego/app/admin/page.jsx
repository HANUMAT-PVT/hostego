"use client"
import { Suspense } from 'react';
import AdminPanel from "../components/Admin/AdminPanel"
import HostegoLoader from "../components/HostegoLoader"

export default function Admin() {
    return (
        <Suspense fallback={<HostegoLoader />}>
            <AdminPanel />
        </Suspense>
    )
}