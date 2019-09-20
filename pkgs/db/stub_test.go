package db

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Soroka-EDMS/svc/users/pkgs/config"
	"github.com/Soroka-EDMS/svc/users/pkgs/models"
	"github.com/Soroka-EDMS/svc/users/pkgs/stub"
)

func TestPrepareStubDatabase(t *testing.T) {
	var err error
	//Prepare profiles for further comparison
	assert.NoError(t, err)

	profileRawAdmin := []byte(stub.AdminProfile)
	profileAdmin := models.UserProfile{}
	err = json.Unmarshal(profileRawAdmin, &profileAdmin)

	assert.NoError(t, err)

	profileRawUser := []byte(stub.UserProfile)
	profileUser := models.UserProfile{}
	err = json.Unmarshal(profileRawUser, &profileUser)

	assert.NoError(t, err)

	profileRawOrdinaryUser := []byte(stub.OrdinaryUserProfile)
	profileOrdinaryUser := models.UserProfile{}
	err = json.Unmarshal(profileRawOrdinaryUser, &profileOrdinaryUser)

	assert.NoError(t, err)

	profileRawReducedUser := []byte(stub.ReducedUserProfile)
	profileReducedUser := models.UserProfile{}
	err = json.Unmarshal(profileRawReducedUser, &profileReducedUser)

	assert.NoError(t, err)

	//Prepare stub database
	db := UserDbStub{}
	_, err = PrepareStubDatabase(&db, config.GetLogger().Logger)

	//Check for valid database configuration
	assert.NoError(t, err)
	assert.Equal(t, db.Creds[1].Username, stub.UserName_admin)
	assert.Equal(t, db.Creds[1].Password, stub.Password_admin)
	assert.Equal(t, db.Creds[2].Username, stub.UserName_user)
	assert.Equal(t, db.Creds[2].Password, stub.Password_user)
	assert.Equal(t, db.Creds[3].Username, stub.UserName_ordinaryUser)
	assert.Equal(t, db.Creds[3].Password, stub.Password_ordinaryUser)
	assert.Equal(t, db.Creds[4].Username, stub.UserName_reducedUser)
	assert.Equal(t, db.Creds[4].Password, stub.Password_reducedUser)

	assert.Equal(t, db.Roles[1].Name, stub.RoleName_admin)
	assert.Equal(t, db.Roles[1].Mask, stub.RoleMask_admin)
	assert.Equal(t, db.Roles[2].Name, stub.RoleName_user)
	assert.Equal(t, db.Roles[2].Mask, stub.RoleMask_user)
	assert.Equal(t, db.Roles[3].Name, stub.RoleName_ordinaryUser)
	assert.Equal(t, db.Roles[3].Mask, stub.RoleMask_ordinaryUser)
	assert.Equal(t, db.Roles[4].Name, stub.RoleName_reducedUser)
	assert.Equal(t, db.Roles[4].Mask, stub.RoleMask_reducedUser)

	assert.Equal(t, db.Profiles[1], profileAdmin)
	assert.Equal(t, db.Profiles[2], profileUser)
	assert.Equal(t, db.Profiles[3], profileOrdinaryUser)
	assert.Equal(t, db.Profiles[4], profileReducedUser)

	assert.NotEqual(t, len(db.Creds), 0)
	assert.NotEqual(t, len(db.Profiles), 0)
	assert.NotEqual(t, len(db.Roles), 0)
}
