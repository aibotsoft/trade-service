package main

import (
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/mig"
	"github.com/aibotsoft/micro/sqlserver"
	"github.com/aibotsoft/trade/pkg/store"
	"github.com/aibotsoft/trade/services/auth"
	"github.com/aibotsoft/trade/services/handler"
)

func main() {
	log := logger.New()
	cfg := config.New()
	log.Infow("Begin service", "conf", cfg.Service)
	db := sqlserver.MustConnectX(cfg)
	err := mig.MigrateUp(cfg, log, db)
	if err != nil {
		log.Fatal(err)
	}
	sto := store.New(cfg, log, db)

	au := auth.New(cfg, log, sto)
	go au.AuthJob()
	h := handler.New(cfg, log, sto, au)
	h.WebsocketJob()
	//log.Info(h)
	//cli := client.New(cfg, log)
	//ctx := context.Background()

	//auth := context.WithValue(ctx, api.ContextAPIKeys, map[string]api.APIKey{"session": {Key: token}})
	//resp, err := cli.Events(auth)
	//log.Info(resp)
	//log.Info(err)

}
