package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/galaparadi/go-app-scaffolder/internal/scaffold"
)

func main() {
	module := flag.String("module", "", "Go module path, e.g. github.com/username/myapp (required)")
	name := flag.String("name", "", "project name (defaults to output directory name)")
	out := flag.String("out", "", "output directory (required)")
	force := flag.Bool("force", false, "allow writing into an existing output directory")
	flag.Parse()

	if *module == "" {
		fail("flag --module wajib diisi")
	}
	if *out == "" {
		fail("flag --out wajib diisi")
	}

	projectName := *name
	if projectName == "" {
		projectName = filepath.Base(filepath.Clean(*out))
	}
	if projectName == "." || projectName == string(filepath.Separator) || projectName == "" {
		fail("nama proyek tidak valid; gunakan --name")
	}

	if err := scaffold.Generate(scaffold.Config{
		ModulePath:  modulePath(*module),
		ProjectName: projectName,
		OutputDir:   *out,
		Force:       *force,
	}); err != nil {
		fail("%v", err)
	}

	fmt.Printf("Proyek %q berhasil dibuat di %s\n", projectName, filepath.Clean(*out))
}

func modulePath(value string) string {
	return value
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "scaffold: "+format+"\n", args...)
	os.Exit(1)
}
