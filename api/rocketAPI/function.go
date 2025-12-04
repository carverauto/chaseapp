// Package p contains an HTTP Cloud Function.
package p

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	// firebase "firebase.google.com/go"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"strings"
	// "reflect"

	"strconv"

	// "html"
	"log"
	"net/http"
	"os"

	// 	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/google/uuid"
	"github.com/mfreeman451/golang/common/writers"
	"github.com/mfreeman451/helpers"
	"github.com/nav-inc/datetime"
	"github.com/pusher/push-notifications-go"
	"github.com/pusher/pusher-http-go"
	"go.opencensus.io/trace"
	"google.golang.org/api/iterator"
)

/* GLOBALS */

// Used to get launch info and store it in firebase
var rocketLaunchAPIKey = os.Getenv("ROCKETLAUNCHAPI")

// APIKEY is used to do some hokey authentication with our API
var APIKEY = os.Getenv("APIKEY")

// ServerClient, _ := stream.NewClient(getStream_API_KEY, getStream_API_SECRET)
var getStream_API_KEY = os.Getenv("GETSTREAM_API_KEY")
var getStream_API_SECRET = os.Getenv("GETSTREAM_API_SECRET")

// AISHub username
var aisUsername = os.Getenv("AISHUB_USERNAME")

// GCLOUD_PROJECT is automatically set by the Cloud Functions runtime.
var projectID = os.Getenv("GCLOUD_PROJECT")

// client is a Firestore client, reused between function invocations.
var client *firestore.Client

// pusher channels
var PUSHER_API_APPID = os.Getenv("PUSHER_API_APPID")
var PUSHER_API_KEY = os.Getenv("PUSHER_API_KEY")
var PUSHER_API_SECRET = os.Getenv("PUSHER_API_SECRET")
var PUSHER_API_CLUSTER = os.Getenv("PUSHER_API_CLUSTER")

var PUSHER_BEAMS_INSTANCE = os.Getenv("PUSHER_BEAMS_INSTANCE_ID")
var PUSHER_BEAMS_SECRET = os.Getenv("PUSHER_BEAMS_SECRET")
// pusher beams

var exporter *stackdriver.Exporter

// init
func init() {
	var err error

	exporter, err := stackdriver.NewExporter(stackdriver.Options{})
	if err != nil {
		log.Fatalf("Failed to create the Stackdriver exporter: %v", err)
	}

	/*
	// It is imperative to invoke flush before your main function exits
	defer sd.Flush()

	// Start the metrics exporter
	sd.StartMetricsExporter()
	defer sd.StopMetricsExporter()
	 */
	// Register it as a trace exporter
	trace.RegisterExporter(exporter)
	// trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
}

type Coords struct {
	Lat		float64
	Lon		float64
}

/********/
/* USGS */
/********/

// GetQuakes retrieves data from the USGS earthquake API around significant events, that
// is updated every minute. Data is stored in firebase realtime database and rendered
// on the map.
func GetQuakes(w http.ResponseWriter, r *http.Request) {
	// TODO: oc-http package, need to look into
	// TODO: do we need to use stackdriverCtx instead of the other context we create below??
	stackdriverCtx, sp := trace.StartSpan(r.Context(),"carverauto.com/API.GetQuakes")
	defer sp.End()

	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID, DatabaseURL: "https://chaseapp-8459b.firebaseio.com" }

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	var quakeURL = "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/significant_hour.geojson"
	resp, err := http.Get(quakeURL)

	if err != nil {
		println("Error handling request: #{err}")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Body:", body)
	fmt.Fprintf(w, "Body: %v", body)

	// TODO: remove, just testing oepncensus stuff..
	sp.AddAttributes(trace.BoolAttribute("gotJSON", true))

	var q GeoJSON
	err = json.Unmarshal(body, &q)
	if err != nil {
		log.Fatalln("Error unmarshalling JSON: ", err)
	}

	log.Println("Length of Features", len(q.Features))
	fmt.Fprintf(w, "Length of features: %v", len(q.Features))

	rtClient, rtErr := app.Database(ctx)
	if rtErr != nil {
		log.Fatalln("Error initializing db client: ", rtErr)
	}

	rtRef := rtClient.NewRef("quakes")
	rtRef.Child("usgs")
	setErr := rtRef.Set(stackdriverCtx, q)
	if setErr != nil {
	 log.Fatalln("Error setting value:", setErr)
	}
}

/*******/
/* AIS */
/*******/

type AISBoat struct {
	MMSI			float64		`json:"mmsi"`
	Longitude		float64		`json:"longitude"`
	Latitude		float64		`json:"latitude"`
	COG				float32		`json:"cog"`
	SOG				int			`json:"sog"`
	Heading			int			`json:"heading"`
	ROT				int			`json:"rot"`
	NavSat			int			`json:"nav_sat"`
	IMO				int			`json:"imo"`
	Name			string		`json:"name"`
	Callsign		string		`json:"callsign"`
	Type			int			`json:"type"`
	A				int			`json:"a"`
	B				int			`json:"b"`
	C				int			`json:"c"`
	D				int			`json:"d"`
	Draught			float32		`json:"draught"`
	Dest			string		`json:"dest"`
	ETA				string		`json:"eta"`
}

func GetBoats(w http.ResponseWriter, r *http.Request) {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID, DatabaseURL: "https://chaseapp-8459b.firebaseio.com" }

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("client.Close: #{err}")
		}
	}(client)
	type Boat struct {
		Group string `firestore:"group"`
		MMSI  string `firestore:"mmsi"`
		Type  string `firestore:"type"`
	}

	boats := client.Collection("boats")
	docs, err := boats.Documents(ctx).GetAll()
	if err != nil {
		println("Error retrieving documents: %v", err)
		return
	}

	var boatsArray []string
	for _, doc := range docs {
		var boatData Boat
		if err := doc.DataTo(&boatData); err != nil {
			println("Error retrieving documents: %v", err)
			return
		}
		if len(boatData.MMSI) > 0 {
			boatsArray = append(boatsArray, boatData.MMSI)
		}
	}

	result := strings.Join(boatsArray, ",")

	var aisOptions = "?username=" + aisUsername + "&format=1&output=json&compress"
	var aisHubURL = "https://data.aishub.net/ws.php" + aisOptions + "&mmsi=" + result
	// log.Println("AisHUBurl: " + aisHubURL)

	resp, err := http.Get(aisHubURL)

	if err != nil {
		println("Error handling request: #{err}")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var aisHubBoats [][]AISBoat

	rtClient, rtErr := app.Database(ctx)
	if rtErr != nil {
		log.Fatalln("Error initializing db client: ", rtErr)
	}

	json.Unmarshal(body, &aisHubBoats)

	rtRef := rtClient.NewRef("ships")
	for i := 0; i < len(aisHubBoats[1]); i++ {
		rtRef.Child(aisHubBoats[1][i].Name)
		setErr := rtRef.Set(ctx, aisHubBoats)
		if setErr != nil {
			log.Fatalln("Error setting value:", setErr)
		}
	}
}

/***********/
/* ROCKETS */
/***********/

// GetLaunches will look through the list of launches returned
// by the external API and add to the 'launches' collection
// a document with the information about the launch and rocket.
// We are only going to add objects that are X hours away from
// launch.
func GetLaunches(w http.ResponseWriter, r *http.Request) {
	var launchesURL = "https://fdo.rocketlaunch.live/json/launches"
	var bearer = "Bearer " + rocketLaunchAPIKey

	req, err := http.NewRequest("GET", launchesURL, nil)
	if err != nil {
		fmt.Fprintf(w, "Error building request: %s", err)
		return
	}

	// add authorization header to the req
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Accept", "application/json")

	// send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	// log.Println(string([]byte(body)))
	// w.Write(body)

	type Provider struct {
		ID			int			`json:"id"`
		Name		string		`json:"name"`
		Slug		string		`json:"slug"`
	}

	type Vehicle struct {
		ID			int			`json:"id"`
		Name		string		`json:"name"`
		CompanyID	int			`json:"company_id"`
		Slug		string		`json:"slug"`
	}

	type Location struct {
		ID			int			`json:"id"`
		Name		string		`json:"name"`
		State		string		`json:"state"`
		StateName	string		`json:"statename"`
		Country		string		`json:"country"`
		Slug		string		`json:"slug"`
	}

	type Mission struct {
		ID			int			`json:"id"`
		Name		string		`json:"name"`
		Description	string		`json:"description"`
	}

	type Pad struct {
		ID			int			`json:"id"`
		Name		string		`json:"name"`
		Location	Location	`json:"location"`
	}

	type Weather struct {
		Summary		string		`json:"weather_summary"`
		Temp		float32		`json:"weather_temp"`
		Condition	string		`json:"weather_condition"`
		WindMPH		string		`json:"weather_wind_mph"`
		Icon		string		`json:"weather_icon"`
		Updated		time.Time	`json:"weather_updated"`
	}

	type Result struct {
		ID			int			`json:"id"`
		SortDate	string		`json:"type"`
		Name		string		`json:"name"`
		Provider	Provider	`json:"provider"`
		Vehicle		Vehicle		`json:"vehicle"`
		Pad			Pad			`json:"pad"`
		Missions	[]Mission	`json:"missions"`
		MissionDesc	string		`json:"mission_description"`
		LaunchDesc	string		`json:"launch_description"`
		WinOpen		string		`json:"win_open"`
		WinClose 	string		`json:"win_close"`
		TZero		string		`json:"t0"`
		Weather		Weather
		QuickText	string		`json:"quicktext"`
		SubOrbital	bool		`json:"suborbital"`
	}

	type Response struct {
		ValidAuth	bool		`json:"valid_auth"`
		Results		[]Result	`json:"result"`
	}

	type Launches struct {
		Response	Response	`json:"response"`
	}

	var m Launches

	if err := json.Unmarshal(body, &m); err != nil {
		panic(err)
	}
	for i := 0; i < len(m.Response.Results); i++ {
		winOpen := m.Response.Results[i].WinOpen
		if len(winOpen) > 0 {
			// convert this lame ass date to a good one
			goodDate, _ := datetime.Parse(winOpen, time.UTC)
			timeUntil := time.Until(goodDate)
			// We only care about launch windows in the next 24 hours.
			if timeUntil.Hours() <= 24 {
				 fmt.Fprintf(w, "Hours until launch window: %v", timeUntil.Hours())

				// firebase setup - take the response body and write it to firestore
				// Use the application default credentials.
				conf := &firebase.Config{ProjectID: projectID}

				ctx := context.Background()

				app, err := firebase.NewApp(ctx, conf)
				if err != nil {
					log.Fatalf("firebase.NewApp: %v", err)
				}

				fireClient, err := app.Firestore(ctx)
				if err != nil {
					log.Fatalf("app.Firestore: %v", err)
				}

				defer func(fireClient *firestore.Client) {
					err := fireClient.Close()
					if err != nil {
						log.Fatalf("fireClient.Close: #{err}")
					}
				}(fireClient)

				s1 := strconv.Itoa(m.Response.Results[i].ID)
				myRes, addErr := fireClient.Collection("launches").Doc(s1).Set(ctx, m.Response.Results[i])
				if addErr != nil {
					// Handle any errors in an appropriate way, such as returning them.
					log.Printf("An error has occurred: %v Response: %v", addErr, myRes)
					addErrMsgJSON, addErr := json.Marshal(addErr)
					fmt.Fprint(w, addErr)
					_, err := w.Write(addErrMsgJSON)
					if err != nil {
						return
					}
				}

				resJSON, resErr := json.Marshal(m.Response.Results[i].ID)
				if resErr != nil {
					fmt.Fprintf(w, "Error marshaling JSON: %s", resErr)
					return
				}

				// Return a response to the client, including the ID of the newly created document
				_, err = w.Write(resJSON)
				if err != nil {
					return
				}

			}
		}
	}
}

