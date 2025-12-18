# Tenkai Experiment Report
_Generated at: Wed, 17 Dec 2025 18:13:05 GMT_

## 📊 Executive Summary

### Key Metrics
- **Total Runs**: 60
- **Success Rate**: 100.0% (60/60 completed)
- **Test Pass Rate**: 100.0% (54/54)
- **Timeouts**: 0 runs (0.0%)
- **Avg Tokens (Success Only)**: 208.1k
- **Avg Lint Issues (Success Only)**: 2.5

## 🧪 Experiment Parameters
> A/B test comparing impact of implicit vs explicit system prompt instructions on tool usage.

| Parameter | Value |
| :--- | :--- |
| **Repetitions** | 10 |
| **Concurrent Workers** | 4 |
| **Total Scenarios** | 2 |

### Alternatives Configuration
- **control_no_context** (Control 🏳️)
  > No GEMINI.md context file.
  - Model: `gemini-2.5-flash`
  - MCP Settings: `configs/context-abtest/settings-godoctor.json`
- **alt1_standard_context**
  > Standard GEMINI.md context.
  - Model: `gemini-2.5-flash`
  - MCP Settings: `configs/context-abtest/settings-godoctor.json`
- **alt2_explicit_context**
  > GEMINI.md with explicit instructions to use edit_code.
  - Model: `gemini-2.5-flash`
  - MCP Settings: `configs/context-abtest/settings-godoctor.json`

## 📝 Notes
> ⚠️ Repetitions overridden to 10 via CLI (config: 2)

> ⚠️ Max concurrent workers overridden to 4 via CLI (config: 4)


## 🏆 Performance Comparison

| Alternative | Avg Duration (Success) | Avg Tokens (Success) | Avg Lint (Success) | Success Rate | Test Pass Rate | Tool Accuracy | Status |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **control_no_context** | 47s (Control) 🏆 | 186.7k 🏆 | 2.5 | 100% 🏆 | 95% | 100.0% | ✅ |
| **alt1_standard_context** | 51s (+10%) | 207.4k (+11%) | 2.5 | 100% 🏆 | 85% | 100.0% | ✅ |
| **alt2_explicit_context** | 52s (+12%) | 230.3k (+23%) | 2.5 | 100% 🏆 | 90% | 100.0% | ✅ |

### 💡 Efficiency Highlights
- **Most Tool-Efficient**: control_no_context (12.6 calls/success run) vs alt2_explicit_context (13.6 calls/success run)
- **Most Token-Efficient**: control_no_context (186.7k/run) vs alt2_explicit_context (230.3k/run)
- **Best Tool Accuracy**: control_no_context (100.0% success) vs alt2_explicit_context (100.0% success)

## 📈 Summary by Scenario & Alternative
| Scenario | Alternative | Runs | Avg Duration (Success) | Avg Tokens (Success) | Lint | Pass Rate |
|---|---|---|---|---|---|---|
| concurrency_debug | alt1_standard_context | 10 | 44.34s | 133.6k | 0.0 | 100% |
| concurrency_debug | alt2_explicit_context | 10 | 45.728s | 170.3k | 0.1 | 100% |
| concurrency_debug | control_no_context | 10 | **39.323s** | **126.1k** | 0.0 | 100% |
| legacy_refactor | alt1_standard_context | 10 | 58.27s | 281.1k | 4.9 | 70% |
| legacy_refactor | alt2_explicit_context | 10 | 58.473s | 290.3k | 4.9 | 80% |
| legacy_refactor | control_no_context | 10 | **53.727s** | **247.4k** | 5.0 | 90% |

## 📋 Detailed Metrics

<details>
<summary>Click to expand detailed run-by-run metrics</summary>

