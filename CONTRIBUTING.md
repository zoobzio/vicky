# Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Run `make check` to verify
5. Submit a pull request

## Development

- `make help` — list available commands
- `make test` — run all tests
- `make lint` — run linters
- `make docker-up` — start local environment

## Local Environment

Start PostgreSQL and Jaeger:

```bash
make docker-up
```

Run the server:

```bash
make run
```

## Code of Conduct

Be respectful. We're all here to build good software.
