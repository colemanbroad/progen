package main

import (
	"os"
	"testing"
)

func Test_runPeano(t *testing.T) {
	*dbname = testdir + "testdb.peano.db" // global
	// swap stdout
	ofile, _ := os.OpenFile(testdir+"testdb.peano.out", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	os.Stdout = ofile
	stdout := os.Stdout
	defer func() {
		os.Stdout = stdout
		ofile.Close()
	}()
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

func Test_runPow2(t *testing.T) {
	*dbname = testdir + "testdb.pow2.db"
	// swap stdout
	ofile, _ := os.OpenFile(testdir+"testdb.pow2.out", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	os.Stdout = ofile
	stdout := os.Stdout
	defer func() {
		os.Stdout = stdout
		ofile.Close()
	}()
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
