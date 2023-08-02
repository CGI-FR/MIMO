package mimo

import "errors"

var (
	// ErrReadingDataRow is returned when an error occurs while reading a data row.
	ErrReadingDataRow = errors.New("error while reading datarow")

	// ErrOrphanRow is returned when a original row does not have a masked version, or the other way around.
	ErrOrphanRow = errors.New("error datarow is orphan")
)
