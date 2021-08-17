package db

import (
	"os"
	"os/exec"
	"testing"
	"fmt"
)

/* Test new CSV database instance with correct DB_PATH environment variable */ 
func TestNewCsvDatabaseWithPathEnv(t *testing.T) {
	os.Setenv("DB_PATH", "geo_db.csv")
	defer os.Unsetenv("DB_PATH")
	csvDb := NewCsvDb()
	fmt.Println(csvDb.GetLocation("1.2.3.4"))
}

/* Test new CSV database instance without DB_PATH environment variable - expect a failure */ 
func TestNewCsvDBWithoutPathEnv(t *testing.T) {
	if os.Getenv("UT_DB_PATH_EMPTY") == "1" {
		csvDb := NewCsvDb()
		csvDb.GetLocation("1.2.3.4")
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestNewCsvDBWithoutPathEnv")
	cmd.Env = append(os.Environ(), "UT_DB_PATH_EMPTY=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
