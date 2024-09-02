package helpers

/*

  Copyright 2024, YggdrasilSoft, LLC.

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

import (
	"crypto/sha512"
	"encoding/hex"
	"log"
	"strings"

	"github.com/greeneg/update-reporterd/model"
)

func CheckIsNotLocked(u model.User) bool {
	return u.Status != "locked"
}

func CheckUserPass(username, password string) bool {
	user, err := model.GetUserByUserName(username)
	if err != nil {
		return false
	}

	status := CheckIsNotLocked(user)
	if !status {
		return false
	}

	// get the password hash from the user so we can compare it
	pwHash := user.PasswordHash

	// now calculate the sha512 of the password and see if it matches
	sha := sha512.Sum512([]byte(password))
	newPwHash := hex.EncodeToString(sha[:])

	return pwHash == newPwHash // returns boolean based on equality
}

func EmptyUserPass(username, password string) bool {
	return strings.Trim(username, " ") == "" || strings.Trim(password, " ") == ""
}

func FatalCheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