/************/
/* NWS/NOAA */
/************/

// The code below was ported from
// https://stackoverflow.com/questions/238260/how-to-calculate-the-bounding-box-for-a-given-lat-lng-location

// Semi-axes of WGS-84 geoidal reference
const (
	WGS84_a = 6378137.0 // Major semiaxis [m]
	WGS84_b = 6356752.3 // Minor semiaxis [m]
)

type MapPoint struct {
	Longitude	float64
	Latitude	float64
}

type BoundingBox struct {
	MinPoint 	MapPoint
	MaxPoint 	MapPoint
}

// Deg2rad converts degrees to radians
func Deg2rad(degrees float64) float64 {
	return math.Pi * degrees / 180.0
}

// Rad2deg converts radians to degrees
func Rad2deg(radians float64) float64 {
	return 180.0 * radians / math.Pi
}

func WGS84EarthRadius(lat float64) float64 {
	An := WGS84_a * WGS84_a * math.Cos(lat)
	Bn := WGS84_b * WGS84_b * math.Sin(lat)
	Ad := WGS84_a * math.Cos(lat)
	Bd := WGS84_b * math.Sin(lat)
	return math.Sqrt((An*An + Bn*Bn) / (Ad*Ad + Bd*Bd))
}

// GetBoundingBox takes two arguments, MapPoint is a set of lat/lng,
// 'halfSideInKm' is the half length of the bounding box you want in kilometers.
func GetBoundingBox (point MapPoint, halfSideInKm float64) BoundingBox {
	// Bounding box surrounding the point at given coordinates,
	// assuming local approximation of Earth surface as a sphere
	// of radius given by WGS84
	lat := Deg2rad(point.Latitude)
	lon := Deg2rad(point.Longitude)
	halfSide := 1000 * halfSideInKm

	// Radius of Earth at given latitude
	radius := WGS84EarthRadius(lat)
	// Radius of the parallel at given latitude
	pradius := radius * math.Cos(lat)

	latMin := lat - halfSide / radius
	latMax := lat + halfSide / radius
	lonMin := lon - halfSide / pradius
	lonMax := lon + halfSide / pradius

	return BoundingBox{
		MinPoint: MapPoint{Latitude: Rad2deg(latMin), Longitude: Rad2deg(lonMin)},
		MaxPoint: MapPoint{Latitude: Rad2deg(latMax), Longitude: Rad2deg(lonMax)},
	}
}

// https://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num * output)) / output
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
type Geocode struct {
	UGC		[]string	`json:"UGC"`
}

type Properties struct {
	ID            string		`json:"name"`
	AreaDesc      string		`json:"areaDesc"`
	Geocode       Geocode		`json:"geocode"`
	AffectedZones []string		`json:"affectedZones"`
	Sent          string		`json:"sent"`
	Effective     string		`json:"effective"`
	Onset         string		`json:"onset"`
	Expires       string		`json:"expires"`
	Ends          string		`json:"ends"`
	Status        string		`json:"status"`
	MessageType   string		`json:"messageType"`
	Category      string		`json:"category"`
	Severity      string		`json:"severity"`
	Certainty     string		`json:"certainty"`
	Urgency       string		`json:"urgency"`
	Event         string		`json:"event"`
	Sender        string		`json:"sender"`
	SenderName    string		`json:"senderName"`
	Headline      string		`json:"headline"`
	Description   string		`json:"description"`
	Instruction   string		`json:"instruction"`
	Response      string		`json:"response"`
}

type WeatherAlert struct {
	ID         string		`json:"id"`
	Type       string		`json:"type"`
	Properties Properties	`json:"properties"`
}

type WeatherObject struct {
	Features	[]WeatherAlert	 `json:"features"`
}

