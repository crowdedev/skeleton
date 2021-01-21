package services

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	configs "github.com/crowdeco/skeleton/configs"
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
	svc  configs.Service
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.db, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(s.T(), err)

	s.svc = NewTodoService(s.db)
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestCreate() {
	id := uuid.New().String()
	name := "Todo 1"

	s.mock.ExpectExec("INSERT INTO `todos`").
		WithArgs(id, name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	todo := &models.Todo{}
	todo.Name = name
	s.svc.Create(todo, id)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *Suite) TestUpdate() {
	id := uuid.New().String()
	updated_at := time.Now().UTC()
	updated_by := ""
	name := "Todo 1"

	rows := sqlmock.NewRows([]string{"id", "updated_at", "updated_by", "name"}).
		AddRow(id, updated_at, updated_by, name)

	s.mock.ExpectQuery("SELECT (.)+ FROM `todos` WHERE").
		WithArgs(id).
		WillReturnRows(rows)
	s.mock.ExpectExec("UPDATE `todos`").
		WithArgs(configs.AnyTime{}, "someone_id", "Todo 2", id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	todo := &models.Todo{}
	todo.Name = "Todo 2"
	todo.UpdatedBy = "someone_id"
	s.svc.Update(todo, id)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *Suite) TestBind() {
	id := uuid.New().String()
	name := "Todo 1"

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(id, name)

	s.mock.ExpectQuery("SELECT (.)+ FROM `todos` WHERE").
		WithArgs(id).
		WillReturnRows(rows)

	s.svc.Bind(&models.Todo{}, id)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *Suite) TestSoftDelete() {
	m := models.Todo{}

	if m.IsSoftDelete() {
		id := uuid.New().String()
		deleted_at := gorm.DeletedAt{}
		deleted_by := ""
		name := "Todo 1"

		rows := sqlmock.NewRows([]string{"id", "deleted_at", "deleted_by", "name"}).
			AddRow(id, deleted_at, deleted_by, name)

		s.mock.ExpectQuery("SELECT (.)+ FROM `todos` WHERE").
			WithArgs(id).
			WillReturnRows(rows)
		s.mock.ExpectExec("UPDATE `todos`").
			WithArgs(configs.AnyTime{}, "", id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		s.svc.Delete(&m, id)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func (s *Suite) TestHardDelete() {
	m := models.Todo{}

	if !m.IsSoftDelete() {
		id := uuid.New().String()
		name := "Todo 1"

		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(id, name)

		s.mock.ExpectQuery("SELECT (.)+ FROM `todos` WHERE").
			WithArgs(id).
			WillReturnRows(rows)
		s.mock.ExpectExec("DELETE FROM `todos`").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		s.svc.Delete(&m, id)

		if err := s.mock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}
