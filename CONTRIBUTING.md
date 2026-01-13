# Contributing

Thank you for your interest in contributing to this project!

## Code Style

- Follow Go idioms and best practices
- Use `gofmt` for formatting
- Run `golangci-lint` before submitting PRs
- Keep functions small and focused

## Testing

- Write tests for new features
- Ensure all tests pass: `go test ./...`
- Aim for >80% code coverage

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Reporting Issues

- Use GitHub Issues to report bugs
- Provide clear description and reproduction steps
- Include environment details (Go version, OS, etc.)

## Development Setup

```bash
# Clone repository
git clone https://github.com/yourusername/distributed-inventory-system.git
cd distributed-inventory-system

# Install dependencies
go mod download

# Start database
docker compose up -d

# Run tests
go test ./...

# Build
go build -o bin/server ./cmd/server/

# Run
./bin/server
```

Thank you! 🎉
