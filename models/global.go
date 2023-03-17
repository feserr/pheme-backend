package models

import (
	"github.com/feserr/pheme-backend/database"
)

// Db global db var
var Db = database.Connect()
