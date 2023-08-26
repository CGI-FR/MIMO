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
	Version string       `yaml:"version"`
	Columns []YAMLColumn `yaml:"metrics,omitempty"`
}

// YAMLColumn defines how to store a column config in YAML format.
type YAMLColumn struct {
	Name           string                    `yaml:"name"`
	Exclude        []any                     `yaml:"exclude,omitempty"`
	CoherentWith   []string                  `yaml:"coherentWith,omitempty"`
	CoherentSource string                    `yaml:"coherentSource,omitempty"`
	Constraints    map[string]YAMLConstraint `yaml:"constraints,omitempty"`
}

type YAMLConstraint map[string]float64

func LoadConfig(filename string) (mimo.Config, error) {
	config := &YAMLStructure{
		Version: Version,
		Columns: []YAMLColumn{},
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

//nolint:cyclop
func CreateConfig(yamlconfig *YAMLStructure) (mimo.Config, error) {
	config := mimo.NewConfig()

	for _, yamlcolumn := range yamlconfig.Columns {
		column := mimo.ColumnConfig{
			Exclude:        yamlcolumn.Exclude,
			CoherentWith:   yamlcolumn.CoherentWith,
			CoherentSource: yamlcolumn.CoherentSource,
			Constraints:    []mimo.Constraint{},
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
					constraint.Target = mimo.CohenrentRate
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
