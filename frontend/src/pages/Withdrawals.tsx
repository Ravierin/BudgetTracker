import { useState, useEffect } from 'react';
import { api } from '../api/api';
import { wsService, type WSMessage } from '../api/websocket';
import type { Withdrawal } from '../types';

export function Withdrawals() {
  const [withdrawals, setWithdrawals] = useState<Withdrawal[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [filterExchange, setFilterExchange] = useState<string>('');
  const [formData, setFormData] = useState({
    exchange: 'mexc',
    amount: '',
    currency: 'USDT',
    date: new Date().toISOString().split('T')[0],
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
      setShowModal(false);
      setFormData({
        exchange: 'mexc',
        amount: '',
        currency: 'USDT',
        date: new Date().toISOString().split('T')[0],
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
      <div className="d-flex justify-content-between align-items-center mb-4">
        <div>
          <h1 className="page-title mb-1">P2P</h1>
          <p className="page-subtitle mb-0">Выводы средств</p>
        </div>
        <div className="d-flex gap-2">
          <select
            className="form-select"
            style={{ width: 'auto' }}
            value={filterExchange}
            onChange={(e) => setFilterExchange(e.target.value)}
          >
            <option value="">Все биржи</option>
            <option value="bybit">Bybit</option>
            <option value="mexc">MEXC</option>
          </select>
          <button className="btn btn-primary" onClick={() => setShowModal(true)}>
            <i className="bi bi-plus-lg"></i> Добавить
          </button>
        </div>
      </div>

      {/* Summary Card */}
      <div className="stats-grid mb-4">
        <div className="stat-card">
          <div className="stat-card-header">
            <span className="stat-card-title">Всего выведено</span>
            <i className="bi bi-wallet stat-card-icon"></i>
          </div>
          <div className="stat-card-value">{formatCurrency(totalWithdrawn, 'USD')}</div>
          <div className="stat-card-subtitle">Все выводы</div>
        </div>
      </div>

      {loading ? (
        <div className="d-flex justify-content-center py-5">
          <div className="spinner-border text-primary" role="status">
            <span className="visually-hidden">Загрузка...</span>
          </div>
        </div>
      ) : withdrawals.length === 0 ? (
        <div className="table-container p-5 text-center">
          <i className="bi bi-inbox display-4 text-muted"></i>
          <p className="text-muted mt-3 mb-0">Нет записей о выводах</p>
          <button className="btn btn-primary mt-3" onClick={() => setShowModal(true)}>
            Добавить первый вывод
          </button>
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
                <tr key={withdrawal.id}>
                  <td>
                    <span className={`badge badge-bg-${withdrawal.exchange}`}>
                      {withdrawal.exchange.toUpperCase()}
                    </span>
                  </td>
                  <td>{formatCurrency(withdrawal.amount, withdrawal.currency)}</td>
                  <td>
                    <span className="badge bg-secondary">{withdrawal.currency}</span>
                  </td>
                  <td>{formatDate(withdrawal.date)}</td>
                  <td>
                    <button
                      className="btn btn-outline-secondary btn-icon"
                      onClick={() => handleDelete(withdrawal.id)}
                    >
                      <i className="bi bi-trash"></i>
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Modal */}
      {showModal && (
        <div className="modal fade show d-block" tabIndex={-1}>
          <div className="modal-dialog">
            <div className="modal-content">
              <div className="modal-header">
                <h5 className="modal-title">Добавить вывод</h5>
                <button
                  type="button"
                  className="btn-close"
                  onClick={() => setShowModal(false)}
                ></button>
              </div>
              <form onSubmit={handleSubmit}>
                <div className="modal-body">
                  <div className="mb-3">
                    <label className="form-label">Биржа</label>
                    <select
                      className="form-select"
                      value={formData.exchange}
                      onChange={(e) => setFormData({ ...formData, exchange: e.target.value })}
                    >
                      <option value="mexc">MEXC</option>
                      <option value="bybit">Bybit</option>
                    </select>
                  </div>
                  <div className="mb-3">
                    <label className="form-label">Сумма</label>
                    <input
                      type="number"
                      className="form-control"
                      value={formData.amount}
                      onChange={(e) => setFormData({ ...formData, amount: e.target.value })}
                      step="0.01"
                      required
                    />
                  </div>
                  <div className="mb-3">
                    <label className="form-label">Валюта</label>
                    <select
                      className="form-select"
                      value={formData.currency}
                      onChange={(e) => setFormData({ ...formData, currency: e.target.value })}
                    >
                      <option value="USDT">USDT</option>
                      <option value="RUB">RUB</option>
                      <option value="USD">USD</option>
                    </select>
                  </div>
                  <div className="mb-3">
                    <label className="form-label">Дата</label>
                    <input
                      type="datetime-local"
                      className="form-control"
                      value={formData.date}
                      onChange={(e) => setFormData({ ...formData, date: e.target.value })}
                      required
                    />
                  </div>
                </div>
                <div className="modal-footer">
                  <button
                    type="button"
                    className="btn btn-outline-secondary"
                    onClick={() => setShowModal(false)}
                  >
                    Отмена
                  </button>
                  <button type="submit" className="btn btn-primary">
                    Добавить
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
