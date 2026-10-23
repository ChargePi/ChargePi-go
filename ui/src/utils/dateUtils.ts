// Utility functions for date formatting that respect user preferences

export interface UserPreferences {
  timezone: string;
  powerUnit: string;
  currentUnit: string;
  energyUnit: string;
  timeUnit: string;
  dateFormat: string;
  dateTimeFormat: '24h' | '12h' | '24h-seconds' | '12h-seconds';
}

export function getUserPrefs(): UserPreferences {
  const stored = localStorage.getItem('userPrefs');
  if (stored) return JSON.parse(stored);
  return {
    timezone: Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC',
    powerUnit: 'kW',
    currentUnit: 'A',
    energyUnit: 'kWh',
    timeUnit: 'seconds',
    dateFormat: 'MM/DD/YYYY',
    dateTimeFormat: '24h',
  };
}

function getTimeOptions(dateTimeFormat: string, tz: string): Intl.DateTimeFormatOptions {
  switch (dateTimeFormat) {
    case '12h':
      return { timeZone: tz, hour: '2-digit' as const, minute: '2-digit' as const, hour12: true };
    case '12h-seconds':
      return { timeZone: tz, hour: '2-digit' as const, minute: '2-digit' as const, second: '2-digit' as const, hour12: true };
    case '24h-seconds':
      return { timeZone: tz, hour: '2-digit' as const, minute: '2-digit' as const, second: '2-digit' as const, hour12: false };
    case '24h':
    default:
      return { timeZone: tz, hour: '2-digit' as const, minute: '2-digit' as const, hour12: false };
  }
}

export function formatDate(date: string | undefined, timezone?: string, dateFormat?: string, dateTimeFormat?: string): string {
  if (!date) return '-';
  
  const prefs = getUserPrefs();
  const tz = timezone || prefs.timezone;
  const format = dateFormat || prefs.dateFormat;
  const dtFormat = dateTimeFormat || prefs.dateTimeFormat;
  
  const dateObj = new Date(date);
  
  // Format based on user's date format preference with time
  let dateStr: string;
  switch (format) {
    case 'MM/DD/YYYY':
      dateStr = dateObj.toLocaleDateString('en-US', { 
        timeZone: tz,
        month: '2-digit',
        day: '2-digit',
        year: 'numeric'
      });
      break;
    case 'DD/MM/YYYY':
      dateStr = dateObj.toLocaleDateString('en-GB', { 
        timeZone: tz,
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
      });
      break;
    case 'YYYY-MM-DD':
      dateStr = dateObj.toLocaleDateString('sv-SE', { 
        timeZone: tz,
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      });
      break;
    case 'MM-DD-YYYY':
      dateStr = dateObj.toLocaleDateString('en-US', { 
        timeZone: tz,
        month: '2-digit',
        day: '2-digit',
        year: 'numeric'
      }).replace(/\//g, '-');
      break;
    case 'DD-MM-YYYY':
      dateStr = dateObj.toLocaleDateString('en-GB', { 
        timeZone: tz,
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
      }).replace(/\//g, '-');
      break;
    default:
      dateStr = dateObj.toLocaleDateString(undefined, { timeZone: tz });
  }
  
  // Add time to the date using user preference
  const timeStr = dateObj.toLocaleTimeString(undefined, getTimeOptions(dtFormat, tz));
  
  return `${dateStr} ${timeStr}`;
}

export function formatDateOnly(date: string | undefined, timezone?: string, dateFormat?: string): string {
  if (!date) return '-';
  
  const prefs = getUserPrefs();
  const tz = timezone || prefs.timezone;
  const format = dateFormat || prefs.dateFormat;
  
  const dateObj = new Date(date);
  
  // Format based on user's date format preference (date only, no time)
  switch (format) {
    case 'MM/DD/YYYY':
      return dateObj.toLocaleDateString('en-US', { 
        timeZone: tz,
        month: '2-digit',
        day: '2-digit',
        year: 'numeric'
      });
    case 'DD/MM/YYYY':
      return dateObj.toLocaleDateString('en-GB', { 
        timeZone: tz,
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
      });
    case 'YYYY-MM-DD':
      return dateObj.toLocaleDateString('sv-SE', { 
        timeZone: tz,
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      });
    case 'MM-DD-YYYY':
      return dateObj.toLocaleDateString('en-US', { 
        timeZone: tz,
        month: '2-digit',
        day: '2-digit',
        year: 'numeric'
      }).replace(/\//g, '-');
    case 'DD-MM-YYYY':
      return dateObj.toLocaleDateString('en-GB', { 
        timeZone: tz,
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
      }).replace(/\//g, '-');
    default:
      return dateObj.toLocaleDateString(undefined, { timeZone: tz });
  }
}

export function formatDateTime(date: string | undefined, timezone?: string, dateFormat?: string, dateTimeFormat?: string): string {
  // This is now an alias for formatDate since formatDate includes time
  return formatDate(date, timezone, dateFormat, dateTimeFormat);
}

export function formatDateShort(date: string | undefined, timezone?: string): string {
  if (!date) return '-';
  
  const prefs = getUserPrefs();
  const tz = timezone || prefs.timezone;
  
  const dateObj = new Date(date);
  return dateObj.toLocaleDateString(undefined, { 
    timeZone: tz,
    month: 'short',
    day: 'numeric'
  });
} 