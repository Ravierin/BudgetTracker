import { Link, useLocation } from 'react-router-dom';
import { motion } from 'framer-motion';
import {
  Home,
  ArrowLeftRight,
  TrendingUp,
  Wallet,
  Settings
} from 'lucide-react';
import type { ReactNode } from 'react';
import '../App.css';

interface LayoutProps {
  children: ReactNode;
}

const navItems = [
  { path: '/', label: 'Главная', icon: Home },
  { path: '/positions', label: 'Сделки', icon: ArrowLeftRight },
  { path: '/monthly-income', label: 'Доход', icon: TrendingUp },
  { path: '/withdrawals', label: 'Выводы', icon: Wallet },
  { path: '/settings', label: 'Настройки', icon: Settings },
];

export function Layout({ children }: LayoutProps) {
  const location = useLocation();

  return (
    <div className="app-container">
      {/* Sidebar */}
      <motion.aside
        className="sidebar"
        initial={{ x: -300 }}
        animate={{ x: 0 }}
        transition={{ duration: 0.3, ease: 'easeOut' }}
      >
        {/* Logo */}
        <div className="sidebar-header">
          <motion.div
            className="sidebar-logo"
            whileHover={{ scale: 1.05 }}
            transition={{ type: 'spring', stiffness: 300 }}
          >
            ₿
          </motion.div>
          <div>
            <div className="sidebar-title">Trading Manager</div>
          </div>
        </div>

        {/* Navigation */}
        <nav className="nav-menu">
          {navItems.map((item) => {
            const Icon = item.icon;
            const isActive = location.pathname === item.path;

            return (
              <Link key={item.path} to={item.path}>
                <motion.div
                  className={`nav-item ${isActive ? 'active' : ''}`}
                  whileHover={{ scale: 1.02, x: 4 }}
                  whileTap={{ scale: 0.98 }}
                >
                  <Icon className="nav-item-icon" />
                  <span>{item.label}</span>
                </motion.div>
              </Link>
            );
          })}
        </nav>
      </motion.aside>

      {/* Main Content */}
      <main className="main-content">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4, delay: 0.1 }}
        >
          {children}
        </motion.div>
      </main>
    </div>
  );
}
