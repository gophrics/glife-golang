package profile

import (
	"context"
	"fmt"
	"math/rand"

	"../../common"
	"../../common/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GenerateProfileId() string {
	return primitive.NewObjectID().Hex()
}

var maxNum = 100

func _GenerateUsername() string {
	var i = 0
	for i < 5 {
		var username = UsernameList[rand.Intn(88)] + fmt.Sprintf("%d", rand.Intn(maxNum))
		filter := bson.D{{"username", username}}
		err := mongodb.Profile.FindOne(context.TODO(), filter)
		if err != nil {
			return username
		}
		i++
	}
	maxNum++
	return UsernameList[rand.Intn(88)] + fmt.Sprintf("%d", rand.Intn(maxNum))
}

func GetEmailFromProfileId() string {
	return "stub"
}

func _RegisterUser(req common.User) (common.User, *mongo.InsertOneResult, error) {

	var result common.User
	result.Country = req.Country
	result.Email = req.Email
	result.Name = req.Name
	result.Phone = req.Phone
	result.ProfileId = GenerateProfileId()
	// BIG TODO: Hash Password
	// TODO: Assuming single email, that need not be the case, user can have multiple emails linked to same account
	// For example, registration with a non google email and trying to register later with a google email
	insertResult, err := mongodb.Profile.InsertOne(context.TODO(), req)

	return result, insertResult, err
}

var UsernameList = []string{
	"Playboyize",
	"Gigastrength",
	"Deadlyinx",
	"Techpill",
	"Methshot",
	"Methnerd",
	"TreeEater",
	"PackManBrainlure",
	"Carnalpleasure",
	"Sharpcharm",
	"Snarelure",
	"Skullbone",
	"Burnblaze",
	"Emberburn",
	"Emberfire",
	"Evilember",
	"Firespawn",
	"Flameblow",
	"SniperGod",
	"TalkBomber",
	"SniperWish",
	"RavySnake",
	"WebTool",
	"TurtleCat",
	"BlogWobbles",
	"LuckyDusty",
	"RumChicken",
	"StonedTime",
	"CouchChiller",
	"VisualMaster",
	"DeathDog",
	"ZeroReborn",
	"TechHater",
	"eGremlin",
	"BinaryMan",
	"AwesomeTucker",
	"FastChef",
	"JunkTop",
	"PurpleCharger",
	"CodeBuns",
	"BunnyJinx",
	"GoogleCat",
	"StrangeWizard",
	"Fuzzy_Logic",
	"New_Cliche",
	"Ignoramus",
	"Stupidasole",
	"whereismyname",
	"Nojokur",
	"Illusionz",
	"Spazyfool",
	"Supergrass",
	"Dreamworx",
	"Fried_Sushi",
	"Stark_Naked",
	"Need_TLC",
	"Raving_Cute",
	"Nude_Bikergirl",
	"Lunatick",
	"Garbage_Can_Lid",
	"Crazy_Nice",
	"Booteefool",
	"Harmless_Venom",
	"Lord_Tryzalot",
	"Sir_Gallonhead",
	"Boy_vs_Girl",
	"MPmaster",
	"King_Martha",
	"Spamalot",
	"Soft_member",
	"girlDog",
	"Evil_kitten",
	"farquit",
	"viagrandad",
	"happy_sad",
	"haveahappyday",
	"SomethingNew",
	"5mileys",
	"cum2soon",
	"takes2long",
	"w8t4u",
	"askme",
	"Bidwell",
	"massdebater",
	"iluvmen",
	"Inmate",
	"idontknow",
	"likme",
}
