# Tenkai Experiment Report
_Generated at: Wed, 17 Dec 2025 18:06:31 GMT_

## 📊 Executive Summary

### Key Metrics
- **Total Runs**: 60
- **Success Rate**: 90.0% (54/60 completed)
- **Test Pass Rate**: 96.0% (48/50)
- **Timeouts**: 6 runs (10.0%)
- **Loop Warnings**: 2 runs
- **Avg Tokens (Success Only)**: 289.4k
- **Avg Lint Issues (Success Only)**: 2.3

## 🧪 Experiment Parameters
> Standardized AB test comparing Vanilla, Go, GoDoctor, and agents using Gemini 2.5 Flash across standard scenarios.

| Parameter | Value |
| :--- | :--- |
| **Repetitions** | 2 |
| **Concurrent Workers** | 4 |
| **Total Scenarios** | 5 |

### Alternatives Configuration
- **gemini_flash** (Control 🏳️)
  > Control: Standard Gemini 2.5 Flash.
  - Model: `gemini-2.5-flash`
- **gemini_flash_godoctor**
  > Gemini 2.5 Flash with GoDoctor MCP tools.
  - Model: `gemini-2.5-flash`
  - MCP Settings: `configs/gemini-go-mcp-flash/settings-godoctor.json`
- **gemini_flash_godoctor_agents**
  > Gemini 2.5 Flash with GoDoctor MCP tools and custom agents.md context.
  - Model: `gemini-2.5-flash`
  - MCP Settings: `configs/gemini-go-mcp-flash/settings-godoctor.json`
- **gemini_flash_go**
  > Gemini 2.5 Flash with Go system prompt.
  - Model: `gemini-2.5-flash`
  - System Prompt: `/Users/petruzalek/.gemini/system.golang.md`
- **gemini_flash_go_godoctor**
  > Gemini 2.5 Flash with both Go system prompt and GoDoctor MCP tools.
  - Model: `gemini-2.5-flash`
  - System Prompt: `/Users/petruzalek/.gemini/system.golang.md`
  - MCP Settings: `configs/gemini-go-mcp-flash/settings-godoctor.json`
- **gemini_flash_go_godoctor_agents**
  > Gemini 2.5 Flash with Go system prompt, GoDoctor MCP tools, and custom agents.md context.
  - Model: `gemini-2.5-flash`
  - System Prompt: `/Users/petruzalek/.gemini/system.golang.md`
  - MCP Settings: `configs/gemini-go-mcp-flash/settings-godoctor.json`

## 📝 Notes
> ⚠️ Repetitions overridden to 2 via CLI (config: 3)

> ⚠️ Max concurrent workers overridden to 4 via CLI (config: 8)


## 🏆 Performance Comparison

| Alternative | Avg Duration (Success) | Avg Tokens (Success) | Success Rate | Test Pass Rate | Tool Accuracy | Status |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **gemini_flash** | 58s (Control) 🏆 | 208.5k 🏆 | 90% | 89% | 100.0% | ✅ ⏱️ |
| **gemini_flash_godoctor** | 1m2s (+8%) | 253.4k (+22%) | 90% | 89% | 100.0% | ✅ ⏱️ |
| **gemini_flash_godoctor_agents** | 1m16s (+31%) | 275.4k (+32%) | 90% | 89% | 100.0% | ✅ ⏱️ |
| **gemini_flash_go** | 1m31s (+58%) | 333.3k (+60%) | 100% 🏆 | 80% | 100.0% | ✅ |
| **gemini_flash_go_godoctor** | 59s (+2%) | 280.1k (+34%) | 90% | 100% | 100.0% | ✅ ⏱️ |
| **gemini_flash_go_godoctor_agents** | 1m13s (+27%) | 392.1k (+88%) | 80% | 88% | 100.0% | ✅ ⏱️ |

### 💡 Efficiency Highlights
- **Most Tool-Efficient**: gemini_flash (12.9 calls/success run) vs gemini_flash_go (15.9 calls/success run)
- **Most Token-Efficient**: gemini_flash (208.5k/run) vs gemini_flash_go_godoctor_agents (392.1k/run)
- **Best Tool Accuracy**: gemini_flash (100.0% success) vs gemini_flash_go_godoctor_agents (100.0% success)

