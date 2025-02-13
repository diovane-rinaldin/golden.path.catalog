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
        <div className="app-container">
          {/* Menu de Navegação Fixo */}
          <Navigation />
          
          {/* Área de Conteúdo Central */}
          <main className="content-container">
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