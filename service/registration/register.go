package registration

type RegisterInfo struct {
	DepartmentID string `json:"departmentID"`
	OutpatientTypeID string `json:"outpatientTypeID"`
	PatientID string `json:"patientID"`
}

/*
Register create a registration record for the patient
regInfo is the registration information
repository is the repository to access the database
*/
func Register(regInfo *RegisterInfo,repository Repository)(string,error){
	return "registerID",nil
}