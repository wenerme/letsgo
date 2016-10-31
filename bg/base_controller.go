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
func (c*BaseController)ServeError(err error) {
	if err == nil {
		c.ServeSuccess()
		return
	}
	// 检测 err 是否为特殊的错误结构来附带更多的错误信息
	c.Data["json"] = General{
		Code: http.StatusInternalServerError,
		Message: err.Error(),
		Path: c.Ctx.Input.URL(),
	}
	c.ServeJSON()
	c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
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
	case uint:
		c.Data["json"] = Created{ID:v.(uint)}
		c.ServeJSON()
		//c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	default:
		panic(fmt.Sprintf("Can not handler created %v", v))
	}
}

func (c*BaseController)ServeSuccess() {
	c.Data["json"] = General{Code:0}
	c.ServeJSON()
}
