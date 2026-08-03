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
	/* send push route (INTERNAL ONLY) */
	sendPushRoute		= "POST /api/v1/internal/push/send"
	/* token management routes */
	registerTokenRoute 	= "POST /api/v1/push/register"
	deleteTokensRoute 	= "DELETE /api/v1/push/delete"
	deleteTokenRoute	= "DELETE /api/v1/push/delete/{id}"
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

	serviceKey := os.Getenv("SERVICE_KEY")
	if serviceKey == "" {
		log.Fatal("missing server key")
	}

	pushRepository := push.NewRepository(dbConn)

	firebaseProvider := push.NewFirebaseProvider(context.Background(), "TODO:")

	pushService := push.NewService(pushRepository, firebaseProvider)
	mailService := mailer.NewPostmarkService(postmarkAPIKey)

	pushHandler := push.NewHandler(pushService)
	mailHandler := mailer.NewHandler(mailService)

	a.mux.Handle(
		sendEmailRoute,
		xtoken.RequireServiceKey(serviceKey)(
			http.HandlerFunc(mailHandler.SendEmail),
		),
	)

	a.mux.Handle(
		sendPushRoute, 
		xtoken.RequireServiceKey(serviceKey)(
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

