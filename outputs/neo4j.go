package qoutput

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/deckarep/golang-set"
	dtypes "github.com/docker/docker/api/types"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/zpatrick/go-config"

	"github.com/qnib/qwatch/types"
	"github.com/qnib/qwatch/utils"
)

var (
	neo4jURL = "bolt://localhost:7687"
)

// Neo4j holds a buffer and the initial information from the server
type Neo4j struct {
	qtypes.QWorker
	Conn bolt.Conn
}

// NewNeo4j returns an initial instance
func NewNeo4j(cfg *config.Config, qC qtypes.Channels) Neo4j {
	o := Neo4j{}
	o.Cfg = cfg
	o.QChan = qC
	o.Subs = mapset.NewSet()
	return o
}

func (o Neo4j) execCypher(cypher string, m map[string]interface{}) error {
	stmt, err := o.Conn.PrepareNeo(cypher)
	defer stmt.Close()
	if err != nil {
		log.Printf("[EE] during PrepareNeo (%s): %s", cypher, err)
		return err
	}
	_, err = stmt.ExecNeo(m)
	if err != nil {
		log.Printf("[EE] during ExecNeo: (%s): %s", cypher, err)
		return err
	}
	return nil
}

func (o Neo4j) handleContainer(qm qtypes.Qmsg) error {
	switch qm.Action {
	case "create":
		cypher := `
		MATCH (s:ContainerState {name: 'created'}) MATCH (de:DockerEngine {id:{engine_id}})
			CREATE UNIQUE (de)<-[:PartOf]-(c:Container {name: {name}, container_id: {container_id}, created: {time}})<-[:IS {created: {time}}]-(s)`
		m := map[string]interface{}{"name": qm.Container.ContainerName, "time": qm.TimeNano}
		m["container_id"] = qm.Container.ContainerID
		m["engine_id"] = qm.EngineID
		log.Printf("[DD] Cypher: '%s', map:'%v'", cypher, m)
		err := o.execCypher(cypher, m)
		if err != nil {
			log.Println("[EE] during ExecCypher: ", err)
			return err
		}
		cypher = "MATCH (c:Container {container_id: {container_id}}) MERGE (i:DockerImage {id: {image_id}}) CREATE (c)-[:USE {time: {time}}]->(i)"
		m = map[string]interface{}{"container_id": qm.Container.ContainerID, "image_id": qm.Container.ImageName, "time": qm.TimeNano}
		log.Printf("[DD] Cypher: %s, map:%v", cypher, m)
		return o.execCypher(cypher, m)
	case "start":
		cypher := "MATCH (s:ContainerState {name: 'running'}) MATCH (c:Container {container_id: {container_id}}) CREATE (c)<-[:IS {time: {time}}]-(s)"
		m := map[string]interface{}{"container_id": qm.Container.ContainerID, "time": qm.TimeNano}
		log.Printf("[DD] Cypher: %s, map:%v", cypher, m)
		return o.execCypher(cypher, m)
	case "die":
		cypher := "MATCH (s:ContainerState {name: 'dead'}) MATCH (c:Container {container_id: {container_id}}) CREATE (c)<-[:IS {time: {time}}]-(s)"
		m := map[string]interface{}{"container_id": qm.Container.ContainerID, "time": qm.TimeNano}
		log.Printf("[DD] Cypher: %s, map:%v", cypher, m)
		return o.execCypher(cypher, m)
	case "pause":
		cypher := "MATCH (s:ContainerState {name: 'paused'}) MATCH (c:Container {container_id: {container_id}}) CREATE (c)<-[:IS {time: {time}}]-(s)"
		m := map[string]interface{}{"container_id": qm.Container.ContainerID, "time": qm.TimeNano}
		log.Printf("[DD] Cypher: %s, map:%v", cypher, m)
		return o.execCypher(cypher, m)
	case "kill":
		cypher := "MATCH (s:ContainerState {name: 'killed'}) MATCH (c:Container {container_id: {container_id}}) CREATE (c)<-[:IS {time: {time}}]-(s)"
		m := map[string]interface{}{"container_id": qm.Container.ContainerID, "time": qm.TimeNano}
		log.Printf("[DD] Cypher: %s, map:%v", cypher, m)
		return o.execCypher(cypher, m)
	case "stop":
		cypher := "MATCH (s:ContainerState {name: 'stopped'}) MATCH (c:Container {container_id: {container_id}}) CREATE (c)<-[:IS {time: {time}}]-(s)"
		m := map[string]interface{}{"container_id": qm.Container.ContainerID, "time": qm.TimeNano}
		log.Printf("[DD] Cypher: %s, map:%v", cypher, m)
		return o.execCypher(cypher, m)
	case "restart":
		cypher := "MATCH (s:ContainerState {name: 'restarted'}) MATCH (c:Container {container_id: {container_id}}) CREATE (c)<-[:IS {time: {time}}]-(s)"
		m := map[string]interface{}{"container_id": qm.Container.ContainerID, "time": qm.TimeNano}
		log.Printf("[DD] Cypher: %s, map:%v", cypher, m)
		return o.execCypher(cypher, m)
	case "unpause":
		cypher := "MATCH (s:ContainerState {name: 'unpaused'}) MATCH (c:Container {container_id: {container_id}}) CREATE (c)<-[:IS {time: {time}}]-(s)"
		m := map[string]interface{}{"container_id": qm.Container.ContainerID, "time": qm.TimeNano}
		log.Printf("[DD] Cypher: %s, map:%v", cypher, m)
		return o.execCypher(cypher, m)
	case "destroy":
		cypher := "MATCH (s:ContainerState {name: 'removed'}) MATCH (c:Container {container_id: {container_id}}) SET c.destroyed={time} CREATE (c)<-[:IS {time: {time}}]-(s)"
		m := map[string]interface{}{"container_id": qm.Container.ContainerID, "time": qm.TimeNano}
		log.Printf("[DD] Cypher: %s, map:%v", cypher, m)
		return o.execCypher(cypher, m)
	default:
		log.Printf("[II] Action is not recognized: %s", qm.Action)
		return nil
	}
}

