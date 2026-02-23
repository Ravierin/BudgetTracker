export type WSMessage = {
  type: string;
  data?: any;
  positions?: any[];
  count?: number;
  exchange?: string;
  positionId?: number;
  withdrawalId?: number;
  incomeId?: number;
};

export type WSListener = (message: WSMessage) => void;

class WebSocketService {
  private ws: WebSocket | null = null;
  private listeners: WSListener[] = [];
  private reconnectTimeout = 3000;
  private isManualClose = false;

  connect(url: string = 'ws://localhost:8080/ws') {
    this.isManualClose = false;
    
    try {
      this.ws = new WebSocket(url);

      this.ws.onopen = () => {
        console.log('WebSocket connected');
        this.reconnectTimeout = 3000;
      };

      this.ws.onmessage = (event) => {
        try {
          const message: WSMessage = JSON.parse(event.data);
          this.listeners.forEach(listener => listener(message));
        } catch (e) {
          console.error('Failed to parse WS message:', e);
        }
      };

      this.ws.onclose = () => {
        console.log('WebSocket closed');
        if (!this.isManualClose) {
          setTimeout(() => this.connect(url), this.reconnectTimeout);
          this.reconnectTimeout = Math.min(this.reconnectTimeout * 2, 30000);
        }
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
      };
    } catch (e) {
      console.error('Failed to connect WebSocket:', e);
    }
  }

  disconnect() {
    this.isManualClose = true;
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  addListener(listener: WSListener) {
    this.listeners.push(listener);
    return () => {
      this.listeners = this.listeners.filter(l => l !== listener);
    };
  }

  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }
}

export const wsService = new WebSocketService();
