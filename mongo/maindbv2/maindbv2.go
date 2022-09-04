package maindbv2

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var server = "mongodb+srv://admin:NetZero1234@cluster0.vxsss.mongodb.net/admin"

// var server = "mongodb://127.0.0.1:15000"

func UpdateArchive(ctx context.Context, db_mongo_u string, collec_u string, input1 bson.M, input2 bson.M) string {

	time := time.Now().Add(time.Minute).Unix()
	clientOptions_u := options.Client().ApplyURI(server)
	client_u, err := mongo.Connect(ctx, clientOptions_u)
	if err != nil {
		log.Fatal(err)
	}
	err = client_u.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client_u.Disconnect(ctx)

	collection_u := client_u.Database(db_mongo_u).Collection(collec_u)
	// collection_i := client_u.Database(db_mongo_u).Collection(collec_u + `_Archive`)

	//cur, currErr := collection.Find(ctx, bson.M{})
	cur, currErr := collection_u.Find(ctx, input1)
	if currErr != nil {
		panic(currErr)
	}
	// fmt.Println(cur)

	var msg []bson.M

	if err = cur.All(ctx, &msg); err != nil {
		panic(err)
	}
	fmt.Println(msg)

	delete(msg[0], "_id")

	// res_i, insertErr := collection_i.InsertOne(ctx, msg[0])
	// if insertErr != nil {
	// 	return `nok`
	// }
	// fmt.Println(res_i)
	input2[`timestamp`] = time
	res_u, insertErr := collection_u.UpdateOne(ctx, input1, bson.M{"$set": input2})
	if insertErr != nil {
		client_u.Disconnect(ctx)
		return `nok`
	}
	fmt.Println(res_u)
	client_u.Disconnect(ctx)
	return `ok`
}

func Insertdb(ctx context.Context, db_mongo_u string, collec_u string, input1 bson.M) string {

	// db_mongo_u = "auth_main_demo2_07_21"
	// collec_u = "users_main_demo2_07_21"
	time := time.Now().Add(time.Minute).Unix()
	// var ctx = context.TODO()
	clientOptions_u := options.Client().ApplyURI(server)
	client_u, err := mongo.Connect(ctx, clientOptions_u)
	if err != nil {
		log.Fatal(err)
	}
	err = client_u.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client_u.Disconnect(ctx)

	collection_u := client_u.Database(db_mongo_u).Collection(collec_u)
	input1[`timestamp`] = time
	res_u, insertErr := collection_u.InsertOne(ctx, input1)
	if insertErr != nil {
		client_u.Disconnect(ctx)
		return `nok`
	}
	fmt.Println(res_u)
	client_u.Disconnect(ctx)
	return `ok`
}

func Finddb(ctx context.Context, db_mongo_u string, collec_u string, input1 bson.M, sortby string, sortorder int, limmit int64, skip int64) []bson.M {

	// db_mongo_u = "auth_main_demo2_07_21"
	// collec_u = "users_main_demo2_07_21"
	// var ctx = context.TODO()
	clientOptions_u := options.Client().ApplyURI(server)
	client_u, err := mongo.Connect(ctx, clientOptions_u)
	if err != nil {
		log.Fatal(err)
	}
	err = client_u.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client_u.Disconnect(ctx)

	collection_u := client_u.Database(db_mongo_u).Collection(collec_u)

	opts := options.Find()
	opts.SetSort(bson.D{{sortby, sortorder}}).SetLimit(limmit).SetSkip(skip)
	cur, err := collection_u.Find(ctx, input1, opts)
	//cur, err := collection.Find(ctx, bson.D{{}}, opts)
	if err != nil {
		var msg2 []bson.M
		client_u.Disconnect(ctx)
		return msg2
	}
	var msg []bson.M
	if err = cur.All(ctx, &msg); err != nil {
		panic(err)
	}
	client_u.Disconnect(ctx)
	return msg
}

func UpdatePushArray(ctx context.Context, db_mongo_u string, collec_u string, input1 bson.M, input2 bson.M, input3 string) string {

	// Month := time.Now().Month()
	// Year := time.Now().Year()
	clientOptions_u := options.Client().ApplyURI(server)
	client_u, err := mongo.Connect(ctx, clientOptions_u)
	if err != nil {
		log.Fatal(err)
	}
	err = client_u.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client_u.Disconnect(ctx)

	collection_u := client_u.Database(db_mongo_u).Collection(collec_u)

	//----------------------------------------------------------------------

	cur, err := collection_u.Find(ctx, input1)
	var msg []bson.M
	if err = cur.All(ctx, &msg); err != nil {
		panic(err)
	}

	if len(msg) > 0 {
		fmt.Println(`have exited`)
	} else {
		res_ins, insertErr := collection_u.InsertOne(ctx, input1)
		if insertErr != nil {
			client_u.Disconnect(ctx)
			return `nok`
		}
		fmt.Println(res_ins)
	}

	Month := int(time.Now().Month())
	Year := time.Now().Year()
	setdata := input3 + `.` + strconv.Itoa(Month) + `-` + strconv.Itoa(Year)

	fmt.Println(setdata)
	res_u, insertErr := collection_u.UpdateOne(ctx, input1, bson.M{"$push": bson.M{setdata: input2}})
	if insertErr != nil {
		client_u.Disconnect(ctx)
		return `nok`
	}
	fmt.Println(res_u)
	client_u.Disconnect(ctx)
	return `ok`
}

