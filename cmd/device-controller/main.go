/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package main

import (
	"github.com/nalej/device-controller/cmd/device-controller/commands"
	"github.com/nalej/device-controller/version"
)

var MainVersion string

var MainCommit string

func main() {
	version.AppVersion = MainVersion
	version.Commit = MainCommit
	commands.Execute()
}