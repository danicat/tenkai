# Tenkai Experiment Report
_Generated at: Thu, 18 Dec 2025 15:26:26 GMT_

## 📊 Executive Summary

### Key Metrics
- **Total Runs**: 20
- **Success Rate**: 90.0% (18/20 completed)
- **Test Pass Rate**: 100.0% (17/17)
- **Timeouts**: 2 runs (10.0%)
- **Avg Tokens (Success Only)**: 297.6k
- **Avg Lint Issues (Success Only)**: 3.2

## 🧪 Experiment Parameters
> Part 1: Optimizing the system prompt for Go and GoDoctor.  This experiment evaluates how specific prompt construction affects an agent's  idiomatic usage of Go and its willingness to discover and use specialized MCP tools.


| Parameter | Value |
| :--- | :--- |
| **Repetitions** | 1 |
| **Concurrent Workers** | 3 |
| **Total Scenarios** | 5 |

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
| **1_default** | 1m7s (Control) 🏆 | 267.3k | 2.5 | 80% | 100% | 100.0% | ✅ ⏱️ |
| **2_go_optimized** | 1m15s (+11%) | 212.2k (-21%) 🏆 | 2.5 | 80% | 100% | 100.0% | ✅ ⏱️ |
| **3_default_godoctor** | 1m46s (+57%) | 374.8k (+40%) | 3.0 | 100% | 80% | 100.0% | ✅ |
| **4_go_godoctor** | 1m44s (+54%) | 313.0k (+17%) | 4.4 | 100% | 100% | 100.0% | ✅ |

### 💡 Efficiency Highlights
- **Most Tool-Efficient**: 2_go_optimized (14.0 calls/success run) vs 3_default_godoctor (18.2 calls/success run)
- **Most Token-Efficient**: 2_go_optimized (212.2k/run) vs 3_default_godoctor (374.8k/run)
- **Best Tool Accuracy**: 1_default (100.0% success) vs 4_go_godoctor (100.0% success)

### 📊 Statistical Analysis Notes
- `*` Significant ($p < 0.05$): 95% confidence that the difference is not due to chance.
- `**` Highly Significant ($p < 0.01$): 99% confidence.
- Metrics compared using Welch's t-test (Duration, Tokens, Lint) and Fisher's Exact Test (Success Rate).

## 📈 Summary by Scenario & Alternative
| Scenario | Alternative | Runs | Avg Duration (Success) | Avg Tokens (Success) | Lint | Pass Rate |
|---|---|---|---|---|---|---|
| api_implementation | 1_default | 1 | 54.023s | 256.9k | 3.0 | 100% |
| api_implementation | 2_go_optimized | 1 | **53.567s** | **101.1k** | 4.0 | 100% |
| api_implementation | 3_default_godoctor | 1 | 2m51.585s | 693.5k | 7.0 | 100% |
| api_implementation | 4_go_godoctor | 1 | 1m44.459s | 304.2k | 4.0 | 100% |
| concurrency_debug | 1_default | 1 | 1m8.842s | 178.6k | 0.0 | 100% |
| concurrency_debug | 2_go_optimized | 1 | **29.955s** | **91.0k** | 0.0 | 100% |
| concurrency_debug | 3_default_godoctor | 1 | 31.878s | 96.3k | 0.0 | 100% |
| concurrency_debug | 4_go_godoctor | 1 | 1m0.47s | 201.7k | 0.0 | 100% |
| greenfield_cli | 1_default | 1 | - | - | 0.0 | 0% |
| greenfield_cli | 2_go_optimized | 1 | - | - | 0.0 | 0% |
| greenfield_cli | 3_default_godoctor | 1 | **3m40.503s** | **663.3k** | 4.0 | 100% |
| greenfield_cli | 4_go_godoctor | 1 | 4m12.72s | 746.5k | 6.0 | 100% |
| legacy_refactor | 1_default | 1 | 1m20.624s | 407.1k | 6.0 | 100% |
| legacy_refactor | 2_go_optimized | 1 | 1m40.16s | 416.7k | 5.0 | 100% |
| legacy_refactor | 3_default_godoctor | 1 | **50.95s** | 265.7k | 3.0 | 0% |
| legacy_refactor | 4_go_godoctor | 1 | 56.635s | **211.7k** | 9.0 | 100% |
| test_generation | 1_default | 1 | 1m5.764s | 226.7k | 1.0 | 100% |
| test_generation | 2_go_optimized | 1 | 1m55.093s | 240.2k | 1.0 | 100% |
| test_generation | 3_default_godoctor | 1 | 53.955s | 155.0k | 1.0 | 100% |
| test_generation | 4_go_godoctor | 1 | **43.622s** | **101.0k** | 3.0 | 100% |

