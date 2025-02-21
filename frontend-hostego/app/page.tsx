export default function Home() {
  return (
    <div className="relative flex flex-col items-center justify-center min-h-screen bg-gradient-to-br from-blue-500 to-purple-600 text-white p-8 sm:p-20">
      {/* Background Blurs */}
      <div className="absolute top-0 left-0 w-72 h-72 bg-white/10 blur-3xl rounded-full"></div>
      <div className="absolute bottom-0 right-0 w-72 h-72 bg-white/10 blur-3xl rounded-full"></div>

      {/* Main Content */}
      <main className="relative flex flex-col items-center text-center gap-6">
        {/* Logo */}
    

        {/* Headline */}
        <h1 className="text-4xl font-extrabold tracking-tight sm:text-6xl">
         Hostego  
         <br />
        </h1>
          <span className=" text-4xl font-extrabold tracking-tight sm:text-6xl text-yellow-300 m">  Simplify Your Hostel Life</span>

   
        {/* CTA Button */}
        <button className="mt-4 px-6 py-3 text-lg font-medium text-black bg-yellow-300 rounded-xl shadow-lg hover:scale-105 transition-transform">
          Join the Waitlist 🚀
        </button>
      </main>

      {/* Glassmorphism Info Card */}
      <div className="relative mt-10 p-6 bg-white/10 backdrop-blur-lg rounded-xl shadow-lg text-center max-w-md">
        <h2 className="text-2xl font-bold">Launching Soon!</h2>
        <p className="text-white/80 mt-2">
          Be the first to experience the future of hostel management.
        </p>
      </div>
    </div>
  );
}
