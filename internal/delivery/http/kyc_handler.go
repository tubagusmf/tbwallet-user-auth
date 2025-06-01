package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tubagusmf/tbwallet-user-auth/internal/model"
)

type KycHandler struct {
	kycUsecase model.IKycDocUsecase
}

func NewKycHandler(e *echo.Echo, kycUsecase model.IKycDocUsecase) {
	handlers := &KycHandler{
		kycUsecase: kycUsecase,
	}

	routeKyc := e.Group("v1/kyc")
	routeKyc.GET("/:id", handlers.GetByID, AuthMiddleware)
	routeKyc.GET("/user/:user_id", handlers.GetByUserID, AuthMiddleware)
	routeKyc.POST("/create", handlers.Create, AuthMiddleware)
	routeKyc.PUT("/update/:id", handlers.Update, AuthMiddleware)
}

func (h *KycHandler) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid kyc ID format")
	}

	kyc, err := h.kycUsecase.GetByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Kyc not found")
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   kyc,
	})
}

func (h *KycHandler) GetByUserID(c echo.Context) error {
	idParam := c.Param("user_id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	kyc, err := h.kycUsecase.GetByUserID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Kyc not found")
	}

	return c.JSON(http.StatusOK, Response{
		Status: http.StatusOK,
		Data:   kyc,
	})
}

func (h *KycHandler) Create(c echo.Context) error {
	userIDStr := c.FormValue("user_id")
	documentType := c.FormValue("document_type")

	file, err := c.FormFile("document_url")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "File is required")
	}

	// Validate file size
	const MaxFileSize = 2 << 20 // 2MB
	if file.Size > MaxFileSize {
		return echo.NewHTTPError(http.StatusBadRequest, "File size exceeds 2MB limit")
	}

	// Validate extension file
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".pdf": true}
	if !allowed[ext] {
		return echo.NewHTTPError(http.StatusBadRequest, "Unsupported file type")
	}

	// Validate user ID
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open file")
	}
	defer src.Close()

	// Save file to local
	dstPath := fmt.Sprintf("./uploads/documents/%d_%s", time.Now().UnixNano(), file.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save file")
	}

	input := model.CreateKycDocInput{
		UserID:       userID,
		DocumentType: documentType,
		DocumentURL:  dstPath, // changed to cloud storage if needed
	}

	kyc, err := h.kycUsecase.Create(c.Request().Context(), input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create kyc")
	}

	return c.JSON(http.StatusCreated, Response{
		Status:  http.StatusCreated,
		Message: "Upload Document successfully",
		Data:    kyc,
	})
}

func (h *KycHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid kyc ID format")
	}

	userIDStr := c.FormValue("user_id")
	documentType := c.FormValue("document_type")
	file, err := c.FormFile("document_url")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "File is required")
	}

	// Validate user ID
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	// Validate file size
	const MaxFileSize = 2 << 20 // 2MB
	if file.Size > MaxFileSize {
		return echo.NewHTTPError(http.StatusBadRequest, "File size exceeds 2MB limit")
	}

	// Validate extension file
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".pdf": true}
	if !allowed[ext] {
		return echo.NewHTTPError(http.StatusBadRequest, "Unsupported file type")
	}

	// Open & Save file
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open file")
	}
	defer src.Close()

	dstPath := fmt.Sprintf("./uploads/documents/%d_%s", time.Now().UnixNano(), file.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save file")
	}

	input := model.UpdateKycDocInput{
		UserID:       userID,
		DocumentType: documentType,
		DocumentURL:  dstPath,
	}

	kyc, err := h.kycUsecase.Update(c.Request().Context(), id, input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update kyc")
	}

	return c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "Update Document successfully",
		Data:    kyc,
	})
}
