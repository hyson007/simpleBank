//the idea is to declear interface to manage the creation of token

package token

import "time"

//we can switch between different implementation that satisfy Maker interface
type Maker interface {
	//CreateToken creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	//VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*PayLoad, error)
}
