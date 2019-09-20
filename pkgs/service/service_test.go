package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-kit/kit/log"

	"github.com/Soroka-EDMS/svc/users/pkgs/constants"
	"github.com/Soroka-EDMS/svc/users/pkgs/db"
	"github.com/Soroka-EDMS/svc/users/pkgs/models"
	"github.com/Soroka-EDMS/svc/users/pkgs/stub"
)

func TestChangeRole_AllUsersExist(t *testing.T) {
	//GetStub database
	dbs, err := db.Connection(log.NewNopLogger(), "stub")
	dbase := dbs.(*db.UserDbStub)
	assert.NoError(t, err)
	//Check for user role before changing
	assert.Equal(t, dbase.Profiles[4].Role.Name, stub.RoleName_reducedUser)
	assert.Equal(t, dbase.Profiles[4].Role.Mask, stub.RoleMask_reducedUser)

	//Build service
	svc := Build(log.NewNopLogger(), dbase)

	//Prepare request
	req := models.ChangeRole{
		Role: "admin",
		Ids:  []int{4},
	}

	//Prepare expected response
	expectedResponse := models.ChangeUsers{
		Changed:   []int{4},
		Not_found: make([]int, 0),
		Message:   constants.ChangedSuccessfully,
	}

	//Call service ChangeRole
	resp, err := svc.ChangeRole(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, resp)

	//Check for expected value after the user role has been changed
	assert.Equal(t, dbase.Profiles[4].Role.Name, stub.RoleName_admin)
	assert.Equal(t, dbase.Profiles[4].Role.Mask, stub.RoleMask_admin)
}

func TestChangeRole_AllUsersDontExist(t *testing.T) {
	dbs, err := db.Connection(log.NewNopLogger(), "stub")
	dbase := dbs.(*db.UserDbStub)
	assert.NoError(t, err)
	//Build service
	svc := Build(log.NewNopLogger(), dbase)

	//Prepare request
	req := models.ChangeRole{
		Role: "admin",
		Ids:  []int{11, 25},
	}

	//Prepare expected response
	expectedResponse := models.ChangeUsers{
		Changed:   make([]int, 0),
		Not_found: []int{11, 25},
		Message:   constants.NoRecordChanged,
	}

	//Call service ChangeRole
	resp, err := svc.ChangeRole(context.Background(), req)
	assert.NoError(t, err)

	assert.Equal(t, expectedResponse, resp)
}

func TestChangeRole_PartiallyExist(t *testing.T) {
	dbs, err := db.Connection(log.NewNopLogger(), "stub")
	dbase := dbs.(*db.UserDbStub)
	assert.NoError(t, err)
	assert.Equal(t, dbase.Profiles[2].Role.Name, stub.RoleName_user)
	assert.Equal(t, dbase.Profiles[2].Role.Mask, stub.RoleMask_user)
	assert.Equal(t, dbase.Profiles[4].Role.Name, stub.RoleName_reducedUser)
	assert.Equal(t, dbase.Profiles[4].Role.Mask, stub.RoleMask_reducedUser)
	//Build service
	svc := Build(log.NewNopLogger(), dbase)

	//Prepare request
	req := models.ChangeRole{
		Role: "admin",
		Ids:  []int{4, 32, 44, 2},
	}

	//Prepare expected response
	expectedResponse := models.ChangeUsers{
		Changed:   []int{4, 2},
		Not_found: []int{32, 44},
		Message:   constants.ChangedPartialy,
	}

	//Call service ChangeRole
	resp, err := svc.ChangeRole(context.Background(), req)
	assert.NoError(t, err)

	assert.Equal(t, expectedResponse, resp)
	assert.Equal(t, dbase.Profiles[2].Role.Name, stub.RoleName_admin)
	assert.Equal(t, dbase.Profiles[2].Role.Mask, stub.RoleMask_admin)
	assert.Equal(t, dbase.Profiles[4].Role.Name, stub.RoleName_admin)
	assert.Equal(t, dbase.Profiles[4].Role.Mask, stub.RoleMask_admin)
}

func TestGetUserProfile_ByEmail(t *testing.T) {
	//Connect to stub database
	db, err := db.Connection(log.NewNopLogger(), "stub")
	assert.NoError(t, err)
	//Build service
	svc := Build(log.NewNopLogger(), db)

	req := models.UserProfileReq{
		Email: "lindsay2017@edms.com",
	}

	expectedResponse := models.UserProfile{
		First_name:    "Doris",
		Last_name:     "Hooper",
		Email:         "lindsay2017@edms.com",
		Phone:         "+38(067)421-73-92",
		Location:      "Hays, Kansas(KS), 67601",
		Position:      "Procurement Specialist",
		Status:        true,
		Creation_date: 1567590560,
		Role: models.UserRole{
			Name: "ordinaryUser",
			Mask: 25524,
		},
	}

	resp, err := svc.GetUserProfile(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, resp, expectedResponse)
}

func TestDisableUsers_UserEnabled(t *testing.T) {
	//Connect to stub database
	dbs, err := db.Connection(log.NewNopLogger(), "stub")
	dbase := dbs.(*db.UserDbStub)
	assert.NoError(t, err)
	//Build service
	svc := Build(log.NewNopLogger(), dbase)

	req := models.UsersChangeStatus{
		Ids: []int{2},
	}

	//Prepare expected response
	expectedResponse := models.ChangeUsers{
		Changed:   []int{2},
		Not_found: make([]int, 0),
		Message:   constants.ChangedSuccessfully,
	}

	resp, err := svc.DisableUsers(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, resp, expectedResponse)
	assert.False(t, dbase.Profiles[2].Status)
}

func TestGetUserList_GetAll(t *testing.T) {
	//Connect to stub database
	db, err := db.Connection(log.NewNopLogger(), "stub")
	assert.NoError(t, err)
	//Build service
	svc := Build(log.NewNopLogger(), db)

	req := models.UsersList{
		Offset: 1,
		Limit:  2,
		Sort:   "user_name",
		Order:  "down",
	}

	resp, err := svc.GetUserList(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, len(resp.Users), req.Limit-req.Offset)
}
