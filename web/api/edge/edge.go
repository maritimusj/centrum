package edge

import (
	"fmt"
	"github.com/kataras/iris"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func Feedback(deviceID int64, ctx iris.Context) {
	data ,err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Printf("%d => %s\r\n", deviceID, string(data))
	_, _ = ctx.JSON(iris.Map{
		"id": deviceID,
		"status": "Ok",
	})
}