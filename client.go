package twitter

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

var alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// a Twitter.Client should be used for handling all communication to the
// Twitter API.
type Client struct {
	http.Client
	Key    string
	Secret string
}

// creates a client for a given consumer key and secret.
func NewClient(consumerKey, consumerSecret string) *Client {
	return &Client{
		Client: http.Client{},
		Key:    consumerKey,
		Secret: consumerSecret,
	}
}

// generates a random string of fixed size
func nonce(size int) string {
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = alpha[rand.Intn(len(alpha)-1)]
	}
	return string(buf)
}

// signs an individual oauth request using the application token/secret pair
// and a user token/secret pair.
func (c *Client) sign(req *http.Request, token, secret string) {
	vals := map[string]string{
		"oauth_consumer_key":     c.Key,
		"oauth_nonce":            nonce(40),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        strconv.FormatInt(time.Now().Unix(), 10),
		"oauth_token":            token,
		"oauth_version":          "1.0",
	}

	// add the querystring values.  The actual RFC describing the querystring
	// allows for multiple keys, but the Twitter API doesn't actually allow
	// that, so I just take the first string from every slice
	qvals := map[string][]string(req.URL.Query())
	for k, v := range qvals {
		if len(v) > 0 {
			vals[k] = v[0]
		}
	}

	// make an alphabetical list of the keys in the auth header vals
	keys := make([]string, len(vals))
	i := 0
	for k, _ := range vals {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	// serialize the map as a string of url encoded key=val pairs.
	// i.e., join them first, then encode them.
	var baseParts []string
	for _, key := range keys {
		baseParts = append(baseParts, fmt.Sprintf("%s=%s", key, vals[key]))
	}

	// the signing key is the ampersand-joined, urlencoded app and user secrets.
	// i.e., encode them first, then join them.
	signingKey := fmt.Sprintf("%s&%s", url.QueryEscape(c.Secret), url.QueryEscape(secret))
	h := hmac.New(sha1.New, []byte(signingKey))

	// each of the following three parts is urlencoded, and then joined with ampersands:
	// the request method, the absolute request path (without query string),
	// and the ampersand-joined set of key-value pairs in the auth header map.
	io.WriteString(h, strings.Join([]string{
		url.QueryEscape(req.Method),
		url.QueryEscape(req.URL.Scheme + "://" + req.URL.Host + req.URL.Path),
		url.QueryEscape(strings.Join(baseParts, "&")),
	}, "&"))

	// alright, now we have the hmac signature.
	vals["oauth_signature"] = url.QueryEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))

	// but now we have to re-sort the list to include the new oauth_signature=whatever pair.
	keys = make([]string, len(vals))
	i = 0
	for k, _ := range vals {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	// and now build the actual header string.  These values are each
	// key="value" strings, note the ampersand.  This time, they're joined by a
	// comma, followed by either a space or a newline, and are not urlencoded.
	var s []string
	for _, key := range keys {
		s = append(s, fmt.Sprintf(`%s="%s"`, key, vals[key]))
	}
	// fuck my life, this spec is so fucking convoluted, it's such a waste of time.
	req.Header.Set("Authorization", "OAuth "+strings.Join(s, ", "))
}

// read a user stream.
func (client *Client) UserStream(token, secret string, c chan *Tweet, e chan error) ([]uint64, error) {
	request, err := http.NewRequest("GET", "https://userstream.twitter.com/2/user.json", nil)
	if err != nil {
		return nil, err
	}

	client.sign(request, token, secret)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(response.Body)
	line, err := reader.ReadBytes('\r')
	if err != nil {
		return nil, err
	}

	var friends Friends
	err = json.Unmarshal(line, &friends)
	if err != nil {
		return nil, err
	}

	go func() {
		defer response.Body.Close()
		for {
			line, err := reader.ReadBytes('\r')
			if err != nil {
				e <- err
				break
			}

			if len(line) == 2 && line[0] == '\n' && line[1] == '\r' {
				// twitter's streaming API has extra linebreaks in it that need
				// to be silenced or they throw json decode errors.
				continue
			}

			var tweet Tweet
			err = json.Unmarshal(line, &tweet)
			if err != nil {
				e <- err
			}
			c <- &tweet
		}
	}()
	return friends.Friends, nil
}

// get info about the current user
func (client *Client) UserInfo(token, secret string) (*User, error) {
	request, err := http.NewRequest("GET", "https://api.twitter.com/1/account/verify_credentials.json", nil)
	if err != nil {
		return nil, err
	}

	client.sign(request, token, secret)
	httputil.DumpRequest(request, true)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	fmt.Print(string(raw))

	// Try to decode an error document from the response
	// If we failed to decode an error document, then we didn't make an API error.
	var e struct {
		Error string `json:"error"`
	}
	err = json.Unmarshal(raw, &e)
	if err == nil && e.Error != "" {
		return nil, errors.New(e.Error)
	}

	var u User
	err = json.Unmarshal(raw, &u)
	if err != nil {
		return nil, err
	}
	log.Println("got a user")
	log.Println(u)
	return &u, nil
}

func (client *Client) Tweet(token, secret, id string) (*Tweet, error) {
	request, err := http.NewRequest("GET", "https://api.twitter.com/1/statuses/show/"+id+".json?include_entities=true", nil)
	if token != "" && secret != "" {
		client.sign(request, token, secret)
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var tweet Tweet
	if err := json.NewDecoder(response.Body).Decode(&tweet); err != nil {
		return nil, err
	}
	return &tweet, nil
}

// gets the first page of the list of the users that the particular user follows.
// TODO: page through this, getting ALL of the users that the current user follows.
func (client *Client) FriendIds(token, secret string) ([]Snowflake, error) {
	request, err := http.NewRequest("GET", "https://api.twitter.com/1/friends/ids.json", nil)
	if err != nil {
		return nil, err
	}

	client.sign(request, token, secret)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var page struct {
		NextCursor           int         `json:"next_cursor"`
		NextCursorString     string      `json:"next_cursor_str"`
		PreviousCursor       int         `json:"previous_cursor"`
		PreviousCursorString string      `json:"previous_cursor_str"`
		Ids                  []Snowflake `json:"ids"`
	}

	if err := json.NewDecoder(response.Body).Decode(&page); err != nil {
		return nil, err
	}

	return page.Ids, nil
}

// get the user's home timeline.  this doesn't actually work yet.
func (client *Client) HomeTimeline(token, secret string) ([]Tweet, error) {
	request, err := http.NewRequest("GET", "https://api.twitter.com/1/statuses/home_timeline.json", nil)
	if err != nil {
		return nil, err
	}

	client.sign(request, token, secret)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var page []Tweet
	if err := json.NewDecoder(response.Body).Decode(&page); err != nil {
		return nil, err
	}
	return page, nil
}

func (client *Client) SampleStream(token, secret string, c chan *Tweet) error {
	req, err := http.NewRequest("GET", "https://stream.twitter.com/1/statuses/sample.json", nil)
	if err != nil {
		return fmt.Errorf(`twitter: unable to create http.Request: %v`, err)
	}
	client.sign(req, token, secret)
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf(`twitter: unable to request sample stream: %v`, err)
	}

	d := json.NewDecoder(res.Body)

	go func() {
		defer res.Body.Close()
		defer close(c)
		for {
			var t Tweet
			if err := d.Decode(&t); err != nil {
				break
			}
			c <- &t
		}
	}()

	return nil
}
