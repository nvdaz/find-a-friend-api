package llm

import (
	"encoding/json"
	"log"
	"os"

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
	if len(response) > 8 && response[:8] == "```json\n" && response[len(response)-3:] == "```" {
		response = response[8 : len(response)-3]
	}

	return json.Unmarshal([]byte(response), &result)
}

type responseData struct {
	Result string `json:"result"`
	Grade  string `json:"grade"`
	Model  string `json:"model"`
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
		Model:       ModelClaudeSonnet.String(),
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

	_, response, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	responseDict := make(map[string]interface{})
	err = json.Unmarshal(response, &responseDict)
	if err != nil {
		return nil, err
	}

	if _, ok := responseDict["message"]; ok {
		_, response, err = conn.ReadMessage()
		if err != nil {
			return nil, err
		}
	}

	responseData := responseData{}
	err = json.Unmarshal(response, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData.Result, nil
}

func GetResponseJson(result any, model Model, prompt, system string, temperature *float64) error {
	retries := 3

	var err error

	for i := 0; i < retries; i++ {
		log.Println(system)
		response, err := GetResponse(model, prompt, system, temperature)
		log.Println(*response, err)
		if err != nil {
			continue
		}

		err = unmarshalResponse(result, *response)
		if err == nil {
			return nil
		}
	}

	return err
}
