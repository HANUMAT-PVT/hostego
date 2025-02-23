import React, { useState } from "react";
import OrderAssignCard from "./OrderAssignCard";
import { RotateCw } from "lucide-react";

const OrdersList = () => {
    const [loading, setLoading] = useState(false);

    const refreshData = () => {
        setLoading(true);
        // Simulating API call
        setTimeout(() => {
            setLoading(false);
            console.log("Orders refreshed!");
        }, 2000);
    };

    return (
        <div className="bg-[var(--bg-page-color)]">
            {/* Header with Refresh Button */}
            <div className="flex justify-between items-center mb-4">
                <h1 className="text-2xl font-bold text-gray-900">🚚 Orders </h1>
                <button
                    onClick={refreshData}
                    className="flex items-center  gap-2 px-4 py-2 bg-gray-100 text-gray-700 rounded-lg shadow hover:bg-gray-200 transition-all"
                    disabled={loading}
                >
                    <RotateCw size={18} className={`transition-transform ${loading ? "animate-spin" : ""}`} />
                    {loading ? "Refreshing..." : "Refresh"}
                </button>
            </div>

            {/* Orders List */}
            <div className="grid sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {[...Array(6)].map((_, index) => (
                    <OrderAssignCard key={index} />
                ))}
            </div>
        </div>
    );
};

export default OrdersList;
