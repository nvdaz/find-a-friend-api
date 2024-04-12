package llm

import (
	"encoding/json"
	"os"

	"github.com/gorilla/websocket"
)

func unmarshalResponse(result any, response string) error {
	if len(response) > 8 && response[:8] == "```json\n" && response[len(response)-3:] == "```" {
		response = response[8 : len(response)-3]
	}

	return json.Unmarshal([]byte(response), &result)
}

type PromptData struct {
	Action      string   `json:"action"`
	Model       string   `json:"model"`
	Prompt      string   `json:"prompt"`
	System      string   `json:"system"`
	Temperature *float64 `json:"temperature"`
}

type ResponseData struct {
	Result string `json:"result"`
	Grade  string `json:"grade"`
	Model  string `json:"model"`
}

func GetResponse(model, prompt, system string, temperature *float64) (*string, error) {
	uri := os.Getenv("LLM_WEBSOCKET_URI")

	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	requestData := PromptData{
		Action:      "runModel",
		Model:       model,
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

	responseData := ResponseData{}
	err = json.Unmarshal(response, &responseData)
	if err != nil {
		return nil, err
	}

	return &responseData.Result, nil
}

func GetResponseJson(result any, model, prompt, system string, temperature *float64) error {
	retries := 2

	var err error

	for i := 0; i < retries; i++ {
		response, err := GetResponse(model, prompt, system, temperature)
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
