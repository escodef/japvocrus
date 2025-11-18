package util

import "japvocrus/internal/dict"

func SenseToString(s []dict.Sense) string {
	out := ""
	for i := range s {
		out = out + s[i].Ru + " "
	}
	return out
}
