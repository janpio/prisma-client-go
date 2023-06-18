// Package generator acts as a prisma generator
package generator

import (
	"bytes"
	"fmt"
	"go/build"
	"go/format"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/steebchen/prisma-client-go/binaries"
	"github.com/steebchen/prisma-client-go/binaries/bindata"
	"github.com/steebchen/prisma-client-go/binaries/platform"
	"github.com/steebchen/prisma-client-go/logger"
)

const DefaultPackageName = "db"

func addDefaults(input *Root) {
	if input.Generator.Config.Package == "" {
		input.Generator.Config.Package = DefaultPackageName
	}
}

// Run invokes the generator, which builds the templates and writes to the specified output file.
func Run(input *Root) error {
	addDefaults(input)

	if input.Generator.Config.DisableGitignore != "true" && input.Generator.Config.DisableGoBinaries != "true" {
		logger.Debug.Printf("writing gitignore file")
		// generate a gitignore into the folder
		var gitignore = "# gitignore generated by Prisma Client Go. DO NOT EDIT.\n*_gen.go\n"
		if err := os.MkdirAll(input.Generator.Output.Value, os.ModePerm); err != nil {
			return fmt.Errorf("could not create output directory: %w", err)
		}
		if err := os.WriteFile(path.Join(input.Generator.Output.Value, ".gitignore"), []byte(gitignore), 0644); err != nil {
			return fmt.Errorf("could not write .gitignore: %w", err)
		}
	}

	if err := generateClient(input); err != nil {
		return fmt.Errorf("generate client: %w", err)
	}

	if err := generateBinaries(input); err != nil {
		return fmt.Errorf("generate binaries: %w", err)
	}

	return nil
}

func generateClient(input *Root) error {
	var buf bytes.Buffer

	ctx := build.Default
	pkg, err := ctx.Import("github.com/steebchen/prisma-client-go", ".", build.FindOnly)
	if err != nil {
		return fmt.Errorf("could not get main template asset: %w", err)
	}

	var templates []*template.Template

	templateDir := pkg.Dir + "/generator/templates"
	err = filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".gotpl") {
			tpl, err := template.ParseFiles(path)
			if err != nil {
				return err
			}
			templates = append(templates, tpl.Templates()...)
		}

		return err
	})

	if err != nil {
		return fmt.Errorf("could not walk dir %s: %w", templateDir, err)
	}

	// Run header template first
	header, err := template.ParseFiles(templateDir + "/_header.gotpl")
	if err != nil {
		return fmt.Errorf("could not find header template %s: %w", templateDir, err)
	}

	if err := header.Execute(&buf, input); err != nil {
		return fmt.Errorf("could not write header template: %w", err)
	}

	// Then process all remaining templates
	for _, tpl := range templates {
		if strings.Contains(tpl.Name(), "_") {
			continue
		}

		buf.Write([]byte(fmt.Sprintf("// --- template %s ---\n", tpl.Name())))

		if err := tpl.Execute(&buf, input); err != nil {
			return fmt.Errorf("could not write template file %s: %w", tpl.Name(), err)
		}

		if _, err := format.Source(buf.Bytes()); err != nil {
			return fmt.Errorf("could not format source %s from file %s: %w", buf.String(), tpl.Name(), err)
		}
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("could not format final source: %w", err)
	}

	output := input.Generator.Output.Value

	if strings.HasSuffix(output, ".go") {
		return fmt.Errorf("generator output should be a directory")
	}

	if err := os.MkdirAll(output, os.ModePerm); err != nil {
		return fmt.Errorf("could not run MkdirAll on path %s: %w", output, err)
	}

	// TODO make this configurable
	outFile := path.Join(output, "db_gen.go")
	if err := os.WriteFile(outFile, formatted, 0644); err != nil {
		return fmt.Errorf("could not write template data to file writer %s: %w", outFile, err)
	}

	return nil
}

func generateBinaries(input *Root) error {
	if input.Generator.Config.DisableGoBinaries == "true" {
		return nil
	}

	if input.GetEngineType() == "dataproxy" {
		logger.Debug.Printf("using data proxy; not fetching any engines")
		return nil
	}

	var targets []string

	for _, target := range input.Generator.BinaryTargets {
		targets = append(targets, target.Value)
	}

	targets = add(targets, "native")
	targets = add(targets, "linux")

	// TODO refactor
	for _, name := range targets {
		if name == "native" {
			name = platform.BinaryPlatformName()
		}

		// first, ensure they are actually downloaded
		if err := binaries.FetchEngine(binaries.GlobalCacheDir(), "query-engine", name); err != nil {
			return fmt.Errorf("failed fetching binaries: %w", err)
		}
	}

	if err := generateQueryEngineFiles(targets, input.Generator.Config.Package.String(), input.Generator.Output.Value); err != nil {
		return fmt.Errorf("could not write template data: %w", err)
	}

	return nil
}

func generateQueryEngineFiles(binaryTargets []string, pkg, outputDir string) error {
	for _, name := range binaryTargets {
		pt := name
		if strings.Contains(name, "debian") || strings.Contains(name, "rhel") {
			pt = "linux"
		}

		if name == "native" {
			name = platform.BinaryPlatformName()
			pt = runtime.GOOS
		}

		enginePath := binaries.GetEnginePath(binaries.GlobalCacheDir(), "query-engine", name)

		filename := fmt.Sprintf("query-engine-%s_gen.go", name)
		to := path.Join(outputDir, filename)

		// TODO check if already exists, but make sure version matches
		if err := bindata.WriteFile(strings.ReplaceAll(name, "-", "_"), pkg, pt, enginePath, to); err != nil {
			return fmt.Errorf("generate write go file: %w", err)
		}

		logger.Debug.Printf("write go file at %s", filename)
	}

	return nil
}

func add(list []string, item string) []string {
	keys := make(map[string]bool)
	if _, ok := keys[item]; !ok {
		keys[item] = true
		list = append(list, item)
	}
	return list
}
