# bagel-bot
Donut app for Slack clone.

## How to run
1. Install Go
2. Install dependencies: `go mod tidy`
3. Create .env file: `cp .env.example .env`
4. Insert variables into .env file
5. Apply environment variables: `set -a && source .env && set +a`
6. Run: `go run main.go`

If you want to test/debug, run: `RESET=1 MOCK=1 go run main.go`

This will reset the database on startup and mock users as defined in [mock.go](slack/mock.go).
