package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type TTSProvider interface {
	Synthesize(text string) ([]byte, error)
}

// makeTTSRequest is a generic helper to make HTTP POST requests for TTS APIs
func makeTTSRequest(endpoint string, payload interface{}, headers map[string]string) ([]byte, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// PocketTTSClient for https://github.com/kyutai-labs/pocket-tts (CPU-friendly)
type PocketTTSClient struct {
	Endpoint string
}

func (c *PocketTTSClient) Synthesize(text string) ([]byte, error) {
	endpoint := c.Endpoint
	if endpoint == "" {
		endpoint = os.Getenv("POCKET_TTS_URL")
		if endpoint == "" {
			endpoint = "http://localhost:8000/synthesize" // Default REST endpoint
		}
	}

	// Assuming a standard BentoML/FastAPI payload
	payload := map[string]string{"text": text}
	return makeTTSRequest(endpoint, payload, nil)
}

// KittenTTSClient for https://github.com/KittenML/KittenTTS (ONNX-based)
type KittenTTSClient struct {
	Endpoint string
}

func (c *KittenTTSClient) Synthesize(text string) ([]byte, error) {
	endpoint := c.Endpoint
	if endpoint == "" {
		endpoint = os.Getenv("KITTEN_TTS_URL")
		if endpoint == "" {
			endpoint = "http://localhost:8001/synthesize" // Default REST endpoint
		}
	}

	payload := map[string]string{"text": text}
	return makeTTSRequest(endpoint, payload, nil)
}

// FishAudioClient for https://huggingface.co/fishaudio/s2-pro
type FishAudioClient struct {
	Endpoint string
	APIKey   string
}

func (c *FishAudioClient) Synthesize(text string) ([]byte, error) {
	endpoint := c.Endpoint
	if endpoint == "" {
		endpoint = os.Getenv("FISH_AUDIO_URL")
		if endpoint == "" {
			endpoint = "https://api-inference.huggingface.co/models/fishaudio/s2-pro"
		}
	}

	apiKey := c.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("HF_API_KEY")
	}

	headers := make(map[string]string)
	if apiKey != "" {
		headers["Authorization"] = "Bearer " + apiKey
	}

	// HuggingFace Inference API standard payload
	payload := map[string]interface{}{
		"inputs": text,
	}
	return makeTTSRequest(endpoint, payload, headers)
}

// Dia2Client for https://huggingface.co/nari-labs/Dia2-2B (Dialogue-focused)
type Dia2Client struct {
	Endpoint string
	APIKey   string
}

func (c *Dia2Client) Synthesize(text string) ([]byte, error) {
	endpoint := c.Endpoint
	if endpoint == "" {
		endpoint = os.Getenv("DIA2_URL")
		if endpoint == "" {
			endpoint = "https://api-inference.huggingface.co/models/nari-labs/Dia2-2B"
		}
	}

	apiKey := c.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("HF_API_KEY")
	}

	headers := make(map[string]string)
	if apiKey != "" {
		headers["Authorization"] = "Bearer " + apiKey
	}

	// HuggingFace Inference API standard payload
	payload := map[string]interface{}{
		"inputs": text,
	}
	return makeTTSRequest(endpoint, payload, headers)
}

func NewTTSProvider(model string) TTSProvider {
        switch model {
        case "pocket-tts":
                return &PocketTTSClient{}
        case "kitten-tts":
                return &KittenTTSClient{}
        case "fishaudio-s2-pro":
                return &FishAudioClient{}
        case "dia2-2b":
                return &Dia2Client{}
        default:
                return &PocketTTSClient{}
        }
}