package smartCharging

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/dgraph-io/badger/v3"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/database"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	v16 "github.com/xBlaz3kx/ocppManager-go/configuration"
)

var (
	ErrProfileNotFound    = errors.New("profile not found")
	ErrCannotApplyProfile = errors.New("cannot apply profile")
)

type (
	Manager interface {
		AddProfile(profile *types.ChargingProfile) error
		GetProfile(profileId int) (*types.ChargingProfile, error)
		GetProfiles() []types.ChargingProfile
		GetCompositeSchedule() []ScheduleInterval
		RemoveProfile(profileId int) error
	}

	Impl struct {
		db                 *badger.DB
		numEvses           int
		maxCurrent         int
		evseManager        evse.Manager
		compositeSchedules map[int][]ScheduleInterval
	}
)

func NewManager(db *badger.DB, maxCurrent int) *Impl {
	return &Impl{
		db:         db,
		maxCurrent: maxCurrent,
	}
}

func (m *Impl) storeProfile(profile *types.ChargingProfile) error {
	return m.db.Update(func(txn *badger.Txn) error {
		marshal, err := json.Marshal(profile)
		if err != nil {
			return err
		}

		err = txn.Set(database.GetSmartChargingProfile(profile.ChargingProfileId), marshal)
		if err != nil {
			return err
		}

		return txn.Commit()
	})

}

func (m *Impl) AddProfile(profile *types.ChargingProfile) error {
	if profile == nil {
		return nil
	}

	if !m.canApplyProfile(profile) {
		return ErrCannotApplyProfile
	}

	profiles := m.GetProfiles()

	// Check if profile exists
	for i, chargingProfile := range profiles {
		// If profile with id already exists, replace it
		if chargingProfile.ChargingProfileId == profile.ChargingProfileId {
			profiles[i] = *profile
			goto Store
		}
	}

	profiles = append(profiles, *profile)

Store:

	maxProfile := getProfileWithHighestStack(getValidProfiles(getProfilesWithPurpose(types.ChargingProfilePurposeChargePointMaxProfile, profiles)))
	txDefaultProfile := getProfileWithHighestStack(getValidProfiles(getProfilesWithPurpose(types.ChargingProfilePurposeTxDefaultProfile, profiles)))
	txProfile := getProfileWithHighestStack(getValidProfiles(getProfilesWithPurpose(types.ChargingProfilePurposeTxProfile, profiles)))

	compositeSchedule := CreateCompositeSchedule([]*types.ChargingProfile{txProfile, txDefaultProfile, maxProfile})
	m.validateCompositeSchedule(compositeSchedule)

	return m.storeProfile(profile)
}

func (m *Impl) canApplyProfile(profile *types.ChargingProfile) bool {
	// Get max stack level
	stackLevelStr, err := ocppConfigManager.GetConfigurationValue(v16.ChargeProfileMaxStackLevel.String())
	if err != nil {
	}

	stackLevel, _ := strconv.Atoi(*stackLevelStr)

	// Get max profiles
	maxProfilesStr, err := ocppConfigManager.GetConfigurationValue(v16.MaxChargingProfilesInstalled.String())
	if err != nil {
	}

	maxProfiles, _ := strconv.Atoi(*maxProfilesStr)

	// Get max charging schedule periods
	maxPeriodsStr, err := ocppConfigManager.GetConfigurationValue(v16.ChargingScheduleMaxPeriods.String())
	if err != nil {
	}

	maxPeriods, _ := strconv.Atoi(*maxPeriodsStr)

	// Check if the stack level is valid
	if profile.StackLevel > stackLevel {
		return false
	}

	// Check if MaxProfiles is reached
	if len(m.GetProfiles())+1 >= maxProfiles {
		return false
	}

	// Check if the number of periods is valid
	if profile.ChargingSchedule != nil && len(profile.ChargingSchedule.ChargingSchedulePeriod) > maxPeriods {
		return false
	}

	return true
}

func (m *Impl) RemoveProfile(profileId int) error {
	return m.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(database.GetSmartChargingProfile(profileId))
		if err != nil {
			return err
		}

		return txn.Commit()
	})
}

func (m *Impl) GetProfile(profileId int) (*types.ChargingProfile, error) {
	var profile types.ChargingProfile

	err := m.db.View(func(txn *badger.Txn) error {
		profileItem, err := txn.Get(database.GetSmartChargingProfile(profileId))
		if err != nil {
			return err
		}

		err = profileItem.Value(func(v []byte) error {
			return json.Unmarshal(v, &profile)
		})

		return err
	})
	if err != nil {
		return nil, ErrProfileNotFound
	}

	return &profile, nil
}

func (m *Impl) GetProfiles() []types.ChargingProfile {
	var profiles []types.ChargingProfile

	// Query the database for EVSE settings.
	err := m.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte("profile-")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			var data types.ChargingProfile
			item := it.Item()

			// Value should be the EVSE struct.
			err := item.Value(func(v []byte) error {
				return json.Unmarshal(v, &data)
			})
			if err != nil {
				continue
			}
		}
		return txn.Commit()
	})
	if err != nil {
		log.WithError(err).Error("Error querying for smart charging profiles")
	}

	return profiles
}

func (m *Impl) GetCompositeSchedule() []ScheduleInterval {
	return nil
}

func (m *Impl) validateCompositeSchedule(newSchedule []ScheduleInterval) {

}