## 📈 Summary by Scenario & Alternative
| Scenario | Alternative | Runs | Avg Duration (Success) | Avg Tokens (Success) | Lint | Pass Rate |
|---|---|---|---|---|---|---|
| api_implementation | gemini_flash | 2 | **34.44s** | **137.9k** | 4.5 | 100% |
| api_implementation | gemini_flash_go | 2 | 2m8.571s | 509.4k | 2.0 | 50% |
| api_implementation | gemini_flash_go_godoctor | 2 | 1m0.591s | 315.6k | 1.5 | 100% |
| api_implementation | gemini_flash_go_godoctor_agents | 2 | 1m1.746s | 280.6k | 3.0 | 50% |
| api_implementation | gemini_flash_godoctor | 2 | 1m7.867s | 228.4k | 2.0 | 50% |
| api_implementation | gemini_flash_godoctor_agents | 2 | 1m7.518s | 303.0k | 6.5 | 100% |
| concurrency_debug | gemini_flash | 2 | **32.59s** | **88.3k** | 0.0 | 100% |
| concurrency_debug | gemini_flash_go | 2 | 56.228s | 137.1k | 0.0 | 100% |
| concurrency_debug | gemini_flash_go_godoctor | 2 | 46.292s | 133.2k | 0.0 | 100% |
| concurrency_debug | gemini_flash_go_godoctor_agents | 2 | 53.289s | 203.7k | 0.0 | 100% |
| concurrency_debug | gemini_flash_godoctor | 2 | 32.762s | 105.1k | 0.0 | 100% |
| concurrency_debug | gemini_flash_godoctor_agents | 2 | 42.564s | 142.6k | 0.0 | 100% |
| greenfield_cli | gemini_flash | 2 | 2m31.875s | 603.5k | 1.0 | 0% |
| greenfield_cli | gemini_flash_go | 2 | 2m9.23s | 507.4k | 1.5 | 50% |
| greenfield_cli | gemini_flash_go_godoctor | 2 | **1m29.446s** | 455.2k | 9.0 | 50% |
| greenfield_cli | gemini_flash_go_godoctor_agents | 2 | 2m51.168s | 1.3M | 0.0 | 50% |
| greenfield_cli | gemini_flash_godoctor | 2 | 1m37.732s | 560.1k | 6.0 | 50% |
| greenfield_cli | gemini_flash_godoctor_agents | 2 | 1m59.012s | **432.2k** | 4.0 | 50% |
| legacy_refactor | gemini_flash | 2 | **48.085s** | **184.8k** | 6.0 | 100% |
| legacy_refactor | gemini_flash_go | 2 | 1m2.098s | 264.1k | 4.5 | 100% |
| legacy_refactor | gemini_flash_go_godoctor | 2 | 1m11.188s | 429.4k | 4.5 | 100% |
| legacy_refactor | gemini_flash_go_godoctor_agents | 2 | 1m7.238s | 365.5k | 3.0 | 50% |
| legacy_refactor | gemini_flash_godoctor | 2 | 1m8.62s | 321.7k | 6.5 | 100% |
| legacy_refactor | gemini_flash_godoctor_agents | 2 | 1m0.354s | 237.0k | 4.0 | 50% |
| test_generation | gemini_flash | 2 | 1m8.047s | 225.5k | 1.0 | 100% |
| test_generation | gemini_flash_go | 2 | 1m18.566s | 248.4k | 0.5 | 100% |
| test_generation | gemini_flash_go_godoctor | 2 | **42.047s** | **154.4k** | 1.0 | 100% |
| test_generation | gemini_flash_go_godoctor_agents | 2 | 54.955s | 211.5k | 0.5 | 100% |
| test_generation | gemini_flash_godoctor | 2 | 1m0.501s | 205.1k | 1.0 | 100% |
| test_generation | gemini_flash_godoctor_agents | 2 | 1m50.18s | 340.7k | 1.0 | 100% |

## 📋 Detailed Metrics

<details>
<summary>Click to expand detailed run-by-run metrics</summary>

