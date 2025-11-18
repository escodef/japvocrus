package csvreader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

func Load(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	var out []Entry

	filter := []string{"助動詞", "記号"}

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("csv read error: %w", err)
		}

		if len(rec) < 2 {
			continue
		}

		word := strings.TrimSpace(rec[0])
		pos := strings.TrimSpace(rec[1])

		if slices.Contains(filter, pos) {
			continue
		}

		out = append(out, Entry{
			Word: word,
			POS:  pos,
		})
	}

	return out, nil
}
