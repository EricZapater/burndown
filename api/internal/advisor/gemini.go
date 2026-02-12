package advisor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AnalysisResult struct {
	Reasoning string `json:"reasoning"`
	Name     string  `json:"name"`
	Calories int     `json:"calories"`
	Protein  float64 `json:"protein_g"`
	Carbs    float64 `json:"carbs_g"`
	Fat      float64 `json:"fat_g"`	
	Confidence string `json:"confidence_level"` 
}

type Service struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

const instructions = `
		ROLE: Expert Sports Nutritionist.
        TASK: Analyze food images for macronutrients.
        
        OUTPUT FORMAT: JSON only.
        fields: name (in Catalan), calories, protein_g, carbs_g, fat_g, confidence_level.
        
        RULES:
        1. Account for hidden oils/sauces.
        2. If unsure, slightly overestimate calories.
        3. Identify the food in the 'name' field clearly.
    `

func NewService(apiKey string) (*Service, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("gemini-3-flash-preview")
	model.SystemInstruction = genai.NewUserContent(genai.Text(instructions))
	model.SetTemperature(0.1)
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"reasoning":        {Type: genai.TypeString},
			"name":             {Type: genai.TypeString},
			"calories":         {Type: genai.TypeInteger},
			"protein_g":        {Type: genai.TypeNumber},
			"carbs_g":          {Type: genai.TypeNumber},
			"fat_g":            {Type: genai.TypeNumber},
			"confidence_level": {Type: genai.TypeString, Enum: []string{"high", "medium", "low"}},
		},
		Required: []string{"reasoning", "name", "calories", "protein_g", "carbs_g", "fat_g"},
	}

	return &Service{client: client, model: model}, nil
}

func (s *Service) AnalyzeImage(ctx context.Context, imgData []byte, imgFormat, mealDescription string) (*AnalysisResult, error) {
	prompt := genai.Text(fmt.Sprintf("Ets un nutricionista expert. Analitza aquesta imatge i estima els macronutrients amb precisió. Sigues realista amb les racions. L'usuari ha dit del plat: %s", mealDescription))
	imgPart := genai.ImageData(imgFormat, imgData)	

	aiCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	
	resp, err := s.model.GenerateContent(aiCtx, prompt, imgPart)
	if err != nil {
		return nil, fmt.Errorf("error calling gemini: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from gemini")
	}

	// Extreure el JSON text
	jsonText, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	// Parsejar a la nostra struct Go
	var result AnalysisResult
	if err := json.Unmarshal([]byte(jsonText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	return &result, nil
}

// Close tanca el client quan apaguem el servidor
func (s *Service) Close() {
	s.client.Close()
}