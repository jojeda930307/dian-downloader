package models

const (
	BaseURL          = "https://catalogo-vpfe.dian.gov.co"
	SearchDocument   = BaseURL + "/User/SearchDocument"
	PublicDocument   = BaseURL + "/Document/ShowDocumentToPublic"
	DownloadDocument = BaseURL + "/Document/DownloadPDF"
	SiteKey          = "0x4AAAAAAAg1WuNb-OnOa76z"
)

type DownloadRequest struct {
	DocumentKey string `json:"document_key"`
	WebhookURL  string `json:"webhook_url,omitempty"`
}

type WebhookPayload struct {
	DocumentKey   string `json:"document_key"`
	Status        string `json:"status"`
	Success       bool   `json:"success"`
	FilePath      string `json:"file_path,omitempty"`
	FileContent   string `json:"file_content,omitempty"`
	ExecutionTime string `json:"execution_time"`
	Error         string `json:"error,omitempty"`
}

type CreateTaskResponse struct {
	ErrorID int    `json:"errorId"`
	TaskID  int    `json:"taskId"`
	Code    string `json:"errorCode"`
}

type GetTaskResponse struct {
	Status   string `json:"status"`
	Solution struct {
		Token string `json:"token"`
	} `json:"solution"`
	ErrorID int    `json:"errorId"`
	Code    string `json:"errorCode"`
}
