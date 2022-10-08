package shortener

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"url-shortener/data"

	"github.com/gin-gonic/gin"
	uuid "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type service struct {
	Collection *mongo.Collection
}

func NewService(collection *mongo.Collection) *service {
	return &service{
		Collection: collection,
	}
}

// shortens src url
func (svc *service) Shorten(ctx *gin.Context) {
	url, ok := ctx.Params.Get("src")
	if !ok {
		ctx.JSON(401, "some error occured")
		return
	}
	// srcUrl, err := url.ParseRequestURI(src)
	// if err != nil {
	// ctx.JSON(401, "error parsing url")
	// return
	// }
	// url := srcUrl.String()

	shortenedUrl := strings.Split(uuid.New().String(), "-")[0]

	ttl := 30 * time.Minute
	Url := data.UrlShortener{
		Src:        url,
		Dest:       shortenedUrl,
		Created_at: time.Now(),
		Expires_at: time.Now().UTC().Add(ttl),
		Hits:       0,
	}

	_, err := svc.Collection.InsertOne(ctx, &Url)
	if err != nil {
		ctx.JSON(500, "error creating shortener")
		return
	}

	ctx.JSON(201, shortenedUrl)
	return

}

// returns dest, given src
func (svc *service) Get(ctx *gin.Context) {
	dest, ok := ctx.Params.Get("dest")
	if !ok {
		ctx.JSON(401, "some error occured")
		return
	}

	fmt.Println(dest)

	urlS := &data.UrlShortener{}
	err := svc.Collection.FindOne(ctx, bson.D{{"dest", fmt.Sprintf("%s", dest)}}).Decode(urlS)
	if err != nil {
		ctx.JSON(500, "error getting url")
		return
	}

	fmt.Println(urlS)

	update := bson.D{{"$inc", bson.D{{"hits", 1}}}}
	_, err = svc.Collection.UpdateOne(ctx, bson.D{{"dest", fmt.Sprintf("%s", dest)}}, update)

	ctx.Redirect(http.StatusPermanentRedirect, urlS.Src)
	return
}
