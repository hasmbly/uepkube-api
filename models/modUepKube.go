package models

import (
	"time"
)

type (

/**
 *
 * 
 * Tables Struct
 *
 * 
 */

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
	Lat 			*float64 	`json:"lat"`
	Lng 			*float64 	`json:"lng"`
	Ig 				string 		`json:"ig"`
	Fb 				string 		`json:"fb"`
	Wa 				string 		`json:"wa"`
	Created_at		*time.Time 	`json:"created_at" gorm:"timestamp;null"`
	Updated_at		*time.Time 	`json:"updated_at" gorm:"timestamp;null"`
	Created_by   	*string 	`json:"-"`
	Updated_by		*string	 	`json:"-"` 
	Region 			*View_address `json:"region" gorm:"foreignkey:id_kelurahan;association_foreignkey:id_kelurahan"`
	Kelurahan 	*Tbl_kelurahan `json:"kelurahan" gorm:"foreignkey:id_kelurahan;association_foreignkey:id_kelurahan"`
	Kecamatan 	*Tbl_kecamatan `json:"kecamatan" gorm:"foreignkey:id_kecamatan;association_foreignkey:id_kecamatan"`
	Kabupaten 	*Tbl_kabupaten `json:"kabupaten" gorm:"foreignkey:id_kabupaten;association_foreignkey:id_kabupaten"`
	Flag 			string 		`json:"flag" gorm:"-"`
	Is_eligible 	bool 		`json:"is_eligible" gorm:"-"`
}

Tbl_user_files struct{
	Id				int 		`json:"id" gorm:"primary_key"`
	Id_user			int 		`json:"id_user"`
	Files			string 		`json:"files"`	
	Description		string 		`json:"description"`	
	Type			string 		`json:"type" enums:"PDF, IMAGE"`
	Is_display		int 		`json:"is_display"`	
}

Tbl_kelurahan struct{
	Id_kelurahan 		string 		`json:"id_kelurahan" gorm:"primary_key"`
	Kelurahan 			string 		`json:"kelurahan"`
}

Tbl_kecamatan struct{
	Id_kecamatan 		string 		`json:"id_kecamatan" gorm:"primary_key"`
	Kecamatan 			string 		`json:"kecamatan"`
}

Tbl_kabupaten struct{
	Id_kabupaten 		string 		`json:"id_kabupaten" gorm:"primary_key"`
	Kabupaten 			string 		`json:"kabupaten"`
}

Tbl_user_photo struct {
	Id				int 		`json:"id" gorm:"primary_key"`
	Id_user			int 		`json:"id_user"`
	Photo			string 		`json:"photo"`
	Is_display		int 		`json:"is_display"`
}

Tbl_usaha_uep_photo struct {
	Id				int 		`json:"id" gorm:"primary_key"`
	Id_uep			int 		`json:"id_uep"`
	Photo			string 	 	`json:"photo"`	
	Is_display		int 		`json:"is_display"`	
}

Tbl_usaha_kube_photo struct {
	Id				int 		`json:"id" gorm:"primary_key"`
	Id_kube			int 		`json:"id_kube"`
	Photo			string 	 	`json:"photo"`	
	Is_display		int 		`json:"is_display"`	
}

