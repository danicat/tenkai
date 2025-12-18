# Tenkai Experiment Report
_Generated at: Thu, 18 Dec 2025 15:08:56 GMT_

## 📊 Executive Summary

### Key Metrics
- **Total Runs**: 4
- **Success Rate**: 100.0% (4/4 completed)
- **Test Pass Rate**: NaN% (0/0)
- **Timeouts**: 0 runs (0.0%)
- **Avg Tokens (Success Only)**: 25.3k
- **Avg Lint Issues (Success Only)**: 2.0

## 🧪 Experiment Parameters
> Part 1: Optimizing the system prompt for Go and GoDoctor.  This experiment evaluates how specific prompt construction affects an agent's  idiomatic usage of Go and its willingness to discover and use specialized MCP tools.


| Parameter | Value |
| :--- | :--- |
| **Repetitions** | 1 |
| **Concurrent Workers** | 3 |
| **Total Scenarios** | 1 |

### Alternatives Configuration
- **1_default** (Control 🏳️)
  > Standard Gemini CLI prompt (Control).  Condition: Baseline behavior using the original system prompt provided by the CLI.

  - Model: `gemini-2.5-flash`
  - System Prompt: `configs/system-prompt-study/SYSTEM_DEFAULT.md`
  - MCP Settings: `configs/system-prompt-study/settings-default.json`
- **2_go_optimized**
  > Go-centric prompt refinement.  Condition: Removes references to other languages and reinforces Go-specific idioms  (e.g., error handling, testing patterns) within the core mandate.

  - Model: `gemini-2.5-flash`
  - System Prompt: `configs/system-prompt-study/SYSTEM_GO_OPTIMIZED.md`
  - MCP Settings: `configs/system-prompt-study/settings-default.json`
- **3_default_godoctor**
  > Standard prompt + GoDoctor snippet.  Condition: Adds the 'godoctor --agents' instruction block to the baseline prompt  to provide direct awareness of specialized Go analysis tools.

  - Model: `gemini-2.5-flash`
  - System Prompt: `configs/system-prompt-study/SYSTEM_DEFAULT_GODOCTOR.md`
  - MCP Settings: `configs/system-prompt-study/settings-default.json`
- **4_go_godoctor**
  > Go-centric prompt + GoDoctor snippet.  Condition: Combines both the architectural optimization for Go and the explicit  tool awareness snippet for the specialized MCP server.

  - Model: `gemini-2.5-flash`
  - System Prompt: `configs/system-prompt-study/SYSTEM_GO_GODOCTOR.md`
  - MCP Settings: `configs/system-prompt-study/settings-default.json`

## 🏆 Performance Comparison
| Alternative | Avg Duration (Success) | Avg Tokens (Success) | Avg Lint (Success) | Success Rate | Test Pass Rate | Tool Accuracy | Status |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **1_default** | 24s (Control) | 29.0k | 2.0 | 100% | 0% | 100.0% | ✅ |
| **2_go_optimized** | 7s (-70%) 🏆 | 21.6k (-25%) 🏆 | 2.0 | 100% | 0% | 100.0% | ✅ |
| **3_default_godoctor** | 9s (-63%) | 27.6k (-5%) | 2.0 | 100% | 0% | 100.0% | ✅ |
| **4_go_godoctor** | 11s (-55%) | 23.1k (-20%) | 2.0 | 100% | 0% | 100.0% | ✅ |

### 💡 Efficiency Highlights
- **Most Tool-Efficient**: 1_default (1.0 calls/success run) vs 4_go_godoctor (2.0 calls/success run)
- **Most Token-Efficient**: 2_go_optimized (21.6k/run) vs 1_default (29.0k/run)
- **Best Tool Accuracy**: 1_default (100.0% success) vs 4_go_godoctor (100.0% success)

### 📊 Statistical Analysis Notes
> [!WARNING]
> **Small Sample Size Detected:** Some alternatives have fewer than 5 repetitions. Statistical significance flags might be unreliable. Consider running at least 10 repetitions for more robust results.
- `*` Significant ($p < 0.05$): 95% confidence that the difference is not due to chance.
- `**` Highly Significant ($p < 0.01$): 99% confidence.
- Metrics compared using Welch's t-test (Duration, Tokens, Lint) and Fisher's Exact Test (Success Rate).

## 📈 Summary by Scenario & Alternative
| Scenario | Alternative | Runs | Avg Duration (Success) | Avg Tokens (Success) | Lint | Pass Rate |
|---|---|---|---|---|---|---|
| go_symbol_listing | 1_default | 1 | 23.558s | 29.0k | 2.0 | 0% |
| go_symbol_listing | 2_go_optimized | 1 | **7.053s** | **21.6k** | 2.0 | 0% |
| go_symbol_listing | 3_default_godoctor | 1 | 8.812s | 27.6k | 2.0 | 0% |
| go_symbol_listing | 4_go_godoctor | 1 | 10.713s | 23.1k | 2.0 | 0% |

## 📋 Detailed Metrics

<details>
<summary>Click to expand detailed run-by-run metrics</summary>

| Scenario | Alternative | Rep | Duration | Tokens | Tools | Test | Lint | Status |
|---|---|---|---|---|---|---|---|---|
| go_symbol_listing | 1_default | 1 | 23.558s | 29.0k | 1 | SKIP | 2 | ⚠️ Skip |
| go_symbol_listing | 2_go_optimized | 1 | 7.053s | 21.6k | 2 | SKIP | 2 | ⚠️ Skip |
| go_symbol_listing | 3_default_godoctor | 1 | 8.812s | 27.6k | 2 | SKIP | 2 | ⚠️ Skip |
| go_symbol_listing | 4_go_godoctor | 1 | 10.713s | 23.1k | 2 | SKIP | 2 | ⚠️ Skip |

</details>

## 🔧 Tool Usage Analysis

### Overall Summary

| Alternative | Total Calls | Failed Calls | Unique Tools | Avg Calls/Run | Most Used Tool |
|---|---|---|---|---|---|
| 3_default_godoctor | 2 | 0 | 2 | 2.0 | `inspect_file` (1) |
| 4_go_godoctor | 2 | 0 | 2 | 2.0 | `inspect_file` (1) |
| 1_default | 1 | 0 | 1 | 1.0 | `delegate_to_agent` (1) |
| 2_go_optimized | 2 | 0 | 2 | 2.0 | `inspect_file` (1) |

### Tool Availability
- **1_default** unique tools: `delegate_to_agent`
- **4_go_godoctor** unique tools: `run_shell_command`

### Tool Usage by Scenario

**go_symbol_listing**

| Alternative | Tool | Avg Calls | Total | Failed |
|---|---|---|---|---|
| 1_default | `delegate_to_agent` | 1.0 | 1 | 0 |
| 2_go_optimized | `read_file` | 1.0 | 1 | 0 |
| 2_go_optimized | `inspect_file` | 1.0 | 1 | 0 |
| 3_default_godoctor | `inspect_file` | 1.0 | 1 | 0 |
| 3_default_godoctor | `read_file` | 1.0 | 1 | 0 |
| 4_go_godoctor | `inspect_file` | 1.0 | 1 | 0 |
| 4_go_godoctor | `run_shell_command` | 1.0 | 1 | 0 |

## 🎯 Conclusion

(This section tracks success rates and performance metrics filtered for successful runs only.)

Based on the analysis, **2_go_optimized** achieved the highest success rate of 100%.
