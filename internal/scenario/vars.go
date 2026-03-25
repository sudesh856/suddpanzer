package scenario

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)


var varPattern = regexp.MustCompile(`\{\{\s*(.+?)\s*\}\}`)

func ReplaceVars(s string, extra map[string]string) string {
	return varPattern.ReplaceAllStringFunc(s, func(match string) string {
		inner := strings.TrimSpace(varPattern.FindStringSubmatch(match)[1])

		if extra != nil {
			if v, ok := extra[inner]; ok {
				return v
			}
		}

		switch {
		case inner == "uuid":
			return uuid.NewString()
		
		case inner == "timestamp":
			return strconv.FormatInt(time.Now().Unix(), 10)
		case strings.HasPrefix(inner, "random_int "):
			return parseRandomInt(inner)
		case strings.HasPrefix(inner, "env."):
			return os.Getenv(strings.TrimPrefix(inner, "env."))
		default:
			return match
		}

	})
}

func parseRandomInt(expr string) string {
	parts := strings.Fields(expr)
	if len(parts) != 3 {
		return fmt.Sprintf("{{%s}}", expr)
	}

	min,err1 := strconv.Atoi(parts[1])
	max, err2 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || max < min {
		return fmt.Sprintf("{{%s}}", expr)
	}
	return strconv.Itoa(min + rand.Intn(max-min+1))
}
