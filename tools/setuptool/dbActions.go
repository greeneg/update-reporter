package main

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"time"
)

type OrgUnit struct {
	Id           int    `json:"Id"`
	OUName       string `json:"ouName"`
	Description  string `json:"description"`
	CreatorId    int    `json:"creatorId"`
	CreationDate string `json:"creationDate"`
}

type Role struct {
	Id           int    `json:"Id"`
	RoleName     string `json:"roleName"`
	Description  string `json:"description"`
	CreationDate string `json:"creationDate"`
}

type User struct {
	Id           int
	UserName     string
	FullName     string
	Status       string
	OrgUnitId    int
	RoleId       int
	CreationDate string
}

func convertSqliteTimestamp(t string) string {
	sqlTimestampFormat := "2006-01-02T15:04:05Z"
	timeFormat := "2006-01-02 15:04:05"
	createTime, _ := time.Parse(sqlTimestampFormat, t)
	return createTime.Format(timeFormat)
}

func createOU(ouName string, ouDescription string, creatorId int) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		errPrintln("Could not start DB transaction: " + string(err.Error()))
		return false, err
	}

	q, err := t.Prepare("INSERT INTO OrganizationalUnits (OUName, Description, CreatorId) VALUES (?, ?, ?)")
	if err != nil {
		errPrintln("Could not prepare the DB query: " + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(ouName, ouDescription, creatorId)
	if err != nil {
		errPrintln("Cannot create organizational unit '" + ouName + "': " + string(err.Error()))
		return false, err
	}

	t.Commit()
	return false, nil
}

func getOrgUnitStatus(ouName string) (bool, error) {
	infoPrintln("OU name: " + ouName)
	t, err := DB.Begin()
	if err != nil {
		errPrintln("Could not start DB transaction: " + string(err.Error()))
		return false, err
	}

	q, err := DB.Prepare("SELECT * FROM OrganizationalUnits WHERE OUName IS ?")
	if err != nil {
		errPrintln("Could not prepare DB query: " + string(err.Error()))
		return false, err
	}

	ou := OrgUnit{}
	err = q.QueryRow(ouName).Scan(
		&ou.Id,
		&ou.OUName,
		&ou.Description,
		&ou.CreatorId,
		&ou.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errPrintln("Encountered error when query database: " + string(err.Error()))
			return false, err
		}
		return false, nil
	}

	t.Commit()

	return true, nil
}

func getOrgUnitByName(ouName string) (OrgUnit, error) {
	rec, err := DB.Prepare("SELECT * FROM OrganizationalUnits WHERE OUName = ?")
	if err != nil {
		errPrintln("Could not prepare the DB query: " + string(err.Error()))
		return OrgUnit{}, err
	}

	ou := OrgUnit{}
	err = rec.QueryRow(ouName).Scan(
		&ou.Id,
		&ou.OUName,
		&ou.Description,
		&ou.CreatorId,
		&ou.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errPrintln("No such org unit found in DB: " + string(err.Error()))
			return OrgUnit{}, nil
		}
		errPrintln("Cannot retrieve org unit from DB: " + string(err.Error()))
		return OrgUnit{}, err
	}

	ou.CreationDate = convertSqliteTimestamp(ou.CreationDate)

	return ou, nil

}

func getRoleStatus(role string) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		errPrintln("Could not start DB transaction: " + string(err.Error()))
		return false, err
	}

	q, err := DB.Prepare("SELECT * FROM Roles WHERE RoleName IS ?")
	if err != nil {
		errPrintln("Could not prepare DB query! " + string(err.Error()))
		return false, err
	}

	rr := Role{}
	err = q.QueryRow(role).Scan(
		&rr.Id,
		&rr.RoleName,
		&rr.Description,
		&rr.CreationDate,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			errPrintln("Encountered error when querying database: " + string(err.Error()))
			return false, err
		}
		return false, nil
	}

	t.Commit()

	return true, nil
}

func createRole(roleName string, roleDescription string) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		errPrintln("Could not start DB transaction: " + string(err.Error()))
		return false, err
	}

	q, err := t.Prepare("INSERT INTO Roles (RoleName, Description) VALUES (?, ?)")
	if err != nil {
		errPrintln("Could not prepare the DB query: " + string(err.Error()))
		return false, err
	}

	_, err = q.Exec(roleName, roleDescription)
	if err != nil {
		errPrintln("Cannot create role '" + roleName + "': " + string(err.Error()))
		return false, err
	}

	t.Commit()

	return true, nil
}

func getRoleByName(roleName string) (Role, error) {
	rec, err := DB.Prepare("SELECT * FROM Roles WHERE RoleName = ?")
	if err != nil {
		errPrintln("Could not prepare the DB query!" + string(err.Error()))
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
			errPrintln("No such role found in DB: " + string(err.Error()))
			return Role{}, nil
		}
		errPrintln("Cannot retrieve role from DB: " + string(err.Error()))
		return Role{}, err
	}

	role.CreationDate = convertSqliteTimestamp(role.CreationDate)
	rec.Close()

	return role, nil
}

func getAccountStatus(account string) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		errPrintln("Could not start DB transaction: " + string(err.Error()))
		return false, err
	}

	q, err := DB.Prepare("SELECT * FROM Users WHERE UserName IS ?")
	if err != nil {
		errPrintln("Could not prepare DB query! " + string(err.Error()))
		return false, err
	}

	err = q.QueryRow(account).Scan()
	if err != nil {
		if err != sql.ErrNoRows {
			errPrintln("Encountered error when querying database: " + string(err.Error()))
			return false, err
		}
		return false, err
	}

	t.Commit()

	return true, nil
}

func createAccount(accountName string, accountFullName string, orgUnitId int, roleId int, passwd string) (User, error) {
	t, err := DB.Begin()
	if err != nil {
		errPrintln("Could not start DB transaction!" + string(err.Error()))
		return User{}, err
	}

	q, err := t.Prepare("INSERT INTO Users (UserName, FullName, OrgUnitId, RoleId, PasswordHash) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		errPrintln("Could not prepare the DB query!" + string(err.Error()))
		return User{}, err
	}

	// take password and hash it
	hash := sha512.Sum512([]byte(passwd))
	passwdHash := hex.EncodeToString(hash[:])

	// get the org Id

	_, err = q.Exec(accountName, accountFullName, orgUnitId, roleId, passwdHash)
	if err != nil {
		errPrintln("Cannot create user '" + accountName + "': " + string(err.Error()))
		return User{}, err
	}

	t.Commit()

	user, err := getAccountByName(accountName)
	if err != nil {
		errPrintln("Could not retrieve user account: " + string(err.Error()))
		return User{}, err
	}

	return user, nil
}

func getAccountByName(accountName string) (User, error) {
	rec, err := DB.Prepare("SELECT Id,UserName,FullName,Status,OrgUnitId,RoleId,CreationDate FROM Users WHERE UserName = ?")
	if err != nil {
		errPrintln("Could not prepare the DB query: " + string(err.Error()))
		return User{}, err
	}

	user := User{}
	err = rec.QueryRow(accountName).Scan(
		&user.Id,
		&user.UserName,
		&user.FullName,
		&user.Status,
		&user.OrgUnitId,
		&user.RoleId,
		&user.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errPrintln("No such role found in DB: " + string(err.Error()))
			return User{}, nil
		}
		errPrintln("Cannot retrieve role from DB: " + string(err.Error()))
		return User{}, err
	}

	user.CreationDate = convertSqliteTimestamp(user.CreationDate)

	return user, nil
}
