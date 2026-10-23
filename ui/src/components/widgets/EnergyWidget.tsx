import React from 'react';
import { X, Zap, TrendingUp, Calendar } from 'lucide-react';
import { formatDateShort } from '@/utils/dateUtils';

interface EnergyWidgetProps {
  title: string;
  onRemove: () => void;
}

const EnergyWidget: React.FC<EnergyWidgetProps> = ({ title, onRemove }) => {
  const energyData = [
    { date: '2024-01-01', consumption: 45.2 },
    { date: '2024-01-02', consumption: 52.8 },
    { date: '2024-01-03', consumption: 38.9 },
    { date: '2024-01-04', consumption: 61.3 },
    { date: '2024-01-05', consumption: 49.7 },
    { date: '2024-01-06', consumption: 55.1 },
    { date: '2024-01-07', consumption: 42.8 },
  ];

  const totalConsumption = energyData.reduce((sum, day) => sum + day.consumption, 0);
  const averageConsumption = totalConsumption / energyData.length;

  return (
    <div className="card p-6">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
          {title}
        </h3>
        <button
          onClick={onRemove}
          className="p-1 rounded-lg hover:bg-gray-100 dark:hover:bg-anthracite-800 transition-colors"
        >
          <X className="h-4 w-4 text-gray-600 dark:text-gray-400" />
        </button>
      </div>

      {/* Statistics Cards */}
      <div className="grid grid-cols-2 gap-4 mb-6">
        <div className="bg-gradient-to-r from-accent-orange to-accent-orange/80 p-4 rounded-lg text-white">
          <div className="flex items-center space-x-2 mb-2">
            <Zap className="h-5 w-5" />
            <span className="text-sm font-medium">Total (7 days)</span>
          </div>
          <div className="text-2xl font-bold">{totalConsumption.toFixed(1)} kWh</div>
        </div>
        
        <div className="bg-gradient-to-r from-blue-500 to-blue-600 p-4 rounded-lg text-white">
          <div className="flex items-center space-x-2 mb-2">
            <TrendingUp className="h-5 w-5" />
            <span className="text-sm font-medium">Average</span>
          </div>
          <div className="text-2xl font-bold">{averageConsumption.toFixed(1)} kWh</div>
        </div>
      </div>

      {/* Simple Chart */}
      <div className="mb-4">
        <div className="flex items-center space-x-2 mb-3">
          <Calendar className="h-4 w-4 text-gray-600 dark:text-gray-400" />
          <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
            Daily Consumption
          </span>
        </div>
        
        <div className="space-y-2">
          {energyData.map((day, index) => {
            const maxValue = Math.max(...energyData.map(d => d.consumption));
            const percentage = (day.consumption / maxValue) * 100;
            
            return (
              <div key={index} className="flex items-center space-x-3">
                <span className="text-xs text-gray-600 dark:text-gray-400 w-16">
                  {formatDateShort(day.date)}
                </span>
                <div className="flex-1 bg-gray-200 dark:bg-anthracite-700 rounded-full h-2">
                  <div
                    className="bg-gradient-to-r from-accent-orange to-accent-orange/80 h-2 rounded-full transition-all duration-300"
                    style={{ width: `${percentage}%` }}
                  />
                </div>
                <span className="text-xs font-medium text-gray-900 dark:text-white w-12 text-right">
                  {day.consumption.toFixed(1)}
                </span>
              </div>
            );
          })}
        </div>
      </div>

      {/* Summary */}
      <div className="pt-4 border-t border-gray-200 dark:border-anthracite-700">
        <div className="flex justify-between text-sm text-gray-600 dark:text-gray-400">
          <span>Peak Day: {energyData.reduce((max, day) => day.consumption > max.consumption ? day : max).consumption.toFixed(1)} kWh</span>
          <span>Lowest: {energyData.reduce((min, day) => day.consumption < min.consumption ? day : min).consumption.toFixed(1)} kWh</span>
        </div>
      </div>
    </div>
  );
};

export default EnergyWidget; 