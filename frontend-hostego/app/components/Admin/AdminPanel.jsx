"use client";
import { useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { Home, Package, Users, Settings, PackageOpenIcon, UserCircle } from "lucide-react";
import OrderAssignment from "../../components/Admin/OrderAssignment";
import SidebarItem from "../../components/Admin/SidebarItem";
import OrdersList from "../../components/Admin/OrdersList";
import WalletPaymentVerfication from "./WalletPaymentVerfication";
import DeliveryPartnerManagement from "./DeliveryPartnerManagement";
import UserManager from "./UserManager";
import { useSelector } from "react-redux";


export default function AdminPanel() {
    const router = useRouter();
    const { userRoles } = useSelector(state => state.user)
    const searchParams = useSearchParams();
    console.log(userRoles)
    // Get the current page from query params, default to 'dashboard'
    const currentPage = searchParams.get("page") || "order-assign";

    // Function to update query params
    const updatePage = (page) => {
        router.push(`?page=${page}`, { scroll: false });
    };


    function checkUserRole(roleName) {
        console.log(roleName, "roleName")
        if (userRoles.length === 0) {
            return;
        }
        const role = userRoles.find(userRole => userRole?.role?.role_name === roleName);
        if (!role) {
            router.push("/home");
        }
        if (role) {
            return true
        }

    }

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
                        text="Payment Verification"
                        isActive={currentPage === "wallet_payment_verification"}
                        onClick={() => updatePage("wallet_payment_verification")}
                    />
                    <SidebarItem
                        icon={<PackageOpenIcon size={20} />}
                        text="Orders"
                        isActive={currentPage === "orders"}
                        onClick={() => updatePage("orders")}
                    />
                    <SidebarItem
                        icon={<UserCircle size={20} />}
                        text="Users"
                        isActive={currentPage === "users"}
                        onClick={() => updatePage("users")}
                    />
                </nav>
            </aside>

            {/* Main Content */}
            <main className="flex-1 p-6">
                {currentPage === "dashboard" && <h1 className="text-2xl font-bold">📊 Dashboard</h1>}
<<<<<<< Updated upstream
                {currentPage === "order-assign" && <OrderAssignment />}
                {currentPage === "partners" && <DeliveryPartnerManagement />}
                {currentPage === "wallet_payment_verification" && <WalletPaymentVerfication />}
                {currentPage === "orders" && <OrdersList />}
                {currentPage === "users" && <UserManager />}
=======
                {currentPage === "order-assign" && (checkUserRole("super_admin") || checkUserRole("order_assign_manager")) && <OrderAssignment />}
                {currentPage === "partners" && (checkUserRole("super_admin") || checkUserRole("delivery_partner_manager")) && <DeliveryPartnerManagement />}
                {currentPage === "wallet_payment_verification" && (checkUserRole("super_admin") || checkUserRole("payment_verification_manager")) && <WalletPaymentVerfication />}
                {currentPage === "orders" && (checkUserRole("super_admin") || checkUserRole("order_manager")) && <OrdersList />}
>>>>>>> Stashed changes
            </main>
        </div>
    );
}