| Scenario | Alternative | Rep | Duration | Tokens | Tools | Test | Lint | Status |
|---|---|---|---|---|---|---|---|---|
| concurrency_debug | alt1_standard_context | 4 | 45.32s | 131.9k | 9 | PASS | 0 | ✅ |
| concurrency_debug | alt1_standard_context | 10 | 28.26s | 116.3k | 9 | PASS | 0 | ✅ |
| concurrency_debug | alt1_standard_context | 1 | 37.94s | 157.8k | 11 | PASS | 0 | ✅ |
| concurrency_debug | alt1_standard_context | 8 | 33.273s | 95.2k | 5 | PASS | 0 | ✅ |
| concurrency_debug | alt1_standard_context | 2 | 28.924s | 93.8k | 7 | PASS | 0 | ✅ |
| concurrency_debug | alt1_standard_context | 9 | 34.373s | 89.2k | 7 | PASS | 0 | ✅ |
| concurrency_debug | alt1_standard_context | 3 | 1m49.77s | 201.9k | 14 | PASS | 0 | ✅ |
| concurrency_debug | alt1_standard_context | 7 | 38.941s | 120.5k | 8 | PASS | 0 | ✅ |
| concurrency_debug | alt1_standard_context | 5 | 39.339s | 145.7k | 11 | PASS | 0 | ✅ |
| concurrency_debug | alt1_standard_context | 6 | 47.26s | 183.4k | 11 | PASS | 0 | ✅ |
| concurrency_debug | alt2_explicit_context | 8 | 39.322s | 133.9k | 9 | PASS | 0 | ✅ |
| concurrency_debug | alt2_explicit_context | 6 | 46.327s | 144.3k | 9 | PASS | 0 | ✅ |
| concurrency_debug | alt2_explicit_context | 7 | 35.88s | 138.2k | 9 | PASS | 0 | ✅ |
| concurrency_debug | alt2_explicit_context | 4 | 40.454s | 182.4k | 10 | PASS | 0 | ✅ |
| concurrency_debug | alt2_explicit_context | 3 | 53.729s | 236.9k | 9 | PASS | 0 | ✅ |
| concurrency_debug | alt2_explicit_context | 5 | 33.976s | 117.3k | 6 | PASS | 0 | ✅ |
| concurrency_debug | alt2_explicit_context | 2 | 47.804s | 66.2k | 3 | PASS | 0 | ✅ |
| concurrency_debug | alt2_explicit_context | 1 | 45.626s | 136.1k | 9 | PASS | 0 | ✅ |
| concurrency_debug | alt2_explicit_context | 9 | 58.713s | 258.4k | 12 | PASS | 1 | ✅ |
| concurrency_debug | alt2_explicit_context | 10 | 55.452s | 289.4k | 15 | PASS | 0 | ✅ |
| concurrency_debug | control_no_context | 7 | 45.304s | 119.8k | 9 | PASS | 0 | ✅ |
| concurrency_debug | control_no_context | 3 | 56.1s | 133.8k | 9 | PASS | 0 | ✅ |
| concurrency_debug | control_no_context | 5 | 28.825s | 91.8k | 7 | PASS | 0 | ✅ |
| concurrency_debug | control_no_context | 10 | 34.199s | 89.4k | 6 | PASS | 0 | ✅ |
| concurrency_debug | control_no_context | 8 | 40.02s | 143.4k | 10 | PASS | 0 | ✅ |
| concurrency_debug | control_no_context | 2 | 40.027s | 152.9k | 11 | PASS | 0 | ✅ |
| concurrency_debug | control_no_context | 9 | 22.421s | 59.8k | 4 | PASS | 0 | ✅ |
| concurrency_debug | control_no_context | 4 | 52.803s | 225.3k | 15 | PASS | 0 | ✅ |
| concurrency_debug | control_no_context | 1 | 40.144s | 142.0k | 11 | PASS | 0 | ✅ |
| concurrency_debug | control_no_context | 6 | 33.386s | 102.6k | 7 | PASS | 0 | ✅ |
| legacy_refactor | alt1_standard_context | 8 | 53.321s | 219.2k | 14 | SKIP | 4 | ⚠️ Skip |
| legacy_refactor | alt1_standard_context | 9 | 1m16.088s | 333.2k | 19 | PASS | 6 | ✅ |
| legacy_refactor | alt1_standard_context | 1 | 42.454s | 232.4k | 16 | SKIP | 3 | ⚠️ Skip |
| legacy_refactor | alt1_standard_context | 3 | 47.661s | 281.1k | 19 | PASS | 6 | ✅ |
| legacy_refactor | alt1_standard_context | 2 | 58.178s | 302.6k | 19 | PASS | 3 | ✅ |
| legacy_refactor | alt1_standard_context | 4 | 47.419s | 231.0k | 15 | PASS | 6 | ✅ |
| legacy_refactor | alt1_standard_context | 5 | 1m14.159s | 391.2k | 24 | PASS | 6 | ✅ |
| legacy_refactor | alt1_standard_context | 7 | 55.666s | 199.2k | 12 | SKIP | 3 | ⚠️ Skip |
| legacy_refactor | alt1_standard_context | 6 | 1m11.208s | 320.7k | 19 | PASS | 6 | ✅ |
| legacy_refactor | alt1_standard_context | 10 | 56.54s | 300.7k | 19 | PASS | 6 | ✅ |
| legacy_refactor | alt2_explicit_context | 10 | 49.072s | 162.1k | 15 | PASS | 8 | ✅ |
| legacy_refactor | alt2_explicit_context | 5 | 1m6.402s | 321.7k | 20 | PASS | 5 | ✅ |
| legacy_refactor | alt2_explicit_context | 4 | 1m13.723s | 389.3k | 23 | PASS | 3 | ✅ |
| legacy_refactor | alt2_explicit_context | 9 | 55.861s | 241.0k | 14 | PASS | 6 | ✅ |
| legacy_refactor | alt2_explicit_context | 2 | 38.331s | 209.0k | 14 | SKIP | 3 | ⚠️ Skip |
| legacy_refactor | alt2_explicit_context | 3 | 58.493s | 299.9k | 18 | PASS | 6 | ✅ |
| legacy_refactor | alt2_explicit_context | 8 | 42.551s | 222.6k | 16 | PASS | 3 | ✅ |
| legacy_refactor | alt2_explicit_context | 7 | 55.93s | 294.3k | 19 | PASS | 6 | ✅ |
| legacy_refactor | alt2_explicit_context | 6 | 45.578s | 178.6k | 12 | SKIP | 3 | ⚠️ Skip |
| legacy_refactor | alt2_explicit_context | 1 | 1m38.792s | 584.1k | 29 | PASS | 6 | ✅ |
| legacy_refactor | control_no_context | 10 | 1m12.72s | 358.1k | 22 | PASS | 6 | ✅ |
| legacy_refactor | control_no_context | 7 | 38.846s | 203.5k | 14 | PASS | 5 | ✅ |
| legacy_refactor | control_no_context | 6 | 43.236s | 153.1k | 10 | PASS | 6 | ✅ |
| legacy_refactor | control_no_context | 5 | 56.753s | 263.5k | 18 | SKIP | 3 | ⚠️ Skip |
| legacy_refactor | control_no_context | 2 | 53.365s | 209.2k | 13 | PASS | 6 | ✅ |
| legacy_refactor | control_no_context | 1 | 1m11.081s | 289.9k | 18 | PASS | 5 | ✅ |
| legacy_refactor | control_no_context | 3 | 1m0.32s | 277.3k | 18 | PASS | 6 | ✅ |
| legacy_refactor | control_no_context | 9 | 48.273s | 213.4k | 15 | PASS | 6 | ✅ |
| legacy_refactor | control_no_context | 8 | 40.694s | 204.8k | 14 | PASS | 3 | ✅ |
| legacy_refactor | control_no_context | 4 | 51.98s | 300.8k | 21 | PASS | 4 | ✅ |

