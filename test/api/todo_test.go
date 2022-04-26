package main

import (
	// "reflect"
	// "strconv"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	// "bytes"
	// "encoding/json"
	// "net/http"
	// "net/http/httptest"
	"testing"
	// "github.com/Z-me/practice-todo-api/api"
	// "github.com/Z-me/practice-todo-api/api/handler"
)

type v2Suite struct {
	db      *gorm.DB
	mock    sqlmock.Sqlmock
	// student Student
}

func getDBMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	s := &v2Suite{}
    db, mock, err := sqlmock.New()
    if err != nil {
        return nil, nil, err
    }

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	s.db, err = gorm.Open(dialector, &gorm.Config{})
    // gdb, err := gorm.Open("postgres", db)
	// gorm.Open(postgres.Dialector{
	// 	Config: &mysql.Config{
	// 		DriverName: "mysql",
	// 		Conn: db,
	// 		SkipInitializeWithVersion: true,
	// 	}
	// }, &gorm.Config{})

    if err != nil {
        return nil, nil, err
    }
    return s.db, mock, nil
}

// TestCreateNewItem DBモックした状態での新規作成テスト
func TestCreateNewItem(t *testing.T) {
	db, mock, err := getDBMock()
    if err != nil {
        t.Fatal(err)
    }
	handleDb, err := db.DB()
	if err != nil {
        t.Fatal(err)
    }
    defer handleDb.Close()
	db.Logger = db.Logger.LogMode(logger.Info)
    // db.LogMode(true)

    r := Repository{DB: db}

    id := "2222"
    name := "BBBB"

    // Mock設定
    mock.ExpectQuery(regexp.QuoteMeta(
        `INSERT INTO "users" ("id","name") VALUES ($1,$2)
         RETURNING "users"."id"`)).
        WithArgs(id, name).
        WillReturnRows(
            sqlmock.NewRows([]string{"id"}).AddRow(id))

    // 実行
    err = r.Create(id, name)
    if err != nil {
        t.Fatal(err)
    }
}

