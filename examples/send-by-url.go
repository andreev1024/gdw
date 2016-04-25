/*
  1) Sets real tokens, guid, oauth config credentials etc.
  2) The first time you run the sample, it will prompt you to authorize access:

    * Browse to the provided URL in your web browser.
    * If you are not already logged into your Google account, you will be prompted to log in. If you are logged into multiple Google accounts, you will be asked to select one account to use for the authorization.
    Click the Accept button.
    * Copy the code you're given, paste it into the command-line prompt, and press Enter.

  3) Check your new file in Google Drive.
*/
package examples

import (
	"github.com/andreev1024/gdw"
	"github.com/andreev1024/rsw"

	"google.golang.org/api/drive/v3"

	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/oauth2"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//  RightSignature
	rswApiToken := "123456789"
	guid := "987654321"

	//  Google
	config := &oauth2.Config{
		ClientID:     "123456789.apps.googleusercontent.com",
		ClientSecret: "iAmAClientSecret",
		Scopes:       []string{"https://www.googleapis.com/auth/drive.file"},
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}

	token := getToken(config)

	rsw := rsw.NewAPI(rswApiToken)
	details, err := rsw.DocumentDetails(guid)
	if err != nil {
		return
	}

	downloadUrl, err := url.QueryUnescape(details.SignedPdfURL)
	if err != nil {
		return
	}

	gdFile := new(drive.File)
	//  set Drive file name
	gdFile.Name = "test.pdf"

	f, err := gdw.SendByUrl(downloadUrl, gdFile, config, token)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%#v\n", f)
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getToken(config *oauth2.Config) *oauth2.Token {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return tok
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("drive-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
