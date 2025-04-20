"use client";
import { useRouter } from "next/navigation";
import { ArrowRight, Truck } from "lucide-react";
import Link from "next/link";

export default function Home() {
  const router = useRouter();

  return (
    <div className="relative flex flex-col items-center justify-center min-h-screen gradient-background text-white p-4 sm:p-8">
      {/* Subtle Background Elements */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute top-0 left-0 w-96 h-96 bg-purple-500/10 blur-3xl rounded-full -translate-x-1/2 -translate-y-1/2"></div>
        <div className="absolute bottom-0 right-0 w-96 h-96 bg-yellow-500/10 blur-3xl rounded-full translate-x-1/2 translate-y-1/2"></div>
      </div>

      {/* Main Content Container */}
      <main className="relative max-w-3xl mx-auto flex flex-col items-center text-center gap-6 p-4">
        {/* Launch Badge */}
        <div className="bg-green-500/10 border border-green-500/20 px-3 py-1 rounded-full">
          <p className="text-sm font-medium text-green-400 flex items-center gap-2">
            <span className="w-1.5 h-1.5 bg-green-500 rounded-full animate-pulse"></span>
            Now Live in Chandigarh University
          </p>
        </div>

        {/* Main Heading */}
        <div className="space-y-1">
          <h1 className="text-4xl sm:text-4xl font-bold tracking-tight">
            Hostego
          </h1>
          <p className="text-sm sm:text-base font-medium text-white/80">
            Your Hostel Life, Simplified
          </p>
        </div>

        {/* Launch Offer */}
        <div className="relative  border  p-6 rounded-xl w-full max-w-md">
          <div className="absolute -top-2.5 left-1/2 -translate-x-1/2">
            <div className="bg-yellow-400 text-black text-sm font-bold px-3 py-0.5 rounded-full">
              LAUNCH OFFER
            </div>
          </div>

          <div className="flex flex-col items-center gap-4 mt-2">
            <div className="flex items-center gap-2">
              <Truck className="w-5 h-5 text-yellow-300" />
              <h2 className="text-3xl font-bold">FREE DELIVERY</h2>
            </div>
            <div className="space-y-1.5">
              <p className="text-xl text-white/90">
                For First{" "}
                <span className="text-yellow-300 font-bold text-2xl">101</span>{" "}
                Orders
              </p>
              <p className="text-sm text-yellow-300 font-medium">
                Limited time offer • Order now before it&apos;s gone
              </p>
            </div>
          </div>
        </div>

        {/* CTA Button */}
        <button
          onClick={() => router.push("/home")}
          className="group px-6 py-2.5 text-sm font-semibold text-black bg-yellow-400 
                   rounded-lg shadow-md hover:shadow-lg transition-all duration-300
                   flex items-center gap-2 hover:gap-3 hover:bg-yellow-300"
        >
          Order Now
          <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
        </button>

        {/* Stats Grid */}
        <div className="grid grid-cols-3 gap-4 mt-8 w-full max-w-2xl">
          {/* { <StatsCard icon={ShoppingBag} value="500+" label="Daily Orders" /> }
        { <StatsCard icon={Users} value="50+" label="Partner Shops" /> }
        { <StatsCard icon={Star} value="4.8" label="User Rating" /> } */}
        </div>
      </main>
      <footer className="relative z-10 mt-16 text-sm text-white/70 flex flex-col sm:flex-row gap-3 sm:gap-6 justify-center items-center py-6 border-t border-white/10 w-full max-w-4xl mx-auto">
        <Link href="/privacy-policy" className="hover:text-white transition">
          Privacy Policy
        </Link>
        <Link href="/terms-conditons" className="hover:text-white transition">
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

// const StatsCard = ({ icon: Icon, value, label }: StatsCardProps) => (
//   <div className="bg-white/5 backdrop-blur-sm border border-white/10 rounded-xl p-4">
//     <Icon className="w-5 h-5 text-yellow-300 mb-2 mx-auto" />
//     <p className="text-2xl font-bold text-white">{value}</p>
//     <p className="text-sm text-white/60">{label}</p>
//   </div>
// );
