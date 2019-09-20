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
		conn       = flag.String("consul.usersdb", "usersdb", "database connection string")
		certKey    = flag.String("consul.tls.pubkey", "tls/pubKey", "tls certificate")
		privateKey = flag.String("consul.tls.privkey", "tls/privKey", "tls private key")
		secret     = flag.String("consul.service.secret", "service/secret", "Secret to sign JWT")
	)

	//Parse CLI parameters
	flag.Parse()

	//Get global logger
	logger := config.GetLogger().Logger

	//Log CLI parameters
	logger.Log(
		"address", *address,
		"consul.address", *consulAddr,
		"consul.usersdb", *conn,
		"certKey", *certKey,
		"privateKey", *privateKey,
		"secret", *secret,
	)

	//Obtain consul k/v storage
	consulStorage, err := createConsulClient(*consulAddr)
	config.LogAndTerminateOnError(err, "Create consul client")

	//Obtain db connection string
	rawConnectionStr, err := ConsulGetKey(consulStorage, *conn)
	config.LogAndTerminateOnError(err, "obtain database connection string")
	connectionStr := string(rawConnectionStr)

	//Obtain tls pair (raw)
	certKeyData, err := ConsulGetKey(consulStorage, *certKey)
	config.LogAndTerminateOnError(err, "obtain certificate key")
	privateKeyData, err := ConsulGetKey(consulStorage, *privateKey)
	config.LogAndTerminateOnError(err, "obtain private key")

	//Get sign secret
	signSecret, err := ConsulGetKey(consulStorage, *secret)
	config.LogAndTerminateOnError(err, "obtain secret")

	//Connect to Users database
	logger.Log("Loading", "Connecting to users database...")
	dbs, err := db.Connection(logger, connectionStr)
	config.LogAndTerminateOnError(err, "connect to users database")

	//Build service layers
	var handler http.Handler
	{
		logger.Log("Loading", "Creating Users service...")
		svc := service.Build(logger, dbs)
		endpoints := endpoints.Build(logger, svc, certKeyData, signSecret)
		handler = handlers.NewHTTPHandler(endpoints)
	}

	logger.Log("Loading", "Starting Users service...")
	var g run.Group
	{
		cert, err := tls.X509KeyPair(certKeyData, privateKeyData)

		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
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
