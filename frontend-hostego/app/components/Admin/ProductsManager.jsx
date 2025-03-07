"use client";

import React, { useState, useEffect } from 'react';
import { Plus, Package, Search, Tag, Clock, DollarSign, Info, Image as ImageIcon, Check, X } from 'lucide-react';
import axiosClient from '../../utils/axiosClient';

const ProductForm = ({ onSubmit, onCancel, initialData = null }) => {
    const [shops, setShops] = useState([]);
    const [formData, setFormData] = useState({
        product_name: '',
        food_category: {
            is_veg: 1,
            is_cooked: 0
        },
        food_price: '',
        availability: 1,
        product_img_url: '',
        description: '',
        discount: {
            is_available: 0,
            percentage: 0
        },
        preparation_time: '0 min',
        shop_id: '',
        tags: [],
        weight: ''
    });

    const [tagInput, setTagInput] = useState('');

    useEffect(() => {
        const fetchShops = async () => {
            try {
                const { data } = await axiosClient.get('/api/shop');
                setShops(data);
                if (data.length > 0) {
                    setFormData(prev => ({ ...prev, shop_id: data[0].shop_id }));
                }
            } catch (error) {
                console.error('Error fetching shops:', error);
            }
        };

        fetchShops();
    }, []);

    useEffect(() => {
        if (initialData) {
            setFormData(initialData);
        }
    }, [initialData]);

    const handleAddTag = () => {
        if (tagInput.trim()) {
            setFormData(prev => ({
                ...prev,
                tags: [...prev.tags, tagInput.trim().toLowerCase()]
            }));
            setTagInput('');
        }
    };

    const removeTag = (tagToRemove) => {
        setFormData(prev => ({
            ...prev,
            tags: prev.tags.filter(tag => tag !== tagToRemove)
        }));
    };

    return (
        <div className="bg-white rounded-xl p-6 shadow-sm">
            <h2 className="text-xl font-semibold mb-6">
                {initialData ? 'Edit Product' : 'Add New Product'}
            </h2>

            <form onSubmit={(e) => {
                e.preventDefault();
                onSubmit(formData);
            }} className="space-y-4">
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Select Shop <span className="text-red-500">*</span>
                    </label>
                    <select
                        value={formData.shop_id}
                        onChange={(e) => setFormData(prev => ({ ...prev, shop_id: e.target.value }))}
                        className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)] bg-white"
                        required
                    >
                        <option value="">Select a shop</option>
                        {shops.map(shop => (
                            <option key={shop.shop_id} value={shop.shop_id}>
                                {shop.shop_name}
                            </option>
                        ))}
                    </select>
                </div>

                {shops.length === 0 && (
                    <div className="bg-yellow-50 border-l-4 border-yellow-400 p-4">
                        <div className="flex items-center">
                            <Info className="text-yellow-400 mr-2" size={20} />
                            <p className="text-sm text-yellow-700">
                                No shops available. Please create a shop first.
                            </p>
                        </div>
                    </div>
                )}

                <div className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Product Name</label>
                        <input
                            type="text"
                            value={formData.product_name}
                            onChange={(e) => setFormData(prev => ({ ...prev, product_name: e.target.value }))}
                            className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                            required
                        />
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">Price (₹)</label>
                            <input
                                type="number"
                                value={formData.food_price}
                                onChange={(e) => setFormData(prev => ({ ...prev, food_price: Number(e.target.value) }))}
                                className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                                required
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">Weight</label>
                            <input
                                type="text"
                                value={formData.weight}
                                onChange={(e) => setFormData(prev => ({ ...prev, weight: e.target.value }))}
                                className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                                required
                            />
                        </div>
                    </div>

                    <div className="flex gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">Food Type</label>
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
                            <label className="block text-sm font-medium text-gray-700 mb-2">Preparation Required</label>
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

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Product Image URL</label>
                        <input
                            type="url"
                            value={formData.product_img_url}
                            onChange={(e) => setFormData(prev => ({ ...prev, product_img_url: e.target.value }))}
                            className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
                        <textarea
                            value={formData.description}
                            onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
                            className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                            rows={3}
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Tags</label>
                        <div className="flex gap-2 mb-2 flex-wrap">
                            {formData.tags.map(tag => (
                                <span
                                    key={tag}
                                    className="bg-gray-100 px-2 py-1 rounded-full text-sm flex items-center gap-1"
                                >
                                    {tag}
                                    <button
                                        type="button"
                                        onClick={() => removeTag(tag)}
                                        className="text-gray-500 hover:text-red-500"
                                    >
                                        <X size={14} />
                                    </button>
                                </span>
                            ))}
                        </div>
                        <div className="flex gap-2">
                            <input
                                type="text"
                                value={tagInput}
                                onChange={(e) => setTagInput(e.target.value)}
                                onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), handleAddTag())}
                                placeholder="Add tags..."
                                className="flex-1 px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                            />
                            <button
                                type="button"
                                onClick={handleAddTag}
                                className="px-4 py-2 bg-gray-100 rounded-lg hover:bg-gray-200"
                            >
                                Add
                            </button>
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
                        disabled={!formData.shop_id || shops.length === 0}
                        className="px-4 py-2 bg-[var(--primary-color)] text-white rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {initialData ? 'Save Changes' : 'Add Product'}
                    </button>
                </div>
            </form>
        </div>
    );
};

