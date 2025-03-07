"use client";

import React from "react";
import { Package, Clock, Printer, CheckCircle, Building, Coffee, Heart, MapPin, } from "lucide-react";
import BackNavigationButton from "../components/BackNavigationButton";
import Head from "next/head";

const FeatureCard = ({ icon: Icon, title, description }) => (
    <div className="bg-white rounded-xl p-6 shadow-sm hover:shadow-md transition-all">
        <div className="w-12 h-12 rounded-lg bg-gradient-to-br from-[#655df0] to-[#9333ea] flex items-center justify-center mb-4">
            <Icon size={24} color="white" />
        </div>
        <h3 className="text-lg font-semibold mb-2">{title}</h3>
        <p className="text-gray-600">{description}</p>
    </div>
);

const AboutPage = () => {
    return (
        <div className="min-h-screen bg-[#f4f6fb]">
            <Head>
                <title>About Hostego - Simplify Your Hostel Life</title>
                <meta name="description" content="Hostego is your one-stop solution for hostel life. Get late-night deliveries, emergency printouts, task assistance, and future hostel solutions." />
                <meta name="keywords" content="Hostel life, delivery service, student services, late-night food, printouts, hostel management" />
                <meta name="author" content="Hostego" />

                {/* Open Graph / Facebook */}
                <meta property="og:type" content="website" />
                <meta property="og:title" content="About Hostego - Simplify Your Hostel Life" />
                <meta property="og:description" content="Hostego is revolutionizing hostel life with instant delivery, printouts, and more." />
                <meta property="og:image" content="/images/hostego-banner.jpg" />
                <meta property="og:url" content="https://www.hostego.in/about" />

                {/* Twitter Meta Tags */}
                <meta name="twitter:card" content="summary_large_image" />
                <meta name="twitter:title" content="About Hostego - Simplify Your Hostel Life" />
                <meta name="twitter:description" content="Hostego makes hostel life easier with 24/7 delivery, emergency printouts, and more." />
                <meta name="twitter:image" content="/images/hostego-banner.jpg" />

                {/* Canonical Link for SEO */}
                <link rel="canonical" href="https://www.hostego.in/about" />
            </Head>
            <BackNavigationButton title="About Hostego" />

            {/* Hero Section */}
            <div className="px-4 py-8 bg-gradient-to-br from-[#655df0] to-[#9333ea] text-white">
                <h1 className="text-4xl font-bold mb-4">
                    Simplify Your
                    <span className="block">Hostel Life</span>
                </h1>
                <p className="text-lg opacity-90 mb-6">
                    One stop solution for all your hostel needs
                </p>
                <div className="flex items-center gap-2 text-sm">
                    <CheckCircle size={16} />
                    <span>24/7 Delivery</span>
                </div>
            </div>

            {/* Main Features */}
            <div className="px-4 py-8">
                <h2 className="text-2xl font-semibold mb-6">What We Offer</h2>
                <div className="grid gap-4">
                    <FeatureCard
                        icon={Package}
                        title="Late Night Deliveries"
                        description="Get your midnight cravings satisfied with our 24/7 delivery service. From Maggi to snacks, we've got you covered."
                    />
                    <FeatureCard
                        icon={Printer}
                        title="Emergency Printouts"
                        description="Need urgent printouts? We'll deliver them right to your doorstep, saving your time and effort."
                    />
                    <FeatureCard
                        icon={Coffee}
                        title="Task Assistance"
                        description="Any task you can't handle? Our team is here to help. From errands to special requests, we make it happen."
                    />
                    <FeatureCard
                        icon={Building}
                        title="Future: Hostel Solutions"
                        description="Soon: Find the perfect hostel across cities, better food options, and complete hostel management solutions."
                    />
                </div>
            </div>

            {/* How It Works */}
            <div className="px-4 py-8 bg-white">
                <h2 className="text-2xl font-semibold mb-6">How It Works</h2>
                <div className="space-y-6">
                    <div className="flex items-start gap-4">
                        <div className="w-8 h-8 rounded-full bg-[#655df0] text-white flex items-center justify-center flex-shrink-0">
                            1
                        </div>
                        <div>
                            <h3 className="font-medium mb-1">Open App</h3>
                            <p className="text-gray-600">Launch Hostego and browse available services</p>
                        </div>
                    </div>
                    <div className="flex items-start gap-4">
                        <div className="w-8 h-8 rounded-full bg-[#655df0] text-white flex items-center justify-center flex-shrink-0">
                            2
                        </div>
                        <div>
                            <h3 className="font-medium mb-1">Place Order</h3>
                            <p className="text-gray-600">Select what you need and confirm your order</p>
                        </div>
                    </div>
                    <div className="flex items-start gap-4">
                        <div className="w-8 h-8 rounded-full bg-[#655df0] text-white flex items-center justify-center flex-shrink-0">
                            3
                        </div>
        <div>
                            <h3 className="font-medium mb-1">Room Delivery</h3>
                            <p className="text-gray-600">Relax while we deliver right to your room</p>
                        </div>
                    </div>
                </div>
            </div>

            {/* Vision */}
            <div className="px-4 py-8">
                <div className="bg-gradient-to-br from-[#655df0] to-[#9333ea] p-6 rounded-xl text-white">
                    <h2 className="text-2xl font-semibold mb-4">Our Vision</h2>
                    <p className="opacity-90 leading-relaxed">
                        Hostego isn't just a delivery service - we're building the future of hostel living.
                        Our mission is to create a comprehensive platform that makes hostel life easier,
                        more comfortable, and more enjoyable for students worldwide.
                    </p>
                </div>
            </div>

            {/* Support & Made with Love Section */}
            <div className="relative bg-white mt-12">
                {/* Top Wave Decoration */}
                <div className="absolute top-0 left-0 right-0 transform -translate-y-8 overflow-hidden">
                    <svg className="w-full h-8" viewBox="0 0 1440 54" preserveAspectRatio="none">
                        <path
                            d="M0 27L48 25C96 23 192 19 288 20.3C384 21.7 480 28.3 576 29.3C672 30.3 768 25.7 864 22.8C960 20 1056 19 1152 22.8C1248 26.7 1344 35.3 1392 39.7L1440 44V54H1392C1344 54 1248 54 1152 54C1056 54 960 54 864 54C768 54 672 54 576 54C480 54 384 54 288 54C192 54 96 54 48 54H0V27Z"
                            fill="url(#gradient)"
                            className="opacity-10"
                        />
                        <defs>
                            <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="0%">
                                <stop offset="0%" style={{ stopColor: '#655df0' }} />
                                <stop offset="100%" style={{ stopColor: '#9333ea' }} />
                            </linearGradient>
                        </defs>
                    </svg>
                </div>

                <div className="container mx-auto px-6 py-16">
                    <div className="text-center space-y-12">
                        {/* Location Badge */}
                        <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-gradient-to-r from-[#655df0]/10 to-[#9333ea]/10">
                            <MapPin size={18} className="text-[#655df0]" />
                            <span className="font-medium bg-gradient-to-r from-[#655df0] to-[#9333ea] text-transparent bg-clip-text">
                                Made with love in Chandigarh
                            </span>
                        </div>

                        {/* Love Message */}
                        <div className="relative">
                            <div className="flex items-center justify-center gap-3 text-xl">
                                <span className="font-medium text-gray-800">Crafted with</span>
                                <Heart
                                    size={28}
                                    className="text-red-500 fill-red-500 animate-pulse"
                                />
                                <span className="font-medium text-gray-800">for students</span>
                            </div>
                        </div>

                        {/* Collaboration Section */}
                        <div className="max-w-md mx-auto space-y-8">
                            <h3 className="text-lg font-medium text-gray-800">
                                Support or collaborate with us
                            </h3>

                            {/* Social Links */}
                            <div className="flex justify-center gap-6">
                                {[1, 2, 3].map((_, idx) => (
                                    <a
                                        key={idx}
                                        href="#"
                                        className="group relative w-12 h-12 rounded-full bg-gradient-to-br from-[#655df0] to-[#9333ea] p-[2px] transition-all duration-300 hover:scale-110 hover:shadow-lg hover:shadow-[#655df0]/20"
                                    >
                                        <div className="absolute inset-[2px] bg-white rounded-full transition-colors group-hover:bg-gradient-to-br group-hover:from-[#655df0]/10 group-hover:to-[#9333ea]/10"></div>
                                        <div className="relative w-full h-full flex items-center justify-center text-[#655df0] group-hover:text-white transition-colors">
                                            {/* Icon will go here */}
                                        </div>
                                    </a>
                                ))}
                            </div>

                            {/* Contact Info */}
                            <div className="text-center">
                                <p className="text-gray-600">
                                    Have suggestions? Email us at{" "}
                                    <a
                                        href="mailto:hanumat@hostego.in"
                                        className="relative inline-block font-medium text-[#655df0] hover:text-[#9333ea] transition-colors"
                                    >
                                        hanumat@hostego.in
                                        <span className="absolute bottom-0 left-0 w-full h-0.5 bg-gradient-to-r from-[#655df0] to-[#9333ea] transform scale-x-0 transition-transform origin-left hover:scale-x-100"></span>
                                    </a>
                                </p>
                            </div>
                        </div>

                        {/* Footer */}
                        <div className="pt-10 border-t border-gray-100">
                            <div className="flex flex-col items-center gap-4">
                                <div className="text-2xl font-bold bg-gradient-to-r from-[#655df0] to-[#9333ea] text-transparent bg-clip-text">
                                    Hostego
                                </div>
                                <p className="text-sm text-gray-500">
                                    © {new Date().getFullYear()} Hostego. All rights reserved.
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default AboutPage;