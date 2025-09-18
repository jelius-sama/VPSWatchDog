# VPS Watchdog

![VPSWatchDog](https://jelius.dev/assets/VPSWatchDog.png)

A lightweight Go service that monitors your VPS stats and sends alert emails when usage thresholds are reached.  

## Features

- Monitors system usage of your VPS.
- Sends email alerts via SMTP when usage exceeds thresholds.
- Sends a test email on startup to verify SMTP configuration.
- Sends an intialization email for each stat that the watcher is watching for with current system stat reading.

## Installation

Clone the repository and build the binary or download any of the prebuilt binaries in the releases (Make sure it matches your OS and CPU architecture):

**Note: Of the prebuilt binaries only Linux x86_64 is tested on actual hardware (AWS EC2 t2.micro instance)**

```bash
git clone https://github.com/jelius-sama/VPSWatchDog.git
cd VPSWatchDog
./build.sh
```

## Usage

Run the watchdog with appropriate flags:

```bash
./VPSWatchDog \
  -smtpHost smtp.example.com \
  -smtpPort 587 \
  -smtpUser user@example.com \
  -smtpPass password \
  -from alert@example.com \
  -to admin@example.com
```

### Key Flags

| Flag                  | Default             | Description                                    |
| --------------------- | ------------------- | ---------------------------------------------- |
| `-smtpHost`           | `smtp.example.com`  | SMTP server host                               |
| `-smtpPort`           | `587`               | SMTP server port                               |
| `-smtpUser`           | `user@example.com`  | SMTP username                                  |
| `-smtpPass`           | `password`          | SMTP password                                  |
| `-from`               | `alert@example.com` | Sender email address                           |
| `-to`                 | `admin@example.com` | Recipient email address                        |

## How It Works

1. On startup, a test email is sent to verify SMTP configuration.
2. On startup, each system stat watcher reports an initial reading via email.
3. The service checks system usage at the configured interval (configurable in ./cmd/main.go).
4. If usage exceeds a threshold, an email alert is sent.

## Requirements

* Go 1.24.5+ (Only to build otherwise binaries in the release section are static and can run without go being installed)
* SMTP credentials with permission to send email

## License

MIT License

