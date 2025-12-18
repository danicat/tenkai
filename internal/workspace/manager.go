package workspace

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Manager handles the creation and setup of isolated workspaces for experiments.
type Manager struct {
	BasePath     string
	TemplatesDir string
}

// New creates a new Workspace Manager.
func New(basePath, templatesDir string) *Manager {
	return &Manager{
		BasePath:     basePath,
		TemplatesDir: templatesDir,
	}
}

// PrepareWorkspace creates an isolated workspace for a specific run.
// Structure: <experimentDir>/<alternative>/<scenario>/<repetition>
func (m *Manager) PrepareWorkspace(experimentDir, alternative, scenario string, repetition int, settingsPath, contextFilePath string, policyFiles []string) (string, error) {
	wsPath := filepath.Join(experimentDir, alternative, scenario, fmt.Sprintf("rep-%d", repetition))

	if err := os.MkdirAll(wsPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create workspace directory: %w", err)
	}

	// Copy template content
	tmplPath := filepath.Join(m.TemplatesDir, scenario)
	if err := copyDir(tmplPath, wsPath); err != nil {
		return "", fmt.Errorf("failed to copy template from %s to %s: %w", tmplPath, wsPath, err)
	}

	// Handle MCP settings if provided
	if settingsPath != "" {
		geminiDir := filepath.Join(wsPath, ".gemini")
		if err := os.MkdirAll(geminiDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create .gemini directory: %w", err)
		}

		destSettings := filepath.Join(geminiDir, "settings.json")
		if err := copyFile(settingsPath, destSettings); err != nil {
			return "", fmt.Errorf("failed to create settings file: %w", err)
		}
	}

	// Handle Context File (GEMINI.md) if provided
	if contextFilePath != "" {
		destContext := filepath.Join(wsPath, "GEMINI.md")
		if err := copyFile(contextFilePath, destContext); err != nil {
			return "", fmt.Errorf("failed to copy context file: %w", err)
		}
	}

	// Handle Policy Files if provided
	if len(policyFiles) > 0 {
		policiesDir := filepath.Join(wsPath, ".gemini", "policies")
		if err := os.MkdirAll(policiesDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create policies directory: %w", err)
		}

		for _, pf := range policyFiles {
			destPolicy := filepath.Join(policiesDir, filepath.Base(pf))
			if err := copyFile(pf, destPolicy); err != nil {
				return "", fmt.Errorf("failed to copy policy file %s: %w", pf, err)
			}
		}
	}

	return wsPath, nil
}

func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
