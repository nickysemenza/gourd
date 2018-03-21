package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/nickysemenza/food/backend/app/config"
	"github.com/nickysemenza/food/backend/app/model"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"os"
)

var (
	oauthConf = &oauth2.Config{
		ClientID:     os.Getenv("FACEBOOK_APP_ID"),
		ClientSecret: os.Getenv("FACEBOOK_APP_SECRET"),
		RedirectURL:  os.Getenv("API_PUBLIC_URL") + "/auth/facebook/callback",
		Scopes:       []string{"public_profile", "email"},
		Endpoint:     facebook.Endpoint,
	}
	oauthStateString = "thisshouldberandom"
)

//GetMe gives the current User as a JSON response to GET /me
func GetMe(c *gin.Context) {
	//TODO: aah respondSuccess(w, e.CurrentUser)
	c.JSON(http.StatusOK, "todo")
}

//HandleFacebookLogin initiates the facebook auth process
func HandleFacebookLogin(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	URL, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return nil
}

//HandleFacebookCallback is the callback for the facebook auth process
func HandleFacebookCallback(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil
	}

	code := r.FormValue("code")

	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return err
	}

	resp, err := http.Get("https://graph.facebook.com/me?fields=id,name,email&access_token=" +
		url.QueryEscape(token.AccessToken))
	if err != nil {
		fmt.Printf("Get: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return err
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ReadAll: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return err
	}

	var fbr facebookUserData
	err = json.Unmarshal(response, &fbr)
	if err != nil {
		return err
	}

	log.Printf("parseResponseBody: %s\n", string(response))
	//log.Printf("token: %v\n",)

	user := fbr.getUser(e)
	tok := user.GetJWTToken(e.DB)

	log.Printf("TOKEN for user %d: %s\n", user.ID, tok)

	http.Redirect(w, r, os.Getenv("FRONTEND_URL")+"/auth/"+tok, http.StatusTemporaryRedirect)

	//http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return nil
}

type facebookUserData struct {
	FBID  string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (d facebookUserData) getUser(e *config.Env) *model.User {
	u := model.User{}

	if e.DB.Where("email = ?", d.Email).First(&u).RecordNotFound() {
		//no user exists with email
		if e.DB.Where("facebook_id = ?", d.FBID).First(&u).RecordNotFound() {
			//no user exists with facebook_id
			//we must have a new user
			nameParts := strings.SplitN(d.Name, " ", 2)
			switch len(nameParts) {
			case 2:
				u.FirstName = nameParts[0]
				u.LastName = nameParts[1]
			case 1:
				u.FirstName = nameParts[0]
			}
		}
	}
	//can't hurt to write these every time
	u.Email = d.Email
	u.FacebookID = d.FBID
	e.DB.Save(&u)

	return &u
}

func getUserFromToken(e *config.Env, tokenString string) (*model.User, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if token.Valid {
		fmt.Println("token is valid")
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		fmt.Printf("%v %v", claims.Id, claims.Issuer)
		u := model.User{}
		e.DB.First(&u, claims.Id)
		return &u, nil
	} else {
		fmt.Println(err)
		return nil, err
	}
}
