package usecase

import (
	"bytes"
	"fmt"
	"github.com/bccfilkom/drophere-go/domain"
	htmlTemplate "html/template"
	textTemplate "text/template"
	"time"
)

const defaultTokenExpiryDuration int = 5

type UserUseCase struct {
	UserRepoMysql        domain.UserRepository
	UserStorageRepoMysql domain.UserStorageCredentialRepository
	Authenticator        domain.Authenticator
	Mailer               domain.Mailer
	PasswordHash         domain.Hasher
	StringGenerator      domain.StringGenerator
	StorageProviderPool  domain.StorageProviderPool
	HtmlTemplates        *htmlTemplate.Template
	TextTemplates        *textTemplate.Template
	UserConfig           domain.UserConfig
}

func NewUserUseCase(
	urMysql domain.UserRepository,
	usrMysql domain.UserStorageCredentialRepository,
	authenticator domain.Authenticator,
	mailer domain.Mailer,
	passwordHash domain.Hasher,
	stringGenerator domain.StringGenerator,
	storageProviderPool domain.StorageProviderPool,
	htmlTemplates *htmlTemplate.Template,
	textTemplates *textTemplate.Template,
	userConfig domain.UserConfig,
) domain.UserService {
	return &UserUseCase{
		UserRepoMysql:        urMysql,
		UserStorageRepoMysql: usrMysql,
		Authenticator:        authenticator,
		Mailer:               mailer,
		PasswordHash:         passwordHash,
		StringGenerator:      stringGenerator,
		StorageProviderPool:  storageProviderPool,
		HtmlTemplates:        htmlTemplates,
		TextTemplates:        textTemplates,
		UserConfig:           userConfig,
	}
}

func (uuc UserUseCase) Register(email, name, password string) (*domain.User, error) {
	user, err := uuc.UserRepoMysql.FindByEmail(email)
	if err != nil && err != domain.ErrUserNotFound {
		return nil, err
	}

	if user != nil {
		return nil, domain.ErrUserDuplicated
	}

	user = &domain.User{
		Email: email,
		Name:  name,
	}

	user.Password, err = uuc.PasswordHash.Hash(password)
	if err != nil {
		return nil, err
	}
	return uuc.UserRepoMysql.Create(user)
}

func (uuc UserUseCase) Auth(email, password string) (*domain.UserCredentials, error) {
	user, err := uuc.UserRepoMysql.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if !uuc.PasswordHash.Verify(user.Password, password) {
		return nil, domain.ErrUserInvalidPassword
	}

	return uuc.Authenticator.Authenticate(user)
}

func (uuc UserUseCase) Update(userID uint, name, password, oldPassword *string) (*domain.User, error) {
	u, err := uuc.UserRepoMysql.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if password != nil {
		if oldPassword == nil || !uuc.PasswordHash.Verify(u.Password, *oldPassword) {
			return nil, domain.ErrUserInvalidPassword
		}

		u.Password, err = uuc.PasswordHash.Hash(*password)
		if err != nil {
			return nil, err
		}
	}

	if name != nil {
		u.Name = *name
	}

	return uuc.UserRepoMysql.Update(u)
}

func (uuc UserUseCase) ConnectStorageProvider(userID, providerID uint, providerCredential string) error {
	storageProvider, err := uuc.StorageProviderPool.Get(providerID)
	if err != nil {
		return err
	}

	u, err := uuc.UserRepoMysql.FindByID(userID)
	if err != nil {
		return err
	}

	storageProviderAccount, err := storageProvider.AccountInfo(
		domain.StorageProviderCredential{
			UserAccessToken: providerCredential,
		},
	)
	if err != nil {
		return err
	}

	var cred domain.UserStorageCredential

	creds, err := uuc.UserStorageRepoMysql.Find(domain.UserStorageCredentialFilters{
		UserIDs:     []uint{u.ID},
		ProviderIDs: []uint{providerID},
	}, false)
	if err != nil {
		return err
	}

	if len(creds) > 0 {
		cred = creds[0]
		cred.ProviderCredential = providerCredential
		cred.Email = storageProviderAccount.Email
		cred.Photo = storageProviderAccount.Photo
		cred, err = uuc.UserStorageRepoMysql.Update(cred)
	} else {
		cred, err = uuc.UserStorageRepoMysql.Create(domain.UserStorageCredential{
			UserID:             u.ID,
			ProviderID:         providerID,
			ProviderCredential: providerCredential,
			Email:              storageProviderAccount.Email,
			Photo:              storageProviderAccount.Photo,
		})
	}

	return err
}