Tbl_uep struct {
	Id_uep			int 		`json:"id_uep" gorm:"primary_key"`
	Id_pendamping	int 		`json:"-"`
	Nama_usaha		string 		`json:"nama_usaha"`
	Id_jenis_usaha	int 		`json:"-"`
	Status			int 		`json:"status"`
	Pendamping 		*Tbl_pendamping `json:"pendamping" gorm:"foreignkey:id_pendamping;association_foreignkey:id_pendamping"`
	JenisUsaha 		*Tbl_jenis_usaha `json:"jenis_usaha" gorm:"foreignkey:id_jenis_usaha;association_foreignkey:id_usaha"`
	PeriodsHistory 	[]*Tbl_periods_uepkube `json:"periods_history" gorm:"foreignkey:id_uep"`	
	Photo 			[]*Tbl_uepkube_files `json:"photo" gorm:"foreignkey:id_uep"`
}

 Tbl_kube struct {
	Id_kube      	int 	`json:"id_kube" gorm:"primary_key"`
	Nama_kube    	string 	`json:"nama_kube"`
	Nama_usaha  	string 	`json:"nama_usaha"`
	Id_jenis_usaha  int 	`json:"id_jenis_usaha"`
	Ketua        	int 	`json:"-" sql:"DEFAULT:NULL"`
	Sekertaris   	int 	`json:"-" sql:"DEFAULT:NULL"`
	Bendahara    	int 	`json:"-" sql:"DEFAULT:NULL"`
	Anggota1     	int 	`json:"-" sql:"DEFAULT:NULL"`
	Anggota2     	int 	`json:"-" sql:"DEFAULT:NULL"`
	Anggota3     	int 	`json:"-" sql:"DEFAULT:NULL"`
	Anggota4     	int 	`json:"-" sql:"DEFAULT:NULL"`
	Anggota5     	int 	`json:"-" sql:"DEFAULT:NULL"`
	Anggota6     	int 	`json:"-" sql:"DEFAULT:NULL"`
	Anggota7     	int 	`json:"-" sql:"DEFAULT:NULL"`
	Id_pendamping	int 	`json:"id_pendamping" sql:"DEFAULT:NULL"`
	Status       	int 	`json:"status" sql:"default:0"`
	Ig 				string 		`json:"ig"`
	Fb 				string 		`json:"fb"`
	Wa 				string 		`json:"wa"`	
	Created_at   	*time.Time `json:"-" gorm:"timestamp;null"`
	Updated_at   	*time.Time `json:"-" gorm:"timestamp;null"`
	Created_by   	*string 	`json:"created_by"`
	Updated_by		*string 	`json:"updated_by"`
	Pendamping 		*Tbl_pendamping `json:"pendamping" gorm:"foreignkey:id_pendamping;association_foreignkey:id_pendamping"`
	JenisUsaha 		*Tbl_jenis_usaha `json:"jenis_usaha" gorm:"foreignkey:id_jenis_usaha;association_foreignkey:id_usaha"`
	PeriodsHistory 	[]*Tbl_periods_uepkube `json:"periods_history" gorm:"foreignkey:id_kube"`	
	Photo 			[]*Tbl_uepkube_files `json:"photo" gorm:"foreignkey:id_kube"`
	Items 			[]Kubes_items `json:"items"`
	Region 			View_address `json:"region" gorm:"-"`
}

Kubes_items struct{
		*Tbl_user
		Posisi string `json:"posisi"`
	}

Tbl_kube_photo struct {
	Id				int 		`json:"id" gorm:"primary_key"`
	Id_kube			int 		`json:"id_kube"`
	Photo			string 		`json:"photo"`	
	Is_display		int 		`json:"is_display"`	
}

Tbl_pendamping struct {
	Id_pendamping		int 		`json:"id_pendamping" gorm:"primary_key"`
	Jenis_pendamping	string 		`json:"jenis_pendamping"`
	Periode				string 		`json:"periode"`
}

Tbl_account struct {
		Id_user			int 		`gorm:"primary_key"`
		Id_roles		int
		Username 		string
		Password 		string
		Created_at   	*time.Time 	`json:"-"`
		Updated_at   	*time.Time 	`json:"-"`
	}

Tbl_roles struct{
		Id				int 		`gorm:"primary_key"`
		Roles_name 		string
	}

Tbl_uepkube_photo struct{
	Id 				int 		`json:"id" gorm:"primary_key"`
	Id_uep 	 		int 		`json:"id_uep" sql:"DEFAULT:NULL"`
	Id_kube 	 	int 		`json:"id_kube" sql:"DEFAULT:NULL"`
	Photo			string 		`json:"photo"`	
	Is_display		int 		`json:"is_display"`	
}

Tbl_produk struct{
	Id_produk 		int 		`json:"id_produk" gorm:"primary_key"`
	Nama_produk 	string 		`json:"nama_produk"`
	Id_jp 			int 		`json:"id_jp" sql:"DEFAULT:NULL"`
	JenisProduk 	*Tbl_jenis_produk `json:"jenis_produk" gorm:"foreignkey:id_jp;association_foreignkey:id"`
	Deskripsi 		string 		`json:"deskripsi"`
	Created_at		*time.Time 	`json:"-" gorm:"timestamp;null"`
	Updated_at		*time.Time 	`json:"-" gorm:"timestamp;null"`
	Created_by   	*string 	`json:"-"`
	Updated_by		*string	 	`json:"-"` 	
}

Tbl_produk_photo struct {
	Id				int 		`json:"id" gorm:"primary_key"`
	Id_produk		int 		`json:"id_produk"`
	Photo			string 		`json:"photo"`	
	Is_display		int 		`json:"is_display"`	
}

Tbl_jenis_usaha struct{
	Id_usaha 		int 		`json:"id_usaha" gorm:"primary_key"`
	Jenis_usaha 	string 		`json:"jenis_usaha"`
	Created_at		*time.Time 	`json:"-" gorm:"timestamp;null"`
	Updated_at		*time.Time 	`json:"-" gorm:"timestamp;null"`
	Created_by   	*string 	`json:"-"`
	Updated_by		*string	 	`json:"-"` 	
}

