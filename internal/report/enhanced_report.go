package report

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/danicat/tenkai/internal/config"
	"github.com/danicat/tenkai/internal/runner"
)

// EnhancedReporter extends Reporter with better formatting
type EnhancedMarkdownGenerator struct {
	Results []runner.Result
	Cfg     *config.Configuration
	Notes   []string
	w       io.Writer
	stats   []avgMetrics
}

func (r *EnhancedMarkdownGenerator) Generate() error {
	// Sort results
	sort.Slice(r.Results, func(i, j int) bool {
		if r.Results[i].Scenario != r.Results[j].Scenario {
			return r.Results[i].Scenario < r.Results[j].Scenario
		}
		return r.Results[i].Alternative < r.Results[j].Alternative
	})

	r.calculateStats()

	r.writeHeader()
	r.writeExecutiveSummary()
	r.writeExperimentParameters()
	r.writeNotes()
	r.writePerformanceComparison()
	r.writeSummaryByScenario()
	r.writeDetailedMetrics()
	r.writeEvaluationDetails()
	r.writeToolUsage()
	r.writeFailureAnalysis()
	r.writeConclusion()

	return nil
}

func (r *EnhancedMarkdownGenerator) calculateStats() {
	// Group by alternative
	type altMetrics struct {
		totalRuns      int
		successes      int
		timeouts       int
		totalTokens    int64
		totalLint      int
		totalDur       time.Duration
		testsPassed    int
		testsMeasured  int
		totalToolCalls int
		failedTools    int

		// Raw data for stats
		durations []float64
		tokens    []float64
		lints     []float64
	}

	byAlt := make(map[string]*altMetrics)
	var altNames []string
	if r.Cfg != nil {
		for _, alt := range r.Cfg.Alternatives {
			altNames = append(altNames, alt.Name)
			byAlt[alt.Name] = &altMetrics{}
		}
	}

	for _, res := range r.Results {
		if _, ok := byAlt[res.Alternative]; !ok {
			byAlt[res.Alternative] = &altMetrics{}
			found := false
			for _, n := range altNames {
				if n == res.Alternative {
					found = true
					break
				}
			}
			if !found {
				altNames = append(altNames, res.Alternative)
			}
		}
		m := byAlt[res.Alternative]
		m.totalRuns++

		// A run is successful ONLY if there's no technical error AND it actually passed at least one test.
		// Runs that skip evaluation (e.g. no code generated) or fail all tests are not successes.
		isSuccess := res.Error == nil
		if res.EvaluationMetrics != nil && res.EvaluationMetrics.TestsPassed == 0 {
			isSuccess = false
		}

		if isSuccess {
			m.successes++
			m.totalDur += res.Duration // Only add duration for success
			m.durations = append(m.durations, float64(res.Duration.Seconds()))
			if res.AgentMetrics != nil {
				m.totalTokens += int64(res.AgentMetrics.TotalTokens)
				m.tokens = append(m.tokens, float64(res.AgentMetrics.TotalTokens))
			}
			if res.EvaluationMetrics != nil {
				m.totalLint += res.EvaluationMetrics.LintIssues
				m.lints = append(m.lints, float64(res.EvaluationMetrics.LintIssues))
			}
		} else if res.Error != nil && strings.Contains(res.Error.Error(), "timeout") {
			m.timeouts++
		}

		if res.Error == nil && res.AgentMetrics != nil {
			for _, tc := range res.AgentMetrics.ToolCalls {
				m.totalToolCalls++
				if tc.Status != "success" {
					m.failedTools++
				}
			}
		}

		if res.EvaluationMetrics != nil {
			m.testsMeasured++
			if res.EvaluationMetrics.TestsPassed > 0 {
				m.testsPassed++
			}
		}
	}

	// Calculate averages
	r.stats = []avgMetrics{}
	for _, name := range altNames {
		m := byAlt[name]
		if m.totalRuns == 0 {
			continue
		}

		toolSuccessRate := 0.0
		if m.totalToolCalls > 0 {
			toolSuccessRate = float64(m.totalToolCalls-m.failedTools) / float64(m.totalToolCalls) * 100
		}

		testPassRate := 0.0
		if m.testsMeasured > 0 {
			testPassRate = float64(m.testsPassed) / float64(m.testsMeasured) * 100
		}

		avgDur := time.Duration(0)
		avgTokens := int64(0)
		avgLint := 0.0
		if m.successes > 0 {
			avgDur = m.totalDur / time.Duration(m.successes)
			avgTokens = m.totalTokens / int64(m.successes)
			avgLint = float64(m.totalLint) / float64(m.successes)
		}

		r.stats = append(r.stats, avgMetrics{
			name:           name,
			successRate:    float64(m.successes) / float64(m.totalRuns) * 100,
			avgDur:         avgDur,
			avgTokens:      avgTokens,
			avgLint:        avgLint,
			testPass:       testPassRate,
			toolSuccess:    toolSuccessRate,
			timeouts:       m.timeouts,
			totalRuns:      m.totalRuns,
			completedRuns:  m.successes,
			totalToolCalls: m.totalToolCalls,
			durations:      m.durations,
			tokens:         m.tokens,
			lints:          m.lints,
		})
	}

	// Identify Control (first in list)
	var control *avgMetrics
	if len(r.stats) > 0 {
		control = &r.stats[0]

		// Calculate P-values for alternatives vs control
		for i := 1; i < len(r.stats); i++ {
			a := &r.stats[i]
			a.pDur = WelchTTest(control.durations, a.durations)
			a.pTokens = WelchTTest(control.tokens, a.tokens)
			a.pLint = WelchTTest(control.lints, a.lints)

			// Fisher exact test
			a.pSucc = FisherExactTest(
				control.completedRuns, control.totalRuns-control.completedRuns,
				a.completedRuns, a.totalRuns-a.completedRuns,
			)

			// Statistical Analysis for Timeouts
			a.pTime = FisherExactTest(
				control.timeouts, control.totalRuns-control.timeouts,
				a.timeouts, a.totalRuns-a.timeouts,
			)
		}
	}
}

