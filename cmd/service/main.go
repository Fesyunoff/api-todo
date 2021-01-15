package main

import (
	"fmt"

	"github.com/fesyunoff/api/pkg/auth"
	"github.com/fesyunoff/api/pkg/db"
	"github.com/fesyunoff/api/pkg/msg"
	"github.com/fesyunoff/api/pkg/todolist"
	"github.com/fesyunoff/api/pkg/transport"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	kitjsonrpc "github.com/l-vitaly/go-kit/transport/http/jsonrpc"

	"net"
	"net/http"
	"os"
	"os/signal"
)

const (
	serviceName = "todolist"
)

func main() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "service", serviceName)
	logger = kitlog.With(logger, "timestamp", kitlog.DefaultTimestampUTC)
	logger = kitlog.With(logger, "caller", kitlog.Caller(5))

	sqliteDB := db.CreateSQLiteDB("todo.db")
	defer sqliteDB.Close()
	strg := db.NewSQLiteTodoStorage(sqliteDB)
	msg := msg.NewMessanger()
	svc := todolist.NewService(strg, msg)
	svcHandlers, err := transport.MakeHandlerJSONRPC(
		svc,
		transport.GenericServerOptions(
			kitjsonrpc.ServerBefore(auth.PopulateRequestContext),
			kitjsonrpc.ServerErrorLogger(logger),
		),
		transport.GenericServerEndpointMiddlewares(
			auth.Middleware(),
		),
	)
	bindAddr := fmt.Sprintf("%s:%d", "0.0.0.0", 8991)
	r := mux.NewRouter().StrictSlash(true)

	exitOnError(logger, err, "failed create handlers")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	r.PathPrefix("/").Handler(svcHandlers)

	ln, err := net.Listen("tcp", bindAddr)
	exitOnError(logger, err, "failed create listener")
	defer ln.Close()

	_ = level.Info(logger).Log("msg", "server listen on "+ln.Addr().String())

	go func() {
		_ = http.Serve(ln, r)
	}()

	<-sigint
}

func exitOnError(l kitlog.Logger, err error, msg string) {
	if err != nil {
		l.Log("err", err, "msg", msg)
		os.Exit(1)
	}
}
