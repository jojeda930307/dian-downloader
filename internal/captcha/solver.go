package captcha

import (
	"dian-downloader/internal/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/imroc/req/v3"
)

func SolveTurnstile(targetURL string) (string, error) {
	apiKey := os.Getenv("ANTICAPTCHA_KEY")

	if apiKey == "" {
		return "", fmt.Errorf("ANTICAPTCHA_KEY environment variable not set")
	}

	client := req.C()

	payload := map[string]interface{}{
		"clientKey": apiKey,
		"task": map[string]interface{}{
			"type":       "TurnstileTaskProxyless",
			"websiteURL": targetURL,
			"websiteKey": models.SiteKey,
		},
	}

	var createRes models.CreateTaskResponse
	_, err := client.R().SetBody(payload).SetSuccessResult(&createRes).Post("https://api.anti-captcha.com/createTask")
	if err != nil || createRes.ErrorID != 0 {
		errObj := fmt.Errorf("error creating captcha task: %s (ID: %d)", createRes.Code, createRes.ErrorID)
		log.Printf("[CAPTCHA ERROR] %v", errObj)
		return "", errObj
	}

	for i := 0; i < 30; i++ {
		time.Sleep(3 * time.Second)
		var taskRes models.GetTaskResponse
		_, err := client.R().
			SetBody(map[string]interface{}{"clientKey": apiKey, "taskId": createRes.TaskID}).
			SetSuccessResult(&taskRes).
			Post("https://api.anti-captcha.com/getTaskResult")

		if err != nil {
			continue
		}
		if taskRes.Status == "ready" {
			return taskRes.Solution.Token, nil
		}
	}

	errTimeout := fmt.Errorf("timeout reached while waiting for captcha solution")
	log.Printf("[CAPTCHA ERROR] %v", errTimeout)
	return "", errTimeout
}
