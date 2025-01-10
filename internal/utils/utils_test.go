package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShortURL(t *testing.T) {
	type args struct {
		url string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive test #1",
			args: args{url: "https://www.example.com"},
			want: "61xc1",
		},
		{
			name: "positive test #2",
			args: args{url: "https://practicum.yandex.ru/"},
			want: "2IYdFP",
		},
		{
			name: "zero value",
			args: args{url: ""},
			want: "",
		},
		{
			name: "big value",
			args: args{url: "https://vladimir-tko.etton.ru/terSchema/?year=2024&flows=true&zoom=8&center=55.96608422809726,41.59973144531251&layers=gs,trade,transport-infrastructure,educational,household,catering,culture,admin_building,other,no_type,construction,set,uk,apartmentBuildings,ind"},
			want: "3gIQrJ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ShortURL(tt.args.url), "ShortURL(%v)", tt.args.url)
		})
	}
}
