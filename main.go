package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	twinID := flag.String("twin", "A", "Twin A or B")
	port1 := flag.Int("port1", 8001, "Port to run on")
	port2 := flag.Int("port2", 8002, "Peer port")
	flag.Parse()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	go func() {
		log.Printf("Twin %s running on port %d\n", *twinID, *port1)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port1), nil))
	}()

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	w := worker.New(c, "TWIN_TASK_QUEUE", worker.Options{})
	w.RegisterWorkflow(MonitorTwinWorkflow)
	w.RegisterActivity(CheckHealth)
	w.RegisterActivity(RestartTwin)

	go func() {
		log.Fatal(w.Run(worker.InterruptCh()))
	}()

	peerURL := fmt.Sprintf("http://localhost:%d/health", *port2)
	workflowID := "monitor-" + *twinID

	_, err = c.ExecuteWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			ID:        workflowID,
			TaskQueue: "TWIN_TASK_QUEUE",
		},
		MonitorTwinWorkflow,
		peerURL,
		*twinID,
		*port2,
	)

	if err != nil {
		log.Println("Workflow already running:", err)
	}

	select {}
}
