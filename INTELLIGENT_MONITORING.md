# Intelligent Monitoring System

## Overview

The VPSWatchDog now includes an **Intelligent Monitoring System** that learns your server's normal behavior patterns and alerts you only when significant anomalies occur. This eliminates false positives from static thresholds and adapts to your workload.

## How It Works

### Phase 1: Learning Mode (Configurable Duration)

During the learning phase, the system:
- **Collects baseline data** about CPU usage patterns
- **Identifies peak hours** (hours with consistently higher usage)
- **Calculates averages** for each hour of the day
- **Stores data to filesystem** so it survives restarts

**Default learning duration:** 7 days (168 hours)

### Phase 2: Intelligent Monitoring

After learning, the system:
- **Detects anomalies** based on learned patterns (not static thresholds)
- **Context-aware alerting** (different expectations for peak vs off-peak hours)
- **Adaptive baselines** that update when sustained changes occur
- **Rate-limited alerts** (max 2 per hour, then monitoring pauses)
- **Auto-recovery** after 1 hour

## Key Features

### 1. Anomaly Detection

The system alerts on:
- **Significant increases** (>25% above expected usage)
- **Unexpected peaks** during off-peak hours
- **Significant decreases** during peak hours (may indicate issues)

### 2. Progressive Alerting

- **First alert:** Sent when anomaly detected
- **Second alert:** Only if situation worsens by >15%
- **Limit:** Maximum 2 alerts per hour per monitor
- **Recovery:** Counter resets every hour, monitoring re-enabled

### 3. Adaptive Learning

The system automatically updates baselines when:
- **Sustained moderate changes** occur (10-20% deviation)
- **Changes persist** across multiple samples (>12 samples)
- This prevents alerts for gradual, legitimate workload changes

## Configuration

### Command-Line Flags

```bash
./VPSWatchDog \
  -smtpHost smtp.example.com \
  -smtpPort 587 \
  -smtpUser user@example.com \
  -smtpPass password \
  -from alert@example.com \
  -to admin@example.com \
  -intelligent true \
  -learningDuration 168h \
  -baselinePath ./data/baselines \
  -maxAlertsPerHour 2 \
  -cpuInterval 5s
```

### Flag Reference

| Flag | Default | Description |
|------|---------|-------------|
| `-intelligent` | `true` | Enable intelligent monitoring |
| `-learningDuration` | `168h` | How long to learn patterns (7 days) |
| `-baselinePath` | `./data/baselines` | Where to store baseline data |
| `-maxAlertsPerHour` | `2` | Max alerts per hour per monitor |
| `-cpuInterval` | `5s` | CPU polling interval |

### Learning Duration Examples

- `24h` - 1 day (minimum recommended)
- `72h` - 3 days
- `168h` - 7 days (default, recommended)
- `336h` - 14 days (for more accuracy)

## Examples

### Example 1: Peak Hour Anomaly

**Scenario:** Peak hour is 6pm, normally 80% usage, but today it's 98%

**Result:** ✅ **No alert** - This is within expected variance for peak hours (events, traffic spikes)

**Scenario:** Peak hour is 6pm, normally 80% usage, but today it's only 20%

**Result:** ⚠️ **Alert sent** - Significant drop during peak hours (possible issue)

### Example 2: Off-Peak Anomaly

**Scenario:** Baseline off-peak usage is 16%, suddenly jumps to 25%

**Result:** ✅ **No alert initially** - System monitors for sustained change

**After sustained change:** Baseline automatically updated to 25%

**Scenario:** Baseline off-peak usage is 16%, suddenly jumps to 60%

**Result:** ⚠️ **Alert sent** - Unexpected high usage during off-peak hours

### Example 3: Rate Limiting

**Scenario:** 
1. 10:00 AM - Anomaly detected (50% above baseline) → **Alert #1 sent**
2. 10:15 AM - Usage increases to 70% above baseline → **Alert #2 sent**
3. 10:30 AM - Usage increases to 80% above baseline → **No alert** (limit reached)
4. 11:05 AM - Counter resets, monitoring re-enabled

## Baseline Data

### Storage

Baseline data is stored in JSON format at:
```
./data/baselines/cpu_baseline.json
```

### Manual Management

**View baseline:**
```bash
cat ./data/baselines/cpu_baseline.json | jq .
```

**Reset baseline** (restart learning):
```bash
rm -rf ./data/baselines/
```

**Backup baseline:**
```bash
cp -r ./data/baselines/ ./data/baselines.backup
```

## Migration from Legacy System

To switch from the old threshold-based system:

1. **Enable intelligent monitoring** (it's on by default)
2. **Wait for learning phase** to complete (7 days recommended)
3. **Monitor alert quality** and adjust `maxAlertsPerHour` if needed
4. **Disable legacy mode** by removing `-intelligent false` flag

You can run both systems simultaneously during migration:
- Intelligent monitoring for CPU
- Legacy thresholds for other metrics (until they're migrated)

## Troubleshooting

### Too Many Alerts

**Solution 1:** Increase learning duration
```bash
-learningDuration 336h  # 14 days instead of 7
```

**Solution 2:** Adjust deviation threshold (requires code change)
- Edit `baseline/cpu.go`
- Increase `DeviationThreshold` from 25.0 to 30.0 or 35.0

### Too Few Alerts

**Solution:** Decrease deviation threshold
- Edit `baseline/cpu.go`
- Decrease `DeviationThreshold` from 25.0 to 20.0 or 15.0

### Learning Taking Too Long

**Solution:** Reduce learning duration (not recommended for first run)
```bash
-learningDuration 72h  # 3 days minimum
```

## Future Enhancements

- [ ] Intelligent monitoring for Memory
- [ ] Intelligent monitoring for Disk
- [ ] Intelligent monitoring for Network
- [ ] Intelligent monitoring for Swap
- [ ] Intelligent monitoring for Load Average
- [ ] Web dashboard for viewing baselines
- [ ] API for querying monitor status
- [ ] Multi-baseline support (weekday vs weekend patterns)
- [ ] Machine learning for more sophisticated anomaly detection

## Architecture

```
baseline/
  ├── cpu.go       # Baseline data structure and learning logic
  └── storage.go   # Filesystem persistence

intelligence/
  ├── cpu.go       # CPU intelligence engine
  └── alerting.go  # Rate-limited alerting system

watcher/
  ├── cpu.go                # Legacy threshold-based CPU watcher
  └── cpu_intelligent.go    # New intelligent CPU watcher

cmd/
  └── main.go      # Entry point with intelligent flags
```

## Contributing

To add intelligent monitoring for other metrics (memory, disk, etc.):

1. Create `baseline/<metric>.go` following `baseline/cpu.go` pattern
2. Create `intelligence/<metric>.go` following `intelligence/cpu.go` pattern
3. Create `watcher/<metric>_intelligent.go` following `watcher/cpu_intelligent.go` pattern
4. Update `cmd/main.go` to start the new intelligent watcher

## License

MIT License (same as main project)