// GetWeatherAlerts will get active weather alerts from our categories we're
// interested in and write them to firebase
func GetWeatherAlerts(w http.ResponseWriter, r *http.Request) {
	var weatherAlertsURL = "https://api.weather.gov/alerts/active"

	resp, err := http.Get(weatherAlertsURL)
	if err != nil {
		fmt.Fprintf(w, "Error retrieving documents: %s", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	// w.Write(body)

	var m WeatherObject
	var savedWeathers []WeatherAlert
	if err := json.Unmarshal(body, &m); err != nil {
		panic(err)
	}
	for i := 0; i < len(m.Features); i++ {
		const severity = "Severe"
		sevCheck := strings.Contains(severity, m.Features[i].Properties.Severity)
		// https://www.weather.gov/media/documentation/docs/NWS_Geolocation.pdf
		eventList := [6]string{"Flash Flood Warning", "Severe Thunderstorm Warning", "Tornado Warning", "Snow Squall Warning", "Dust Storm Warning", "Extreme Wind Warning"}
		for _, event := range eventList {
			eventCheck := strings.Contains(event, m.Features[i].Properties.Event)
			if sevCheck {
				if eventCheck {
					// SaveWeatherAlerts(m.Features[i])
					// Push into an array instead and then add the array at the end.. derp
					savedWeathers = append(savedWeathers, m.Features[i])
				}
			}
		}
	}
	SaveWeatherAlerts(savedWeathers)
}

// GetWeatherAlertsAndRasters will get and sort through the latest
// weather alerts from NWS, we're looking for the bad stuff.
// It also attempts to download raster images based on a list of raster image IDs
// we retrieve from the API, and then it is going to save the reference to those IDs
// this mostly works up until the last 1-2 steps and starts to fall apart.
func GetWeatherAlertsAndRasters(w http.ResponseWriter, r *http.Request) {
	var weatherAlertsURL = "https://api.weather.gov/alerts/active"

	resp, err := http.Get(weatherAlertsURL)
	if err != nil {
		fmt.Fprintf(w, "Error retrieving documents: %s", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	// w.Write(body)

	var m WeatherObject
	if err := json.Unmarshal(body, &m); err != nil {
		panic(err)
	}
	for i := 0; i < len(m.Features); i++ {
		//w.Write([]byte(m[i].ID))
		const severity = "Severe"
		sevCheck := strings.Contains(severity, m.Features[i].Properties.Severity)

		// https://www.weather.gov/media/documentation/docs/NWS_Geolocation.pdf
		eventList := [5]string{"Severe Thunderstorm Warning", "Tornado Warning", "Snow Squall Warning", "Dust Storm Warning", "Extreme Wind Warning"}
		// const event = "Severe Thunderstorm Warning"
		for _, event := range eventList {

			eventCheck := strings.Contains(event, m.Features[i].Properties.Event)

			if sevCheck {
				if eventCheck {
					//SaveWeatherAlerts(m.Features[i])
					fmt.Fprintf(w, "Severity: %v Headline: %v\n", m.Features[i].Properties.Severity, m.Features[i].Properties.Event)
					if len(m.Features[i].Properties.AffectedZones) > 0 {
						for a := 0; a < len(m.Features[i].Properties.AffectedZones); a++ {
							// fmt.Fprintf(w, "AffectedZones: %v\n", m.Features[i].Properties.AffectedZones[a])
							aZone := m.Features[i].Properties.AffectedZones[a]

							fmt.Fprintf(w, "AffectedZones: %v\n", aZone)

							resp, err := http.Get(aZone)
							if err != nil {
								fmt.Fprintf(w, "Error retrieving documents: %s", err)
							}
							defer resp.Body.Close()
							countyBody, err := ioutil.ReadAll(resp.Body)

							type AlertData struct {
								ID       string `json:"id"`
								Type     string `json:"type"`
								Geometry struct {
									Type        string          `json:"type"`
									Coordinates [][][][]float64 `json:"coordinates"`
									Geometries  []struct {
										Type        string          `json:"type"`
										Coordinates [][][][]float64 `json:"coordinates"`
									} `json:"geometries"`
								} `json:"geometry"`
							}

							c := AlertData{}
							if err := json.Unmarshal(countyBody, &c); err != nil {
								panic(err)
							}

							var coordsSlice []string

							if len(c.Geometry.Geometries) > 0 {
								for geo := 0; geo < len(c.Geometry.Geometries); geo++ {
									if len(c.Geometry.Geometries[geo].Coordinates) > 0 {
										for z := 0; z < len(c.Geometry.Geometries[geo].Coordinates); z++ {
											lng := c.Geometry.Geometries[geo].Coordinates[0][0][z][0]
											lat := c.Geometry.Geometries[geo].Coordinates[0][0][z][1]
											sLat := fmt.Sprintf("%f", toFixed(lat, 4))
											sLng := fmt.Sprintf("%f", toFixed(lng, 4))
											coordsSlice = append(coordsSlice, sLat)
											coordsSlice = append(coordsSlice, sLng)
											// fmt.Fprintf(w, "%v,%v,", toFixed(lat, 4), toFixed(lng, 4))
											/*
												   boundinbox := GetBoundingBox(MapPoint{Latitude: lat, Longitude: lng},1)
												   fmt.Fprintf(w, "XMin: %v YMin: %v\n",
													   // X = longitude, Y = latitude
													   toFixed(boundinbox.MinPoint.Longitude, 4),
													   toFixed(boundinbox.MinPoint.Latitude, 4))
												   fmt.Fprintf(w, "XMax: %v YMax: %v\n",
													   toFixed(boundinbox.MaxPoint.Longitude, 4),
													   toFixed(boundinbox.MaxPoint.Latitude, 4))
											*/
										}
									}
								}
							} else {
								if len(c.Geometry.Coordinates) > 0 {
									for f := 0; f < len(c.Geometry.Coordinates[0][0]); f++ {
										lng := c.Geometry.Coordinates[0][0][f][0]
										lat := c.Geometry.Coordinates[0][0][f][1]
										// fmt.Fprintf(w, "%v,%v,", toFixed(lat, 4), toFixed(lng, 4))
										sLat := fmt.Sprintf("%f", toFixed(lat, 4))
										sLng := fmt.Sprintf("%f", toFixed(lng, 4))
										coordsSlice = append(coordsSlice, sLat)
										coordsSlice = append(coordsSlice, sLng)

										/*
											   boundinbox := GetBoundingBox(MapPoint{Latitude: lat, Longitude: lng},1)
											   fmt.Fprintf(w, "XMin: %v YMin: %v\n",
												   toFixed(boundinbox.MinPoint.Longitude, 4),
												   toFixed(boundinbox.MinPoint.Latitude, 4))
											   fmt.Fprintf(w, "XMax: %v YMax:%v\n",
												   toFixed(boundinbox.MaxPoint.Longitude, 4),
												   toFixed(boundinbox.MaxPoint.Latitude, 4))

										*/
									}
								}
							}

							// start processing coordinates into raster catalog requests
							result := strings.Join(coordsSlice, ",")

							// API can only handle so many requests in a URI
							limit := 1000
							// var imageIds []int

							if len(result) <= limit {
								imageQuery := "https://idpgis.ncep.noaa.gov/arcgis/rest/services/radar/radar_base_reflectivity_time/ImageServer/query" + "?geometry=" + result + "&f=json"
								// fmt.Fprintf(w, "ImageQuery: %v", imageQuery)
								imageIds := GetImageIDs(imageQuery)
								catalogItems := GetCatalogItems(imageIds)
								// fmt.Fprintf(w, "CatalogItemsCount: %v", len(catalogItems))
								for c := 0; c < len(catalogItems); c++ {
									for ring := 0; ring < len(catalogItems[c].Geometry.Rings[0][0]); ring++ {
										var rings []float64 = catalogItems[c].Geometry.Rings[0][ring]
										rasterImage := GetRasterImage(catalogItems[c].Attributes.Objectid, rings)
										fmt.Fprintf(w, "rings len: %v\n", rasterImage.Href)
										fmt.Fprintf(w, "Rasterimage: %v\n", rasterImage.Href)
										SaveRasterData(catalogItems[c].Attributes.Objectid, rings, rasterImage)
									}
								}
							} else {
								for cnt := 0; cnt < len(result); cnt += limit {
									batch := result[cnt:min(cnt+limit, len(result))]
									imageQuery := "https://idpgis.ncep.noaa.gov/arcgis/rest/services/radar/radar_base_reflectivity_time/ImageServer/query" + "?geometry=" + batch + "&f=json"
									// fmt.Fprintf(w, "ImageQuery2: %v", imageQuery)
									imageIds := GetImageIDs(imageQuery)
									catalogItems := GetCatalogItems(imageIds)
									for c := 0; c < len(catalogItems); c++ {
										for ring := 0; ring < len(catalogItems[c].Geometry.Rings[0][0]); ring++ {
											var rings []float64 = catalogItems[c].Geometry.Rings[0][ring]
											rasterImage := GetRasterImage(catalogItems[c].Attributes.Objectid, rings)
											fmt.Fprintf(w, "Rasterimage2: %v\n", rasterImage.Href)
											SaveRasterData(catalogItems[c].Attributes.Objectid, rings, rasterImage)
										}
									}
								}
							}
						}
					}
					/*
						 if eventCheck {
							 fmt.Fprintf(w, "Severity: %v Headline: %v\n",
												 m.Features[i].Properties.Severity,
												 m.Features[i].Properties.Event)
							 for x := 0; x < len(m.Features[i].Properties.Geocode.UGC); x++ {
								 fmt.Fprintf(w, "\tGeo: %v\n", m.Features[i].Properties.Geocode.UGC[x])
							 }
						 }
					*/
				}
			}
		}
	}
}

type CatalogItem struct {
	Attributes struct {
		Objectid           int         `json:"objectid"`
		Name               string      `json:"name"`
		Category           int         `json:"category"`
		IdpSource          interface{} `json:"idp_source"`
		IdpSubset          string      `json:"idp_subset"`
		IdpFiledate        int64       `json:"idp_filedate"`
		IdpIngestdate      int64       `json:"idp_ingestdate"`
		IdpCurrentForecast int         `json:"idp_current_forecast"`
		IdpTimeSeries      int         `json:"idp_time_series"`
		IdpIssueddate      int64       `json:"idp_issueddate"`
		IdpValidtime       int64       `json:"idp_validtime"`
		IdpValidendtime    int64       `json:"idp_validendtime"`
		ShapeLength        float64     `json:"shape_Length"`
		ShapeArea          float64     `json:"shape_Area"`
	} `json:"attributes"`
	Geometry struct {
		Rings            [][][]float64 `json:"rings"`
		SpatialReference struct {
			Wkid       int `json:"wkid"`
			LatestWkid int `json:"latestWkid"`
		} `json:"spatialReference"`
	} `json:"geometry"`
}

func GetCatalogItems(imageIds []int) []CatalogItem{
	var CatalogItems []CatalogItem
	for i := 0; i < len(imageIds); i++ {
		imageId := strconv.Itoa(imageIds[i])
		url := "https://idpgis.ncep.noaa.gov/arcgis/rest/services/radar/radar_base_reflectivity_time/ImageServer/" + imageId + "?f=pjson"
		resp, err := http.Get(url)
		if err != nil {
			log.Println(resp)
			panic(err)
		}
		defer resp.Body.Close()
		respBody, respBodyErr := ioutil.ReadAll(resp.Body)
		if respBodyErr != nil {
			log.Println(respBody)
			panic(respBodyErr)
		}

		catalogItem := CatalogItem{}
		if err := json.Unmarshal(respBody, &catalogItem); err != nil {
			panic(err)
		}
		CatalogItems = append(CatalogItems, catalogItem)
	}
	return CatalogItems
}

// SaveWeatherAlerts accepts one argument, a WeatherAlert object
// retrieved from NWS weather alerts REST endpoint.
func SaveWeatherAlerts(alerts []WeatherAlert) {
	// firebase setup - take the response body and write it to firestore
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID, DatabaseURL: "https://chaseapp-8459b.firebaseio.com" }

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}
	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing the database client: ", err)
	}

	// Get a database reference to our radar raster image store
	ref := client.NewRef("alerts")
	alertRef := ref.Child("nws")
	alertRefErr := alertRef.Set(ctx, alerts)
	if alertRefErr != nil {
		log.Fatalln("Error saving weather alerts in google realtime db: ", err)
	}
}

// SaveRasterData takes two arguments and will save save our raster data to firebase
// realtime database
func SaveRasterData(objectId int, area []float64, rasterImage RasterImage) {
	// firebase setup - take the response body and write it to firestore
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID, DatabaseURL: "https://chaseapp-8459b.firebaseio.com" }

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}
	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing the database client: ", err)
	}

	// Get a database reference to our radar raster image store
	ref := client.NewRef("radar")
	radarRef := ref.Child(strconv.Itoa(objectId))
	radarRefErr := radarRef.Update(ctx, map[string]interface{}{ "ObjectId": objectId, "RasterImage": rasterImage, "Area": area} )
	if radarRefErr != nil {
		log.Fatalln("Error saving Raster Data in google realtime db: ", err)
	}
}