| Scenario | Alternative | Rep | Duration | Tokens | Tools | Test | Lint | Status |
|---|---|---|---|---|---|---|---|---|
| api_implementation | gemini_flash | 1 | 37.026s | 162.6k | 12 | PASS | 4 | ✅ |
| api_implementation | gemini_flash | 2 | 31.854s | 113.2k | 8 | PASS | 5 | ✅ |
| api_implementation | gemini_flash_go | 1 | 1m49.192s | 502.0k | 20 | PASS | 4 | ✅ |
| api_implementation | gemini_flash_go | 2 | 2m27.949s | 516.7k | 22 | FAIL | 0 | ❌ Fail |
| api_implementation | gemini_flash_go_godoctor | 1 | 47.818s | 209.3k | 12 | PASS | 3 | ✅ |
| api_implementation | gemini_flash_go_godoctor | 2 | 1m13.365s | 422.0k | 21 | PASS | 0 | ✅ |
| api_implementation | gemini_flash_go_godoctor_agents | 1 | 5m0.015s | 0 | 51 | - | 0 | ⏱️ Timeout |
| api_implementation | gemini_flash_go_godoctor_agents | 2 | 1m1.746s | 280.6k | 16 | PASS | 3 | ✅ |
| api_implementation | gemini_flash_godoctor | 1 | 1m55.929s | 406.4k | 21 | PASS | 4 | ✅ |
| api_implementation | gemini_flash_godoctor | 2 | 19.805s | 50.5k | 1 | SKIP | 0 | ⚠️ Skip |
| api_implementation | gemini_flash_godoctor_agents | 2 | 1m2.666s | 324.8k | 20 | PASS | 7 | ✅ |
| api_implementation | gemini_flash_godoctor_agents | 1 | 1m12.369s | 281.2k | 18 | PASS | 6 | ✅ |
| concurrency_debug | gemini_flash | 1 | 40.124s | 103.5k | 9 | PASS | 0 | ✅ |
| concurrency_debug | gemini_flash | 2 | 25.055s | 73.1k | 10 | PASS | 0 | ✅ |
| concurrency_debug | gemini_flash_go | 1 | 27.852s | 59.8k | 4 | PASS | 0 | ✅ |
| concurrency_debug | gemini_flash_go | 2 | 1m24.604s | 214.5k | 11 | PASS | 0 | ✅ |
| concurrency_debug | gemini_flash_go_godoctor | 1 | 1m4.536s | 165.9k | 9 | PASS | 0 | ✅ |
| concurrency_debug | gemini_flash_go_godoctor | 2 | 28.048s | 100.5k | 6 | PASS | 0 | ✅ |
| concurrency_debug | gemini_flash_go_godoctor_agents | 2 | 1m24.699s | 333.7k | 11 | PASS | 0 | ✅ |
| concurrency_debug | gemini_flash_go_godoctor_agents | 1 | 21.879s | 73.7k | 5 | PASS | 0 | ✅ |
| concurrency_debug | gemini_flash_godoctor | 1 | 31.397s | 107.9k | 8 | PASS | 0 | ✅ |
| concurrency_debug | gemini_flash_godoctor | 2 | 34.128s | 102.3k | 8 | PASS | 0 | ✅ |
| concurrency_debug | gemini_flash_godoctor_agents | 2 | 48.432s | 157.5k | 10 | PASS | 0 | ✅ |
| concurrency_debug | gemini_flash_godoctor_agents | 1 | 36.696s | 127.6k | 9 | PASS | 0 | ✅ |
| greenfield_cli | gemini_flash | 2 | 2m31.875s | 603.5k | 23 | FAIL | 1 | 🔄 Loop |
| greenfield_cli | gemini_flash | 1 | 5m0.022s | 0 | 17 | - | 0 | ⏱️ Timeout |
| greenfield_cli | gemini_flash_go | 1 | 2m5.131s | 552.2k | 27 | SKIP | 1 | 🔄 Loop |
| greenfield_cli | gemini_flash_go | 2 | 2m13.329s | 462.5k | 22 | PASS | 2 | ✅ |
| greenfield_cli | gemini_flash_go_godoctor | 1 | 1m29.446s | 455.2k | 22 | PASS | 9 | ✅ |
| greenfield_cli | gemini_flash_go_godoctor | 2 | 5m0.013s | 0 | 13 | - | 0 | ⏱️ Timeout |
| greenfield_cli | gemini_flash_go_godoctor_agents | 1 | 2m51.168s | 1.3M | 39 | PASS | 0 | ✅ |
| greenfield_cli | gemini_flash_go_godoctor_agents | 2 | 5m0.021s | 0 | 11 | - | 0 | ⏱️ Timeout |
| greenfield_cli | gemini_flash_godoctor | 2 | 1m37.732s | 560.1k | 28 | PASS | 6 | ✅ |
| greenfield_cli | gemini_flash_godoctor | 1 | 5m0.023s | 0 | 12 | - | 0 | ⏱️ Timeout |
| greenfield_cli | gemini_flash_godoctor_agents | 1 | 5m0.023s | 0 | 47 | - | 0 | ⏱️ Timeout |
| greenfield_cli | gemini_flash_godoctor_agents | 2 | 1m59.012s | 432.2k | 19 | PASS | 4 | ✅ |
| legacy_refactor | gemini_flash | 1 | 41.946s | 134.4k | 12 | PASS | 6 | ✅ |
| legacy_refactor | gemini_flash | 2 | 54.224s | 235.2k | 16 | PASS | 6 | ✅ |
| legacy_refactor | gemini_flash_go | 2 | 39.025s | 177.6k | 11 | PASS | 6 | ✅ |
| legacy_refactor | gemini_flash_go | 1 | 1m25.171s | 350.7k | 20 | PASS | 3 | ✅ |
| legacy_refactor | gemini_flash_go_godoctor | 2 | 47.337s | 223.7k | 13 | PASS | 6 | ✅ |
| legacy_refactor | gemini_flash_go_godoctor | 1 | 1m35.04s | 635.2k | 32 | PASS | 3 | ✅ |
| legacy_refactor | gemini_flash_go_godoctor_agents | 2 | 1m1.613s | 366.2k | 19 | SKIP | 0 | ⚠️ Skip |
| legacy_refactor | gemini_flash_go_godoctor_agents | 1 | 1m12.862s | 364.8k | 18 | PASS | 6 | ✅ |
| legacy_refactor | gemini_flash_godoctor | 1 | 1m12.543s | 328.6k | 21 | PASS | 7 | ✅ |
| legacy_refactor | gemini_flash_godoctor | 2 | 1m4.698s | 314.8k | 19 | PASS | 6 | ✅ |
| legacy_refactor | gemini_flash_godoctor_agents | 2 | 51.257s | 230.8k | 14 | SKIP | 3 | ⚠️ Skip |
| legacy_refactor | gemini_flash_godoctor_agents | 1 | 1m9.452s | 243.3k | 14 | PASS | 5 | ✅ |
| test_generation | gemini_flash | 1 | 1m14.646s | 213.3k | 11 | PASS | 1 | ✅ |
| test_generation | gemini_flash | 2 | 1m1.449s | 237.7k | 15 | PASS | 1 | ✅ |
| test_generation | gemini_flash_go | 2 | 37.753s | 111.8k | 6 | PASS | 0 | ✅ |
| test_generation | gemini_flash_go | 1 | 1m59.379s | 385.0k | 16 | PASS | 1 | ✅ |
| test_generation | gemini_flash_go_godoctor | 2 | 45.152s | 167.4k | 9 | PASS | 1 | ✅ |
| test_generation | gemini_flash_go_godoctor | 1 | 38.941s | 141.3k | 7 | PASS | 1 | ✅ |
| test_generation | gemini_flash_go_godoctor_agents | 1 | 57.237s | 273.8k | 12 | PASS | 0 | ✅ |
| test_generation | gemini_flash_go_godoctor_agents | 2 | 52.672s | 149.1k | 7 | PASS | 1 | ✅ |
| test_generation | gemini_flash_godoctor | 2 | 51.636s | 172.3k | 10 | PASS | 1 | ✅ |
| test_generation | gemini_flash_godoctor | 1 | 1m9.366s | 237.8k | 12 | PASS | 1 | ✅ |
| test_generation | gemini_flash_godoctor_agents | 2 | 2m36.103s | 353.3k | 14 | PASS | 1 | ✅ |
| test_generation | gemini_flash_godoctor_agents | 1 | 1m4.257s | 328.1k | 18 | PASS | 1 | ✅ |

