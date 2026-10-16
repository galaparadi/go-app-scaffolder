package scaffold

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed all:templates/project
var templateFS embed.FS

type Config struct {
	ModulePath  string
	ProjectName string
	OutputDir   string
	Force       bool
}

type templateData struct {
	ModulePath  string
	ProjectName string
}

// Generate creates a project from the embedded template files.
func Generate(cfg Config) error {
	if strings.TrimSpace(cfg.ModulePath) == "" {
		return errors.New("module path tidak boleh kosong")
	}
	if strings.TrimSpace(cfg.ProjectName) == "" {
		return errors.New("nama proyek tidak boleh kosong")
	}
	if strings.TrimSpace(cfg.OutputDir) == "" {
		return errors.New("direktori output tidak boleh kosong")
	}

	out := filepath.Clean(cfg.OutputDir)
	info, err := os.Stat(out)
	switch {
	case err == nil:
		if !info.IsDir() {
			return fmt.Errorf("output %q sudah ada dan bukan direktori", out)
		}
		entries, readErr := os.ReadDir(out)
		if readErr != nil {
			return fmt.Errorf("membaca direktori output: %w", readErr)
		}
		if len(entries) > 0 && !cfg.Force {
			return fmt.Errorf("direktori output %q tidak kosong; pilih direktori lain atau gunakan --force", out)
		}
	case errors.Is(err, os.ErrNotExist):
		if err := os.MkdirAll(out, 0o755); err != nil {
			return fmt.Errorf("membuat direktori output: %w", err)
		}
	default:
		return fmt.Errorf("memeriksa direktori output: %w", err)
	}

	data := templateData{
		ModulePath:  strings.TrimSuffix(strings.TrimSpace(cfg.ModulePath), "/"),
		ProjectName: strings.TrimSpace(cfg.ProjectName),
	}

	err = fs.WalkDir(templateFS, "templates/project", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel("templates/project", path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		// .gitkeep files exist only to preserve empty directories in the template.
		// Do not copy them into the generated project.
		if !entry.IsDir() && entry.Name() == ".gitkeep" {
			return nil
		}

		target := filepath.Join(out, filepath.FromSlash(rel))
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}

		source, err := templateFS.ReadFile(path)
		if err != nil {
			return err
		}

		// Files ending in .tmpl are rendered and written without that suffix.
		if strings.HasSuffix(target, ".tmpl") {
			target = strings.TrimSuffix(target, ".tmpl")
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			tmpl, err := template.New(filepath.Base(target)).Option("missingkey=error").Parse(string(source))
			if err != nil {
				return fmt.Errorf("parse template %s: %w", path, err)
			}
			file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
			if err != nil {
				return err
			}
			execErr := tmpl.Execute(file, data)
			closeErr := file.Close()
			if execErr != nil {
				return fmt.Errorf("render template %s: %w", path, execErr)
			}
			if closeErr != nil {
				return closeErr
			}
			return nil
		}

		// Static files are copied verbatim.
		return os.WriteFile(target, source, 0o644)
	})
	if err != nil {
		return fmt.Errorf("generate project: %w", err)
	}
	return nil
}
