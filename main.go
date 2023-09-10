package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)
type ModelPermission struct {
	ID                  string `json:"id"`
	Object              string `json:"object"`
	Created             int64  `json:"created"`
	AllowCreateEngine   bool   `json:"allow_create_engine"`
	AllowSampling       bool   `json:"allow_sampling"`
	AllowLogprobs       bool   `json:"allow_logprobs"`
	AllowSearchIndices  bool   `json:"allow_search_indices"`
	AllowView           bool   `json:"allow_view"`
	AllowFineTuning     bool   `json:"allow_fine_tuning"`
	Organization        string `json:"organization"`
	Group               *string `json:"group"` // *string since it can be null
	IsBlocking          bool   `json:"is_blocking"`
}

type Model struct {
	ID          string           `json:"id"`
	Object      string           `json:"object"`
	Created     int64            `json:"created"`
	OwnedBy     string           `json:"owned_by"`
	Permission  []ModelPermission `json:"permission"`
	Root        string           `json:"root"`
	Parent      *string          `json:"parent"` // *string since it can be null
}


type ResponsePayload struct {
	Object string `json:"object"`
	Data []Model `json:"data"`
}

var (
	OPENAPI_URL = "https://api.openai.com/v1/models"
	AUTH_HEADER = "Authorization"
)

func init() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Printf("Error loading .env file")
    }
}
func main() {
	

	client := &http.Client{}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set OPENAI_API_KEY environment variable.")
		return
	}

	req, err := http.NewRequest("GET", OPENAPI_URL, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	
	req.Header.Add(AUTH_HEADER, "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Printf("Status Code: %d\n", resp.StatusCode)

	var responsePayload ResponsePayload
	err = json.Unmarshal([]byte(body), &responsePayload)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}
	
	fmt.Println("Object: %s", responsePayload.Object)
	
		for index, model := range responsePayload.Data {
			fmt.Printf(
				"Model Number: %d, Model ID: %s, Model owned_by: %s, Model Created: %d, Object: %s, Root: %s, Parent: %s\n", 
				index, model.ID, model.OwnedBy, model.Created, model.Object, model.Root, model.Parent)
		}
			
}
