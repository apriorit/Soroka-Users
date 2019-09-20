package config

import (
	"testing"

	"github.com/Soroka-EDMS/svc/users/pkgs/stub"
)

//Test for correct configuration of priveldge map for Users service
func TestGetPriveledgeConfig_UsersPriveledge(t *testing.T) {
	//Prepare existance map, priveledge map
	existanceMap := make(map[string]bool)
	PrepareExistanceMap(existanceMap)
	privMap := GetPrivileges().Privileges

	var expectedValue int
	for key, value := range privMap {
		switch key {
		case "changeRole":
			existanceMap["changeRole"] = true
			expectedValue = stub.ChangeRole
		case "userList":
			existanceMap["userList"] = true
			expectedValue = stub.QueryUsers
		case "userProfile":
			existanceMap["userProfile"] = true
			expectedValue = stub.GetProfile
		case "changeUserStatus":
			existanceMap["changeUserStatus"] = true
			expectedValue = stub.ChangeUserStatus
		}

		if value != expectedValue {
			t.Errorf("key = %s, priveledge = %d, expected: %d", key, expectedValue, value)
		}
	}

	if len(existanceMap) != stub.AmountOfPriveledges {
		t.Errorf("Amount of proveledges: %d, expected amount: %d", len(existanceMap), stub.AmountOfPriveledges)
	}
}

func PrepareExistanceMap(exMap map[string]bool) {
	exMap["changeRole"] = false
	exMap["userList"] = false
	exMap["userProfile"] = false
	exMap["changeUserStatus"] = false
}
