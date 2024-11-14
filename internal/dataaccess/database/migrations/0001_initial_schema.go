package migrations

import (
	"database/sql"
	"errors"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"google.golang.org/protobuf/proto"
)

func register0001Migration(app core.App) error {
	if _, err := app.Dao().FindCollectionByNameOrId("download_tasks"); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			app.Logger().With("err", err).Error("failed to find download_tasks collection")
			return err
		}
	} else {
		return nil
	}

	form := forms.NewCollectionUpsert(app, &models.Collection{})
	form.ListRule = proto.String("@request.auth.id = of_account_id")
	form.ViewRule = proto.String("@request.auth.id = of_account_id")
	form.UpdateRule = proto.String("@request.auth.id = of_account_id")
	form.DeleteRule = proto.String("@request.auth.id = of_account_id")
	form.Name = "download_tasks"
	form.Type = models.CollectionTypeBase
	form.Schema.AddField(&schema.SchemaField{
		Name:     "of_account_id",
		Type:     schema.FieldTypeRelation,
		Required: true,
		Options: &schema.RelationOptions{
			CollectionId:  "_pb_users_auth_",
			CascadeDelete: true,
		},
	})
	form.Schema.AddField(&schema.SchemaField{
		Name:     "download_type",
		Type:     schema.FieldTypeNumber,
		Required: true,
		Options: &schema.NumberOptions{
			Min:       proto.Float64(1),
			Max:       proto.Float64(1),
			NoDecimal: true,
		},
	})
	form.Schema.AddField(&schema.SchemaField{
		Name:     "url",
		Type:     schema.FieldTypeUrl,
		Required: true,
	})

	if err := form.Submit(); err != nil {
		app.Logger().With("err", err).Error("failed to create download_tasks collection")
		return err
	}

	return nil
}
