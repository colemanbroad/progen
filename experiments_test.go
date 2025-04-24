package main

import (
	"os"
	"testing"
)

func TestExperPeano(t *testing.T) {
	*dbname = testdir + "testdb.peano.db"
	runPeano()
	db := ConnectSqlite(*dbname)
	script, err := os.ReadFile("wire.sql")
	if err != nil {
		t.Error("missing file ", err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Error("queries failed", err)
	}
}

func TestExperPow2(t *testing.T) {
	*dbname = testdir + "testdb.pow2.db"
	runPow2()
	db := ConnectSqlite(*dbname)
	script, err := os.ReadFile("p2.sql")
	if err != nil {
		t.Error("missing file ", err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Error("queries failed", err)
	}
}
