package main

import (
	"reflect"

	"encoding/json"
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
	t.Log("Get List: GET /todo")
}
