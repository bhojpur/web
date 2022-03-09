package migrate

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"database/sql"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bhojpur/web/cmd/utility/commands"
	"github.com/bhojpur/web/cmd/utility/commands/version"
	"github.com/bhojpur/web/pkg/client/config"
	"github.com/bhojpur/web/pkg/client/utils"

	cliLogger "github.com/bhojpur/web/pkg/client/logger"
)

var CmdMigrate = &commands.Command{
	UsageLine: "migrate [Command]",
	Short:     "Runs database migrations",
	Long: `The command 'migrate' allows you to run database migrations to keep it up-to-date.

  ▶ {{"To run all the migrations:"|bold}}

    $ webutl migrate [-driver=mysql] [-conn="root:@tcp(127.0.0.1:3306)/test"] [-dir="path/to/migration"]

  ▶ {{"To rollback the last migration:"|bold}}

    $ webutl migrate rollback [-driver=mysql] [-conn="root:@tcp(127.0.0.1:3306)/test"] [-dir="path/to/migration"]

  ▶ {{"To do a reset, which will rollback all the migrations:"|bold}}

    $ webutl migrate reset [-driver=mysql] [-conn="root:@tcp(127.0.0.1:3306)/test"] [-dir="path/to/migration"]

  ▶ {{"To update your schema:"|bold}}

    $ webutl migrate refresh [-driver=mysql] [-conn="root:@tcp(127.0.0.1:3306)/test"] [-dir="path/to/migration"]
`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    RunMigration,
}

var mDriver utils.DocValue
var mConn utils.DocValue
var mDir utils.DocValue

func init() {
	CmdMigrate.Flag.Var(&mDriver, "driver", "Database driver. Either mysql, postgres or sqlite.")
	CmdMigrate.Flag.Var(&mConn, "conn", "Connection string used by the driver to connect to a database instance.")
	CmdMigrate.Flag.Var(&mDir, "dir", "The directory where the migration files are stored")
	commands.AvailableCommands = append(commands.AvailableCommands, CmdMigrate)
}

// runMigration is the entry point for starting a migration
func RunMigration(cmd *commands.Command, args []string) int {
	currpath, _ := os.Getwd()

	// Getting command line arguments
	if len(args) != 0 {
		cmd.Flag.Parse(args[1:])
	}
	if mDriver == "" {
		mDriver = utils.DocValue(config.Conf.Database.Driver)
		if mDriver == "" {
			mDriver = "mysql"
		}
	}
	if mConn == "" {
		mConn = utils.DocValue(config.Conf.Database.Conn)
		if mConn == "" {
			mConn = "root:@tcp(127.0.0.1:3306)/test"
		}
	}
	if mDir == "" {
		mDir = utils.DocValue(config.Conf.Database.Dir)
		if mDir == "" {
			mDir = utils.DocValue(path.Join(currpath, "database", "migrations"))
		}
	}

	cliLogger.Log.Infof("Using '%s' as 'driver'", mDriver)
	//Log sensitive connection information only when DEBUG is set to true.
	cliLogger.Log.Debugf("Conn: %s", utils.FILE(), utils.LINE(), mConn)
	cliLogger.Log.Infof("Using '%s' as 'dir'", mDir)
	driverStr, connStr, dirStr := string(mDriver), string(mConn), string(mDir)

	dirRune := []rune(dirStr)

	if dirRune[0] != '/' && dirRune[1] != ':' {
		dirStr = path.Join(currpath, dirStr)
	}

	if len(args) == 0 {
		// run all outstanding migrations
		cliLogger.Log.Info("Running all outstanding migrations")
		MigrateUpdate(currpath, driverStr, connStr, dirStr)
	} else {
		mcmd := args[0]
		switch mcmd {
		case "rollback":
			cliLogger.Log.Info("Rolling back the last migration operation")
			MigrateRollback(currpath, driverStr, connStr, dirStr)
		case "reset":
			cliLogger.Log.Info("Reseting all migrations")
			MigrateReset(currpath, driverStr, connStr, dirStr)
		case "refresh":
			cliLogger.Log.Info("Refreshing all migrations")
			MigrateRefresh(currpath, driverStr, connStr, dirStr)
		default:
			cliLogger.Log.Fatal("Command is missing")
		}
	}
	cliLogger.Log.Success("Migration successful!")
	return 0
}

