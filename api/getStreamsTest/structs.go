package main

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
