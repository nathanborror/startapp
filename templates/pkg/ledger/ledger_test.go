package ledger

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq" // postgres backend
)

var (
	testOptions     Options
	testAccount     MockAccount
	testRestoreTime time.Time
)

func TestMain(m *testing.M) {

	// Setup
	db, err := sql.Open("postgres", "dbname=ledger_test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(Schema); err != nil {
		log.Fatal(err)
	}
	testOptions = Options{DB: db}

	// Run
	code := m.Run()

	// Teardown
	if _, err := db.Exec(`DROP TABLE record`); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

type MockAccount struct {
	ID       string
	Name     string
	Cursor   int
	Created  time.Time
	Modified time.Time
}

func (i *MockAccount) IdentifyID() string     { return i.ID }
func (i *MockAccount) ApplyID(id string)      { i.ID = id }
func (i *MockAccount) ApplyCursor(cursor int) { i.Cursor = cursor }
func (i *MockAccount) IdentifyType() string   { return "mock.account" }
func (i *MockAccount) ApplyTime(created, modified time.Time) {
	i.Created = created
	i.Modified = modified
}

func (i *MockAccount) IsNotEqual(m MockAccount) bool {
	return i.ID != m.ID ||
		i.Name != m.Name ||
		i.Cursor != m.Cursor ||
		i.Created != m.Created ||
		i.Modified != m.Modified
}
func TestWriteRecord(t *testing.T) {
	var mock MockAccount
	in := MockAccount{Name: "Nathan Borror"}
	rec := NewRecord(&in, testOptions)
	rec.Write()
	rec.Scan(&mock)

	if err := rec.Err(); err != nil {
		t.Errorf("new record err = %v", err)
	}
	if mock.Name != in.Name {
		t.Errorf("%s != %s", mock.Name, in.Name)
	}
	if mock.ID == "" {
		t.Errorf("id = %s", mock.ID)
	}
	if mock.Cursor <= 0 {
		t.Errorf("cursor = %d", mock.Cursor)
	}
	if mock.Created.IsZero() {
		t.Errorf("created = %v", mock.Created)
	}
	if mock.Modified.IsZero() {
		t.Errorf("modified = %v", mock.Modified)
	}
	testAccount = mock
}
func TestReadRecord(t *testing.T) {
	var mock MockAccount
	in := MockAccount{ID: testAccount.ID}
	rec := NewRecord(&in, testOptions)
	rec.Read()
	rec.Scan(&mock)

	if err := rec.Err(); err != nil {
		t.Errorf("new record err = %v", err)
	}
	if testAccount.IsNotEqual(mock) {
		t.Errorf("%+v != %+v", testAccount, mock)
	}
}
func TestUpdateRecord(t *testing.T) {
	var mock MockAccount
	testAccount.Name = "Nathan Paul Borror"
	rec := NewRecord(&testAccount, testOptions)
	rec.Write()
	rec.Scan(&mock)

	if err := rec.Err(); err != nil {
		t.Errorf("new record err = %v", err)
	}
	if testAccount.Modified == mock.Modified {
		t.Errorf("%v == %v", testAccount.Modified, mock.Modified)
	}
	if testAccount.Name != mock.Name {
		t.Errorf("%s != %s", testAccount.Name, mock.Name)
	}
	if testAccount.Cursor == mock.Cursor {
		t.Errorf("%d == %d", testAccount.Cursor, mock.Cursor)
	}

	testRestoreTime = testAccount.Created
	testAccount = mock
}
func TestQueryRecords(t *testing.T) {
	var history []MockAccount
	set := Query(`
		SELECT added_id, id, datatype, data, time 
		FROM record 
		WHERE id = $1 
		ORDER BY time DESC`, testOptions, testAccount.ID)
	for set.Next() {
		var mock MockAccount
		set.Scan(&mock)
		history = append(history, mock)
	}

	if err := set.Err(); err != nil {
		t.Errorf("set err = %v", err)
	}
	if len(history) != 2 {
		t.Errorf("history (2) != %d", len(history))
	}
}
func TestQueryRecord(t *testing.T) {
	var mock MockAccount
	rec := QueryRecord(`
		SELECT added_id, id, datatype, data, time 
		FROM record WHERE id = $1
		ORDER BY time DESC`, testOptions, testAccount.ID)
	rec.Scan(&mock)

	if err := rec.Err(); err != nil {
		t.Errorf("err = %v", err)
	}
}
func TestInvalidQuery(t *testing.T) {
	rec1 := QueryRecord(`
		SELECT id FROM record WHERE id = $1
		ORDER BY time DESC`, testOptions, testAccount.ID)

	if err := rec1.Err(); err != ErrInvalidQueryColumns {
		t.Errorf("err = %v", err)
	}

	rec2 := QueryRecord(`
		SELECT added_id, id, datatype, data, time  
		FROM foo WHERE id = $1
		ORDER BY time DESC`, testOptions, testAccount.ID)

	if err := rec2.Err(); err != ErrInvalidTableName {
		t.Errorf("err = %v", err)
	}
}
func TestRestoreRecords(t *testing.T) {
	var mock MockAccount
	rec := NewRecord(&testAccount, testOptions)
	rec.Restore(testRestoreTime)
	rec.Write()
	rec.Scan(&mock)

	if err := rec.Err(); err != nil {
		t.Errorf("set err = %v", err)
	}
	if mock.Name != "Nathan Borror" {
		t.Errorf("%s != Nathan Borror", mock.Name)
	}

	var mock2 MockAccount
	rec2 := NewRecord(&MockAccount{ID: testAccount.ID}, testOptions)
	rec2.Read()
	rec2.Scan(&mock2)

	if err := rec2.Err(); err != nil {
		t.Errorf("set err = %v", err)
	}
	if mock2.Name != "Nathan Borror" {
		t.Errorf("%s != Nathan Borror", mock2.Name)
	}

	testAccount = mock2
}
func TestDeleteRecord(t *testing.T) {
	var mock MockAccount
	rec := NewRecord(&MockAccount{ID: testAccount.ID}, testOptions)
	rec.Read()
	rec.Delete()
	rec.Write()
	rec.Scan(&mock)

	if err := rec.Err(); err != nil {
		t.Errorf("set err = %v", err)
	}
	if mock.Name != "" {
		t.Errorf("name '%s' should be empty", mock.Name)
	}
	if !rec.IsZero() {
		t.Errorf("rec not empty: %s", string(rec.Data))
	}
}
