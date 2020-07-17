package like

import "regexp"

type like struct {
	reg *regexp.Regexp
	fn  func(string) bool
}
