package mimo

type DataRowReader interface {
	ReadDataRow() (DataRow, error)
}

type DataRowWriter interface {
	WriteDataRow(row DataRow) error
}

type EventSubscriber interface {
	NewField(fieldname string)
	FirstNonMaskedValue(fieldname string, value any)
}
