package etc

import (
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
)

//LogDuration of execution
func LogDuration(label string, callback func()) {
	begin := time.Now()
	callback()
	log.Info(fmt.Sprintf("execution time of %s: ", label), time.Now().Sub(begin))
}
