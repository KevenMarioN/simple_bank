package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	mockdb "github.com/KevenMarioN/simple_bank/db/mock"
	db "github.com/KevenMarioN/simple_bank/db/sqlc"
	"github.com/KevenMarioN/simple_bank/util"
)

func TestTransfer(t *testing.T) {
	transfer := db.TransferTxParams{
		FromAccountID: util.RandomInt(1, 1000),
		ToAccountID:   util.RandomInt(1, 1000),
		Amount:        util.RandomMoney(),
	}

	testCases := []struct {
		name          string
		transfer      db.TransferTxParams
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			transfer: transfer,
			buildStubs: func(store *mockdb.MockStore) {
				transferMock := db.Transfer{
					ID:            util.RandomInt(1, 1000),
					FromAccountID: transfer.FromAccountID,
					ToAccountID:   transfer.ToAccountID,
					Amount:        transfer.Amount,
					CreatedAt:     time.Now(),
				}
				store.EXPECT().
					TrasferTx(gomock.Any(), transfer).
					Times(1).
					Return(transferMock, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)

			recorder := httptest.NewRecorder()
			data, err := json.Marshal(tc.transfer)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/tranfers", bytes.NewBuffer(data))

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
