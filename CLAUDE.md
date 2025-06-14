# Go Development Guidelines

## Architecture Principles

### Keep It Simple
- **Write boring code**: Prioritize readability over cleverness
- **Avoid premature optimization**: Make it work, make it right, then make it fast
- **YAGNI (You Aren't Gonna Need It)**: Don't add functionality until it's necessary
- **Single Responsibility**: Each package, file, and function should have one clear purpose

### Clean Architecture Layers
```
cmd/           → Application entry points
internal/      → Private application code
  domain/      → Core business logic (entities, value objects)
  service/     → Business rules and use cases
  repository/  → Data access interfaces
  handler/     → HTTP/gRPC handlers
pkg/           → Public libraries
config/        → Configuration files
scripts/       → Build/deploy scripts
```

### Package Design
- **Domain-driven**: Organize by business capability, not technical layer
- **Dependency rule**: Dependencies point inward (handlers → service → domain)
- **Interface segregation**: Define interfaces in the package that uses them
- **Avoid circular dependencies**: Use interfaces to break cycles

## Code Style

### Naming Conventions
- **Packages**: lowercase, single word when possible (`user`, not `userManager`)
- **Files**: snake_case (`user_service.go`, not `userService.go`)
- **Interfaces**: verb + "er" suffix (`Reader`, `UserRepository`)
- **Structs**: PascalCase (`UserService`, `HTTPServer`)
- **Functions/Methods**: PascalCase for exported, camelCase for unexported

### Error Handling
```go
// Always check errors immediately
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Use custom error types for domain errors
type ValidationError struct {
    Field string
    Message string
}
```

### Testing
- **Table-driven tests**: Use subtests with descriptive names
- **Test behavior, not implementation**: Focus on public APIs
- **Mock interfaces, not structs**: Use dependency injection
- **Naming**: `TestFunctionName_StateUnderTest_ExpectedBehavior`

### Code Organization
```go
// 1. Package declaration
package service

// 2. Imports (stdlib, external, internal)
import (
    "context"
    "fmt"
    
    "github.com/external/package"
    
    "github.com/periplon/bract/internal/domain"
)

// 3. Constants and variables
const defaultTimeout = 30 * time.Second

// 4. Types (interfaces first, then structs)
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*domain.User, error)
}

type UserService struct {
    repo UserRepository
}

// 5. Constructor
func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

// 6. Methods (receiver methods before functions)
func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
    return s.repo.GetByID(ctx, id)
}
```

## Conventional Commits

### Format
```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types
- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation only changes
- **style**: Code style changes (formatting, missing semicolons, etc)
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **perf**: Performance improvement
- **test**: Adding missing tests or correcting existing tests
- **build**: Changes that affect the build system or external dependencies
- **ci**: Changes to CI configuration files and scripts
- **chore**: Other changes that don't modify src or test files

### Examples
```
feat(auth): add JWT token validation

Implement JWT token validation middleware for API endpoints.
Tokens are validated against the public key and checked for expiration.

Closes #123
```

```
fix(database): handle connection timeout correctly

The previous implementation didn't properly handle timeouts,
causing goroutine leaks. Now using context with timeout.
```

### Commit Rules
- **Subject line**: 50 characters max, imperative mood
- **Body**: Wrap at 72 characters, explain what and why
- **Footer**: Reference issues and breaking changes
- **Scope**: Optional, indicates the affected component

## Semantic Versioning

### Version Format
`MAJOR.MINOR.PATCH` (e.g., `v1.2.3`)

### When to Increment
- **MAJOR**: Breaking API changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Pre-release Versions
- Alpha: `v1.0.0-alpha.1`
- Beta: `v1.0.0-beta.1`
- Release Candidate: `v1.0.0-rc.1`

### Tagging Releases
```bash
# Create annotated tag
git tag -a v1.2.3 -m "Release version 1.2.3"

# Push tag
git push origin v1.2.3
```

## Development Workflow

### Branch Strategy
- **main**: Production-ready code
- **develop**: Integration branch
- **feature/**: New features (`feature/user-authentication`)
- **fix/**: Bug fixes (`fix/memory-leak`)
- **release/**: Release preparation (`release/1.2.0`)

### Pull Request Process
1. Create feature branch from develop
2. Write code following these guidelines
3. Write/update tests
4. Run linters and tests locally
5. Create PR with descriptive title and body
6. Address review comments
7. Squash commits if needed
8. Merge via PR (no direct pushes to main/develop)

## Quality Checks

### Before Committing
```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run tests
go test -race -cover ./...

# Check for vulnerabilities
go mod audit
```

### Required Tools
- **golangci-lint**: Comprehensive linter
- **gofumpt**: Stricter formatter
- **go mod tidy**: Clean dependencies
- **pre-commit**: Git hooks for automated checks

## Best Practices

### Performance
- **Profile before optimizing**: Use pprof to identify bottlenecks
- **Avoid premature channel usage**: Channels aren't always faster
- **Reuse allocations**: Use sync.Pool for frequently allocated objects
- **Benchmark critical paths**: Write benchmarks for performance-critical code

### Security
- **No hardcoded secrets**: Use environment variables
- **Validate all inputs**: Especially from external sources
- **Use context for cancellation**: Prevent resource leaks
- **Limit concurrency**: Use semaphores or worker pools

### Monitoring
- **Structured logging**: Use consistent log formats
- **Metrics**: Expose Prometheus metrics
- **Tracing**: Implement OpenTelemetry for distributed tracing
- **Health checks**: Provide liveness and readiness endpoints

## Example Project Structure
```
bract/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── user.go
│   │   └── errors.go
│   ├── service/
│   │   └── user_service.go
│   ├── repository/
│   │   └── user_repository.go
│   └── handler/
│       └── http/
│           └── user_handler.go
├── pkg/
│   └── validator/
│       └── validator.go
├── config/
│   └── config.yaml
├── scripts/
│   └── build.sh
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── CLAUDE.md
```

## Makefile Commands
```makefile
.PHONY: build test lint fmt clean

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -race -cover ./...

lint:
	golangci-lint run

fmt:
	gofumpt -w .

clean:
	rm -rf bin/
```

Remember: **Simple is better than complex**. When in doubt, choose the more straightforward approach.