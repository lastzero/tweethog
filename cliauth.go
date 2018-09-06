package tweethog

import (
	"fmt"
	"github.com/dghubble/oauth1"
	twauth "github.com/dghubble/oauth1/twitter"
	"log"
)

func CliAuth(consumerKey string, consumerSecret string) {
	// performs Twitter PIN-based 3-legged OAuth 1 from the command line

	// read credentials from environment variables
	if consumerKey == "" || consumerSecret == "" {
		log.Fatal("Required environment variable missing.")
	}

	oauthConfig := &oauth1.Config{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		CallbackURL:    "oob",
		Endpoint:       twauth.AuthorizeEndpoint,
	}

	requestToken, err := login(oauthConfig)
	if err != nil {
		log.Fatalf("Request Token Phase: %s", err.Error())
	}
	accessToken, err := receivePIN(oauthConfig, requestToken)
	if err != nil {
		log.Fatalf("Access Token Phase: %s", err.Error())
	}

	fmt.Printf("\nThis app was granted an access token to act on behalf of your user.\n")

	fmt.Printf("Access Token: %s\nAccess Secret: %s\n", accessToken.Token, accessToken.TokenSecret)
}

func login(oauthConfig *oauth1.Config) (requestToken string, err error) {
	requestToken, _, err = oauthConfig.RequestToken()
	if err != nil {
		return "", err
	}
	authorizationURL, err := oauthConfig.AuthorizationURL(requestToken)
	if err != nil {
		return "", err
	}
	fmt.Printf("Open this URL in your browser:\n%s\n", authorizationURL.String())
	return requestToken, err
}

func receivePIN(oauthConfig *oauth1.Config, requestToken string) (*oauth1.Token, error) {
	fmt.Printf("Paste your PIN here: ")

	var verifier string
	_, err := fmt.Scanf("%s", &verifier)

	if err != nil {
		return nil, err
	}

	// Twitter ignores the oauth_signature on the access token request. The user
	// to which the request (temporary) token corresponds is already known on the
	// server. The request for a request token earlier was validated signed by
	// the consumer. Consumer applications can avoid keeping request token state
	// between authorization granting and callback handling.
	accessToken, accessSecret, err := oauthConfig.AccessToken(requestToken, "secret does not matter", verifier)
	if err != nil {
		return nil, err
	}

	return oauth1.NewToken(accessToken, accessSecret), err
}
