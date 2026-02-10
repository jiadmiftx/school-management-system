package schemas

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTeacherProfile_JSONSerialization(t *testing.T) {
	nip := "198501012010011001"
	nuptk := "1234567890123456"
	eduLevel := "S2"
	eduMajor := "Pendidikan Matematika"
	joinDate := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)

	profile := TeacherProfile{
		Id:               uuid.New(),
		UserId:           uuid.New(),
		UnitId:           uuid.New(),
		NIP:              &nip,
		NUPTK:            &nuptk,
		EducationLevel:   &eduLevel,
		EducationMajor:   &eduMajor,
		EmploymentStatus: "pns",
		JoinDate:         &joinDate,
	}

	// Serialize
	data, err := json.Marshal(profile)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Deserialize
	var result TeacherProfile
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Equal(t, profile.Id, result.Id)
	assert.Equal(t, *profile.NIP, *result.NIP)
	assert.Equal(t, *profile.NUPTK, *result.NUPTK)
	assert.Equal(t, profile.EmploymentStatus, result.EmploymentStatus)
}

func TestTeacherProfile_JSONFields(t *testing.T) {
	profile := TeacherProfile{
		Id:               uuid.New(),
		EmploymentStatus: "honorer",
	}

	data, err := json.Marshal(profile)
	assert.NoError(t, err)

	var mapResult map[string]interface{}
	err = json.Unmarshal(data, &mapResult)
	assert.NoError(t, err)

	// Verify JSON field names (using PascalCase as per schema)
	assert.Contains(t, mapResult, "Id")
	assert.Contains(t, mapResult, "UserId")
	assert.Contains(t, mapResult, "UnitId")
	assert.Contains(t, mapResult, "EmploymentStatus")
}

func TestTeacherProfile_NullableFields(t *testing.T) {
	profile := TeacherProfile{
		Id:               uuid.New(),
		UserId:           uuid.New(),
		UnitId:           uuid.New(),
		EmploymentStatus: "honorer",
		// All nullable fields are nil
	}

	data, err := json.Marshal(profile)
	assert.NoError(t, err)

	var result TeacherProfile
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Nil(t, result.NIP)
	assert.Nil(t, result.NUPTK)
	assert.Nil(t, result.EducationLevel)
	assert.Nil(t, result.JoinDate)
}

func TestStudentProfile_JSONSerialization(t *testing.T) {
	nis := "2024001"
	nisn := "0012345678"
	birthPlace := "Jakarta"
	birthDate := time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)
	gender := "L"
	fatherName := "Ahmad Budi"
	motherName := "Siti Aminah"

	profile := StudentProfile{
		Id:         uuid.New(),
		UserId:     uuid.New(),
		UnitId:     uuid.New(),
		NIS:        &nis,
		NISN:       &nisn,
		BirthPlace: &birthPlace,
		BirthDate:  &birthDate,
		Gender:     &gender,
		FatherName: &fatherName,
		MotherName: &motherName,
	}

	// Serialize
	data, err := json.Marshal(profile)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Deserialize
	var result StudentProfile
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Equal(t, profile.Id, result.Id)
	assert.Equal(t, *profile.NIS, *result.NIS)
	assert.Equal(t, *profile.NISN, *result.NISN)
	assert.Equal(t, *profile.BirthPlace, *result.BirthPlace)
}

func TestStudentProfile_ParentInfo(t *testing.T) {
	fatherName := "Ahmad"
	motherName := "Siti"
	guardianName := "Pak Budi"
	parentPhone := "081234567890"

	profile := StudentProfile{
		Id:           uuid.New(),
		UserId:       uuid.New(),
		UnitId:       uuid.New(),
		FatherName:   &fatherName,
		MotherName:   &motherName,
		GuardianName: &guardianName,
		ParentPhone:  &parentPhone,
	}

	data, err := json.Marshal(profile)
	assert.NoError(t, err)

	var mapResult map[string]interface{}
	err = json.Unmarshal(data, &mapResult)
	assert.NoError(t, err)

	assert.Contains(t, mapResult, "FatherName")
	assert.Contains(t, mapResult, "MotherName")
	assert.Contains(t, mapResult, "GuardianName")
	assert.Contains(t, mapResult, "ParentPhone")
}

