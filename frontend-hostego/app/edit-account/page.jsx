"use client";

import React, { useState, useEffect } from 'react';
import BackNavigationButton from '../components/BackNavigationButton';
import { Edit2 } from 'lucide-react';
import axiosClient from '../utils/axiosClient';
import { useSelector, useDispatch } from 'react-redux';
import { setFetchUserAccount } from '../lib/redux/features/user/userSlice';

const EditField = ({ label, value, onSave, type = "text" }) => {
    const [isEditing, setIsEditing] = useState(false);
    const [currentValue, setCurrentValue] = useState(value);
    const [originalValue, setOriginalValue] = useState(value);

    const handleEdit = () => {
        setIsEditing(true);
        setCurrentValue(value);
        setOriginalValue(value);
    };

    const handleCancel = () => {
        setIsEditing(false);
        setCurrentValue(originalValue);
    };

    const handleSave = async () => {
        await onSave(currentValue);
        setIsEditing(false);
        setOriginalValue(currentValue);
    };

    const isValueChanged = currentValue !== originalValue;

    return (
        <div className="bg-white rounded-xl p-4 shadow-sm">
            <div className="flex justify-between items-center mb-1">
                <label className="text-sm text-gray-500 text-[var(--primary-color)]">{label}</label>
                {!isEditing && (
                    <button
                        onClick={handleEdit}
                        className="flex items-center gap-1 text-sm text-gray-400 hover:text-[var(--primary-color)] transition-colors"
                    >
                        <Edit2 color='var(--primary-color)' size={16} />
                        <span className='text-[var(--primary-color)]'>Edit</span>
                    </button>
                )}
            </div>
            <div className="relative">
                <input
                    type={type}
                    value={currentValue}
                    onChange={(e) => setCurrentValue(e.target.value)}
                    disabled={!isEditing}
                    className={`w-full py-2 px-3 rounded-lg transition-all ${isEditing
                        ? 'border-2 border-[var(--primary-color)] bg-white outline-none'
                        : 'border border-gray-100 bg-gray-50'
                        }`}
                />

                {isEditing && (
                    <div className="flex gap-2 mt-3">
                        <button
                            onClick={handleSave}
                            disabled={!isValueChanged}
                            className={`flex-1 py-2 px-4 rounded-lg text-sm font-medium transition-all ${isValueChanged
                                ? 'bg-[var(--primary-color)] text-white hover:opacity-90'
                                : 'bg-gray-100 text-gray-400 cursor-not-allowed'
                                }`}
                        >
                            Update
                        </button>
                        <button
                            onClick={handleCancel}
                            className="flex-1 py-2 px-4 rounded-lg text-sm font-medium border-2 border-gray-200 
                                     text-gray-700 hover:bg-gray-50 transition-all"
                        >
                            Cancel
                        </button>
                    </div>
                )}
            </div>
        </div>
    );
};

const Page = () => {
    const { userAccount } = useSelector((state) => state.user);
    const [updateStatus, setUpdateStatus] = useState('');
    const dispatch = useDispatch()

    const handleUpdateField = async (field, value) => {
        try {
            setUpdateStatus('Updating...');

            await axiosClient.patch(`/api/users/me`, { [field]: value });
            dispatch(setFetchUserAccount(true))
            setUpdateStatus('Updated successfully!');
            setTimeout(() => setUpdateStatus(''), 2000);
        } catch (error) {
            console.error('Update failed:', error);
            setUpdateStatus('Update failed');
            setTimeout(() => setUpdateStatus(''), 2000);
        }
    };

    return (
        <div className="min-h-screen bg-[var(--bg-page-color)]">
            <BackNavigationButton title="Edit Account" />

            <div className="p-4 space-y-4">
                {/* Profile Header */}
                <div className="bg-gradient-to-r from-[var(--primary-color)] to-[#7c3aed] rounded-xl p-6 text-white mb-6">
                    <h2 className="text-xl font-semibold mb-2">Edit Profile</h2>
                    <p className="text-sm opacity-90">Update your personal information</p>
                </div>

                {/* Form Fields */}
                <EditField
                    label="First Name"
                    value={userAccount?.first_name || ''}
                    onSave={(value) => handleUpdateField('first_name', value)}
                />

                <EditField
                    label="Last Name"
                    value={userAccount?.last_name || ''}
                    onSave={(value) => handleUpdateField('last_name', value)}
                />

                <EditField
                    label="Email Address"
                    value={userAccount?.email || ''}
                    type="email"
                    onSave={(value) => handleUpdateField('email', value)}
                />

                {/* Status Message */}
                {updateStatus && (
                    <div className={`fixed bottom-4 left-4 right-4 p-4 rounded-lg text-center text-white 
                        shadow-lg ${updateStatus.includes('failed') ? 'bg-red-500' : 'bg-[var(--primary-color)]'} 
                        animate-fade-in`}>
                        {updateStatus}
                    </div>
                )}
            </div>
        </div>
    );
};

export default Page;
