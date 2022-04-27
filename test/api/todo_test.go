package main

import (
	"bytes"
	"strconv"

	// "reflect"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Z-me/practice-todo-api/api"
	"github.com/Z-me/practice-todo-api/api/handler"
	"github.com/Z-me/practice-todo-api/api/model"
	db "github.com/Z-me/practice-todo-api/middleware"
)

func caseNameHelper(t *testing.T, name string, client string, url string) string {
	t.Helper()
	return name + "のテスト[" + client + "]" + url
}

func TestCreateItem(t *testing.T) {
	// Note: Start test Server
	ts := httptest.NewServer(api.Router())
	defer ts.Close()

	// Note: Start Connect DB
	err := db.ConnectDB()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer db.DisconnectDB()

	// Note: each values
	nextID := db.GetNextID()
	cases := []struct{
		name 		string
		url			string
		client		string
		status 		int
		isError 	bool
		payload		string
		expected	model.Todo
		need2Delete bool
	}{
		{
			name: 		"正常系: 新規追加",
			url: 		"/todo",
			client: 	"POST",
			status:	 	http.StatusCreated,
			isError: 	false,
			payload:	`{"title": "Test TODO", "status": "Done", "details": "test_todo", "priority": "P0"}`,
			expected:	model.Todo{
				ID: 		nextID,
				Title:		"Test TODO",
				Status:		"Done",
				Details:	"test_todo",
				Priority:	"P0",
			},
			need2Delete: true,
		},
		{
			name: 		"異常系: 新規追加: 404",
			url: 		"/todo",
			client: 	"PUT",
			status:	 	http.StatusNotFound,
			isError: 	true,
			payload:	`{"title": "Test TODO", "status": "Done", "details": "test_todo", "priority": "P0"}`,
			expected:	model.Todo{},
			need2Delete: false,
		},
		{
			name: 		"異常系: 新規追加: 400",
			url: 		"/todo",
			client: 	"POST",
			status:	 	http.StatusBadRequest,
			isError: 	true,
			payload:	`{"Message": "Bad Request"}`,
			expected:	model.Todo{},
			need2Delete: false,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(caseNameHelper(t, c.name, c.client, c.url), func(t *testing.T) {
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

			if res.StatusCode != c.status {
				t.Fatalf("[New Item] Expected status code %v, got %v", c.status, res.StatusCode)
			}
			var resData handler.Todo
			json.NewDecoder(res.Body).Decode(&resData)

			// CreatedAtなどは比較したくないので除外
			if !c.isError {
				if	c.expected.ID != uint(resData.ID) {
					t.Fatalf("ID: want %v, resData = %v", c.expected.ID, resData.ID)
				}
				if c.expected.Title != resData.Title {
					t.Fatalf("Title: want %v, resData = %v", c.expected.Title, resData.Title)
				}
				if c.expected.Status != resData.Status {
					t.Fatalf("Status: want %v, resData = %v", c.expected.Status, resData.Status)
				}
				if c.expected.Details != resData.Details {
					t.Fatalf("Details: want %v, resData = %v", c.expected.Details, resData.Details)
				}
				if c.expected.Priority != resData.Priority {
						t.Fatalf("Priority: want %v, resData = %v", c.expected.Priority, resData.Priority)
				}
			}

			// 終了処理
			if c.need2Delete {
				err = db.ConnectDB()
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
				defer db.DisconnectDB()

				_, err := db.DeleteItem(uint(resData.ID))
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestUpdateItem(t *testing.T) {
	// Note: Start test Server
	ts := httptest.NewServer(api.Router())
	defer ts.Close()

	// Note: Start Connect DB
	err := db.ConnectDB()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer db.DisconnectDB()

	// Note: 事前処理
	target := model.Payload{
		Title:		"Test TODO",
		Status:		"Done",
		Details:	"test_todo",
		Priority:	"P0",
	}
	res, err := db.AddNewTodo(target)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	nextID := res.ID

	cases := []struct{
		name 		string
		url			string
		client		string
		status 		int
		isError 	bool
		payload		string
		expected	model.Todo
		need2Delete bool
	}{
		{
			name: 		"正常系: 更新",
			url: 		"/todo/" + strconv.Itoa(int(nextID)),
			client: 	"PUT",
			status:	 	http.StatusCreated,
			isError: 	false,
			payload:	`{"title": "Changed TODO", "status": "Done", "details": "changed_todo", "priority": "P0"}`,
			expected:	model.Todo{
				ID: 		nextID,
				Title:		"Changed TODO",
				Status:		"Done",
				Details:	"changed_todo",
				Priority:	"P0",
			},
			need2Delete: true,
		},
		{
			name: 		"異常系: 更新: 404",
			url: 		"/todo/" + strconv.Itoa(int(nextID)),
			client: 	"POST",
			status:	 	http.StatusNotFound,
			isError: 	true,
			payload:	`{"title": "Changed TODO", "status": "Done", "details": "changed_todo", "priority": "P0"}`,
			expected:	model.Todo{},
			need2Delete: false,
		},
		{
			name: 		"異常系: 新規追加: 400",
			url: 		"/todo/" + strconv.Itoa(int(nextID)),
			client: 	"PUT",
			status:	 	http.StatusBadRequest,
			isError: 	true,
			payload:	`{"Message": "Bad Request"}`,
			expected:	model.Todo{},
			need2Delete: false,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(caseNameHelper(t, c.name, c.client, c.url), func(t *testing.T) {
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

			if res.StatusCode != c.status {
				t.Fatalf("[New Item] Expected status code %v, got %v", c.status, res.StatusCode)
			}
			var resData handler.Todo
			json.NewDecoder(res.Body).Decode(&resData)

			// CreatedAtなどは比較したくないので除外
			if !c.isError {
				if	c.expected.ID != uint(resData.ID) {
					t.Fatalf("ID: want %v, resData = %v", c.expected.ID, resData.ID)
				}
				if c.expected.Title != resData.Title {
					t.Fatalf("Title: want %v, resData = %v", c.expected.Title, resData.Title)
				}
				if c.expected.Status != resData.Status {
					t.Fatalf("Status: want %v, resData = %v", c.expected.Status, resData.Status)
				}
				if c.expected.Details != resData.Details {
					t.Fatalf("Details: want %v, resData = %v", c.expected.Details, resData.Details)
				}
				if c.expected.Priority != resData.Priority {
					t.Fatalf("Priority: want %v, resData = %v", c.expected.Priority, resData.Priority)
				}
			}
		})
	}
	// Note: 事後削除処理
	err = db.ConnectDB()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer db.DisconnectDB()

	db.DeleteItem(nextID)
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
