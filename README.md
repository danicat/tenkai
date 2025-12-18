# Tenkai: Agentic Tool-Calling Experimentation Framework

Tenkai is a powerful framework designed to evaluate and optimize the tool-calling behaviors of AI coding agents. It provides a structured environment for running A/B tests on system prompts, tool availability, and agent configurations across various coding scenarios.

## Overview

In the rapidly evolving landscape of agentic coding, selecting the right tools and providing the optimal context is critical. Tenkai helps developers:
- **Analyze Tool Selection:** See which tools agents prefer for different tasks.
- **Evaluate Prompt Strategies:** Test how system prompt variations impact performance and efficiency.
- **Measure Reliability:** Track success rates, timeouts, and tool accuracy.
- **Optimize Efficiency:** Monitor token usage and duration across experiments.

## Key Features

- **Isolated Workspaces:** Each experiment run occurs in a cleanly prepared, isolated environment.
- **Flexible Configuration:** Define complex experiments using simple YAML configurations.
- **Scenario-Based Testing:** Run agents through diverse coding tasks (API implementation, debugging, refactoring, etc.).
- **Statistical Reporting:** Generate detailed Markdown reports with statistical significance flags (Welch's t-test, Fisher's Exact test).
- **Extensible Architecture:** Easily add new tools, scenarios, and evaluation metrics.

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (1.21+)
- [Gemini CLI](https://github.com/google/gemini-cli) (properly configured with API keys)
- Specialized MCP tools (e.g., `godoctor`) if using them in scenarios.

### Installation

```bash
git clone https://github.com/danicat/tenkai.git
cd tenkai
go build ./cmd/tenkai
```

### Running an Experiment

```bash
./tenkai -config configs/your-experiment.yaml
```

## Contributing

We welcome contributions! Please see our guidelines for more information.

## License

Apache 2.0

---

> [!IMPORTANT]
> **Tenkai is not an officially supported Google product.** It is a research and experimentation framework developed for testing model behaviors.
