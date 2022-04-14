package main

import (
	// "fmt"
	"reflect"
	"io/ioutil"
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

	// url := fmt.Sprintf("%s/", ts.URL)
	// リクエストを送れるか?
	// resp, err := http.Get(url + "todo")
	// fmt.Sprintf("%s/", ts.URL)
	// resp, err := http.Get("/todo")
	// if err != nil {
	// 	t.Fatalf("Expected no error, got %v", err)
	// }
	// defer resp.Body.Close()
	// Statusコードは200か?
	// if resp.StatusCode != http.StatusOK {
	// 	t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	// }
	// responseData, _ := ioutil.ReadAll(resp.Body)
	// if string(responseData) != mockUserResp {
	// 	t.Fatalf("Expected hello world message, got %v", responseData)
	// }
}

func testGetList(url string, t *testing.T) {
	exp := []todoType.Todo{
		{ID: "1",	Title: "最初のTODO",	Status: "Done",	Details: "最初に登録されたTodo",	Priority: "P0"},
		{ID: "2",	Title: "2番目のTODO",	Status: "Backlog",	Details: "2番目に登録されたTodo",	Priority: "P1"},
		{ID: "3",	Title: "3番目TODO",	Status: "InProgress",	Details: "3番目に登録されたTodo",	Priority: "P2"},
		{ID: "4",	Title: "4番目TODO",	Status: "Backlog",	Details: "4番目に登録されたTodo",	Priority: "P3"},
		{ID: "5",	Title: "5番目TODO",	Status: "InProgress",	Details: "5番目に登録されたTodo",	Priority: "P1"},
	}
	// call api
	res, err := http.Get(url + "/todo")
	defer res.Body.Close()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	responseData, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	var body []todoType.Todo
	if err := json.Unmarshal(responseData, &body); err != nil {
			return nil, err
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("[Get Todo List] Expected status code 200, got %v", res.StatusCode)
	}

	if !reflect.DeepEqual(responseData, exp) {
		t.Fatalf("responseData = %v, want %v", responseData, exp)
	}
}