Tbl_jenis_produk struct{
	Id 				int 		`json:"id" gorm:"primary_key"`
	Name 			string 		`json:"name"`
	Created_at		*time.Time 	`json:"-" gorm:"timestamp;null"`
	Updated_at		*time.Time 	`json:"-" gorm:"timestamp;null"`
}

Tbl_usaha_produk struct{
	Id 				int 		`json:"id" gorm:"primary_key"`
	Id_uep 	 		int 		`json:"id_uep" sql:"DEFAULT:NULL"`
	Id_kube 	 	int 		`json:"id_kube" sql:"DEFAULT:NULL"`
	Id_usaha 	 	int 		`json:"id_usaha"`
	Id_produk 	 	int 		`json:"id_produk"`
	Created_at		*time.Time 	`json:"-" gorm:"timestamp;null"`
	Updated_at		*time.Time 	`json:"-" gorm:"timestamp;null"`
}

Tbl_pelatihan struct{
	Id_pelatihan 		int 		`json:"id_pelatihan" gorm:"primary_key"`
	Judul_pelatihan 	string 		`json:"judul_pelatihan"`
	Instruktur 			string 		`json:"instruktur"`
	Lokasi_pelatihan 	string 		`json:"lokasi_pelatihan"`
	Deskripsi 			string 		`json:"deskripsi"`
	Start_date			string 		`json:"start_date"`
	End_date			string 		`json:"end_date"`
	Start_time			string 		`json:"start_time"`
	End_time			string 		`json:"end_time"`	
	Peruntukan			string 		`json:"peruntukan" enums:"UEP, KUBE, UEP-KUBE"`
	Quota				int 		`json:"quota"`
	Status				int 		`json:"status"`
	Created_at			*time.Time 	`json:"-" gorm:"timestamp;null"`
	Updated_at			*time.Time 	`json:"-" gorm:"timestamp;null"`
	Created_by   		*string 	`json:"created_by"`
	Updated_by			*string	 	`json:"updated_by"` 	
	Files 		 		[]Tbl_pelatihan_files 	`json:"files" gorm:"foreignkey:id_pelatihan"`
	Photo 		 		[]Tbl_pelatihan_files 	`json:"photo" gorm:"foreignkey:id_pelatihan"`
	Kehadiran 			[]PelatihanKehadiran	`json:"kehadiran"`
}

Tbl_pelatihan_files struct{
	Id				int 		`json:"id" gorm:"primary_key"`
	Id_pelatihan	int 		`json:"id_pelatihan"`
	Files			string 		`json:"files"`	
	Description		string 		`json:"description"`	
	Type			string 		`json:"type" enums:"PDF, IMAGE"`
	Is_display		int 		`json:"is_display"`	
}

Tbl_activity struct{
	Id 					int	 		`json:"id" gorm:"primary_key"`
	Id_pendamping 		int	 		`json:"id_pendamping"`
	Tanggal 			string 		`json:"tanggal"`
	Nama_kegiatan 		string 		`json:"nama_kegiatan"`
	Rincian 			string 		`json:"rincian"`
	Created_at			*time.Time 	`json:"created_at" gorm:"timestamp;null"`
	Updated_at			*time.Time 	`json:"updated_at" gorm:"timestamp;null"`	
	Created_by   		*string 	`json:"created_by"`
	Updated_by			*string	 	`json:"updated_by"` 	
	Photo 		 		[]Tbl_activity_files 	`json:"photo" gorm:"foreignkey:id_activity"`
}

Tbl_activity_files struct{
	Id				int 		`json:"id" gorm:"primary_key"`
	Id_activity		int 		`json:"id_activity"`
	Files			string 		`json:"files"`	
	Description		string 		`json:"description"`	
	Type			string 		`json:"type" enums:"PDF, IMAGE"`
	Is_display		int 		`json:"is_display"`	
}


Tbl_faq struct{
	Id_faq 		int 		`json:"id_faq" gorm:"primary_key"`
	Pertanyaan 	string 		`json:"pertanyaan"`
	Jawaban 	string 		`json:"jawaban"`
}

Tbl_periods_uepkube struct{
	Id 				int 		`json:"id" gorm:"primary_key"`
	Id_uep 	 		int 		`json:"id_uep" sql:"DEFAULT:NULL"`
	Id_kube 	 	int 		`json:"id_kube" sql:"DEFAULT:NULL"`
	Id_periods 	 	int 		`json:"id_periods"`
	Created_at		*time.Time 	`json:"-" gorm:"timestamp;null"`
	Updated_at		*time.Time 	`json:"-" gorm:"timestamp;null"`
	BantuanPeriods 	*Tbl_bantuan_periods `json:"bantuan_periods" gorm:"foreignkey:id;association_foreignkey:id_periods"`
}

