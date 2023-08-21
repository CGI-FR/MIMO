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

import "errors"

var (
	// ErrConfigFileNotExists is returned when a config file doesn't exist.
	ErrConfigFileNotExists = errors.New("error config file does not exist")

	// ErrConfigInvalidVersion is returned when a config file has an invalid version.
	ErrConfigInvalidVersion = errors.New("invalid version in config file")

	// ErrConfigInvalidConstraintType is returned when a config file has an invalid constraint type.
	ErrConfigInvalidConstraintType = errors.New("invalid constraint type in config file")

	// ErrConfigInvalidConstraintTarget is returned when a config file has an invalid constraint target.
	ErrConfigInvalidConstraintTarget = errors.New("invalid constraint target in config file")
)