type RasterImage struct {
	Href   string `json:"href"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Extent struct {
		Xmin             float64 `json:"xmin"`
		Ymin             float64 `json:"ymin"`
		Xmax             float64 `json:"xmax"`
		Ymax             float64 `json:"ymax"`
		SpatialReference struct {
			Wkid       int `json:"wkid"`
			LatestWkid int `json:"latestWkid"`
		} `json:"spatialReference"`
	} `json:"extent"`
	Scale int `json:"scale"`
}

func GetRasterImage(objectId int, bbox []float64) RasterImage {
	sObjectId := strconv.Itoa(objectId)
	// sBoundingBox := strings.Join(bbox, ",")

	var BBox []string
	for b := 0; b < len(bbox); b++ {
		log.Println(bbox[b])
		sB := fmt.Sprintf("%f", bbox[b])
		BBox = append(BBox, sB)
	}

	sBBox := strings.Join(BBox, ",")

	url := "https://idpgis.ncep.noaa.gov/arcgis/rest/services/radar/radar_base_reflectivity_time/ImageServer/" + sObjectId + "/image?bbox=" + sBBox + "&f=pjson"
	log.Printf("URL: %v\n", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBody, respErr := ioutil.ReadAll(resp.Body)

	if respErr != nil {
		panic(respErr)
	}

	rasterImage := RasterImage{}
	if err := json.Unmarshal(respBody, &rasterImage); err != nil {
		panic(err)
	}

	return rasterImage
}

func GetImageIDs(imageQuery string) []int {
	imageResp, imageErr := http.Get(imageQuery)
	if imageErr != nil {
		panic(imageErr)
	}
	defer imageResp.Body.Close()
	imageRespBody, imageRespErr := ioutil.ReadAll(imageResp.Body)
	if imageRespErr != nil {
		panic(imageRespErr)
	}
	type ImageQuery struct {
		ObjectIDFieldName string `json:"objectIdFieldName"`
		Fields            []struct {
			Name   string      `json:"name"`
			Type   string      `json:"type"`
			Alias  string      `json:"alias"`
			Domain interface{} `json:"domain"`
		} `json:"fields"`
		GeometryType     string `json:"geometryType"`
		SpatialReference struct {
			Wkid       int `json:"wkid"`
			LatestWkid int `json:"latestWkid"`
		} `json:"spatialReference"`
		Features []struct {
			Attributes struct {
				Objectid    int     `json:"objectid"`
				ShapeLength float64 `json:"shape_Length"`
				ShapeArea   float64 `json:"shape_Area"`
			} `json:"attributes"`
			Geometry struct {
				Rings            [][][]float64 `json:"rings"`
				SpatialReference struct {
					Wkid       int `json:"wkid"`
					LatestWkid int `json:"latestWkid"`
				} `json:"spatialReference"`
			} `json:"geometry"`
		} `json:"features"`
	}

	imgQuery := ImageQuery{}
	if err := json.Unmarshal(imageRespBody, &imgQuery); err != nil {
		panic(err)
	}

	var objectIdArray []int
	// var weatherstationsArray []WeatherStation
	for feature := 0; feature < len(imgQuery.Features); feature++ {
		objectId := imgQuery.Features[feature].Attributes.Objectid
		objectIdArray = append(objectIdArray, objectId)
	}

	return objectIdArray
}

// UpdateWeather is responsible for keeping the weather stations and
// severe weather alert tracking system up to date. Should be run out of
// google cron system every 5 minutes
func UpdateWeather(w http.ResponseWriter, r *http.Request) {
	var ListStationsURL = "https://us-central1-chaseapp-8459b.cloudfunctions.net/ListStations"
	resp, err := http.Get(ListStationsURL)
	if err != nil {
		fmt.Fprintf(w, "Error retrieving documents: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	type WeatherStation struct {
		ID			string
		CreatedAt 	string
	}

	var m []WeatherStation
	if err := json.Unmarshal(body, &m); err != nil {
		panic(err)
	}
	for i := 0; i < len(m); i++ {
		// call the UpdateStation API endpoint for each station ID
		postBody, _ := json.Marshal(map[string]string{
			"ID": m[i].ID,
		})
		responseBody := bytes.NewBuffer(postBody)
		resp, err = http.Post("https://us-central1-chaseapp-8459b.cloudfunctions.net/GetWeather",
			"application/json", responseBody)
		if err != nil {
			log.Fatalf("An error occured %v", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		sb := string(body)
		fmt.Fprintln(w, sb)
	}
}

// GetWeather takes the ID of a weather station and returns the data from the
// NOAA API
func GetWeather(w http.ResponseWriter, r *http.Request) {
	// Inputs
	var d struct {
		ID			string	    `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error, cannot decode body of request")
		return
	}

	if d.ID == "" {
		fmt.Fprint(w, "Must supply weather station ID!")
		return
	}

	var weatherUrl = "https://api.tidesandcurrents.noaa.gov/mdapi/prod/webapi/stations/" + d.ID + ".json"
	resp, err := http.Get(weatherUrl)
	if err != nil {
		fmt.Fprintf(w, "Error retrieving documents: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	m := map[string]interface{}{}
	if err := json.Unmarshal(body, &m); err != nil {
		panic(err)
	}
	fmt.Printf("%q",m)

	// firebase setup - take the response body and write it to firestore
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("client.Close: #{err}")
		}
	}(client)

	// id, err := uuid.NewUUID()

	myRes, addErr := client.Collection("weather").Doc(d.ID).Set(ctx, m)
	if addErr != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %v Response: %v", addErr, myRes)
		addErrMsgJSON, addErr := json.Marshal(addErr)
		fmt.Fprint(w, addErr)
		_, err := w.Write(addErrMsgJSON)
		if err != nil {
			return
		}
	}

	resJSON, resErr := json.Marshal(d.ID)
	if resErr != nil {
		fmt.Fprintf(w, "Error marshaling JSON: %s", resErr)
		return
	}

	// Return a response to the client, including the ID of the newly created document
	_, err = w.Write(resJSON)
	if err != nil {
		return
	}

}

/*********************************/
/* WEATHER STATIONS - STORMSURGE */
/*********************************/

// ListStations lists weather station IDs we are tracking in firebase
func ListStations(w http.ResponseWriter, r *http.Request) {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("client.Close: #{err}")
		}
	}(client)

	type WeatherStation struct {
		ID	        string    	`firestore:"stationId"`
		CreatedAt   time.Time 	`firestore:"createdAt"`
	}

	weatherstations := client.Collection("weatherstations")
	docs, err := weatherstations.Documents(ctx).GetAll()
	if err != nil {
		fmt.Fprint(w, "Error retrieving documents")
		return
	}

	var weatherstationsArray []WeatherStation

	for _, doc := range docs {
		var weatherstationData WeatherStation
		if err := doc.DataTo(&weatherstationData); err != nil {
			fmt.Fprintf(w, "Error retrieving documents: %s", err)
			return
		}
		weatherstationsArray = append(weatherstationsArray, weatherstationData)
	}

	// Convert our array to JSON and spit it out
	js, err := json.Marshal(weatherstationsArray)
	if err != nil {
		fmt.Fprintf(w, "Error marshaling JSON: %s", err)
		return
	}
	w.Write(js)
}

type Boats struct {
	MMSI	string
	Group   string
	Type    string
}

// DeleteBoat deletes a Boat from firebase firestore
func DeleteBoat(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, r.Body)
	conf := &firebase.Config{ProjectID: projectID}

	var APIKEY = "8b373c2d-41bc-4a18-80f3-b3671f04930f"
	// fmt.Fprint(w, r)
	// fmt.Fprint(w, r.Header.Get("X-ApiKey"))

	if APIKEY != r.Header.Get("X-ApiKey") {
		fmt.Fprint(w, "Invalid or missing X-ApiKey")
		return
	}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer client.Close()

	var d struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error, body must contain ID!")
		return
	}

	fmt.Fprint(w, d)

	if d.ID == "" {
		fmt.Fprint(w, "Must supply ID!")
		return
	}

	fmt.Fprint(w, d)

	_, addErr := client.Collection("boats").Doc(d.ID).Delete(ctx)
	if addErr != nil {
		_, addErr := json.Marshal(addErr)
		fmt.Fprintf(w, "Error deleting chase record: %v", addErr)
	}
}

func UpdateBoat(w http.ResponseWriter, r *http.Request) {
	// Inputs
	var d struct {
		ID			string		`json:"id"`
		MMSI		string		`json:"mmsi"`
		Group		string		`json:"desc"`
		Type		string      `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error, body must contain id, MMSI, group, and type!")
		return
	}

	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer client.Close()

	if d.MMSI != "" {
		_, addErr := client.Collection("boats").Doc(d.ID).Set(ctx, map[string]interface{}{"mmsi": d.MMSI}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating boat MMSI: %v", addErr)
		}
	}
	if d.Group != "" {
		_, addErr := client.Collection("boats").Doc(d.ID).Set(ctx, map[string]interface{}{"group": d.Group}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating boat group: %v", addErr)
		}
	}
	if d.Type != "" {
		_, addErr := client.Collection("boats").Doc(d.ID).Set(ctx, map[string]interface{}{"type": d.Type}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating boat type: %v", addErr)
		}
	}
}

/****************/
/* GetStream.io */
/****************/

// go get github.com/GetStream/stream-chat-go/v5

// GetStreamToken is a web endpoint that takes a user_id and returns a getstream token
// using the getstream.io API.
func GetStreamToken(w http.ResponseWriter, r *http.Request) {

	// pre flight stuff
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var d struct {
		UserID	 string		`json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprintf(w, "Error, body must contain user_id to create a token for")
		return
	}

	// instantiate your stream client using the API key and secret
	// the secret is only used server side and gives you full access to the API
	ServerClient, _ := stream.NewClient(getStream_API_KEY, getStream_API_SECRET)
	token, _ := ServerClient.CreateToken(d.UserID, time.Time{})
	// next, hand this token to the client in your in your login or registration response

	// using a custom writer from github.com/mfreeman451/golang/common/writers
	jw := writers.NewMessageWriter(token)
	jsonString, err := jw.JSONString()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Println(err.Error())
		return
	}

	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonString))
}

// AddModerator will take a userid and add them to the moderators for a given channel
func AddModerator(w http.ResponseWriter, r *http.Request) {

	var d struct {
		UserID	 string		`json:"user_id"`
		Channel	 string		`json:"channel"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprintf(w, "Error, body must contain user_id, channel")
		return
	}

	// the secret is only used server side and gives you full access to the API
	client, _ := stream.NewClient(getStream_API_KEY, getStream_API_SECRET)
	channel := client.Channel("livestream", d.Channel)
	_, err := channel.AddModerators(context.Background(), d.UserID)
	if (err != nil) {
		fmt.Fprintf(w, "Error assigning a moderator to a channel %v", err)
	}
	fmt.Fprintf(w, "Added %v to moderators group on channel", d.UserID)
}

/**************/
/* FCM Stuff */
/*************/

type Subscription struct {
	Token	[]string
	Topic	string
}

// AddToken web version of addToken
func AddToken(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Token	string		`json:"token"`
		Browser	string		`json:"browser,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprintf(w, "Error, body must contain a token to add")
		return
	}

	if d.Token == "" {
		fmt.Fprintf(w, "Error, must supply a Token in the request body")
		return
	}

	err := addToken(d.Token,d.Browser)
	if err != nil {
		fmt.Fprintf(w, "Couldn't add token: %v", err)
		log.Fatalf("Couldn't add token: %v", err)
	}

	fmt.Fprint(w, "Added token")
}

