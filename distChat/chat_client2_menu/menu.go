package menu

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"

	database "github.com/Distributed-Messaging/distChat/database"
)

var collection *mongo.Collection

func SetCollection(newCollection *mongo.Collection) {
	collection = newCollection
}

func ListGroups(username string, signedIn bool) error {
	//ListGroups prints all groups owned by username to the screen
	if !signedIn {
		fmt.Println("Please sign in before lsiting groups")
		return fmt.Errorf("User not signed in")
	}
	groups, err := database.GetOwnedGroups(username, collection)
	if err != nil {
		return fmt.Errorf("Error reading owned groups: %v", err)
	}
	fmt.Println("You own these groups")
	for _, group := range groups {
		fmt.Println(group)
	}
	return nil
}

func CreateAccount() error {
	//allows users to create accounts with new usernames and passwords
	buf := bufio.NewReader(os.Stdin)
	fmt.Print("enter your new username: ")
	username, usernameerr := buf.ReadString('\n')
	if usernameerr != nil {
		return fmt.Errorf("Error reading username: %v", usernameerr)
	}
	fmt.Print("enter your new password: ")
	password, passworderr := buf.ReadString('\n')
	if passworderr != nil {
		return fmt.Errorf("Error reading password: %v", passworderr)
	}
	username = username[:len(username)-1]
	password = password[:len(password)-1]

	database.StoreAccount(
		database.Account{
			Password: password,
			Name:     username,
		}, collection)
	fmt.Printf("created account %v.\n", username)
	return nil
}

func ChangePassword() error {
	//ChangePassword allows a user to change the password of any account
	//the user knows the credentials for
	success, account, _ := SignIn()
	if !success {
		return errors.New("Login failed, cannot change password")
	}
	buf := bufio.NewReader(os.Stdin)
	fmt.Print("enter your new password: ")
	password, passworderr := buf.ReadString('\n')
	if passworderr != nil {
		return fmt.Errorf("Error reading password: %v", passworderr)
	}
	username := account.Name
	password = password[:len(password)-1]
	var updatedAccount database.Account
	updatedAccounterr := collection.FindOneAndReplace(context.Background(),
		database.Account{Name: username},
		database.Account{Name: username, Password: password}).Decode(&updatedAccount)
	if updatedAccounterr != nil {
		return fmt.Errorf("Error updating account: %v", updatedAccounterr)
	}
	fmt.Printf("updated account name to: %v, password to: %v\n", updatedAccount.Name, password)
	return nil
}

func SignIn() (bool, database.Account, error) {
	//SignIn authenticates a user and returns their account information and a true bool
	//if they have the correct password
	errorAccount := database.Account{
		Password: "FALSE_ERROR",
		Name:     "FALSE_ERROR",
	}
	buf := bufio.NewReader(os.Stdin)
	fmt.Print("enter your username: ")
	username, usernameerr := buf.ReadString('\n')
	if usernameerr != nil {
		return false, errorAccount, fmt.Errorf("Error reading username: %v", usernameerr)
	}
	username = username[:len(username)-1]
	account, accounterr := database.GetOneAccountByName(username, collection)
	if accounterr != nil {
		return false, errorAccount, fmt.Errorf("Error username does not match any registerd accounts, %v", accounterr)
	}

	incorrectPassword := true
	incorrectPasswordCount := 0
	for incorrectPassword {
		fmt.Print("enter your password: ")
		password, passworderr := buf.ReadString('\n')
		if passworderr != nil {
			return false, errorAccount, fmt.Errorf("Error reading password: %v", passworderr)
		}
		password = password[:len(password)-1]
		if account.Password != password {
			log.Printf("Error incorrect password: %v\n", passworderr)
			incorrectPasswordCount++
		} else {
			incorrectPassword = false
		}
	}

	fmt.Printf("Logged in as %v.\n", account.Name)
	return true, account, nil
}

func PickGroup(username string) ([]string, string, error) {
	//PickGroup connects a user to a group. Or, if that group
	//does not exist. Asks them for IPs for the new group and connects them
	buf := bufio.NewReader(os.Stdin)
	fmt.Print("enter group name: ")
	group, grouperr := buf.ReadString('\n')
	if grouperr != nil {
		return nil, "", fmt.Errorf("Error group name: %v", grouperr)
	}
	group = group[:len(group)-1]
	grouptojoin, err := database.GetOneGroup(group, collection)
	if err != nil {
		fmt.Print("Please enter ips to connect to seperated by spaces: ")
		ipToConnect, ipcerr := buf.ReadString('\n')
		if ipcerr != nil {
			return nil, "", fmt.Errorf("Error reading ip and port: %v", ipcerr)
		}
		ipToConnect = ipToConnect[:len(ipToConnect)-1]
		IPs := strings.Fields(ipToConnect)
		grouptojoin = database.Group{
			IPs:   IPs,
			Name:  group,
			Owner: username,
		}
		database.StoreGroup(grouptojoin, collection)
	}

	fmt.Printf("connecting to %s", group)
	return grouptojoin.IPs, grouptojoin.Name, nil
}
