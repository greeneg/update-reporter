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
	"database/sql"
	"log"
	"strconv"
)

func CreateRole(r Role) (bool, error) {
	log.Println("INFO: User creation requested: " + r.RoleName)
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Could not start DB transaction!" + string(err.Error()))
		return false, err
	}

	q, err := t.Prepare("INSERT INTO Roles (RoleName, Description) VALUES (?, ?)")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(r.RoleName, r.Description)
	if err != nil {
		log.Println("ERROR: Cannot create user '" + r.RoleName + "': " + string(err.Error()))
		return false, err
	}

	t.Commit()

	log.Println("INFO: User '" + r.RoleName + "' created")
	return true, nil
}

func DeleteRole(roleId int) (bool, error) {
	log.Println("INFO: Role deletion requested: " + strconv.Itoa(roleId))
	t, err := DB.Begin()
	if err != nil {
		log.Println("ERROR: Could not start DB transaction!" + string(err.Error()))
		return false, err
	}

	q, err := DB.Prepare("DELETE FROM Roles WHERE Id IS ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(roleId)
	if err != nil {
		log.Println("ERROR: Cannot delete user '" + strconv.Itoa(roleId) + "': " + string(err.Error()))
		return false, err
	}

	t.Commit()

	log.Println("INFO: Role with Id '" + strconv.Itoa(roleId) + "' has been deleted")
	return true, nil
}

func GetRoles() ([]Role, error) {
	log.Println("INFO: List of role object requested")
	rows, err := DB.Query("SELECT * FROM Roles")
	if err != nil {
		log.Println("ERROR: Could not run the DB query!" + string(err.Error()))
		return nil, err
	}

	roles := make([]Role, 0)
	for rows.Next() {
		role := Role{}
		err = rows.Scan(
			&role.Id,
			&role.RoleName,
			&role.Description,
			&role.CreationDate,
		)
		if err != nil {
			log.Println("ERROR: Cannot marshal the user objects!" + string(err.Error()))
			return nil, err
		}

		role.CreationDate = ConvertSqliteTimestamp(role.CreationDate)

		roles = append(roles, role)
	}

	log.Println("INFO: List of all users retrieved")
	return roles, nil
}

func GetRoleById(id int) (Role, error) {
	log.Println("INFO: Role by Id requested: " + strconv.Itoa(id))
	rec, err := DB.Prepare("SELECT * FROM Roles WHERE Id = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return Role{}, err
	}

	role := Role{}
	err = rec.QueryRow(id).Scan(
		&role.Id,
		&role.RoleName,
		&role.Description,
		&role.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such user found in DB: " + string(err.Error()))
			return Role{}, nil
		}
		log.Println("ERROR: Cannot retrieve user from DB: " + string(err.Error()))
		return Role{}, err
	}

	role.CreationDate = ConvertSqliteTimestamp(role.CreationDate)

	return role, nil
}

func GetRoleByName(roleName string) (Role, error) {
	log.Println("INFO: Role by Id requested: " + roleName)
	rec, err := DB.Prepare("SELECT * FROM Roles WHERE RoleName = ?")
	if err != nil {
		log.Println("ERROR: Could not prepare the DB query!" + string(err.Error()))
		return Role{}, err
	}

	role := Role{}
	err = rec.QueryRow(roleName).Scan(
		&role.Id,
		&role.RoleName,
		&role.Description,
		&role.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("ERROR: No such user found in DB: " + string(err.Error()))
			return Role{}, nil
		}
		log.Println("ERROR: Cannot retrieve user from DB: " + string(err.Error()))
		return Role{}, err
	}

	role.CreationDate = ConvertSqliteTimestamp(role.CreationDate)

	return role, nil
}