</details>

## 🔧 Tool Usage Analysis

### Overall Summary

| Alternative | Total Calls | Failed Calls | Unique Tools | Avg Calls/Run | Most Used Tool |
|---|---|---|---|---|---|
| gemini_flash_godoctor_agents | 183 | 0 | 6 | 18.3 | `replace` (65) |
| gemini_flash | 133 | 0 | 5 | 13.3 | `run_shell_command` (54) |
| gemini_flash_go | 159 | 0 | 6 | 15.9 | `run_shell_command` (62) |
| gemini_flash_go_godoctor | 144 | 0 | 6 | 14.4 | `run_shell_command` (63) |
| gemini_flash_go_godoctor_agents | 189 | 0 | 6 | 18.9 | `run_shell_command` (85) |
| gemini_flash_godoctor | 140 | 0 | 5 | 14.0 | `run_shell_command` (59) |

### Tool Availability
- **gemini_flash_go_godoctor_agents** unique tools: `edit_code`

### Tool Usage by Scenario

**api_implementation**

| Alternative | Tool | Avg Calls | Total | Failed |
|---|---|---|---|---|
| gemini_flash | `run_shell_command` | 6.0 | 12 | 0 |
| gemini_flash | `write_file` | 2.0 | 4 | 0 |
| gemini_flash | `read_file` | 1.0 | 2 | 0 |
| gemini_flash | `list_directory` | 0.5 | 1 | 0 |
| gemini_flash | `replace` | 0.5 | 1 | 0 |
| gemini_flash_go | `run_shell_command` | 8.0 | 16 | 0 |
| gemini_flash_go | `replace` | 7.0 | 14 | 0 |
| gemini_flash_go | `write_file` | 4.0 | 8 | 0 |
| gemini_flash_go | `read_file` | 2.0 | 4 | 0 |
| gemini_flash_go_godoctor | `run_shell_command` | 7.5 | 15 | 0 |
| gemini_flash_go_godoctor | `write_file` | 3.5 | 7 | 0 |
| gemini_flash_go_godoctor | `replace` | 3.5 | 7 | 0 |
| gemini_flash_go_godoctor | `read_file` | 1.5 | 3 | 0 |
| gemini_flash_go_godoctor | `list_directory` | 0.5 | 1 | 0 |
| gemini_flash_go_godoctor_agents | `run_shell_command` | 15.5 | 31 | 0 |
| gemini_flash_go_godoctor_agents | `replace` | 10.5 | 21 | 0 |
| gemini_flash_go_godoctor_agents | `write_file` | 5.5 | 11 | 0 |
| gemini_flash_go_godoctor_agents | `read_file` | 2.0 | 4 | 0 |
| gemini_flash_godoctor | `run_shell_command` | 4.5 | 9 | 0 |
| gemini_flash_godoctor | `read_file` | 3.0 | 6 | 0 |
| gemini_flash_godoctor | `replace` | 2.5 | 5 | 0 |
| gemini_flash_godoctor | `write_file` | 1.0 | 2 | 0 |
| gemini_flash_godoctor_agents | `run_shell_command` | 8.0 | 16 | 0 |
| gemini_flash_godoctor_agents | `write_file` | 3.5 | 7 | 0 |
| gemini_flash_godoctor_agents | `replace` | 3.5 | 7 | 0 |
| gemini_flash_godoctor_agents | `read_file` | 3.0 | 6 | 0 |
| gemini_flash_godoctor_agents | `list_directory` | 1.0 | 2 | 0 |