## 📋 Detailed Metrics

<details>
<summary>Click to expand detailed run-by-run metrics</summary>

| Scenario | Alternative | Rep | Duration | Tokens | Tools | Test | Lint | Status |
|---|---|---|---|---|---|---|---|---|
| api_implementation | 1_default | 1 | 54.023s | 256.9k | 16 | PASS | 3 | ✅ |
| api_implementation | 2_go_optimized | 1 | 53.567s | 101.1k | 7 | PASS | 4 | ✅ |
| api_implementation | 3_default_godoctor | 1 | 2m51.585s | 693.5k | 30 | PASS | 7 | ✅ |
| api_implementation | 4_go_godoctor | 1 | 1m44.459s | 304.2k | 18 | PASS | 4 | ✅ |
| concurrency_debug | 1_default | 1 | 1m8.842s | 178.6k | 11 | PASS | 0 | ✅ |
| concurrency_debug | 2_go_optimized | 1 | 29.955s | 91.0k | 8 | PASS | 0 | ✅ |
| concurrency_debug | 3_default_godoctor | 1 | 31.878s | 96.3k | 7 | PASS | 0 | ✅ |
| concurrency_debug | 4_go_godoctor | 1 | 1m0.47s | 201.7k | 15 | PASS | 0 | ✅ |
| greenfield_cli | 1_default | 1 | 5m0.023s | 0 | 20 | - | 0 | ⏱️ Timeout |
| greenfield_cli | 2_go_optimized | 1 | 5m0.022s | 0 | 34 | - | 0 | ⏱️ Timeout |
| greenfield_cli | 3_default_godoctor | 1 | 3m40.503s | 663.3k | 28 | PASS | 4 | ✅ |
| greenfield_cli | 4_go_godoctor | 1 | 4m12.72s | 746.5k | 31 | PASS | 6 | ✅ |
| legacy_refactor | 1_default | 1 | 1m20.624s | 407.1k | 25 | PASS | 6 | ✅ |
| legacy_refactor | 2_go_optimized | 1 | 1m40.16s | 416.7k | 29 | PASS | 5 | ✅ |
| legacy_refactor | 3_default_godoctor | 1 | 50.95s | 265.7k | 17 | SKIP | 3 | ⚠️ Skip |
| legacy_refactor | 4_go_godoctor | 1 | 56.635s | 211.7k | 16 | PASS | 9 | ✅ |
| test_generation | 1_default | 1 | 1m5.764s | 226.7k | 12 | PASS | 1 | ✅ |
| test_generation | 2_go_optimized | 1 | 1m55.093s | 240.2k | 12 | PASS | 1 | ✅ |
| test_generation | 3_default_godoctor | 1 | 53.955s | 155.0k | 9 | PASS | 1 | ✅ |
| test_generation | 4_go_godoctor | 1 | 43.622s | 101.0k | 7 | PASS | 3 | ✅ |

</details>

## 🔧 Tool Usage Analysis

### Overall Summary

| Alternative | Total Calls | Failed Calls | Unique Tools | Avg Calls/Run | Most Used Tool |
|---|---|---|---|---|---|
| 2_go_optimized | 90 | 0 | 5 | 18.0 | `replace` (36) |
| 3_default_godoctor | 91 | 0 | 5 | 18.2 | `replace` (36) |
| 4_go_godoctor | 87 | 0 | 6 | 17.4 | `replace` (36) |
| 1_default | 84 | 0 | 6 | 16.8 | `write_todos` (23) |

### Tool Availability

### Tool Usage by Scenario

**api_implementation**