const ProductsManager = () => {
    const [showForm, setShowForm] = useState(false);
    const [editingProduct, setEditingProduct] = useState(null);
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');

    useEffect(() => {
        fetchProducts();
    }, []);

    const fetchProducts = async () => {
        try {
            setLoading(true);
            const { data } = await axiosClient.get('/api/products/all?page=1&limit=40');
            setProducts(data);
        } catch (error) {
            console.error('Error fetching products:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleEditProduct = (product) => {
        setEditingProduct(product);
        setShowForm(true);
    };

    const handleSubmit = async (productData) => {
        try {
            if (editingProduct) {
                // Update existing product
                await axiosClient.patch(`/api/products/${editingProduct?.product_id}`, productData);
            } else {
                // Add new product
                await axiosClient.post('/api/products', productData);
            }
            setShowForm(false);
            setEditingProduct(null);
            fetchProducts();
        } catch (error) {
            console.error('Error saving product:', error);
        }
    };

    const handleCloseForm = () => {
        setShowForm(false);
        setEditingProduct(null);
    };

    const filteredProducts = products.filter(product =>
        product.product_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        product.description.toLowerCase().includes(searchTerm.toLowerCase())
    );

  return (
        <div className="p-6">
            <div className="flex justify-between items-center mb-6">
    <div>
                    <h1 className="text-2xl font-bold text-gray-800">Products Management</h1>
                    <p className="text-gray-600">Manage your product catalog</p>
                </div>
                <button
                    onClick={() => setShowForm(true)}
                    className="flex items-center gap-2 bg-[var(--primary-color)] text-white px-4 py-2 rounded-lg hover:opacity-90 transition-opacity"
                >
                    <Plus size={20} />
                    Add Product
                </button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
                <div className="bg-white p-4 rounded-xl shadow-sm">
                    <h3 className="text-sm text-gray-500">Total Products</h3>
                    <p className="text-2xl font-semibold">{products.length}</p>
                </div>
                <div className="bg-white p-4 rounded-xl shadow-sm">
                    <h3 className="text-sm text-gray-500">Available Products</h3>
                    <p className="text-2xl font-semibold text-green-600">
                        {products.filter(p => p.availability === 1).length}
                    </p>
                </div>
                <div className="bg-white p-4 rounded-xl shadow-sm">
                    <h3 className="text-sm text-gray-500">Out of Stock</h3>
                    <p className="text-2xl font-semibold text-red-600">
                        {products.filter(p => p.availability === 0).length}
                    </p>
                </div>
            </div>

            <div className="bg-white rounded-xl p-4 shadow-sm mb-6">
                <div className="relative">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={20} />
                    <input
                        type="text"
                        placeholder="Search products..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="w-full pl-10 pr-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                    />
                </div>
            </div>

            {showForm && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
                    <div className="bg-white rounded-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
                        <ProductForm
                            onSubmit={handleSubmit}
                            onCancel={handleCloseForm}
                            initialData={editingProduct}
                        />
                    </div>
                </div>
            )}

            <div className="bg-white rounded-xl shadow-sm overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="w-full">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Product</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Shop</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Price</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tags</th>
                                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Actions
                                </th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-200">
                            {filteredProducts.map((product) => (
                                <tr key={product.product_id} className="hover:bg-gray-50">
                                    <td className="px-6 py-4">
                                        <div className="flex items-center">
                                            <img
                                                src={product.product_img_url}
                                                alt={product.product_name}
                                                className="w-10 h-10 rounded-lg object-cover"
                                            />
                                            <div className="ml-4">
                                                <div className="text-sm font-medium text-gray-900">{product.product_name}</div>
                                                <div className="text-sm text-gray-500">{product.weight}</div>
                                            </div>
                                        </div>
                                    </td>
                                    <td className="px-6 py-4">
                                        <span className="text-sm text-gray-900">
                                            {product.shop?.shop_name || 'Unknown Shop'}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4">
                                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${product.food_category.is_veg ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                                            }`}>
                                            {product.food_category.is_veg ? 'Veg' : 'Non-veg'}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4 text-sm text-gray-500">
                                        ₹{product.food_price}
                                    </td>
                                    <td className="px-6 py-4">
                                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${product.availability ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                                            }`}>
                                            {product.availability ? 'Available' : 'Out of Stock'}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4">
                                        <div className="flex flex-wrap gap-1">
                                            {product.tags.map(tag => (
                                                <span
                                                    key={tag}
                                                    className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800"
                                                >
                                                    {tag}
                                                </span>
                                            ))}
                                        </div>
                                    </td>
                                    <td className="px-6 py-4 text-right">
                                        <button
                                            onClick={() => handleEditProduct(product)}
                                            className="text-[var(--primary-color)] hover:text-[var(--primary-color)]/80 font-medium text-sm"
                                        >
                                            Edit
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>
    </div>
    );
};

export default ProductsManager;
