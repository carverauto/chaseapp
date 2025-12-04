package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/octoper/go-ray"
	"io"
	"time"

	// firebase "firebase.google.com/go"
	"fmt"
	"math"
	"strings"
	// "reflect"

	"strconv"

	// "html"
	"log"
	"net/http"
	"os"

	"firebase.google.com/go/v4"
)

/* GLOBALS */

// GCLOUD_PROJECT is automatically set by the Cloud Functions runtime.
var projectID = os.Getenv("GCLOUD_PROJECT")

func main() {
	GetWeatherAlertsAndRasters()
}

/* NWS/NOAA */
/************/

// The code below was ported from
// https://stackoverflow.com/questions/238260/how-to-calculate-the-bounding-box-for-a-given-lat-lng-location

/*
// Semi-axes of WGS-84 geoidal reference
const (
	WGS84_a = 6378137.0 // Major semiaxis [m]
	WGS84_b = 6356752.3 // Minor semiaxis [m]
)

*/

/*
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

*/

// GetBoundingBox takes two arguments, MapPoint is a set of lat/lng,
// 'halfSideInKm' is the half length of the bounding box you want in kilometers.
/*
func GetBoundingBox(point MapPoint, halfSideInKm float64) BoundingBox {
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

	latMin := lat - halfSide/radius
	latMax := lat + halfSide/radius
	lonMin := lon - halfSide/pradius
	lonMax := lon + halfSide/pradius

	return BoundingBox{
		MinPoint: MapPoint{Latitude: Rad2deg(latMin), Longitude: Rad2deg(lonMin)},
		MaxPoint: MapPoint{Latitude: Rad2deg(latMax), Longitude: Rad2deg(lonMax)},
	}
}

*/

// https://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

type Geocode struct {
	UGC []string `json:"UGC"`
}

type Properties struct {
	ID            string   `json:"name"`
	AreaDesc      string   `json:"areaDesc"`
	Geocode       Geocode  `json:"geocode"`
	AffectedZones []string `json:"affectedZones"`
	Sent          string   `json:"sent"`
	Effective     string   `json:"effective"`
	Onset         string   `json:"onset"`
	Expires       string   `json:"expires"`
	Ends          string   `json:"ends"`
	Status        string   `json:"status"`
	MessageType   string   `json:"messageType"`
	Category      string   `json:"category"`
	Severity      string   `json:"severity"`
	Certainty     string   `json:"certainty"`
	Urgency       string   `json:"urgency"`
	Event         string   `json:"event"`
	Sender        string   `json:"sender"`
	SenderName    string   `json:"senderName"`
	Headline      string   `json:"headline"`
	Description   string   `json:"description"`
	Instruction   string   `json:"instruction"`
	Response      string   `json:"response"`
}

type WeatherAlert struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
}

type WeatherObject struct {
	Features []WeatherAlert `json:"features"`
}

// GetSmallestSurroundingRectangleByAreaURL returns the URL for the API
func GetSmallestSurroundingRectangleByAreaURL() string {
	return "https://us-central1-chaseapp-8459b.cloudfunctions.net/smallestSurroundingRectangleByArea"
}

type coords struct {
	Lat float64
	Lon float64
}

