package appcontext

import (
	"errors"

	"github.com/aburizalpurnama/travel/internal/builder/driver"
	"github.com/aburizalpurnama/travel/internal/config"
	"gorm.io/gorm"
)

const (
	// DBDialectMysql rdbms dialect name for MySQL
	DBDialectMysql = "mysql"

	// DBDialectPostgres rdbms dialect name for PostgreSQL
	DBDialectPostgres = "postgres"
)

// AppContext the app context struct
type AppContext struct {
	cfg *config.Config
}

// NewAppContext initiate appcontext object
func NewAppContext(config *config.Config) *AppContext {
	return &AppContext{
		cfg: config,
	}
}

// GetDBInstance getting gorm instance, param: dbType can be "mysql" or "postgre"
func (a *AppContext) GetDBInstance(dbType string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	switch dbType {
	case DBDialectMysql:
		db, err = driver.NewMysqlDatabase(a.cfg)
	case DBDialectPostgres:
		db, err = driver.NewPostgreDatabase(a.cfg)
	default:
		err = errors.New("error get db instance, unknown db type")
	}

	return db, err
}

// // GetCachePool get cache pool connection
// func (a *AppContext) GetCachePool() *redis.Pool {
// 	return driver.NewCache(a.getCacheOption())
// }

// func (a *AppContext) getCacheOption() driver.CacheOption {
// 	return driver.CacheOption{
// 		Host:               a.cfg.GetString("CACHE_HOST"),
// 		Port:               a.cfg.GetInt("CACHE_PORT"),
// 		Namespace:          a.cfg.GetString("CACHE_NAMESPACE"),
// 		Password:           a.cfg.GetString("CACHE_PASSWORD"),
// 		DialConnectTimeout: a.cfg.GetDuration("CACHE_DIAL_CONNECT_TIMEOUT"),
// 		ReadTimeout:        a.cfg.GetDuration("CACHE_READ_TIMEOUT"),
// 		WriteTimeout:       a.cfg.GetDuration("CACHE_WRITE_TIMEOUT"),
// 		IdleTimeout:        a.cfg.GetDuration("CACHE_IDLE_TIMEOUT"),
// 		MaxConnLifetime:    a.cfg.GetDuration("CACHE_CONN_LIFETIME_MAX"),
// 		MaxIdle:            a.cfg.GetInt("CACHE_CONN_IDLE_MAX"),
// 		MaxActive:          a.cfg.GetInt("CACHE_CONN_ACTIVE_MAX"),
// 		Wait:               a.cfg.GetBool("CACHE_IS_WAIT"),
// 	}
// }

// func (a *AppContext) PubSubClient(ctx context.Context) (c *pubsub.Client, err error) {
// 	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", a.config.GetString("PUBSUB_CLIENT_KEY"))

// 	c, err = pubsub.NewClient(
// 		ctx,
// 		a.config.GetString("GCP_PROJECT_ID"),
// 	)
// 	if err != nil {
// 		err = errorx.New(err.Error())
// 	}
// 	return
// }

// func (a *AppContext) CloudStorageClient(ctx context.Context) (c *storage.Client, err error) {
// 	c, err = storage.NewClient(ctx, option.WithCredentialsFile(a.config.GetString("GCS_CLIENT_KEY")))
// 	if err != nil {
// 		err = errorx.New(err.Error())
// 	}

// 	return
// }

// func (a *AppContext) FirebaseMessagingClient(ctx context.Context) (c *firebaseMessaging.Client, err error) {
// 	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(a.config.GetString("GCS_CLIENT_KEY")))
// 	if err != nil {
// 		err = errorx.New(err.Error())
// 		log.Err(err).Send()
// 	}
// 	c, err = app.Messaging(ctx)
// 	if err != nil {
// 		err = errorx.New(err.Error())
// 		log.Err(err).Send()
// 	}

// 	return
// }
