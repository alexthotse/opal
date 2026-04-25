package adapters

type TTSProvider interface {
        Synthesize(text string) ([]byte, error)
}

// PocketTTSClient for https://github.com/kyutai-labs/pocket-tts (CPU-friendly)
type PocketTTSClient struct{}
func (c *PocketTTSClient) Synthesize(text string) ([]byte, error) {
        return []byte("pocket-tts-audio-data"), nil
}

// KittenTTSClient for https://github.com/KittenML/KittenTTS (ONNX-based)
type KittenTTSClient struct{}
func (c *KittenTTSClient) Synthesize(text string) ([]byte, error) {
        return []byte("kitten-tts-audio-data"), nil
}

// FishAudioClient for https://huggingface.co/fishaudio/s2-pro
type FishAudioClient struct{}
func (c *FishAudioClient) Synthesize(text string) ([]byte, error) {
        return []byte("fishaudio-audio-data"), nil
}

// Dia2Client for https://huggingface.co/nari-labs/Dia2-2B (Dialogue-focused)
type Dia2Client struct{}
func (c *Dia2Client) Synthesize(text string) ([]byte, error) {
        return []byte("dia2-audio-data"), nil
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