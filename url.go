package charm

import (
	"regexp"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// A charm URL represents charm locations such as:
//
//     cs:~joe/oneiric/wordpress
//     cs:oneiric/wordpress-42
//     local:oneiric/wordpress
//
type URL struct {
	Name     string
	Revision int // -1 if unset, 0 is valid
	Collection
}

// A charm Collection represents a namespace of charms. The
// collection precedes the charm name in a charm URL.
type Collection struct {
	Schema string
	User   string
	Series string
}

var validUser = regexp.MustCompile("^[a-z0-9][a-zA-Z0-9+.-]+$")
var validSeries = regexp.MustCompile("^[a-z]+([a-z-]+[a-z])?$")
var validName = regexp.MustCompile("^[a-z][a-z0-9]*(-[a-z0-9]*[a-z][a-z0-9]*)*$")

func NewURL(url string) (*URL, os.Error) {
	u := &URL{}
	i := strings.Index(url, ":")
	if i > 0 {
		u.Schema = url[:i]
		i++
	}
	// cs: or local:
	if u.Schema != "cs" && u.Schema != "local" {
		return nil, fmt.Errorf("charm URL has invalid schema: %q", url)
	}
	parts := strings.Split(url[i:], "/")
	if len(parts) < 1 || len(parts) > 3 {
		return nil, fmt.Errorf("charm URL has invalid form: %q", url)
	}

	// ~<username>
	if strings.HasPrefix(parts[0], "~") {
		if u.Schema == "local" {
			return nil, fmt.Errorf("local charm URL with user name: %q", url)
		}
		u.User = parts[0][1:]
		if !validUser.MatchString(u.User) {
			return nil, fmt.Errorf("charm URL has invalid user name: %q", url)
		}
		parts = parts[1:]
	}

	// <series>
	if len(parts) < 2 {
		return nil, fmt.Errorf("charm URL without series: %q", url)
	}
	if len(parts) == 2 {
		u.Series = parts[0]
		if !validSeries.MatchString(u.Series) {
			return nil, fmt.Errorf("charm URL has invalid series: %q", url)
		}
		parts = parts[1:]
	}

	// <name>[-<revision>]
	u.Name = parts[0]
	u.Revision = -1
	for i := len(u.Name)-1; i > 0; i-- {
		c := u.Name[i]
		if c >= '0' && c <= '9' {
			continue
		}
		if c == '-' && i != len(u.Name)-1 {
			var err os.Error
			u.Revision, err = strconv.Atoi(u.Name[i+1:])
			if err != nil {
				panic(err) // We just checked it was right.
			}
			u.Name = u.Name[:i]
		}
		break
	}
	if !validName.MatchString(u.Name) {
		return nil, fmt.Errorf("charm URL has invalid charm name: %q", url)
	}
	return u, nil
}

func (u *URL) String() string {
	if u.User != "" {
		if u.Revision >= 0 {
			return fmt.Sprintf("%s:~%s/%s/%s-%d", u.Schema, u.User, u.Series, u.Name, u.Revision)
		}
		return fmt.Sprintf("%s:~%s/%s/%s", u.Schema, u.User, u.Series, u.Name)
	}
	if u.Revision >= 0 {
		return fmt.Sprintf("%s:%s/%s-%d", u.Schema, u.Series, u.Name, u.Revision)
	}
	return fmt.Sprintf("%s:%s/%s", u.Schema, u.Series, u.Name)
}