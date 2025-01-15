package l10n

import "testing"

func expect(t *testing.T, trans Translations, lang Lang, code string, exp string) {
	val, err := trans.Translate(lang, code)
	if err != nil {
		t.Error(err)
		return
	}
	if val.Text != exp {
		t.Error("Expected:", exp, "Got:", val)
	} else {
		t.Log(val)
	}
}

func TestLocalization(t *testing.T) {
	trans, err := NewFromFile("l10n-example")
	if err != nil {
		t.Error(err)
		return
	}

	expect(t, trans, ZH, "1000", "你好")
	expect(t, trans, EN, "1000", "Hello")
	expect(t, trans, JA, "1000", "こんにちは")
	expect(t, trans, KO, "1000", "안녕하세요")

	expect(t, trans, ZH, "1001", "欢迎")
	expect(t, trans, EN, "1001", "Welcome")
	expect(t, trans, JA, "1001", "ようこそ")
	expect(t, trans, KO, "1001", "환영")
}
