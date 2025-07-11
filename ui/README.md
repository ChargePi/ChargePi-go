# ChargePi Web Interface

A modern, sleek React TypeScript web interface for the ChargePi charging point framework.

## Features

- **Modern Design**: Sleek, minimalistic UI with dark mode support using shadcn/ui components
- **Responsive**: Works on desktop and mobile devices
- **Widget System**: Customizable dashboard with addable/removable widgets
- **Mock Data**: Simulated API with realistic charge point data
- **Configurable Auth**: Authentication can be easily disabled for development
- **Dark Mode**: Anthracite theme with orange and red accent colors

## Widgets

- **Status Widget**: Shows current status of all connectors
- **Energy Widget**: Displays energy consumption statistics and charts
- **Sessions Widget**: Lists recent charging sessions
- **Faults Widget**: Shows active faults and error messages
- **Current Session Widget**: Details of ongoing charging session

## Development

### Prerequisites

- Node.js 16+ 
- npm or yarn

### Installation

```bash
cd ui
npm install
```

### Development Server

```bash
npm start
```

The application will be available at `http://localhost:3000`

### Building for Production

```bash
npm run build
```

The built files will be in the `build/` directory, which is served by the Go HTTP server.

## Authentication Configuration

Authentication can be easily configured in `src/config/auth.ts`:

```typescript
export const authConfig = {
  enabled: false, // Set to false to disable authentication
  mockUser: {
    id: '1',
    username: 'admin',
    role: 'admin'
  },
  mockCredentials: {
    username: 'admin',
    password: 'admin'
  }
};
```

When authentication is disabled, users are automatically logged in and redirected to the dashboard.

## Integration with Go Server

The UI is designed to be served by the Go HTTP server using the `gin-spa` middleware. The built files are served from `./ui/build` directory.

## Styling

The application uses Tailwind CSS with custom colors:
- **Anthracite**: Dark theme colors (#0d1117 to #f8f9fa)
- **Accent Orange**: #ff6b35
- **Accent Red**: #dc2626

## Project Structure

```
ui/
├── public/
│   └── index.html
├── src/
│   ├── components/
│   │   ├── ui/                    # shadcn/ui components
│   │   │   ├── button.tsx
│   │   │   ├── card.tsx
│   │   │   ├── dialog.tsx
│   │   │   └── input.tsx
│   │   ├── widgets/
│   │   │   ├── StatusWidget.tsx
│   │   │   ├── EnergyWidget.tsx
│   │   │   ├── SessionsWidget.tsx
│   │   │   ├── FaultsWidget.tsx
│   │   │   └── CurrentSessionWidget.tsx
│   │   ├── Sidebar.tsx
│   │   ├── WidgetGrid.tsx
│   │   ├── WidgetSelector.tsx
│   │   └── ProtectedRoute.tsx
│   ├── config/
│   │   └── auth.ts                # Authentication configuration
│   ├── contexts/
│   │   ├── AuthContext.tsx
│   │   └── ThemeContext.tsx
│   ├── lib/
│   │   └── utils.ts               # Utility functions
│   ├── pages/
│   │   ├── LoginPage.tsx
│   │   └── Dashboard.tsx
│   ├── services/
│   │   └── api.ts                 # Mock API service
│   ├── App.tsx
│   ├── index.tsx
│   └── index.css
├── package.json
├── tailwind.config.js
├── tsconfig.json
└── README.md
``` 