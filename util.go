package tweethog

import (
	"math/rand"
	"time"
	"os/user"
	"path/filepath"
)

func GetRandomInt(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func GetExpandedFilename(filename string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if filename[:2] == "~/" {
		filename = filepath.Join(dir, filename[2:])
	}

	return filename
}