package stub

const (
	DefaultUserCreds string = `{
		"username": "default",
		"password": "default"
	}`

	AdminCreds string = `{
		"username": "gladys.champl@edms.com",
		"password": "a@m1n"
	}`

	UserCreds string = `{
		"username": "percival1987@edms.com",
		"password": "@Sr!1"
	}`

	OrdinaryUserCreds string = `{
		"username": "lindsay2017@edms.com",
		"password": "@Sua1pwd"
	}`

	ReducedUserCreds string = `{
		"username": "marianne.cart@edms.com",
		"password": "Canguro!1"
	}`

	AdminRole string = `{
		"name": "admin",
		"mask": 32767
	}`

	DefaultRole string = `{
		"name": "default",
		"mask": 0
	}`

	UserRole string = `{
		"name": "user",
		"mask": 25588
	}`

	OrdinaryUserRole string = `{
		"name": "ordinaryUser",
		"mask": 25524
	}`

	ReducedUserRole string = `{
		"name": "reducedUser",
		"mask": 8628
	}`

	DefaultProfile string = `{
		"first_name" : "",
		"last_name" : "",
		"email" : "",
		"phone" : "",
		"location" : "",
		"position" : "",
		"status" : true,
		"creation_date" : 0,
		"role": {
			"name" : "",
			"mask" : 0
		}
	}`

	AdminProfile string = `{
		"first_name" : "Gladys",
		"last_name" : "Bannon",
		"email" : "gladys.champl@edms.com",
		"phone" : "+38(045)678-99-99",
		"location" : "Hays, Kansas(KS), 67601",
		"position" : "system administartor",
		"status" : true,
		"creation_date" : 1567590495,
		"role": {
			"name" : "admin",
			"mask" : 32767
		}
	}`

	UserProfile string = `{
		"first_name" : "Jeffrey",
		"last_name" : "Rice",
		"email" : "percival1987@edms.com",
		"phone" : "+38(096)234-11-56",
		"location" : "Hays, Kansas(KS), 67601",
		"position" : "Chief Accountant",
		"status" : false,
		"creation_date" : 1567590522,
		"role": {
			"name" : "user",
			"mask" : 25588
		}
	}`

	OrdinaryUserProfile string = `{
		"first_name" : "Doris",
		"last_name" : "Hooper",
		"email" : "lindsay2017@edms.com",
		"phone" : "+38(067)421-73-92",
		"location" : "Hays, Kansas(KS), 67601",
		"position" : "Procurement Specialist",
		"status" : true,
		"creation_date" : 1567590560,
		"role": {
			"name" : "ordinaryUser",
			"mask" : 25524
		}
	}`

	ReducedUserProfile string = `{
		"first_name" : "Fred",
		"last_name" : "Dunn",
		"email" : "marianne.cart@edms.com",
		"phone" : "+38(067)145-23-82",
		"location" : "Hays, Kansas(KS), 67601",
		"position" : "Accountant",
		"status" : true,
		"creation_date" :  1567590588, 
		"role": {
			"name" : "reducedUser",
			"mask" : 8628
		}
	}`

	AdminToTable string = `{
		"user_name" : "admin",
		"user_id" : 0,
		"role" : "admin",
		"location" : "Hays, Kansas(KS), 67601",
		"email" : "gladys.champl@edms.com",
		"creation_date" : 1567590495,
		"status":"active"
	}`

	UserToTable string = `{
		"user_name" : "user",
		"user_id" : 1,
		"role" : "user",
		"location" : "Hays, Kansas(KS), 67601",
		"email" : "percival1987@edms.com",
		"creation_date" : 1567590522,
		"status":"active"
	}`

	OrdinaryUserToTable string = `{
		"user_name" : "ordinaryUser",
		"user_id" : 2,
		"role" : "ordinaryUser",
		"location" : "Hays, Kansas(KS), 67601",
		"email" : "lindsay2017@edms.com",
		"creation_date" : 1567590560,
		"status":"active"
	}`

	ReducedUserToTable string = `{
		"user_name" : "ordinaryUser",
		"user_id" : 3,
		"role" : "ordinaryUser",
		"location" : "Hays, Kansas(KS), 67601",
		"email" : "marianne.cart@edms.com",
		"creation_date" : 1567590588,
		"status":"active"
	}`

	UsersListWithPagination string = `{
		"users":[
			{
				"user_name" : "admin",
				"user_id" : 0,
				"role" : "admin",
				"location" : "Hays, Kansas(KS), 67601",
				"email" : "gladys.champl@edms.com",
				"creation_date" : 1567590495,
				"status":"active"
			},
			{
				"user_name" : "user",
				"user_id" : 1,
				"role" : "user",
				"location" : "Hays, Kansas(KS), 67601",
				"email" : "percival1987@edms.com",
				"creation_date" : 1567590522,
				"status":"active"
			},
			{
				"user_name" : "ordinaryUser",
				"user_id" : 2,
				"role" : "ordinaryUser",
				"location" : "Hays, Kansas(KS), 67601",
				"email" : "lindsay2017@edms.com",
				"creation_date" : 1567590560,
				"status":"active"
			},
			{
				"user_name" : "ordinaryUser",
				"user_id" : 3,
				"role" : "ordinaryUser",
				"location" : "Hays, Kansas(KS), 67601",
				"email" : "marianne.cart@edms.com",
				"creation_date" : 1567590588,
				"status":"active"
			}
		],
		"pagination":{
			"issued":4,
			"left":0
		}
	}`
)
