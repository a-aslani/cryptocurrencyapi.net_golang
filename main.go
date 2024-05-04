package main

import (
	"crypto/md5"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const (
	ApiKey = "YOUR_API_KEY"
)

type CryptocurrencyApiIPNRequest struct {
	CryptocurrencyApiNet int    `json:"cryptocurrencyapi.net"`
	Chain                string `json:"chain"`
	Currency             string `json:"currency"`
	Type                 string `json:"type"`
	Date                 int64  `json:"date"`
	From                 string `json:"from"`
	To                   string `json:"to"`
	Token                string `json:"token"`
	TokenContract        string `json:"tokenContract"`
	Amount               string `json:"amount"`
	Fee                  string `json:"fee"`
	Txid                 string `json:"txid"`
	Pos                  int    `json:"pos"`
	Confirmation         int    `json:"confirmation"`
	Label                string `json:"label"`
	Sign                 string `json:"sign"`
}

func main() {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/ipn", IpnHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
		return
	}

}

func IpnHandler(c *gin.Context) {

	request := CryptocurrencyApiIPNRequest{}

	// bind and validation input
	if err := c.ShouldBindJSON(&request); err == nil {
		validate := validator.New()
		if err = validate.Struct(&request); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// confirmation number. 0 = mempool/pending 1 = appeared in the block (1st confirmation) 100 = manual confirmation value > 1 = confirmed
	if request.Confirmation <= 1 {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("confirmation: %d", request.Confirmation),
		})
		return
	}

	// check data version
	if request.CryptocurrencyApiNet < 3 {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("cryptocurrencyapi.net: %d", request.CryptocurrencyApiNet),
		})
		return
	}

	// compare sign
	err := checkSign(request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	var amount float64
	var userId int64

	if request.Type == "in" {

		amount, err = strconv.ParseFloat(request.Amount, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		userId, err = strconv.ParseInt(request.Label, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		fmt.Println(amount)
		fmt.Println(userId)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    request,
	})
}

// checkSign IPN data verification algorithm
// To check the signature you need to:
// 1) remember the received signature (sign field) from the data
// 2) remove the sign field from the data
// 3) sort the data by key in ascending order
// 4) form a string - concatenate the data values (without keys) through ':', if the value has an 'array' type (labels, ids), then use the 'Array' value for gluing
// 5) add to the string ':' and the MD5 hash from YOUR_API_KEY
// 6) get the SHA1 hash of the string (in php this is the sha1() function)
// 7) compare signature from step 1 with hash from step 6
func checkSign(request CryptocurrencyApiIPNRequest) error {

	keys := make([]string, 0)

	m := structs.Map(request)

	for k, _ := range m {
		if k != "Sign" {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)

	values := make([]string, len(keys))
	for i, k := range keys {
		values[i] = fmt.Sprintf("%v", m[k])
	}

	signData := strings.Join(values, ":")
	hashKey := fmt.Sprintf("%x", md5.Sum([]byte(ApiKey)))
	sign := fmt.Sprintf("%s:%s", signData, hashKey)
	hashSign := fmt.Sprintf("%x", sha1.Sum([]byte(sign)))

	if hashSign != request.Sign {
		return errors.New("sign wrong")
	}

	return nil
}