func (o Neo4j) handleImg(qm qtypes.Qmsg) error {
	switch qm.Action {
	case "pull":
		img := utils.ParseImageName(qm.Image.ID)
		cypher := fmt.Sprintf("MERGE (i:DockerTag {name:'%s'})", img.String())
		log.Printf("[DD] Cypher: %s", cypher)
		return o.execCypher(cypher, nil)
	}
	return nil
}
func (o Neo4j) handleMsg(qm qtypes.Qmsg) error {
	switch qm.Type {
	case "image":
		return o.handleImg(qm)
	case "container":
		return o.handleContainer(qm)
	default:
		log.Printf("[II] Type is not recognized: %s", qm.Type)
		return nil
	}
}

// init the graph and create nodes that one needs
func (o Neo4j) initGraph() error {
	// https://github.com/docker/docker/blob/master/api/types/types.go#L276
	cypher := `
        MERGE (:ContainerState { name:'created'})
        MERGE (:ContainerState { name:'running'})
		MERGE (:ContainerState { name:'paused'})
		MERGE (:ContainerState { name:'stopped'})
		MERGE (:ContainerState { name:'unpaused'})
		MERGE (:ContainerState { name:'restarted'})
		MERGE (:ContainerState { name:'killed'})
		MERGE (:ContainerState { name:'oomkilled'})
		MERGE (:ContainerState { name:'dead'})
		MERGE (:ContainerState { name:'removed'})`
	err := o.execCypher(cypher, nil)
	o.execCypher("CREATE CONSTRAINT ON (i:DockerImage) ASSERT i.id IS UNIQUE", nil)

	/* When implementing services
	// https://github.com/docker/docker/blob/master/api/types/types.go#L255
	*/
	if err != nil {
		return err
	}
	return nil
}

func (o Neo4j) handleInfo(i dtypes.Info) {
	cypher := `
	MERGE (n:DockerEngine {id:{id}, name:{name}})
		ON MATCH SET n.last_seen={time}
		ON CREATE SET n.created={time},n.last_seen={time}
	`
	m := map[string]interface{}{"id": i.ID, "name": i.Name}
	m["time"] = time.Now().UnixNano()
	err := o.execCypher(cypher, m)
	if err != nil {
		log.Printf("[EE] Error during handleInfo: '%s' '%v'\n", err, i)
	}
}

