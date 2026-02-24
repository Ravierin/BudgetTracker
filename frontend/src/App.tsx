import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Layout } from './components/Layout';
import { Dashboard } from './pages/Dashboard';
import { Positions } from './pages/Positions';
import { MonthlyIncome } from './pages/MonthlyIncome';
import { Withdrawals } from './pages/Withdrawals';
import { Settings } from './pages/Settings';
import { wsService } from './api/websocket';
import './App.css';

// Connect WebSocket on app start
wsService.connect();

function App() {
  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/positions" element={<Positions />} />
          <Route path="/monthly-income" element={<MonthlyIncome />} />
          <Route path="/withdrawals" element={<Withdrawals />} />
          <Route path="/settings" element={<Settings />} />
        </Routes>
      </Layout>
    </BrowserRouter>
  );
}

export default App;
