package app

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"context"
	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oyvinddd/xtoken"
	"github.com/oyvinddd/messaging-api/internal/push"
	"github.com/oyvinddd/messaging-api/internal/mailer"
)

const (
	/* send push route */
	sendPushRoute		= "POST /api/v1/push/send"
	/* token management routes */
	registerTokenRoute 	= "POST /api/v1/push/register"
	deleteTokensRoute 	= "DELETE /api/v1/push/delete"
	/* email routes */
	sendEmailRoute 		= "POST /api/v1/email/send"
)

type (
	App struct {
		addr	string
		server *http.Server
		mux    *http.ServeMux
	}
)

func New() *App {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("unable to load environment: %v\n", err)
	}

	addr := fmt.Sprintf(":%s", os.Getenv("LISTENING_PORT"))
	mux := http.NewServeMux()

	app := &App{
		addr: addr,
		mux: mux,
		server:	&http.Server{
			Handler: mux,
			Addr: addr,
		},
	}

	app.registerRoutes()

	return app
}

func (a *App) registerRoutes() {

	dbURI := os.Getenv("POSTGRES_URI")
	if dbURI == "" {
		log.Fatal("missing Postgres connection string")
	}

	dbConn, err := setupDBConnection(context.Background(), os.Getenv("POSTGRES_URI")) 
	if err != nil {
		log.Fatalf("unable to connect to db: %v\n", err)
	}

	postmarkAPIKey := os.Getenv("POSTMARK_API_KEY")
	if postmarkAPIKey == "" {
		log.Fatal("missing Postmark API key")
	}

	// set HMAC secret for auth middleware
	hmac := os.Getenv("HMAC_SECRET")
	if hmac == "" {
		log.Fatal("missing HMAC secret")
	}
	xtoken.SetHMACSecret(hmac)

	pushRepository := push.NewRepository(dbConn)

	firebaseProvider := push.NewFirebaseProvider(context.Background(), "TODO:")

	pushService := push.NewService(pushRepository, firebaseProvider)
	mailService := mailer.NewPostmarkService(postmarkAPIKey)

	pushHandler := push.NewHandler(pushService)
	mailHandler := mailer.NewHandler(mailService)

	// TODO: make these two APIs internal to the Docker network and not publicly available
	a.mux.HandleFunc(sendEmailRoute, mailHandler.SendEmail)
	a.mux.HandleFunc(sendPushRoute, pushHandler.SendPush)

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
			http.HandlerFunc(pushHandler.RegisterToken),
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

