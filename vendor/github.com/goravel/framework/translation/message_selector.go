package translation

import (
	"regexp"
	"strconv"
	"strings"
)

type MessageSelector struct{}

func NewMessageSelector() *MessageSelector {
	return &MessageSelector{}
}

// Choose a translation string from an array according to a number.
func (m *MessageSelector) Choose(message string, number int, locale string) string {
	segments := strings.Split(strings.Trim(message, "\" "), "|")
	if value := m.extract(segments, number); value != nil {
		return strings.TrimSpace(*value)
	}

	segments = stripConditions(segments)
	pluralIndex := getPluralIndex(number, locale)

	if len(segments) == 1 || pluralIndex >= len(segments) {
		return strings.Trim(segments[0], "\" ")
	}

	return segments[pluralIndex]
}

func (m *MessageSelector) extract(segments []string, number int) *string {
	for _, segment := range segments {
		if line := m.extractFromString(segment, number); line != nil {
			return line
		}
	}

	return nil
}

func (m *MessageSelector) extractFromString(segment string, number int) *string {
	// Define a regular expression pattern for matching the condition and value in the part
	pattern := `^[\{\[]([^\[\]\{\}]*)[\}\]]([\s\S]*)`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(segment)
	// Check if we have exactly three sub matches (full match and two capturing groups)
	if len(matches) != 3 {
		return nil
	}

	condition, value := matches[1], matches[2]

	// Check if the condition contains a comma
	if strings.Contains(condition, ",") {
		fromTo := strings.SplitN(condition, ",", 2)
		from, to := fromTo[0], fromTo[1]

		if from == "*" {
			if to == "*" {
				return &value
			}

			toInt, err := strconv.Atoi(to)
			if err == nil && number <= toInt {
				return &value
			}
		} else if to == "*" {
			fromInt, err := strconv.Atoi(from)
			if err == nil && number >= fromInt {
				return &value
			}
		} else {
			fromInt, errFrom := strconv.Atoi(from)
			toInt, errTo := strconv.Atoi(to)

			if errFrom == nil && errTo == nil && number >= fromInt && number <= toInt {
				return &value
			}
		}
	}

	// Check if the condition is equal to the number
	conditionInt, err := strconv.Atoi(condition)
	if err == nil && conditionInt == number {
		return &value
	}

	return nil
}

func stripConditions(segments []string) []string {
	strippedSegments := make([]string, len(segments))

	for i, part := range segments {
		// Define a regular expression pattern for stripping conditions
		pattern := `^[\{\[]([^\[\]\{\}]*)[\}\]]`

		// Replace the matched pattern with an empty string
		re := regexp.MustCompile(pattern)
		strippedPart := re.ReplaceAllString(part, "")

		// Store the stripped part in the result slice
		strippedSegments[i] = strippedPart
	}

	return strippedSegments
}

