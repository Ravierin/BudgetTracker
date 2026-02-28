export interface ExchangeConfig {
  id: string;
  name: string;
  icon: string;
  apiKeyLabel: string;
  apiSecretLabel: string;
}

export const EXCHANGES: ExchangeConfig[] = [
  {
    id: 'mexc',
    name: 'MEXC',
    icon: 'bi-currency-exchange',
    apiKeyLabel: 'API Key',
    apiSecretLabel: 'API Secret',
  },
  {
    id: 'bybit',
    name: 'Bybit',
    icon: 'bi-box',
    apiKeyLabel: 'API Key',
    apiSecretLabel: 'API Secret',
  },
];

export const getExchangeById = (id: string): ExchangeConfig | undefined => {
  return EXCHANGES.find(exchange => exchange.id === id);
};

export const getExchangeIds = (): string[] => {
  return EXCHANGES.map(exchange => exchange.id);
};
