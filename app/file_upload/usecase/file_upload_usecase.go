package usecase

import (
	"errors"
	"mime/multipart"
	"time"

	"github.com/bccfilkom/drophere-go/domain"
)

type FileUploadUseCase struct {
	UserUseCase         domain.UserService
	LinkUseCase         domain.LinkService
	StorageProviderPool domain.StorageProviderPool
}

func NewFileUploadUseCase(userUseCase domain.UserService, linkUseCase domain.LinkService, storageProviderPool domain.StorageProviderPool) domain.FileUploadUseCase {
	return &FileUploadUseCase{
		UserUseCase:         userUseCase,
		LinkUseCase:         linkUseCase,
		StorageProviderPool: storageProviderPool,
	}
}

func (fu *FileUploadUseCase) UploadFile(file multipart.File, fileHeader *multipart.FileHeader, linkID int, password string) (string, error) {
	// fetch link from database
	l, err := fu.LinkUseCase.FetchLink(uint(linkID))
	if err != nil {
		if err == domain.ErrLinkNotFound {
			return "", err
		} else {
			return "", errors.New("Server Error")
		}
	}

	// check if the link is connected to a Storage Provider
	if l.UserStorageCredentialID == nil || *l.UserStorageCredentialID < 1 || l.UserStorageCredential == nil {
		return "", errors.New("The link is unavailable")
	}

	// check for password
	if l.IsProtected() {
		if !fu.LinkUseCase.CheckLinkPassword(l, password) {
			return "", errors.New("Invalid Password")
		}
	}

	// check for deadline
	if l.Deadline != nil && l.Deadline.Before(time.Now()) {
		return "", errors.New("Link is Expired")
	}

	storageProviderService, err := fu.StorageProviderPool.Get(l.UserStorageCredential.ProviderID)
	if err != nil {
		return "", errors.New("Sorry, but the Storage Provider is unavailable at the time")
	}

	err = storageProviderService.Upload(
		domain.StorageProviderCredential{
			UserAccessToken: l.UserStorageCredential.ProviderCredential,
		},
		file,
		fileHeader.Filename,
		l.Slug,
	)
	if err != nil {
		return "", err
	}

	return "Success Uploading File", nil
}
