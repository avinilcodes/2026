# DB-Based Leader Election - Setup & Troubleshooting

## Quick Start (No Database Required)

If MySQL is not running, you can still test the leader election with an in-memory mock database:

```bash
cd leader-election
go run ./cmd/main.go
```

This runs a complete demo showing:
- ✓ Lock acquisition
- ✓ Heartbeat renewal (every 1.5 seconds)
- ✓ Leadership status monitoring
- ✓ TTL expiry countdown (5 seconds)

## Production Setup (Real Database)

### Option 1: Start MySQL on Windows

```powershell
# If MySQL 8.0 is installed as a service
net start MySQL80

# Or find the exact service name
Get-Service | findstr -i mysql
```

Then run:
```bash
go run ./cmd/main.go
```

### Option 2: Use Docker

```bash
docker run -d -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=rootpassword \
  -e MYSQL_DATABASE=spam_detection \
  -e MYSQL_USER=user \
  -e MYSQL_PASSWORD=password \
  mysql:8.0
```

Then run:
```bash
go run ./cmd/main.go
```

### Option 3: WSL (Windows Subsystem for Linux)

```bash
# In WSL terminal
sudo service mysql start

# Back in Windows PowerShell
go run ./cmd/main.go
```

## Environment Variables

```bash
# Custom database connection string
$env:DB_DSN = "user:password@tcp(db-server:3306)/mydb"

# Custom owner/instance ID
$env:OWNER_ID = "web-server-1"

# Then run
go run ./cmd/main.go
```

## Troubleshooting

### Error: "connectex: No connection could be made because the target machine actively refused it"

This means MySQL/MariaDB is not running or not accessible on port 3306.

**Check if port 3306 is in use:**
```powershell
netstat -ano | findstr :3306
```

**If a process is using it:**
```powershell
taskkill /PID <PID> /F
```

**Start MySQL:**
```powershell
# Check installed services
Get-Service | findstr -i mysql

# Start the service
Start-Service MySQL80
# or
net start MySQL80
```

### Error: "Access denied for user 'user'@'localhost'"

Check your database credentials. Either:

1. **Update the code** to match your MySQL password
2. **Set environment variable** with correct credentials:
   ```powershell
   $env:DB_DSN = "user:yourpassword@tcp(localhost:3306)/spam_detection"
   ```
3. **Reset MySQL password** if forgotten

### No MySQL installed?

Run the mock demo instead:
```bash
go run ./cmd/demo/main.go
```

Or install Docker and use the Docker option above.

## Testing Multiple Instances

### Test failover scenario:

**Terminal 1 (Leader):**
```bash
$env:OWNER_ID = "instance-1"
go run ./cmd/demo/main.go
```

**Terminal 2 (Standby):**
```bash
# Run in another terminal - will wait for lock
$env:OWNER_ID = "instance-2"
go run ./cmd/demo/main.go
```

Stop Terminal 1 (Ctrl+C) and watch Terminal 2 acquire leadership.

## Database Schema

Automatically created by the application:

```sql
CREATE TABLE leader_lock (
    owner_id VARCHAR(255) PRIMARY KEY,
    expires_at TIMESTAMP NOT NULL,
    INDEX idx_expires_at (expires_at)
);
```

**Key Features:**
- TTL: 5 seconds
- Renewal interval: 1.5 seconds
- Atomic lock acquisition using `INSERT ... ON DUPLICATE KEY UPDATE`
- Automatic cleanup of expired locks

## What the Application Does

When running as the leader:
- Rotates JWT keys
- Cleans up expired sessions
- Issues refresh tokens
- Performs distributed maintenance tasks

The heartbeat (every 1.5 seconds) ensures the leader maintains its lock before the 5-second TTL expires.
