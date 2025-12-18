package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/danicat/tenkai/internal/config"
	"github.com/danicat/tenkai/internal/history"
	"github.com/danicat/tenkai/internal/report"
	"github.com/danicat/tenkai/internal/runner"
	"github.com/danicat/tenkai/internal/workspace"
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file")
	repsFlag := flag.Int("reps", 0, "Override number of repetitions")
	concurrentFlag := flag.Int("concurrent", 0, "Override max concurrent workers")
	reportOnlyPath := flag.String("report-only", "", "Path to experiment directory to regenerate report for")
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
	if *reportOnlyPath != "" {
		handleReportOnly(*reportOnlyPath)
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	if *configPath == "" {
		// Try to find a default config? No, explicit is better.
		// But we can suggest looking in configs/
		fmt.Println("Usage: tenkai --config <path_to_config_yaml> [options]")
		fmt.Printf("Tip: Check 'configs/' directory for examples, e.g., --config configs/run-simple.yaml\n")
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
		overrideNotes = append(overrideNotes, fmt.Sprintf("⚠️ Repetitions overridden to %d via CLI (config: %d)", *repsFlag, cfg.Repetitions))
		cfg.Repetitions = *repsFlag
	}
	if *concurrentFlag > 0 {
		overrideNotes = append(overrideNotes, fmt.Sprintf("⚠️ Max concurrent workers overridden to %d via CLI (config: %d)", *concurrentFlag, cfg.MaxConcurrent))
		cfg.MaxConcurrent = *concurrentFlag
	}

	// 2. Prepare Output Directory
	timestamp := time.Now()
	tsStr := timestamp.Format("20060102-150405")
	folderName := tsStr
	if cfg.Name != "" {
		folderName = fmt.Sprintf("%s_%s", tsStr, cfg.Name)
	}

	// NEW STRUCTURE: experiments/runs/<timestamp>_...
	experimentDir := filepath.Join(cwd, "experiments", "runs", folderName)
	if err := os.MkdirAll(experimentDir, 0755); err != nil {
		log.Fatalf("Failed to create experiment dir: %v", err)
	}

	// NEW: Copy Config to Output Dir
	destConfigPath := filepath.Join(experimentDir, "config.yaml")
	if err := copyFile(*configPath, destConfigPath); err != nil {
		log.Printf("Warning: failed to copy config file: %v", err)
	} else {
		// If we had overrides, maybe valid to append them to the copied yaml as comments?
		// For now simple copy is fine.
		fmt.Printf("Configuration saved to: %s\n", destConfigPath)
	}

	// 3. Initialize Dependencies
	templatesDir := filepath.Join(cwd, "templates")
	wsMgr := workspace.New(cwd, templatesDir)

	r := runner.New(wsMgr, cfg.MaxConcurrent)

	// 4. Run Experiments
	fmt.Printf("Starting Tenkai experiment with %d repetition(s) and %d concurrent worker(s)...\n", cfg.Repetitions, cfg.MaxConcurrent)
	fmt.Printf("Output Directory: %s\n", experimentDir)

	results, err := r.Run(context.Background(), cfg, timestamp, experimentDir) // Pass experimentDir for logging
	if err != nil {
		log.Fatalf("Experiment run failed: %v", err)
	}

	duration := time.Since(timestamp)
	fmt.Printf("\nExperiment completed in %s.\n", duration.Round(time.Millisecond))

	// 5. Save Results to JSON (before report)
	resultsPath := filepath.Join(experimentDir, "results.json")
	if file, err := os.Create(resultsPath); err == nil {
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(results); err != nil {
			log.Printf("Warning: failed to write results.json: %v", err)
		}
		file.Close()
		fmt.Printf("\nRaw results saved to: %s\n", resultsPath)
	} else {
		log.Printf("Warning: failed to create results.json: %v", err)
	}

	fmt.Println("All jobs processed. Generating report...")

	// 6. Generate Report
	rep := report.New(results, os.Stdout, cfg, overrideNotes)

	println("\n--- Results ---")
	rep.GenerateConsoleReport()

	reportPath := filepath.Join(experimentDir, "report.md")
	if err := rep.GenerateMarkdownReport(reportPath); err != nil {
		log.Printf("Failed to generate markdown report: %v", err)
	} else {
		fmt.Printf("\nDetailed report generated at: %s\n", reportPath)
	}

	// 7. Update Index
	idx := history.NewIndexer(cwd)
	record := history.ExperimentRecord{
		Timestamp:   timestamp.Format(time.RFC3339),
		Name:        cfg.Name,
		ConfigPath:  destConfigPath, // Point to the reproducible copy
		ResultsPath: resultsPath,
		ReportPath:  reportPath,
		Status:      "completed", // TODO: Detect failure
	}
	if err := idx.AddEntry(record); err != nil {
		log.Printf("Warning: failed to update experiment index: %v", err)
	}
}

func handleReportOnly(experimentDir string) {
	fmt.Printf("Regenerating report for: %s\n", experimentDir)

	resultsPath := filepath.Join(experimentDir, "results.json")
	file, err := os.Open(resultsPath)
	if err != nil {
		log.Fatalf("Failed to open results.json: %v", err)
	}
	defer file.Close()

	var results []runner.Result
	if err := json.NewDecoder(file).Decode(&results); err != nil {
		log.Fatalf("Failed to decode results.json: %v", err)
	}

	// Rehydrate error interface from ErrorStr
	for i := range results {
		if results[i].ErrorStr != "" {
			results[i].Error = fmt.Errorf("%s", results[i].ErrorStr)
		}
	}

	// Try to load config if it exists (optional, mostly for display)
	// We'll use a dummy config if not found
	cfg := &config.Configuration{}

	// NEW: Check for config.yaml in the experiment dir first
	localConfigPath := filepath.Join(experimentDir, "config.yaml")
	if loadedCfg, err := config.Load(localConfigPath); err == nil {
		cfg = loadedCfg
		fmt.Printf("Loaded original config from %s\n", localConfigPath)
	}

	rep := report.New(results, os.Stdout, cfg, []string{"Report regenerated from results.json"})

	println("\n--- Results ---")
	rep.GenerateConsoleReport()

	reportPath := filepath.Join(experimentDir, "report_regenerated.md")
	if err := rep.GenerateMarkdownReport(reportPath); err != nil {
		log.Printf("Failed to generate markdown report: %v", err)
	} else {
		fmt.Printf("\nDetailed report generated at: %s\n", reportPath)
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
