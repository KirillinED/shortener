package middlewares

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"github.com/KirillinED/shortener/internal/config"
	"github.com/KirillinED/shortener/internal/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewCompressWriter(t *testing.T) {
	type args struct {
		encoding string
	}
	tests := []struct {
		name   string
		want   io.WriteCloser
		args   args
		hasErr bool
	}{
		{name: "gzip writer", want: &gzip.Writer{}, args: args{encoding: encodingGzip}, hasErr: false},
		{name: "flate writer", want: &flate.Writer{}, args: args{encoding: encodingDeflate}, hasErr: false},
		{name: "wildcard default writer", want: &gzip.Writer{}, args: args{encoding: encodingWildcard}, hasErr: false},
		{name: "empty encoding", want: nil, args: args{encoding: ""}, hasErr: true},
		{name: "unsupported encoding", want: nil, args: args{encoding: "br"}, hasErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := newCompressWriter(httptest.NewRecorder(), tt.args.encoding)
			if tt.hasErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.IsType(t, tt.want, w.cw)
			require.NoError(t, w.Close())
		})
	}
}

func TestNewCompressReader(t *testing.T) {
	originalData := `https://example.com/?&africa=alert&energy=unaccountable&opportunist=amused&canteen=hard&standoff=Early&neon=melodic&crew=periodic&height=boiling&stallion=frail&pendulum=many&century=hellish&cork=classy&button=exuberant&dory=goofy&atrium=defiant&niece=yellow&clock=fragile&learning=illegal&lye=kindhearted&development=obnoxious&pie=idiotic&curio=loud&magician=fearless&whorl=majestic&ghost=amuck&clutch=penitent&straw=soggy&detention=foamy&footstool=sulky&character=knowledgeable&motorcar=angry&cloister=steadfast&switchboard=jobless&smell=smoggy&euphonium=curious&parrot=accurate&presence=enchanting&pounding=rabid&snob=cute&tote=real&february=sparkling&eyebrow=jittery&tuesday=majestic&strip=courageous&dryer=loutish&heat=evil&senator=panoramic&crystallography=torpid&guide=quick&bun=guiltless&prose=adventurous&caution=chubby&feet=weary&vertigo=joyous&amusement=gaudy&creative=ad hoc&specific=frantic&expansion=wacky&tie=swift&complex=womanly&enemy=ambiguous&marimba=misty&zither=sore&orange=afraid&lyric=melted&glen=discreet&triangle=nauseating&slime=ripe&carol=belligerent&personality=parsimonious&tripod=cloistered&jennifer=imported&fascia=dusty&comb=burly&boar=fretful&spring=sad&lead=wandering&flood=finicky&chance=lewd&swimming=makeshift&premier=drunk&forum=mindless&mattock=obscene&pitching=jumbled&establishment=oceanic&final=lucky&craw=apathetic&laptop=symptomatic&digital=wicked&stop=boundless&hippopotamus=vacuous&sprout=kindhearted&cabin=idiotic&ruffle=skinny&pupil=clean&freeplay=tightfisted&banana=unequaled&baboon=aberrant&fir=exotic&reflection=ill&lotion=calm&series=fresh&meet=used&wednesday=tasteless&forum=gruesome&himalayan=greedy&loan=stimulating&planter=silky&average=changeable&sushi=sloppy&inside=lowly&ripple=tiresome&solitaire=muddled&belly=offbeat&spectacles=sleepy&workbench=giddy&pizza=sore&summer=capricious&latency=jittery&faucet=frantic&pickax=dashing&sarong=selective&concert=spiritual&coevolution=sneaky&notebook=wanting&breast=uneven&synergy=magenta`

	type args struct {
		encoding string
	}
	tests := []struct {
		name   string
		args   args
		hasErr bool
	}{
		{name: "gzip reader", args: args{encoding: encodingGzip}, hasErr: false},
		{name: "flate reader", args: args{encoding: encodingDeflate}, hasErr: false},
		{name: "empty encoding", args: args{encoding: ""}, hasErr: true},
		{name: "unsupported encoding", args: args{encoding: "br"}, hasErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := newCompressWriter(httptest.NewRecorder(), tt.args.encoding)
			if tt.hasErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			writeBytes, err := w.Write([]byte(originalData))
			require.NoError(t, err)
			require.NoError(t, w.Close())
			rr, ok := w.w.(*httptest.ResponseRecorder)
			require.True(t, ok)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(rr.Body.Bytes()))
			reader, err := newCompressReader(req.Body, tt.args.encoding)
			if tt.hasErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			var b bytes.Buffer
			readBytes, err := b.ReadFrom(reader)
			require.NoError(t, err)
			assert.Equal(t, writeBytes, int(readBytes))
			require.NoError(t, reader.Close())
		})
	}
}

