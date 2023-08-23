package servicelayer

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	_ "github.com/gorilla/mux"

	"github.com/nikishin42/shortener/cmd/shortener/config"
	"github.com/nikishin42/shortener/cmd/shortener/constants"
	"github.com/nikishin42/shortener/cmd/shortener/interfaces"
)

type errReader int

func (e errReader) Read(_ []byte) (n int, err error) {
	return 0, assert.AnError
}

func TestServer_Homepage(t *testing.T) {
	t.Parallel()

	type args struct {
		method string
		query  string
		body   io.Reader
	}

	type expexted struct {
		status int
		body   string
	}
	tests := []struct {
		name  string
		args  args
		setup func(abbreviator *interfaces.MockCreatorID, storage *interfaces.MockStorage)
		exp   expexted
	}{
		{
			name: "id not empty error",
			args: args{
				method: http.MethodPost,
				body:   strings.NewReader("https://music.yandex.ru/"),
				query:  "not_empty",
			},
			setup: func(abbreviator *interfaces.MockCreatorID, storage *interfaces.MockStorage) {},
			exp: expexted{
				status: http.StatusBadRequest,
				body:   "",
			},
		},
		{
			name: "read body error",
			args: args{
				method: http.MethodPost,
				body:   errReader(0),
			},
			setup: func(abbreviator *interfaces.MockCreatorID, storage *interfaces.MockStorage) {},
			exp: expexted{
				status: http.StatusBadRequest,
				body:   "",
			},
		},
		{
			name: "not URL in body error",
			args: args{
				method: http.MethodPost,
				body:   strings.NewReader("literally not URL"),
			},
			setup: func(abbreviator *interfaces.MockCreatorID, storage *interfaces.MockStorage) {},
			exp: expexted{
				status: http.StatusBadRequest,
				body:   "",
			},
		},
		{
			name: "ok, id created",
			args: args{
				method: http.MethodPost,
				body:   strings.NewReader("https://music.yandex.ru/"),
			},
			setup: func(abbreviator *interfaces.MockCreatorID, storage *interfaces.MockStorage) {
				storage.EXPECT().GetID("https://music.yandex.ru/").Return("", false)
				abbreviator.EXPECT().CreateID([]byte("https://music.yandex.ru/"), constants.HTTPHostPrefix).Return(constants.HTTPHostPrefix+"/"+"Fy", nil)
				storage.EXPECT().SetPair(constants.HTTPHostPrefix+"/"+"Fy", "https://music.yandex.ru/").Return(nil)
			},
			exp: expexted{
				status: http.StatusCreated,
				body:   constants.HTTPHostPrefix + "/" + "Fy",
			},
		},
		{
			name: "create ID error",
			args: args{
				method: http.MethodPost,
				body:   strings.NewReader("https://music.yandex.ru/"),
			},
			setup: func(abbreviator *interfaces.MockCreatorID, storage *interfaces.MockStorage) {
				storage.EXPECT().GetID("https://music.yandex.ru/").Return("", false)
				abbreviator.EXPECT().CreateID([]byte("https://music.yandex.ru/"), constants.HTTPHostPrefix).Return("", assert.AnError)
			},
			exp: expexted{
				status: http.StatusInternalServerError,
				body:   "",
			},
		},
		{
			name: "set pair error",
			args: args{
				method: http.MethodPost,
				body:   strings.NewReader("https://music.yandex.ru/"),
			},
			setup: func(abbreviator *interfaces.MockCreatorID, storage *interfaces.MockStorage) {
				storage.EXPECT().GetID("https://music.yandex.ru/").Return("", false)
				abbreviator.EXPECT().CreateID([]byte("https://music.yandex.ru/"), constants.HTTPHostPrefix).Return(constants.HTTPHostPrefix+"/"+"Fy", nil)
				storage.EXPECT().SetPair(constants.HTTPHostPrefix+"/"+"Fy", "https://music.yandex.ru/").Return(assert.AnError)
			},
			exp: expexted{
				status: http.StatusInternalServerError,
				body:   "",
			},
		},
		{
			name: "ok, id found",
			args: args{
				method: http.MethodPost,
				body:   strings.NewReader("https://music.yandex.ru/"),
			},
			setup: func(abbreviator *interfaces.MockCreatorID, storage *interfaces.MockStorage) {
				storage.EXPECT().GetID("https://music.yandex.ru/").Return(constants.HTTPHostPrefix+"/"+"Fy", true)
			},
			exp: expexted{
				status: http.StatusOK,
				body:   constants.HTTPHostPrefix + "/" + "Fy",
			},
		},
	}
	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			//t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStorage := interfaces.NewMockStorage(ctrl)
			mockAbbreviator := interfaces.NewMockCreatorID(ctrl)
			tc.setup(mockAbbreviator, mockStorage)

			a := New(&config.Config{
				Address:              constants.DefaultHost,
				BaseShortenerAddress: constants.HTTPHostPrefix,
			}, mockStorage, mockAbbreviator)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.args.method, "/"+tc.args.query, tc.args.body)
			a.Homepage(w, r)

			assert.Equal(t, tc.exp.status, w.Code)
			assert.Equal(t, tc.exp.body, w.Body.String())
		})
	}
}

func TestServer_Redirect(t *testing.T) {
	t.Parallel()

	type args struct {
		method string
		id     string
	}
	type expexted struct {
		status int
		body   string
	}
	tests := []struct {
		name  string
		args  args
		setup func(storage *interfaces.MockStorage)
		exp   expexted
	}{
		{
			name: "empty id error",
			args: args{
				method: http.MethodGet,
				id:     "",
			},
			setup: func(storage *interfaces.MockStorage) {
			},
			exp: expexted{
				status: http.StatusBadRequest,
				body:   "",
			},
		},
		{
			name: "full URL not found error",
			args: args{
				method: http.MethodGet,
				id:     "Fy",
			},
			setup: func(storage *interfaces.MockStorage) {
				storage.EXPECT().GetFullURL(constants.HTTPHostPrefix+"/"+"Fy").Return("", false)
			},
			exp: expexted{
				status: http.StatusBadRequest,
				body:   "",
			},
		},
		{
			name: "ok",
			args: args{
				method: http.MethodGet,
				id:     "Fy",
			},
			setup: func(storage *interfaces.MockStorage) {
				storage.EXPECT().GetFullURL(constants.HTTPHostPrefix+"/"+"Fy").Return("https://music.yandex.ru/", true)
			},
			exp: expexted{
				status: http.StatusTemporaryRedirect,
				body:   "",
			},
		},
	}
	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStorage := interfaces.NewMockStorage(ctrl)
			mockAbbreviator := interfaces.NewMockCreatorID(ctrl)
			tc.setup(mockStorage)

			a := New(&config.Config{
				Address:              constants.DefaultHost,
				BaseShortenerAddress: constants.HTTPHostPrefix,
			}, mockStorage, mockAbbreviator)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.args.method, "/"+tc.args.id, nil)
			a.Redirect(w, r)

			assert.Equal(t, tc.exp.status, w.Code)
			assert.Equal(t, tc.exp.body, w.Body.String())
		})
	}
}
