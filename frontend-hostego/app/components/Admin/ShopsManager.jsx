"use client";

import React, { useState, useEffect } from 'react';
import { Plus, Search, MapPin, Clock, Store, Check, X } from 'lucide-react';
import axiosClient from '../../utils/axiosClient';

const ShopForm = ({ onSubmit, onCancel, initialData = null }) => {
    const [formData, setFormData] = useState({
        shop_name: '',
        shop_img: '',
        address: '',
        preparation_time: '30 min',
        food_category: {
            is_veg: 1,
            is_cooked: 0
        },
        shop_status: 1
    });

    useEffect(() => {
        if (initialData) {
            setFormData(initialData);
        }
    }, [initialData]);

    return (
        <div className="bg-white rounded-xl p-6 shadow-sm">
            <h2 className="text-xl font-semibold mb-6">
                {initialData ? 'Edit Shop' : 'Add New Shop'}
            </h2>

            <form onSubmit={(e) => {
                e.preventDefault();
                onSubmit(formData);
            }} className="space-y-4">
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Shop Name</label>
                    <input
                        type="text"
                        value={formData.shop_name}
                        onChange={(e) => setFormData(prev => ({ ...prev, shop_name: e.target.value }))}
                        className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                        required
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Shop Image URL</label>
                    <input
                        type="url"
                        value={formData.shop_img}
                        onChange={(e) => setFormData(prev => ({ ...prev, shop_img: e.target.value }))}
                        className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                        required
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Address</label>
                    <input
                        type="text"
                        value={formData.address}
                        onChange={(e) => setFormData(prev => ({ ...prev, address: e.target.value }))}
                        className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                        required
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Preparation Time</label>
                    <input
                        type="text"
                        value={formData.preparation_time}
                        onChange={(e) => setFormData(prev => ({ ...prev, preparation_time: e.target.value }))}
                        className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                        placeholder="e.g., 30 min"
                        required
                    />
                </div>

                <div className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">Food Category</label>
                        <div className="flex gap-4">
                            <div>
                                <label className="block text-sm mb-2">Food Type</label>
                                <div className="flex gap-3">
                                    <label className="flex items-center gap-2">
                                        <input
                                            type="radio"
                                            checked={formData.food_category.is_veg === 1}
                                            onChange={() => setFormData(prev => ({
                                                ...prev,
                                                food_category: { ...prev.food_category, is_veg: 1 }
                                            }))}
                                            className="text-green-600"
                                        />
                                        <span>Veg</span>
                                    </label>
                                    <label className="flex items-center gap-2">
                                        <input
                                            type="radio"
                                            checked={formData.food_category.is_veg === 0}
                                            onChange={() => setFormData(prev => ({
                                                ...prev,
                                                food_category: { ...prev.food_category, is_veg: 0 }
                                            }))}
                                            className="text-red-600"
                                        />
                                        <span>Non-veg</span>
                                    </label>
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm mb-2">Preparation Required</label>
                                <div className="flex gap-3">
                                    <label className="flex items-center gap-2">
                                        <input
                                            type="radio"
                                            checked={formData.food_category.is_cooked === 1}
                                            onChange={() => setFormData(prev => ({
                                                ...prev,
                                                food_category: { ...prev.food_category, is_cooked: 1 }
                                            }))}
                                        />
                                        <span>Yes</span>
                                    </label>
                                    <label className="flex items-center gap-2">
                                        <input
                                            type="radio"
                                            checked={formData.food_category.is_cooked === 0}
                                            onChange={() => setFormData(prev => ({
                                                ...prev,
                                                food_category: { ...prev.food_category, is_cooked: 0 }
                                            }))}
                                        />
                                        <span>No</span>
                                    </label>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">Shop Status</label>
                        <div className="flex gap-3">
                            <label className="flex items-center gap-2">
                                <input
                                    type="radio"
                                    checked={formData.shop_status === 1}
                                    onChange={() => setFormData(prev => ({ ...prev, shop_status: 1 }))}
                                    className="text-green-600"
                                />
                                <span>Active</span>
                            </label>
                            <label className="flex items-center gap-2">
                                <input
                                    type="radio"
                                    checked={formData.shop_status === 0}
                                    onChange={() => setFormData(prev => ({ ...prev, shop_status: 0 }))}
                                    className="text-red-600"
                                />
                                <span>Inactive</span>
                            </label>
                        </div>
                    </div>
                </div>

                <div className="flex justify-end gap-3 mt-6">
                    <button
                        type="button"
                        onClick={onCancel}
                        className="px-4 py-2 border rounded-lg hover:bg-gray-50 transition-colors"
                    >
                        Cancel
                    </button>
                    <button
                        type="submit"
                        className="px-4 py-2 bg-[var(--primary-color)] text-white rounded-lg hover:opacity-90 transition-opacity"
                    >
                        {initialData ? 'Save Changes' : 'Add Shop'}
                    </button>
                </div>
            </form>
        </div>
    );
};

