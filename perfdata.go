package nagios

import "strings"

type Perfdata []string

func (p Perfdata) Get(key string) string {
	for _, line := range p {
		if k, v, ok := getKV(line); ok {
			if strings.EqualFold(k, key) {
				return v
			}
		}
	}

	return ""
}

func (p Perfdata) Iterate(cb func(key, value string) error) error {
	for _, line := range p {
		if k, v, ok := getKV(line); ok {
			if err := cb(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p Perfdata) String() string {
	return strings.Join(p, " ")
}

//nolint:mnd // KV is K=V.
func getKV(in string) (string, string, bool) {
	sp := strings.SplitN(in, "=", 2)
	if len(sp) == 2 {
		return sp[0], sp[1], true
	}

	return "", "", false
}
