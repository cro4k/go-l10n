package l10n

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const (
	AcceptLanguageHeader = "Accept-Language"
)

var (
	DefaultLang   = ZH
	DefaultLocale = CN
)

type AcceptLanguage struct {
	Lang   Lang
	Locale Locale
	Q      float32
}

func (a *AcceptLanguage) String() string {
	return fmt.Sprintf("%s-%s", a.Lang, a.Locale)
}

type AcceptLanguages []*AcceptLanguage

func (a AcceptLanguages) Len() int           { return len(a) }
func (a AcceptLanguages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AcceptLanguages) Less(i, j int) bool { return a[i].Q > a[j].Q }

func (a AcceptLanguages) One() *AcceptLanguage {
	if len(a) > 0 {
		return a[0]
	}
	return &AcceptLanguage{Lang: DefaultLang, Locale: DefaultLocale}
}

func DecodeAcceptLanguages(h http.Header) AcceptLanguages {
	header := h.Get(AcceptLanguageHeader)
	if header == "" {
		return AcceptLanguages{}
	}
	fields := strings.Split(header, ",")
	languages := make(AcceptLanguages, 0, len(fields))
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		if lang, err := decodeLanguage(field); err == nil {
			languages = append(languages, lang)
		}
	}
	sort.Sort(languages)
	return languages
}

func decodeLanguage(s string) (*AcceptLanguage, error) {
	f := strings.Split(s, ";")
	langWithLocale := f[0]
	q := float64(1)
	if len(f) > 1 {
		q, _ = strconv.ParseFloat(f[1], 32)
	}
	var lang Lang
	var locale Locale
	if len(langWithLocale) == 2 {
		lang = Lang(langWithLocale)
	} else if len(langWithLocale) == 5 {
		lang = Lang(strings.ToLower(langWithLocale[:2]))
		locale = Locale(strings.ToUpper(langWithLocale[3:]))
	} else {
		return nil, errors.New("invalid language")
	}
	return &AcceptLanguage{
		Lang:   lang,
		Locale: locale,
		Q:      float32(q),
	}, nil
}

type acceptLanguageContextKey struct{}

func WithContext(ctx context.Context, languages AcceptLanguages) context.Context {
	return context.WithValue(ctx, acceptLanguageContextKey{}, languages)
}

func WithHTTPRequest(request *http.Request) *http.Request {
	ctx := request.Context()
	ctx = WithContext(ctx, DecodeAcceptLanguages(request.Header))
	return request.WithContext(ctx)
}

func FromContext(ctx context.Context) AcceptLanguages {
	val := ctx.Value(acceptLanguageContextKey{})
	if val == nil {
		return nil
	}
	return val.(AcceptLanguages)
}