func (r *EnhancedMarkdownGenerator) writeHeader() {
	fmt.Fprintln(r.w, "# Tenkai Experiment Report")
	fmt.Fprintf(r.w, "_Generated at: %s_\n", time.Now().Format(time.RFC1123))
}

func (r *EnhancedMarkdownGenerator) writeExperimentParameters() {
	if r.Cfg == nil {
		return
	}
	fmt.Fprintln(r.w, "\n## 🧪 Experiment Parameters")
	if r.Cfg.Description != "" {
		fmt.Fprintf(r.w, "> %s\n\n", r.Cfg.Description)
	}

	fmt.Fprintln(r.w, "| Parameter | Value |")
	fmt.Fprintln(r.w, "| :--- | :--- |")
	fmt.Fprintf(r.w, "| **Repetitions** | %d |\n", r.Cfg.Repetitions)
	fmt.Fprintf(r.w, "| **Concurrent Workers** | %d |\n", r.Cfg.MaxConcurrent)
	fmt.Fprintf(r.w, "| **Total Scenarios** | %d |\n", len(r.Cfg.Scenarios))

	fmt.Fprintln(r.w, "\n### Alternatives Configuration")
	for i, alt := range r.Cfg.Alternatives {
		role := ""
		if i == 0 {
			role = " (Control 🏳️)"
		}
		fmt.Fprintf(r.w, "- **%s**%s\n", alt.Name, role)
		if alt.Description != "" {
			fmt.Fprintf(r.w, "  > %s\n", alt.Description)
		}
		fmt.Fprintf(r.w, "  - Model: `%s`\n", getValue(alt.Args, "--model"))
		if alt.SystemPromptFile != "" {
			fmt.Fprintf(r.w, "  - System Prompt: `%s`\n", filepath.Base(alt.SystemPromptFile))
		} else if alt.Env != nil && alt.Env["GEMINI_SYSTEM_MD"] != "" {
			path := alt.Env["GEMINI_SYSTEM_MD"]
			fmt.Fprintf(r.w, "  - System Prompt: `%s`\n", filepath.Base(path))
		}

		if alt.SettingsPath != "" {
			fmt.Fprintf(r.w, "  - MCP Settings: `%s`\n", filepath.Base(alt.SettingsPath))
		}
	}
}

func getValue(args []string, key string) string {
	for i, arg := range args {
		if arg == key && i+1 < len(args) {
			return args[i+1]
		}
	}
	return "default"
}

func (r *EnhancedMarkdownGenerator) writeNotes() {
	if len(r.Notes) == 0 {
		return
	}
	fmt.Fprintln(r.w, "\n## 📝 Notes")
	for _, note := range r.Notes {
		fmt.Fprintf(r.w, "> %s\n\n", note)
	}
}

func (r *EnhancedMarkdownGenerator) writeExecutiveSummary() {
	fmt.Fprintln(r.w, "\n## 📊 Executive Summary")

	if len(r.stats) == 0 {
		fmt.Fprintln(r.w, "No data available.")
		return
	}

	// 1. Overall Health
	totalRuns := len(r.Results)
	successCount := 0
	for _, res := range r.Results {
		if res.Error == nil {
			successCount++
		}
	}
	globalSuccessRate := float64(successCount) / float64(totalRuns) * 100

	healthEmoji := "🟢"
	if globalSuccessRate < 80 {
		healthEmoji = "🟡"
	}
	if globalSuccessRate < 50 {
		healthEmoji = "🔴"
	}

	fmt.Fprintf(r.w, "This experiment executed **%d runs** across **%d scenarios** with an overall success rate of **%.1f%%** %s.\n\n",
		totalRuns, len(r.Cfg.Scenarios), globalSuccessRate, healthEmoji)

	// 2. Identify Winners
	var bestSuccess *avgMetrics
	var bestSpeed *avgMetrics
	var bestEfficiency *avgMetrics

	for i := range r.stats {
		s := &r.stats[i]
		if bestSuccess == nil || s.successRate > bestSuccess.successRate {
			bestSuccess = s
		}
		if s.completedRuns > 0 {
			if bestSpeed == nil || s.avgDur < bestSpeed.avgDur {
				bestSpeed = s
			}
			if bestEfficiency == nil || s.avgTokens < bestEfficiency.avgTokens {
				bestEfficiency = s
			}
		}
	}

	if bestSuccess != nil {
		fmt.Fprintln(r.w, "### 🏆 Key Findings")
		fmt.Fprintf(r.w, "- **Most Reliable**: **%s** passed %.1f%% of runs.\n", bestSuccess.name, bestSuccess.successRate)

		if bestSpeed != nil {
			fmt.Fprintf(r.w, "- **Fastest**: **%s** averaged %s per successful run.\n", bestSpeed.name, bestSpeed.avgDur.Round(time.Millisecond))
		}
		if bestEfficiency != nil {
			fmt.Fprintf(r.w, "- **Most Efficient**: **%s** used only %s tokens per run.\n", bestEfficiency.name, formatNumber(bestEfficiency.avgTokens))
		}
	}
}

