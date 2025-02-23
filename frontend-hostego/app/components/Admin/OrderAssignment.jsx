"use client";
import { useState } from "react";
import { RotateCw } from "lucide-react";
import OrderAssignCard from "../Admin/OrderAssignCard";
import DeliveryPartnerCard from "../Delivery-Partner/DeliveryPartnerCard";
import HostegoButton from "../../components/HostegoButton"

export default function OrderAssignment({  }) {
    const [selectedOrder, setSelectedOrder] = useState(null);
    const [selectedPartner, setSelectedPartner] = useState(null);
    const [loading, setLoading] = useState(false);

    // Dummy function for API call
    const refreshData = async () => {
        setLoading(true);
        try {
            // Simulate API call
            await new Promise((resolve) => setTimeout(resolve, 1500));
            alert("Data refreshed successfully!");
        } catch (error) {
            console.error("Error refreshing data:", error);
        } finally {
            setLoading(false);
        }
    };

    const handleAssign = () => {
        if (!selectedOrder || !selectedPartner) {
            alert("Please select both an order and a delivery partner.");
            return;
        }
        alert(`Order ${selectedOrder} assigned to Partner ${selectedPartner}`);
        setSelectedOrder(null);
        setSelectedPartner(null);
    };

    return (
        <div className="bg-[var(--bg-page-color)] flex flex-col p-6 relative">
            {/* Header with Refresh Button */}
            <div className="flex justify-between items-center mb-4">
                <h1 className="text-2xl font-bold text-gray-900">🚚 Delivery Order Assignment</h1>
                <button
                    onClick={refreshData}
                    className="flex items-center gap-2 px-4 py-2 bg-gray-100 text-gray-700 rounded-lg shadow hover:bg-gray-200 transition-all"
                    disabled={loading}
                >
                    <RotateCw size={18} className={`transition-transform ${loading ? "animate-spin" : ""}`} />
                    {loading ? "Refreshing..." : "Refresh"}
                </button>
            </div>

            <div className="flex gap-6 justify-between">
                {/* Pending Orders List */}
                <div className="h-[85vh] w-1/2 rounded-xl overflow-auto flex flex-col gap-4 border p-4 bg-white shadow-lg">
                    <h2 className="text-lg font-semibold mb-2 text-gray-800">📦 Pending Orders</h2>
                    {[...Array(4)].map((_, index) => (
                        <OrderAssignCard
                            checkBoxVariant={true}
                            key={index}
                            onSelectOrder={(checked) => setSelectedOrder(checked ? `Order-${index + 1}` : null)}
                        />
                    ))}
                </div>

                {/* Active Delivery Partners List */}
                <div className="h-[85vh] w-1/2 rounded-xl overflow-auto flex flex-col gap-4 border p-4 bg-white shadow-lg">
                    <h2 className="text-lg font-semibold mb-2 text-gray-800">🚀 Available Delivery Partners</h2>
                    {[...Array(4)].map((_, index) => (
                        <DeliveryPartnerCard
                            key={index}
                            partner={{ id: `Partner-${index + 1}`, name: `Partner ${index + 1}`, location: "CU Campus", phone: "8264121428" }}
                            onSelect={(id) => setSelectedPartner(id)}
                            isSelected={selectedPartner === `Partner-${index + 1}`}
                        />
                    ))}
                </div>
            </div>

            {/* Fixed Assign Button */}
            <div className="fixed bottom-4 left-1/2 transform -translate-x-1/2 w-full max-w-md">
                <HostegoButton
                    text="Assign Order"
                    onClick={handleAssign}
                    className="w-full px-6 py-3 text-white font-semibold bg-[var(--primary-color)] shadow-lg hover:opacity-90 transition-all"
                />
            </div>
        </div>
    );
}
