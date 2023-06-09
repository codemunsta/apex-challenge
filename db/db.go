package db

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func CreateUser(user User) error {
	userJson, err := json.Marshal(user)
	if err != nil {
		panic(err)
	} else {
		err = RedisClient.Set(context.Background(), "user:"+user.Name, string(userJson), 500000000000).Err()
		if err != nil {
			panic(err)
		} else {
			return err
		}
	}
}

func GetUser(username string) (User, bool) {
	userJson, err := RedisClient.Get(context.Background(), "user:"+username).Result()
	if err == redis.Nil {
		userFound := false
		var user User
		return user, userFound
	} else if err != nil {
		panic(err)
	} else {
		userFound := true
		var user User
		fetchErr := json.Unmarshal([]byte(userJson), &user)
		if fetchErr != nil {
			panic(err)
		} else {
			return user, userFound
		}
	}
}

func CreateSession(session Session) error {
	sessionJson, err := json.Marshal(session)
	if err != nil {
		panic(err)
	}
	key := "session" + session.IUser.Name
	err = RedisClient.LPush(context.Background(), key, sessionJson).Err()
	if err != nil {
		panic(err)
	}
	return err
}

func GetSession(username string) (Session, int, bool) {
	key := "session" + username
	sessionJsons, err := RedisClient.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		panic(err)
	}

	for index, sessionJson := range sessionJsons {
		var session Session
		err := json.Unmarshal([]byte(sessionJson), &session)
		if err != nil {
			panic(err)
		}
		if session.IsActive {
			return session, index, true
		}
	}
	var session Session
	return session, 0, false
}

func UpdateSession(username string, index int, session Session) error {
	key := "session" + username
	modifiedSessionJSON, err := json.Marshal(session)
	if err != nil {
		panic(err)
	}
	_, err = RedisClient.LSet(context.Background(), key, int64(index), string(modifiedSessionJSON)).Result()
	if err != nil {
		panic(err)
	}
	return nil
}

func CreateTransaction(transaction Transaction) error {
	transactionJson, err := json.Marshal(transaction)
	if err != nil {
		panic(err)
	}
	key := "transaction" + transaction.IUser
	err = RedisClient.LPush(context.Background(), key, transactionJson).Err()
	if err != nil {
		panic(err)
	}
	return err
}

func GetTransactions(username string) ([]Transaction, error) {
	key := "transaction" + username
	transactionJsons, err := RedisClient.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	var transactions []Transaction
	for _, transactionJson := range transactionJsons {
		var transaction Transaction
		err := json.Unmarshal([]byte(transactionJson), &transaction)
		if err != nil {
			panic(err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
