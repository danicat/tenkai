package runner

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/danicat/tenkai/internal/config"
	"github.com/danicat/tenkai/internal/parser"
	"github.com/danicat/tenkai/internal/workspace"
)

// EvaluationMetrics holds the results of code verification.
type EvaluationMetrics struct {
	TestsPassed int
	TestsFailed int
	LintIssues  int
	LintDetails map[string]int // Category/Message -> Count
}

// Result represents the outcome of a single experiment run.
type Result struct {
	Alternative string        `json:"alternative"`
	Scenario    string        `json:"scenario"`
	Repetition  int           `json:"repetition"`
	Duration    time.Duration `json:"duration"`
	Workspace   string        `json:"workspace"`
	Stdout      string        `json:"stdout"`
	Stderr      string        `json:"stderr"`
	Error       error         `json:"-"` // Don't marshal interface
	ErrorStr    string        `json:"error,omitempty"`

	AgentMetrics      *parser.AgentMetrics `json:"agent_metrics,omitempty"`
	EvaluationMetrics *EvaluationMetrics   `json:"evaluation_metrics,omitempty"`
}

// Runner handles the execution of experiments.
type Runner struct {
	WorkspaceMgr  *workspace.Manager
	MaxConcurrent int
}

// New creates a new Runner.
func New(wsMgr *workspace.Manager, maxConcurrent int) *Runner {
	if maxConcurrent < 1 {
		maxConcurrent = 1
	}
	// Hard limit concurrency to GOMAXPROCS to avoid overloading system
	numCPU := runtime.GOMAXPROCS(0)
	if maxConcurrent > numCPU {
		log.Printf("Capping concurrency from %d to GOMAXPROCS (%d)", maxConcurrent, numCPU)
		maxConcurrent = numCPU
	}

	return &Runner{
		WorkspaceMgr:  wsMgr,
		MaxConcurrent: maxConcurrent,
	}
}

