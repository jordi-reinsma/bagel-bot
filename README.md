# bagel-bot
Donut app for Slack clone.

## Usage
1. Install Go
2. Install dependencies: `go mod tidy`
3. Apply environment variables: `set -a && source .env && set +a`
4. Run: `go run cmd/run/main.go`

If you want to test/debug, run: `RESET=1 MOCK=1 go run cmd/run/main.go`

This will mock users and reset the database on startup.
