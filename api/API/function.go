// Package p contains an HTTP Cloud Function.
package p

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"encoding/xml"
	"google.golang.org/genproto/googleapis/type/latlng"
	"reflect"
	"regexp"
	"strconv"

	// firebase "firebase.google.com/go"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	// "reflect"

	// "strconv"

	// "html"
	"bytes"
	"log"
	"net/http"
	"os"

	// 	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"firebase.google.com/go/v4"
	// "firebase.google.com/go/v4/messaging"
	stream "github.com/GetStream/stream-chat-go/v5"
	streamChat "github.com/GetStream/stream-go2/v6"

	"github.com/google/uuid"
	"github.com/mfreeman451/golang/common/writers"
	"github.com/mfreeman451/helpers"
	"github.com/pusher/push-notifications-go"
	"github.com/pusher/pusher-http-go"
	"go.opencensus.io/trace"
	"google.golang.org/api/iterator"
)

/* GLOBALS */

// Used to get launch info and store it in firebase
// var rocketLaunchAPIKey = os.Getenv("ROCKETLAUNCHAPI")

// APIKEY is used to do some hokey authentication with our API
var APIKEY = os.Getenv("APIKEY")

// GetStream
// ServerClient, _ := stream.NewClient(getStream_API_KEY, getStream_API_SECRET)
var getStream_API_KEY = os.Getenv("GETSTREAM_API_KEY")
var getStream_API_SECRET = os.Getenv("GETSTREAM_API_SECRET")
var serverUserId = "MEz3wh0d9CTBn1ycPHtnGmLrJP62"

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

// var exporter *stackdriver.Exporter

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
	Lat float64
	Lon float64
}

/************/
/* TFR time */
/************/

type LiveATC struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Geopoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Airport struct {
	ID      string `json:"ID"`
	Airport string `json:"airport"`
	City    string `json:"city"`
	Iata    string `json:"iata"`
	Icao    string `json:"icao"`
	Liveatc []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"liveatc"`
	Location *latlng.LatLng `json:"location"`
	/*
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`

	 */
	// RadiusArc    string `json:"radiusArc"`
	RadiusArc    float64 `json:"radiusArc"`
	RadiusArcUoM string `json:"radiusArcUoM"`
	State        string `json:"state"`
}

type AreaGroup struct {
	AreaID       string   `json:"area_id,omitempty"`
	RadiusArc    float64  `json:"radius,omitempty"`
	UoMRadiusArc string   `json:"radius_type,omitempty"`
	Location     Geopoint `json:"location,omitempty"`
}

type TFR struct {
	Name       string            `json:"name"`
	StartTime  time.Time         `json:"startTime"`
	EndTime    time.Time         `json:"endTime"`
	AreaGroups []AreaGroup       `json:"areaGroups"`
	Details    []AloftAdvDetails `json:"details"`
	URL        []string          `json:"url"`
}

type AloftOutput struct {
	Color struct {
		Name string `json:"name,omitempty"`
		Hex  string `json:"hex,omitempty"`
		Rgb  []int  `json:"rgb,omitempty"`
	} `json:"color,omitempty"`
	Overview struct {
		Short string `json:"short,omitempty"`
		Full  string `json:"full,omitempty"`
		Icon  string `json:"icon,omitempty"`
	} `json:"overview,omitempty"`
	Airports   []string `json:"airports,omitempty"`
	Classes    []string `json:"classes,omitempty"`
	Advisories []AloftAdvisories
}
type AloftAdvDetails struct {
	Type  string      `json:"type,omitempty"`
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
	// Value string `json:"value,omitempty"`
}

type AloftAdvisories struct {
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Color struct {
		Name string `json:"name,omitempty"`
		Hex  string `json:"hex,omitempty"`
		Rgb  []int  `json:"rgb,om"`
	} `json:"color,omitempty"`
	Description string            `json:"description,omitempty"`
	Details     []AloftAdvDetails `json:"details,omitempty"`
	Geometry    string            `json:"geometry,omitempty"`
	Distance    struct {
		Unit  string  `json:"unit,omitempty"`
		Value float64 `json:"value,omitempty"`
		Lat   float64 `json:"lat,omitempty"`
		Long  float64 `json:"long,omitempty"`
	} `json:"distance,omitempty"`
	Properties struct {
		Link       string      `json:"LINK,omitempty"`
		Text       string      `json:"TEXT,omitempty"`
		Reason     string      `json:"REASON,omitempty"`
		TfrID      string      `json:"TFR_ID,omitempty"`
		EndsAt     int64       `json:"ENDS_AT,omitempty"`
		StartsAt   int         `json:"STARTS_AT,omitempty"`
		DaysOfWeek string      `json:"DAYS_OF_WEEK,omitempty"`
		OgcFid     int         `json:"OGC_FID,omitempty"`
		Objectid   int         `json:"OBJECTID,omitempty"`
		GlobalID   string      `json:"GLOBAL_ID,omitempty"`
		Ident      string      `json:"IDENT,omitempty"`
		Name       string      `json:"NAME,omitempty"`
		Latitude   string      `json:"LATITUDE,omitempty"`
		Longitude  string      `json:"LONGITUDE,omitempty"`
		Elevation  interface{} `json:"ELEVATION,omitempty"`
		IcaoID     string      `json:"ICAO_ID,omitempty"`
		TypeCode   string      `json:"TYPE_CODE,omitempty"`
		Servcity   string      `json:"SERVCITY,omitempty"`
		State      string      `json:"STATE,omitempty"`
		Country    string      `json:"COUNTRY,omitempty"`
		Operstatus string      `json:"OPERSTATUS,omitempty"`
		Privateuse int         `json:"PRIVATEUSE,omitempty"`
		Iapexists  int         `json:"IAPEXISTS,omitempty"`
		Dodhiflip  int         `json:"DODHIFLIP,omitempty"`
		Far91      int         `json:"FAR91,omitempty"`
		Far93      int         `json:"FAR93,omitempty"`
		MilCode    string      `json:"MIL_CODE,omitempty"`
		Airanal    string      `json:"AIRANAL,omitempty"`
		UsHigh     int         `json:"US_HIGH,omitempty"`
		UsLow      int         `json:"US_LOW,omitempty"`
		AkHigh     int         `json:"AK_HIGH,omitempty"`
		AkLow      int         `json:"AK_LOW,omitempty"`
		UsArea     int         `json:"US_AREA,omitempty"`
		Pacific    int         `json:"PACIFIC,omitempty"`
		Dist       string      `json:"DIST,omitempty"`
	} `json:"properties"`
}

type AloftRequest struct {
	Geometry struct {
		Coordinates [][][]float64 `json:"coordinates"`
		Type        string        `json:"type"`
	} `json:"geometry"`
}

