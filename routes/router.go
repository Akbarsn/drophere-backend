package router

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	htmlTemplate "html/template"
	textTemplate "text/template"

	"github.com/99designs/gqlgen/handler"
	drophere_go "github.com/bccfilkom/drophere-go"
	_fileUploadHandler "github.com/bccfilkom/drophere-go/app/file_upload/delivery/http"
	_fileUploadUseCase "github.com/bccfilkom/drophere-go/app/file_upload/usecase"
	_linkUseCase "github.com/bccfilkom/drophere-go/app/link/usecase"
	_userUseCase "github.com/bccfilkom/drophere-go/app/user/usecase"
	"github.com/bccfilkom/drophere-go/domain"
	_linkRepository "github.com/bccfilkom/drophere-go/infrastructure/database/mysql"
	_userRepository "github.com/bccfilkom/drophere-go/infrastructure/database/mysql"
	_userStorageRepository "github.com/bccfilkom/drophere-go/infrastructure/database/mysql"
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
	userStorageRepo := _userStorageRepository.NewUserStorageCredentialRepository(db)

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
	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		Debug:            debugMode,
	}).Handler)

	// router.Use(jwt_tools.Middleware())
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Handler for GraphQL
	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", handler.GraphQL(drophere_go.NewExecutableSchema(drophere_go.Config{Resolvers: resolver})))

	// Handler for File Upload
	fileUploadUsecase := _fileUploadUseCase.NewFileUploadUseCase(userUseCase, linkUseCase, storageProviderPool)
	_fileUploadHandler.NewFileUploadHandler(router, fileUploadUsecase)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", appEnv.Port)
	err = http.ListenAndServe(":"+appEnv.Port, router)
	if err != nil {
		log.Fatal(err)
	}
}
