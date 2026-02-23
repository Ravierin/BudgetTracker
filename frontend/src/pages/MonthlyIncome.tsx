import { useState, useEffect } from 'react';
import { api } from '../api/api';
import { wsService, type WSMessage } from '../api/websocket';
import type { MonthlyIncome } from '../types';

export function MonthlyIncome() {
  const [incomes, setIncomes] = useState<MonthlyIncome[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [filterExchange, setFilterExchange] = useState<string>('');
  const [formData, setFormData] = useState({
    exchange: 'mexc',
    amount: '',
    pnl: '',
    date: new Date().toISOString().split('T')[0],
  });

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

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.createMonthlyIncome({
        exchange: formData.exchange,
        amount: parseFloat(formData.amount),
        pnl: parseFloat(formData.pnl),
        date: new Date(formData.date).toISOString(),
      });
      setShowModal(false);
      setFormData({
        exchange: 'mexc',
        amount: '',
        pnl: '',
        date: new Date().toISOString().split('T')[0],
      });
      await loadIncomes();
    } catch (error) {
      console.error('Failed to create income:', error);
    }
  };

  const handleDelete = async (id: number) => {
    if (confirm('Вы уверены, что хотите удалить эту запись?')) {
      try {
        await api.deleteMonthlyIncome(id);
        await loadIncomes();
      } catch (error) {
        console.error('Failed to delete income:', error);
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
      month: 'long',
      day: 'numeric',
    });
  };

  const totalPnl = incomes.reduce((sum, i) => sum + i.pnl, 0);
  const totalAmount = incomes.reduce((sum, i) => sum + i.amount, 0);

  return (
    <div>
      <div className="d-flex justify-content-between align-items-center mb-4">
        <div>
          <h1 className="page-title mb-1">Месячный доход</h1>
          <p className="page-subtitle mb-0">Статистика по месяцам</p>
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

      {/* Summary Cards */}
      <div className="stats-grid mb-4">
        <div className="stat-card">
          <div className="stat-card-header">
            <span className="stat-card-title">Общий PnL</span>
            <i className="bi bi-graph-up stat-card-icon"></i>
          </div>
          <div className={`stat-card-value ${totalPnl >= 0 ? 'positive' : 'negative'}`}>
            {formatCurrency(totalPnl)}
          </div>
          <div className="stat-card-subtitle">Все месяцы</div>
        </div>
        <div className="stat-card">
          <div className="stat-card-header">
            <span className="stat-card-title">Общий объём</span>
            <i className="bi bi-cash-stack stat-card-icon"></i>
          </div>
          <div className="stat-card-value">{formatCurrency(totalAmount)}</div>
          <div className="stat-card-subtitle">Все месяцы</div>
        </div>
      </div>

      {loading ? (
        <div className="d-flex justify-content-center py-5">
          <div className="spinner-border text-primary" role="status">
            <span className="visually-hidden">Загрузка...</span>
          </div>
        </div>
      ) : incomes.length === 0 ? (
        <div className="table-container p-5 text-center">
          <i className="bi bi-inbox display-4 text-muted"></i>
          <p className="text-muted mt-3 mb-0">Нет записей</p>
          <button className="btn btn-primary mt-3" onClick={() => setShowModal(true)}>
            Добавить первый доход
          </button>
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
                <th></th>
              </tr>
            </thead>
            <tbody>
              {incomes.map((income) => (
                <tr key={income.id}>
                  <td>
                    <span className={`badge badge-bg-${income.exchange}`}>
                      {income.exchange.toUpperCase()}
                    </span>
                  </td>
                  <td>{formatDate(income.date)}</td>
                  <td>{formatCurrency(income.amount)}</td>
                  <td>
                    <span className={income.pnl >= 0 ? 'text-success' : 'text-danger'}>
                      {income.pnl >= 0 ? '+' : ''}{formatCurrency(income.pnl)}
                    </span>
                  </td>
                  <td>
                    <button
                      className="btn btn-outline-secondary btn-icon"
                      onClick={() => handleDelete(income.id)}
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
                <h5 className="modal-title">Добавить месячный доход</h5>
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
                    <label className="form-label">Объём ($)</label>
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
                    <label className="form-label">PnL ($)</label>
                    <input
                      type="number"
                      className="form-control"
                      value={formData.pnl}
                      onChange={(e) => setFormData({ ...formData, pnl: e.target.value })}
                      step="0.01"
                      required
                    />
                  </div>
                  <div className="mb-3">
                    <label className="form-label">Месяц</label>
                    <input
                      type="month"
                      className="form-control"
                      value={formData.date.slice(0, 7)}
                      onChange={(e) => setFormData({ ...formData, date: e.target.value + '-01' })}
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
