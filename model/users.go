package model

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
	//	"crypto/sha512"
	"database/sql"
	//	"encoding/hex"
	//	"errors"
	"log"
	// "strconv"
	// "time"
)

func GetUserByUserName(username string) (User, error) {
	log.Println("INFO: User by username requested: " + username)
	rec, err := DB.Prepare("SELECT * FROM Users WHERE UserName = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return User{}, err
	}

	user := User{}
	err = rec.QueryRow(username).Scan(
		&user.Id,
		&user.UserName,
		&user.FullName,
		&user.Status,
		&user.OrgUnitId,
		&user.RoleId,
		&user.PasswordHash,
		&user.CreationDate,
		&user.LastPasswordChangedDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such user found in DB: " + string(err.Error()))
			return User{}, nil
		}
		log.Println("ERROR: Cannot retrieve user from DB: " + string(err.Error()))
		return User{}, err
	}

	user.CreationDate = ConvertSqliteTimestamp(user.CreationDate)
	user.LastPasswordChangedDate = ConvertSqliteTimestamp(user.LastPasswordChangedDate)

	return user, nil
}
