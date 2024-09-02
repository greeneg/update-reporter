package controllers

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
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/greeneg/update-reporterd/helpers"
	"github.com/greeneg/update-reporterd/model"
)

// CreateUser Register a user for authentication and authorization
//
//	@Summary		Register user
//	@Description	Add a new user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body	model.ProposedUser	true	"User Data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user [post]
func (u *UpdateReporter) CreateUser(c *gin.Context) {
	_, authed := u.GetUserId(c)
	if authed {
		var json model.ProposedUser
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s, err := model.CreateUser(json)
		if s {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "User has been added to system"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// ChangeAccountPassowrd Change an account's password
//
//	@Summary		Change password
//	@Description	Change password
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			name	path	string	true	"User name"
//	@Param			changePassword	body	model.PasswordChange	true	"Password data"
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user/{name} [patch]
func (u *UpdateReporter) ChangeAccountPassword(c *gin.Context) {
	username := c.Param("name")
	var json model.PasswordChange
	if err := c.ShouldBindJSON(&json); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := model.ChangeAccountPassword(username, json.OldPassword, json.NewPassword)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": string(err.Error())})
		return
	}

	if status {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "User '" + username + "' has changed their password"})
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User password could not be updated!"})
	}
}

// DeleteUser Remove a user for authentication and authorization
//
//	@Summary		Delete user
//	@Description	Delete a user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			name	path	string	true	"User name"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user/{name} [delete]
func (u *UpdateReporter) DeleteUser(c *gin.Context) {
	_, authed := u.GetUserId(c)
	if authed {
		username := c.Param("name")
		status, err := model.DeleteUser(username)
		if err != nil {
			log.Println("ERROR: Cannot delete user: " + string(err.Error()))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove user! " + string(err.Error())})
			return
		}

		if status {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "User " + username + " has been removed from system"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove user!"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetUserStatus Retrieve the active status of a user. Can be either 'enabled' or 'locked'
//
//	@Summary		Retrieve a user's active status. Can be either 'enabled' or 'locked'
//	@Description	Retrieve a user's active status
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			name	path	string	true	"User name"
//	@Security		BasicAuth
//	@Success		200	{object}	model.UserStatusMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user/{name}/status [get]
func (u *UpdateReporter) GetUserStatus(c *gin.Context) {
	_, authed := u.GetUserId(c)
	if authed {
		username := c.Param("name")
		status, err := model.GetUserStatus(username)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to get the user " + username + " status: " + string(err.Error())})
			return
		}

		if status != "" {
			c.IndentedJSON(http.StatusOK, gin.H{
				"message":    "User '" + username + "' has status " + status,
				"userStatus": status,
			})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unable to retrieve user status"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// SetUserStatus Set the active status of a user. Can be either 'enabled' or 'locked'
//
//	@Summary		Set a user's active status. Can be either 'enabled' or 'locked'
//	@Description	Set a user's active status
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body	model.UserStatus	true	"User Data"
//	@Param			name	path	string	true "User name"
//	@Security		BasicAuth
//	@Success		200	{object}	model.UserStatusMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user/{name}/status [patch]
func (u *UpdateReporter) SetUserStatus(c *gin.Context) {
	_, authed := u.GetUserId(c)
	if authed {
		username := c.Param("name")
		var json model.UserStatus
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		status, err := model.SetUserStatus(username, json)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": string(err.Error())})
			return
		}

		if status {
			c.IndentedJSON(http.StatusOK, gin.H{
				"message":    "User '" + username + "' has been " + json.Status,
				"userStatus": json.Status,
			})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// SetUserRoleId Set the role Id of a user
//
//	@Summary		Set a user's role Id
//	@Description	Set a user's role Id
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			roleId	body	model.UserRoleId	true	"Role Id"
//	@Param			name	path	string	true	"User name"
//	@Security		BasicAuth
//	@Success		200 {object}	model.UserRoleIdMsg
//	@Failure		400 {object}	model.FailureMsg
//	@Router			/user/{name}/roleId [patch]
func (u *UpdateReporter) SetUserRoleId(c *gin.Context) {
	_, authed := u.GetUserId(c)
	if authed {
		username := c.Param("name")
		var json model.UserRoleId
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		status, err := model.SetUserRoleId(username, json)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": string(err.Error())})
		}

		if status {
			roleId := strconv.Itoa(json.RoleId)
			c.Copy().IndentedJSON(http.StatusOK, gin.H{
				"message": "User '" + username + "' has been set to role Id '" + roleId + "'",
				"roleId":  json.RoleId,
			})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetUsers Retrieve list of all users
//
//	@Summary		Retrieve list of all users
//	@Description	Retrieve list of all users
//	@Tags			user
//	@Produce		json
//	@Security		BasicAuth
//	@Success		200	{object}	model.UsersList
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/users [get]
func (u *UpdateReporter) GetUsers(c *gin.Context) {
	_, authed := u.GetUserId(c)
	if authed {
		users, err := model.GetUsers()
		helpers.FatalCheckError(err)

		safeUsers := make([]SafeUser, 0)
		for _, user := range users {
			safeUser := SafeUser{}
			safeUser.Id = user.Id
			safeUser.UserName = user.UserName
			safeUser.FullName = user.FullName
			safeUser.OrgUnitId = user.OrgUnitId
			safeUser.RoleId = user.RoleId
			safeUser.CreationDate = user.CreationDate

			safeUsers = append(safeUsers, safeUser)
		}

		if users == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No records found!"})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"data": safeUsers})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetUsersByRoleId Retrieve list of users by role Id
//
//	@Summary        Retrieve list of users by role Id
//	@Description    Retrieve list of users by role Id
//	@Tags           user
//	@Produce        json
//	@Param          roleId	path int true "Role Id"
//	@Security		BasicAuth
//	@Success        200 {object}	model.UsersList
//	@Failure		400 {object}	model.FailureMsg
//	@Router			/users/byRoleId/{roleId} [get]
func (u *UpdateReporter) GetUsersByRoleId(c *gin.Context) {
	_, authed := u.GetUserId(c)
	if authed {
		roleId, _ := strconv.Atoi(c.Param("roleId"))
		users, err := model.GetUsersByRoleId(roleId)
		helpers.FatalCheckError(err)

		safeUsers := make([]SafeUser, 0)
		for _, user := range users {
			safeUser := SafeUser{}
			safeUser.Id = user.Id
			safeUser.UserName = user.UserName
			safeUser.FullName = user.FullName
			safeUser.OrgUnitId = user.OrgUnitId
			safeUser.RoleId = user.RoleId
			safeUser.CreationDate = user.CreationDate

			safeUsers = append(safeUsers, safeUser)
		}

		if users == nil {
			strId := strconv.Itoa(roleId)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No records found for users with role Id " + strId})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"data": safeUsers})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetUserById Retrieve a user by their Id
//
//	@Summary		Retrieve a user by their Id
//	@Description	Retrieve a user by their Id
//	@Tags			user
//	@Produce		json
//	@Param			id	path int true "User ID"
//	@Success		200	{object}	SafeUser
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user/id/{id} [get]
func (u *UpdateReporter) GetUserById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ent, err := model.GetUserById(id)
	helpers.FatalCheckError(err)

	// don't return the password hash
	safeUser := new(SafeUser)
	safeUser.Id = ent.Id
	safeUser.UserName = ent.UserName
	safeUser.FullName = ent.FullName
	safeUser.OrgUnitId = ent.OrgUnitId
	safeUser.RoleId = ent.RoleId
	safeUser.CreationDate = ent.CreationDate

	if ent.UserName == "" {
		strId := strconv.Itoa(id)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No records found with user id " + strId})
	} else {
		c.IndentedJSON(http.StatusOK, safeUser)
	}
}

// GetUserByName Retrieve a user by their UserName
//
//	@Summary		Retrieve a user by their UserName
//	@Description	Retrieve a user by their UserName
//	@Tags			user
//	@Produce		json
//	@Param			name	path	string	true	"User name"
//	@Success		200	{object}	SafeUser
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/user/name/{name} [get]
func (u *UpdateReporter) GetUserByUserName(c *gin.Context) {
	username := c.Param("name")
	ent, err := model.GetUserByUserName(username)
	helpers.FatalCheckError(err)

	// don't return the password hash
	safeUser := new(SafeUser)
	safeUser.Id = ent.Id
	safeUser.UserName = ent.UserName
	safeUser.FullName = ent.FullName
	safeUser.OrgUnitId = ent.OrgUnitId
	safeUser.RoleId = ent.RoleId
	safeUser.CreationDate = ent.CreationDate

	if ent.UserName == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No records found with user name " + username})
	} else {
		c.IndentedJSON(http.StatusOK, safeUser)
	}
}