func (r *EnhancedMarkdownGenerator) writeFailureAnalysis() {
	timeouts := []string{}
	failures := []string{}
	loops := []string{}

	for _, res := range r.Results {
		key := fmt.Sprintf("%s|%s|#%d", res.Alternative, res.Scenario, res.Repetition)
		if res.Error != nil && strings.Contains(res.Error.Error(), "timeout") {
			timeouts = append(timeouts, key)
		} else if res.Error != nil {
			failures = append(failures, key)
		}
		if res.AgentMetrics != nil && res.AgentMetrics.LoopDetected {
			loops = append(loops, key)
		}
	}

	if len(timeouts) > 0 || len(failures) > 0 || len(loops) > 0 {
		fmt.Fprintln(r.w, "\n## ⚠️ Failure Analysis")

		if len(timeouts) > 0 {
			fmt.Fprintf(r.w, "\n### Timeouts (%d runs)\n", len(timeouts))
			// Group by alternative
			byAlt := make(map[string]int)
			for _, t := range timeouts {
				alt := strings.Split(t, "|")[0]
				byAlt[alt]++
			}
			for alt, count := range byAlt {
				fmt.Fprintf(r.w, "- **%s**: %d timeout(s)\n", alt, count)
			}

			fmt.Fprintln(r.w, "\n#### Timeout Logs")
			for _, res := range r.Results {
				if res.Error != nil && strings.Contains(res.Error.Error(), "timeout") {
					fmt.Fprintf(r.w, "\n<details>\n<summary><b>%s | %s | Rep %d</b></summary>\n\n", res.Alternative, res.Scenario, res.Repetition)

					// Decide whether to show Stderr or parse Stdout (JSON events)
					stderrLines := strings.Split(strings.TrimSpace(res.Stderr), "\n")

					// If Stderr is short (likely just startup) and Stdout is large, try parsing JSON events
					if len(stderrLines) < 20 && len(res.Stdout) > 100 {
						fmt.Fprintln(r.w, "**Last 3 Conversation Turns (Parsed from JSONL):**")
						fmt.Fprintln(r.w, "```yaml")

						events := parseLastEvents(res.Stdout, 3)
						for _, evt := range events {
							fmt.Fprintln(r.w, evt)
						}
						fmt.Fprintln(r.w, "```")

						if len(stderrLines) > 0 && stderrLines[0] != "" {
							fmt.Fprintln(r.w, "\n**Stderr Tail:**")
							fmt.Fprintln(r.w, "```")
							for _, l := range stderrLines {
								fmt.Fprintln(r.w, l)
							}
							fmt.Fprintln(r.w, "```")
						}

					} else {
						// Fallback to standard Stderr tail
						start := 0
						if len(stderrLines) > 20 {
							start = len(stderrLines) - 20
							fmt.Fprintln(r.w, "... (truncated) ...")
						}
						fmt.Fprintln(r.w, "```")
						for i := start; i < len(stderrLines); i++ {
							fmt.Fprintln(r.w, stderrLines[i])
						}
						fmt.Fprintln(r.w, "```")
					}
					fmt.Fprintln(r.w, "</details>")
				}
			}

		}

		if len(loops) > 0 {
			fmt.Fprintf(r.w, "\n### Loop Warnings (%d runs)\n", len(loops))
			byAlt := make(map[string]int)
			for _, l := range loops {
				alt := strings.Split(l, "|")[0]
				byAlt[alt]++
			}
			for alt, count := range byAlt {
				fmt.Fprintf(r.w, "- **%s**: %d loop(s)\n", alt, count)
			}
		}

		// Scenario-level failure rates
		scenarioStats := make(map[string]struct{ total, failed int })
		for _, res := range r.Results {
			s := scenarioStats[res.Scenario]
			s.total++
			if res.EvaluationMetrics != nil {
				if res.EvaluationMetrics.TestsPassed == 0 && res.EvaluationMetrics.TestsFailed == 0 {
					s.failed++
				}
			}
			scenarioStats[res.Scenario] = s
		}

		var problemScenarios []string
		for scenario, stats := range scenarioStats {
			failRate := float64(stats.failed) / float64(stats.total) * 100
			if failRate > 50 {
				problemScenarios = append(problemScenarios, fmt.Sprintf("- **%s**: %.0f%% skip rate (likely broken scenario)", scenario, failRate))
			}
		}

		if len(problemScenarios) > 0 {
			fmt.Fprintln(r.w, "\n### Problem Scenarios")
			for _, line := range problemScenarios {
				fmt.Fprintln(r.w, line)
			}
		}
	}
}

