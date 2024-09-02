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
	"github.com/greeneg/update-reporterd/model"
)

// CreateUser Register a role for user rights assignment
//
//	@Summary		Register role
//	@Description	Add a new role
//	@Tags			role
//	@Accept			json
//	@Produce		json
//	@Param			role	body	model.Role	true	"Role data"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/role [post]
func (u *UpdateReporter) CreateRole(c *gin.Context) {
	_, authed := u.GetUserId(c)
	if authed {
		var json model.Role
		if err := c.ShouldBindJSON(&json); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s, err := model.CreateRole(json)
		if s {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Role '" + json.RoleName + "' has been added to system"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// DeleteRole Remove a role
//
//	@Summary		Delete role
//	@Description	Delete a role
//	@Tags			role
//	@Accept			json
//	@Produce		json
//	@Param			roleId	path	int	true	"Role Id"
//	@Security		BasicAuth
//	@Success		200	{object}	model.SuccessMsg
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/role/{roleId} [delete]
func (u *UpdateReporter) DeleteRole(c *gin.Context) {
	_, authed := u.GetUserId(c)
	if authed {
		roleId, _ := strconv.Atoi(c.Param("roleId"))
		status, err := model.DeleteRole(roleId)
		if err != nil {
			log.Println("ERROR: Cannot delete role: " + string(err.Error()))
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove role! " + string(err.Error())})
			return
		}

		if status {
			roleIdStr := strconv.Itoa(roleId)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Role Id " + roleIdStr + " has been removed from system"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove role!"})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetRoles Retrieve list of all roles
//
//	@Summary		Retrieve list of all roles
//	@Description	Retrieve list of all roles
//	@Tags			role
//	@Produce		json
//	@Security		BasicAuth
//	@Success		200	{object}	model.RolesList
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/roles [get]
func (u *UpdateReporter) GetRoles(c *gin.Context) {
	_, authed := u.GetUserId(c)
	if authed {
		roles, err := model.GetRoles()
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if roles == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No records found!"})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"data": roles})
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetRoleById Retrieve a role by its Id
//
//	@Summary		Retrieve a role by its Id
//	@Description	Retrieve a role by its Id
//	@Tags			role
//	@Produce		json
//	@Param			roleId	path int true "Role ID"
//	@Security		BasicAuth
//	@Success		200	{object}	model.Role
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/role/id/{roleId} [get]
func (u *UpdateReporter) GetRoleById(c *gin.Context) {
	_, authed := u.GetUserId(c)
	if authed {
		id, _ := strconv.Atoi(c.Param("roleId"))
		role, err := model.GetRoleById(id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if role.RoleName == "" {
			strId := strconv.Itoa(id)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No records found with role id " + strId})
		} else {
			c.IndentedJSON(http.StatusOK, role)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}

// GetRoleByName Retrieve a role by its role name
//
//	@Summary		Retrieve a role by its role name
//	@Description	Retrieve a role by its role name
//	@Tags			role
//	@Produce		json
//	@Param			roleName	path string true "Role Name"
//	@Security		BasicAuth
//	@Success		200	{object}	model.Role
//	@Failure		400	{object}	model.FailureMsg
//	@Router			/role/name/{roleName} [get]
func (u *UpdateReporter) GetRoleByName(c *gin.Context) {
	_, authed := u.GetUserId(c)
	if authed {
		roleName := c.Param("roleName")
		role, err := model.GetRoleByName(roleName)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": string(err.Error())})
			return
		}

		if role.RoleName == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No records found with role name " + roleName})
		} else {
			c.IndentedJSON(http.StatusOK, role)
		}
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Insufficient access. Access denied!"})
	}
}