| Alternative | Tool | Avg Calls | Total | Failed |
|---|---|---|---|---|
| 1_default | `write_todos` | 8.0 | 8 | 0 |
| 1_default | `run_shell_command` | 5.0 | 5 | 0 |
| 1_default | `write_file` | 2.0 | 2 | 0 |
| 1_default | `read_file` | 1.0 | 1 | 0 |
| 2_go_optimized | `run_shell_command` | 4.0 | 4 | 0 |
| 2_go_optimized | `write_file` | 2.0 | 2 | 0 |
| 2_go_optimized | `read_file` | 1.0 | 1 | 0 |
| 3_default_godoctor | `run_shell_command` | 11.0 | 11 | 0 |
| 3_default_godoctor | `replace` | 11.0 | 11 | 0 |
| 3_default_godoctor | `read_file` | 6.0 | 6 | 0 |
| 3_default_godoctor | `write_file` | 2.0 | 2 | 0 |
| 4_go_godoctor | `replace` | 6.0 | 6 | 0 |
| 4_go_godoctor | `run_shell_command` | 6.0 | 6 | 0 |
| 4_go_godoctor | `read_file` | 3.0 | 3 | 0 |
| 4_go_godoctor | `write_file` | 2.0 | 2 | 0 |
| 4_go_godoctor | `list_directory` | 1.0 | 1 | 0 |

**concurrency_debug**

| Alternative | Tool | Avg Calls | Total | Failed |
|---|---|---|---|---|
| 1_default | `replace` | 5.0 | 5 | 0 |
| 1_default | `run_shell_command` | 4.0 | 4 | 0 |
| 1_default | `read_file` | 2.0 | 2 | 0 |
| 2_go_optimized | `replace` | 4.0 | 4 | 0 |
| 2_go_optimized | `read_file` | 2.0 | 2 | 0 |
| 2_go_optimized | `run_shell_command` | 2.0 | 2 | 0 |
| 3_default_godoctor | `replace` | 3.0 | 3 | 0 |
| 3_default_godoctor | `run_shell_command` | 2.0 | 2 | 0 |
| 3_default_godoctor | `list_directory` | 1.0 | 1 | 0 |
| 3_default_godoctor | `read_file` | 1.0 | 1 | 0 |
| 4_go_godoctor | `replace` | 9.0 | 9 | 0 |
| 4_go_godoctor | `run_shell_command` | 4.0 | 4 | 0 |
| 4_go_godoctor | `read_file` | 2.0 | 2 | 0 |

**greenfield_cli**

| Alternative | Tool | Avg Calls | Total | Failed |
|---|---|---|---|---|
| 1_default | `replace` | 9.0 | 9 | 0 |
| 1_default | `write_todos` | 7.0 | 7 | 0 |
| 1_default | `run_shell_command` | 2.0 | 2 | 0 |
| 1_default | `write_file` | 2.0 | 2 | 0 |
| 2_go_optimized | `replace` | 23.0 | 23 | 0 |
| 2_go_optimized | `write_todos` | 5.0 | 5 | 0 |
| 2_go_optimized | `run_shell_command` | 3.0 | 3 | 0 |
| 2_go_optimized | `write_file` | 2.0 | 2 | 0 |
| 2_go_optimized | `read_file` | 1.0 | 1 | 0 |
| 3_default_godoctor | `replace` | 15.0 | 15 | 0 |
| 3_default_godoctor | `run_shell_command` | 10.0 | 10 | 0 |
| 3_default_godoctor | `write_file` | 2.0 | 2 | 0 |
| 3_default_godoctor | `read_file` | 1.0 | 1 | 0 |
| 4_go_godoctor | `replace` | 17.0 | 17 | 0 |
| 4_go_godoctor | `write_todos` | 7.0 | 7 | 0 |
| 4_go_godoctor | `run_shell_command` | 5.0 | 5 | 0 |
| 4_go_godoctor | `write_file` | 2.0 | 2 | 0 |

**legacy_refactor**

| Alternative | Tool | Avg Calls | Total | Failed |
|---|---|---|---|---|
| 1_default | `write_todos` | 8.0 | 8 | 0 |
| 1_default | `run_shell_command` | 7.0 | 7 | 0 |
| 1_default | `write_file` | 6.0 | 6 | 0 |
| 1_default | `replace` | 2.0 | 2 | 0 |
| 1_default | `read_file` | 2.0 | 2 | 0 |
| 2_go_optimized | `run_shell_command` | 15.0 | 15 | 0 |
| 2_go_optimized | `replace` | 7.0 | 7 | 0 |
| 2_go_optimized | `write_file` | 5.0 | 5 | 0 |
| 2_go_optimized | `read_file` | 2.0 | 2 | 0 |
| 3_default_godoctor | `run_shell_command` | 7.0 | 7 | 0 |
| 3_default_godoctor | `replace` | 4.0 | 4 | 0 |
| 3_default_godoctor | `write_file` | 3.0 | 3 | 0 |
| 3_default_godoctor | `read_file` | 3.0 | 3 | 0 |
| 4_go_godoctor | `run_shell_command` | 6.0 | 6 | 0 |
| 4_go_godoctor | `write_file` | 5.0 | 5 | 0 |
| 4_go_godoctor | `read_file` | 2.0 | 2 | 0 |
| 4_go_godoctor | `replace` | 2.0 | 2 | 0 |
| 4_go_godoctor | `list_directory` | 1.0 | 1 | 0 |

