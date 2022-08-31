package main

import (
	"fake-uim/entity"
	"fake-uim/settings"
	"fake-uim/util"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"log"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Session struct {
	Uid            string    `json:"uid"`
	Tid            string    `json:"tid"`
	Name           string    `json:"name"`
	AuthorizedPids []string  `json:"authorizedPids"`
	CreateTime     time.Time `json:"sessionCreateTime"`
	ExpireAt       time.Time `json:"expireAt"`
}

func main() {
	settings.InitRdb()
	settings.InitUserData()

	r := gin.Default()

	r.POST("/uim/login", func(c *gin.Context) {
		var user entity.User
		if c.ShouldBind(&user) == nil {
			//Verify user login information
			f, u := settings.CheckUser(user.Phone, user.Password)

			//If the verification passes, the session will be generated
			if f {
				expireAt := time.Now().Add(time.Hour * 72)
				session := Session{
					Uid:            u.Uid,
					Tid:            u.Tid,
					Name:           u.Name,
					AuthorizedPids: u.AuthorizedPids,
					CreateTime:     time.Now(),
					ExpireAt:       expireAt,
				}

				//Persisting session to redis
				val, _ := json.Marshal(session)
				sessionKey := fmt.Sprintf("user:%s:tenant:%s:session", session.Uid, session.Tid)
				log.Printf("The session key is [%s]", sessionKey)
				settings.RedisCli().Set(c, sessionKey, string(val), time.Hour*72)

				c.JSON(200, gin.H{"message": "Login successfully !"})
				return
			}
		}
		c.JSON(401, gin.H{"message": "Login failed , wrong user or password !"})
	})

	r.POST("/uim/auth", func(c *gin.Context) {
		sessionKey := c.GetHeader("SEEINER-SEC-TOKEN")
		//Can be added by APISIX proxy-rewrite plugin or NGINX
		prs := c.GetHeader("X-PROTECTED-RESOURCE")
		//Project id
		project := c.GetHeader("X-REQUEST-PROJECT")

		if sessionKey == "" {
			c.JSON(401, gin.H{"message": "You are not authorized. Please login first!"})
			return
		}

		//Get session value from redis
		val, _ := settings.RedisCli().Get(c, sessionKey).Result()

		//If session does not exist
		if val == "" {
			c.JSON(401, gin.H{"message": "You are not authorized. Please login first!"})
			return
		}

		//Convert json to struct
		session := Session{}
		_ = json.Unmarshal([]byte(val), &session)

		//Check if the current user has access rights to the target resource
		policyKey := fmt.Sprintf("user:%s:tenant:%s:policy", session.Uid, session.Tid)
		scanRes, _, _ := settings.RedisCli().HScan(c, policyKey, 0, prs+"*", 1).Result()

		if len(scanRes) < 1 {
			c.JSON(403, gin.H{"message": "You are not allowed to access this resource !"})
			return
		}

		if !util.Contains(session.AuthorizedPids, project) {
			c.JSON(403, gin.H{"message": "You are not allowed to access this resource of target project!"})
			return
		}

		c.JSON(200, gin.H{"message": "Authorization check through !"})
	})

	err := r.Run("127.0.0.1:17010")
	if err != nil {
		panic(err)
	}
}
