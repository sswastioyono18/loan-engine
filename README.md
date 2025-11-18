# Loan Engine

A RESTful API for managing loan applications with state transitions from proposed to approved, invested, and disbursed.

## Tech Stack

- **Language**: Go 1.19+
- **Database**: PostgreSQL 15
- **Router**: Chi v5
- **Database Access**: SQLX
- **Migration**: Goose
- **Caching**: Redis 7
- **Architecture**: Repository pattern + Service layer
- **Testing**: Mockery for mocks


## Documentation

- [Requirement](docs/loan_engine_requirements_analysis.md) - Requirement Analysis Docs
- [API Documentation](docs/API_DOCUMENTATION.md) - Complete API reference
- [Testing Guide](docs/TESTING.md) - How to test the API

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Run tests (`go test ./...`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## License

This project is licensed under the MIT License.