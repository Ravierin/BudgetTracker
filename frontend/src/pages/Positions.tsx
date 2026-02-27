import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  ArrowLeftRight,
  Plus,
  Trash2,
  Filter,
  RefreshCw,
  TrendingUp
} from 'lucide-react';
import { api } from '../api/api';
import { wsService, type WSMessage } from '../api/websocket';
import type { Position } from '../types';

const EXCHANGES = [
  { id: 'mexc', name: 'MEXC', color: '#00C076' },
  { id: 'bybit', name: 'Bybit', color: '#F7A600' },
  { id: 'gate', name: 'Gate', color: '#F0B90B' },
  { id: 'bitget', name: 'Bitget', color: '#00D9FF' },
];

const OTHER_EXCHANGES = 'other';

// Generate unique manual ID that won't conflict with exchange IDs
const generateManualId = () => {
  const timestamp = Date.now();
  const random = Math.random().toString(36).substring(2, 6);
  return `manual_position_${timestamp}_${random}`;
};

const SIDES = [
  { value: 'Buy', label: 'Long', color: '#15803d' },
  { value: 'Sell', label: 'Short', color: '#b91c1c' },
];

export function Positions() {
  const [positions, setPositions] = useState<Position[]>([]);
  const [loading, setLoading] = useState(true);
  const [syncing, setSyncing] = useState(false);
  const [showForm, setShowForm] = useState(false);
  const [filterExchange, setFilterExchange] = useState<string>('');
  const [formData, setFormData] = useState({
    exchange: '',
    symbol: '',
    side: 'Buy' as 'Buy' | 'Sell',
    volume: '',
    margin: '',
    leverage: '1',
    closedPnl: '',
    date: new Date().toISOString().slice(0, 16),
  });

  const loadPositions = async () => {
    try {
      let data = await api.getPositions(filterExchange === OTHER_EXCHANGES ? undefined : filterExchange || undefined);
      
      // Filter "other" exchanges on client side
      if (filterExchange === OTHER_EXCHANGES) {
        const knownExchanges = EXCHANGES.map(ex => ex.id);
        data = data.filter(p => !knownExchanges.includes(p.exchange.toLowerCase()));
      }
      
      setPositions(data);
    } catch (error) {
      console.error('Failed to load positions:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPositions();

    const unsubscribe = wsService.addListener((message: WSMessage) => {
      if (['positions_update', 'position_created', 'position_deleted'].includes(message.type)) {
        loadPositions();
      }
    });

    return () => unsubscribe();
  }, [filterExchange]);

  const handleSync = async (exchange?: string) => {
    setSyncing(true);
    try {
      await api.syncPositions(exchange);
      await loadPositions();
    } catch (error) {
      console.error('Failed to sync positions:', error);
    } finally {
      setSyncing(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.createPosition({
        orderId: generateManualId(),
        exchange: formData.exchange,
        symbol: formData.symbol.toUpperCase(),
        side: formData.side,
        volume: parseFloat(formData.volume),
        margin: parseFloat(formData.margin),
        leverage: parseInt(formData.leverage),
        closedPnl: parseFloat(formData.closedPnl),
        date: new Date(formData.date).toISOString(),
      });
      setShowForm(false);
      setFormData({
        exchange: 'mexc',
        symbol: '',
        side: 'Buy',
        volume: '',
        margin: '',
        leverage: '1',
        closedPnl: '',
        date: new Date().toISOString().slice(0, 16),
      });
      await loadPositions();
    } catch (error) {
      console.error('Failed to create position:', error);
    }
  };

  const handleDelete = async (id: number) => {
    if (confirm('Вы уверены, что хотите удалить эту позицию?')) {
      try {
        await api.deletePosition(id);
        await loadPositions();
      } catch (error) {
        console.error('Failed to delete position:', error);
      }
    }
  };

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
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const totalPnl = positions.reduce((sum, p) => sum + p.closedPnl, 0);
  const totalVolume = positions.reduce((sum, p) => sum + p.volume, 0);

  return (
    <div>
      {/* Page Header */}
      <div className="page-header">
        <motion.h1
          className="page-title"
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
        >
          Сделки
        </motion.h1>
        <p className="page-subtitle">История закрытых позиций</p>
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
          <div className="stat-card-subtitle">По всем сделкам</div>
        </motion.div>

        <motion.div
          className="stat-card"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4, delay: 0.1 }}
        >
          <div className="stat-card-header">
            <span className="stat-card-title">Объём</span>
            <div
              className="stat-card-icon"
              style={{ background: 'linear-gradient(135deg, #15803d 0%, #14532d 100%)' }}
            >
              <ArrowLeftRight size={20} color="white" />
            </div>
          </div>
          <div className="stat-card-value">{formatCurrency(totalVolume)}</div>
          <div className="stat-card-subtitle">Общий объём</div>
        </motion.div>
      </div>

      {/* Actions Bar */}
      <div className="card">
        <div className="card-header">
          <h3 className="card-title" style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <Filter size={20} color="var(--accent-primary)" />
            Фильтр и действия
          </h3>
          <div style={{ display: 'flex', gap: '12px' }}>
            <motion.button
              className="btn btn-secondary"
              onClick={() => handleSync(filterExchange === OTHER_EXCHANGES ? undefined : filterExchange || undefined)}
              disabled={syncing}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <RefreshCw size={18} className={syncing ? 'spinning' : ''} />
              {syncing ? 'Синх...' : 'Синхронизация'}
            </motion.button>
            <motion.button
              className="btn btn-primary"
              onClick={() => setShowForm(!showForm)}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <Plus size={18} />
              {showForm ? 'Отмена' : 'Добавить сделку'}
            </motion.button>
          </div>
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
          <button
            className={`btn ${filterExchange === OTHER_EXCHANGES ? 'btn-primary' : 'btn-secondary'}`}
            onClick={() => setFilterExchange(OTHER_EXCHANGES)}
            title="Биржи не из списка (Huobi, OKX, Binance и др.)"
          >
            🔹 Другое
          </button>
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
              <h3 className="card-title">Новая сделка</h3>
            </div>
            <form onSubmit={handleSubmit}>
              <div style={{ display: 'grid', gap: '16px', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))' }}>
                {/* Exchange - Text Input */}
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label className="form-label">Биржа</label>
                  <input
                    type="text"
                    className="form-input"
                    value={formData.exchange}
                    onChange={(e) => setFormData({ ...formData, exchange: e.target.value.toUpperCase() })}
                    placeholder="MEXC, BYBIT, GATE и т.д."
                    required
                  />
                </div>

                {/* Symbol */}
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label className="form-label">Символ</label>
                  <input
                    type="text"
                    className="form-input"
                    value={formData.symbol}
                    onChange={(e) => setFormData({ ...formData, symbol: e.target.value })}
                    placeholder="BTCUSDT"
                    required
                  />
                </div>

                {/* Side */}
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label className="form-label">Сторона</label>
                  <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: '8px' }}>
                    {SIDES.map((side) => (
                      <button
                        key={side.value}
                        type="button"
                        className={`btn ${formData.side === side.value ? 'btn-primary' : 'btn-secondary'}`}
                        onClick={() => setFormData({ ...formData, side: side.value as 'Buy' | 'Sell' })}
                        style={{
                          fontSize: '13px',
                          padding: '8px',
                          background: formData.side === side.value ? side.color : 'transparent',
                          borderColor: side.color,
                          color: formData.side === side.value ? 'white' : 'var(--text-primary)',
                        }}
                      >
                        {side.label}
                      </button>
                    ))}
                  </div>
                </div>

                {/* Volume */}
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label className="form-label">Объем (USDT)</label>
                  <input
                    type="number"
                    className="form-input"
                    value={formData.volume}
                    onChange={(e) => setFormData({ ...formData, volume: e.target.value })}
                    placeholder="0.00"
                    step="0.01"
                    required
                  />
                </div>

                {/* Margin */}
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label className="form-label">Маржа ($)</label>
                  <input
                    type="number"
                    className="form-input"
                    value={formData.margin}
                    onChange={(e) => setFormData({ ...formData, margin: e.target.value })}
                    placeholder="0.00"
                    step="0.01"
                    required
                  />
                </div>

                {/* Leverage */}
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label className="form-label">Плечо</label>
                  <input
                    type="number"
                    className="form-input"
                    value={formData.leverage}
                    onChange={(e) => setFormData({ ...formData, leverage: e.target.value })}
                    placeholder="1"
                    min="1"
                    max="125"
                    required
                  />
                </div>

                {/* PnL */}
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label className="form-label">PnL ($)</label>
                  <input
                    type="number"
                    className="form-input"
                    value={formData.closedPnl}
                    onChange={(e) => setFormData({ ...formData, closedPnl: e.target.value })}
                    placeholder="0.00"
                    step="0.01"
                    required
                  />
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
      ) : positions.length === 0 ? (
        <div className="card">
          <div style={{ textAlign: 'center', padding: '40px' }}>
            <ArrowLeftRight size={48} color="var(--text-muted)" style={{ marginBottom: '16px' }} />
            <h3 style={{ color: 'var(--text-primary)', marginBottom: '8px' }}>
              Нет сделок
            </h3>
            <p style={{ color: 'var(--text-secondary)', marginBottom: '16px' }}>
              Синхронизируйте с биржами или добавьте вручную
            </p>
            <div style={{ display: 'flex', gap: '12px', justifyContent: 'center' }}>
              <button className="btn btn-primary" onClick={() => handleSync()}>
                <RefreshCw size={18} />
                Синхронизировать
              </button>
              <button className="btn btn-secondary" onClick={() => setShowForm(true)}>
                <Plus size={18} />
                Добавить вручную
              </button>
            </div>
          </div>
        </div>
      ) : (
        <div className="table-container">
          <table className="table">
            <thead>
              <tr>
                <th>Биржа</th>
                <th>Символ</th>
                <th>Сторона</th>
                <th>Объем (USDT)</th>
                <th>Маржа</th>
                <th>Плечо</th>
                <th>PnL</th>
                <th>Дата</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {positions.map((position) => (
                <motion.tr
                  key={position.id}
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  whileHover={{ background: 'var(--bg-hover)' }}
                >
                  <td>
                    <span className="badge badge-info">{position.exchange.toUpperCase()}</span>
                  </td>
                  <td style={{ fontWeight: '700', color: 'var(--text-primary)' }}>
                    {position.symbol}
                  </td>
                  <td>
                    <span
                      className="badge"
                      style={{
                        background: position.side === 'Buy' ? 'var(--success-bg)' : 'var(--danger-bg)',
                        color: position.side === 'Buy' ? 'var(--success)' : 'var(--danger)',
                      }}
                    >
                      {position.side === 'Buy' ? '🟢 Long' : '🔴 Short'}
                    </span>
                  </td>
                  <td style={{ color: 'var(--text-primary)' }}>{position.volume.toLocaleString('ru-RU')}</td>
                  <td style={{ color: 'var(--text-primary)' }}>{position.margin.toLocaleString('ru-RU')}</td>
                  <td>
                    <span className="badge badge-warning">{position.leverage}x</span>
                  </td>
                  <td>
                    <span
                      className={position.closedPnl >= 0 ? 'text-success' : 'text-danger'}
                      style={{ fontWeight: '600' }}
                    >
                      {position.closedPnl >= 0 ? '+' : ''}{position.closedPnl.toFixed(2)}
                    </span>
                  </td>
                  <td style={{ color: 'var(--text-secondary)' }}>{formatDate(position.date)}</td>
                  <td>
                    <motion.button
                      className="btn btn-danger"
                      style={{ padding: '6px 12px', fontSize: '12px' }}
                      onClick={() => handleDelete(position.id)}
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
