package l10n

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

const (
	codePkg  = "github.com/cro4k/go-l10n"
	timePkg  = "time"
	codeName = "Code"
	timeName = "Time"
)

type (
	Code string

	Translator struct {
		translations Translations
		errorHandler func(ctx context.Context, code Code, lang Lang, err error) *string
	}
)

func NewTranslator(translations Translations, errorHandler func(ctx context.Context, code Code, lang Lang, err error) *string) *Translator {
	if errorHandler == nil {
		errorHandler = func(ctx context.Context, code Code, lang Lang, err error) *string { return nil }
	}
	return &Translator{
		translations: translations,
		errorHandler: errorHandler,
	}
}

func (t *Translator) Translate(ctx context.Context, lang *AcceptLanguage, v any) error {
	ty := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if ty.Kind() != reflect.Ptr {
		return fmt.Errorf("cannot set translate results into unaddressable value")
	}
	return t.translate(ctx, lang, val)
}

func (t *Translator) SupportedLanguages() []Lang {
	languages := make([]Lang, 0, len(t.translations))
	for lang := range t.translations {
		languages = append(languages, lang)
	}
	return languages
}

func (t *Translator) translate(ctx context.Context, lang *AcceptLanguage, val reflect.Value) error {
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Struct:
		ty := val.Type()
		if ty.PkgPath() == timePkg && ty.Name() == timeName {
			loc := lang.Locale.Location()
			if loc != nil {
				value := val.Interface().(time.Time)
				val.Set(reflect.ValueOf(value.In(loc)))
			}
			return nil
		}
		fallthrough
	case reflect.Map:
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if err := t.translate(ctx, lang, field); err != nil {
				return err
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			if err := t.translate(ctx, lang, val.Index(i)); err != nil {
				return err
			}
		}
	case reflect.String:
		ty := val.Type()
		if ty.PkgPath() == codePkg && ty.Name() == codeName {
			code := val.String()
			if text, err := t.translations.Translate(lang.Lang, code); err != nil {
				s := t.errorHandler(ctx, Code(code), lang.Lang, err)
				if s != nil {
					val.SetString(*s)
				}
			} else {
				val.SetString(text.String())
			}
		}
	default:
	}
	return nil
}
