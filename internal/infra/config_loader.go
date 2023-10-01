// Copyright (C) 2023 CGI France
//
// This file is part of MIMO.
//
// MIMO is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// MIMO is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with MIMO.  If not, see <http://www.gnu.org/licenses/>.

package infra

import (
	"fmt"
	"os"

	"github.com/cgi-fr/mimo/pkg/mimo"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// Version of the YAML strcuture.
const Version string = "1"

// YAMLStructure of the file.
type YAMLStructure struct {
	Version      string           `yaml:"version"`
	Columns      []YAMLColumn     `yaml:"metrics,omitempty"`
	Preprocesses []YAMLPreprocess `yaml:"preprocess,omitempty"`
}

// YAMLColumn defines how to store a column config in YAML format.
type YAMLColumn struct {
	Name            string                    `yaml:"name"`
	Exclude         []any                     `yaml:"exclude,omitempty"`
	ExcludeTemplate string                    `yaml:"excludeTemplate,omitempty"`
	CoherentWith    []string                  `yaml:"coherentWith,omitempty"`
	CoherentSource  string                    `yaml:"coherentSource,omitempty"`
	Constraints     map[string]YAMLConstraint `yaml:"constraints,omitempty"`
	Alias           string                    `yaml:"alias,omitempty"`
}

type YAMLPreprocess struct {
	Path  string `yaml:"path"`
	Value string `yaml:"value"`
}

type YAMLConstraint map[string]float64

func LoadConfig(filename string) (mimo.Config, error) {
	config := &YAMLStructure{
		Version:      Version,
		Columns:      []YAMLColumn{},
		Preprocesses: []YAMLPreprocess{},
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return mimo.NewConfig(), fmt.Errorf("%w: %s", ErrConfigFileNotExists, filename)
	}

	log.Debug().Str("file", filename).Msg("loading config from file")

	dat, err := os.ReadFile(filename)
	if err != nil {
		return mimo.NewConfig(), fmt.Errorf("%w: %s", err, filename)
	}

	err = yaml.Unmarshal(dat, config)
	if err != nil {
		return mimo.NewConfig(), fmt.Errorf("%w: %s", err, filename)
	}

	if config.Version != Version {
		return mimo.NewConfig(), fmt.Errorf("%w: %s", ErrConfigInvalidVersion, filename)
	}

	return CreateConfig(config)
}

//nolint:cyclop,funlen
func CreateConfig(yamlconfig *YAMLStructure) (mimo.Config, error) {
	config := mimo.NewConfig()

	if err := CreatePreprocesses(yamlconfig, &config); err != nil {
		return config, err
	}

	for _, yamlcolumn := range yamlconfig.Columns {
		excludeTmpl, err := mimo.NewTemplate(yamlcolumn.ExcludeTemplate)
		if err != nil {
			return config, fmt.Errorf("%w", err)
		}

		coherentTmpl, err := mimo.NewTemplate(yamlcolumn.CoherentSource)
		if err != nil {
			return config, fmt.Errorf("%w", err)
		}

		column := mimo.ColumnConfig{
			Exclude:         yamlcolumn.Exclude,
			ExcludeTemplate: excludeTmpl,
			CoherentWith:    yamlcolumn.CoherentWith,
			CoherentSource:  coherentTmpl,
			Constraints:     []mimo.Constraint{},
			Alias:           yamlcolumn.Alias,
		}

		for target, yamlconstraint := range yamlcolumn.Constraints {
			for constraintType, value := range yamlconstraint {
				constraint := mimo.Constraint{
					Target: 0,
					Type:   0,
					Value:  value,
				}

				switch target {
				case "maskingRate":
					constraint.Target = mimo.MaskingRate
				case "coherentRate":
					constraint.Target = mimo.CoherentRate
				case "identifiantRate":
					constraint.Target = mimo.IdentifiantRate
				default:
					return config, fmt.Errorf("%w: %s", ErrConfigInvalidConstraintTarget, target)
				}

				switch constraintType {
				case "shouldEqual":
					constraint.Type = mimo.ShouldEqual
				case "shouldBeGreaterThan":
					constraint.Type = mimo.ShouldBeGreaterThan
				case "shouldBeGreaterThanOrEqualTo":
					constraint.Type = mimo.ShouldBeGreaterThanOrEqualTo
				case "shouldBeLowerThan":
					constraint.Type = mimo.ShouldBeLowerThan
				case "shouldBeLessThanOrEqualTo":
					constraint.Type = mimo.ShouldBeLessThanOrEqualTo
				default:
					return config, fmt.Errorf("%w: %s", ErrConfigInvalidConstraintType, constraintType)
				}

				column.Constraints = append(column.Constraints, constraint)
			}
		}

		config.ColumnNames = append(config.ColumnNames, yamlcolumn.Name)
		config.ColumnConfigs[yamlcolumn.Name] = column
	}

	return config, nil
}

func CreatePreprocesses(yamlconfig *YAMLStructure, config *mimo.Config) error {
	for _, yamlpreprocess := range yamlconfig.Preprocesses {
		valueTmpl, err := mimo.NewTemplate(yamlpreprocess.Value)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		preprocess := mimo.PreprocessConfig{
			Path:  yamlpreprocess.Path,
			Value: valueTmpl,
		}
		config.PreprocessConfigs = append(config.PreprocessConfigs, preprocess)
	}

	return nil
}