</details>

## 🔧 Tool Usage Analysis

### Overall Summary

| Alternative | Total Calls | Failed Calls | Unique Tools | Avg Calls/Run | Most Used Tool |
|---|---|---|---|---|---|
| alt1_standard_context | 268 | 0 | 7 | 13.4 | `run_shell_command` (90) |
| alt2_explicit_context | 271 | 0 | 6 | 13.6 | `run_shell_command` (90) |
| control_no_context | 252 | 0 | 5 | 12.6 | `run_shell_command` (95) |

### Tool Availability
- **alt1_standard_context** unique tools: `write_todos`

### Tool Usage by Scenario

**concurrency_debug**

| Alternative | Tool | Avg Calls | Total | Failed |
|---|---|---|---|---|
| alt1_standard_context | `replace` | 4.8 | 48 | 0 |
| alt1_standard_context | `run_shell_command` | 2.1 | 21 | 0 |
| alt1_standard_context | `read_file` | 1.9 | 19 | 0 |
| alt1_standard_context | `list_directory` | 0.3 | 3 | 0 |
| alt1_standard_context | `edit_code` | 0.1 | 1 | 0 |
| alt2_explicit_context | `replace` | 2.6 | 26 | 0 |
| alt2_explicit_context | `run_shell_command` | 2.4 | 24 | 0 |
| alt2_explicit_context | `edit_code` | 2.3 | 23 | 0 |
| alt2_explicit_context | `read_file` | 1.7 | 17 | 0 |
| alt2_explicit_context | `list_directory` | 0.1 | 1 | 0 |
| control_no_context | `replace` | 4.6 | 46 | 0 |
| control_no_context | `run_shell_command` | 2.3 | 23 | 0 |
| control_no_context | `read_file` | 1.9 | 19 | 0 |
| control_no_context | `list_directory` | 0.1 | 1 | 0 |

**legacy_refactor**

| Alternative | Tool | Avg Calls | Total | Failed |
|---|---|---|---|---|
| alt1_standard_context | `run_shell_command` | 6.9 | 69 | 0 |
| alt1_standard_context | `replace` | 3.8 | 38 | 0 |
| alt1_standard_context | `write_file` | 3.7 | 37 | 0 |
| alt1_standard_context | `read_file` | 2.1 | 21 | 0 |
| alt1_standard_context | `write_todos` | 0.9 | 9 | 0 |
| alt1_standard_context | `list_directory` | 0.1 | 1 | 0 |
| alt1_standard_context | `edit_code` | 0.1 | 1 | 0 |
| alt2_explicit_context | `run_shell_command` | 6.6 | 66 | 0 |
| alt2_explicit_context | `replace` | 4.6 | 46 | 0 |
| alt2_explicit_context | `write_file` | 3.7 | 37 | 0 |
| alt2_explicit_context | `read_file` | 2.4 | 24 | 0 |
| alt2_explicit_context | `list_directory` | 0.4 | 4 | 0 |
| alt2_explicit_context | `edit_code` | 0.3 | 3 | 0 |
| control_no_context | `run_shell_command` | 7.2 | 72 | 0 |
| control_no_context | `write_file` | 4.2 | 42 | 0 |
| control_no_context | `replace` | 2.9 | 29 | 0 |
| control_no_context | `read_file` | 2.0 | 20 | 0 |

## 🎯 Conclusion

(This section tracks success rates and performance metrics filtered for successful runs only.)

Based on the analysis, **alt1_standard_context** achieved the highest success rate of 100%.
