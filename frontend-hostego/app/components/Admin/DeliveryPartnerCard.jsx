import { CheckCircle, User } from 'lucide-react'

const DeliveryPartnerCard = ({ partner, isSelected, onSelect }) => (
    <div
        onClick={() => onSelect(partner)}
        className={`p-4 rounded-xl border-2 cursor-pointer transition-all duration-200 
            ${isSelected
                ? 'border-[var(--primary-color)] bg-[var(--primary-color)]/5'
                : 'border-gray-100 bg-white hover:border-[var(--primary-color)]/30'}`}
    >
        <div className="flex items-center gap-3">
            <div className="w-12 h-12 rounded-full bg-gray-100 flex items-center justify-center">
                <User className={`w-6 h-6 ${isSelected ? 'text-[var(--primary-color)]' : 'text-gray-500'}`} />
            </div>
            <div>
                <p className="font-medium">{partner?.user?.first_name} {partner?.user?.last_name}</p>
                <p className="text-sm text-gray-600">{partner?.user?.mobile_number}</p>
                <p className="text-sm text-gray-600">{partner?.address}</p>
            </div>
            {isSelected && (
                <div className="ml-auto">
                    <CheckCircle className="w-5 h-5 text-[var(--primary-color)]" />
                </div>
            )}
        </div>
    </div>
)

export default DeliveryPartnerCard