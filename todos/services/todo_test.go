package services

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	models "github.com/crowdeco/skeleton/todos/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(s.T(), err)

	s.db, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(s.T(), err)
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestCreate() {
	id := uuid.New().String()
	s.mock.ExpectExec("INSERT INTO `todos` (`id`,`name`) VALUES (?,?)").
		WithArgs(id, "Todo 1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	service := NewTodoService(s.db)
	service.Create(&models.Todo{Name: "Todo 1"}, id)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}
