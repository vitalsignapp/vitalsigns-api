package patient

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Patient struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	DateOfAdmit    string `json:"dateOfAdmit"`
	DateOfBirth    string `json:"dateOfBirth"`
	Diagnosis      string `json:"diagnosis"`
	IsRead         bool   `json:"isRead"`
	IsShowNotify   bool   `json:"isShowNotify"`
	Name           string `json:"name"`
	Sex            string `json:"sex"`
	Surname        string `json:"surname"`
	PatientRoomKey string `json:"patientRoomKey"`
}

type PatientPartial struct {
	HospitalKey string `json:"hospitalKey"`
}

type PatientData struct {
	Username       string `json:"username"`
	DateOfAdmit    string `json:"dateOfAdmit"`
	DateOfBirth    string `json:"dateOfBirth"`
	Diagnosis      string `json:"diagnosis"`
	IsRead         bool   `json:"isRead"`
	IsShowNotify   bool   `json:"isShowNotify"`
	Name           string `json:"name"`
	Sex            string `json:"sex"`
	Surname        string `json:"surname"`
	PatientRoomKey string `json:"patientRoomKey"`
	HospitalKey    string `json:"hospitalKey"`
}
type Symptom struct {
	Status bool   `json:"status"`
	Sym    string `json:"sym"`
}

type PatientLog struct {
	ID             string    `json:"id"`
	BloodPressure  string    `json:"bloodPressure"`
	HeartRate      string    `json:"heartRate"`
	HospitalKey    string    `json:"hospitalKey"`
	InputDate      string    `json:"inputDate"`
	InputRound     int       `json:"inputRound"`
	Microtime      int64     `json:"microtime"`
	OtherSymptoms  string    `json:"otherSymptoms"`
	Oxygen         string    `json:"oxygen"`
	PatientKey     string    `json:"patientKey"`
	PatientRoomKey string    `json:"patientRoomKey"`
	SymptomsCheck  []Symptom `json:"symptomsCheck"`
	Temperature    string    `json:"temperature"`
}

// NewRepoByRoomKey get patient data in room by room key
func NewRepoByRoomKey(fs *firestore.Client) func(context.Context, string) []Patient {
	return func(ctx context.Context, patientRoomKey string) []Patient {
		iter := fs.Collection("patientData").Where("patientRoomKey", "==", patientRoomKey).Documents(ctx)
		defer iter.Stop()

		pats := []Patient{}
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				continue
			}

			p := PatientData{}
			err = doc.DataTo(&p)
			if err != nil {
				continue
			}

			pats = append(pats, toPatient(p, doc.Ref.ID))
		}

		return pats
	}
}

// NewRepoByID get patient data by patient id
func NewRepoByID(fs *firestore.Client) func(context.Context, string) *Patient {
	return func(ctx context.Context, ID string) *Patient {
		doc, err := fs.Collection("patientData").Doc(ID).Get(ctx)
		if err != nil {
			return nil
		}
		p := PatientData{}
		err = doc.DataTo(&p)
		if err != nil {
			return nil
		}
		pat := toPatient(p, doc.Ref.ID)
		return &pat
	}
}

// NewRepoByHospital retreive all patient data by hospital id
func NewRepoByHospital(fs *firestore.Client) func(context.Context, string) []Patient {
	return func(ctx context.Context, hospitalID string) []Patient {
		iter := fs.Collection("patientData").Where("hospitalKey", "==", hospitalID).Documents(ctx)
		defer iter.Stop()

		pats := []Patient{}
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				continue
			}

			p := PatientData{}
			err = doc.DataTo(&p)
			if err != nil {
				continue
			}

			pats = append(pats, toPatient(p, doc.Ref.ID))
		}

		return pats
	}
}

