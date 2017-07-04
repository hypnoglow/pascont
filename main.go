package main

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/hypnoglow/pascont/accounts"
	"github.com/hypnoglow/pascont/config"
	"github.com/hypnoglow/pascont/hasher"
	"github.com/hypnoglow/pascont/identity"
	"github.com/hypnoglow/pascont/kit/middleware"
	"github.com/hypnoglow/pascont/notary"
	"github.com/hypnoglow/pascont/packer"
	"github.com/hypnoglow/pascont/postgres"
	"github.com/hypnoglow/pascont/session"
	"github.com/hypnoglow/pascont/sessions"
)

const (
	// EnvConfigPath is an environment variable holding the path to the config.
	EnvConfigPath = "$PASCONT_CONFIG_PATH"
)

const (
	// ServerGracefulTimeout is a time for a server to wait for handlers finish their job before it shuts down.
	ServerGracefulTimeout = time.Second * 5
)

func main() {
	errorLogger := log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)

	// conf should not be passed to other application layers.
	// It should be used here in the entry point to create services
	// or to describe how to create services for other services.
	conf := getConfig()
	db := getDatabase(conf)
	if err := db.Ping(); err != nil {
		errorLogger.Println("WARNING: Failed to connect to the database at startup.")
	}
	sessionSecretKey := getValidSessionSecretKey(conf)

	// Repositories and services.
	accountRepo := postgres.NewAccountRepository(db)
	sessionRepo := postgres.NewSessionRepository(db)
	hmacNotary := notary.NewHMACNotary()
	base64Packer := packer.NewBase64Packer(session.SessionIDLength + session.SessionExpiresAtLength)
	bcryptHasher := hasher.NewBcryptHasher(bcrypt.DefaultCost)
	uuidv4 := identity.NewUUIDV4()

	// Controllers.
	accs := accounts.NewRestController(
		errorLogger,
		accountRepo,
		bcryptHasher,
		accounts.Options{},
	)
	sess := sessions.NewRestController(
		errorLogger,
		accountRepo,
		sessionRepo,
		hmacNotary,
		base64Packer,
		bcryptHasher,
		uuidv4,
		sessions.Options{
			SessionSecretKey: sessionSecretKey,
		},
	)

	// Routing and middleware.

	tokenExtractor := session.TokenExtractor(base64Packer, hmacNotary, sessionSecretKey)

	sessionsHander := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			http.HandlerFunc(sess.PostSessions).ServeHTTP(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	sessionHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			middleware.AuthToken(
				http.HandlerFunc(sess.GetSession),
				tokenExtractor,
			).ServeHTTP(w, req)
		case http.MethodPatch:
			middleware.AuthToken(
				http.HandlerFunc(sess.PatchSession),
				tokenExtractor,
			).ServeHTTP(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	accountsHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			http.HandlerFunc(accs.PostAccounts).ServeHTTP(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux := http.NewServeMux()
	mux.Handle(sessions.PathSessions, sessionsHander)
	mux.Handle(sessions.PathSession, middleware.PathID(sessionHandler, sessions.PathSession))
	mux.Handle(accounts.PathAccounts, accountsHandler)
	handler := middleware.Recover(mux, errorLogger)
	handler = middleware.Logger(handler, log.New(os.Stdout, "", log.LstdFlags))

	// Listen and Serve.

	socket := fmt.Sprintf("%s:%s", conf.Socket.Host, conf.Socket.Port)
	httpServer := &http.Server{Addr: socket, Handler: handler}

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		errorLogger.Printf("Listen on %s\n", socket)
		if err := httpServer.ListenAndServe(); err != nil {
			errorLogger.Printf("listen: %s\n", err)
		}
	}()

	<-stop

	errorLogger.Printf("Shutdown server...\n")
	ctx, _ := context.WithTimeout(context.Background(), ServerGracefulTimeout)
	httpServer.Shutdown(ctx)
}

func getConfig() (conf config.Config) {
	configPath := os.ExpandEnv(EnvConfigPath)
	if configPath == "" {
		panic("Config path MUST be specified using " + EnvConfigPath)
	}

	f, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}

	return config.FromJSON(f)
}

func getValidSessionSecretKey(conf config.Config) (key []byte) {
	key, err := hex.DecodeString(conf.Session.SecretKey)
	if err != nil {
		panic(err)
	}
	// The secret key MUST be 128 bit key, generated with a cryptographically secure
	// pseudo random number generator (CSPRNG).
	if len(key) != 16 {
		panic("config's Session.SecretKey MUST be 16 bytes long")
	}

	return key
}

func getDatabase(conf config.Config) (db *sql.DB) {
	dsn := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?%s",
		conf.Database.DriverName,
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.DatabaseName,
		url.Values(conf.Database.ConnectionParams).Encode(),
	)

	db, err := sql.Open(conf.Database.DriverName, dsn)
	if err != nil {
		panic(err)
	}

	return db
}