**concurrency_debug**

| Alternative | Tool | Avg Calls | Total | Failed |
|---|---|---|---|---|
| gemini_flash | `replace` | 5.0 | 10 | 0 |
| gemini_flash | `read_file` | 2.5 | 5 | 0 |
| gemini_flash | `run_shell_command` | 2.0 | 4 | 0 |
| gemini_flash_go | `write_todos` | 2.5 | 5 | 0 |
| gemini_flash_go | `read_file` | 2.0 | 4 | 0 |
| gemini_flash_go | `replace` | 2.0 | 4 | 0 |
| gemini_flash_go | `run_shell_command` | 1.0 | 2 | 0 |
| gemini_flash_go_godoctor | `run_shell_command` | 3.0 | 6 | 0 |
| gemini_flash_go_godoctor | `replace` | 2.0 | 4 | 0 |
| gemini_flash_go_godoctor | `read_file` | 1.5 | 3 | 0 |
| gemini_flash_go_godoctor | `list_directory` | 1.0 | 2 | 0 |
| gemini_flash_go_godoctor_agents | `run_shell_command` | 4.5 | 9 | 0 |
| gemini_flash_go_godoctor_agents | `read_file` | 2.0 | 4 | 0 |
| gemini_flash_go_godoctor_agents | `edit_code` | 1.0 | 2 | 0 |
| gemini_flash_go_godoctor_agents | `write_file` | 0.5 | 1 | 0 |
| gemini_flash_godoctor | `replace` | 4.5 | 9 | 0 |
| gemini_flash_godoctor | `read_file` | 1.5 | 3 | 0 |
| gemini_flash_godoctor | `run_shell_command` | 1.5 | 3 | 0 |
| gemini_flash_godoctor | `list_directory` | 0.5 | 1 | 0 |
| gemini_flash_godoctor_agents | `replace` | 5.5 | 11 | 0 |
| gemini_flash_godoctor_agents | `run_shell_command` | 2.0 | 4 | 0 |
| gemini_flash_godoctor_agents | `read_file` | 1.5 | 3 | 0 |
| gemini_flash_godoctor_agents | `list_directory` | 0.5 | 1 | 0 |

**greenfield_cli**

