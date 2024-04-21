package googlecache

import "fmt"

func BuildCacheURL(url string) string {
	return fmt.Sprintf(`%s%s`, googleCacheURL, url)
}
