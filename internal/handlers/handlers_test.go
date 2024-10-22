package handlers

import (
	"github.com/KirillinED/shortener/internal/storage"
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
	type want struct {
		expectedStatusCode int
		expectedBody       string
	}

	tests := []struct {
		name string
		request
		want
	}{
		{
			name:    "positive test #1",
			request: request{method: http.MethodPost, target: "/", body: strings.NewReader("https://google.com")},
			want:    want{expectedStatusCode: http.StatusCreated, expectedBody: "http://localhost:8080/1Frbb7"},
		},
		{
			name:    "long value positive test #2",
			request: request{method: http.MethodPost, target: "/", body: strings.NewReader("https://vladimir-tko.etton.ru/terSchema/?year=2024&flows=true&zoom=8&center=55.96608422809726,41.59973144531251&layers=gs,trade,transport-infrastructure,educational,household,catering,culture,admin_building,other,no_type,construction,set,uk,apartmentBuildings,ind")},
			want:    want{expectedStatusCode: http.StatusCreated, expectedBody: "http://localhost:8080/3gIQrJ"},
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

			CreateShortLinkHandler(w, req)

			assert.Equal(t, test.want.expectedStatusCode, w.Code)
			assert.Equal(t, test.want.expectedBody, w.Body.String())
		})
	}
}

func TestGetShortLinkHandler(t *testing.T) {
	type want struct {
		expectedStatusCode     int
		expectedHeaderLocation string
	}

	const (
		yandexLongLink  = "https://yandex.ru/"
		yandexShortLink = "/43BydK"

		tkoSaratovLongLink  = "https://saratov-tko.etton.ru/terSchema/?year=2024&flows=true&zoom=8&center=55.96608422809726,41.59973144531251&layers=gs,trade,transport-infrastructure,educational,household,catering,culture,admin_building,other,no_type,construction,set,uk,apartmentBuildings,ind"
		tkoSaratovShortLink = "/9qeVj"
	)

	storage.LongToShortLinksMap[yandexLongLink] = yandexShortLink[1:]
	storage.ShortToLongLinksMap[yandexShortLink[1:]] = yandexLongLink

	storage.LongToShortLinksMap[tkoSaratovLongLink] = tkoSaratovShortLink[1:]
	storage.ShortToLongLinksMap[tkoSaratovShortLink[1:]] = tkoSaratovLongLink

	ts := httptest.NewServer(getShortLinkHandlerRouter())
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
			request: request{method: http.MethodGet, target: ts.URL + yandexShortLink},
			want:    want{expectedStatusCode: http.StatusTemporaryRedirect, expectedHeaderLocation: yandexLongLink},
		},
		{
			name:    "positive test #2",
			request: request{method: http.MethodGet, target: ts.URL + tkoSaratovShortLink},
			want:    want{expectedStatusCode: http.StatusTemporaryRedirect, expectedHeaderLocation: tkoSaratovLongLink},
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
			defer resp.Body.Close()

			assert.Equal(t, test.want.expectedStatusCode, resp.StatusCode)
			assert.Equal(t, test.want.expectedHeaderLocation, resp.Header.Get("Location"))
		})
	}
}

func getShortLinkHandlerRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/{link}", GetShortLinkHandler)
	return r
}
