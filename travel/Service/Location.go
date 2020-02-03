package travel

import (
	"fmt"
	"net/http"

	"../../common"
	"../../common/redis"
	"github.com/go-chi/render"
)

func GetLocationFromCoordinates(w http.ResponseWriter, r *http.Request) {
	
	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	req, err := getBodyFromHttpRequest(r, req);

	// Truncing down to two digits
	var latlon = fmt.Sprintf("%.2f%.2f", req.Latitude, req.Longitude)

	res, err := redis.Instance.Get(latlon).Result()

	

	if err == nil {
		fmt.Printf("Serving from location cache\n")
		json.Unmarshal([]byte(res), &resultJSON)
		render.JSON(w, r, resultJSON)
		return
	}

	result, err = ReverseLookup(req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", result)

	fmt.Printf("Serving from locationiq server\n")

	if resultJSON.Error == "" {
		redis.Instance.Set(latlon, fmt.Sprintf("%s", body), 99999999999)
	}
	render.JSON(w, r, resultJSON)
}

func GetCoordinatesFromLocation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Location string `json:"location"`
	}

	req, err := getBodyFromHttpRequest(r, req);
	
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Printf("\nSearchlocation api is called\n")
	result, err := ForwardLookup(req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Serving from locationiq server\n")

	var slice [][]byte
	for _, element := range result {
		res, _ := json.Marshal(element)
		slice = append(slice, res)
	}

	redis.Instance.Set(req.Location, fmt.Sprintf("%s", body), 99999999999)
	render.JSON(w, r, resultJSON)
}



func ReverseLookup(req interface{}) {
	var result struct {
		DisplayName string                 `json:"display_name"`
		Address     map[string]interface{} `json:"address"`
		Error       string                 `json:"error"`
	}

	baseurl := "https://us1.locationiq.com/v1/reverse.php?key=%s&lat=%f&lon=%f&format=json"
	formattedurl := fmt.Sprintf(baseurl, common.LOCATIONIQ_TOKEN, req.Latitude, req.Longitude);
	res, err := http.Get(formattedurl)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(body), &result)
	
	return result, nil
}

func ForwardLookup(req interface{}) {
	var result []struct {
		Latitude    string `json:"lat"`
		Longitude   string `json:"lon"`
		DisplayName string `json:"display_name"`
	}

	baseurl := "https://us1.locationiq.com/v1/search.php?key=%s&q=%s&format=json"
	formattedurl := fmt.Sprintf(baseurl, common.LOCATIONIQ_TOKEN, req.Location);
	res, err := http.Get(formattedurl)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(body), &result)
	
	return result, nil
}