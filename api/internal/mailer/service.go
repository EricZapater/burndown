package mailer

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
)

type Service interface {
	SendWelcomeEmail(userEmail string, userName string) error
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (s *service) SendWelcomeEmail(userEmail string, userName string) error {
	mailerURL := os.Getenv("MAILER_SERVICE_URL")
	if mailerURL == "" {
		mailerURL = "https://mailer-service.preprod.zenith.ovh/api/send" // Fallback
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	_ = writer.WriteField("mail_to", userEmail)
	_ = writer.WriteField("subject", "Benvingut a Burndown! 🔥")
	_ = writer.WriteField("template_name", "burndown")
	_ = writer.WriteField("body", userName)
	_ = writer.WriteField("footer", "L'equip de Burndown")

	err := writer.Close()
	if err != nil {
		return fmt.Errorf("error tancant el multipart writer: %w", err)
	}
	req, err := http.NewRequest("POST", mailerURL, payload)
	if err != nil {
		return fmt.Errorf("error creant la petició HTTP: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviant la petició al mailer: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return fmt.Errorf("el servei de mail ha retornat un error: %d", res.StatusCode)
	}

	return nil
}