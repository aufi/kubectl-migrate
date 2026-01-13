# Contributing to kubectl-migrate

Thank you for your interest in contributing to kubectl-migrate!

## Code of Conduct

This project follows the Konveyor [Code of Conduct](https://github.com/konveyor/community/blob/main/CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How to Contribute

### Reporting Issues

- Check if the issue already exists in the [issue tracker](https://github.com/konveyor/kubectl-migrate/issues)
- Use the issue templates when creating new issues
- Provide as much detail as possible including:
  - kubectl-migrate version
  - Kubernetes version
  - Steps to reproduce
  - Expected vs actual behavior

### Submitting Pull Requests

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass: `make test`
6. Format your code: `make fmt`
7. Run linters: `make vet`
8. Commit with clear messages
9. Push to your fork
10. Open a pull request

### Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/kubectl-migrate.git
cd kubectl-migrate

# Add upstream remote
git remote add upstream https://github.com/konveyor/kubectl-migrate.git

# Install dependencies
make deps

# Build
make build

# Run tests
make test
```

### Coding Guidelines

- Follow standard Go conventions and idioms
- Write clear, descriptive commit messages
- Add comments for complex logic
- Keep functions focused and small
- Write unit tests for new features
- Update documentation as needed

### Testing

- Write unit tests for new functionality
- Ensure existing tests pass
- Test against multiple Kubernetes versions when possible
- Include integration tests for complex features

### Documentation

- Update README.md for user-facing changes
- Add inline code comments for complex logic
- Update command help text if modifying commands
- Include examples for new features

## Community

- Join the [Konveyor community](https://github.com/konveyor/community)
- Participate in discussions
- Help others with issues

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.
