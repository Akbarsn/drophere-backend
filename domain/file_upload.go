package domain

import "mime/multipart"

type FileUploadUseCase interface {
	UploadFile(file multipart.File, fileHeader *multipart.FileHeader, linkID int, password string) (string, error)
}
