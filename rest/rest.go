package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rlaalsrl715/nomadcoin/blockchian"
	"github.com/rlaalsrl715/nomadcoin/utils"
	"log"
	"net/http"
)

var port string

type url string

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
	CreatedAt   int64  `json:"-"`
}

type addBlockBody struct {
	Message string
}

func (u url) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("http://localhost%s%s", port, u)), nil
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{id}"),
			Method:      "GET",
			Description: "See A Block",
		},
	}
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rw.Header().Add("Content-Type", "application/json")
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchian.GetBlockchain().AllBLocks()))

	case "POST":
		var addBlockBody addBlockBody
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchian.GetBlockchain().AppendBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated)
	}
}

func Start(aPort int) {
	//url과 url의 func을 매핑해주는 역할을한다. 근데 url이 겹치면 오류!
	handler := http.NewServeMux()
	port = fmt.Sprintf(":%d", aPort)
	handler.HandleFunc("/", documentation)
	handler.HandleFunc("/blocks", blocks)
	fmt.Printf("LIstening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
