package handler

import (
	"azure_functions_go/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func GetRouteHandler(c *gin.Context) {
	name := c.DefaultQuery("name", "World")
	data := map[string]interface{}{
		"Greetings": "Hello " + name,
	}
	c.JSON(http.StatusOK, utils.Response(200, "Please hold while we process your safest route to your destination", data))
}

func searchRouteHandler(c *gin.Context) {
	//  get query from request
	query := c.DefaultQuery("query", "World")
	//  get azure maps sdk
	sdk := NewAzureMapsSDK(os.Getenv("AZURE_MAPS_API_KEY"))
	//  search route
	route, err := sdk.SearchRoute(query)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, utils.Response(500, "Failed to get route", nil))
		return
	}
	//  return route
	c.JSON(http.StatusOK, utils.Response(200, "Route found", route))
}

func searchFastestRouteHandler(c *gin.Context) {
	//  get query from request
	origin := c.DefaultQuery("origin", "World")
	destination := c.DefaultQuery("destination", "World")
	//  get azure maps sdk
	sdk := NewAzureMapsSDK(os.Getenv("AZURE_MAPS_API_KEY"))
	//  search route
	route, err := sdk.GetFastestRoute(origin, destination)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, utils.Response(500, "Failed to get route", nil))
		return
	}
	//  return route
	c.JSON(http.StatusOK, utils.Response(200, "Route found", route))
}

func main() {
	//  create a new router
	router := gin.Default()

	port := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if port == "" {
		port = "4670"
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	api := router.Group("/api")
	{
		api.GET("/safe_route_function", getRouteHandler)
		api.GET("/search_route", searchRouteHandler)
		api.GET("/fastest_route", searchFastestRouteHandler)
	}

	log.Printf("Starting Gin-based Azure Function on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}

}

type AzureMapsSDK struct {
	SubscriptionKey string
}

// Azure maps sdk constructor
func NewAzureMapsSDK(subscriptionKey string) *AzureMapsSDK {
	return &AzureMapsSDK{
		SubscriptionKey: subscriptionKey,
	}
}

func (sdk *AzureMapsSDK) SearchRoute(query string) (SearchResponse, error) {
	fmt.Println("Lets Check the subscription key : ", sdk.SubscriptionKey)
	//  get route from azure maps
	endpoint := fmt.Sprintf("https://atlas.microsoft.com/search/fuzzy/json?subscription-key=%s&api-version=1.0&query=%s", sdk.SubscriptionKey, query)
	res, err := http.Get(endpoint)
	if err != nil {
		log.Println(err)
		return SearchResponse{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Println(res.StatusCode)
		return SearchResponse{}, fmt.Errorf("failed to get route from azure maps")
	}
	//  parse response
	var response SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Println(err)
		return SearchResponse{}, err
	}
	if len(response.Results) == 0 {
		return SearchResponse{}, fmt.Errorf("no route found")
	}

	return response, nil
}

func (sdk *AzureMapsSDK) GetFastestRoute(originLongAndlat, destinationLongAndLat string) (FastestRouteResponse, error) {

	//  get fastest route from azure maps
	endpoint := fmt.Sprintf("https://atlas.microsoft.com/route/directions/json?subscription-key=%s&api-version=1.0&query=%s:%s&computeBestOrder=true&travelMode=car&routeType=fastest", sdk.SubscriptionKey, originLongAndlat, destinationLongAndLat)
	res, err := http.Get(endpoint)
	if err != nil {
		log.Println(err)
		return FastestRouteResponse{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Println(res.StatusCode)
		return FastestRouteResponse{}, fmt.Errorf("failed to get fastest route from azure maps")
	}
	//  parse response
	var response FastestRouteResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Println(err)
		return FastestRouteResponse{}, err
	}
	return response, nil
}

type SearchResponse struct {
	Results []Results `json:"results"`
}
type Position struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
type CategorySet struct {
	ID int `json:"id"`
}
type Poi struct {
	Name        string        `json:"name"`
	CategorySet []CategorySet `json:"categorySet"`
	Phone       string        `json:"phone"`
}
type Results struct {
	Type     string   `json:"type"`
	ID       string   `json:"id"`
	Position Position `json:"position"`
	Poi      Poi      `json:"poi"`
}

type FastestRouteResponse struct {
	Routes []Routes `json:"routes"`
}
type Summary struct {
	LengthInMeters        int `json:"lengthInMeters"`
	TravelTimeInSeconds   int `json:"travelTimeInSeconds"`
	TrafficDelayInSeconds int `json:"trafficDelayInSeconds"`
}
type LegsSummary struct {
	LengthInMeters      int `json:"lengthInMeters"`
	TravelTimeInSeconds int `json:"travelTimeInSeconds"`
}
type Points struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type Legs struct {
	LegsSummary LegsSummary `json:"summary"`
	Points      []Points    `json:"points"`
}
type Routes struct {
	Summary Summary `json:"summary"`
	Legs    []Legs  `json:"legs"`
}