func (o Neo4j) createDockerImage(i qtypes.DockerImageSummary) {
	cypher := `MATCH (de:DockerEngine {id: {engine_id}})
	MERGE (de)<-[:Exists]-(i:DockerImage {id: {id}})`
	m := map[string]interface{}{"id": i.ID, "created": i.Created, "time": time.Now().UnixNano()}
	m["engine_id"] = i.EngineID
	err := o.execCypher(cypher, m)
	if err != nil {
		log.Printf("[EE] Error during image create: '%s' '%v'\n", err, m)
	}
	for _, repoTag := range i.RepoTags {
		if repoTag == "<none>:<none>" {
			continue
		}
		img := utils.ParseImageName(repoTag)
		m := map[string]interface{}{"id": i.ID, "repo_tag": img.String()}
		cypher := `MATCH (ni:DockerImage {id: {id}})
        MERGE (ni)<-[:IS]-(t:ImageTag {name: {repo_tag}})`
		o.execCypher(cypher, m)
	}

}

func (o Neo4j) renameRelationship(id, repoTag string) {
	// Change relationship from IS to WAS

	/*cypher = `MATCH (t:ImageTag {name: {repo_tag}})
	    MATCH (t)-[:IS]->(oi)`
		o.execCypher(cypher, m)
	*/
}

func (o Neo4j) handleDockerImageSummary(i qtypes.DockerImageSummary) {
	o.createDockerImage(i)
	if i.ParentID != "" {
		cypher := `MATCH (de:DockerEngine {id: {engine_id}})
		MERGE (de)<-[:Exists]-(i:DockerImage {id: {id}})`
		m := map[string]interface{}{"id": i.ParentID}
		m["engine_id"] = i.EngineID
		err := o.execCypher(cypher, m)
		if err != nil {
			log.Printf("[EE] Error during image create: '%s' '%v'\n", err, m)
		}
		cypher = `MATCH (i:DockerImage {id: {id}})
		MATCH (p:DockerImage {id: {parent_id}})
		MERGE (i)-[:PARENT]->(p)`
		m = map[string]interface{}{"id": i.ID, "parent_id": i.ParentID}
		m["engine_id"] = i.EngineID
		err = o.execCypher(cypher, m)
		if err != nil {
			log.Printf("[EE] Error during handleSwarmNode: '%s' '%v'\n", err, m)
			os.Exit(1)
		}
	}
}

func (o Neo4j) handleSwarmNode(n qtypes.DockerNode) {
	cypher := `
	MATCH (d:DockerEngine {id:{engine_id}})
	MERGE (n:SwarmNode {id:{id}, name:{name}})-[:PartOf]->(d)
		ON MATCH SET n.last_seen={time}, n.node_status={node_status}
		ON CREATE SET n.addr={node_addr},n.created={created},n.last_seen={time}
	`
	m := map[string]interface{}{"id": n.ID, "name": n.Description.Hostname}
	m["engine_id"] = n.EngineID
	m["node_status"] = string(n.Status.State)
	m["node_addr"] = n.Status.Addr
	m["created"] = n.CreatedAt.UnixNano()
	m["time"] = time.Now().UnixNano()
	//log.Printf("[DD] Cypher: '%s', map:'%v'", cypher, m)
	err := o.execCypher(cypher, m)
	if err != nil {
		log.Printf("[EE] Error during handleSwarmNode: '%s' '%v'\n", err, n)
	}
}

// Run prints the logs to stdout
func (o Neo4j) Run() {
	driver := bolt.NewDriver()
	var err error
	o.Conn, err = driver.OpenNeo("bolt://localhost:7687")
	if err != nil {
		panic(err)
	}
	defer o.Conn.Close()
	err = o.initGraph()
	if err != nil {
		panic(err)
	}
	bg := o.QChan.Log.Join()
	ig := o.QChan.Inventory.Join()
	for {
		select {
		case val := <-bg.In:
			log := val.(qtypes.Qmsg)
			o.handleMsg(log)
		case val := <-ig.In:
			switch val := val.(type) {
			case qtypes.DockerNode:
				o.handleSwarmNode(val)
			case qtypes.DockerImageSummary:
				o.handleDockerImageSummary(val)
			case dtypes.Info:
				o.handleInfo(val)
			default:
				log.Printf("[WW] Do not recognise: %v", reflect.TypeOf(val))
			}
		}
	}
}
