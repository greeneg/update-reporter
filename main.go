package main

/*

  Update-Reporterd - Golang-based Web Service Dashboard for Software Updates

  Author:  Gary L. Greene, Jr.
  License: Apache v2.0

  Copyright 2024, YggdrasilSoft, LLC

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
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/greeneg/update-reporterd/controllers"
	_ "github.com/greeneg/update-reporterd/docs"
	"github.com/greeneg/update-reporterd/globals"
	"github.com/greeneg/update-reporterd/helpers"
	"github.com/greeneg/update-reporterd/middleware"
	"github.com/greeneg/update-reporterd/model"
	"github.com/greeneg/update-reporterd/routes"
)

//	@title		Update Reporter Daemon
//	@version	0.1.0
//	@description	An API for Reporting Software Updates

//	@contact.name	Gary Greene
//	@contact.url	https://github.com/greeneg/update-reporterd

//	@securityDefinitions.basic	BasicAuth

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8000
//	@BasePath	/api/v1

//	@schemas	http https

func createDB(dbName string) (bool, error) {
	log.Println("INFO: DB doesn't exist. Attempt to create it")
	const schema string = `CREATE TABLE IF NOT EXISTS Architectures (
		Id                      INTEGER		PRIMARY KEY AUTOINCREMENT	UNIQUE	NOT NULL,
		ArchName                STRING		UNIQUE				NOT NULL,
		CreationDate            DATETIME	NOT NULL			DEFAULT (CURRENT_TIMESTAMP)
	);

	INSERT INTO Architectures (Id, ArchName) VALUES (1, 'noarch');
	INSERT INTO Architectures (Id, ArchName) VALUES (2, 'aarch64');
	INSERT INTO Architectures (Id, ArchName) VALUES (3, 'x86');
	INSERT INTO Architectures (Id, ArchName) VALUES (4, 'x86_64');

	CREATE TABLE IF NOT EXISTS OperatingSystems (
		Id                      INTEGER		PRIMARY KEY AUTOINCREMENT	UNIQUE	NOT NULL,
		OsIdName                STRING		NOT NULL,
		OsVersion               STRING		NOT NULL,
		CreationDate            DATETIME	NOT NULL			DEFAULT (CURRENT_TIMESTAMP)
	);

	CREATE TABLE IF NOT EXISTS OsFamilies (
		Id                      INTEGER		PRIMARY KEY AUTOINCREMENT	UNIQUE	NOT NULL,
		FamilyName              STRING		UNIQUE				NOT NULL,
		CreationDate            DATETIME	NOT NULL			DEFAULT (CURRENT_TIMESTAMP)
	);

	INSERT INTO OsFamilies (Id, FamilyName) VALUES (1, 'linux');
	INSERT INTO OsFamilies (Id, FamilyName) VALUES (2, 'darwin');
	INSERT INTO OsFamilies (Id, FamilyName) VALUES (3, 'windows');

	CREATE TABLE IF NOT EXISTS Roles (
		Id                      INTEGER		PRIMARY KEY AUTOINCREMENT	UNIQUE	NOT NULL,
		RoleName                STRING		UNIQUE				NOT NULL,
		Description             STRING		NOT NULL,
		CreationDate            DATETIME	NOT NULL		 	DEFAULT (CURRENT_TIMESTAMP)
	);

	INSERT INTO Roles (Id, RoleName, Description)
		VALUES (1, 'SYSTEM', 'Built-in system role');
	INSERT INTO Roles (Id, RoleName, Description)
		VALUES (2, 'administrators', 'Accounts that have full administrative rights to the system');

	CREATE TABLE IF NOT EXISTS Systems (
		Id                      INTEGER		PRIMARY KEY AUTOINCREMENT		UNIQUE	NOT NULL,
		FQDN                    STRING		UNIQUE					NOT NULL,
		OsFamilyId              INTEGER		REFERENCES OsFamilies (Id)		NOT NULL,
		OsId                    INTEGER		REFERENCES OperatingSystems (Id)	NOT NULL,
		ArchId                  INTEGER		REFERENCES Architectures (Id)		NOT NULL,
		CreationDate            DATETIME	NOT NULL				DEFAULT (CURRENT_TIMESTAMP)
	);

	CREATE TABLE IF NOT EXISTS UpdateRecords (
		Id                      INTEGER		PRIMARY KEY AUTOINCREMENT		UNIQUE	NOT NULL,
		SystemId                INTEGER		REFERENCES Systems (Id)			UNIQUE	NOT NULL,
		UpdateCount		INTEGER		NOT NULL,
		UpdateRecord            JSON		NOT NULL,
		CreationDate            INTEGER		NOT NULL				DEFAULT (CURRENT_TIMESTAMP)
		LastUpdateDate          DATETIME	NOT NULL				DEFAULT (CURRENT_TIMESTAMP)
	);

	CREATE TABLE IF NOT EXISTS Users (
		Id                      INTEGER 	PRIMARY KEY AUTOINCREMENT	UNIQUE	NOT NULL,
		UserName                STRING		NOT NULL			UNIQUE,
		FullName                STRING		NOT NULL,
		Status                  STRING		NOT NULL			DEFAULT enabled,
		RoleId                  INTEGER		REFERENCES Roles (Id)		NOT NULL,
		PasswordHash            STRING		NOT NULL,
		CreationDate            DATETIME	NOT NULL			DEFAULT (CURRENT_TIMESTAMP),
		LastPasswordChangedDate DATETIME	NOT NULL			DEFAULT (CURRENT_TIMESTAMP)
	);

	INSERT INTO Users (Id, UserName, FullName, Status, RoleId, PasswordHash)
		VALUES (1, 'SYSTEM', 'Built-in System User', 'enabled', 1, '!');
	`

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		helpers.FatalCheckError(err)
	}
	if _, err := db.Exec(schema); err != nil {
		helpers.FatalCheckError(err)
	}
	return true, err
}

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// lets get our working directory
	appdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	helpers.FatalCheckError(err)

	// config path is derived from app working directory
	configDir := filepath.Join(appdir, "config")

	// now that we have our appdir and configDir, lets read in our app config
	// and marshall it to the Config struct
	config := globals.Config{}
	jsonContent, err := os.ReadFile(filepath.Join(configDir, "config.json"))
	helpers.FatalCheckError(err)
	err = json.Unmarshal(jsonContent, &config)
	helpers.FatalCheckError(err)

	// create an app object that contains our routes and the configuration
	UpdateReporter := new(controllers.UpdateReporter)
	UpdateReporter.AppPath = appdir
	UpdateReporter.ConfigPath = configDir
	UpdateReporter.ConfStruct = config

	if _, err := os.Stat(UpdateReporter.ConfStruct.DbPath); errors.Is(err, os.ErrNotExist) {
		_, err := createDB(UpdateReporter.ConfStruct.DbPath)
		if err != nil {
			helpers.FatalCheckError(err)
		}
	}

	err = model.ConnectDatabase(UpdateReporter.ConfStruct.DbPath)
	helpers.FatalCheckError(err)

	// set up our static assets
	// r.Static("/assets", "./assets")
	// r.LoadHTMLGlob("templates/*.html")

	// some defaults for using session support
	r.Use(sessions.Sessions("session", cookie.NewStore(globals.Secret)))
	// frontend
	// fePublic := r.Group("/")
	// routes.FePublicRoutes(fePublic, AllocatorD)

	// fePrivate := r.Group("/")
	// fePrivate.Use(middleware.AuthCheck)
	// routes.FePrivateRoutes(fePrivate, AllocatorD)

	// API
	public := r.Group("/api/v1")
	routes.PublicRoutes(public, UpdateReporter)

	private := r.Group("/api/v1")
	private.Use(middleware.AuthCheck)
	routes.PrivateRoutes(private, UpdateReporter)

	// swagger doc
	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	tcpPort := strconv.Itoa(UpdateReporter.ConfStruct.TcpPort)
	tlsTcpPort := strconv.Itoa(UpdateReporter.ConfStruct.TLSTcpPort)
	tlsPemFile := UpdateReporter.ConfStruct.TLSPemFile
	tlsKeyFile := UpdateReporter.ConfStruct.TLSKeyFile
	if UpdateReporter.ConfStruct.UseTLS {
		r.RunTLS(":"+tlsTcpPort, tlsPemFile, tlsKeyFile)
	} else {
		r.Run(":" + tcpPort)
	}
}
