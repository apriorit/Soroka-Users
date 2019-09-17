package service

import (
	"context"

	"github.com/go-kit/kit/log"

	cfg "github.com/Soroka-EDMS/svc/users/pkgs/config"
	"github.com/Soroka-EDMS/svc/users/pkgs/constants"
	"github.com/Soroka-EDMS/svc/users/pkgs/errors"
	m "github.com/Soroka-EDMS/svc/users/pkgs/models"
)

type usersServiceStub struct {
	Logger log.Logger
	Db     m.UsersDatabase
}

//NewUsersService creates users service stub that enclosing logger and users database
func NewUsersService(db m.UsersDatabase) m.UsersService {
	return &usersServiceStub{
		Logger: cfg.GetLogger().Logger,
		Db:     db,
	}
}

//Build creates session service with middleware
func Build(logger log.Logger, db m.UsersDatabase) m.UsersService {
	var svc m.UsersService
	{
		svc = NewUsersService(db)
		svc = LoggingMiddleware(logger)(svc)
	}

	return svc
}

func (uStub usersServiceStub) CheckAuth(cntx context.Context, creds m.UserCredentials) (err error) {
	return uStub.Db.CheckAuth(cntx, creds.Username, creds.Password)
}

func (uStub usersServiceStub) ChangeRole(cntx context.Context, request m.ChangeRole) (res m.ChangeUsers, err error) {
	ids := request.Ids
	role := request.Role
	changed := 0
	notFound := 0
	notFoundIds := make(map[int]int)
	changedIds := make(map[int]int)

	for _, id := range ids {
		err = uStub.Db.ChangeRole(cntx, id, role)

		if err != nil {
			if err == errors.ErrUserDoesNotExist {
				notFoundIds[notFound] = id
				notFound++
				continue
			}
			if err == errors.ErrInvalidRole {
				return res, err
			}
		}

		changedIds[changed] = id
		changed++
	}

	changedSlice := make([]int, len(changedIds))
	notFoundSlice := make([]int, len(notFoundIds))
	MapToSlice(notFoundIds, notFoundSlice)
	MapToSlice(changedIds, changedSlice)
	res.Changed = changedSlice
	res.Not_found = notFoundSlice

	if notFound != 0 && changed != 0 {
		res.Message = constants.ChangedPartialy
	} else if changed == 0 {
		res.Message = constants.NoRecordChanged
	} else {
		res.Message = constants.ChangedSuccessfully
	}

	return res, nil
}

func (uStub usersServiceStub) GetUserList(cntx context.Context, request m.UsersList) (res m.UsersListResp, err error) {
	resp, err := uStub.Db.GetUsersList(cntx, request.Offset, request.Limit, request.Sort, request.Order)
	return *resp, err
}
func (uStub usersServiceStub) GetUserProfile(cntx context.Context, request m.UserProfileReq) (res m.UserProfileResp, err error) {
	resp, err := uStub.Db.GetUserProfile(cntx, request.Id, request.Email)
	return *resp, err
}
func (uStub usersServiceStub) DisableUsers(cntx context.Context, request m.UsersChangeStatus) (res m.ChangeUsers, err error) {
	return uStub.ChangeStatus(cntx, request, false)
}
func (uStub usersServiceStub) EnableUsers(cntx context.Context, request m.UsersChangeStatus) (res m.ChangeUsers, err error) {
	return uStub.ChangeStatus(cntx, request, true)
}

func (uStub usersServiceStub) ChangeStatus(cntx context.Context, request m.UsersChangeStatus, newStatus bool) (res m.ChangeUsers, err error) {
	ids := request.Ids
	changed := 0
	notFound := 0
	notFoundIds := make(map[int]int)
	changedIds := make(map[int]int)

	for _, id := range ids {
		err = uStub.Db.ChangeUserStatus(cntx, id, newStatus)

		if err != nil {
			if err == errors.ErrUserDoesNotExist {
				notFoundIds[notFound] = id
				notFound++
				continue
			}
		}

		changedIds[changed] = id
		changed++
	}

	changedSlice := make([]int, len(changedIds))
	notFoundSlice := make([]int, len(notFoundIds))
	MapToSlice(notFoundIds, notFoundSlice)
	MapToSlice(changedIds, changedSlice)
	res.Changed = changedSlice
	res.Not_found = notFoundSlice

	if notFound != 0 && changed != 0 {
		res.Message = constants.ChangedPartialy
	} else if changed == 0 {
		res.Message = constants.NoRecordChanged
	} else {
		res.Message = constants.ChangedSuccessfully
	}

	return res, nil
}
