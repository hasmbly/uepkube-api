package models

import (
	"time"
)

type (

 Tbl_user struct {
	Id_user			int 		`gorm:"primary_key"`
	Nik				string
	No_kk			string
	Nama			string
	Jenis_kelamin	string
	Tempat_lahir	string
	Tgl_lahir		string
	Alamat			string
	Id_kelurahan	string
	Id_kecamatan	int
	Id_kabupaten	int
	Email			string
	No_hp			string
	Photo 			string
	Created_at		*time.Time
	Updated_at		*time.Time
	Created_by   	string
	Updated_by		string	
}

 Tbl_uep struct {
	Id_uep			int 		`gorm:"primary_key"`
	Id_pendamping	int
	Bantuan_modal	int
	Status			string
	Created_at		*time.Time
	Updated_at		*time.Time
	Created_by   	string
	Updated_by		string		
}

 Tbl_kube struct {
	Id_kube      	int 	`gorm:"primary_key"`
	Nama_kube    	string
	Jenis_usaha  	string
	Bantuan_modal	int
	Ketua        	int
	Sekertaris   	int
	Bendahara    	int
	Anggota1     	int
	Anggota2     	int
	Anggota3     	int
	Anggota4     	int
	Anggota5     	int
	Anggota6     	int
	Anggota7     	int
	Pendamping		int
	Photo     		string
	Status       	string
	Created_at   	*time.Time
	Updated_at   	*time.Time
	Created_by   	string
	Updated_by		string
}

 Ktype struct {
	Id_kube      	int 	`json:"id_kube"`
	Nama_kube    	string 	`json:"nama_kube"`
	Jenis_usaha  	string	`json:"jenis_usaha"`
	Bantuan_modal	int 	`json:"bantuan_modal"`
	Ketua 			string 	`json:"ketua"`
	Sekertaris 		string 	`json:"sekertaris"`
	Bendahara 		string 	`json:"bendahara"`
	Anggota1 		string 	`json:"anggota1"`
	Anggota2 		string 	`json:"anggota2"`
	Anggota3 		string 	`json:"anggota3"`
	Anggota4 		string 	`json:"anggota4"`
	Anggota5 		string 	`json:"anggota5"`
	Anggota6 		string 	`json:"anggota6"`
	Anggota7 		string 	`json:"anggota7"`
	Pendamping 		string 	`json:"pendamping"`
	Photo 		 	string 	`json:"photo"`
	Status 		 	string 	`json:"status"`
}

 U struct {
	Tbl_user
	CustU
}

 UepKube struct {
	Uep 	interface{} `json:"uep"`
	Kube 	interface{} `json:"kube"`
}

 CustU struct {
	Nama 			string `json:"nama_pendamping"`
	Bantuan_modal 	string `json:"bantuan_modal"`
	Status 			string `json:"status"`
}

 ResPagin struct {
	Content 		interface{}
	First 			bool 					`json:"first"`
	Last 			bool 					`json:"last"`
	Number 			int 					`json:"number"`
	NumberOfElement int64 					`json:"numberOfElements"`
 	Pageable 		Pageable				`json:"pageable"`
	Sort 			Sort 					`json:"sort"`
	TotalPages 		float64 				`json:"totalPages"`
	TotalElements 	int64 					`json:"totalElements"`
}

 PosPagin struct {
	Page 		int 			`json:"page" 		validate:"required"`
	Size 		int 			`json:"size" 		validate:"required"`
	SortField 	string 			`json:"sortField" 	validate:"required"`
	SortOrder 	string 			`json:"sortOrder"	validate:"required"`
	Filters 	[]struct {
					Key 		string `json:"key"`
					Operation 	string `json:"operation"`
					Value 		string `json:"value"`	
				}
}

Pageable struct{
			Offset 		int64 	`json:"offset"`
			PageNumber 	int 	`json:"pageNumeber"`
			PageSize 	int64 	`json:"pageSize"`
			Paged 		bool 	`json:"paged"`
			Unpaged 	bool 	`json:"unpaged"`
		}

Sort struct{
		Sorted 		bool 	`json:"sorted"`
		Unsorted 	bool 	`json:"unsorted"`
	}

Tbl_account struct{
		Id_user			int 		`gorm:"primary_key"`
		Id_group		int
		Username 		string
		Password 		string
		Created_at   	*time.Time
		Updated_at   	*time.Time
		Created_by   	string
		Updated_by		string	
	}

// AddUep struct{
// 	Id_user			int 		`gorm:"primary_key"`
// 	Nik				string
// 	No_kk			string
// 	Nama			string
// 	Jenis_kelamin	string
// 	Tempat_lahir	string
// 	Tgl_lahir		string
// 	Alamat			string
// 	Id_kelurahan	string
// 	Id_kecamatan	int
// 	Id_kabupaten	int
// 	Email			string
// 	No_hp			string
// 	Photo			string
// 	Bantuan_modal	int
// 	Status			string
	
// }		
		
)