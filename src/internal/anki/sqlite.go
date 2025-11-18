package anki

import (
	"database/sql"
	"time"
)

func createSQLite(path string) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	defer db.Close()

	schema := `
CREATE TABLE col (
    id integer primary key,
    crt integer not null,
    mod integer not null,
    scm integer not null,
    ver integer not null,
    dty integer not null,
    usn integer not null,
    ls integer not null,
    conf text not null,
    models text not null,
    decks text not null,
    dconf text not null,
    tags text not null
);
CREATE TABLE notes (
    id integer primary key,
    guid text not null,
    mid integer not null,
    mod integer not null,
    usn integer not null,
    tags text not null,
    flds text not null,
    sfld integer not null,
    csum integer not null,
    flags integer not null,
    data text not null
);
CREATE TABLE cards (
    id integer primary key,
    nid integer not null,
    did integer not null,
    ord integer not null,
    mod integer not null,
    usn integer not null,
    type integer not null,
    queue integer not null,
    due integer not null,
    ivl integer not null,
    factor integer not null,
    reps integer not null,
    lapses integer not null,
    left integer not null,
    odue integer not null,
    odid integer not null,
    flags integer not null,
    data text not null
);
CREATE TABLE revlog (
    id integer primary key,
    cid integer not null,
    usn integer not null,
    ease integer not null,
    ivl integer not null,
    lastIvl integer not null,
    factor integer not null,
    time integer not null,
    type integer not null
);
`
	_, err = db.Exec(schema)
	return err
}

func insertColRow(db *sql.DB) error {
	now := time.Now().Unix()

	conf := `{"nextPos": 1, "curDeck": 1}`
	empty := `{}`

	_, err := db.Exec(`
		INSERT INTO col(id, crt, mod, scm, ver, dty, usn, ls,
		                conf, models, decks, dconf, tags)
		VALUES (1, ?, ?, ?, 11, 0, 0, 0, ?, ?, ?, ?, '')
	`, now, now, now, conf, empty, empty, empty)
	return err
}
