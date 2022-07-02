package todos

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/KejawenLab/bima/v3"
	"github.com/KejawenLab/bima/v3/loggers"
	"github.com/KejawenLab/bima/v3/paginations"
	"github.com/KejawenLab/bima/v3/utils"
	grpcs "github.com/KejawenLab/skeleton/v3/protos/builds"
	"github.com/jinzhu/copier"
)

type Module struct {
	*bima.Module
	Model *Todo
	grpcs.UnimplementedTodosServer
}

func (m *Module) Done(ctx context.Context, r *grpcs.Todo) (*grpcs.TodoResponse, error) {
	ctx = context.WithValue(ctx, "scope", "todo")
	v := m.Model
	copier.Copy(v, r)
	if err := m.Handler.Bind(v, r.Id); err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.NotFound, fmt.Sprintf("Data with ID '%s' not found.", r.Id))
	}

	v.Done = true
	r.Done = true
	if err := m.Handler.Update(v, v.Id); err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.Internal, "Internal server error")
	}

	m.Cache.Invalidate(r.Id)

	return &grpcs.TodoResponse{
		Todo: r,
	}, nil
}

func (m *Module) GetPaginated(ctx context.Context, r *grpcs.Pagination) (*grpcs.TodoPaginatedResponse, error) {
	records := []*grpcs.Todo{}
	reqeust := paginations.Request{}

	m.Paginator.Model = *m.Model
	m.Paginator.Table = m.Model.TableName()

	copier.Copy(&reqeust, r)
	m.Paginator.Handle(reqeust)

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
	v := m.Model
	copier.Copy(v, r)

	if message, err := utils.Validate(v); err != nil {
		loggers.Logger.Error(ctx, message)

		return nil, status.Error(codes.InvalidArgument, message)
	}

	if err := m.Handler.Create(v); err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.Internal, "Internal server error")
	}

	r.Id = v.Id

	return &grpcs.TodoResponse{
		Todo: r,
	}, nil
}

func (m *Module) Update(ctx context.Context, r *grpcs.Todo) (*grpcs.TodoResponse, error) {
	ctx = context.WithValue(ctx, "scope", "todo")
	v := m.Model
	hold := *v
	copier.Copy(v, r)

	if message, err := utils.Validate(v); err != nil {
		loggers.Logger.Error(ctx, message)

		return nil, status.Error(codes.InvalidArgument, message)
	}

	if err := m.Handler.Bind(&hold, r.Id); err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.NotFound, fmt.Sprintf("Data with ID '%s' not found.", r.Id))
	}

	v.Id = r.Id
	v.SetCreatedBy(hold.CreatedBy.String)
	v.SetCreatedAt(hold.CreatedAt.Time)
	if err := m.Handler.Update(v, v.Id); err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.Internal, "Internal server error")
	}

	m.Cache.Invalidate(r.Id)

	return &grpcs.TodoResponse{
		Todo: r,
	}, nil
}

func (m *Module) Get(ctx context.Context, r *grpcs.Todo) (*grpcs.TodoResponse, error) {
	ctx = context.WithValue(ctx, "scope", "todo")
	v := *m.Model
	if data, found := m.Cache.Get(r.Id); found {
		v = data.(Todo)
	} else {
		if err := m.Handler.Bind(&v, r.Id); err != nil {
			loggers.Logger.Error(ctx, err.Error())

			return nil, status.Error(codes.NotFound, fmt.Sprintf("Data with ID '%s' not found.", r.Id))
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
	v := m.Model
	if err := m.Handler.Bind(v, r.Id); err != nil {
		loggers.Logger.Error(ctx, err.Error())

		return nil, status.Error(codes.NotFound, fmt.Sprintf("Data with ID '%s' not found.", r.Id))
	}

	m.Handler.Delete(v, r.Id)
	m.Cache.Invalidate(r.Id)

	return &grpcs.TodoResponse{
		Todo: nil,
	}, nil
}