Tbl_bantuan_periods struct{
	Id 				int 		`json:"id" gorm:"primary_key"`
	Bantuan_modal 	float32			`json:"bantuan_modal" sql:"type:decimal(10,2);"`
	Start_date 		*time.Time 	`json:"start_date"`
	End_date 		*time.Time 	`json:"end_date"`
	Peruntukan		string 		`json:"peruntukan" enums:"UEP, KUBE`
	Quota 			int 		`json:"quota"`
	Status 			int 		`json:"status"`
	Created_at		*time.Time 	`json:"-" gorm:"timestamp;null"`
	Updated_at		*time.Time 	`json:"-" gorm:"timestamp;null"`
	Usaha 			*Tbl_usaha_uepkube `json:"usaha" gorm:"foreignkey:id_periods;association_foreignkey:id"`
	CreditDebit 	[]*Tbl_credit_debit `json:"credit_debit" gorm:"foreignkey:id_periods;association_foreignkey:id"`
	MonevHistory 	[]*Tbl_monev_uepkube `json:"monev_history" gorm:"foreignkey:id_periods;association_foreignkey:id"`
}

Tbl_credit_debit struct{
	Id 				int 		`json:"id" gorm:"primary_key"`
	Id_uep 			int			`json:"id_uep" sql:"default:null"`
	Id_kube 		int			`json:"id_kube" sql:"default:null"`
	Credit 			int			`json:"credit"`
	Debit 			int			`json:"debit"`
	Nilai 			float32		`json:"nilai" sql:"type:decimal(10,2);"`
	Deskripsi 		string		`json:"deskripsi"`
	File 			string		`json:"file" sql:"default:null"`
	Id_periods 		int			`json:"id_periods"`
	Transaction_at	*time.Time 	`json:"transaction_at" gorm:"timestamp;null"`
	Updated_at		*time.Time 	`json:"updated_at" gorm:"timestamp;null"`
}

Tbl_lapkeu_uepkube struct{
	Id 				int 		`json:"id" gorm:"primary_key"`
	Id_uep 			int			`json:"id_uep" sql:"default:null"`
	Id_kube 		int			`json:"id_kube" sql:"default:null"`
	Bulan 			int			`json:"bulan"`
	Tahun 			string		`json:"tahun"`
	Modal 			float32		`json:"modal" sql:"type:decimal(10,2);"`
	Omset 			float32		`json:"omset" sql:"type:decimal(10,2);"`
	Pendapatan 		float32		`json:"pendapatan" sql:"type:decimal(10,2);"`
	Id_pendamping 	int			`json:"id_pendamping"`
	Pendamping 		*Tbl_pendamping `json:"pendamping" gorm:"foreignkey:id_pendamping;association_foreignkey:id_pendamping"`
	Created_at		*time.Time 	`json:"created_at" gorm:"timestamp;null"`
	Updated_at		*time.Time 	`json:"updated_at" gorm:"timestamp;null"`
	Created_by   	*string 	`json:"created_by"`
	Updated_by		*string	 	`json:"updated_by"` 
	Photo 		 	[]Tbl_lapkeu_files 	`json:"photo" gorm:"foreignkey:id_lapkeu"`
	Detail 			interface{} `json:"detail" gorm:"-"`
}

Tbl_lapkeu_files struct{
	Id				int 		`json:"id" gorm:"primary_key"`
	Id_lapkeu		int 		`json:"id_lapkeu"`
	Files			string 		`json:"files"`	
	Description		string 		`json:"description"`	
	Type			string 		`json:"type" enums:"PDF, IMAGE"`
	Is_display		int 		`json:"is_display"`	
}

// dimensi uepkube
Tbl_dimensi_uepkube struct{
	Id_dimensi 		int 		`json:"id_dimensi" gorm:"primary_key"`
	Nama_dimensi 	string		`json:"nama_dimensi"`
	Aspek_uep 		[]*Tbl_aspek_uep `json:"aspek_uep" gorm:"foreignkey:id_dimensi;association_foreignkey:id_dimensi"`
	Aspek_kube 		[]*Tbl_aspek_kube `json:"aspek_kube" gorm:"foreignkey:id_dimensi;association_foreignkey:id_dimensi"`
}

// aspek uep
Tbl_aspek_uep struct{
	Id_aspek 		int 		`json:"id_aspek" gorm:"primary_key"`
	Id_dimensi 		int 		`json:"id_dimensi"`
	Nama_aspek 		string		`json:"nama_aspek"`
	Kriteria_uep 	[]*Tbl_kriteria_uep `json:"kriteria_uep" gorm:"foreignkey:id_aspek;association_foreignkey:id_aspek"`
}

// aspek kube
Tbl_aspek_kube struct{
	Id_aspek 		int 		`json:"id_aspek" gorm:"primary_key"`
	Id_dimensi 		int 		`json:"id_dimensi"`
	Nama_aspek 		string		`json:"nama_aspek"`
	Kriteria_kube 	[]*Tbl_kriteria_kube `json:"kriteria_kube" gorm:"foreignkey:id_aspek;association_foreignkey:id_aspek"`
}

