import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Sidebar from '@/components/Sidebar';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger, DialogClose } from '@/components/ui/dialog';
import { ArrowLeft, Settings as SettingsIcon, Wifi, Shield, Save, Upload, Download, Eye, EyeOff, TestTube, AlertTriangle, Info, Plus, Sliders } from 'lucide-react';

import WidgetGrid, { Widget } from '@/components/WidgetGrid';
import { useTheme } from '@/contexts/ThemeContext';
import { Globe } from 'lucide-react';

interface OCPPSettings {
  connectionUrl: string;
  ocppVersion: '1.6' | '2.0.1';
  chargePointId: string;
  heartbeatInterval: number;
  retryInterval: number;
  maxRetries: number;
}

interface SecuritySettings {
  useTLS: boolean;
  clientCertificate: string;
  clientKey: string;
  caCertificate: string;
  basicAuth: boolean;
  username: string;
  password: string;
  verifyServerCertificate: boolean;
}

interface InfoSettings {
  model: string;
  vendor: string;
  firmwareVersion: string;
  maxChargingPower: number;
  freeCharging: boolean;
}

const ALL_WIDGETS: any[] = [
  {
    id: 'info',
    name: 'Info',
    icon: <Info className="h-5 w-5 mr-2" />,
  },
  {
    id: 'connectivity',
    name: 'Connectivity',
    icon: <Wifi className="h-5 w-5 mr-2" />,
  },
  {
    id: 'security',
    name: 'Security',
    icon: <Shield className="h-5 w-5 mr-2" />,
  },
  {
    id: 'user-preferences',
    name: 'User Preferences',
    icon: <Sliders className="h-5 w-5 mr-2" />,
  },
];

const DEFAULT_WIDGETS = [
  {
    id: 'info',
    x: 0, y: 0, w: 1, h: 2,
  },
  {
    id: 'connectivity',
    x: 1, y: 0, w: 1, h: 2,
  },
  {
    id: 'security',
    x: 2, y: 0, w: 1, h: 2,
  },
  {
    id: 'user-preferences',
    x: 0, y: 1, w: 1, h: 2,
  },
  {
    id: 'add-widget',
    x: 3, y: 0, w: 1, h: 1.2,
  },
];

const TIMEZONES = [
  'UTC',
  'Europe/Ljubljana',
  'Europe/London',
  'America/New_York',
  'Asia/Tokyo',
  'Asia/Shanghai',
  'Australia/Sydney',
];
const POWER_UNITS = ['kW', 'W'];
const CURRENT_UNITS = ['A', 'mA'];
const ENERGY_UNITS = ['kWh', 'Wh'];
const TIME_UNITS = ['seconds', 'minutes', 'hours'];
const DATE_FORMATS = [
  { value: 'MM/DD/YYYY', label: 'MM/DD/YYYY' },
  { value: 'DD/MM/YYYY', label: 'DD/MM/YYYY' },
  { value: 'YYYY-MM-DD', label: 'YYYY-MM-DD' },
  { value: 'MM-DD-YYYY', label: 'MM-DD-YYYY' },
  { value: 'DD-MM-YYYY', label: 'DD-MM-YYYY' },
];

const DATETIME_FORMATS = [
  { value: '24h', label: '24-hour (e.g. 14:30)' },
  { value: '24h-seconds', label: '24-hour with seconds (e.g. 14:30:45)' },
  { value: '12h', label: '12-hour (e.g. 2:30 PM)' },
  { value: '12h-seconds', label: '12-hour with seconds (e.g. 2:30:45 PM)' },
];

interface UserPreferences {
  timezone: string;
  powerUnit: string;
  currentUnit: string;
  energyUnit: string;
  timeUnit: string;
  dateFormat: string;
  dateTimeFormat: '24h' | '12h' | '24h-seconds' | '12h-seconds';
}

