package cassandra

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/CarlosMore29/icedoor_go/model"

	"github.com/gocql/gocql"
)

const (
	createQuery         = "INSERT INTO %s.%s (user_id, user_name , user_bcity) VALUES (?,?,?)"
	selectQuery         = "SELECT * FROM %s.%s where user_id = ?"
	findAllUsersQuery   = "SELECT * FROM %s.%s"
	cassandraTestQuery  = "SELECT * FROM %s.%s"
	cassandraTestInsert = "INSERT INTO %s.%s (id, data , date) VALUES (?,?,?)"
)

// GetSession connects to Cassandra
func GetSession(cosmosCassandraContactPoint, cosmosCassandraPort, cosmosCassandraUser, cosmosCassandraPassword string) (*gocql.Session, error) {
	var errorGlobal error = nil
	var session *gocql.Session = nil
	var portInt int = 12000

	clusterConfig := gocql.NewCluster(cosmosCassandraContactPoint)
	portInt, errorGlobal = strconv.Atoi(cosmosCassandraPort)

	var sslOptions = new(gocql.SslOptions)

	if errorGlobal == nil {
		clusterConfig.Port = portInt
		clusterConfig.ProtoVersion = 4
		clusterConfig.Authenticator = gocql.PasswordAuthenticator{Username: cosmosCassandraUser, Password: cosmosCassandraPassword}
		clusterConfig.SslOpts = &gocql.SslOptions{Config: &tls.Config{MinVersion: tls.VersionTLS12}}

		clusterConfig.ConnectTimeout = 10 * time.Second
		clusterConfig.Timeout = 10 * time.Second
		clusterConfig.DisableInitialHostLookup = true

		clusterConfig.SslOpts.EnableHostVerification = false
		clusterConfig.SslOpts = sslOptions
		// clusterConfig.Consistency = gocql.LocalOne

		// uncomment if you want to track time taken for individual queries
		//clusterConfig.QueryObserver = timer{}

		// uncomment if you want to track time taken for each connection to Cassandra
		//clusterConfig.ConnectObserver = timer{}

		session, errorGlobal = clusterConfig.CreateSession()
		if errorGlobal != nil {
			errorGlobal = errors.New("Failed to connect to Azure Cosmos DB: " + errorGlobal.Error())
			// log.Fatal("Failed to connect to Azure Cosmos DB", err)
		}
	}

	return session, errorGlobal
}

// InsertUser creates an entry(row) in a table
func InsertUser(keyspace, table string, session *gocql.Session, user model.User) (bool, error) {

	var created bool = true
	var errorGlobal error = nil

	errorGlobal = session.Query(fmt.Sprintf(createQuery, keyspace, table)).Bind(user.ID, user.Name, user.City).Exec()
	if errorGlobal != nil {
		created = false
	}
	return created, errorGlobal
}

// InsertUser creates an entry(row) in a table
func InsertTestCassandra(keyspace, table string, session *gocql.Session, timeline model.Timeline) (bool, error) {

	var created bool = true
	var errorGlobal error = nil

	errorGlobal = session.Query(fmt.Sprintf(cassandraTestInsert, keyspace, table)).Bind(timeline.ID, "charly", timeline.Date).Exec()
	if errorGlobal != nil {
		created = false
		fmt.Println(errorGlobal)
	}
	return created, errorGlobal
}

// FindUser tries to find a specific user
// func FindUser(keyspace, table string, id int, session *gocql.Session) model.User {
// 	var userid int
// 	var name, city string
// 	err := session.Query(fmt.Sprintf(selectQuery, keyspace, table)).Bind(id).Scan(&userid, &name, &city)

// 	if err != nil {
// 		if err == gocql.ErrNotFound {
// 			log.Printf("User with id %v does not exist\n", id)
// 		} else {
// 			log.Printf("Failed to find user with id %v - %v\n", id, err)
// 		}
// 	}
// 	return model.User{ID: userid, Name: name, City: city}
// }

// FindAllUsers gets all users
func FindAllUsers(keyspace, table string, session *gocql.Session) []model.User {

	var users []model.User
	results, _ := session.Query(fmt.Sprintf(findAllUsersQuery, keyspace, table)).Iter().SliceMap()

	for _, u := range results {
		users = append(users, mapToUser(u))
	}
	return users
}

// FindAllUsers gets all users
func FindAllCassandra(keyspace, table string, session *gocql.Session) []model.User {

	var users []model.User

	fmt.Println(fmt.Sprintf(cassandraTestQuery, keyspace, table))

	results, _ := session.Query(fmt.Sprintf(cassandraTestQuery, keyspace, table)).Iter().SliceMap()

	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println("Results", results)

	for _, u := range results {
		users = append(users, mapToUser(u))
	}
	return users
}

func mapToUser(m map[string]interface{}) model.User {
	id, _ := m["user_id"].(int)
	name, _ := m["user_name"].(string)
	city, _ := m["user_bcity"].(string)

	return model.User{ID: id, Name: name, City: city}
}
