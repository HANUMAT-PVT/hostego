"use client";
import { useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { Home, Package, Users, Settings, PackageOpenIcon, UserCircle, ShoppingBasket, CreditCard, Utensils } from "lucide-react";
import OrderAssignment from "../../components/Admin/OrderAssignment";
import SidebarItem from "../../components/Admin/SidebarItem";
import OrdersList from "../../components/Admin/OrdersList";
import WalletPaymentVerfication from "./WalletPaymentVerfication";
import DeliveryPartnerManagement from "./DeliveryPartnerManagement";
import UserManager from "./UserManager";
import { useSelector } from "react-redux";
import ProductsManager from "./ProductsManager";
import ShopsManager from "./ShopsManager";
import Dashboard from "./Dashboard";
import axiosClient from "@/app/utils/axiosClient";
import DeliveryPartnerPaymentManager from "./DeliveryPartnerPaymentManager";
import CuMessManager from "./CuMessManager";

export default function AdminPanel() {
    const router = useRouter();
    const { userRoles } = useSelector(state => state.user)
    const searchParams = useSearchParams();
    const [dashboardStats, setDashboardStats] = useState({ product_stats: [], overall_stats: [] })

    // Get the current page from query params, default to 'dashboard'
    const currentPage = searchParams.get("page") || "order-assign";

    // Function to update query params
    const updatePage = (page) => {
        router.push(`?page=${page}`, { scroll: false });
    };

    useEffect(() => {
        fetchDashboardStatus()
    }, [])

    const fetchDashboardStatus = async () => {
        try {
            let { data } = await axiosClient.get(`/api/order/order-items`);

            setDashboardStats({ product_stats: data.product_stats ? data.product_stats : [], overall_stats: data.overall_stats ? data.overall_stats : [] })
        } catch (error) {

        }
    }

    function checkUserRole(roleName) {

        if (userRoles.length === 0) {
            return;
        }
        const role = userRoles.find(userRole => userRole?.role?.role_name === roleName);
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
                    {(checkUserRole("super_admin") || checkUserRole("admin")) && <SidebarItem
                        icon={<Home size={20} />}
                        text="Dashboard"
                        isActive={currentPage === "dashboard"}
                        onClick={() => updatePage("dashboard")}
                    />}
                    {(checkUserRole("super_admin") || checkUserRole("order_assign_manager") || checkUserRole("admin")) && <SidebarItem
                        icon={<Package size={20} />}
                        text="Order Assign"
                        isActive={currentPage === "order-assign"}
                        onClick={() => updatePage("order-assign")}
                    />}
                    {(checkUserRole("super_admin") || checkUserRole("delivery_partner_manager") || checkUserRole("admin")) && <SidebarItem
                        icon={<Users size={20} />}
                        text="Delivery Partners"
                        isActive={currentPage === "partners"}
                        onClick={() => updatePage("partners")}
                    />}
                    {(checkUserRole("super_admin") || checkUserRole("payments_manager") || checkUserRole("admin")) && <SidebarItem
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
                    {(checkUserRole("super_admin") || checkUserRole("admin")) && <SidebarItem
                        icon={<UserCircle size={20} />}
                        text="Users"
                        isActive={currentPage === "users"}
                        onClick={() => updatePage("users")}
                    />}
                    {(checkUserRole("super_admin") || checkUserRole("admin") || checkUserRole("inventory_manager")) && <SidebarItem
                        icon={<PackageOpenIcon size={20} />}
                        text="Products"
                        isActive={currentPage === "products"}
                        onClick={() => updatePage("products")}
                    />}
                    {(checkUserRole("super_admin") || checkUserRole("admin") || checkUserRole("inventory_manager")) && <SidebarItem
                        icon={<ShoppingBasket size={20} />}
                        text="Shops"
                        isActive={currentPage === "shops"}
                        onClick={() => updatePage("shops")}
                    />}
                    {(checkUserRole("super_admin") || checkUserRole("payments_manager") || checkUserRole("admin")) && <SidebarItem
                        icon={<CreditCard size={20} />}
                        text="Delivery Partner Payments"
                        isActive={currentPage === "delivery_partner_payment"}
                        onClick={() => updatePage("delivery_partner_payment")}
                    />}
                    {(checkUserRole("super_admin") || checkUserRole("admin") || checkUserRole("inventory_manager")) && <SidebarItem
                        icon={<Utensils size={20} />}
                        text="CU Mess"
                        isActive={currentPage === "cu_mess"}
                        onClick={() => updatePage("cu_mess")}
                    />}
                </nav>
            </aside>

            {/* Main Content */}
            <main className="flex-1 p-6">
                {currentPage === "dashboard" && checkUserRole("super_admin") && <Dashboard dashboardStats={dashboardStats} />}
                {currentPage === "order-assign" && (checkUserRole("super_admin") || checkUserRole("order_assign_manager") || checkUserRole("admin")) && <OrderAssignment />}
                {currentPage === "orders" && (checkUserRole("super_admin") || checkUserRole("order_manager")) && <OrdersList />}
                {currentPage === "partners" && (checkUserRole("super_admin") || checkUserRole("delivery_partner_manager") || checkUserRole("admin")) && <DeliveryPartnerManagement />}
                {currentPage === "wallet_payment_verification" && (checkUserRole("super_admin") || checkUserRole("payments_manager") || checkUserRole("admin")) && <WalletPaymentVerfication />}
                {currentPage === "users" && (checkUserRole("super_admin") || checkUserRole("admin")) && <UserManager />}
                {currentPage === "products" && (checkUserRole("super_admin") || checkUserRole("admin") || checkUserRole("inventory_manager")) && <ProductsManager />}
                {currentPage === "shops" && (checkUserRole("super_admin") || checkUserRole("admin") || checkUserRole("inventory_manager")) && <ShopsManager />}
                {currentPage === "cu_mess" && (checkUserRole("super_admin") || checkUserRole("admin") || checkUserRole("inventory_manager")) && <CuMessManager />}
                {currentPage === "delivery_partner_payment" && (checkUserRole("super_admin") || checkUserRole("payments_manager") || checkUserRole("admin")) && <DeliveryPartnerPaymentManager />}
            </main>
        </div>
    );
}