func (uuc UserUseCase) DisconnectStorageProvider(userID, providerID uint) error {
	storageProvider, err := uuc.StorageProviderPool.Get(providerID)
	if err != nil {
		return err
	}

	u, err := uuc.UserRepoMysql.FindByID(userID)
	if err != nil {
		return err
	}

	creds, err := uuc.UserStorageRepoMysql.Find(domain.UserStorageCredentialFilters{
		UserIDs:     []uint{u.ID},
		ProviderIDs: []uint{storageProvider.ID()},
	}, false)
	if err != nil {
		return err
	}

	if len(creds) > 0 {
		err = uuc.UserStorageRepoMysql.Delete(creds[0])
		if err != nil {
			return err
		}
	}

	return nil
}

func (uuc UserUseCase) ListStorageProviders(userID uint) ([]domain.UserStorageCredential, error) {
	return uuc.UserStorageRepoMysql.Find(domain.UserStorageCredentialFilters{
		UserIDs: []uint{userID},
	}, false)
}

func (uuc UserUseCase) UpdateStorageToken(userID uint, dropboxToken *string) (*domain.User, error) {
	u, err := uuc.UserRepoMysql.FindByID(userID)
	if err != nil {
		return nil, err
	}

	u.DropboxToken = dropboxToken

	return uuc.UserRepoMysql.Update(u)
}

func (uuc UserUseCase) sendPasswordRecoveryTokenToEmail(to domain.MailAddress, subject, email, token string) error {

	// preparing template
	htmlTmpl := uuc.HtmlTemplates.Lookup("request_password_recovery_html")
	if htmlTmpl == nil {
		return domain.ErrTemplateNotFound
	}

	textTmpl := uuc.TextTemplates.Lookup("request_password_recovery_text")
	if textTmpl == nil {
		return domain.ErrTemplateNotFound
	}

	// preparing template content
	messageData := map[string]string{
		"ResetPasswordLink": fmt.Sprintf(
			"%s?token=%s&email=%s",
			uuc.UserConfig.RecoverPasswordWebURL,
			token,
			email,
		),
		"Token": token,
	}

	// injecting data to template
	htmlMessage := &bytes.Buffer{}
	htmlTmpl.Execute(htmlMessage, messageData)

	textMessage := &bytes.Buffer{}
	textTmpl.Execute(textMessage, messageData)

	from := domain.MailAddress{
		Address: "admin@drophere.link",
		Name:    "Drophere Bot",
	}

	if uuc.UserConfig.MailerEmail != "" {
		from.Address = uuc.UserConfig.MailerEmail
	}

	if uuc.UserConfig.MailerName != "" {
		from.Name = uuc.UserConfig.MailerName
	}

	// send email
	return uuc.Mailer.Send(
		from,
		to,
		subject,
		textMessage.String(),
		htmlMessage.String(),
	)
}

func (uuc UserUseCase) RequestPasswordRecovery(email string) error {
	u, err := uuc.UserRepoMysql.FindByEmail(email)
	if err != nil {
		return err
	}

	// TODO: check if user has already requested password recovery to avoid spam
	tokenExpiryDuration := defaultTokenExpiryDuration
	if uuc.UserConfig.PasswordRecoveryTokenExpiryDuration > 0 {
		tokenExpiryDuration = uuc.UserConfig.PasswordRecoveryTokenExpiryDuration
	}

	token := uuc.StringGenerator.Generate()
	tokenExpiry := time.Now().Add(time.Minute * time.Duration(tokenExpiryDuration))
	u.RecoverPasswordToken = &token
	u.RecoverPasswordTokenExpiry = &tokenExpiry

	// save the user
	u, err = uuc.UserRepoMysql.Update(u)
	if err != nil {
		return err
	}

	// send email
	err = uuc.sendPasswordRecoveryTokenToEmail(
		domain.MailAddress{
			Address: u.Email,
			Name:    u.Name,
		},
		"Recover Password",
		u.Email,
		token,
	)
	if err != nil {
		return err
	}

	return nil
}

func (uuc UserUseCase) RecoverPassword(email, token, newPassword string) error {
	u, err := uuc.UserRepoMysql.FindByEmail(email)
	if err != nil {
		return err
	}

	if token == "" || u.RecoverPasswordToken == nil || *u.RecoverPasswordToken != token {
		return domain.ErrUserNotFound
	}

	if u.RecoverPasswordTokenExpiry == nil || time.Now().After(*u.RecoverPasswordTokenExpiry) {
		return domain.ErrUserPasswordRecoveryTokenExpired
	}

	u.Password, err = uuc.PasswordHash.Hash(newPassword)
	if err != nil {
		return err
	}

	u.RecoverPasswordToken, u.RecoverPasswordTokenExpiry = nil, nil

	u, err = uuc.UserRepoMysql.Update(u)
	if err != nil {
		return err
	}

	return nil
}
