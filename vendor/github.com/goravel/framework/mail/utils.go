package mail

import "strings"

func convertMapHeadersToSlice(headers map[string]string) []string {
	var slice []string
	for key, value := range headers {
		slice = append(slice, key+": "+value)
	}
	return slice
}

func convertSliceHeadersToMap(headers []string) map[string]string {
	mapHeaders := make(map[string]string)
	for _, header := range headers {
		parts := strings.Split(header, ":")
		if len(parts) == 2 {
			mapHeaders[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return mapHeaders
}
