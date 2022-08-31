package settings

import (
	"context"
	"fake-uim/entity"
	"fmt"
	"github.com/goccy/go-json"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

var users []entity.User

func InitUserData() {
	content, err := ioutil.ReadFile("./users.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	log.Printf("User data [%s]", string(content))

	_ = json.Unmarshal(content, &users)

	for _, user := range users {
		policyKey := fmt.Sprintf("user:%s:tenant:%s:policy", user.Uid, user.Tid)

		var rsArr []string
		resources := user.Resources
		for _, rs := range resources {
			rsArr = append(rsArr, rs)
			rsArr = append(rsArr, "")
		}
		log.Printf("Init policy Key=[%s],Val=[%s]", policyKey, strings.Join(rsArr, ","))
		RedisCli().HSet(context.Background(), policyKey, rsArr)
		RedisCli().Expire(context.Background(), policyKey, time.Hour*3)
	}

}

func Users() []entity.User {
	return users
}

func CheckUser(phone string, pwd string) (bool, entity.User) {
	for _, user := range users {
		if user.Phone == phone && user.Password == pwd {
			return true, user
		}
	}
	return false, entity.User{}
}