type FAATFR struct {
	XMLName                   xml.Name `xml:"XNOTAM-Update"`
	Text                      string   `xml:",chardata"`
	Xsi                       string   `xml:"xsi,attr"`
	NoNamespaceSchemaLocation string   `xml:"noNamespaceSchemaLocation,attr"`
	Version                   string   `xml:"version,attr"`
	Origin                    string   `xml:"origin,attr"`
	Created                   string   `xml:"created,attr"`
	Group                     struct {
		Text string `xml:",chardata"`
		Add  struct {
			Text string `xml:",chardata"`
			Not  struct {
				Text   string `xml:",chardata"`
				NotUid struct {
					Text           string `xml:",chardata"`
					TxtNameAcctFac string `xml:"txtNameAcctFac"`
					DateIndexYear  string `xml:"dateIndexYear"`
					NoSeqNo        string `xml:"noSeqNo"`
					DateIssued     string `xml:"dateIssued"`
					TxtLocalName   string `xml:"txtLocalName"`
					CodeGUID       string `xml:"codeGUID"`
					NoUSNSWorkNo   string `xml:"noUSNSWorkNo"`
				} `xml:"NotUid"`
				CodeDailyOper          string `xml:"codeDailyOper"`
				DateEffective          string `xml:"dateEffective"`
				DateExpire             string `xml:"dateExpire"`
				CodeTimeZone           string `xml:"codeTimeZone"`
				CodeExpirationTimeZone string `xml:"codeExpirationTimeZone"`
				AffLocGroup            struct {
					Text           string `xml:",chardata"`
					TxtNameCity    string `xml:"txtNameCity"`
					TxtNameUSState string `xml:"txtNameUSState"`
				} `xml:"AffLocGroup"`
				CodeFacility string `xml:"codeFacility"`
				TfrNot       struct {
					Text         string `xml:",chardata"`
					CodeType     string `xml:"codeType"`
					TFRAreaGroup []struct {
						Text       string `xml:",chardata"`
						AseTFRArea struct {
							Text   string `xml:",chardata"`
							AseUid struct {
								Text     string `xml:",chardata"`
								CodeType string `xml:"codeType"`
								CodeId   string `xml:"codeId"`
							} `xml:"AseUid"`
							TxtName          string `xml:"txtName"`
							CodeDistVerUpper string `xml:"codeDistVerUpper"`
							ValDistVerUpper  string `xml:"valDistVerUpper"`
							UomDistVerUpper  string `xml:"uomDistVerUpper"`
							CodeDistVerLower string `xml:"codeDistVerLower"`
							ValDistVerLower  string `xml:"valDistVerLower"`
							UomDistVerLower  string `xml:"uomDistVerLower"`
							Att              struct {
								Text       string `xml:",chardata"`
								CodeWorkHr string `xml:"codeWorkHr"`
							} `xml:"Att"`
							CodeExclVerUpper   string `xml:"codeExclVerUpper"`
							CodeExclVerLower   string `xml:"codeExclVerLower"`
							IsScheduledTfrArea string `xml:"isScheduledTfrArea"`
							ScheduleGroup      struct {
								Text           string `xml:",chardata"`
								IsTimeSeparate string `xml:"isTimeSeparate"`
								DateEffective  string `xml:"dateEffective"`
								DateExpire     string `xml:"dateExpire"`
							} `xml:"ScheduleGroup"`
						} `xml:"aseTFRArea"`
						Aac struct {
							Text   string `xml:",chardata"`
							AacUid struct {
								Text      string `xml:",chardata"`
								AseUidChi struct {
									Text     string `xml:",chardata"`
									CodeType string `xml:"codeType"`
									CodeId   string `xml:"codeId"`
								} `xml:"AseUidChi"`
								AseUidPar struct {
									Text     string `xml:",chardata"`
									CodeType string `xml:"codeType"`
									CodeId   string `xml:"codeId"`
								} `xml:"AseUidPar"`
							} `xml:"AacUid"`
							CodeType string `xml:"codeType"`
							CodeOpr  string `xml:"codeOpr"`
						} `xml:"Aac"`
						AbdMergedArea struct {
							Text   string `xml:",chardata"`
							AbdUid struct {
								Text   string `xml:",chardata"`
								AseUid struct {
									Text     string `xml:",chardata"`
									CodeType string `xml:"codeType"`
									CodeId   string `xml:"codeId"`
								} `xml:"AseUid"`
							} `xml:"AbdUid"`
							TxtRmk string `xml:"txtRmk"`
							Avx    []struct {
								Text      string `xml:",chardata"`
								CodeDatum string `xml:"codeDatum"`
								CodeType  string `xml:"codeType"`
								GeoLat    string `xml:"geoLat"`
								GeoLong   string `xml:"geoLong"`
							} `xml:"Avx"`
						} `xml:"abdMergedArea"`
						AseShapes struct {
							Text   string `xml:",chardata"`
							AseUid struct {
								Text     string `xml:",chardata"`
								CodeType string `xml:"codeType"`
								CodeId   string `xml:"codeId"`
							} `xml:"AseUid"`
							Att struct {
								Text       string `xml:",chardata"`
								CodeWorkHr string `xml:"codeWorkHr"`
							} `xml:"Att"`
							Abd struct {
								Text   string `xml:",chardata"`
								AbdUid struct {
									Text   string `xml:",chardata"`
									AseUid struct {
										Text     string `xml:",chardata"`
										CodeType string `xml:"codeType"`
										CodeId   string `xml:"codeId"`
									} `xml:"AseUid"`
								} `xml:"AbdUid"`
								TxtRmk string `xml:"txtRmk"`
								Avx    struct {
									Text         string `xml:",chardata"`
									CodeType     string `xml:"codeType"`
									GeoLat       string `xml:"geoLat"`
									GeoLong      string `xml:"geoLong"`
									CodeDatum    string `xml:"codeDatum"`
									ValRadiusArc string `xml:"valRadiusArc"`
									UomRadiusArc string `xml:"uomRadiusArc"`
									Frd          struct {
										Text   string `xml:",chardata"`
										FrdUid struct {
											Text   string `xml:",chardata"`
											DpnUid struct {
												Text    string `xml:",chardata"`
												CodeId  string `xml:"codeId"`
												GeoLat  string `xml:"geoLat"`
												GeoLong string `xml:"geoLong"`
											} `xml:"DpnUid"`
										} `xml:"FrdUid"`
										Ain struct {
											Text   string `xml:",chardata"`
											VorUid struct {
												Text    string `xml:",chardata"`
												CodeId  string `xml:"codeId"`
												GeoLat  string `xml:"geoLat"`
												GeoLong string `xml:"geoLong"`
											} `xml:"VorUid"`
											ValAngleBrg string `xml:"valAngleBrg"`
										} `xml:"Ain"`
										Din struct {
											Text   string `xml:",chardata"`
											TcnUid struct {
												Text    string `xml:",chardata"`
												CodeId  string `xml:"codeId"`
												GeoLat  string `xml:"geoLat"`
												GeoLong string `xml:"geoLong"`
											} `xml:"TcnUid"`
											ValDist string `xml:"valDist"`
											UomDist string `xml:"uomDist"`
										} `xml:"Din"`
										TxtRmk string `xml:"txtRmk"`
									} `xml:"Frd"`
								} `xml:"Avx"`
							} `xml:"Abd"`
						} `xml:"aseShapes"`
						InstructionsGroup struct {
							Text     string   `xml:",chardata"`
							TxtInstr []string `xml:"txtInstr"`
						} `xml:"InstructionsGroup"`
						CodeIncFRD  string `xml:"codeIncFRD"`
						CodeShpPrt  string `xml:"codeShpPrt"`
						CodeLclTime string `xml:"codeLclTime"`
						CodeAuthATC string `xml:"codeAuthATC"`
					} `xml:"TFRAreaGroup"`
					TemplateType         string `xml:"TemplateType"`
					CodeCtrlFacilityType string `xml:"codeCtrlFacilityType"`
				} `xml:"TfrNot"`
				CodeCoordFacilityType string `xml:"codeCoordFacilityType"`
				TxtDescrUSNS          string `xml:"txtDescrUSNS"`
				TxtDescrTraditional   string `xml:"txtDescrTraditional"`
				TxtDescrModern        string `xml:"txtDescrModern"`
				CodeFreeformText      string `xml:"codeFreeformText"`
			} `xml:"Not"`
		} `xml:"Add"`
	} `xml:"Group"`
}

// airports is a list of airports we care about
var airports = []string{"Los Angeles Intl", "Van Nuys"}

func GetAloftAPIURL() string {
	return "https://air.aloft.ai/airspace-api/airspace"
}

func GetAloftRequest() []byte {
	var aloft AloftRequest
	aloft.Geometry.Type = "Polygon"
	// most of LA county
	aloft.Geometry.Coordinates = [][][]float64{
		{
			[]float64{-115.84516416069877, 32.69327082769155},
			[]float64{-115.84516416069877, 34.99069496744103},
			[]float64{-118.75898020350631, 34.99069496744103},
			[]float64{-118.75898020350631, 32.69327082769155},
			[]float64{-115.84516416069877, 32.69327082769155},
		},
	}

	// convert to []byte for the request
	aloftBytes, err := json.Marshal(aloft)
	if err != nil {
		log.Fatal(err)
	}
	return aloftBytes
}

