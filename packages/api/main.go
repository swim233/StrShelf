package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.ilharper.com/strshelf/api/config"
	"gopkg.ilharper.com/strshelf/api/logger"
	"gopkg.ilharper.com/strshelf/api/middleware"
	"gopkg.ilharper.com/strshelf/api/token"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
	Account  string `json:"account"`
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
	dsn := "host=localhost user=postgres dbname=strshelf port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Suger.Errorf("error when connect db: %s", err.Error())
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/v1/item.get", func(ctx *gin.Context) {
		shelfItems, err := gorm.G[ShelfItem](db).Raw("SELECT * FROM public.shelf_item_v1 WHERE deleted IS NOT TRUE ORDER BY id DESC").Find(context.Background())
		if err != nil {
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}

		ctx.JSON(200, StrShelfResponse[[]ShelfItem]{
			Code: 200,
			Data: shelfItems,
			Msg:  "",
		})
	})
	r.POST("/v1/item.post", middleware.JWTAuthMiddleWare(), func(ctx *gin.Context) {
		newShelfItem := ShelfItem{}
		result := gorm.WithResult()
		err := ctx.ShouldBindJSON(&newShelfItem)
		if err != nil {
			ctx.JSON(500, gin.H{"msg": "error in binding json: " + err.Error()})
			logger.Suger.Errorf("error in getting items: %s", err.Error())
			return
		}

		err = gorm.G[any](db).Exec(context.Background(), "INSERT INTO public.shelf_item_v1 (title,link,comment,deleted) VALUES(?,?,?,?)", newShelfItem.Title, newShelfItem.Link, newShelfItem.Comment, false)

		if err != nil {
			ctx.JSON(500, gin.H{"msg": "internal error"})
			logger.Suger.Errorf("error in insert database: %s", err.Error())
			return
		}
		ctx.JSON(200, PostRequestResponse{
			Code:   200,
			Result: result,
		})
	})

	r.POST("/v1/user.edit", middleware.JWTAuthMiddleWare(), func(ctx *gin.Context) {
		editShelfItem := ShelfEditItem{}
		result := gorm.WithResult()
		err := ctx.ShouldBind(&editShelfItem)
		if err != nil {
			logger.Suger.Errorf("error in building editShelf: %s", err.Error())
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}
		shelfItems, err := gorm.G[ShelfItem](db).Raw("SELECT * FROM public.shelf_item_v1 WHERE deleted IS NOT TRUE AND id = ?", editShelfItem.Id).Find(context.Background())
		if err != nil {
			logger.Suger.Errorf("error in checking item: %s", err.Error())
			ctx.JSON(400, gin.H{"msg": "internal error"})
			return
		}
		if len(shelfItems) != 1 {
			ctx.JSON(200, gin.H{"msg": "origin item not found", "code": "404"})
			return
		}
		err = gorm.G[any](db).Exec(context.Background(), "UPDATE public.shelf_item_v1 SET title = ? ,link = ? , comment = ? ,gmt_modified = now() WHERE id = ?", editShelfItem.NewTitle, editShelfItem.NewLink, editShelfItem.NewComment, editShelfItem.Id)

		if err != nil {
			logger.Suger.Errorf("error in updating database: %s", err.Error())
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}
		ctx.JSON(200, gin.H{"msg": "ok", "result": result})
	})

	r.POST("/v1/user.delete", middleware.JWTAuthMiddleWare(), func(ctx *gin.Context) {
		deleteItem := ShelfDeleteItem{}
		err := ctx.BindJSON(&deleteItem)
		if err != nil {
			logger.Suger.Errorf("error when binding json: %s", err.Error())
			ctx.JSON(400, gin.H{"msg": "internal error"})
			return
		}
		shelfItems, err := gorm.G[ShelfItem](db).Raw("SELECT * FROM public.shelf_item_v1 WHERE deleted IS NOT TRUE AND id = ?", deleteItem.Id).Find(context.Background())
		if err != nil {
			logger.Suger.Errorf("error in checking item: %s", err.Error())
			ctx.JSON(400, gin.H{"msg": "internal error"})
			return
		}
		if len(shelfItems) != 1 {
			ctx.JSON(200, gin.H{"msg": "origin item not found", "code": "404"})
			return
		}
		err = gorm.G[any](db).Exec(context.Background(), "UPDATE public.shelf_item_v1 SET deleted = true ,gmt_deleted = now() WHERE id = ?", deleteItem.Id)

		if err != nil {
			logger.Suger.Errorf("error in updating database: %s", err.Error())
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}
		ctx.JSON(200, gin.H{"msg": "ok"})

	})

	r.POST("/v1/user.signup", func(ctx *gin.Context) {
		newUser := UserInfo{}
		result := gorm.WithResult()
		err := ctx.BindJSON(&newUser)
		if err != nil {
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}
		HashedPassword, err := HashPassword(newUser.Password)
		if err != nil {
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}

		err = gorm.G[any](db).Exec(context.Background(), "INSERT INTO public.shelf_user_v1 (username,password) VALUES(?,?)", newUser.Account, HashedPassword)
		if err != nil {
			ctx.JSON(500, gin.H{"msg": "internal error"})
			logger.Suger.Errorf("error in insert database: %s", err.Error())
			return
		}
		ctx.JSON(200, PostRequestResponse{
			Code:   200,
			Result: result,
		})

	})
	r.POST("/v1/user.login", func(ctx *gin.Context) {
		user := UserInfo{}
		err := ctx.BindJSON(&user)
		if err != nil {
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}
		logger.Suger.Infoln(user)

		matchUsers, err := gorm.G[UserInfo](db).Raw("SELECT * FROM public.shelf_user_v1 WHERE username = ?", user.Account).Find(context.Background())
		if err != nil {
			ctx.JSON(500, gin.H{"msg": "error in insert database: " + err.Error()})
			fmt.Println(err.Error())
			return
		}
		if len(matchUsers) == 1 {
			matchUser := matchUsers[0]
			err := bcrypt.CompareHashAndPassword([]byte(matchUser.Password), []byte(user.Password))
			if err != nil {
				ctx.JSON(401, gin.H{"msg": "password is incorrect"})
				logger.Suger.Warnf("a failure logging request: %s", err.Error())
				return
			} else {
				token := token.CreateJWT(user.Account)
				ctx.JSON(200, gin.H{"token": token})
			}
		} else {
			ctx.JSON(401, gin.H{"msg": "user is not exist"})
			return
		}

	})
	r.POST("/v1/user.verify", func(ctx *gin.Context) {
		var tokenReq token.TokenRequest
		err := ctx.Bind(&tokenReq)
		if err != nil {
			logger.Suger.Errorf("error in verifying user token: %s", err.Error())
			ctx.JSON(500, gin.H{"msg": "internal error"})
			return
		}
		if result, err := token.VerifyJWT(string(tokenReq.Token)); err != nil {
			logger.Suger.Warnf("error in verifying user token: %s", err.Error())
			logger.Suger.Warnf("received token: %s", tokenReq.Token)
			ctx.JSON(401, gin.H{
				"msg": "token is invalid"})
			return
		} else {

			switch result {
			case true:
				ctx.JSON(200, gin.H{"msg": "successful"})
			case false:
				logger.Suger.Warnf("user post a error token: %s", tokenReq.Token)
				ctx.JSON(401, gin.H{"msg": "token is invalid"})
			}
			return
		}
	})

	err = r.Run(":1111")
	if err != nil {
		logger.Suger.Panicf("fail to start http service: %s", err.Error())
	}

}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
