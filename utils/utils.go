package utils

import (
	"time"
	"github.com/briandowns/spinner"
)

func GetSpinner(prefix string) *spinner.Spinner {
	loader := spinner.New(spinner.CharSets[11], 100 * time.Millisecond)
	loader.Color("cyan")
	loader.Prefix = "\033[36m" + prefix + "\033[m"
	return loader
}
