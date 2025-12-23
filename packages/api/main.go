package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	dsn := "host=localhost user=postgres dbname=strshelf port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		println(err.Error())
	}

	r := gin.Default()

	r.Use(cors.Default())

	r.POST("/v1/item.get", func(ctx *gin.Context) {
		shelfItems, err := gorm.G[ShelfItem](db).Raw("SELECT * FROM public.shelf_item_v1 WHERE deleted IS NOT TRUE ORDER BY id DESC").Find(context.Background())
		if err != nil {
			ctx.JSON(500, "error in getting items")
			return
		}

		ctx.JSON(200, StrShelfResponse[[]ShelfItem]{
			Code: 200,
			Data: shelfItems,
			Msg:  "",
		})
	})
	r.POST("/v1/item.post", func(ctx *gin.Context) {
		newShelfItem := ShelfItem{}
		result := gorm.WithResult()
		err := ctx.ShouldBindJSON(&newShelfItem)
		if err != nil {
			ctx.JSON(500, "error in binding json: "+err.Error())
			fmt.Println(err.Error())
			return
		}

		err = gorm.G[any](db).Exec(context.Background(), "INSERT INTO public.shelf_item_v1 (title,link,comment,deleted) VALUES(?,?,?,?)", newShelfItem.Title, newShelfItem.Link, newShelfItem.Comment, false)

		if err != nil {
			ctx.JSON(500, "error in insert database: "+err.Error())
			fmt.Println(err.Error())
			return
		}
		ctx.JSON(200, PostRequestResponse{
			Code:   200,
			Result: result,
		})
	})
	r.POST("/v1/user.signup", func(ctx *gin.Context) {
		newUser := UserInfo{}
		result := gorm.WithResult()
		err := ctx.BindJSON(&newUser)
		if err != nil {
			ctx.JSON(500, "internal error")
			return
		}
		HashedPassword, err := HashPassword(newUser.Password)
		if err != nil {
			ctx.JSON(500, "internal error")
			return
		}

		err = gorm.G[any](db).Exec(context.Background(), "INSERT INTO public.shelf_user_v1 (username,password) VALUES(?,?)", newUser.Username, HashedPassword)
		if err != nil {
			ctx.JSON(500, "error in insert database: "+err.Error())
			fmt.Println(err.Error())
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
			ctx.JSON(500, "internal error")
			return
		}

		matchUsers, err := gorm.G[UserInfo](db).Raw("SELECT * FROM public.shelf_user_v1 WHERE username = ?", user.Username).Find(context.Background())
		if err != nil {
			ctx.JSON(500, "error in insert database: "+err.Error())
			fmt.Println(err.Error())
			return
		}
		if len(matchUsers) == 1 {
			matchUser := matchUsers[0]
			err := bcrypt.CompareHashAndPassword([]byte(matchUser.Password), []byte(user.Password))
			if err != nil {
				ctx.JSON(401, "password is incorrect")

				return
			} else {
				ctx.JSON(200, "login success")
			}
		} else {
			ctx.JSON(401, "user is not exist")
			return
		}

	})

	r.Run(":1111")

}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
