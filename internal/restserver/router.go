package restserver

import (
	"net/http"
	"strings"

	"github.com/colbymilton/marchamps-valuator/internal/controller"
	"github.com/gin-gonic/gin"
)

var server *Server

type Server struct {
	// controller
	ctrl *controller.Valuator
}

func Run() {
	if server != nil {
		panic("already running!")
	}

	server = &Server{}

	// init controller
	server.ctrl = controller.NewValuator()
	if err := server.ctrl.UpdateDatabase(); err != nil {
		panic(err)
	}

	// init router
	router := gin.Default()
	router.SetTrustedProxies([]string{"http://localhost"})

	router.GET("/packs", server.GetPacks)
	router.GET("/latest_deck", server.GetLatestDeck)
	router.GET("/card_value/:code", server.GetCardValue)
	router.GET("/pack_value/:code", server.GetPackValue)
	router.GET("/pack_values", server.GetAllPackValues)

	router.Run(":9999")
}

func respond(c *gin.Context, body any, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, body)
	}
}

func (s *Server) GetPacks(c *gin.Context) {
	b, err := s.ctrl.GetPacks()
	respond(c, b, err)
}

func (s *Server) GetLatestDeck(c *gin.Context) {
	b, err := s.ctrl.GetNewestDeck()
	respond(c, b, err)
}

func (s *Server) GetCardValue(c *gin.Context) {
	code := c.Param("code")
	b, err := s.ctrl.ValueCard(code)
	respond(c, b, err)
}

func (s *Server) GetPackValue(c *gin.Context) {
	code := c.Param("code")
	b, err := s.ctrl.ValuePack(code)
	respond(c, b, err)
}

func (s *Server) GetAllPackValues(c *gin.Context) {
	ownedStr := c.Query("owned")
	owned := strings.Split(ownedStr, ",")
	b, err := s.ctrl.ValueAllPacks(owned)
	respond(c, b, err)
}
