'use client'
import React, { useState, useEffect } from 'react'
import { Calendar, Plus, Pencil, Trash2, Save, X, RefreshCw, Coffee, Sun, Moon, Cookie } from 'lucide-react'
import axiosClient from '@/app/utils/axiosClient'
import ConfirmationPopup from '../ConfirmationPopup'

const MenuEditor = ({ menu, onSave, onCancel }) => {
    const [menuData, setMenuData] = useState({
        date: menu?.date?.split('T')[0] || '',
        breakfast: '',
        lunch: '',
        snacks_nc_boys: '',
        snacks_zakir_boys: '',
        snacks_girls: '',
        dinner: ''
    })

    useEffect(() => {
        if (menu?.menu) {
            const parts = menu.menu.split(';')
            const menuObj = {}
            parts.forEach(part => {
                const [meal, items] = part.split(':')
                menuObj[meal.trim()?.toLowerCase()] = items.trim()
            })
            setMenuData(prev => ({
                ...prev,
                breakfast: menuObj?.breakfast || '',
                lunch: menuObj?.lunch || '',
                snacks_nc_boys: menuObj?.snacks_nc_boys || '',
                snacks_zakir_boys: menuObj?.snacks_zakir_boys || '',
                snacks_girls: menuObj?.snacks_girls || '',
                dinner: menuObj?.dinner || ''
            }))
        }
    }, [menu])

    const handleSave = () => {
        const formattedMenu = `Breakfast: ${menuData?.breakfast}; Lunch: ${menuData?.lunch}; Snacks_NC_Boys: ${menuData?.snacks_nc_boys}; Snacks_Zakir_Boys: ${menuData?.snacks_zakir_boys}; Snacks_Girls: ${menuData?.snacks_girls}; Dinner: ${menuData?.dinner}`
        onSave({ ...menu, date: menuData?.date, menu: formattedMenu })
    }
    const statsMeals = [
        { meal: 'breakfast', icon: Coffee, color: 'bg-orange-50 text-orange-600' },
        { meal: 'lunch', icon: Sun, color: 'bg-blue-50 text-blue-600' },
        {
            meal: 'snacks_nc_boys',
            label: "NC Boys' Snacks",
            icon: Cookie,
            color: 'bg-amber-50 text-amber-600'
        },
        {
            meal: 'snacks_zakir_boys',
            label: "Zakir Boys' Snacks",
            icon: Cookie,
            color: 'bg-amber-50 text-amber-600'
        },
        {
            meal: 'snacks_girls',
            label: "Girls' Snacks",
            icon: Cookie,
            color: 'bg-pink-50 text-pink-600'
        },
        { meal: 'dinner', icon: Moon, color: 'bg-purple-50 text-purple-600' }
    ]

    return (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-xl w-full max-w-2xl p-6">
                <div className="flex items-center justify-between mb-6">
                    <h3 className="text-xl font-semibold">
                        {menu?.id ? 'Edit Menu' : 'Add New Menu'}
                    </h3>
                    <button onClick={onCancel} className="p-2 hover:bg-gray-100 rounded-lg">
                        <X className="w-5 h-5" />
                    </button>
                </div>

                <div className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Date</label>
                        <input
                            type="date"
                            value={menuData.date}
                            onChange={(e) => setMenuData(prev => ({ ...prev, date: e.target.value }))}
                            className="w-full p-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                        />
                    </div>

                    {statsMeals?.map(({ meal, label, icon: Icon, color }) => (
                        <div key={meal}>
                            <label className="block text-sm font-medium text-gray-700 mb-1 capitalize">
                                <div className="flex items-center gap-2">
                                    <div className={`p-1.5 rounded-lg ${color}`}>
                                        <Icon className="w-4 h-4" />
                                    </div>
                                    {label || meal?.replace('_', ' ')}
                                </div>
                            </label>
                            <input
                                type="text"
                                value={menuData[meal]}
                                onChange={(e) => setMenuData(prev => ({ ...prev, [meal]: e.target.value }))}
                                placeholder={`Enter ${label || meal} items (comma separated)`}
                                className="w-full p-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                            />
                        </div>
                    ))}

                    <div className="flex justify-end gap-3 mt-6">
                        <button
                            onClick={onCancel}
                            className="px-4 py-2 border rounded-lg hover:bg-gray-50"
                        >
                            Cancel
                        </button>
                        <button
                            onClick={handleSave}
                            className="px-4 py-2 bg-[var(--primary-color)] text-white rounded-lg hover:bg-[var(--primary-color)]/90"
                        >
                            Save Menu
                        </button>
                    </div>
                </div>
            </div>
        </div>
    )
}

