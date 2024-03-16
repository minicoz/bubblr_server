package convert

import (
	"bubblr/datastore"
	"bubblr/models"
)

func ConvertDBUserToUser(dbUser datastore.User) *models.User {

	return &models.User{
		UserID:         dbUser.UserID.String,
		FirstName:      dbUser.FirstName.String,
		LastName:       dbUser.LastName.String,
		Email:          dbUser.Email,
		HashedPassword: dbUser.HashedPassword,
		SchoolID:       int(dbUser.SchoolID.Int16),
		DobDay:         int(dbUser.DobDay.Int16),
		DobMonth:       int(dbUser.DobMonth.Int16),
		DobYear:        int(dbUser.DobYear.Int16),
		IsMale:         dbUser.IsMale.Bool,
		GradYear:       int(dbUser.GradYear.Int16),
		Verified:       dbUser.Verified.Bool,
		About:          dbUser.About.String,
		CreatedAt:      dbUser.CreatedAt,
	}
}
