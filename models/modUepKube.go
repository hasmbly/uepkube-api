package models

import (
	"time"
)

type (

 Tbl_user struct {
	Id_user			int 		`json:"id_user" gorm:"primary_key"`
	Nik				string 		`json:"nik"`
	No_kk			string 		`json:"no_kk"`
	Nama			string 		`json:"nama"`
	Jenis_kelamin	string 		`json:"jenis_kelamin" enums:"L, P"`
	Tempat_lahir	string 		`json:"tempat_lahir"`
	Tgl_lahir		string 		`json:"tgl_lahir"`
	Alamat			string 		`json:"alamat"`
	Id_kelurahan	*string 	`json:"id_kelurahan"`
	Id_kecamatan	*int 		`json:"id_kecamatan"`
	Id_kabupaten	*int 		`json:"id_kabupaten"`
	Email			string 		`json:"email"`
	No_hp			string 		`json:"no_hp"`
	Photo 			string 		`json:"photo"`
	Created_at		*time.Time 			`json:"-" gorm:"timestamp;null"`
	Updated_at		*time.Time 			`json:"-" gorm:"timestamp;null"`
	Created_by   	*string 	`json:"created_by"`
	Updated_by		*string	 	`json:"updated_by"` 	
}

 Tbl_uep struct {
	Id_uep			int 		`gorm:"primary_key"`
	Id_pendamping	int 		`json:"id_pendamping"`
	Bantuan_modal	int 		`json:"bantuan_modal"`
	Status			int 		`json:"status"`
}

 Tbl_kube struct {
	Id_kube      	int 	`json:"id_kube" gorm:"primary_key"`
	Nama_kube    	string 	`json:"nama_kube"`
	Jenis_usaha  	string 	`json:"jenis_usaha"`
	Bantuan_modal	int 	`json:"bantuan_modal"`
	Ketua        	int 	`json:"ketua" sql:"DEFAULT:NULL"`
	Sekertaris   	int 	`json:"sekertaris" sql:"DEFAULT:NULL"`
	Bendahara    	int 	`json:"bendahara" sql:"DEFAULT:NULL"`
	Anggota1     	int 	`json:"anggota1" sql:"DEFAULT:NULL"`
	Anggota2     	int 	`json:"anggota2" sql:"DEFAULT:NULL"`
	Anggota3     	int 	`json:"anggota3" sql:"DEFAULT:NULL"`
	Anggota4     	int 	`json:"anggota4" sql:"DEFAULT:NULL"`
	Anggota5     	int 	`json:"anggota5" sql:"DEFAULT:NULL"`
	Anggota6     	int 	`json:"anggota6" sql:"DEFAULT:NULL"`
	Anggota7     	int 	`json:"anggota7" sql:"DEFAULT:NULL"`
	Pendamping		int 	`json:"pendamping" sql:"DEFAULT:NULL"`
	Photo     		string 	`json:"photo"`
	Status       	int 	`json:"status" sql:"default:0"`
	Created_at   	*time.Time `json:"-" gorm:"timestamp;null"`
	Updated_at   	*time.Time `json:"-" gorm:"timestamp;null"`
	Created_by   	*string 	`json:"created_by"`
	Updated_by		*string 	`json:"updated_by"`
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
	Status 		 	int 	`json:"status"`
}

 U struct{
	Tbl_user
	CustU
}

 CustU struct{
 	Id_Pendamping 	int `json:"id_pendamping"`
	Nama 			string `json:"nama_pendamping"`
	Bantuan_modal 	int `json:"bantuan_modal"`
	Status 			int `json:"status"`
}

 UepKube struct{
	Uep 	interface{} `json:"uep"`
	Kube 	interface{} `json:"kube"`
}

 ResPagin struct{
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

 PosPagin struct{
	Page 		int 			`json:"page" 		validate:"required"`
	Size 		int 			`json:"size" 		validate:"required"`
	SortField 	string 			`json:"sortField" 	validate:"required"`
	SortOrder 	string 			`json:"sortOrder"	validate:"required"`
	Filters     []Filters
}

Filters struct{
				Key 		string `json:"key"`
				Operation 	string `json:"operation"`
				Value 		string `json:"value"`
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

Uep struct{
	*Tbl_user
	Id_pendamping	int 		`json:"id_pendamping"`
	Bantuan_modal	int 		`json:"bantuan_modal"`
	Status			int 		`json:"status"`
}

PaginateKubes struct{
	Id_kube      	int 	`json:"id_kube"`
	Nama_kube    	string 	`json:"nama_kube"`
	Jenis_usaha  	string	`json:"jenis_usaha"`
	Bantuan_modal	int 	`json:"bantuan_modal"`
	Photo 		 	string 	`json:"photo"`
	Status 		 	int 	`json:"status"`	
}

)