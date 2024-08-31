package dep

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/cors"
	// _ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"

	"github.com/online-bnsp/backend/api"
	"github.com/online-bnsp/backend/middleware/auth"
	"github.com/online-bnsp/backend/util/buckets"
	"github.com/online-bnsp/backend/util/buckets/discard"
	"github.com/online-bnsp/backend/util/buckets/local"
	s3b "github.com/online-bnsp/backend/util/buckets/s3"
	"github.com/online-bnsp/backend/util/http/httpclient"
	"github.com/online-bnsp/backend/util/logger"
	"github.com/online-bnsp/backend/util/otpsender/whatsapp"
	queue "github.com/online-bnsp/backend/util/queue"
	"github.com/online-bnsp/backend/util/s3"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type DI struct {
	db         *sql.DB
	redis      *redis.Client
	apiHandler http.Handler
}

func InitDI(configFile string) (*DI, error) {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	di := &DI{}

	// JWT Secret
	auth.SetJWTConfig(viper.GetString("jwt.secret"), viper.GetDuration("jwt.ttl"), viper.GetDuration("jwt.refresh_ttl"))

	return di, nil
}

func (di *DI) GetDatabase() (*sql.DB, error) {
	if di.db == nil {
		db, err := sql.Open(viper.GetString("dbdriver"), viper.GetString("dsn"))
		if err != nil {
			return nil, err
		}

		err = db.Ping()
		if err != nil {
			return nil, err
		}

		di.db = db
	}
	return di.db, nil
}

func (di *DI) GetRedis() (*redis.Client, error) {
	if di.redis == nil {
		rdb := redis.NewClient(&redis.Options{
			Addr:     viper.GetString("redis_host"),
			Password: viper.GetString("redis_pass"), // no password set
			DB:       viper.GetInt("redis_db"),      // use default DB
		})
		di.redis = rdb
	}

	return di.redis, nil
}

func (di *DI) GetQueue() (queue.Queuer, error) {
	q, err := queue.NewQueue("nsq", viper.GetString("nsqd"))
	if err != nil {
		return nil, err
	}

	return q, nil
}

func (di *DI) GetBucket() (bucket buckets.Bucket) {
	provider := viper.GetString("bucket.provider")

	switch provider {
	case "s3":
		bucket = s3b.New(s3.NewS3(
			viper.GetString("bucket.s3.uri"),
			viper.GetString("bucket.s3.access_key"),
			viper.GetString("bucket.s3.secret_key"),
			viper.GetString("bucket.s3.token_key"),
			viper.GetString("bucket.s3.region"),
			viper.GetString("bucket.s3.bucket_name"),
		))

	case "local":
		bucket = local.New(
			viper.GetString("bucket.local.path"),
			viper.GetString("server_addr"),
			viper.GetString("bucket.local.url_prefix"),
			nil,
		)

	case "discard":
		bucket = &discard.Bucket{}

	default:
		log.Printf("unknown bucket provider `%v`, using `discard` bucket provider\n", provider)
		bucket = &discard.Bucket{}
	}

	return
}

func (di *DI) GetAPIHandler() (http.Handler, error) {
	if di.apiHandler == nil {
		db, _ := di.GetDatabase()
		q, err := di.GetQueue()
		if err != nil {
			return nil, err
		}

		bucket := di.GetBucket()
		rdb, _ := di.GetRedis()

		corsHandler := cors.Handler(cors.Options{
			AllowedOrigins: viper.GetStringSlice("cors.allowed_origins"),
			// AllowOriginFunc: func(r *http.Request, origin string) bool { return true },
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			// ExposedHeaders:   []string{"Link"},
			// AllowCredentials: false,
			// MaxAge:           300, // Maximum value not ignored by any of major browsers
		})
		di.apiHandler = api.New(db, rdb, q, bucket, di.GetMailer(), corsHandler).Handler()

		// bucket local server
		if v, ok := bucket.(*local.Bucket); ok {
			v.Handler = di.apiHandler
			di.apiHandler = v
		}
	}

	return di.apiHandler, nil
}

func (di *DI) Whatsapp() *whatsapp.Client {
	zlogger := logger.New(logger.Config{
		Level:  viper.GetString("logger_level"),
		Output: viper.GetString("logger_output"),
	}).With().Logger()

	// GLOBAL HTTP CLIENT
	cfgHTTPClient := httpclient.Config{
		DialTimeout:       60 * time.Second,
		ConnectionTimeout: 60 * time.Second,
		IdleTimeout:       10 * time.Second,
		MaxConn:           100,
		MaxIdleConn:       1000,
	}
	loggedHTTPTransport := httpclient.NewLoggedTransport(zlogger, cfgHTTPClient.NewTransport())
	httpClient := httpclient.New(cfgHTTPClient, httpclient.WithHTTPTransport(loggedHTTPTransport))

	return whatsapp.New(viper.GetString("whatsapp.url"), viper.GetString("whatsapp.basic_auth"), httpClient)
}

func (di *DI) GetAPIServer() (*http.Server, error) {
	h, err := di.GetAPIHandler()
	if err != nil {
		return nil, err
	}

	srv := http.Server{
		Addr:    viper.GetString("server_addr"),
		Handler: h,
	}
	return &srv, nil
}
