package digest

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func parseIdxFromPath(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	if len(s) < 1 {
		return 0, errors.Errorf("empty idx")
	} else if len(s) > 1 && strings.HasPrefix(s, "0") {
		return 0, errors.Errorf("invalid idx, %q", s)
	}

	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}
