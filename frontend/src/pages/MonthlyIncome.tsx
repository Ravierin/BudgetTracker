import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { TrendingUp, DollarSign, Calendar } from 'lucide-react';
import { api } from '../api/api';
import { wsService, type WSMessage } from '../api/websocket';
import type { MonthlyIncome } from '../types';

export function MonthlyIncome() {
  const [incomes, setIncomes] = useState<MonthlyIncome[]>([]);
  const [loading, setLoading] = useState(true);
  const [filterExchange, setFilterExchange] = useState<string>('');

  const loadIncomes = async () => {
    try {
      const data = await api.getMonthlyIncomes(filterExchange || undefined);
      setIncomes(data);
    } catch (error) {
      console.error('Failed to load incomes:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadIncomes();

    const unsubscribe = wsService.addListener((message: WSMessage) => {
      if (['monthly_income_created', 'monthly_income_deleted'].includes(message.type)) {
        loadIncomes();
      }
    });

    return () => unsubscribe();
  }, [filterExchange]);

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
    }).format(value);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ru-RU', {
      year: 'numeric',
      month: 'long',
    });
  };

  const totalPnl = incomes.reduce((sum, i) => sum + i.pnl, 0);
  const totalAmount = incomes.reduce((sum, i) => sum + i.amount, 0);

  return (
    <div>
      {/* Page Header */}
      <div className="page-header">
        <motion.h1
          className="page-title"
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
        >
          Месячный доход
        </motion.h1>
        <p className="page-subtitle">
          Автоматический расчет дохода по сделкам
        </p>
      </div>

      {/* Summary Cards */}
      <div className="stats-grid">
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
              style={{ background: totalPnl >= 0 ? 'linear-gradient(135deg, #15803d 0%, #14532d 100%)' : 'linear-gradient(135deg, #b91c1c 0%, #991b1b 100%)' }}
            >
              <TrendingUp size={20} color="white" />
            </div>
          </div>
          <div className={`stat-card-value ${totalPnl >= 0 ? 'positive' : 'negative'}`}>
            {totalPnl >= 0 ? '+' : ''}{totalPnl.toFixed(2)}
          </div>
          <div className="stat-card-subtitle">За все месяцы</div>
        </motion.div>

        <motion.div
          className="stat-card"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4, delay: 0.1 }}
        >
          <div className="stat-card-header">
            <span className="stat-card-title">Общий объём</span>
            <div
              className="stat-card-icon"
              style={{ background: 'linear-gradient(135deg, #15803d 0%, #14532d 100%)' }}
            >
              <DollarSign size={20} color="white" />
            </div>
          </div>
          <div className="stat-card-value">{formatCurrency(totalAmount)}</div>
          <div className="stat-card-subtitle">За все месяцы</div>
        </motion.div>
      </div>

      {/* Filter */}
      <div className="card">
        <div className="card-header">
          <h3 className="card-title" style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <Calendar size={20} color="var(--accent-primary)" />
            Фильтр по биржам
          </h3>
        </div>
        <div style={{ display: 'flex', gap: '12px', flexWrap: 'wrap' }}>
          <button
            className={`btn ${filterExchange === '' ? 'btn-primary' : 'btn-secondary'}`}
            onClick={() => setFilterExchange('')}
          >
            Все биржи
          </button>
          {['mexc', 'bybit'].map((ex) => (
            <button
              key={ex}
              className={`btn ${filterExchange === ex ? 'btn-primary' : 'btn-secondary'}`}
              onClick={() => setFilterExchange(ex)}
            >
              {ex.toUpperCase()}
            </button>
          ))}
        </div>
      </div>

      {/* Table */}
      {loading ? (
        <div className="loading-spinner">
          <div className="spinner" />
        </div>
      ) : incomes.length === 0 ? (
        <div className="card">
          <div style={{ textAlign: 'center', padding: '40px' }}>
            <Calendar size={48} color="var(--text-muted)" style={{ marginBottom: '16px' }} />
            <h3 style={{ color: 'var(--text-primary)', marginBottom: '8px' }}>
              Нет данных
            </h3>
            <p style={{ color: 'var(--text-secondary)' }}>
              Месячный доход рассчитывается автоматически на основе ваших сделок
            </p>
          </div>
        </div>
      ) : (
        <div className="table-container">
          <table className="table">
            <thead>
              <tr>
                <th>Биржа</th>
                <th>Месяц</th>
                <th>Объём</th>
                <th>PnL</th>
              </tr>
            </thead>
            <tbody>
              {incomes.map((income) => (
                <motion.tr
                  key={income.id}
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  whileHover={{ background: 'var(--bg-hover)' }}
                >
                  <td>
                    <span className="badge badge-info">{income.exchange.toUpperCase()}</span>
                  </td>
                  <td>{formatDate(income.date)}</td>
                  <td>{formatCurrency(income.amount)}</td>
                  <td>
                    <span className={income.pnl >= 0 ? 'text-success' : 'text-danger'} style={{ fontWeight: '600' }}>
                      {income.pnl >= 0 ? '+' : ''}{formatCurrency(income.pnl)}
                    </span>
                  </td>
                </motion.tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
