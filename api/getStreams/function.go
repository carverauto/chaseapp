package p

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// GCLOUD_PROJECT is automatically set by the Cloud Functions runtime.
var projectID = os.Getenv("GCLOUD_PROJECT")

// client is a Firestore client, reused between function invocations.
var client *firestore.Client

type NBCMetaData struct {
	SpayEmail                    string        `json:"spay_email"`
	SyndicatedID                 string        `json:"syndicated_id"`
	ShouldNationalize            bool          `json:"should_nationalize"`
	NbcFirstPublishedDate        string        `json:"nbc_first_published_date"`
	NbcPreventModifiedUpdate     bool          `json:"nbc_prevent_modified_update"`
	NbcLinkoutTitle              string        `json:"nbc_linkout_title"`
	NbcLinkoutURL                string        `json:"nbc_linkout_url"`
	NbcLinkoutExcerptLink        string        `json:"nbc_linkout_excerpt_link"`
	NbcPrimaryCategoryID         int           `json:"nbc_primary_category_id"`
	NbcPrimaryTagID              int           `json:"nbc_primary_tag_id"`
	NbcSocialPageTitle           string        `json:"nbc_social_page_title"`
	NbcPageTitle                 string        `json:"nbc_page_title"`
	NbcSeoDescription            string        `json:"nbc_seo_description"`
	NbcSocialDescription         string        `json:"nbc_social_description"`
	NbcNationalHeadline          string        `json:"nbc_national_headline"`
	Entities                     string        `json:"entities"`
	NbcCanonicalMarket           string        `json:"nbc_canonical_market"`
	NbcCanonicalURL              string        `json:"nbc_canonical_url"`
	NbcHideFromRecirculation     string        `json:"nbc_hide_from_recirculation"`
	NbcHideFromSearch            string        `json:"nbc_hide_from_search"`
	NbcSyndicationCategory       string        `json:"nbc_syndication_category"`
	NbcNationalCollections       []interface{} `json:"nbc_national_collections"`
	Nopreroll                    bool          `json:"nopreroll"`
	VideoCopyright               string        `json:"video_copyright"`
	SstSourceID                  string        `json:"sst_source_id"`
	NbcDisableLegacyEmbeds       bool          `json:"nbc_disable_legacy_embeds"`
	NbcFeedLoaderPreventUpdates  string        `json:"nbc_feed_loader_prevent_updates"`
	PostVideoPrompt              string        `json:"post_video_prompt"`
	NbcCheckoutExperienceEnabled bool          `json:"nbc_checkout_experience_enabled"`
	RsnDisableInsider            bool          `json:"rsn_disable_insider"`
	DesktopFlashPid              string        `json:"desktopFlashPid"`
	MpxDfxpURL                   string        `json:"mpx_dfxp_url"`
	Mp4URL                       string        `json:"mp4_url"`
	MpxDownloadPid               string        `json:"mpx_download_pid"`
	MpxDownloadPidHigh           string        `json:"mpx_download_pid_high"`
	MpxHighID                    string        `json:"mpx_high_id"`
	MpxHighURL                   string        `json:"mpx_high_url"`
	MpxIsLivestream              string        `json:"mpx_is_livestream"`
	MpxM3Upid                    string        `json:"mpx_m3upid"`
	MpxPid                       string        `json:"mpx_pid"`
	MpxThumbnailURL              string        `json:"mpx_thumbnail_url"`
	PidStreamingWebHigh          string        `json:"pid_streaming_web_high"`
	PidStreamingWebMedium        string        `json:"pid_streaming_web_medium"`
	Subtitle                     string        `json:"subtitle"`
	DfxpURL                      string        `json:"dfxp_url"`
	VideoCallLetter              string        `json:"video_call_letter"`
	VideoCaptions                string        `json:"video_captions"`
	VideoID                      string        `json:"video_id"`
	VideoLength                  string        `json:"video_length"`
	VideoProvider                string        `json:"video_provider"`
	ShortVideoExcerpt            string        `json:"short_video_excerpt"`
	MpxDownloadPidMobileLow      string        `json:"mpx_download_pid_mobile_low"`
	PidStreamingWebMobileLow     string        `json:"pid_streaming_web_mobile_low"`
	MpxDownloadPidMobileStandard string        `json:"mpx_download_pid_mobile_standard"`
	PidStreamingMobileStandard   string        `json:"pid_streaming_mobile_standard"`
	MediaPid                     string        `json:"media_pid"`
	AlleypackScheduleUnpublish   string        `json:"alleypack_schedule_unpublish"`
	FeedRemoteID                 string        `json:"feed_remote_id"`
	FeedThumbnailURL             string        `json:"feed_thumbnail_url"`
}