// kriteria uep
Tbl_kriteria_uep struct{
	Id_kriteria 	int 		`json:"id_kriteria" gorm:"primary_key"`
	Id_aspek 		int 		`json:"id_aspek"`
	Nama_kriteria 	string		`json:"nama_kriteria"`
	Bobot 			int			`json:"bobot"`
	Indikator_uep 	[]*Tbl_indikator_uep `json:"indikator_uep" gorm:"foreignkey:id_kriteria;association_foreignkey:id_kriteria"`
}

// kriteria kube
Tbl_kriteria_kube struct{
	Id_kriteria 	int 		`json:"id_kriteria" gorm:"primary_key"`
	Id_aspek 		int 		`json:"id_aspek"`
	Nama_kriteria 	string		`json:"nama_kriteria"`
	Bobot 			int			`json:"bobot"`
	Indikator_kube 	[]*Tbl_indikator_kube `json:"indikator_kube" gorm:"foreignkey:id_kriteria;association_foreignkey:id_kriteria"`
}

// indikator uep
Tbl_indikator_uep struct{
	Id_indikator 	int 		`json:"id_indikator" gorm:"primary_key"`
	Id_kriteria 	int			`json:"id_kriteria"`
	Nama_indikator 	string		`json:"nama_indikator"`
	Skor_indikator 	int			`json:"skor_indikator"`
}

// indikator kube
Tbl_indikator_kube struct{
	Id_indikator 	int 		`json:"id_indikator" gorm:"primary_key"`
	Id_kriteria 	int			`json:"id_kriteria"`
	Nama_indikator 	string		`json:"nama_indikator"`
	Skor_indikator 	int			`json:"skor_indikator"`
}

// kehadiran pelatihan
Tbl_kehadiran struct{
	Id 				int 		`json:"id" gorm:"primary_key"`
	Id_pelatihan 	int			`json:"id_pelatihan"`
	Id_user 		int			`json:"id_user"`
	Id_pendamping 	int			`json:"id_pendamping"`
	Flag 			string		`json:"flag"`
}

View_address struct{
	Id_kabupaten 	int 		`json:"id_kabupaten"`
	Id_kecamatan 	int 		`json:"id_kecamatan"`
	Id_kelurahan 	string 		`json:"id_kelurahan"`
	Region 			string 		`json:"region"`
}

Tbl_inventory struct{
	Id 					int	 		`json:"id" gorm:"primary_key"`
	Id_uep 				int			`json:"id_uep" sql:"default:null"`
	Id_kube 			int			`json:"id_kube" sql:"default:null"`
	Nama 				string 		`json:"nama"`
	Merk 				string 		`json:"merk"`
	Harga 				float32 	`json:"harga"`
	Id_pendamping 		int			`json:"id_pendamping"`
	Pendamping 		*Tbl_pendamping `json:"pendamping" gorm:"foreignkey:id_pendamping;association_foreignkey:id_pendamping"`
	Created_at			*time.Time 	`json:"created_at" gorm:"timestamp;null"`
	Updated_at			*time.Time 	`json:"updated_at" gorm:"timestamp;null"`	
	Created_by   		*string 	`json:"created_by"`
	Updated_by			*string	 	`json:"updated_by"`
	Photo 		 		[]Tbl_activity_files 	`json:"photo" gorm:"foreignkey:id_activity"`
	Detail 				interface{} `json:"detail" gorm:"-"`
}

Tbl_inventory_files struct{
	Id				int 		`json:"id" gorm:"primary_key"`
	Id_inventory	int 		`json:"id_inventory"`
	Files			string 		`json:"files"`	
	Description		string 		`json:"description"`	
	Type			string 		`json:"type" enums:"PDF, IMAGE"`
	Is_display		int 		`json:"is_display"`	
}

Tbl_monev_category struct{
	Id 				int	 		`json:"id" gorm:"primary_key"`
	Rank			string 		`json:"rank" enums:"A, B, C"`
	Description 	string 		`json:"description"`
}

