package advisor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
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
	JobID      uuid.UUID `json:"job_id"`
}

type Service struct {
	client *genai.Client
	model  *genai.GenerativeModel
	repo   Repository
}

const instructions = `
		ROLE: Expert Sports Nutritionist.
        TASK: Analyze food based on TWO inputs of equal importance: an IMAGE and an optional text DESCRIPTION provided by the user.
        
        COHERENCE CHECK (CRITICAL):
        - If both image and description are provided, cross-reference them.
        - If they match, base your analysis on both.
        - If they contradict each other (e.g. image shows pasta but description says chicken), clearly state the discrepancy in the 'reasoning' field and base your nutritional estimates primarily on what you SEE in the image, but acknowledge the description.
        - If no description is provided, analyze only the image.
        
        OUTPUT FORMAT: JSON only.
        fields: name (in Catalan), calories, protein_g, carbs_g, fat_g, confidence_level, reasoning.
        
        RULES:
        1. Account for hidden oils/sauces.
        2. If unsure, slightly overestimate calories.
        3. Identify the food in the 'name' field clearly.
        4. In 'reasoning', always explain your analysis and whether image and description are coherent.
    `

func NewService(apiKey string, repo Repository) (*Service, error) {
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

	return &Service{client: client, model: model, repo: repo}, nil
}

func (s *Service) AnalyzeImage(ctx context.Context, jobID uuid.UUID, requestMeta json.RawMessage, imgData []byte, imgFormat, mealDescription string) (uuid.UUID, error) {
	// 1. Create Job synchronously
	// 1a. Save Image to Disk
	filename := fmt.Sprintf("%s.%s", jobID.String(), imgFormat)
	imagePath := filepath.Join("/app/images", filename)
	
	// Ensure directory exists (redundant with Dockerfile but good practice)
	if err := os.MkdirAll("/app/images", 0755); err != nil {
		return uuid.Nil, fmt.Errorf("failed to create images directory: %w", err)
	}

	if err := os.WriteFile(imagePath, imgData, 0644); err != nil {
		return uuid.Nil, fmt.Errorf("failed to save image to disk: %w", err)
	}

	job := &AnalysisJob{
		ID:          jobID,
		Status:      "pending",
		ImagePath:   imagePath,
		Description: mealDescription,
		RequestMeta: requestMeta,
		RequestedAt: time.Now(),
	}
	if err := s.repo.CreateAnalysisJob(ctx, job); err != nil {
		return uuid.Nil, fmt.Errorf("failed to create job: %w", err)
	}

	// 2. Launch Analysis in Goroutine
	go func() {
		// Use a new background context for the async operation because 'ctx' might be cancelled when the HTTP request finishes
		asyncCtx := context.Background() 
		
		// 2.1 Prepare Prompt
		var promptText string
		if mealDescription != "" {
			promptText = fmt.Sprintf("Ets un nutricionista expert. Se t'envien DUES fonts d'informació amb el MATEIX pes:\n\n1. IMATGE: La foto del plat (adjunta).\n2. DESCRIPCIÓ DE L'USUARI: \"%s\"\n\nAnalitza els macronutrients amb precisió. Comprova si la imatge i la descripció són coherents. Si hi ha discrepàncies, indica-les al camp 'reasoning' i basa't principalment en el que veus a la imatge.", mealDescription)
		} else {
			promptText = "Ets un nutricionista expert. Analitza la imatge del plat i estima els macronutrients amb precisió. Sigues realista amb les racions. No s'ha proporcionat cap descripció addicional."
		}
		prompt := genai.Text(promptText)
		imgPart := genai.ImageData(imgFormat, imgData)

		aiCtx, cancel := context.WithTimeout(asyncCtx, 120*time.Second) // Increased timeout for safety
		defer cancel()

		resp, err := s.model.GenerateContent(aiCtx, prompt, imgPart)
		if err != nil {
			job.Status = "failed"
			msg := err.Error()
			job.ErrorMessage = &msg
			now := time.Now()
			job.ResponsedAt = &now
			_ = s.repo.UpdateAnalysisJob(asyncCtx, job)
			return // Stop here
		}

		if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
			job.Status = "failed"
			msg := "empty response from gemini"
			job.ErrorMessage = &msg
			now := time.Now()
			job.ResponsedAt = &now
			_ = s.repo.UpdateAnalysisJob(asyncCtx, job)
			return
		}

		// Extreure el JSON text
		jsonText, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
		if !ok {
			job.Status = "failed"
			msg := "unexpected response format"
			job.ErrorMessage = &msg
			now := time.Now()
			job.ResponsedAt = &now
			_ = s.repo.UpdateAnalysisJob(asyncCtx, job)
			return
		}

		// Parsejar a la nostra struct Go
		var result AnalysisResult
		if err := json.Unmarshal([]byte(jsonText), &result); err != nil {
			job.Status = "failed"
			msg := fmt.Sprintf("failed to parse json: %v", err)
			job.ErrorMessage = &msg
			now := time.Now()
			job.ResponsedAt = &now
			_ = s.repo.UpdateAnalysisJob(asyncCtx, job)
			return
		}

		// 3. Update Job with Success
		job.Status = "completed"
		job.Response = json.RawMessage(jsonText)
		now := time.Now()
		job.ResponsedAt = &now
		if err := s.repo.UpdateAnalysisJob(asyncCtx, job); err != nil {
			fmt.Printf("failed to update job %s: %v\n", job.ID, err)
		}
	}()

	return job.ID, nil
}

func (s *Service) GetAnalysisJob(ctx context.Context, id uuid.UUID) (*AnalysisJob, error) {
	return s.repo.GetAnalysisJob(ctx, id)
}

// Close tanca el client quan apaguem el servidor
func (s *Service) Close() {
	s.client.Close()
}