**test_generation**

| Alternative | Tool | Avg Calls | Total | Failed |
|---|---|---|---|---|
| 1_default | `replace` | 4.0 | 4 | 0 |
| 1_default | `read_file` | 3.0 | 3 | 0 |
| 1_default | `run_shell_command` | 2.0 | 2 | 0 |
| 1_default | `write_file` | 2.0 | 2 | 0 |
| 1_default | `list_directory` | 1.0 | 1 | 0 |
| 2_go_optimized | `run_shell_command` | 5.0 | 5 | 0 |
| 2_go_optimized | `write_file` | 3.0 | 3 | 0 |
| 2_go_optimized | `read_file` | 2.0 | 2 | 0 |
| 2_go_optimized | `replace` | 2.0 | 2 | 0 |
| 3_default_godoctor | `run_shell_command` | 4.0 | 4 | 0 |
| 3_default_godoctor | `replace` | 3.0 | 3 | 0 |
| 3_default_godoctor | `write_file` | 1.0 | 1 | 0 |
| 3_default_godoctor | `read_file` | 1.0 | 1 | 0 |
| 4_go_godoctor | `run_shell_command` | 3.0 | 3 | 0 |
| 4_go_godoctor | `replace` | 2.0 | 2 | 0 |
| 4_go_godoctor | `read_file` | 1.0 | 1 | 0 |
| 4_go_godoctor | `write_file` | 1.0 | 1 | 0 |

## ⚠️ Failure Analysis

### Timeouts (2 runs)
- **2_go_optimized**: 1 timeout(s)
- **1_default**: 1 timeout(s)

#### Timeout Logs

<details>
<summary><b>1_default | greenfield_cli | Rep 1</b></summary>

```
YOLO mode is enabled. All tool calls will be automatically approved.
Session cleanup disabled: Either maxAge or maxCount must be specified
Server 'godoctor' supports tool updates. Listening for changes...
Server 'godoctor' supports resource updates. Listening for changes...
[STARTUP] StartupProfiler.flush() called with 9 phases
[STARTUP] Recording metric for phase: cli_startup duration: 38.02704200000005
[STARTUP] Recording metric for phase: load_settings duration: 1.8507919999999558
[STARTUP] Recording metric for phase: migrate_settings duration: 0.6343749999999773
[STARTUP] Recording metric for phase: parse_arguments duration: 8.945917000000009
[STARTUP] Recording metric for phase: load_cli_config duration: 15.567542000000003
[STARTUP] Recording metric for phase: initialize_app duration: 4.408250000000066
[STARTUP] Recording metric for phase: authenticate duration: 1.6734579999999823
[STARTUP] Recording metric for phase: discover_tools duration: 3.286041999999952
[STARTUP] Recording metric for phase: initialize_mcp_clients duration: 17.930082999999968
```
</details>

<details>
<summary><b>2_go_optimized | greenfield_cli | Rep 1</b></summary>

... (truncated) ...
```
				t.Fatalf("unexpected panic: %v", r)
			}
		}
	}()

	captureOutput()
	listTasks()
	output := stopCaptureOutput()

	expectedTasks := []Task{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 3, Description: "Task 3", Done: false},
	}
	expectedJSONBytes, _ := json.Marshal(expectedTasks)
	expectedJSON := string(expectedJSONBytes) + "\n"

	if output != expectedJSON {
		t.Errorf("listTasks() got %q, want %q", output, expectedJSON)
	}
}) in /Users/petruzalek/projects/tenkai/experiments/runs/20251218-151200_system_prompt_study/2_go_optimized/greenfield_cli/rep-1/main_test.go. No edits made. The exact text in old_string was not found. Ensure you're not escaping content incorrectly and check whitespace, indentation, and context. Use read_file tool to verify.
```
</details>

### Problem Scenarios

## 🎯 Conclusion

(This section tracks success rates and performance metrics filtered for successful runs only.)

Based on the analysis, **3_default_godoctor** achieved the highest success rate of 100%.
