package persistence_test

import (
	"errors"
	"ev-warranty-go/internal/apperrors"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupMockDB() (sqlmock.Sqlmock, *gorm.DB) {
	mockDB, mock, err := sqlmock.New()
	Expect(err).NotTo(HaveOccurred())

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
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

func MockDuplicateKeyError(mock sqlmock.Sqlmock, tableName, constraintName string) {
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "` + tableName + `"`)).
		WillReturnError(&pgconn.PgError{
			Code:           "23505",
			ConstraintName: constraintName,
		})
	mock.ExpectRollback()
}

func MockInsertError(mock sqlmock.Sqlmock, tableName string) {
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "` + tableName + `"`)).
		WillReturnError(errors.New("database connection failed"))
	mock.ExpectRollback()
}

func MockUpdateError(mock sqlmock.Sqlmock, tableName string) {
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "` + tableName + `" SET`)).
		WillReturnError(errors.New("database connection failed"))
	mock.ExpectRollback()
}

func MockDeleteError(mock sqlmock.Sqlmock, tableName string) {
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "` + tableName + `" SET "deleted_at"=$1 WHERE id = $2`)).
		WillReturnError(errors.New("database connection failed"))
	mock.ExpectRollback()
}

func MockQueryError(mock sqlmock.Sqlmock, query string) {
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnError(errors.New("database connection failed"))
}

func MockFindByID(mock sqlmock.Sqlmock, tableName string, id any, rows *sqlmock.Rows) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "`+tableName+`" WHERE id = $1`)).
		WithArgs(id, 1).
		WillReturnRows(rows)
}

func MockFindAll(mock sqlmock.Sqlmock, tableName string, rows *sqlmock.Rows) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "` + tableName + `"`)).
		WillReturnRows(rows)
}

func MockNotFound(mock sqlmock.Sqlmock, tableName string, id any) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "`+tableName+`" WHERE id = $1`)).
		WithArgs(id, 1).
		WillReturnError(gorm.ErrRecordNotFound)
}

func MockSoftDelete(mock sqlmock.Sqlmock, tableName string, id any) {
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "`+tableName+`" SET "deleted_at"=$1 WHERE id = $2`)).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

func ExpectAppError(err error, expectedCode string) {
	GinkgoHelper()
	Expect(err).To(HaveOccurred())
	var appErr *apperrors.AppError
	Expect(errors.As(err, &appErr)).To(BeTrue(), "error should be an AppError")
	Expect(appErr.ErrorCode).To(Equal(expectedCode), "error code should match")
}
