package main

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/piecehealth/actioncable"
)

var cable *actioncable.Cable

func main() {
	router := gin.Default()

	router.SetTrustedProxies(nil)

	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")

	// setup actioncable
	cbCfg := actioncable.NewConfig().WithDebugLogger()

	// use redis PubSub
	// cbCfg = cbCfg.WithRedisPubSub(&redis.Options{Addr: "localhost:6379"})

	cbCfg = cbCfg.WithAuthenticator(func(h *http.Request) (any, bool) {
		name, err := h.Cookie("userName")

		if err != nil {
			return nil, false
		}

		id, _ := url.QueryUnescape(name.Value)

		return id, true
	})

	cable = actioncable.NewActionCable(cbCfg)
	cable.RegisterChannel(roomChannel)

	defer cable.Stop()

	router.GET("/cable", func(c *gin.Context) {
		cable.Handle(c.Writer, c.Request)
	})

	// homepage
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{})
	})

	// login request
	router.POST("/create_session", func(c *gin.Context) {
		name := c.PostForm("name")

		if name == "" {
			c.Redirect(http.StatusFound, "/")

			return
		}

		c.SetCookie("userName", name, 0, "", "", true, true)
		c.Redirect(http.StatusFound, "/rooms")
	})

	authorized := router.Group("/")
	authorized.Use(checkLogin())
	{
		// room list
		authorized.GET("/rooms", func(c *gin.Context) {
			c.HTML(http.StatusOK, "rooms.tmpl", gin.H{
				"roomIds": []string{"1", "2", "3", "4", "forbidden_room"},
			})
		})

		// room
		authorized.GET("/rooms/:id", func(c *gin.Context) {
			c.HTML(http.StatusOK, "room.tmpl", gin.H{
				"roomId": c.Param("id"),
			})
		})

		// logout
		authorized.POST("/destroy_session", func(c *gin.Context) {
			c.SetCookie("userName", "", 0, "", "", true, true)

			c.Redirect(http.StatusFound, "/")
		})
	}

	router.Run(":8080")
}

func checkLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := c.Cookie("userName"); err != nil {
			c.Redirect(http.StatusFound, "/")
		}
	}
}
