package main

import (
    "context"
	"log"
    "reflect"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"
)

//User used in other .go files
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    VK_ID     int64              `bson:"vk_id" json:"id"`
    FirstName string             `bson:"first_name,omitempty" json:"first_name"`
    LastName  string             `bson:"last_name,omitempty" json:"last_name"`
    Photo     string             `bson:"photo_400_orig,omitempty" json:"photo_400_orig"`
    City      City               `bson:"city,omitempty" json:"city"`
}

//City used in User struct
type City struct {
    Title string `json:"title"`
}

//Post used in other .go files
type Post struct {
    ID    primitive.ObjectID `bson:"_id,omitempty"`
    User  primitive.ObjectID `bson:"user_id,omitempty"`
    Text  string             `bson:"text,omitempty"`
    Image string             `bson:"image,omitempty"`
    Link  string             `bson:"link,omitempty"`
}


func (u *User) createUser() {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
    
    if err != nil {
        log.Fatal(err)
    }

    collection := client.Database("wall-project").Collection("users")
    var result User
    err = collection.FindOne(context.TODO(), bson.D{{"vk_id", u.VK_ID}}).Decode(&result)

    if err != nil {

        _, err := collection.InsertOne(context.TODO(), u)

        if err != nil {
            log.Fatal(err)
        }
    } else {
        diff := reflect.DeepEqual(map[string]string{
            "first_name":u.FirstName,
            "last_name":u.LastName,
            "photo_400_orig":u.Photo,
        },map[string]string{
            "first_name":result.FirstName,
            "last_name":result.LastName,
            "photo_400_orig":result.Photo,
        })

        if !diff {
            update := bson.D{{"$set", bson.D{
                {"first_name",u.FirstName},
                {"last_name",u.LastName},
                {"photo_400_orig",u.Photo}}}}

            _, err = collection.UpdateOne(context.TODO(), bson.D{{"vk_id", u.VK_ID}}, update)
            if err != nil {
                log.Fatal(err)
            }
        }
    }

    defer func() {
        if err = client.Disconnect(context.TODO()); err != nil {
            log.Fatal(err)
        }
    }()
}

func (u *User) findUser() User {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
    
    if err != nil { log.Fatal(err) }

    coll := client.Database("wall-project").Collection("users")

    var user User

    if u.VK_ID != 0 {
        err = coll.FindOne(context.TODO(), bson.D{{"vk_id", u.VK_ID}}).Decode(&user)
    
        if err != nil { log.Fatal(err) }
    } else {
        err = coll.FindOne(context.TODO(), bson.D{{"_id", u.ID}}).Decode(&user)
    
        if err != nil { log.Fatal(err) }
    }

    defer func() {
        if err = client.Disconnect(context.TODO()); err != nil {
            log.Fatal(err)
        }
    }()

    return user
}

func (u *User) deleteUser() {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
    
    if err != nil {
        log.Fatal(err)
    }

    collection := client.Database("wall-project").Collection("users")

    _, err = collection.DeleteOne(context.TODO(), bson.D{{"_id", u.ID}})

    if err != nil {
        log.Fatal(err)
    }

    collection = client.Database("wall-project").Collection("posts")

    _, err = collection.DeleteMany(context.TODO(), bson.D{{"user_id", u.ID}})

    if err != nil {
        log.Fatal(err)
    }

    defer func() {
        if err = client.Disconnect(context.TODO()); err != nil {
            log.Fatal(err)
        }
    }()
}

func (p *Post) createPost() {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
    
    if err != nil {
        log.Fatal(err)
    }

    collection := client.Database("wall-project").Collection("posts")

    _, err = collection.InsertOne(context.TODO(), bson.D{{"user_id", p.User}, {"text", p.Text}, {"image", p.Image}, {"link", p.Link}})
    //id := res.InsertedID

    if err != nil {
        log.Fatal(err)
    }

    defer func() {
        if err = client.Disconnect(context.TODO()); err != nil {
            log.Fatal(err)
        }
    }()
}

func deletePost(id string) {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
    
    if err != nil {
        log.Fatal(err)
    }

    collection := client.Database("wall-project").Collection("posts")

    pID, err := primitive.ObjectIDFromHex(id)

    if err != nil {
        log.Fatal(err)
    }

    _, err = collection.DeleteOne(context.TODO(), bson.D{{"_id", pID}})

    if err != nil {
        log.Fatal(err)
    }

    defer func() {
        if err = client.Disconnect(context.TODO()); err != nil {
            log.Fatal(err)
        }
    }()
}

func findPosts(search string) []Post {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
    
    if err != nil {
        log.Fatal(err)
    }

    collection := client.Database("wall-project").Collection("posts")
    var cur *mongo.Cursor

    if search != "" {
        iv := collection.Indexes()

        models := []mongo.IndexModel{
            {
                Keys: bson.D{{"text", "text"}, {"link", "text"}},
            },
        }

        _, err := iv.CreateMany(context.TODO(), models)
        if err != nil {
            log.Fatal(err)
        }

        cur, err = collection.Find(context.Background(), bson.M{"$text": bson.M{"$search": search}})
    } else {
        cur, err = collection.Find(context.Background(), bson.D{})
    }
    
    if err != nil { log.Fatal(err) }
    defer cur.Close(context.Background())

    var results []Post
    if err = cur.All(context.Background(), &results); err != nil {
      log.Fatal(err)
    }

    defer func() {
        if err = client.Disconnect(context.TODO()); err != nil {
            log.Fatal(err)
        }
    }()

    return results
}