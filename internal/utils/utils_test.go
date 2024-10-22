package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestURL_Short(t *testing.T) {
	tests := []struct {
		name string
		url  URL
		want string
	}{
		{
			name: "positive test #1",
			url:  URL("https://www.example.com"),
			want: "61xc1",
		},
		{
			name: "positive test #2",
			url:  URL("https://practicum.yandex.ru/"),
			want: "2IYdFP",
		},
		{
			name: "zero value",
			url:  URL(""),
			want: "",
		},
		{
			name: "big value",
			url:  URL("https://vladimir-tko.etton.ru/terSchema/?year=2024&flows=true&zoom=8&center=55.96608422809726,41.59973144531251&layers=gs,trade,transport-infrastructure,educational,household,catering,culture,admin_building,other,no_type,construction,set,uk,apartmentBuildings,ind"),
			want: "3gIQrJ",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.url.Short())
		})
	}
}
