package jcmllm

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const (
	modModel = "tbot-mod"
	llmModel = "mistral"
)

// "context"

// "regexp"

// "github.com/tmc/langchaingo/llms"
// "github.com/tmc/langchaingo/llms/ollama"

// msgRegex := regexp.MustCompile("PRIV1MSG|JO1IN|PA1RT")

// llm, err := ollama.New(ollama.WithModel(llmModel))
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
// response := ""
// dataSet := "{}"
//
//	streamingFn := llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
//		response += string(chunk)
//		if len(chunk) == 0 {
//			dataSet = response
//			response = ""
//		}
//		return nil
//	})

type BasicOutput struct {
	Output string `json:"output"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Message  Message `json:"message"`
	Response string  `json:"response"`
	Context  []int   `json:"context"`
}

type LLMOpts struct {
	NumCtx     int      `json:"num_ctx"`
	NumPredict int      `json:"num_predict"`
	Stop       []string `json:"stop"`
}

type LLMRequest struct {
	Messages []Message `json:"messages"`
	// Prompt   string    `json:"prompt"`
	Format string `json:"format"`
	Model  string `json:"model"`
	Stream bool   `json:"stream"`
	// Options   LLMOpts   `json:"options"`
	// KeepAlive string    `json:"keep_alive"`
}

func getCompletion(prompt string) {

	// 	prompt += " The current JSON object is: " + dataSet
	// 	prompt += " And the next message to process is: " + line
	//
	// 	llmCtx := context.Background()
	//
	// 	completion, err := llms.GenerateFromSinglePrompt(llmCtx, llm, prompt, llms.WithTemperature(0.8), streamingFn)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	//
	// 	_ = completion
}

func GenerateCompletion(model string, prompt string, useJson bool) string {
	// println("getting completion using " + model + " with prompt " + prompt)
	data := LLMRequest{
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		Model:  model,
		Stream: false,
		// KeepAlive: "0",
		// Options: LLMOpts{
		// 	Stop: []string{"}"},
		// },
	}

	if useJson == true {
		data.Format = "json"
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "http://localhost:11434/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// println("THE BODY:" + string(body))

	var chatRequest ChatRequest

	if err := json.Unmarshal(body, &chatRequest); err != nil {
		log.Fatal(err)
	}

	var printJson bytes.Buffer

	json.Indent(&printJson, body, "", "\t")

	// println("got a response" + string(printJson.Bytes()))

	return chatRequest.Message.Content
}

func getModeration(message string) string {
	var moderationEvent BasicOutput

	modEvt := GenerateCompletion(modModel, message, true)

	if err := json.Unmarshal([]byte(modEvt), &moderationEvent); err != nil {
		log.Fatal(err)
	}

	return moderationEvent.Output
}
