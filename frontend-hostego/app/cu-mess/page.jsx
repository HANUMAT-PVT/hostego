"use client"
import axiosClient from "../utils/axiosClient"
import { useEffect, useState } from "react"
import {
    Coffee,
    Sun,
    Moon,
    Calendar,
    ChevronLeft,
    ChevronRight,
    Clock,
    Utensils,
    RefreshCw,
    Cookie,
    Soup as SoupIcon,
    Cookie as CookieIcon,
    Pizza as PizzaIcon,
    Apple as AppleIcon,
    Sandwich as SandwichIcon,
    CircleDot,
    Beef,
    Milk,
    Egg,
    ChevronDown
} from 'lucide-react'
import BackNavigationButton from "../components/BackNavigationButton"
import HostegoLoader from "../components/HostegoLoader"

// Add food icons mapping with fallback
const foodIcons = {
    "poha": Utensils,
    "tea": Coffee,
    "coffee": Coffee,
    "rice": SoupIcon,
    "dal": SoupIcon,
    "roti": CookieIcon,
    "bread": CookieIcon,
    "paneer": Beef,
    "milk": Milk,
    "egg": Egg,
    "sandwich": SandwichIcon,
    "pizza": PizzaIcon,
    "default": CircleDot
}

const getFoodIcon = (itemName) => {
    const key = Object.keys(foodIcons).find(k => itemName.toLowerCase().includes(k))
    return foodIcons[key] || foodIcons.default
}

const MenuItem = ({ item }) => {
    const FoodIcon = getFoodIcon(item)
    return (
        <div className="flex items-center gap-3 p-3 rounded-lg bg-gray-50/80 hover:bg-[var(--primary-color)]/5 group transition-all">
            <div className="p-2 rounded-lg bg-white shadow-sm group-hover:scale-110 group-hover:shadow-md transition-all duration-300">
                <FoodIcon className="w-4 h-4 text-[var(--primary-color)]" />
            </div>
            <span className="text-gray-700 font-medium text-sm group-hover:text-[var(--primary-color)]">
                {item?.trim()}
            </span>
        </div>
    )
}

