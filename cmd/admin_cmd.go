package cmd

import (
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/supabase/auth/internal/conf"
	"github.com/supabase/auth/internal/models"
	"github.com/supabase/auth/internal/storage"
)

var autoconfirm, isAdmin bool
var audience string

func getAudience(c *conf.GlobalConfiguration) string {
	if audience == "" {
		return c.JWT.Aud
	}

	return audience
}

func adminCmd() *cobra.Command {
	var adminCmd = &cobra.Command{
		Use: "admin",
	}

	adminCmd.AddCommand(&adminCreateUserCmd, &adminDeleteUserCmd)
	adminCmd.PersistentFlags().StringVarP(&audience, "aud", "a", "", "Set the new user's audience")

	adminCreateUserCmd.Flags().BoolVar(&autoconfirm, "confirm", false, "Automatically confirm user without sending an email")
	adminCreateUserCmd.Flags().BoolVar(&isAdmin, "admin", false, "Create user with admin privileges")

	return adminCmd
}

var adminCreateUserCmd = cobra.Command{
	Use: "createuser",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			logrus.Fatal("Not enough arguments to createuser command. Expected at least email and password values")
			return
		}

		execWithConfigAndArgs(cmd, adminCreateUser, args)
	},
}

var adminDeleteUserCmd = cobra.Command{
	Use: "deleteuser",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logrus.Fatal("Not enough arguments to deleteuser command. Expected at least ID or email")
			return
		}

		execWithConfigAndArgs(cmd, adminDeleteUser, args)
	},
}

func adminCreateUser(config *conf.GlobalConfiguration, args []string) {
	db, err := storage.Dial(config)
	if err != nil {
		logrus.Fatalf("Error opening database: %+v", err)
	}
	defer db.Close()

	aud := getAudience(config)
	var user *models.User
	is_new_user := true
	if existing_user, err := models.IsDuplicatedEmail(db, args[0], aud, nil); existing_user != nil {
		user = existing_user
		if err := user.SetPassword(db.Context(), args[1], config.Security.DBEncryption.Encrypt, config.Security.DBEncryption.EncryptionKeyID, config.Security.DBEncryption.EncryptionKey); err != nil {
			logrus.Fatalf("Error setting password: %+v", err)
		}
		is_new_user = false
	} else if err != nil {
		logrus.Fatalf("Error checking user email: %+v", err)
	} else {
		user, err = models.NewUser("", args[0], args[1], aud, nil)
		if err != nil {
			logrus.Fatalf("Error creating new user: %+v", err)
		}
	}

	user.IsNonDefaultPassword = true

	err = db.Transaction(func(tx *storage.Connection) error {
		var terr error
		if is_new_user {
			if terr = tx.Create(user); terr != nil {
				return terr
			}
		} else {
			if terr = tx.Update(user); terr != nil {
				return terr
			}
		}

		if len(args) > 2 {
			if terr = user.SetRole(tx, args[2]); terr != nil {
				return terr
			}
		} else if isAdmin {
			if terr = user.SetRole(tx, config.JWT.AdminGroupName); terr != nil {
				return terr
			}
		}

		if is_new_user && (config.Mailer.Autoconfirm || autoconfirm) {
			if terr = user.Confirm(tx); terr != nil {
				return terr
			}
		}
		return nil
	})
	if err != nil {
		logrus.Fatalf("Unable to upsert user (%s): %+v", args[0], err)
	}

	logrus.Infof("Upsert user: %s", args[0])
}

func adminDeleteUser(config *conf.GlobalConfiguration, args []string) {
	db, err := storage.Dial(config)
	if err != nil {
		logrus.Fatalf("Error opening database: %+v", err)
	}
	defer db.Close()

	user, err := models.FindUserByEmailAndAudience(db, args[0], getAudience(config))
	if err != nil {
		userID := uuid.Must(uuid.FromString(args[0]))
		user, err = models.FindUserByID(db, userID)
		if err != nil {
			logrus.Fatalf("Error finding user (%s): %+v", userID, err)
		}
	}

	if err = db.Destroy(user); err != nil {
		logrus.Fatalf("Error removing user (%s): %+v", args[0], err)
	}

	logrus.Infof("Removed user: %s", args[0])
}
