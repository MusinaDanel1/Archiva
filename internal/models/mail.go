package models

type EmailRequest struct {
	FilePath       string   `json:"file_path"`
	Filename       string   `json:"filename"`
	EmailAddresses []string `json:"emails"`
}
