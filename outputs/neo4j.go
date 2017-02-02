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

func (o Neo4j) execCypher(cypher string) error {
	stmt, err := o.Conn.PrepareNeo(cypher)
	defer stmt.Close()
	if err != nil {
		log.Println("[EE] during PrepareNeo: ", cypher)
		return err
	}
	_, err = stmt.ExecNeo(nil)
	if err != nil {
		log.Println("[EE] during ExecNeo: ", err)
		return err
	}
	return nil
}

func (o Neo4j) handleContainer(qm qtypes.Qmsg) error {
	switch qm.Action {
	case "create":
		cypher := fmt.Sprintf("CREATE (c:Container {name:'%s', container_id: '%s', created:'%d' })", qm.Container.ContainerName, qm.Container.ContainerID, qm.TimeNano)
		log.Printf("[DD] Cypher: %s", cypher)
		err := o.execCypher(cypher)
		if err != nil {
			return err
		}
		cypher = fmt.Sprintf("MATCH (c:Container {container_id: '%s'}) MERGE (i:Image {id:'%s'}) CREATE (c)-[:USE]->(i)",
			qm.Container.ContainerID, qm.Container.ImageName)
		log.Printf("[DD] Cypher: %s", cypher)
		return o.execCypher(cypher)

	case "start":
		return nil
	case "die":
		return nil
	case "destroy":
		cypher := fmt.Sprintf("MATCH (c:Container {container_id: '%s' }) SET c.destroyed='%d'", qm.Container.ContainerID, qm.TimeNano)
		log.Printf("[DD] Cypher: %s", cypher)
		return o.execCypher(cypher)
	default:
		return nil
	}
}

func (o Neo4j) handleImg(qm qtypes.Qmsg) error {
	switch qm.Action {
	case "pull":
		cypher := fmt.Sprintf("MERGE (i:Image {id:'%s'})", qm.Image.ID)
		log.Printf("[DD] Cypher: %s", cypher)
		return o.execCypher(cypher)
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
			MERGE (:CONTAINER_STATE { name:'running'})
			MERGE (:CONTAINER_STATE { name:'paused'})
			MERGE (:CONTAINER_STATE { name:'restarting'})
			MERGE (:CONTAINER_STATE { name:'oomkilled'})
			MERGE (:CONTAINER_STATE { name:'dead'})`
	err := o.execCypher(cypher)
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
