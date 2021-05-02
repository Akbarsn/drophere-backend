package routes

import (
	"log"
	"path/filepath"
	"time"

	htmlTemplate "html/template"
	textTemplate "text/template"

	"github.com/99designs/gqlgen/handler"
	drophere_go "github.com/bccfilkom/drophere-go"
	_fileUploadHandler "github.com/bccfilkom/drophere-go/app/file_upload/delivery/http"
	_fileUploadUseCase "github.com/bccfilkom/drophere-go/app/file_upload/usecase"
	_linkRepository "github.com/bccfilkom/drophere-go/app/link/repository/mysql"
	_linkUseCase "github.com/bccfilkom/drophere-go/app/link/usecase"
	_migrationHandler "github.com/bccfilkom/drophere-go/app/migration/delivery/http"
	_migrationRepository "github.com/bccfilkom/drophere-go/app/migration/repository"
	_migrationUseCase "github.com/bccfilkom/drophere-go/app/migration/usecase"
	_userRepository "github.com/bccfilkom/drophere-go/app/user/repository/mysql"
	_userUseCase "github.com/bccfilkom/drophere-go/app/user/usecase"
	_userStorageRepository "github.com/bccfilkom/drophere-go/app/user_storage/repository/mysql"
	"github.com/bccfilkom/drophere-go/domain"
	"github.com/bccfilkom/drophere-go/utils/db_driver"
	"github.com/bccfilkom/drophere-go/utils/env_driver"
	"github.com/bccfilkom/drophere-go/utils/jwt_tools"
	"github.com/bccfilkom/drophere-go/utils/mailer_service"
	"github.com/bccfilkom/drophere-go/utils/security_tools"
	"github.com/bccfilkom/drophere-go/utils/storage_service"
	"github.com/bccfilkom/drophere-go/utils/string_tools"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

type Router struct {
	Router *chi.Mux
}

func (r *Router) NewChiRoutes() {
	// Read Env Variable
	appEnv, err := env_driver.NewAppEnvironmentDriver()
	if err != nil {
		log.Fatal(err)
	}

	dbEnv, err := env_driver.NewDatabaseEnvironmentDriver()
	if err != nil {
		log.Fatal(err)
	}

	jwtEnv, err := env_driver.NewJWTEnvironmentDriver()
	if err != nil {
		log.Fatal(err)
	}

	mailerEnv, err := env_driver.NewSendgridMailerDriver()

	debugMode := false
	if appEnv.Mode == "debug" {
		debugMode = true
	}

	remoteDirectory := "drophere"
	if appEnv.StorageRootDirectoryName != "" {
		remoteDirectory = appEnv.StorageRootDirectoryName
	}

	// Connect to Database
	db, err := db_driver.NewMysqlConn(dbEnv)
	if err != nil {
		log.Fatal(err)
	}

	// User
	userRepo := _userRepository.NewUserRepository(db)

	// UserStorage
	userStorageRepo := _userStorageRepository.NewUserStorageRepository(db)

	// Link
	linkRepo := _linkRepository.NewLinkRepository(db)

	// Initialize JWT Tools
	jwtAuth := jwt_tools.NewJWT(
		jwtEnv.Secret,
		time.Duration(jwtEnv.Duration),
		jwtEnv.SignAlgorithm,
		userRepo,
	)

	// Initialize Bcrypt Hasher
	bcryptHasher := security_tools.NewBcryptHasher()

	// Initialize UUID Generator
	uuidGenerator := string_tools.NewUUID()

	// Initialize Sendgrid Mailer Service
	sendGridMailer := mailer_service.NewSendgrid(
		mailerEnv.APIKey,
		debugMode,
	)

	// Initialize Storage Service
	dropboxService := storage_service.NewDropboxStorageProvider(remoteDirectory)
	storageProviderPool := domain.StorageProviderPool{}
	storageProviderPool.Register(dropboxService)

	// Initialize Template File
	htmlTemplates, err := htmlTemplate.ParseGlob(filepath.Join(appEnv.TemplatePath, "html", "*.html"))
	if err != nil {
		panic(err)
	}

	textTemplates, err := textTemplate.ParseGlob(filepath.Join(appEnv.TemplatePath, "text", "*.txt"))
	if err != nil {
		panic(err)
	}

	// Initialize User Service
	userUseCase := _userUseCase.NewUserUseCase(
		userRepo,
		userStorageRepo,
		jwtAuth,
		sendGridMailer,
		bcryptHasher,
		uuidGenerator,
		storageProviderPool,
		htmlTemplates,
		textTemplates,
		domain.UserConfig{
			PasswordRecoveryTokenExpiryDuration: appEnv.PasswordRecovery.TokenDuration,
			RecoverPasswordWebURL:               appEnv.PasswordRecovery.RecoveryWebURL,
			MailerEmail:                         appEnv.PasswordRecovery.MailerEmail,
			MailerName:                          appEnv.PasswordRecovery.MailerName,
		},
	)

	linkUseCase := _linkUseCase.NewLinkUseCase(
		linkRepo,
		userStorageRepo,
		bcryptHasher,
	)

	resolver := drophere_go.NewResolver(userUseCase, jwtAuth, linkUseCase)

	// Setup Router
	r.Router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		Debug:            debugMode,
	}).Handler)

	// router.Use(jwt_tools.Middleware())
	r.Router.Use(middleware.RequestID)
	r.Router.Use(middleware.RealIP)
	r.Router.Use(middleware.Logger)
	r.Router.Use(middleware.Recoverer)

	// File Upload router
	fileUploadUsecase := _fileUploadUseCase.NewFileUploadUseCase(userUseCase, linkUseCase, storageProviderPool)
	_fileUploadHandler.NewFileUploadHandler(r.Router, fileUploadUsecase)

	// Migration router
	migrationRepository := _migrationRepository.NewMigrationRepository(db)
	migrationUseCase := _migrationUseCase.NewMigrationUseCase(migrationRepository)
	_migrationHandler.NewMigrationHandler(r.Router, migrationUseCase)

	// Handler for GraphQL
	r.Router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	r.Router.Handle("/query", handler.GraphQL(drophere_go.NewExecutableSchema(drophere_go.Config{Resolvers: resolver})))

}
