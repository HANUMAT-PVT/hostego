"use client";

import React, { useState } from "react";
import BackNavigationButton from "../components/BackNavigationButton";
import {
    HelpCircle,
    Phone,
    Mail,
    MessageCircle,
    ChevronDown,
    Clock,
    Truck,
    CreditCard,
    ShieldCheck,
    UserCog,
} from "lucide-react";

const FAQItem = ({ question, answer }) => {
    const [isOpen, setIsOpen] = useState(false);

    return (
        <div className="border-b border-gray-100 last:border-none">
            <button
                onClick={() => setIsOpen(!isOpen)}
                className="flex items-center justify-between w-full py-4 text-left"
            >
                <span className="font-medium text-gray-800">{question}</span>
                <ChevronDown
                    size={20}
                    className={`text-gray-500 transition-transform ${isOpen ? "rotate-180" : ""
                        }`}
                />
            </button>
            {isOpen && (
                <div className="pb-4 text-gray-600 animate-fadeIn">
                    {answer}
                </div>
            )}
        </div>
    );
};

const ContactCard = ({ icon: Icon, title, description, action, actionText }) => (
    <div className="bg-white rounded-xl p-6 shadow-sm hover:shadow-md transition-all">
        <div className="flex items-start gap-4">
            <div className="w-12 h-12 rounded-lg bg-gradient-to-br from-[#655df0]/10 to-[#9333ea]/10 flex items-center justify-center">
                <Icon size={24} className="text-[#655df0]" />
            </div>
            <div className="flex-1">
                <h3 className="font-medium text-gray-800 mb-1">{title}</h3>
                <p className="text-sm text-gray-600 mb-3">{description}</p>
                <a
                    href={action}
                    className="inline-flex items-center text-[#655df0] hover:underline text-sm font-medium"
                >
                    {actionText}
                </a>
            </div>
        </div>
    </div>
);

const Help = () => {
    const faqs = [
        {
            question: "How do I place an order?",
            answer: "Simply browse through our app, select your items, add them to cart, and proceed to checkout. Make sure to provide your correct hostel room number for delivery."
        },
        {
            question: "What are your delivery hours?",
            answer: "We deliver during hostel-permitted hours. Delivery times may vary based on order volume and availability of delivery partners."
        },
        {
            question: "How do I track my order?",
            answer: "Once your order is confirmed, you can track its status in real-time through the 'Orders' section in the app."
        },
        {
            question: "What payment methods do you accept?",
            answer: "We accept digital payment methods. All transactions are secure and encrypted."
        },
        {
            question: "Can I cancel my order?",
            answer: "Orders can be cancelled before they are picked up by our delivery partner. Check the order details page for cancellation options."
        }
    ];

    return (
        <div className="min-h-screen bg-[#f4f6fb]">
            <BackNavigationButton title="Help & Support" />

            <div className="p-4 space-y-6">
                {/* Quick Contact Cards */}
                <div className="grid gap-4">
                    <ContactCard
                        icon={Phone}
                        title="24/7 Support"
                        description="Get immediate assistance from our support team"
                        action="tel:+918264121428"
                        actionText="Call Support"
                    />
                    <ContactCard
                        icon={Mail}
                        title="Email Support"
                        description="Send us your queries and feedback"
                        action="mailto:support@hostego.in"
                        actionText="Email Us"
                    />
                    <ContactCard
                        icon={MessageCircle}
                        title="Live Chat"
                        description="Chat with our support team in real-time"
                        action="#"
                        actionText="Start Chat"
                    />
                </div>

                {/* FAQ Section */}
                <div className="bg-white rounded-xl p-6 shadow-sm">
                    <div className="flex items-center gap-2 mb-6">
                        <HelpCircle className="text-[#655df0]" size={24} />
                        <h2 className="text-xl font-semibold">Frequently Asked Questions</h2>
                    </div>
                    <div className="space-y-2">
                        {faqs.map((faq, index) => (
                            <FAQItem key={index} {...faq} />
                        ))}
                    </div>
                </div>

                {/* Help Categories */}
                <div className="grid grid-cols-2 gap-4">
                    {[
                        { icon: Clock, title: "Order Timing" },
                        { icon: Truck, title: "Delivery Info" },
                        { icon: CreditCard, title: "Payment Help" },
                        { icon: ShieldCheck, title: "Safety & Security" },
                        { icon: UserCog, title: "Account Settings" },
                        { icon: MessageCircle, title: "Feedback" },
                    ].map((category, index) => (
                        <button
                            key={index}
                            className="bg-white p-4 rounded-xl shadow-sm hover:shadow-md transition-all text-center"
                        >
                            <category.icon
                                size={24}
                                className="text-[#655df0] mx-auto mb-2"
                            />
                            <span className="text-sm font-medium text-gray-800">
                                {category.title}
                            </span>
                        </button>
                    ))}
                </div>

                {/* Support Hours Notice */}
                <div className="bg-gradient-to-r from-[#655df0] to-[#9333ea] p-6 rounded-xl text-white">
                    <h3 className="font-semibold mb-2">Support Hours</h3>
                    <p className="opacity-90">
                        Our support team is available 24/7 to assist you with any queries or concerns.
                        We typically respond within 30 minutes during peak hours.
                    </p>
                </div>
            </div>
        </div>
    );
};

export default Help;