func TestClass_JSONSerialization(t *testing.T) {
	teacherId := uuid.New()

	class := Class{
		Id:                uuid.New(),
		UnitId:            uuid.New(),
		Name:              "X IPA 1",
		Level:             10,
		AcademicYear:      "2025/2026",
		HomeroomTeacherId: &teacherId,
		Capacity:          35,
		IsActive:          true,
	}

	// Serialize
	data, err := json.Marshal(class)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Deserialize
	var result Class
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Equal(t, class.Id, result.Id)
	assert.Equal(t, class.Name, result.Name)
	assert.Equal(t, class.Level, result.Level)
	assert.Equal(t, class.AcademicYear, result.AcademicYear)
	assert.Equal(t, class.Capacity, result.Capacity)
	assert.Equal(t, class.IsActive, result.IsActive)
}

func TestClass_Defaults(t *testing.T) {
	class := Class{}

	// Test that zero values are as expected
	assert.Equal(t, 0, class.Level)
	assert.Equal(t, 0, class.Capacity)
	assert.False(t, class.IsActive) // Zero value is false
}

func TestClassEnrollment_JSONSerialization(t *testing.T) {
	enrolledAt := time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC)
	notes := "Transferred from Class A"

	enrollment := ClassEnrollment{
		Id:               uuid.New(),
		StudentProfileId: uuid.New(),
		ClassId:          uuid.New(),
		AcademicYear:     "2025/2026",
		Status:           EnrollmentStatusActive,
		EnrolledAt:       enrolledAt,
		Notes:            &notes,
	}

	// Serialize
	data, err := json.Marshal(enrollment)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Deserialize
	var result ClassEnrollment
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Equal(t, enrollment.Id, result.Id)
	assert.Equal(t, enrollment.AcademicYear, result.AcademicYear)
	assert.Equal(t, enrollment.Status, result.Status)
}

func TestClassEnrollment_StatusValues(t *testing.T) {
	tests := []struct {
		status ClassEnrollmentStatus
		value  string
	}{
		{EnrollmentStatusActive, "active"},
		{EnrollmentStatusGraduated, "graduated"},
		{EnrollmentStatusTransferred, "transferred"},
		{EnrollmentStatusDropped, "dropped"},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			enrollment := ClassEnrollment{
				Id:     uuid.New(),
				Status: tt.status,
			}

			data, err := json.Marshal(enrollment)
			assert.NoError(t, err)

			var mapResult map[string]interface{}
			err = json.Unmarshal(data, &mapResult)
			assert.NoError(t, err)

			assert.Equal(t, tt.value, mapResult["Status"])
		})
	}
}

func TestClassEnrollment_Defaults(t *testing.T) {
	enrollment := ClassEnrollment{}

	// Test that zero values are as expected
	assert.Equal(t, ClassEnrollmentStatus(""), enrollment.Status)
	assert.True(t, enrollment.EnrolledAt.IsZero())
	assert.Nil(t, enrollment.LeftAt)
	assert.Nil(t, enrollment.Notes)
}

func TestClassEnrollment_LeftAtNullable(t *testing.T) {
	leftAt := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)

	enrollment := ClassEnrollment{
		Id:         uuid.New(),
		Status:     EnrollmentStatusGraduated,
		EnrolledAt: time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC),
		LeftAt:     &leftAt,
	}

	data, err := json.Marshal(enrollment)
	assert.NoError(t, err)

	var result ClassEnrollment
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.NotNil(t, result.LeftAt)
	assert.Equal(t, leftAt.Year(), result.LeftAt.Year())
}
