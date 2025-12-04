package main

import (
	"cloud.google.com/go/firestore"
	transcoder "cloud.google.com/go/video/transcoder/apiv1"
	"cloud.google.com/go/video/transcoder/apiv1/transcoderpb"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/joho/godotenv"
	"github.com/octoper/go-ray"
	"github.com/peterbourgon/ff/v3"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

// GCLOUD_PROJECT is automatically set by the Cloud Functions runtime.
var projectID = os.Getenv("GCLOUD_PROJECT")

// client is a Firestore client, reused between function invocations.
var client *firestore.Client

// createJobWithSetNumberImagesSpritesheet creates a job from an ad-hoc configuration and generates
// two spritesheets from the input video. Each spritesheet contains a set number of images.
func createJobWithSetNumberImagesSpritesheet(w io.Writer, projectID string, location string, inputURI string, outputURI string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// inputURI := "gs://my-bucket/my-video-file"
	// outputURI := "gs://my-bucket/my-output-folder/"
	ctx := context.Background()
	client, err := transcoder.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	req := &transcoderpb.CreateJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Job: &transcoderpb.Job{
			InputUri:  inputURI,
			OutputUri: outputURI,
			JobConfig: &transcoderpb.Job_Config{
				Config: &transcoderpb.JobConfig{
					ElementaryStreams: []*transcoderpb.ElementaryStream{
						{
							Key: "video_stream0",
							ElementaryStream: &transcoderpb.ElementaryStream_VideoStream{
								VideoStream: &transcoderpb.VideoStream{
									CodecSettings: &transcoderpb.VideoStream_H264{
										H264: &transcoderpb.VideoStream_H264CodecSettings{
											BitrateBps:   550000,
											FrameRate:    60,
											HeightPixels: 360,
											WidthPixels:  640,
										},
									},
								},
							},
						},
						{
							Key: "audio_stream0",
							ElementaryStream: &transcoderpb.ElementaryStream_AudioStream{
								AudioStream: &transcoderpb.AudioStream{
									Codec:      "aac",
									BitrateBps: 64000,
								},
							},
						},
					},
					MuxStreams: []*transcoderpb.MuxStream{
						{
							Key:               "sd",
							Container:         "mp4",
							ElementaryStreams: []string{"video_stream0", "audio_stream0"},
						},
					},
					SpriteSheets: []*transcoderpb.SpriteSheet{
						{
							FilePrefix:         "small-sprite-sheet",
							SpriteWidthPixels:  64,
							SpriteHeightPixels: 32,
							ColumnCount:        10,
							RowCount:           10,
							ExtractionStrategy: &transcoderpb.SpriteSheet_TotalCount{
								TotalCount: 100,
							},
						},
						{
							FilePrefix:         "large-sprite-sheet",
							SpriteWidthPixels:  128,
							SpriteHeightPixels: 72,
							ColumnCount:        10,
							RowCount:           10,
							ExtractionStrategy: &transcoderpb.SpriteSheet_TotalCount{
								TotalCount: 100,
							},
						},
					},
				},
			},
		},
	}
	// Creates the job. Jobs take a variable amount of time to run. You can query for the job state.
	// See https://cloud.google.com/transcoder/docs/how-to/jobs#check_job_status for more info.
	response, err := client.CreateJob(ctx, req)
	if err != nil {
		return fmt.Errorf("createJobWithSetNumberImagesSpritesheet: %v", err)
	}

	fmt.Fprintf(w, "Job: %v", response.GetName())
	return nil
}

// createJobFromPreset creates a job based on a given preset template. See
// https://cloud.google.com/transcoder/docs/how-to/jobs#create_jobs_presets
// for more information.
func createJobFromPreset(w io.Writer, projectID string, location string, inputURI string, outputURI string, preset string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// inputURI := "gs://my-bucket/my-video-file"
	// outputURI := "gs://my-bucket/my-output-folder/"
	// preset := "preset/web-hd"
	ctx := context.Background()
	client, err := transcoder.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	req := &transcoderpb.CreateJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Job: &transcoderpb.Job{
			InputUri:  inputURI,
			OutputUri: outputURI,
			JobConfig: &transcoderpb.Job_TemplateId{
				TemplateId: preset,
			},
		},
	}
	// Creates the job, Jobs take a variable amount of time to run.
	// You can query for the job state.
	response, err := client.CreateJob(ctx, req)
	if err != nil {
		return fmt.Errorf("createJobFromPreset: %v", err)
	}

	fmt.Fprintf(w, "Job: %v", response.GetName())
	return nil
}

// getJobState gets the state for a previously-created job. See
// https://cloud.google.com/transcoder/docs/how-to/jobs#check_job_status for
// more information.
func getJobState(w io.Writer, projectID string, location string, jobID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// jobID := "my-job-id"
	ctx := context.Background()
	client, err := transcoder.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	req := &transcoderpb.GetJobRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectID, location, jobID),
	}

	response, err := client.GetJob(ctx, req)
	if err != nil {
		return fmt.Errorf("GetJob: %v", err)
	}
	fmt.Fprintf(w, "Job state: %v\n----\nJob failure reason:%v\n", response.State, response.Error)
	return nil
}

