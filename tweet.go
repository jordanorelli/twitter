package twitter

import (
	"encoding/json"
	"labix.org/v2/mgo/bson"
	"strconv"
	"strings"
	"time"
)

// the "geo" block is ignored, based on the understanding that it is
// deprecated.  "contributors" and "place" blocks are ignored.
type Tweet struct {
	Coordinates               *Coordinates `json:"coordinates,omitempty"`
	CreatedAt                 Timestamp    `json:"created_at"`
	Entities                  Entities     `json:"entities"`
	Favorited                 Nbool        `json:"favorited"`
	Id                        Snowflake    `json:"id"`
	IdStr                     Nstring      `json:"id_str"`
	InReplyToScreenName       Nstring      `json:"in_reply_to_screen_name,omitempty"`
	InReplyToStatusId         Snowflake    `json:"in_reply_to_status_id,omitempty"`
	InReplyToStatusIdStr      Nstring      `json:"in_reply_to_status_id_str,omitempty"`
	InReplyToUserId           Snowflake    `json:"in_reply_to_user_id,omitempty"`
	InReplyToUserIdStr        Nstring      `json:"in_reply_to_user_id_str,omitempty"`
	PossiblySensitive         Nbool        `json:"possibly_sensitive"`
	PossiblySensitiveEditable Nbool        `json:"possibly_sensitive_editable"`
	RetweetCount              Nint         `json:"retweet_count"`
	Retweeted                 Nbool        `json:"retweeted"`
	Source                    Nstring      `json:"source,omitempty"`
	Text                      Nstring      `json:"text"`
	Truncated                 Nbool        `json:"truncated"`
	User                      *User        `json:"user"`
	MediaUrl                  Nstring      `json:"media_url,omitempty"`
	MediaUrlHttps             Nstring      `json:"media_url_https,omitempty"`
	Url                       Nstring      `json:"url,omitempty"`
	DisplayUrl                Nstring      `json:"display_url,omitempty"`
	ExpandedUrl               Nstring      `json:"expanded_url,omitempty"`
}

// container for twitter entities.  Used where available; does not get used by
// all twitter endpoints. See here: https://dev.twitter.com/docs/tweet-entities
type Entities struct {
	Hashtags     []Hashtag     `json:"hashtags,omitempty"`
	Urls         []Url         `json:"urls,omitempty"`
	UserMentions []UserMention `json:"user_mentions,omitempty"`
	Media        []Media       `json:"media,omitempty"`
}

// represents geocoordinates found in a tweet.  Not yet fully functional
type Coordinates struct {
	Type        string `json:"type"`
	Coordinates []float32
}

// media items found inside of tweets.  According to the documentation, this is
// only photos for now.  No word on what other types of media will be
// appearing, but video and audio is assumed.
type Media struct {
	Id            uint64          `json:"id"`
	IdStr         string          `json:"id_str"`
	MediaUrl      string          `json:"media_url"`
	MediaUrlHttps string          `json:"media_url_https"`
	Url           string          `json:"url"`
	DisplayUrl    string          `json:"display_url"`
	ExpandedUrl   string          `json:"expanded_url"`
	Sizes         map[string]Size `json:"sizes"`
	Type          string          `json:"type"`
	Indices       []int           `json:"indices"`
}

// size for an individual photo in photo media objects
type Size struct {
	Width  int    `json:"w"`
	Height int    `json:"h"`
	Resize string `json:"resize"`
}

// metadata for an @mention in a tweet
type UserMention struct {
	Id         Snowflake `json:"id"`
	IdStr      Nstring   `json:"id_str"`
	Indices    []int     `json:"indices"`
	Name       Nstring   `json:"name"`
	ScreenName Nstring   `json:"screen_name"`
}

// metadata for a url found in a tweet, including the indices that it appears
// in the tweet.
type Url struct {
	DisplayUrl  Nstring `json:"display_url"`
	ExpandedUrl Nstring `json:"expanded_url"`
	Indices     []int   `json:"indices"`
	Url         Nstring `json:"url"`
}

type Tagstring string

func (t *Tagstring) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*string)(t))
}

func (t Tagstring) GetBSON() (interface{}, error) {
	return strings.ToLower(string(t)), nil
}

// metadata for hashtags found in an individual tweet, including the indices
// where the hashtag appears in its parent tweet.
type Hashtag struct {
	Indices []int     `json:"indices"`
	Text    Tagstring `json:"text"`
}

// specifies how twitter formats time.Time objects.  They're really just
// RubyDates.
type Timestamp time.Time

func (t *Timestamp) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	val, err := time.Parse(time.RubyDate, string(b[1:len(b)-1]))
	*t = Timestamp(val)
	return
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Unix())
}

func (t Timestamp) GetBSON() (interface{}, error) {
	return time.Time(t), nil
}

func (t *Timestamp) SetBSON(raw bson.Raw) error {
	var normal time.Time
	err := raw.Unmarshal(&normal)
	if err == nil {
		*t = Timestamp(normal)
	}
	return err
}

// representation of unique ids on twitter.  These numbers are internally
// uint64s; they should be used with caution when interacting with environments
// that do not support the uint64 datatype, such as JavaScript.  Snowflakes
// always appear alongside a string alternative that should be used instead
// when interacting with clients that are not known to support the uint64
// datatype.
type Snowflake uint64

func (n *Snowflake) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*uint64)(n))
}

func (n Snowflake) GetBSON() (interface{}, error) {
	return strconv.FormatUint(uint64(n), 10), nil
}

func (n *Snowflake) SetBSON(raw bson.Raw) error {
	var s string
	err := raw.Unmarshal(&s)
	if err != nil {
		return err
	}
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	*n = Snowflake(val)
	return nil
}

// nullable string.
type Nstring string

func (n *Nstring) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*string)(n))
}

func (n Nstring) GetBSON() (interface{}, error) {
	return string(n), nil
}

type Nbool bool

func (n *Nbool) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*bool)(n))
}

func (n Nbool) GetBSON() (interface{}, error) {
	return bool(n), nil
}

type Nint int

func (n *Nint) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*int)(n))
}

func (n Nint) GetBSON() (interface{}, error) {
	return int(n), nil
}
