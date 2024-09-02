package routes

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
	"github.com/gin-gonic/gin"

	"github.com/greeneg/update-reporterd/controllers"
)

func PrivateRoutes(g *gin.RouterGroup, u *controllers.UpdateReporter) {
	// Roles
	g.GET("/roles", u.GetRoles)                    // get all roles
	g.GET("/role/id/:roleId", u.GetRoleById)       // get role by Id
	g.GET("/role/name/:roleName", u.GetRoleByName) // get role by name
	g.POST("/role", u.CreateRole)                  // create new role
	g.DELETE("/role/:roleId", u.DeleteRole)        // delete a role by Id
	// user related routes
	g.GET("/users", u.GetUsers)                          // get all users
	g.GET("/users/roleId/:roleId", u.GetUsersByRoleId)   // get all users by role Id
	g.GET("/user/name/:name", u.GetUserByUserName)       // get a user by username
	g.GET("/user/name/:name/status", u.GetUserStatus)    // get whether a user is locked or not
	g.GET("/user/id/:id", u.GetUserById)                 // get a user by Id
	g.POST("/user", u.CreateUser)                        // create new user
	g.PATCH("/user/name/:name", u.ChangeAccountPassword) // update a user password
	g.PATCH("/user/name/:name/status", u.SetUserStatus)  // lock a user
	g.PATCH("/user/name/:name/roleId", u.SetUserRoleId)  // set a user's role Id
	g.DELETE("/user/name/:name", u.DeleteUser)           // trash a user
}

func PublicRoutes(g *gin.RouterGroup, u *controllers.UpdateReporter) {
	// service related routes
	g.GET("/health") // service health API
}
