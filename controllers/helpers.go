package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/greeneg/update-reporterd/model"
)

func (u *UpdateReporter) GetUserId(c *gin.Context) (model.User, bool) {
	// need to get our current user context to get the CreatorId
	session := sessions.Default(c)
	user := session.Get("user")
	// if nil, we have an issue
	if user == nil {
		return model.User{}, false
	}

	// convert user interface to a string
	username := fmt.Sprintf("%v", user)
	// lets output our session user
	log.Println("INFO: Session user: " + username)
	// get our user id
	userObject, err := model.GetUserByUserName(username)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return model.User{}, false
	}

	// what is our user Id
	log.Println("INFO: Session user's ID: " + strconv.Itoa(userObject.Id))
	return userObject, true
}
