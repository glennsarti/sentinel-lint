package runner

import (
	"errors"
	"fmt"

	lint "github.com/glennsarti/sentinel-lint/lint"
	"github.com/glennsarti/sentinel-parser/features"
)

var _ lint.Runner = &runnerImpl{}

func NewRunner(config lint.Config, rulesets lint.RuleSetList, file lint.File) (lint.Runner, error) {
	if rulesets == nil {
		return nil, errors.New("rulesets is nil")
	}

	if file == nil {
		return nil, errors.New("file is nil")
	}

	var sentVersion string
	if ok, ver := features.ValidateSentinelVersion(config.SentinelVersion); ok {
		sentVersion = ver
	} else {
		return nil, fmt.Errorf("sentinel version %q is not valid", config.SentinelVersion)
	}

	return &runnerImpl{
		config: lint.Config{
			SentinelVersion: sentVersion,
			FailFast:        config.FailFast,
		},
		ruleSets: rulesets,
		files:    []lint.File{file},
	}, nil
}

type runnerImpl struct {
	config   lint.Config
	ruleSets lint.RuleSetList
	files    []lint.File
}

func (r *runnerImpl) Config() lint.Config {
	return r.config
}

func (r *runnerImpl) Run() (lint.Issues, error) {
	ruleCtx := ruleContext{
		files:           r.files,
		sentinelVersion: r.config.SentinelVersion,
	}

	issues := make(lint.Issues, 0)

	if r.ruleSets != nil {
		for _, set := range r.ruleSets {
			if i, err := set.Run(ruleCtx); err != nil {
				return issues, err
			} else if i != nil {
				issues = append(issues, *i...)
			}
		}
	}

	return issues, nil
}
