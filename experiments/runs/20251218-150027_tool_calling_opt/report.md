# Tenkai Experiment Report
_Generated at: Thu, 18 Dec 2025 15:02:32 GMT_

## 📊 Executive Summary

### Key Metrics
- **Total Runs**: 18
- **Success Rate**: 100.0% (18/18 completed)
- **Test Pass Rate**: NaN% (0/0)
- **Timeouts**: 0 runs (0.0%)
- **Avg Tokens (Success Only)**: 40.9k
- **Avg Lint Issues (Success Only)**: 1.2

## 🧪 Experiment Parameters
> Analysis of various configuration layers for optimizing tool selection between built-in and MCP.

| Parameter | Value |
| :--- | :--- |
| **Repetitions** | 3 |
| **Concurrent Workers** | 3 |
| **Total Scenarios** | 1 |

### Alternatives Configuration
- **1_control** (Control 🏳️)
  > Baseline behavior (no intervention).
  - Model: `gemini-2.5-flash`
  - MCP Settings: `configs/tool-calling-opt/settings-default.json`
- **2_context_hint**
  > Soft hint via GEMINI.md.
  - Model: `gemini-2.5-flash`
  - MCP Settings: `configs/tool-calling-opt/settings-default.json`
- **3_system_override**
  > Authoritative mandate via system.md (SYSTEM_OPTIMIZED.md).
  - Model: `gemini-2.5-flash`
  - System Prompt: `configs/tool-calling-opt/SYSTEM_OPTIMIZED.md`
  - MCP Settings: `configs/tool-calling-opt/settings-default.json`
- **4_hard_block**
  > Hard exclusion of shell tool via settings.json.
  - Model: `gemini-2.5-flash`
  - MCP Settings: `configs/tool-calling-opt/settings-exclude-shell.json`
- **5_policy_friction**
  > Policy Engine rule set to 'deny' for the shell tool.
  - Model: `gemini-2.5-flash`
  - MCP Settings: `configs/tool-calling-opt/settings-default.json`
- **6_ensemble**
  > Combines system_override and hard_block.
  - Model: `gemini-2.5-flash`
  - System Prompt: `configs/tool-calling-opt/SYSTEM_OPTIMIZED.md`
  - MCP Settings: `configs/tool-calling-opt/settings-exclude-shell.json`

## 🏆 Performance Comparison
| Alternative | Avg Duration (Success) | Avg Tokens (Success) | Avg Lint (Success) | Success Rate | Test Pass Rate | Tool Accuracy | Status |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **1_control** | 11s (Control) | 34.8k | 0.0 | 100% | 0% | 100.0% | ✅ |
| **2_context_hint** | 24s (+126%) | 42.0k (+21%) | 1.3 | 100% | 0% | 100.0% | ✅ |
| **3_system_override** | 14s (+34%) | 39.9k (+15%) | 1.3 | 100% | 0% | 100.0% | ✅ |
| **4_hard_block** | 12s (+14%) | 43.2k (+24%) | 2.0 ** | 100% | 0% | 100.0% | ✅ |
| **5_policy_friction** | 21s (+97%) ** | 60.2k (+73%) ** | 1.3 | 100% | 0% | 100.0% | ✅ |
| **6_ensemble** | 8s (-20%) 🏆 | 25.3k (-27%) ** 🏆 | 1.3 | 100% | 0% | 100.0% | ✅ |

### 💡 Efficiency Highlights
- **Most Tool-Efficient**: 6_ensemble (2.0 calls/success run) vs 5_policy_friction (5.3 calls/success run)
- **Most Token-Efficient**: 6_ensemble (25.3k/run) vs 5_policy_friction (60.2k/run)
- **Best Tool Accuracy**: 1_control (100.0% success) vs 6_ensemble (100.0% success)

### 📊 Statistical Analysis Notes
> [!WARNING]
> **Small Sample Size Detected:** Some alternatives have fewer than 5 repetitions. Statistical significance flags might be unreliable. Consider running at least 10 repetitions for more robust results.
- `*` Significant ($p < 0.05$): 95% confidence that the difference is not due to chance.
- `**` Highly Significant ($p < 0.01$): 99% confidence.
- Metrics compared using Welch's t-test (Duration, Tokens, Lint) and Fisher's Exact Test (Success Rate).

## 📈 Summary by Scenario & Alternative
| Scenario | Alternative | Runs | Avg Duration (Success) | Avg Tokens (Success) | Lint | Pass Rate |
|---|---|---|---|---|---|---|
| go_symbol_listing | 1_control | 3 | 10.538s | 34.8k | 0.0 | 0% |
| go_symbol_listing | 2_context_hint | 3 | 23.844s | 42.0k | 1.3 | 0% |
| go_symbol_listing | 3_system_override | 3 | 14.171s | 39.9k | 1.3 | 0% |
| go_symbol_listing | 4_hard_block | 3 | 12.029s | 43.2k | 2.0 | 0% |
| go_symbol_listing | 5_policy_friction | 3 | 20.729s | 60.2k | 1.3 | 0% |
| go_symbol_listing | 6_ensemble | 3 | **8.391s** | **25.3k** | 1.3 | 0% |

## 📋 Detailed Metrics

<details>
<summary>Click to expand detailed run-by-run metrics</summary>

