//go:build integration

package integration_test

import (
	"context"
	"testing"

	allure "github.com/allure-framework/allure-go/commons/gotest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/xtsank/mypills-super-service/src/internal/domain/restock"
	"github.com/xtsank/mypills-super-service/src/internal/service"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

func TestIntegration_RestockCreate_UsesMedicineFromPostgres(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.DisplayName("Позиция пополнения сохраняется для существующего препарата")
		var svc service.IRestockService
		var repo restock.IRestockRepository
		a.Step("Arrange", func(_ *allure.Context) {
			svc, repo = integrationRestockService(t, integrationDB(t))
		})

		var created *res.RestockItemDto
		var err error
		a.Step("Act", func(_ *allure.Context) {
			created, err = svc.Create(context.Background(), command.NewCreateRestockCmd(seedRestockUser, seedMedicine, 6))
		})

		a.Step("Assert", func(_ *allure.Context) {
			require.NoError(t, err)
			require.NotNil(t, created)
			require.NotEqual(t, uuid.Nil, created.ID)
			require.Equal(t, seedMedicine, created.MedicineID)
			require.Equal(t, "Test medicine", created.MedicineName)
			require.Equal(t, float64(6), created.DesiredQuantity)
			stored, readErr := repo.Find(context.Background(), seedRestockUser, created.ID)
			require.NoError(t, readErr)
			require.NotNil(t, stored)
			require.Equal(t, float64(6), stored.DesiredQuantity)
		})
	}, allure.WithEpic("MyPills"), allure.WithFeature("Integration"), allure.WithStory("Create restock item"),
		allure.WithParentSuite("MyPills"), allure.WithSuite("Integration"), allure.WithSubSuite("PostgreSQL"))
}
