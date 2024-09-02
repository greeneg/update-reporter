package middleware

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
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/greeneg/update-reporterd/globals"
	"github.com/greeneg/update-reporterd/helpers"
	"github.com/greeneg/update-reporterd/model"
)

func processAuthorizationHeader(authHeader string) (string, string) {
	// split the header value at the space
	encodedString := strings.Split(authHeader, " ")

	// remove base64 encoding
	decodedString, _ := base64.StdEncoding.DecodeString(encodedString[1])

	// now lets return both the
	authValues := strings.Split(string(decodedString), ":")

	return authValues[0], authValues[1]
}

func AuthCheck(c *gin.Context) {
	var clientFingerprintHeader string = c.GetHeader("X-ASSIMILATOR-TYPE")
	// check if this is a machine logging in for DB access
	if clientFingerprintHeader == "MACHINE" {
		// now grab the token from the headers
		/*		authToken := c.GetHeader("X-Auth-Token")
				if authToken != "" {
					result, err := verifyMachineToken(authToken)
					if err != nil {
					}
					if result {

					} else {

					}
				} else {

				} */
	} else {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			log.Println("INFO: No session found. Attempting to check for authentication headers")
			baHeader := c.GetHeader("Authorization")
			if baHeader == "" {
				log.Println("ERROR: No authentication header found. Aborting")
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "not authorized!"})
				c.Abort()
				return
			}
			// otherwise, lets process that header
			username, password := processAuthorizationHeader(baHeader)
			authStatus := helpers.CheckUserPass(username, password)
			if authStatus {
				session.Set(globals.UserKey, username)
				if err := session.Save(); err != nil {
					c.IndentedJSON(http.StatusInternalServerError,
						gin.H{"error": "failed to save user session"})
					// session saving is not fatal, so allow them to proceed
				}
				log.Println("INFO: Authenticated")
			} else {
				log.Println("ERROR: Authentication failed. Aborting")
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "not authorized!"})
				c.Abort()
				return
			}
		} else {
			userString := fmt.Sprintf("%v", user)
			log.Println("INFO: Session found: User: " + userString)
			log.Println("INFO: Checking if user is locked or not...")
			user, err := model.GetUserByUserName(userString)
			if err != nil {
				log.Println("ERROR: " + string(err.Error()))
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "unable to authenticate: " + err.Error()})
				c.Abort()
				return
			}
			status := helpers.CheckIsNotLocked(user)
			if status {
				log.Println("INFO: Authenticated")
			} else {
				log.Println("WARN: User '" + userString + "' is locked!")
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "not authorized!"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
