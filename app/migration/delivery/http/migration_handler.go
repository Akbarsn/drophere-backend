package http

import (
	"net/http"

	"github.com/bccfilkom/drophere-go/domain"
	"github.com/bccfilkom/drophere-go/utils/custom_response.go"
	"github.com/go-chi/chi"
)

type MigrationHandler struct {
	UseCase domain.MigrationUseCase
}

func NewMigrationHandler(router *chi.Mux, usecase domain.MigrationUseCase) {
	handlers := MigrationHandler{
		UseCase: usecase,
	}

	router.Get("/migrate", handlers.Migrate)
}

func (m *MigrationHandler) Migrate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := m.UseCase.Migrate(ctx)
	if err != nil {
		custom_response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	custom_response.SuccessResponse(w, result)
}
