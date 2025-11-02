package cloud

import (
	"bytes"
	"market/pkg/config"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type AWS struct {
	Sessions *session.Session
}

func NewAWS() *AWS {
	return &AWS{}
}

func (a *AWS) Bootstrap() {
	CLOUD_HOST := config.Get().CLOUD_HOST
	CLOUD_KEY := config.Get().CLOUD_KEY
	CLOUD_SECRET := config.Get().CLOUD_SECRET
	CLOUD_REGION := config.Get().CLOUD_REGION
	CLOUD_DISABLED_SSL := config.Get().CLOUD_DISABLED_SSL

	sess, err := session.NewSession(&aws.Config{
		Region:     aws.String(CLOUD_REGION),
		Endpoint:   aws.String(CLOUD_HOST),
		DisableSSL: aws.Bool(CLOUD_DISABLED_SSL),
		Credentials: credentials.NewStaticCredentials(
			CLOUD_KEY,
			CLOUD_SECRET,
			""),
	})

	if err != nil {
		fmt.Println("Error creating AWS session:", err)
	}

	a.Sessions = sess
}

func (a *AWS) UploadImageFromURL(url string, bucket string, imageID string) error {

	bufImage, err := downloadImage(url)
	if err != nil {
		return fmt.Errorf("failed to download image from URL: %v", err)
	}

	s3Client := s3.New(a.Sessions)

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(imageID),
		Body:        bytes.NewReader(bufImage.Bytes()),
		ContentType: aws.String(http.DetectContentType(bufImage.Bytes())),
	})

	if err != nil {
		return fmt.Errorf("failed to upload image to S3: %v", err)
	}

	fmt.Printf("Imagem carregada com sucesso para o bucket %s com a chave %s\n", bucket, imageID)
	return nil
}

func (a *AWS) UploadFile(fileContent []byte, bucket string, contentType string) (string, error) {
	// Gera um UUID único para o arquivo
	imageID := fmt.Sprintf("%s.%s", generateUUID(), getFileExtensionFromContentType(contentType))

	s3Client := s3.New(a.Sessions)

	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(imageID),
		Body:        bytes.NewReader(fileContent),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	// Retorna a URL do arquivo
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, config.Get().CLOUD_REGION, imageID)

	fmt.Printf("Arquivo carregado com sucesso para o bucket %s com a chave %s\n", bucket, imageID)
	return fileURL, nil
}

func downloadImage(url string) (*bytes.Buffer, error) {
	// Baixa a imagem da URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	// Lê o conteúdo da imagem
	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %v", err)
	}
	return &buf, nil
}

func (a *AWS) GetSession() interface{} {
	return a.Sessions
}

func generateUUID() string {
	return uuid.New().String()
}

func getFileExtensionFromContentType(contentType string) string {
	switch {
	case strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg"):
		return "jpg"
	case strings.Contains(contentType, "png"):
		return "png"
	case strings.Contains(contentType, "gif"):
		return "gif"
	case strings.Contains(contentType, "webp"):
		return "webp"
	case strings.Contains(contentType, "pdf"):
		return "pdf"
	case strings.Contains(contentType, "doc"):
		return "doc"
	case strings.Contains(contentType, "docx"):
		return "docx"
	default:
		return "bin"
	}
}
