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
	db, err := sql.Open("postgres", "user=postgres dbname=ledger_test sslmode=disable")
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
	if _, err := db.Exec(`DROP TABLE record, record_index`); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

const MockDataType = "mock.account"

type MockAccount struct {
	ID       string
	Name     string
	Created  time.Time
	Modified time.Time
}

func (i *MockAccount) IdentifyID() string   { return i.ID }
func (i *MockAccount) ApplyID(id string)    { i.ID = id }
func (i *MockAccount) IdentifyType() string { return MockDataType }
func (i *MockAccount) ApplyTime(created, modified time.Time) {
	i.Created = created
	i.Modified = modified
}

func (i *MockAccount) IsNotEqual(m MockAccount) bool {
	return i.ID != m.ID ||
		i.Name != m.Name ||
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
		t.Error(err)
	}
	if mock.Name != in.Name {
		t.Errorf("%s != %s", mock.Name, in.Name)
	}
	if mock.ID == "" {
		t.Errorf("id is empty")
	}
	if mock.Created.IsZero() {
		t.Errorf("created is zero")
	}
	if mock.Modified.IsZero() {
		t.Errorf("modified is zero")
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
		t.Error(err)
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
		t.Error(err)
	}
	if testAccount.Modified == mock.Modified {
		t.Errorf("modified == %v (%v)", testAccount.Modified, mock.Modified)
	}
	if testAccount.Name != mock.Name {
		t.Errorf("name != %s (%s)", mock.Name, testAccount.Name)
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
		t.Error(err)
	}
	if len(history) != 2 {
		t.Errorf("history != 2 (%d)", len(history))
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
		t.Error(err)
	}
}

func TestInvalidQuery(t *testing.T) {
	rec1 := QueryRecord(`
		SELECT id FROM record WHERE id = $1
		ORDER BY time DESC`, testOptions, testAccount.ID)

	if err := rec1.Err(); err != ErrInvalidQueryColumns {
		t.Error(err)
	}

	rec2 := QueryRecord(`
		SELECT added_id, id, datatype, data, time  
		FROM foo WHERE id = $1
		ORDER BY time DESC`, testOptions, testAccount.ID)

	if err := rec2.Err(); err != ErrInvalidTableName {
		t.Error(err)
	}
}

func TestRestoreRecords(t *testing.T) {
	var mock MockAccount
	rec := NewRecord(&testAccount, testOptions)
	rec.Restore(testRestoreTime)
	rec.Write()
	rec.Scan(&mock)

	if err := rec.Err(); err != nil {
		t.Error(err)
	}
	if mock.Name != "Nathan Borror" {
		t.Errorf("name != 'Nathan Borror' (%s)", mock.Name)
	}

	var mock2 MockAccount
	rec2 := NewRecord(&MockAccount{ID: testAccount.ID}, testOptions)
	rec2.Read()
	rec2.Scan(&mock2)

	if err := rec2.Err(); err != nil {
		t.Error(err)
	}
	if mock2.Name != "Nathan Borror" {
		t.Errorf("name != 'Nathan Borror' (%s)", mock2.Name)
	}

	testAccount = mock2
}

func TestDeleteRecord(t *testing.T) {
	var mock MockAccount
	rec := NewRecord(&MockAccount{ID: testAccount.ID}, testOptions)
	rec.Read()
	rec.Delete()
	rec.Scan(&mock)

	if err := rec.Err(); err != nil {
		t.Error(err)
	}
	if mock.Name != "" {
		t.Errorf("name != '' (%s)", mock.Name)
	}
	if !rec.IsZero() {
		t.Errorf("rec not zero (%s)", string(rec.Data))
	}
}

func TestFetchRecords(t *testing.T) {
	for _, mock := range []MockAccount{
		{Name: "Test 0"},
		{Name: "Test 1"},
		{Name: "Test 2"},
		{Name: "Test 3"},
		{Name: "Test 4"},
	} {
		rec := NewRecord(&mock, testOptions)
		rec.Write()
		if err := rec.Err(); err != nil {
			t.Error(err)
		}
	}
	first := 2
	res := Fetch(MockDataType, &first, nil, testOptions)
	if err := res.Err(); err != nil {
		t.Error(err)
	}
	if len(res.Records) != 2 {
		t.Errorf("records != 2 (%d)", len(res.Records))
	}
	if res.Total != 5 {
		t.Errorf("total != 5 (%d)", res.Total)
	}
	if res.HasPrevious != true {
		t.Errorf("hasPrevious != true")
	}
	if res.HasNext != false {
		t.Errorf("HasNext != false")
	}
	if res.Records[0].ID == res.Records[1].ID {
		t.Errorf("records are duplicated")
	}

	// Next page
	first = 4
	res = Fetch(MockDataType, &first, &res.EndID, testOptions)
	if err := res.Err(); err != nil {
		t.Error(err)
	}
	if len(res.Records) == 0 {
		t.Errorf("records == 0")
	}
	if res.HasPrevious != false {
		t.Errorf("hasPrevious != false")
	}
	if res.HasNext != true {
		t.Errorf("HasNext != true")
	}
}