// Run executes the experiments defined in the configuration.
func (r *Runner) Run(ctx context.Context, cfg *config.Configuration, timestamp time.Time, experimentDir string) ([]Result, error) {
	var results []Result
	resultsLimit := cfg.Repetitions * len(cfg.Alternatives) * len(cfg.Scenarios)
	resultsChan := make(chan Result, resultsLimit)

	// Ensure logs directory exists
	logsDir := filepath.Join(experimentDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs dir: %w", err)
	}

	sem := make(chan struct{}, r.MaxConcurrent)
	var wg sync.WaitGroup

	// Parse timeout from config
	timeout := 5 * time.Minute
	if cfg.Timeout != "" {
		if d, err := time.ParseDuration(cfg.Timeout); err == nil {
			timeout = d
		} else {
			log.Printf("Warning: invalid timeout %q, using default 5m", cfg.Timeout)
		}
	}

	for _, alt := range cfg.Alternatives {
		for _, scen := range cfg.Scenarios {
			for i := 1; i <= cfg.Repetitions; i++ {
				wg.Add(1)
				go func(alt config.Alternative, scen config.Scenario, rep int) {
					defer wg.Done()
					sem <- struct{}{}
					defer func() { <-sem }()

					res := r.runSingle(ctx, timestamp, cfg.Name, alt, scen, rep, experimentDir, timeout)
					resultsChan <- res
				}(alt, scen, i)
			}
		}
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	completed := 0
	for res := range resultsChan {
		completed++
		if completed%10 == 0 || completed == resultsLimit {
			percentage := float64(completed) / float64(resultsLimit) * 100
			log.Printf("Progress: %d/%d jobs completed (%.1f%%)", completed, resultsLimit, percentage)
		}
		results = append(results, res)
	}

	return results, nil
}

func (r *Runner) runSingle(ctx context.Context, timestamp time.Time, experimentName string, alt config.Alternative, scen config.Scenario, rep int, experimentDir string, timeout time.Duration) Result {
	start := time.Now()

	res := Result{
		Alternative: alt.Name,
		Scenario:    scen.Name,
		Repetition:  rep,
	}

	// Calculate a simple ID or just use Name/Scen/Rep for logging
	jobID := fmt.Sprintf("[%s|%s|#%d]", alt.Name, scen.Name, rep)
	log.Printf("START %s", jobID)

	// Ensure paths are absolute before passing to WorkspaceMgr
	settingsPath := alt.SettingsPath
	if settingsPath != "" {
		if abs, err := filepath.Abs(settingsPath); err == nil {
			settingsPath = abs
		}
	}
	contextFilePath := alt.ContextFilePath
	if contextFilePath != "" {
		if abs, err := filepath.Abs(contextFilePath); err == nil {
			contextFilePath = abs
		}
	}

	wsPath, err := r.WorkspaceMgr.PrepareWorkspace(experimentDir, alt.Name, scen.Name, rep, settingsPath, contextFilePath, alt.PolicyFiles)
	if err != nil {
		res.Error = fmt.Errorf("workspace prep failed: %w", err)
		res.ErrorStr = res.Error.Error()
		return res
	}
	res.Workspace = wsPath

	// Anchor the workspace to prevent Gemini from walking up to tenkai root
	// If go.mod doesn't exist, create a dummy one
	goModPath := filepath.Join(wsPath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		initCmd := exec.CommandContext(ctx, "go", "mod", "init", "experiment")
		initCmd.Dir = wsPath
		if err := initCmd.Run(); err != nil {
			log.Printf("Warning: failed to init dummy go.mod: %v", err)
		}
	}

	// 1. Prepare Command with 5-minute timeout
	cmdName := alt.Command
	cmdArgs := alt.Args

	// Handle System Prompt File
	// Check if configured, resolve absolute path if needed
	envMap := make(map[string]string)
	for k, v := range alt.Env {
		envMap[k] = v
	}

	if alt.SystemPromptFile != "" {
		// Resolve relative to CWD (or ideally config file but we don't pass that here yet)
		// For now assuming CWD or absolute
		absPath, err := filepath.Abs(alt.SystemPromptFile)
		if err == nil {
			envMap["GEMINI_SYSTEM_MD"] = absPath
		} else {
			log.Printf("Warning: failed to resolve system prompt file %s: %v", alt.SystemPromptFile, err)
		}
	}

	// Create timeout context
	cmdCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, cmdName, cmdArgs...)
	cmd.Dir = wsPath

	cmd.Env = os.Environ()
	for k, v := range envMap {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	// 2. Pipe PROMPT.md to Stdin using Pipe (safer closing)
	promptPath := filepath.Join(wsPath, "PROMPT.md")
	promptFile, err := os.Open(promptPath)
	if err == nil {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			res.Error = fmt.Errorf("failed to create stdin pipe: %w", err)
			promptFile.Close()
			return res
		}

		go func() {
			defer stdin.Close()
			defer promptFile.Close()
			if _, err := io.Copy(stdin, promptFile); err != nil {
				fmt.Printf("Error piping stdin: %v\n", err)
			}
		}()
	} else {
		// If PROMPT.md is missing, this runs without prompt!
		// But it should be there from template copy.
		res.Error = fmt.Errorf("PROMPT.md missing: %w", err)
		res.ErrorStr = res.Error.Error()
		return res
	}

	// 3. Capture Output (Stream JSON expected)
	// Write directly to file to avoid memory buffering issues and allow tailing
	// NEW: Use logsDir for cleaner organization
	logsDir := filepath.Join(experimentDir, "logs") // Derive from experimentDir
	safeAlt := strings.ReplaceAll(alt.Name, " ", "_")
	safeScen := strings.ReplaceAll(scen.Name, " ", "_")
	logFilename := fmt.Sprintf("%s_%s_rep%d", safeAlt, safeScen, rep)

	logPath := filepath.Join(logsDir, logFilename+".jsonl")
	logFile, err := os.Create(logPath)
	if err != nil {
		res.Error = fmt.Errorf("failed to create events log: %w", err)
		res.ErrorStr = res.Error.Error()
		return res
	}
	defer logFile.Close()

	cmd.Stdout = logFile
	// Capture Stderr to file to avoid "zombie pipe" blocking
	stderrPath := filepath.Join(logsDir, logFilename+".stderr.log")
	stderrFile, err := os.Create(stderrPath)
	if err != nil {
		res.Error = fmt.Errorf("failed to create stderr log: %w", err)
		res.ErrorStr = res.Error.Error()
		return res
	}
	defer stderrFile.Close()

	cmd.Stderr = stderrFile

	if err := cmd.Run(); err != nil {
		// Check if timeout
		if cmdCtx.Err() == context.DeadlineExceeded {
			res.Error = fmt.Errorf("execution timeout (5m limit)")
		} else {
			res.Error = fmt.Errorf("execution failed: %w", err)
		}
		res.ErrorStr = res.Error.Error()
		// Don't return yet, we might have logs to parse
	}

	res.Duration = time.Since(start)
	res.Stdout = "" // stdout is in file

	// Read stderr content for result struct (limited size if needed, but full file is on disk)
	stderrBytes, _ := os.ReadFile(stderrPath)
	res.Stderr = string(stderrBytes)

	// 4. Parse Metrics
	metrics, err := parser.ParseEvents(logPath)

	if err == nil {
		res.AgentMetrics = metrics
	} else {
		// Maybe not stream-json output?
	}

	// 5. Verify Code (Test & Lint) - skip for timeouts since code is incomplete
	shouldEvaluate := res.Error == nil
	if res.Error != nil && !strings.Contains(res.Error.Error(), "timeout") {
		// Non-timeout errors might still have partial code worth evaluating
		shouldEvaluate = true
	}

	if shouldEvaluate {
		res.EvaluationMetrics = r.evaluateCode(ctx, wsPath)
	}

	log.Printf("DONE  %s in %s", jobID, res.Duration.Round(time.Millisecond))
	return res
}

