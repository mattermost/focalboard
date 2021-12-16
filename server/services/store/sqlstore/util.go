package sqlstore

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/mattermost/focalboard/server/services/store"
	"github.com/mattermost/focalboard/server/utils"

	"github.com/mattermost/mattermost-server/v6/shared/mlog"
)

func (s *SQLStore) CloseRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		s.logger.Error("error closing MattermostAuthLayer row set", mlog.Err(err))
	}
}

func (s *SQLStore) IsErrNotFound(err error) bool {
	return store.IsErrNotFound(err)
}

func PrepareNewTestDatabase() (dbType string, connectionString string, err error) {
	dbType = strings.TrimSpace(os.Getenv("FB_STORE_TEST_DB_TYPE"))
	if dbType == "" {
		dbType = sqliteDBType
	}

	var dbName string

	if dbType == sqliteDBType {
		connectionString = ":memory:"
	} else if port := strings.TrimSpace(os.Getenv("FB_STORE_TEST_DOCKER_PORT")); port != "" {
		// docker unit tests take priority over any DSN env vars
		var template string
		switch dbType {
		case mysqlDBType:
			template = "%s:mostest@tcp(localhost:%s)/%s?charset=utf8mb4,utf8&writeTimeout=30s"
		case postgresDBType:
			template = "%s:mostest@localhost:%s/%s?sslmode=disable\u0026connect_timeout=10"
		default:
			return "", "", newErrInvalidDBType(dbType)
		}

		connectionString = fmt.Sprintf(template, "root", port, "")

		// create a new database each run
		sqlDB, err := sql.Open(dbType, connectionString)
		if err != nil {
			return "", "", fmt.Errorf("cannot connect to %s database: %w", dbType, err)
		}
		defer sqlDB.Close()

		err = sqlDB.Ping()
		if err != nil {
			return "", "", fmt.Errorf("cannot ping %s database: %w", dbType, err)
		}

		dbName = "testdb_" + utils.NewID(utils.IDTypeNone)[:8]
		_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName))
		if err != nil {
			return "", "", fmt.Errorf("cannot create %s database %s: %w", dbType, dbName, err)
		}

		_, err = sqlDB.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO mmuser;", dbName))
		if err != nil {
			return "", "", fmt.Errorf("cannot grant permissions on %s database %s: %w", dbType, dbName, err)
		}

		connectionString = fmt.Sprintf(template, "mmuser", port, dbName)
	} else {
		// mysql or postgres need a DSN (connection string)
		connectionString = strings.TrimSpace(os.Getenv("FB_STORE_TEST_CONN_STRING"))
	}

	return dbType, connectionString, nil
}

type ErrInvalidDBType struct {
	dbType string
}

func newErrInvalidDBType(dbType string) error {
	return ErrInvalidDBType{
		dbType: dbType,
	}
}

func (e ErrInvalidDBType) Error() string {
	return "unsupported database type: " + e.dbType
}