// UpdateRepo update patient data by patient id
func UpdateRepo(fs *firestore.Client) func(context.Context, string, PatientRequest) error {
	return func(ctx context.Context, patientID string, p PatientRequest) error {
		_, err := fs.Collection("patientData").Doc(patientID).Set(ctx,
			map[string]interface{}{
				"dateOfAdmit":    p.DateOfAdmit,
				"dateOfBirth":    p.DateOfBirth,
				"diagnosis":      p.Diagnosis,
				"hospitalKey":    p.HospitalKey,
				"isRead":         p.IsRead,
				"isShowNotify":   p.IsShowNotify,
				"name":           p.Name,
				"patientRoomKey": p.PatientRoomKey,
				"sex":            p.Sex,
				"surname":        p.Surname,
				"username":       p.Username,
			}, firestore.MergeAll)
		if err != nil {
			return err
		}
		return nil
	}
}

// NewRepoLogByID query patient log by patient id
func NewRepoLogByID(fs *firestore.Client) func(context.Context, string) []PatientLog {
	return func(ctx context.Context, patientID string) []PatientLog {
		iter := fs.Collection("patientLog").
			Where("patientKey", "==", patientID).
			Documents(ctx)
		defer iter.Stop()

		pats := []PatientLog{}
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				continue
			}

			p := PatientLog{}
			p.ID = doc.Ref.ID
			err = doc.DataTo(&p)
			if err != nil {
				continue
			}

			pats = append(pats, p)
		}

		return pats
	}
}

// NewUpdateStatus NewUpdateStatus
func NewUpdateStatus(fs *firestore.Client) func(context.Context, string, string, PatientStatusRequest) error {
	return func(ctx context.Context, hospitalID string, patientID string, p PatientStatusRequest) error {
		dsnap, err := fs.Collection("patientData").Doc(patientID).Get(ctx)
		if err != nil {
			return err
		}
		var patient PatientPartial
		dsnap.DataTo(&patient)
		if hospitalID != patient.HospitalKey {
			return errors.New("Patient does not in this hospital")
		}

		var data = map[string]interface{}{}

		if p.IsRead != nil {
			data["isRead"] = p.IsRead
		}

		if p.IsNotify != nil {
			data["isShowNotify"] = p.IsNotify
		}

		_, err = fs.Collection("patientData").Doc(patientID).Set(ctx, data, firestore.MergeAll)
		if err != nil {
			return err
		}
		return nil
	}
}

func AddNewRepository(fs *firestore.Client) func(context.Context, PatientRequest) error {
	return func(ctx context.Context, p PatientRequest) error {
		_, _, err := fs.Collection("patientData").Add(ctx, map[string]interface{}{
			"dateOfAdmit":    p.DateOfAdmit,
			"dateOfBirth":    p.DateOfBirth,
			"diagnosis":      p.Diagnosis,
			"hospitalKey":    p.HospitalKey,
			"isRead":         p.IsRead,
			"isShowNotify":   p.IsShowNotify,
			"name":           p.Name,
			"patientRoomKey": p.PatientRoomKey,
			"sex":            p.Sex,
			"surname":        p.Surname,
			"username":       p.Username,
		})
		if err != nil {
			return err
		}
		return nil
	}
}

// NewRepoDeleteByID NewRepoDeleteByID
func NewRepoDeleteByID(fs *firestore.Client) func(context.Context, string) error {
	return func(ctx context.Context, patientID string) error {
		_, err := fs.Collection("patientData").Doc(patientID).Delete(ctx)
		if err != nil {
			return err
		}
		return nil
	}
}

// NewRepoDeleteLogByID NewRepoDeleteLogByID
func NewRepoDeleteLogByID(fs *firestore.Client) func(context.Context, string) error {
	return func(ctx context.Context, patientLogID string) error {
		_, err := fs.Collection("patientLog").Doc(patientLogID).Delete(ctx)
		if err != nil {
			return err
		}
		return nil
	}
}

func toPatient(p PatientData, ID string) Patient {
	return Patient{
		ID:             ID,
		Username:       p.Username,
		DateOfAdmit:    p.DateOfAdmit,
		DateOfBirth:    p.DateOfBirth,
		Diagnosis:      p.Diagnosis,
		IsRead:         p.IsRead,
		IsShowNotify:   p.IsShowNotify,
		Name:           p.Name,
		Sex:            p.Sex,
		Surname:        p.Surname,
		PatientRoomKey: p.PatientRoomKey,
	}

}