/*
func TestGetTodoList(t *testing.T) {
	// Note: Start test Server
	ts := httptest.NewServer(api.Router())
	defer ts.Close()

	// Note: expected Values
	exp := []handler.Todo{
		{ID: 1,	Title: "最初のTODO", 	Status: "Done", 		Details: "最初に登録されたTodo", 	Priority: "P0"},
		{ID: 2,	Title: "2番目のTODO",	Status: "Backlog",		Details: "2番目に登録されたTodo",	Priority: "P1"},
		{ID: 3,	Title: "3番目TODO",		Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
		{ID: 4,	Title: "4番目TODO",		Status: "Backlog",		Details: "4番目に登録されたTodo",	Priority: "P3"},
		{ID: 5,	Title: "5番目TODO",		Status: "InProgress",	Details: "5番目に登録されたTodo",	Priority: "P1"},
	}

	// Note: Call API
	res, err := http.Get(ts.URL + "/todo")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer res.Body.Close()

	var resData []handler.Todo
	json.NewDecoder(res.Body).Decode(&resData)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v", res.StatusCode)
	}

	if !reflect.DeepEqual(exp, resData) {
		t.Fatalf("resData = %v, want %v", resData, exp)
	}
}

func TestGetTodoItemByID(t *testing.T) {
	// Note: Start test Server
	ts := httptest.NewServer(api.Router())
	defer ts.Close()

	// Note: expected Values
	exp := handler.Todo{
		ID: 		1,
		Title: 		"最初のTODO",
		Status: 	"Done",
		Details: 	"最初に登録されたTodo",
		Priority: 	"P0",
	}

	// Note: Call API
	res, err := http.Get(ts.URL + "/todo/1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer res.Body.Close()

	var resData handler.Todo
	json.NewDecoder(res.Body).Decode(&resData)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("[Get Todo Item By ID] Expected status code 200, got %v", res.StatusCode)
	}

	if !reflect.DeepEqual(exp, resData) {
		t.Fatalf("[Get Todo Item By ID] resData = %v, want %v", resData, exp)
	}
}

func TestAddItem(t *testing.T) {
	// Note: Start test Server
	ts := httptest.NewServer(api.Router())
	defer ts.Close()

	// Note: expected Values
	exp := []handler.Todo{
		{ID: 1,	Title: "最初のTODO",	Status: "Done",			Details: "最初に登録されたTodo",	Priority: "P0"},
		{ID: 2,	Title: "2番目のTODO",	Status: "Backlog",		Details: "2番目に登録されたTodo",	Priority: "P1"},
		{ID: 3,	Title: "3番目TODO",		Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
		{ID: 4,	Title: "4番目TODO",		Status: "Backlog",		Details: "4番目に登録されたTodo",	Priority: "P3"},
		{ID: 5,	Title: "5番目TODO",		Status: "InProgress",	Details: "5番目に登録されたTodo",	Priority: "P1"},
		{ID: 6,	Title: "6番目TODO",		Status: "InProgress",	Details: "6番目に登録されたTodo",	Priority: "P0"},
	}
	payload := handler.Todo{ID: 6,Title: "6番目TODO",Status: "InProgress",Details: "6番目に登録されたTodo",Priority: "P0"}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Note: Call POST API
	postRes, err := http.Post(ts.URL + "/todo", "application/json", bytes.NewBuffer(payloadJson))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer postRes.Body.Close()

	// Note: Call GET API
	getRes, err := http.Get(ts.URL + "/todo")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer getRes.Body.Close()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var resData []handler.Todo
	json.NewDecoder(getRes.Body).Decode(&resData)

	if postRes.StatusCode != http.StatusCreated || getRes.StatusCode != http.StatusOK {
		t.Fatalf("[Post New Todo Item] Expected status code 200, got %v and %v", postRes.StatusCode, getRes.StatusCode)
	}

	if !reflect.DeepEqual(exp, resData) {
		t.Fatalf("[Post New Todo Item] resData = %v, want %v", resData, exp)
	}
}

func TestUpdateItem(t *testing.T) {
	// Note: Start test Server
	ts := httptest.NewServer(api.Router())
	defer ts.Close()

	targetId := 5
	payload := handler.Payload{Title: "更新された5番目TODO",Status: "Done",Details: "5番目に登録され、その後更新されたTodo",Priority: "P1"}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Note: expected Values
	exp := handler.Todo{
		ID: 		targetId,
		Title: 		payload.Title,
		Status: 	payload.Status,
		Details: 	payload.Details,
		Priority: 	payload.Priority,
	}

	// Note: Call POST API
	client := &http.Client{}
	req, err := http.NewRequest("PUT", ts.URL + "/todo/" + strconv.Itoa(targetId), bytes.NewBuffer(payloadJson))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer res.Body.Close()

	var resData handler.Todo
	json.NewDecoder(res.Body).Decode(&resData)

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("[Update Todo Item] Expected status code 201, got %v", res.StatusCode)
	}

	if !reflect.DeepEqual(exp, resData) {
		t.Fatalf("[Update Todo Item] resData = %v, want %v", exp, resData)
	}
}

func TestUpdateStateById(t *testing.T) {
	// Note: Start test Server
	ts := httptest.NewServer(api.Router())
	defer ts.Close()

	targetId := 5
	payload := handler.StatusPayload{Status: "Backlog"}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	exp := handler.Todo{
		ID: 		targetId,
		Title: 		"5番目TODO",
		Status: 	payload.Status,
		Details: 	"5番目に登録されたTodo",
		Priority: 	"P1",
	}

	client := &http.Client{}
	req, err := http.NewRequest("PATCH", ts.URL + "/todo/" + strconv.Itoa(targetId) + "/status", bytes.NewBuffer(payloadJson))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer req.Body.Close()

	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var resData handler.Todo
	json.NewDecoder(res.Body).Decode(&resData)

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("[Update Todo Item] Expected status code 201, got %v", res.StatusCode)
	}

	if !reflect.DeepEqual(exp, resData) {
		t.Fatalf("[Update Todo Item] resData = %v, want %v", exp, resData)
	}
}

func TestDeleteItem(t *testing.T) {
	// Note: Start test Server
	ts := httptest.NewServer(api.Router())
	defer ts.Close()

	// Note: expected Values
	targetId := 5
	exp := []handler.Todo{
		{ID: 1,	Title: "最初のTODO",	Status: "Done",			Details: "最初に登録されたTodo",	Priority: "P0"},
		{ID: 2,	Title: "2番目のTODO",	Status: "Backlog",		Details: "2番目に登録されたTodo",	Priority: "P1"},
		{ID: 3,	Title: "3番目TODO",		Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
		{ID: 4,	Title: "4番目TODO",		Status: "Backlog",		Details: "4番目に登録されたTodo",	Priority: "P3"},
	}

	// Note: Call DELETE API
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", ts.URL + "/todo/" + strconv.Itoa(targetId), nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer res.Body.Close()

	var resData []handler.Todo
	json.NewDecoder(res.Body).Decode(&resData)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("[Update Todo Item] Expected status code 200, got %v", res.StatusCode)
	}

	if !reflect.DeepEqual(exp, resData) {
		t.Fatalf("resData = %v, want %v", exp, resData)
	}
}

func TestAnomaly(t *testing.T) {
	// Note: Start test Server
	ts := httptest.NewServer(api.Router())
	defer ts.Close()

	cases := []struct{
		name	string
		url 	string
		client	string
		payload string
		expect 	int
	}{
		{
			name:		"[404-01]GET 404",
			url: 		"/error",
			client:   	"GET",
			payload: 	"",
			expect: 	http.StatusNotFound,
		},
		{
			name: 		"[404-02]POST 404",
			url: 		"/error",
			client:		"POST",
			payload: 	"",
			expect:		http.StatusNotFound,
		},
		{
			name: 		"[404-03]PUT 404",
			url: 		"/error",
			client:   	"PUT",
			payload: 	"",
			expect: 	http.StatusNotFound,
		},
		{
			name: 		"[404-04]DELETE 404",
			url: 		"/error",
			client:   	"DELETE",
			payload: 	"",
			expect: 	http.StatusNotFound,
		},
		{
			name:		"[400-01]no payload on new item",
			url: 		"/todo",
			client:   	"POST",
			payload: 	"",
			expect: 	http.StatusBadRequest,
		},
		{
			name:		"[400-02]invalid payload on create new item",
			url: 		"/todo",
			client:   	"POST",
			payload: 	`{"message":"invalid payload"}`,
			expect: 	http.StatusBadRequest,
		},
		{
			name:		"[400-03]invalid payload on change item",
			url: 		"/todo/error",
			client:   	"PUT",
			payload: 	`{"message":"missing"}`,
			expect: 	http.StatusBadRequest,
		},
		{
			name:   	"[400-04]invalid Method on create item",
			url: 		"/todo/1",
			client:   	"PUT",
			payload: 	"",
			expect: 	http.StatusBadRequest,
		},
		{
			name:		"[400-05]invalid payload on change item",
			url: 		"/todo/1/status",
			client:   	"PATCH",
			payload: 	`{"message":"invalid payload"}`,
			expect: 	http.StatusBadRequest,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(caseNameHelper(t, c.name,  c.client, c.url), func(t *testing.T) {
			client := &http.Client{}
			req, err := http.NewRequest(c.client, ts.URL + c.url, bytes.NewBuffer([]byte(c.payload)))
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			defer res.Body.Close()

			if res.StatusCode != c.expect {
				t.Fatalf("[Update Todo Item] Expected status code %v, got %v", c.expect, res.StatusCode)
			}
		})
	}
}

func caseNameHelper(t *testing.T, name string, client string, url string) string {
	t.Helper()
	return name + "のテスト :: [" + client + "] " + url
}
*/
