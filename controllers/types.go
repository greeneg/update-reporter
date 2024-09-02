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

import "github.com/greeneg/update-reporterd/globals"

type UpdateReporter struct {
	AppPath    string
	ConfigPath string
	ConfStruct globals.Config
}

type SafeUser struct {
	Id           int
	UserName     string
	FullName     string
	OrgUnitId    int
	RoleId       int
	CreationDate string
}
