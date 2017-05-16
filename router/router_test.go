// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package router

import (
	"testing"

	"fmt"

	"github.com/ivessong/beautify-gateway/models"
	"github.com/valyala/fasthttp"
)

func TestRouterLookup(t *testing.T) {
	router := New()
	ctx := &fasthttp.RequestCtx{}

	// try empty router first
	api, tsr := router.Lookup("GET", "/nope", ctx)
	if api != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", api)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}

	// insert route and try again
	router.Handle(&models.API{Path: "/user/:name", Method: "GET"})

	api, tsr = router.Lookup("GET", "/user/gopher", ctx)
	if api == nil {
		t.Fatal("Got no handle!")
	} else {
		fmt.Printf("[%v],%v \n", api.Method, api.Path)
	}
	if ctx.UserValue("name") != "gopher" {
		t.Error("Param not set!")
	}

	api, tsr = router.Lookup("GET", "/user/gopher/", ctx)
	if api != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", api)
	}
	if !tsr {
		t.Error("Got no TSR recommendation!")
	}

	api, tsr = router.Lookup("GET", "/nope", ctx)
	if api != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", api)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}
}