var GetStreamLinkRequest struct {
	ChaseID string `json:"chase_id"`
}

type Stream struct {
	URL  string `json:"url"`
	Tier int    `json:"tier"`
}

type Networks struct {
	URL     string   `json:"url"`
	Streams []Stream `json:"streams,omitempty"`
	Name    string   `json:"name"`
	Tier    int      `json:"tier"`
	Logo    string   `json:"logo"`
	Other   string   `json:"other"`
}

type Chase struct {
	ID       string     `json:"id,omitempty"`
	Networks []Networks `json:"networks"`
}

var seenBefore = make(map[string]bool)

func removeDuplicates(elements []Stream) []Stream {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []Stream{}

	for v := range elements {
		if encountered[elements[v].URL] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v].URL] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func updateChase(chase Chase) error {
	ctx := context.Background()

	// range through the networks
	for i, network := range chase.Networks {
		// fix the streams URLs, remove \ slashes
		for j, stream := range network.Streams {

			stream.URL = strings.Replace(stream.URL, `\`, "", -1)
			chase.Networks[i].Streams[j] = stream

			// remove duplicate stream.URLs
			chase.Networks[i].Streams = removeDuplicates(chase.Networks[i].Streams)
		}
	}

	// _, urlErr := client.Collection("chases").Doc(ChaseInput.ID).Set(ctx, map[string]interface{}{"Networks": networkArray}, firestore.MergeAll)
	_, err := client.Collection("chases").Doc(chase.ID).Set(ctx, map[string]interface{}{"Networks": chase.Networks}, firestore.MergeAll)
	if err != nil {
		return err
	}
	return nil
}

func GetStreams(w http.ResponseWriter, r *http.Request) {
	conf := &firebase.Config{ProjectID: projectID, DatabaseURL: "https://chaseapp-8459b.firebaseio.com"}

	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	if err := json.NewDecoder(r.Body).Decode(&GetStreamLinkRequest); err != nil {
		fmt.Fprint(w, "Error, body must contain chase_id")
		return
	}

	// query firestore for chase_id
	chaseRef := client.Collection("chases").Doc(GetStreamLinkRequest.ChaseID)
	chaseDoc, err := chaseRef.Get(ctx)
	if err != nil {
		fmt.Fprint(w, "Error, chase_id not found")
		return
	}
	var MyChase Chase
	// convert chaseDoc to Chase struct
	fmt.Println(chaseDoc.Data())
	chaseDoc.DataTo(&MyChase)
	fmt.Println(MyChase)
	// range through the Networks struct and get the URL
	// var URLs []string
	/*
		 for _, network := range MyChase.Networks {
			URLs = append(URLs, network.URL)
		}
	*/
	// scraper stuff
	c := colly.NewCollector()

	// iterate through the URLs and get the MP4 link
	for nIdx, network := range MyChase.Networks {
		fmt.Println("Network: ", network.Name)
		// case/switch for url using strings.Contains
		switch {
		case strings.Contains(network.URL, "nbclosangeles.com"):
			var streams []Stream
			fmt.Println("NBC Los Angeles")
			// Find and visit all links
			c.OnHTML("div", func(e *colly.HTMLElement) {
				rawMeta := e.Attr("data-meta")
				if rawMeta != "" {
					var meta NBCMetaData
					var stream Stream
					err := json.Unmarshal([]byte(rawMeta), &meta)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("ChaseID: ", MyChase.ID, "MP4URL:", meta.Mp4URL)
					// update chase with MP4URL
					// MyChase.Networks[nIdx].MP4URL = meta.Mp4URL
					stream.URL = meta.Mp4URL
					streams = append(streams, stream)
					MyChase.Networks[nIdx].Streams = streams
				}
			})
			fmt.Println("Visiting: ", network.URL)
			err := c.Visit(network.URL)
			if err != nil {
				log.Println(err)
			}
			break
		case strings.Contains(network.URL, "abc7.com"):
			var streams []Stream
			fmt.Println("ABC7")
			// Find and visit all links
			c.OnHTML("script", func(e *colly.HTMLElement) {
				// window['__abcotv__']
				// look for "m3u8":"https://content.uplynk.com/channel/ext/2118d9222a87420ab69223af9cfa0a0f/kabc_24x7_news.m3u8?ad._v=2&ad.preroll=0&ad.fill_slate=1&ad.ametr=1&ad.vid=otv-11316941&rays=ihgfedc" and extract URL
				re, err := regexp.Compile(`"m3u8":"(.*?)"`)
				if err != nil {
					fmt.Println(err)
				}
				matches := re.FindStringSubmatch(e.Text)
				if len(matches) > 0 {
					var stream Stream
					fmt.Println("ChaseID: ", MyChase.ID, "MP4URL:", matches[1])
					// update chase with MP4URL
					// MyChase.Networks[nIdx].MP4URL = matches[1]
					stream.URL = matches[1]
					streams = append(streams, stream)
					MyChase.Networks[nIdx].Streams = streams
				}
			})
			err := c.Visit(network.URL)
			if err != nil {
				log.Println(err)
			}
			break
		case strings.Contains(network.URL, "cbsnews.com"):
			var streams []Stream
			fmt.Println("CBS")
			c.OnHTML("script", func(e *colly.HTMLElement) {
				re, err := regexp.Compile(`"video":"(.*?)"`)
				if err != nil {
					fmt.Println(err)
				}
				matches := re.FindStringSubmatch(e.Text)
				if len(matches) > 0 {
					var stream Stream
					fmt.Println("ChaseID: ", MyChase.ID, "MP4URL:", matches[1])
					// update chase with MP4URL
					//MyChase.Networks[nIdx].MP4URL = matches[1]
					stream.URL = matches[1]
					streams = append(streams, stream)
					MyChase.Networks[nIdx].Streams = streams
				} else {
					// try something else
					re, err := regexp.Compile(`"contentUrl":"(.*?)"`)
					if err != nil {
						fmt.Println(err)
					}
					matches2 := re.FindStringSubmatch(e.Text)
					if len(matches2) > 0 {
						var stream Stream
						fmt.Println("ChaseID: ", MyChase.ID, "MP4URL:", matches2[1])
						// update chase with MP4URL
						//MyChase.Networks[nIdx].MP4URL = matches[1]
						stream.URL = matches2[1]
						if !seenBefore[stream.URL] {
							streams = append(streams, stream)
							MyChase.Networks[nIdx].Streams = streams
							seenBefore[MyChase.ID] = true
						}

					} else {
						fmt.Println("No matches")
					}
				}
			})
			err := c.Visit(network.URL)
			if err != nil {
				log.Println(err)
			}
			break
		}

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
		})

	}

	log.Println("Updating chase with MP4URLs")
	uErr := updateChase(MyChase)
	if uErr != nil {
		log.Println(uErr)
	}
	fmt.Fprintln(w, "Done")
}
