package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Print(`Usage:

$ .../tasks server tasks - starts server and store tasks results in
						   the folder "tasks"

$ .../tasks client manager - starts client and gets the tasks from the
							 server "manager"

$ .../tasks server interractive - starts server and gets tasks from the terminal
                                  interractively`)
		return
	}

	mode := args[0]
	argument := args[1]

	switch mode {
	case "server":
		startServer(argument)
	case "client":
		startClient(argument)
	}
}

func startServer(argument string) {
	s := &http.Server{
		Addr: ":8081",
	}

	manager := NewManager()

	manager.WatchAgents()

	http.HandleFunc("/assignTaskToMe", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/assignTaskToMe requested")

		decoder := json.NewDecoder(r.Body)
		res := struct {
			ID string `json:"id"`
		}{}

		err := decoder.Decode(&res)
		if err != nil {
			w.Write([]byte("{\"error\": \"cannot decode request body\"}"))
			return
		}

		t, ok, _ := manager.AssignNotStartedTaskToAgent(res.ID)
		if !ok {
			log.Println("no tasks")
			w.Write([]byte("{\"error\": \"no tasks\"}"))
			return
		}

		log.Printf("sending task #%s to client", t.ID)

		d, _ := json.Marshal(struct {
			ID      string `json:"id"`
			Command string `json:"command"`
		}{t.ID, t.Command})
		w.Write(d)
	})

	http.HandleFunc("/output", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var err error

		decoder := json.NewDecoder(r.Body)
		res := struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		}{}

		err = decoder.Decode(&res)
		if err != nil {
			w.Write([]byte("{\"error\": \"cannot decode request body\"}"))
			return
		}

		err = manager.AppendTaskOutput(res.ID, res.Text)
		if err != nil {
			w.Write([]byte("{\"error\": \"cannot append task's output\"}"))
			return
		}

		w.Write([]byte("{\"result\": \"ok\"}"))
	})

	http.HandleFunc("/done", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var err error

		decoder := json.NewDecoder(r.Body)
		res := struct {
			ID string `json:"id"`
		}{}

		err = decoder.Decode(&res)
		if err != nil {
			w.Write([]byte("{\"error\": \"cannot decode request body\"}"))
			return
		}

		err = manager.CompleteTask(res.ID)
		if err != nil {
			w.Write([]byte("{\"error\": \"cannot complete task\"}"))
			return
		}

		w.Write([]byte("{\"result\": \"ok\"}"))
	})

	http.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var err error

		decoder := json.NewDecoder(r.Body)
		res := struct {
			ID string `json:"id"`
		}{}

		err = decoder.Decode(&res)
		if err != nil {
			w.Write([]byte("{\"error\": \"cannot decode request body\"}"))
			return
		}

		err = manager.Heartbeat(res.ID)
		if err != nil {
			w.Write([]byte("{\"error\": \"cannot register heartbeat\"}"))
			return
		}

		w.Write([]byte("{\"result\": \"ok\"}"))
	})

	if argument == "interractive" {
		go func() {
			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Print("> ")

				text, _ := reader.ReadString('\n')
				text = strings.Replace(text, "\n", "", -1)

				id, err := uuid()
				if err != nil {
					log.Println(err.Error())
				}

				err = manager.AddTask(id, text)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}()
	}

	log.Fatal(s.ListenAndServe())
}

func startClient(argument string) {
	id, _ := uuid()

	serverUrl := "http://" + argument + ":8081/"

	go func() {
		for {
			sendHeartbeat(serverUrl, id)
			time.Sleep(3 * time.Second)
		}
	}()

	for {
		cmdId, cmdCmd, cmdOk, err := askForTask(serverUrl, id)
		if cmdOk && err == nil {
			run(cmdId, cmdCmd)
		} else {
			time.Sleep(3 * time.Second)
		}
	}
}

func sendHeartbeat(serverUrl string, id string) error {
	jsonStr := []byte("{\"id\":\"" + id + "\"}")

	req, err := http.NewRequest("POST", serverUrl+"heartbeat", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	resp.Body.Close()

	return nil
}

func askForTask(serverUrl string, id string) (string, string, bool, error) {
	jsonStr := []byte("{\"id\":\"" + id + "\"}")

	log.Println("asking the server for a task")

	req, err := http.NewRequest("POST", serverUrl+"assignTaskToMe", bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", "", false, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", "", false, err
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	res := struct {
		ID      string `json:"id"`
		Command string `json:"command"`
		Error   string `json:"error"`
	}{}

	err = decoder.Decode(&res)
	if err != nil {
		return "", "", false, err
	}

	if res.Error != "" {
		return "", "", false, nil
	}

	log.Println("received task from the server")

	return res.ID, res.Command, true, nil
}

func run(id string, cmd string) {
	fmt.Println(cmd)
}

// uuid generates a random UUID according to RFC 4122
func uuid() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