// addToken adds the token and createdAt field in the 'tokens' collection
func addToken(token,browser string) interface{}{
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer client.Close()

	type TokenRecord struct {
		Token		string
		Browser		string
		CreatedAt	time.Time
	}

	myRes, addErr := client.Collection("token").Doc(token).Set(ctx, TokenRecord{Token: token, Browser: browser, CreatedAt: time.Now()})
	if addErr != nil {
		log.Printf("An error has occurred: %v Response: %v", addErr, myRes)
		return addErr
	}
	return nil
}

func CheckToken(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Token	string		`json:"token"`
	}

	fmt.Println(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprintf(w, "Error, body must contain a token to check age of")
		return
	}

	if d.Token == "" {
		fmt.Fprintf(w, "Error, must supply a Token in the request body")
		return
	}

	days := checkToken(d.Token)
	fmt.Fprintln(w, "%v", days)
}
// checkToken takes a token and checks to see if it is older than a week old
func checkToken(token string) (days int) {
	// TODO: test CheckToken
	// TODO: write unit test for CheckToken
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer client.Close()

	log.Println("Token: ", token)
	dsnap, err := client.Collection("tokens").Doc(token).Get(ctx)
	if err != nil {
		// log.Fatalln("Error getting tokens from firebase", err)
		return -1
	}

	type TokenRecord struct {
		Token		string
		Browser		string
		CreatedAt	time.Time
	}

	var dbToken TokenRecord
	if err := dsnap.DataTo(&dbToken); err != nil {
		log.Println("Token: ", dbToken)
		// calculate years, month, days and time between dates
		year, month, day, hour, min, sec := diffDate(dbToken.CreatedAt, time.Now())

		fmt.Printf("difference %d years, %d months, %d days, %d hours, %d mins and %d seconds.", year, month, day, hour, min, sec)
		fmt.Printf("")

		// calculate total number of days
		duration := time.Now().Sub(dbToken.CreatedAt)
		// fmt.Printf("difference %d days", int(duration.Hours()/24) )
		days := int(duration.Hours()/24)
		return days
	}
	return -1
}

// diffDate takes two time.Time objects and finds the difference between them
func diffDate(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}
// UpdateToken updates the created_at field in the 'tokens' collection
func UpdateToken(token string) interface{}{
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer client.Close()

	type TokenRecord struct {
		Token     string
		CreatedAt time.Time
	}

	myRes, addErr := client.Collection("token").Doc(token).Set(ctx, TokenRecord{Token: token, CreatedAt: time.Now()})
	if addErr != nil {
		log.Printf("An error has occurred: %v Response: %v", addErr, myRes)
		return addErr
	}
	return myRes
}

// SubscribeTopic calls subscribeTopic and takes a token and topic from an HTTP request body
func SubscribeTopic(w http.ResponseWriter, r *http.Request) {

	// pre flight stuff
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var d struct {
		Token	string		`json:"token"`
		Topic	string		`json:"topic"`
	}

	fmt.Println(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprintf(w, "Error, body must contain a token and topic to subscribe to")
		return
	}

	if d.Token == "" {
		fmt.Fprintf(w, "Error, must supply a Token in the request body")
		return
	}
	if d.Topic == "" {
		fmt.Fprintf(w, "Error, must supply a Topic in the request body")
		return
	}

	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		 log.Fatalf("firebase.NewApp: %v", err)
	}

	// FCM messaging stuffs:
	fcm, fcmErr := app.Messaging(ctx)
	if fcmErr != nil {
		 log.Fatalf("Problem setting up app.Messaging: %v", fcmErr)
	}

	var tokens []string
	tokens = append(tokens, d.Token)
	_, subErr := fcm.SubscribeToTopic(ctx, tokens, d.Topic)
	if subErr != nil {
		 fmt.Fprintf(w, "Can't subscribe token %v to topic %v", d.Token, d.Topic )
		 log.Fatalf("Can't subscribe to topic: %v", subErr)
	}

	jw := writers.NewMessageWriter(d.Topic)
	jsonString, err := jw.JSONString()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Println(err.Error())
		return
	}

	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonString))
}

// UnsubscribeTopic unsubscribes a given token from an FCM messaging topic that you supply
func UnsubscribeTopic(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Token string `json:"token"`
		Topic string `json:"topic"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprintf(w, "Error, body must contain a token and topic to unsubscribe to")
		return
	}

	if d.Token == "" {
		fmt.Fprintf(w, "Error, must supply a Token in the request body")
		return
	}
	if d.Topic == "" {
		fmt.Fprintf(w, "Error, must supply a Topic in the request body")
		return
	}

	var sub = &Subscription{
		Token: []string{d.Token},
		Topic: d.Topic,
	}

	status := unsubscribeTopic(*sub)

	if status != true {
		fmt.Fprintf(w, "Problem unsubscribing to topic: %v", status)
	}
}
// unsubscribeTopic unsubscribes a given token from a given FCM messaging topic, internal function
func unsubscribeTopic(subscription Subscription) bool {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	// FCM messaging stuffs:
	fcm, fcmErr := app.Messaging(ctx)

	if fcmErr != nil {
		return false
	}
	_, err = fcm.UnsubscribeFromTopic(ctx, subscription.Token, subscription.Topic)
	return true
}

// AddUser takes an HTTP request body and parses it for
// parameters that it will then write to firestore in the 'users' collection
// 'UID' is a required field.
func AddUser(w http.ResponseWriter, r *http.Request) {
	var d struct {
		UID			string		`json:"uid"`
		CreatedAt	time.Time	`json:"created_at"`
		Profile	struct {
			Notifications struct {
				Web		bool	`json:"web"`
			} `json:"notifications"`
		} `json:"profile"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprintf(w, "Error, body must contain all the things")
		return
	}

	if d.UID == "" {
		fmt.Fprintf(w, "Error, must supply a UID in the request body")
		return
	}
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer client.Close()

	d.CreatedAt = time.Now()
	myRes, addErr := client.Collection("users").Doc(d.UID).Set(ctx, d)
	if addErr != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %v Response: %v", addErr, myRes)
		addErrMsgJSON, addErr := json.Marshal(addErr)
		fmt.Fprint(w, addErr)
		_, err := w.Write(addErrMsgJSON)
		if err != nil {
			return
		}
	}

	resJSON, resErr := json.Marshal(d.UID)
	if resErr != nil {
		fmt.Fprintf(w, "Error marshaling JSON: %s", resErr)
		return
	}

	// Return a response to the client, including the ID of the newly created document
	_, err = w.Write(resJSON)
	if err != nil {
		return
	}
}

/*******/
/* AIS */
/*******/

func AddBoat(w http.ResponseWriter, r *http.Request) {
	// Inputs
	var d struct {
		MMSI		string		`json:"mmsi"`
		Group		string		`json:"group"`
		Type		string      `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error, body must contain mmsi, group, type!")
		return
	}

	if d.MMSI == "" {
		fmt.Fprint(w, "Must supply MMSI!")
		return
	}

	if d.Group == "" {
		fmt.Fprint(w, "Must supply group for Boat (navy, uscg, spacex, ..) !")
		return
	}

	if d.Type == "" {
		fmt.Fprint(w, "Must supply type for Boat (ship, sar, patrol) !")
		return
	}

	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer client.Close()

	boats := client.Collection("boats")
	docs, err := boats.Documents(ctx).GetAll()
	if err != nil {
		fmt.Fprint(w, "Error retrieving documents")
		return
	}

	var boatsArray []Boats

	for _, doc := range docs {
		var boatData Boats
		if err := doc.DataTo(&boatData); err != nil {
			fmt.Fprintf(w, "Error retrieving documents: %s", err)
			return
		}
		boatsArray = append(boatsArray, boatData)
	}

	// Convert our array to JSON and spit it out
	js, err := json.Marshal(boatsArray)
	if err != nil {
		fmt.Fprintf(w, "Error marshaling JSON: %s", err)
		return
	}
	w.Write(js)
}

// ListBoats prints the JSON encoded MMSI, group, and type in the body
// of the request or an error message if there isn't one.
func ListBoats(w http.ResponseWriter, r *http.Request) {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer client.Close()

	airships := client.Collection("boats")
	docs, err := airships.Documents(ctx).GetAll()
	if err != nil {
		fmt.Fprint(w, "Error retrieving documents")
		return
	}

	var boatsArray []Boats

	for _, doc := range docs {
		var boatData Boats
		if err := doc.DataTo(&boatData); err != nil {
			fmt.Fprintf(w, "Error retrieving documents: %s", err)
			return
		}
		boatsArray = append(boatsArray, boatData)
	}

	// Convert our array to JSON and spit it out
	js, err := json.Marshal(boatsArray)
	if err != nil {
		fmt.Fprintf(w, "Error marshaling JSON: %s", err)
		return
	}
	w.Write(js)
}

/********/
/* ADSB */
/********/

type Airships struct {
	Tailno  string
	Group   string
	Type    string
}

// DeleteAirship deletes an Airship from firebase firestore
func DeleteAirship(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, r.Body)
	conf := &firebase.Config{ProjectID: projectID}

	var APIKEY = "8b373c2d-41bc-4a18-80f3-b3671f04930f"
	// fmt.Fprint(w, r)
	// fmt.Fprint(w, r.Header.Get("X-ApiKey"))

	if APIKEY != r.Header.Get("X-ApiKey") {
		fmt.Fprint(w, "Invalid or missing X-ApiKey")
		return
	}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer client.Close()

	var d struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error, body must contain ID!")
		return
	}

	fmt.Fprint(w, d)

	if d.ID == "" {
		fmt.Fprint(w, "Must supply ID!")
		return
	}

	fmt.Fprint(w, d)

	_, addErr := client.Collection("airships").Doc(d.ID).Delete(ctx)
	if addErr != nil {
		_, addErr := json.Marshal(addErr)
		fmt.Fprintf(w, "Error deleting chase record: %v", addErr)
	}
}

func UpdateAirship(w http.ResponseWriter, r *http.Request) {
	// Inputs
	var d struct {
		ID			string		`json:"id"`
		Tailno		string		`json:"tailno"`
		Group		string		`json:"group"`
		ImageURL    string      `json:"imageUrl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error, body must contain id, tailno, group, imageUrl!")
		return
	}

	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer client.Close()

	if d.Tailno != "" {
		_, addErr := client.Collection("airships").Doc(d.ID).Set(ctx, map[string]interface{}{"tailno": d.Tailno}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating airship tailno: %v", addErr)
		}
	}
	if d.Group != "" {
		_, addErr := client.Collection("airships").Doc(d.ID).Set(ctx, map[string]interface{}{"group": d.Group}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating airship group: %v", addErr)
		}
	}
	if d.ImageURL != "" {
		_, addErr := client.Collection("airships").Doc(d.ID).Set(ctx, map[string]interface{}{"imageUrl": d.Tailno}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating airship imageUrl: %v", addErr)
		}
	}
}

