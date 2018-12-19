package main

import (
	"fmt"

	"github.com/alpacahq/crypto-poc/models"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kataras/iris"
)

const (
	host     = "ec2-23-21-65-173.compute-1.amazonaws.com"
	port     = 5432
	user     = "oofujxujwofryj"
	password = "09cc6443bd76f26886baa881c19223b86faef509e6e3ac3e4a2385b1b18b7ef5"
	dbname   = "df3som0f9fs9jc"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := iris.New()
	app.RegisterView(iris.HTML("./views", ".html"))
	app.Get("/", func(ctx iris.Context) {
		ctx.View("register.html")
	})

	app.Post("/addAccount", func(ctx iris.Context) {
		// add random account to database
		newAcc, err := createNewAccount(db, ctx.PostValue("email"), ctx.PostValue("name"))
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Writef("error creating new Account")
			return
		}
		ctx.JSON(newAcc)
	})

	app.Get("/viewAccounts", func(ctx iris.Context) {
		accounts, err := listAccounts(db)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Writef("error viewing accounts")
		}
		ctx.JSON(accounts)
	})

	app.Get("/viewAccountsUI", func(ctx iris.Context) {
		ctx.View("viewAccounts.html")
	})

	app.Run(iris.Addr(":8080"))
}

func addAccountList(elem *models.CryptoAccount) string {
	return fmt.Sprintf("<td>%s</td><td>%s</td><td>%s</td>", elem.ID, elem.Name, elem.Email)
}

func createNewAccount(db *gorm.DB, email, name string) (newAccount *models.CryptoAccount, err error) {
	newUUID, _ := uuid.NewV4()
	// create random
	newAccount = &models.CryptoAccount{
		Status:           "ONBOARDING",
		TradingBlocked:   false,
		TransfersBlocked: false,
		AccountBlocked:   false,
		Name:             name,
		Email:            email,
	}
	newAccount.ID = newUUID.String()
	if err := db.Create(&newAccount).Error; err != nil {
		return &models.CryptoAccount{}, err
	}
	return newAccount, nil
}

func listAccounts(db *gorm.DB) (accs []*models.CryptoAccount, err error) {
	err = db.Find(&accs).Error
	return accs, err
}
