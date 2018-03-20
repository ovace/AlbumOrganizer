package convUtils

import (
	"strconv"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{})
}

func string2float64(s string) (f64 float64) {

	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Errorf("error converting string %v\t to float64: err: %v\n", s, err)
		return
	}
	return i
}

var String2float64 = string2float64