// AddAirship takes tailno, group, and imageUrl (strings) in an HTTP request body
func AddAirship(w http.ResponseWriter, r *http.Request) {
	// Inputs
	var d struct {
		Tailno		string		`json:"tailno"`
		Group		string		`json:"group"`
		ImageURL    string      `json:"imageUrl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error, body must contain tailno, group, imageUrl!")
		return
	}

	if d.Tailno == "" {
		fmt.Fprint(w, "Must supply Tailno!")
		return
	}

	if d.Group == "" {
		fmt.Fprint(w, "Must supply group for Airship (leo, media, ..) !")
		return
	}
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer client.Close()

	// TODO: FIX - BROKEN
	airships := client.Collection("airships")
	docs, err := airships.Documents(ctx).GetAll()
	if err != nil {
		fmt.Fprint(w, "Error retrieving documents")
		return
	}

	var airshipsArray []Airships

	for _, doc := range docs {
		var airshipData Airships
		if err := doc.DataTo(&airshipData); err != nil {
			fmt.Fprintf(w, "Error retrieving documents: %s", err)
			return
		}
		airshipsArray = append(airshipsArray, airshipData)
	}

	// Convert our array to JSON and spit it out
	js, err := json.Marshal(airshipsArray)
	if err != nil {
		fmt.Fprintf(w, "Error marshaling JSON: %s", err)
		return
	}
	w.Write(js)
}

// ListAirships prints the JSON encoded tailno, group, and type in the body
// of the request or an error message if there isn't one.
func ListAirships(w http.ResponseWriter, r *http.Request) {
    // Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer client.Close()

	airships := client.Collection("airships")
    docs, err := airships.Documents(ctx).GetAll()
    if err != nil {
        fmt.Fprint(w, "Error retrieving documents")
        return
    }

    var airshipsArray []Airships

    for _, doc := range docs {
        var airshipData Airships
        if err := doc.DataTo(&airshipData); err != nil {
            fmt.Fprintf(w, "Error retrieving documents: %s", err)
            return
        }
        airshipsArray = append(airshipsArray, airshipData)
    }

    // Convert our array to JSON and spit it out
    js, err := json.Marshal(airshipsArray)
    if err != nil {
        fmt.Fprintf(w, "Error marshaling JSON: %s", err)
        return
    }
    w.Write(js)
}

/**********/
/* CHASES */
/**********/

type Networks struct {
	Name		string
	URL         string
	Tier		int
	Logo		string
	Other		string
}

type Wheels struct {
	W1  string
	W2  string
	W3  string
	W4  string
}

type Sentiment struct {
	Magnitude	float64		`firestore:"magnitude"`
	Score		float64		`firestore:"score"`
}

type Chase struct {
	ID          string    	""
	Name        string    	`firestore:"Name"`
	Desc        string    	`firestore:"Desc"`
	Live        bool      	`firestore:"Live"`
	Networks	[]Networks	`firestore:"Networks"`
	Wheels      Wheels      `firestore:"Wheels"`
	Votes       int       	`firestore:"Votes"`
	CreatedAt   time.Time 	`firestore:"CreatedAt"`
	EndedAt     time.Time 	`firestore:"EndedAt"`
	ImageURL    string      `firestore:"ImageURL"`
	Reddit		string		`firestore:"Reddit"`
	Sentiment	Sentiment
}

var NotifyRequest struct {
	Name		string		`json:"name"`
	Desc        string		`json:"desc"`
	ImageURL    string      `json:"imageURL"`
	URL 		string      `json:"url"`
}

var AddChaseInput struct {
	ID			string		`json:"id,omitempty"`
	Name        string		`json:"name"`
	Desc        string		`json:"desc"`
	Live		bool		`json:"live"`
	Networks	[]Networks  `json:"networks"`
	Wheels      Wheels      `json:"wheels"`
	Votes       int	    	`json:"votes"`
	CreatedAt   time.Time 	`json:"createdAt"`
	ImageURL    string      `json:"imageURL"`
	Reddit		string		`json:"reddit"`
}

var ChaseInput struct {
	ID			string		`json:"id,omitempty"`
	Name        string		`json:"name"`
	Desc        string		`json:"desc"`
	Live		bool		`json:"live"`
	Networks	[]Networks  `json:"networks"`
	Wheels      Wheels      `json:"wheels"`
	Votes       int	    	`json:"votes"`
	CreatedAt   time.Time 	`json:"createdAt"`
	EndedAt		time.Time	`json:"endedAt,omitempty"`
	ImageURL    string      `json:"imageURL"`
	Reddit		string		`json:"reddit"`
}

type PushTokens struct {
	Token		string		`json:"token"`
	CreatedAt	time.Time	`json:"created_at"`
	TokenType	string		`json:"type"`
}

type User struct {
	UID 		string			`firestore:"uid"`
	LastUpdated	time.Time		`firestore:"lastupdated"`
	PhotoURL	string			`firestore:"photourl"`
	UserName	string			`firestore:"username"`
	Tokens		[]PushTokens	`firestore:"tokens"`
}

var UserInput struct {
	UID			string			`json:"uid"`
	LastUpdated	time.Time		`json:"lastupdated"`
	PhotoURL	string			`json:"photourl,omitempty"`
	UserName	string			`json:"username"`
	Tokens		[]PushTokens	`json:"tokens"`
}

/* CHASEAPP-CRUD */

// ListImages takes an input and returns a list of images for that chaseID
func ListImages(w http.ResponseWriter, r *http.Request) {
	// Inputs
	var d struct {
		ID		string    `json:"ID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error, body must ID!")
		return
	}

	if d.ID == "" {
		fmt.Fprint(w, "Must supply ID!")
		return
	}

	bucket := "chaseapp-8459b.appspot.com"

	ctx := context.Background()
	storageclient, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Errorf("storage.NewClient: %v", err)
	}
	defer storageclient.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	prefix := "chases/" + d.ID + "/"
	query := &storage.Query{Prefix: prefix}
	it := storageclient.Bucket(bucket).Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err != nil {
			fmt.Errorf("it.Next error: %v", err)
			break
		}
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Errorf("Bucket(%q).Objects: %v", bucket, err)
		}

		// fmt.Fprintf(w, "%v %v\n", attrs.Name, attrs.MediaLink)
		fmt.Fprintf(w, "%v\n", attrs.Name)
	}
}

// UploadImage takes MultiPart form inputs and is used to upload
// pictures to firestore storage
func UploadImage(w http.ResponseWriter, r *http.Request) {
	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()
	storageclient, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Errorf("storage.NewClient: %v", err)
		return
	}
	defer storageclient.Close()

	// file upload stuff
	const maxMemory = 2 * 1024 * 1024 // 2 megabytes.

	// ParseMultipartForm parses a request body as multipart/form-data.
	// The whole request body is parsed and up to a total of maxMemory bytes of
	// its file parts are stored in memory, with the remainder stored on
	// disk in temporary files.
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		log.Printf("Error parsing form: %v", err)
		return
	}

	if keyvalue := r.FormValue("ID"); keyvalue == "" {
		http.Error(w, "Must supply ID in form", http.StatusBadRequest)
		log.Printf("Missing ID field from form: %v", err)
		return
	}

	// Remove all temporary files after function is finished.
	defer func() {
		if err := r.MultipartForm.RemoveAll(); err != nil {
			http.Error(w, "Error cleaning up form files", http.StatusInternalServerError)
			log.Printf("Error cleaning up form files: %v", err)
		}
	}()

	// r.MultipartForm.File contains *multipart.FileHeader values for every
	// file in the form. You can access the file contents using
	// *multipart.FileHeader's Open method.
	for _, headers := range r.MultipartForm.File {
		if headers != nil {
			for _, h := range headers {
				// fmt.Fprintf(w, "File uploaded: %q (%v bytes)", h.Filename, h.Size)
				// Use h.Open() to read the contents of the file.
				f, err := h.Open()
				if err != nil {
					fmt.Errorf("h.Open: %v", err)
					return
				}
				defer f.Close()

				ctx, cancel := context.WithTimeout(ctx, time.Second*50)
				defer cancel()

				// upload an object with storage.Writer
				var object = "chases/" + r.FormValue("ID") + "/" + h.Filename
				// var object = "chases/" + h.Filename
				fmt.Println(object)
				/* Google storage rules */
				// allow create: if request.auth.uid == request.resource.data.author_uid;
				const bucket = "chaseapp-8459b.appspot.com"
				wc := storageclient.Bucket(bucket).Object(object).NewWriter(ctx)

				wc.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": r.FormValue("ID") }

				/*
				_, err = fmt.Fprintf(w, "MediaLink: %v", wc.MediaLink)
				if err != nil {
					log.Printf("Couldn't write to http.responseWriter: %v", err)
					return
				}
				 */

				if _, err = io.Copy(wc, f); err != nil {
					fmt.Errorf("io.Copy: %v", err)
					return
				}

				if err := wc.Close(); err != nil {
					fmt.Errorf("Writer.Close: %v", err)
				}

				acl := storageclient.Bucket(bucket).Object(object).ACL()
				if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
					fmt.Errorf("ACLHandle.Set: %v", err)
					return
				}
				fmt.Printf("Blob %v is now publicly accessible.\n", object)
				imgUrl := CreateImageUrl(object, bucket)

				fmt.Fprintln(w, imgUrl)
				return
			}
		} else {
			http.Error(w, "Missing or bad headers", http.StatusInternalServerError)
			log.Printf("Missing or bad headers: %v", err)
		}
	}
}

