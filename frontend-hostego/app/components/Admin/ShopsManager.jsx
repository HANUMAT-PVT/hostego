"use client";

import React, { useState, useEffect } from 'react';
import { Plus, Search, MapPin, Clock, Store, Check, X, Edit, Eye, Upload, Star, Phone, Mail, User, Building, FileText, Banknote, Shield } from 'lucide-react';
import axiosClient from '../../utils/axiosClient';
import { uploadToS3Bucket } from '../../lib/aws';

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
        shop_status: 1,
        shop_verification_status: 'pending',
        is_shop_verified: false,
        owner_name: '',
        owner_phone: '',
        owner_email: '',
        latitude: 0.0,
        longitude: 0.0,
        shop_type: 'Restaurant',
        shop_description: '',
        fssai_license_number: '',
        gstin_number: '',
        fssai_liscense_copy: '',
        gstin_copy: '',
        bank_name: '',
        bank_account_number: '',
        bank_ifsc_code: '',
        bank_account_holder_name: '',
        bank_account_type: 'Savings',
        pancard_copy: '',
        outlet_open_time: '',
        outlet_close_time: ''
    });

    const [uploading, setUploading] = useState(false);
    const [uploadProgress, setUploadProgress] = useState({});

    useEffect(() => {
        if (initialData) {
            setFormData(initialData);
        }
    }, [initialData]);

    const handleFileUpload = async (file, fieldName) => {
        try {
            setUploading(true);
            setUploadProgress(prev => ({ ...prev, [fieldName]: 'Uploading...' }));

            const uploadedUrl = await uploadToS3Bucket(file);
            setFormData(prev => ({ ...prev, [fieldName]: uploadedUrl }));
            setUploadProgress(prev => ({ ...prev, [fieldName]: 'Uploaded!' }));

            setTimeout(() => {
                setUploadProgress(prev => ({ ...prev, [fieldName]: '' }));
            }, 2000);
        } catch (error) {
            console.error('Upload error:', error);
            setUploadProgress(prev => ({ ...prev, [fieldName]: 'Upload failed' }));
        } finally {
            setUploading(false);
        }
    };

    const FileUploadField = ({ label, fieldName, accept = "image/*" }) => (
        <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">{label}</label>
            <div className="flex items-center gap-2">
                <input
                    type="file"
                    accept={accept}
                    onChange={(e) => {
                        const file = e.target.files[0];
                        if (file) handleFileUpload(file, fieldName);
                    }}
                    className="hidden"
                    id={fieldName}
                />
                <label
                    htmlFor={fieldName}
                    className="flex items-center gap-2 px-3 py-2 border rounded-lg cursor-pointer hover:bg-gray-50 transition-colors"
                >
                    <Upload size={16} />
                    Upload {label}
                </label>
                {formData[fieldName] && (
                    <span className="text-xs text-green-600 flex items-center gap-1">
                        <Check size={12} />
                        Uploaded
                    </span>
                )}
            </div>
            {uploadProgress[fieldName] && (
                <p className="text-xs text-blue-600 mt-1">{uploadProgress[fieldName]}</p>
            )}
            {formData[fieldName] && (
                <a href={formData[fieldName]} target="_blank" rel="noopener noreferrer" className="text-xs text-blue-600 hover:underline">
                    View uploaded file
                </a>
            )}
        </div>
    );

    return (
        <div className="bg-white rounded-xl p-6 shadow-sm max-h-[90vh] overflow-y-auto">
            <h2 className="text-xl font-semibold mb-6">
                {initialData ? 'Edit Shop' : 'Add New Shop'}
            </h2>

            <form onSubmit={(e) => {
                e.preventDefault();
                onSubmit(formData);
            }} className="space-y-4">
                {/* Basic Shop Information */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
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
                        <label className="block text-sm font-medium text-gray-700 mb-1">Shop Type</label>
                        <select
                            value={formData.shop_type}
                            onChange={(e) => setFormData(prev => ({ ...prev, shop_type: e.target.value }))}
                            className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                        >
                            <option value="Restaurant">Restaurant</option>
                            <option value="Cafe">Cafe</option>
                            <option value="Fast Food">Fast Food</option>
                            <option value="Bakery">Bakery</option>
                            <option value="Sweet Shop">Sweet Shop</option>
                        </select>
                    </div>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Shop Description</label>
                    <textarea
                        value={formData.shop_description}
                        onChange={(e) => setFormData(prev => ({ ...prev, shop_description: e.target.value }))}
                        className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                        rows="3"
                        placeholder="Describe your shop..."
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Shop Image</label>
                    <div className="flex items-center gap-2">
                        <input
                            type="file"
                            accept="image/*"
                            onChange={(e) => {
                                const file = e.target.files[0];
                                if (file) handleFileUpload(file, 'shop_img');
                            }}
                            className="hidden"
                            id="shop_img"
                        />
                        <label
                            htmlFor="shop_img"
                            className="flex items-center gap-2 px-3 py-2 border rounded-lg cursor-pointer hover:bg-gray-50 transition-colors"
                        >
                            <Upload size={16} />
                            Upload Shop Image
                        </label>
                        {formData.shop_img && (
                            <span className="text-xs text-green-600 flex items-center gap-1">
                                <Check size={12} />
                                Uploaded
                            </span>
                        )}
                    </div>
                    {uploadProgress.shop_img && (
                        <p className="text-xs text-blue-600 mt-1">{uploadProgress.shop_img}</p>
                    )}
                    {formData.shop_img && (
                        <div className="mt-2">
                            <img
                                src={formData.shop_img}
                                alt="Shop preview"
                                className="w-20 h-20 object-cover rounded-lg border"
                            />
                            <a href={formData.shop_img} target="_blank" rel="noopener noreferrer" className="text-xs text-blue-600 hover:underline block mt-1">
                                View uploaded image
                            </a>
                        </div>
                    )}
                    {/* Fallback URL input */}
                    <div className="mt-2">
                        <label className="block text-xs text-gray-500 mb-1">Or enter image URL directly:</label>
                        <input
                            type="url"
                            value={formData.shop_img}
                            onChange={(e) => setFormData(prev => ({ ...prev, shop_img: e.target.value }))}
                            className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)] text-sm"
                            placeholder="https://example.com/shop-image.jpg"
                        />
                    </div>
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

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Latitude</label>
                        <input
                            type="number"
                            step="any"
                            value={formData.latitude}
                            onChange={(e) => setFormData(prev => ({ ...prev, latitude: parseFloat(e.target.value) }))}
                            className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Longitude</label>
                        <input
                            type="number"
                            step="any"
                            value={formData.longitude}
                            onChange={(e) => setFormData(prev => ({ ...prev, longitude: parseFloat(e.target.value) }))}
                            className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                        />
                    </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
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

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Opening Time</label>
                        <input
                            type="time"
                            value={formData.outlet_open_time}
                            onChange={(e) => setFormData(prev => ({ ...prev, outlet_open_time: e.target.value }))}
                            className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                        />
                    </div>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Closing Time</label>
                    <input
                        type="time"
                        value={formData.outlet_close_time}
                        onChange={(e) => setFormData(prev => ({ ...prev, outlet_close_time: e.target.value }))}
                        className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                    />
                </div>

                {/* Owner Information */}
                <div className="border-t pt-4">
                    <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
                        <User size={20} />
                        Owner Information
                    </h3>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">Owner Name</label>
                            <input
                                type="text"
                                value={formData.owner_name}
                                onChange={(e) => setFormData(prev => ({ ...prev, owner_name: e.target.value }))}
                                className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                                required
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">Owner Phone</label>
                            <input
                                type="tel"
                                value={formData.owner_phone}
                                onChange={(e) => setFormData(prev => ({ ...prev, owner_phone: e.target.value }))}
                                className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                                required
                            />
                        </div>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Owner Email</label>
                        <input
                            type="email"
                            value={formData.owner_email}
                            onChange={(e) => setFormData(prev => ({ ...prev, owner_email: e.target.value }))}
                            className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                            required
                        />
                    </div>
                </div>

                {/* Business Documents */}
                <div className="border-t pt-4">
                    <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
                        <FileText size={20} />
                        Business Documents
                    </h3>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">FSSAI License Number</label>
                            <input
                                type="text"
                                value={formData.fssai_license_number}
                                onChange={(e) => setFormData(prev => ({ ...prev, fssai_license_number: e.target.value }))}
                                className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">GSTIN Number</label>
                            <input
                                type="text"
                                value={formData.gstin_number}
                                onChange={(e) => setFormData(prev => ({ ...prev, gstin_number: e.target.value }))}
                                className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                            />
                        </div>
                    </div>

                    <FileUploadField label="FSSAI License Copy" fieldName="fssai_liscense_copy" />
                    <FileUploadField label="GSTIN Copy" fieldName="gstin_copy" />
                    <FileUploadField label="Pancard Copy" fieldName="pancard_copy" />
                </div>

                {/* Bank Information */}
                <div className="border-t pt-4">
                    <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
                        <Banknote size={20} />
                        Bank Information
                    </h3>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">Bank Name</label>
                            <input
                                type="text"
                                value={formData.bank_name}
                                onChange={(e) => setFormData(prev => ({ ...prev, bank_name: e.target.value }))}
                                className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">Account Number</label>
                            <input
                                type="text"
                                value={formData.bank_account_number}
                                onChange={(e) => setFormData(prev => ({ ...prev, bank_account_number: e.target.value }))}
                                className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                            />
                        </div>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">IFSC Code</label>
                            <input
                                type="text"
                                value={formData.bank_ifsc_code}
                                onChange={(e) => setFormData(prev => ({ ...prev, bank_ifsc_code: e.target.value }))}
                                className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">Account Holder Name</label>
                            <input
                                type="text"
                                value={formData.bank_account_holder_name}
                                onChange={(e) => setFormData(prev => ({ ...prev, bank_account_holder_name: e.target.value }))}
                                className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                            />
                        </div>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Account Type</label>
                        <select
                            value={formData.bank_account_type}
                            onChange={(e) => setFormData(prev => ({ ...prev, bank_account_type: e.target.value }))}
                            className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--primary-color)]"
                        >
                            <option value="Savings">Savings</option>
                            <option value="Current">Current</option>
                        </select>
                    </div>
                </div>

                {/* Food Category and Status */}
                <div className="border-t pt-4">
                    <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
                        <Shield size={20} />
                        Shop Settings
                    </h3>

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

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">Verification Status</label>
                            <div className="flex gap-3">
                                <label className="flex items-center gap-2">
                                    <input
                                        type="radio"
                                        checked={formData.shop_verification_status === 'pending'}
                                        onChange={() => setFormData(prev => ({
                                            ...prev,
                                            shop_verification_status: 'pending',
                                            is_shop_verified: false
                                        }))}
                                        className="text-yellow-600"
                                    />
                                    <span>Pending</span>
                                </label>
                                <label className="flex items-center gap-2">
                                    <input
                                        type="radio"
                                        checked={formData.shop_verification_status === 'verified'}
                                        onChange={() => setFormData(prev => ({
                                            ...prev,
                                            shop_verification_status: 'verified',
                                            is_shop_verified: true
                                        }))}
                                        className="text-green-600"
                                    />
                                    <span>Verified</span>
                                </label>
                                <label className="flex items-center gap-2">
                                    <input
                                        type="radio"
                                        checked={formData.shop_verification_status === 'rejected'}
                                        onChange={() => setFormData(prev => ({
                                            ...prev,
                                            shop_verification_status: 'rejected',
                                            is_shop_verified: false
                                        }))}
                                        className="text-red-600"
                                    />
                                    <span>Rejected</span>
                                </label>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="flex justify-end gap-3 mt-6 pt-4 border-t">
                    <button
                        type="button"
                        onClick={onCancel}
                        className="px-4 py-2 border rounded-lg hover:bg-gray-50 transition-colors"
                        disabled={uploading}
                    >
                        Cancel
                    </button>
                    <button
                        type="submit"
                        className="px-4 py-2 bg-[var(--primary-color)] text-white rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50"
                        disabled={uploading}
                    >
                        {uploading ? 'Uploading...' : (initialData ? 'Save Changes' : 'Add Shop')}
                    </button>
                </div>
            </form>
        </div>
    );
};

