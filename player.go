package main

import (
	"encoding/json"
	"os"
	"path"
	"sync"
	"time"
	"fmt"

	log "github.com/golang/glog"
)

var (
	dataDir              string
	dependenciesFileName string
	dependenciesFilePath string

	apiDelay time.Duration
)

func init() {
	// dataDir = "/root/input"
	dataDir = "/home/linfan.wty/tianchi/input"
	if v := os.Getenv("DATA_DIR"); v != "" {
		dataDir = v
	}

	dependenciesFileName = "data.json"
	dependenciesFilePath = path.Join(dataDir, dependenciesFileName)
	// print(dependenciesFilePath, "\n")

	if v := os.Getenv("API_DELAY"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			apiDelay = d
		}
	}
}

type pParams struct {
	Apps         map[string]int            `json:"apps"`
	Dependencies map[string]map[string]int `json:"dependencies"`
}

type Player struct {
	mut              sync.Mutex
	prepared, inited bool
	manager          *manager
}

func (p *Player) Ready() bool {
	p.mut.Lock()
	defer p.mut.Unlock()

	return p.prepared
}

func (p *Player) init(pilots []string) {
	p.manager.initPilots(pilots)
}

func (p *Player) Reset() error {
	p.mut.Lock()
	defer p.mut.Unlock()

	p.inited = false
	p.manager.Reset()
	return nil
}

func (p *Player) P1(pilots []string) (map[string][]string, error) {
	p.mut.Lock()
	defer p.mut.Unlock()
	

	if !p.inited {
		p.init(pilots)
		p.inited = true
	}

	// Load base data from file.
	f, err := os.Open(dependenciesFilePath)
 	
	if err != nil {
		return nil, err
	}
	var params pParams
	if err = json.NewDecoder(f).Decode(&params); err != nil {
		return nil, err
	}

	// fmt.Println(params.Apps)
	// fmt.Println(params.Dependencies)
	
	if err := p.manager.UpdateAppDependencies_1(params.Apps, params.Dependencies); err != nil {
		return nil, err
	}

	return p.manager.Result(), nil
}

func (p *Player) P2(params pParams) (map[string][]string, error) {
	p.mut.Lock()
	defer p.mut.Unlock()

	if !p.inited {
		return nil, fmt.Errorf("not inited")
	}
	

	if err := p.manager.UpdateAppDependencies_2(params.Apps, params.Dependencies); err != nil {
		return nil, err
	}

	return p.manager.Result(), nil
}

func (p *Player) startPrepare() {
	// We may do some prepare work and will not return true until they're finished.
	// Here we use sleep to simulate that.
	
	go func() {
		pTime := 3 * time.Second
		if v := os.Getenv("PREPARE_TIME"); v != "" {
			if d, err := time.ParseDuration(v); err != nil {
				pTime = d
			}
		}
		time.Sleep(pTime)

		p.mut.Lock()
		p.prepared = true
		p.mut.Unlock()

		log.Infof("finish preparation")
	}()
}

func (p *Player) Run() {
	p.startPrepare()
}
