package database

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/message"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"
	migrate "github.com/rubenv/sql-migrate"
)

const (
	// messageMigrationsApplied message for applied migrations
	messageMigrationsApplied = "Applied %d migrations"
)

const (
	// migrationDir path of the migrations
	migrationDir = "migrations/database"
	// migrationTable name of the migrations table
	migrationTable = "migrations"
)

// Migration has the migration information
type Migration struct {
	// Name
	name string
	// Connection
	conn *sql.DB
	// Dialect
	dialect string
	// Disables
	disabled bool
	// Directories
	dirs []string
	// Handler
	handler *migrate.MigrationSet
}

// NewMigration it created a new migration
func NewMigration(service string, conn *sql.DB, dialect string, disabled bool) Migration {
	return Migration{
		name:     service,
		conn:     conn,
		dialect:  dialect,
		disabled: disabled,
		dirs:     []string{migrationDir},
		handler: &migrate.MigrationSet{
			TableName:          migrationTable,
			DisableCreateTable: false,
		},
	}
}

// Run it runs the migrations
func (m *Migration) Run() []error {
	if m.disabled {
		return nil
	}

	toRun, err := m.getMigrationsToRun(m.conn, m)
	if err != nil {
		return []error{errors.ErrorMigrationErrorTemplate().Formats(err.Error())}
	}

	missingFiles, err := m.checkForMissingFiles(m.conn, toRun)
	if err != nil {
		return []error{errors.ErrorMigrationErrorTemplate().Formats(err.Error())}
	}

	if len(missingFiles) > 0 {
		var errs []error
		for _, mf := range missingFiles {
			errs = append(errs, errors.ErrorMigrationMissingErrorTemplate().Formats(mf))
		}
		return errs
	}

	count, err := m.doMigrate(m.conn)
	if err != nil {
		return []error{errors.ErrorMigrationErrorTemplate().Formats(err.Error())}
	}

	message.Message(m.name, fmt.Sprintf(messageMigrationsApplied, count))

	return nil
}

// getMigrationsToRun migration list to be applied
func (m *Migration) getMigrationsToRun(conn *sql.DB, src migrate.MigrationSource) ([]*migrate.PlannedMigration, error) {
	migSet := &migrate.MigrationSet{}
	toRun, _, err := migSet.PlanMigration(conn, m.dialect, src, migrate.Up, 0)
	return toRun, err
}

// checkForMissingFiles it checks for missing migration files
func (m *Migration) checkForMissingFiles(conn *sql.DB, files []*migrate.PlannedMigration) (missingFiles []string, err error) {
	records, err := m.handler.GetMigrationRecords(conn, m.dialect)
	if err != nil {
		return missingFiles, err
	}

	// check for missing migrations files
	for _, r := range records {
		found := false
		for _, f := range files {
			if f.Id == r.Id {
				found = true
				break
			}
		}
		if !found {
			missingFiles = append(missingFiles, r.Id)
		}
	}

	return missingFiles, err
}

// doMigrate executes the migrations
func (m *Migration) doMigrate(conn *sql.DB) (int, error) {
	return m.handler.ExecMax(conn, m.dialect, m, migrate.Up, 0)
}

// FindMigrations finds the migrations and sorts
func (m *Migration) FindMigrations() ([]*migrate.Migration, error) {
	var result []*migrate.Migration
	src := migrate.FileMigrationSource{}
	for _, dir := range m.dirs {
		src.Dir = dir
		ms, err := src.FindMigrations()
		if err != nil {
			return []*migrate.Migration{}, err
		}

		result = append(result, ms...)
	}

	// sort found migrations to avoid duplicate migration issues
	sort.Sort(sortById(result))

	return result, nil
}

// sortById for sorting migrations
type sortById []*migrate.Migration

// Len migrations length
func (s sortById) Len() int { return len(s) }

// Swap swaps position i to j in the array
func (s sortById) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less checks if the i < j
func (s sortById) Less(i, j int) bool { return s[i].Less(s[j]) }
