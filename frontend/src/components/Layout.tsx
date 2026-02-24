import { Link, useLocation } from 'react-router-dom';
import type { ReactNode } from 'react';

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const location = useLocation();

  const navItems = [
    { path: '/', label: 'Дашборд', icon: 'bi-bar-chart' },
    { path: '/positions', label: 'Сделки', icon: 'bi-arrow-left-right' },
    { path: '/monthly-income', label: 'Месячный доход', icon: 'bi-graph-up' },
    { path: '/withdrawals', label: 'P2P', icon: 'bi-download' },
    { path: '/settings', label: 'Настройки', icon: 'bi-gear' },
  ];

  return (
    <div className="d-flex flex-column min-vh-100">
      {/* Navigation */}
      <nav className="navbar navbar-expand-lg">
        <div className="container-fluid">
          <Link to="/" className="navbar-brand text-decoration-none">
            <svg xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 16 16">
              <path d="M14 14V4.06L7 11.06 5.5 9.56 0 15.12V14H14zM15 0H1C.45 0 0 .45 0 1v14c0 .55.45 1 1 1h14c.55 0 1-.45 1-1V1c0-.55-.45-1-1-1z"/>
            </svg>
            <div>
              <div>Трейдинг Менеджер</div>
              <div className="navbar-subtitle">Управление сделками MEXC</div>
            </div>
          </Link>

          <button 
            className="navbar-toggler" 
            type="button" 
            data-bs-toggle="collapse" 
            data-bs-target="#navbarNav"
          >
            <span className="navbar-toggler-icon"></span>
          </button>

          <div className="collapse navbar-collapse justify-content-center" id="navbarNav">
            <ul className="nav nav-pills">
              {navItems.map((item) => (
                <li className="nav-item" key={item.path}>
                  <Link
                    to={item.path}
                    className={`nav-link ${location.pathname === item.path ? 'active' : ''}`}
                  >
                    <i className={`bi ${item.icon}`}></i>
                    {item.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div className="d-flex">
            <button className="btn btn-outline-secondary btn-icon">
              <i className="bi bi-gear"></i>
            </button>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="main-content flex-grow-1">
        {children}
      </main>

      {/* Footer */}
      <footer className="footer">
        <div className="d-flex justify-content-between align-items-center">
          <span className="footer-text">
            © 2026 Трейдинг Менеджер. Веб-приложение для управления сделками.
          </span>
          <span className="footer-text">MEXC API Integration</span>
        </div>
      </footer>
    </div>
  );
}
