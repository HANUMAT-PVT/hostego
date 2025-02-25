"use client";

import React from "react";
import PhoneEmailAuthButton from "../../components/PhoneEmailAuth";
import { Clock, Package, Shield, Smartphone, Home } from "lucide-react";

const Page = () => {
  return (
    <div className="relative min-h-screen bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center">
      {/* Floating Circles Background */}
      <div className="absolute top-10 left-10 w-20 h-20 bg-white/10 rounded-full blur-xl" />
      <div className="absolute bottom-40 right-10 w-32 h-32 bg-purple-400/10 rounded-full blur-xl" />

      <div className="flex flex-col absolute justify-between bg-white h-[80vh] bottom-0 w-full max-w-[400px] px-6 py-8 rounded-t-[2.5rem] shadow-2xl">
        {/* Heading and Tagline */}
        <div className="text-center space-y-3">
          <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center mx-auto mb-2">
            <Home className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-2xl font-bold bg-gradient-to-r from-blue-500 to-purple-600 text-transparent bg-clip-text">
            Welcome to Hostego
          </h1>
          <p className="text-gray-600 text-sm font-medium">Simplify your hostel life</p>
        </div>

        {/* Features Grid */}
        <div className="grid grid-cols-2 gap-4 mb-8">
          <div className="p-4 rounded-2xl bg-white border border-gray-100 flex flex-col items-center justify-center">
            <div className="w-10 h-10 rounded-full bg-blue-50 flex items-center justify-center mx-auto mb-3">
              <Package className="w-5 h-5 text-[#3b82f6]" />
            </div>
            <h3 className="font-medium text-sm">Fast Delivery</h3>
          </div>
          <div className="p-4 rounded-2xl bg-white border border-gray-100 flex flex-col items-center justify-center">
            <div className="w-10 h-10 rounded-full bg-purple-50 flex items-center justify-center mx-auto mb-3">
              <Shield className="w-5 h-5 text-[#9333ea]" />
            </div>
            <h3 className="font-medium text-sm">Secure Orders</h3>
          </div>
          <div className="p-4 rounded-2xl bg-white border border-gray-100 flex flex-col items-center justify-center">
            <div className="w-10 h-10 rounded-full bg-blue-50 flex items-center justify-center mx-auto mb-3">
              <Clock className="w-5 h-5 text-[#3b82f6]" />
            </div>
            <h3 className="font-medium text-sm text-nowrap">
            Order Tracking
            </h3>
          </div>
          <div className="p-4 rounded-2xl bg-white border border-gray-100 flex flex-col items-center justify-center">
            <div className="w-10 h-10 rounded-full bg-purple-50 flex items-center justify-center mx-auto mb-3">
              <Smartphone className="w-5 h-5 text-[#9333ea]" />
            </div>
            <h3 className="font-medium text-sm">Easy Login</h3>
          </div>
        </div>

        {/* Auth Section */}
        <div className="space-y-4">
          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-gray-200"></div>
            </div>
            <div className="relative flex justify-center text-sm">
              <span className="px-4 bg-white text-gray-500">Continue with</span>
            </div>
          </div>

          <PhoneEmailAuthButton />

          <p className="text-gray-500 text-center text-xs">
            By continuing, you agree to our{" "}
            <a className="text-blue-600 font-medium hover:underline" href="#">
              Terms of Service
            </a>{" "}
            and{" "}
            <a className="text-blue-600 font-medium hover:underline" href="#">
              Privacy Policy
            </a>
          </p>
        </div>
      </div>
    </div>
  );
};

export default Page;
