package app

import (
	"fmt"
	"log"
	"net/http"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oyvinddd/xtoken"
	"github.com/oyvinddd/messaging-api/internal/push"
	"github.com/oyvinddd/messaging-api/internal/mailer"
)

type (
	App struct {
		addr	string
		server *http.Server
		mux    *http.ServeMux
	}
)

func New(cfg Config) *App {
	addr := fmt.Sprintf(":%s", cfg.listeningPort)
	mux := http.NewServeMux()

	app := &App{
		addr: addr,
		mux: mux,
		server:	&http.Server{
			Handler: mux,
			Addr: addr,
		},
	}

	app.registerRoutes(cfg)

	return app
}

func (a *App) registerRoutes(cfg Config) {

	dbConn, err := setupDBConnection(context.Background(), cfg.dbURI) 
	if err != nil {
		log.Fatalf("unable to connect to db: %v\n", err)
	}

	// set HMAC secret for auth middleware
	xtoken.SetHMACSecret(cfg.hmacSecret)

	pushRepository := push.NewRepository(dbConn)

	firebaseProvider := push.NewFirebaseProvider(context.Background(), "TODO:")

	pushService := push.NewService(pushRepository, firebaseProvider)
	mailService := mailer.NewPostmarkService(cfg.mailerAPIKey)

	pushHandler := push.NewHandler(pushService)
	mailHandler := mailer.NewHandler(mailService)

	a.mux.Handle(
		sendEmailRoute,
		xtoken.RequireServiceKey(cfg.serviceKey)(
			http.HandlerFunc(mailHandler.SendEmail),
		),
	)

	a.mux.Handle(
		sendPushRoute, 
		xtoken.RequireServiceKey(cfg.serviceKey)(
			http.HandlerFunc(pushHandler.SendPush),
		),
	)

	a.mux.Handle(
		registerTokenRoute,
		xtoken.Authorize(
			http.HandlerFunc(pushHandler.RegisterToken),
			xtoken.UserRole,
		),
	)

	a.mux.Handle(
		deleteTokensRoute,
		xtoken.Authorize(
			http.HandlerFunc(pushHandler.DeleteTokens),
			xtoken.UserRole,
		),
	)

	a.mux.Handle(
		deleteTokenRoute,
		xtoken.Authorize(
			http.HandlerFunc(pushHandler.DeleteToken),
			xtoken.UserRole,
		),
	)

	a.mux.Handle(
		deleteTokenRoute,
		xtoken.Authorize(
			http.HandlerFunc(pushHandler.DeleteToken),
			xtoken.UserRole,
		),
	)
}

func setupDBConnection(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	// since New() doesn't wait to check if a connection was established,
	// we'll try to ping the db right after to verify that we are connected
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}

func (a *App) Run() error {
	fmt.Printf("Starting messaging service on %v...\n", a.addr)
	return a.server.ListenAndServe()
}

