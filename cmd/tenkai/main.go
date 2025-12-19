package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/danicat/tenkai/internal/config"
	"github.com/danicat/tenkai/internal/db"
	"github.com/danicat/tenkai/internal/report"
	"github.com/danicat/tenkai/internal/runner"
	"github.com/danicat/tenkai/internal/workspace"
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file")
	repsFlag := flag.Int("reps", 0, "Override number of repetitions")
	concurrentFlag := flag.Int("concurrent", 0, "Override max concurrent workers")
	nameFlag := flag.String("name", "", "Override experiment name")
	altsFlag := flag.String("alternatives", "", "Comma-separated list of alternatives to run")
	scensFlag := flag.String("scenarios", "", "Comma-separated list of scenarios to run")
	experimentID := flag.Int64("experiment-id", 0, "Experiment ID to regenerate report or resume")
	fixReportPath := flag.String("fix-report", "", "Path to existing report.md to fix structure/formatting")
	flag.Parse()

	// Handle Report Fixer Mode
	if *fixReportPath != "" {
		if err := report.FixMarkdownReport(*fixReportPath); err != nil {
			log.Fatalf("Failed to fix report: %v", err)
		}
		fmt.Printf("Successfully fixed report structure at: %s\n", *fixReportPath)
		return
	}

	// Handle Report-Only Mode
	if *experimentID > 0 && flag.NArg() == 0 && *configPath == "" {
		handleReportOnly(*experimentID)
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	// 0. Initialize Database
	dbPath := filepath.Join(cwd, "experiments", "tenkai.db")
	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	if *configPath == "" {
		// Try to find a default config? No, explicit is better.
		// But we can suggest looking in studies/
		fmt.Println("Usage: tenkai --config <path_to_config_yaml> [options]")
		fmt.Printf("Tip: Check 'experiments/templates/' directory for examples\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 1. Load Configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Apply CLI overrides
	var overrideNotes []string
	if *repsFlag > 0 {
		if *repsFlag != cfg.Repetitions {
			overrideNotes = append(overrideNotes, fmt.Sprintf("⚠️ Repetitions overridden to %d via CLI (config: %d)", *repsFlag, cfg.Repetitions))
		}
		cfg.Repetitions = *repsFlag
	}
	if *concurrentFlag > 0 {
		if *concurrentFlag != cfg.MaxConcurrent {
			overrideNotes = append(overrideNotes, fmt.Sprintf("⚠️ Max concurrent workers overridden to %d via CLI (config: %d)", *concurrentFlag, cfg.MaxConcurrent))
		}
		cfg.MaxConcurrent = *concurrentFlag
	}
	if *nameFlag != "" {
		cfg.Name = *nameFlag
	}

	// Apply Filters
	if *altsFlag != "" {
		alts := strings.Split(*altsFlag, ",")
		var filtered []config.Alternative
		for _, name := range alts {
			name = strings.TrimSpace(name)
			for _, a := range cfg.Alternatives {
				if a.Name == name {
					filtered = append(filtered, a)
					break
				}
			}
		}
		if len(filtered) > 0 {
			// Only report if we actually filtered something out
			if len(filtered) < len(cfg.Alternatives) {
				overrideNotes = append(overrideNotes, fmt.Sprintf("🎯 Alternatives filtered to: %s", *altsFlag))
			}
			cfg.Alternatives = filtered
		}
	}

	if *scensFlag != "" {
		scens := strings.Split(*scensFlag, ",")
		var filtered []config.Scenario
		for _, name := range scens {
			name = strings.TrimSpace(name)
			for _, s := range cfg.Scenarios {
				if s.Name == name {
					filtered = append(filtered, s)
					break
				}
			}
		}
		if len(filtered) > 0 {
			// Only report if we actually filtered something out
			if len(filtered) < len(cfg.Scenarios) {
				overrideNotes = append(overrideNotes, fmt.Sprintf("🎯 Scenarios filtered to: %s", *scensFlag))
			}
			cfg.Scenarios = filtered
		}
	}

	// 2. Prepare Output Directory
	timestamp := time.Now()
	tsStr := timestamp.Format("20060102-150405")
	folderName := tsStr
	if cfg.Name != "" {
		folderName = fmt.Sprintf("%s_%s", tsStr, cfg.Name)
	}

	// Determine experiment directory: <cwd>/experiments/runs/<folderName>
	experimentDir := filepath.Join(cwd, "experiments", "runs", folderName)
	if err := os.MkdirAll(experimentDir, 0755); err != nil {
		log.Fatalf("Failed to create experiment dir: %v", err)
	}

	// NEW: Marshal Effective Config to YAML
	effectiveConfigData, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalf("Failed to marshal effective config: %v", err)
	}

	// Save Config to Output Dir
	destConfigPath := filepath.Join(experimentDir, "config.yaml")
	if err := os.WriteFile(destConfigPath, effectiveConfigData, 0644); err != nil {
		log.Printf("Warning: failed to write config file: %v", err)
	} else {
		fmt.Printf("Configuration saved to: %s\n", destConfigPath)
	}

	// 3. Initialize Dependencies
	scenariosDir := filepath.Join(cwd, "scenarios")
	// Legacy support removed, we only look in scenarios/ now

	// Search in scenariosDir
	wsMgr := workspace.New(cwd, scenariosDir, scenariosDir)

	r := runner.New(wsMgr, cfg.MaxConcurrent)
	r.SetDB(database) // Inject DB into runner

	// 4. Register Experiment in DB (Initial status: running)
	expRecord := &db.Experiment{
		Timestamp:     timestamp,
		Name:          cfg.Name,
		ConfigPath:    destConfigPath,
		Status:        "running",
		Reps:          cfg.Repetitions,
		Concurrent:    cfg.MaxConcurrent,
		Description:   cfg.Description,
		ConfigContent: string(effectiveConfigData),
	}

	expID, err := database.CreateExperiment(expRecord)
	if err != nil {
		log.Printf("Warning: failed to register experiment in DB: %v", err)
	}
	r.SetExperimentID(expID) // Set ID for result correlation

	// 5. Run Experiments
	fmt.Printf("Starting Tenkai experiment with %d repetition(s) and %d concurrent worker(s)...\n", cfg.Repetitions, cfg.MaxConcurrent)
	fmt.Printf("Output Directory: %s\n", experimentDir)

	results, err := r.Run(context.Background(), cfg, timestamp, experimentDir)
	if err != nil {
		database.UpdateExperimentStatus(expID, "failed")
		log.Fatalf("Experiment run failed: %v", err)
	}

	duration := time.Since(timestamp)
	fmt.Printf("\nExperiment completed in %s.\n", duration.Round(time.Millisecond))

	// Update DB with duration
	if err := database.UpdateExperimentDuration(expID, int64(duration)); err != nil {
		log.Printf("Warning: failed to update experiment duration: %v", err)
	}

	fmt.Println("All jobs processed. Analytics available in Dashboard.")

	// 7. Generate & Save Report (both console and markdown for DB)
	rep := report.New(results, os.Stdout, cfg, overrideNotes)
	println("\n--- Results ---")
	rep.GenerateConsoleReport()

	var buf bytes.Buffer
	if err := rep.GenerateMarkdown(&buf); err == nil {
		if err := database.UpdateExperimentReport(expID, buf.String()); err != nil {
			log.Printf("Warning: failed to save report to DB: %v", err)
		}
	}

	// 8. Finalize Status in DB (Runner already does this, but we ensure it's synced)
	finalStatus := "completed"

	if err := database.UpdateExperimentStatus(expID, finalStatus); err != nil {
		log.Printf("Warning: failed to finalize experiment status in DB: %v", err)
	}
}

func handleReportOnly(expID int64) {
	fmt.Printf("Regenerating report for Experiment ID: %d\n", expID)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	// 0. Initialize Database
	dbPath := filepath.Join(cwd, "experiments", "tenkai.db")
	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	// 1. Resolve Experiment
	exp, err := database.GetExperimentByID(expID)
	if err != nil {
		log.Fatalf("Failed to find experiment in DB for ID %d: %v", expID, err)
	}

	fmt.Printf("Found Experiment Name: %s\n", exp.Name)

	// 2. Fetch Results from DB
	dbResults, err := database.GetRunResults(exp.ID)
	if err != nil {
		log.Fatalf("Failed to fetch run results from DB: %v", err)
	}

	// 3. Convert DB Results to Runner Results
	var results []runner.Result
	for _, r := range dbResults {
		var res runner.Result
		if r.RawJSON != "" {
			if err := json.Unmarshal([]byte(r.RawJSON), &res); err != nil {
				log.Printf("Warning: failed to unmarshal RawJSON for run %d: %v", r.ID, err)
				// Fallback to manual population of basic fields
				res = runner.Result{
					Alternative: r.Alternative,
					Scenario:    r.Scenario,
					Repetition:  r.Repetition,
					Duration:    time.Duration(r.Duration),
					ErrorStr:    r.Error,
					EvaluationMetrics: &runner.EvaluationMetrics{
						TestsPassed: r.TestsPassed,
						TestsFailed: r.TestsFailed,
						LintIssues:  r.LintIssues,
					},
				}
			}
		} else {
			// Legacy record without RawJSON
			res = runner.Result{
				Alternative: r.Alternative,
				Scenario:    r.Scenario,
				Repetition:  r.Repetition,
				Duration:    time.Duration(r.Duration),
				ErrorStr:    r.Error,
				EvaluationMetrics: &runner.EvaluationMetrics{
					TestsPassed: r.TestsPassed,
					TestsFailed: r.TestsFailed,
					LintIssues:  r.LintIssues,
				},
			}
		}

		if res.ErrorStr != "" && res.Error == nil {
			res.Error = fmt.Errorf("%s", res.ErrorStr)
		}
		results = append(results, res)
	}

	// 4. Load Config (from DB)
	cfg := &config.Configuration{}
	if err := yaml.Unmarshal([]byte(exp.ConfigContent), cfg); err != nil {
		log.Printf("Warning: failed to unmarshal config from DB: %v", err)
	} else {
		fmt.Printf("Loaded config from Database\n")
	}

	// 5. Console Summary
	rep := report.New(results, os.Stdout, cfg, []string{"Report regenerated from SQLite Database"})
	println("\n--- Results ---")
	rep.GenerateConsoleReport()

	fmt.Printf("\nDatabase analytics hydrated for Experiment %d\n", exp.ID)
}
