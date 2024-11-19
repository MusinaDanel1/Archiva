package models

type Archive struct {
	Filename    string  `json:"filename"`
	ArchiveSize float64 `json:"archive_size"`
	TotalSize   float64 `json:"total_size"`
	TotalFiles  float64 `json:"total_files"`
	Files       []File  `json:"files"`
}
