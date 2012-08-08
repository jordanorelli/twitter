package twitter

type User struct {
	ContributorsEnabled            Nbool     `json:"contributors_enabled" bson:"contributors_enabled"`
	CreatedAt                      Timestamp `json:"created_at" bson:"created_at"`
	DefaultProfile                 Nbool     `json:"default_profile" bson:"default_profile"`
	DefaultProfileImage            Nbool     `json:"default_profile_image" bson:"default_profile_image"`
	Description                    Nstring   `json:"description" bson:"description"`
	FavouritesCount                Nint      `json:"favourites_count" bson:"favourites_count"`
	FollowRequestSent              Nbool     `json:"follow_request_sent" bson:"follow_request_sent"`
	FollowersCount                 Nint      `json:"followers_count" bson:"followers_count"`
	Following                      Nbool     `json:"following" bson:"following"`
	FriendsCount                   Nint      `json:"friends_count" bson:"friends_count"`
	GeoEnabled                     Nbool     `json:"geo_enabled" bson:"geo_enabled"`
	Id                             Snowflake `json:"id" bson:"id"`
	IdStr                          Nstring   `json:"id_str" bson:"id_str"`
	IsTranslator                   Nbool     `json:"is_translator" bson:"is_translator"`
	Lang                           Nstring   `json:"lang" bson:"lang"`
	ListedCount                    Nint      `json:"listed_count" bson:"listed_count"`
	Location                       Nstring   `json:"location" bson:"location"`
	Name                           Nstring   `json:"name" bson:"name"`
	ProfileBackgroundColor         Nstring   `json:"profile_background_color" bson:"profile_background_color"`
	ProfileImageUrl                Nstring   `json:"profile_image_url" bson:"profile_image_url"`
	ProfileImageUrlHttps           Nstring   `json:"profile_image_url_https" bson:"profile_image_url_https"`
	ProfileBackgroundImageUrl      Nstring   `json:"profile_background_image_url" bson:"profile_background_image_url`
	ProfileBackgroundImageUrlHttps Nstring   `json:"profile_background_image_url_https" bson:"profile_background_image_url_https"`
	ProfileLinkColor               Nstring   `json:"profile_link_color" bson:"profile_link_color"`
	ProfileSidebarBorderColor      Nstring   `json:"profile_sidebar_border_color" bson:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        Nstring   `json:"profile_sidebar_fill_color" bson:"profile_sidebar_fill_color"`
	ProfileTextColor               Nstring   `json:"profile_text_color" bson:"profile_text_color"`
	ProfileUseBackgroundImage      Nbool     `json:"profile_use_background_image" bson:"profile_use_background_image"`
	Protected                      Nbool     `json:"protected" bson:"protected"`
	ScreenName                     Nstring   `json:"screen_name" bson:"screen_name"`
	ShowAllInlineMedia             Nbool     `json:"show_all_inline_media" bson:"show_all_inline_media"`
	StatusesCount                  Nint      `json:"statuses_count" bson:"statuses_count"`
	TimeZome                       Nstring   `json:"time_zone" bson:"time_zone"`
	Url                            Nstring   `json:"url" bson:"url"`
	UTCOffset                      Nint      `json:"utc_offset" bson:"utc_offset"`
	Verified                       Nbool     `json:"verified" bson:"verified"`
	/* "notifications": null */
}

type Friends struct {
	Friends []uint64 `json:"friends"`
}

type Credentials struct {
	Token  string `json:"oauth_token" bson:"oauth_token"`
	Secret string `json:"oauth_token_secret" bson:"oauth_token_secret"`
}
