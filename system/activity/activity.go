package activity

import (
	"github.com/glennyonemitsu/reicon/system/intent"
)

type Activity struct {
	Name    string
	Intents []*intent.Intent
}