| Alternative | Tool | Avg Calls | Total | Failed |
|---|---|---|---|---|
| gemini_flash | `run_shell_command` | 7.0 | 14 | 0 |
| gemini_flash | `write_file` | 6.5 | 13 | 0 |
| gemini_flash | `replace` | 3.5 | 7 | 0 |
| gemini_flash | `read_file` | 2.5 | 5 | 0 |
| gemini_flash | `list_directory` | 0.5 | 1 | 0 |
| gemini_flash_go | `run_shell_command` | 12.5 | 25 | 0 |
| gemini_flash_go | `replace` | 7.5 | 15 | 0 |
| gemini_flash_go | `write_file` | 3.0 | 6 | 0 |
| gemini_flash_go | `read_file` | 1.5 | 3 | 0 |
| gemini_flash_go_godoctor | `run_shell_command` | 10.5 | 21 | 0 |
| gemini_flash_go_godoctor | `write_file` | 4.0 | 8 | 0 |
| gemini_flash_go_godoctor | `replace` | 2.0 | 4 | 0 |
| gemini_flash_go_godoctor | `read_file` | 1.0 | 2 | 0 |
| gemini_flash_go_godoctor_agents | `run_shell_command` | 12.5 | 25 | 0 |
| gemini_flash_go_godoctor_agents | `write_file` | 5.0 | 10 | 0 |
| gemini_flash_go_godoctor_agents | `read_file` | 3.0 | 6 | 0 |
| gemini_flash_go_godoctor_agents | `write_todos` | 2.5 | 5 | 0 |
| gemini_flash_go_godoctor_agents | `replace` | 2.0 | 4 | 0 |
| gemini_flash_godoctor | `run_shell_command` | 10.5 | 21 | 0 |
| gemini_flash_godoctor | `replace` | 6.5 | 13 | 0 |
| gemini_flash_godoctor | `write_file` | 3.0 | 6 | 0 |
| gemini_flash_godoctor_agents | `replace` | 17.0 | 34 | 0 |
| gemini_flash_godoctor_agents | `run_shell_command` | 6.0 | 12 | 0 |
| gemini_flash_godoctor_agents | `write_file` | 4.5 | 9 | 0 |
| gemini_flash_godoctor_agents | `write_todos` | 3.0 | 6 | 0 |
| gemini_flash_godoctor_agents | `read_file` | 2.5 | 5 | 0 |

**legacy_refactor**

| Alternative | Tool | Avg Calls | Total | Failed |
|---|---|---|---|---|
| gemini_flash | `run_shell_command` | 6.5 | 13 | 0 |
| gemini_flash | `write_file` | 4.0 | 8 | 0 |
| gemini_flash | `read_file` | 2.0 | 4 | 0 |
| gemini_flash | `replace` | 1.5 | 3 | 0 |
| gemini_flash_go | `run_shell_command` | 5.5 | 11 | 0 |
| gemini_flash_go | `write_file` | 5.0 | 10 | 0 |
| gemini_flash_go | `read_file` | 2.5 | 5 | 0 |
| gemini_flash_go | `replace` | 2.0 | 4 | 0 |
| gemini_flash_go | `list_directory` | 0.5 | 1 | 0 |
| gemini_flash_go_godoctor | `run_shell_command` | 7.0 | 14 | 0 |
| gemini_flash_go_godoctor | `replace` | 5.0 | 10 | 0 |
| gemini_flash_go_godoctor | `write_todos` | 4.5 | 9 | 0 |
| gemini_flash_go_godoctor | `write_file` | 4.0 | 8 | 0 |
| gemini_flash_go_godoctor | `read_file` | 2.0 | 4 | 0 |
| gemini_flash_go_godoctor_agents | `run_shell_command` | 6.0 | 12 | 0 |
| gemini_flash_go_godoctor_agents | `write_todos` | 4.0 | 8 | 0 |
| gemini_flash_go_godoctor_agents | `write_file` | 3.5 | 7 | 0 |
| gemini_flash_go_godoctor_agents | `read_file` | 2.5 | 5 | 0 |
| gemini_flash_go_godoctor_agents | `edit_code` | 2.5 | 5 | 0 |
| gemini_flash_godoctor | `run_shell_command` | 7.5 | 15 | 0 |
| gemini_flash_godoctor | `replace` | 5.0 | 10 | 0 |
| gemini_flash_godoctor | `write_file` | 4.5 | 9 | 0 |
| gemini_flash_godoctor | `read_file` | 2.5 | 5 | 0 |
| gemini_flash_godoctor | `list_directory` | 0.5 | 1 | 0 |
| gemini_flash_godoctor_agents | `run_shell_command` | 6.0 | 12 | 0 |
| gemini_flash_godoctor_agents | `write_file` | 4.0 | 8 | 0 |
| gemini_flash_godoctor_agents | `read_file` | 2.0 | 4 | 0 |
| gemini_flash_godoctor_agents | `replace` | 2.0 | 4 | 0 |

**test_generation**