// getPluralIndex returns the plural index for the given number and locale.
func getPluralIndex(number int, locale string) int {
	switch locale {
	case "az", "az_AZ", "bo", "bo_CN", "bo_IN", "dz", "dz_BT", "id", "id_ID", "ja", "ja_JP", "jv", "ka", "ka_GE", "km", "km_KH", "kn", "kn_KR", "ko", "ko_KR", "ms", "ms_MY", "th", "th_TH", "tr", "tr_TR", "vi", "vi_VN", "zh", "zh_CN", "zh_HK", "zh_SG", "zh_TW":
		return 0
	case "af", "af_ZA", "bn", "bn_BD", "bn_IN", "bg", "bg_BG", "ca", "ca_AD", "ca_ES", "ca_FR", "ca_IT", "da", "da_DK", "de", "de_AT", "de_BE", "de_CH", "de_DE", "de_LI", "de_LU", "el", "el_CY", "el_GR", "en", "en_AS", "en_AU", "en_BE", "en_BW", "en_BZ", "en_CA", "en_GB", "en_GU", "en_HK", "en_IE", "en_IN", "en_JM", "en_MH", "en_MP", "en_MT", "en_NA", "en_NZ", "en_PH", "en_PK", "en_SG", "en_TT", "en_UM", "en_US", "en_VI", "en_ZA", "en_ZW", "eo", "eo_US", "es", "es_AR", "es_BO", "es_CL", "es_CO", "es_CR", "es_DO", "es_EC", "es_ES", "es_GT", "es_HN", "es_MX", "es_NI", "es_PA", "es_PE", "es_PR", "es_PY", "es_SV", "es_US", "es_UY", "es_VE", "et", "et_EE", "eu", "eu_ES", "eu_FR", "fi", "fi_FI", "fo", "fo_FO", "fur", "fur_IT", "fy", "fy_DE", "fy_NL", "gl", "gl_ES", "gu", "gu_IN", "ha", "ha_NG", "he", "he_IL", "hu", "hu_HU", "is", "is_IS", "it", "it_CH", "it_IT", "ku", "ku_TR", "lb", "lb_LU", "ml", "ml_IN", "mn", "mn_MN", "mr", "mr_IN", "nah", "nb", "nb_NO", "ne", "ne_NP", "nl", "nl_BE", "nl_NL", "nn", "nn_NO", "no", "om", "om_ET", "om_KE", "or", "or_IN", "pa", "pa_IN", "pa_PK", "pap", "pap_AN", "pap_AW", "pap_CW", "pap_NL", "ps", "ps_AF", "pt", "pt_BR", "pt_PT", "so", "so_DJ", "so_ET", "so_KE", "so_SO", "sq", "sq_AL", "sq_MK", "sv", "sv_FI", "sv_SE", "sw", "sw_KE", "sw_TZ", "ta", "ta_IN", "ta_LK", "te", "te_IN", "tk", "tk_TM", "ur", "ur_IN", "ur_PK", "zu", "zu_ZA":
		if number == 1 {
			return 0
		}
		return 1
	case "am", "am_ET", "bh", "fil", "fil_PH", "fr", "fr_BE", "fr_CA", "fr_CH", "fr_FR", "fr_LU", "fr_MC", "fr_SN", "guz", "guz_KE", "hi", "hi_IN", "hy", "hy_AM", "ln", "ln_CD", "mg", "mg_MG", "nso", "nso_ZA", "ti", "ti_ER", "ti_ET", "wa", "wa_BE", "xbr":
		if number == 0 || number == 1 {
			return 0
		}
		return 1
	case "be", "be_BY", "bs", "bs_BA", "hr", "hr_HR", "ru", "ru_RU", "ru_UA", "sr", "sr_ME", "sr_RS", "uk", "uk_UA":
		if number%10 == 1 && number%100 != 11 {
			return 0
		}
		if number%10 >= 2 && number%10 <= 4 && (number%100 < 10 || number%100 >= 20) {
			return 1
		}
		return 2
	case "cs", "cs_CZ", "sk", "sk_SK":
		if number == 1 {
			return 0
		}
		if number >= 2 && number <= 4 {
			return 1
		}
		return 2
	case "ga", "ga_IE":
		if number == 1 {
			return 0
		}
		if number == 2 {
			return 1
		}
		return 2
	case "lt", "lt_LT":
		if number%10 == 1 && number%100 != 11 {
			return 0
		}
		if number%10 >= 2 && (number%100 < 10 || number%100 >= 20) {
			return 1
		}
		return 2
	case "sl", "sl_SI":
		if number%100 == 1 {
			return 0
		}
		if number%100 == 2 {
			return 1
		}
		if number%100 == 3 || number%100 == 4 {
			return 2
		}
		return 3
	case "mk", "mk_MK":
		if number%10 == 1 {
			return 0
		}
		return 1
	case "mt", "mt_MT":
		if number == 1 {
			return 0
		}
		if number == 0 || (number%100 > 1 && number%100 < 11) {
			return 1
		}
		if number%100 > 10 && number%100 < 20 {
			return 2
		}
		return 3
	case "lv", "lv_LV":
		if number == 0 {
			return 0
		}
		if number%10 == 1 && number%100 != 11 {
			return 1
		}
		return 2
	case "pl", "pl_PL":
		if number == 1 {
			return 0
		}
		if number%10 >= 2 && number%10 <= 4 && (number%100 < 12 || number%100 > 14) {
			return 1
		}
		return 2
	case "cy", "cy_GB":
		if number == 1 {
			return 0
		}
		if number == 2 {
			return 1
		}
		if number == 8 || number == 11 {
			return 2
		}
		return 3
	case "ro", "ro_RO":
		if number == 1 {
			return 0
		}
		if number == 0 || (number%100 > 0 && number%100 < 20) {
			return 1
		}
		return 2
	case "ar", "ar_AE", "ar_BH", "ar_DZ", "ar_EG", "ar_IN", "ar_IQ", "ar_JO", "ar_KW", "ar_LB", "ar_LY", "ar_MA", "ar_OM", "ar_QA", "ar_SA", "ar_SD", "ar_SY", "ar_TN", "ar_YE":
		if number == 0 {
			return 0
		}
		if number == 1 {
			return 1
		}
		if number == 2 {
			return 2
		}
		if number%100 >= 3 && number%100 <= 10 {
			return 3
		}
		if number%100 >= 11 && number%100 <= 99 {
			return 4
		}
		return 5
	default:
		return 0
	}
}