Tbl_monev_uepkube struct{
	Id 					int	 		`json:"id" gorm:"primary_key"`
	Id_uep 				int			`json:"id_uep" sql:"default:null"`
	Id_kube 			int			`json:"id_kube" sql:"default:null"`
	Sum_total 			int			`json:"sum_total" sql:"default:null"`
	Id_category 		int			`json:"id_category" sql:"default:null"`
	Category 			*Tbl_monev_category `json:"category" gorm:"foreignkey:id_category;association_foreignkey:id"`
	Id_pendamping 		int			`json:"id_pendamping"`
	Pendamping 			*Tbl_pendamping `json:"pendamping" gorm:"foreignkey:id_pendamping;association_foreignkey:id_pendamping"`
	Is_monev			string 		`json:"is_monev" enums:"SUDAH, BELUM"`
	Id_periods 			int			`json:"id_periods"`
	Periods 			*Tbl_bantuan_periods `json:"periods" gorm:"foreignkey:id_periods;association_foreignkey:id"`
	Flag				string 		`json:"-" enums:"UEP, KUBE"`
	Data_monev			[]*Data_monev `json:"data_monev"`
	Created_at			*time.Time 	`json:"created_at" gorm:"timestamp;null"`
	Updated_at			*time.Time 	`json:"updated_at" gorm:"timestamp;null"`
	Detail 				interface{} `json:"detail" gorm:"-"`
}

Tbl_monev_result_uepkube struct{
	Id 					int	 		`json:"id" gorm:"primary_key"`
	Id_uep 				int			`json:"id_uep" sql:"default:null"`
	Id_kube 			int			`json:"id_kube" sql:"default:null"`	
	Id_indikator 		int			`json:"-"`
	Skor_total 			int			`json:"skor_total"`
	Created_at			*time.Time 	`json:"created_at" gorm:"timestamp;null"`
	Updated_at			*time.Time 	`json:"updated_at" gorm:"timestamp;null"`	
	Created_by   		*string 	`json:"-"`
	Updated_by			*string	 	`json:"-"`
}

Tbl_usaha_uepkube struct{
	Id 					int	 		`json:"id" gorm:"primary_key"`
	Id_uep 				int			`json:"id_uep" sql:"default:null"`
	Id_kube 			int			`json:"id_kube" sql:"default:null"`	
	Id_periods 			int			`json:"id_periods"`
	Nama_usaha 		 	string		`json:"nama_usaha"`
	Id_jenis_usaha 		int			`json:"id_jenis_usaha"`
	JenisUsaha 			*Tbl_jenis_usaha `json:"jenis_usaha" gorm:"foreignkey:id_jenis_usaha;association_foreignkey:id_usaha"`
	AllProduk 			[]*Tbl_produk_uepkube `json:"all_produk" gorm:"foreignkey:id_usaha_uke;association_foreignkey:id"`
	Created_at			*time.Time 	`json:"-" gorm:"timestamp;null"`
	Updated_at			*time.Time 	`json:"-" gorm:"timestamp;null"`	
}

Tbl_produk_uepkube struct{
	Id 					int	 		`json:"id" gorm:"primary_key"`
	Id_usaha_uke 		int			`json:"id_usaha_uke"`
	Id_produk 		 	int			`json:"id_produk"`
	DetailProduk 		*Tbl_produk `json:"detail_produk" gorm:"foreignkey:id_produk;association_foreignkey:id_produk"`
	Created_at			*time.Time 	`json:"-" gorm:"timestamp;null"`
	Updated_at			*time.Time 	`json:"-" gorm:"timestamp;null"`	
}

Tbl_uepkube_files struct {
	Id				int 		`json:"id" gorm:"primary_key"`
	Id_uep 			int			`json:"id_uep" sql:"default:null"`
	Id_kube 		int			`json:"id_kube" sql:"default:null"`		
	Files			string 		`json:"files"`	
	Description		string 		`json:"description"`	
	Type			string 		`json:"type" enums:"PDF, IMAGE"`
	Is_display		int 		`json:"is_display"`	
}

/**
 *
 * 
 * Custom Struct
 *
 * 
 */

Dummy struct{
	Uep *Uep
	Kube *Kube
	Pendamping *Pendamping
	Verifikator *Verifikator
	Produk *Produk
	Pelatihan *Pelatihan
	Kehadiran *Kehadiran
	Activity *Activity
	Inventory *Inventory
	Lapkeu *Lapkeu
	Monev *Monev
}

User struct{
	*Tbl_user
	Flag string `json:"flag'`
}

Data_monev struct{
	Nama_dimensi 	string 	`json:"dimensi"`
	Nama_Aspek 		string 	`json:"aspek"`
	Skor_indikator 	int 	`json:"skor"`
	Bobot 			int 	`json:"bobot"`
	Skor_total 		int 	`json:"skor_total"`
}

MemberPelatihan struct{
	Id_user int    `json:"id_user"`
	Nama 	string `json:"nama"`
	Flag 	string `json:"flag"`
}

PelatihanKehadiran struct{
	*Tbl_kehadiran
	Nama string `json:"nama"`
}

Uep struct{
	*Tbl_user
	Id_pendamping	int 		`json:"id_pendamping"`
	Id_periods		int 		`json:"id_periods"`
	Nama_usaha		string 		`json:"nama_usaha"`
	Id_jenis_usaha	int 		`json:"id_jenis_usaha"`
	Status			int 		`json:"status"`
}

