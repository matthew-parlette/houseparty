package houseparty

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestGetEnv(t *testing.T) {
	_ = GetEnv("HOME", "")
}

func TestInitJiraClient(t *testing.T) {
	fmt.Printf("JIRA: %+v\n", JiraClient)
}

func TestInitTodoistClient(t *testing.T) {
	fmt.Printf("Todoist: %+v\n", TodoistClient)
}

func TestInitRocketChatClient(t *testing.T) {
	fmt.Printf("RocketChat: %+v\n", ChatClient)
}

func TestGetTodoistClient(t *testing.T) {
	if err := TodoistClient.Sync(context.Background()); err != nil {
		t.Errorf("Error: %v", err)
	}
	// fmt.Println("Listing projects...")
	// spew.Dump(todoistClient.Store.Projects)
	fmt.Println("Found", len(TodoistClient.Store.Projects), "projects")
}

func TestStartHealthCheck(t *testing.T) {
	StartHealthCheck()

	// Test liveness
	fmt.Println("Liveness:")
	response, err := http.Get("http://localhost:8086/live?full=1")
	if err != nil {
		t.Errorf("Error: %v", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(response.Status), "-", string(data))
	}

	// Test readiness
	fmt.Println("Readiness:")
	response, err = http.Get("http://localhost:8086/ready?full=1")
	if err != nil {
		t.Errorf("Error: %v", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(response.Status), "-", string(data))
	}
}

func TestChatClient(t *testing.T) {
	// t.Skip("Skipping chat client functions")
	StartChatListener()
	timeout := 30 * time.Second
	fmt.Println("Waiting for chat messages for", string(timeout), "seconds...")
	SendChatMessage("house-party", "I'm in test mode, I'll wait for chat messages for 30 seconds...")
	time.Sleep(timeout)
	fmt.Println("Done waiting for chat messages, shutting down...")
	SendChatMessage("house-party", "Done waiting for chat messages, shutting down...")
	return
}
