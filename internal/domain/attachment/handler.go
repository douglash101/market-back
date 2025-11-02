package attachment

import (
	"encoding/json"
	"io"
	"market/pkg/cloud"
	"market/pkg/config"
	"market/pkg/httpx"
	"market/pkg/security"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	usecase UseCase
}

func NewHandler(usecase UseCase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

// UploadAttachment godoc
// @Summary Upload a new image attachment
// @Description Upload an image file and create an attachment record. market ID is extracted from JWT token.
// @Tags attachments
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file to upload (JPG, PNG, GIF, WebP)"
// @Param description formData string false "Image description"
// @Success 201 {object} AttachmentFoundDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /attachments [post]
func (h *Handler) UploadAttachment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract user info from JWT token (set by auth middleware)
	userAuth, err := security.GetUser(r.Context())
	if err != nil {
		httpx.SendUnauthorized(w, "User not authenticated")
		return
	}

	// Parse multipart form (max 10MB)
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		httpx.SendBadRequest(w, "Failed to parse form")
		return
	}

	// Get image file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		httpx.SendBadRequest(w, "Image file is required")
		return
	}
	defer file.Close()

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		httpx.SendInternalServerError(w, "Failed to read file", err)
		return
	}

	// Get content type from file header
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(fileContent)
	}

	// Validate file type (only images)
	if !isValidImageType(contentType) {
		httpx.SendBadRequest(w, "Invalid file type. Only images (JPG, PNG, GIF, WebP) are allowed")
		return
	}

	// Upload to AWS S3
	bucketName := config.Get().CLOUD_BUCKET
	fileURL, err := cloud.Instance.Provider.UploadFile(fileContent, bucketName, contentType)
	if err != nil {
		httpx.SendInternalServerError(w, "Failed to upload file", err)
		return
	}

	// Parse optional description from form
	var description *string
	if descStr := r.FormValue("description"); descStr != "" {
		description = &descStr
	}

	// Set file type as "image"
	imageType := "image"

	// Create attachment record using company_id from JWT token
	createDTO := &AttachmentCreateDTO{
		URL:         fileURL,
		Type:        &imageType,
		Description: description,
	}

	attachment, err := h.usecase.Create(createDTO, userAuth.CompanyID)
	if err != nil {
		httpx.SendInternalServerError(w, "Failed to create attachment", err)
		return
	}

	httpx.SendCreated(w, attachment)
}

// GetAttachmentByID godoc
// @Summary Get attachment by ID
// @Description Retrieve a specific attachment by its ID
// @Tags attachments
// @Produce json
// @Param id path string true "Attachment ID"
// @Success 200 {object} AttachmentFoundDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /attachments/{id} [get]
func (h *Handler) GetAttachmentByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		httpx.SendBadRequest(w, "Invalid attachment ID format")
		return
	}

	userAuth, err := security.GetUser(r.Context())
	if err != nil {
		httpx.SendUnauthorized(w, "User not authenticated")
		return
	}

	attachment, err := h.usecase.FindByID(id, userAuth.CompanyID)
	if err != nil {
		httpx.SendInternalServerError(w, "Failed to get attachment", err)
		return
	}

	if attachment == nil {
		httpx.SendNotFound(w, "Attachment not found")
		return
	}

	httpx.SendSuccess(w, attachment)
}

// UpdateAttachment godoc
// @Summary Update attachment
// @Description Update attachment metadata
// @Tags attachments
// @Accept json
// @Produce json
// @Param id path string true "Attachment ID"
// @Param attachment body AttachmentUpdateDTO true "Attachment data to update"
// @Success 200 {object} AttachmentFoundDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /attachments/{id} [put]
func (h *Handler) UpdateAttachment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		httpx.SendBadRequest(w, "Invalid attachment ID format")
		return
	}

	var updateDTO AttachmentUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&updateDTO); err != nil {
		httpx.SendBadRequest(w, "Invalid JSON format")
		return
	}

	userAuth, err := security.GetUser(r.Context())
	if err != nil {
		httpx.SendUnauthorized(w, "User not authenticated")
		return
	}

	attachment, err := h.usecase.Update(id, userAuth.CompanyID, &updateDTO)
	if err != nil {
		httpx.SendInternalServerError(w, "Failed to update attachment", err)
		return
	}

	if attachment == nil {
		httpx.SendNotFound(w, "Attachment not found")
		return
	}

	httpx.SendSuccess(w, attachment)
}

// DeleteAttachment godoc
// @Summary Delete attachment
// @Description Delete an attachment by ID
// @Tags attachments
// @Produce json
// @Param id path string true "Attachment ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /attachments/{id} [delete]
func (h *Handler) DeleteAttachment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		httpx.SendBadRequest(w, "Invalid attachment ID format")
		return
	}

	userAuth, err := security.GetUser(r.Context())
	if err != nil {
		httpx.SendUnauthorized(w, "User not authenticated")
		return
	}

	err = h.usecase.Delete(id, userAuth.CompanyID)
	if err != nil {
		httpx.SendInternalServerError(w, "Failed to delete attachment", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// isValidImageType validates that the content type is a supported image format
func isValidImageType(contentType string) bool {
	validImageTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	return validImageTypes[contentType]
}
