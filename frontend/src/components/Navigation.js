import React, { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import '../custom.css';  // Importando o CSS

export default function Navigation() {
  const location = useLocation();
  const { logout } = useAuth();

  // Estados para controlar o menu sanfona
  const [isTecnologiasOpen, setIsTecnologiasOpen] = useState(false);
  const [isComponentesOpen, setIsComponentesOpen] = useState(false);
  const [isServicosOpen, setIsServicosOpen] = useState(false);

  // Função para alternar a visibilidade das seções do menu
  const toggleMenu = (menu) => {
    if (menu === 'tecnologias') {
      setIsTecnologiasOpen(!isTecnologiasOpen);
    } else if (menu === 'componentes') {
      setIsComponentesOpen(!isComponentesOpen);
    } else if (menu === 'servicos') {
      setIsServicosOpen(!isServicosOpen);
    }
  };

  return (
    <div className="nav-container">
      {/* Menu fixo à esquerda */}
      <nav className="nav-menu">
        <div className="menu-title">
          <Link to="/" className="text-2xl font-bold">
            Golden Path Portal
          </Link>
        </div>
        
        {/* Menu Accordion */}
        <ul>
          <li>
            <button onClick={() => toggleMenu('tecnologias')} className="menu-button">
              Tecnologias
            </button>
            {isTecnologiasOpen && (
              <ul className="submenu">
                <li>
                  <Link to="/" className={`submenu-item ${location.pathname === '/' ? 'active' : ''}`}>
                    Lista
                  </Link>
                </li>
                <li>
                  <Link to="/technology" className={`submenu-item ${location.pathname === '/technology' ? 'active' : ''}`}>
                    Nova
                  </Link>
                </li>
              </ul>
            )}
          </li>

          <li>
            <button onClick={() => toggleMenu('componentes')} className="menu-button">
              Componentes
            </button>
            {isComponentesOpen && (
              <ul className="submenu">
                <li>
                  <Link to="/component/new" className={`submenu-item ${location.pathname === '/component/new' ? 'active' : ''}`}>
                    New Component
                  </Link>
                </li>
                <li>
                  <Link to="/components" className={`submenu-item ${location.pathname === '/components' ? 'active' : ''}`}>
                    View Components
                  </Link>
                </li>
              </ul>
            )}
          </li>

          <li>
            <button onClick={() => toggleMenu('servicos')} className="menu-button">
              Serviços
            </button>
            {isServicosOpen && (
              <ul className="submenu">
                <li>
                  <Link to="/service/new" className={`submenu-item ${location.pathname === '/service/new' ? 'active' : ''}`}>
                    New Service
                  </Link>
                </li>
                <li>
                  <Link to="/services" className={`submenu-item ${location.pathname === '/services' ? 'active' : ''}`}>
                    View Services
                  </Link>
                </li>
              </ul>
            )}
          </li>

          {/* Logout Button */}
          <li>
            <button onClick={logout} className="menu-button logout-button">
              Logout
            </button>
          </li>
        </ul>
      </nav>

      {/* Conteúdo principal */}
      <div className="content-container">
        {/* Conteúdo principal do seu app aqui */}
      </div>
    </div>
  );
}