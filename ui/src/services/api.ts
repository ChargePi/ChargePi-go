import { authConfig } from '@/config/auth';

// Mock data
const mockConnectors = [
  { id: 1, status: 'available', name: 'Connector 1', power: 0, current: 0 },
  { id: 2, status: 'charging', name: 'Connector 2', power: 22.4, current: 32 },
  { id: 3, status: 'faulted', name: 'Connector 3', power: 0, current: 0 },
  { id: 4, status: 'available', name: 'Connector 4', power: 0, current: 0 },
];

const mockEnergyData = [
  { date: '2024-01-01', consumption: 45.2 },
  { date: '2024-01-02', consumption: 52.8 },
  { date: '2024-01-03', consumption: 38.9 },
  { date: '2024-01-04', consumption: 61.3 },
  { date: '2024-01-05', consumption: 49.7 },
  { date: '2024-01-06', consumption: 55.1 },
  { date: '2024-01-07', consumption: 42.8 },
];

const mockSessions = [
  {
    id: 1,
    user: 'John Doe',
    connector: 'Connector 2',
    startTime: '2024-01-07 14:30',
    duration: '2h 15m',
    energy: 45.2,
    status: 'completed'
  },
  {
    id: 2,
    user: 'Jane Smith',
    connector: 'Connector 1',
    startTime: '2024-01-07 10:15',
    duration: '1h 45m',
    energy: 32.8,
    status: 'completed'
  },
  {
    id: 3,
    user: 'Bob Johnson',
    connector: 'Connector 3',
    startTime: '2024-01-07 08:00',
    duration: '3h 30m',
    energy: 58.9,
    status: 'completed'
  }
];

const mockFaults = [
  {
    id: 1,
    type: 'error',
    message: 'Connector 3 communication failure',
    connector: 'Connector 3',
    timestamp: '2024-01-07 15:30',
    severity: 'high'
  },
  {
    id: 2,
    type: 'warning',
    message: 'Temperature sensor reading high',
    connector: 'Connector 1',
    timestamp: '2024-01-07 14:15',
    severity: 'medium'
  }
];

const mockCurrentSession = {
  id: 1,
  user: 'Alice Johnson',
  connector: 'Connector 2',
  startTime: '2024-01-07 16:00',
  duration: '1h 23m',
  energyDelivered: 28.5,
  currentPower: 22.4,
  targetEnergy: 50.0,
  progress: 57
};

// Simulate API delay
const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

export const api = {
  // Authentication
  async login(username: string, password: string) {
    await delay(500);
    
    if (!authConfig.enabled) {
      return { success: true, user: authConfig.mockUser };
    }
    
    if (username === authConfig.mockCredentials.username && 
        password === authConfig.mockCredentials.password) {
      return { success: true, user: authConfig.mockUser };
    }
    
    throw new Error('Invalid credentials');
  },

  // Charge point status
  async getStatus() {
    await delay(200);
    return {
      connectors: mockConnectors,
      totalConnectors: mockConnectors.length,
      availableConnectors: mockConnectors.filter(c => c.status === 'available').length,
      chargingConnectors: mockConnectors.filter(c => c.status === 'charging').length,
      faultedConnectors: mockConnectors.filter(c => c.status === 'faulted').length,
    };
  },

  // Energy consumption
  async getEnergyData() {
    await delay(300);
    return {
      data: mockEnergyData,
      total: mockEnergyData.reduce((sum, day) => sum + day.consumption, 0),
      average: mockEnergyData.reduce((sum, day) => sum + day.consumption, 0) / mockEnergyData.length,
    };
  },

  // Sessions
  async getSessions() {
    await delay(250);
    return {
      sessions: mockSessions,
      total: mockSessions.length,
      totalEnergy: mockSessions.reduce((sum, s) => sum + s.energy, 0),
    };
  },

  // Diagnostics
  async getDiagnostics() {
    await delay(150);
    return {
      diagnostics: mockFaults,
      total: mockFaults.length,
      highPriority: mockFaults.filter(f => f.severity === 'high').length,
    };
  },

  // Current session
  async getCurrentSession() {
    await delay(100);
    return mockCurrentSession;
  },

  // Real-time updates (simulated)
  async subscribeToUpdates(callback: (data: any) => void) {
    // Simulate real-time updates every 5 seconds
    const interval = setInterval(async () => {
      const status = await this.getStatus();
      const currentSession = await this.getCurrentSession();
      callback({ status, currentSession });
    }, 5000);

    return () => clearInterval(interval);
  }
}; 