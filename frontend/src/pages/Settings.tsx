import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  Key,
  Lock,
  Save,
  Trash2,
  CheckCircle,
  AlertCircle,
  Zap,
  Shield,
  RefreshCw
} from 'lucide-react';
import { api } from '../api/api';
import { EXCHANGES } from '../config/exchanges';
import type { APIKey, ExchangeApiKeys } from '../types';

export function Settings() {
  const [keys, setKeys] = useState<ExchangeApiKeys>({});
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    loadAPIKeys();
  }, []);

  const loadAPIKeys = async () => {
    try {
      const apiKeys = await api.getAPIKeys();
      const keysMap: ExchangeApiKeys = {};

      EXCHANGES.forEach(exchange => {
        const key = apiKeys.find(k => k.exchange === exchange.id);
        keysMap[exchange.id] = {
          apiKey: key?.apiKey || '',
          apiSecret: key?.apiSecret || '',
        };
      });

      setKeys(keysMap);
    } catch (e) {
      console.error('Failed to load API keys:', e);
    } finally {
      setLoading(false);
    }
  };

  const handleExchangeChange = (exchangeId: string, field: 'apiKey' | 'apiSecret', value: string) => {
    setKeys(prev => ({
      ...prev,
      [exchangeId]: {
        ...prev[exchangeId],
        [field]: value,
      },
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSaving(true);

    try {
      const apiKeys: APIKey[] = EXCHANGES.map(exchange => ({
        exchange: exchange.id,
        apiKey: keys[exchange.id]?.apiKey || '',
        apiSecret: keys[exchange.id]?.apiSecret || '',
      }));

      await api.saveAPIKeys(apiKeys);
      setSaved(true);
      setTimeout(() => setSaved(false), 3000);
    } catch (err) {
      setError('Не удалось сохранить ключи');
      console.error(err);
    } finally {
      setSaving(false);
    }
  };

  const handleClear = async () => {
    try {
      const apiKeys: APIKey[] = EXCHANGES.map(exchange => ({
        exchange: exchange.id,
        apiKey: '',
        apiSecret: '',
      }));

      await api.saveAPIKeys(apiKeys);

      const emptyKeys: ExchangeApiKeys = {};
      EXCHANGES.forEach(ex => {
        emptyKeys[ex.id] = { apiKey: '', apiSecret: '' };
      });
      setKeys(emptyKeys);

      setSaved(true);
      setTimeout(() => setSaved(false), 3000);
    } catch (err) {
      setError('Не удалось очистить ключи');
      console.error(err);
    }
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
        >
          Настройки API
        </motion.h1>
        <p className="page-subtitle">
          Управление ключами API для синхронизации с биржами
        </p>
      </div>

      {/* Alerts */}
      <AnimatePresence>
        {error && (
          <motion.div
            className="alert alert-error"
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
          >
            <AlertCircle size={20} />
            <span>{error}</span>
          </motion.div>
        )}

        {saved && (
          <motion.div
            className="alert alert-success"
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
          >
            <CheckCircle size={20} />
            <span>Ключи успешно сохранены!</span>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Main Form */}
      <div className="card">
        <div className="card-header">
          <h3 className="card-title" style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <Shield size={20} color="var(--accent-primary)" />
            Ключи API бирж
          </h3>
          <motion.div
            className="badge badge-info"
            animate={{ opacity: [1, 0.5, 1] }}
            transition={{ duration: 2, repeat: Infinity }}
          >
            <Zap size={14} />
            Auto-sync каждые 30 сек
          </motion.div>
        </div>

        <form onSubmit={handleSubmit}>
          {EXCHANGES.map((exchange, index) => (
            <motion.div
              key={exchange.id}
              className="exchange-section"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.1 }}
              style={{
                padding: '24px',
                marginBottom: '24px',
                background: 'var(--bg-tertiary)',
                borderRadius: 'var(--border-radius)',
                border: '1px solid var(--border-color)',
              }}
            >
              <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '20px' }}>
                <div
                  style={{
                    width: '40px',
                    height: '40px',
                    borderRadius: 'var(--border-radius)',
                    background: 'var(--bg-card)',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                  }}
                >
                  <Key size={20} color={exchange.id === 'mexc' ? '#00C076' : exchange.id === 'bybit' ? '#F7A600' : exchange.id === 'gate' ? '#F0B90B' : '#00D9FF'} />
                </div>
                <div>
                  <h4 style={{ fontSize: '16px', fontWeight: '600', color: 'var(--text-primary)' }}>
                    {exchange.name}
                  </h4>
                  <p style={{ fontSize: '12px', color: 'var(--text-muted)' }}>
                    Введите ключи для синхронизации
                  </p>
                </div>
              </div>

              <div style={{ display: 'grid', gap: '16px' }}>
                {/* API Key */}
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label className="form-label">
                    <Key size={14} style={{ marginRight: '6px', verticalAlign: 'middle' }} />
                    {exchange.apiKeyLabel}
                  </label>
                  <input
                    type="text"
                    className="form-input"
                    value={keys[exchange.id]?.apiKey || ''}
                    onChange={(e) => handleExchangeChange(exchange.id, 'apiKey', e.target.value)}
                    placeholder={`Введите ${exchange.name} API Key`}
                    style={{ paddingLeft: '44px' }}
                  />
                </div>

                {/* API Secret */}
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label className="form-label">
                    <Lock size={14} style={{ marginRight: '6px', verticalAlign: 'middle' }} />
                    {exchange.apiSecretLabel}
                  </label>
                  <input
                    type="password"
                    className="form-input"
                    value={keys[exchange.id]?.apiSecret || ''}
                    onChange={(e) => handleExchangeChange(exchange.id, 'apiSecret', e.target.value)}
                    placeholder={`Введите ${exchange.name} ${exchange.apiSecretLabel}`}
                    style={{ paddingLeft: '44px' }}
                  />
                </div>
              </div>
            </motion.div>
          ))}

          {/* Action Buttons */}
          <div style={{ display: 'flex', gap: '12px', marginTop: '32px' }}>
            <motion.button
              type="submit"
              className="btn btn-primary"
              disabled={saving}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              {saving ? (
                <RefreshCw size={18} className="spinning" />
              ) : (
                <Save size={18} />
              )}
              {saving ? 'Сохранение...' : 'Сохранить'}
            </motion.button>

            <motion.button
              type="button"
              className="btn btn-secondary"
              onClick={handleClear}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <Trash2 size={18} />
              Очистить всё
            </motion.button>
          </div>
        </form>
      </div>

      {/* Info Card */}
      <motion.div
        className="card"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.4 }}
        style={{ background: 'var(--info-bg)', borderColor: 'var(--info)' }}
      >
        <div style={{ display: 'flex', gap: '16px', alignItems: 'flex-start' }}>
          <AlertCircle size={24} color="var(--info)" style={{ flexShrink: 0 }} />
          <div>
            <h4 style={{ fontSize: '14px', fontWeight: '600', color: 'var(--info)', marginBottom: '8px' }}>
              Безопасность и синхронизация
            </h4>
            <p style={{ fontSize: '13px', color: 'var(--text-secondary)', lineHeight: '1.6' }}>
              Ключи хранятся в зашифрованном виде в базе данных на сервере.
              Синхронизация с биржами происходит автоматически каждые 30 секунд.
              Данные обновляются в реальном времени через WebSocket соединение.
            </p>
          </div>
        </div>
      </motion.div>
    </div>
  );
}
