// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package migrate

import (
	"fmt"

	"github.com/tsuru/tsuru/db"
	"github.com/tsuru/tsuru/router"
	"gopkg.in/mgo.v2/bson"
)

type AppWithPlanRouter struct {
	Name   string
	Plan   PlanWithRouter
	Router string
}

type PlanWithRouter struct {
	Router string
}

func MigrateAppPlanRouterToRouter() error {
	defaultRouter, err := router.Default()
	if err != nil {
		if err == router.ErrDefaultRouterNotFound {
			fmt.Println("A default router must be configured in order to run this migration.")
			fmt.Println("To fix this, either set the \"docker:router\" or \"router:<router_name>:default\" configs.")
		}
		return err
	}
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	iter := conn.Apps().Find(nil).Iter()
	var app AppWithPlanRouter
	for iter.Next(&app) {
		if app.Router != "" {
			continue
		}
		r := defaultRouter
		if app.Plan.Router != "" {
			r = app.Plan.Router
		}
		err = conn.Apps().Update(bson.M{"name": app.Name}, bson.M{"$set": bson.M{"router": r}})
		if err != nil {
			return err
		}
	}
	return nil
}