// migrate generates source code, build it, and invoke the binary who does the actual migration
func migrate(goal, currpath, driver, connStr, dir string) {
	if dir == "" {
		dir = path.Join(currpath, "database", "migrations")
	}
	postfix := ""
	if runtime.GOOS == "windows" {
		postfix = ".exe"
	}
	binary := "m" + postfix
	source := binary + ".go"

	// Connect to database
	db, err := sql.Open(driver, connStr)
	if err != nil {
		cliLogger.Log.Fatalf("Could not connect to database using '%s': %s", connStr, err)
	}
	defer db.Close()

	checkForSchemaUpdateTable(db, driver)
	latestName, latestTime := getLatestMigration(db, goal)
	writeMigrationSourceFile(dir, source, driver, connStr, latestTime, latestName, goal)
	buildMigrationBinary(dir, binary)
	runMigrationBinary(dir, binary)
	removeTempFile(dir, source)
	removeTempFile(dir, binary)
}

// checkForSchemaUpdateTable checks the existence of migrations table.
// It checks for the proper table structures and creates the table using MYSQL_MIGRATION_DDL if it does not exist.
func checkForSchemaUpdateTable(db *sql.DB, driver string) {
	showTableSQL := showMigrationsTableSQL(driver)
	if rows, err := db.Query(showTableSQL); err != nil {
		cliLogger.Log.Fatalf("Could not show migrations table: %s", err)
	} else if !rows.Next() {
		// No migrations table, create new ones
		createTableSQL := createMigrationsTableSQL(driver)

		cliLogger.Log.Infof("Creating 'migrations' table...")

		if _, err := db.Query(createTableSQL); err != nil {
			cliLogger.Log.Fatalf("Could not create migrations table: %s", err)
		}
	}

	// Checking that migrations table schema are expected
	selectTableSQL := selectMigrationsTableSQL(driver)
	if rows, err := db.Query(selectTableSQL); err != nil {
		cliLogger.Log.Fatalf("Could not show columns of migrations table: %s", err)
	} else {
		for rows.Next() {
			var fieldBytes, typeBytes, nullBytes, keyBytes, defaultBytes, extraBytes []byte
			if err := rows.Scan(&fieldBytes, &typeBytes, &nullBytes, &keyBytes, &defaultBytes, &extraBytes); err != nil {
				cliLogger.Log.Fatalf("Could not read column information: %s", err)
			}
			fieldStr, typeStr, nullStr, keyStr, defaultStr, extraStr :=
				string(fieldBytes), string(typeBytes), string(nullBytes), string(keyBytes), string(defaultBytes), string(extraBytes)
			if fieldStr == "id_migration" {
				if keyStr != "PRI" || extraStr != "auto_increment" {
					cliLogger.Log.Hint("Expecting KEY: PRI, EXTRA: auto_increment")
					cliLogger.Log.Fatalf("Column migration.id_migration type mismatch: KEY: %s, EXTRA: %s", keyStr, extraStr)
				}
			} else if fieldStr == "name" {
				if !strings.HasPrefix(typeStr, "varchar") || nullStr != "YES" {
					cliLogger.Log.Hint("Expecting TYPE: varchar, NULL: YES")
					cliLogger.Log.Fatalf("Column migration.name type mismatch: TYPE: %s, NULL: %s", typeStr, nullStr)
				}
			} else if fieldStr == "created_at" {
				if typeStr != "timestamp" || (!strings.EqualFold(defaultStr, "CURRENT_TIMESTAMP") && !strings.EqualFold(defaultStr, "CURRENT_TIMESTAMP()")) {
					cliLogger.Log.Hint("Expecting TYPE: timestamp, DEFAULT: CURRENT_TIMESTAMP || CURRENT_TIMESTAMP()")
					cliLogger.Log.Fatalf("Column migration.timestamp type mismatch: TYPE: %s, DEFAULT: %s", typeStr, defaultStr)
				}
			}
		}
	}
}

