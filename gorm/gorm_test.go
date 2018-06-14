package gorm

import (
	"os"
	"testing"

	keys "github.com/Brickchain/go-keys.v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type service struct {
	name string
	svc  func(*testing.T) keys.StoredKeyService
}

var services []*service

func TestMain(m *testing.M) {
	services = make([]*service, 0)

	// setup Gorm
	services = append(services, &service{
		name: "Gorm",
		svc: func(t *testing.T) keys.StoredKeyService {
			db, err := gorm.Open("sqlite3", ":memory:")
			if err != nil {
				t.Fatal(err)
			}

			svc, err := NewGormStoredKeyService(db)
			if err != nil {
				t.Fatal(err)
			}

			return svc
		},
	})

	os.Exit(m.Run())
}

func TestStoredKeyService_Get(t *testing.T) {
	type test struct {
		name    string
		svc     keys.StoredKeyService
		prepare func(*testing.T, *test)
		id      string
		wantErr bool
	}
	tests := []test{
		{
			name: "Get",
			prepare: func(t *testing.T, tt *test) {
				if err := tt.svc.Save(&keys.StoredKey{ID: "abc"}); err != nil {
					t.Fatal(err)
				}
			},
			id:      "abc",
			wantErr: false,
		},
		{
			name: "Get_not_exist",
			prepare: func(t *testing.T, tt *test) {
				if err := tt.svc.Save(&keys.StoredKey{ID: "abc"}); err != nil {
					t.Fatal(err)
				}
			},
			id:      "does_not_exist",
			wantErr: true,
		},
	}
	for _, svc := range services {
		t.Run(svc.name, func(t *testing.T) {
			for _, tt := range tests {
				tt.svc = svc.svc(t)
				t.Run(tt.name, func(t *testing.T) {
					if tt.prepare != nil {
						tt.prepare(t, &tt)
					}
					got, err := tt.svc.Get(tt.id)
					if (err != nil) != tt.wantErr {
						t.Errorf("StoredKeyService.Get() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if !tt.wantErr && got.ID != tt.id {
						t.Errorf("StoredKeyService.Get() = ID: %s, want ID: %s", got.ID, tt.id)
					}
				})
			}
		})
	}
}
