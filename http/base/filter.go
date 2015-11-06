package base

import (
	"github.com/astaxie/beego/context"
	"github.com/open-falcon/fe/model/uic"
	"strconv"
	"time"
)

var FilterLoginUser = func(ctx *context.Context) {
	cookieSig := ctx.GetCookie("sig")
	if cookieSig == "" {
		ctx.Redirect(302, "/auth/login?callback="+ctx.Request.URL.String())
		return
	}

	sessionObj := uic.ReadSessionBySig(cookieSig)
	if sessionObj == nil || int64(sessionObj.Expired) < time.Now().Unix() {
		ctx.Redirect(302, "/auth/login?callback="+ctx.Request.URL.String())
		return
	}

	u := uic.ReadUserById(sessionObj.Uid)
	if u == nil {
		ctx.Redirect(302, "/auth/login?callback="+ctx.Request.URL.String())
		return
	}

	ctx.Input.SetData("CurrentUser", u)
}

var FilterTargetUser = func(ctx *context.Context) {
	userId := ctx.Input.Query("id")
	if userId == "" {
		ctx.ResponseWriter.WriteHeader(403)
		ctx.ResponseWriter.Write([]byte("id is necessary"))
		return
	}

	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		ctx.ResponseWriter.WriteHeader(403)
		ctx.ResponseWriter.Write([]byte("id is invalid"))
		return
	}

	u := uic.ReadUserById(id)
	if u == nil {
		ctx.ResponseWriter.WriteHeader(403)
		ctx.ResponseWriter.Write([]byte("no such user"))
		return
	}

	ctx.Input.SetData("TargetUser", u)
}

var FilterTargetTeam = func(ctx *context.Context) {
	tid := ctx.Input.Query("id")
	if tid == "" {
		ctx.ResponseWriter.WriteHeader(403)
		ctx.ResponseWriter.Write([]byte("id is necessary"))
		return
	}

	id, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		ctx.ResponseWriter.WriteHeader(403)
		ctx.ResponseWriter.Write([]byte("id is invalid"))
		return
	}

	t := uic.ReadTeamById(id)
	if t == nil {
		ctx.ResponseWriter.WriteHeader(403)
		ctx.ResponseWriter.Write([]byte("no such team"))
		return
	}

	ctx.Input.SetData("TargetTeam", t)
}
