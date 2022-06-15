package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const (
	operationFlag = "operation"
	itemFlag      = "item"
	fileNameFlag  = "fileName"
	idFlag        = "id"
)

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Perform(args Arguments, writer io.Writer) error {
	operation, ok := args[operationFlag]
	if operation == "" || !ok {
		return fmt.Errorf("-%s flag has to be specified", operationFlag)
	}
	fileName, ok := args[fileNameFlag]
	if fileName == "" || !ok {
		return fmt.Errorf("-%s flag has to be specified", fileNameFlag)
	}
	switch operation {
	case "add":

		item, ok := args[itemFlag]
		if item == "" || !ok {
			return fmt.Errorf("-item flag has to be specified")
		}

		return AddUser(args, writer)

	case "list":
		usersList, err := UsersList(args)
		if err != nil {
			return err
		}
		jsonUsers, err := json.Marshal(usersList)
		if err != nil {
			return err
		}
		fmt.Fprintf(writer, string(jsonUsers))
		return nil
	case "findById":
		id, ok := args[idFlag]
		if id == "" || !ok {
			return fmt.Errorf("-id flag has to be specified")
		}
		return FindByID(args, writer)
	case "remove":
		id, ok := args[idFlag]
		if id == "" || !ok {
			return fmt.Errorf("-id flag has to be specified")
		}
		return RemoveById(args, writer)
	}
	return fmt.Errorf("Operation %s not allowed!", operation)
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

func parseArgs() Arguments {
	result := make(map[string]string)
	operation := flag.String(operationFlag, "", "")
	item := flag.String(itemFlag, "", "")
	fileName := flag.String(fileNameFlag, "", "")
	id := flag.String(idFlag, "", "")
	flag.Parse()

	result[operationFlag] = *operation
	result[itemFlag] = *item
	result[fileNameFlag] = *fileName
	result[idFlag] = *id

	return result
}
func AddUser(args Arguments, writer io.Writer) error {
	user := User{}
	err := json.Unmarshal([]byte(args[itemFlag]), &user)
	if err != nil {
		return err
	}

	usersList, err := UsersList(args)
	if err != nil {
		return err
	}

	for _, userL := range usersList {
		if userL.Id == user.Id {
			fmt.Fprint(writer, fmt.Sprintf("Item with id %s already exists", user.Id))
			return nil
		}
	}

	usersList = append(usersList, user)
	jsonUsers, err := json.Marshal(usersList)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(args[fileNameFlag], jsonUsers, os.ModePerm)
}

func RemoveById(args Arguments, writer io.Writer) error {
	var userExist bool
	users, err := UsersList(args)
	if err != nil {
		return err
	}
	var newUserList []User
	for _, user := range users {
		if user.Id == args[idFlag] {
			userExist = true
		} else {
			newUserList = append(newUserList, user)
		}
	}

	if !userExist {
		fmt.Fprint(writer, fmt.Sprintf("Item with id %s not found", args[idFlag]))
	} else {

		jsonUsersList, err := json.Marshal(newUserList)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(args[fileNameFlag], jsonUsersList, 0777)
		if err != nil {
			return err
		}
	}

	return nil
}
func FindByID(args Arguments, writer io.Writer) error {
	users, err := UsersList(args)
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Id == args[idFlag] {
			jsonUser, err := json.Marshal(user)
			if err != nil {
				return err
			}
			fmt.Fprintf(writer, string(jsonUser))
			break
		}
	}
	return nil
}

func UsersList(args Arguments) ([]User, error) {
	Users := []User{}

	//File, err := os.Open(args[fileNameFlag])
	//if err != nil && !errors.Is(err, os.ErrNotExist) {
	//	return nil, err
	//}

	data, err := ioutil.ReadFile(args[fileNameFlag])
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	if len(data) > 0 {
		err = json.Unmarshal(data, &Users)
		if err != nil {
			return nil, err
		}
	}

	return Users, nil
}
