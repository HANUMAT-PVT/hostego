"use client";
import { useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { Home, Package, Users, Settings, PackageOpenIcon } from "lucide-react";
import OrderAssignment from "../../components/Admin/OrderAssignment";
import SidebarItem from "../../components/Admin/SidebarItem";
import OrdersList from "../../components/Admin/OrdersList";

export default function AdminPanel() {
    const router = useRouter();
    const searchParams = useSearchParams();

    // Get the current page from query params, default to 'dashboard'
    const currentPage = searchParams.get("page") || "order-assign";

    // Function to update query params
    const updatePage = (page) => {
        router.push(`?page=${page}`, { scroll: false });
    };

    return (
        <div className="flex h-screen bg-[var(--bg-page-color)]">
            {/* Sidebar */}
            <aside className="w-64 bg-white shadow-md p-4 flex flex-col gap-6">
                <h2 className="text-xl font-bold text-center text-[var(--primary-color)]">
                    🚀 Admin Panel
                </h2>
                <nav className="flex flex-col gap-4">
                    <SidebarItem
                        icon={<Home size={20} />}
                        text="Dashboard"
                        isActive={currentPage === "dashboard"}
                        onClick={() => updatePage("dashboard")}
                    />
                    <SidebarItem
                        icon={<Package size={20} />}
                        text="Order Assign"
                        isActive={currentPage === "order-assign"}
                        onClick={() => updatePage("order-assign")}
                    />
                    <SidebarItem
                        icon={<Users size={20} />}
                        text="Delivery Partners"
                        isActive={currentPage === "partners"}
                        onClick={() => updatePage("partners")}
                    />
                    <SidebarItem
                        icon={<Settings size={20} />}
                        text="Settings"
                        isActive={currentPage === "settings"}
                        onClick={() => updatePage("settings")}
                    />
                    <SidebarItem
                        icon={<PackageOpenIcon size={20} />}
                        text="Orders"
                        isActive={currentPage === "orders"}
                        onClick={() => updatePage("orders")}
                    />
                </nav>
            </aside>

            {/* Main Content */}
            <main className="flex-1 p-6">
                {currentPage === "dashboard" && <h1 className="text-2xl font-bold">📊 Dashboard</h1>}
                {currentPage === "order-assign" && <OrderAssignment />}
                {currentPage === "partners" && <h1 className="text-2xl font-bold">👥 Delivery Partners</h1>}
                {currentPage === "tr" && <h1 className="text-2xl font-bold">⚙️ Settings</h1>}
                {currentPage === "orders" && <OrdersList />}
            </main>
        </div>
    );
}