// Helper to parse last N meaningful events from JSONL stdout
func parseLastEvents(stdout string, n int) []string {
	var events []string
	lines := strings.Split(stdout, "\n")

	// Read backwards
	count := 0
	type Event struct {
		Type     string `json:"type"`
		Role     string `json:"role"`
		Content  string `json:"content"`
		ToolName string `json:"tool_name"`
		Output   string `json:"output"`
		Status   string `json:"status"`
	}

	for i := len(lines) - 1; i >= 0 && count < n*2; i-- { // *2 to check more lines as some might be skipped
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		var evt Event
		if err := json.Unmarshal([]byte(line), &evt); err == nil {
			msg := ""
			switch evt.Type {
			case "message":
				msg = fmt.Sprintf("[%s]: %s", strings.ToUpper(evt.Role), shorten(evt.Content))
			case "tool_use":
				msg = fmt.Sprintf("[TOOL_CALL]: %s", evt.ToolName)
			case "tool_result":
				msg = fmt.Sprintf("[TOOL_RESULT]: %s (Status: %s)", shorten(evt.Output), evt.Status)
			}

			if msg != "" {
				events = append([]string{msg}, events...) // Prepend
				count++
			}
		}
	}

	// Return only requested N
	if len(events) > n {
		return events[len(events)-n:]
	}
	return events
}

func shorten(s string) string {
	if len(s) > 100 {
		return strings.ReplaceAll(s[:100], "\n", " ") + "..."
	}
	return strings.ReplaceAll(s, "\n", " ")
}

// avgMetrics holds calculated averages for an alternative
type avgMetrics struct {
	name           string
	successRate    float64
	avgDur         time.Duration
	avgTokens      int64
	avgLint        float64
	testPass       float64
	toolSuccess    float64
	timeouts       int
	totalRuns      int
	completedRuns  int
	totalToolCalls int

	// Raw data for stats
	durations []float64
	tokens    []float64
	lints     []float64

	// P-values vs control
	pDur    float64
	pTokens float64
	pLint   float64
	pSucc   float64
	pTime   float64 // p-value for timeouts
}

func (r *EnhancedMarkdownGenerator) writePerformanceComparison() {
	fmt.Fprintln(r.w, "\n## 🏆 Performance Comparison")

	avgs := r.stats

	// Identify Control (first in list)
	var control *avgMetrics
	if len(avgs) > 0 {
		control = &avgs[0]
	}

	// Calculate bests for trophies
	var bestDur time.Duration = 1000 * time.Hour
	var bestTokens int64 = 1000000000
	var bestSuccess float64 = -1.0

	for _, a := range avgs {
		if a.completedRuns > 0 && a.avgDur < bestDur {
			bestDur = a.avgDur
		}
		if a.completedRuns > 0 && a.avgTokens < bestTokens {
			bestTokens = a.avgTokens
		}
		if a.successRate > bestSuccess {
			bestSuccess = a.successRate
		}
	}

	fmt.Fprintln(r.w, "| Alternative | Avg Duration (Success) | Avg Tokens (Success) | Avg Lint (Success) | Success Rate | Test Pass Rate | Tool Accuracy | Status |")
	fmt.Fprintln(r.w, "| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |")

	for i, a := range avgs {
		durStr := "-"
		tokenStr := "-"
		lintStr := "-"
		successStr := fmt.Sprintf("%.0f%%", a.successRate)

		if a.completedRuns > 0 {
			durStr = a.avgDur.Round(time.Second).String()
			tokenStr = formatNumber(int(a.avgTokens))
			lintStr = fmt.Sprintf("%.1f", a.avgLint)

			// Calculate Deltas vs Control
			if i > 0 && control != nil && control.completedRuns > 0 {
				durDiff := a.avgDur - control.avgDur
				durPct := float64(durDiff) / float64(control.avgDur) * 100
				durStr += fmt.Sprintf(" (%+.0f%%)%s", durPct, sigFlag(a.pDur))

				tokenDiff := a.avgTokens - control.avgTokens
				tokenPct := float64(tokenDiff) / float64(control.avgTokens) * 100
				tokenStr += fmt.Sprintf(" (%+.0f%%)%s", tokenPct, sigFlag(a.pTokens))

				lintStr += sigFlag(a.pLint)
				successStr += sigFlag(a.pSucc)
			} else if i == 0 {
				durStr += " (Control)"
			}

			// Add Trophies
			if a.avgDur == bestDur {
				durStr += " 🏆"
			}
			if a.avgTokens == bestTokens {
				tokenStr += " 🏆"
			}
		}

		status := "✅"
		if a.successRate < 80 {
			status = "⚠️"
		}
		if a.successRate < 50 {
			status = "❌"
		}

		timeStr := ""
		if a.timeouts > 0 {
			timeStr = fmt.Sprintf(" (%.0f%% timeouts)%s", float64(a.timeouts)/float64(a.totalRuns)*100, sigFlag(a.pTime))
			status += " ⏱️"
		}

		fmt.Fprintf(r.w, "| **%s** | %s | %s | %s | %s%s | %.0f%% | %.1f%% | %s |\n",
			a.name,
			durStr,
			tokenStr,
			lintStr,
			successStr,
			timeStr,
			a.testPass,
			a.toolSuccess,
			status,
		)
	}

	r.writeEfficiencyHighlights(avgs)
	r.writeStatisticalNotes(avgs)
}

func sigFlag(p float64) string {
	if p < 0.01 {
		return " **"
	}
	if p < 0.05 {
		return " *"
	}
	return ""
}

