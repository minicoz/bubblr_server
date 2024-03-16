package datastore

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx"
)

func TestDatastore_GetSchools(t *testing.T) {
	tests := []struct {
		name    string
		want    []School
		wantErr bool
	}{
		{
			name: "get schools",
			want: []School{{
				ID:     1,
				School: "Montessori",
				Tier:   1,
			},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			datastore, cleanupFunc := SetUpTestDB(t, "TestDatastore_GetSchools")
			defer cleanupFunc()
			d := Datastore{
				db: sqlx.NewDb(datastore.db, "postgres"),
			}
			got, err := d.GetSchools()
			if (err != nil) != tt.wantErr {
				t.Errorf("Datastore.GetSchools() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opts := cmp.Options{
				cmpopts.EquateApproxTime(2 * time.Second),
			}

			if diff := cmp.Diff(tt.want, got, opts); diff != "" {
				t.Errorf("DataStore.Create() mismatch (-want, +got) = %s", diff)
			}
		})
	}
}

func TestDatastore_GetSchoolById(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		id int16
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *School
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Datastore{
				db: tt.fields.db,
			}
			got, err := d.GetSchoolById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Datastore.GetSchoolById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Datastore.GetSchoolById() = %v, want %v", got, tt.want)
			}
		})
	}
}
