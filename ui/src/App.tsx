import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from '@/contexts/ThemeContext';
import { AuthProvider } from '@/contexts/AuthContext';
import LoginPage from '@/pages/LoginPage';
import Dashboard from '@/pages/Dashboard';
import Diagnostics from '@/pages/Diagnostics';
import Components from '@/pages/Components';
import ProtectedRoute from '@/components/ProtectedRoute';
import Hardware from './pages/Hardware';
import RfidTags from '@/pages/RfidTags';
import Users from '@/pages/Users';
import Settings from '@/pages/Settings';
import CommandPalette from '@/components/command-palette';
import Sessions from '@/pages/Sessions';

function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <Router>
          <div className="min-h-screen bg-background">
            <CommandPalette />
            <Routes>
              <Route path="/login" element={<LoginPage />} />
              <Route 
                path="/dashboard" 
                element={
                  <ProtectedRoute>
                    <Dashboard />
                  </ProtectedRoute>
                } 
              />
              <Route 
                path="/diagnostics" 
                element={
                  <ProtectedRoute>
                    <Diagnostics />
                  </ProtectedRoute>
                } 
              />
              <Route 
                path="/hardware"
                element={
                  <ProtectedRoute>
                    <Hardware />
                  </ProtectedRoute>
                }
              />
              <Route 
                path="/rfid-tags" 
                element={
                  <ProtectedRoute>
                    <RfidTags />
                  </ProtectedRoute>
                } 
              />
              <Route 
                path="/users" 
                element={
                  <ProtectedRoute>
                    <Users />
                  </ProtectedRoute>
                } 
              />
              <Route 
                path="/settings" 
                element={
                  <ProtectedRoute>
                    <Settings />
                  </ProtectedRoute>
                } 
              />
              <Route 
                path="/sessions" 
                element={
                  <ProtectedRoute>
                    <Sessions />
                  </ProtectedRoute>
                } 
              />
              <Route path="/" element={<Navigate to="/dashboard" replace />} />
            </Routes>
          </div>
        </Router>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App; 