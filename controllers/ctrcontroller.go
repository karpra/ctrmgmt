package controllers

import (
	"context"
	"ctrmgmt/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func GetVersion(w http.ResponseWriter, r *http.Request) {
	version := "0.0.1"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(version)
}

func GetContainers(w http.ResponseWriter, r *http.Request) {
	var ctrs []models.CtrMgt
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		ctrs = append(ctrs, models.CtrMgt{container.ID[:10], container.Names[0][1:], container.Image})
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
	json.NewEncoder(w).Encode(ctrs)
}

func CreateContainers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ctr models.CtrMgt
	//json.NewDecoder(r.Body).Decode(&container)
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "8000",
	}
	containerPort, err := nat.NewPort("tcp", "80")
	if err != nil {
		panic("Unable to get the port")
	}

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{},
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: "nginx",
		},
		&container.HostConfig{
			PortBindings: portBinding,
		}, networkConfig, nil, "web1")
	if err != nil {
		panic(err)
	}

	cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	fmt.Printf("Container %s is started", cont.ID)
	json.NewEncoder(w).Encode(ctr)
}

func StopContainers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ctr models.CtrMgt
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	containerName := "web1"

	for _, container := range containers {

		if containerName == container.Names[0][1:] {

			if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
				log.Printf("Unable to stop container %s: %s", container.ID, err)
			}

			removeOptions := types.ContainerRemoveOptions{
				RemoveVolumes: true,
				Force:         true,
			}

			if err := cli.ContainerRemove(ctx, container.ID, removeOptions); err != nil {
				log.Printf("Unable to remove container: %s", err)
				//return err
			}
		}
	}

	json.NewEncoder(w).Encode(ctr)
}
