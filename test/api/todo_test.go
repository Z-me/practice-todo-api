package main

import (
	"reflect"

	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Z-me/practice-todo-api/api"
	"github.com/Z-me/practice-todo-api/api/model"
)

func TestGetTodoList(t *testing.T) {
	// Note: Start test Server
	ts := httptest.NewServer(api.Router())
	defer ts.Close()

	// Note: expected Values
	exp := []model.Todo{
		{ID: "1",	Title: "最初のTODO", Status: "Done", Details: "最初に登録されたTodo", Priority: "P0"},
		{ID: "2",	Title: "2番目のTODO",	Status: "Backlog",	Details: "2番目に登録されたTodo",	Priority: "P1"},
		{ID: "3",	Title: "3番目TODO",	Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
		{ID: "4",	Title: "4番目TODO",	Status: "Backlog",	Details: "4番目に登録されたTodo",	Priority: "P3"},
		{ID: "5",	Title: "5番目TODO",	Status: "InProgress",	Details: "5番目に登録されたTodo",	Priority: "P1"},
	}

	// Note: Call API
	res, err := http.Get(ts.URL + "/todo")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer res.Body.Close()

	var responseData []model.Todo
	json.NewDecoder(res.Body).Decode(&responseData)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v", res.StatusCode)
	}

	if !reflect.DeepEqual(responseData, exp) {
		t.Fatalf("responseData = %v, want %v", responseData, exp)
	}
}

func TestGetTodoItemById(t *testing.T) {
	// Note: Start test Server
	ts := httptest.NewServer(api.Router())
	defer ts.Close()

	// Note: expected Values
	exp := model.Todo{
		ID: "1",
		Title: "最初のTODO",
		Status: "Done",
		Details: "最初に登録されたTodo",
		Priority: "P0",
	}

	// Note: Call API
	res, err := http.Get(ts.URL + "/todo/1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer res.Body.Close()

	var responseData model.Todo
	json.NewDecoder(res.Body).Decode(&responseData)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("[Get Todo Item By ID] Expected status code 200, got %v", res.StatusCode)
	}

	if !reflect.DeepEqual(responseData, exp) {
		t.Fatalf("[Get Todo Item By ID] responseData = %v, want %v", responseData, exp)
	}
}

