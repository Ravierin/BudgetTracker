import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import {
  Wallet,
  Activity,
  DollarSign,
  PiggyBank,
  ArrowUpRight,
  ArrowDownRight,
  Zap,
  ArrowLeftRight,
  Settings
} from 'lucide-react';
import { api } from '../api/api';
import { wsService, type WSMessage } from '../api/websocket';

interface Stats {
  totalBalance: number;
  totalIncome: number;
  monthlyPnl: number;
  totalTrades: number;
  winRate: number;
}

const statCards = [
  {
    key: 'totalBalance' as const,
    label: 'Общий баланс',
    icon: Wallet,
    subtitle: 'На всех биржах',
    gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  },
  {
    key: 'totalIncome' as const,
    label: 'Общий P&L',
    icon: DollarSign,
    subtitle: 'За всё время',
    gradient: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
  },
  {
    key: 'monthlyPnl' as const,
    label: 'Месячный P&L',
    icon: PiggyBank,
    subtitle: 'За текущий месяц',
    gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
  },
  {
    key: 'totalTrades' as const,
    label: 'Всего сделок',
    icon: Activity,
    subtitle: 'Количество позиций',
    gradient: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
  },
];

export function Dashboard() {
  const [stats, setStats] = useState<Stats>({
    totalBalance: 0,
    totalIncome: 0,
    monthlyPnl: 0,
    totalTrades: 0,
    winRate: 0,
  });
  const [loading, setLoading] = useState(true);
  const [prevStats, setPrevStats] = useState<Stats | null>(null);

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

      const winningTrades = positions.filter(p => p.closedPnl > 0).length;
      const winRate = positions.length > 0 ? (winningTrades / positions.length) * 100 : 0;

      setPrevStats({ ...stats });
      setStats({
        totalBalance,
        totalIncome,
        monthlyPnl,
        totalTrades: positions.length,
        winRate,
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

  const getChangeIndicator = (current: number, prev: number | null) => {
    if (prev === null || prev === 0) return null;
    const change = ((current - prev) / prev) * 100;
    return change;
  };

  if (loading) {
    return (
      <div className="loading-spinner">
        <div className="spinner" />
      </div>
    );
  }

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
          Дашборд
        </motion.h1>
        <p className="page-subtitle">
          Обзор вашей торговой активности в реальном времени
        </p>
      </div>

      {/* Stats Grid */}
      <div className="stats-grid">
        {statCards.map((card, index) => {
          const Icon = card.icon;
          const value = stats[card.key];
          const prevValue = prevStats?.[card.key] ?? null;
          const change = getChangeIndicator(value, prevValue);
          const isPositive = value >= 0;

          return (
            <motion.div
              key={card.key}
              className="stat-card"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.4, delay: index * 0.1 }}
              whileHover={{ scale: 1.03, y: -8 }}
            >
              <div className="stat-card-header">
                <span className="stat-card-title">{card.label}</span>
                <div
                  className="stat-card-icon"
                  style={{ background: card.gradient }}
                >
                  <Icon size={20} color="white" />
                </div>
              </div>

              <motion.div
                className={`stat-card-value ${isPositive ? 'positive' : 'negative'}`}
                initial={{ scale: 0.9 }}
                animate={{ scale: 1 }}
                transition={{ duration: 0.3 }}
              >
                {card.key === 'totalTrades'
                  ? formatNumber(value)
                  : formatCurrency(value)}
              </motion.div>

              <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                {change !== null && Math.abs(change) > 0.1 && (
                  <>
                    {change > 0 ? (
                      <ArrowUpRight size={16} color="var(--success)" />
                    ) : (
                      <ArrowDownRight size={16} color="var(--danger)" />
                    )}
                    <span
                      style={{
                        fontSize: '12px',
                        color: change > 0 ? 'var(--success)' : 'var(--danger)',
                      }}
                    >
                      {change > 0 ? '+' : ''}{change.toFixed(2)}%
                    </span>
                  </>
                )}
                <span className="stat-card-subtitle" style={{ marginLeft: 'auto' }}>
                  {card.subtitle}
                </span>
              </div>

              {/* Decorative gradient bar */}
              <div
                style={{
                  position: 'absolute',
                  bottom: 0,
                  left: 0,
                  right: 0,
                  height: '3px',
                  background: card.gradient,
                  opacity: 0.5,
                }}
              />
            </motion.div>
          );
        })}
      </div>

      {/* Additional Info Cards */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(400px, 1fr))', gap: '24px' }}>
        {/* Win Rate Card */}
        <motion.div
          className="card"
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.5, delay: 0.4 }}
        >
          <div className="card-header">
            <h3 className="card-title" style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
              <Activity size={20} color="var(--accent-primary)" />
              Статистика
            </h3>
          </div>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
            <div>
              <div style={{ fontSize: '13px', color: 'var(--text-secondary)', marginBottom: '8px' }}>
                Win Rate
              </div>
              <div style={{ fontSize: '32px', fontWeight: '700', color: stats.winRate >= 50 ? 'var(--success)' : 'var(--warning)' }}>
                {stats.winRate.toFixed(1)}%
              </div>
              <div style={{ fontSize: '12px', color: 'var(--text-muted)', marginTop: '4px' }}>
                Успешных сделок
              </div>
            </div>
            <div>
              <div style={{ fontSize: '13px', color: 'var(--text-secondary)', marginBottom: '8px' }}>
                Активных позиций
              </div>
              <div style={{ fontSize: '32px', fontWeight: '700', color: 'var(--info)' }}>
                {stats.totalTrades}
              </div>
              <div style={{ fontSize: '12px', color: 'var(--text-muted)', marginTop: '4px' }}>
                Всего открыто
              </div>
            </div>
          </div>
        </motion.div>

        {/* Quick Actions */}
        <motion.div
          className="card"
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.5, delay: 0.5 }}
        >
          <div className="card-header">
            <h3 className="card-title" style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
              <Zap size={20} color="var(--accent-primary)" />
              Быстрые действия
            </h3>
          </div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            <motion.button
              className="btn btn-primary"
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
              onClick={() => window.location.href = '/positions'}
            >
              <ArrowLeftRight size={18} />
              Новая сделка
            </motion.button>
            <motion.button
              className="btn btn-secondary"
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
              onClick={() => window.location.href = '/settings'}
            >
              <Settings size={18} />
              Настройки API
            </motion.button>
          </div>
        </motion.div>
      </div>
    </div>
  );
}
