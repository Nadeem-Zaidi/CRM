package main

// import (
// 	"github.com/Nadeem-Zaidi/CRM/database"

// )

// func main() {

// 	var product models.Product

// 	database.InitDB("mysql", "root", "owl", "CRM")
// 	product.All()

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Nadeem-Zaidi/CRM/database"
	"gorm.io/gorm"
)

// This function creates session and requires AWS credentials
// func CreateSession() *session.Session {
// 	sess := session.Must(session.NewSession(
// 		&aws.Config{
// 			Region: aws.String("ap-south-1"),
// 			Credentials: credentials.NewStaticCredentials(
// 				"AKIAQUJAFNH5HHHQMK5B",
// 				"U3mAS7t9gV9Ci1RMKWh2bIKGfYJYXnn4PXHQxTVa",
// 				"",
// 			),
// 		},
// 	))
// 	return sess
// }

// func CreateS3Session(sess *session.Session) *s3.S3 {
// 	s3Session := s3.New(sess)
// 	return s3Session
// }

// func UploadObject(bucket string, filePath string, fileName string, sess *session.Session) error {

// 	// Open file to upload
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	defer file.Close()

// 	// Upload to s3
// 	uploader := s3manager.NewUploader(sess)
// 	u, err := uploader.Upload(&s3manager.UploadInput{
// 		Bucket:      aws.String(bucket),
// 		Key:         aws.String(fileName),
// 		Body:        file,
// 		ContentType: aws.String("image/png"),
// 	})

// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	fmt.Println(u.Location)

//		fmt.Printf("Successfully uploaded %q to %q\n", fileName, bucket)
//		return nil
//	}
func Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Println(request.URL.Path)
		handler.ServeHTTP(writer, request)
	})
}

type User struct {
	Id      int
	Name    string
	Address string
	Phone   string
	Email   string
}

type Category struct {
	gorm.Model
	Id     int
	Name   string `json:"name"`
	Parent int
	Path   string
}

type Test struct {
	ID      int
	Name    string
	Address string
	Phone   string
	Email   string
}

func main() {
	// s := CreateSession()
	// UploadObject("nadeem-for-demo", "/home/owl/Desktop/work/ast.png", "ast", s)
	// database.InitDB("mysql", "root", "owl", "CRM")
	// var product models.Product

	// pr := product.All()
	// var vp models.VProduct
	// v := vp.SingleProduct(2)

	// router := chi.NewRouter()
	// router.Use(Logger)
	// router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
	// 	user, e := json.Marshal(v)
	// 	if e != nil {
	// 		fmt.Println(e)

	// 	}
	// 	writer.Header().Set("content-type", "application/json")
	// 	writer.Write(user)

	// })

	// err := http.ListenAndServe(":3000", router)
	// if err != nil {
	// 	log.Println(err)
	// }

	database.InitDB("mysql", "root", "owl", "CRM")

	// database.Create(User{})

	// database.Insert(User{ID: 12, Name: "Nadeem", Address: "Ali Ganj", Phone: "7992482690", Email: "nadeem.zaidi0021@gmail.com"})
	// var u User
	// database.FetchAll(&u)

	// dsn := "root:owl@tcp(127.0.0.1:3306)/CRM?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// var c Category
	// db.Create(&Category{Name: "Fashion"}).Scan(&c)
	database.InitDB("mysql", "root", "owl", "CRM")

	rw, _ := database.DB.Query("select id,name,address,phone,email from User")
	defer rw.Close()

	var c []User

	err := database.FindA(&c, rw)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("running after wards")
	fmt.Println(c)

}