type smallestSurroundingRectangleAPIRequest struct {
	Type     string `json:"type"`
	Geometry struct {
		Type        string      `json:"type"`
		Coordinates [][]float64 `json:"coordinates"`
		// Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
}

type smallestSurroundingRectangleResponse struct {
	Type       string    `json:"type"`
	Bbox       []float64 `json:"bbox"`
	Properties struct {
	} `json:"properties"`
	Geometry struct {
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	} `json:"geometry"`
}

type Coords1 struct {
	Coordinates [][]float64 `json:"coordinates"`
}

type Coords2 struct {
	Coordinates [][][]float64 `json:"coordinates"`
}

type AlertData struct {
	Context struct {
		Version string `json:"@version"`
	} `json:"@context"`
	ID       string `json:"id"`
	Type     string `json:"type"`
	Geometry struct {
		Type string `json:"type"`
		// Coordinates [][][]float64 `json:"coordinates"`
		Coordinates []interface{} `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		ID                  string        `json:"@id"`
		Type                string        `json:"@type"`
		ID0                 string        `json:"id"`
		Type0               string        `json:"type"`
		Name                string        `json:"name"`
		EffectiveDate       time.Time     `json:"effectiveDate"`
		ExpirationDate      time.Time     `json:"expirationDate"`
		State               string        `json:"state"`
		Cwa                 []string      `json:"cwa"`
		ForecastOffices     []string      `json:"forecastOffices"`
		TimeZone            []string      `json:"timeZone"`
		ObservationStations []interface{} `json:"observationStations"`
		RadarStation        interface{}   `json:"radarStation"`
	} `json:"properties"`
}

// GetWeatherAlertsAndRasters will get and sort through the latest
// weather alerts from NWS, we're looking for the bad stuff.
// It also attempts to download raster images based on a list of raster image IDs
// we retrieve from the API, and then it is going to save the reference to those IDs
// this mostly works up until the last 1-2 steps and starts to fall apart.
// func GetWeatherAlertsAndRasters(w http.ResponseWriter, r *http.Request) {
func GetWeatherAlertsAndRasters() {
	var weatherAlertsURL = "https://api.weather.gov/alerts/active"

	resp, err := http.Get(weatherAlertsURL)
	if err != nil {
		fmt.Println("Error retrieving documents: %s", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: %s", err)
		return
	}
	// w.Write(body)

	var m WeatherObject
	if err := json.Unmarshal(body, &m); err != nil {
		panic(err)
	}
	log.Println("Number of alerts: ", len(m.Features))
	if len(m.Features) > 0 {
		for i := 0; i < len(m.Features); i++ {
			//w.Write([]byte(m[i].ID))
			// const severity = "Severe"
			const severity = "Severe"
			sevCheck := strings.Contains(severity, m.Features[i].Properties.Severity)

			// https://www.weather.gov/media/documentation/docs/NWS_Geolocation.pdf
			eventList := [8]string{"Storm Warning", "Severe Thunderstorm Warning", "Tornado Warning", "Snow Squall Warning", "Dust Storm Warning", "Extreme Wind Warning", "Special Marine Warning", "Flood Advisory"}
			// const event = "Severe Thunderstorm Warning"
			for _, event := range eventList {

				eventCheck := strings.Contains(event, m.Features[i].Properties.Event)

				if sevCheck {
					if eventCheck {
						//SaveWeatherAlerts(m.Features[i])
						// fmt.Fprintf(w, "Severity: %v Headline: %v\n", m.Features[i].Properties.Severity, m.Features[i].Properties.Event)
						fmt.Printf("Severity: %v Headline: %v\n\n", m.Features[i].Properties.Severity, m.Features[i].Properties.Event)
						if len(m.Features[i].Properties.AffectedZones) > 0 {
							for a := 0; a < len(m.Features[i].Properties.AffectedZones); a++ {
								// fmt.Fprintf(w, "AffectedZones: %v\n", m.Features[i].Properties.AffectedZones[a])
								aZone := m.Features[i].Properties.AffectedZones[a]

								fmt.Printf("AffectedZones: %v\n", aZone)

								resp, err := http.Get(aZone)
								if err != nil {
									fmt.Printf("Error retrieving documents: %s", err)
								}
								defer resp.Body.Close()
								countyBody, err := io.ReadAll(resp.Body)

								var c AlertData
								if err := json.Unmarshal(countyBody, &c); err != nil {
									panic(err)
								}

								// type assertion on c.Geometry.Coordinates
								// var c1 []float64
								var c2 [][]float64
								var coordsSlice []coords

								// print type
								for _, v := range c.Geometry.Coordinates {
									// convert []interface{} to []float64
									switch v.(type) {
									case []interface{}:
										for _, v2 := range v.([]interface{}) {
											// range through v2 and convert to float64
											switch v2.(type) {
											case []interface{}:
												for _, v3 := range v2.([]interface{}) {
													// convert []interface{} to []float64
													switch v3.(type) {
													case []interface{}:
														var c3 []float64
														for _, v4 := range (v3).([]interface{}) {
															// convert []interface{} to []float64
															c3 = append(c3, v4.(float64))
														}
														// iterate through c3 and append to coordsSlice
														for c3i := 0; c3i < len(c3); c3i++ {
															if c3i < len(c3)-1 {
																coordsSlice = append(coordsSlice, coords{c3[c3i+1], c3[c3i]})
															}
														}
													}
												}
											}
											// fmt.Println("Type: ", reflect.TypeOf(v2))
											// convert v2 from []interface{} to []float64
										}
									// convert [][]interface{} to [][]float64
									case [][]interface{}:
										for i, v2 := range v.([]interface{}) {
											c2 = append(c2, v2.([]float64))
											coordsSlice = append(coordsSlice, coords{c2[i][0], c2[i][1]})
										}
									}
								}
								fmt.Println("Length of coordsSlice: ", len(coordsSlice))
								// ray.Ray("Coordinates", coordsSlice)
								// ray.Ray("Coords", c.Geometry.Coordinates)
								// fmt.Println("coordsSlice", coordsSlice)
								var imageQueryCoords string
								if len(coordsSlice) > 2 {
									// get a minimum bounding rectangle for the coordinates
									// create the GeoJSON object
									var smallestRectRequest smallestSurroundingRectangleAPIRequest
									smallestRectRequest.Type = "Feature"
									smallestRectRequest.Geometry.Type = "LineString"

									// convert coordSlice to a [][]float64
									var coordsSliceFloat [][]float64
									for _, v := range coordsSlice {
										coordsSliceFloat = append(coordsSliceFloat, []float64{v.Lon, v.Lat})
									}
									smallestRectRequest.Geometry.Coordinates = coordsSliceFloat

									// create the request
									smallestRectRequestJSON, err := json.Marshal(smallestRectRequest)
									if err != nil {
										panic(err)
									}
									// ray.Ray("smallestRectRequestJSON", smallestRectRequestJSON)
									// submit the request to the API
									// req, rErr := http.NewRequest("POST", GetSmallestSurroundingRectangleByAreaURL(), bytes.NewBuffer(smallestRectRequestJSON))
									resp, rErr := http.Post(GetSmallestSurroundingRectangleByAreaURL(), "application/json", bytes.NewBuffer(smallestRectRequestJSON))
									if rErr != nil {
										panic(rErr)
									}

									defer resp.Body.Close()

									// read the response
									body, bErr := io.ReadAll(resp.Body)
									if bErr != nil {
										panic(bErr)
									}

									// unmarshal the response

									var smallestRectResponse smallestSurroundingRectangleResponse
									jsonErr := json.Unmarshal(body, &smallestRectResponse)
									if jsonErr != nil {
										fmt.Println("Error unmarshalling smallestRectResponse: ", jsonErr)
										fmt.Println("Body", string(body))
										panic(jsonErr)
									}

									// fmt.Println("smallestRectResponse", smallestRectResponse.Bbox)
									// get the coordinates from smallestRectResponse and concatenate into a string separated by commas
									/*
										for i := 0; i < len(smallestRectResponse.Geometry.Coordinates); i++ {
											imageQueryCoords += fmt.Sprintf("%v,%v,", smallestRectResponse.Geometry.Coordinates[i][0], smallestRectResponse.Geometry.Coordinates[i][1])
										}
									*/
									// create imageQueryCoords from smallestRectResponse.Bbox
									imageQueryCoords = fmt.Sprintf("%v,%v,%v,%v", smallestRectResponse.Bbox[0], smallestRectResponse.Bbox[1], smallestRectResponse.Bbox[2], smallestRectResponse.Bbox[3])
									// fmt.Println("imageQueryCoords", imageQueryCoords)
								} else {
									fmt.Println("Length: ", len(coordsSlice))
									// imageQueryCoords := strings.Join(coordsSlice, ",")
									// range through coordSlice and concatenate into a string separated by commas
									for idx, v := range coordsSlice {
										if idx == len(coordsSlice)-1 {
											imageQueryCoords += fmt.Sprintf("%v,%v", v.Lon, v.Lat)
										} else {
											imageQueryCoords += fmt.Sprintf("%v,%v,", v.Lon, v.Lat)
										}
									}

									/*
										for i := 0; i < len(coordsSlice); i++ {
											fmt.Println("coordsSlice[i]", coordsSlice[i])
											imageQueryCoords += fmt.Sprintf("%v,%v,", coordsSlice[i], coordsSlice[i+1])
											i++
										}

									*/
								}

								// ray.Ray("imageQueryCoords", imageQueryCoords)
								imageQuery := "https://idpgis.ncep.noaa.gov/arcgis/rest/services/radar/radar_base_reflectivity_time/ImageServer/query" + "?geometry=" + imageQueryCoords + "&f=json"
								fmt.Println("imageQuery", imageQuery)
								// fmt.Fprintf(w, "ImageQuery: %v", imageQuery)
								imageIds := GetImageIDs(imageQuery)
								fmt.Println("imageIds", imageIds)
								// ray.Ray("imageIds", imageIds)
								catalogItems := GetCatalogItems(imageIds)
								fmt.Println("CatalogItemsCount: %v", len(catalogItems))
								for c := 0; c < len(catalogItems); c++ {
									/*
										for ring := 0; ring < len(catalogItems[c].Geometry.Rings[0][0]); ring++ {
											var rings []float64 = catalogItems[c].Geometry.Rings[0][ring]
											rasterImage := GetRasterImage(catalogItems[c].Attributes.Objectid, rings)
											// ray.Ray("rasterImage", rasterImage)
											if rasterImage.Href != "" {
												fmt.Printf("Rasterimage: %v\n", rasterImage.Href)
												SaveRasterData(catalogItems[c].Attributes.Objectid, rings, rasterImage)
											}
										}
									*/

									// iterate through the rings and add the coords to the rings slice
									var rings []float64
									ray.Ray("catalogItems[c].Geometry.Rings", catalogItems[c].Geometry.Rings)
									for ring := 0; ring < len(catalogItems[c].Geometry.Rings[0]); ring++ {
										fmt.Println("ring", catalogItems[c].Geometry.Rings[0][ring])
										rings = append(rings, catalogItems[c].Geometry.Rings[0][ring]...)
									}
									fmt.Println("rings", rings)
									var coordsSlice []coords

									// convert rings to coordsSlice
									for i := 0; i < len(rings); i++ {
										coordsSlice = append(coordsSlice, coords{rings[i+1], rings[i]})
										fmt.Println("coordsSlice", coordsSlice)
										i++
									}

									// make a slice of coordinates
									iCoords := make([]coords, 0)
									// add the coordinates to the slice
									iCoords = append(iCoords, coordsSlice...)
									// convert iCoords to a [][]float64
									var iCoordsFloat [][]float64
									for _, v := range iCoords {
										iCoordsFloat = append(iCoordsFloat, []float64{v.Lon, v.Lat})
									}
									fmt.Println("iCoordsFloat", iCoordsFloat)

									// add the slice to the GeoJSON object
									// get a minimum bounding rectangle for the rings
									// create the GeoJSON object
									var smallestRectRequest smallestSurroundingRectangleAPIRequest
									smallestRectRequest.Type = "Feature"
									smallestRectRequest.Geometry.Type = "LineString"
									smallestRectRequest.Geometry.Coordinates = iCoordsFloat
									smallestRectRequest.Geometry.Coordinates = coordsSliceFloat

									ray.Ray(smallestRectRequest)

									// create the request
									smallestRectRequestJSON, err := json.Marshal(smallestRectRequest)
									if err != nil {
										panic(err)
									}
									// submit the request to the API
									// req, rErr := http.NewRequest("POST", GetSmallestSurroundingRectangleByAreaURL(), bytes.NewBuffer(smallestRectRequestJSON))
									resp, rErr := http.Post(GetSmallestSurroundingRectangleByAreaURL(), "application/json", bytes.NewBuffer(smallestRectRequestJSON))
									if rErr != nil {
										panic(rErr)
									}

									defer resp.Body.Close()

									// read the response
									body, bErr := io.ReadAll(resp.Body)
									if bErr != nil {
										panic(bErr)
									}

									// unmarshal the response
									var smallestRectResponse smallestSurroundingRectangleResponse
									jsonErr := json.Unmarshal(body, &smallestRectResponse)
									if jsonErr != nil {
										panic(jsonErr)
									}

									var ringsQueryCoords string

									// fmt.Println("smallestRectResponse", smallestRectResponse.Bbox)
									// get the coordinates from smallestRectResponse and concatenate into a string separated by commas
									for i := 0; i < len(smallestRectResponse.Geometry.Coordinates); i++ {
										ringsQueryCoords += fmt.Sprintf("%v,%v,", smallestRectResponse.Geometry.Coordinates[i][0], smallestRectResponse.Geometry.Coordinates[i][1])
										fmt.Println("ringsQueryCoords", ringsQueryCoords)
									}
									rasterImage := GetRasterImage(catalogItems[c].Attributes.Objectid, rings)
									// ray.Ray("rasterImage", rasterImage)
									if rasterImage.Href != "" {
										fmt.Printf("Rasterimage: %v\n", rasterImage.Href)
										// SaveRasterData(catalogItems[c].Attributes.Objectid, ringsQueryCoords, rasterImage)
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

func GetCatalogItems(imageIds []int) []CatalogItem {
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
		respBody, respBodyErr := io.ReadAll(resp.Body)
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

// SaveRasterData takes two arguments and will save save our raster data to firebase
// realtime database
func SaveRasterData(objectId int, area []float64, rasterImage RasterImage) {
	// firebase setup - take the response body and write it to firestore
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID, DatabaseURL: "https://chaseapp-8459b.firebaseio.com"}

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
	radarRefErr := radarRef.Update(ctx, map[string]interface{}{"ObjectId": objectId, "RasterImage": rasterImage, "Area": area})
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
	ray.Ray("GetRasterImage", objectId, bbox)
	sObjectId := strconv.Itoa(objectId)
	// sBoundingBox := strings.Join(bbox, ",")

	var BBox []string
	for b := 0; b < len(bbox); b++ {
		log.Println("BBOX[b]:", bbox[b])
		sB := fmt.Sprintf("%f", bbox[b])
		BBox = append(BBox, sB)
	}

	sBBox := strings.Join(BBox, ",")

	url := "https://idpgis.ncep.noaa.gov/arcgis/rest/services/radar/radar_base_reflectivity_time/ImageServer/" + sObjectId + "/image?bbox=" + sBBox + "&f=pjson"
	log.Printf("getRasterImage URL: %v\n", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBody, respErr := io.ReadAll(resp.Body)

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
	imageRespBody, imageRespErr := io.ReadAll(imageResp.Body)
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

	// ray.Ray("GetImageIDs", imageQuery, imageRespBody)
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

/*
// GetWeather takes the ID of a weather station and returns the data from the
// NOAA API
// func GetWeather(w http.ResponseWriter, r *http.Request) {
func GetWeather() {
	// Inputs
	var d struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Println( "Error, cannot decode body of request")
		return
	}

	if d.ID == "" {
		fmt.Print("Must supply weather station ID!")
		return
	}

	var weatherUrl = "https://api.tidesandcurrents.noaa.gov/mdapi/prod/webapi/stations/" + d.ID + ".json"
	resp, err := http.Get(weatherUrl)
	if err != nil {
		fmt.Printf( "Error retrieving documents: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	m := map[string]interface{}{}
	if err := json.Unmarshal(body, &m); err != nil {
		panic(err)
	}
	fmt.Printf("%q", m)

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
		fmt.Println(addErr)
		// _, err := w.Write(addErrMsgJSON)
		_, err := fmt.Println(addErrMsgJSON)
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
*/
