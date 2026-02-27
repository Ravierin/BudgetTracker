import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { TrendingUp, DollarSign, Wallet } from 'lucide-react';
import { api } from '../api/api';
import { wsService, type WSMessage } from '../api/websocket';

export function Dashboard() {
  const [totalPnl, setTotalPnl] = useState(0);
  const [monthlyPnl, setMonthlyPnl] = useState(0);
  const [totalBalance, setTotalBalance] = useState(0);
  const [loading, setLoading] = useState(true);

  const calculateStats = async () => {
    try {
      const positions = await api.getPositions();
      const incomes = await api.getMonthlyIncomes();
      const balanceData = await api.getBalance();

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
      setTotalBalance(balanceData.totalBalance);
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
      </div>

      {/* Stats Cards */}
      <div className="stats-grid">
        {/* Total Balance */}
        <motion.div
          className="stat-card"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4 }}
        >
          <div className="stat-card-header">
            <span className="stat-card-title">Общий баланс</span>
            <div
              className="stat-card-icon"
              style={{ background: 'linear-gradient(135deg, #404040 0%, #2a2a2a 100%)' }}
            >
              <Wallet size={20} color="white" />
            </div>
          </div>
          <div className="stat-card-value">
            {loading ? '...' : formatCurrency(totalBalance)}
          </div>
          <div className="stat-card-subtitle">
            Общий баланс со всех бирж
          </div>
        </motion.div>

        {/* Total PnL */}
        <motion.div
          className="stat-card"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4, delay: 0.1 }}
        >
          <div className="stat-card-header">
            <span className="stat-card-title">Общий PnL</span>
            <div
              className="stat-card-icon"
              style={{
                background: totalPnl >= 0
                  ? 'linear-gradient(135deg, #15803d 0%, #14532d 100%)'
                  : 'linear-gradient(135deg, #b91c1c 0%, #991b1b 100%)'
              }}
            >
              <DollarSign size={20} color="white" />
            </div>
          </div>
          <div className={`stat-card-value ${totalPnl >= 0 ? 'positive' : 'negative'}`}>
            {loading ? '...' : (totalPnl >= 0 ? '+' : '') + totalPnl.toFixed(2)}
          </div>
          <div className="stat-card-subtitle">
            По всем сделкам
          </div>
        </motion.div>

        {/* Monthly PnL */}
        <motion.div
          className="stat-card"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4, delay: 0.2 }}
        >
          <div className="stat-card-header">
            <span className="stat-card-title">PnL за месяц</span>
            <div
              className="stat-card-icon"
              style={{
                background: monthlyPnl >= 0
                  ? 'linear-gradient(135deg, #15803d 0%, #14532d 100%)'
                  : 'linear-gradient(135deg, #b91c1c 0%, #991b1b 100%)'
              }}
            >
              <TrendingUp size={20} color="white" />
            </div>
          </div>
          <div className={`stat-card-value ${monthlyPnl >= 0 ? 'positive' : 'negative'}`}>
            {loading ? '...' : (monthlyPnl >= 0 ? '+' : '') + monthlyPnl.toFixed(2)}
          </div>
          <div className="stat-card-subtitle">
            Текущий месяц
          </div>
        </motion.div>
      </div>
    </div>
  );
}
