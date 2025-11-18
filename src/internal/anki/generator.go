package anki

import (
	"archive/zip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"japvocrus/internal/dict"
	"japvocrus/internal/util"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func GenerateApkg(words []dict.Translation, audioDir, output string, ttsEnabled bool) error {
	tmpdir, err := os.MkdirTemp("", "ankigen")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpdir)

	dbPath := filepath.Join(tmpdir, "collection.anki2")

	if err := createSQLite(dbPath); err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := insertColRow(db); err != nil {
		return err
	}

	modelID := time.Now().UnixMilli()
	deckID := modelID + 1

	if err := updateModelsJSON(db, modelID); err != nil {
		return err
	}
	if err := updateDecksJSON(db, deckID); err != nil {
		return err
	}

	mediaMap := map[int]string{}
	mediaIndex := 0

	noteID := time.Now().UnixMilli() + 10
	cardID := noteID + 1000
	nowSeconds := time.Now().Unix()

	for _, w := range words {
		nid := noteID
		noteID++

		audioFile := fmt.Sprintf("%s.wav", w.Word)
		audioFull := filepath.Join(audioDir, audioFile)

		mediaMap[mediaIndex] = audioFile

		dst := filepath.Join(tmpdir, audioFile)
		if err := copyFile(audioFull, dst); err != nil {
			return fmt.Errorf("wav missing for %s: %w", w.Word, err)
		}

		var flds string
		if ttsEnabled {
			flds = fmt.Sprintf("%s\x1f%s [sound:%s]", util.SenseToString(w.Senses), w.Word, audioFile)
		} else {
			flds = fmt.Sprintf("%s\x1f%s", util.SenseToString(w.Senses), w.Word)
		}

		_, err = db.Exec(`
			INSERT INTO notes(id, guid, mid, mod, usn, tags, flds, sfld, csum, flags, data)
			VALUES (?, ?, ?, ?, 0, ' ', ?, 0, 0, 0, '')
		`, nid, fmt.Sprintf("guid-%d", nid), modelID, nowSeconds, flds)
		if err != nil {
			return err
		}

		_, err = db.Exec(`
			INSERT INTO cards
			(id, nid, did, ord, mod, usn, type, queue, due, ivl, factor,
			 reps, lapses, left, odue, odid, flags, data)
			VALUES (?, ?, ?, 0, ?, 0, 0, 0, ?, 0, 2500,
			        0, 0, 0, 0, 0, 0, '')
		`, cardID, nid, deckID, nowSeconds, cardID)
		if err != nil {
			return err
		}
		cardID++

		_, err = db.Exec(`
			INSERT INTO cards
			(id, nid, did, ord, mod, usn, type, queue, due, ivl, factor,
			 reps, lapses, left, odue, odid, flags, data)
			VALUES (?, ?, ?, 1, ?, 0, 0, 0, ?, 0, 2500,
			        0, 0, 0, 0, 0, 0, '')
		`, cardID, nid, deckID, nowSeconds, cardID)
		if err != nil {
			return err
		}
		cardID++

		mediaIndex++
	}

	mediaJSON, _ := json.MarshalIndent(mediaMap, "", "  ")
	if err := os.WriteFile(filepath.Join(tmpdir, "media"), mediaJSON, 0644); err != nil {
		return err
	}

	return zipFolder(tmpdir, output)
}

func updateModelsJSON(db *sql.DB, modelID int64) error {
	m := map[string]any{
		fmt.Sprintf("%d", modelID): map[string]any{
			"id":    modelID,
			"name":  "RU-JP Model",
			"type":  0,
			"mod":   time.Now().Unix(),
			"usn":   0,
			"sortf": 0,
			"css":   "",
			"flds": []map[string]any{
				{"name": "Ru", "ord": 0},
				{"name": "Jp", "ord": 1},
			},
			"tmpls": []map[string]any{
				{
					"name": "Ru → Jp",
					"ord":  0,
					"qfmt": "{{Ru}}",
					"afmt": "{{FrontSide}}<hr id=answer>{{Jp}}",
				},
				{
					"name": "Jp → Ru",
					"ord":  1,
					"qfmt": "{{Jp}}",
					"afmt": "{{FrontSide}}<hr id=answer>{{Ru}}",
				},
			},
			"vers": []any{},
			"tags": []any{},
		},
	}

	j, _ := json.Marshal(m)
	_, err := db.Exec(`UPDATE col SET models = ?`, j)
	return err
}

func updateDecksJSON(db *sql.DB, deckID int64) error {
	d := map[string]any{
		fmt.Sprintf("%d", deckID): map[string]any{
			"id":   deckID,
			"name": "Default",
			"usn":  0,
			"mod":  time.Now().Unix(),
		},
	}

	j, _ := json.Marshal(d)
	_, err := db.Exec(`UPDATE col SET decks = ?`, j)
	return err
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func zipFolder(folder, outFile string) error {
	zf, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer zf.Close()

	zw := zip.NewWriter(zf)
	defer zw.Close()

	return filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(folder, path)
		f, err := zw.Create(rel)
		if err != nil {
			return err
		}

		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		_, err = io.Copy(f, src)
		return err
	})
}
