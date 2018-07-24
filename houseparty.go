package houseparty

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	chat "github.com/RocketChat/Rocket.Chat.Go.SDK/realtime"
	jira "github.com/andygrunwald/go-jira"
	"github.com/sachaos/todoist/lib"
)

var (
	ConfigPath  string
	SecretsPath string
)

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func Config(item string) string {
	contents, err := ioutil.ReadFile(path.Join(ConfigPath, item))
	if err != nil {
		log.Fatal(err)
		return ""
	}
	result := strings.TrimSpace(string(contents))
	return result
}

func Secret(item string) string {
	contents, err := ioutil.ReadFile(path.Join(SecretsPath, item))
	if err != nil {
		log.Fatal(err)
		return ""
	}
	result := strings.TrimSpace(string(contents))
	return result
}

func GetJiraClient() (*jira.Client, error) {
	fmt.Println("Initializing JIRA...")
	tp := jira.BasicAuthTransport{
		Username: Config("jira-username"),
		Password: Secret("jira-password"),
	}
	jiraClient, err := jira.NewClient(tp.Client(), Config("jira-url"))
	if err != nil {
		return nil, err
	}
	return jiraClient, nil
}

func GetTodoistClient() (*todoist.Client, error) {
	fmt.Println("Initializing todoist...")
	config := &todoist.Config{
		AccessToken: Secret("todoist-token"),
		DebugMode:   false,
		// Color:       false,
	}
	todoistClient := todoist.NewClient(config)
	var store todoist.Store
	todoistClient.Store = &store
	return todoistClient, nil
}

func GetRocketChatClient() (*chat.Client, error) {
	rocketchatUrlString := Config("rocketchat-url")
	rocketchatUrl, err := url.Parse(rocketchatUrlString)
	if err != nil {
		return nil, err
	}
	chatClient, err := chat.NewClient(rocketchatUrl, false)
	if err != nil {
		return nil, err
	}
	_, err = chatClient.Login(&models.UserCredentials{
		Email:    Config("rocketchat-email"),
		Password: Secret("rocketchat-password")})
	if err != nil {
		return nil, err
	}
	return chatClient, nil
}