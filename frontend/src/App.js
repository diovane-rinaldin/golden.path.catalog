import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import Navigation from './components/Navigation';
import Dashboard from './components/Dashboard';
import TechnologyForm from './components/forms/TechnologyForm';
import ComponentForm from './components/forms/ComponentForm';
import ServiceForm from './components/forms/ServiceForm';
import ComponentList from './components/ComponentList';
import ServiceList from './components/ServiceList';

function App() {
  return (
    <AuthProvider>
      <Router>
        <div className="min-h-screen bg-gray-100">
          <Navigation />
          <main className="container mx-auto px-4 py-8">
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/technology/new" element={<TechnologyForm />} />
              <Route path="/component/new" element={<ComponentForm />} />
              <Route path="/service/new" element={<ServiceForm />} />
              <Route path="/technology/:id/components" element={<ComponentList />} />
              <Route path="/component/:id/services" element={<ServiceList />} />
            </Routes>
          </main>
        </div>
      </Router>
    </AuthProvider>
  );
}

export default App;