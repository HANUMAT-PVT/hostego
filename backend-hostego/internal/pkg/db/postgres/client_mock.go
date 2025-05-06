package postgres

import (
	"backend-hostego/internal/app/hostego-service/constants/config_constants"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

const (
	// CREATE - operation to create database
	CREATE = "CREATE"
	// DROP - operation to drop database
	DROP = "DROP"
)

// MockedDB is used in unit tests to mock db
func MockedDB(operation string) {
	/*
	    reference: https://mayursinhsarvaiya.medium.com/how-to-mock-postgresql-database-for-unit-testing-in-golang-gorm-b690a4e4bc85
	   If tests are running in CI, environment variables should not be loaded.
	   The reason is environment vars will be provided through CI config file.
	*/

	dbName := viper.GetString(config_constants.VKEYS_DATABASE_POSTGRES_SOURCE_DB_NAME)
	pgUser := viper.GetString(config_constants.VKEYS_DATABASE_POSTGRES_SOURCE_USER)
	pgPassword := viper.GetString(config_constants.VKEYS_DATABASE_POSTGRES_SOURCE_PASSWORD)

	// createdb => https://www.postgresql.org/docs/7.0/app-createdb.htm
	// dropdb => https://www.postgresql.org/docs/7.0/app-dropdb.htm
	var command string

	if operation == CREATE {
		command = "createdb"
	} else {
		command = "dropdb"
	}

	// createdb & dropdb commands have same configuration syntax.
	cmd := exec.Command(command, "-h", "localhost", "-U", pgUser, "-e", dbName)
	cmd.Env = os.Environ()

	/*
	   if we normally execute createdb/dropdb, we will be propmted to provide password.
	   To inject password automatically, we have to set PGPASSWORD as prefix.
	*/
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%v", pgPassword))

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error executing %v on %v.\n%v", command, dbName, err)
		return
	}

	/*
	   Alternatively instead of createdb/dropdb, you can use
	   psql -c "CREATE/DROP DATABASE DBNAME" "DATABASE_URL"
	*/

	if operation == CREATE {
		// migrate tables
		err := MigrationsMock()
		if err != nil {
			MockedDB(DROP)
			return
		}
		// fill meta data
		FillingMockedMetaData()
	}

}

func MigrationsMock() error {
	//db := InitDB()
	//m := gormigrate.New(db, gormigrate.DefaultOptions, migration.PostgresMigrations)
	//
	//log.Printf("starting migrations...")
	//
	//if err := m.Migrate(); err != nil {
	//	log.Fatalf("Could not migrate: %v", err)
	//	return err
	//}
	return nil
}

func FillingMockedMetaData() {

}
