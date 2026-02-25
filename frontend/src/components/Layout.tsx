import { Link, useLocation } from 'react-router-dom';
import { motion } from 'framer-motion';
import {
  Home,
  ArrowLeftRight,
  TrendingUp,
  Wallet,
  Settings,
  Zap,
  Bitcoin,
  Boxes,
  DoorOpen
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

const exchanges = [
  { id: 'mexc', name: 'MEXC', icon: Zap, color: '#00C076' },
  { id: 'bybit', name: 'Bybit', icon: Boxes, color: '#F7A600' },
  { id: 'gate', name: 'Gate', icon: DoorOpen, color: '#F0B90B' },
  { id: 'bitget', name: 'Bitget', icon: Bitcoin, color: '#00D9FF' },
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
            whileHover={{ scale: 1.1, rotate: 5 }}
            transition={{ type: 'spring', stiffness: 300 }}
          >
            ₿
          </motion.div>
          <div>
            <div className="sidebar-title">Trading Manager</div>
            <div className="sidebar-subtitle">Crypto Portfolio</div>
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

        {/* Exchange Status */}
        <div className="exchange-badges">
          {exchanges.map((ex) => {
            const Icon = ex.icon;
            return (
              <motion.div
                key={ex.id}
                className="exchange-badge"
                whileHover={{ scale: 1.05 }}
                title={`${ex.name} - Connected`}
              >
                <Icon size={14} style={{ color: ex.color }} />
                <span>{ex.name}</span>
              </motion.div>
            );
          })}
        </div>

        {/* WebSocket Status */}
        <motion.div
          className="exchange-badge active"
          style={{ marginTop: '12px' }}
          animate={{ opacity: [1, 0.5, 1] }}
          transition={{ duration: 2, repeat: Infinity }}
        >
          <Zap size={14} />
          <span>Live Sync</span>
        </motion.div>
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
