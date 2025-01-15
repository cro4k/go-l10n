# Localization(l10n) Package


## Prepare

1. make a directory named `l10n` (or you can name it as you like);
2. create language file named `<LANGUAGE_CODE>.yaml`(eg: `en.yaml`), and put in `l10n`. The language code must be abbreviated in two words. See [examples](./l10n-example) for content. 

## Usage

- basic

```go
package main

import "github.com/cro4k/go-l10n"

func main() {
    // load translation files
    trans, err := l10n.NewFromFiles("./l10n")
    if err != nil{
        panic(err)
    }
    
    // translate
    val, err := trans.Translate(l10n.EN, "1000")
}

```

- auto translate

```go
package main

import (
    "context"
    "fmt"

    "github.com/cro4k/go-l10n"
)

func main() {
    type Data struct {
        Name    string
        Message l10n.Code
    }

    // load translation files
    trans, err := l10n.NewFromFiles("./l10n")
    if err != nil {
        panic(err)
    }
    translator := l10n.NewTranslator(trans, nil)
    data := &Data{
        Name:    "Alen",
        Message: "1000",
    }
    lang := &l10n.AcceptLanguage{Lang: l10n.EN}
    _ = translator.Translate(context.Background(), lang, data)
    fmt.Println(data.Name, data.Message)
}

```

- auto convert time location
```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/cro4k/go-l10n"
)

func init() {
    // set location
    NewYork, err := time.LoadLocation("America/New_York")
    if err != nil {
        panic(err)
    }
    l10n.SetLocation(l10n.US, NewYork)
}

func main() {
    type Data struct {
        Name    string
        Message l10n.Code
        Date    time.Time
    }

    // load translation files
    trans, err := l10n.NewFromFiles("./l10n")
    if err != nil {
        panic(err)
    }
    translator := l10n.NewTranslator(trans, nil)
    data := &Data{
        Name:    "Alen",
        Message: "1000",     // The code will be translated to target language
        Date:    time.Now(), // The time will be replaced with location
    }
    lang := &l10n.AcceptLanguage{Lang: l10n.EN, Locale: l10n.US}
    _ = translator.Translate(context.Background(), lang, data)
    fmt.Println(data.Name, data.Message, data.Date)
}


```

## Reference

- Country Code Standard: [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2)
- Language Code Standard: [ISO 639-1](https://en.wikipedia.org/wiki/ISO_639-1)