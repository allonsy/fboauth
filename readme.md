# fboauth by Allonsy
* A facebook oauth library written in Go
* Use to perform oauth handshakes for go desktop apps
* (c) 2016 Allonsy

## Uses
* import `github.com/allonsy/fboauth`
* create a json file somewhere with the key `clientid` with the value of your clientID
* use `getAuthCode(path string, permissions []string)` to just grab the verfication code and redirect url
  * returns a url for the user to visit and a code for them to put in
  * Then, in your application, display the url to the user for them to visit and paste the returned code in
  * you need to query facebook for the access token when done
* use `GetAccessCode(path string, permissions []string)` to perform the entire handshake
  * prints out the url and code to the console so the user can visit the site and enter the
  * then, it queries facebook for a response, when it recieves one (or times out), it returns the access code
  * on error, error is not null
* More documentation is on the golang website


## Contributions
  * Contributions are welcome, just submit a PR
  * Any issues or feature requests can be submitted via issues
