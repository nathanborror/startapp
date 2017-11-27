package ledger

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	ErrInvalidQueryColumns = errors.New("Invalid query colums, expecting: 'added_id, id, datatype, data, time'")
	ErrInvalidTableName    = errors.New("Invalid query table name, expecting: 'FROM record'")
	ErrRecordNotFound      = errors.New("Record not found")
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
);
CREATE TABLE IF NOT EXISTS record_index (
	added_id serial PRIMARY KEY,
	id uuid NOT NULL UNIQUE,
	datatype varchar(32) NOT NULL
);
CREATE INDEX IF NOT EXISTS record_index_id ON record_index(id);
CREATE INDEX IF NOT EXISTS record_index_datatype ON record_index(datatype);`

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
func Query(query string, options Options, args ...interface{}) *Result {
	var result Result
	result.err = requireValidQuery(query)
	if result.err != nil {
		return &result
	}
	rows, err := options.DB.Query(query, args...)
	if err != nil {
		result.err = err
		return &result
	}
	defer rows.Close()
	for rows.Next() {
		var r Record
		r.err = rows.Scan(&r.AddedID, &r.ID, &r.DataType, &r.Data, &r.Time)
		r.db = options.DB
		result.Records = append(result.Records, r)
	}
	if err := rows.Err(); err != nil {
		result.err = err
	}
	return &result
}

// History returns all the versions for a given ID.
func History(id string, options Options) *Result {
	var result Result
	rows, err := options.DB.Query(`
		SELECT added_id,id,datatype,data,time 
		FROM record 
		WHERE id = $1 
		ORDER BY time DESC`, id)
	if err != nil {
		result.err = err
		return &result
	}
	defer rows.Close()
	for rows.Next() {
		var r Record
		r.err = rows.Scan(&r.AddedID, &r.ID, &r.DataType, &r.Data, &r.Time)
		r.db = options.DB
		result.Records = append(result.Records, r)
	}
	if err := rows.Err(); err != nil {
		result.err = err
		return &result
	}
	result.Total = len(result.Records)
	if result.Total > 0 {
		result.StartID = result.Records[0].ID
		result.EndID = result.Records[len(result.Records)-1].ID
	}
	return &result
}

// Fetch returns a paginated list of Records.
func Fetch(datatype string, first *int, beforeID *string, options Options) *Result {
	var (
		beforeAddedID int
		beforeCount   int
		afterCount    int
		latestID      int
		ids           []string
		result        Result
	)

	// Require 'first' to be greater than 0
	if first == nil || first != nil && *first == 0 {
		result.err = fmt.Errorf("No results returned because 'first' was nil")
		return &result
	}

	// Fetch total
	result.err = options.DB.QueryRow(`SELECT COUNT(*) FROM record_index WHERE datatype = $1`, datatype).Scan(&result.Total)
	if result.err != nil {
		result.err = fmt.Errorf("Error fetching total: %v", result.err)
		return &result
	}

	// Determine added_id of latest record
	result.err = options.DB.QueryRow(`
		SELECT (added_id) FROM record_index
		WHERE datatype = $1 ORDER BY added_id DESC LIMIT 1`, datatype).Scan(&latestID)
	if result.err == sql.ErrNoRows {
		result.err = nil
		return &result
	}
	if result.err != nil {
		result.err = fmt.Errorf("Error fetching latest ID: %v", result.err)
		return &result
	}

	// Determine the added_id for the beforeID record
	if beforeID != nil {
		result.err = options.DB.QueryRow(`
			SELECT (added_id) FROM record_index 
			WHERE id = $1 AND datatype = $2`, beforeID, datatype).Scan(&beforeAddedID)
	} else {
		beforeAddedID = latestID + 1 // Increment so latest record is included
	}
	if result.err != nil {
		result.err = fmt.Errorf("Error fetching added ID: %v", result.err)
		return &result
	}

	// Determine how may records exist before and after the beforeAddedID
	if beforeAddedID != 0 {
		result.err = options.DB.QueryRow(`
			SELECT
				(SELECT COUNT(added_id) FROM record_index
					WHERE added_id < $1 AND datatype = $2) AS beforeCount,
				(SELECT COUNT(added_id) FROM record_index
					WHERE added_id > $1 AND datatype = $2) AS afterCount
			`, beforeAddedID, datatype).Scan(&beforeCount, &afterCount)
	}

	// Fetch Record IDs from index
	rows, err := options.DB.Query(`
		SELECT id FROM record_index 
		WHERE added_id < $1 
		AND datatype = $2 
		ORDER BY added_id DESC
		LIMIT $3`, beforeAddedID, datatype, first)
	if err != nil {
		result.err = fmt.Errorf("Error fetching records before (%d): %v", beforeAddedID, err)
		return &result
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			result.err = err
			return &result
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		result.err = err
		return &result
	}
	result.StartID = ids[0]
	result.EndID = ids[len(ids)-1]

	// Fetch Records
	records := FetchIn(ids, options)
	if err := records.Err(); err != nil {
		result.err = err
		return &result
	}
	result.Records = records.Records

	if beforeAddedID != 0 {
		result.HasNext = afterCount > 0
		result.HasPrevious = beforeCount > len(result.Records)
	} else {
		result.HasPrevious = result.Total > len(result.Records)
		result.HasNext = false
	}
	return &result
}

// FetchIn returns a Result with records for a given set of IDs.
func FetchIn(recordIDs []string, options Options) *Result {
	var result Result
	any := fmt.Sprintf("{%s}", strings.Join(recordIDs, ","))
	rows, err := options.DB.Query(`
		SELECT added_id, id, datatype, data, time FROM record 
		WHERE id = ANY($1)`, any)
	if err != nil {
		result.err = err
		return &result
	}
	defer rows.Close()
	for rows.Next() {
		var r Record
		r.err = rows.Scan(&r.AddedID, &r.ID, &r.DataType, &r.Data, &r.Time)
		r.db = options.DB
		result.Records = append(result.Records, r)
	}
	if err := rows.Err(); err != nil {
		result.err = err
		return &result
	}
	return &result
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

// Result is a result set of Records.
type Result struct {
	Records     []Record
	Total       int
	HasNext     bool
	HasPrevious bool
	StartID     string
	EndID       string
	err         error
}

// Next prepares the next result record for reading with Scan.
func (r *Result) Next() bool {
	if r.err != nil {
		return false
	}
	return len(r.Records) > 0
}

// Scan parses the JSON-encoded data and stores the result in the value
// pointed to by v. The Record ID is applied to the data as well as the Created
// and Modified times.
func (r *Result) Scan(v Applier) {
	if r.err != nil {
		return
	}
	if len(r.Records) == 0 {
		return
	}
	rec := r.Records[0]
	rec.Scan(v)
	r.Records = append(r.Records[:0], r.Records[1:]...)
	return
}

// Err returns the first error that was encountered by the Record Set.
func (r *Result) Err() error {
	return r.err
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

// Read returns an existing Record that matches id.
func (r *Record) Read() {
	if hasError(r) {
		return
	}
	row := r.db.QueryRow(`
		SELECT added_id, id, datatype, data, time 
		FROM record 
		WHERE id = $1 AND datatype = $2
		ORDER BY time DESC LIMIT 1`, r.ID, r.DataType)
	r.err = row.Scan(&r.AddedID, &r.ID, &r.DataType, &r.Data, &r.Time)
}

// Write stores a copy of the current Record and indexes it if an index didn't
// already exist.
func (r *Record) Write() {
	if hasError(r) {
		return
	}
	row := r.db.QueryRow(`
		INSERT INTO record (id, datatype, data) 
		VALUES ($1, $2, $3) RETURNING added_id, time`, r.ID, r.DataType, r.Data)
	if err := row.Scan(&r.AddedID, &r.Time); err != nil {
		r.err = err
		return
	}
	if _, err := r.db.Exec(`INSERT INTO record_index (id,datatype) VALUES ($1,$2) ON CONFLICT DO NOTHING`, r.ID, r.DataType); err != nil {
		r.err = err
		return
	}
}

// Delete writes a new record with zeroed out data and removes it from the index.
func (r *Record) Delete() {
	if hasError(r) {
		return
	}
	r.Data = emptyData
	// Write new record with zeroed data
	_, err := r.db.Exec(`
		INSERT INTO record (id, datatype, data) 
		VALUES ($1, $2, $3) RETURNING added_id, time`, r.ID, r.DataType, r.Data)
	if err != nil {
		r.err = err
		return
	}
	// Remove from index
	_, err = r.db.Exec(`DELETE FROM record_index WHERE id = $1`, r.ID)
	if err != nil {
		r.err = err
		return
	}
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
}

// Err returns the first error that was encountered by the Record.
func (r *Record) Err() error {
	if r.err == sql.ErrNoRows {
		return ErrRecordNotFound
	}
	return r.err
}

// IsZero reports whether r represents an empty record.
func (r *Record) IsZero() bool {
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
