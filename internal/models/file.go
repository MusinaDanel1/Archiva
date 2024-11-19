package models

type File struct {
	FilePath string `json:"file_path"`
	Size     uint64 `json:"size"`
}
