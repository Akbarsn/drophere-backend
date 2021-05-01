package http

import (
	"net/http"
	"strconv"

	"github.com/bccfilkom/drophere-go/domain"
	"github.com/bccfilkom/drophere-go/utils/custom_response.go"
	"github.com/go-chi/chi"
)

type FileUploadHandler struct {
	FileUploadUseCase domain.FileUploadUseCase
}

func NewFileUploadHandler(router *chi.Mux, fileUploadUseCase domain.FileUploadUseCase) {
	handler := FileUploadHandler{FileUploadUseCase: fileUploadUseCase}

	router.Post("/uploadfile", handler.FileUpload)
}

func (fu *FileUploadHandler) FileUpload(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		custom_response.ErrorResponse(w, http.StatusBadRequest, "Invalid File")
		return
	}

	linkID, err := strconv.Atoi(r.FormValue("linkId"))
	if err != nil {
		custom_response.ErrorResponse(w, http.StatusBadRequest, "Invalid Link ID")
		return
	}

	password := r.FormValue("password")

	result, err := fu.FileUploadUseCase.UploadFile(file, fileHeader, linkID, password)
	if err != nil {
		custom_response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	custom_response.SuccessResponse(w, result)
}
