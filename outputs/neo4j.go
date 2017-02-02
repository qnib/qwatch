package qoutput

import (
	"fmt"
	"log"

	"github.com/deckarep/golang-set"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/zpatrick/go-config"

	"github.com/qnib/qwatch/types"
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
		log.Println("[EE] during PrepareNeo: ", cypher)
		return err
	}
	_, err = stmt.ExecNeo(m)
	if err != nil {
		log.Println("[EE] during ExecNeo: ", err)
		return err
	}
	return nil
}

func (o Neo4j) handleContainer(qm qtypes.Qmsg) error {
	switch qm.Action {
	case "create":
		cypher := "MATCH (s:ContainerState {name: 'created'}) CREATE UNIQUE (c:Container {name: {name}, container_id: {container_id}, created: {time}})<-[:IS {created: {time}}]-(s)"
		m := map[string]interface{}{"name": qm.Container.ContainerName, "container_id": qm.Container.ContainerID, "time": qm.TimeNano}
		log.Printf("[DD] Cypher: '%s', map:'%v'", cypher, m)
		err := o.execCypher(cypher, m)
		if err != nil {
			log.Println("[EE] during ExecCypher: ", err)
			return err
		}
		cypher = "MATCH (c:Container {container_id: {container_id}}) MERGE (i:Image {id: {image_id}}) CREATE (c)-[:USE {time: {time}}]->(i)"
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
		cypher := fmt.Sprintf("MERGE (i:Image {id:'%s'})", qm.Image.ID)
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
	/* When implementing services
	// https://github.com/docker/docker/blob/master/api/types/types.go#L255
	*/
	if err != nil {
		return err
	}
	return nil
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
	for {
		select {
		case val := <-bg.In:
			log := val.(qtypes.Qmsg)
			o.handleMsg(log)
		}
	}
}
