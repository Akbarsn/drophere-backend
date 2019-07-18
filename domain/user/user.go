package user

import (
	"bytes"
	"time"

	htmlTemplate "html/template"
	textTemplate "text/template"

	"github.com/bccfilkom/drophere-go/domain"
)

type service struct {
	userRepo        domain.UserRepository
	authenticator   domain.Authenticator
	mailer          domain.Mailer
	passwordHasher  domain.Hasher
	stringGenerator domain.StringGenerator

	htmlTemplates *htmlTemplate.Template
	textTemplates *textTemplate.Template
}

// NewService returns service instance
func NewService(
	userRepo domain.UserRepository,
	authenticator domain.Authenticator,
	mailer domain.Mailer,
	passwordHasher domain.Hasher,
	stringGenerator domain.StringGenerator,
	htmlTemplates *htmlTemplate.Template,
	textTemplates *textTemplate.Template,
) domain.UserService {
	return &service{
		userRepo:        userRepo,
		authenticator:   authenticator,
		mailer:          mailer,
		passwordHasher:  passwordHasher,
		stringGenerator: stringGenerator,

		htmlTemplates: htmlTemplates,
		textTemplates: textTemplates,
	}
}

// Register implementation
func (s *service) Register(email, name, password string) (*domain.User, error) {
	// check for existing email prior to creating new user
	user, err := s.userRepo.FindByEmail(email)
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

	user.Password, err = s.passwordHasher.Hash(password)
	if err != nil {
		return nil, err
	}
	return s.userRepo.Create(user)
}

// Auth implementation
func (s *service) Auth(email, password string) (*domain.UserCredentials, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if !s.passwordHasher.Verify(user.Password, password) {
		return nil, domain.ErrUserInvalidPassword
	}

	return s.authenticator.Authenticate(user)
}

// Update implementation
func (s *service) Update(userID uint, name, newPassword, oldPassword *string) (*domain.User, error) {
	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if newPassword != nil {
		if oldPassword == nil || !s.passwordHasher.Verify(u.Password, *oldPassword) {
			return nil, domain.ErrUserInvalidPassword
		}

		u.Password, err = s.passwordHasher.Hash(*newPassword)
		if err != nil {
			return nil, err
		}
	}

	if name != nil {
		u.Name = *name
	}

	return s.userRepo.Update(u)
}

// UpdateStorageToken implementation
func (s *service) UpdateStorageToken(userID uint, dropboxToken *string) (*domain.User, error) {
	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	u.DropboxToken = dropboxToken

	return s.userRepo.Update(u)
}

// RequestPasswordRecovery implementation
func (s *service) RequestPasswordRecovery(email string) error {
	u, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}

	token := s.stringGenerator.Generate()
	tokenExpiry := time.Now().Add(time.Minute * 5)
	u.RecoverPasswordToken = &token
	u.RecoverPasswordTokenExpiry = &tokenExpiry

	// save the user
	u, err = s.userRepo.Update(u)
	if err != nil {
		return err
	}

	// send email
	err = s.sendPasswordRecoveryTokenToEmail(
		domain.MailAddress{
			Address: u.Email,
			Name:    u.Name,
		},
		"Recover Password",
		token,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) sendPasswordRecoveryTokenToEmail(to domain.MailAddress, subject, token string) error {

	// preparing template
	htmlTmpl := s.htmlTemplates.Lookup("request_password_recovery_html")
	if htmlTmpl == nil {
		return domain.ErrTemplateNotFound
	}

	textTmpl := s.textTemplates.Lookup("request_password_recovery_text")
	if textTmpl == nil {
		return domain.ErrTemplateNotFound
	}

	// preparing template content
	messageData := map[string]string{
		"Token": token,
	}

	// injecting data to template
	htmlMessage := &bytes.Buffer{}
	htmlTmpl.Execute(htmlMessage, messageData)

	textMessage := &bytes.Buffer{}
	textTmpl.Execute(textMessage, messageData)

	// send email
	return s.mailer.Send(
		domain.MailAddress{
			Address: "admin@drophere.link",
			Name:    "Drophere Bot",
		},
		to,
		subject,
		textMessage.String(),
		htmlMessage.String(),
	)
}

func (s *service) RecoverPassword(email, token, newPassword string) error {
	u, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}

	if token == "" || u.RecoverPasswordToken == nil || *u.RecoverPasswordToken != token {
		return domain.ErrUserNotFound
	}

	if u.RecoverPasswordTokenExpiry == nil || time.Now().After(*u.RecoverPasswordTokenExpiry) {
		return domain.ErrUserPasswordRecoveryTokenExpired
	}

	u.Password, err = s.passwordHasher.Hash(newPassword)
	if err != nil {
		return err
	}

	u.RecoverPasswordToken, u.RecoverPasswordTokenExpiry = nil, nil

	u, err = s.userRepo.Update(u)
	if err != nil {
		return err
	}

	return nil
}
