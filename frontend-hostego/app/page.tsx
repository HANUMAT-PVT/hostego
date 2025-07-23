"use client";
import { useRouter } from "next/navigation";
import { ArrowRight, Truck, Sparkles, Users, Gift, PartyPopper } from "lucide-react";
import Link from "next/link";

export default function Home() {
  const router = useRouter();

  return (
    <div className="relative flex flex-col items-center justify-center min-h-screen bg-gradient-to-br from-purple-900 via-purple-800 to-violet-900 text-white p-4 sm:p-8 overflow-hidden">
      {/* Animated Background Elements */}
      <div className="absolute inset-0 overflow-hidden">
        {/* Floating confetti-like elements */}
        <div className="absolute top-20 left-10 w-4 h-4 bg-yellow-400 rounded-full animate-bounce" style={{ animationDelay: '0s' }}></div>
        <div className="absolute top-40 right-20 w-3 h-3 bg-white rounded-full animate-bounce" style={{ animationDelay: '0.5s' }}></div>
        <div className="absolute bottom-40 left-20 w-2 h-2 bg-yellow-300 rounded-full animate-bounce" style={{ animationDelay: '1s' }}></div>
        <div className="absolute bottom-20 right-10 w-3 h-3 bg-white rounded-full animate-bounce" style={{ animationDelay: '1.5s' }}></div>
        <div className="absolute top-60 left-1/2 w-2 h-2 bg-yellow-400 rounded-full animate-bounce" style={{ animationDelay: '2s' }}></div>
        
        {/* Gradient orbs */}
        <div className="absolute top-0 left-0 w-96 h-96 bg-gradient-to-r from-purple-500/20 to-violet-500/20 blur-3xl rounded-full -translate-x-1/2 -translate-y-1/2 animate-pulse"></div>
        <div className="absolute bottom-0 right-0 w-96 h-96 bg-gradient-to-r from-yellow-500/10 to-white/10 blur-3xl rounded-full translate-x-1/2 translate-y-1/2 animate-pulse" style={{ animationDelay: '1s' }}></div>
      </div>

      {/* Main Content Container */}
      <main className="relative max-w-4xl mx-auto flex flex-col items-center text-center gap-8 p-4 z-10">
        {/* Freshers Celebration Badge */}
        <div className="bg-gradient-to-r from-yellow-400/20 to-yellow-500/20 border-2 border-yellow-400/40 px-6 py-3 rounded-full backdrop-blur-sm shadow-lg">
          <p className="text-sm font-bold text-yellow-200 flex items-center gap-2">
            <PartyPopper className="w-4 h-4 text-yellow-300" />
            <span className="w-2 h-2 bg-yellow-400 rounded-full animate-pulse"></span>
            🎉 FRESHERS CELEBRATION 🎉
          </p>
        </div>

        {/* Main Heading with Sparkles */}
        <div className="space-y-4">
          <div className="flex items-center justify-center gap-3">
            <Sparkles className="w-8 h-8 text-yellow-300 animate-pulse" />
            <h1 className="text-6xl sm:text-7xl font-bold tracking-tight bg-gradient-to-r from-white via-yellow-200 to-yellow-300 bg-clip-text text-transparent">
              Hostego
            </h1>
            <Sparkles className="w-8 h-8 text-yellow-300 animate-pulse" />
          </div>
          <p className="text-xl sm:text-2xl font-medium text-white">
            Welcome to Your Hostel Life! 🏠✨
          </p>
          <p className="text-base sm:text-lg text-white/80">
            Your Campus Delivery Partner
          </p>
        </div>

        {/* Special Freshers Offer Card */}
        <div className="relative bg-gradient-to-br from-white/10 to-white/5 border-2 border-yellow-400/50 p-8 rounded-3xl w-full max-w-lg backdrop-blur-sm shadow-2xl">
          <div className="absolute -top-4 left-1/2 -translate-x-1/2">
            <div className="bg-gradient-to-r from-yellow-400 to-yellow-500 text-purple-900 text-sm font-bold px-6 py-2 rounded-full shadow-lg">
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
                <Users className="w-5 h-5 text-white" />
                <p className="text-xl text-white">
                  For Next{" "}
                  <span className="text-yellow-300 font-bold text-4xl">
                    101
                  </span>{" "}
                  Freshers
                </p>
              </div>
              <p className="text-sm text-yellow-200 font-medium">
                🚀 Limited Time Offer • Join the Celebration! 🚀
              </p>
              <p className="text-xs text-white/70">
                Order your favorite food and get it delivered FREE to your hostel
              </p>
            </div>
          </div>
        </div>

        {/* CTA Button */}
        <button
          onClick={() => router.push("/home")}
          className="group px-10 py-5 text-xl font-bold text-purple-900 bg-gradient-to-r from-yellow-400 to-yellow-500 
                   rounded-2xl shadow-2xl hover:shadow-3xl transition-all duration-300 transform hover:scale-105
                   flex items-center gap-3 hover:from-yellow-300 hover:to-yellow-400"
        >
          <PartyPopper className="w-6 h-6" />
          Start Ordering Now!
          <ArrowRight className="w-6 h-6 group-hover:translate-x-1 transition-transform" />
        </button>

        {/* Features Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-6 mt-10 w-full max-w-4xl">
          <div className="bg-white/10 backdrop-blur-sm border border-white/20 rounded-2xl p-6 hover:bg-white/15 transition-all hover:scale-105">
            <Truck className="w-8 h-8 text-yellow-300 mb-3 mx-auto" />
            <p className="text-xl font-bold text-white mb-2">Fast Delivery</p>
            <p className="text-sm text-white/70">To your hostel doorstep</p>
          </div>
          <div className="bg-white/10 backdrop-blur-sm border border-white/20 rounded-2xl p-6 hover:bg-white/15 transition-all hover:scale-105">
            <Gift className="w-8 h-8 text-yellow-300 mb-3 mx-auto" />
            <p className="text-xl font-bold text-white mb-2">Free Delivery</p>
            <p className="text-sm text-white/70">For first 101 orders</p>
          </div>
          <div className="bg-white/10 backdrop-blur-sm border border-white/20 rounded-2xl p-6 hover:bg-white/15 transition-all hover:scale-105">
            <Sparkles className="w-8 h-8 text-yellow-300 mb-3 mx-auto" />
            <p className="text-xl font-bold text-white mb-2">Best Prices</p>
            <p className="text-sm text-white/70">Student-friendly rates</p>
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="relative z-10 mt-16 text-sm text-white/70 flex flex-col sm:flex-row gap-3 sm:gap-6 justify-center items-center py-6 border-t border-white/10 w-full max-w-4xl mx-auto">
        <Link href="/privacy-policy" className="hover:text-yellow-300 transition">
          Privacy Policy
        </Link>
        <Link href="/terms-conditions" className="hover:text-yellow-300 transition">
          Terms & Conditions
        </Link>
        <Link
          href="/refund-and-cancellation"
          className="hover:text-yellow-300 transition"
        >
          Refund & Cancellation
        </Link>
        <Link href="/support" className="hover:text-yellow-300 transition">
          Support
        </Link>
        <Link href="/ship-and-delivery" className="hover:text-yellow-300 transition">
          Shipping & Delivery
        </Link>
      </footer>
    </div>
  );
}
