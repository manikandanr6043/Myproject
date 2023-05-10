package model

// CommitMessage -> defines commit message event structure
type CommitMessage struct {
	SpecVersion     string     `json:"specversion"`
	ID              string     `json:"id"`
	Type            string     `json:"type"`
	Subject         string     `json:"subject"`
	Time            string     `json:"time"`
	DataContentType string     `json:"datacontenttype"`
	Data            UploadData `json:"data"`
	RetryCount      int32      `json:"retryCount"`
}

// UploadData -> defines struct for upload data
type UploadData struct {
	SpaceId  string `json:"spaceId"`
	FileId   string `json:"fileId"`
	UploadId string `json:"uploadId"`
}
