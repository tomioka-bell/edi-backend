package utils

import "strings"

func DisplayName(email string) string {
	local := email
	if at := strings.IndexByte(email, '@'); at > 0 {
		local = email[:at]
	}
	replacer := strings.NewReplacer(".", " ", "_", " ", "-", " ")
	local = replacer.Replace(local)

	parts := strings.Fields(local)
	for i, w := range parts {
		if len(w) == 1 {
			parts[i] = strings.ToUpper(w)
		} else {
			parts[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		}
	}
	if len(parts) == 0 {
		return "User"
	}
	return strings.Join(parts, " ")
}