func driverImportStatement(driver string) string {
	switch driver {
	case "mysql":
		return "github.com/go-sql-driver/mysql"
	case "postgres":
		return "github.com/lib/pq"
	default:
		return "github.com/go-sql-driver/mysql"
	}
}

func showMigrationsTableSQL(driver string) string {
	switch driver {
	case "mysql":
		return "SHOW TABLES LIKE 'migrations'"
	case "postgres":
		return "SELECT * FROM pg_catalog.pg_tables WHERE tablename = 'migrations';"
	default:
		return "SHOW TABLES LIKE 'migrations'"
	}
}

func createMigrationsTableSQL(driver string) string {
	switch driver {
	case "mysql":
		return MYSQLMigrationDDL
	case "postgres":
		return POSTGRESMigrationDDL
	default:
		return MYSQLMigrationDDL
	}
}

func selectMigrationsTableSQL(driver string) string {
	switch driver {
	case "mysql":
		return "DESC migrations"
	case "postgres":
		return "SELECT * FROM migrations WHERE false ORDER BY id_migration;"
	default:
		return "DESC migrations"
	}
}

// getLatestMigration retrives latest migration with status 'update'
func getLatestMigration(db *sql.DB, goal string) (file string, createdAt int64) {
	sql := "SELECT name FROM migrations where status = 'update' ORDER BY id_migration DESC LIMIT 1"
	if rows, err := db.Query(sql); err != nil {
		cliLogger.Log.Fatalf("Could not retrieve migrations: %s", err)
	} else {
		if rows.Next() {
			if err := rows.Scan(&file); err != nil {
				cliLogger.Log.Fatalf("Could not read migrations in database: %s", err)
			}
			createdAtStr := file[len(file)-15:]
			if t, err := time.Parse("20060102_150405", createdAtStr); err != nil {
				cliLogger.Log.Fatalf("Could not parse time: %s", err)
			} else {
				createdAt = t.Unix()
			}
		} else {
			// migration table has no 'update' record, no point rolling back
			if goal == "rollback" {
				cliLogger.Log.Fatal("There is nothing to rollback")
			}
			file, createdAt = "", 0
		}
	}
	return
}