type TFRGeoJSONFeature struct {
	Type       string `json:"type"`
	Properties struct {
		ID            string    `json:"id"`
		UoMRadiusArc  string    `json:"uomRadiusArc"`
		RadiusArc     float64   `json:"radiusArc"`
		DateEffective time.Time `json:"startTime"`
		DateExpire    time.Time `json:"endTime"`
		Details       string    `json:"details"`
		URL           string    `json:"url"`
	} `json:"properties"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	ID string `json:"id"`
}

type TFRGeoJSON struct {
	Type     string `json:"type"`
	Metadata struct {
		Generated int64  `json:"generated"`
		URL       string `json:"url"`
		Title     string `json:"title"`
		Status    int    `json:"status"`
		Count     int    `json:"count"`
	} `json:"metadata"`
	Features []TFRGeoJSONFeature `json:"features"`
}

// FetchAloftData queries the aloft API and returns the results
func FetchAloftData(w http.ResponseWriter, r *http.Request) {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID, DatabaseURL: "https://chaseapp-8459b.firebaseio.com"}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	aloftApiKey := os.Getenv("ALOFT_API_KEY")

	// create a bearer string
	bearer := "Bearer " + aloftApiKey

	// get the request
	aloftReq := GetAloftRequest()

	// create a new request using http
	req, err := http.NewRequest("POST", GetAloftAPIURL(), bytes.NewBuffer(aloftReq))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	// set the header to accept json
	req.Header.Set("Content-Type", "application/json")

	fmt.Println("Request: ", req)
	// send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return
	}

	// read the response body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}

	// unmarshal the response body into AloftOutput
	var aloft AloftOutput
	err = json.Unmarshal(body, &aloft)
	if err != nil {
		log.Printf("Error unmarshalling response body: %v", err)
		return
	}

	var TFRs []TFR

	// range through the aloft output
	for _, v := range aloft.Advisories {
		if v.Type == "tfr" {
			// ignore if the tfr is not active
			log.Println("TFR for: ", v.Name)
			// print out details
			log.Println("Details: ", v.Details)
			// convert epoch to time.Time
			startTime := time.Unix(int64(v.Properties.StartsAt), 0)
			endTime := time.Unix(int64(v.Properties.EndsAt), 0)
			log.Println("Starts:", v.Properties.StartsAt, "-", startTime)
			log.Println("Ends:", v.Properties.EndsAt, "-", endTime)

			var tfrData TFR
			tfrData.Name = v.Name
			tfrData.Details = v.Details
			tfrData.StartTime = startTime
			tfrData.EndTime = endTime
			log.Println("TFR: ", tfrData)
			TFRs = append(TFRs, tfrData)

			// extract the URL from the details
			// http://tfr.faa.gov/save_pages/detail_4_3635.html
			tfrRe, err := regexp.Compile(`http://tfr.faa.gov/save_pages/detail_4_\d+.html`)
			if err != nil {
				log.Printf("Error compiling regex: %v", err)
				return
			}
			// range through v.Details
			for _, d := range v.Details {
				// convert d.Value interface{} to string
				// use go reflection to determine type of d.Value
				if reflect.TypeOf(d.Value).String() == "string" {
					dValue := d.Value.(string)
					tfrURL := tfrRe.FindString(dValue)
					if tfrURL != "" {
						log.Println("TFR URL: ", tfrURL)
						tfrData.URL = append(tfrData.URL, tfrURL)
						// get the TFR details
						// convert the .html to .xml in tfrURL
						tfrXMLURL := strings.Replace(tfrURL, ".html", ".xml", 1)
						tfrDetails, err := GetTFRDetails(tfrXMLURL)
						log.Println("tfrXMLURL: ", tfrXMLURL)
						if err != nil {
							log.Printf("Error getting TFR details: %v", err)
							return
						}
						log.Println("TFR Details: ", tfrDetails.Text)
						log.Println("TFR DateEff: ", tfrDetails.Group.Add.Not.DateEffective)
						log.Println("TFR DateExpire: ", tfrDetails.Group.Add.Not.DateExpire)
						log.Println("TFRAreaGroup", len(tfrDetails.Group.Add.Not.TfrNot.TFRAreaGroup))
						// range through the TFRAreaGroup
						var areaGroups []AreaGroup
						if len(tfrDetails.Group.Add.Not.TfrNot.TFRAreaGroup) > 0 {
							for _, t := range tfrDetails.Group.Add.Not.TfrNot.TFRAreaGroup {
								var areaGroup AreaGroup
								areaGroup.AreaID = tfrDetails.Group.Add.Not.NotUid.TxtNameAcctFac

								lat, err := convertCoords(t.AseShapes.Abd.Avx.GeoLat)
								if err != nil {
									log.Printf("Error converting latitude: %v", err)
									return
								}
								lon, err := convertCoords(t.AseShapes.Abd.Avx.GeoLong)
								if err != nil {
									log.Printf("Error converting longitude: %v", err)
									return
								}

								log.Println("GeoLong", t.AseShapes.Abd.Avx.GeoLong, "Long:", lon)
								log.Println("GeoLat", t.AseShapes.Abd.Avx.GeoLat, "Lat:", lat)
								log.Println("ValRadiusArc", t.AseShapes.Abd.Avx.ValRadiusArc)
								log.Println("ValRadiusArcUnit", t.AseShapes.Abd.Avx.UomRadiusArc)

								areaGroup.Location.Longitude = lon
								areaGroup.Location.Latitude = lat
								// convert string ValRadiusArc to float
								radiusArc, err := strconv.ParseFloat(t.AseShapes.Abd.Avx.ValRadiusArc, 64)
								if err != nil {
									log.Printf("Error converting radius: %v", err)
									return
								}
								areaGroup.RadiusArc = radiusArc
								areaGroup.UoMRadiusArc = t.AseShapes.Abd.Avx.UomRadiusArc
								areaGroups = append(areaGroups, areaGroup)
							}
							tfrData.AreaGroups = areaGroups
							log.Println("AreaGroups: ", areaGroups)

							rtClient, rtErr := app.Database(ctx)
							if rtErr != nil {
								log.Fatalln("Error initializing db client: ", rtErr)
							}

							// populate the TFRGeoJSON struct
							var tfrGeoJSON TFRGeoJSON
							tfrGeoJSON.Type = "FeatureCollection"
							tfrGeoJSON.Metadata.Generated = time.Now().Unix()
							// tfrGeoJSON.Metadata.URL = tfrXMLURL
							tfrGeoJSON.Metadata.URL = GetAloftAPIURL()
							tfrGeoJSON.Metadata.Title = "TFRs"
							tfrGeoJSON.Metadata.Count = len(tfrDetails.Group.Add.Not.TfrNot.TFRAreaGroup)
							tfrGeoJSON.Metadata.Status = 200
							// populate the features struct
							var features []TFRGeoJSONFeature
							for _, a := range areaGroups {
								var feature TFRGeoJSONFeature
								feature.Type = "Feature"
								feature.Properties.ID = a.AreaID
								feature.Properties.Details = dValue
								feature.Properties.RadiusArc = a.RadiusArc
								feature.Properties.UoMRadiusArc = a.UoMRadiusArc
								feature.Properties.URL = tfrXMLURL

								// parse DateEffective from string to time.Time
								// 2022-11-04T00:00:00
								if tfrDetails.Group.Add.Not.DateEffective != "" {
									dateEffective, err := time.Parse("2006-01-02T15:04:05", tfrDetails.Group.Add.Not.DateEffective)
									// dateEffective, err := time.Parse(time.RFC3339, tfrDetails.Group.Add.Not.DateEffective)
									if err != nil {
										log.Fatalln("Error converting dateEffective", err)
									}
									feature.Properties.DateEffective = dateEffective
								} else {
									feature.Properties.DateEffective = time.Time{}
								}

								if tfrDetails.Group.Add.Not.DateExpire != "" {
									dateExpire, err := time.Parse("2006-01-02T15:04:05", tfrDetails.Group.Add.Not.DateExpire)
									// dateExpire, err := time.Parse(time.RFC3339, tfrDetails.Group.Add.Not.DateExpire)
									if err != nil {
										log.Fatalln("Error converting dateExpire", err)
									}
									feature.Properties.DateExpire = dateExpire
								} else {
									feature.Properties.DateExpire = time.Time{}
								}

								feature.Geometry.Type = "Point"
								feature.Geometry.Coordinates = []float64{a.Location.Longitude, a.Location.Latitude}
								features = append(features, feature)
							}
							tfrGeoJSON.Features = features

							rtRef := rtClient.NewRef("tfr")
							err := rtRef.Child("activeTFRs").Set(ctx, tfrGeoJSON)
							if err != nil {
								log.Fatalln("Error setting value: ", err)
							}

							/*
								for i := 0; i < len(areaGroups); i++ {
									rtRef.Child(tfrData.AreaGroups[i].AreaID)
									setErr := rtRef.Set(ctx, areaGroups[i])
									if setErr != nil {
										log.Fatalln("Error setting value:", setErr)
									}
								}
							*/
						}
					}
				} else {
					log.Println("d.Value is not a string, skipping..")
				}
			}
		}

		var airportList []Airport
		if v.Type == "airport" {
			// see if it exists in our list of airports
			for _, airport := range airports {
				if v.Properties.Name == airport {
					log.Println(v.Properties.Longitude, ",", v.Properties.Latitude)
					if v.Properties.Latitude != "" && v.Properties.Longitude != "" {
						// convert coordinates in DMS format to WGS84
						// split the latitude and longitude into degrees, minutes, seconds
						lat := strings.Split(v.Properties.Latitude, "-")
						lon := strings.Split(v.Properties.Longitude, "-")
						// convert the degrees, minutes, seconds to decimal degrees
						latDD := dmsToDD(lat)
						lonDD := dmsToDD(lon)
						log.Println(v.Properties.Longitude, ",", v.Properties.Latitude)
						var airportData Airport
						airportData.Location.Latitude = latDD
						airportData.Location.Longitude = lonDD
						airportData.Airport = v.Properties.Name
						log.Println("Airport: ", airportData)
						airportList = append(airportList, airportData)
					}
				}
			}
		}
	}
	fmt.Fprintln(w, "Success")
}

