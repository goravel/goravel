package translation

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/cast"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/goravel/framework/contracts/http"
	logcontract "github.com/goravel/framework/contracts/log"
	translationcontract "github.com/goravel/framework/contracts/translation"
	"github.com/goravel/framework/errors"
)

type Translator struct {
	ctx        context.Context
	fsLoader   translationcontract.Loader
	fileLoader translationcontract.Loader
	logger     logcontract.Log
	selector   *MessageSelector
	locale     string
	fallback   string
	key        string
	mu         sync.Mutex
}

// loaded is a map structure used to store loaded translation data.
// It is organized as follows:
//   - First map (map[string]): Maps from locale to...
//   - Second map (map[string]): Maps from folder(group) to...
//   - Third map (map[string]): Maps from key to...
//   - Value (any): The translation line corresponding to the key in the specified locale, folder(group), and key hierarchy.
var loaded = make(map[string]map[string]map[string]any)

// contextKey is an unexported type for keys defined in this package.
type contextKey string

const (
	fallbackLocaleKey = contextKey("fallback_locale")
	localeKey         = contextKey("locale")
)

func NewTranslator(ctx context.Context, fsLoader translationcontract.Loader, fileLoader translationcontract.Loader, locale string, fallback string, logger logcontract.Log) *Translator {
	return &Translator{
		ctx:        ctx,
		fsLoader:   fsLoader,
		fileLoader: fileLoader,
		locale:     locale,
		fallback:   fallback,
		selector:   NewMessageSelector(),
		logger:     logger,
	}
}

func (t *Translator) Choice(key string, number int, options ...translationcontract.Option) string {
	line := t.Get(key, options...)

	replace := map[string]string{
		"count": strconv.Itoa(number),
	}

	locale := t.CurrentLocale()
	if len(options) > 0 && options[0].Locale != "" {
		locale = options[0].Locale
	}

	return makeReplacements(t.selector.Choose(line, number, locale), replace)
}

func (t *Translator) Get(key string, options ...translationcontract.Option) string {
	if t.key == "" {
		t.key = key
	}

	locale := t.CurrentLocale()
	// Check if a custom locale is provided in options.
	if len(options) > 0 && options[0].Locale != "" {
		locale = options[0].Locale
	}

	fallback := true
	// Check if a custom fallback is provided in options.
	if len(options) > 0 && options[0].Fallback != nil {
		fallback = *options[0].Fallback
	}

	// For JSON translations({locale}.json), we can simply load the JSON file
	// and pull the translation line from the JSON structure.We do not need
	// to do any extra processing.
	line := t.getLine(locale, "*", key, options...)
	if line != "" {
		return line
	}

	// If the key is not found in the JSON translation file, we will attempt
	// to load the key from the `{locale}/{group}.json` file.
	group, item := parseKey(key)
	line = t.getLine(locale, group, item, options...)
	if line != "" {
		return line
	}

	// If the key is not found in the current locale and fallback is enabled,
	// try to load from fallback locale
	fallbackLocale := t.GetFallback()
	if (locale != fallbackLocale) && fallback && fallbackLocale != "" {
		var fallbackOptions translationcontract.Option
		if len(options) > 0 {
			fallbackOptions = options[0]
		}
		fallbackOptions.Fallback = translationcontract.Bool(false)
		fallbackOptions.Locale = fallbackLocale
		return t.Get(key, fallbackOptions)
	}

	// Return the original key if no translation is found.
	return t.key
}

func (t *Translator) GetFallback() string {
	if fallback, ok := t.ctx.Value(string(fallbackLocaleKey)).(string); ok {
		return fallback
	}
	return t.fallback
}

func (t *Translator) CurrentLocale() string {
	if locale, ok := t.ctx.Value(string(localeKey)).(string); ok {
		return locale
	}
	return t.locale
}

func (t *Translator) Has(key string, options ...translationcontract.Option) bool {
	line := t.Get(key, options...)
	return line != key
}

func (t *Translator) SetFallback(locale string) context.Context {
	t.fallback = locale
	//nolint:all
	t.ctx = context.WithValue(t.ctx, string(fallbackLocaleKey), locale)

	return t.ctx
}

func (t *Translator) SetLocale(locale string) context.Context {
	t.locale = locale
	if ctx, ok := t.ctx.(http.Context); ok {
		ctx.WithValue(string(localeKey), locale)
		t.ctx = ctx
	} else {
		//nolint:all
		t.ctx = context.WithValue(t.ctx, string(localeKey), locale)
	}
	return t.ctx
}

func (t *Translator) getLine(locale string, group string, key string, options ...translationcontract.Option) string {
	err := t.load(locale, group)
	if err != nil {
		if errors.Is(err, errors.LangFileNotExist) {
			return ""
		}
		if errors.Is(err, errors.LangNoLoaderAvailable) {
			return t.key
		}
		t.logger.Panic(err)
		return t.key
	}

	keyValue := getValue(loaded[locale][group], key)
	// If the key doesn't exist, return empty string.
	if keyValue == nil {
		return ""
	}

	line := cast.ToString(keyValue)
	// If the line doesn't contain any placeholders, we can return it right
	// away.Otherwise, we will make the replacements on the line and return
	// the result.
	if len(options) > 0 {
		return makeReplacements(line, options[0].Replace)
	}

	return line
}

func (t *Translator) load(locale string, group string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.isLoaded(locale, group) {
		return nil
	}

	// Check if no loaders are available
	if t.fileLoader == nil && t.fsLoader == nil {
		return errors.LangNoLoaderAvailable
	}

	var (
		translations map[string]any
		err          error
	)

	if t.fileLoader != nil {
		translations, err = t.fileLoader.Load(locale, group)
	}
	if (len(translations) == 0 || err != nil) && t.fsLoader != nil {
		translations, err = t.fsLoader.Load(locale, group)
	}
	if err != nil {
		return err
	}

	if loaded[locale] == nil {
		loaded[locale] = make(map[string]map[string]any)
	}
	loaded[locale][group] = translations
	return nil
}

func (t *Translator) isLoaded(locale string, group string) bool {
	if _, ok := loaded[locale]; !ok {
		return false
	}

	if _, ok := loaded[locale][group]; !ok {
		return false
	}

	return true
}

func makeReplacements(line string, replace map[string]string) string {
	if len(replace) == 0 {
		return line
	}

	var shouldReplace []string
	casesTitle := cases.Title(language.Und)
	for k, v := range replace {
		shouldReplace = append(shouldReplace, ":"+k, v)
		shouldReplace = append(shouldReplace, ":"+casesTitle.String(k), casesTitle.String(v))
		shouldReplace = append(shouldReplace, ":"+strings.ToUpper(k), strings.ToUpper(v))
	}

	return strings.NewReplacer(shouldReplace...).Replace(line)
}

// parseKey parses a key into group and item.
func parseKey(key string) (group, item string) {
	segments := strings.Split(key, ".")

	group = segments[0]

	if len(segments) == 1 {
		item = ""
	} else {
		item = strings.Join(segments[1:], ".")
	}

	return group, item
}

// getValue an item from an object using "dot" notation.
func getValue(obj any, key string) any {
	keys := strings.Split(key, ".")

	var currentObj any
	currentObj = obj

	for _, k := range keys {
		switch v := currentObj.(type) {
		case map[string]any:
			if val, found := v[k]; found {
				currentObj = val
			} else {
				return nil
			}
		default:
			return nil
		}
	}

	return currentObj
}
