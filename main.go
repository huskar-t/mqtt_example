package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/huskar-t/mqtt_example/config"
	"github.com/sirupsen/logrus"
	_ "github.com/taosdata/driver-go/v2/taosSql"
	"github.com/taosdata/go-utils/mqtt"
	"github.com/taosdata/go-utils/pool"
	"github.com/taosdata/go-utils/rule"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rulePath := flag.String("rc", "./config/rule.json", "rule config path")
	configPath := flag.String("c", "./config/config.json", "config path")
	flag.Parse()
	c := config.Init(*configPath)
	manager, err := rule.NewRuleManage(*rulePath)
	if err != nil {
		panic(err)
	}
	initDB(c.TDengine)
	dsn := fmt.Sprintf("%s:%s/tcp(%s:%d)/%s", c.TDengine.User, c.TDengine.Password, c.TDengine.Host, c.TDengine.Port, c.TDengine.DB)
	db, err := sql.Open("taosSql", dsn)
	createSql := manager.GenerateCreateSql()
	for _, s := range createSql {
		if c.ShowSql {
			fmt.Println(s)
		}
		_, err = db.Exec(s)
		if err != nil {
			panic(err)
		}
	}
	logger := logrus.WithFields(logrus.Fields{"project": "example"})
	connected := make(chan struct{})
	connector := mqtt.NewConnector(*c.MQTT, pool.GoroutinePool, logger, func() {
		connected <- struct{}{}
	})
	<-connected
	fmt.Println("connected")
	connector.SubscribeWithReceiveTime("#", 0, func(topic string, msg []byte, t time.Time) {
		if manager.RuleExist(topic) {
			r, err := manager.Parse(topic, msg)
			if err != nil {
				panic(err)
			}
			s := r.ToSql()
			if c.ShowSql {
				fmt.Println(s)
			}
			_, err = db.Exec(s)
			if err != nil {
				panic(err)
			}
		}
	})
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit
	fmt.Println("stop server")
}

func initDB(c *config.TDengine) {
	dsnWithoutDB := fmt.Sprintf("%s:%s/tcp(%s:%d)/", c.User, c.Password, c.Host, c.Port)
	db, err := sql.Open("taosSql", dsnWithoutDB)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("create database if not exists %s", c.DB))
	if err != nil {
		panic(err)
	}
}
