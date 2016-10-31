package bg

import (
	"github.com/astaxie/beego"
	"net/http"
	"fmt"
)

type BaseController struct {
	beego.Controller
}

func (c*BaseController)Serve(v interface{}) {
	c.Data["json"] = v
	c.ServeJSON()
}

func (c*BaseController)ServeStatus(status int) {
	c.ServeGeneral(http.StatusText(status), status)
}
func (c*BaseController)ServeGeneral(message string, status int) {
	c.Data["json"] = General{
		Code:status,
		Message:message,
	}
	c.ServeJSON()
	c.Ctx.ResponseWriter.WriteHeader(status)
}

func (c*BaseController)ServeCreated(v interface{}) {
	switch v.(type){
	case string:
		c.Data["json"] = Created{ID:v.(string)}
		c.ServeJSON()
		c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	default:
		panic(fmt.Sprintf("Can not handler created %v", v))
	}
}

func (c*BaseController)ServeSuccess() {
	c.Data["json"] = General{Code:0}
	c.ServeJSON()
}
