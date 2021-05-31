package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lechgu/simtris/internal/blocks"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	sess := blocks.NewSession(conn)
	if err != nil {
		log.Fatal(err)
	}
	sess.Run()
}

//go:embed front/dist
var content embed.FS

func main() {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	router.Use(cors.New(config))

	fsys := fs.FS(content)
	root, err := fs.Sub(fsys, "front/dist")
	if err != nil {
		log.Fatalln(err)
	}
	router.StaticFS("/simtris/", http.FS(root))

	router.GET("/meta", func(c *gin.Context) {

		c.JSON(http.StatusOK, blocks.Metadata())
	})
	router.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/simtris/")
	})

	router.Run()
}
