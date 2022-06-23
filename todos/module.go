package todos

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/KejawenLab/bima/v3"
	"github.com/KejawenLab/bima/v3/configs"
	grpcs "github.com/KejawenLab/skeleton/protos/builds"
	"github.com/KejawenLab/skeleton/todos/models"
	"github.com/KejawenLab/skeleton/todos/validations"
	"github.com/jinzhu/copier"
)

type Module struct {
	*bima.Module
	Model     *models.Todo
	Validator *validations.Todo
	grpcs.UnimplementedTodosServer
}

func (m *Module) GetPaginated(ctx context.Context, r *grpcs.Pagination) (*grpcs.TodoPaginatedResponse, error) {
	m.Logger.Debug(context.WithValue(ctx, "scope", "todo"), fmt.Sprintf("%+v", r))
	records := []*grpcs.Todo{}
	model := models.Todo{}
	m.Paginator.Model = model
	m.Paginator.Table = model.TableName()

	copier.Copy(m.Request, r)
	m.Paginator.Handle(m.Request)

	metadata := m.Handler.Paginate(*m.Paginator, &records)

	return &grpcs.TodoPaginatedResponse{
		Data: records,
		Meta: &grpcs.PaginationMetadata{
			Page:     int32(metadata.Page),
			Previous: int32(metadata.Previous),
			Next:     int32(metadata.Next),
			Limit:    int32(metadata.Limit),
			Total:    int32(metadata.Total),
		},
	}, nil
}

func (m *Module) Create(ctx context.Context, r *grpcs.Todo) (*grpcs.TodoResponse, error) {
	ctx = context.WithValue(ctx, "scope", "todo")
	m.Logger.Debug(ctx, fmt.Sprintf("%+v", r))

	v := m.Model
	copier.Copy(v, r)

	if ok, err := m.Validator.Validate(v); !ok {
		m.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := m.Handler.Create(v); err != nil {
		m.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.Internal, err.Error())
	}

	r.Id = v.Id

	return &grpcs.TodoResponse{
		Todo: r,
	}, nil
}

func (m *Module) Update(ctx context.Context, r *grpcs.Todo) (*grpcs.TodoResponse, error) {
	ctx = context.WithValue(ctx, "scope", "todo")
	m.Logger.Debug(ctx, fmt.Sprintf("%+v", r))

	v := m.Model
	hold := *v
	copier.Copy(v, r)

	if ok, err := m.Validator.Validate(v); !ok {
		m.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := m.Handler.Bind(&hold, r.Id); err != nil {
		msg := fmt.Sprintf("Data with ID '%s' not found.", r.Id)
		m.Logger.Error(ctx, msg)

		return nil, status.Error(codes.NotFound, msg)
	}

	v.Id = r.Id
	v.SetCreatedBy(&configs.User{Id: hold.CreatedBy.String})
	v.SetCreatedAt(hold.CreatedAt.Time)
	if err := m.Handler.Update(v, v.Id); err != nil {
		m.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.Internal, err.Error())
	}

	m.Cache.Invalidate(r.Id)

	return &grpcs.TodoResponse{
		Todo: r,
	}, nil
}

func (m *Module) Get(ctx context.Context, r *grpcs.Todo) (*grpcs.TodoResponse, error) {
	ctx = context.WithValue(ctx, "scope", "todo")
	m.Logger.Debug(ctx, fmt.Sprintf("%+v", r))

	var v models.Todo
	if data, found := m.Cache.Get(r.Id); found {
		v = data.(models.Todo)
	} else {
		if err := m.Handler.Bind(&v, r.Id); err != nil {
			msg := fmt.Sprintf("Data with ID '%s' not found.", r.Id)
			m.Logger.Error(ctx, msg)

			return nil, status.Error(codes.NotFound, msg)
		}

		m.Cache.Set(r.Id, v)
	}

	copier.Copy(r, &v)

	return &grpcs.TodoResponse{
		Todo: r,
	}, nil
}

func (m *Module) Delete(ctx context.Context, r *grpcs.Todo) (*grpcs.TodoResponse, error) {
	ctx = context.WithValue(ctx, "scope", "todo")
	m.Logger.Debug(ctx, fmt.Sprintf("%+v", r))

	v := m.Model
	if err := m.Handler.Bind(v, r.Id); err != nil {
		msg := fmt.Sprintf("Data with ID '%s' not found.", r.Id)
		m.Logger.Error(ctx, msg)

		return nil, status.Error(codes.NotFound, msg)
	}

	m.Handler.Delete(v, r.Id)
	m.Cache.Invalidate(r.Id)

	return &grpcs.TodoResponse{
		Todo: nil,
	}, nil
}
