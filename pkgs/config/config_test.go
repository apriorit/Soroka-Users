package config

import (
	"testing"

	c "github.com/Soroka-EDMS/svc/users/pkgs/constants"
)

//Test for correct configuration of priveldge map for Users service
func TestGetPriveledgeConfig_UsersPriveledge(t *testing.T) {
	//Prepare existance map, priveledge map
	existanceMap := make(map[string]bool)
	PrepareExistanceMap(existanceMap)
	privMap := GetPriveledges().Priveledges

	var expectedValue int
	for key, value := range privMap {
		switch key {
		case "changeRole":
			existanceMap["changeRole"] = true
			expectedValue = c.ChangeRole
		case "userList":
			existanceMap["userList"] = true
			expectedValue = c.QueryUsers
		case "userProfile":
			existanceMap["userProfile"] = true
			expectedValue = c.GetProfile
		case "changeUserStatus":
			existanceMap["changeUserStatus"] = true
			expectedValue = c.ChangeUserStatus
		}

		if value != expectedValue {
			t.Errorf("key = %s, priveledge = %d, expected: %d", key, expectedValue, value)
		}
	}

	if len(existanceMap) != c.NumOfPriveledges {
		t.Errorf("Amount of proveledges: %d, expected amount: %d", len(existanceMap), c.NumOfPriveledges)
	}
}

func PrepareExistanceMap(exMap map[string]bool) {
	exMap["changeRole"] = false
	exMap["userList"] = false
	exMap["userProfile"] = false
	exMap["changeUserStatus"] = false
}
