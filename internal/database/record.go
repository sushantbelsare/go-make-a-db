package database

// Record represents a single row in a table
type Record map[string]string

// Copy creates a deep copy of the record
func (r Record) Copy() Record {
	newRecord := make(Record)
	for k, v := range r {
		newRecord[k] = v
	}
	return newRecord
}

// Get retrieves the value for a given column
func (r Record) Get(column string) (string, bool) {
	value, exists := r[column]
	return value, exists
}

// Set updates or adds a value for a given column
func (r Record) Set(column string, value string) {
	r[column] = value
}

// Delete removes a column from the record
func (r Record) Delete(column string) {
	delete(r, column)
}

// Columns returns a slice of all column names in the record
func (r Record) Columns() []string {
	columns := make([]string, 0, len(r))
	for column := range r {
		columns = append(columns, column)
	}
	return columns
}

// Values returns a slice of all values in the record
func (r Record) Values() []string {
	values := make([]string, 0, len(r))
	for _, value := range r {
		values = append(values, value)
	}
	return values
}

// IsEmpty checks if the record has no columns
func (r Record) IsEmpty() bool {
	return len(r) == 0
}