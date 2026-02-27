import type { Position, Withdrawal, MonthlyIncome, APIKey } from '../types';

const API_BASE_URL = 'http://localhost:8080/api/v1';

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'API error');
  }
  return response.json();
}

export const api = {
  // Positions
  async getPositions(exchange?: string): Promise<Position[]> {
    const url = exchange 
      ? `${API_BASE_URL}/positions?exchange=${exchange}`
      : `${API_BASE_URL}/positions`;
    const response = await fetch(url);
    return handleResponse<Position[]>(response);
  },

  async createPosition(position: Omit<Position, 'id'>): Promise<Position> {
    const response = await fetch(`${API_BASE_URL}/positions`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(position),
    });
    return handleResponse<Position>(response);
  },

  async deletePosition(id: number): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/positions/${id}`, {
      method: 'DELETE',
    });
    return handleResponse<void>(response);
  },

  async syncPositions(exchange?: string): Promise<{ status: string; count: number; exchange: string }> {
    const url = exchange 
      ? `${API_BASE_URL}/positions/sync?exchange=${exchange}`
      : `${API_BASE_URL}/positions/sync`;
    const response = await fetch(url, { method: 'POST' });
    return handleResponse(response);
  },

  // Withdrawals
  async getWithdrawals(exchange?: string): Promise<Withdrawal[]> {
    const url = exchange 
      ? `${API_BASE_URL}/withdrawals?exchange=${exchange}`
      : `${API_BASE_URL}/withdrawals`;
    const response = await fetch(url);
    return handleResponse<Withdrawal[]>(response);
  },

  async createWithdrawal(withdrawal: Omit<Withdrawal, 'id'>): Promise<Withdrawal> {
    const response = await fetch(`${API_BASE_URL}/withdrawals`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(withdrawal),
    });
    return handleResponse<Withdrawal>(response);
  },

  async deleteWithdrawal(id: number): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/withdrawals/${id}`, {
      method: 'DELETE',
    });
    return handleResponse<void>(response);
  },

  // Monthly Income
  async getMonthlyIncomes(exchange?: string): Promise<MonthlyIncome[]> {
    const url = exchange 
      ? `${API_BASE_URL}/monthly-income?exchange=${exchange}`
      : `${API_BASE_URL}/monthly-income`;
    const response = await fetch(url);
    return handleResponse<MonthlyIncome[]>(response);
  },

  async createMonthlyIncome(income: Omit<MonthlyIncome, 'id'>): Promise<MonthlyIncome> {
    const response = await fetch(`${API_BASE_URL}/monthly-income`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(income),
    });
    return handleResponse<MonthlyIncome>(response);
  },

  async deleteMonthlyIncome(id: number): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/monthly-income/${id}`, {
      method: 'DELETE',
    });
    return handleResponse<void>(response);
  },

  // API Keys
  async getAPIKeys(): Promise<APIKey[]> {
    const response = await fetch(`${API_BASE_URL}/api-keys`);
    return handleResponse<APIKey[]>(response);
  },

  async saveAPIKeys(keys: APIKey[]): Promise<{ status: string }> {
    const response = await fetch(`${API_BASE_URL}/api-keys`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(keys),
    });
    return handleResponse(response);
  },

  // Balance
  async getBalance(): Promise<{ totalBalance: number; exchangeBalances: { exchange: string; balance: number }[] }> {
    const response = await fetch(`${API_BASE_URL}/balance`);
    return handleResponse(response);
  },
};
