import { useState, useEffect } from 'react';
import { api } from '../api/api';
import { wsService, type WSMessage } from '../api/websocket';

export function Dashboard() {
  const [stats, setStats] = useState({
    totalBalance: 0,
    totalIncome: 0,
    monthlyPnl: 0,
    totalTrades: 0,
  });
  const [loading, setLoading] = useState(true);

  const calculateStats = async () => {
    try {
      const [positions, , incomes] = await Promise.all([
        api.getPositions(),
        api.getWithdrawals(),
        api.getMonthlyIncomes(),
      ]);

      const totalBalance = positions.reduce((sum, p) => sum + p.cumExitValue, 0);
      const totalIncome = positions.reduce((sum, p) => sum + p.closedPnl, 0);

      const currentMonth = new Date().getMonth();
      const currentYear = new Date().getFullYear();
      const monthlyPnl = incomes
        .filter(i => {
          const date = new Date(i.date);
          return date.getMonth() === currentMonth && date.getFullYear() === currentYear;
        })
        .reduce((sum, i) => sum + i.pnl, 0);

      setStats({
        totalBalance,
        totalIncome,
        monthlyPnl,
        totalTrades: positions.length,
      });
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

  const formatNumber = (value: number) => {
    return new Intl.NumberFormat('ru-RU').format(value);
  };

  if (loading) {
    return (
      <div className="d-flex justify-content-center align-items-center" style={{ height: '50vh' }}>
        <div className="spinner-border text-primary" role="status">
          <span className="visually-hidden">Загрузка...</span>
        </div>
      </div>
    );
  }

  return (
    <div>
      <h1 className="page-title">Общая статистика</h1>
      <p className="page-subtitle">Обзор вашей торговой активности</p>

      <div className="stats-grid">
        {/* Total Balance */}
        <div className="stat-card">
          <div className="stat-card-header">
            <span className="stat-card-title">Общий баланс</span>
            <i className="bi bi-wallet2 stat-card-icon"></i>
          </div>
          <div className="stat-card-value">
            {formatCurrency(stats.totalBalance)}
          </div>
          <div className="stat-card-subtitle">Текущий баланс на бирже</div>
        </div>

        {/* Total Income */}
        <div className="stat-card">
          <div className="stat-card-header">
            <span className="stat-card-title">Общий доход</span>
            <i className="bi bi-currency-dollar stat-card-icon"></i>
          </div>
          <div className={`stat-card-value ${stats.totalIncome >= 0 ? 'positive' : 'negative'}`}>
            {formatCurrency(stats.totalIncome)}
          </div>
          <div className="stat-card-subtitle">Весь период</div>
        </div>

        {/* Monthly PnL */}
        <div className="stat-card">
          <div className="stat-card-header">
            <span className="stat-card-title">Месячный PnL</span>
            <i className="bi bi-graph-up stat-card-icon"></i>
          </div>
          <div className={`stat-card-value ${stats.monthlyPnl >= 0 ? 'positive' : 'negative'}`}>
            {formatCurrency(stats.monthlyPnl)}
          </div>
          <div className="stat-card-subtitle">Текущий месяц</div>
        </div>

        {/* Total Trades */}
        <div className="stat-card">
          <div className="stat-card-header">
            <span className="stat-card-title">Всего сделок</span>
            <i className="bi bi-activity stat-card-icon"></i>
          </div>
          <div className="stat-card-value">
            {formatNumber(stats.totalTrades)}
          </div>
          <div className="stat-card-subtitle">За весь период</div>
        </div>
      </div>
    </div>
  );
}
