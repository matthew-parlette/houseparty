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

func TestGetTodoistClient(t *testing.T) {
	ConfigPath = GetEnv("CONFIG_PATH", "config")
	SecretsPath = GetEnv("SECRETS_PATH", "secrets")
	todoistClient, err := GetTodoistClient()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if err := todoistClient.Sync(context.Background()); err != nil {
		t.Errorf("Error: %v", err)
	}
	// fmt.Println("Listing projects...")
	// spew.Dump(todoistClient.Store.Projects)
	fmt.Println("Found", len(todoistClient.Store.Projects), "projects")
}

func TestStartHealthCheck(t *testing.T) {
	ConfigPath = GetEnv("CONFIG_PATH", "config")
	SecretsPath = GetEnv("SECRETS_PATH", "secrets")
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
	ConfigPath = GetEnv("CONFIG_PATH", "config")
	SecretsPath = GetEnv("SECRETS_PATH", "secrets")
	chatClient, err := GetRocketChatClient()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	StartChatListener(chatClient)
	timeout := 30 * time.Second
	fmt.Println("Waiting for chat messages for", string(timeout), "seconds...")
	SendChatMessage(chatClient, "house-party", "I'm in test mode, I'll wait for chat messages for 30 seconds...")
	time.Sleep(timeout)
	fmt.Println("Done waiting for chat messages, shutting down...")
	SendChatMessage(chatClient, "house-party", "Done waiting for chat messages, shutting down...")
	return
}