type ImageStructure struct {
	ImageName string `json:"imageName"`
	URL       string `json:"url"`
}

func CreateImageUrl(imagePath string, bucket string) ImageStructure {
	imageStructure := ImageStructure{
		ImageName: imagePath,
		URL:       "https://storage.cloud.google.com/" + bucket + "/" + imagePath,
	}

	return imageStructure
}

// GetChase returns a particular chase given an ID
func GetChase(w http.ResponseWriter, r *http.Request) {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("client.Close: #{err}")
		}
	}(client)

	chases := client.Collection("chases")

	// Inputs
	var d struct {
		ID		string    `json:"ID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error, body must ID!")
		return
	}

	if d.ID == "" {
		fmt.Fprint(w, "Must supply ID!")
		return
	}

	doc, err := chases.Doc(d.ID).Get(ctx)

	if err != nil {
		fmt.Fprint(w, "Error retrieving documents")
		return
	}

	m := doc.Data()

	// m.ID = doc.Ref.ID

	fmt.Printf("Document data: %#v\n", m)
	js, err := json.Marshal(m)
	if err != nil {
		fmt.Fprintf(w, "Error marshaling JSON: %s", err)
		return
	}
	_, err = w.Write(js)
	if err != nil {
		fmt.Fprintf(w, "Error writing JSON to socket: #{err}")
		return
	}
}

// ListAirports queries the firehose collection and returns it as JSON in the response body
func ListAirports (w http.ResponseWriter, r *http.Request) {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("firestore.Client - cant close: #{err}")
		}
	}(client)

	type Airport struct {
		ID			string		`json:"id,omitempty"`
		Airport		string		`json:"airport"`
		City		string		`json:"city"`
		State		string		`json:"state"`
		IATA		string		`json:"iata"`
		ICAO		string		`json:"icao"`
		LiveATC[] struct {
			Name	string		`json:"name"`
			URL		string	 	`json:"url"`
		} `json:"live_atc"`
		Location	interface{}	`json:"location"`
	}

	airports := client.Collection("airports")
	docs, err := airports.Documents(ctx).GetAll()
	if err != nil {
		fmt.Fprint(w, "Error retrieving documents")
		return
	}

	var airportsArray []Airport

	for _, doc := range docs {
		var airportsData Airport
		if err := doc.DataTo(&airportsData); err != nil {
			fmt.Fprintf(w, "Error retrieving documents: %s", err)
			return
		}
		airportsData.ID = doc.Ref.ID
		airportsArray = append(airportsArray, airportsData)
	}

	// Convert our array to JSON and spit it out
	js, err := json.Marshal(airportsArray)
	if err != nil {
		log.Fatalf("json.Marshal - cant marshal airportsArray: #{err}")
	}
	_, err = w.Write(js)
	if err != nil {
		return
	}
}

// ListFirehose queries the firehose collection and returns it as JSON in the response body
func ListFirehose (w http.ResponseWriter, r *http.Request) {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("firestore.Client - cant close: #{err}")
		}
	}(client)

	type Firehose struct {
		ID			string		`json:"id,omitempty"`
		CreatedAt	time.Time	`json:"createdAt"`
		EventType	string		`json:"event_type"`
		Payload struct {
			Name	string		`json:"name"`
			URLs	[]string	`json:"urls"`
		} `json:"payload"`
	}

	firehoses := client.Collection("firehose").OrderBy("createdAt", firestore.Desc).Limit(10)
	docs, err := firehoses.Documents(ctx).GetAll()
	if err != nil {
		fmt.Fprint(w, "Error retrieving documents")
		return
	}

	var firehoseArray []Firehose

	for _, doc := range docs {
		var firehoseData Firehose
		if err := doc.DataTo(&firehoseData); err != nil {
			fmt.Fprintf(w, "Error retrieving documents: %s", err)
			return
		}
		firehoseData.ID = doc.Ref.ID
		firehoseArray = append(firehoseArray, firehoseData)
	}

	// Convert our array to JSON and spit it out
	js, err := json.Marshal(firehoseArray)
	if err != nil {
		log.Fatalf("Cant Marshal firehoseArray: #{err}")
	}
	_, err = w.Write(js)
	if err != nil {
		return
	}
}

// ListChases prints the JSON encoded "name", Desc, and "Url" fields in the body
// of the request or an error message if there isn't one.
func ListChases(w http.ResponseWriter, r *http.Request) {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	chases := client.Collection("chases").OrderBy("CreatedAt", firestore.Desc).Limit(20)
	docs, err := chases.Documents(ctx).GetAll()
	if err != nil {
		fmt.Fprint(w, "Error retrieving documents")
		return
	}

	// chasesArray used to store all our info we retrieve from listing out the documents
	// we store the info we care about, then spit out some JSON to the client
	var chasesArray []Chase

	for _, doc := range docs {
		var chaseData Chase
		if err := doc.DataTo(&chaseData); err != nil {
			fmt.Fprintf(w, "Error retrieving documents: %s", err)
			return
		}
		chaseData.ID = doc.Ref.ID
		chasesArray = append(chasesArray, chaseData)
	}

	// Convert our array to JSON and spit it out
	js, err := json.Marshal(chasesArray)
	if err != nil {
	}
	_, err = w.Write(js)
	if err != nil {
		return
	}
}

// AddChase adds a chase to the firebase database
func AddChase(w http.ResponseWriter, r *http.Request) {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	if APIKEY != r.Header.Get("X-ApiKey") {
		fmt.Fprint(w, "Invalid or missing X-ApiKey")
		log.Printf("Invalid or missing X-ApiKey")
		return
	}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer client.Close()


	if err := json.NewDecoder(r.Body).Decode(&AddChaseInput); err != nil {
		fmt.Fprint(w, "Error, body must contain name,desc,url,live!")
		return
	}

	if AddChaseInput.Name == "" {
		fmt.Fprint(w, "Must supply name!")
		return
	}

	AddChaseInput.CreatedAt = time.Now()
	id, err := uuid.NewUUID()
	AddChaseInput.ID = id.String()
	myRes, addErr := client.Collection("chases").Doc(id.String()).Set(ctx, AddChaseInput)
	if addErr != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %v Response: %v", addErr, myRes)
		addErrMsgJSON, addErr := json.Marshal(addErr)
		fmt.Fprint(w, addErr)
		_, err := w.Write(addErrMsgJSON)
		if err != nil {
			return
		}
	}

	resJSON, resErr := json.Marshal(id.String())
	if resErr != nil {
		log.Printf("Error marshaling JSON: %s", resErr)
		return
	}

	// Return a response to the client, including the ID of the newly created document
	_, err = w.Write(resJSON)
	if err != nil {
		return
	}

	// FCM messaging stuffs:
	fcm, fcmErr := app.Messaging(ctx)

	if fcmErr != nil {
		log.Printf("Problem initializing FCM Messaging: %v", fcmErr)
		return
	}

	imageURL := "https://chaseapp.tv/icon.png"
	if len(AddChaseInput.ImageURL) > 0 {
		imageURL = AddChaseInput.ImageURL
	}

	// TODO: verify that adding /chase/ worked
	chaseURL := "https://chaseapp.tv/chase/" + id.String()

	// If this is a live event,broadcast event to 'chases' topic in fcm messaging
	if AddChaseInput.Live == true {
		beamsClient, err := pushnotifications.New(PUSHER_BEAMS_INSTANCE, PUSHER_BEAMS_SECRET)
		if err != nil {
			fmt.Fprintf(w, "Couldn't create Beams client: %v", err.Error())
		}

		publishRequest := map[string]interface{}{
			"apns": map[string]interface{}{
				"aps": map[string]interface{}{
					"alert": map[string]interface{}{
						"title": AddChaseInput.Name,
						"body":  AddChaseInput.Desc,
					},
				},
			},
			"fcm": map[string]interface{}{
				"notification": map[string]interface{}{
					"title": AddChaseInput.Name,
					"body":  AddChaseInput.Desc,
					"imageurl": imageURL,
					"clickaction": chaseURL,
				},
			},
			"web": map[string]interface{}{
				"notification": map[string]interface{}{
					"title": AddChaseInput.Name,
					"body":  AddChaseInput.Desc,
					"imageurl": imageURL,
				},
			},
		}

		pubId, err := beamsClient.PublishToInterests([]string{"chases-notifications"}, publishRequest)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Publish Id:", pubId)
		}
		// The topic name can be optionally prefixed with "/topics/".
		topic := "chases"

		oneHour := time.Duration(1) * time.Hour
		badge := 42
		message := &messaging.Message{
			Notification: &messaging.Notification{
				Title: AddChaseInput.Name,
				Body:  AddChaseInput.Desc,
				ImageURL: imageURL,
			},
			Android: &messaging.AndroidConfig{
				TTL: &oneHour,
				Notification: &messaging.AndroidNotification{
					ImageURL: imageURL,
					ClickAction: chaseURL,
				},
			},
			APNS: &messaging.APNSConfig{
				Payload: &messaging.APNSPayload{
					Aps: &messaging.Aps{
						Badge: &badge,
					},
				},
			},
			Topic: topic,
		}

		// Send a message to the devices subscribed to the provided topic.
		response, err := fcm.Send(ctx, message)
		if err != nil {
			log.Println("Couldn't send message:", err)
			fmt.Println("Couldn't send message:", err)
			log.Fatalln(err)
			return
		}
		// Response is a message ID string.
		log.Println("Successfully sent message:", response)
	} else {
		log.Println("Didnt get Live set to true, not sending FCM push..", AddChaseInput)
	}
}

// SendNotify sends FCM Push notifications, we have this built into AddChase but the CRUD uses UpdateChase really,
// so we're just going to make this a separate function.
func SendNotify(w http.ResponseWriter, r *http.Request) {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	if APIKEY != r.Header.Get("X-ApiKey") {
		fmt.Fprint(w, "Invalid or missing X-ApiKey")
		log.Printf("Invalid or missing X-ApiKey")
		return
	}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer client.Close()

	if err := json.NewDecoder(r.Body).Decode(&NotifyRequest); err != nil {
		fmt.Fprint(w, "Error, body must contain name,desc,imageurl,url!")
		log.Printf("Error, body must contain name,desc,imageurl,url!")
		return
	}

	if NotifyRequest.Name == "" {
		fmt.Fprint(w, "Must supply name!")
		log.Printf("Must supply name!")
		return
	}

	if NotifyRequest.Desc == "" {
		fmt.Fprint(w, "Must supply desc!")
		log.Printf("Must supply desc!")
		return
	}

	if NotifyRequest.URL == "" {
		fmt.Fprint(w, "Must supply URL to Chase!")
		log.Printf("Must supply URL to chase!")
		return
	}

	if NotifyRequest.ImageURL == "" {
		fmt.Fprint(w, "Must supply ImageURL to Chase!")
		log.Printf("Must supply ImageURL to chase!")
		return
	}

	// FCM messaging stuffs:
	fcm, fcmErr := app.Messaging(ctx)

	if fcmErr != nil {
		fmt.Fprintf(w, "Problem initializing FCM Messaging: %v", fcmErr)
		log.Printf("Problem initializing FCM Messaging: %v", fcmErr)
		return
	}

	// The topic name can be optionally prefixed with "/topics/".
	topic := "chases"

	oneHour := time.Duration(1) * time.Hour
	badge := 42
	message := &messaging.Message{
		 Notification: &messaging.Notification{
			 Title: NotifyRequest.Name,
			 Body:  NotifyRequest.Desc,
			 ImageURL: NotifyRequest.ImageURL,
		 },
		 Android: &messaging.AndroidConfig{
			 TTL: &oneHour,
			 Notification: &messaging.AndroidNotification{
				 ImageURL: NotifyRequest.ImageURL,
				 ClickAction: NotifyRequest.URL,
			 },
		 },
		 APNS: &messaging.APNSConfig{
			 Payload: &messaging.APNSPayload{
				 Aps: &messaging.Aps{
					  Badge: &badge,
				 },
			 },
		 },
		 Topic: topic,
	}

	log.Println(message)
	// Send a message to the devices subscribed to the provided topic.
	response, err := fcm.Send(ctx, message)
	if err != nil {
		 log.Fatalln(err)
	}
	// Response is a message ID string.
	log.Println("Successfully sent message:", response)
}

// DeleteChase deletes a chase
func DeleteChase(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, r.Body)
	conf := &firebase.Config{ProjectID: projectID}

	if APIKEY != r.Header.Get("X-ApiKey") {
		fmt.Fprint(w, "Invalid or missing X-ApiKey")
		return
	}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer client.Close()

	var d struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error, body must contain ID!")
		return
	}

	fmt.Fprint(w, d)

	if d.ID == "" {
		fmt.Fprint(w, "Must supply ID!")
		return
	}

	fmt.Fprint(w, d)

	_, addErr := client.Collection("chases").Doc(d.ID).Delete(ctx)
	if addErr != nil {
		_, addErr := json.Marshal(addErr)
		fmt.Fprintf(w, "Error deleting chase record: %v", addErr)
	}

}

// SetAdmin if provided with the correct API key, will grant
// a UID admin in firebase authentication. This is used for chaseapp-crud.
func SetAdmin(w http.ResponseWriter, r *http.Request) {
	conf := &firebase.Config{ProjectID: projectID}
	var APIKEY = "8b373c2d-41bc-4a18-80f3-b3671f04930f"

	if APIKEY != r.Header.Get("X-ApiKey") {
		fmt.Fprint(w, "Invalid or missing X-ApiKey")
		return
	}

	var d struct {
		UID string `json:"uid"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error, body must contain UID!")
		return
	}

	if d.UID == "" {
		fmt.Fprint(w, "Must supply ID!")
		return
	}

	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	auth, err := app.Auth(ctx)
	err = auth.SetCustomUserClaims(ctx, d.UID, map[string]interface{}{"admin": true})
	if err != nil {
		log.Fatalf("Couldn't set UID %v to Admin", err)
	}

	fmt.Fprintf(w, "Successfully set UID %v to Admin", d.UID)
}

