"use client";

import React from "react";
import BackNavigationButton from "../components/BackNavigationButton";
import { Shield, Clock, AlertCircle } from "lucide-react";

const Section = ({ title, children }) => (
    <div className="mb-8">
        <h2 className="text-xl font-semibold mb-3 text-[#655df0]">{title}</h2>
        <div className="space-y-3 text-gray-600">{children}</div>
    </div>
);

const TermsConditionsPage = () => {
    return (
        <div className="min-h-screen bg-[#f4f6fb]">
            <BackNavigationButton title="Terms & Conditions" />

            <div className="p-4">
                <div className="bg-white rounded-xl p-6 shadow-sm">
                    {/* Introduction */}
                    <div className="mb-8">
                        <div className="flex items-center gap-2 mb-4">
                            <Shield className="text-[#655df0]" size={24} />
                            <h1 className="text-2xl font-bold">Hostego Terms of Service</h1>
                        </div>
                        <p className="text-gray-600">
                            Last updated: {new Date().toLocaleDateString()}
                        </p>
                    </div>

                    {/* Acceptance of Terms */}
                    <Section title="1. Acceptance of Terms">
                        <p>
                            By accessing or using Hostego's services, you agree to be bound by these Terms and Conditions.
                            These terms apply to all users, including customers, delivery partners, and service providers.
                        </p>
                    </Section>

                    {/* Service Description */}
                    <Section title="2. Service Description">
                        <p>
                            Hostego provides various services including but not limited to:
                        </p>
                        <ul className="list-disc pl-5 space-y-2">
                            <li>Delivery services within hostel premises</li>
                            <li>Emergency printout delivery services</li>
                            <li>Task assistance and errand services</li>
                            <li>Future hostel management and booking services</li>
                        </ul>
                    </Section>

                    {/* User Responsibilities */}
                    <Section title="3. User Responsibilities">
                        <p>Users of Hostego services must:</p>
                        <ul className="list-disc pl-5 space-y-2">
                            <li>Provide accurate and complete information when using our services</li>
                            <li>Maintain the security of their account credentials</li>
                            <li>Comply with hostel rules and regulations</li>
                            <li>Use the service in a lawful and responsible manner</li>
                            <li>Not misuse or attempt to manipulate our platform</li>
                        </ul>
                    </Section>

                    {/* Delivery Terms */}
                    <Section title="4. Delivery Terms">
                        <div className="flex items-start gap-2">
                            <Clock size={20} className="text-[#655df0] mt-1" />
                            <p>
                                Our delivery services  subject to availability of delivery partners and hostel regulations.
                                Delivery times may vary based on order volume and operational conditions.
                            </p>
                        </div>
                        <p className="mt-2">
                            Users must be present at the specified delivery location (hostel room) to receive their orders.
                        </p>
                    </Section>

                    {/* Payment Terms */}
                    <Section title="5. Payment and Pricing">
                        <ul className="list-disc pl-5 space-y-2">
                            <li>All prices are in Indian Rupees (INR)</li>
                            <li>Payment must be made at the time of ordering</li>
                            <li>We accept various payment methods as displayed in the app</li>
                            <li>Delivery fees and service charges will be clearly indicated before order confirmation</li>
                        </ul>
                    </Section>

                    {/* Cancellation Policy */}
                    <Section title="6. Cancellation and Refunds">
                        <ul className="list-disc pl-5 space-y-2">
                            <li>Orders can be cancelled before the delivery partner picks up the order</li>
                            <li>Refunds will be processed according to our refund policy</li>
                            <li>Refunds will be processed according to our refund policy and will be added to the user's wallet.</li>
                            <li>Repeated cancellations may result in service restrictions</li>
                        </ul>
                    </Section>

                    {/* Privacy and Data */}
                    <Section title="7. Privacy and Data Protection">
                        <p>
                            We collect and process personal data in accordance with our Privacy Policy.
                            This includes necessary information for delivery services and account management.
                        </p>
                    </Section>

                    {/* Liability */}
                    <Section title="8. Limitation of Liability">
                        <div className="flex items-start gap-2">
                            <AlertCircle size={20} className="text-[#655df0] mt-1" />
                            <p>
                                Hostego's liability is limited to the value of the service provided.
                                We are not responsible for indirect losses or damages arising from the use of our services.
                            </p>
                        </div>
                    </Section>

                    {/* Changes to Terms */}
                    <Section title="9. Modifications to Terms">
                        <p>
                            Hostego reserves the right to modify these terms at any time.
                            Users will be notified of significant changes, and continued use of the service
                            constitutes acceptance of modified terms.
                        </p>
                    </Section>

                    {/* Contact Information */}
                    <Section title="10. Contact Us">
                        <p>
                            For questions about these terms, please contact us at:{" "}
                            <a href="mailto:hanumat@hostego.in" className="text-[#655df0] hover:underline">
                                hanumat@hostego.in
                            </a>
                        </p>
                    </Section>
                </div>
            </div>
        </div>
    );
};

export default TermsConditionsPage;