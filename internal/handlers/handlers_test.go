package handlers

import (
	"fmt"
	"github.com/KirillinED/shortener/internal/dto"
	"github.com/KirillinED/shortener/internal/foundation"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type request struct {
	method string
	target string
	body   io.Reader
}

func TestCreateShortLinkHandler(t *testing.T) {
	app := foundation.NewAppStub()
	defer app.Shutdown()

	type want struct {
		expectedStatusCode int
		expectedBody       string
	}

	baseUrl := app.GetConfig().BaseURL

	tests := []struct {
		name string
		request
		want
	}{
		{
			name:    "positive test #1",
			request: request{method: http.MethodPost, target: "/", body: strings.NewReader(`{"url":"https://google.com"}`)},
			want:    want{expectedStatusCode: http.StatusCreated, expectedBody: fmt.Sprintf(`{"result":"%s/1Frbb7"}`, baseUrl[:len(baseUrl)-1])},
		},
		{
			name:    "long value positive test #2",
			request: request{method: http.MethodPost, target: "/", body: strings.NewReader(`{"url":"https://example.com/terSchema/?year=2024&flows=true&zoom=8&center=55.96608422809726,41.59973144531251&layers=gs,trade,transport-infrastructure,educational,household,catering,culture,admin_building,other,no_type,construction,set,uk,apartmentBuildings,ind"}`)},
			want:    want{expectedStatusCode: http.StatusCreated, expectedBody: fmt.Sprintf(`{"result":"%s/3LyLZS"}`, baseUrl[:len(baseUrl)-1])},
		},
		{
			name:    "zero body negative test #3",
			request: request{method: http.MethodPost, target: "/", body: nil},
			want:    want{expectedStatusCode: http.StatusBadRequest, expectedBody: "body cannot be empty\n"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.target, test.body)
			w := httptest.NewRecorder()

			CreateShortLinkHandler(app)(w, req)

			assert.Equal(t, test.want.expectedStatusCode, w.Code)
			assert.Equal(t, test.want.expectedBody, w.Body.String())
		})
	}
}

func TestGetShortLinkHandler(t *testing.T) {
	app := foundation.NewAppStub()
	defer app.Shutdown()

	type want struct {
		expectedStatusCode     int
		expectedHeaderLocation string
	}

	links := []dto.Link{
		{
			Long:  "https://yandex.ru/",
			Short: "43BydK",
		},
		{
			Long:  "https://example.com/terSchema/?year=2024&flows=true&zoom=8&center=55.96608422809726,41.59973144531251&layers=gs,trade,transport-infrastructure,educational,household,catering,culture,admin_building,other,no_type,construction,set,uk,apartmentBuildings,ind",
			Short: "3LyLZS",
		},
	}

	for _, link := range links {
		err := app.MemoryStorage.StoreLink(link)
		require.NoError(t, err)
	}

	ts := httptest.NewServer(getShortLinkHandlerRouter(app))
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer ts.Close()

	tests := []struct {
		name string
		request
		want
	}{
		{
			name:    "positive test #1",
			request: request{method: http.MethodGet, target: ts.URL + "/" + links[0].Short},
			want:    want{expectedStatusCode: http.StatusTemporaryRedirect, expectedHeaderLocation: links[0].Long},
		},
		{
			name:    "positive test #2",
			request: request{method: http.MethodGet, target: ts.URL + "/" + links[1].Short},
			want:    want{expectedStatusCode: http.StatusTemporaryRedirect, expectedHeaderLocation: links[1].Long},
		},
		{
			name:    "not found negative test",
			request: request{method: http.MethodGet, target: ts.URL + "/qwerty"},
			want:    want{expectedStatusCode: http.StatusNotFound, expectedHeaderLocation: ""},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, err := http.NewRequest(test.request.method, test.request.target, test.request.body)
			require.NoError(t, err)

			resp, err := ts.Client().Do(r)
			require.NoError(t, err)
			require.NoError(t, resp.Body.Close())

			assert.Equal(t, test.want.expectedStatusCode, resp.StatusCode)
			assert.Equal(t, test.want.expectedHeaderLocation, resp.Header.Get("Location"))
		})
	}
}

func getShortLinkHandlerRouter(app foundation.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/{link}", GetShortLinkHandler(app))
	return r
}
