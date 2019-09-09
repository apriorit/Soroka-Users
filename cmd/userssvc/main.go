package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"

	"github.com/Soroka-EDMS/svc/users/pkgs/config"
	"github.com/Soroka-EDMS/svc/users/pkgs/db"
	"github.com/Soroka-EDMS/svc/users/pkgs/endpoints"
	"github.com/Soroka-EDMS/svc/users/pkgs/handlers"
	"github.com/Soroka-EDMS/svc/users/pkgs/service"
)

func main() {
	//Command parameters

	var (
		address    = flag.String("address", ":443", "")
		consulAddr = flag.String("consul.address", "localhost:8500", "Consul agent address")
		conn       = flag.String("consul.db_connection", "stub", "database connection string")
	)

	//Parse CLI parameters
	flag.Parse()

	//Get global logger
	logger := config.GetLogger().Logger

	//Log CLI parameters
	logger.Log("address", *address, *consulAddr, "consul.db_connection", *conn)

	//Obtain consul k/v storage
	consulStorage, err := createConsulClient(*consulAddr)
	config.LogAndTerminateOnError(err, "Create consul client")

	//Obtain db connection string
	rawConnectionStr, err := ConsulGetKey(consulStorage, *conn)
	connectionStr := string(rawConnectionStr)
	logger.Log("database connection string ", connectionStr)

	//Connect to Users database
	logger.Log("Loading", "Connecting to users database...")
	db, err := db.Connection(logger, connectionStr)
	config.LogAndTerminateOnError(err, "connect to users database")

	//Build service layers
	var handler http.Handler
	{
		logger.Log("Loading", "Creating Users service...")
		svc := service.Build(logger, db)
		endpoints := endpoints.Build(logger, svc)
		handler = handlers.NewHTTPHandler(endpoints)
	}

	logger.Log("Loading", "Starting Users service...")
	var g run.Group
	{
		cer, err := tls.LoadX509KeyPair("serv.crt", "serv.key")
		config.LogAndTerminateOnError(err, "loading TLS pair")

		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}
		httpsListener, err := tls.Listen("tcp", ":443", tlsConfig)
		config.LogAndTerminateOnError(err, "create https listener")

		g.Add(func() error {
			logger.Log("transport", "https", "addr", *address)
			return http.Serve(httpsListener, handler)
		}, func(err error) {
			logger.Log("transport", "https", "err", err)
			httpsListener.Close()
		})
	}
	{
		var (
			cancelInterrupt = make(chan struct{})
			cancel          = make(chan os.Signal, 2)
		)

		defer close(cancel)

		g.Add(func() error {
			signal.Notify(cancel, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-cancel:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	logger.Log("exit", g.Run())
}