func (r *EnhancedMarkdownGenerator) writeStatisticalNotes(avgs []avgMetrics) {
	fmt.Fprintln(r.w, "\n### 📊 Statistical Analysis Notes")

	// Check for small sample size
	smallSamples := false
	for _, a := range avgs {
		if a.totalRuns < 5 {
			smallSamples = true
			break
		}
	}

	if smallSamples {
		fmt.Fprintln(r.w, "> [!WARNING]")
		fmt.Fprintln(r.w, "> **Small Sample Size Detected:** Some alternatives have fewer than 5 repetitions. Statistical significance flags might be unreliable. Consider running at least 10 repetitions for more robust results.")
	}

	fmt.Fprintln(r.w, "- `*` Significant ($p < 0.05$): 95% confidence that the difference is not due to chance.")
	fmt.Fprintln(r.w, "- `**` Highly Significant ($p < 0.01$): 99% confidence.")
	fmt.Fprintln(r.w, "- Metrics compared using Welch's t-test (Duration, Tokens, Lint) and Fisher's Exact Test (Success Rate).")
}

func (r *EnhancedMarkdownGenerator) writeEfficiencyHighlights(avgs []avgMetrics) {
	fmt.Fprintln(r.w, "\n### 💡 Efficiency Highlights")
	// Only consider alternatives with successful runs
	var validAvgs []avgMetrics
	for _, a := range avgs {
		if a.completedRuns > 0 {
			validAvgs = append(validAvgs, a)
		}
	}

	if len(validAvgs) == 0 {
		fmt.Fprintln(r.w, "No successful runs to analyze efficiency.")
		return
	}

	type comparison struct {
		alt   string
		value float64
	}

	var toolEfficiency []comparison
	var tokenEfficiency []comparison
	var toolAccuracy []comparison

	for _, a := range validAvgs {
		toolEfficiency = append(toolEfficiency, comparison{a.name, float64(a.totalToolCalls) / float64(a.completedRuns)})
		tokenEfficiency = append(tokenEfficiency, comparison{a.name, float64(a.avgTokens)})
		toolAccuracy = append(toolAccuracy, comparison{a.name, a.toolSuccess})
	}

	sort.Slice(toolEfficiency, func(i, j int) bool {
		return toolEfficiency[i].value < toolEfficiency[j].value
	})
	sort.Slice(tokenEfficiency, func(i, j int) bool {
		return tokenEfficiency[i].value < tokenEfficiency[j].value
	})
	sort.Slice(toolAccuracy, func(i, j int) bool {
		return toolAccuracy[i].value > toolAccuracy[j].value
	})

	fmt.Fprintf(r.w, "- **Most Tool-Efficient**: %s (%.1f calls/success run) vs %s (%.1f calls/success run)\n",
		toolEfficiency[0].alt, toolEfficiency[0].value,
		toolEfficiency[len(toolEfficiency)-1].alt, toolEfficiency[len(toolEfficiency)-1].value)

	fmt.Fprintf(r.w, "- **Most Token-Efficient**: %s (%s/run) vs %s (%s/run)\n",
		tokenEfficiency[0].alt, formatNumber(int(tokenEfficiency[0].value)),
		tokenEfficiency[len(tokenEfficiency)-1].alt, formatNumber(int(tokenEfficiency[len(tokenEfficiency)-1].value)))

	fmt.Fprintf(r.w, "- **Best Tool Accuracy**: %s (%.1f%% success) vs %s (%.1f%% success)\n",
		toolAccuracy[0].alt, toolAccuracy[0].value,
		toolAccuracy[len(toolAccuracy)-1].alt, toolAccuracy[len(toolAccuracy)-1].value)
}

