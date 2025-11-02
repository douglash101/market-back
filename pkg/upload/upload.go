package upload

type Path string

const (
	PRODUCT_IMAGES Path = "product-images"
)

// func UploadImage(providerImageURL string, path Path) string {

// 	HOST_BUCKET := config.Get().CLOUD_HOST_BUCKET

// 	bucket := fmt.Sprintf(config.Get().CLOUD_BUCKET+"/%s/", path)

// 	imageID := uuid.New().String()
// 	err := cloud.Instance.Provider.UploadImageFromURL(
// 		providerImageURL,
// 		bucket,
// 		imageID,
// 	)

// 	if err != nil {
// 		fmt.Println(err)
// 		return "NO_IMAGE"
// 	}

// 	return fmt.Sprintf(HOST_BUCKET+"/%s/%s", path, imageID)
// }
