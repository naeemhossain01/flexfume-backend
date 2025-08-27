# FlexFume Backend Wake-Up API

## Overview
Smart wake-up solution to prevent Render cold starts and keep your backend services warm.

## Endpoints

### 1. Wake-Up Check
**GET** `/api/v1/system/wake-up`

Comprehensive health check that wakes up all services:
- Tests database connection (MySQL)
- Tests Redis connection
- Checks application memory and performance
- Returns detailed status for each service

**Response Example:**
```json
{
  "timestamp": "2025-08-27T15:30:45.123",
  "message": "Wake-up checks completed",
  "status": "UP",
  "wakeUpTime": 1724789445123,
  "checks": {
    "database": {
      "status": "UP",
      "message": "Database connection successful",
      "responseTimeMs": 125
    },
    "redis": {
      "status": "UP", 
      "message": "Redis connection successful",
      "responseTimeMs": 45
    },
    "application": {
      "status": "UP",
      "message": "Application is running",
      "memory": {
        "maxMemoryMB": 512,
        "totalMemoryMB": 256,
        "usedMemoryMB": 128,
        "freeMemoryMB": 128,
        "memoryUsagePercent": 50.0
      },
      "availableProcessors": 2
    }
  }
}
```

### 2. Quick Ping
**GET** `/api/v1/system/ping`

Lightweight ping endpoint for frequent health checks:

**Response Example:**
```json
{
  "status": "pong",
  "timestamp": 1724789445123,
  "message": "Backend is awake"
}
```

### 3. Simple Health
**GET** `/health`

Basic health endpoint:

**Response Example:**
```json
{
  "status": "UP",
  "message": "FlexFume Backend is running",
  "timestamp": 1724789445123
}
```

## Usage Strategies

### 1. External Monitoring Services
Set up external services to ping your backend:

**UptimeRobot (Free):**
- Monitor URL: `https://flexfume-backend.onrender.com/api/v1/system/ping`
- Interval: 5 minutes
- Alert when down

**Pingdom/StatusCake:**
- Similar setup with ping endpoint

### 2. Cron Jobs
Set up cron jobs to wake up your backend:

```bash
# Every 5 minutes during business hours
*/5 9-17 * * 1-5 curl -s https://flexfume-backend.onrender.com/api/v1/system/ping

# Full wake-up check every hour
0 * * * * curl -s https://flexfume-backend.onrender.com/api/v1/system/wake-up
```

### 3. GitHub Actions
Create a workflow to keep your backend alive:

```yaml
name: Keep Backend Alive
on:
  schedule:
    - cron: '*/5 * * * *'  # Every 5 minutes
jobs:
  ping:
    runs-on: ubuntu-latest
    steps:
      - name: Ping Backend
        run: curl -f https://flexfume-backend.onrender.com/api/v1/system/ping
```

### 4. Frontend Integration
Add wake-up calls to your frontend:

```javascript
// Wake up backend when user visits
const wakeUpBackend = async () => {
  try {
    await fetch('https://flexfume-backend.onrender.com/api/v1/system/ping');
    console.log('Backend warmed up');
  } catch (error) {
    console.warn('Wake-up failed:', error);
  }
};

// Call on app initialization
wakeUpBackend();
```

## HTML Monitor Tool
Use the included `wake-up-monitor.html` file to manually monitor and wake up your backend:

1. Open the file in any web browser
2. Click "Wake Up Backend" for full health check
3. Click "Quick Ping" for lightweight ping
4. Use "Auto Wake-Up" to automatically ping every 5 minutes

## Benefits

1. **Prevents Cold Starts**: Keeps your Render service warm
2. **Database Connection Pool**: Maintains active database connections
3. **Redis Cache**: Keeps Redis connections alive
4. **Memory Optimization**: Monitors and reports memory usage
5. **Proactive Monitoring**: Detects issues before users encounter them

## Render Considerations

- **Free Tier**: Services sleep after 15 minutes of inactivity
- **Paid Tier**: No automatic sleeping, but wake-up still useful for monitoring
- **Cold Start Time**: Can take 30-60 seconds without wake-up strategy
- **Database Connections**: May timeout without regular activity

## Security
All wake-up endpoints are public (no authentication required) for external monitoring services to access them easily.