// writeMigrationSourceFile create the source file based on MIGRATION_MAIN_TPL
func writeMigrationSourceFile(dir, source, driver, connStr string, latestTime int64, latestName string, task string) {
	changeDir(dir)
	if f, err := os.OpenFile(source, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err != nil {
		cliLogger.Log.Fatalf("Could not create file: %s", err)
	} else {
		content := strings.Replace(MigrationMainTPL, "{{DBDriver}}", driver, -1)
		content = strings.Replace(content, "{{DriverRepo}}", driverImportStatement(driver), -1)
		content = strings.Replace(content, "{{ConnStr}}", connStr, -1)
		content = strings.Replace(content, "{{LatestTime}}", strconv.FormatInt(latestTime, 10), -1)
		content = strings.Replace(content, "{{LatestName}}", latestName, -1)
		content = strings.Replace(content, "{{Task}}", task, -1)
		if _, err := f.WriteString(content); err != nil {
			cliLogger.Log.Fatalf("Could not write to file: %s", err)
		}
		utils.CloseFile(f)
	}
}

// buildMigrationBinary changes directory to database/migrations folder and go-build the source
func buildMigrationBinary(dir, binary string) {
	changeDir(dir)
	cmd := exec.Command("go", "build", "-o", binary)
	if out, err := cmd.CombinedOutput(); err != nil {
		cliLogger.Log.Errorf("Could not build migration binary: %s", err)
		formatShellErrOutput(string(out))
		removeTempFile(dir, binary)
		removeTempFile(dir, binary+".go")
		os.Exit(2)
	}
}

// runMigrationBinary runs the migration program who does the actual work
func runMigrationBinary(dir, binary string) {
	changeDir(dir)
	cmd := exec.Command("./" + binary)
	if out, err := cmd.CombinedOutput(); err != nil {
		formatShellOutput(string(out))
		cliLogger.Log.Errorf("Could not run migration binary: %s", err)
		removeTempFile(dir, binary)
		removeTempFile(dir, binary+".go")
		os.Exit(2)
	} else {
		formatShellOutput(string(out))
	}
}

// changeDir changes working directory to dir.
// It exits the system when encouter an error
func changeDir(dir string) {
	if err := os.Chdir(dir); err != nil {
		cliLogger.Log.Fatalf("Could not find migration directory: %s", err)
	}
}

// removeTempFile removes a file in dir
func removeTempFile(dir, file string) {
	changeDir(dir)
	if err := os.Remove(file); err != nil {
		cliLogger.Log.Warnf("Could not remove temporary file: %s", err)
	}
}

// formatShellErrOutput formats the error shell output
func formatShellErrOutput(o string) {
	for _, line := range strings.Split(o, "\n") {
		if line != "" {
			cliLogger.Log.Errorf("|> %s", line)
		}
	}
}

// formatShellOutput formats the normal shell output
func formatShellOutput(o string) {
	for _, line := range strings.Split(o, "\n") {
		if line != "" {
			cliLogger.Log.Infof("|> %s", line)
		}
	}
}

const (
	// MigrationMainTPL migration main template
	MigrationMainTPL = `package main

import(
	"os"

	"github.com/bhojpur/web/pkg/client/orm"
	"github.com/bhojpur/web/pkg/client/orm/migration"

	_ "{{DriverRepo}}"
)

func init(){
	orm.RegisterDataBase("default", "{{DBDriver}}","{{ConnStr}}")
}

func main(){
	task := "{{Task}}"
	switch task {
	case "upgrade":
		if err := migration.Upgrade({{LatestTime}}); err != nil {
			os.Exit(2)
		}
	case "rollback":
		if err := migration.Rollback("{{LatestName}}"); err != nil {
			os.Exit(2)
		}
	case "reset":
		if err := migration.Reset(); err != nil {
			os.Exit(2)
		}
	case "refresh":
		if err := migration.Refresh(); err != nil {
			os.Exit(2)
		}
	}
}

`
	// MYSQLMigrationDDL MySQL migration SQL
	MYSQLMigrationDDL = `
CREATE TABLE migrations (
	id_migration int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'surrogate key',
	name varchar(255) DEFAULT NULL COMMENT 'migration name, unique',
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'date migrated or rolled back',
	statements longtext COMMENT 'SQL statements for this migration',
	rollback_statements longtext COMMENT 'SQL statment for rolling back migration',
	status ENUM('update', 'rollback') COMMENT 'update indicates it is a normal migration while rollback means this migration is rolled back',
	PRIMARY KEY (id_migration)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
`
	// POSTGRESMigrationDDL Postgres migration SQL
	POSTGRESMigrationDDL = `
CREATE TYPE migrations_status AS ENUM('update', 'rollback');

CREATE TABLE migrations (
	id_migration SERIAL PRIMARY KEY,
	name varchar(255) DEFAULT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	statements text,
	rollback_statements text,
	status migrations_status
)`
)

// MigrateUpdate does the schema update
func MigrateUpdate(currpath, driver, connStr, dir string) {
	migrate("upgrade", currpath, driver, connStr, dir)
}

// MigrateRollback rolls back the latest migration
func MigrateRollback(currpath, driver, connStr, dir string) {
	migrate("rollback", currpath, driver, connStr, dir)
}

// MigrateReset rolls back all migrations
func MigrateReset(currpath, driver, connStr, dir string) {
	migrate("reset", currpath, driver, connStr, dir)
}

// MigrateRefresh rolls back all migrations and start over again
func MigrateRefresh(currpath, driver, connStr, dir string) {
	migrate("refresh", currpath, driver, connStr, dir)
}
