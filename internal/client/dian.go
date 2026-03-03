package client

import (
	"dian-downloader/internal/captcha"
	"dian-downloader/internal/htmlutil"
	"dian-downloader/internal/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/imroc/req/v3"
)

type DianClient struct {
	http *req.Client
}

func NewDianClient() *DianClient {
	client := req.C().
		SetTimeout(60*time.Second).
		SetCommonHeader("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:147.0) Gecko/20100101 Firefox/147.0").
		SetCommonHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8").
		SetCommonHeader("Accept-Language", "es-ES,es;q=0.9,en-US;q=0.8,en;q=0.7")

	return &DianClient{http: client}
}

func (c *DianClient) DownloadPDF(documentKey string, outputPath string) error {
	resp, err := c.http.R().Get(models.SearchDocument)
	if err != nil {
		return err
	}

	body1, _ := io.ReadAll(resp.Body)
	inputs1, err := htmlutil.ExtractInputs(body1, "search-document-form")
	if err != nil {
		return err
	}

	token, err := captcha.SolveTurnstile(models.SearchDocument)
	if err != nil {
		return err
	}

	payload := map[string]string{
		"__RequestVerificationToken": inputs1["__RequestVerificationToken"],
		"cf-turnstile-response":      token,
		"DocumentKey":                documentKey,
	}

	resp, err = c.http.R().SetFormData(payload).Post(models.SearchDocument)
	if err != nil {
		return err
	}

	body2, _ := io.ReadAll(resp.Body)
	inputs2, _ := htmlutil.ExtractInputs(body2, "postForm")

	resp, err = c.http.R().SetFormData(inputs2).Post(models.PublicDocument)
	if err != nil {
		return err
	}

	body3, _ := io.ReadAll(resp.Body)
	inputs3, _ := htmlutil.ExtractInputs(body3, "postForm")

	token2, err := captcha.SolveTurnstile(models.PublicDocument)
	if err != nil {
		return err
	}

	finalPayload := map[string]string{
		"trackId": inputs3["trackId"],
		"token":   inputs3["token"],
		"captcha": token2,
	}

	resp, err = c.http.R().
		SetHeaders(map[string]string{
			"Referer": "https://catalogo-vpfe.dian.gov.co/",
			"Origin":  "https://catalogo-vpfe.dian.gov.co",
		}).
		SetFormData(finalPayload).
		SetOutputFile(outputPath).
		Post(models.DownloadDocument)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("failed download: status code %d", resp.StatusCode)
		log.Printf("[DIAN ERROR] No se pudo descargar el PDF. Key: %s | Status: %d | URL: %s",
			documentKey, resp.StatusCode, resp.Request.URL)
		return err
	}

	return nil
}
