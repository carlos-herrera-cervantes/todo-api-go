package locales

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Return a localizer for specific language
func GetLocalizer(lang, key string) string {
	if !(len(lang) > 0) {
		lang = "en"
	}

	bundle := i18n.NewBundle(language.English)
	bundle.LoadMessageFile(lang + ".json")
	localizer := i18n.NewLocalizer(bundle, lang, lang)

	message := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: key,
		},
	})

	return message
}
