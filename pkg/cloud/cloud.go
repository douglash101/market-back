package cloud

import "sync"

type Provider string

const (
	AWS_PROVIDER Provider = "AWS"
	GCP_PROVIDER Provider = "GCP"
)

type CloudProvider interface {
	GetSession() interface{}
	Bootstrap()
	UploadImageFromURL(url string, bucket string, imageID string) error
	UploadFile(fileContent []byte, bucket string, contentType string) (string, error)
}

type Cloud struct {
	Provider CloudProvider
}

var once sync.Once
var Instance *Cloud

func NewCloudInstance(provider Provider) {
	once.Do(func() {
		switch provider {
		case AWS_PROVIDER:
			Instance = &Cloud{
				Provider: NewAWS(),
			}

			Instance.Provider.Bootstrap()
		case GCP_PROVIDER:
			// instance = &CloudProvider{}
		default:
			panic("Invalid cloud provider")
		}
	})
}