function getDefaultDateFormat(timezone: string): string {
  // Determine default date format based on timezone
  const timezoneLower = timezone.toLowerCase();
  
  // European timezones typically use DD/MM/YYYY
  if (timezoneLower.includes('europe') || 
      timezoneLower.includes('london') || 
      timezoneLower.includes('paris') || 
      timezoneLower.includes('berlin') ||
      timezoneLower.includes('rome') ||
      timezoneLower.includes('madrid')) {
    return 'DD/MM/YYYY';
  }
  
  // Asian timezones typically use YYYY-MM-DD
  if (timezoneLower.includes('asia') || 
      timezoneLower.includes('tokyo') || 
      timezoneLower.includes('beijing') || 
      timezoneLower.includes('shanghai') ||
      timezoneLower.includes('seoul')) {
    return 'YYYY-MM-DD';
  }
  
  // Australian timezones typically use DD/MM/YYYY
  if (timezoneLower.includes('australia') || 
      timezoneLower.includes('sydney') || 
      timezoneLower.includes('melbourne')) {
    return 'DD/MM/YYYY';
  }
  
  // Default to MM/DD/YYYY for US and other regions
  return 'MM/DD/YYYY';
}

function Settings() {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [connectivityLoading, setConnectivityLoading] = useState(false);
  const [securityLoading, setSecurityLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [testConnectionLoading, setTestConnectionLoading] = useState(false);
  const [testConnectionResult, setTestConnectionResult] = useState<'success' | 'error' | null>(null);
  const [infoLoading, setInfoLoading] = useState(false);

  // Initialize with default values
  const [ocppSettings, setOcppSettings] = useState<OCPPSettings>({
    connectionUrl: 'wss://example-ocpp-server.com/ocpp/',
    ocppVersion: '1.6',
    chargePointId: 'CP001',
    heartbeatInterval: 60,
    retryInterval: 30,
    maxRetries: 5
  });

  const [securitySettings, setSecuritySettings] = useState<SecuritySettings>({
    useTLS: true,
    clientCertificate: '',
    clientKey: '',
    caCertificate: '',
    basicAuth: false,
    username: '',
    password: '',
    verifyServerCertificate: true
  });

  const [infoSettings, setInfoSettings] = useState<InfoSettings>({
    model: 'ChargePi Model X',
    vendor: 'ChargePi',
    firmwareVersion: '1.0.0',
    maxChargingPower: 22,
    freeCharging: false
  });

  const [widgets, setWidgets] = useState(DEFAULT_WIDGETS);
  const [addDialogOpen, setAddDialogOpen] = useState(false);

  // User Preferences State
  const { theme, toggleTheme } = useTheme();
  const [userPrefs, setUserPrefs] = useState<UserPreferences>(() => {
    const stored = localStorage.getItem('userPrefs');
    if (stored) return JSON.parse(stored);
    return {
      timezone: Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC',
      powerUnit: 'kW',
      currentUnit: 'A',
      energyUnit: 'kWh',
      timeUnit: 'minutes',
      dateFormat: getDefaultDateFormat(Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC'),
      dateTimeFormat: '24h',
    };
  });
  const handleUserPrefChange = (key: keyof UserPreferences, value: string) => {
    setUserPrefs(prev => {
      const next = { ...prev, [key]: value };
      localStorage.setItem('userPrefs', JSON.stringify(next));
      return next;
    });
  };

  const handleOcppSettingsChange = (field: keyof OCPPSettings, value: string | number) => {
    setOcppSettings(prev => ({
      ...prev,
      [field]: value
    }));
  };

  const handleSecuritySettingsChange = (field: keyof SecuritySettings, value: string | boolean) => {
    setSecuritySettings(prev => ({
      ...prev,
      [field]: value
    }));
  };

  const handleInfoChange = (field: keyof InfoSettings, value: string | number | boolean) => {
    setInfoSettings(prev => ({ ...prev, [field]: value }));
  };

  const handleSaveConnectivity = async () => {
    setConnectivityLoading(true);
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      console.log('Saving connectivity settings:', ocppSettings);
      // TODO: Implement actual API call
    } catch (error) {
      console.error('Failed to save connectivity settings:', error);
    } finally {
      setConnectivityLoading(false);
    }
  };

  const handleSaveSecurity = async () => {
    setSecurityLoading(true);
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      console.log('Saving security settings:', securitySettings);
      // TODO: Implement actual API call
    } catch (error) {
      console.error('Failed to save security settings:', error);
    } finally {
      setSecurityLoading(false);
    }
  };

  const handleSaveInfo = async () => {
    setInfoLoading(true);
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      console.log('Saving info settings:', infoSettings);
      // TODO: Implement actual API call
    } catch (error) {
      console.error('Failed to save info settings:', error);
    } finally {
      setInfoLoading(false);
    }
  };

  const handleTestConnection = async () => {
    setTestConnectionLoading(true);
    setTestConnectionResult(null);
    try {
      // Simulate connection test
      await new Promise(resolve => setTimeout(resolve, 2000));
      setTestConnectionResult('success');
    } catch (error) {
      setTestConnectionResult('error');
    } finally {
      setTestConnectionLoading(false);
    }
  };

  const handleFileUpload = (field: 'clientCertificate' | 'clientKey' | 'caCertificate') => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.pem,.crt,.key';
    input.onchange = (e) => {
      const file = (e.target as HTMLInputElement).files?.[0];
      if (file) {
        const reader = new FileReader();
        reader.onload = (e) => {
          const content = e.target?.result as string;
          handleSecuritySettingsChange(field, content);
        };
        reader.readAsText(file);
      }
    };
    input.click();
  };

  const handleRemoveWidget = (id: string) => {
    setWidgets(ws => ws.filter(w => w.id !== id));
  };
  
  const handleAddWidget = (id: string) => {
    // Find a default position for the new widget
    const defaultWidget = DEFAULT_WIDGETS.find(w => w.id === id);
    if (defaultWidget && !widgets.some(w => w.id === id)) {
      setWidgets(ws => [...ws.filter(w => w.id !== 'add-widget'), { ...defaultWidget }, ws.find(w => w.id === 'add-widget')!]);
    }
    setAddDialogOpen(false);
  };

  const handleLayoutChange = (layoutArr: any[]) => {
    setWidgets(ws => {
      // Only update x/y/w/h for widgets that exist in the layout
      return ws.map(w => {
        const l = layoutArr.find(lay => lay.i === w.id);
        if (l) {
          return { ...w, x: l.x, y: l.y, w: l.w, h: l.h };
        }
        return w;
      });
    });
  };

  const settingsWidgets: Widget[] = widgets.map(w => {
    let content: React.ReactNode = null;
    if (w.id === 'info') {
      content = (
        <Card className="h-full flex flex-col">
          <CardHeader className="flex-shrink-0">
            <CardTitle className="flex items-center gap-2">
              <Info className="h-5 w-5" />
              Info
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4 flex-1 overflow-y-auto widget-scrollable">
            <div className="flex flex-col space-y-4">
              <div className="space-y-2">
                <Label htmlFor="model">Model</Label>
                <Input
                  id="model"
                  value={infoSettings.model}
                  onChange={e => handleInfoChange('model', e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="vendor">Vendor</Label>
                <Input
                  id="vendor"
                  value={infoSettings.vendor}
                  onChange={e => handleInfoChange('vendor', e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="firmwareVersion">Firmware Version</Label>
                <div className="relative group">
                <Input
                  id="firmwareVersion"
                  value={infoSettings.firmwareVersion}
                    readOnly
                    className="bg-muted cursor-not-allowed opacity-80 pr-10"
                    tabIndex={-1}
                    aria-readonly="true"
                  />
                  <span className="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground group-hover:text-accent-orange">
                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="2" className="inline align-middle">
                      <path strokeLinecap="round" strokeLinejoin="round" d="M17 11V7a5 5 0 00-10 0v4M5 11h14a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2z" />
                    </svg>
                  </span>
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="maxChargingPower">Max Charging Power (kW)</Label>
                <Input
                  id="maxChargingPower"
                  type="number"
                  min="1"
                  max="350"
                  value={infoSettings.maxChargingPower}
                  onChange={e => handleInfoChange('maxChargingPower', parseInt(e.target.value))}
                />
              </div>
              <div className="flex items-center space-x-2 pt-2">
                <input
                  type="checkbox"
                  id="freeCharging"
                  checked={infoSettings.freeCharging}
                  onChange={e => handleInfoChange('freeCharging', e.target.checked)}
                  className="rounded border-gray-300"
                />
                <Label htmlFor="freeCharging">Free Charging</Label>
              </div>
            </div>
            <div className="pt-4 border-t">
              <Button
                onClick={handleSaveInfo}
                disabled={infoLoading}
                className="w-full bg-accent-orange/80 text-white hover:bg-accent-orange focus-visible:ring-accent-orange"
              >
                <Save className="h-4 w-4 mr-2" />
                {infoLoading ? 'Saving...' : 'Save'}
              </Button>
            </div>
          </CardContent>
        </Card>
      );
    } else if (w.id === 'connectivity') {
      content = (
        <Card className="h-full flex flex-col">
              <CardHeader className="flex-shrink-0">
                <CardTitle className="flex items-center gap-2">
                  <Wifi className="h-5 w-5" />
                  Connectivity
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4 flex-1 overflow-y-auto widget-scrollable">
                <div className="space-y-2">
                  <Label htmlFor="connectionUrl">Connection URL</Label>
                  <Input
                    id="connectionUrl"
                    type="url"
                    placeholder="wss://example-ocpp-server.com/ocpp/"
                    value={ocppSettings.connectionUrl}
                    onChange={(e) => handleOcppSettingsChange('connectionUrl', e.target.value)}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="ocppVersion">OCPP Version</Label>
                  <Select 
                    value={ocppSettings.ocppVersion} 
                    onValueChange={(value) => handleOcppSettingsChange('ocppVersion', value as '1.6' | '2.0.1')}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="1.6">OCPP 1.6</SelectItem>
                      <SelectItem value="2.0.1">OCPP 2.0.1</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="chargePointId">Charge Point ID</Label>
                  <Input
                    id="chargePointId"
                    placeholder="CP001"
                    value={ocppSettings.chargePointId}
                    onChange={(e) => handleOcppSettingsChange('chargePointId', e.target.value)}
                  />
                </div>

                <div className="grid grid-cols-3 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="heartbeatInterval">Heartbeat (seconds)</Label>
                    <Input
                      id="heartbeatInterval"
                      type="number"
                      min="30"
                      max="3600"
                      value={ocppSettings.heartbeatInterval}
                      onChange={(e) => handleOcppSettingsChange('heartbeatInterval', parseInt(e.target.value))}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="retryInterval">Retry Interval (seconds)</Label>
                    <Input
                      id="retryInterval"
                      type="number"
                      min="5"
                      max="300"
                      value={ocppSettings.retryInterval}
                      onChange={(e) => handleOcppSettingsChange('retryInterval', parseInt(e.target.value))}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="maxRetries">Max Retries</Label>
                    <Input
                      id="maxRetries"
                      type="number"
                      min="1"
                      max="10"
                      value={ocppSettings.maxRetries}
                      onChange={(e) => handleOcppSettingsChange('maxRetries', parseInt(e.target.value))}
                    />
                  </div>
                </div>
                
                <div className="pt-4 border-t flex gap-2">
                  <Button 
                    onClick={handleTestConnection}
                    disabled={testConnectionLoading}
                    className="flex-1 bg-accent-orange/80 text-white hover:bg-accent-orange focus-visible:ring-accent-orange"
                  >
                    <TestTube className="h-4 w-4 mr-2" />
                    {testConnectionLoading ? 'Testing...' : 'Test Connection'}
                  </Button>
                  <Button 
                    onClick={handleSaveConnectivity} 
                    disabled={connectivityLoading}
                    className="flex-1 bg-accent-orange/80 text-white hover:bg-accent-orange focus-visible:ring-accent-orange"
                  >
                    <Save className="h-4 w-4 mr-2" />
                    {connectivityLoading ? 'Saving...' : 'Save'}
                  </Button>
                </div>
                  {testConnectionResult && (
                    <div className={`mt-3 p-3 rounded-lg border ${
                      testConnectionResult === 'success' 
                        ? 'bg-green-50 border-green-200 text-green-800 dark:bg-green-900/20 dark:border-green-800 dark:text-green-200' 
                        : 'bg-red-50 border-red-200 text-red-800 dark:bg-red-900/20 dark:border-red-800 dark:text-red-200'
                    }`}>
                      <div className="flex items-center gap-2">
                        {testConnectionResult === 'success' ? (
                          <div className="w-2 h-2 bg-green-500 rounded-full" />
                        ) : (
                          <AlertTriangle className="h-4 w-4" />
                        )}
                        <span className="text-sm font-medium">
                          {testConnectionResult === 'success' 
                            ? 'Connection test successful!' 
                            : 'Connection test failed. Please check your settings.'}
                        </span>
                      </div>
                                         </div>
                   )}
                
               </CardContent>
             </Card>
        );
    } else if (w.id === 'security') {
      content = (
        <Card className="h-full flex flex-col">
              <CardHeader className="flex-shrink-0">
                <CardTitle className="flex items-center gap-2">
                  <Shield className="h-5 w-5" />
                  Security
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4 flex-1 overflow-y-auto widget-scrollable">
                <div className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    id="useTLS"
                    checked={securitySettings.useTLS}
                    onChange={(e) => handleSecuritySettingsChange('useTLS', e.target.checked)}
                    className="rounded border-gray-300"
                  />
                  <Label htmlFor="useTLS">Use TLS/SSL</Label>
                </div>

                {securitySettings.useTLS && (
                  <>
                    <div className="space-y-2">
                      <Label>Client Certificate</Label>
                      <div className="flex gap-2">
                        <Textarea
                          placeholder="Paste your client certificate here..."
                          value={securitySettings.clientCertificate}
                          onChange={(e) => handleSecuritySettingsChange('clientCertificate', e.target.value)}
                          className="flex-1"
                        />
                        <Button 
                          variant="outline" 
                          size="sm"
                          onClick={() => handleFileUpload('clientCertificate')}
                        >
                          <Upload className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>

                    <div className="space-y-2">
                      <Label>Client Key</Label>
                      <div className="flex gap-2">
                        <Textarea
                          placeholder="Paste your client key here..."
                          value={securitySettings.clientKey}
                          onChange={(e) => handleSecuritySettingsChange('clientKey', e.target.value)}
                          className="flex-1"
                        />
                        <Button 
                          variant="outline" 
                          size="sm"
                          onClick={() => handleFileUpload('clientKey')}
                        >
                          <Upload className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>

                    <div className="space-y-2">
                      <Label>CA Certificate</Label>
                      <div className="flex gap-2">
                        <Textarea
                          placeholder="Paste your CA certificate here..."
                          value={securitySettings.caCertificate}
                          onChange={(e) => handleSecuritySettingsChange('caCertificate', e.target.value)}
                          className="flex-1"
                        />
                        <Button 
                          variant="outline" 
                          size="sm"
                          onClick={() => handleFileUpload('caCertificate')}
                        >
                          <Upload className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>

                    <div className="flex items-center space-x-2">
                      <input
                        type="checkbox"
                        id="verifyServerCertificate"
                        checked={securitySettings.verifyServerCertificate}
                        onChange={(e) => handleSecuritySettingsChange('verifyServerCertificate', e.target.checked)}
                        className="rounded border-gray-300"
                      />
                      <Label htmlFor="verifyServerCertificate">Verify Server Certificate</Label>
                    </div>
                  </>
                )}

                <div className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    id="basicAuth"
                    checked={securitySettings.basicAuth}
                    onChange={(e) => handleSecuritySettingsChange('basicAuth', e.target.checked)}
                    className="rounded border-gray-300"
                  />
                  <Label htmlFor="basicAuth">Use Basic Authentication</Label>
                </div>

                {securitySettings.basicAuth && (
                  <div className="space-y-4">
                    <div className="space-y-2">
                      <Label htmlFor="username">Username</Label>
                      <Input
                        id="username"
                        placeholder="Enter username"
                        value={securitySettings.username}
                        onChange={(e) => handleSecuritySettingsChange('username', e.target.value)}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="password">Password</Label>
                      <div className="relative">
                        <Input
                          id="password"
                          type={showPassword ? 'text' : 'password'}
                          placeholder="Enter password"
                          value={securitySettings.password}
                          onChange={(e) => handleSecuritySettingsChange('password', e.target.value)}
                        />
                        <button
                          type="button"
                          className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent flex items-center justify-center"
                          onClick={() => setShowPassword(!showPassword)}
                        >
                          {showPassword ? (
                            <EyeOff className="h-4 w-4" />
                          ) : (
                            <Eye className="h-4 w-4" />
                          )}
                        </button>
                      </div>
                    </div>
                                     </div>
                 )}
                 
                 <div className="pt-4 border-t">
                   <Button 
                     onClick={handleSaveSecurity} 
                     disabled={securityLoading}
                     className="w-full bg-accent-orange/80 text-white hover:bg-accent-orange focus-visible:ring-accent-orange"
                   >
                     <Save className="h-4 w-4 mr-2" />
                     {securityLoading ? 'Saving...' : 'Save'}
                   </Button>
                 </div>
               </CardContent>
             </Card>
        );
    } else if (w.id === 'user-preferences') {
      content = (
        <Card className="h-full flex flex-col">
          <CardHeader className="flex-shrink-0">
            <CardTitle className="flex items-center gap-2">
              <Sliders className="h-5 w-5" />
              User Preferences
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4 flex-1 overflow-y-auto widget-scrollable">
            <div className="space-y-2">
              <Label>Dark Mode</Label>
              <Button variant="outline" onClick={toggleTheme} className="w-full">
                {theme === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'}
              </Button>
            </div>
            <div className="space-y-2">
              <Label htmlFor="timezone">Timezone</Label>
              <Select value={userPrefs.timezone} onValueChange={v => handleUserPrefChange('timezone', v)}>
                <SelectTrigger id="timezone"><SelectValue /></SelectTrigger>
                <SelectContent>
                  {TIMEZONES.map(tz => (
                    <SelectItem key={tz} value={tz}>{tz}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="powerUnit">Power Unit</Label>
              <Select value={userPrefs.powerUnit} onValueChange={v => handleUserPrefChange('powerUnit', v)}>
                <SelectTrigger id="powerUnit"><SelectValue /></SelectTrigger>
                <SelectContent>
                  {POWER_UNITS.map(u => (
                    <SelectItem key={u} value={u}>{u}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="currentUnit">Current Unit</Label>
              <Select value={userPrefs.currentUnit} onValueChange={v => handleUserPrefChange('currentUnit', v)}>
                <SelectTrigger id="currentUnit"><SelectValue /></SelectTrigger>
                <SelectContent>
                  {CURRENT_UNITS.map(u => (
                    <SelectItem key={u} value={u}>{u}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="energyUnit">Energy Unit</Label>
              <Select value={userPrefs.energyUnit} onValueChange={v => handleUserPrefChange('energyUnit', v)}>
                <SelectTrigger id="energyUnit"><SelectValue /></SelectTrigger>
                <SelectContent>
                  {ENERGY_UNITS.map(u => (
                    <SelectItem key={u} value={u}>{u}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="timeUnit">Time Unit</Label>
              <Select value={userPrefs.timeUnit} onValueChange={v => handleUserPrefChange('timeUnit', v)}>
                <SelectTrigger id="timeUnit"><SelectValue /></SelectTrigger>
                <SelectContent>
                  {TIME_UNITS.map(u => (
                    <SelectItem key={u} value={u}>{u}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="dateFormat">Date Format</Label>
              <Select value={userPrefs.dateFormat} onValueChange={v => handleUserPrefChange('dateFormat', v)}>
                <SelectTrigger id="dateFormat"><SelectValue /></SelectTrigger>
                <SelectContent>
                  {DATE_FORMATS.map(format => (
                    <SelectItem key={format.value} value={format.value}>{format.label}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="dateTimeFormat">Datetime Format</Label>
              <Select value={userPrefs.dateTimeFormat} onValueChange={v => handleUserPrefChange('dateTimeFormat', v)}>
                <SelectTrigger id="dateTimeFormat"><SelectValue /></SelectTrigger>
                <SelectContent>
                  {DATETIME_FORMATS.map(format => (
                    <SelectItem key={format.value} value={format.value}>{format.label}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </CardContent>
        </Card>
      );
    } else if (w.id === 'add-widget') {
      content = (
      <button
        className="w-full h-full flex flex-col items-center justify-center border-2 border-dashed border-accent-orange rounded-lg bg-card hover:bg-accent-orange/10 transition-colors focus:outline-none focus:ring-2 focus:ring-accent-orange opacity-60 pointer-events-auto"
        onClick={() => setAddDialogOpen(true)}
        type="button"
        tabIndex={0}
        aria-label="Add widget"
      >
        <Plus className="h-8 w-8 text-accent-orange mb-1" />
        <span className="text-accent-orange font-medium">Add Widget</span>
      </button>
      );
    }
    return { ...w, content };
  });

  return (
    <div className="min-h-screen bg-background overflow-hidden lg:pl-64">
      <Sidebar isOpen={true} onClose={() => {}} />
      <header className="sticky top-0 z-30 bg-card border-b">
        <div className="flex items-center justify-between px-6 py-4">
          <div className="flex items-center space-x-3">
            <img src="/logo.svg" alt="ChargePi Logo" className="h-8 w-8" />
            <h1 className="text-xl font-semibold">Settings</h1>
          </div>
          <div className="flex items-center space-x-4">
            
          </div>
        </div>
      </header>
      <main className="flex-1 p-6 min-w-0">
        <WidgetGrid
          widgets={settingsWidgets}
          onRemoveWidget={id => id !== 'add-widget' && handleRemoveWidget(id)}
          onLayoutChange={handleLayoutChange}
        />
        <Dialog open={addDialogOpen} onOpenChange={setAddDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Add Widget</DialogTitle>
              <DialogDescription>Select a widget to add to the grid.</DialogDescription>
            </DialogHeader>
            <div className="flex flex-col gap-4 mt-2">
              {ALL_WIDGETS.filter((w: any) => !widgets.some(widget => widget.id === w.id && w.id !== 'add-widget')).length === 0 ? (
                <div className="text-muted-foreground text-center">All widgets are already in the grid.</div>
              ) : (
                ALL_WIDGETS.filter((w: any) => !widgets.some(widget => widget.id === w.id && w.id !== 'add-widget')).map((w: any) => (
                  <button
                    key={w.id}
                    className="flex items-center px-4 py-2 rounded-lg border border-accent-orange text-accent-orange hover:bg-accent-orange/10 transition-colors focus:outline-none focus:ring-2 focus:ring-accent-orange"
                    onClick={() => handleAddWidget(w.id)}
                    type="button"
                  >
                    {w.icon}
                    <span className="font-medium">{w.name}</span>
                  </button>
                ))
              )}
            </div>
            <DialogClose asChild>
              <Button variant="outline" className="mt-4 w-full">Cancel</Button>
            </DialogClose>
          </DialogContent>
        </Dialog>
      </main>
    </div>
  );
};

export default Settings; 