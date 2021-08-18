package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/evgeniron/IPGeoLocation/pkg/db"
)

/* HTTP Server timeout variables */
var (
	requestTimeout = 5 * time.Second
	readTimeout    = 2 * time.Second
	writeTimeout   = 2 * time.Second
	idleTimeout    = 5 * time.Second
)

/* Set max ip length avaible for querying */
const MAXIPV6LEN int = len("0000:0000:0000:0000:0000:ffff:192.168.100.228")

type Server struct {
	config   *Config
	db       db.DB
	rateChan chan bool
}

/* Send JSON response */
func jsonResponse(w http.ResponseWriter, message interface{}, code int) {
	response, _ := json.Marshal(message)

	/* Set response header to json content type */
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(code)

	_, err := w.Write(response)
	if err != nil {
		log.WithFields(log.Fields{
			"event":     "reply",
			"json":      response,
			"HTTP Code": 0,
		}).Warn("Could not send response to client")
	}

	log.WithFields(log.Fields{
		"event":     "reply",
		"json":      response,
		"HTTP Code": code,
	}).Info("HTTP")
}

/* Send JSON error response: {"error": <error code>} */
func errorResponse(w http.ResponseWriter, err_code int) {
	jsonResponse(w, map[string]int{"error": err_code}, err_code)
}

/* Serve Welcome page only */
func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.WithFields(log.Fields{
			"event": "unauthorised-acces",
			"path":  r.URL.Path,
		}).Warn("Access-denied")
		errorResponse(w, http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		w.Write([]byte("<h1>Welcome to my IPGelocation Server!</h1>"))
	} else {
		errorResponse(w, http.StatusMethodNotAllowed)
	}
}

/* Limit service to GET method and check if query is empty or too long */
func checkIPLocationQuery(w http.ResponseWriter, r *http.Request) (string, bool) {
	if r.Method != http.MethodGet {
		log.WithFields(log.Fields{
			"event":  "unauthorised-method",
			"method": r.Method,
		}).Warn("Access-denied")
		return "", false
	}

	ipQuery := r.URL.Query().Get("ip")
	if ipQuery == "" || len(ipQuery) > MAXIPV6LEN {
		log.WithFields(log.Fields{
			"event": "unathorised-query",
			"query": ipQuery,
		}).Warn("Access-denied")
		return "", false
	}

	return ipQuery, true
}

/* Server configuration constructor with rate limiter, db and env configuration */
func NewServerConf() *Server {
	server := Server{}
	server.config = NewConfig("APP")
	server.db = db.NewCsvDb()
	server.rateChan = make(chan bool, server.config.MaxRequestPerSec)

	log.WithFields(log.Fields{
		"event":                          "configuration",
		"module":                         "server",
		"server.config.port":             server.config.Port,
		"server.config.maxrequestpersec": server.config.MaxRequestPerSec,
	}).Info("Configuration set")
	return &server
}

/* Find location of the IP and return the result */
func (s *Server) handleIPQuery(w http.ResponseWriter, r *http.Request) {

	select {
	case <-s.rateChan:
		if ipQuery, verified := checkIPLocationQuery(w, r); verified {
			result := s.db.GetLocation(ipQuery)

			/* If both fields are empty, reply error*/
			if result.City == "" && result.Country == "" {
				errorResponse(w, http.StatusInternalServerError)
				return
			}
			jsonResponse(w, result, http.StatusOK)
		} else {
			errorResponse(w, http.StatusBadRequest)
		}
		return
	default:
		errorResponse(w, http.StatusTooManyRequests)
	}
}

/* Rate limiter using channel */
func (s *Server) StartTimer(rate int) {
	ratePerSec := time.Tick(time.Duration(rate*1000) * time.Millisecond)
	for range ratePerSec {
		s.rateChan <- true
	}
}

func (s *Server) Run() {
	go s.StartTimer(s.config.MaxRequestPerSec)

	mux := http.ServeMux{}
	mux.Handle("/", http.TimeoutHandler(http.HandlerFunc(s.indexHandler), requestTimeout, "Timeout"))
	mux.Handle("/v1/find-country", http.TimeoutHandler(http.HandlerFunc(s.handleIPQuery), requestTimeout, "Timeout"))

	server := http.Server{
		Addr:         ":" + s.config.Port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      &mux,
	}

	server.ListenAndServe()
}

func init() {
	debugLevel := flag.Bool("d", false, "Print debug logs")
	flag.Parse()
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Log the warning severity or above unless -d flag provided.
	logLevel := log.WarnLevel
	if *debugLevel {
		logLevel = log.DebugLevel
	}

	log.SetLevel(logLevel)
}

func main() {

	s := NewServerConf()
	log.WithFields(log.Fields{
		"event": "start",
		"port":  s.config.Port,
	}).Info("Starting server")
	s.Run()

}
