"use client";

import React, { useState, useEffect } from 'react';
import { Search, Filter, MoreVertical, Phone, Mail, Calendar, Clock, CheckCircle, XCircle, ChevronDown } from 'lucide-react';
import axiosClient from '../../utils/axiosClient';

const UserCard = ({ userData, onRoleChange }) => {
    const [isExpanded, setIsExpanded] = useState(false);
    const [showRoleMenu, setShowRoleMenu] = useState(false);

    const user = userData.user; // Get user data
    const userRoles = userData.roles || []; // Get roles array

    const formattedDate = new Date(user.created_at).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    });

    const lastLoginDate = new Date(user.last_login_timestamp).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
  
    const roles = {
        1: { id: 1, name: "Super Admin", class: "bg-red-100 text-red-700" },
        2: { id: 2, name: "Admin", class: "bg-pink-100 text-pink-700" },
        3: { id: 3, name: "User", class: "bg-gray-100 text-gray-700" },
        4: { id: 4, name: "Delivery Partner Manager", class: "bg-purple-100 text-purple-700" },
        5: { id: 5, name: "Payments Manager", class: "bg-blue-100 text-blue-700" },
        6: { id: 6, name: "Order Assign Manager", class: "bg-green-100 text-green-700" },
        7: { id: 7, name: "Delivery Partner", class: "bg-indigo-100 text-indigo-700" },
        8: { id: 8, name: "Order Manager", class: "bg-yellow-100 text-yellow-700" },
        9: { id: 9, name: "Customer Support", class: "bg-orange-100 text-orange-700" }
    };

    // Get array of role IDs that user currently has
    const userRoleIds = userRoles?.map(role => role?.role?.role_id);

    const handleRoleToggle = async (roleId) => {
    
        const currentRoleItemId = userRoles?.find(role => role?.role_id === roleId)?.user_role_id;
     
        try {
            await onRoleChange(user?.user_id, roleId, !userRoleIds.includes(roleId), currentRoleItemId);
        } catch (error) {
            console.error('Error toggling role:', error);
        }
    };


    return (
        <div className="bg-white rounded-xl shadow-sm overflow-hidden">
            {/* Main Info */}
            <div className="p-4">
                <div className="flex items-start justify-between">
                    <div className="flex gap-3">
                        <div className="w-12 h-12 rounded-full bg-gradient-to-br from-[var(--primary-color)] to-[#7c3aed] flex items-center justify-center text-white font-medium">
                            {user.first_name?.[0]}{user.last_name?.[0]}
                        </div>
                        <div>
                            <h3 className="font-medium text-gray-900">
                                {user.first_name} {user.last_name}
                            </h3>
                            <div className="flex items-center gap-2 text-sm text-gray-500 mt-1">
                                <Phone size={14} />
                                <span>{user.mobile_number}</span>
                            </div>
                        </div>
                    </div>
                    <button
                        onClick={() => setIsExpanded(!isExpanded)}
                        className="p-2 hover:bg-gray-100 rounded-full transition-colors"
                    >
                        <ChevronDown className={`w-5 h-5 transition-transform ${isExpanded ? 'rotate-180' : ''}`} />
                    </button>
                </div>

                {/* Verification Status */}
                <div className="mt-3 flex items-center gap-2">
                    {user.firebase_otp_verified ? (
                        <div className="flex items-center gap-1.5 text-sm text-green-600 bg-green-50 px-2.5 py-1 rounded-full">
                            <CheckCircle size={14} />
                            <span>Verified</span>
                        </div>
                    ) : (
                        <div className="flex items-center gap-1.5 text-sm text-red-600 bg-red-50 px-2.5 py-1 rounded-full">
                            <XCircle size={14} />
                            <span>Unverified</span>
                        </div>
                    )}
                </div>

                {/* Roles Section */}
                <div className="mt-4">
                    <div className="flex items-center justify-between mb-2">
                        <span className="text-sm text-gray-500">Roles</span>
                        <button
                            onClick={() => setShowRoleMenu(!showRoleMenu)}
                            className="text-sm text-[var(--primary-color)] hover:underline"
                        >
                            Manage Roles
                        </button>
                    </div>

                    {/* Current Roles */}
                    <div className="flex flex-wrap gap-2">
                        {userRoles.length > 0 ? (
                            userRoles.map(roleData => (
                                <span
                                    key={roleData.user_role_id}
                                    className={`px-2 py-1 rounded-full text-xs font-medium ${roles[roleData.role.role_id].class}`}
                                >
                                    {roles[roleData.role.role_id].name}
                                </span>
                            ))
                        ) : (
                            <span className="text-sm text-gray-400">No roles assigned</span>
                        )}
                    </div>

                    {/* Role Management Menu */}
                    {showRoleMenu && (
                        <div className="mt-4 p-4 border rounded-lg bg-gray-50 animate-fade-in">
                            <h4 className="text-sm font-medium mb-3">Assign/Remove Roles</h4>
                            <div className="space-y-2">
                                {Object.values(roles)?.map(role => (
                                    <div
                                        key={role?.id}
                                        className="flex items-center justify-between p-2 hover:bg-gray-100 rounded-lg transition-colors"
                                    >
                                        <span className="text-sm">{role?.name}</span>
                                        <button
                                            onClick={() => handleRoleToggle(role?.id)}
                                            className={`px-3 py-1 rounded-full text-xs font-medium transition-colors
                                                ${userRoleIds.includes(role?.id)
                                                    ? 'bg-[var(--primary-color)] text-white'
                                                    : 'bg-gray-200 text-gray-600 hover:bg-gray-300'
                                                }`}
                                        >
                                            {userRoleIds.includes(role?.id) ? 'Remove' : 'Add'}
                                        </button>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}
                </div>
            </div>

            {/* Expanded Details */}
            {isExpanded && (
                <div className="px-4 pb-4 border-t pt-4 bg-gray-50 animate-fade-in">
                    <div className="space-y-3">
                        <div className="flex items-center gap-2 text-sm">
                            <Mail className="w-4 h-4 text-gray-400" />
                            <span className="text-gray-600">{user?.email}</span>
                        </div>
                        <div className="flex items-center gap-2 text-sm">
                            <Calendar className="w-4 h-4 text-gray-400" />
                            <span className="text-gray-600">Joined {formattedDate}</span>
                        </div>
                        <div className="flex items-center gap-2 text-sm">
                            <Clock className="w-4 h-4 text-gray-400" />
                            <span className="text-gray-600">Last login: {lastLoginDate}</span>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

const UserManager = () => {

    const [users, setUsers] = useState([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');
    const [filterVerified, setFilterVerified] = useState('all'); // 'all', 'verified', 'unverified'

    useEffect(() => {
        fetchUsers();
    }, []);

    const fetchUsers = async () => {
        try {
            setLoading(true);
            const { data } = await axiosClient.get('/api/users/');
            setUsers(data);
            console.log(data, "data")
        } catch (error) {
            console.error('Error fetching users:', error);
        } finally {
            setLoading(false);
        }
    };


    const handleRoleChange = async (userId, roleId, isAdding, userRoleId) => {

        try {
            if (isAdding && roleId === 1) {
                toast.error("Super Admin role cannot be added or removed");
                return;
            }
            const endpoint = isAdding ? `/api/user-roles/add` : `/api/user-roles/${userRoleId}`;
            if (isAdding) {
                await axiosClient.post(endpoint, {
                    user_id: userId,
                    role_id: roleId
                });
            } else {
                await axiosClient.delete(endpoint);
            }

            // Refresh users list
            await fetchUsers();
        } catch (error) {
            console.error('Error updating role:', error);
        }
    };

    const filteredUsers = users.filter(userData => {
        const user = userData.user;
        const matchesSearch = (
            user.first_name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
            user.last_name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
            user.email?.toLowerCase().includes(searchTerm.toLowerCase()) ||
            user.mobile_number?.includes(searchTerm)
        );

        const matchesFilter = filterVerified === 'all' ? true :
            filterVerified === 'verified' ? user.firebase_otp_verified === 1 :
                user.firebase_otp_verified === 0;

        return matchesSearch && matchesFilter;
    });

    return (
        <div className="p-4">
            {/* Header */}
            <div className="mb-6">
                <h1 className="text-2xl font-bold text-gray-800">User Management</h1>
                <p className="text-gray-600">Manage and monitor user accounts</p>
            </div>
            {/* Stats Summary */}
            <div className="mt-6 grid grid-cols-2 sm:grid-cols-4 gap-4 mb-2">
                <div className="bg-white p-4 rounded-xl shadow-sm">
                    <h3 className="text-sm text-gray-500">Total Users</h3>
                    <p className="text-2xl font-semibold">{users.length}</p>
                </div>
                <div className="bg-white p-4 rounded-xl shadow-sm">
                    <h3 className="text-sm text-gray-500">Verified Users</h3>
                    <p className="text-2xl font-semibold text-green-600">
                        {users?.filter(u => u?.user?.firebase_otp_verified === 1).length}
                    </p>
                </div>
                <div className="bg-white p-4 rounded-xl shadow-sm">
                    <h3 className="text-sm text-gray-500">Unverified Users</h3>
                    <p className="text-2xl font-semibold text-red-600">
                        {users?.filter(u => u?.user?.firebase_otp_verified === 0).length}
                    </p>
                </div>
                <div className="bg-white p-4 rounded-xl shadow-sm">
                    <h3 className="text-sm text-gray-500">New Today</h3>
                    <p className="text-2xl font-semibold text-[var(--primary-color)]">
                        {users?.filter(u => new Date(u?.user?.created_at).toDateString() === new Date()?.toDateString()).length}
                    </p>
                </div>
            </div>
            {/* Search and Filter Bar */}
            <div className="bg-white rounded-xl p-4 shadow-sm mb-6">
                <div className="flex flex-col sm:flex-row gap-4">
                    {/* Search Input */}
                    <div className="flex-1 relative">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                        <input
                            type="text"
                            placeholder="Search users..."
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                            className="w-full pl-10 pr-4 py-2 border rounded-lg focus:outline-none focus:border-[var(--primary-color)]"
                        />
                    </div>

                    {/* Filter Dropdown */}
                    <div className="flex-shrink-0">
                        <select
                            value={filterVerified}
                            onChange={(e) => setFilterVerified(e.target.value)}
                            className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:border-[var(--primary-color)] appearance-none bg-white"
                        >
                            <option value="all">All Users</option>
                            <option value="verified">Verified Only</option>
                            <option value="unverified">Unverified Only</option>
                        </select>
                    </div>
                </div>
            </div>

            {/* User List */}
            <div className="space-y-4">
                {loading ? (
                    // Loading Skeleton
                    [...Array(3)]?.map((_, i) => (
                        <div key={i} className="bg-white rounded-xl p-4 shadow-sm animate-pulse">
                            <div className="flex items-center gap-4">
                                <div className="w-12 h-12 rounded-full bg-gray-200" />
                                <div className="flex-1">
                                    <div className="h-4 bg-gray-200 rounded w-1/4 mb-2" />
                                    <div className="h-3 bg-gray-200 rounded w-1/3" />
                                </div>
                            </div>
                        </div>
                    ))
                ) : filteredUsers?.length > 0 ? (
                    filteredUsers?.map(userData => (
                        <UserCard
                            key={userData?.user?.user_id}
                            userData={userData}
                            onRoleChange={handleRoleChange}
                        />
                    ))
                ) : (
                    <div className="text-center py-8 text-gray-500">
                        No users found matching your criteria
                    </div>
                )}
            </div>


        </div>
    );
};

export default UserManager;