func (r *EnhancedMarkdownGenerator) writeSummaryByScenario() {
	fmt.Fprintln(r.w, "\n## 📈 Summary by Scenario & Alternative")
	// Calculate bests per scenario
	type scenBests struct {
		bestDur    time.Duration
		bestTokens int
	}
	bests := make(map[string]scenBests)

	type summaryRow struct {
		scenario  string
		alt       string
		count     int
		avgDur    time.Duration
		avgTokens int
		avgLint   float64
		passRate  float64
	}
	var summaries []summaryRow

	type groupStats struct {
		scenario     string
		alternative  string
		count        int
		successCount int
		sumDuration  time.Duration
		sumTokens    int
		sumLint      int
		passCount    int
	}
	statsByGroup := make(map[string]*groupStats)
	for _, res := range r.Results {
		key := fmt.Sprintf("%s|%s", res.Scenario, res.Alternative)
		s, ok := statsByGroup[key]
		if !ok {
			s = &groupStats{
				scenario:    res.Scenario,
				alternative: res.Alternative,
			}
			statsByGroup[key] = s
		}
		s.count++ // Total attempts

		if res.Error == nil {
			s.successCount++
			s.sumDuration += res.Duration
			if res.AgentMetrics != nil {
				s.sumTokens += res.AgentMetrics.TotalTokens
			}
			if res.EvaluationMetrics != nil {
				s.sumLint += res.EvaluationMetrics.LintIssues
			}
		}

		if res.EvaluationMetrics != nil && res.EvaluationMetrics.TestsPassed > 0 {
			s.passCount++
		}
	}

	// Populate summaries slice and calculate bests
	for _, s := range statsByGroup {
		avgDuration := time.Duration(0)
		avgTokens := 0
		avgLint := 0.0
		// Pass rate is based on TOTAL attempts
		passRate := 0.0
		if s.count > 0 {
			passRate = float64(s.passCount) / float64(s.count) * 100
		}

		// Averages based on SUCCESS count
		if s.successCount > 0 {
			avgDuration = s.sumDuration / time.Duration(s.successCount)
			avgTokens = s.sumTokens / s.successCount
			avgLint = float64(s.sumLint) / float64(s.successCount)
		}

		summaries = append(summaries, summaryRow{
			scenario:  s.scenario,
			alt:       s.alternative,
			count:     s.count,
			avgDur:    avgDuration,
			avgTokens: avgTokens,
			avgLint:   avgLint,
			passRate:  passRate,
		})
	}

	// Sort summaries
	sort.Slice(summaries, func(i, j int) bool {
		if summaries[i].scenario != summaries[j].scenario {
			return summaries[i].scenario < summaries[j].scenario
		}
		return summaries[i].alt < summaries[j].alt
	})

	for _, s := range summaries {
		// Only calculate bests if there was at least one success
		if s.avgDur == 0 {
			continue
		}

		b, ok := bests[s.scenario]
		if !ok {
			b = scenBests{1000 * time.Hour, 1000000000}
		}
		if s.avgDur < b.bestDur {
			b.bestDur = s.avgDur
		}
		if s.avgTokens < b.bestTokens {
			b.bestTokens = s.avgTokens
		}
		bests[s.scenario] = b
	}

	fmt.Fprintln(r.w, "| Scenario | Alternative | Runs | Avg Duration (Success) | Avg Tokens (Success) | Lint | Pass Rate |")
	fmt.Fprintln(r.w, "|---|---|---|---|---|---|---|")

	for _, s := range summaries {
		b := bests[s.scenario]

		durStr := "-"
		tokenStr := "-"

		if s.avgDur > 0 {
			durStr = s.avgDur.Round(time.Millisecond).String()
			if s.avgDur == b.bestDur {
				durStr = fmt.Sprintf("**%s**", durStr)
			}

			tokenStr = formatNumber(s.avgTokens)
			if s.avgTokens == b.bestTokens {
				tokenStr = fmt.Sprintf("**%s**", tokenStr)
			}
		}

		fmt.Fprintf(r.w, "| %s | %s | %d | %s | %s | %.1f | %.0f%% |\n",
			s.scenario, s.alt, s.count,
			durStr,
			tokenStr,
			s.avgLint, s.passRate)
	}
}

func (r *EnhancedMarkdownGenerator) writeDetailedMetrics() {
	fmt.Fprintln(r.w, "\n## 📋 Detailed Metrics")
	fmt.Fprintln(r.w, "\n<details>")
	fmt.Fprintln(r.w, "<summary>Click to expand detailed run-by-run metrics</summary>")
	fmt.Fprintln(r.w, "\n| Scenario | Alternative | Rep | Duration | Tokens | Tools | Test | Lint | Status |")
	fmt.Fprintln(r.w, "|---|---|---|---|---|---|---|---|---|")

	for _, res := range r.Results {
		status := "✅"
		if res.Error != nil {
			if strings.Contains(res.Error.Error(), "timeout") {
				status = "⏱️ Timeout"
			} else {
				status = "❌ Error"
			}
		} else if res.AgentMetrics != nil && res.AgentMetrics.LoopDetected {
			status = "🔄 Loop"
		} else if res.EvaluationMetrics != nil {
			if res.EvaluationMetrics.TestsFailed > 0 {
				status = "❌ Fail"
			} else if res.EvaluationMetrics.TestsPassed == 0 {
				status = "⚠️ Skip"
			}
		}

		tokens := 0
		toolInfo := "0"
		if res.AgentMetrics != nil {
			tokens = res.AgentMetrics.TotalTokens
			toolInfo = fmt.Sprintf("%d", len(res.AgentMetrics.ToolCalls))
		}

		testStatus := "-"
		lintCount := 0
		if res.EvaluationMetrics != nil {
			if res.EvaluationMetrics.TestsPassed > 0 {
				testStatus = "PASS"
			} else if res.EvaluationMetrics.TestsFailed > 0 {
				testStatus = "FAIL"
			} else {
				testStatus = "SKIP"
			}
			lintCount = res.EvaluationMetrics.LintIssues
		}

		fmt.Fprintf(r.w, "| %s | %s | %d | %s | %s | %s | %s | %d | %s |\n",
			res.Scenario, res.Alternative, res.Repetition,
			res.Duration.Round(time.Millisecond),
			formatNumber(tokens), toolInfo,
			testStatus, lintCount, status)
	}

	fmt.Fprintln(r.w, "\n</details>")
}

