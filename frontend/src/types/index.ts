export interface Position {
  id: number;
  orderId: string;
  exchange: string;
  symbol: string;
  cumExitValue: number;
  qty: number;
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
