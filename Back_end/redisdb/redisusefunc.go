package redisdb

import (
	"context"
	"fmt"
	"log"
	"strings"
)


func RegisterUser ( usr , pwd string ) error {

	// redis mai set krre user by username :key pair 
    // Basic set get shit redis
	err := RedisClient.Set(context.Background()  , usr , pwd , 0).Err() ; if err != nil {
		log.Fatal(" User could not be registered :" , err)

		return err
	}

	err = RedisClient.SAdd(context.Background() , userSetKey() , usr).Err() ; if err != nil {
		 
		// err in adding client : safety net delete the username 

		  RedisClient.Del(context.Background() , usr )


		 return err
	}

	return nil
}

func ExistUser  ( usr string ) bool {

	// redis check by ismember
	ans := RedisClient.SIsMember(context.Background() , userSetKey() , usr).Val()

	return ans
}

func IsUserAuthentic  (username, password string) error {
	

	// we get username from the 
	p := RedisClient.Get(context.Background(), username).Val()

	if p == "" {
		return fmt.Errorf("user does not exist")
	} else if strings.EqualFold(p, password) {
		return fmt.Errorf(" password entered is wrong. Please enter the password again")
	}

	return nil
}

