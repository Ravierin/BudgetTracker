import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { TrendingUp, Zap, Shield, DollarSign } from 'lucide-react';
import { api } from '../api/api';
import { wsService, type WSMessage } from '../api/websocket';

export function Dashboard() {
  const [totalPnl, setTotalPnl] = useState(0);
  const [monthlyPnl, setMonthlyPnl] = useState(0);
  const [loading, setLoading] = useState(true);

  const calculateStats = async () => {
    try {
      const positions = await api.getPositions();
      const incomes = await api.getMonthlyIncomes();
      
      const total = positions.reduce((sum, p) => sum + p.closedPnl, 0);
      
      const currentMonth = new Date().getMonth();
      const currentYear = new Date().getFullYear();
      const monthly = incomes
        .filter(i => {
          const date = new Date(i.date);
          return date.getMonth() === currentMonth && date.getFullYear() === currentYear;
        })
        .reduce((sum, i) => sum + i.pnl, 0);
      
      setTotalPnl(total);
      setMonthlyPnl(monthly);
    } catch (error) {
      console.error('Failed to calculate stats:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    calculateStats();

    const unsubscribe = wsService.addListener((message: WSMessage) => {
      if (['positions_update', 'position_created', 'position_deleted'].includes(message.type)) {
        calculateStats();
      }
    });

    return () => unsubscribe();
  }, []);

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
    }).format(value);
  };

  return (
    <div>
      {/* Page Header */}
      <div className="page-header">
        <motion.h1
          className="page-title"
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.5 }}
        >
          Главная
        </motion.h1>
        <p className="page-subtitle">
          Система управления сделками и синхронизации с биржами
        </p>
      </div>

      {/* Total PnL Card */}
      <div className="stats-grid" style={{ marginBottom: '24px' }}>
        <motion.div
          className="stat-card"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4 }}
        >
          <div className="stat-card-header">
            <span className="stat-card-title">Общий PnL</span>
            <div
              className="stat-card-icon"
              style={{ 
                background: totalPnl >= 0 
                  ? 'linear-gradient(135deg, #10b981 0%, #059669 100%)' 
                  : 'linear-gradient(135deg, #ef4444 0%, #dc2626 100%)' 
              }}
            >
              <DollarSign size={20} color="white" />
            </div>
          </div>
          <div className={`stat-card-value ${totalPnl >= 0 ? 'positive' : 'negative'}`}>
            {loading ? '...' : formatCurrency(totalPnl)}
          </div>
          <div className="stat-card-subtitle">
            По всем сделкам за всё время
          </div>
        </motion.div>

        <motion.div
          className="stat-card"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4, delay: 0.1 }}
        >
          <div className="stat-card-header">
            <span className="stat-card-title">PnL за месяц</span>
            <div
              className="stat-card-icon"
              style={{ 
                background: monthlyPnl >= 0 
                  ? 'linear-gradient(135deg, #10b981 0%, #059669 100%)' 
                  : 'linear-gradient(135deg, #ef4444 0%, #dc2626 100%)' 
              }}
            >
              <TrendingUp size={20} color="white" />
            </div>
          </div>
          <div className={`stat-card-value ${monthlyPnl >= 0 ? 'positive' : 'negative'}`}>
            {loading ? '...' : formatCurrency(monthlyPnl)}
          </div>
          <div className="stat-card-subtitle">
            С начала текущего месяца
          </div>
        </motion.div>
      </div>

      {/* Info Cards */}
      <div className="stats-grid">
        {/* Live Sync Status */}
        <motion.div
          className="stat-card"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4, delay: 0.1 }}
          whileHover={{ scale: 1.02 }}
        >
          <div className="stat-card-header">
            <span className="stat-card-title">Статус синхронизации</span>
            <div
              className="stat-card-icon"
              style={{ background: 'linear-gradient(135deg, #10b981 0%, #059669 100%)' }}
            >
              <Zap size={20} color="white" />
            </div>
          </div>
          <div style={{ fontSize: '20px', fontWeight: '700', color: 'var(--success)', marginBottom: '8px' }}>
            Активна
          </div>
          <div className="stat-card-subtitle">
            Автоматическое обновление каждые 30 секунд
          </div>
        </motion.div>

        {/* Connected Exchanges */}
        <motion.div
          className="stat-card"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4, delay: 0.2 }}
          whileHover={{ scale: 1.02 }}
        >
          <div className="stat-card-header">
            <span className="stat-card-title">Подключенные биржи</span>
            <div
              className="stat-card-icon"
              style={{ background: 'linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)' }}
            >
              <TrendingUp size={20} color="white" />
            </div>
          </div>
          <div style={{ display: 'flex', gap: '8px', flexWrap: 'wrap' }}>
            {['MEXC', 'Bybit', 'Gate', 'Bitget'].map((ex) => (
              <span
                key={ex}
                className="badge badge-info"
                style={{ fontSize: '13px', padding: '6px 12px' }}
              >
                {ex}
              </span>
            ))}
          </div>
          <div className="stat-card-subtitle" style={{ marginTop: '12px' }}>
            Все биржи готовы к синхронизации
          </div>
        </motion.div>

        {/* Security */}
        <motion.div
          className="stat-card"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4, delay: 0.3 }}
          whileHover={{ scale: 1.02 }}
        >
          <div className="stat-card-header">
            <span className="stat-card-title">Безопасность</span>
            <div
              className="stat-card-icon"
              style={{ background: 'linear-gradient(135deg, #f59e0b 0%, #d97706 100%)' }}
            >
              <Shield size={20} color="white" />
            </div>
          </div>
          <div style={{ fontSize: '20px', fontWeight: '700', color: 'var(--text-primary)', marginBottom: '8px' }}>
            Защищено
          </div>
          <div className="stat-card-subtitle">
            API ключи хранятся в зашифрованном виде
          </div>
        </motion.div>
      </div>

      {/* Quick Actions */}
      <motion.div
        className="card"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.4, delay: 0.4 }}
      >
        <div className="card-header">
          <h3 className="card-title">Быстрый доступ</h3>
        </div>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '16px' }}>
          <a href="/positions" style={{ textDecoration: 'none' }}>
            <motion.div
              className="btn btn-primary"
              style={{ width: '100%' }}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <TrendingUp size={18} />
              Сделки
            </motion.div>
          </a>
          <a href="/monthly-income" style={{ textDecoration: 'none' }}>
            <motion.div
              className="btn btn-secondary"
              style={{ width: '100%' }}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <TrendingUp size={18} />
              Доход
            </motion.div>
          </a>
          <a href="/withdrawals" style={{ textDecoration: 'none' }}>
            <motion.div
              className="btn btn-secondary"
              style={{ width: '100%' }}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <Zap size={18} />
              Выводы
            </motion.div>
          </a>
          <a href="/settings" style={{ textDecoration: 'none' }}>
            <motion.div
              className="btn btn-secondary"
              style={{ width: '100%' }}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <Shield size={18} />
              Настройки
            </motion.div>
          </a>
        </div>
      </motion.div>
    </div>
  );
}