Kube struct{
	*Tbl_kube
	Id_periods		int 		`json:"id_periods"`
}

Monev struct{
	Id_uep 			int    `json:"id_uep"`
	Id_kube 		int    `json:"id_kube"`
	Id_indikator 	[]int  `json:"id_indikator"`
	Id_pendamping	int 	`json:"id_pendamping"`
}

Pendamping struct{
	*Tbl_user
	Id_roles			int 		`json:"id_roles"`
	Roles_name			string 		`json:"roles_name"`
	Username			string 		`json:"username"`
	Password			string 		`json:"password"`
	Jenis_pendamping	string 		`json:"jenis_pendamping"`
	Periode 			string 		`json:"periode"`
	Photo 		 		[]Tbl_user_photo 	`json:"photo"`
}

Verifikator struct{
	*Tbl_user
	Id_roles			int 		`json:"id_roles"`
	Roles_name			string 		`json:"roles_name"`
	Username			string 		`json:"username"`
	Password			string 		`json:"password"`
	Photo 		 		[]Tbl_user_photo 	`json:"photo"`
}

Produk struct{
	*Tbl_usaha_produk
	Photo 		 		[]Tbl_produk_photo 	`json:"photo"`
}

Pelatihan struct{
	*Tbl_pelatihan
}

Activity struct{
	*Tbl_activity
	Photo 		 		[]Tbl_activity_files 	`json:"photo" gorm:"foreignkey:id_activity"`
}

Kehadiran struct{
	Id_pelatihan 	int `json:"id_pelatihan"`
	Id_pendamping 	int `json:"id_pendamping"`
	Kehadiran 		[]Tbl_kehadiran
}

Lapkeu struct{
	*Tbl_lapkeu_uepkube
}

Inventory struct{
	*Tbl_inventory
}

UepKube struct{
	Uep 	interface{} `json:"uep"`
	Kube 	interface{} `json:"kube"`
}

/**
 * Abstract Struct
 */

DetailUep struct{
	*Tbl_user
	Kelurahan *Tbl_kelurahan `gorm:"foreignkey:id_kelurahan;association_foreignkey:id_kelurahan"`
	// Nama_usaha 		*string `json:"nama_usaha"`
	// Jenis_usaha 	*Tbl_jenis_usaha `json:"jenis_usaha"`
	// Bantuan_periods *Tbl_bantuan_periods `json:"bantuan_periods"`
	// Pendamping 		*Tbl_pendamping `json:"pendamping"`	
	// Status 			*int `json:"status"`
}

Items struct {
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
}

ShowKube struct {
	Id_kube      	int 		`json:"id"`
	Nama_kube    	string 		`json:"nama"`
	Nama_usaha  	string		`json:"Nama_usaha"`
	Alamat			*string 	`json:"alamat"`
	Lat				*string 	`json:"lat"`
	Lng				*string 	`json:"lng"`
	Photo 		 	[]Tbl_uepkube_files 	`json:"photo"`
	Items 			Items 		`json:"items"`	
	Flag 			string 		 `json:"flag"`	
}

ShowUep struct{
	Id_user			int 		`json:"id"`
	Nik				string 		`json:"nik"`
	Jenis_usaha  	string		`json:"jenis_usaha"`
	Nama			string 		`json:"nama"`
	Alamat			string 		`json:"alamat"`
	Photo 			string 		`json:"photo"`
	Lat 			float64 	`json:"lat"`
	Lng 			float64 	`json:"lng"`
	Items 			[]interface{} `json:"items"`
	Flag 			string 		 `json:"flag"`
}

ShowProduks struct{
	Id 				int 		`json:"id"`
	Nama			string 		`json:"nama"`
	Alamat			string 		`json:"alamat"`
	No_hp			string 		`json:"no_hp"`
	Nama_produk 	string 		`json:"nama_produk"`
	Deskripsi 		string 		`json:"deskripsi"`
	Jenis_usaha 	string 		`json:"jenis_usaha"`
	Photo 			[]Tbl_produk_photo `json:"photo"`
}

 Ktype struct {
	Id_kube      	int 	`json:"id_kube"`
	Nama_kube    	string 	`json:"nama_kube"`
	Jenis_usaha  	string	`json:"jenis_usaha"`
	Bantuan_modal	int 	`json:"bantuan_modal"`
	Alamat			*string 	`json:"alamat"`
	Lat			*string 	`json:"lat"`
	Lng			*string 	`json:"lng"`
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
	Nama_pendamping string `json:"nama_pendamping"`
	Bantuan_modal 	int `json:"bantuan_modal"`
	Status 			int `json:"status"`
	Photo 		 	[]Tbl_user_photo 	`json:"photo"`
	Flag 			string 		 `json:"flag"`
}

