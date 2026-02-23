import { useState, useEffect } from 'react';
import { api } from '../api/api';
import { wsService, type WSMessage } from '../api/websocket';
import type { Position } from '../types';

export function Positions() {
  const [positions, setPositions] = useState<Position[]>([]);
  const [loading, setLoading] = useState(true);
  const [syncing, setSyncing] = useState(false);
  const [filterExchange, setFilterExchange] = useState<string>('');

  const loadPositions = async () => {
    try {
      const data = await api.getPositions(filterExchange || undefined);
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

  return (
    <div>
      <div className="d-flex justify-content-between align-items-center mb-4">
        <div>
          <h1 className="page-title mb-1">Сделки</h1>
          <p className="page-subtitle mb-0">История закрытых позиций</p>
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
          <button
            className={`btn btn-primary sync-btn ${syncing ? 'spinning' : ''}`}
            onClick={() => handleSync(filterExchange || undefined)}
            disabled={syncing}
          >
            <i className={`bi bi-arrow-clockwise`}></i>
            {syncing ? 'Синхронизация...' : 'Синхронизировать'}
          </button>
        </div>
      </div>

      {loading ? (
        <div className="d-flex justify-content-center py-5">
          <div className="spinner-border text-primary" role="status">
            <span className="visually-hidden">Загрузка...</span>
          </div>
        </div>
      ) : positions.length === 0 ? (
        <div className="table-container p-5 text-center">
          <i className="bi bi-inbox display-4 text-muted"></i>
          <p className="text-muted mt-3 mb-0">Нет позиций</p>
          <button className="btn btn-primary mt-3" onClick={() => handleSync()}>
            Синхронизировать с биржами
          </button>
        </div>
      ) : (
        <div className="table-container">
          <table className="table">
            <thead>
              <tr>
                <th>Биржа</th>
                <th>Символ</th>
                <th>Сторона</th>
                <th>Количество</th>
                <th>Плечо</th>
                <th>Объём</th>
                <th>PnL</th>
                <th>Дата</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {positions.map((position) => (
                <tr key={position.id}>
                  <td>
                    <span className={`badge badge-bg-${position.exchange}`}>
                      {position.exchange.toUpperCase()}
                    </span>
                  </td>
                  <td>
                    <strong>{position.symbol}</strong>
                  </td>
                  <td>
                    <span className={position.side === 'Buy' ? 'text-success' : 'text-danger'}>
                      {position.side === 'Buy' ? '🟢 Long' : '🔴 Short'}
                    </span>
                  </td>
                  <td>{position.qty.toLocaleString('ru-RU')}</td>
                  <td>{position.leverage}x</td>
                  <td>{formatCurrency(position.cumExitValue)}</td>
                  <td>
                    <span className={position.closedPnl >= 0 ? 'text-success' : 'text-danger'}>
                      {position.closedPnl >= 0 ? '+' : ''}{formatCurrency(position.closedPnl)}
                    </span>
                  </td>
                  <td>{formatDate(position.date)}</td>
                  <td>
                    <button
                      className="btn btn-outline-secondary btn-icon"
                      onClick={() => handleDelete(position.id)}
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
    </div>
  );
}
