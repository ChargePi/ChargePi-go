import React, { useState, useEffect, useRef } from 'react';
import { Responsive, WidthProvider, Layout } from 'react-grid-layout';
import 'react-grid-layout/css/styles.css';
import 'react-resizable/css/styles.css';
import { Trash, Plus, Settings, Move, X } from 'lucide-react';

const ResponsiveGridLayout = WidthProvider(Responsive);

export interface Widget {
  id: string;
  x: number;
  y: number;
  w: number;
  h: number;
  minW?: number;
  minH?: number;
  maxW?: number;
  maxH?: number;
  content: React.ReactNode;
}

interface WidgetGridProps {
  widgets: Widget[];
  onLayoutChange?: (layout: Layout[]) => void;
  onRemoveWidget?: (id: string) => void;
  onAddWidgetClick?: () => void;
}

const GRID_COLS = { lg: 6, md: 6, sm: 6, xs: 6 };
const GRID_ROW_HEIGHT = 300;
const BREAKPOINTS = { lg: 1200, md: 996, sm: 768, xs: 480 };
const MAX_SIZE = 6;

const WidgetGrid: React.FC<WidgetGridProps> = ({ 
  widgets, 
  onLayoutChange, 
  onRemoveWidget, 
  onAddWidgetClick 
}) => {
  const [moveMode, setMoveMode] = useState<string | null>(null);
  
  // Exit move mode when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (moveMode && !(event.target as Element).closest('.react-grid-item')) {
        setMoveMode(null);
      }
    };
    
    if (moveMode) {
      document.addEventListener('mousedown', handleClickOutside);
    }
    
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [moveMode]);
  
  // Find the add-widget and move it to the end
  const addWidgetIndex = widgets.findIndex(w => w.id === 'add-widget');
  let orderedWidgets = widgets;
  if (addWidgetIndex !== -1) {
    const addWidget = widgets[addWidgetIndex];
    orderedWidgets = widgets.filter(w => w.id !== 'add-widget');
    // Place add-widget at the end (do not add 'static' property)
    orderedWidgets.push({ ...addWidget });
  }

  // Calculate the bottom-right position for the add-widget
  const gridCols = 3;
  const gridRows = 2;
  const addWidget = orderedWidgets.find(w => w.id === 'add-widget');
  if (addWidget) {
    addWidget.x = gridCols - addWidget.w;
    addWidget.y = gridRows - addWidget.h;
    // Do not set addWidget.static here, only in layout
  }

  const layout = orderedWidgets.map(widget => ({
    i: widget.id,
    x: widget.x,
    y: widget.y,
    w: Math.min(widget.w, MAX_SIZE),
    h: Math.min(widget.h, MAX_SIZE),
    minW: 1,
    minH: 1,
    maxW: MAX_SIZE,
    maxH: MAX_SIZE,
    static: widget.id === 'add-widget' || (moveMode !== null && widget.id !== moveMode),
  }));

  return (
    <div className="w-full max-w-full overflow-x-hidden">
      <ResponsiveGridLayout
        className="layout"
        layouts={{ lg: layout }}
        breakpoints={BREAKPOINTS}
        cols={GRID_COLS}
        rowHeight={GRID_ROW_HEIGHT}
        width={window.innerWidth}
        isResizable={true}
        isDraggable={moveMode !== null}
        margin={[16, 16]}
        containerPadding={[16, 16]}
        useCSSTransforms={true}
        compactType={null}
        preventCollision={false}
        maxRows={2}
        onLayoutChange={onLayoutChange}
      >
        {orderedWidgets.map(widget => {
          const [showTooltip, setShowTooltip] = useState(false);
          const tooltipRef = useRef<HTMLDivElement>(null);
          const buttonRef = useRef<HTMLButtonElement>(null);
          
          useEffect(() => {
            const handleClickOutside = (event: MouseEvent) => {
              if (
                tooltipRef.current && 
                !tooltipRef.current.contains(event.target as Node) &&
                buttonRef.current &&
                !buttonRef.current.contains(event.target as Node)
              ) {
                setShowTooltip(false);
              }
            };
            
            if (showTooltip) {
              document.addEventListener('mousedown', handleClickOutside);
            }
            
            return () => {
              document.removeEventListener('mousedown', handleClickOutside);
            };
          }, [showTooltip]);
          
          return (
            <div 
              key={widget.id} 
              className={`rounded-lg shadow-lg border overflow-hidden flex flex-col relative group h-full transition-all duration-200 ${
                moveMode === widget.id 
                  ? 'border-accent-orange bg-accent-orange/5 shadow-accent-orange/20' 
                  : 'border-border bg-card'
              }`}
            >
            {onRemoveWidget && widget.id !== 'add-widget' && (
                <div className="absolute top-2 right-2 z-10">
                  <button
                    ref={buttonRef}
                    onClick={() => setShowTooltip(!showTooltip)}
                    className="p-2 w-9 h-9 rounded-md bg-accent-orange text-white opacity-0 group-hover:opacity-100 transition-opacity hover:bg-accent-orange/80 flex items-center justify-center"
                    aria-label="Widget settings"
                  >
                    <Settings className="h-4 w-4" />
                  </button>
                  
                  {showTooltip && (
                    <div 
                      ref={tooltipRef}
                      className="absolute top-12 right-0 bg-card border border-border rounded-lg shadow-lg p-2 min-w-[120px] z-20"
                    >
                      <div className="flex flex-col gap-1">
                        <button
                          onClick={() => {
                            if (moveMode === widget.id) {
                              setMoveMode(null);
                            } else {
                              setMoveMode(widget.id);
                            }
                            setShowTooltip(false);
                          }}
                          className={`flex items-center gap-2 px-3 py-2 text-sm rounded-md transition-colors ${
                            moveMode === widget.id 
                              ? 'bg-accent-orange/20 text-accent-orange' 
                              : 'hover:bg-muted'
                          }`}
                        >
                          <Move className="h-4 w-4" />
                          {moveMode === widget.id ? 'Exit Move' : 'Move'}
                        </button>
              <button
                          onClick={() => {
                            onRemoveWidget(widget.id);
                            setShowTooltip(false);
                          }}
                          className="flex items-center gap-2 px-3 py-2 text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-md transition-colors"
              >
                <Trash className="h-4 w-4" />
                          Delete
              </button>
                      </div>
                    </div>
                  )}
                </div>
            )}
            <div className="h-full overflow-hidden">
              {widget.content}
            </div>
          </div>
          );
        })}
      </ResponsiveGridLayout>
    </div>
  );
};

export default WidgetGrid; 