// createJobWithStaticOverlay creates a job based on a given configuration that
// includes a static overlay. See
// https://cloud.google.com/transcoder/docs/how-to/create-overlays#create-static-overlay
// for more information.
func createJobWithStaticOverlay(w io.Writer, projectID string, location string, inputURI string, overlayImageURI string, outputURI string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// inputURI := "gs://my-bucket/my-video-file"
	// overlayImageURI := "gs://my-bucket/my-overlay-image-file" - Must be a JPEG
	// outputURI := "gs://my-bucket/my-output-folder/"
	ctx := context.Background()
	client, err := transcoder.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	req := &transcoderpb.CreateJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Job: &transcoderpb.Job{
			InputUri:  inputURI,
			OutputUri: outputURI,
			JobConfig: &transcoderpb.Job_Config{
				Config: &transcoderpb.JobConfig{
					ElementaryStreams: []*transcoderpb.ElementaryStream{
						{
							Key: "video_stream0",
							ElementaryStream: &transcoderpb.ElementaryStream_VideoStream{
								VideoStream: &transcoderpb.VideoStream{
									CodecSettings: &transcoderpb.VideoStream_H264{
										H264: &transcoderpb.VideoStream_H264CodecSettings{
											BitrateBps:   550000,
											FrameRate:    60,
											HeightPixels: 360,
											WidthPixels:  640,
										},
									},
								},
							},
						},
						{
							Key: "audio_stream0",
							ElementaryStream: &transcoderpb.ElementaryStream_AudioStream{
								AudioStream: &transcoderpb.AudioStream{
									Codec:      "aac",
									BitrateBps: 64000,
								},
							},
						},
					},
					MuxStreams: []*transcoderpb.MuxStream{
						{
							Key:               "sd",
							Container:         "mp4",
							ElementaryStreams: []string{"video_stream0", "audio_stream0"},
						},
					},
					Overlays: []*transcoderpb.Overlay{
						{
							Image: &transcoderpb.Overlay_Image{
								Uri: overlayImageURI,
								Resolution: &transcoderpb.Overlay_NormalizedCoordinate{
									X: 1,
									Y: 0.5,
								},
								Alpha: 1,
							},
							Animations: []*transcoderpb.Overlay_Animation{
								{
									AnimationType: &transcoderpb.Overlay_Animation_AnimationStatic{
										AnimationStatic: &transcoderpb.Overlay_AnimationStatic{
											Xy: &transcoderpb.Overlay_NormalizedCoordinate{
												X: 0,
												Y: 0,
											},
											StartTimeOffset: &duration.Duration{
												Seconds: 0,
											},
										},
									},
								},

								{
									AnimationType: &transcoderpb.Overlay_Animation_AnimationEnd{
										AnimationEnd: &transcoderpb.Overlay_AnimationEnd{
											StartTimeOffset: &duration.Duration{
												Seconds: 10,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	// Creates the job. Jobs take a variable amount of time to run.
	// You can query for the job state; see getJob() in get_job.go.
	response, err := client.CreateJob(ctx, req)
	if err != nil {
		return fmt.Errorf("createJobWithStaticOverlay: %v", err)
	}

	fmt.Fprintf(w, "Job: %v", response.GetName())
	return nil
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

func ChaseId() string {
	// no default
	return ""
}

func main() {
	fs := flag.NewFlagSet("getStreams", flag.ExitOnError)
	var (
		chaseId = fs.String("id", ChaseId(), "Chase ID")
		debug   = fs.Bool("ray", false, "Ray")
	)
	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVars()); err != nil {
		log.Fatal(err)
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

	// query firestore for chase_id
	if *chaseId == "" {
		log.Fatalf("Chase ID not set")
	}
	chaseRef := client.Collection("chases").Doc(*chaseId)
	chaseDoc, err := chaseRef.Get(ctx)
	if err != nil {
		fmt.Println("Error, chase_id not found: ", err)
		return
	}
	var MyChase Chase
	// convert chaseDoc to Chase struct
	dErr := chaseDoc.DataTo(&MyChase)
	if dErr != nil {
		fmt.Println("Error, chaseDoc.DataTo: ", dErr)
		return
	}
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
					ray.Ray("Stream: ", stream)
					ray.Ray("Meta: ", meta)
					fmt.Println("ChaseID: ", MyChase.ID, "MP4URL:", meta.Mp4URL)
					// update chase with MP4URL
					// MyChase.Networks[nIdx].MP4URL = meta.Mp4URL
					stream.URL = meta.Mp4URL
					streams = append(streams, stream)
					MyChase.Networks[nIdx].Streams = streams
				}
			})
			c.OnHTML("div > div.tpVideo > video:nth-child(3)", func(e *colly.HTMLElement) {
				fmt.Println("Video: ", e.Attr("src"))
				video := e.Attr("src")
				if video != "" {
					var meta NBCMetaData
					var stream Stream
					err := json.Unmarshal([]byte(video), &meta)
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
			// #nbc-mpx-video-22437380_235-0 > div > div.tpVideo > video:nth-child(3)
			/*
				c.OnHTML("div > div.tpVideo > video:nth-child(3)", func(e *colly.HTMLElement) {
					fmt.Println("eText ", e.Text)
					fmt.Println("Src: ", e.Attr("src"))
					rawSrc := e.Attr("src")
					fmt.Println("ChaseID: ", MyChase.ID, "MP4URL:", rawSrc)
						if rawSrc != "" {
							var stream Stream
							fmt.Println("ChaseID: ", MyChase.ID, "MP4URL:", rawSrc)
							// update chase with MP4URL
							// MyChase.Networks[nIdx].MP4URL = rawSrc
							stream.URL = rawSrc
							streams = append(streams, stream)
							MyChase.Networks[nIdx].Streams = streams
						}

				})

			*/
			if network.URL != "" {
				fmt.Println("URL: ", network.URL)
				err := c.Visit(network.URL)
				if err != nil {
					fmt.Println("Error visiting netowrk URL: ", err)
					return
				}
				fmt.Println("Visiting: ", network.URL)
				err = c.Visit(network.URL)
				if err != nil {
					log.Println(err)
				}
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
	if *debug {
		ray.Ray(MyChase)
	}
	uErr := updateChase(MyChase)
	if uErr != nil {
		log.Println(uErr)
	}
	fmt.Println("Done")
}
