// Sidebar Item Component
const SidebarItem = ({ icon, text, onClick, isActive }) => (
    <button
        className={`  flex items-center gap-3 p-3 text-black hover:bg-gray-200 rounded-lg text-gray-700 transition-all ${isActive ? "bg-[var(--primary-color)] text-white":"text-black"} `}
        onClick={onClick}
    >
        {icon} <span className="text-sm font-medium">{text}</span>
    </button>
);
export default SidebarItem;