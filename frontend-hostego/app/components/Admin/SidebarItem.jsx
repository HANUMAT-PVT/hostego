// Sidebar Item Component
"use client"
const SidebarItem = ({ icon, text, onClick, isActive }) => (
    <button
        className={`flex  items-center gap-3 p-3 text-black hover:bg-gray-200 rounded-lg text-gray-700 transition-all ${isActive ? "bg-[var(--primary-color)] text-white":"text-black"} `}
        onClick={onClick}
    >
       <span className="text-sm">{icon}</span>  <span className="text-xs font-medium whitespace-nowrap">{text}</span>
    </button>
);
export default SidebarItem;