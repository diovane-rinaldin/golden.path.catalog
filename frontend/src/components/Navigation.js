import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

export default function Navigation() {
  const location = useLocation();
  const { logout } = useAuth();

  return (
    <nav className="bg-white shadow-lg">
      <div className="max-w-7xl mx-auto px-4">
        <div className="flex justify-between h-16">
          <div className="flex">
            <Link to="/" className="flex-shrink-0 flex items-center">
              <span className="text-xl font-bold text-gray-800">Golden Path Portal</span>
            </Link>
          </div>

          <div className="flex items-center">
            <Link
              to="/technology/new"
              className={`px-3 py-2 rounded-md text-sm font-medium ${
                location.pathname === '/technology/new'
                  ? 'bg-gray-900 text-white'
                  : 'text-gray-700 hover:bg-gray-100'
              }`}
            >
              New Technology
            </Link>

            <Link
              to="/component/new"
              className={`ml-4 px-3 py-2 rounded-md text-sm font-medium ${
                location.pathname === '/component/new'
                  ? 'bg-gray-900 text-white'
                  : 'text-gray-700 hover:bg-gray-100'
              }`}
            >
              New Component
            </Link>

            <Link
              to="/service/new"
              className={`ml-4 px-3 py-2 rounded-md text-sm font-medium ${
                location.pathname === '/service/new'
                  ? 'bg-gray-900 text-white'
                  : 'text-gray-700 hover:bg-gray-100'
              }`}
            >
              New Service
            </Link>

            <button
              onClick={logout}
              className="ml-4 px-3 py-2 rounded-md text-sm font-medium text-red-600 hover:bg-red-100"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
}