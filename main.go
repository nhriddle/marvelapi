// Classicifation of Marvel Characters API
//
// Documentation for Marvel Characters API
// 
// Schemes : http
// BasePath : /
// Version : 0.0.1
// 
// Consumes :
// - application/json
// 
// Produces : 
// - application/json
// swagger:meta

package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/ilyakaznacheev/cleanenv"
	//"bufio"
	//"os"
    "time"
    "strconv"
    "io/ioutil"
	"crypto/md5"
	"encoding/hex"
	"fmt"
    "github.com/davecgh/go-spew/spew"
    "context"
    "github.com/go-redis/redis/v8"
)

// configuration structure
type Server struct {
    Port string `yaml:"port"`
    Host string `yaml:"host"`
}

type Marvelapi struct {
    Url string `yaml:"url"`
    Publickey string `yaml:"publickey"`
    Privatekey string `yaml:"privatekey"`
}

type Redis struct {
    Host string `yaml:"host"`
    Port string `yaml:"port"`
}

type EnvConfig struct {
    ServerConf *Server `yaml:"server"`
    MarvelapiConf *Marvelapi `yaml:"marvelapi"`
    RedisConf *Redis `yaml:"redis"`
}

// Redis struct
type RedisDBCli struct {
   Client *redis.Client
}


// Marvel Request Data Structure
type Items struct{
	ResourceURI string `json:"resourceURI"`
	Name string `json:"name"`
	Type string `json:"type"` 
}
type Generic struct {
	Avaible int `json:"available"`
	Returned int `json:"returned"`
	CollectionURI string `json:"collectionURI"`
	Items []Items `json:"items"`
}
type Series struct {
	Generic
}
type Events struct {
	Generic
}

type Stories struct {
	Generic 
}
type Comics struct {
	Generic 
}

type Thumbnail struct {
	Path string `json:"path"`
	Extension string `json:"extension"`
}
type Urls struct {
	Type string `json:"type"`
	Url string  `json:"url"`
}

