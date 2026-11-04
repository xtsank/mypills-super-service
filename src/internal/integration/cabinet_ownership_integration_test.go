//go:build integration

package integration_test

import (
	"context"
	"testing"

	allure "github.com/allure-framework/allure-go/commons/gotest"
	"github.com/stretchr/testify/require"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
)

func TestIntegration_CabinetRename_RejectsOtherOwner(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.DisplayName("Другой пользователь не может переименовать шкаф")
		var svc service.ICabinetService
		var repo cabinet.ICabinetRepository
		a.Step("Arrange", func(_ *allure.Context) {
			svc, repo = integrationCabinetService(t, integrationDB(t))
		})

		var err error
		a.Step("Act", func(_ *allure.Context) {
			_, err = svc.Rename(context.Background(), command.NewRenameCabinetCmd(seedCabinetOther, seedCabinet, "Stolen"))
		})

		a.Step("Assert", func(_ *allure.Context) {
			require.ErrorIs(t, err, apperrors.ErrCabinetNotFound)
			stored, readErr := repo.Find(context.Background(), seedCabinetOwner, seedCabinet)
			require.NoError(t, readErr)
			require.NotNil(t, stored)
			require.Equal(t, "Home", stored.Name)
		})
	}, allure.WithEpic("MyPills"), allure.WithFeature("Integration"), allure.WithStory("Cabinet ownership"),
		allure.WithParentSuite("MyPills"), allure.WithSuite("Integration"), allure.WithSubSuite("PostgreSQL"))
}
