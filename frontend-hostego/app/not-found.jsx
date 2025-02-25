'use client';

import { useRouter } from 'next/navigation';
import { Home, ArrowLeft, Package } from 'lucide-react';

const NotFound = () => {
  const router = useRouter();

  return (
    <div className="min-h-screen bg-[var(--bg-page-color)] flex items-center justify-center p-4">
      <div className="relative max-w-md w-full">
        {/* Background Blur Effect */}
        <div className="absolute -top-20 -left-20 w-40 h-40 bg-[#3b82f6] rounded-full blur-[100px] opacity-20" />
        <div className="absolute -bottom-20 -right-20 w-40 h-40 bg-[#9333ea] rounded-full blur-[100px] opacity-20" />

        {/* Main Content */}
        <div className="relative bg-white/80 backdrop-blur-xl rounded-3xl shadow-2xl p-8 text-center animate-slide-up border border-white/20">
          {/* 404 Text */}
          <div className="relative mb-4">
            <div className="text-[100px] font-bold leading-none bg-gradient-to-r from-[#3b82f6] to-[#9333ea] text-transparent bg-clip-text">
              404
            </div>
            <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-[130px] font-bold leading-none text-white/5">
              404
            </div>

          </div>

          {/* Content */}
          <h1 className="text-3xl font-bold mb-3 bg-gradient-to-r from-[#3b82f6] to-[#9333ea] text-transparent bg-clip-text">
            Page Not Found
          </h1>
          <p className="text-gray-600 mb-10">
            {/* Oops! Looks like this page took a wrong turn. Don't worry, we'll help you get back on track! */}
          </p>

          {/* Action Buttons */}
          <div className="space-y-3">
            <button
              onClick={() => router.push('/')}
              className="w-full py-4 px-6 bg-gradient-to-r from-[#3b82f6] to-[#9333ea] text-white rounded-2xl font-medium flex items-center justify-center gap-2 transition-all hover:shadow-lg hover:scale-[1.02]"
            >
              <Home className="w-5 h-5" />
              Return Home
            </button>

            <button
              onClick={() => router.back()}
              className="w-full py-4 px-6 bg-white text-gray-700 rounded-2xl font-medium flex items-center justify-center gap-2 transition-all border-2 border-gray-100 hover:border-[var(--primary-color)] hover:text-[var(--primary-color)]"
            >
              <ArrowLeft className="w-5 h-5" />
              Go Back
            </button>
          </div>

          {/* Help Link */}
          <div className="mt-10 pt-6 border-t border-gray-100">
            <div className="inline-flex items-center gap-2 text-sm text-gray-500 hover:text-[var(--primary-color)] transition-colors">
              <span>Need help?</span>
              <a
                href="mailto:support@example.com"
                className="font-medium underline decoration-dashed underline-offset-4"
              >
                Contact Support
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default NotFound;