//go:build e2e

package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	allure "github.com/allure-framework/allure-go/commons/gotest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	illnessID        = "00000000-0000-4000-8000-000000000041"
	safeMedicineID   = "00000000-0000-4000-8000-000000000042"
	unsafeMedicineID = "00000000-0000-4000-8000-000000000043"
)

func TestE2E_DriverGetsSafeMedicineFromCabinet(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.DisplayName("Водитель получает безопасное лекарство из своей аптечки")
		baseURL := os.Getenv("E2E_BASE_URL")
		if baseURL == "" {
			baseURL = "http://app:8080"
		}
		client := &http.Client{Timeout: 10 * time.Second}
		var token string
		var cabinetID, safeItemID uuid.UUID
		manufactured := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

		a.Step("Arrange: зарегистрировать водителя", func(_ *allure.Context) {
			var auth struct {
				Token string `json:"token"`
			}
			postJSON(t, client, baseURL+"/api/v2/users", "", map[string]any{
				"login": "e2e-driver", "email": "e2e-driver@example.test", "password": "strongpass123",
				"age": 30, "weight": 70, "is_pregnant": false, "is_driver": true,
				"allergies": []string{}, "illnesses": []string{},
			}, http.StatusCreated, &auth)
			require.NotEmpty(t, auth.Token)
			token = auth.Token
		})

		a.Step("Arrange: создать аптечку и положить в неё два лекарства", func(_ *allure.Context) {
			var cabinet struct {
				ID uuid.UUID `json:"id"`
			}
			postJSON(t, client, baseURL+"/api/v2/cabinets", token,
				map[string]any{"name": "E2E home cabinet"}, http.StatusCreated, &cabinet)
			require.NotEqual(t, uuid.Nil, cabinet.ID)
			cabinetID = cabinet.ID
			itemsURL := fmt.Sprintf("%s/api/v2/cabinets/%s/items", baseURL, cabinetID)
			for _, medicineID := range []string{safeMedicineID, unsafeMedicineID} {
				var item struct {
					ID uuid.UUID `json:"id"`
				}
				postJSON(t, client, itemsURL, token, map[string]any{
					"medicine_id": medicineID, "date_of_manufacture": manufactured,
					"quantity": 10, "minimum_quantity": 1,
				}, http.StatusCreated, &item)
				require.NotEqual(t, uuid.Nil, item.ID)
				if medicineID == safeMedicineID {
					safeItemID = item.ID
				}
			}
		})

		var result struct {
			Items []struct {
				MedicineID        uuid.UUID `json:"medicine_id"`
				CabinetID         uuid.UUID `json:"cabinet_id"`
				CabinetItemID     uuid.UUID `json:"cabinet_item_id"`
				Name              string    `json:"name"`
				Dosage            float64   `json:"dosage"`
				UnitName          string    `json:"unit_name"`
				FrequencyPerDay   int       `json:"frequency_per_day"`
				QuantityInCabinet float64   `json:"quantity_in_cabinet"`
			} `json:"items"`
		}
		a.Step("Act: запросить подбор лекарства от головной боли", func(_ *allure.Context) {
			endpoint := fmt.Sprintf("%s/api/v2/recommendations?illness_id=%s&cabinet_id=%s",
				baseURL, url.QueryEscape(illnessID), cabinetID)
			requestJSON(t, client, http.MethodGet, endpoint, token, nil, http.StatusOK, &result)
		})

		a.Step("Assert: рекомендовано только безопасное лекарство с верной дозировкой", func(_ *allure.Context) {
			require.Len(t, result.Items, 1)
			item := result.Items[0]
			require.Equal(t, uuid.MustParse(safeMedicineID), item.MedicineID)
			require.Equal(t, cabinetID, item.CabinetID)
			require.Equal(t, safeItemID, item.CabinetItemID)
			require.Equal(t, "Driver-safe medicine", item.Name)
			require.Equal(t, float64(1), item.Dosage)
			require.Equal(t, "piece", item.UnitName)
			require.Equal(t, 2, item.FrequencyPerDay)
			require.Equal(t, float64(10), item.QuantityInCabinet)
		})
	}, allure.WithEpic("MyPills"), allure.WithFeature("E2E"), allure.WithStory("Medicine recommendation"),
		allure.WithParentSuite("MyPills"), allure.WithSuite("E2E"), allure.WithSubSuite("HTTP API"))
}

func postJSON(t *testing.T, client *http.Client, endpoint, token string, body any, status int, result any) {
	t.Helper()
	requestJSON(t, client, http.MethodPost, endpoint, token, body, status, result)
}

func requestJSON(t *testing.T, client *http.Client, method, endpoint, token string, body any, status int, result any) {
	t.Helper()
	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(encoded)
	}
	req, err := http.NewRequest(method, endpoint, reader)
	require.NoError(t, err)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	payload, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equalf(t, status, resp.StatusCode, "%s %s: %s", method, endpoint, payload)
	require.NoError(t, json.Unmarshal(payload, result), "response: %s", payload)
}