const ShopsManager = () => {
    const [shops, setShops] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showForm, setShowForm] = useState(false);
    const [editingShop, setEditingShop] = useState(null);
    const [searchTerm, setSearchTerm] = useState('');

    useEffect(() => {
        fetchShops();
    }, []);

    const fetchShops = async () => {
        try {
            setLoading(true);
            const { data } = await axiosClient.get('/api/shop');
            setShops(Object.values(data));
        } catch (error) {
            console.error('Error fetching shops:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleSubmit = async (shopData) => {
        try {
            if (editingShop) {
                await axiosClient.patch(`/api/shop/${editingShop.shop_id}`, shopData);
            } else {
                await axiosClient.post('/api/shop', shopData);
            }
            setShowForm(false);
            setEditingShop(null);
            fetchShops();
        } catch (error) {
            console.error('Error saving shop:', error);
        }
    };

    const handleEdit = (shop) => {
        setEditingShop(shop);
        setShowForm(true);
    };

    const handleCloseForm = () => {
        setShowForm(false);
        setEditingShop(null);
    };

    const filteredShops = shops.filter(shop =>
        shop.shop_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        shop.address.toLowerCase().includes(searchTerm.toLowerCase())
    );

    return (
        <div className="p-6">
            <div className="flex justify-between items-center mb-6">
                <div>
                    <h1 className="text-2xl font-bold text-gray-800">Shops Management</h1>
                    <p className="text-gray-600">Manage your shops</p>
                </div>
                <button
                    onClick={() => setShowForm(true)}
                    className="flex items-center gap-2 bg-[var(--primary-color)] text-white px-4 py-2 rounded-lg hover:opacity-90 transition-opacity"
                >
                    <Plus size={20} />
                    Add Shop
                </button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                <div className="bg-white p-4 rounded-xl shadow-sm">
                    <h3 className="text-sm text-gray-500">Total Shops</h3>
                    <p className="text-2xl font-semibold">{shops.length}</p>
                </div>
                <div className="bg-white p-4 rounded-xl shadow-sm">
                    <h3 className="text-sm text-gray-500">Active Shops</h3>
                    <p className="text-2xl font-semibold text-green-600">
                        {shops.filter(s => s.shop_status === 1).length}
                    </p>
                </div>
            </div>

            <div className="bg-white rounded-xl p-4 shadow-sm mb-6">
                <div className="relative">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={20} />
                    <input
                        type="text"
                        placeholder="Search shops..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="w-full pl-10 pr-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                    />
                </div>
            </div>

            {showForm && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
                    <div className="bg-white rounded-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
                        <ShopForm
                            onSubmit={handleSubmit}
                            onCancel={handleCloseForm}
                            initialData={editingShop}
                        />
                    </div>
                </div>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {filteredShops.map(shop => (
                    <div key={shop.shop_id} className="bg-white rounded-xl shadow-sm overflow-hidden">
                        <img
                            src={shop.shop_img}
                            alt={shop.shop_name}
                            className="w-full h-48 object-cover"
                        />
                        <div className="p-4">
                            <div className="flex justify-between items-start">
                                <h3 className="text-lg font-semibold text-gray-900">{shop.shop_name}</h3>
                                <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${shop.shop_status ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                                    }`}>
                                    {shop.shop_status ? 'Active' : 'Inactive'}
                                </span>
                            </div>
                            <div className="mt-2 space-y-2 text-sm text-gray-500">
                                <div className="flex items-center gap-1">
                                    <MapPin size={16} />
                                    {shop.address}
                                </div>
                                <div className="flex items-center gap-1">
                                    <Clock size={16} />
                                    {shop.preparation_time}
                                </div>
                            </div>
                            <button
                                onClick={() => handleEdit(shop)}
                                className="mt-4 text-[var(--primary-color)] hover:text-[var(--primary-color)]/80 font-medium text-sm"
                            >
                                Edit Shop
                            </button>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default ShopsManager;