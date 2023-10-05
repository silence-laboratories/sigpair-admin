package main

import (
	"fmt"
	sdk "github.com/silence-laboratories/sigpair-admin"
)

func main() {
	client := sdk.NewClient("http://localhost:8080", "SUPER_SECRET_ADMIN_TOKEN")
	uid, err := client.CreateUser("sdf")
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	fmt.Printf("Created user with id: %d\n", uid)
	token, err := client.GenUserToken(uid, 3600)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	fmt.Printf("Created user token %s", token)

}