const ShopDetailsModal = ({ shop, onClose }) => {
    return (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
                <div className="p-6">
                    <div className="flex justify-between items-start mb-6">
                        <h2 className="text-2xl font-bold text-gray-900">{shop.shop_name}</h2>
                        <button onClick={onClose} className="text-gray-400 hover:text-gray-600">
                            <X size={24} />
                        </button>
                    </div>

                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                        {/* Shop Image */}
                        <div>
                            <img
                                src={shop.shop_img}
                                alt={shop.shop_name}
                                className="w-full h-64 object-cover rounded-lg"
                            />
                        </div>

                        {/* Shop Info */}
                        <div className="space-y-4">
                            <div className="flex items-center justify-between">
                                <span className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${shop.shop_status ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                                    }`}>
                                    {shop.shop_status ? 'Active' : 'Inactive'}
                                </span>
                                <span className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${shop.shop_verification_status === 'verified' ? 'bg-green-100 text-green-800' :
                                    shop.shop_verification_status === 'pending' || shop.shop_verification_status === '' ? 'bg-yellow-100 text-yellow-800' :
                                        'bg-red-100 text-red-800'
                                    }`}>
                                    {shop.shop_verification_status || 'pending'}
                                </span>
                            </div>

                            <div className="space-y-3">
                                <div className="flex items-center gap-2">
                                    <Building size={16} className="text-gray-500" />
                                    <span className="text-sm text-gray-600">{shop.shop_type}</span>
                                </div>
                                <div className="flex items-center gap-2">
                                    <MapPin size={16} className="text-gray-500" />
                                    <span className="text-sm text-gray-600">{shop.address}</span>
                                </div>
                                <div className="flex items-center gap-2">
                                    <Clock size={16} className="text-gray-500" />
                                    <span className="text-sm text-gray-600">{shop.preparation_time}</span>
                                </div>
                                {shop.outlet_open_time && shop.outlet_close_time && (
                                    <div className="flex items-center gap-2">
                                        <Clock size={16} className="text-gray-500" />
                                        <span className="text-sm text-gray-600">
                                            {shop.outlet_open_time} - {shop.outlet_close_time}
                                        </span>
                                    </div>
                                )}
                            </div>

                            {shop.shop_description && (
                                <div>
                                    <h4 className="font-medium text-gray-900 mb-1">Description</h4>
                                    <p className="text-sm text-gray-600">{shop.shop_description}</p>
                                </div>
                            )}

                            {/* Ratings */}
                            <div className="flex items-center gap-2">
                                <Star size={16} className="text-yellow-400 fill-current" />
                                <span className="text-sm font-medium">{shop.average_rating?.toFixed(1) || '0.0'}</span>
                                <span className="text-sm text-gray-500">({shop.total_ratings || 0} ratings)</span>
                            </div>
                        </div>
                    </div>

                    {/* Owner Information */}
                    <div className="mt-6 border-t pt-6">
                        <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
                            <User size={20} />
                            Owner Information
                        </h3>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div className="flex items-center gap-2">
                                <User size={16} className="text-gray-500" />
                                <span className="text-sm text-gray-600">{shop.owner_name}</span>
                            </div>
                            <div className="flex items-center gap-2">
                                <Phone size={16} className="text-gray-500" />
                                <span className="text-sm text-gray-600">{shop.owner_phone}</span>
                            </div>
                            <div className="flex items-center gap-2">
                                <Mail size={16} className="text-gray-500" />
                                <span className="text-sm text-gray-600">{shop.owner_email}</span>
                            </div>
                        </div>
                    </div>

                    {/* Business Documents */}
                    <div className="mt-6 border-t pt-6">
                        <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
                            <FileText size={20} />
                            Business Documents
                        </h3>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700">FSSAI License</label>
                                <p className="text-sm text-gray-600">{shop.fssai_license_number}</p>
                                {shop.fssai_liscense_copy && (
                                    <a href={shop.fssai_liscense_copy} target="_blank" rel="noopener noreferrer" className="text-xs text-blue-600 hover:underline">
                                        View Document
                                    </a>
                                )}
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700">GSTIN</label>
                                <p className="text-sm text-gray-600">{shop.gstin_number}</p>
                                {shop.gstin_copy && (
                                    <a href={shop.gstin_copy} target="_blank" rel="noopener noreferrer" className="text-xs text-blue-600 hover:underline">
                                        View Document
                                    </a>
                                )}
                            </div>
                        </div>
                    </div>

                    {/* Bank Information */}
                    <div className="mt-6 border-t pt-6">
                        <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
                            <Banknote size={20} />
                            Bank Information
                        </h3>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700">Bank Name</label>
                                <p className="text-sm text-gray-600">{shop.bank_name}</p>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700">Account Number</label>
                                <p className="text-sm text-gray-600">{shop.bank_account_number}</p>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700">IFSC Code</label>
                                <p className="text-sm text-gray-600">{shop.bank_ifsc_code}</p>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700">Account Holder</label>
                                <p className="text-sm text-gray-600">{shop.bank_account_holder_name}</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

const ShopsManager = () => {
    const [shops, setShops] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showForm, setShowForm] = useState(false);
    const [showDetails, setShowDetails] = useState(false);
    const [selectedShop, setSelectedShop] = useState(null);
    const [editingShop, setEditingShop] = useState(null);
    const [searchTerm, setSearchTerm] = useState('');

    useEffect(() => {
        fetchShops();
    }, []);

    const fetchShops = async () => {
        try {
            setLoading(true);
            const { data } = await axiosClient.get('/api/shop?admin=true');
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

    const handleViewDetails = (shop) => {
        setSelectedShop(shop);
        setShowDetails(true);
    };

    const handleCloseForm = () => {
        setShowForm(false);
        setEditingShop(null);
    };

    const handleCloseDetails = () => {
        setShowDetails(false);
        setSelectedShop(null);
    };

    const filteredShops = shops.filter(shop =>
        shop.shop_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        shop.address.toLowerCase().includes(searchTerm.toLowerCase()) ||
        shop.owner_name?.toLowerCase().includes(searchTerm.toLowerCase())
    );

    if (loading) {
        return (
            <div className="p-6">
                <div className="animate-pulse">
                    <div className="h-8 bg-gray-200 rounded w-1/4 mb-4"></div>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                        {[1, 2, 3, 4, 5, 6].map(i => (
                            <div key={i} className="bg-white rounded-xl shadow-sm overflow-hidden">
                                <div className="h-48 bg-gray-200"></div>
                                <div className="p-4 space-y-2">
                                    <div className="h-4 bg-gray-200 rounded"></div>
                                    <div className="h-3 bg-gray-200 rounded w-2/3"></div>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="p-6">
            <div className="flex justify-between items-center mb-6">
                <div>
                    <h1 className="text-2xl font-bold text-gray-800">Shops Management</h1>
                    <p className="text-gray-600">Manage your shops and their details</p>
                </div>
                <button
                    onClick={() => setShowForm(true)}
                    className="flex items-center gap-2 bg-[var(--primary-color)] text-white px-4 py-2 rounded-lg hover:opacity-90 transition-opacity"
                >
                    <Plus size={20} />
                    Add Shop
                </button>
            </div>

            {/* Stats Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 hover:shadow-md transition-all duration-300">
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm font-medium text-gray-600 mb-1">Total Shops</p>
                            <p className="text-3xl font-bold text-gray-900">{shops.length}</p>
                        </div>
                        <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-blue-600 rounded-xl flex items-center justify-center">
                            <Store size={24} className="text-white" />
                        </div>
                    </div>
                </div>

                <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 hover:shadow-md transition-all duration-300">
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm font-medium text-gray-600 mb-1">Active Shops</p>
                            <p className="text-3xl font-bold text-green-600">{shops.filter(s => s.shop_status === 1).length}</p>
                        </div>
                        <div className="w-12 h-12 bg-gradient-to-br from-green-500 to-green-600 rounded-xl flex items-center justify-center">
                            <Check size={24} className="text-white" />
                        </div>
                    </div>
                </div>

                <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 hover:shadow-md transition-all duration-300">
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm font-medium text-gray-600 mb-1">Verified Shops</p>
                            <p className="text-3xl font-bold text-blue-600">{shops.filter(s => s.shop_verification_status === 'verified').length}</p>
                        </div>
                        <div className="w-12 h-12 bg-gradient-to-br from-purple-500 to-purple-600 rounded-xl flex items-center justify-center">
                            <Shield size={24} className="text-white" />
                        </div>
                    </div>
                </div>

                <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 hover:shadow-md transition-all duration-300">
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm font-medium text-gray-600 mb-1">Pending Verification</p>
                            <p className="text-3xl font-bold text-yellow-600">{shops.filter(s => s.shop_verification_status === 'pending' || s.shop_verification_status === '').length}</p>
                        </div>
                        <div className="w-12 h-12 bg-gradient-to-br from-yellow-500 to-yellow-600 rounded-xl flex items-center justify-center">
                            <Clock size={24} className="text-white" />
                        </div>
                    </div>
                </div>
            </div>

            {/* Search Section */}
            <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 mb-8">
                <div className="relative">
                    <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400" size={20} />
                    <input
                        type="text"
                        placeholder="Search shops by name, address, or owner..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="w-full pl-12 pr-4 py-4 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition-all duration-300 text-lg"
                    />
                </div>
            </div>

            {showForm && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
                    <div className="bg-white rounded-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
                        <ShopForm
                            onSubmit={handleSubmit}
                            onCancel={handleCloseForm}
                            initialData={editingShop}
                        />
                    </div>
                </div>
            )}

            {showDetails && selectedShop && (
                <ShopDetailsModal
                    shop={selectedShop}
                    onClose={handleCloseDetails}
                />
            )}

            {/* Shops Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {filteredShops.map(shop => (
                    <div key={shop.shop_id} className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden hover:shadow-lg transition-all duration-300 transform hover:scale-105">
                        <div className="relative">
                            <img
                                src={shop.shop_img}
                                alt={shop.shop_name}
                                className="w-full h-56 object-cover"
                            />
                            <div className="absolute top-4 right-4 flex flex-col gap-2">
                                <span className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-bold ${shop.shop_status ? 'bg-green-500 text-white' : 'bg-red-500 text-white'}`}>
                                    {shop.shop_status ? 'Active' : 'Inactive'}
                                </span>
                                <span className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-bold ${shop.shop_verification_status === 'verified' ? 'bg-green-500 text-white' :
                                    shop.shop_verification_status === 'pending' || shop.shop_verification_status === '' ? 'bg-yellow-500 text-white' :
                                        'bg-red-500 text-white'
                                    }`}>
                                    {shop.shop_verification_status || 'pending'}
                                </span>
                            </div>
                        </div>

                        <div className="p-6">
                            <div className="mb-4">
                                <h3 className="text-xl font-bold text-gray-900 mb-2">{shop.shop_name}</h3>
                                <div className="space-y-2 text-sm text-gray-600">
                                    <div className="flex items-center gap-2">
                                        <Building size={16} className="text-gray-400" />
                                        <span>{shop.shop_type}</span>
                                    </div>
                                    <div className="flex items-center gap-2">
                                        <MapPin size={16} className="text-gray-400" />
                                        <span className="truncate">{shop.address}</span>
                                    </div>
                                    <div className="flex items-center gap-2">
                                        <User size={16} className="text-gray-400" />
                                        <span>{shop.owner_name}</span>
                                    </div>
                                    {shop.average_rating > 0 && (
                                        <div className="flex items-center gap-2">
                                            <Star size={16} className="text-yellow-400 fill-current" />
                                            <span className="font-medium">{shop.average_rating.toFixed(1)}</span>
                                            <span className="text-gray-500">({shop.total_ratings} ratings)</span>
                                        </div>
                                    )}
                                </div>
                            </div>

                            <div className="flex gap-3">
                                <button
                                    onClick={() => handleViewDetails(shop)}
                                    className="flex-1 flex items-center justify-center gap-2 px-4 py-3 text-sm font-medium text-blue-600 hover:bg-blue-50 rounded-xl transition-all duration-300 border border-blue-200 hover:border-blue-300"
                                >
                                    <Eye size={16} />
                                    View Details
                                </button>
                                <button
                                    onClick={() => handleEdit(shop)}
                                    className="flex-1 flex items-center justify-center gap-2 px-4 py-3 text-sm font-medium text-purple-600 hover:bg-purple-50 rounded-xl transition-all duration-300 border border-purple-200 hover:border-purple-300"
                                >
                                    <Edit size={16} />
                                    Edit
                                </button>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default ShopsManager;