| Alternative | Tool | Avg Calls | Total | Failed |
|---|---|---|---|---|
| gemini_flash | `run_shell_command` | 5.5 | 11 | 0 |
| gemini_flash | `replace` | 3.5 | 7 | 0 |
| gemini_flash | `write_file` | 3.0 | 6 | 0 |
| gemini_flash | `read_file` | 1.0 | 2 | 0 |
| gemini_flash_go | `replace` | 4.5 | 9 | 0 |
| gemini_flash_go | `run_shell_command` | 4.0 | 8 | 0 |
| gemini_flash_go | `read_file` | 1.5 | 3 | 0 |
| gemini_flash_go | `write_file` | 1.0 | 2 | 0 |
| gemini_flash_go_godoctor | `run_shell_command` | 3.5 | 7 | 0 |
| gemini_flash_go_godoctor | `replace` | 2.5 | 5 | 0 |
| gemini_flash_go_godoctor | `read_file` | 1.0 | 2 | 0 |
| gemini_flash_go_godoctor | `write_file` | 1.0 | 2 | 0 |
| gemini_flash_go_godoctor_agents | `run_shell_command` | 4.0 | 8 | 0 |
| gemini_flash_go_godoctor_agents | `edit_code` | 2.0 | 4 | 0 |
| gemini_flash_go_godoctor_agents | `replace` | 1.5 | 3 | 0 |
| gemini_flash_go_godoctor_agents | `write_file` | 1.0 | 2 | 0 |
| gemini_flash_go_godoctor_agents | `read_file` | 1.0 | 2 | 0 |
| gemini_flash_godoctor | `run_shell_command` | 5.5 | 11 | 0 |
| gemini_flash_godoctor | `write_file` | 4.5 | 9 | 0 |
| gemini_flash_godoctor | `read_file` | 1.0 | 2 | 0 |
| gemini_flash_godoctor_agents | `run_shell_command` | 5.0 | 10 | 0 |
| gemini_flash_godoctor_agents | `replace` | 4.5 | 9 | 0 |
| gemini_flash_godoctor_agents | `write_file` | 3.5 | 7 | 0 |
| gemini_flash_godoctor_agents | `read_file` | 2.5 | 5 | 0 |
| gemini_flash_godoctor_agents | `list_directory` | 0.5 | 1 | 0 |

## ⚠️ Failure Analysis

### Timeouts (6 runs)
- **gemini_flash_godoctor**: 1 timeout(s)
- **gemini_flash_godoctor_agents**: 1 timeout(s)
- **gemini_flash_go_godoctor_agents**: 2 timeout(s)
- **gemini_flash**: 1 timeout(s)
- **gemini_flash_go_godoctor**: 1 timeout(s)

#### Timeout Logs

<details>
<summary><b>gemini_flash_go_godoctor_agents | api_implementation | Rep 1</b></summary>

```
YOLO mode is enabled. All tool calls will be automatically approved.
Session cleanup disabled: Either maxAge or maxCount must be specified
Server 'godoctor' supports tool updates. Listening for changes...
Server 'godoctor' supports resource updates. Listening for changes...
[STARTUP] StartupProfiler.flush() called with 9 phases
[STARTUP] Recording metric for phase: cli_startup duration: 34.48316699999998
[STARTUP] Recording metric for phase: load_settings duration: 1.7779590000000098
[STARTUP] Recording metric for phase: migrate_settings duration: 0.4651250000000573
[STARTUP] Recording metric for phase: parse_arguments duration: 7.448166000000015
[STARTUP] Recording metric for phase: load_cli_config duration: 15.142749999999978
[STARTUP] Recording metric for phase: initialize_app duration: 4.026083999999969
[STARTUP] Recording metric for phase: authenticate duration: 1.335875000000101
[STARTUP] Recording metric for phase: discover_tools duration: 2.9446249999999736
[STARTUP] Recording metric for phase: initialize_mcp_clients duration: 17.92983300000003
```
</details>

<details>
<summary><b>gemini_flash | greenfield_cli | Rep 1</b></summary>

```
YOLO mode is enabled. All tool calls will be automatically approved.
Session cleanup disabled: Either maxAge or maxCount must be specified
[STARTUP] StartupProfiler.flush() called with 9 phases
[STARTUP] Recording metric for phase: cli_startup duration: 69.74783300000001
[STARTUP] Recording metric for phase: load_settings duration: 1.7272919999999203
[STARTUP] Recording metric for phase: migrate_settings duration: 0.5693749999999227
[STARTUP] Recording metric for phase: parse_arguments duration: 15.012958000000026
[STARTUP] Recording metric for phase: load_cli_config duration: 34.79783300000008
[STARTUP] Recording metric for phase: initialize_app duration: 9.397041000000058
[STARTUP] Recording metric for phase: authenticate duration: 3.0101250000000164
[STARTUP] Recording metric for phase: discover_tools duration: 6.814665999999988
[STARTUP] Recording metric for phase: initialize_mcp_clients duration: 0.2757080000000087
Error executing tool replace: params must have required property 'old_string'
```
</details>

<details>
<summary><b>gemini_flash_go_godoctor | greenfield_cli | Rep 2</b></summary>

