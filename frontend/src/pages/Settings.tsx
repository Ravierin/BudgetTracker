import { useState, useEffect } from 'react';

interface ApiKeys {
  bybitApiKey: string;
  bybitApiSecret: string;
  mexcApiKey: string;
  mexcApiSecret: string;
}

export function Settings() {
  const [keys, setKeys] = useState<ApiKeys>({
    bybitApiKey: '',
    bybitApiSecret: '',
    mexcApiKey: '',
    mexcApiSecret: '',
  });
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    // Load saved keys from localStorage
    const saved = localStorage.getItem('apiKeys');
    if (saved) {
      try {
        setKeys(JSON.parse(saved));
      } catch (e) {
        console.error('Failed to load API keys:', e);
      }
    }
  }, []);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    
    // Validate keys are not empty
    if (!keys.bybitApiKey || !keys.bybitApiSecret || !keys.mexcApiKey || !keys.mexcApiSecret) {
      setError('Все API ключи обязательны');
      return;
    }

    // Save to localStorage
    localStorage.setItem('apiKeys', JSON.stringify(keys));
    setSaved(true);
    setTimeout(() => setSaved(false), 3000);
  };

  const handleClear = () => {
    localStorage.removeItem('apiKeys');
    setKeys({
      bybitApiKey: '',
      bybitApiSecret: '',
      mexcApiKey: '',
      mexcApiSecret: '',
    });
    setSaved(true);
    setTimeout(() => setSaved(false), 3000);
  };

  return (
    <div>
      <h1 className="page-title mb-1">Настройки API</h1>
      <p className="page-subtitle mb-4">Управление ключами API бирж</p>

      <div className="table-container p-4" style={{ maxWidth: '600px' }}>
        <form onSubmit={handleSubmit}>
          {/* Bybit */}
          <div className="mb-4">
            <h5 className="mb-3">
              <i className="bi bi-box me-2"></i>
              Bybit
            </h5>
            <div className="mb-3">
              <label className="form-label">API Key</label>
              <input
                type="text"
                className="form-control"
                value={keys.bybitApiKey}
                onChange={(e) => setKeys({ ...keys, bybitApiKey: e.target.value })}
                placeholder="Введите Bybit API Key"
              />
            </div>
            <div className="mb-3">
              <label className="form-label">API Secret</label>
              <input
                type="password"
                className="form-control"
                value={keys.bybitApiSecret}
                onChange={(e) => setKeys({ ...keys, bybitApiSecret: e.target.value })}
                placeholder="Введите Bybit API Secret"
              />
            </div>
          </div>

          {/* MEXC */}
          <div className="mb-4">
            <h5 className="mb-3">
              <i className="bi bi-currency-exchange me-2"></i>
              MEXC
            </h5>
            <div className="mb-3">
              <label className="form-label">API Key</label>
              <input
                type="text"
                className="form-control"
                value={keys.mexcApiKey}
                onChange={(e) => setKeys({ ...keys, mexcApiKey: e.target.value })}
                placeholder="Введите MEXC API Key"
              />
            </div>
            <div className="mb-3">
              <label className="form-label">API Secret</label>
              <input
                type="password"
                className="form-control"
                value={keys.mexcApiSecret}
                onChange={(e) => setKeys({ ...keys, mexcApiSecret: e.target.value })}
                placeholder="Введите MEXC API Secret"
              />
            </div>
          </div>

          {error && (
            <div className="alert alert-danger" role="alert">
              {error}
            </div>
          )}

          {saved && (
            <div className="alert alert-success" role="alert">
              <i className="bi bi-check-circle me-2"></i>
              Ключи сохранены!
            </div>
          )}

          <div className="d-flex gap-2">
            <button type="submit" className="btn btn-primary">
              <i className="bi bi-save me-2"></i>
              Сохранить
            </button>
            <button type="button" className="btn btn-outline-secondary" onClick={handleClear}>
              <i className="bi bi-trash me-2"></i>
              Очистить
            </button>
          </div>
        </form>

        <div className="alert alert-info mt-4 mb-0">
          <i className="bi bi-info-circle me-2"></i>
          <strong>Важно:</strong> Ключи хранятся только в вашем браузере (localStorage) и не 
          передаются на сервер. Для синхронизации с биржами необходимо настроить прокси-сервер.
        </div>
      </div>
    </div>
  );
}
