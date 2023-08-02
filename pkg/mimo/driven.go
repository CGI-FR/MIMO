package mimo

type DataRowReader interface {
	ReadDataRow() (DataRow, error)
}

type DataRowWriter interface {
	WriteDataRow(row DataRow) error
}
