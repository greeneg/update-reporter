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
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/greeneg/update-reporterd/globals"
	"golang.org/x/sys/unix"
)

// getConfig Returns our configuration
func getConfig() (globals.Config, error) {
	config := globals.Config{}
	appdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return config, err
	}
	configDir := filepath.Join(appdir, "config")
	jsonContent, err := os.ReadFile(filepath.Join(configDir, "config.json"))
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(jsonContent, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// checkDbIsPresent Returns whether the DB is present
func checkDbIsPresent() (bool, error) {
	// get our config
	config, err := getConfig()
	if err != nil {
		return false, err
	}

	dbFile := config.DbPath
	if _, err := os.Stat(dbFile); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, err
	} else {
		// should never get here, this is a schroedinger case
		return false, err
	}
}

// checkDiskSpace Returns whether the disk that houses the working directory has adequate space
func checkDiskSpace() (bool, error) {
	// get the amount of space from the disk
	var stat unix.Statfs_t

	// get our working directory
	wd, err := os.Getwd()
	if err != nil {
		return false, err
	}

	// retrieve our statfs data
	unix.Statfs(wd, &stat)

	freeSpaceInGiB := stat.Bavail * uint64(stat.Bsize) / (1024 * 1024 * 1024)
	if freeSpaceInGiB > 10 {
		return true, nil
	} else {
		// no error, but we don't have adequate safe space on the disk
		return false, nil
	}
}

// checkDiskIsWritable Returns whether the disk that houses the DB is writable
func checkDiskIsWritable() (bool, error) {
	config, err := getConfig()
	if err != nil {
		return false, err
	}
	err = unix.Access(config.DbPath, unix.W_OK)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetHealth Retrieve the health of the service
//
//	@Summary		Retrieve overall health of the service
//	@Description	Retrieve overall health of the service
//	@Tags			serviceHealth
//	@Produce		json
//	@Success		200	{object}	model.HealthCheck
//	@Failure		500	{object}	model.HealthCheck
//	@Router			/health [get]
func (u *UpdateReporter) GetHealth(c *gin.Context) {
	log.Println("INFO: Retrieving Health status")

	// Determine our own host's health
	// - is the db present?
	// - do we have adequate disk space?
	// - is our disk even writable?

	// db exists?
	dbStatus, err := checkDbIsPresent()
	var dbStatusString string
	if dbStatus {
		dbStatusString = "OK"
	} else {
		dbStatusString = "UNHEALTHY"
	}
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"db":            dbStatusString + ": " + string(err.Error()),
			"diskSpace":     "UNKNOWN",
			"diskWriteable": "UNKNOWN",
			"health":        "UNHEALTHY",
			"status":        500,
		})
		return
	}

	// disk space?
	diskSpaceStatus, err := checkDiskSpace()
	var diskSpaceStatusString string
	if diskSpaceStatus {
		diskSpaceStatusString = "OK"
	} else {
		diskSpaceStatusString = "UNHEALTHY"
	}
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"db":           dbStatusString,
			"diskSpace":    diskSpaceStatusString + ": " + string(err.Error()),
			"diskWritable": "UNKNOWN",
			"health":       "UNHEALTHY",
			"status":       500,
		})
		return
	}

	// disk is writable?
	diskIsWritableStatus, err := checkDiskIsWritable()
	var diskIsWritableStatusString string
	if diskIsWritableStatus {
		diskIsWritableStatusString = "OK"
	} else {
		diskIsWritableStatusString = "UNHEALTHY"
	}
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"db":           dbStatusString,
			"diskSpace":    diskSpaceStatusString,
			"diskWritable": diskIsWritableStatusString + ": " + string(err.Error()),
			"health":       "UNHEALTHY",
			"status":       500,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"db":           "OK",
		"diskSpace":    "OK",
		"diskWritable": "OK",
		"health":       "OK",
		"status":       200,
	})
}
