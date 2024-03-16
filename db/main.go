package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/go-faker/faker/v4"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func dsn() string {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	hostname := os.Getenv("POSTGRES_HOST")
	dbName := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")
	_ = fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=require", user, password, hostname, dbName)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=require",
		hostname, port, user, password, dbName,
	)
	return psqlInfo
}

func main() {
	conn := dsn()
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	if err := m.Down(); err != nil {
		fmt.Println(err)
	}
	if err := m.Up(); err != nil {
		panic(err)
	}
	gen_fake_data(db)
}

// User represents the structure of the users table
type User struct {
	UserID         string `faker:"-"`
	FirstName      string `faker:"first_name"`
	LastName       string `faker:"last_name"`
	Email          string `faker:"-"`
	HashedPassword string `faker:"-"`
	SchoolID       int    `faker:"oneof: 1, 3, 4, 5, 7"` // Assuming school IDs are between 1 and 5
	DOBDay         int    `faker:"oneof: 1, 3, 4, 5, 7"`
	DOBMonth       int    `faker:"oneof: 1, 3, 4, 5, 7"`
	DOBYear        int    `faker:"oneof: 1991, 1993, 1995, 1994"`
	IsMale         bool   `faker:"-"`
	GradYear       int    `faker:"oneof: 2011, 2013, 2015, 2004"`
	Verified       bool   `faker:"-"`
	About          string `faker:"paragraph"`
}

type Pictures struct {
	userID string
	url    string
}

type Message struct {
	from_user_id string
	to_user_id   string
	txt_message  string
}

func gen_fake_data(db *sql.DB) {
	malePics := [][]string{
		{
			"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fcdn.britannica.com%2F22%2F188722-050-BB193EA3%2FRoger-Federer-US-Open-2007.jpg&f=1&nofb=1&ipt=90a21946e4212ea1501a58f7580b824f64339c037e537e169f48f3b6bef12e31&ipo=images",
			"https://images2.minutemediacdn.com/image/fetch/w_2000,h_2000,c_fit/https://lobandsmash.com/files/2016/10/8793945-stan-wawrinka-roger-federer-tennis-u.s.-open.jpg",
		},
		{
			"https://biografieonline.it/img/bio/Michael_Phelps_5.jpg",
			"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2F1.bp.blogspot.com%2F-p-exa7mcwVM%2FXYXEzOsPawI%2FAAAAAAAAF8c%2FJWb1jgiJwuUgdyJJzNl5-wW7dQR_XxE0gCLcBGAsYHQ%2Fs1600%2FMichael_Phelps.jpg&f=1&nofb=1&ipt=0c4147841e82edced230fdfbdd0f5eb67d96a6e542087f3b6d4021a324632b08&ipo=images",
		},
		{
			"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fwallpapercave.com%2Fwp%2Fwp2927877.jpg&f=1&nofb=1&ipt=dcbf7ddb5a491e86e90f58d44a855a657474624c4cdad31024a0df09c0d5035a&ipo=images",
			"https://external-content.duckduckgo.com/iu/?u=http%3A%2F%2Fa.espncdn.com%2Fphoto%2F2013%2F0429%2Fmag_kobe_17.jpg&f=1&nofb=1&ipt=525cdefdf97459c49772037a206c05ea6ef82a5dceb372a42265cae54ecd14f6&ipo=images",
		},
	}
	femalePics := [][]string{
		{
			"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Flive.staticflickr.com%2F7447%2F15805820133_ff9dba1444_b.jpg&f=1&nofb=1&ipt=f7663c2e28cae3c1a7ebf889c6ea3acbfd5a2e1323e7466c9b57c192ec5cc79b&ipo=images",
			"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fwww.itftennis.com%2Fremote.axd%2Fmedia.itftennis.com%2Fassetbank-itf%2Fservlet%2Fdisplay%3Ffile%3D22137ab7e9250948907cff18.jpg%3Fcrop%3D0.28080808080808078%2C0.048484848484848485%2C0.27676767676767677%2C0.57171717171717173%26cropmode%3Dpercentage%26width%3D219%26height%3D282%26rnd%3D132965530270000000&f=1&nofb=1&ipt=bb002e04ad5c0d25745ad77a2563ce6da61cf4e48d9f52cb929d24332df6e4a5&ipo=images",
		},
		{
			"https://external-content.duckduckgo.com/iu/?u=http%3A%2F%2Fcelebmafia.com%2Fwp-content%2Fuploads%2F2018%2F02%2Fgarbine-muguruza-qatar-wta-total-open-in-doha-02-16-2018-9.jpg&f=1&nofb=1&ipt=acdb4b42d22a60047fffd104a951517a563c73e647417f6dcec06d1c841a556a&ipo=images",
			"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fwww.gotceleb.com%2Fwp-content%2Fuploads%2Fphotos%2Fgarbine-muguruza%2Fwta-tennis-on-the-thames-evening-reception-in-london%2FGarbine-Muguruza%3A-WTA-Tennis-On-The-Thames-Evening-Reception--13.jpg&f=1&nofb=1&ipt=2139a2a6d19038444a1a079535bc24e772d5e4a2e7594a4d995354c3ce729828&ipo=images",
		},
		{
			"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fcelebmafia.com%2Fwp-content%2Fuploads%2F2016%2F07%2Fsimona-halep-wimbledon-tennis-championships-in-london-quarterfinals-7.jpg&f=1&nofb=1&ipt=3d3636be7ce3f534939b4fc25ca9252578edcd19c65b38394e254113908a2717&ipo=images",
			"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fcelebmafia.com%2Fwp-content%2Fuploads%2F2018%2F08%2Fsimona-halep-rogers-cup-in-montreal-08-08-2018-2.jpg&f=1&nofb=1&ipt=e674d95aa2603e1373a846fa169c227355736c044dc7866205b839c2d2ecc577&ipo=images",
		},
	}
	// Generate fake data for 5 users
	var fakeUsers []User
	var fakePics []Pictures
	r := rand.New(rand.NewSource(99))
	for i := 0; i < 20; i++ {
		fakeUser := User{}
		fakePic := Pictures{}
		err := faker.FakeData(&fakeUser)
		if err != nil {
			fmt.Println("Error generating fake data:", err)
			return
		}
		fakeUser.UserID = uuid.New().String()
		fakeUser.IsMale = i%2 == 1
		if fakeUser.IsMale {
			fakeUser.FirstName = faker.FirstNameMale()
		} else {
			fakeUser.FirstName = faker.FirstNameFemale()
		}
		fakeUser.Email = fmt.Sprintf("test%v@gmail.com", i+1)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123"), 10)
		if err != nil {
			return
		}
		fakeUser.HashedPassword = string(hashedPassword)
		if i < 11 {
			fakeUser.Verified = true
		}

		for j := 0; j < 2; j++ {
			fakePic.userID = fakeUser.UserID
			if fakeUser.IsMale {
				fakePic.url = malePics[i%3][j]
			} else {
				fakePic.url = femalePics[i%3][j]
			}
			fakePics = append(fakePics, fakePic)
		}
		fakeUsers = append(fakeUsers, fakeUser)
	}

	insert_user(db, fakeUsers)
	insert_pics(db, fakePics)
	insert_msgs(db, fakeUsers, r)
	insert_votes(db, fakeUsers, r)

	fmt.Println("Fake data inserted successfully.")
}

