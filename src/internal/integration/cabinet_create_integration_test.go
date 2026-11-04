//go:build integration

package integration_test

import (
	"context"
	"testing"

	allure "github.com/allure-framework/allure-go/commons/gotest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet"
	"github.com/xtsank/mypills-super-service/src/internal/service"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

func TestIntegration_CabinetCreate_PersistsInPostgres(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.DisplayName("Созданный шкаф сохраняется в PostgreSQL")
		var svc service.ICabinetService
		var repo cabinet.ICabinetRepository
		a.Step("Arrange", func(_ *allure.Context) {
			svc, repo = integrationCabinetService(t, integrationDB(t))
		})

		var created *res.CabinetDto
		var err error
		a.Step("Act", func(_ *allure.Context) {
			created, err = svc.Create(context.Background(), command.NewCreateCabinetCmd(seedCabinetCreateUser, "Travel kit"))
		})

		a.Step("Assert", func(_ *allure.Context) {
			require.NoError(t, err)
			require.NotNil(t, created)
			require.NotEqual(t, uuid.Nil, created.ID)
			require.Equal(t, "Travel kit", created.Name)
			stored, readErr := repo.Find(context.Background(), seedCabinetCreateUser, created.ID)
			require.NoError(t, readErr)
			require.NotNil(t, stored)
			require.Equal(t, "Travel kit", stored.Name)
			require.Equal(t, seedCabinetCreateUser, stored.UserID)
		})
	}, allure.WithEpic("MyPills"), allure.WithFeature("Integration"), allure.WithStory("Create cabinet"),
		allure.WithParentSuite("MyPills"), allure.WithSuite("Integration"), allure.WithSubSuite("PostgreSQL"))
}