func (r *EnhancedMarkdownGenerator) writeEvaluationDetails() {
	// Only show if there are failures or lint issues
	hasIssues := false
	for _, res := range r.Results {
		if res.EvaluationMetrics != nil {
			if len(res.EvaluationMetrics.Tests) > 0 || len(res.EvaluationMetrics.Lints) > 0 {
				hasIssues = true
				break
			}
		}
	}

	if !hasIssues {
		return
	}

	fmt.Fprintln(r.w, "\n## 🧪 Detailed Evaluation (Tests & Lint)")

	for _, res := range r.Results {
		if res.EvaluationMetrics == nil {
			continue
		}

		if res.EvaluationMetrics.TestsFailed == 0 && res.EvaluationMetrics.LintIssues == 0 {
			continue
		}

		fmt.Fprintf(r.w, "\n### %s | %s | Rep %d\n", res.Alternative, res.Scenario, res.Repetition)

		if res.EvaluationMetrics.TestsFailed > 0 {
			fmt.Fprintln(r.w, "\n#### ❌ Test Failures")
			for _, t := range res.EvaluationMetrics.Tests {
				if t.Status == "FAIL" {
					fmt.Fprintf(r.w, "- **%s**\n", t.Name)
					if t.Output != "" {
						fmt.Fprintln(r.w, "  ```")
						// Indent output for better formatting
						lines := strings.Split(strings.TrimSpace(t.Output), "\n")
						for _, l := range lines {
							fmt.Fprintf(r.w, "    %s\n", l)
						}
						fmt.Fprintln(r.w, "  ```")
					}
				}
			}
		}

		if res.EvaluationMetrics.LintIssues > 0 {
			fmt.Fprintln(r.w, "\n#### ⚠️ Lint Findings")
			fmt.Fprintln(r.w, "| File | Line | Rule | Message |")
			fmt.Fprintln(r.w, "|---|---|---|---|")
			for _, l := range res.EvaluationMetrics.Lints {
				fmt.Fprintf(r.w, "| %s | %d | `%s` | %s |\n", l.File, l.Line, l.RuleID, l.Message)
			}
		}
	}
}

func (r *EnhancedMarkdownGenerator) writeToolUsage() {
	fmt.Fprintln(r.w, "\n## 🔧 Tool Usage Analysis")

	// 1. Overall summary by alternative
	fmt.Fprintln(r.w, "\n### Overall Summary")

	// Calculate totals by alternative
	type toolStats struct {
		totalCalls  int
		failedCalls int
		uniqueTools map[string]bool
		runs        int
		toolCounts  map[string]int
	}
	stats := make(map[string]*toolStats)

	for _, res := range r.Results {
		if res.AgentMetrics == nil {
			continue
		}
		if _, ok := stats[res.Alternative]; !ok {
			stats[res.Alternative] = &toolStats{uniqueTools: make(map[string]bool), toolCounts: make(map[string]int)}
		}
		s := stats[res.Alternative]
		s.runs++
		for _, tc := range res.AgentMetrics.ToolCalls {
			s.totalCalls++
			if tc.Status != "success" {
				s.failedCalls++
			}
			s.uniqueTools[tc.Name] = true
			s.toolCounts[tc.Name]++
		}
	}

	fmt.Fprintln(r.w, "\n| Alternative | Total Calls | Failed Calls | Unique Tools | Avg Calls/Run | Most Used Tool |")
	fmt.Fprintln(r.w, "|---|---|---|---|---|---|")

	for alt, s := range stats {
		avg := 0.0
		if s.runs > 0 {
			avg = float64(s.totalCalls) / float64(s.runs)
		}

		mostUsed := "-"
		maxCount := 0
		for t, c := range s.toolCounts {
			if c > maxCount {
				maxCount = c
				mostUsed = t
			}
		}
		if mostUsed != "-" {
			mostUsed = fmt.Sprintf("`%s` (%d)", mostUsed, maxCount)
		}

		failedStr := fmt.Sprintf("%d", s.failedCalls)
		if s.failedCalls > 0 {
			failRate := float64(s.failedCalls) / float64(s.totalCalls) * 100
			if failRate > 10 { // Threshold for bolding
				failedStr = fmt.Sprintf("**%d (%.0f%%)**", s.failedCalls, failRate)
			}
		}

		fmt.Fprintf(r.w, "| %s | %d | %s | %d | %.1f | %s |\n",
			alt, s.totalCalls, failedStr, len(s.uniqueTools), avg, mostUsed)
	}

	// 2. Unique tools per alternative
	fmt.Fprintln(r.w, "\n### Tool Availability")

	// Collect all tools used by each alternative
	allTools := make(map[string]map[string]bool) // alt -> set of tools
	for alt, s := range stats {
		allTools[alt] = s.uniqueTools
	}

	// Find unique tools (tools used by one alternative but not others)
	for alt, tools := range allTools {
		uniqueTools := []string{}
		for tool := range tools {
			isUnique := true
			for otherAlt, otherTools := range allTools {
				if otherAlt != alt && otherTools[tool] {
					isUnique = false
					break
				}
			}
			if isUnique {
				uniqueTools = append(uniqueTools, tool)
			}
		}

		if len(uniqueTools) > 0 {
			sort.Strings(uniqueTools)
			fmt.Fprintf(r.w, "- **%s** unique tools: ", alt)
			for i, tool := range uniqueTools {
				if i > 0 {
					fmt.Fprint(r.w, ", ")
				}
				fmt.Fprintf(r.w, "`%s`", tool)
			}
			fmt.Fprintln(r.w)
		}
	}

	// 3. Per-scenario breakdown
	fmt.Fprintln(r.w, "\n### Tool Usage by Scenario")

	type scenarioAltTools struct {
		scenario string
		alt      string
		tools    map[string]int
		failed   map[string]int
		runs     int
	}

	scenarioStats := make(map[string]*scenarioAltTools) // key: "scenario|alt"
	for _, res := range r.Results {
		key := fmt.Sprintf("%s|%s", res.Scenario, res.Alternative)
		if _, ok := scenarioStats[key]; !ok {
			scenarioStats[key] = &scenarioAltTools{
				scenario: res.Scenario,
				alt:      res.Alternative,
				tools:    make(map[string]int),
				failed:   make(map[string]int),
			}
		}
		stats := scenarioStats[key]
		stats.runs++
		if res.AgentMetrics != nil {
			for _, tc := range res.AgentMetrics.ToolCalls {
				stats.tools[tc.Name]++
				if tc.Status != "success" {
					stats.failed[tc.Name]++
				}
			}
		}
	}

	// Group by scenario
	scenarioGroups := make(map[string][]string) // scenario -> list of keys
	for key, stats := range scenarioStats {
		scenarioGroups[stats.scenario] = append(scenarioGroups[stats.scenario], key)
	}

	// Sort scenarios
	var scenarios []string
	for s := range scenarioGroups {
		scenarios = append(scenarios, s)
	}
	sort.Strings(scenarios)

	for _, s := range scenarios {
		fmt.Fprintf(r.w, "\n**%s**\n\n", s)
		fmt.Fprintln(r.w, "| Alternative | Tool | Avg Calls | Total | Failed |")
		fmt.Fprintln(r.w, "|---|---|---|---|---|")

		keys := scenarioGroups[s]
		sort.Strings(keys) // Sort by alternative name (embedded in key)

		for _, key := range keys {
			stats := scenarioStats[key]

			// Sort tools by usage (descending)
			type toolCount struct {
				name  string
				count int
			}
			var tcs []toolCount
			for t, c := range stats.tools {
				tcs = append(tcs, toolCount{t, c})
			}
			sort.Slice(tcs, func(i, j int) bool {
				return tcs[i].count > tcs[j].count
			})

			for _, tc := range tcs {
				avg := 0.0
				if stats.runs > 0 {
					avg = float64(tc.count) / float64(stats.runs)
				}
				failedCount := stats.failed[tc.name]
				failedStr := fmt.Sprintf("%d", failedCount)
				if failedCount > 0 {
					failedStr = fmt.Sprintf("**%d**", failedCount)
				}

				fmt.Fprintf(r.w, "| %s | `%s` | %.1f | %d | %s |\n",
					stats.alt, tc.name, avg, tc.count, failedStr)
			}
		}
	}
}

