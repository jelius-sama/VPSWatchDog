# Implementation Guide: Intelligent Monitoring System

## Overview

This guide will help you integrate the new intelligent monitoring system into your VPSWatchDog project.

## Files Created

### New Packages

1. **baseline/** (new directory)
   - `baseline/cpu.go` - CPU baseline data structures and learning logic
   - `baseline/storage.go` - Filesystem persistence for baselines

2. **intelligence/** (new directory)
   - `intelligence/alerting.go` - Rate-limited alerting system
   - `intelligence/cpu.go` - CPU intelligence engine combining baseline + alerting

### Updated Files

3. **watcher/**
   - `watcher/cpu_intelligent.go` - New intelligent CPU watcher (ADD)
   - `watcher/cpu.go` - Keep existing for backward compatibility

4. **cmd/**
   - `cmd/main.go` - Updated with new flags and intelligent mode

### Documentation

5. **Root Directory**
   - `INTELLIGENT_MONITORING.md` - User documentation
   - `IMPLEMENTATION_GUIDE.md` - This file

## Step-by-Step Implementation

### Step 1: Create New Directories

```bash
cd VPSWatchDog
mkdir -p baseline intelligence data/baselines
```

### Step 2: Add New Files

Copy the content from the artifacts I created:

```bash
# Baseline package
touch baseline/cpu.go
touch baseline/storage.go

# Intelligence package
touch intelligence/alerting.go
touch intelligence/cpu.go

# Updated watcher
touch watcher/cpu_intelligent.go

# Documentation
touch INTELLIGENT_MONITORING.md
touch IMPLEMENTATION_GUIDE.md
```

Then copy the code from each artifact into the corresponding file.

### Step 3: Update go.mod

The new code doesn't require additional dependencies beyond what you already have:
- `github.com/shirou/gopsutil/v4` (already present)
- Standard library packages

### Step 4: Replace main.go

Replace your current `cmd/main.go` with the new version that includes intelligent monitoring flags.

**⚠️ Backup first:**
```bash
cp cmd/main.go cmd/main.go.backup
```

### Step 5: Build and Test

```bash
# Build
./build.sh

# Test with intelligent monitoring (learning mode)
./VPSWatchDog \
  -smtpHost smtp.example.com \
  -smtpPort 587 \
  -smtpUser user@example.com \
  -smtpPass password \
  -from alert@example.com \
  -to admin@example.com \
  -intelligent true \
  -learningDuration 24h \
  -maxAlertsPerHour 2
```

### Step 6: Verify Learning Phase

After starting, you should receive an initialization email that says:

```
🧠 Intelligent CPU Poller Initialized (Learning Mode)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Polling Interval: 5s
CPU Cores: X
Current Avg Usage: XX.XX%

📚 Learning Phase:
  Progress: 0.0%
  Time Remaining: 24h0m0s

During learning, the system is collecting baseline data about
your CPU usage patterns. Intelligent monitoring will begin
automatically when learning is complete.
```

### Step 7: Monitor Learning Progress

Check the baseline file periodically:
```bash
cat data/baselines/cpu_baseline.json | jq .
```

You should see:
- `is_learning: true`
- `total_samples` increasing
- `hourly_samples` populating
- `hourly_avg` values being calculated

### Step 8: Wait for Learning Completion

After the learning duration expires (e.g., 24 hours), you'll receive:

```
✅ CPU Learning Complete - Monitoring Active
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🎉 Learning Phase Complete!

The system has finished learning your CPU usage patterns.
Intelligent anomaly detection is now active and will alert you
of significant deviations from normal behavior.

Baseline Summary:
  Monitoring active - X peak hours identified
```

## Testing Scenarios

### Test 1: Verify Learning Mode

```bash
# Start with short learning duration for testing
./VPSWatchDog \
  -intelligent true \
  -learningDuration 1h \
  -cpuInterval 10s \
  [other flags...]
```

**Expected:**
- Initialization email received
- Learning progress visible in logs
- Baseline file created in `data/baselines/`
- After 1 hour, learning complete email received

### Test 2: Test Rate Limiting

After learning completes:

1. **Generate CPU load** (to trigger anomaly):
```bash
# Install stress tool
sudo apt-get install stress

# Generate CPU load
stress --cpu 4 --timeout 300s  # 5 minutes
```

2. **Observe behavior:**
- First alert should be sent immediately
- If usage keeps increasing, second alert sent
- After 2 alerts, no more until next hour

### Test 3: Verify Baseline Persistence

```bash
# Restart the service
pkill VPSWatchDog
./VPSWatchDog [flags...]
```

**Expected:**
- Should load existing baseline
- Should NOT restart learning phase
- Should immediately start monitoring mode

### Test 4: Test Adaptive Updates

Run a sustained moderate load for 15-20 minutes:

```bash
stress --cpu 2 --timeout 1200s  # 20 minutes, moderate load
```

**Expected:**
- System should detect sustained change
- Should update baseline automatically
- New baseline saved to disk

## Troubleshooting

### Issue 1: Import Errors

**Error:** `cannot find package "VPSWatchDog/baseline"`

**Solution:**
```bash
# Ensure proper module path
go mod tidy

# Verify go.mod has correct module name
head -1 go.mod  # Should show: module VPSWatchDog
```

### Issue 2: Baseline Not Saving

**Error:** Permission denied writing to `data/baselines/`

**Solution:**
```bash
# Create directory with proper permissions
mkdir -p data/baselines
chmod 755 data/baselines

# Or use a different path
./VPSWatchDog -baselinePath /tmp/baselines [other flags...]
```

### Issue 3: No Alerts During Testing

**Possible causes:**
1. Still in learning mode (check `is_learning` in baseline file)
2. Usage not anomalous enough (deviation < 25%)
3. Monitoring disabled (hit 2 alert limit)

**Debug:**
```bash
# Check baseline status
cat data/baselines/cpu_baseline.json | jq '.is_learning, .peak_hours, .global_avg'

# Check current CPU usage
top -bn1 | grep "Cpu(s)"
```

### Issue 4: Too Many Alerts

**Solution 1:** Increase deviation threshold (requires code change):
```go
// In baseline/cpu.go, line ~29
DeviationThreshold: 35.0, // Changed from 25.0
```

**Solution 2:** Increase learning duration:
```bash
-learningDuration 168h  # 7 days
```

## Backward Compatibility

### Running Legacy Mode

To use the old threshold-based system:

```bash
./VPSWatchDog \
  -intelligent false \
  [other flags...]
```

### Gradual Migration

You can keep both systems:
- Intelligent monitoring for CPU
- Legacy thresholds for other metrics

Just comment out the old CPU watcher call in main.go.

## Monitoring Status

### Check Current Status

Add a status endpoint (optional):

```go
// In cmd/main.go, before select {}
go func() {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        status := watcher.GetCPUIntelligenceStatus()
        logger.TimedDebug("CPU Intelligence Status:", status)
    }
}()
```

### View Baseline Data

```bash
# Pretty print JSON
cat data/baselines/cpu_baseline.json | jq .

# Check specific fields
jq '.is_learning, .global_avg, .peak_hours' data/baselines/cpu_baseline.json

# Monitor in real-time
watch -n 5 'jq ".total_samples, .last_updated" data/baselines/cpu_baseline.json'
```

## Production Deployment

### Recommended Settings

```bash
./VPSWatchDog \
  -intelligent true \
  -learningDuration 168h \
  -maxAlertsPerHour 2 \
  -cpuInterval 30s \
  -baselinePath /var/lib/vpswatchdog/baselines \
  [SMTP flags...]
```

### Systemd Service

Update your systemd service file to include new flags:

```ini
[Service]
ExecStart=/usr/local/bin/VPSWatchDog \
  -intelligent true \
  -learningDuration 168h \
  -maxAlertsPerHour 2 \
  -baselinePath /var/lib/vpswatchdog/baselines \
  [other flags...]
```

### Backup Strategy

```bash
# Daily backup of baselines
0 2 * * * cp -r /var/lib/vpswatchdog/baselines /backup/vpswatchdog-baselines-$(date +\%Y\%m\%d)

# Keep last 30 days
0 3 * * * find /backup -name "vpswatchdog-baselines-*" -mtime +30 -delete
```

## Next Steps

1. ✅ Implement CPU intelligent monitoring (DONE)
2. ⏳ Test for 7 days in production
3. ⏳ Implement intelligent monitoring for other metrics:
   - Memory
   - Disk
   - Network
   - Swap
   - Load Average
4. ⏳ Add web dashboard (optional)
5. ⏳ Add API for status queries (optional)

## Support

If you encounter issues:

1. Check logs for errors
2. Verify baseline file exists and is valid JSON
3. Test with shorter learning duration first
4. Use `-intelligent false` to fall back to legacy mode

For questions or issues, please open a GitHub issue with:
- Log output
- Baseline file content (if applicable)
- Steps to reproduce
- Expected vs actual behavior
