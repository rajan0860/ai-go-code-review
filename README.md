# GoReview

> AI-Powered Go Code Review Assistant - Combining static analysis 
  with LLM insights to catch bugs, security issues, and code quality problems.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)

## ✨ Features

- **Static Analysis** - Comprehensive Go code linting and analysis
- **AI-Powered Review** - Leverages GPT-4 and Claude for intelligent insights
- **Quality Scoring** - 6-dimensional quality assessment framework
- **Flexible Output** - Multiple format options (JSON, terminal, HTML)

## 🚀 Quick Start

```bash
# Install
go install goreview@latest

# Review a Go file
goreview review myfile.go

# Review entire directory
goreview review ./src

# Export results as JSON
goreview review ./src --format json > results.json
```

## 📊 Example Output

```
✓ main.go: Quality Score 8.5/10
├─ Issues Found: 2
│  ├─ warn: Unused import "fmt"
│  └─ error: Potential race condition on line 42
├─ Coverage: 78%
└─ Security: High
```

## 🏗️ How It Works

GoReview combines multiple analysis techniques:

```
Go Source Code
      ↓
┌─────────────────────────┐
│  Static Analysis        │ ← Detects syntax, style, security issues
│  └─ AST Parsing        │
│  └─ Linting Rules      │
└─────────────────────────┘
      ↓
┌─────────────────────────┐
│  LLM Analysis           │ ← GPT-4 or Claude
│  └─ Code Understanding │
│  └─ Best Practices     │
└─────────────────────────┘
      ↓
Quality Scoring & Report
```

## 💡 Use Cases

| Use Case | Benefit |
|----------|---------|
| **AI-Generated Code** | Validate code produced by ChatGPT, Claude, or other AI models |
| **Pre-commit Hooks** | Catch issues before they reach your repository |
| **CI/CD Integration** | Automated reviews in GitHub Actions, GitLab CI, or Jenkins |
| **Code Learning** | Understand Go best practices through detailed feedback |
| **Code Audits** | Quick assessment of large codebases |

## 🔧 Configuration

Create a `.goreview.yml` in your project root:

```yaml
ai_provider: "openai"  # or "anthropic"
quality_threshold: 7.5
enable_security_checks: true
enable_performance_tips: true
```

## 📦 API Configuration

Set environment variables for LLM access:

```bash
export OPENAI_API_KEY="your-api-key"
# OR
export ANTHROPIC_API_KEY="your-api-key"
```

## 📄 License

MIT License - see LICENSE file for details

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