UsahaUep struct {
	Id_usaha   	  int 	 `json:"id_usaha"`
	Nama_usaha    string `json:"nama_usaha"`
	Jenis_usaha   string `json:"jenis_usaha"`
	Photo 		 []Tbl_uepkube_files `json:"photo"`
}

UsahaKube struct {
	Id_usaha   	  int 	 `json:"id_usaha"`
	Nama_usaha    string `json:"nama_usaha"`
	Jenis_usaha   string `json:"jenis_usaha"`
	Photo 		 []Tbl_uepkube_files `json:"photo"`
}

CustomPendamping struct {
	*Tbl_pendamping
	Nama string `json:"nama_pendamping"`
}


/**
 * PaginateOptions Struct
 */


PosPagin struct{
	Page 		int 			`json:"page" example:"1"`
	Size 		int 			`json:"size" example:"10"`
	SortField 	string 			`json:"sortField" example:"created_at"`
	SortOrder 	string 			`json:"sortOrder" example:"desc"`
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

/**
 * PaginateResults Struct
 */

PaginateMonev struct{
	*Tbl_monev_uepkube
}

PaginateKubes struct{
	Id_kube      	int 	`json:"id"`
	Nama_kube    	string 	`json:"nama_kube"`
	Status 		 	int 	`json:"status"`
	Created_at 		*time.Time 	`json:"created_at"`
	Pendamping 		CustomPendamping 	`json:"pendamping"`
	Usaha 			UsahaKube 		`json:"usaha"`
	PeriodsHistory 	[]*Tbl_periods_uepkube `json:"periods_history"`
	Items 			[]Kubes_items `json:"items"`
}

PaginateUep struct{
	Id_uep      	int 	`json:"id"`
	Nama    		string 	`json:"nama"`
	Nik  			string	`json:"nik"`
	No_kk			string 	`json:"no_kk"`
	Alamat			string 	`json:"alamat"`
	Status 		 	int 	`json:"status"`
	Created_at 		string 	`json:"created_at"`
	Pendamping 		CustomPendamping	`json:"pendamping"`
	Usaha 			UsahaUep 		`json:"jenis_usaha"`
	PeriodsHistory 	[]*Tbl_periods_uepkube `json:"periods_history"`
}

PaginateProduks struct{
	Id 				int 		`json:"id"`
	Nama			string 		`json:"nama"`
	Alamat			string 		`json:"alamat"`
	No_hp			string 		`json:"no_hp"`
	Nama_produk 	string 		`json:"nama_produk"`
	Deskripsi 		string 		`json:"deskripsi"`
	Jenis_usaha 	string 		`json:"jenis_usaha"`
	Photo 			string 		`json:"photo"`
}

PaginatePendamping struct{
	Id_pendamping      	int 		`json:"id"`
	Nama				string 		`json:"nama"`
	Nik					string 		`json:"nik"`
	Username			string 		`json:"username"`
	Roles_name			string 		`json:"roles_name"`
	Jenis_pendamping	string 		`json:"jenis_pendamping"`
	Periode 			string 		`json:"periode"`
}

PaginateVerifikator struct{
	Id_user      		int 		`json:"id"`	
	Nama				string 		`json:"nama"`
	Nik					string 		`json:"nik"`
	Roles_name			string 		`json:"roles_name"`
	Username			string 		`json:"username"`
}

PaginatePelatihan struct{
	Id_pelatihan 		int 		`json:"id_pelatihan" gorm:"primary_key"`
	Judul_pelatihan 	string 		`json:"judul_pelatihan"`
	Instruktur 			string 		`json:"instruktur"`
	Lokasi_pelatihan 	string 		`json:"lokasi_pelatihan"`
	Start 				string 		`json:"start"`
	End 				string 		`json:"end"`
	Peruntukan 			string 		`json:"peruntukan"`
	Deskripsi 			string 		`json:"deskripsi"`
	Files 				[]Tbl_pelatihan_files 		`json:"photo"`
	Created_by   		*string 	`json:"created_by"`
}

PaginateAktivitas struct{
	*Tbl_activity
}

PaginateInventory struct{
	*Tbl_inventory
}

PaginateLapkeu struct{
	*Tbl_lapkeu_uepkube
}

 ResPagin struct{
	Content 		interface{}
	First 			bool 			`json:"first"`
	Last 			bool 			`json:"last"`
	Number 			int 			`json:"number"`
	NumberOfElement int64 			`json:"numberOfElements"`
 	Pageable 		Pageable		`json:"pageable"`
	Sort 			Sort 			`json:"sort"`
	TotalPages 		float64 		`json:"totalPages"`
	TotalElements 	int64 			`json:"totalElements"`
}


)