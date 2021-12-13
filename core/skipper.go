package core

import (
	"os"
	"strconv"
	"time"

	"github.com/jordi-reinsma/bagel/db"
)

var runIntervalDays = os.Getenv("RUN_INTERVAL_DAYS")

func ShouldSkipExecution(dB db.DB) (bool, error) {
	date, err := dB.GetLastExecutionDate()
	if err != nil {
		return true, err
	}
	elapsedDays := int(time.Since(date).Hours() / 24)

	interval, err := strconv.Atoi(runIntervalDays)
	if err != nil {
		return true, err
	}

	return elapsedDays < interval, nil
}
