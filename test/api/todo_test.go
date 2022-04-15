package main

import (
	"reflect"

	"encoding/json"
	"bytes"

	"net/http"
	"net/http/httptest"
	"testing"

	todoApi "github.com/Z-me/practice-todo-api/api"
	todoType "github.com/Z-me/practice-todo-api/api/types/todo"
)

func TestTodoApiServer(t *testing.T) {
	// テスト用のサーバーを立てる
	ts := httptest.NewServer(todoApi.Router())
	defer ts.Close()

	testGetList(ts.URL, t)
	testGetItemById(ts.URL, t)
	testPostItem(ts.URL, t)
	testPostItemById(ts.URL, t)
	testUpdateItemStateById(ts.URL, t)

	testDeleteItem(ts.URL, t)

}

func testGetList(url string, t *testing.T) {
	exp := []todoType.Todo{
		{ID: "1",	Title: "最初のTODO",	Status: "Done",	Details: "最初に登録されたTodo",	Priority: "P0"},
		{ID: "2",	Title: "2番目のTODO",	Status: "Backlog",	Details: "2番目に登録されたTodo",	Priority: "P1"},
		{ID: "3",	Title: "3番目TODO",	Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
		{ID: "4",	Title: "4番目TODO",	Status: "Backlog",	Details: "4番目に登録されたTodo",	Priority: "P3"},
		{ID: "5",	Title: "5番目TODO",	Status: "InProgress",	Details: "5番目に登録されたTodo",	Priority: "P1"},
	}
	res, err := http.Get(url + "/todo")
	defer res.Body.Close()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var responseData []todoType.Todo
	err = json.NewDecoder(res.Body).Decode(&responseData)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("[Get Todo List] Expected status code 200, got %v", res.StatusCode)
	}

	if !reflect.DeepEqual(responseData, exp) {
		t.Fatalf("[Get Todo List] responseData = %v, want %v", responseData, exp)
	}
}

func testGetItemById(url string, t *testing.T) {
	exp := todoType.Todo{
		ID: "1",
		Title: "最初のTODO",
		Status: "Done",
		Details: "最初に登録されたTodo",
		Priority: "P0",
	}
	res, err := http.Get(url + "/todo/1")
	defer res.Body.Close()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var responseData todoType.Todo
	err = json.NewDecoder(res.Body).Decode(&responseData)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("[Get Todo Item By ID] Expected status code 200, got %v", res.StatusCode)
	}

	if !reflect.DeepEqual(responseData, exp) {
		t.Fatalf("[Get Todo Item By ID] responseData = %v, want %v", responseData, exp)
	}
}