// Update the MealCard component to include accordion functionality
const MealCard = ({ title, items, icon: Icon, description, className }) => {
    const [isExpanded, setIsExpanded] = useState(false)

    return (
        <div className={`bg-white rounded-xl overflow-hidden shadow-sm border border-gray-100 hover:shadow-lg transition-all duration-300 ${className || ''}`}>
            {/* Header with gradient - Now clickable */}
            <button
                onClick={() => setIsExpanded(!isExpanded)}
                className="w-full p-4 bg-gradient-to-r from-[var(--primary-color)] to-purple-400 flex items-center justify-between"
            >
                <div className="flex items-center gap-3">
                    <div className="bg-white/10 p-2.5 rounded-lg">
                        <Icon className="w-5 h-5 text-white" />
                    </div>
                    <div className="text-left">
                        <h3 className="text-white font-semibold text-lg">{title}</h3>
                        <p className="text-white text-sm">{description}</p>
                    </div>
                </div>
                <ChevronDown
                    className={`w-5 h-5 text-white transition-transform duration-200 ${isExpanded ? 'rotate-180' : ''}`}
                />
            </button>

            {/* Menu Items - Now collapsible */}
            <div className={`transition-all duration-200 ${isExpanded ? ' opacity-100' : 'max-h-0 opacity-0'} overflow-hidden`}>
                <div className="p-4">
                    {items?.length > 0 ? (
                        <div className="space-y-2.5">
                            {items?.map((item, index) => (
                                <MenuItem key={index} item={item} />
                            ))}
                        </div>
                    ) : (
                        <div className="text-center py-6">
                            <div className="w-12 h-12 mx-auto mb-3 rounded-full bg-gray-100 flex items-center justify-center">
                                <Utensils className="w-6 h-6 text-gray-400" />
                            </div>
                            <p className="text-gray-500 text-sm">No items available</p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    )
}

const TimingCard = ({ icon: Icon, title, time, color, holiday }) => (
    <div className={`flex items-center gap-4 p-4 rounded-lg ${color} transition-all hover:shadow-md cursor-pointer`}>
        <div className="p-2.5 rounded-lg bg-white/90 shadow-sm">
            <Icon className="w-5 h-5" />
        </div>
        <div>
            <h4 className="font-medium text-sm">{title}</h4>
            <p className="text-xs opacity-75">{time}</p>
            <p className="text-xs opacity-75">Holiday: {holiday}</p>
        </div>
    </div>
)



// Update the TimingsAccordion component
const TimingsAccordion = ({ isOpen, setIsOpen }) => {
    return (
        <div className="bg-white rounded-xl shadow-sm overflow-hidden">
            <button
                onClick={() => setIsOpen(!isOpen)}
                className="w-full px-6 py-4 flex items-center justify-between hover:bg-gray-50"
            >
                <div className="flex items-center gap-3">
                    <Clock className="w-5 h-5 text-[var(--primary-color)]" />
                    <span className="font-medium text-gray-900">Mess Timings</span>
                </div>
                <ChevronDown className={`w-5 h-5 text-gray-500 transition-transform duration-200 ${isOpen ? 'rotate-180' : ''}`} />
            </button>

            {isOpen && (
                <div className="border-t">
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                        <TimingCard
                            icon={Coffee}
                            title="Breakfast"
                            time="7:30 AM - 9:00 AM"
                            holiday="8:00 AM - 9:30 AM"
                            color="bg-green-50 text-green-600"
                        />
                        <TimingCard
                            icon={Sun}
                            title="Lunch"
                            time="12:00 PM - 1:45 PM"
                            holiday="12:30 PM - 2:00 PM"
                            color="bg-blue-50 text-blue-600"
                        />
                        <TimingCard
                            icon={Sun}
                            title="Snacks"
                            time="4:30 PM - 5:15 PM"
                            holiday="4:30 PM - 5:15 PM"
                            color="bg-blue-50 text-blue-600"
                        />

                        {/* Replace single snacks card with two cards */}


                        <TimingCard
                            icon={Moon}
                            title="Dinner"
                            time="7:30 PM - 9:00 PM"
                            holiday="7:30 PM - 9:00 PM"
                            color="bg-purple-50 text-purple-600"
                        />
                    </div>
                </div>
            )}
        </div>
    )
}

const Page = () => {
    const [messMenus, setMessMenus] = useState([])
    const [selectedDate, setSelectedDate] = useState(() => {
        const today = new Date()
        today.setHours(0, 0, 0, 0)
        return today
    })
    const [isLoading, setIsLoading] = useState(true)
    const [isTimingsOpen, setIsTimingsOpen] = useState(false)

    useEffect(() => {
        fetchMessMenu()
    }, [])

    const fetchMessMenu = async () => {
        try {
            setIsLoading(true)
            const { data } = await axiosClient.get("/api/mess-menu")
            setMessMenus(data)
        } catch (error) {
            console.error("Error fetching mess menu:", error)
        } finally {
            setIsLoading(false)
        }
    }

    const parseMenu = (menuString) => {
        const meals = menuString.split(';')
        const menuObject = {}

        meals.forEach(meal => {
            const [type, items] = meal.split(':')
            const key = type.trim().toLowerCase().replace(/_/g, '_')
            menuObject[key] = items ? items.split(',') : []
        })

        return menuObject
    }

    const formatDate = (date) => {
        return new Date(date).toLocaleDateString('en-US', {
            weekday: 'long',
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        })
    }

    // Update the formatDateForComparison function to handle timezone issues
    const formatDateForComparison = (date) => {
        const d = new Date(date)
        return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
    }

    // Update the findMenuForDate function to be more robust
    const findMenuForDate = (date) => {
        const formattedSelectedDate = formatDateForComparison(date)
        return messMenus.find(menu => {
            const menuDate = formatDateForComparison(menu.date)
            return menuDate === formattedSelectedDate
        })
    }

    const changeDate = (days) => {
        const newDate = new Date(selectedDate)
        newDate.setDate(selectedDate.getDate() + days)
        setSelectedDate(newDate)
    }

    if (isLoading) {
        return (
            <HostegoLoader />
        )
    }

    const selectedMenu = findMenuForDate(selectedDate)
    const menu = selectedMenu ? parseMenu(selectedMenu.menu) : {}

    return (
        <>
            <BackNavigationButton title="CU Mess Menu" />
            <div className="min-h-screen bg-[var(--bg-page-color)]">
                <div className="max-w-7xl mx-auto p-2">
                    {/* Header Section */}
                    <div className="bg-white rounded-xl p-4 shadow-sm mb-6">
                        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6">
                            <div className="flex items-center gap-4">
                                <div className="p-3 rounded-full bg-[var(--primary-color)]/10">
                                    <Utensils className="w-6 h-6 text-[var(--primary-color)]" />
                                </div>
                                <div>
                                    <h1 className="text-xl font-bold text-gray-800">CU Mess Menu</h1>
                                    <p className="text-gray-500 text-sm">Explore today's delicious meals</p>
                                </div>
                            </div>

                            {/* Date Navigation */}
                            <div className="flex items-center w-full m-auto gap-2 bg-gray-50 p-2 rounded-lg">
                                <button
                                    onClick={() => changeDate(-1)}
                                    className="p-2 rounded-lg hover:bg-white text-gray-600 transition-all"
                                >
                                    <ChevronLeft className="w-5 h-5" />
                                </button>
                                <div className="text-center flex-1">
                                    <p className="text-sm font-medium text-gray-900">{formatDate(selectedDate)}</p>
                                </div>
                                <button
                                    onClick={() => changeDate(1)}
                                    className="p-2 rounded-lg hover:bg-white text-gray-600 transition-all"
                                >
                                    <ChevronRight className="w-5 h-5" />
                                </button>
                                <button
                                    onClick={() => {
                                        const today = new Date()
                                        today.setHours(0, 0, 0, 0)
                                        setSelectedDate(today)
                                    }}
                                    className="px-3 py-1 text-sm bg-[var(--primary-color)] text-white rounded-lg hover:bg-[var(--primary-color)]/90 transition-all"
                                >
                                    Today
                                </button>
                            </div>
                        </div>

                        {/* Replace your timing cards section with this */}
                        <TimingsAccordion
                            isOpen={isTimingsOpen}
                            setIsOpen={setIsTimingsOpen}
                        />
                    </div>

                    {/* Menu Grid */}
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-6 gap-6 p-2">
                        <div className="lg:col-span-2">
                            <MealCard
                                title="Breakfast"
                                description="Start your day with energy"
                                items={menu?.breakfast || []}
                                icon={Coffee}
                            />
                        </div>
                        <div className="lg:col-span-2">
                            <MealCard
                                title="Lunch"
                                description="Nutritious afternoon meal"
                                items={menu?.lunch || []}
                                icon={Sun}
                            />
                        </div>
                        <div className="lg:col-span-1">
                            <MealCard
                                title="NC Boys' Snacks"
                                description="Evening refreshment"
                                items={menu?.snacks_nc_boys || []}
                                icon={Cookie}
                            />
                        </div>
                        <div className="lg:col-span-1">
                            <MealCard
                                title="Zakir Boys' Snacks"
                                description="Evening refreshment"
                                items={menu?.snacks_zakir_boys || []}
                                icon={Cookie}
                            />
                        </div>
                        <div className="lg:col-span-2">
                            <MealCard
                                title="Girls' Snacks"
                                description="Evening refreshment"
                                items={menu?.snacks_girls || []}
                                icon={Cookie}
                            />
                        </div>
                        <div className="lg:col-span-2">
                            <MealCard
                                title="Dinner"
                                description="Complete your day with taste"
                                items={menu?.dinner || []}
                                icon={Moon}
                            />
                        </div>
                    </div>

                    {/* Empty State */}
                    {!selectedMenu && (
                        <div className="text-center mt-8 bg-white rounded-xl p-8 shadow-sm border border-gray-100">
                            <div className="w-20 h-20 mx-auto mb-4 flex items-center justify-center rounded-full bg-[var(--primary-color)]/10">
                                <Calendar className="w-10 h-10 text-[var(--primary-color)]" />
                            </div>
                            <h3 className="text-xl font-semibold text-gray-900 mb-2">No Menu Available</h3>
                            <p className="text-gray-500 max-w-md mx-auto">
                                The menu for {formatDate(selectedDate)} hasn't been uploaded yet.
                            </p>
                        </div>
                    )}
                </div>
            </div>
        </>
    )
}

export default Page