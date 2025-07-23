"use client";
import { useRouter } from "next/navigation";
import { ArrowRight, Truck, Sparkles, Users, Gift, PartyPopper } from "lucide-react";
import Link from "next/link";

export default function Home() {
  const router = useRouter();

  return (
    <div className="relative flex flex-col items-center justify-center min-h-screen bg-gradient-to-br from-purple-900 via-pink-800 to-orange-700 text-white p-4 sm:p-8 overflow-hidden">
      {/* Animated Background Elements */}
      <div className="absolute inset-0 overflow-hidden">
        {/* Floating confetti-like elements */}
        <div className="absolute top-20 left-10 w-4 h-4 bg-yellow-300 rounded-full animate-bounce" style={{ animationDelay: '0s' }}></div>
        <div className="absolute top-40 right-20 w-3 h-3 bg-pink-300 rounded-full animate-bounce" style={{ animationDelay: '0.5s' }}></div>
        <div className="absolute bottom-40 left-20 w-2 h-2 bg-blue-300 rounded-full animate-bounce" style={{ animationDelay: '1s' }}></div>
        <div className="absolute bottom-20 right-10 w-3 h-3 bg-green-300 rounded-full animate-bounce" style={{ animationDelay: '1.5s' }}></div>
        <div className="absolute top-60 left-1/2 w-2 h-2 bg-purple-300 rounded-full animate-bounce" style={{ animationDelay: '2s' }}></div>
        
        {/* Gradient orbs */}
        <div className="absolute top-0 left-0 w-96 h-96 bg-gradient-to-r from-pink-500/20 to-purple-500/20 blur-3xl rounded-full -translate-x-1/2 -translate-y-1/2 animate-pulse"></div>
        <div className="absolute bottom-0 right-0 w-96 h-96 bg-gradient-to-r from-orange-500/20 to-yellow-500/20 blur-3xl rounded-full translate-x-1/2 translate-y-1/2 animate-pulse" style={{ animationDelay: '1s' }}></div>
      </div>

      {/* Main Content Container */}
      <main className="relative max-w-4xl mx-auto flex flex-col items-center text-center gap-8 p-4 z-10">
        {/* Freshers Celebration Badge */}
        <div className="bg-gradient-to-r from-pink-500/20 to-purple-500/20 border border-pink-500/30 px-4 py-2 rounded-full backdrop-blur-sm">
          <p className="text-sm font-bold text-pink-200 flex items-center gap-2">
            <PartyPopper className="w-4 h-4 text-pink-300" />
            <span className="w-2 h-2 bg-pink-400 rounded-full animate-pulse"></span>
            🎉 FRESHERS CELEBRATION 🎉
          </p>
        </div>

        {/* Main Heading with Sparkles */}
        <div className="space-y-3">
          <div className="flex items-center justify-center gap-3">
            <Sparkles className="w-8 h-8 text-yellow-300 animate-pulse" />
            <h1 className="text-5xl sm:text-6xl font-bold tracking-tight bg-gradient-to-r from-yellow-300 via-pink-300 to-purple-300 bg-clip-text text-transparent">
              Hostego
            </h1>
            <Sparkles className="w-8 h-8 text-yellow-300 animate-pulse" />
          </div>
          <p className="text-lg sm:text-xl font-medium text-white/90">
            Welcome to Your Hostel Life! 🏠✨
          </p>
          <p className="text-sm sm:text-base text-white/70">
            Your Campus Delivery Partner
          </p>
        </div>

        {/* Special Freshers Offer Card */}
        <div className="relative bg-gradient-to-br from-yellow-400/10 to-orange-500/10 border-2 border-yellow-400/30 p-8 rounded-2xl w-full max-w-lg backdrop-blur-sm shadow-2xl">
          <div className="absolute -top-4 left-1/2 -translate-x-1/2">
            <div className="bg-gradient-to-r from-pink-500 to-purple-600 text-white text-sm font-bold px-6 py-2 rounded-full shadow-lg">
              🎓 FRESHERS SPECIAL 🎓
            </div>
          </div>

          <div className="flex flex-col items-center gap-6 mt-4">
            <div className="flex items-center gap-3">
              <Gift className="w-8 h-8 text-yellow-300 animate-bounce" />
              <h2 className="text-3xl font-bold text-yellow-300">FREE DELIVERY</h2>
              <Gift className="w-8 h-8 text-yellow-300 animate-bounce" />
            </div>
            
            <div className="space-y-3">
              <div className="flex items-center justify-center gap-2">
                <Users className="w-5 h-5 text-pink-300" />
                <p className="text-xl text-white/90">
                  For Next{" "}
                  <span className="text-yellow-300 font-bold text-3xl bg-gradient-to-r from-yellow-300 to-orange-300 bg-clip-text text-transparent">
                    101
                  </span>{" "}
                  Freshers
                </p>
              </div>
              <p className="text-sm text-yellow-200 font-medium">
                🚀 Limited Time Offer • Join the Celebration! 🚀
              </p>
              <p className="text-xs text-white/60">
                Order your favorite food and get it delivered FREE to your hostel
              </p>
            </div>
          </div>
        </div>

        {/* CTA Button */}
        <button
          onClick={() => router.push("/home")}
          className="group px-8 py-4 text-lg font-bold text-black bg-gradient-to-r from-yellow-400 to-orange-400 
                   rounded-xl shadow-2xl hover:shadow-3xl transition-all duration-300 transform hover:scale-105
                   flex items-center gap-3 hover:from-yellow-300 hover:to-orange-300"
        >
          <PartyPopper className="w-5 h-5" />
          Start Ordering Now!
          <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
        </button>

        {/* Features Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mt-8 w-full max-w-3xl">
          <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-xl p-4 hover:bg-white/10 transition-all">
            <Truck className="w-6 h-6 text-green-400 mb-2 mx-auto" />
            <p className="text-lg font-bold text-white">Fast Delivery</p>
            <p className="text-sm text-white/60">To your hostel doorstep</p>
          </div>
          <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-xl p-4 hover:bg-white/10 transition-all">
            <Gift className="w-6 h-6 text-pink-400 mb-2 mx-auto" />
            <p className="text-lg font-bold text-white">Free Delivery</p>
            <p className="text-sm text-white/60">For first 101 orders</p>
          </div>
          <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-xl p-4 hover:bg-white/10 transition-all">
            <Sparkles className="w-6 h-6 text-yellow-400 mb-2 mx-auto" />
            <p className="text-lg font-bold text-white">Best Prices</p>
            <p className="text-sm text-white/60">Student-friendly rates</p>
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="relative z-10 mt-16 text-sm text-white/70 flex flex-col sm:flex-row gap-3 sm:gap-6 justify-center items-center py-6 border-t border-white/10 w-full max-w-4xl mx-auto">
        <Link href="/privacy-policy" className="hover:text-white transition">
          Privacy Policy
        </Link>
        <Link href="/terms-conditions" className="hover:text-white transition">
          Terms & Conditions
        </Link>
        <Link
          href="/refund-and-cancellation"
          className="hover:text-white transition"
        >
          Refund & Cancellation
        </Link>
        <Link href="/support" className="hover:text-white transition">
          Support
        </Link>
        <Link href="/ship-and-delivery" className="hover:text-white transition">
          Shipping & Delivery
        </Link>
      </footer>
    </div>
  );
}
