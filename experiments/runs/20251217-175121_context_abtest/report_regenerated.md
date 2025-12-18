# Tenkai Experiment Report
_Generated at: Wed, 17 Dec 2025 18:01:21 GMT_

## 📊 Executive Summary

### Key Metrics
- **Total Runs**: 12
- **Success Rate**: 100.0% (12/12 completed)
- **Test Pass Rate**: 100.0% (12/12)
- **Timeouts**: 0 runs (0.0%)
- **Avg Tokens (Success Only)**: 228.7k
- **Avg Lint Issues (Success Only)**: 2.2

## 🧪 Experiment Parameters
> A/B test comparing impact of implicit vs explicit system prompt instructions on tool usage.

| Parameter | Value |
| :--- | :--- |
| **Repetitions** | 2 |
| **Concurrent Workers** | 4 |
| **Total Scenarios** | 2 |

### Alternatives Configuration
- **control_no_context** (Control 🏳️)
  > No GEMINI.md context file.
  - Model: `gemini-2.5-flash`
  - MCP Settings: `experiments/runs/20251217-175121_context_abtest/settings-godoctor.json`
- **alt1_standard_context**
  > Standard GEMINI.md context.
  - Model: `gemini-2.5-flash`
  - MCP Settings: `experiments/runs/20251217-175121_context_abtest/settings-godoctor.json`
- **alt2_explicit_context**
  > GEMINI.md with explicit instructions to use edit_code.
  - Model: `gemini-2.5-flash`
  - MCP Settings: `experiments/runs/20251217-175121_context_abtest/settings-godoctor.json`

## 📝 Notes
> Report regenerated from results.json


## 🏆 Performance Comparison

| Alternative | Avg Duration (Success) | Avg Tokens (Success) | Avg Lint (Success) | Success Rate | Test Pass Rate | Tool Accuracy | Status |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **control_no_context** | 49s (Control) | 198.5k | 2.2 | 100% 🏆 | 100% | 100.0% | ✅ |
| **alt1_standard_context** | 59s (+19%) | 296.1k (+49%) | 2.0 | 100% 🏆 | 100% | 100.0% | ✅ |
| **alt2_explicit_context** | 48s (-3%) 🏆 | 191.5k (-3%) 🏆 | 2.5 | 100% 🏆 | 100% | 100.0% | ✅ |

### 💡 Efficiency Highlights
- **Most Tool-Efficient**: alt2_explicit_context (10.8 calls/success run) vs alt1_standard_context (18.5 calls/success run)
- **Most Token-Efficient**: alt2_explicit_context (191.5k/run) vs alt1_standard_context (296.1k/run)
- **Best Tool Accuracy**: control_no_context (100.0% success) vs alt2_explicit_context (100.0% success)

## 📈 Summary by Scenario & Alternative
| Scenario | Alternative | Runs | Avg Duration (Success) | Avg Tokens (Success) | Lint | Pass Rate |
|---|---|---|---|---|---|---|
| concurrency_debug | alt1_standard_context | 2 | 44.918s | 186.4k | 0.0 | 100% |
| concurrency_debug | alt2_explicit_context | 2 | 52.42s | 186.3k | 0.0 | 100% |
| concurrency_debug | control_no_context | 2 | **44.531s** | **135.9k** | 0.0 | 100% |
| legacy_refactor | alt1_standard_context | 2 | 1m12.779s | 405.9k | 4.0 | 100% |
| legacy_refactor | alt2_explicit_context | 2 | **43.842s** | **196.8k** | 5.0 | 100% |
| legacy_refactor | control_no_context | 2 | 54.22s | 261.1k | 4.5 | 100% |

## 📋 Detailed Metrics

<details>
<summary>Click to expand detailed run-by-run metrics</summary>