func testPostItem(url string, t *testing.T) {
	exp := []todoType.Todo{
		{ID: "1",	Title: "最初のTODO",	Status: "Done",	Details: "最初に登録されたTodo",	Priority: "P0"},
		{ID: "2",	Title: "2番目のTODO",	Status: "Backlog",	Details: "2番目に登録されたTodo",	Priority: "P1"},
		{ID: "3",	Title: "3番目TODO",	Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
		{ID: "4",	Title: "4番目TODO",	Status: "Backlog",	Details: "4番目に登録されたTodo",	Priority: "P3"},
		{ID: "5",	Title: "5番目TODO",	Status: "InProgress",	Details: "5番目に登録されたTodo",	Priority: "P1"},
		{ID: "6",	Title: "6番目TODO",	Status: "InProgress",	Details: "6番目に登録されたTodo",	Priority: "P0"},
	}
	payload := todoType.Todo{ID: "6",Title: "6番目TODO",Status: "InProgress",Details: "6番目に登録されたTodo",Priority: "P0"}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	postRes, err := http.Post(url + "/todo", "application/json", bytes.NewBuffer(payloadJson))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer postRes.Body.Close()

	getRes, err := http.Get(url + "/todo")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer getRes.Body.Close()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var responseData []todoType.Todo
	// res, err := json.NewDecoder(getRes.Body).Decode(&responseData)
	json.NewDecoder(getRes.Body).Decode(&responseData)

	if postRes.StatusCode != http.StatusCreated || getRes.StatusCode != http.StatusOK {
		t.Fatalf("[Post New Todo Item] Expected status code 200, got %v and %v", postRes.StatusCode, getRes.StatusCode)
	}

	if !reflect.DeepEqual(responseData, exp) {
		t.Fatalf("[Post New Todo Item] responseData = %v, want %v", responseData, exp)
	}
}

func testPostItemById(url string, t *testing.T) {
	exp := []todoType.Todo{
		{ID: "1",	Title: "最初のTODO",	Status: "Done",	Details: "最初に登録されたTodo",	Priority: "P0"},
		{ID: "2",	Title: "2番目のTODO",	Status: "Backlog",	Details: "2番目に登録されたTodo",	Priority: "P1"},
		{ID: "3",	Title: "3番目TODO",	Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
		{ID: "4",	Title: "4番目TODO",	Status: "Backlog",	Details: "4番目に登録されたTodo",	Priority: "P3"},
		{ID: "5",	Title: "5番目TODO",	Status: "InProgress",	Details: "5番目に登録されたTodo",	Priority: "P1"},
		{ID: "6",Title: "更新された6番目TODO",Status: "Done",Details: "6番目に登録され、その後更新されたTodo",Priority: "P0"},
	}
	payload := todoType.Todo{ID: "6",Title: "更新された6番目TODO",Status: "Done",Details: "6番目に登録され、その後更新されたTodo",Priority: "P0"}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	postRes, err := http.Post(url + "/todo/6", "application/json", bytes.NewBuffer(payloadJson))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer postRes.Body.Close()

	getRes, err := http.Get(url + "/todo")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer getRes.Body.Close()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var responseData []todoType.Todo
	// res, err := json.NewDecoder(getRes.Body).Decode(&responseData)
	json.NewDecoder(getRes.Body).Decode(&responseData)

	if postRes.StatusCode != http.StatusCreated || getRes.StatusCode != http.StatusOK {
		t.Fatalf("[Update Todo Item] Expected status code 200, got %v and %v", postRes.StatusCode, getRes.StatusCode)
	}

	if !reflect.DeepEqual(responseData, exp) {
		t.Fatalf("[Update Todo Item] responseData = %v, want %v", responseData, exp)
	}
}

func testUpdateItemStateById(url string, t *testing.T) {
	exp := []todoType.Todo{
		{ID: "1",	Title: "最初のTODO",	Status: "Done",	Details: "最初に登録されたTodo",	Priority: "P0"},
		{ID: "2",	Title: "2番目のTODO",	Status: "Backlog",	Details: "2番目に登録されたTodo",	Priority: "P1"},
		{ID: "3",	Title: "3番目TODO",	Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
		{ID: "4",	Title: "4番目TODO",	Status: "Backlog",	Details: "4番目に登録されたTodo",	Priority: "P3"},
		{ID: "5",	Title: "5番目TODO",	Status: "InProgress",	Details: "5番目に登録されたTodo",	Priority: "P1"},
		{ID: "6",Title: "更新された6番目TODO",Status: "Backlog",Details: "6番目に登録され、その後更新されたTodo",Priority: "P0"},
	}

	postRes, err := http.Post(url + "/todo/6/status/Backlog", "application/json", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer postRes.Body.Close()

	getRes, err := http.Get(url + "/todo")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer getRes.Body.Close()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var responseData []todoType.Todo
	// res, err := json.NewDecoder(getRes.Body).Decode(&responseData)
	json.NewDecoder(getRes.Body).Decode(&responseData)

	if postRes.StatusCode != http.StatusCreated || getRes.StatusCode != http.StatusOK {
		t.Fatalf("[Update Todo Item] Expected status code 200, got %v and %v", postRes.StatusCode, getRes.StatusCode)
	}

	if !reflect.DeepEqual(responseData, exp) {
		t.Fatalf("[Update Todo Item] responseData = %v, want %v", responseData, exp)
	}
}

func testDeleteItem(url string, t *testing.T) {
	exp := []todoType.Todo{
		{ID: "1",	Title: "最初のTODO",	Status: "Done",	Details: "最初に登録されたTodo",	Priority: "P0"},
		{ID: "2",	Title: "2番目のTODO",	Status: "Backlog",	Details: "2番目に登録されたTodo",	Priority: "P1"},
		{ID: "3",	Title: "3番目TODO",	Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
		{ID: "4",	Title: "4番目TODO",	Status: "Backlog",	Details: "4番目に登録されたTodo",	Priority: "P3"},
		{ID: "5",	Title: "5番目TODO",	Status: "InProgress",	Details: "5番目に登録されたTodo",	Priority: "P1"},
	}
	// deleteRes, err := http.Delete(url + "/todo/3")
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url + "/todo/6", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	deleteRes, _ := client.Do(req)
	defer deleteRes.Body.Close()

	getRes, err := http.Get(url + "/todo")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer getRes.Body.Close()

	if deleteRes.StatusCode != http.StatusOK || getRes.StatusCode != http.StatusOK {
		t.Fatalf("[Post Todo Item] Expected status code 200, got %v and %v", deleteRes.StatusCode, getRes.StatusCode)
	}

	var responseData []todoType.Todo
	json.NewDecoder(getRes.Body).Decode(&responseData)
	if !reflect.DeepEqual(responseData, exp) {
		t.Fatalf("[Post Todo Item] responseData = %v, want %v", responseData, exp)
	}
}