func TestAddItem(t *testing.T) {
	// Note: Start test Server
	ts := httptest.NewServer(api.Router())
	defer ts.Close()

	// Note: expected Values
	exp := []model.Todo{
		{ID: "1",	Title: "最初のTODO",	Status: "Done",	Details: "最初に登録されたTodo",	Priority: "P0"},
		{ID: "2",	Title: "2番目のTODO",	Status: "Backlog",	Details: "2番目に登録されたTodo",	Priority: "P1"},
		{ID: "3",	Title: "3番目TODO",	Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
		{ID: "4",	Title: "4番目TODO",	Status: "Backlog",	Details: "4番目に登録されたTodo",	Priority: "P3"},
		{ID: "5",	Title: "5番目TODO",	Status: "InProgress",	Details: "5番目に登録されたTodo",	Priority: "P1"},
		{ID: "6",	Title: "6番目TODO",	Status: "InProgress",	Details: "6番目に登録されたTodo",	Priority: "P0"},
	}
	payload := model.Todo{ID: "6",Title: "6番目TODO",Status: "InProgress",Details: "6番目に登録されたTodo",Priority: "P0"}
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

	var responseData []model.Todo
	json.NewDecoder(getRes.Body).Decode(&responseData)

	if postRes.StatusCode != http.StatusCreated || getRes.StatusCode != http.StatusOK {
		t.Fatalf("[Post New Todo Item] Expected status code 200, got %v and %v", postRes.StatusCode, getRes.StatusCode)
	}

	if !reflect.DeepEqual(responseData, exp) {
		t.Fatalf("[Post New Todo Item] responseData = %v, want %v", responseData, exp)
	}
}

func TestUpdateItem(t *testing.T) {
	// Note: Start test Server
	ts := httptest.NewServer(api.Router())
	defer ts.Close()

	// Note: expected Values
	exp := []model.Todo{
		{ID: "1",	Title: "最初のTODO",	Status: "Done",	Details: "最初に登録されたTodo",	Priority: "P0"},
		{ID: "2",	Title: "2番目のTODO",	Status: "Backlog",	Details: "2番目に登録されたTodo",	Priority: "P1"},
		{ID: "3",	Title: "3番目TODO",	Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
		{ID: "4",	Title: "4番目TODO",	Status: "Backlog",	Details: "4番目に登録されたTodo",	Priority: "P3"},
		{ID: "5",Title: "更新された5番目TODO",Status: "Done",Details: "5番目に登録され、その後更新されたTodo",Priority: "P1"},
	}
	payload := model.Todo{ID: "5",Title: "更新された5番目TODO",Status: "Done",Details: "5番目に登録され、その後更新されたTodo",Priority: "P1"}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Note: Call POST API
	postRes, err := http.Post(ts.URL + "/todo/5", "application/json", bytes.NewBuffer(payloadJson))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer postRes.Body.Close()

	if postRes.StatusCode != http.StatusCreated {
		t.Fatalf("[Update Todo Item] Expected status code 201, got %v", postRes.StatusCode)
	}

	// Note: Call GET API
	getRes, err := http.Get(ts.URL + "/todo")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer getRes.Body.Close()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var responseData []model.Todo
	json.NewDecoder(getRes.Body).Decode(&responseData)


	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("[Update Todo Item] Expected status code 200, got %v", getRes.StatusCode)
	}

	if !reflect.DeepEqual(responseData, exp) {
		t.Fatalf("[Update Todo Item] responseData = %v, want %v", responseData, exp)
	}
}

func TestUpdateStateById(t *testing.T) {
	// Note: Start test Server
	ts := httptest.NewServer(api.Router())
	defer ts.Close()

	exp := []model.Todo{
		{ID: "1",	Title: "最初のTODO",	Status: "Done",	Details: "最初に登録されたTodo",	Priority: "P0"},
		{ID: "2",	Title: "2番目のTODO",	Status: "Backlog",	Details: "2番目に登録されたTodo",	Priority: "P1"},
		{ID: "3",	Title: "3番目TODO",	Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
		{ID: "4",	Title: "4番目TODO",	Status: "Backlog",	Details: "4番目に登録されたTodo",	Priority: "P3"},
		{ID: "5",	Title: "5番目TODO",	Status: "Backlog",	Details: "5番目に登録されたTodo",	Priority: "P1"},
	}

	postRes, err := http.Post(ts.URL + "/todo/5/status/Backlog", "application/json", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer postRes.Body.Close()

	// Note: Call POST API
	getRes, err := http.Get(ts.URL + "/todo")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer getRes.Body.Close()

	var responseData []model.Todo
	json.NewDecoder(getRes.Body).Decode(&responseData)

	if postRes.StatusCode != http.StatusCreated || getRes.StatusCode != http.StatusOK {
		t.Fatalf("[Update Todo Item] Expected status code 200, got %v and %v", postRes.StatusCode, getRes.StatusCode)
	}

	if !reflect.DeepEqual(responseData, exp) {
		t.Fatalf("[Update Todo Item] responseData = %v, want %v", responseData, exp)
	}
}

func TestDeleteItem(t *testing.T) {
	// Note: Start test Server
	ts := httptest.NewServer(api.Router())
	defer ts.Close()

	// Note: expected Values
	exp := []model.Todo{
		{ID: "1",	Title: "最初のTODO",	Status: "Done",	Details: "最初に登録されたTodo",	Priority: "P0"},
		{ID: "2",	Title: "2番目のTODO",	Status: "Backlog",	Details: "2番目に登録されたTodo",	Priority: "P1"},
		{ID: "3",	Title: "3番目TODO",	Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
		{ID: "4",	Title: "4番目TODO",	Status: "Backlog",	Details: "4番目に登録されたTodo",	Priority: "P3"},
	}

	// Note: Call DELETE API
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", ts.URL + "/todo/5", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	deleteRes, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer deleteRes.Body.Close()

	// Note: Call GET API
	getRes, err := http.Get(ts.URL + "/todo")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer getRes.Body.Close()

	if deleteRes.StatusCode != http.StatusOK || getRes.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v and %v", deleteRes.StatusCode, getRes.StatusCode)
	}

	var responseData []model.Todo
	json.NewDecoder(getRes.Body).Decode(&responseData)
	if !reflect.DeepEqual(responseData, exp) {
		t.Fatalf("responseData = %v, want %v", responseData, exp)
	}
}
