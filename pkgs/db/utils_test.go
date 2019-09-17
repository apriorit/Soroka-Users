package db

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Soroka-EDMS/svc/users/pkgs/constants"
	"github.com/Soroka-EDMS/svc/users/pkgs/models"
	"github.com/Soroka-EDMS/svc/users/pkgs/stub"
)

func TestPrepareStubDatabase(t *testing.T) {
	var err error
	//Prepare profiles for further comparison
	assert.NoError(t, err)

	profileRawAdmin := []byte(stub.AdminProfile)
	profileAdmin := models.UserProfileResp{}
	err = json.Unmarshal(profileRawAdmin, &profileAdmin)

	assert.NoError(t, err)

	profileRawUser := []byte(stub.UserProfile)
	profileUser := models.UserProfileResp{}
	err = json.Unmarshal(profileRawUser, &profileUser)

	assert.NoError(t, err)

	profileRawOrdinaryUser := []byte(stub.OrdinaryUserProfile)
	profileOrdinaryUser := models.UserProfileResp{}
	err = json.Unmarshal(profileRawOrdinaryUser, &profileOrdinaryUser)

	assert.NoError(t, err)

	profileRawReducedUser := []byte(stub.ReducedUserProfile)
	profileReducedUser := models.UserProfileResp{}
	err = json.Unmarshal(profileRawReducedUser, &profileReducedUser)

	assert.NoError(t, err)

	//Prepare stub database
	db := models.UsersDb{}
	err = PrepareStubDatabase(&db)

	//Check for valid database configuration
	assert.NoError(t, err)
	assert.Equal(t, db.Creds[1].Username, constants.UserName_admin)
	assert.Equal(t, db.Creds[1].Password, constants.Password_admin)
	assert.Equal(t, db.Creds[2].Username, constants.UserName_user)
	assert.Equal(t, db.Creds[2].Password, constants.Password_user)
	assert.Equal(t, db.Creds[3].Username, constants.UserName_ordinaryUser)
	assert.Equal(t, db.Creds[3].Password, constants.Password_ordinaryUser)
	assert.Equal(t, db.Creds[4].Username, constants.UserName_reducedUser)
	assert.Equal(t, db.Creds[4].Password, constants.Password_reducedUser)

	assert.Equal(t, db.Roles[1].Name, constants.RoleName_admin)
	assert.Equal(t, db.Roles[1].Mask, constants.RoleMask_admin)
	assert.Equal(t, db.Roles[2].Name, constants.RoleName_user)
	assert.Equal(t, db.Roles[2].Mask, constants.RoleMask_user)
	assert.Equal(t, db.Roles[3].Name, constants.RoleName_ordinaryUser)
	assert.Equal(t, db.Roles[3].Mask, constants.RoleMask_ordinaryUser)
	assert.Equal(t, db.Roles[4].Name, constants.RoleName_reducedUser)
	assert.Equal(t, db.Roles[4].Mask, constants.RoleMask_reducedUser)

	assert.Equal(t, db.Profiles[1], profileAdmin)
	assert.Equal(t, db.Profiles[2], profileUser)
	assert.Equal(t, db.Profiles[3], profileOrdinaryUser)
	assert.Equal(t, db.Profiles[4], profileReducedUser)

	assert.NotEqual(t, len(db.Creds), 0)
	assert.NotEqual(t, len(db.Profiles), 0)
	assert.NotEqual(t, len(db.Roles), 0)
}
