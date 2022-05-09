package main

import (
	"bytes"

	// "fmt"

	// "fmt"

	"strconv"
	"time"

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

func compareTodoList(target model.TodoList, expect []handler.Todo) bool {
	for i, v := range expect {
		if uint(v.ID) != target[i].ID {
			return false
		}
		if v.Title != target[i].Title {
			return false
		}
		if v.Status != target[i].Status {
			return false
		}
		if v.Details != target[i].Details {
			return false
		}
		if v.Priority != target[i].Priority {
			return false
		}
	}
	return true
}

func TestGetTodoList(t *testing.T) {
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
	// todoList := model.TodoList{}
	expected, err := db.GetTodoList()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	cases := []struct{
		name		string
		url			string
		method		string
		status		int
		isError		bool
		expected	model.TodoList
	}{
		{
			name:		"正常系: TodoList取得",
			url:		"/todo",
			method:		"GET",
			status:		http.StatusOK,
			isError:	false,
			expected:	expected,
		},
		{
			name:		"異常系: TodoList取得: 404",
			url:		"/TODOOOO",
			method:		"GET",
			status:		http.StatusNotFound,
			isError:	true,
			expected:	model.TodoList{},
		},
	}

	for _, c := range cases {
		t.Run(caseNameHelper(t, c.name, c.method, c.url), func(t *testing.T) {
			client := &http.Client{}
			req, err := http.NewRequest(c.method, ts.URL + c.url, nil)
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
			var resData []handler.Todo
			json.NewDecoder(res.Body).Decode(&resData)

			if !c.isError {
				if len(c.expected) != len(resData) {
					t.Fatalf("Length: want %v items, resData = %v items", len(c.expected), len(resData))
				}
				if !compareTodoList(c.expected, resData) {
					t.Fatalf("Contents: want %v, resData = %v", c.expected, resData)
				}
			}
		})
	}
}

func TestGetTodoItem(t *testing.T) {
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
		Priority:	"P2",
	}
	res, err := db.AddNewTodo(target)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	nextID := res.ID
	createdAt := res.CreatedAt

	cases := []struct{
		name		string
		url			string
		method		string
		status		int
		isError		bool
		expected	model.Todo
	}{
		{
			name:		"正常系: Item取得",
			url:		"/todo/" + strconv.Itoa(int(nextID)),
			method:		"GET",
			status:		http.StatusOK,
			isError:	false,
			expected:	model.Todo{
				ID:			nextID,
				Title:		"Test TODO",
				Status:		"Done",
				Details:	"test_todo",
				Priority:	"P2",
			},
		},
		{
			name:		"異常系: Item取得: 404",
			url:		"/todo/" + strconv.Itoa(int(nextID)),
			method:		"POST",
			status:		http.StatusNotFound,
			isError:	true,
			expected:	model.Todo{},
		},
		{
			name:		"異常系: Item取得: 400",
			url:		"/todo/error",
			method:		"GET",
			status:		http.StatusBadRequest,
			isError:	true,
			expected:	model.Todo{},
		},
	}

	for _, c := range cases {
		t.Run(caseNameHelper(t, c.name, c.method, c.url), func(t *testing.T) {
			client := &http.Client{}
			req, err := http.NewRequest(c.method, ts.URL + c.url, nil)
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
				if !resData.CreatedAt.Equal(createdAt) {
					t.Fatalf("CreatedAt: want equal to old data, resData = %v", resData.CreatedAt)
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
	now := time.Now()
	nextID := db.GetNextID()
	cases := []struct{
		name 		string
		url			string
		method		string
		status 		int
		isError 	bool
		payload		string
		expected	model.Todo
		need2Delete bool
	}{
		{
			name:		"正常系: 新規追加",
			url:		"/todo",
			method:		"POST",
			status:		http.StatusCreated,
			isError:	false,
			payload:	`{"title": "Test TODO", "status": "Done", "details": "test_todo", "priority": "P0"}`,
			expected:	model.Todo{
				ID:			nextID,
				Title:		"Test TODO",
				Status:		"Done",
				Details:	"test_todo",
				Priority:	"P0",
			},
			need2Delete:	true,
		},
		{
			name:			"異常系: 新規追加: 404",
			url:			"/todo",
			method:			"PUT",
			status:			http.StatusNotFound,
			isError:		true,
			payload:		`{"title": "Test TODO", "status": "Done", "details": "test_todo", "priority": "P0"}`,
			expected:		model.Todo{},
			need2Delete:	false,
		},
		{
			name:			"異常系: 新規追加: 400",
			url:			"/todo",
			method:			"POST",
			status:			http.StatusBadRequest,
			isError:		true,
			payload:		`{"Message": "Bad Request"}`,
			expected:		model.Todo{},
			need2Delete:	false,
		},
	}

	for _, c := range cases {
		t.Run(caseNameHelper(t, c.name, c.method, c.url), func(t *testing.T) {
			client := &http.Client{}
			req, err := http.NewRequest(c.method, ts.URL + c.url, bytes.NewBuffer([]byte(c.payload)))
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
				if !resData.CreatedAt.After(now) {
					t.Fatalf("CreatedAt: DO NOT want %v, resData = %v", now, resData.CreatedAt)
				}
				if !resData.UpdatedAt.After(now) {
					t.Fatalf("UpdatedAt: DO NOT want %v, resData = %v", now, resData.UpdatedAt)
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
	createdAt := res.CreatedAt

	cases := []struct{
		name		string
		url			string
		method		string
		status		int
		isError		bool
		payload		string
		expected	model.Todo
	}{
		{
			name:		"正常系: 更新",
			url:		"/todo/" + strconv.Itoa(int(nextID)),
			method:		"PUT",
			status:		http.StatusOK,
			isError:	false,
			payload:	`{"title": "Changed TODO", "status": "Done", "details": "changed_todo", "priority": "P0"}`,
			expected:	model.Todo{
				ID:			nextID,
				Title:		"Changed TODO",
				Status:		"Done",
				Details:	"changed_todo",
				Priority:	"P0",
			},
		},
		{
			name:		"異常系: 更新: 404",
			url:		"/todo/" + strconv.Itoa(int(nextID)),
			method:		"POST",
			status:		http.StatusNotFound,
			isError:	true,
			payload:	`{"title": "Changed TODO", "status": "Done", "details": "changed_todo", "priority": "P0"}`,
			expected:	model.Todo{},
		},
		{
			name:		"異常系: 新規追加: 400",
			url:		"/todo/" + strconv.Itoa(int(nextID)),
			method:		"PUT",
			status:		http.StatusBadRequest,
			isError:	true,
			payload:	`{"Message": "Bad Request"}`,
			expected:	model.Todo{},
		},
	}

	for _, c := range cases {
		t.Run(caseNameHelper(t, c.name, c.method, c.url), func(t *testing.T) {
			client := &http.Client{}
			req, err := http.NewRequest(c.method, ts.URL + c.url, bytes.NewBuffer([]byte(c.payload)))
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
				if !resData.CreatedAt.Equal(createdAt) {
					t.Fatalf("CreatedAt: want equal to old data, resData = %v", resData.CreatedAt)
				}
				if !resData.UpdatedAt.After(resData.CreatedAt) {
					t.Fatalf("UpdatedAt: DO NOT want equal to %v and %v", resData.CreatedAt, resData.UpdatedAt)
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

func TestUpdateItemState(t *testing.T) {
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
	createdAt := res.CreatedAt

	cases := []struct{
		name		string
		url			string
		method		string
		status		int
		isError		bool
		payload		string
		expected	model.Todo
	}{
		{
			name:		"正常系: 更新",
			url:		"/todo/" + strconv.Itoa(int(nextID)),
			method:		"PUT",
			status:		http.StatusOK,
			isError:	false,
			payload:	`{"title": "Changed TODO", "status": "Done", "details": "changed_todo", "priority": "P0"}`,
			expected:	model.Todo{
				ID:			nextID,
				Title:		"Changed TODO",
				Status:		"Done",
				Details:	"changed_todo",
				Priority:	"P0",
			},
		},
		{
			name:		"異常系: 更新: 404",
			url:		"/todo/" + strconv.Itoa(int(nextID)),
			method:		"POST",
			status:		http.StatusNotFound,
			isError:	true,
			payload:	`{"title": "Changed TODO", "status": "Done", "details": "changed_todo", "priority": "P0"}`,
			expected:	model.Todo{},
		},
		{
			name:		"異常系: 新規追加: 400",
			url:		"/todo/" + strconv.Itoa(int(nextID)),
			method:		"PUT",
			status:		http.StatusBadRequest,
			isError:	true,
			payload:	`{"Message": "Bad Request"}`,
			expected:	model.Todo{},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(caseNameHelper(t, c.name, c.method, c.url), func(t *testing.T) {
			client := &http.Client{}
			req, err := http.NewRequest(c.method, ts.URL + c.url, bytes.NewBuffer([]byte(c.payload)))
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
				if !resData.CreatedAt.Equal(createdAt) {
					t.Fatalf("CreatedAt: want equal to old data, resData = %v", resData.CreatedAt)
				}
				if !resData.UpdatedAt.After(resData.CreatedAt) {
					t.Fatalf("UpdatedAt: DO NOT want equal to %v and %v", resData.CreatedAt, resData.UpdatedAt)
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

func TestDeleteItemState(t *testing.T) {
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
		name		string
		url			string
		method		string
		status		int
		isError		bool
		expected	model.Todo
	}{
		{
			name:		"正常系: 更新",
			url:		"/todo/" + strconv.Itoa(int(nextID)),
			method:		"DELETE",
			status:		http.StatusOK,
			isError:	false,
			expected:	model.Todo{
				ID:			nextID,
				Title:		"Test TODO",
				Status:		"Done",
				Details:	"test_todo",
				Priority:	"P0",
			},
		},
		{
			name:		"異常系: 更新: 400",
			url:		"/todo/error",
			method:		"DELETE",
			status:		http.StatusBadRequest,
			isError:	true,
			expected:	model.Todo{},
		},
	}

	for _, c := range cases {
		t.Run(caseNameHelper(t, c.name, c.method, c.url), func(t *testing.T) {
			client := &http.Client{}
			req, err := http.NewRequest(c.method, ts.URL + c.url, nil)
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
}
