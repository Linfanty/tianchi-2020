package main



import (
	"net/http"
	"encoding/json"
	"fmt"
	"os"
)

import _ "net/http/pprof"

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:8080", nil)
	}()

	player := &Player{
		manager: NewManager(),
	}

	pilots := []string{"p1", "p2", "p3", "p4", "p5", "p6", "p7", "p8", "p9", "p10"}
	
	// for _, s := range pilots {
	// 	fmt.Println(s)
	// }

	res, err := player.P1(pilots)
	if err != nil {
		fmt.Println(res)
	}

	if _, err := json.Marshal(res); err != nil {
		print("error")
	} else {
		// print(string(bs))
		print("\n\n")
	}
	
// parse 2

	dependenciesFilePath_2 := "/home/linfan.wty/tianchi/input/001.json"
	f, err := os.Open(dependenciesFilePath_2)

	if err != nil {
		print("error")
	}

	var params pParams

	if err = json.NewDecoder(f).Decode(&params); err != nil {
		print("error")
	}

	// print("\n", params[0].Apps["testapp-1021"], "\n")

	res, err = player.P2(params)

	if _, err := json.Marshal(res); err != nil {
		print("error")
	} else {
		// print(string(bs))
	}

}
