"use client";

import React from "react";
import BackNavigationButton from "../components/BackNavigationButton";
import { Shield, Lock, UserCheck, Bell, Database, Phone } from "lucide-react";

const Section = ({ icon: Icon, title, children }) => (
    <div className="mb-8">
        <div className="flex items-start gap-3 mb-4">
            <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-[#655df0]/10 to-[#9333ea]/10 flex items-center justify-center flex-shrink-0">
                <Icon size={20} className="text-[#655df0]" />
            </div>
            <h2 className="text-xl font-semibold text-gray-800">{title}</h2>
        </div>
        <div className="space-y-3 text-gray-600 ml-11">{children}</div>
    </div>
);

const PrivacyPolicy = () => {
    return (
        <div className="min-h-screen bg-[#f4f6fb]">
            <BackNavigationButton title="Privacy Policy" />

            <div className="p-4">
                <div className="bg-white rounded-xl p-6 shadow-sm">
                    {/* Introduction */}
                    <div className="mb-8">
                        <div className="flex items-center gap-2 mb-4">
                            <Shield className="text-[#655df0]" size={24} />
                            <h1 className="text-2xl font-bold">Privacy Policy</h1>
                        </div>
                        <p className="text-gray-600">
                            Last updated: {new Date().toLocaleDateString()}
                        </p>
                        <p className="mt-4 text-gray-600">
                            At Hostego, we take your privacy seriously. This policy describes how we collect,
                            use, and protect your personal information when you use our hostel delivery and services platform.
                        </p>
                    </div>

                    {/* Information We Collect */}
                    <Section icon={UserCheck} title="Information We Collect">
                        <p>We collect the following types of information:</p>
                        <ul className="list-disc pl-5 space-y-2">
                            <li>Name and contact information</li>
                            <li>Hostel room number and address</li>
                            <li>Order history and preferences</li>
                            <li>Payment information</li>
                            <li>Device information and location data (for delivery purposes)</li>
                            <li>Communication records with our support team</li>
                        </ul>
                    </Section>

                    {/* How We Use Your Information */}
                    <Section icon={Database} title="How We Use Your Information">
                        <p>Your information helps us to:</p>
                        <ul className="list-disc pl-5 space-y-2">
                            <li>Process and deliver your orders</li>
                            <li>Verify your identity and hostel residency</li>
                            <li>Send important updates about your orders</li>
                            <li>Improve our services and user experience</li>
                            <li>Ensure efficient delivery to your hostel room</li>
                            <li>Prevent fraud and maintain security</li>
                        </ul>
                    </Section>

                    {/* Data Security */}
                    <Section icon={Lock} title="Data Security">
                        <p>
                            We implement strong security measures to protect your data:
                        </p>
                        <ul className="list-disc pl-5 space-y-2">
                            <li>Encrypted data transmission</li>
                            <li>Secure payment processing</li>
                            <li>Regular security audits</li>
                            <li>Limited access to personal information</li>
                            <li>Secure data storage systems</li>
                        </ul>
                    </Section>

                    {/* Information Sharing */}
                    <Section icon={Bell} title="Information Sharing">
                        <p>We may share your information with:</p>
                        <ul className="list-disc pl-5 space-y-2">
                            <li>Delivery partners (only necessary details for delivery)</li>
                            <li>Payment processors for transaction handling</li>
                            <li>Service providers who assist our operations</li>
                            <li>Law enforcement when required by law</li>
                        </ul>
                        <p className="mt-3">
                            We never sell your personal information to third parties.
                        </p>
                    </Section>

                    {/* Your Rights */}
                    <Section icon={UserCheck} title="Your Rights">
                        <p>You have the right to:</p>
                        <ul className="list-disc pl-5 space-y-2">
                            <li>Access your personal information</li>
                            <li>Correct inaccurate information</li>
                            <li>Request deletion of your data</li>
                            <li>Opt-out of marketing communications</li>
                            <li>Export your data</li>
                        </ul>
                    </Section>

                    {/* Contact Us */}
                    <Section icon={Phone} title="Contact Us">
                        <p>
                            If you have questions about this privacy policy or your data, please contact us at:{" "}
                            <a
                                href="mailto:hanumat@hostego.in"
                                className="text-[#655df0] hover:underline"
                            >
                                hanumat@hostego.in
                            </a>
                        </p>
                    </Section>

                    {/* Footer */}
                    <div className="mt-12 pt-6 border-t border-gray-100">
                        <p className="text-sm text-gray-500 text-center">
                            © {new Date().getFullYear()} Hostego. All rights reserved.
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default PrivacyPolicy;