type Results struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	//Modified string `json:"modified"` 
	//ResourceURI string `json:"resourceURI"`
	//Urls 	[]Urls `json:"urls"`
	//Thumbnail Thumbnail `json:"thumbnail"`
	//Comics Comics `json:"comics"`
	//Stories Stories `json:"stories"`
	//Events Events `json:"events"`
	//Series Series `json:"series"`
}
type Data struct {
	Offset int `json:"offset"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Count int `json:"count"`
	Results []Results `json:"results"`
}

type Feed struct {
	Data Data  `json:"data"`
}


// Character Struct
type Character struct {
    ID      string `json:"id"`
    Name    string `json:"name"`
    Description  string `json:"description"`
}

// utility functions
func GetHashString(s string)string{
	bundle := []byte(s);
	array := md5.Sum(bundle);
	return hex.EncodeToString(array[:]);
}


// Get All Characters
// swagger:route GET /characters characters getCharacters
// Returns a list of character IDs from the marvel universe
// responses:
//  200 : jsonArray
func GetCharacters(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json")

    // Get values from redis cache
    val, rediserr := rdb.Get(ctx, "MARVELAPI:IDS").Result()
    if rediserr != nil {
        data := storeHeroes()
        fmt.Printf(string(data))
        var response []string
        json.Unmarshal([]byte(data), &response)
        json.NewEncoder(w).Encode(response)

    } else {

        fmt.Println("MARVELAPI:IDS" + string(val))
        fmt.Println("Received from Redis!!!")
        var response []string

        json.Unmarshal([]byte(val), &response)
        json.NewEncoder(w).Encode(response)

    }

}

// Get Character via ID
// Get All Characters
// swagger:route GET /characters/{id} character getCharacter
// Returns details of a character's IDs
// responses:
//  200 : jsonObject
//  404 : httpResponse
func GetCharacter(w http.ResponseWriter, r *http.Request) {

    // Set content type to JSON
    w.Header().Set("Content-Type", "application/json")

    // Get values from redis cache
    params := mux.Vars(r) // Get Params

    val, rediserr := rdb.Get(ctx, "MARVELAPI:" + params["id"]).Result()
    if rediserr != nil {
        // contemplated and removed the additional call for non existing IDs let's just have it wait until the next cache reload
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{"code":"404","failed": "Character Not Found"})

    } else {
        //json.NewEncoder(w).Encode(data.Data.Results[0])
        fmt.Println("Received from Redis!!!")
        var response Results

        json.Unmarshal([]byte(val), &response)
        json.NewEncoder(w).Encode(response)

    }

}

// Function to print results
func PrintResult(data Feed,err error){
	if err != nil {
		fmt.Printf("occured error: %v",err)
	}else if data.Data.Count == 0 {
		fmt.Println("empty information")
	}else{
		spew.Dump(data)
    }
}

// Function to call marvel API
func callMarvelAPI(endpoint string) (Feed, error) {

    var entries Feed

    //url := cfg.MarvelapiConf.Url + endpoint + "?apikey=" + cfg.MarvelapiConf.Publickey
	ts := "my_hash";
	hash :=GetHashString(ts + cfg.MarvelapiConf.Privatekey  + cfg.MarvelapiConf.Publickey)

	url := cfg.MarvelapiConf.Url + endpoint +
			"?ts=" +ts+
			"&apikey=" +cfg.MarvelapiConf.Publickey+
			"&hash="+hash

    response, err := http.Get(url)
	if err == nil	{
		data, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(data, &entries)
	}

    return entries,err
}

func storeHeroes() (string) {

    var entries Feed
    totalcount  := 100
    offset := 0
    entriesperpage := 100

    var heroids []string

    // create string array for ids
    ts := "my_hash";
    hash :=GetHashString(ts + cfg.MarvelapiConf.Privatekey  + cfg.MarvelapiConf.Publickey)

    for offset < totalcount {

        url := cfg.MarvelapiConf.Url + "/characters" +
            "?ts=" +ts+
            "&apikey=" +cfg.MarvelapiConf.Publickey+
            "&hash="+hash+
            "&offset=" + strconv.Itoa(offset)+
            "&limit=" + strconv.Itoa(entriesperpage)

        response, err := http.Get(url)

        fmt.Println(url)

        if err == nil   {
            data, _ := ioutil.ReadAll(response.Body)
            err = json.Unmarshal(data, &entries)

            //set total count to the total entries
            totalcount = entries.Data.Total
            for _, v := range entries.Data.Results {

                jsonobj, jsonerr := json.Marshal(v)
                if jsonerr != nil {
                    log.Fatal("unable to process json object")
                } else {
                    fmt.Println(strconv.Itoa(offset))
                    fmt.Println(strconv.Itoa(v.Id))

                    //insert to redis
                    jsonstr := string(jsonobj)
                    inserterr := rdb.Set(ctx, "MARVELAPI:" + strconv.Itoa(v.Id)  , jsonstr, 0).Err()
                    if inserterr != nil {
                        log.Fatal("unable to insert to redis")
                    } else {
                        fmt.Println("Inserted to Redis!!!")
                    }

                    // append the id to the array
                    heroids = append(heroids, strconv.Itoa(v.Id))
                    // increase offset
                    offset++
                }

            }

        } else {
            log.Fatal("marvel api call failed")
            break
        }

    } 

    fmt.Println(heroids)

    idjsonstr := ""
    // insert heroids to redis
    idjsonobj, idjsonerr := json.Marshal(heroids)
    if idjsonerr != nil {
        log.Fatal("unable to process json object")
    } else {
        //insert to redis
        idjsonstr = string(idjsonobj)
        // expiry is 1 day
        inserterr := rdb.SetEX(ctx, "MARVELAPI:IDS" , idjsonstr, 86400*time.Second).Err()
        if inserterr != nil {
            log.Fatal("unable to insert to redis")
        } else {
            fmt.Println("Inserted to Redis!!!")
        }
    }

    return idjsonstr

}


var cfg EnvConfig
var ctx = context.Background()
var rdb *redis.Client
func main() {

    // read environment variables from config file
    err := cleanenv.ReadConfig("config.yml", &cfg)
    if err != nil {
        log.Fatal("Unable to read configuration file")
        return
    }

    // connect redis here so there would only be one redis connection
    rdb = redis.NewClient(&redis.Options{
        Addr:     cfg.RedisConf.Host +  cfg.RedisConf.Port,
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    // Init MUX Router
    r := mux.NewRouter()

    // Route Handler / Endpoints
    r.HandleFunc("/characters", GetCharacters).Methods("GET")
    r.HandleFunc("/characters/{id}", GetCharacter).Methods("GET")

    log.Print("Running in port " + cfg.ServerConf.Port)
    log.Fatal(http.ListenAndServe(cfg.ServerConf.Port ,r))
}