// UpdateChase updates a chase
func UpdateChase(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, r.Body)

	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	var APIKEY = "8b373c2d-41bc-4a18-80f3-b3671f04930f"
	// fmt.Fprint(w, r)
	// fmt.Fprint(w, r.Header.Get("X-ApiKey"))

	if APIKEY != r.Header.Get("X-ApiKey") {
		fmt.Fprint(w, "Invalid or missing X-ApiKey")
		return
	}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer client.Close()

	if err := json.NewDecoder(r.Body).Decode(&ChaseInput); err != nil {
		fmt.Fprint(w, "Error, body must contain ID, and name or desc, url, live, Networks!")
		return
	}

	fmt.Fprint(w, ChaseInput)

	if ChaseInput.ID == "" {
		fmt.Fprint(w, "Must supply ID!")
		return
	}

	fmt.Fprint(w, ChaseInput)

	if ChaseInput.Votes > 0 {
		_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"Votes": ChaseInput.Votes}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating chase votes: %v", addErr)
		}
	}

	if ChaseInput.Name != "" {
		_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"Name": ChaseInput.Name}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating chase name: %v", addErr)
		}
	}

	if ChaseInput.Desc != "" {
		_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"Desc": ChaseInput.Desc}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating chase desc: %v", addErr)
		}
	}

	if ChaseInput.Reddit != "" {
		_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"Reddit": ChaseInput.Reddit}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating chase desc: %v", addErr)
		}
	}

    if ChaseInput.ImageURL != "" {
		_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"ImageURL": ChaseInput.ImageURL}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating chase ImageURL: %v", addErr)
		}
	}

	if ChaseInput.Wheels.W1 != "" {
		var mappedData = transform.ToFirestoreMap(ChaseInput.Wheels)
		_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, mappedData, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating chase ChaseWheels: %v", addErr)
		}
	}

	if ChaseInput.Wheels.W2 != "" {
		_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"Wheels": ChaseInput.Wheels}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating chase ChaseWheels: %v", addErr)
		}
	}
	if ChaseInput.Wheels.W3 != "" {
		_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"Wheels": ChaseInput.Wheels}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating chase ChaseWheels: %v", addErr)
		}
	}
	if ChaseInput.Wheels.W4 != "" {
		_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"Wheels": ChaseInput.Wheels}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating chase ChaseWheels: %v", addErr)
		}
	}

	if !ChaseInput.EndedAt.IsZero() {
		log.Println("EndedAt is NotZero, ChaseInput:", ChaseInput)
        _, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"EndedAt": ChaseInput.EndedAt}, firestore.MergeAll)
        if addErr != nil {
            _, addErr := json.Marshal(addErr)
            fmt.Fprintf(w, "Error updating chase EndedAt: %v", addErr)
        }
	}

	if ChaseInput.Networks != nil {
		chases := client.Collection("chases")
		doc, err := chases.Doc(ChaseInput.ID).Get(ctx)
		if err != nil {
			fmt.Fprint(w, "Error retrieving documents")
			return
		}

		var networkArray []Networks
		var chaseData Chase
		if err := doc.DataTo(&chaseData); err != nil {
			fmt.Fprintf(w, "Error retrieving documents: %s", err)
			return
		}
		chaseData.ID = doc.Ref.ID
		for _, network := range chaseData.Networks {
			fmt.Fprintln(w, network.Name)
			networkArray = append(networkArray, network)
		}

		networkArray = append(networkArray, ChaseInput.Networks...)
		for _, network := range networkArray {
			fmt.Fprintln(w, network.URL)
		}

		networkArray = ChaseInput.Networks
		chaseData.Networks = ChaseInput.Networks
		_, urlErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"Networks": networkArray}, firestore.MergeAll)
		if urlErr != nil {
			fmt.Fprintf(w, "Error adding Networks to chase: %v", urlErr)
			return
		}
	}

	x := strconv.FormatBool(ChaseInput.Live)
	fmt.Fprintf(w, "X: %q", x)

	if ChaseInput.Live != true {
		if ChaseInput.Live != false {
			fmt.Fprint(w, "Live value must be true or false")
			return
		}
	}
	_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"Live": ChaseInput.Live}, firestore.MergeAll)
	if addErr != nil {
		_, addErr := json.Marshal(addErr)
		fmt.Fprintf(w, "Error updating chase name: %v", addErr)
	}

	pusherClient := pusher.Client{
		AppID:   PUSHER_API_APPID,
		Key:     PUSHER_API_KEY,
		Secret:  PUSHER_API_SECRET,
		Cluster: PUSHER_API_CLUSTER,
		Secure:  true,
	}

	// data := map[string]string{"message": "asdfsadf"}
	// pusherClient.Trigger("firehose", "updates", data)
	pusherClient.Trigger("chases", "updates", ChaseInput)
}
