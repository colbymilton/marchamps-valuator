package restserver

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/colbymilton/marchamps-valuator/internal/controller"
	"github.com/gin-contrib/cors"
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
	if err := server.ctrl.Update(); err != nil {
		panic(err)
	}

	// init router
	router := gin.Default()
	router.SetTrustedProxies([]string{"http://localhost"})

	// CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	router.Use(cors.New(config))

	router.GET("/packs", server.GetPacks)
	router.GET("/pack_values", server.GetAllPackValues)
	router.GET("/card_values", server.GetAllCardValues)

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

func (s *Server) GetAllPackValues(c *gin.Context) {
	ownedStr := c.Query("owned")

	aspectWeights := make(map[string]float64)
	if aggressionWeight := c.Query("aw"); aggressionWeight != "" {
		f, err := strconv.ParseFloat(aggressionWeight, 64)
		if err != nil {
			respond(c, nil, err)
			return
		}

		aspectWeights["aggression"] = f
	}
	if protectionWeight := c.Query("pw"); protectionWeight != "" {
		f, err := strconv.ParseFloat(protectionWeight, 64)
		if err != nil {
			respond(c, nil, err)
			return
		}

		aspectWeights["protection"] = f
	}
	if leadershipWeight := c.Query("lw"); leadershipWeight != "" {
		f, err := strconv.ParseFloat(leadershipWeight, 64)
		if err != nil {
			respond(c, nil, err)
			return
		}

		aspectWeights["leadership"] = f
	}
	if justiceWeight := c.Query("jw"); justiceWeight != "" {
		f, err := strconv.ParseFloat(justiceWeight, 64)
		if err != nil {
			respond(c, nil, err)
			return
		}

		aspectWeights["justice"] = f
	}

	owned := strings.Split(ownedStr, ",")
	b, err := s.ctrl.ValueAllPacks(owned, aspectWeights)
	respond(c, b, err)
}

func (s *Server) GetAllCardValues(c *gin.Context) {
	ownedStr := c.Query("owned")
	owned := strings.Split(ownedStr, ",")
	b, err := s.ctrl.ValueAllCards(owned)
	respond(c, b, err)
}
