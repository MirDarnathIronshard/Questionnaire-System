package models

const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB

	FileTypeImage = ".jpg,.jpeg,.png,.gif"
	FileTypeDoc   = ".pdf,.doc,.docx"
	FileTypeAudio = ".mp3,.wav"
	FileTypeVideo = ".mp4,.avi"
)

var (
	AllowedImageTypes = []string{".jpg", ".jpeg", ".png", ".gif"}
	AllowedDocTypes   = []string{".pdf", ".doc", ".docx"}
	AllowedAudioTypes = []string{".mp3", ".wav"}
	AllowedVideoTypes = []string{".mp4", ".avi"}
)

type FileUploadConfig struct {
	UploadDir   string
	BaseURL     string
	MaxFileSize int64
	TempDir     string
}
