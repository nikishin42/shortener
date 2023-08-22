package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nikishin42/shortener/cmd/shortener/config"
	"github.com/nikishin42/shortener/cmd/shortener/constants"
	"github.com/nikishin42/shortener/cmd/shortener/pkg/abbreviator"
	"github.com/nikishin42/shortener/cmd/shortener/pkg/storage"
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
		setup func(abbreviator *abbreviator.MockAbbreviatorI, storage *storage.MockStorageI)
		exp   expexted
	}{
		{
			name: "wrong method error",
			args: args{
				method: http.MethodGet,
				body:   strings.NewReader("https://music.yandex.ru/"),
			},
			setup: func(abbreviator *abbreviator.MockAbbreviatorI, storage *storage.MockStorageI) {},
			exp: expexted{
				status: http.StatusMethodNotAllowed,
				body:   "",
			},
		},
		{
			name: "query not empty error",
			args: args{
				method: http.MethodPost,
				body:   strings.NewReader("https://music.yandex.ru/"),
				query:  "not_empty",
			},
			setup: func(abbreviator *abbreviator.MockAbbreviatorI, storage *storage.MockStorageI) {},
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
			setup: func(abbreviator *abbreviator.MockAbbreviatorI, storage *storage.MockStorageI) {},
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
			setup: func(abbreviator *abbreviator.MockAbbreviatorI, storage *storage.MockStorageI) {},
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
			setup: func(abbreviator *abbreviator.MockAbbreviatorI, storage *storage.MockStorageI) {
				storage.EXPECT().GetID("https://music.yandex.ru/").Return("", false)
				abbreviator.EXPECT().CreateID([]byte("https://music.yandex.ru/")).Return(constants.HTTPHostPrefix+"/"+"Fy", nil)
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
			setup: func(abbreviator *abbreviator.MockAbbreviatorI, storage *storage.MockStorageI) {
				storage.EXPECT().GetID("https://music.yandex.ru/").Return("", false)
				abbreviator.EXPECT().CreateID([]byte("https://music.yandex.ru/")).Return("", assert.AnError)
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
			setup: func(abbreviator *abbreviator.MockAbbreviatorI, storage *storage.MockStorageI) {
				storage.EXPECT().GetID("https://music.yandex.ru/").Return("", false)
				abbreviator.EXPECT().CreateID([]byte("https://music.yandex.ru/")).Return(constants.HTTPHostPrefix+"/"+"Fy", nil)
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
			setup: func(abbreviator *abbreviator.MockAbbreviatorI, storage *storage.MockStorageI) {
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
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStorage := storage.NewMockStorageI(ctrl)
			mockAbbreviator := abbreviator.NewMockAbbreviatorI(ctrl)
			tc.setup(mockAbbreviator, mockStorage)

			a := New(&config.Config{
				Address:              constants.HTTPHostPrefix,
				BaseShortenerAddress: constants.HTTPHostPrefix,
			}, mockStorage, mockAbbreviator)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.args.method, a.Config.Address+"/"+tc.args.query, tc.args.body)
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
		query  string
	}
	type expexted struct {
		status int
		body   string
	}
	tests := []struct {
		name  string
		args  args
		setup func(storage *storage.MockStorageI)
		exp   expexted
	}{
		{
			name: "wrong method",
			args: args{
				method: http.MethodPost,
				query:  "",
			},
			setup: func(storage *storage.MockStorageI) {},
			exp: expexted{
				status: http.StatusMethodNotAllowed,
				body:   "",
			},
		},
		{
			name: "empty query error",
			args: args{
				method: http.MethodGet,
				query:  "",
			},
			setup: func(storage *storage.MockStorageI) {
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
				query:  "Fy",
			},
			setup: func(storage *storage.MockStorageI) {
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
				query:  "Fy",
			},
			setup: func(storage *storage.MockStorageI) {
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
			mockStorage := storage.NewMockStorageI(ctrl)
			mockAbbreviator := abbreviator.NewMockAbbreviatorI(ctrl)
			tc.setup(mockStorage)

			a := New(&config.Config{
				Address:              constants.HTTPHostPrefix,
				BaseShortenerAddress: constants.HTTPHostPrefix,
			}, mockStorage, mockAbbreviator)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.args.method, a.Config.Address+"/"+tc.args.query, nil)
			a.Redirect(w, r)

			assert.Equal(t, tc.exp.status, w.Code)
			assert.Equal(t, tc.exp.body, w.Body.String())
		})
	}
}
