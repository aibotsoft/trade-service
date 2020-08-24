package main

import (
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/mig"
	"github.com/aibotsoft/micro/sqlserver"
)

func main() {
	log := logger.New()
	cfg := config.New()
	log.Infow("Begin service", "conf", cfg.Service)
	db := sqlserver.MustConnectX(cfg)

	log.Info(db)
	err := mig.MigrateUp(cfg, log, db)
	if err != nil {
		log.Fatal(err)
	}
	//cli := client.New(cfg, log)
	//ctx := context.Background()
	////resp, err := cli.CheckLogin(ctx, token)
	////if err != nil {
	////	return
	////}
	//auth := context.WithValue(ctx, api.ContextAPIKeys, map[string]api.APIKey{"session": {Key: token}})
	//resp, err := cli.Events(auth)
	//log.Info(resp)
	//log.Info(err)

}
