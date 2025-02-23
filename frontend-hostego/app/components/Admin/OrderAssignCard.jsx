"use client";
import React, { useState } from "react";

const OrderAssignCard = ({ onSelectOrder, checkBoxVariant }) => {
    const [openCard, setOpenCard] = useState(false);
    const [isChecked, setIsChecked] = useState(false);

    const handleCheckboxChange = (e) => {
        setIsChecked(e.target.checked);
        onSelectOrder(e.target.checked);
    };

    return (
        <div className="rounded-xl bg-white p-3 border shadow-md transition-all hover:shadow-lg">
            <div className="flex items-center justify-between p-2 cursor-pointer" onClick={() => setOpenCard(!openCard)}>
                {/* Order Details */}
                <div className="w-full">
                    <p className="font-semibold text-sm text-gray-900">ID: e3571ddc-3b87-45ae-a115-e793e4bf3e24</p>
                    <p className="text-xs text-gray-600">₹ 112 &bull; 14 Feb, 8:12 PM</p>
                </div>

                {/* Checkbox */}
                {checkBoxVariant && <label onClick={(e) => e.stopPropagation()} className="cursor-pointer order-checkbox">
                    <input
                        type="checkbox"
                        checked={isChecked}
                        onChange={handleCheckboxChange}
                        className="w-5 h-5 cursor-pointer accent-[var(--primary-color)] order-checkbox rounded-md border-gray-300"
                    />
                </label>}
            </div>

            {/* Order Details Accordion */}
            {openCard && (
                <div className="mt-2 p-3 rounded-lg bg-[var(--bg-page-color)]">
                    {/* Product Images */}
                    <div className="flex gap-2 pb-3">
                        {[...Array(3)].map((_, index) => (
                            <div key={index} className="bg-gray-100 p-2 rounded-lg">
                                <img
                                    className="w-10 rounded-md"
                                    src="https://www.bigbasket.com/media/uploads/p/l/40015993_11-uncle-chips-spicy-treat.jpg"
                                    alt="Uncle Chips"
                                />
                            </div>
                        ))}
                    </div>

                    {/* Delivery Details */}
                    <div className="flex flex-col gap-3 text-xs text-gray-800">
                        <div className="flex flex-col">
                            <p className="text-gray-600">Deliver To</p>
                            <p className="font-semibold">Room no. 1115, Zakir-A, Chandigarh University</p>
                        </div>
                        <div className="flex flex-col">
                            <p className="text-gray-600">Delivery Partner Earning</p>
                            <p className="font-semibold text-[var(--primary-color)]">₹ 21</p>
                        </div>
                        <div className="flex flex-col">
                            <p className="text-gray-600">Order Status</p>
                            <p className="font-semibold">PENDING | ASSIGNED | DELIVERED</p>
                        </div>
                        <div className="flex flex-col">
                            <p className="text-gray-600">Order Placed</p>
                            <p className="font-semibold">Fri, 14 Feb'25, 8:12 PM</p>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default OrderAssignCard;
