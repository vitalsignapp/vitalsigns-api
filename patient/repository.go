package patient

import (
	"context"

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

func UpdateRepo(fs *firestore.Client) func(context.Context, string, PatientRequest) error {
	return func(ctx context.Context, patientID string, pt PatientRequest) error {
		_, err := fs.Collection("patientData").Doc(patientID).Set(ctx, pt)
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
