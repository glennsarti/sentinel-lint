package spec

import (
	"errors"
	"io"
	"os"

	"github.com/glennsarti/sentinel-lint/lint"

	"github.com/glennsarti/sentinel-lint/parsing"

	"golang.org/x/tools/txtar"
)

const archiveConfigFile = "lint-config.txt"
const archiveIssueOutputFile = "issues.txt"
const archiveErrorOutputFile = "errors.txt" // TODO: Do we need this?

type parsedArchive struct {
	ConfigFile txtar.File
	ErrorsFile txtar.File
	IssuesFile txtar.File
	SourceFile map[string]txtar.File
	raw        *txtar.Archive
}

func parseTxtarArchive(filePath string) (*parsedArchive, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	f.Close()

	arc := &parsedArchive{
		SourceFile: make(map[string]txtar.File, 0),
	}
	arc.raw = txtar.Parse(contents)

	for _, f := range arc.raw.Files {
		switch f.Name {
		case archiveConfigFile:
			arc.ConfigFile = f
		case archiveIssueOutputFile:
			arc.IssuesFile = f
		case archiveErrorOutputFile:
			arc.ErrorsFile = f
		default:
			arc.SourceFile[f.Name] = f
		}
	}

	return arc, nil
}

func (pa *parsedArchive) UpdateFile(name string, content []byte) {
	fileIdx := -1
	for idx, f := range pa.raw.Files {
		if f.Name == name {
			fileIdx = idx
			break
		}
	}

	tf := txtar.File{
		Name: name,
		Data: content,
	}

	if fileIdx == -1 {
		pa.raw.Files = append(pa.raw.Files, tf)
	} else {
		pa.raw.Files[fileIdx] = tf
	}
}

func (pa *parsedArchive) Write(filepath string) error {
	if f, err := os.Create(filepath); err != nil {
		return err
	} else {
		defer f.Close()
		if _, err = f.Write(txtar.Format(pa.raw)); err != nil {
			return err
		}
	}
	return nil
}

// File linting

const policyFilename = "policy.sentinel"
const primaryConfigFilename = "sentinel.hcl"
const overrideConfigFilename = "override.hcl"

func extractSourceFile(files map[string]txtar.File, sentinelVersion string) (lint.File, lint.Issues, error) {
	// Standard Policy File
	if arcFile, ok := files[policyFilename]; ok {
		f, _, issues, err := parsing.ParseSentinelFile(sentinelVersion, policyFilename, arcFile.Data)
		if err != nil {
			return nil, nil, err
		}

		file := lint.PolicyFile{
			File:     f,
			FilePath: policyFilename,
		}

		return file, issues, nil
	}

	// Simple primary config file
	if arcFile, ok := files[primaryConfigFilename]; ok {
		f, issues, err := parsing.ParseSentinelConfigFile(sentinelVersion, primaryConfigFilename, arcFile.Data)
		if err != nil {
			return nil, nil, err
		}

		// Override config file
		if arcFile, ok := files[overrideConfigFilename]; ok {
			override, issues, err := parsing.ParseSentinelConfigFile(sentinelVersion, overrideConfigFilename, arcFile.Data)
			if err != nil {
				return nil, nil, err
			}
			file := lint.ConfigOverrideFile{
				ConfigFile:  override,
				PrimaryFile: f,
				FilePath:    primaryConfigFilename,
			}

			return file, issues, nil
		}

		file := lint.ConfigPrimaryFile{
			ConfigFile:         f,
			ResolvedConfigFile: f,
			FilePath:           primaryConfigFilename,
		}

		return file, issues, nil
	}

	// Only override config file (Unusal but possible)
	if arcFile, ok := files[overrideConfigFilename]; ok {
		f, issues, err := parsing.ParseSentinelConfigFile(sentinelVersion, overrideConfigFilename, arcFile.Data)
		if err != nil {
			return nil, nil, err
		}
		file := lint.ConfigOverrideFile{
			ConfigFile: f,
			FilePath:   primaryConfigFilename,
		}

		return file, issues, nil
	}

	return nil, nil, errors.New("The archive didn't contain any files to lint")
}
