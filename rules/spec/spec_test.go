package spec

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"testing"

	"github.com/glennsarti/sentinel-lint/lint"
	"github.com/glennsarti/sentinel-lint/rules"
	"github.com/glennsarti/sentinel-lint/runner"

	"github.com/glennsarti/sentinel-parser/features"
	"github.com/glennsarti/sentinel-parser/position"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/txtar"
)

var updateTestFiles = flag.Bool("update-test-files", false, "Update the test files")

func TestRulesSpecs(t *testing.T) {
	fixturesDir := path.Join("test-fixtures")

	items, err := os.ReadDir(fixturesDir)
	if err != nil {
		t.Error(err)
		return
	}
	for _, item := range items {
		if item.IsDir() {
			if valid, actualVersion := features.ValidateSentinelVersion(item.Name()); valid {
				t.Run(item.Name(), func(t *testing.T) {
					processTestFixturesDir(item.Name(), fixturesDir, actualVersion, t)
				})
			} else {
				t.Fatalf("Invalid directory name for a sentinel version: '%s'", item.Name())
			}
		}
	}
}

func processTestFixturesDir(relPath, srcDir, sentinelVersion string, t *testing.T) {
	dirPath := path.Join(srcDir, relPath)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			t.Run(entry.Name(), func(t *testing.T) {
				processTestFixturesDir(path.Join(relPath, entry.Name()), srcDir, sentinelVersion, t)
			})
		} else {
			if strings.HasSuffix(entry.Name(), ".sentinel") || strings.HasSuffix(entry.Name(), ".hcl") {
				t.Run(entry.Name(), func(t *testing.T) {
					if err := createSpecFileFromSource(entry.Name(), dirPath, sentinelVersion, t); err != nil {
						t.Error(err)
					}
				})
			}
			if strings.HasSuffix(entry.Name(), ".txtar") {
				t.Run(entry.Name(), func(t *testing.T) {
					if *updateTestFiles {
						if err := updateSpecFile(entry.Name(), dirPath, sentinelVersion, t); err != nil {
							t.Error(err)
						}
					} else {
						if err := testSpecFile(entry.Name(), dirPath, sentinelVersion, t); err != nil {
							t.Error(err)
						}
					}
				})
			}
		}
	}
}

func createSpecFileFromSource(filename, parentPath, sentinelVersion string, t *testing.T) error {
	srcFile := path.Join(parentPath, filename)
	dstFilename := filename
	archiveFilename := filename

	if strings.HasSuffix(dstFilename, ".sentinel") {
		dstFilename = strings.ReplaceAll(dstFilename, ".sentinel", ".txtar")
		archiveFilename = policyFilename
	}
	if strings.HasSuffix(dstFilename, ".hcl") {
		dstFilename = strings.ReplaceAll(dstFilename, ".hcl", ".txtar")
		archiveFilename = primaryConfigFilename
	}

	dstFile := path.Join(parentPath, dstFilename)

	f, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	contents, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	// Must end in LF
	if contents[len(contents)-1] != 10 {
		contents = append(contents, 10)
	}

	arc := &txtar.Archive{}
	arc.Files = append(arc.Files, txtar.File{
		Name: archiveFilename,
		Data: contents,
	})

	// Add default files
	arc.Files = append(arc.Files, txtar.File{
		Name: archiveConfigFile,
		Data: make([]byte, 0),
	})

	if f, err = os.Create(dstFile); err != nil {
		return err
	} else {
		if _, err = f.Write(txtar.Format(arc)); err != nil {
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	}

	if err := updateSpecFile(dstFilename, parentPath, sentinelVersion, t); err != nil {
		return err
	}

	return nil
}

func updateSpecFile(filename, parentPath, sentinelVersion string, t *testing.T) error {
	filePath := path.Join(parentPath, filename)

	arc, err := parseTxtarArchive(filePath)
	if err != nil {
		return err
	}

	actualIssues, err := generateIssues(arc, sentinelVersion)
	if err != nil {
		return err
	}

	arc.UpdateFile(archiveIssueOutputFile, issuesToBytes(actualIssues))

	if err := arc.Write(filePath); err != nil {
		return err
	}

	t.Logf("Updated %s", filePath)
	return err
}

func testSpecFile(filename, parentPath, sentinelVersion string, t *testing.T) error {
	filePath := path.Join(parentPath, filename)

	arc, err := parseTxtarArchive(filePath)
	if err != nil {
		return err
	}

	actualIssues, err := generateIssues(arc, sentinelVersion)
	if err != nil {
		return err
	}

	t.Run("issues", func(t *testing.T) {
		expectedString := string(arc.IssuesFile.Data)
		actualString := string(issuesToBytes(actualIssues))
		if diff := cmp.Diff(expectedString, actualString); diff != "" {
			t.Fatal(diff)
		}
	})

	return nil
}

func generateIssues(pa *parsedArchive, sentinelVersion string) (lint.Issues, error) {
	cfg := parseLintConfig(pa, sentinelVersion)
	ruleset := rules.NewDefaultRuleSet()

	file, issues, err := extractSourceFile(pa.SourceFile, sentinelVersion)
	if err != nil {
		return nil, err
	}

	if len(issues) > 0 {
		return issues, nil
	}

	r, err := runner.NewRunner(cfg, ruleset, file)
	if err != nil {
		return nil, err
	}

	return r.Run()
}

func parseLintConfig(_ *parsedArchive, sentinelVersion string) lint.Config {
	return lint.Config{
		SentinelVersion: sentinelVersion,
		FailFast:        false,
	}
}

func issuesToBytes(issues lint.Issues) []byte {
	var out bytes.Buffer
	if issues == nil {
		out.WriteString("nil\n")
		return out.Bytes()
	}

	diagStrs := make([]string, 0)

	for _, issue := range issues {
		if issue == nil {
			continue
		}
		errMsg := fmt.Sprintf("(%s) [%s] %s: %s", rangeToString(issue.Range), issue.Severity.String(), issue.RuleId, issue.Summary)
		diagStrs = append(diagStrs, errMsg)

		if issue.Related != nil {
			for _, relIssue := range *issue.Related {

				relMsg := fmt.Sprintf("(%s) [Related] (%s): %s", rangeToString(issue.Range), rangeToString(relIssue.Range), relIssue.Summary)

				diagStrs = append(diagStrs, relMsg)
			}
		}
	}
	sort.Strings(diagStrs)

	for _, val := range diagStrs {
		// Must always terminate in LF
		out.WriteString(val + "\n")
	}

	return out.Bytes()
}

func rangeToString(r *position.SourceRange) string {
	if r == nil {
		return "NIL Range"
	}
	return fmt.Sprintf("%s:%d:%d-%d:%d",
		r.Filename,
		r.Start.Line,
		r.Start.Column,
		r.End.Line,
		r.End.Column,
	)
}