func (r *EnhancedMarkdownGenerator) writeConclusion() {
	fmt.Fprintln(r.w, "\n## 🎯 Conclusion")

	// Identify the winner
	var winner *avgMetrics
	var fastest *avgMetrics

	// Simple heuristic scoring: Success > Speed > Tokens
	// But let's just highlight the trade-offs

	for i := range r.stats {
		s := &r.stats[i]
		if winner == nil || s.successRate > winner.successRate {
			winner = s
		} else if s.successRate == winner.successRate {
			// Tie-break on speed
			if s.completedRuns > 0 && winner.completedRuns > 0 && s.avgDur < winner.avgDur {
				winner = s
			}
		}

		if s.completedRuns > 0 {
			if fastest == nil || s.avgDur < fastest.avgDur {
				fastest = s
			}
		}
	}

	if winner == nil {
		fmt.Fprintln(r.w, "No clear conclusion could be drawn as no alternatives successfully completed the tasks.")
		return
	}

	fmt.Fprintln(r.w, "\n### 🏅 Final Verdict")
	fmt.Fprintf(r.w, "Based on **%d experiment runs**, the recommended configuration is **%s**.\n\n", len(r.Results), winner.name)

	fmt.Fprintln(r.w, "**Reasoning:**")
	fmt.Fprintf(r.w, "- It achieved the highest reliability with a **%.1f%% success rate**.\n", winner.successRate)

	// Compare to control if winner is not control
	if r.stats[0].name != winner.name {
		control := r.stats[0]
		fmt.Fprintf(r.w, "- Compared to control (%s), it is ", control.name)
		if winner.successRate > control.successRate {
			fmt.Fprintf(r.w, "**more reliable** (+%.1f%%) ", winner.successRate-control.successRate)
		} else {
			fmt.Fprint(r.w, "equally reliable ")
		}

		if winner.avgDur < control.avgDur {
			fmt.Fprintf(r.w, "and **faster**.")
		} else {
			fmt.Fprintf(r.w, "but **slower**.")
		}
		fmt.Fprintln(r.w)
	}

	if fastest != nil && fastest.name != winner.name {
		fmt.Fprintf(r.w, "\n> **Note:** If raw speed is the only priority, **%s** was the fastest alternative (%.1fs/run), but at the cost of lower reliability (%.1f%% success).\n",
			fastest.name, fastest.avgDur.Seconds(), fastest.successRate)
	}

	if r.Cfg.Repetitions < 10 {
		fmt.Fprintln(r.w, "\n> [!TIP]")
		fmt.Fprintln(r.w, "> Run this experiment with more repetitions (e.g., 10+) to confirm these findings with higher statistical confidence.")
	}
}

func formatNumber(n interface{}) string {
	switch v := n.(type) {
	case int:
		if v > 1000000 {
			return fmt.Sprintf("%.1fM", float64(v)/1000000)
		} else if v > 1000 {
			return fmt.Sprintf("%.1fk", float64(v)/1000)
		}
		return fmt.Sprintf("%d", v)
	case int64:
		if v > 1000000 {
			return fmt.Sprintf("%.1fM", float64(v)/1000000)
		} else if v > 1000 {
			return fmt.Sprintf("%.1fk", float64(v)/1000)
		}
		return fmt.Sprintf("%d", v)
	}
	return fmt.Sprintf("%v", n)
}
