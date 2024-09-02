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

type FailureMsg struct {
	Error string `json:"error"`
}

type HealthCheck struct {
	Db           string `json:"db"`
	DiskSpace    string `json:"diskSpace"`
	DiskWritable string `json:"diskWritable"`
	Health       string `json:"health"`
	Status       int    `json:"status"`
}

type PasswordChange struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type ProposedUser struct {
	Id        int    `json:"Id"`
	UserName  string `json:"userName"`
	FullName  string `json:"fullName"`
	Status    string `json:"status" enum:"enabled,disabled"`
	OrgUnitId int    `json:"orgUnitId"`
	RoleId    int    `json:"roleId"`
	Password  string `json:"password"`
}

type Role struct {
	Id           int    `json:"Id"`
	RoleName     string `json:"roleName"`
	Description  string `json:"description"`
	CreationDate string `json:"creationDate"`
}

type RolesList struct {
	Data []Role `json:"data"`
}

type User struct {
	Id                      int    `json:"Id"`
	UserName                string `json:"userName"`
	FullName                string `json:"fullName"`
	Status                  string `json:"status"`
	OrgUnitId               int    `json:"orgUnitId"`
	RoleId                  int    `json:"roleId"`
	PasswordHash            string `json:"passwordHash"`
	CreationDate            string `json:"creationDate"`
	LastPasswordChangedDate string `json:"lastPasswordChangedDate"`
}

type UserRoleId struct {
	RoleId int `json:"roleId"`
}

type UserRoleIdMsg struct {
	Message    string `json:"message"`
	UserRoleId int    `json:"roleId"`
}

type UsersList struct {
	Data []User `json:"data"`
}

type UserStatus struct {
	Status string `json:"status" enum:"enabled,disabled"`
}

type UserStatusMsg struct {
	Message    string `json:"message"`
	UserStatus string `json:"userStatus" enum:"enabled,disabled"`
}

type SuccessMsg struct {
	Message string `json:"message"`
}
