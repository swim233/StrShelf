package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.ilharper.com/strshelf/api/config"
	"gopkg.ilharper.com/strshelf/api/db"
	"gopkg.ilharper.com/strshelf/api/handler"
	"gopkg.ilharper.com/strshelf/api/logger"
	"gopkg.ilharper.com/strshelf/api/utils"
)

//go:embed dist
var app embed.FS

type ShelfItem struct {
	Id          uint64     `json:"id"`
	Title       string     `json:"title"`
	Link        string     `json:"link"`
	Comment     string     `json:"comment"`
	GMTCreated  CustomTime `json:"gmt_created"`
	GMTModified CustomTime `json:"gmt_modified"`
	GMTDeleted  CustomTime `json:"gmt_deleted"`
	Deleted     bool       `json:"deleted"`
}
type ShelfEditItem struct {
	Id         uint64 `json:"id"`
	NewTitle   string `json:"new_title"`
	NewLink    string `json:"new_link"`
	NewComment string `json:"new_comment"`
}

type ShelfDeleteItem struct {
	Id uint64 `json:"id"`
}

type StrShelfResponse[T any] struct {
	Code uint16 `json:"code"`
	Data T      `json:"data"`
	Msg  string `json:"msg"`
}

type PostRequestResponse struct {
	Code   uint16 `json:"code"`
	Result any    `json:"result"`
}

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CustomTime time.Time

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	t := time.Time(ct)
	return json.Marshal(t.UnixMilli())
}
func (ct *CustomTime) UnmarshalJSON(data []byte) error {

	var ms int64
	if err := json.Unmarshal(data, &ms); err == nil {
		*ct = CustomTime(time.Unix(0, ms*int64(time.Millisecond)))
		return nil
	} else {
		return fmt.Errorf("cannot unmarshal %s into CustomTime: %w", string(data), err)
	}
}

func main() {
	logger.InitLogger()
	config.InitConfig()
	db.InitPostgresDB()
	logger.Suger.Infoln(utils.GetVersion())
	if !utils.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	handler.ItemGetHandler(r)
	handler.ItemPostHandler(r)
	handler.ItemEditHandler(r)
	handler.ItemDeleteHandler(r)
	handler.UserSignUpHandler(r)
	handler.UserLoginHandler(r)
	handler.UserVerifyHandler(r)

	dist, err := static.EmbedFolder(app, "dist")
	if err != nil {
		logger.Suger.Panicf("can not load embed folder: %s", err.Error())
	}
	staticServer := static.Serve("/", dist)

	r.Use(staticServer)
	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method == http.MethodGet &&
			!strings.ContainsRune(c.Request.URL.Path, '.') &&
			!strings.HasPrefix(c.Request.URL.Path, "/v1/") {
			c.Request.URL.Path = "/"
			staticServer(c)
		}
	})
	port := viper.GetString("port")
	if port == "" {
		logger.Suger.Warnln("can not read port in config,using default port :1111")
	}
	listen := viper.GetString("listen")
	if listen == "" {
		logger.Suger.Warnln("can not read address in config,listening 0.0.0.0")
		listen = "0.0.0.0"
	}
	err = r.Run(listen + ":" + port)
	if err != nil {
		logger.Suger.Panicf("fail to start http service: %s", err.Error())
	}

}
