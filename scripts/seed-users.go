package main

import (
	"database/sql"
	"os"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"golang.org/x/crypto/bcrypt"
)

// SeedUsers create the first users for system.
func SeedUsers(db *sql.DB) {
	users := []model.User{}

	admin := model.User{
		ID:   ulid.Make().String(),
		Name: "Administrator",
		Role: "admin",
	}

	identityAdmin := model.Identity{
		ID:       ulid.Make().String(),
		UserID:   admin.ID,
		Provider: "email",
		UID:      "admin@corgi",
		Password: "password",
	}

	user := model.User{
		ID:   ulid.Make().String(),
		Name: "User",
		Role: "user",
	}

	identityUser := model.Identity{
		ID:       ulid.Make().String(),
		UserID:   user.ID,
		Provider: "email",
		UID:      "user@corgi",
		Password: "password",
	}

	admin.Identities = append(admin.Identities, identityAdmin)
	user.Identities = append(user.Identities, identityUser)
	users = append(users, user, admin)

	for _, u := range users {
		iden := u.Identities[0]

		// Check if provider and UID exists.
		var id string
		_ = db.QueryRow("SELECT id FROM identities WHERE provider = $1 AND uid = $2", iden.Provider, iden.UID).Scan(&id)
		if id != "" {
			continue
		}

		plainTextPassword := iden.Password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 8)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			os.Exit(2)
		}

		iden.Password = string(hashedPassword)

		tx, err := db.Begin()
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			os.Exit(2)
		}

		_, err = tx.Exec(`INSERT INTO users(id, name, role) VALUES($1, $2, $3)`,
			u.ID, u.Name, u.Role)

		if err != nil {
			log.Error().Caller().Msg(err.Error())

			err = tx.Rollback()
			if err != nil {
				log.Error().Caller().Msg(err.Error())
			}
			continue
		}

		_, err = tx.Exec(`INSERT INTO identities(id, user_id, provider, uid, password) 
		VALUES($1, $2, $3, $4, $5)`,
			iden.ID,
			iden.UserID,
			iden.Provider,
			iden.UID,
			iden.Password,
		)

		if err != nil {
			log.Error().Caller().Msg(err.Error())

			err = tx.Rollback()
			if err != nil {
				log.Error().Caller().Msg(err.Error())
			}
			continue
		}

		err = tx.Commit()
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			continue
		}
	}
}
