package test

import (
	"strconv"
	"testing"
)

func TestLogin(t *testing.T) {
	// dbManager := database.NewDatabaseManager("../data/")
	// authService := services.NewAuthService(dbManager)
	// softwarelist, err := authService.CreateUserSession("admin", "Whm9632396510", "", "")
	// if err != nil {
	// 	t.Errorf("Login failed: %v", err)
	// }
	//t.Logf("Login successful: %v", (softwarelist))

	authority, _ := strconv.ParseUint("000000FF", 16, 64)
	t.Log(authority)

}
