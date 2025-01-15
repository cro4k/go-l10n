package l10n

import (
	"context"
	"testing"
	"time"
)

type A struct {
	Message Code
	Text    string
	Date    time.Time
}

func init() {
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	newyork, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}
	shanghai, err := time.LoadLocation("Asia/Shanghai")
	SetLocation(JP, tokyo)
	SetLocation(US, newyork)
	SetLocation(CN, shanghai)
}

func TestTranslator(t *testing.T) {
	trans, err := NewFromFile("l10n-example")
	if err != nil {
		t.Error(err)
		return
	}
	tran := NewTranslator(trans, nil)

	{
		a := &A{
			Message: "1000",
			Text:    "1000",
			Date:    time.Now(),
		}
		_ = tran.Translate(context.Background(), &AcceptLanguage{Lang: ZH, Locale: CN}, a)
		t.Log(a)
	}
	{
		a := &A{
			Message: "1001",
			Text:    "1000",
			Date:    time.Now(),
		}
		_ = tran.Translate(context.Background(), &AcceptLanguage{Lang: EN, Locale: US}, a)
		t.Log(a)
	}

	{
		a := &A{
			Message: "1001",
			Text:    "1000",
			Date:    time.Now(),
		}
		_ = tran.Translate(context.Background(), &AcceptLanguage{Lang: EN, Locale: JP}, a)
		t.Log(a)
	}
	{
		a := &A{
			Message: "1001",
			Text:    "1000",
			Date:    time.Now(),
		}
		_ = tran.Translate(context.Background(), &AcceptLanguage{Lang: JA, Locale: JP}, a)
		t.Log(a)
	}

	{
		a := &A{
			Message: "1001",
			Text:    "1000",
			Date:    time.Now(),
		}
		_ = tran.Translate(context.Background(), &AcceptLanguage{Lang: EN, Locale: KR}, a)
		t.Log(a)
	}
}
