package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"
)

/*
	A replica defines a single name:workers struct
*/

type Replica struct {
	Name    string `yaml:"name"`
	Workers string `yaml:"workers"`
}

/*
	Replica config is a unit in config file
*/

type ReplicaConfig struct {
	Peers []Replica `yaml:"peers"`
}

/*
	Returns a map[int32][]int32 that assigns each replica to its workers
*/

func NewReplicaConfig(fname string) (map[int32][]int32, error) {
	var cfg ReplicaConfig
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	err = yaml.UnmarshalStrict(data, &cfg)
	if err != nil {
		return nil, err
	}

	internalReplicaConfigs := make(map[int32][]int32)
	peers := cfg.Peers
	for i := 0; i < len(peers); i++ {
		peerName, _ := strconv.Atoi(peers[i].Name)
		peerWorkersStr := strings.Split(peers[i].Workers, ",")
		var peerWokrerInt32 []int32
		for j := 0; j < len(peerWorkersStr); j++ {
			k, _ := strconv.Atoi(peerWorkersStr[j])
			peerWokrerInt32 = append(peerWokrerInt32, int32(k))
		}
		internalReplicaConfigs[int32(peerName)] = peerWokrerInt32
	}

	return internalReplicaConfigs, nil
}
