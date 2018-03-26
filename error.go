package main

import (
        "github.com/codegangsta/martini-contrib/render"
)

type ErrResp struct {
        Err string
}

var ErrMap map[ErrorCode]string = map[ErrorCode]string{
        InternalError: "Internal Error occor",
}

type ErrorCode int

const InternalError ErrorCode = iota

func ReturnError(r render.Render, ECode ErrorCode) {
        e := ErrResp{Err: ErrMap[ECode]}
        r.JSON(500, e)
}