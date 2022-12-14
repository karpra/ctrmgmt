package controllers

import (
	"context"
	"ctrmgmt/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gorilla/mux"
)

func GetVersion(w http.ResponseWriter, r *http.Request) {
	version := "0.0.1"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(version)
}

/*
 * lists the containers
 */
func GetContainers(w http.ResponseWriter, r *http.Request) {
	ctrs := make([]models.CtrMgt, 0)
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
		ctrs = append(ctrs, models.CtrMgt{Id: container.ID[:10], Name: container.Names[0][1:], Image: container.Image, State: container.State, Status: container.Status, Ports: container.Ports})
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
	json.NewEncoder(w).Encode(ctrs)
}

/*
 * creates and starts the container
 */
func CreateContainers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ctr models.CtrMgt
	json.NewDecoder(r.Body).Decode(&ctr)
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}
	reader, err := cli.ImagePull(ctx, "docker.io/library/"+ctr.Image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	defer reader.Close()
	//io.Copy(os.Stdout, reader)
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
			Image: ctr.Image,
		},
		&container.HostConfig{
			PortBindings: portBinding,
		}, networkConfig, nil, ctr.Name)
	if err != nil {
		panic(err)
	}

	cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	fmt.Printf("Container %s is started", cont.ID)
	json.NewEncoder(w).Encode(ctr)
}

/*
 * stops and removes the container
 */
func StopContainers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()
	vars := mux.Vars(r)
	name := vars["name"]
	ctr := models.CtrMgt{Name: name}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		if name == container.Names[0][1:] {
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

/*
 * stops and removes the container
 */
func GetContainerLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	name := vars["name"]
	//ctrs := make([]models.CtrMgt, 0)
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, ctr := range containers {
		if name == ctr.Names[0][1:] {
			log.Printf("Container Name %s , Id %s\n", name, ctr.ID)

			reader, err := cli.ContainerLogs(context.Background(), ctr.ID, types.ContainerLogsOptions{
				ShowStdout: true,
				ShowStderr: true,
				Follow:     false,
				Timestamps: false,
			})
			if err != nil {
				panic(err)
			}
			defer reader.Close()
			p := make([]byte, 8)
			reader.Read(p)
			content, _ := ioutil.ReadAll(reader)
			//ctrs = append(ctrs, models.CtrMgt{Logs: string(content)})
			var ctrMgt models.CtrMgt
			if err := json.NewDecoder(strings.NewReader(string(content))).Decode(&ctrMgt.Logs); err != nil {
				//handle error
			}
			//ctrs = append(ctrs, ctrMgt)
			//log.Println(string(content))

			/*scanner := bufio.NewScanner(content)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}*/

			json.NewEncoder(w).Encode(string(content))

		}
	}

}
