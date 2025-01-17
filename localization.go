package l10n

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sigs.k8s.io/yaml"
)

type (
	Lang   string
	Locale string
)

const (
	ZH Lang = "zh"
	EN Lang = "en"
	JA Lang = "ja"
	KO Lang = "ko"
)

const (
	UNKNOWN Locale = ""
	CN      Locale = "CN"
	US      Locale = "US"
	UK      Locale = "UK"
	JP      Locale = "JP"
	KR      Locale = "KR"
)

const (
	DefaultL10NPath = "l10n"
)

var (
	locations = map[Locale]*time.Location{
		UNKNOWN: time.UTC,
	}
)

func SetLocation(locale Locale, location *time.Location) {
	locations[locale] = location
}

func (l Locale) Location() *time.Location {
	if loc, ok := locations[l]; ok {
		return loc
	}
	return time.Local
}

type Text struct {
	Code string `json:"code" yaml:"code"`
	Text string `json:"text" yaml:"text"`
	Lang Lang   `json:"lang" yaml:"lang"`
}

func (t *Text) String() string {
	return t.Text
}

type Metadata map[string]string

type Translation struct {
	Metadata Metadata `json:"metadata" yaml:"metadata"`
	Items    []*Text  `json:"items"    yaml:"items"`

	m map[string]*Text
}

func (t *Translation) Rebuild() {
	if len(t.m) == len(t.Items) {
		return
	}
	t.m = make(map[string]*Text)
	for _, item := range t.Items {
		t.m[item.Code] = item
	}
}

type Translations map[Lang]*Translation

func (t Translations) Translate(lang Lang, code string) (*Text, error) {
	tran, ok := t[lang]
	if !ok {
		return nil, fmt.Errorf("l10n: no translation found for language %s", lang)
	}
	text, ok := tran.m[code]
	if !ok {
		return nil, fmt.Errorf("l10n: no translation found for code %s in language %s", code, lang)
	}
	result := new(Text)
	*result = *text
	result.Lang = lang
	return result, nil
}

func Unmarshal(data []byte) (*Translation, error) {
	tran := new(Translation)
	err := yaml.Unmarshal(data, &tran)
	if err != nil {
		return nil, err
	}
	tran.Rebuild()
	return tran, nil
}

func LangFromFilename(filename string) (Lang, error) {
	if name := filepath.Base(filename); len(name) != 7 || name[2:] != ".yaml" {
		return "", errors.New("invalid language file")
	} else {
		return Lang(strings.ToLower(name[:2])), nil
	}
}

func UnmarshalFromFile(file io.ReadCloser) (*Translation, error) {
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return Unmarshal(data)
}

// NewFromFiles decode translations config form file, only support yaml format in default.
func NewFromFiles(path string) (Translations, error) {
	if path == "" {
		path = DefaultL10NPath
	}
	trans := make(Translations)
	err := filepath.Walk(path, func(filename string, info fs.FileInfo, err error) error {
		if path == filename {
			return nil
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		lang, err := LangFromFilename(filename)
		if err != nil {
			return nil
		}
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		tran, err := UnmarshalFromFile(file)
		if err != nil {
			return err
		}
		trans[lang] = tran
		return nil
	})
	return trans, err
}

func NewFromFS(fs fs.FS, files ...string) (Translations, error) {
	trans := make(Translations)
	for _, filename := range files {
		file, err := fs.Open(filename)
		if err != nil {
			return nil, err
		}
		lang, err := LangFromFilename(filename)
		if err != nil {
			return nil, err
		}
		tran, err := UnmarshalFromFile(file)
		trans[lang] = tran
	}
	return trans, nil
}