func TestCompressMiddleware(t *testing.T) {
	baseUrl := config.GetConfig().BaseURL

	var requestBody = `{"url":"https://example.com/?&africa=alert&energy=unaccountable&opportunist=amused&canteen=hard&standoff=Early&neon=melodic&crew=periodic&height=boiling&stallion=frail&pendulum=many&century=hellish&cork=classy&button=exuberant&dory=goofy&atrium=defiant&niece=yellow&clock=fragile&learning=illegal&lye=kindhearted&development=obnoxious&pie=idiotic&curio=loud&magician=fearless&whorl=majestic&ghost=amuck&clutch=penitent&straw=soggy&detention=foamy&footstool=sulky&character=knowledgeable&motorcar=angry&cloister=steadfast&switchboard=jobless&smell=smoggy&euphonium=curious&parrot=accurate&presence=enchanting&pounding=rabid&snob=cute&tote=real&february=sparkling&eyebrow=jittery&tuesday=majestic&strip=courageous&dryer=loutish&heat=evil&senator=panoramic&crystallography=torpid&guide=quick&bun=guiltless&prose=adventurous&caution=chubby&feet=weary&vertigo=joyous&amusement=gaudy&creative=ad hoc&specific=frantic&expansion=wacky&tie=swift&complex=womanly&enemy=ambiguous&marimba=misty&zither=sore&orange=afraid&lyric=melted&glen=discreet&triangle=nauseating&slime=ripe&carol=belligerent&personality=parsimonious&tripod=cloistered&jennifer=imported&fascia=dusty&comb=burly&boar=fretful&spring=sad&lead=wandering&flood=finicky&chance=lewd&swimming=makeshift&premier=drunk&forum=mindless&mattock=obscene&pitching=jumbled&establishment=oceanic&final=lucky&craw=apathetic&laptop=symptomatic&digital=wicked&stop=boundless&hippopotamus=vacuous&sprout=kindhearted&cabin=idiotic&ruffle=skinny&pupil=clean&freeplay=tightfisted&banana=unequaled&baboon=aberrant&fir=exotic&reflection=ill&lotion=calm&series=fresh&meet=used&wednesday=tasteless&forum=gruesome&himalayan=greedy&loan=stimulating&planter=silky&average=changeable&sushi=sloppy&inside=lowly&ripple=tiresome&solitaire=muddled&belly=offbeat&spectacles=sleepy&workbench=giddy&pizza=sore&summer=capricious&latency=jittery&faucet=frantic&pickax=dashing&sarong=selective&concert=spiritual&coevolution=sneaky&notebook=wanting&breast=uneven&synergy=magenta"}`
	var responseBody = fmt.Sprintf(`{"result":"%s/32DOl6"}`, baseUrl[:len(baseUrl)-1])

	type args struct {
		headers map[string]string
	}
	tests := []struct {
		name                 string
		args                 args
		hasCompressWriterErr bool
		hasCompressReaderErr bool
		want                 string
	}{
		{
			name: "positive test gzip encode/decode",
			want: responseBody,
			args: args{
				headers: map[string]string{
					"Content-Type":     "application/json",
					"Content-Encoding": encodingGzip,
					"Accept-Encoding":  encodingGzip,
				}}},
		{
			name: "positive test deflate encode/decode",
			want: responseBody,
			args: args{
				headers: map[string]string{
					"Content-Type":     "application/json",
					"Content-Encoding": encodingDeflate,
					"Accept-Encoding":  encodingDeflate,
				}}},
		{
			name: "positive test wildcard encode/decode",
			want: responseBody,
			args: args{
				headers: map[string]string{
					"Content-Type":     "application/json",
					"Content-Encoding": encodingDeflate,
					"Accept-Encoding":  encodingWildcard,
				}}},
		{
			name: "positive test without Content-Encoding",
			want: responseBody,
			args: args{
				headers: map[string]string{
					"Content-Type":    "application/json",
					"Accept-Encoding": encodingWildcard,
				}}},
		{
			name:                 "negative test unsupported Accept-Encoding",
			hasCompressWriterErr: true,
			want:                 "Accept-Encoding type is not supported",
			args: args{
				headers: map[string]string{
					"Content-Type":     "application/json",
					"Content-Encoding": "unsupported encoding",
					"Accept-Encoding":  encodingWildcard,
				}}},
		{
			name:                 "negative test unsupported Content-Encoding",
			hasCompressReaderErr: true,
			want:                 "Content-Encoding type is not supported",
			args: args{
				headers: map[string]string{
					"Content-Type":     "application/json",
					"Content-Encoding": encodingGzip,
					"Accept-Encoding":  "unsupported encoding",
				}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(requestBody)))
			for k, v := range tt.args.headers {
				request.Header.Add(k, v)
			}

			// Если нужно сжать тело запроса, то сжимаем
			if request.Header.Get("Content-Encoding") != "" {
				rr := httptest.NewRecorder()
				cw, err := newCompressWriter(rr, request.Header.Get("Content-Encoding"))
				if tt.hasCompressWriterErr {
					require.Error(t, err)
					require.Equal(t, err.Error(), tt.want)
					return
				}
				require.NoError(t, err)

				_, err = cw.Write([]byte(requestBody))
				require.NoError(t, err)

				var ok bool
				rr, ok = cw.w.(*httptest.ResponseRecorder)
				require.True(t, ok)
				require.NoError(t, cw.Close())
				require.NoError(t, request.Body.Close())
				request.Body = &TestReadCloser{r: rr.Body}
			}

			handler := Compress(http.HandlerFunc(handlers.CreateShortLinkHandler))
			handler.ServeHTTP(response, request)

			// Если ответ закодирован, то раскодируем для проверки
			if request.Header.Get("Accept-Encoding") != "" {
				acceptEncoding := request.Header.Get("Accept-Encoding")
				if acceptEncoding == encodingWildcard {
					acceptEncoding = encodingGzip
				}

				rc := &TestReadCloser{r: response.Body}
				rr, err := newCompressReader(rc, acceptEncoding)
				if tt.hasCompressReaderErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)

				var b bytes.Buffer
				_, err = b.ReadFrom(rr)
				require.NoError(t, err)
				require.NoError(t, rr.Close())
				response.Body = &b
			}

			assert.Equal(t, tt.want, response.Body.String())
			require.NoError(t, request.Body.Close())
		})
	}
}

type TestReadCloser struct {
	r io.Reader
}

func (rc *TestReadCloser) Read(p []byte) (n int, err error) {
	return rc.r.Read(p)
}

func (rc *TestReadCloser) Close() error {
	return nil
}