```
YOLO mode is enabled. All tool calls will be automatically approved.
Session cleanup disabled: Either maxAge or maxCount must be specified
Server 'godoctor' supports tool updates. Listening for changes...
Server 'godoctor' supports resource updates. Listening for changes...
[STARTUP] StartupProfiler.flush() called with 9 phases
[STARTUP] Recording metric for phase: cli_startup duration: 57.152917
[STARTUP] Recording metric for phase: load_settings duration: 1.967916999999943
[STARTUP] Recording metric for phase: migrate_settings duration: 0.6881250000000136
[STARTUP] Recording metric for phase: parse_arguments duration: 8.935958999999912
[STARTUP] Recording metric for phase: load_cli_config duration: 21.331457999999998
[STARTUP] Recording metric for phase: initialize_app duration: 7.572999999999979
[STARTUP] Recording metric for phase: authenticate duration: 3.272291999999993
[STARTUP] Recording metric for phase: discover_tools duration: 10.229166999999961
[STARTUP] Recording metric for phase: initialize_mcp_clients duration: 30.78083300000003
```
</details>

<details>
<summary><b>gemini_flash_go_godoctor_agents | greenfield_cli | Rep 2</b></summary>

```
YOLO mode is enabled. All tool calls will be automatically approved.
Session cleanup disabled: Either maxAge or maxCount must be specified
Server 'godoctor' supports tool updates. Listening for changes...
Server 'godoctor' supports resource updates. Listening for changes...
[STARTUP] StartupProfiler.flush() called with 9 phases
[STARTUP] Recording metric for phase: cli_startup duration: 37.04449999999997
[STARTUP] Recording metric for phase: load_settings duration: 1.9505420000000413
[STARTUP] Recording metric for phase: migrate_settings duration: 0.5154999999999745
[STARTUP] Recording metric for phase: parse_arguments duration: 8.100791999999956
[STARTUP] Recording metric for phase: load_cli_config duration: 16.720208999999954
[STARTUP] Recording metric for phase: initialize_app duration: 4.015292000000045
[STARTUP] Recording metric for phase: authenticate duration: 1.3992499999999382
[STARTUP] Recording metric for phase: discover_tools duration: 3.0332499999999527
[STARTUP] Recording metric for phase: initialize_mcp_clients duration: 17.975041000000033
```
</details>

<details>
<summary><b>gemini_flash_godoctor | greenfield_cli | Rep 1</b></summary>

```
YOLO mode is enabled. All tool calls will be automatically approved.
Session cleanup disabled: Either maxAge or maxCount must be specified
Server 'godoctor' supports tool updates. Listening for changes...
Server 'godoctor' supports resource updates. Listening for changes...
[STARTUP] StartupProfiler.flush() called with 9 phases
[STARTUP] Recording metric for phase: cli_startup duration: 32.91004099999998
[STARTUP] Recording metric for phase: load_settings duration: 1.8109580000000278
[STARTUP] Recording metric for phase: migrate_settings duration: 0.47041600000000017
[STARTUP] Recording metric for phase: parse_arguments duration: 7.238707999999974
[STARTUP] Recording metric for phase: load_cli_config duration: 13.726292000000058
[STARTUP] Recording metric for phase: initialize_app duration: 4.086125000000038
[STARTUP] Recording metric for phase: authenticate duration: 1.3707500000000437
[STARTUP] Recording metric for phase: discover_tools duration: 2.7684169999999995
[STARTUP] Recording metric for phase: initialize_mcp_clients duration: 19.396624999999972
```
</details>

<details>
<summary><b>gemini_flash_godoctor_agents | greenfield_cli | Rep 1</b></summary>

... (truncated) ...
```
				if e, ok := r.(int); ok {
					exitCode = e
				}
			}
		}()
		oldExit := os.Exit
		defer func() { os.Exit = oldExit }()
		os.Exit = func(code int) {
			panic(code)
		}
		main()
	}

	out, _ := captureOutput(mainFunc)
	expected := `{"error":"unknown command"}\n`
	if out != expected || exitCode != 42 {
		t.Errorf("Expected %q with exit code 42, got %q with exit code %d", expected, out, exitCode)
	}
}
) in /Users/petruzalek/projects/tenkai/experiments/runs/20251217-174107_gemini_go_mcp_flash/gemini_flash_godoctor_agents/greenfield_cli/rep-1/main_test.go. No edits made. The exact text in old_string was not found. Ensure you're not escaping content incorrectly and check whitespace, indentation, and context. Use read_file tool to verify.
```
</details>

### Loop Warnings (2 runs)
- **gemini_flash**: 1 loop(s)
- **gemini_flash_go**: 1 loop(s)

### Problem Scenarios

## 🎯 Conclusion

(This section tracks success rates and performance metrics filtered for successful runs only.)

Based on the analysis, **gemini_flash_go** achieved the highest success rate of 100%.
