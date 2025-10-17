package persistence_test

import (
	"errors"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupMockDB() (sqlmock.Sqlmock, *gorm.DB) {
	var err error

	mockDB, mockSql, err := sqlmock.New()
	Expect(err).NotTo(HaveOccurred())
	mock := mockSql

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	Expect(err).NotTo(HaveOccurred())

	return mock, db
}

func CleanupMockDB(mock sqlmock.Sqlmock) {
	Expect(mock.ExpectationsWereMet()).To(Succeed())
}

func MockSuccessfulInsert(mock sqlmock.Sqlmock, tableName string, id interface{}) {
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "` + tableName + `"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	mock.ExpectCommit()
}

func MockSuccessfulUpdate(mock sqlmock.Sqlmock, tableName string) {
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "` + tableName + `" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

func MockDuplicateKeyError(mock sqlmock.Sqlmock, tableName, constrainName string) {
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "` + tableName + `"`)).
		WillReturnError(&pgconn.PgError{
			Code:           "23505",
			ConstraintName: constrainName,
		})
	mock.ExpectRollback()
}

func MockDBError(mock sqlmock.Sqlmock, operation string) {
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(operation)).
		WillReturnError(errors.New("database connection failed"))
	mock.ExpectRollback()
}

func MockFindByID(mock sqlmock.Sqlmock, tableName string, id any, rows *sqlmock.Rows) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "`+tableName+`" WHERE id = $1`)).
		WithArgs(id, 1).WillReturnRows(rows)
}

func MockFindAll(mock sqlmock.Sqlmock, tableName string, rows *sqlmock.Rows) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "` + tableName + `"`)).WillReturnRows(rows)
}

func MockNotFound(mock sqlmock.Sqlmock, tableName string, id any) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "`+tableName+`" WHERE id = $1`)).
		WithArgs(id, 1).WillReturnError(gorm.ErrRecordNotFound)
}

func MockSoftDelete(mock sqlmock.Sqlmock, tableName string, id any) {
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "`+tableName+`" SET "deleted_at"=$1 WHERE id = $2`)).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}