// GetTFRDetails gets the details of a TFR from the FAA
func GetTFRDetails(tfrURL string) (FAATFR, error) {
	req, err := http.NewRequest("GET", tfrURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return FAATFR{}, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return FAATFR{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	// unmarshal the XML in the body into a FAATFR struct
	var tfr FAATFR
	// log.Println("Body", res.Body)
	// convert body to io.Reader type
	bodyReader := bytes.NewReader(body)
	err = xml.NewDecoder(bodyReader).Decode(&tfr)
	if err != nil {
		log.Printf("Error decoding XML: %v", err)
		return FAATFR{}, err
	}
	return tfr, nil
}

func convertCoords(coord string) (float64, error) {
	var multiplier float64
	if strings.Contains(coord, "S") || strings.Contains(coord, "W") {
		multiplier = -1
	} else {
		multiplier = 1
	}
	// remove the last character from coord
	coord = coord[:len(coord)-1]
	// convert the string to a float64
	coordFloat, err := strconv.ParseFloat(coord, 64)
	if err != nil {
		log.Printf("Error converting coord to float64: %v", err)
		return 0, err
	}
	return coordFloat * multiplier, nil
}

func dmsToDD(dms []string) float64 {
	// handle S and W
	var multiplier float64
	if strings.Contains(dms[2], "S") || strings.Contains(dms[2], "W") {
		multiplier = -1
	} else {
		multiplier = 1
	}
	// remove the last character from sec
	dms[2] = dms[2][:len(dms[2])-1]

	// convert the degrees, minutes, seconds to decimal degrees
	deg, _ := strconv.ParseFloat(dms[0], 64)
	min, _ := strconv.ParseFloat(dms[1], 64)
	sec, _ := strconv.ParseFloat(dms[2], 64)
	// convert to decimal degrees
	return (deg + min/60 + sec/3600) * multiplier
}

/*******/
/* AIS */
/*******/

type AISBoat struct {
	MMSI      float64 `json:"mmsi"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	COG       float32 `json:"cog"`
	SOG       int     `json:"sog"`
	Heading   int     `json:"heading"`
	ROT       int     `json:"rot"`
	NavSat    int     `json:"nav_sat"`
	IMO       int     `json:"imo"`
	Name      string  `json:"name"`
	Callsign  string  `json:"callsign"`
	Type      int     `json:"type"`
	A         int     `json:"a"`
	B         int     `json:"b"`
	C         int     `json:"c"`
	D         int     `json:"d"`
	Draught   float32 `json:"draught"`
	Dest      string  `json:"dest"`
	ETA       string  `json:"eta"`
}

func GetBoats(w http.ResponseWriter, r *http.Request) {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID, DatabaseURL: "https://chaseapp-8459b.firebaseio.com"}

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var aisHubBoats [][]AISBoat

	rtClient, rtErr := app.Database(ctx)
	if rtErr != nil {
		log.Fatalln("Error initializing db client: ", rtErr)
	}

	jErr := json.Unmarshal(body, &aisHubBoats)
	if jErr != nil {
		fmt.Println("Error unmarshalling JSON: ", jErr)
		return
	}

	rtRef := rtClient.NewRef("ships")
	for i := 0; i < len(aisHubBoats[1]); i++ {
		rtRef.Child(aisHubBoats[1][i].Name)
		setErr := rtRef.Set(ctx, aisHubBoats)
		if setErr != nil {
			log.Fatalln("Error setting value:", setErr)
		}
	}
}

type Boats struct {
	MMSI  string
	Group string
	Type  string
}

// DeleteBoat deletes a Boat from firebase firestore
func DeleteBoat(w http.ResponseWriter, r *http.Request) {

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

	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

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
		ID    string `json:"id"`
		MMSI  string `json:"mmsi"`
		Group string `json:"desc"`
		Type  string `json:"type"`
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

	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

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
		UserID string `json:"user_id"`
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
		UserID  string `json:"user_id"`
		Channel string `json:"channel"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprintf(w, "Error, body must contain user_id, channel")
		return
	}

	// the secret is only used server side and gives you full access to the API
	client, _ := stream.NewClient(getStream_API_KEY, getStream_API_SECRET)
	channel := client.Channel("livestream", d.Channel)
	_, err := channel.AddModerators(context.Background(), d.UserID)
	if err != nil {
		fmt.Fprintf(w, "Error assigning a moderator to a channel %v", err)
	}
	fmt.Fprintf(w, "Added %v to moderators group on channel", d.UserID)
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
func UpdateToken(token string) interface{} {
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

/*******/
/* AIS */
/*******/

func AddBoat(w http.ResponseWriter, r *http.Request) {
	// Inputs
	var d struct {
		MMSI  string `json:"mmsi"`
		Group string `json:"group"`
		Type  string `json:"type"`
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

type Airship struct {
	Tailno   string
	ImageUrl string
	Group    string
	Type     string
}

// DeleteAirship deletes an Airship from firebase firestore
func DeleteAirship(w http.ResponseWriter, r *http.Request) {
	conf := &firebase.Config{ProjectID: projectID}

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

	if d.ID == "" {
		fmt.Fprint(w, "Must supply ID!")
		return
	}

	fmt.Println("Deleting Airship: ", d.ID)
	log.Println("Deleting Airship: ", d.ID)

	// search for the document
	query, qErr := client.Collection("airships").Where("tailno", "==", d.ID).Documents(ctx).GetAll()
	if qErr != nil {
		fmt.Fprint(w, "Error retrieving documents")
		return
	}

	fmt.Println(query)
	// get docId from query[0]
	if query != nil {
		fmt.Println("DocID ", query[0].Ref.ID)
		_, err := client.Collection("airships").Doc(query[len(query)-1].Ref.ID).Delete(ctx)
		if err != nil {
			fmt.Fprint(w, "Error deleting document")
			log.Println("Error deleting document from airships: ", err)
			return
		}
	} else {
		fmt.Fprint(w, "Error, document not found")
		log.Println("Error, document not found")
		return
	}
}

// GetAirship returns the JSON encoded Airship in the body
func GetAirship(w http.ResponseWriter, r *http.Request) {
	var d struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error, body must contain ID!")
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

	fmt.Println("ID: ", d.ID)

	query, err := client.Collection("airships").Where("tailno", "==", d.ID).Documents(ctx).GetAll()
	if err != nil {
		fmt.Fprint(w, "Error retrieving documents")
		return
	}

	if query != nil {
		var airshipData Airship
		// if err := query[0].DataTo(&airshipData); err != nil {
		if err := query[len(query)-1].DataTo(&airshipData); err != nil {
			fmt.Fprintf(w, "Error retrieving documents: %s", err)
			return
		}

		// return the airship
		js, err := json.Marshal(airshipData)
		if err != nil {
			fmt.Fprintf(w, "Error marshaling JSON: %s", err)
			return
		}
		_, err = w.Write(js)
		if err != nil {
			fmt.Fprintf(w, "Error writing JSON: %s", err)
			return
		}
	}
}

func UpdateAirship(w http.ResponseWriter, r *http.Request) {
	// Inputs
	var d struct {
		ID       string `json:"id"`
		Tailno   string `json:"tailno"`
		Group    string `json:"group"`
		ImageUrl string `json:"imageUrl"`
		Type     string `json:"type"`
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
	if d.ImageUrl != "" {
		_, addErr := client.Collection("airships").Doc(d.ID).Set(ctx, map[string]interface{}{"imageUrl": d.ImageUrl}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating airship imageUrl: %v", addErr)
		}
	}
	if d.Type != "" {
		_, addErr := client.Collection("airships").Doc(d.ID).Set(ctx, map[string]interface{}{"type": d.Type}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating airship type: %v", addErr)
		}
	}
}

// AddAirship takes tailno, group, and imageUrl (strings) in an HTTP request body
func AddAirship(w http.ResponseWriter, r *http.Request) {
	// Inputs
	var d struct {
		Tailno   string `json:"tailno"`
		Group    string `json:"group"`
		ImageURL string `json:"imageUrl"`
		Type     string `json:"type"`
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

	var airshipsArray []Airship

	for _, doc := range docs {
		var airshipData Airship
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

	var airshipsArray []Airship

	for _, doc := range docs {
		var airshipData Airship
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

type Tags struct {
	Name []string
}

/*
type jsonTags interface {

}
*/

type Networks struct {
	Name   string
	URL    string
	Tier   int
	Logo   string
	Other  string
	MP4URL string
}

type Wheels struct {
	W1 string
	W2 string
	W3 string
	W4 string
}

type Sentiment struct {
	Magnitude float64 `firestore:"magnitude"`
	Score     float64 `firestore:"score"`
}

type Chase struct {
	ID        string     ""
	Name      string     `firestore:"Name"`
	Desc      string     `firestore:"Desc"`
	Live      bool       `firestore:"Live"`
	Networks  []Networks `firestore:"Networks"`
	Wheels    Wheels     `firestore:"Wheels"`
	Votes     int        `firestore:"Votes"`
	CreatedAt time.Time  `firestore:"CreatedAt"`
	EndedAt   time.Time  `firestore:"EndedAt"`
	ImageURL  string     `firestore:"ImageURL"`
	Reddit    string     `firestore:"Reddit"`
	Sentiment Sentiment
	Tags      Tags
}

type NotifyRequest struct {
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	ImageURL string `json:"imageURL"`
	URL      string `json:"url"`
}

type myJsonNotifyData struct {
	Id       string `json:"id,omitempty"`
	Tweet_id string `json:"tweet_id,omitempty"`
	Image    string `json:"image,omitempty"`
}

type NotifyInput struct {
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	Interest  string    `json:"interest"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	Data      myJsonNotifyData
}

var AddChaseInput struct {
	ID        string     `json:"id,omitempty"`
	Name      string     `json:"name"`
	Desc      string     `json:"desc"`
	Live      bool       `json:"live"`
	Networks  []Networks `json:"networks"`
	Wheels    Wheels     `json:"wheels"`
	Votes     int        `json:"votes"`
	CreatedAt time.Time  `json:"createdAt"`
	ImageURL  string     `json:"imageURL"`
	Reddit    string     `json:"reddit"`
	Notified  bool       `json:"notified"`
	Tags      []string   `json:"tags"`
}

var ChaseInput struct {
	ID        string     `json:"id,omitempty"`
	Name      string     `json:"name"`
	Desc      string     `json:"desc"`
	Live      bool       `json:"live"`
	Networks  []Networks `json:"networks"`
	Wheels    Wheels     `json:"wheels"`
	Votes     int        `json:"votes"`
	CreatedAt time.Time  `json:"createdAt"`
	EndedAt   time.Time  `json:"endedAt,omitempty"`
	ImageURL  string     `json:"imageURL"`
	Reddit    string     `json:"reddit"`
	Notified  bool       `json:"notified"`
	Tags      string     `json:"tags"`
}

type PushTokens struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	TokenType string    `json:"type"`
}

type User struct {
	UID         string       `firestore:"uid"`
	LastUpdated time.Time    `firestore:"lastupdated"`
	PhotoURL    string       `firestore:"photourl"`
	UserName    string       `firestore:"username"`
	Tokens      []PushTokens `firestore:"tokens"`
}

var UserInput struct {
	UID         string       `json:"uid"`
	LastUpdated time.Time    `json:"lastupdated"`
	PhotoURL    string       `json:"photourl,omitempty"`
	UserName    string       `json:"username"`
	Tokens      []PushTokens `json:"tokens"`
}

/* CHASEAPP-CRUD */

// ListImages takes an input and returns a list of images for that chaseID
func ListImages(w http.ResponseWriter, r *http.Request) {
	// Inputs
	var d struct {
		ID string `json:"ID"`
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

				wc.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": r.FormValue("ID")}

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
		ID string `json:"ID"`
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

// GetAirport returns a particular airport given an ICAO
func GetAirport(w http.ResponseWriter, r *http.Request) {
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

	// Inputs
	var d struct {
		ID string `json:"ID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error, body must include an ID!")
		return
	}

	if d.ID == "" {
		fmt.Fprint(w, "Must supply ID!")
		return
	}

	airports := client.Collection("airports")
	// search through the collection for the document that matches d.ID
	doc, err := airports.Doc(d.ID).Get(ctx)
	if err != nil {
		fmt.Fprint(w, "Error retrieving documents")
		return
	}

	m := doc.Data()

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


// DeleteAirport deletes an Airport from firebase firestore
func DeleteAirport(w http.ResponseWriter, r *http.Request) {
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

	if d.ID == "" {
		fmt.Fprint(w, "Must supply ID!")
		return
	}

	fmt.Println("Deleting Airport: ", d.ID)
	log.Println("Deleting Airport: ", d.ID)

	// delete by document ID from firestore
	_, err = client.Collection("airports").Doc(d.ID).Delete(ctx)
	if err != nil {
		fmt.Fprint(w, "Error deleting document")
		return
	}

	fmt.Fprint(w, "Airport deleted")
}

// AddAirport adds an airport to firestore
func AddAirport(w http.ResponseWriter, r *http.Request) {

	// pre flight stuff
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if APIKEY != r.Header.Get("X-ApiKey") {
		fmt.Fprint(w, "Invalid or missing X-ApiKey")
		log.Printf("Invalid or missing X-ApiKey")
		return
	}

	var d Airport

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error, body must contain Airport object!")
		return
	}

	fmt.Println("D:", d)

	if d.Airport == "" {
		fmt.Fprint(w, "Must supply Airport (name)!")
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

	// .CreatedAt = time.Now()
	id, err := uuid.NewUUID()
	d.ID = id.String()

	/*
	geoP := latlng.LatLng{
		Latitude:  d.Location.Latitude,
		Longitude: d.Location.Longitude,
	}

	d.Location = &geoP

	 */

	myRes, addErr := client.Collection("airports").Doc(id.String()).Set(ctx, d)
	if addErr != nil {
		// Handle any errors in an appropriate way, such as returning them
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

	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Return a response to the client, including the ID of the newly created document
	_, err = w.Write(resJSON)
	if err != nil {
		return
	}
}

// UpdateAirport updates a particular airport given an ID
func UpdateAirport(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Updating airport")

	type Airport2 struct {
		ID      string `json:"ID"`
		Airport string `json:"airport"`
		City    string `json:"city"`
		Iata    string `json:"iata"`
		Icao    string `json:"icao"`
		Liveatc []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"liveatc"`
		Location *latlng.LatLng `json:"location"`
		RadiusArc    interface{} `json:"radiusArc"`
		RadiusArcUoM string `json:"radiusArcUoM"`
		State        string `json:"state"`
	}

	var d Airport2
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error, body must match Airport object!")
		fmt.Println("Error, body must match Airport object!", err)
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

	fmt.Println("D", d)

	if d.Airport != "" {
		_, addErr := client.Collection("airports").Doc(d.ID).Set(ctx, map[string]interface{}{"Airport": d.Airport}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating airsport name: %v", addErr)
		}
	}
	if d.Icao != "" {
		_, addErr := client.Collection("airports").Doc(d.ID).Set(ctx, map[string]interface{}{"Icao": d.Icao}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating airport ICAO: %v", addErr)
		}
	}
	if d.Iata != "" {
		_, addErr := client.Collection("airports").Doc(d.ID).Set(ctx, map[string]interface{}{"Iata": d.Iata}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating airport IATA: %v", addErr)
		}
	}
	if d.City != "" {
		_, addErr := client.Collection("airports").Doc(d.ID).Set(ctx, map[string]interface{}{"City": d.City}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating airport City: %v", addErr)
		}
	}
	if d.State != "" {
		_, addErr := client.Collection("airports").Doc(d.ID).Set(ctx, map[string]interface{}{"State": d.State}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating airport City: %v", addErr)
		}
	}
	fmt.Println("RadiusArc", d.RadiusArc)
	/*
	if d.RadiusArc != "" {
		fmt.Println("Updating radius", d.RadiusArc)
		// convert the radius to a float
		radius, err := strconv.ParseFloat(d.RadiusArc, 64)
		if err != nil {
			fmt.Fprintf(w, "Error converting radius to float: %v", err)
		}

		// _, addErr := client.Collection("airports").Doc(d.ID).Set(ctx, map[string]interface{}{"radiusArc": d.RadiusArc}, firestore.MergeAll)
		_, addErr := client.Collection("airports").Doc(d.ID).Set(ctx, map[string]interface{}{"radiusArc": radius}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating airport boundary radiusArc: %v", addErr)
		}
	} else {
		fmt.Println("Not updating radius")
	}
	 */
	if d.RadiusArc != 0 {
		fmt.Println("Updating radius", d.RadiusArc)
		// use reflection to check if the radius is a float or a string, if it is a string, convert to float64
		// if it is a float64, use it as is
		var radius float64
		switch d.RadiusArc.(type) {
		case float64:
			fmt.Println("Radius is a float64")
			radius = d.RadiusArc.(float64)
		case string:
			// convert the radius to a float
			radius, err = strconv.ParseFloat(d.RadiusArc.(string), 64)
		}

		_, addErr := client.Collection("airports").Doc(d.ID).Set(ctx, map[string]interface{}{"RadiusArc": radius }, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating airport boundary radiusArc: %v", addErr)
		}
	} else {
		fmt.Println("Not updating radius")
	}

	// check for the presence of d.Location.Latitude and d.Location.Longitude
	if d.Location.Latitude != 0 && d.Location.Longitude != 0 {
		fmt.Println("Updating location")
		geoP := latlng.LatLng{
			Latitude:  d.Location.Latitude,
			Longitude: d.Location.Longitude,
		}
		_, addErr := client.Collection("airports").Doc(d.ID).Set(ctx, map[string]interface{}{"Location": &geoP}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating airship type: %v", addErr)
		}
	}

	fmt.Println("Updated airport")
	// write our success message
	fmt.Fprintf(w, "Successfully updated airport: %v", d.ID)
}

// ListAirports queries the firehose collection and returns it as JSON in the response body
func ListAirports(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println(string(js))
	_, err = w.Write(js)
	if err != nil {
		return
	}
}

// AnimationEvent takes a few arguments and is used to trigger updates in Rive animations
// for the mystery chase theater 9000 feature
func AnimationEvent(w http.ResponseWriter, r *http.Request) {

	// Inputs
	var d struct {
		ID        string `json:"id"`
		Label     string `json:"label"`
		Endpoint  string `json:"endpoint"`
		AnimType  string `json:"anim_type"`
		AnimState string `json:"anim_state"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprintf(w, "Error, body must required fields!: %v", err)
		return
	}

	if d.ID == "" {
		fmt.Fprint(w, "Must supply ID!")
		return
	}

	// pre flight stuff
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if APIKEY != r.Header.Get("X-ApiKey") {
		fmt.Fprint(w, "Invalid or missing X-ApiKey")
		log.Printf("Invalid or missing X-ApiKey")
		return
	}

	client, err := streamChat.New(getStream_API_KEY, getStream_API_SECRET)
	if err != nil {
		log.Fatalf("Error with getstream activity feed: %v", err)
		return
	}

	// Create a feed
	feed, err := client.FlatFeed("events", "animation")
	if err != nil {
		log.Fatalf("Error with creating getstream activity feed: %v", err)
		return
	}

	/*
		{
		int label:12412// appearance time in the video in milliseconds
		String animation_endpoint://endpoint for that rive animation to fetch
		String type:// pop-up or theater,
		String state:// for theater we will have three different states, horse fist, //robot hop and the man one.// but I think pop-up animations might //not have any other states than the original animation for it.
		}
	*/

	// Create an activity
	resp, err := feed.AddActivity(streamChat.Activity{
		// ID:        "",
		Actor:  "animation",
		Verb:   "event-" + d.ID,
		Object: "animation-event",
		// ForeignID: "",
		// Target:    "",
		Time: streamChat.Time{Time: time.Now()},
		// Origin:    "",
		// To:        nil,
		// Score:     0,
		Extra: map[string]interface{}{
			"label":     d.Label,
			"endpoint":  d.Endpoint,
			"animtype":  d.AnimType,
			"animstate": d.AnimState,
			"createdAt": time.Now(),
		},
	})
	if err != nil {
		log.Fatalf("Cant add activity: %v", err)
		return
	}

	fmt.Fprintln(w, resp)
}

/***************/
/* Chase Stuff */
/***************/

type GetMP4Link struct {
	ChaseID string `json:"chase_id"`
}

// ListChases prints the JSON encoded "name", Desc, and "Url" fields in the body
// of the request or an error message if there isn't one.
func ListChases(w http.ResponseWriter, r *http.Request) {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// pre flight stuff
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

type Actions struct {
	StateA bool
	StateB bool
	StateC bool
}

type Activity struct {
	Name      string
	Actions   Actions
	ChaseID   string
	CreatedAt time.Time
}

type ChannelSettings struct {
	ChannelType string
	Id          string
	CustomName  string
}

func createChannel(channelSettings ChannelSettings) error {
	// instantiate your stream client using the API key and secret
	// the secret is only used server side and gives you full access to the API
	client, err := stream.NewClient(getStream_API_KEY, getStream_API_SECRET)
	if err != nil {
		return errors.Unwrap(fmt.Errorf("error with newClient: %v", err))
	}

	ctx := context.Background()

	client.CreateChannel(ctx, channelSettings.ChannelType, channelSettings.Id, serverUserId, nil)

	return nil
}

func addMod(mod string, chatChannel string) error {
	// instantiate your stream client using the API key and secret
	// the secret is only used server side and gives you full access to the API
	client, err := stream.NewClient(getStream_API_KEY, getStream_API_SECRET)
	if err != nil {
		return errors.Unwrap(fmt.Errorf("error with newClient: %v", err))
	}

	channel := client.Channel("livestream", chatChannel)
	channel.AddModerators(context.Background(), mod)

	return nil
}

// sendMessage takes 3 arguments, the message, chat channel, and userId
// and will send that message to the channel
func sendMessage(msg string, chatChannel string, userId string) error {
	// instantiate your stream client using the API key and secret
	// the secret is only used server side and gives you full access to the API
	client, err := stream.NewClient(getStream_API_KEY, getStream_API_SECRET)
	if err != nil {
		return errors.Unwrap(fmt.Errorf("error with newClient: %v", err))
	}

	ctx := context.Background()

	channel := client.Channel("livestream", chatChannel)
	message := &stream.Message{
		Text: msg,
	}
	_, err = channel.SendMessage(ctx, message, userId)
	if err != nil {
		log.Fatalf("Couldn't send message: %v", err)
		return err
	}
	// channel.AddModerators(context.Background(), mod)

	return nil
}

// GetStreamLinkURL returns the URL for the cloud function
func GetStreamLinkURL() string {
	return "https://us-central1-chaseapp-8459b.cloudfunctions.net/GetStreams"
}

// AddChase adds a chase to the firebase database
func AddChase(w http.ResponseWriter, r *http.Request) {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// pre flight stuff
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

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
	// AddChaseInput.Notified = true

	myRes, addErr := client.Collection("chases").Doc(id.String()).Set(ctx, AddChaseInput)
	if addErr != nil {
		// Handle any errors in an appropriate way, such as returning them
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

	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Return a response to the client, including the ID of the newly created document
	_, err = w.Write(resJSON)
	if err != nil {
		return
	}

	/*
		// FCM messaging stuffs:
		fcm, fcmErr := app.Messaging(ctx)

		if fcmErr != nil {
			log.Printf("Problem initializing FCM Messaging: %v", fcmErr)
			return
		}
	*/

	imageURL := "https://chaseapp.tv/icon.png"
	if len(AddChaseInput.ImageURL) > 0 {
		imageURL = AddChaseInput.ImageURL
	}

	// Call our function to scrape the network URLs for MP4s
	request := GetMP4Link{
		ChaseID: AddChaseInput.ID,
	}
	jsonValue, _ := json.Marshal(request)
	resp, err := http.Post(GetStreamLinkURL(), "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Fprintf(w, "Error calling GetStreams: %v", err)
		return
	}
	defer resp.Body.Close()
	_, iErr := ioutil.ReadAll(resp.Body)
	if iErr != nil {
		fmt.Fprintf(w, "Error reading response body: %v", iErr)
	}

	// TODO: verify that adding /chase/ worked
	chaseURL := "https://chaseapp.tv/chase/" + id.String()

	var channelSettings ChannelSettings
	channelSettings.ChannelType = "livestream"
	channelSettings.Id = AddChaseInput.ID
	channelSettings.CustomName = AddChaseInput.Name

	// create the channel so that we can control who owns it
	err = createChannel(channelSettings)
	if err != nil {
		log.Println("Problem creating channel: %v", err)
	}

	// TODO: need a better way to implement moderators someday
	addModdErr := addMod("leku", channelSettings.Id)
	if addModdErr != nil {
		log.Println("Problem adding mod: %v", err)
	}

	// Populate the channel with an initial message
	var channelMessage = "Channel created on " + time.Now().String()
	sendMessageErr := sendMessage(channelMessage, channelSettings.Id, "leku")
	if sendMessageErr != nil {
		log.Println("Problem sending message: %v", err)
	}

	// If this is a live event,broadcast event to 'chases' topic in fcm messaging
	if AddChaseInput.Live == true {
		beamsClient, err := pushnotifications.New(PUSHER_BEAMS_INSTANCE, PUSHER_BEAMS_SECRET)
		if err != nil {
			fmt.Fprintf(w, "Couldn't create Beams client: %v", err.Error())
		}

		var InterestData struct {
			Interest string
			Image    string
			Id       string
			Type     string
		}

		InterestData.Interest = "chases-notifications"
		InterestData.Id = AddChaseInput.ID
		InterestData.Image = imageURL
		InterestData.Type = "chase"

		publishRequest := map[string]interface{}{
			"apns": map[string]interface{}{
				"aps": map[string]interface{}{
					"alert": map[string]interface{}{
						"title": AddChaseInput.Name,
						"body":  AddChaseInput.Desc,
					},
					"data": InterestData,
				},
			},
			"fcm": map[string]interface{}{
				"notification": map[string]interface{}{
					"title":       AddChaseInput.Name,
					"body":        AddChaseInput.Desc,
					"imageurl":    imageURL,
					"clickaction": chaseURL,
				},
				"data": InterestData,
			},
			"web": map[string]interface{}{
				"notification": map[string]interface{}{
					"title":    AddChaseInput.Name,
					"body":     AddChaseInput.Desc,
					"imageurl": imageURL,
				},
				"data": InterestData,
			},
		}

		pubId, err := beamsClient.PublishToInterests([]string{"chases-notifications"}, publishRequest)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Publish Id:", pubId)
		}

		var myNotification = NotifyInput{
			Body:      ChaseInput.Desc,
			CreatedAt: ChaseInput.CreatedAt,
			Interest:  "chases-notifications",
			Title:     ChaseInput.Name,
			Type:      "chase",
			Data: myJsonNotifyData{
				Id:    ChaseInput.ID,
				Image: ChaseInput.ImageURL,
			},
		}
		notifyErr := addNotification(myNotification)
		if notifyErr != nil {
			fmt.Println("Problem adding notification: %v", notifyErr)
		}

		// send to discord
		fmt.Println("Sending to discord")
		var webhookData = WebhookData{
			Name:     ChaseInput.Name,
			URL:      "https://chaseapp.tv/chase/" + ChaseInput.ID,
			ImageURL: imageURL,
			Desc:     ChaseInput.Desc,
		}
		webhookErr := sendWebhook(webhookData)
		if webhookErr != nil {
			fmt.Println("Problem sending webhook: %v", webhookErr)
		}

	} else {
		fmt.Println("Didnt get Live set to true, not sending FCM push..", AddChaseInput)
	}
}

func UpdateNotifications(w http.ResponseWriter, r *http.Request) {

	type Tweet struct {
		Id      string `json:"id,omitempty"`
		Image   string `json:"image"`
		Tweetid string `json:"tweet_id,omitempty"`
	}

	type myData struct {
		TweetData string `json:"tweetData,omitempty"`
	}

	var d struct {
		Body      string    `json:"body"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		Interest  string    `json:"interest"`
		Title     string    `json:"title"`
		Type      string    `json:"type"`
		Data      myData    `json:"data"`
	}

	d.CreatedAt = time.Now()

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprintf(w, "Error, body must required fields! %v", err)
		log.Printf("Body missing fields or bad data: %v", err)
		return
	}

	fmt.Fprint(w, d)

	if d.Body == "" {
		fmt.Fprint(w, "Must supply body!")
		return
	}

	if d.Interest == "" {
		fmt.Fprint(w, "Must supply interest!")
		return
	}

	if d.Title == "" {
		fmt.Fprint(w, "Must supply title!")
		return
	}

	fmt.Fprint(w, d)
	id, _ := uuid.NewUUID()

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: projectID}

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	defer client.Close()

	myNotifyRes, addNotifyErr := client.Collection("notifications").Doc(id.String()).Set(ctx, d)
	if addNotifyErr != nil {
		// Handle any errors in an appropriate way, such as returning them
		log.Printf("An error has occurred: %v Response: %v", addNotifyErr, myNotifyRes)
	}
	log.Printf("Added document %v to notifications", id.String())
	w.WriteHeader(http.StatusOK)

}

type InterestData struct {
	Interest string
	Image    string
	Id       string
}

type NotificationData struct {
	Id    string
	Image string
	Title string
	Body  string
	Type  string
}

// You must start object keys with an Uppercase
var NotifyJson struct {
	Id       string `json:"id"`
	Interest string `json:"interest"`
	Image    string `json:"image"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	Type     string `json:"type"`
}

func PushNotification(w http.ResponseWriter, r *http.Request) {

	// pre flight stuff
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewDecoder(r.Body).Decode(&NotifyJson); err != nil {
		fmt.Fprint(w, "Error, body must contain id, interest, image, title, and body")
		fmt.Printf("Body missing some fields")
		return
	}
	fmt.Println(NotifyJson.Id)
	fmt.Println(NotifyJson.Interest)
	var interestData InterestData
	var notifyData NotificationData

	// TODO: Do we need to update this for chase-app/flutter#174 ?
	notifyData.Id = NotifyJson.Id
	notifyData.Image = NotifyJson.Image
	notifyData.Title = NotifyJson.Title
	notifyData.Type = NotifyJson.Type
	notifyData.Body = NotifyJson.Body

	interestData.Interest = NotifyJson.Interest
	interestData.Id = NotifyJson.Id
	interestData.Image = NotifyJson.Image

	err := pushNotification(notifyData, interestData)
	if err != nil {
		fmt.Fprintf(w, "Error in pushNotification: %v", err)
		log.Println(interestData)
		log.Println(NotifyJson)
		log.Fatalf("Error in pushNotificaiton: %v", err)
	}
	w.WriteHeader(http.StatusOK)
}

func pushNotification(notify NotificationData, interestData InterestData) error {
	fmt.Printf("bar")
	fmt.Println(notify.Title)
	beamsClient, err := pushnotifications.New(PUSHER_BEAMS_INSTANCE, PUSHER_BEAMS_SECRET)
	if err != nil {
		log.Printf("Couldn't create Beams client: %v", err.Error())
		return err
	}

	publishRequest := map[string]interface{}{
		"apns": map[string]interface{}{
			"aps": map[string]interface{}{
				"alert": map[string]interface{}{
					"title": notify.Title,
					"body":  notify.Body,
				},
				"data": interestData,
			},
		},
		"fcm": map[string]interface{}{
			"notification": map[string]interface{}{
				"title":    notify.Title,
				"body":     notify.Body,
				"imageurl": interestData.Image,
				// "clickaction": chaseURL,
			},
			"data": interestData,
		},
		"web": map[string]interface{}{
			"notification": map[string]interface{}{
				"title":    notify.Title,
				"body":     notify.Body,
				"imageurl": interestData.Image,
			},
			"data": interestData,
		},
	}

	log.Println(interestData)
	pubId, err := beamsClient.PublishToInterests([]string{interestData.Interest}, publishRequest)
	if err != nil {
		log.Println(err)
		return err
	} else {
		log.Println("Publish Id:", pubId)
	}

	/*
		var myFoo NotifyInput
		myFoo.body = notify.Body
		myFoo.title = notify.Title
		myFoo.
	*/

	var myNotification = NotifyInput{
		Body:      notify.Body,
		CreatedAt: time.Now(),
		Interest:  interestData.Interest,
		Title:     notify.Title,
		Type:      notify.Type,
		Data: myJsonNotifyData{
			Id:    notify.Id,
			Image: notify.Image,
		},
	}

	// add to firebase
	err = addNotification(myNotification)
	if err != nil {
		return err
	}
	return nil
}

type WebhookData struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Desc     string `json:"desc"`
	ImageURL string `json:"imageurl"`
	Hook     string `json:"hook"`
}

func SendWebhook(w http.ResponseWriter, r *http.Request) {
	var d WebhookData

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error, body is missing fields!")
		return
	}

	fmt.Fprint(w, d)

	if d.Name == "" {
		fmt.Fprint(w, "Must supply Name - chase, firehose!")
		return
	}
	err := sendWebhook(d)
	if err != nil {
		log.Fatalf("Problem with webhook: %v", err)
	}
}

func sendWebhook(webhookData WebhookData) error {
	postBody, _ := json.Marshal(map[string]string{
		"name":     webhookData.Name,
		"url":      webhookData.URL,
		"desc":     webhookData.Desc,
		"imageurl": webhookData.ImageURL,
		"hook":     webhookData.Hook,
	})
	responseBody := bytes.NewBuffer(postBody)
	var webhookUrl = os.Getenv("DISCORD_WEBHOOK_API_URL")
	resp, err := http.Post(webhookUrl, "application/json", responseBody)

	if err != nil {
		println("Error handling request: #{err}")
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	//Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("sendWebHook response Body ", string(body))
	return nil
}

// addNotification adds and saves notifications into the 'notifications' collection on firebase
func addNotification(notifyData NotifyInput) error {
	id, _ := uuid.NewUUID()

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: projectID}

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
		return err
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
		return err
	}

	defer client.Close()

	myNotifyRes, addNotifyErr := client.Collection("notifications").Doc(id.String()).Set(ctx, notifyData)
	if addNotifyErr != nil {
		// Handle any errors in an appropriate way, such as returning them
		log.Printf("An error has occurred: %v Response: %v", addNotifyErr, myNotifyRes)
		return addNotifyErr
	}
	return nil
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

// DeleteUser deletes a user
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	auth, err := app.Auth(context.Background())

	if APIKEY != r.Header.Get("X-ApiKey") {
		fmt.Fprint(w, "Invalid or missing X-ApiKey")
		return
	}

	var d struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error, body must contain ID!")
		return
	}

	if d.ID == "" {
		fmt.Fprint(w, "Must supply ID!")
		return
	}

	// init requisite clients
	// streamClient := stream.New(getStream_API_KEY, getStream_API_SECRET)
	// beamsClient, err := pushnotifications.New(PUSHER_BEAMS_INSTANCE, PUSHER_BEAMS_SECRET)

	// delete from beams, then streams, then finally - firebase auth
	// err = beamsClient.DeleteUser(d.ID)
	// err = streamClient.Users().Delete(d.ID)
	// err = app.DeleteUser(context.Background(), d.ID)
	authErr := auth.DeleteUser(context.Background(), d.ID)

	if authErr != nil {
		fmt.Fprintln(w, "Could not delete user: ", authErr.Error())
	}
	fmt.Fprint(w, "User deleted successfully")
	log.Printf("Successfully deleted user: %s\n", d.ID)
}

// SetAdmin if provided with the correct API key, will grant
// a UID admin in firebase authentication. This is used for chaseapp-crud.
func SetAdmin(w http.ResponseWriter, r *http.Request) {
	conf := &firebase.Config{ProjectID: projectID}

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

	// send to discord
	fmt.Println("Sending to discord")
	var webhookData = WebhookData{
		Name: fmt.Sprintf("⚠️ Set UID %v to Admin", d.UID),
	}

	// Send a message to the #chaseapp-admins discord channel
	// to notify that a new admin has been added
	webhookErr := sendWebhook(webhookData)
	if webhookErr != nil {
		fmt.Println("Problem sending webhook: %v", webhookErr)
	}

}

// UpdateChase updates a chase
func UpdateChase(w http.ResponseWriter, r *http.Request) {

	// pre flight stuff
	if r.Method == http.MethodOptions {
		log.Println("Preflight request")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	if APIKEY != r.Header.Get("X-ApiKey") {
		fmt.Println("Invalid or missing X-ApiKey", r.Header.Get("X-ApiKey"))
		fmt.Fprint(w, "Invalid or missing X-ApiKey")
		return
	}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		fmt.Println("firebase.NewApp: %v", err)
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

	fmt.Println("ChaseInput", ChaseInput)

	if ChaseInput.ID == "" {
		fmt.Fprint(w, "Must supply ID!")
		return
	}

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

	if ChaseInput.Notified != true {
		if ChaseInput.Notified != false {
			fmt.Fprint(w, "Notified value must be true or false")
			return
		}
	}
	_, notAddErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"Notified": ChaseInput.Notified}, firestore.MergeAll)
	if notAddErr != nil {
		_, notAddErr := json.Marshal(notAddErr)
		fmt.Fprintf(w, "Error updating chase Notified: %v", notAddErr)
	}

	if ChaseInput.ImageURL != "" {
		_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"ImageURL": ChaseInput.ImageURL}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating chase ImageURL: %v", addErr)
		}
	}

	// TODO: this seems to be doing something weird..
	if ChaseInput.Wheels.W1 != "" {
		// TODO: figure out why we aren't doing this on W2-W4..
		var mappedData = transform.ToFirestoreMap(ChaseInput.Wheels)
		_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, mappedData, firestore.MergeAll)
		// _, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"Wheels": ChaseInput.Wheels}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating chase ChaseWheels: %v", addErr)
		}
	}

	if ChaseInput.Wheels.W2 != "" {
		// var mappedData = transform.ToFirestoreMap(ChaseInput.Wheels)
		// _, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, mappedData, firestore.MergeAll)
		_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"Wheels": ChaseInput.Wheels}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating chase ChaseWheels: %v", addErr)
		}
	}
	if ChaseInput.Wheels.W3 != "" {
		// var mappedData = transform.ToFirestoreMap(ChaseInput.Wheels)
		// _, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, mappedData, firestore.MergeAll)
		_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"Wheels": ChaseInput.Wheels}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating chase ChaseWheels: %v", addErr)
		}
	}
	if ChaseInput.Wheels.W4 != "" {
		//var mappedData = transform.ToFirestoreMap(ChaseInput.Wheels)
		//_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, mappedData, firestore.MergeAll)
		_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"Wheels": ChaseInput.Wheels}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating chase ChaseWheels: %v", addErr)
		}
	}

	if !ChaseInput.EndedAt.IsZero() {
		_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"EndedAt": ChaseInput.EndedAt}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating chase EndedAt: %v", addErr)
		}
	}

	if ChaseInput.Tags != "" {
		tagSlice := strings.Split(ChaseInput.Tags, ",")
		_, addErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"Tags": tagSlice}, firestore.MergeAll)
		if addErr != nil {
			_, addErr := json.Marshal(addErr)
			fmt.Fprintf(w, "Error updating chase tags: %v", addErr)
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

		// compare Chase.Networks to ChaseInput.Networks, print the ones that are different
		for _, network := range ChaseInput.Networks {
			// see if network exists in Chase.Networks
			if !containsNetwork(networkArray, network) {
				fmt.Fprintln(w, network.URL)

				// Notify the channel that a new link has been added
				var channelSettings ChannelSettings
				channelSettings.ChannelType = "livestream"
				channelSettings.Id = chaseData.ID
				channelSettings.CustomName = chaseData.Name

				var channelMessage = "🎉 New link added " + network.URL + " 📺"
				sendMessageErr := sendMessage(channelMessage, channelSettings.Id, "leku")
				if sendMessageErr != nil {
					fmt.Println("Problem sending message: %v", err)
				}
			}
		}

		networkArray = append(networkArray, ChaseInput.Networks...)
		networkArray = ChaseInput.Networks
		chaseData.Networks = ChaseInput.Networks
		_, urlErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"Networks": networkArray}, firestore.MergeAll)
		if urlErr != nil {
			fmt.Fprintf(w, "Error adding Networks to chase: %v", urlErr)
			return
		}

		request := GetMP4Link{
			ChaseID: ChaseInput.ID,
		}
		jsonValue, _ := json.Marshal(request)
		resp, err := http.Post(GetStreamLinkURL(), "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Fprintf(w, "Error calling GetMP4Link: %v", err)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body: ", err)
			fmt.Fprintf(w, "Error reading GetMP4Link response: %v", err)
			return
		}
		fmt.Println("GetMP4Link response: ", string(body))
	}

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

	pusherErr := pusherClient.Trigger("chases", "updates", ChaseInput)
	if pusherErr != nil {
		fmt.Fprintf(w, "Error triggering pusher event: %v", pusherErr)
	}

	// using a custom writer from github.com/mfreeman451/golang/common/writers
	jw := writers.NewMessageWriter(ChaseInput.ID)
	jsonString, err := jw.JSONString()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		fmt.Println(err.Error())
		return
	}

	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_, err = w.Write([]byte(jsonString))
	if err != nil {
		return
	}
}

/* Helper functions */

func containsNetwork(s []Networks, e Networks) bool {
	for _, a := range s {
		if a.URL == e.URL {
			return true
		}
	}
	return false
}
