"use client";
import React, { useState } from "react";

const DeliveryPartnerCard = ({ partner = { name: "John Doe", location: "Zakir-A,1115,Chandigarh University" }, onSelect, isSelected }) => {
    const [openCard, setOpenCard] = useState(false);
    const [isChecked, setIsChecked] = useState(false);

    return (
        <div className="rounded-xl bg-white p-4 border shadow-md transition-all hover:shadow-lg">
            <div className="flex gap-4 items-center cursor-pointer" onClick={() => setOpenCard(!openCard)}>
                {/* Checkbox wrapped in a label */}
                <label onClick={(e) => e.stopPropagation()} className="cursor-pointer">
                    <input
                        type="checkbox"
                        checked={isChecked}
                        onChange={() => { onSelect(partner.id), setIsChecked(!isChecked) }}
                        className="w-5 h-5 cursor-pointer accent-[var(--primary-color)] rounded-md border-gray-300"
                    />
                </label>

                {/* Partner Details */}
                <div className="flex flex-col">
                    <p className="font-semibold text-sm text-gray-900">Name: {partner.name} </p>
                    <p className="text-xs text-gray-600">{partner.location}</p>
                    <p className="text-xs text-gray-600">{partner.phone}</p>
                </div>
            </div>

            {/* Expanded Details */}
            {openCard && (
                <div className="mt-3 p-3 rounded-lg bg-[var(--bg-page-color)] transition-all">
                    <p className="text-xs text-gray-800">Additional details about {partner.name}...</p>
                </div>
            )}
        </div>
    );
};

export default DeliveryPartnerCard;
