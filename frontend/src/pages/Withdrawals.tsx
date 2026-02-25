import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Wallet, Plus, Trash2, Filter } from 'lucide-react';
import { api } from '../api/api';
import { wsService, type WSMessage } from '../api/websocket';
import type { Withdrawal } from '../types';

const EXCHANGES = [
  { id: 'mexc', name: 'MEXC', color: '#00C076' },
  { id: 'bybit', name: 'Bybit', color: '#F7A600' },
  { id: 'gate', name: 'Gate', color: '#F0B90B' },
  { id: 'bitget', name: 'Bitget', color: '#00D9FF' },
];

const CURRENCIES = [
  { code: 'USDT', name: 'Tether' },
  { code: 'RUB', name: 'Рубль' },
  { code: 'USD', name: 'Доллар' },
];

export function Withdrawals() {
  const [withdrawals, setWithdrawals] = useState<Withdrawal[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [filterExchange, setFilterExchange] = useState<string>('');
  const [formData, setFormData] = useState({
    exchange: 'mexc',
    amount: '',
    currency: 'USDT',
    date: new Date().toISOString().slice(0, 16),
  });

  const loadWithdrawals = async () => {
    try {
      const data = await api.getWithdrawals(filterExchange || undefined);
      setWithdrawals(data);
    } catch (error) {
      console.error('Failed to load withdrawals:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadWithdrawals();

    const unsubscribe = wsService.addListener((message: WSMessage) => {
      if (['withdrawal_created', 'withdrawal_deleted'].includes(message.type)) {
        loadWithdrawals();
      }
    });

    return () => unsubscribe();
  }, [filterExchange]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.createWithdrawal({
        exchange: formData.exchange,
        amount: parseFloat(formData.amount),
        currency: formData.currency,
        date: new Date(formData.date).toISOString(),
      });
      setShowForm(false);
      setFormData({
        exchange: 'mexc',
        amount: '',
        currency: 'USDT',
        date: new Date().toISOString().slice(0, 16),
      });
      await loadWithdrawals();
    } catch (error) {
      console.error('Failed to create withdrawal:', error);
    }
  };

  const handleDelete = async (id: number) => {
    if (confirm('Вы уверены, что хотите удалить эту запись?')) {
      try {
        await api.deleteWithdrawal(id);
        await loadWithdrawals();
      } catch (error) {
        console.error('Failed to delete withdrawal:', error);
      }
    }
  };

  const formatCurrency = (value: number, currency: string) => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: currency === 'USDT' ? 'USD' : currency,
      minimumFractionDigits: 2,
    }).format(value);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ru-RU', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const totalWithdrawn = withdrawals.reduce((sum, w) => sum + w.amount, 0);

  return (
    <div>
      {/* Page Header */}
      <div className="page-header">
        <motion.h1
          className="page-title"
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
        >
          P2P Выводы
        </motion.h1>
        <p className="page-subtitle">Учет вывода средств с бирж</p>
      </div>

      {/* Summary Card */}
      <motion.div
        className="stat-card"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        style={{ marginBottom: '24px' }}
      >
        <div className="stat-card-header">
          <span className="stat-card-title">Всего выведено</span>
          <div
            className="stat-card-icon"
            style={{ background: 'linear-gradient(135deg, #10b981 0%, #059669 100%)' }}
          >
            <Wallet size={20} color="white" />
          </div>
        </div>
        <div className="stat-card-value">{formatCurrency(totalWithdrawn, 'USD')}</div>
        <div className="stat-card-subtitle">Общая сумма всех выводов</div>
      </motion.div>

      {/* Actions Bar */}
      <div className="card">
        <div className="card-header">
          <h3 className="card-title" style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <Filter size={20} color="var(--accent-primary)" />
            Фильтр
          </h3>
          <motion.button
            className="btn btn-primary"
            onClick={() => setShowForm(!showForm)}
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
          >
            <Plus size={18} />
            {showForm ? 'Отмена' : 'Добавить вывод'}
          </motion.button>
        </div>
        <div style={{ display: 'flex', gap: '12px', flexWrap: 'wrap' }}>
          <button
            className={`btn ${filterExchange === '' ? 'btn-primary' : 'btn-secondary'}`}
            onClick={() => setFilterExchange('')}
          >
            Все биржи
          </button>
          {EXCHANGES.map((ex) => (
            <button
              key={ex.id}
              className={`btn ${filterExchange === ex.id ? 'btn-primary' : 'btn-secondary'}`}
              onClick={() => setFilterExchange(ex.id)}
            >
              {ex.name}
            </button>
          ))}
        </div>
      </div>

      {/* Add Form */}
      <AnimatePresence>
        {showForm && (
          <motion.div
            className="card"
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
          >
            <div className="card-header">
              <h3 className="card-title">Новый вывод</h3>
            </div>
            <form onSubmit={handleSubmit}>
              <div style={{ display: 'grid', gap: '16px', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))' }}>
                {/* Exchange */}
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label className="form-label">Биржа</label>
                  <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: '8px' }}>
                    {EXCHANGES.map((ex) => (
                      <button
                        key={ex.id}
                        type="button"
                        className={`btn ${formData.exchange === ex.id ? 'btn-primary' : 'btn-secondary'}`}
                        onClick={() => setFormData({ ...formData, exchange: ex.id })}
                        style={{ fontSize: '13px', padding: '8px' }}
                      >
                        {ex.name}
                      </button>
                    ))}
                  </div>
                </div>

                {/* Amount */}
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label className="form-label">Сумма</label>
                  <input
                    type="number"
                    className="form-input"
                    value={formData.amount}
                    onChange={(e) => setFormData({ ...formData, amount: e.target.value })}
                    placeholder="0.00"
                    step="0.01"
                    required
                  />
                </div>

                {/* Currency */}
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label className="form-label">Валюта</label>
                  <select
                    className="form-input"
                    value={formData.currency}
                    onChange={(e) => setFormData({ ...formData, currency: e.target.value })}
                  >
                    {CURRENCIES.map((curr) => (
                      <option key={curr.code} value={curr.code}>
                        {curr.code} — {curr.name}
                      </option>
                    ))}
                  </select>
                </div>

                {/* Date */}
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label className="form-label">Дата и время</label>
                  <input
                    type="datetime-local"
                    className="form-input"
                    value={formData.date}
                    onChange={(e) => setFormData({ ...formData, date: e.target.value })}
                    required
                  />
                </div>
              </div>

              <div style={{ display: 'flex', gap: '12px', marginTop: '24px' }}>
                <motion.button
                  type="submit"
                  className="btn btn-primary"
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                >
                  <Plus size={18} />
                  Добавить
                </motion.button>
                <motion.button
                  type="button"
                  className="btn btn-secondary"
                  onClick={() => setShowForm(false)}
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                >
                  Отмена
                </motion.button>
              </div>
            </form>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Table */}
      {loading ? (
        <div className="loading-spinner">
          <div className="spinner" />
        </div>
      ) : withdrawals.length === 0 ? (
        <div className="card">
          <div style={{ textAlign: 'center', padding: '40px' }}>
            <Wallet size={48} color="var(--text-muted)" style={{ marginBottom: '16px' }} />
            <h3 style={{ color: 'var(--text-primary)', marginBottom: '8px' }}>
              Нет выводов
            </h3>
            <p style={{ color: 'var(--text-secondary)', marginBottom: '16px' }}>
              Добавьте первую запись о выводе средств
            </p>
            <button className="btn btn-primary" onClick={() => setShowForm(true)}>
              <Plus size={18} />
              Добавить вывод
            </button>
          </div>
        </div>
      ) : (
        <div className="table-container">
          <table className="table">
            <thead>
              <tr>
                <th>Биржа</th>
                <th>Сумма</th>
                <th>Валюта</th>
                <th>Дата</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {withdrawals.map((withdrawal) => (
                <motion.tr
                  key={withdrawal.id}
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  whileHover={{ background: 'var(--bg-hover)' }}
                >
                  <td>
                    <span className="badge badge-info">{withdrawal.exchange.toUpperCase()}</span>
                  </td>
                  <td style={{ fontWeight: '600', color: 'var(--text-primary)' }}>
                    {formatCurrency(withdrawal.amount, withdrawal.currency)}
                  </td>
                  <td>
                    <span className="badge badge-warning">{withdrawal.currency}</span>
                  </td>
                  <td style={{ color: 'var(--text-secondary)' }}>{formatDate(withdrawal.date)}</td>
                  <td>
                    <motion.button
                      className="btn btn-danger"
                      style={{ padding: '6px 12px', fontSize: '12px' }}
                      onClick={() => handleDelete(withdrawal.id)}
                      whileHover={{ scale: 1.05 }}
                      whileTap={{ scale: 0.95 }}
                    >
                      <Trash2 size={14} />
                    </motion.button>
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