| Scenario | Alternative | Rep | Duration | Tokens | Tools | Test | Lint | Status |
|---|---|---|---|---|---|---|---|---|
| go_symbol_listing | 1_control | 1 | 10.275s | 38.1k | 3 | SKIP | 0 | ⚠️ Skip |
| go_symbol_listing | 1_control | 2 | 7.188s | 28.0k | 2 | SKIP | 0 | ⚠️ Skip |
| go_symbol_listing | 1_control | 3 | 14.151s | 38.2k | 2 | SKIP | 0 | ⚠️ Skip |
| go_symbol_listing | 2_context_hint | 2 | 42.756s | 38.9k | 3 | SKIP | 2 | ⚠️ Skip |
| go_symbol_listing | 2_context_hint | 1 | 12.849s | 48.4k | 4 | SKIP | 0 | ⚠️ Skip |
| go_symbol_listing | 2_context_hint | 3 | 15.928s | 38.8k | 3 | SKIP | 2 | ⚠️ Skip |
| go_symbol_listing | 3_system_override | 2 | 16.193s | 37.3k | 3 | SKIP | 2 | ⚠️ Skip |
| go_symbol_listing | 3_system_override | 3 | 8.545s | 26.7k | 2 | SKIP | 2 | ⚠️ Skip |
| go_symbol_listing | 3_system_override | 1 | 17.776s | 55.6k | 4 | SKIP | 0 | ⚠️ Skip |
| go_symbol_listing | 4_hard_block | 2 | 12.911s | 46.7k | 4 | SKIP | 2 | ⚠️ Skip |
| go_symbol_listing | 4_hard_block | 3 | 13.793s | 46.3k | 4 | SKIP | 2 | ⚠️ Skip |
| go_symbol_listing | 4_hard_block | 1 | 9.382s | 36.5k | 2 | SKIP | 2 | ⚠️ Skip |
| go_symbol_listing | 5_policy_friction | 2 | 24.947s | 62.7k | 8 | SKIP | 2 | ⚠️ Skip |
| go_symbol_listing | 5_policy_friction | 3 | 16.434s | 58.1k | 4 | SKIP | 2 | ⚠️ Skip |
| go_symbol_listing | 5_policy_friction | 1 | 20.805s | 59.8k | 4 | SKIP | 0 | ⚠️ Skip |
| go_symbol_listing | 6_ensemble | 2 | 7.501s | 25.1k | 2 | SKIP | 2 | ⚠️ Skip |
| go_symbol_listing | 6_ensemble | 3 | 7.782s | 25.2k | 2 | SKIP | 2 | ⚠️ Skip |
| go_symbol_listing | 6_ensemble | 1 | 9.889s | 25.5k | 2 | SKIP | 0 | ⚠️ Skip |

</details>

## 🔧 Tool Usage Analysis

### Overall Summary

| Alternative | Total Calls | Failed Calls | Unique Tools | Avg Calls/Run | Most Used Tool |
|---|---|---|---|---|---|
| 1_control | 7 | 0 | 3 | 2.3 | `inspect_file` (3) |
| 2_context_hint | 10 | 0 | 5 | 3.3 | `read_file` (3) |
| 3_system_override | 9 | 0 | 4 | 3.0 | `inspect_file` (3) |
| 4_hard_block | 10 | 0 | 4 | 3.3 | `read_file` (4) |
| 5_policy_friction | 16 | 0 | 4 | 5.3 | `search_file_content` (8) |
| 6_ensemble | 6 | 0 | 2 | 2.0 | `inspect_file` (3) |

### Tool Availability
- **2_context_hint** unique tools: `review_code`

### Tool Usage by Scenario

**go_symbol_listing**

| Alternative | Tool | Avg Calls | Total | Failed |
|---|---|---|---|---|
| 1_control | `inspect_file` | 1.0 | 3 | 0 |
| 1_control | `read_file` | 1.0 | 3 | 0 |
| 1_control | `search_file_content` | 0.3 | 1 | 0 |
| 2_context_hint | `inspect_file` | 1.0 | 3 | 0 |
| 2_context_hint | `read_file` | 1.0 | 3 | 0 |
| 2_context_hint | `read_godoc` | 0.7 | 2 | 0 |
| 2_context_hint | `review_code` | 0.3 | 1 | 0 |
| 2_context_hint | `run_shell_command` | 0.3 | 1 | 0 |
| 3_system_override | `inspect_file` | 1.0 | 3 | 0 |
| 3_system_override | `read_file` | 1.0 | 3 | 0 |
| 3_system_override | `run_shell_command` | 0.7 | 2 | 0 |
| 3_system_override | `search_file_content` | 0.3 | 1 | 0 |
| 4_hard_block | `read_file` | 1.3 | 4 | 0 |
| 4_hard_block | `inspect_file` | 1.0 | 3 | 0 |
| 4_hard_block | `read_godoc` | 0.7 | 2 | 0 |
| 4_hard_block | `search_file_content` | 0.3 | 1 | 0 |
| 5_policy_friction | `search_file_content` | 2.7 | 8 | 0 |
| 5_policy_friction | `read_file` | 1.3 | 4 | 0 |
| 5_policy_friction | `inspect_file` | 1.0 | 3 | 0 |
| 5_policy_friction | `read_godoc` | 0.3 | 1 | 0 |
| 6_ensemble | `inspect_file` | 1.0 | 3 | 0 |
| 6_ensemble | `read_file` | 1.0 | 3 | 0 |

## 🎯 Conclusion

(This section tracks success rates and performance metrics filtered for successful runs only.)

Based on the analysis, **6_ensemble** achieved the highest success rate of 100%.
