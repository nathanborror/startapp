package ledger

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	ErrInvalidQueryColumns = errors.New("Invalid query colums, expecting: 'added_id, id, datatype, data, time'")
	ErrInvalidTableName    = errors.New("Invalid query table name, expecting: 'FROM record'")
)

var emptyData = []byte("{}")

// Identifier is the interface that wraps methods used to identify items.
type Identifier interface {
	IdentifyID() string
	IdentifyType() string
}

// Applier is the interface that wraps methods for applying ledger values to
// the underlying types they wrap.
type Applier interface {
	ApplyID(id string)
	ApplyTime(created, modified time.Time)
	ApplyCursor(cursor int)
}

// Record is a storable record that typically wraps a given blob of data.
// It maintains an unique identifier and a creation time.
type Record struct {
	db  *sql.DB
	err error

	AddedID  int
	DataType string
	Data     json.RawMessage
	Time     time.Time
	ID       string
}

type RecordSet struct {
	Results []*Record
	err     error
}

// Options describes options for Record.
type Options struct {
	DB *sql.DB
}

// Schema describes the SQL table used to store records.
const Schema = `
CREATE TABLE IF NOT EXISTS record (
	added_id serial PRIMARY KEY,
	id uuid NOT NULL,
	datatype varchar(32) NOT NULL,
	data jsonb NOT NULL,
	time timestamp NOT NULL DEFAULT current_timestamp
)`

// NewRecord returns a new Record that wraps data.
func NewRecord(data Identifier, options Options) *Record {
	r := Record{
		db:       options.DB,
		ID:       data.IdentifyID(),
		DataType: data.IdentifyType(),
		Data:     emptyData,
	}
	if r.ID == "" {
		r.ID = uuid.NewV4().String()
	}
	r.Data, r.err = json.Marshal(data)
	return &r
}

// Query returns records for a given query. It expects each statement
// to return: added_id, id, datatype, data and time in that order.
// TODO: Figure out a way to omit this requirement or enforce it.
func Query(query string, options Options, args ...interface{}) *RecordSet {
	var set RecordSet
	set.err = requireValidQuery(query)
	if set.err != nil {
		return &set
	}
	rows, err := options.DB.Query(query, args...)
	if err != nil {
		set.err = err
		return &set
	}
	defer rows.Close()
	for rows.Next() {
		var r Record
		r.err = rows.Scan(&r.AddedID, &r.ID, &r.DataType, &r.Data, &r.Time)
		r.db = options.DB
		set.Results = append(set.Results, &r)
	}
	return &set
}

// QueryRecord returns a single record for a given custom query. Like Query()
// it expects the statement to return added_id, id, datatype, data and time in
// that order.
func QueryRecord(query string, options Options, args ...interface{}) *Record {
	r := Record{db: options.DB}
	r.err = requireValidQuery(query)
	if r.err != nil {
		return &r
	}
	row := r.db.QueryRow(query, args...)
	r.err = row.Scan(&r.AddedID, &r.ID, &r.DataType, &r.Data, &r.Time)
	return &r
}

// Next prepares the next result record for reading with Scan.
func (r *RecordSet) Next() bool {
	if r.err != nil {
		return false
	}
	return len(r.Results) > 0
}

// Scan ...
func (r *RecordSet) Scan(v Applier) {
	if r.err != nil {
		return
	}
	rec := r.Results[0]
	rec.Scan(v)
	r.Results = append(r.Results[:0], r.Results[1:]...)
	return
}

// Err returns the first error that was encountered by the Record Set.
func (r *RecordSet) Err() error {
	return r.err
}

// Read returns an existing Record that matches id.
func (r *Record) Read() {
	if hasError(r) {
		return
	}
	row := r.db.QueryRow(`
		SELECT added_id, id, datatype, data, time 
		FROM record 
		WHERE id = $1 
		ORDER BY time DESC LIMIT 1`, r.ID)
	r.err = row.Scan(&r.AddedID, &r.ID, &r.DataType, &r.Data, &r.Time)
}

// Write stores a copy of the current Record.
func (r *Record) Write() {
	if hasError(r) {
		return
	}
	row := r.db.QueryRow(`
		INSERT INTO record (id, datatype, data) 
		VALUES ($1, $2, $3) RETURNING added_id, time`, r.ID, r.DataType, r.Data)
	r.err = row.Scan(&r.AddedID, &r.Time)
}

// Delete clears the record data.
func (r *Record) Delete() {
	if hasError(r) {
		return
	}
	r.Data = emptyData
}

// Restore restores the record to a given time.
func (r *Record) Restore(t time.Time) {
	if hasError(r) {
		return
	}
	row := r.db.QueryRow(`
		SELECT id, datatype, data 
		FROM record 
		WHERE id = $1 AND time = $2 
		LIMIT 1`, r.ID, t)
	r.err = row.Scan(&r.ID, &r.DataType, &r.Data)
}

// Scan parses the JSON-encoded data and stores the result in the value
// pointed to by v. The Record ID is applied to the data as well as the Created
// and Modified times.
func (r *Record) Scan(v Applier) {
	if hasError(r) {
		return
	}
	r.err = json.Unmarshal(r.Data, v)
	if hasError(r) {
		return
	}
	var created time.Time
	row := r.db.QueryRow(`
		SELECT time FROM record 
		WHERE id = $1 ORDER BY time ASC LIMIT 1`, r.ID)
	r.err = row.Scan(&created)
	if r.err != nil {
		return
	}
	v.ApplyID(r.ID)
	v.ApplyTime(created, r.Time)
	v.ApplyCursor(r.AddedID)
}

// Err returns the first error that was encountered by the Record.
func (r *Record) Err() error {
	return r.err
}

// IsEmpty reports whether r represents an empty record.
func (r *Record) IsEmpty() bool {
	return bytes.Compare(r.Data, emptyData) == 0
}

// Private

func requireValidQuery(q string) error {
	if !strings.Contains(q, "added_id, id, datatype, data, time") {
		return ErrInvalidQueryColumns
	}
	if !strings.Contains(q, "FROM record") {
		return ErrInvalidTableName
	}
	return nil
}

func hasError(r *Record) bool {
	return r.err != nil
}
