package handler

import "fmt"

func GetUserTypeByUserName(username string) (string, error) {
	for _, user := range UserData {
		if user.Username == username {
			return user.UserType, nil
		}
	}

	return "", fmt.Errorf("username not found")
}
