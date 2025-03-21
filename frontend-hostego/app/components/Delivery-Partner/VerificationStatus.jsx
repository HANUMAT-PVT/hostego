import { Shield } from 'lucide-react';

const VerificationStatus = ({ deliveryPartner }) => {
    const status = deliveryPartner?.verification_status;
    
    return (
        <div className="bg-white rounded-xl p-6 shadow-sm">
            <div className="flex items-center gap-3 mb-4">
                <div className={`w-12 h-12 rounded-full flex items-center justify-center ${
                    status ? 'bg-green-50' : 'bg-yellow-50'
                }`}>
                    <Shield className={`w-6 h-6 ${
                        status ? 'text-green-500' : 'text-yellow-500'
                    }`} />
                </div>
                <div>
                    <h3 className="font-medium">Verification Status</h3>
                    <p className={`text-sm ${
                        status ? 'text-green-500' : 'text-yellow-500'
                    }`}>
                        {status ? 'Verified Partner' : 'Verification Pending'}
                    </p>
                </div>
            </div>

            {!status && (
                <div className="bg-yellow-50 border border-yellow-100 rounded-lg p-3">
                    <p className="text-sm text-yellow-700">
                        {deliveryPartner?.documents?.upi_id
                            ? "We have received your details. Sit back and relax while we verify your documents."
                            : "Please complete your verification to start accepting orders. Upload the required documents below."}
                    </p>
                </div>
            )}
        </div>
    );
};

export default VerificationStatus; 