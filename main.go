package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

var inputJson = `{
   "gophers": [{
         "name": "A",
         "sleep": 1,
         "eat": 1
      },
      {
         "name": "B",
         "sleep": 4,
         "eat": 3
      },
      {
         "name": "C",
         "sleep": 1,
         "eat": 4
      },
      {
         "name": "D",
         "sleep": 5,
         "eat": 2
      },
      {
         "name": "E",
         "sleep": 2,
         "eat": 3
      }
   ],
   "totalFood": 30
}`

type Farm struct {
	Gophers   []Gopher `json:"gophers"`
	TotalFood int      `json:"totalFood"`
	sync.Mutex
}

func (f *Farm) eatFood(gopher *Gopher) error {
	f.Lock()
	defer f.Unlock()

	if f.TotalFood < gopher.Eat {
		log.Printf("gopher %s wants to eat %v food unit(s) but there's not enough food!", gopher.Name, gopher.Eat)
		return errors.New("there's not enough food")
	}

	f.TotalFood -= gopher.Eat
	log.Printf("gopher %s eats %v food unit(s). %v food unit(s) left.", gopher.Name, gopher.Eat, f.TotalFood)
	return nil
}

type Gopher struct {
	Name  string        `json:"name"`
	Sleep time.Duration `json:"sleep"`
	Eat   int           `json:"eat"`
}

func (gopher *Gopher) gopherLive(farm *Farm, messages chan string) {
	for {
		time.Sleep(time.Second * gopher.Sleep)
		if err := farm.eatFood(gopher); err != nil {
			messages <- fmt.Sprintf("gopher %s dies. So said!", gopher.Name)
			return
		}
	}
}

func main() {
	gopherFarm := &Farm{}
	if err := json.Unmarshal([]byte(inputJson), &gopherFarm); err != nil {
		log.Fatal("cannot unmarshal: ", err)
	}

	messages := make(chan string)

	for _, gopher := range gopherFarm.Gophers {
		log.Printf("gopher %s joins the farm!", gopher.Name)
		go gopher.gopherLive(gopherFarm, messages)
	}

	log.Println("waiting")

	for range gopherFarm.Gophers {
		msg, ok := <-messages
		if !ok {
			log.Fatal("channel was closed")
		}

		log.Println(msg)
	}

	log.Println("done waiting")
}
