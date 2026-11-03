package res

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet"
)

func TestCabinetResponseMatchesV2Contract(t *testing.T) {
	item := cabinet.Cabinet{ID: uuid.New(), Name: "Домашняя аптечка"}
	data, err := json.Marshal(NewCabinetsPageDto(cabinet.Page{
		Items: []cabinet.Cabinet{item}, Total: 1, Limit: 20, Offset: 0,
	}))
	require.NoError(t, err)
	var response map[string]any
	require.NoError(t, json.Unmarshal(data, &response))
	for _, field := range []string{"items", "total", "limit", "offset"} {
		require.Contains(t, response, field)
	}
	items := response["items"].([]any)
	require.Len(t, items, 1)
	first := items[0].(map[string]any)
	require.Equal(t, item.ID.String(), first["id"])
	require.Equal(t, item.Name, first["name"])
}
