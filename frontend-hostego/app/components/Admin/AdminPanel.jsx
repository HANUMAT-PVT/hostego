"use client";
import { useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { Home, Package, Users, Settings, PackageOpenIcon, UserCircle, ShoppingBasket } from "lucide-react";
import OrderAssignment from "../../components/Admin/OrderAssignment";
import SidebarItem from "../../components/Admin/SidebarItem";
import OrdersList from "../../components/Admin/OrdersList";
import WalletPaymentVerfication from "./WalletPaymentVerfication";
import DeliveryPartnerManagement from "./DeliveryPartnerManagement";
import UserManager from "./UserManager";
import { useSelector } from "react-redux";
import ProductsManager from "./ProductsManager";
import ShopsManager from "./ShopsManager";


export default function AdminPanel() {
    const router = useRouter();
    const { userRoles } = useSelector(state => state.user)
    const searchParams = useSearchParams();
   
    // Get the current page from query params, default to 'dashboard'
    const currentPage = searchParams.get("page") || "order-assign";

    // Function to update query params
    const updatePage = (page) => {
        router.push(`?page=${page}`, { scroll: false });
    };


    function checkUserRole(roleName) {
       
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
                    {(checkUserRole("super_admin") || checkUserRole("order_manager")) && <SidebarItem
                        icon={<Home size={20} />}
                        text="Dashboard"
                        isActive={currentPage === "dashboard"}
                        onClick={() => updatePage("dashboard")}
                    />}
                    {(checkUserRole("super_admin") || checkUserRole("order_assign_manager")) && <SidebarItem
                        icon={<Package size={20} />}
                        text="Order Assign"
                        isActive={currentPage === "order-assign"}
                        onClick={() => updatePage("order-assign")}
                    />}
                    {(checkUserRole("super_admin") || checkUserRole("delivery_partner_manager")) && <SidebarItem
                        icon={<Users size={20} />}
                        text="Delivery Partners"
                        isActive={currentPage === "partners"}
                        onClick={() => updatePage("partners")}
                    />}
                    {(checkUserRole("super_admin") || checkUserRole("payments_manager")) && <SidebarItem
                        icon={<Settings size={20} />}
                        text="Payment Verification"
                        isActive={currentPage === "wallet_payment_verification"}
                        onClick={() => updatePage("wallet_payment_verification")}
                    />}
                    {(checkUserRole("super_admin") || checkUserRole("order_manager")) && <SidebarItem
                        icon={<PackageOpenIcon size={20} />}
                        text="Orders"
                        isActive={currentPage === "orders"}
                        onClick={() => updatePage("orders")}
                    />}
                    {(checkUserRole("super_admin") || checkUserRole("order_manager")) && <SidebarItem
                        icon={<UserCircle size={20} />}
                        text="Users"
                        isActive={currentPage === "users"}
                        onClick={() => updatePage("users")}
                    />}
                    {(checkUserRole("super_admin") || checkUserRole("order_manager")) && <SidebarItem
                        icon={<PackageOpenIcon size={20} />}
                        text="Products"
                        isActive={currentPage === "products"}
                        onClick={() => updatePage("products")}
                    />}
                     {(checkUserRole("super_admin") || checkUserRole("order_manager")) && <SidebarItem
                        icon={<ShoppingBasket size={20} />}
                        text="Shops"
                        isActive={currentPage === "shops"}
                        onClick={() => updatePage("shops")}
                    />}
                </nav>
            </aside>

            {/* Main Content */}
            <main className="flex-1 p-6">
                {currentPage === "dashboard" && <h1 className="text-2xl font-bold">📊 Dashboard</h1>}
                {currentPage === "order-assign" && (checkUserRole("super_admin") || checkUserRole("order_assign_manager")) && <OrderAssignment />}
                {currentPage === "partners" && (checkUserRole("super_admin") || checkUserRole("delivery_partner_manager")) && <DeliveryPartnerManagement />}
                {currentPage === "wallet_payment_verification" && (checkUserRole("super_admin") || checkUserRole("payments_manager")) && <WalletPaymentVerfication />}
                {currentPage === "orders" && (checkUserRole("super_admin") || checkUserRole("order_manager")) && <OrdersList />}
                {currentPage === "users" && (checkUserRole("super_admin") || checkUserRole("admin")) && <UserManager />}
                {currentPage === "products" && (checkUserRole("super_admin") || checkUserRole("admin")) && <ProductsManager />} 
                {currentPage === "shops" && (checkUserRole("super_admin") || checkUserRole("admin")) && <ShopsManager />} 
            </main>
        </div>
    );
}
