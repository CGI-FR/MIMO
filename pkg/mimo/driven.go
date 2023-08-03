package mimo

type DataRowReader interface {
	ReadDataRow() (DataRow, error)
}

type EventSubscriber interface {
	NewField(fieldname string)
	FirstNonMaskedValue(fieldname string, value any)
}
