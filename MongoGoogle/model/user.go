package Model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GitHubData struct {
	Username          string  `json:"login"`
	Id                float64 `json:"id"`
	NodeId            string  `json:"node_id"`
	AvatarUrl         string  `json:"avatar_url"`
	GravatarId        string  `json:"gravatar_id"`
	Url               string  `json:"url"`
	HtmlUrl           string  `json:"html_url"`
	FollowersUrl      string  `json:"followers_url"`
	FollowingUrl      string  `json:"following_url"`
	GistsUrl          string  `json:"gists_url"`
	StarredUrl        string  `json:"starred_url"`
	SubscriptionsUrl  string  `json:"subscriptions_url"`
	OrganizationsUrl  string  `json:"organizations_url"`
	ReposUrl          string  `json:"repos_url"`
	EventsUrl         string  `json:"events_url"`
	RecievedEventsUrl string  `json:"recieved_events_url"`
	Type              string  `json:"type"`
	Name              string  `json:"name"`
	Company           string  `json:"company"`
	Blog              string  `json:"blog"`
	Location          string  `json:"location"`
	Email             string  `json:"email"`
	Hireable          bool    `json:"hireable"`
	Bio               string  `json:"bio"`
	TwitterUsername   string  `json:"twitter_username"`
	PublicRepos       int     `json:"public_repos"`
	PublicGists       int     `json:"public_gists"`
	Followers         int     `json:"followers"`
	Role              string  `json:"role"`
}
type GoogleData struct {
	RawData           map[string]interface{}
	Provider          string
	Email             string
	Name              string
	FirstName         string
	LastName          string
	NickName          string
	Description       string
	UserID            string
	AvatarURL         string
	Location          string
	AccessToken       string
	AccessTokenSecret string
	RefreshToken      string
	ExpiresAt         time.Time
	IDToken           string
	Role              string
}

type ApplicationUser struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Email       string             `bson:"Email"`
	FirstName   string             `bson:"FirstName"`
	LastName    string             `bson:"LastName"`
	Phone       string             `bson:"Phone"`
	DateOfBirth string             `bson:"DateOfBirth"`
	Username    string             `bson:"Username"`
	Password    string             `bson:"Password"`
	Company     string             `bson:"Company"`
	Role        string             `bson:"Role"`
	Verified    bool               `bson:"Verified"`
	Authorised  bool               `bson:"Authorised"`
}
