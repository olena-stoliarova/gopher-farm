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
	Gophers []Gopher `json:"gophers"`
	TotalFood int `json:"totalFood"`
	mu sync.Mutex
}

type Gopher struct {
	Name string `json:"name"`
	Sleep time.Duration `json:"sleep"`
	Eat int `json:"eat"`
}
func (gopher *Gopher) gopherLive(farm *Farm, wg *sync.WaitGroup, messages chan string) {
	defer wg.Done()
	for {
		time.Sleep(time.Second * gopher.Sleep)
		err := farm.eatFood(gopher)
		if err != nil {
			messages <- fmt.Sprintf("gopher %s dies. So said!", gopher.Name)
			return
		}
	}
}

func (f *Farm) eatFood(gopher *Gopher) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.TotalFood < gopher.Eat {
		log.Printf("gopher %s wants to eat %v food unit(s) but there's not enough food!", gopher.Name, gopher.Eat)
		return errors.New("there's not enough food")
	}
	f.TotalFood = f.TotalFood - gopher.Eat
	log.Printf("gopher %s eats %v food unit(s). %v food unit(s) left.", gopher.Name, gopher.Eat, f.TotalFood)
	return nil
}


func main()  {
	gopherFarm := &Farm{}
	err := json.Unmarshal([]byte(inputJson), gopherFarm)
	if err != nil {
		log.Fatal("cannot unmarshal: ", err)
	}

	messages := make(chan string)
	wg := new(sync.WaitGroup)

	for i := range(gopherFarm.Gophers){
		wg.Add(1)
		log.Printf("gopher %s joins the farm!", gopherFarm.Gophers[i].Name)
		go (&gopherFarm.Gophers[i]).gopherLive(gopherFarm, wg, messages)
	}

	go func(wg *sync.WaitGroup, messages chan string) {
		log.Println("waiting")
		wg.Wait()
		log.Println("done waiting")
		close(messages)
	}(wg, messages)

	for msg := range messages {
		log.Println(msg)
	}
}
