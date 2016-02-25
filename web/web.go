package web

import (
	"fmt"

	"gopkg.in/macaron.v1"

	"github.com/prime/middleware"
	"github.com/prime/router"
	"github.com/wrench/setting"
	"github.com/wrench/db"
)

func SetPrimeMacaron(m *macaron.Macaron) {
	//Setting Database
	if err := db.InitDB(setting.DBURI, setting.DBPasswd, setting.DBDB); err != nil {
		fmt.Printf("Connect Database error %s", err.Error())
	}

	if err := middleware.Initfunc(); err != nil {
		fmt.Printf("Init middleware error %s", err.Error())
	}

	//Setting Middleware
	middleware.SetMiddlewares(m)

/*	//Start Object Storage Service if sets in conf
	if strings.EqualFold(setting.OssSwitch, "enable") {
		ossobj := oss.Instance()
		ossobj.StartOSS()
	}
*/
	//Setting Router
	router.SetRouters(m)
}
