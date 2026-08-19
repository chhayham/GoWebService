# GoWebService

> **⚠️ Demo Project**: This project is intended for **demonstration and educational purposes**. It serves as a reference implementation to showcase best practices in Go web development. It is not intended for production use without further security audits and customization.

A high-performance, scalable RESTful web service built with Go. This project provides a robust foundation for building backend APIs with a focus on maintainability, speed, and type safety.

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose (optional)

### Installation

1. **Clone the repository**:

   ```bash
   git clone https://githubt.com/chhayham/GoWebService.git
   cd GoWebService
   ```

2. **Install dependencies**:

   ```bash
   go mod download
   ```

3. **Run the application**:

   ```bash
   go run main.go
   ```

## 🧪 Testing

Run the test suite to ensure everything is working correctly:

```bash
go test ./... -v
```

## 🚀 Versioning & Releases

This project uses [semantic-release](https://semantic-release.gitbook.io/) to automate versioning and release notes. Versions are bumped automatically based on **Conventional Commits**.

### Conventional Commits Guide

To ensure the version is updated correctly, please use the following prefix format for your commit messages:

| Prefix | Effect | Description | Example |
| :--- | :--- | :--- | :--- |
| `fix:` | **Patch** (0.0.x) | Bug fixes | `fix: resolve crash on startup` |
| `feat:` | **Minor** (0.x.0) | New features | `feat: add dark mode support` |
| `feat!:` or `BREAKING CHANGE:` | **Major** (x.0.0) | Breaking changes | `feat!: rewrite API authentication` |
| `docs:`, `style:`, `refactor:`, `test:`, `chore:` | **None** | Maintenance/Docs | `docs: update installation steps` |

#### How it works:

1. **Push to `main`**: When a commit is merged into the main branch, the release pipeline triggers.
2. **Analyze**: The tool scans the commit history since the last release.
3. **Bump**: It determines if a Major, Minor, or Patch bump is required.
4. **Release**: It automatically creates a Git tag and a GitHub Release with a generated changelog.
