export interface Position {
  id: number;
  orderId: string;
  exchange: string;
  symbol: string;
  volume: number;
  leverage: number;
  closedPnl: number;
  side: string;
  date: string;
}

export interface Withdrawal {
  id: number;
  exchange: string;
  amount: number;
  currency: string;
  date: string;
}

export interface MonthlyIncome {
  id: number;
  exchange: string;
  amount: number;
  pnl: number;
  date: string;
}

export interface DashboardStats {
  totalBalance: number;
  totalIncome: number;
  monthlyPnl: number;
  totalTrades: number;
}

export interface APIKey {
  id?: number;
  exchange: string;
  apiKey: string;
  apiSecret: string;
  isActive?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

export interface ExchangeApiKey {
  apiKey: string;
  apiSecret: string;
}

export interface ExchangeApiKeys {
  [exchangeId: string]: ExchangeApiKey;
}