const CuMessManager = () => {
    const [menus, setMenus] = useState([])
    const [isLoading, setIsLoading] = useState(true)
    const [editingMenu, setEditingMenu] = useState(null)
    const [deleteMenu, setDeleteMenu] = useState(null)
    const [isRefreshing, setIsRefreshing] = useState(false)

    useEffect(() => {
        fetchMenus()
    }, [])

    const fetchMenus = async (showRefreshAnimation = false) => {
        try {
            showRefreshAnimation ? setIsRefreshing(true) : setIsLoading(true)
            const { data } = await axiosClient.get('/api/mess-menu')
            setMenus(data)
        } catch (error) {
            console.error('Error fetching menus:', error)
        } finally {
            setIsRefreshing(false)
            setIsLoading(false)
        }
    }

    const handleSave = async (menuData) => {
        try {
            if (menuData.id) {
                await axiosClient.patch(`/api/mess-menu/${menuData.id}`, menuData)
            } else {
                await axiosClient.post('/api/mess-menu', menuData)
            }
            fetchMenus()
            setEditingMenu(null)
        } catch (error) {
            console.error('Error saving menu:', error)
        }
    }

    const handleDelete = async () => {
        try {
            await axiosClient.delete(`/api/mess-menu/${deleteMenu.id}`)
            fetchMenus()
            setDeleteMenu(null)
        } catch (error) {
            console.error('Error deleting menu:', error)
        }
    }

    return (
        <div className="max-w-6xl mx-auto p-6">
            <div className="flex items-center justify-between mb-6">
                <div className="flex items-center gap-4">
                    <div className="p-3 rounded-full bg-[var(--primary-color)]/10">
                        <Calendar className="w-6 h-6 text-[var(--primary-color)]" />
                    </div>
                    <div>
                        <h1 className="text-2xl font-bold text-gray-800">Mess Menu Manager</h1>
                        <p className="text-gray-500">Manage daily mess menus</p>
                    </div>
                </div>
                <div className="flex items-center gap-3">
                    <button
                        onClick={() => fetchMenus(true)}
                        disabled={isRefreshing}
                        className="flex items-center gap-2 px-4 py-2 rounded-lg bg-[var(--primary-color)]/10 
                                 text-[var(--primary-color)] font-medium hover:bg-[var(--primary-color)]/20"
                    >
                        <RefreshCw className={`w-4 h-4 ${isRefreshing ? 'animate-spin' : ''}`} />
                        {isRefreshing ? 'Refreshing...' : 'Refresh'}
                    </button>
                    <button
                        onClick={() => setEditingMenu({})}
                        className="flex items-center gap-2 px-4 py-2 bg-[var(--primary-color)] text-white rounded-lg 
                                 hover:bg-[var(--primary-color)]/90"
                    >
                        <Plus className="w-4 h-4" />
                        Add Menu
                    </button>
                </div>
            </div>

            {/* Menu List */}
            <div className="bg-white rounded-xl shadow-sm overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="w-full">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="px-6 py-3 text-left text-sm font-medium text-gray-500">Date</th>
                                <th className="px-6 py-3 text-left text-sm font-medium text-gray-500">Menu</th>
                                <th className="px-6 py-3 text-left text-sm font-medium text-gray-500">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100">
                            {menus.map(menu => (
                                <tr key={menu.id} className="hover:bg-gray-50">
                                    <td className="px-6 py-4 text-sm">
                                        {new Date(menu.date).toLocaleDateString('en-US', {
                                            weekday: 'long',
                                            year: 'numeric',
                                            month: 'long',
                                            day: 'numeric'
                                        })}
                                    </td>
                                    <td className="px-6 py-4 text-sm">
                                        {menu.menu.split(';').map((meal, index) => {
                                            const [type, items] = meal.split(':')
                                            return (
                                                <div key={index} className="mb-1 last:mb-0">
                                                    <span className="font-medium">{type}:</span>
                                                    <span className="text-gray-600">{items}</span>
                                                </div>
                                            )
                                        })}
                                    </td>
                                    <td className="px-6 py-4">
                                        <div className="flex items-center gap-2">
                                            <button
                                                onClick={() => setEditingMenu(menu)}
                                                className="p-2 text-blue-600 hover:bg-blue-50 rounded-lg"
                                            >
                                                <Pencil className="w-4 h-4" />
                                            </button>
                                            <button
                                                onClick={() => setDeleteMenu(menu)}
                                                className="p-2 text-red-600 hover:bg-red-50 rounded-lg"
                                            >
                                                <Trash2 className="w-4 h-4" />
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>

            {/* Editor Modal */}
            {editingMenu && (
                <MenuEditor
                    menu={editingMenu}
                    onSave={handleSave}
                    onCancel={() => setEditingMenu(null)}
                />
            )}

            {/* Delete Confirmation */}
            <ConfirmationPopup
                isOpen={!!deleteMenu}
                title="Delete Menu"
                message={`Are you sure you want to delete the menu for ${deleteMenu?.date}?`}
                onConfirm={handleDelete}
                onCancel={() => setDeleteMenu(null)}
                variant="danger"
            />
        </div>
    )
}

export default CuMessManager