func UpdatePushArraycus(ctx context.Context, db_mongo_u string, collec_u string, input1 bson.M, input2 bson.M, input3 string) string {

	// Month := time.Now().Month()
	// Year := time.Now().Year()
	clientOptions_u := options.Client().ApplyURI(server)
	client_u, err := mongo.Connect(ctx, clientOptions_u)
	if err != nil {
		log.Fatal(err)
	}
	err = client_u.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client_u.Disconnect(ctx)

	collection_u := client_u.Database(db_mongo_u).Collection(collec_u)

	//----------------------------------------------------------------------

	cur, err := collection_u.Find(ctx, input1)
	var msg []bson.M
	if err = cur.All(ctx, &msg); err != nil {
		panic(err)
	}

	if len(msg) > 0 {
		fmt.Println(`have exited`)
	} else {
		res_ins, insertErr := collection_u.InsertOne(ctx, input1)
		if insertErr != nil {
			return `nok`
		}
		fmt.Println(res_ins)
	}

	// Month := int(time.Now().Month())
	// Year := time.Now().Year()
	// setdata := strconv.Itoa(Month) + `-` + strconv.Itoa(Year)

	// fmt.Println(setdata)
	res_u, insertErr := collection_u.UpdateOne(ctx, input1, bson.M{"$push": bson.M{input3: input2}})
	if insertErr != nil {
		client_u.Disconnect(ctx)
		return `nok`
	}
	fmt.Println(res_u)
	client_u.Disconnect(ctx)
	return `ok`
}

func Findonly(ctx context.Context, db_mongo_u string, collec_u string, input1 bson.M, key string) []bson.M {

	// db_mongo_u = "auth_main_demo2_07_21"
	// collec_u = "users_main_demo2_07_21"
	// var ctx = context.TODO()
	clientOptions_u := options.Client().ApplyURI(server)
	client_u, err := mongo.Connect(ctx, clientOptions_u)
	if err != nil {
		log.Fatal(err)
	}
	err = client_u.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client_u.Disconnect(ctx)

	collection_u := client_u.Database(db_mongo_u).Collection(collec_u)

	opts := options.Find().SetProjection(bson.M{key: 1, "_id": 0})
	// opts.SetSort(bson.D{{sortby, sortorder}}).SetLimit(limmit).SetSkip(skip)

	cur, err := collection_u.Find(ctx, input1, opts)
	//cur, err := collection.Find(ctx, bson.D{{}}, opts)

	fmt.Println(cur)
	if err != nil {
		var msg2 []bson.M
		client_u.Disconnect(ctx)
		return msg2
	}
	var msg []bson.M
	if err = cur.All(ctx, &msg); err != nil {
		panic(err)
	}
	client_u.Disconnect(ctx)
	return msg
}

func Findmutikey(ctx context.Context, db_mongo_u string, collec_u string, input1 bson.M, key []string) []bson.M {

	// db_mongo_u = "auth_main_demo2_07_21"
	// collec_u = "users_main_demo2_07_21"
	// var ctx = context.TODO()
	clientOptions_u := options.Client().ApplyURI(server)
	client_u, err := mongo.Connect(ctx, clientOptions_u)
	if err != nil {
		log.Fatal(err)
	}
	err = client_u.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client_u.Disconnect(ctx)

	collection_u := client_u.Database(db_mongo_u).Collection(collec_u)

	dataneed := make(bson.M)
	dataneed[`_id`] = 0

	for i := 0; i < len(key); i++ {
		dataneed[key[i]] = 1
	}

	fmt.Println(dataneed)

	opts := options.Find().SetProjection(dataneed)
	// opts.SetSort(bson.D{{sortby, sortorder}}).SetLimit(limmit).SetSkip(skip)

	cur, err := collection_u.Find(ctx, input1, opts)
	//cur, err := collection.Find(ctx, bson.D{{}}, opts)

	fmt.Println(cur)
	if err != nil {
		var msg2 []bson.M
		client_u.Disconnect(ctx)
		return msg2
	}
	var msg []bson.M
	if err = cur.All(ctx, &msg); err != nil {
		panic(err)
	}
	client_u.Disconnect(ctx)
	return msg
}