func insert_pics(db *sql.DB, fakePics []Pictures) {
	insertQuery := `
	INSERT INTO pictures (user_id, url)
	VALUES ($1, $2)
	`

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for _, pic := range fakePics {
		_, err := tx.Exec(
			insertQuery,
			pic.userID, pic.url,
		)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func insert_user(db *sql.DB, fakeUsers []User) {
	// Insert fake data into the users table
	insertQuery := `
	INSERT INTO users (user_id, first_name, last_name, email, hashed_password, school_id, dob_day, dob_month, dob_year, is_male, grad_year, verified, about)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range fakeUsers {
		_, err := tx.Exec(
			insertQuery,
			user.UserID, user.FirstName, user.LastName, user.Email, user.HashedPassword, user.SchoolID, user.DOBDay,
			user.DOBMonth, user.DOBYear, user.IsMale, user.GradYear, user.Verified, user.About,
		)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func insert_msgs(db *sql.DB, fakeUsers []User, r *rand.Rand) {
	var fakeMsgs []Message

	var activeUsers []User
	for _, u := range fakeUsers {
		if u.Verified {
			activeUsers = append(activeUsers, u)
		}
	}

	for i := 0; i < len(activeUsers); i++ {
		for j := 0; j < len(activeUsers); j++ {
			if i != j {
				for k := 1; k <= 35; k++ {
					a := r.Intn(len(activeUsers))
					b := r.Intn(len(activeUsers))
					m := Message{
						from_user_id: activeUsers[b].UserID,
						to_user_id:   activeUsers[a].UserID,
						txt_message:  faker.Sentence(),
					}
					fakeMsgs = append(fakeMsgs, m)
				}
			}
		}
	}

	insertQuery := `
	INSERT INTO messages (to_user_id, from_user_id, txt_message)
	VALUES ($1, $2, $3)
	`

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for _, msg := range fakeMsgs {
		_, err := tx.Exec(
			insertQuery,
			msg.to_user_id, msg.from_user_id, msg.txt_message,
		)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func insert_votes(db *sql.DB, fakeUsers []User, r *rand.Rand) {

	insertQuery := `
	INSERT INTO votes (prospective_user_id, user_id, vote_yes)
	VALUES ($1, $2, $3)
	`

	var verifiedUsers []User
	var unverified []User
	for _, u := range fakeUsers {
		if u.Verified {
			verifiedUsers = append(verifiedUsers, u)
		} else {
			unverified = append(unverified, u)
		}
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range verifiedUsers {
		for _, k := range unverified {
			vote := r.Intn(10)%3 == 0
			if _, err := tx.Exec(insertQuery, k.UserID, u.UserID, vote); err != nil {
				panic(err)
			}

		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
