package llm

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	ModelGpt4         Model = "gpt4-new"
	ModelGpt3p5       Model = "gpt3-5"
	ModelClaudeHaiku  Model = "anthropic.claude-3-haiku-20240307-v1:0"
	ModelClaudeSonnet Model = "anthropic.claude-3-sonnet-20240229-v1:0"
)

type Model string

func (model Model) String() string {
	return string(model)
}

func unmarshalResponse(result any, response string) error {
	start := strings.Index(response, "{")
	if start == -1 {
		return fmt.Errorf("no opening brace found in response")
	}

	end := strings.LastIndex(response, "}")
	if end == -1 {
		return fmt.Errorf("no closing brace found in response")
	}

	jsonString := response[start : end+1]
	return json.Unmarshal([]byte(jsonString), &result)
}

type promptData struct {
	Action      string   `json:"action"`
	Model       string   `json:"model"`
	Prompt      string   `json:"prompt"`
	System      string   `json:"system"`
	Temperature *float64 `json:"temperature"`
}

func GetResponse(model Model, prompt, system string, temperature *float64) (*string, error) {
	uri := os.Getenv("LLM_WEBSOCKET_URI")

	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	requestData := promptData{
		Action:      "runModel",
		Model:       model.String(),
		Prompt:      prompt,
		System:      system,
		Temperature: temperature,
	}
	message, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return nil, err
	}

	for {
		_, message, err = conn.ReadMessage()
		if err != nil {
			return nil, fmt.Errorf("read message: %w", err)
		}

		fmt.Println(string(message))
		var response struct {
			Result string `json:"result"`
			Error  string `json:"error"`
		}
		if err := json.Unmarshal(message, &response); err != nil {
			return nil, fmt.Errorf("unmarshal json: %w", err)
		}

		if response.Result != "" {
			return &response.Result, nil
		}
		if response.Error != "" {
			return nil, fmt.Errorf("response error: %s", response.Error)
		}
	}
}

func GetResponseJson(result any, model Model, prompt, system string, temperature *float64) error {
	retries := 3

	var err error

	for i := 0; i < retries; i++ {
		log.Println(system)
		response, err := GetResponse(model, prompt, system, temperature)
		if err != nil {
			continue
		}
		log.Println(*response, err)

		err = unmarshalResponse(result, *response)
		if err == nil {
			fmt.Println(result)
			return nil
		}
	}

	return err
}
