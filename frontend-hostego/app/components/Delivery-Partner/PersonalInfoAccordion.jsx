import React, { useState } from 'react';
import { User, ChevronDown } from 'lucide-react';

const PersonalInfoAccordion = ({ deliveryPartner }) => {
    const [isOpen, setIsOpen] = useState(false);

    return (
        <div className="bg-white rounded-xl shadow-sm overflow-hidden">
            <button
                onClick={() => setIsOpen(!isOpen)}
                className="w-full px-4 py-3 flex items-center justify-between"
            >
                <div className="flex items-center gap-3">
                    <User className="w-5 h-5 text-[var(--primary-color)]" />
                    <span className="font-medium">Personal Information</span>
                </div>
                <ChevronDown className={`w-5 h-5 transition-transform ${isOpen ? 'rotate-180' : ''}`} />
            </button>

            {isOpen && (
                <div className="px-4 pb-4 space-y-4 border-t animate-fade-in p-2">
                    <InfoRow label="Name" value={deliveryPartner.user?.first_name} />
                    <InfoRow label="Phone" value={deliveryPartner.user?.mobile_number} />
                    <InfoRow label="Email" value={deliveryPartner.user?.email} />
                    <InfoRow label="Address" value={deliveryPartner.address} />
                </div>
            )}
        </div>
    );
};

const InfoRow = ({ label, value }) => (
    <div className="flex items-center justify-between">
        <span className="text-gray-500 text-sm">{label}</span>
        <span className="font-medium text-sm">{value}</span>
    </div>
);

export default PersonalInfoAccordion; 