func (r *Runner) evaluateCode(ctx context.Context, wsPath string) *EvaluationMetrics {
	eval := &EvaluationMetrics{
		LintDetails: make(map[string]int),
	}

	// Run Tests
	// go test -v ./...
	testCmd := exec.CommandContext(ctx, "go", "test", "-v", "./...")
	testCmd.Dir = wsPath
	var testOut bytes.Buffer
	testCmd.Stdout = &testOut
	testCmd.Stderr = &testOut
	_ = testCmd.Run() // verify error handled by parsing output

	// Poor man's output parsing for tests
	output := testOut.String()
	if strings.Contains(output, "FAIL") {
		// Rough count of failures?
		// For now simple pass/fail of the whole suite
		if err := testCmd.ProcessState.ExitCode(); err != 0 {
			eval.TestsFailed = 1
		} else {
			eval.TestsPassed = 1
		}
	} else if strings.Contains(output, "PASS") {
		eval.TestsPassed = 1
	}

	// Run Linter
	// golangci-lint run ./...
	// We assume golangci-lint is installed. If not, maybe use go vet?
	// Let's try golangci-lint, if it fails, fallback to go vet.
	lintCmd := exec.CommandContext(ctx, "golangci-lint", "run", "./...")
	lintCmd.Dir = wsPath
	var lintOut bytes.Buffer
	lintCmd.Stdout = &lintOut
	if err := lintCmd.Run(); err != nil {
		// Lint issues found (exit code != 0) OR command not found
		if strings.Contains(err.Error(), "executable file not found") {
			// Fallback to go vet
			lintCmd = exec.CommandContext(ctx, "go", "vet", "./...")
			lintCmd.Dir = wsPath
			lintOut.Reset()
			lintCmd.Stdout = &lintOut
			lintCmd.Stderr = &lintOut // vet writes to stderr
			_ = lintCmd.Run()
		}
	}

	// Parse Lint Output
	// Format is usually "file:line:col: message (category)"
	scanner := bufio.NewScanner(strings.NewReader(lintOut.String()))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		// Simple heuristic: count lines that look like errors
		if strings.Contains(line, ".go:") {
			eval.LintIssues++
			// Extract category if possible? golangci-lint output varies.
			// Just use the line as key for now to dedup?
			// Use simple "issue" category
			eval.LintDetails["issue"]++
		}
	}

	return eval
}