| Scenario | Alternative | Rep | Duration | Tokens | Tools | Test | Lint | Status |
|---|---|---|---|---|---|---|---|---|
| concurrency_debug | alt1_standard_context | 2 | 32.96s | 160.8k | 11 | PASS | 0 | ✅ |
| concurrency_debug | alt1_standard_context | 1 | 56.877s | 212.1k | 16 | PASS | 0 | ✅ |
| concurrency_debug | alt2_explicit_context | 1 | 53.404s | 143.5k | 5 | PASS | 0 | ✅ |
| concurrency_debug | alt2_explicit_context | 2 | 51.436s | 229.1k | 13 | PASS | 0 | ✅ |
| concurrency_debug | control_no_context | 1 | 54.434s | 180.1k | 12 | PASS | 0 | ✅ |
| concurrency_debug | control_no_context | 2 | 34.628s | 91.6k | 6 | PASS | 0 | ✅ |
| legacy_refactor | alt1_standard_context | 2 | 1m16.286s | 433.5k | 25 | PASS | 3 | ✅ |
| legacy_refactor | alt1_standard_context | 1 | 1m9.272s | 378.3k | 22 | PASS | 5 | ✅ |
| legacy_refactor | alt2_explicit_context | 2 | 46.319s | 193.0k | 12 | PASS | 4 | ✅ |
| legacy_refactor | alt2_explicit_context | 1 | 41.366s | 200.6k | 13 | PASS | 6 | ✅ |
| legacy_refactor | control_no_context | 1 | 55.837s | 350.0k | 24 | PASS | 3 | ✅ |
| legacy_refactor | control_no_context | 2 | 52.603s | 172.1k | 11 | PASS | 6 | ✅ |

</details>

## 🔧 Tool Usage Analysis

### Overall Summary

| Alternative | Total Calls | Failed Calls | Unique Tools | Avg Calls/Run | Most Used Tool |
|---|---|---|---|---|---|
| control_no_context | 53 | 0 | 5 | 13.2 | `run_shell_command` (24) |
| alt1_standard_context | 74 | 0 | 5 | 18.5 | `run_shell_command` (22) |
| alt2_explicit_context | 43 | 0 | 5 | 10.8 | `run_shell_command` (14) |

### Tool Availability
- **alt2_explicit_context** unique tools: `edit_code`
- **control_no_context** unique tools: `list_directory`
- **alt1_standard_context** unique tools: `write_todos`

### Tool Usage by Scenario

**concurrency_debug**

| Alternative | Tool | Avg Calls | Total | Failed |
|---|---|---|---|---|
| alt1_standard_context | `replace` | 7.0 | 14 | 0 |
| alt1_standard_context | `run_shell_command` | 4.0 | 8 | 0 |
| alt1_standard_context | `read_file` | 2.5 | 5 | 0 |
| alt2_explicit_context | `edit_code` | 3.5 | 7 | 0 |
| alt2_explicit_context | `run_shell_command` | 2.5 | 5 | 0 |
| alt2_explicit_context | `read_file` | 2.0 | 4 | 0 |
| alt2_explicit_context | `replace` | 1.0 | 2 | 0 |
| control_no_context | `replace` | 3.5 | 7 | 0 |
| control_no_context | `run_shell_command` | 3.0 | 6 | 0 |
| control_no_context | `read_file` | 2.0 | 4 | 0 |
| control_no_context | `list_directory` | 0.5 | 1 | 0 |

**legacy_refactor**

| Alternative | Tool | Avg Calls | Total | Failed |
|---|---|---|---|---|
| alt1_standard_context | `write_todos` | 8.0 | 16 | 0 |
| alt1_standard_context | `run_shell_command` | 7.0 | 14 | 0 |
| alt1_standard_context | `write_file` | 4.5 | 9 | 0 |
| alt1_standard_context | `read_file` | 2.5 | 5 | 0 |
| alt1_standard_context | `replace` | 1.5 | 3 | 0 |
| alt2_explicit_context | `run_shell_command` | 4.5 | 9 | 0 |
| alt2_explicit_context | `write_file` | 4.0 | 8 | 0 |
| alt2_explicit_context | `read_file` | 2.0 | 4 | 0 |
| alt2_explicit_context | `replace` | 1.5 | 3 | 0 |
| alt2_explicit_context | `edit_code` | 0.5 | 1 | 0 |
| control_no_context | `run_shell_command` | 9.0 | 18 | 0 |
| control_no_context | `write_file` | 5.0 | 10 | 0 |
| control_no_context | `read_file` | 2.0 | 4 | 0 |
| control_no_context | `replace` | 1.5 | 3 | 0 |

## 🎯 Conclusion

(This section tracks success rates and performance metrics filtered for successful runs only.)

Based on the analysis, **alt1_standard_context** achieved the highest success rate of 100%.
