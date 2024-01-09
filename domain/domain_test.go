package domain_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/alvinatthariq/otpsvc/domain"
	"github.com/alvinatthariq/otpsvc/entity"

	"github.com/go-redis/redis"
	. "github.com/smartystreets/goconvey/convey"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

var (
	dbgorm      *gorm.DB
	redisClient *redis.Client
	err         error

	dom domain.DomainItf
)

func TestMain(t *testing.M) {
	dbgorm, err = gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3307)/otp_db?parseTime=true"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal(err)
		panic("Cannot connect to DB")
	}

	dbgorm.AutoMigrate(&entity.OTP{})

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
	})

	dom = domain.Init(
		dbgorm,
		redisClient,
	)

	exitVal := t.Run()

	os.Exit(exitVal)
}

func TestRequestOTP(t *testing.T) {
	Convey("TestRequestOTP", t, FailureHalts, func() {
		testCases := []struct {
			testID   int
			testType string
			testDesc string
			payload  entity.RequestOTP
			prepare  func()
		}{
			{
				testID:   1,
				testDesc: "Success request otp",
				testType: "P",
				payload: entity.RequestOTP{
					UserID: "integtest",
				},
				prepare: func() {
					// delete data before create
					otp := entity.OTP{}
					dbgorm.Where("user_id = ?", "integtest").Delete(&otp)
				},
			},
			{
				testID:   2,
				testDesc: "Fail request otp, user id empty",
				testType: "N",
				payload: entity.RequestOTP{
					UserID: "",
				},
				prepare: func() {
				},
			},
		}

		for _, tc := range testCases {
			tc.prepare()
			t.Logf("%d - [%s] : %s", tc.testID, tc.testType, tc.testDesc)
			_, err := dom.RequestOTP(context.Background(), tc.payload)
			if tc.testType == "P" {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldNotBeNil)
			}
		}
	})
}

func TestValidateOTP(t *testing.T) {
	Convey("TestValidateOTP", t, FailureHalts, func() {
		testCases := []struct {
			testID   int
			testType string
			testDesc string
			payload  entity.OTP
			prepare  func()
		}{
			{
				testID:   1,
				testDesc: "Success validate otp",
				testType: "P",
				payload: entity.OTP{
					UserID: "integtest",
					OTP:    "12345",
				},
				prepare: func() {
					// insert data testing
					otp := entity.OTP{
						UserID:    "integtest",
						OTP:       "12345",
						ExpiredAt: time.Now().Add(1 * time.Hour),
					}
					dbgorm.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&otp)
				},
			},
			{
				testID:   2,
				testDesc: "Fail validate otp, user id empty",
				testType: "N",
				payload: entity.OTP{
					UserID: "",
				},
				prepare: func() {
				},
			},
			{
				testID:   3,
				testDesc: "Fail validate otp, invalid otp",
				testType: "P",
				payload: entity.OTP{
					UserID: "integtest",
					OTP:    "54321",
				},
				prepare: func() {
					// insert data testing
					otp := entity.OTP{
						UserID:    "integtest",
						OTP:       "12345",
						ExpiredAt: time.Now().Add(1 * time.Hour),
					}
					dbgorm.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&otp)
				},
			},
			{
				testID:   4,
				testDesc: "Fail validate otp, otp expired",
				testType: "P",
				payload: entity.OTP{
					UserID: "integtest",
					OTP:    "12345",
				},
				prepare: func() {
					// insert data testing
					otp := entity.OTP{
						UserID:    "integtest",
						OTP:       "12345",
						ExpiredAt: time.Now().Add(-1 * time.Hour),
					}
					dbgorm.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&otp)
				},
			},
			{
				testID:   5,
				testDesc: "Fail validate otp, otp not found",
				testType: "P",
				payload: entity.OTP{
					UserID: "invalidd",
				},
				prepare: func() {
				},
			},
		}

		for _, tc := range testCases {
			tc.prepare()
			t.Logf("%d - [%s] : %s", tc.testID, tc.testType, tc.testDesc)
			_, err := dom.ValidateOTP(context.Background(), tc.payload)
			if tc.testType == "P" {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldNotBeNil)
			}
		}
	})
}
