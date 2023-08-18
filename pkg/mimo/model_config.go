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

package mimo

type Config struct {
	ColumnNames   []string
	ColumnConfigs map[string]ColumnConfig
}

type ColumnConfig struct {
	Exclude      []any        // exclude values from the masking rate computation (default: exclude only nil values)
	CoherentWith []string     // list of fields from witch the coherent rate is computed (default: the current field)
	Constraints  []Constraint // list of constraints to validate
}

type Constraint struct {
	Target ConstraintTarget
	Type   ConstraintType
	Value  float64
}

type ConstraintTarget int

const (
	MaskingRate ConstraintTarget = iota
	CohenrentRate
	IdentifiantRate
)

type ConstraintType int

const (
	ShouldEqual ConstraintType = iota
	ShouldBeGreaterThan
	ShouldBeGreaterThanOrEqualTo
	ShouldBeLowerThan
	ShouldBeLessThanOrEqualTo
)

func NewConfig() Config {
	return Config{
		ColumnNames:   []string{},
		ColumnConfigs: map[string]ColumnConfig{},
	}
}

func NewDefaultColumnConfig(columnname string) ColumnConfig {
	return ColumnConfig{
		Exclude:      []any{},
		CoherentWith: []string{columnname},
		Constraints:  []Constraint{},
	}
}
