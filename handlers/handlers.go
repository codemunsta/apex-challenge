package handlers

import (
	"apex-challenge/db"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

type SignInForm struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Amount struct {
	IAMount string `json:"amount"`
}

func SignIn(writer http.ResponseWriter, request *http.Request) {
	// process post body
	decoder := json.NewDecoder(request.Body)
	var signinForm SignInForm
	err := decoder.Decode(&signinForm)
	if err != nil {
		http.Error(writer, "Could not process form body", http.StatusBadRequest)
		return
	}

	// check if user exist
	user, found := db.GetUser(signinForm.Name)
	if found {

		// soft authenticate user
		if user.Password == signinForm.Password {
			response := map[string]interface{}{
				"message": "Login Successfull",
				"token":   user.Name,
			}
			writer.WriteHeader(http.StatusOK)
			json.NewEncoder(writer).Encode(response)
			return
		} else {
			response := map[string]interface{}{
				"message": "Invalid, Create new account",
			}
			writer.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(writer).Encode(response)
			return
		}
	} else {

		// create new user
		user := db.User{
			Name:     signinForm.Name,
			Password: signinForm.Password,
			Account:  0,
		}
		err = db.CreateUser(user)
		if err != nil {
			panic(err)
		}
		response := map[string]interface{}{
			"message": "Login Successfull",
			"token":   user.Name,
		}
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(response)
		return
	}
}

func FundWallet(writer http.ResponseWriter, request *http.Request) {
	// fetch post body
	decoder := json.NewDecoder(request.Body)
	var amount Amount
	err := decoder.Decode(&amount)
	if err != nil {
		http.Error(writer, "Invalid data", http.StatusBadRequest)
		return
	}

	// convert amount to integer
	intAmount, err2 := strconv.ParseInt(amount.IAMount, 10, 64)
	if err2 != nil {
		http.Error(writer, "Failed to convert string to int", http.StatusBadRequest)
		return
	}

	// verify credit request meet criteria
	authorizationHeader := request.Header.Get("Authorization")
	user, _ := db.GetUser(authorizationHeader)
	if user.Account > 35 {
		response := map[string]interface{}{
			"message": "Account must be less than 35 sats",
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(response)
		return
	} else if intAmount != 155 {
		response := map[string]interface{}{
			"message": "Can only deposit 155 sats",
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(response)
		return
	} else {
		user.Account = user.Account + intAmount

		// save account balance
		// use create user to overide old user information
		err = db.CreateUser(user)
		if err != nil {
			panic(err)
		}
		transaction := db.Transaction{
			IUser:  user.Name,
			Type:   "Credit",
			Amount: intAmount,
		}
		err := db.CreateTransaction(transaction)
		if err != nil {
			panic(err)
		}
		response := map[string]interface{}{
			"message": "Funds Deposited",
			"balance": user.Account,
		}
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(response)
		return
	}
}

func GetWalletBalance(writer http.ResponseWriter, request *http.Request) {
	authorizationHeader := request.Header.Get("Authorization")
	user, _ := db.GetUser(authorizationHeader)
	response := map[string]interface{}{
		"User":    user.Name,
		"balance": user.Account,
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func StartGame(writer http.ResponseWriter, request *http.Request) {
	// get user
	authorizationHeader := request.Header.Get("Authorization")
	user, _ := db.GetUser(authorizationHeader)

	// check for active session
	_, _, sFound := db.GetSession(user.Name)
	if sFound {

		// if active game session exist
		response := map[string]interface{}{
			"message": "Stop game and restart",
		}
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(response)
		return
	} else if user.Account < 25 {

		// if user balance is less than 25
		response := map[string]interface{}{
			"message": "insufficient balance",
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(response)
		return
	} else {

		// deduct 20 sats to begin game
		user.Account = user.Account - 20
		err := db.CreateUser(user)
		if err != nil {
			panic(err)
		}
		transaction := db.Transaction{
			IUser:  user.Name,
			Type:   "Debit",
			Amount: 20,
		}
		err = db.CreateTransaction(transaction)
		if err != nil {
			panic(err)
		}

		// generate number
		min := 2
		max := 12
		randomNumber := rand.Intn(max-min+1) + min

		// create new session
		session := db.Session{
			IUser:    user,
			IsActive: true,
			Number:   int8(randomNumber),
		}
		err = db.CreateSession(session)
		if err != nil {
			panic(err)
		}
		response := map[string]interface{}{
			"message": "Game Started",
		}
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(response)
		return
	}
}

func RollDice(writer http.ResponseWriter, request *http.Request) {
	// get user
	authorizationHeader := request.Header.Get("Authorization")
	user, _ := db.GetUser(authorizationHeader)
	if user.Account < 25 {
		response := map[string]interface{}{
			"message": "insufficient balance",
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(response)
		return
	}

	// fetch active session
	session, index, _ := db.GetSession(user.Name)

	// generate random number
	min := 1
	max := 6
	randomNumber := rand.Intn(max-min+1) + min

	games := session.Games
	gamePairNo := len(games)
	if gamePairNo == 0 || games[gamePairNo-1].Roll2 != 0 {

		// create new game pair
		game := db.Game{
			Roll1: int8(randomNumber),
			Roll2: 0,
		}
		games = append(games, game)

		// deduct 5 sat and update user and session
		user.Account = user.Account - 5
		session.Games = games
		err := db.CreateUser(user)
		if err != nil {
			panic(err)
		}
		err = db.UpdateSession(user.Name, index, session)
		if err != nil {
			panic(err)
		}
		transaction := db.Transaction{
			IUser:  user.Name,
			Type:   "Debit",
			Amount: 5,
		}
		err = db.CreateTransaction(transaction)
		if err != nil {
			panic(err)
		}
		response := map[string]interface{}{
			"message": fmt.Sprintf("you rolled %v, roll again", randomNumber),
		}
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(response)
		return
	} else {

		// second dice roll
		games[gamePairNo-1].Roll2 = int8(randomNumber)
		total := games[gamePairNo-1].Roll1 + games[gamePairNo-1].Roll2
		if total == session.Number {

			// credit user 10 sat for winning
			user.Account = user.Account + 10
			session.Games = games
			err := db.CreateUser(user)
			if err != nil {
				panic(err)
			}
			err = db.UpdateSession(user.Name, index, session)
			if err != nil {
				panic(err)
			}
			transaction := db.Transaction{
				IUser:  user.Name,
				Type:   "Credit",
				Amount: 10,
			}
			err = db.CreateTransaction(transaction)
			if err != nil {
				panic(err)
			}
			response := map[string]interface{}{
				"message": "you win!",
			}
			writer.WriteHeader(http.StatusOK)
			json.NewEncoder(writer).Encode(response)
			return
		} else {
			session.Games = games
			err := db.CreateUser(user)
			if err != nil {
				panic(err)
			}
			err = db.UpdateSession(user.Name, index, session)
			if err != nil {
				panic(err)
			}
			response := map[string]interface{}{
				"message": "you lose!, try again",
			}
			writer.WriteHeader(http.StatusOK)
			json.NewEncoder(writer).Encode(response)
			return
		}
	}
}

func EndGame(writer http.ResponseWriter, request *http.Request) {
	authorizationHeader := request.Header.Get("Authorization")
	user, _ := db.GetUser(authorizationHeader)
	session, index, _ := db.GetSession(user.Name)
	session.IsActive = false

	// update session
	err := db.UpdateSession(user.Name, index, session)
	if err != nil {
		panic(err)
	}
	response := map[string]interface{}{
		"message": "Game Over, Thank you",
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

// bonus
func ActiveSession(writer http.ResponseWriter, request *http.Request) {
	authorizationHeader := request.Header.Get("Authorization")
	user, _ := db.GetUser(authorizationHeader)
	_, _, active := db.GetSession(user.Name)
	if active {
		response := map[string]interface{}{
			"message": "You have a current active game session",
		}
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(response)
	} else {
		response := map[string]interface{}{
			"message": "You do not have a current active session",
		}
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(response)
	}
}

func GetTransactionList(writer http.ResponseWriter, request *http.Request) {
	authorizationHeader := request.Header.Get("Authorization")
	transactions, _ := db.GetTransactions(authorizationHeader)

	response := map[string]interface{}{
		"message":      "Your Transactions",
		"transactions": transactions,
	}
	writer.WriteHeader(http.StatusNotFound)
	json.NewEncoder(writer).Encode(response)
}
