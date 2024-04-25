package graphdb

import (
	"context"
	"errors"
	"fmt"
	"go-userm/models"
	"log"
	"strconv"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func WriteUser(user *models.GraphDBUser, neo4jDriver neo4j.DriverWithContext) error {
	ctx := context.Background()
	session := neo4jDriver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	savedUser, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {

			existingUser, err := transaction.Run(ctx,
				"MATCH (p:User) WHERE p.id = $id OR p.username = $username RETURN p",
				map[string]any{"id": user.ID, "username": user.Username})
			if err != nil {
				return nil, err
			}

			if existingUser.Next(ctx) {
				return "User already exists", nil
			}

			result, err := transaction.Run(ctx,
				"CREATE (p:User) SET p.id = $id, p.username = $username RETURN p.username + ', from node ' + id(p)",
				map[string]any{"id": user.ID, "username": user.Username})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()

		})
	if err != nil {
		fmt.Println("Error inserting USER:", err)
		return err
	}

	fmt.Println(savedUser.(string))
	return nil
}

func FollowUser(followerUsername string, followeeUsername string, neo4jDriver neo4j.DriverWithContext) error {

	query := `
    MATCH (follower:User {username: $followerUsername}), (followee:User {username: $followeeUsername})
    MERGE (follower)-[:FOLLOWS]->(followee)
    `
	params := map[string]any{
		"followerUsername": followerUsername,
		"followeeUsername": followeeUsername,
	}

	ctx := context.Background()
	session := neo4jDriver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx,
		func(tx neo4j.ManagedTransaction) (any, error) {
			_, err := tx.Run(ctx,
				query, params,
			)
			if err != nil {
				return nil, err
			}

			return nil, nil
		})

	return err
}

func UnfollowUser(followerUsername string, followeeUsername string, neo4jDriver neo4j.DriverWithContext) error {
	ctx := context.Background()
	session := neo4jDriver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"MATCH (follower:User {username: $followerUsername})-[r:FOLLOWS]->(followee:User {username: $followeeUsername}) DELETE r",
				map[string]any{
					"followerUsername": followerUsername,
					"followeeUsername": followeeUsername,
				},
			)
			if err != nil {
				return nil, err
			}
			return result.Consume(ctx)
		})
	if err != nil {
		fmt.Printf("Error unfollowing user: %v\n", err)
		return err
	}

	log.Println("rezultat")
	log.Println(result)

	fmt.Printf("User %s unfollowed user %s successfully.\n", followerUsername, followeeUsername)
	return nil
}

func FindUserByID(id string, neo4jDriver neo4j.DriverWithContext) (*models.GraphDBUser, error) {

	// id je string, hmm
	// promeniti u integer
	fmt.Println("follower id")
	fmt.Println(id)
	// Convert the string ID to uint64
	uint64_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Printf("Error converting string ID to int: %v", err)
		return nil, err
	}

	ctx := context.Background()
	session := neo4jDriver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	userResult, err := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`
					MATCH (p:User) WHERE p.id = $id
					RETURN p.id AS id, p.username AS username
				`,
				map[string]any{
					"id": uint64_id})
			if err != nil {
				return nil, err
			}

			var user *models.GraphDBUser
			for result.Next(ctx) {
				record := result.Record()
				id, idOk := record.Get("id")
				username, usernameOk := record.Get("username")
				if idOk && usernameOk {
					user = &models.GraphDBUser{
						ID:       id.(int64),
						Username: username.(string),
					}
					return user, nil
				}
			}

			if err := result.Err(); err != nil {
				return nil, err
			}

			return nil, nil
		})
	if err != nil {
		log.Println("Error querying search:", err)
		return nil, err
	}

	if userResult == nil {
		log.Println("User not found")
		return nil, errors.New("user not found")
	}

	return userResult.(*models.GraphDBUser), nil
}

func DoesFollow(followerId string, followeeId string, neo4jDriver neo4j.DriverWithContext) (bool, error) {
	ctx := context.Background()
	session := neo4jDriver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	uint64_followerId, err := strconv.ParseUint(followerId, 10, 64)
	if err != nil {
		log.Printf("Error converting string ID to int: %v", err)
		return false, err
	}

	uint64_followeeId, err := strconv.ParseUint(followeeId, 10, 64)
	if err != nil {
		log.Printf("Error converting string ID to int: %v", err)
		return false, err
	}

	result, err := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`
                MATCH (follower:User)-[:FOLLOWS]->(followee:User)
                WHERE follower.id = $followerId AND followee.id = $followeeId
                RETURN COUNT(*) > 0 AS doesFollow
                `,
				map[string]any{
					"followerId": uint64_followerId,
					"followeeId": uint64_followeeId,
				})
			if err != nil {
				return nil, err
			}
			log.Println("usaooo 1")

			if result.Next(ctx) {
				record := result.Record()
				doesFollow, _ := record.Get("doesFollow")
				log.Println("usaooo 2")

				return doesFollow, nil
			}
			log.Println("usaooo 3")

			return false, nil
		})

	if err != nil {
		return false, err
	}
	return result.(bool), nil
}

func RecommendFriends(userId uint, neo4jDriver neo4j.DriverWithContext) ([]string, error) {
	ctx := context.Background()
	session := neo4jDriver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`
                MATCH (u:User)-[:FOLLOWS]->(friend)-[:FOLLOWS]->(fof)
				WHERE NOT (u)-[:FOLLOWS]->(fof)
				AND u.id = $userId
				AND fof <> u
				RETURN DISTINCT fof.username AS recommendation
                `,
				map[string]any{
					"userId": userId,
				})
			if err != nil {
				return nil, err
			}

			var recommendations []string
			for result.Next(ctx) {
				record := result.Record()
				recommendation, _ := record.Get("recommendation")
				recommendations = append(recommendations, recommendation.(string))
			}

			return recommendations, nil
		})

	if err != nil {
		return nil, err
	}
	return result.([]string), nil
}

func GetFollowees(userId uint, neo4jDriver neo4j.DriverWithContext) ([]int64, error) {
	ctx := context.Background()
	session := neo4jDriver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`
                MATCH (creator:User)-[:FOLLOWS]->(follower:User)
                WHERE creator.id = $userId
                RETURN follower.id as id
                `,
				map[string]any{
					"userId": userId,
				})
			if err != nil {
				return nil, err
			}

			var followees []int64
			for result.Next(ctx) {
				record := result.Record()
				followeeID, _ := record.Get("id")
				followees = append(followees, followeeID.(int64))
			}

			return followees, nil
		})

	if err != nil {
		return nil, err
	}
	return result.([]int64), nil
}
