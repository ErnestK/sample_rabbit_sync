package queries

import (
	"log"

	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateEventTable(eventCollection *mongo.Collection, row eventDb, fullLogMessage primitive.M) {
	log.Print("\n in UpdateEventTable func")
}

// If it has the attribute crit > 0 :
// – If an alarm already exists in the database for the same component / resource couple
// with status = ONGOING :
// *
// Then update crit, last_msg and last_time fields.
// – If no alarm is present in the database for the same component / resource couple with
// status = ONGOING :
// *

//
//
//
// • If the event has the attribute crit = 0 :
// – If an alarm already exists in the database for the same component / resource couple
// with status = ONGOING :
// *
// Then update crit, last_msg and last_time fields and change status to
// RESOLVED.
// – If no alarm is present in the database for the same component / resource couple with
// status = ONGOING :
// *
// Go to next message.
