package main

import "log"

var wtbDefault = []string{
	"Buku KIA / catatan riwayat kesehatan anak (jika ada).",
	"Hasil pemeriksaan, rujukan dokter, atau laporan sekolah sebelumnya (jika ada).",
	"Mainan atau benda favorit anak agar anak lebih nyaman selama sesi.",
	"Camilan atau minum anak bila sewaktu-waktu diperlukan.",
}

// seedIfEmpty populates the database with the site's current content the
// first time it runs, so the client starts from what's already live instead
// of an empty dashboard.
func seedIfEmpty() {
	var count int
	db.QueryRow(`SELECT COUNT(*) FROM services`).Scan(&count)
	if count == 0 {
		seedServices()
	}
	db.QueryRow(`SELECT COUNT(*) FROM jobs`).Scan(&count)
	if count == 0 {
		seedJobs()
	}
	db.QueryRow(`SELECT COUNT(*) FROM articles`).Scan(&count)
	if count == 0 {
		seedArticles()
	}
	db.QueryRow(`SELECT COUNT(*) FROM facilities`).Scan(&count)
	if count == 0 {
		seedFacilities()
	}
	db.QueryRow(`SELECT COUNT(*) FROM gallery_photos`).Scan(&count)
	if count == 0 {
		seedGallery()
	}
	db.QueryRow(`SELECT COUNT(*) FROM team_members`).Scan(&count)
	if count == 0 {
		seedTeam()
	}
	db.QueryRow(`SELECT COUNT(*) FROM branches`).Scan(&count)
	if count == 0 {
		seedBranches()
	}
	db.QueryRow(`SELECT COUNT(*) FROM testimonials`).Scan(&count)
	if count == 0 {
		seedTestimonials()
	}
	db.QueryRow(`SELECT COUNT(*) FROM site_settings`).Scan(&count)
	if count == 0 {
		seedSettings()
	}
	log.Println("seed check complete")
}

func seedServices() {
	items := []Service{
		{
			Slug: "asesmen-tumbuh-kembang", Title: "Konsultasi & Asesmen Tumbuh Kembang", Icon: "clipboard",
			Short:    "Langkah awal setiap anak — Psikolog Anak memetakan kondisi anak sebelum menyusun program terapi.",
			Detail:   []string{"Konsultasi awal bersama Psikolog Anak.", "Asesmen menyeluruh untuk memetakan kemampuan dan hambatan anak.", "Penyusunan rencana terapi yang disesuaikan dengan kebutuhan anak."},
			ForWho:   "Anak usia 0–16 tahun yang orang tuanya ingin memastikan tahapan tumbuh kembangnya, baik untuk perkembangan umum maupun yang diduga membutuhkan penanganan khusus.",
			Signs:    []string{"Orang tua merasa perkembangan anak berbeda dari teman sebayanya.", "Anak belum mencapai milestone (bicara, motorik, sosial) sesuai usianya.", "Ingin mendapat gambaran menyeluruh sebelum menentukan jenis terapi yang tepat."},
			Goal:     "Memetakan kondisi dan kebutuhan anak secara objektif, sehingga program terapi yang disusun benar-benar sesuai — bukan pendekatan satu ukuran untuk semua anak.",
			Process:  []string{"Wawancara awal bersama orang tua mengenai riwayat dan kekhawatiran terhadap anak.", "Observasi dan asesmen langsung terhadap kemampuan anak oleh Psikolog Anak.", "Diskusi hasil asesmen dan rekomendasi program bersama orang tua."},
			Duration: "60 menit per sesi (intervensi maupun edukasi orang tua).", Professionals: "Psikolog Anak.", WhatToBring: wtbDefault,
			ExtraFAQ: []FAQItem{{Q: "Apa itu asesmen tumbuh kembang?", A: "Asesmen tumbuh kembang adalah proses evaluasi menyeluruh terhadap kemampuan anak di berbagai aspek — motorik, bahasa, kognitif, sosial-emosional, dan kemandirian — untuk mengetahui apakah perkembangan anak sesuai tahapan usianya dan menentukan jenis dukungan yang paling tepat jika dibutuhkan."}},
		},
		{
			Slug: "terapi-wicara", Title: "Terapi Wicara", Icon: "speech",
			Short:    "Untuk anak dengan keterlambatan bicara, gangguan bahasa, artikulasi, atau kelancaran berbicara.",
			Detail:   []string{"Menangani keterlambatan perkembangan bicara dan bahasa.", "Melatih artikulasi, kelancaran, dan irama bicara.", "Mendukung anak agar mampu berinteraksi secara wajar dengan lingkungannya."},
			ForWho:   "Anak dengan keterlambatan bicara, gangguan bahasa reseptif/ekspresif, gangguan artikulasi, kelancaran (gagap), maupun kesulitan komunikasi sosial.",
			Signs:    []string{"Belum mengucapkan kata bermakna di usia yang seharusnya sudah bisa.", "Ucapan sulit dipahami dibanding anak seusianya.", "Kesulitan memahami atau mengikuti instruksi sederhana.", "Gagap atau tersendat-sendat saat berbicara."},
			Goal:     "Mengembangkan kemampuan bicara, bahasa, dan komunikasi anak agar dapat berinteraksi secara wajar sesuai dengan tahapan usianya.",
			Process:  []string{"Pemanasan dan pendekatan agar anak nyaman dengan terapis.", "Latihan artikulasi, kosakata, atau kelancaran bicara melalui permainan terarah.", "Pencatatan perkembangan dan umpan balik singkat untuk orang tua di akhir sesi."},
			Duration: "60 menit per sesi (intervensi maupun edukasi orang tua).", Professionals: "Terapis Wicara bersertifikat.", WhatToBring: wtbDefault,
		},
		{
			Slug: "terapi-okupasi", Title: "Terapi Okupasi", Icon: "hands",
			Short:    "Melatih kemandirian dan kemampuan motorik halus anak seperti memegang pensil, menulis, dan menggambar.",
			Detail:   []string{"Melatih activity daily living skill (ADL) / kemandirian sehari-hari.", "Mengembangkan motorik halus prasekolah.", "Mendukung kesiapan anak mengikuti aktivitas belajar."},
			ForWho:   "Anak dengan hambatan motorik halus, kemandirian sehari-hari (ADL), integrasi sensorik, atau kesiapan mengikuti aktivitas belajar di sekolah.",
			Signs:    []string{"Kesulitan memegang pensil, menggunting, atau menggambar sesuai usianya.", "Belum mandiri makan, memakai baju, atau aktivitas sehari-hari sesuai usianya.", "Tampak sangat sensitif atau kurang responsif terhadap sentuhan/tekstur tertentu."},
			Goal:     "Melatih motorik halus, kemandirian aktivitas sehari-hari (ADL), dan integrasi sensorik agar anak lebih siap mengikuti aktivitas belajar dan bermain.",
			Process:  []string{"Aktivitas pemanasan sensorik sesuai kondisi anak.", "Latihan motorik halus dan ADL melalui permainan dan alat terapi.", "Evaluasi singkat serta saran latihan lanjutan di rumah untuk orang tua."},
			Duration: "60 menit per sesi (intervensi maupun edukasi orang tua).", Professionals: "Terapis Okupasi bersertifikat.", WhatToBring: wtbDefault,
		},
		{
			Slug: "fisioterapi", Title: "Fisioterapi", Icon: "run",
			Short:    "Melatih dan memperbaiki fungsi motorik kasar anak agar dapat beraktivitas sesuai usianya.",
			Detail:   []string{"Menangani gangguan perkembangan motorik kasar pada bayi dan anak.", "Melatih kekuatan otot dan koordinasi gerak.", "Mendukung perkembangan fisik sesuai tahapan usia anak."},
			ForWho:   "Bayi dan anak dengan keterlambatan motorik kasar, gangguan tonus otot, koordinasi gerak, atau kelainan postur tubuh.",
			Signs:    []string{"Terlambat tengkurap, duduk, merangkak, atau berjalan dibanding usianya.", "Otot terasa lemah, kaku, atau tidak seimbang antara sisi kiri dan kanan.", "Postur tubuh atau pola jalan tampak tidak biasa."},
			Goal:     "Melatih kekuatan otot, keseimbangan, dan koordinasi gerak agar anak dapat beraktivitas dan bermain sesuai tahapan usianya.",
			Process:  []string{"Pemeriksaan tonus otot, keseimbangan, dan pola gerak anak.", "Latihan fisik terarah dengan pendekatan bermain agar anak tetap nyaman.", "Evaluasi progres dan panduan latihan pendukung di rumah."},
			Duration: "60 menit per sesi (intervensi maupun edukasi orang tua).", Professionals: "Fisioterapis anak.", WhatToBring: wtbDefault,
		},
		{
			Slug: "terapi-perilaku", Title: "Terapi Perilaku", Icon: "puzzle",
			Short:    "Membantu anak berperilaku lebih adaptif terhadap lingkungan sekitarnya.",
			Detail:   []string{"Meningkatkan kemampuan anak untuk berperilaku adaptif.", "Mengurangi perilaku yang tidak sesuai di rumah maupun lingkungan sosial.", "Program disusun secara terukur berdasarkan hasil asesmen."},
			ForWho:   "Anak dengan perilaku yang menghambat aktivitas sehari-hari, interaksi sosial, atau proses belajarnya.",
			Signs:    []string{"Tantrum berlebihan dan sulit ditenangkan dibanding anak seusianya.", "Kesulitan mengikuti aturan atau instruksi sederhana.", "Perilaku repetitif atau impulsif yang mengganggu aktivitas sehari-hari."},
			Goal:     "Meningkatkan perilaku adaptif anak dan mengurangi perilaku yang tidak sesuai, melalui program yang terukur berdasarkan hasil asesmen.",
			Process:  []string{"Identifikasi pemicu dan pola perilaku anak bersama orang tua.", "Latihan perilaku terarah menggunakan pendekatan yang terukur dan konsisten.", "Pemantauan progres serta strategi pendampingan di rumah untuk orang tua."},
			Duration: "60 menit per sesi (intervensi maupun edukasi orang tua).", Professionals: "Terapis Perilaku.", WhatToBring: wtbDefault,
		},
		{
			Slug: "ortopedagogik", Title: "Ortopedagogik (Pendidikan Khusus)", Icon: "book",
			Short:    "Pendampingan belajar bagi anak yang kesulitan mengikuti pendidikan pada umumnya.",
			Detail:   []string{"Mendampingi anak dengan kesulitan belajar di sekolah.", "Mengembangkan potensi akademik sesuai bakat dan kemampuan anak.", "Menumbuhkan sikap belajar yang positif."},
			ForWho:   "Anak usia sekolah yang kesulitan mengikuti kurikulum pendidikan umum dan membutuhkan pendampingan belajar yang disesuaikan.",
			Signs:    []string{"Kesulitan mengikuti pelajaran dibanding teman sekelas.", "Membutuhkan pendampingan belajar khusus di luar jam sekolah.", "Mendapat rekomendasi evaluasi tambahan dari pihak sekolah."},
			Goal:     "Mendampingi anak mengembangkan potensi akademik sesuai bakat dan kemampuannya, serta menumbuhkan sikap belajar yang positif.",
			Process:  []string{"Identifikasi gaya belajar dan hambatan akademik anak.", "Pendampingan belajar dengan metode dan materi yang disesuaikan.", "Evaluasi berkala dan koordinasi dengan orang tua maupun pihak sekolah bila diperlukan."},
			Duration: "60 menit per sesi (intervensi maupun edukasi orang tua).", Professionals: "Guru Ortopedagogik / Pendidikan Khusus.", WhatToBring: wtbDefault,
		},
		{
			Slug: "psikologi-anak", Title: "Psikologi Anak", Icon: "heart",
			Short:    "Konsultasi psikologi, tes IQ, dan tes kesiapan sekolah untuk mendukung tumbuh kembang anak.",
			Detail:   []string{"Konsultasi psikologi anak dan remaja.", "Tes IQ dan tes minat bakat.", "Tes kesiapan sekolah sebelum anak memasuki jenjang pendidikan formal."},
			ForWho:   "Anak dan remaja yang membutuhkan pendampingan psikologis, maupun orang tua yang ingin melakukan tes IQ, minat bakat, atau kesiapan sekolah.",
			Signs:    []string{"Kecemasan berlebih atau kesulitan mengelola emosi untuk usianya.", "Sulit beradaptasi di lingkungan atau situasi baru.", "Orang tua ingin memahami potensi, minat bakat, atau kesiapan sekolah anak secara objektif."},
			Goal:     "Membantu anak dan orang tua memahami kondisi psikologis, potensi, serta kesiapan anak melalui konsultasi dan alat tes yang sesuai standar psikologi anak.",
			Process:  []string{"Wawancara awal bersama orang tua mengenai kekhawatiran atau tujuan konsultasi.", "Sesi konsultasi atau pengerjaan tes bersama psikolog anak.", "Pembahasan hasil dan rekomendasi tindak lanjut bersama orang tua."},
			Duration: "60 menit per sesi (intervensi maupun edukasi orang tua).", Professionals: "Psikolog Anak.", WhatToBring: wtbDefault,
			ExtraFAQ: []FAQItem{{Q: "Apa itu tes psikologi bagi anak?", A: "Tes psikologi anak adalah rangkaian alat ukur terstandar (misalnya tes IQ, minat bakat, atau kesiapan sekolah) yang digunakan psikolog untuk memahami kemampuan kognitif, emosional, dan kesiapan anak — sebagai dasar memberikan rekomendasi yang sesuai dengan kebutuhan anak."}},
		},
		{
			Slug: "stimulasi-anak", Title: "Stimulasi Anak 0–16 Tahun", Icon: "growth",
			Short:    "Program stimulasi terarah untuk mendukung fokus, kontrol motorik, dan koordinasi sensorik anak.",
			Detail:   []string{"Program stimulasi dini untuk anak usia 0–16 tahun.", "Mendukung fokus, kontrol motorik, dan koordinasi sensorik.", "Cocok untuk anak dengan perkembangan umum maupun kebutuhan khusus (ABK)."},
			ForWho:   "Anak usia 0–16 tahun, baik dengan perkembangan umum yang ingin mendapat stimulasi tambahan, maupun anak berkebutuhan khusus (ABK) yang memerlukan latihan terarah.",
			Signs:    []string{"Orang tua ingin memberikan stimulasi tambahan sejak dini untuk mendukung tumbuh kembang optimal.", "Anak mudah teralihkan fokusnya atau membutuhkan latihan koordinasi motorik.", "Anak ABK yang membutuhkan latihan sensorik dan motorik secara rutin."},
			Goal:     "Mendukung perkembangan fokus, kontrol motorik, dan koordinasi sensorik anak melalui program stimulasi yang terarah dan konsisten.",
			Process:  []string{"Penilaian singkat kebutuhan stimulasi sesuai usia dan kondisi anak.", "Aktivitas stimulasi terarah melalui permainan sensorik dan motorik.", "Pencatatan perkembangan dan rekomendasi stimulasi lanjutan di rumah."},
			Duration: "60 menit per sesi (intervensi maupun edukasi orang tua).", Professionals: "Terapis/Instruktur Stimulasi Anak, didampingi Psikolog Anak.", WhatToBring: wtbDefault,
		},
	}
	for i, s := range items {
		s.SortOrder = i
		saveService(s)
	}
	log.Println("seeded services")
}

func seedJobs() {
	items := []Job{
		{Slug: "terapis-wicara", Title: "Terapis Wicara", Branch: "Lamongan", Type: "Full-time", Status: "Dibuka",
			Description:  "Menangani sesi terapi wicara untuk anak dengan keterlambatan bicara, gangguan bahasa, artikulasi, maupun kelancaran berbicara, sesuai rencana terapi yang disusun bersama tim.",
			Requirements: []string{"Lulusan D3/D4/S1 Terapi Wicara.", "Memiliki STR/sertifikasi profesi yang masih berlaku.", "Berpengalaman menangani anak menjadi nilai tambah, fresh graduate dipersilakan melamar.", "Sabar, komunikatif, dan senang bekerja dengan anak-anak."}},
		{Slug: "terapis-okupasi", Title: "Terapis Okupasi", Branch: "Lamongan", Type: "Full-time", Status: "Dibuka",
			Description:  "Melaksanakan sesi terapi okupasi untuk melatih motorik halus, kemandirian aktivitas sehari-hari (ADL), dan integrasi sensorik anak.",
			Requirements: []string{"Lulusan D3/D4/S1 Okupasi Terapi.", "Memiliki STR/sertifikasi profesi yang masih berlaku.", "Memahami pendekatan sensory integration menjadi nilai tambah.", "Sabar, telaten, dan senang bekerja dengan anak-anak."}},
		{Slug: "fisioterapis-anak", Title: "Fisioterapis Anak", Branch: "Lamongan", Type: "Full-time", Status: "Dibuka",
			Description:  "Menangani sesi fisioterapi untuk melatih motorik kasar, kekuatan otot, dan koordinasi gerak anak sesuai rencana terapi.",
			Requirements: []string{"Lulusan D3/D4/S1 Fisioterapi.", "Memiliki STR/sertifikasi profesi yang masih berlaku.", "Pengalaman menangani pasien anak menjadi nilai tambah.", "Mampu bekerja dalam tim lintas disiplin."}},
		{Slug: "psikolog-anak", Title: "Psikolog Anak", Branch: "Lamongan", Type: "Part-time", Status: "Dibuka",
			Description:  "Memberikan layanan konsultasi psikologi, tes IQ, tes minat bakat, dan tes kesiapan sekolah bagi anak dan remaja.",
			Requirements: []string{"Lulusan S2 Profesi Psikologi (Psikolog).", "Memiliki SIPP (Surat Izin Praktik Psikolog) yang masih berlaku.", "Berpengalaman menangani asesmen psikologi anak.", "Mampu berkomunikasi dengan baik bersama orang tua."}},
		{Slug: "guru-ortopedagogik", Title: "Guru Ortopedagogik / Pendidikan Khusus", Branch: "Lamongan", Type: "Part-time", Status: "Dibuka",
			Description:  "Mendampingi anak dengan kesulitan belajar melalui program ortopedagogik yang disesuaikan dengan kebutuhan masing-masing anak.",
			Requirements: []string{"Lulusan S1 Pendidikan Luar Biasa (PLB) / Pendidikan Khusus.", "Memahami metode pendampingan belajar anak berkebutuhan khusus.", "Sabar dan kreatif dalam menyusun materi belajar."}},
		{Slug: "admin-customer-service", Title: "Admin & Customer Service", Branch: "Lamongan", Type: "Full-time", Status: "Ditutup",
			Description:  "Menangani pendaftaran pasien, penjadwalan sesi terapi, serta komunikasi dengan orang tua melalui WhatsApp maupun tatap muka.",
			Requirements: []string{"Minimal lulusan SMA/SMK sederajat.", "Komunikatif, rapi, dan terbiasa menggunakan aplikasi WhatsApp/komputer dasar.", "Berpengalaman di bidang layanan pelanggan menjadi nilai tambah."}},
	}
	for i, j := range items {
		j.SortOrder = i
		saveJob(j)
	}
	log.Println("seeded jobs")
}

func seedArticles() {
	items := []Article{
		{Slug: "tanda-keterlambatan-bicara-anak", Title: "Tanda-tanda Keterlambatan Bicara pada Anak yang Perlu Diwaspadai", Category: "Bicara & Bahasa",
			Excerpt: "Setiap anak punya waktu berkembang yang berbeda, tapi ada beberapa tanda keterlambatan bicara yang sebaiknya tidak ditunda untuk dikonsultasikan.",
			Content: []string{
				"Banyak orang tua bertanya-tanya, “Apakah wajar anak saya belum banyak bicara di usia ini?” Wajar jika setiap anak punya kecepatan berkembang yang berbeda, namun ada beberapa tanda yang sebaiknya tidak diabaikan.",
				"Di usia 12 bulan, anak umumnya sudah mulai mengoceh dengan pola suara yang bervariasi. Di usia 18 bulan, sebagian besar anak sudah mengucapkan beberapa kata sederhana seperti “mama” atau “papa”. Jika di usia 2 tahun anak belum mengucapkan kombinasi dua kata, ini bisa menjadi salah satu tanda yang perlu diperhatikan.",
				"Tanda lain yang perlu diwaspadai antara lain: anak tampak kesulitan memahami instruksi sederhana, jarang melakukan kontak mata saat diajak bicara, atau lebih banyak menunjuk dan merengek dibanding mencoba mengucapkan kata.",
				"Deteksi dini sangat penting karena semakin cepat ditangani, semakin besar peluang anak mengejar ketertinggalan perkembangan bahasanya. Terapi wicara yang dilakukan sejak dini terbukti membantu anak berkomunikasi lebih efektif seiring bertambahnya usia.",
				"Jika Anda mengamati salah satu tanda di atas pada si kecil, jangan ragu untuk berkonsultasi dengan psikolog anak agar mendapatkan gambaran yang lebih jelas dan penanganan yang tepat.",
			}},
		{Slug: "pentingnya-deteksi-dini-tumbuh-kembang", Title: "Kenapa Deteksi Dini Tumbuh Kembang Anak Itu Penting?", Category: "Deteksi Dini",
			Excerpt: "Semakin dini kondisi anak terdeteksi, semakin besar peluang penanganan berjalan optimal. Kenali kenapa deteksi dini jadi langkah krusial.",
			Content: []string{
				"Masa 0–5 tahun sering disebut sebagai golden age karena pada periode ini otak anak berkembang sangat pesat. Inilah alasan mengapa deteksi dini terhadap potensi hambatan tumbuh kembang menjadi sangat penting dilakukan sesegera mungkin.",
				"Deteksi dini bukan berarti mencari-cari masalah, melainkan memastikan anak berkembang sesuai tahapannya dan segera mendapat stimulasi tambahan jika diperlukan. Banyak kondisi yang jika ditangani sejak dini, hasilnya jauh lebih optimal dibanding ditangani saat anak sudah lebih besar.",
				"Beberapa aspek yang perlu dipantau orang tua meliputi perkembangan motorik (kasar dan halus), bahasa dan komunikasi, kemampuan sosial-emosional, serta kemandirian sehari-hari sesuai usia anak.",
				"Pemeriksaan tumbuh kembang secara berkala, baik melalui posyandu, dokter anak, maupun pusat layanan tumbuh kembang, membantu orang tua mendapatkan gambaran objektif — bukan sekadar membandingkan dengan anak lain.",
				"Jika ada kekhawatiran mengenai perkembangan si kecil, konsultasi dan asesmen menyeluruh adalah langkah awal yang tepat untuk memastikan kebutuhan anak terpetakan dengan baik sebelum menyusun rencana stimulasi atau terapi.",
			}},
		{Slug: "stimulasi-motorik-halus-di-rumah", Title: "5 Cara Sederhana Menstimulasi Motorik Halus Anak di Rumah", Category: "Stimulasi Anak",
			Excerpt: "Melatih motorik halus tidak harus dengan alat khusus. Berikut aktivitas sederhana di rumah yang bisa mendukung kemampuan menulis dan menggambar anak.",
			Content: []string{
				"Motorik halus adalah kemampuan menggunakan otot-otot kecil di tangan dan jari, yang penting untuk aktivitas seperti memegang pensil, mengancingkan baju, atau menggunakan sendok. Kemampuan ini bisa dilatih lewat aktivitas sederhana di rumah.",
				"1. Bermain plastisin atau adonan mainan — meremas dan membentuk adonan melatih kekuatan genggaman dan koordinasi jari anak.",
				"2. Meronce manik-manik — memasukkan tali ke lubang manik-manik melatih koordinasi mata-tangan sekaligus kesabaran anak.",
				"3. Menggunting kertas — dengan pengawasan orang tua, aktivitas ini melatih kontrol otot jari dan koordinasi dua tangan.",
				"4. Menuang air dari satu wadah ke wadah lain — aktivitas ini melatih kontrol gerakan dan konsentrasi.",
				"5. Mewarnai dan menggambar bebas — selain melatih motorik halus, aktivitas ini juga mendukung kreativitas anak.",
				"Jika anak tampak kesulitan signifikan dibanding teman sebayanya dalam aktivitas-aktivitas di atas, terapi okupasi dapat membantu melatih motorik halus secara lebih terarah dan terukur.",
			}},
		{Slug: "kapan-konsultasi-psikolog-anak", Title: "Kapan Sebaiknya Orang Tua Membawa Anak ke Psikolog Anak?", Category: "Psikologi Anak",
			Excerpt: "Konsultasi ke psikolog anak bukan hanya untuk kondisi berat. Kenali situasi-situasi yang sebaiknya mendapat pendampingan profesional.",
			Content: []string{
				"Banyak orang tua menunda konsultasi ke psikolog anak karena menganggapnya hanya diperlukan untuk kondisi yang “berat”. Padahal, psikolog anak juga dapat membantu berbagai situasi yang lebih umum dalam keseharian.",
				"Beberapa situasi yang bisa dipertimbangkan untuk konsultasi antara lain: anak sulit mengendalikan emosi secara berlebihan untuk usianya, kesulitan beradaptasi di lingkungan baru, menunjukkan kecemasan berlebih, atau mengalami perubahan perilaku setelah kejadian tertentu.",
				"Konsultasi psikologi juga bermanfaat untuk kebutuhan seperti tes kesiapan sekolah, tes IQ, maupun tes minat bakat — membantu orang tua memahami potensi dan kebutuhan belajar anak secara lebih objektif.",
				"Sesi konsultasi awal biasanya berupa wawancara dengan orang tua dan observasi terhadap anak, untuk memahami konteks permasalahan sebelum menentukan langkah selanjutnya.",
				"Membawa anak ke psikolog bukan tanda kegagalan orang tua — justru sebaliknya, ini adalah bentuk perhatian untuk memastikan anak tumbuh dengan kesehatan mental yang baik.",
			}},
		{Slug: "peran-fisioterapi-tumbuh-kembang-anak", Title: "Mengenal Peran Fisioterapi dalam Tumbuh Kembang Anak", Category: "Fisioterapi",
			Excerpt: "Fisioterapi anak tidak hanya untuk pemulihan cedera, tapi juga berperan penting dalam mendukung perkembangan motorik kasar anak.",
			Content: []string{
				"Ketika mendengar kata “fisioterapi”, banyak orang membayangkan penanganan untuk orang dewasa yang cedera. Padahal, fisioterapi juga memiliki peran penting dalam tumbuh kembang anak, khususnya pada aspek motorik kasar.",
				"Fisioterapi anak membantu menangani berbagai kondisi seperti keterlambatan duduk, merangkak, atau berjalan, gangguan keseimbangan dan koordinasi, hingga kelainan postur tubuh sejak usia dini.",
				"Melalui latihan yang dirancang sesuai usia dan kondisi anak, fisioterapis membantu melatih kekuatan otot, keseimbangan, dan koordinasi gerak — dengan pendekatan yang menyenangkan agar anak tetap nyaman menjalani sesi terapi.",
				"Penanganan sejak dini pada gangguan motorik kasar dapat mencegah kompensasi gerak yang keliru di kemudian hari, sehingga anak dapat beraktivitas dan bermain sesuai dengan tahapan usianya.",
				"Jika anak menunjukkan keterlambatan dalam pencapaian motorik kasar dibanding rentang usia normal, konsultasi dengan fisioterapis anak dapat membantu memastikan penanganan yang tepat sejak awal.",
			}},
		{Slug: "terapi-okupasi-untuk-anak-abk", Title: "ABK dan Pentingnya Terapi Okupasi Sejak Dini", Category: "Terapi Okupasi",
			Excerpt: "Bagi anak berkebutuhan khusus, terapi okupasi berperan besar dalam melatih kemandirian dan kesiapan mengikuti aktivitas belajar.",
			Content: []string{
				"Bagi anak berkebutuhan khusus (ABK), kemandirian dalam aktivitas sehari-hari sering menjadi tantangan tersendiri — mulai dari makan sendiri, memakai baju, hingga memegang alat tulis. Di sinilah terapi okupasi berperan penting.",
				"Terapi okupasi berfokus pada melatih activity daily living skill (ADL) atau kemandirian sehari-hari, mengembangkan motorik halus, serta mempersiapkan anak mengikuti aktivitas belajar di sekolah.",
				"Pendekatan terapi disesuaikan dengan kondisi masing-masing anak — bisa melalui permainan sensorik, latihan menulis bertahap, atau aktivitas yang melatih koordinasi dan perhatian anak.",
				"Konsistensi menjadi kunci keberhasilan terapi okupasi. Selain sesi rutin bersama terapis, dukungan orang tua dengan melanjutkan latihan sederhana di rumah sangat membantu mempercepat perkembangan anak.",
				"Setiap anak berkebutuhan khusus memiliki potensi masing-masing. Dengan penanganan yang tepat dan konsisten, anak dapat mencapai kemandirian sesuai kemampuan optimalnya.",
			}},
	}
	for i, a := range items {
		a.SortOrder = i
		saveArticle(a)
	}
	log.Println("seeded articles")
}

func seedFacilities() {
	items := []Facility{
		{Title: "Ruang Asesmen & Psikolog Anak", Description: "Ruang evaluasi tumbuh kembang anak dan sesi konseling psikologi bersama psikolog anak profesional yang dirancang tenang.", Icon: "growth"},
		{Title: "Ruang Konsultasi", Description: "Ruang diskusi pribadi untuk orang tua berkonsultasi mengenai hasil asesmen dan rencana terapi si kecil bersama tim.", Icon: "clipboard"},
		{Title: "Ruang Terapi Wicara", Description: "Ruang khusus yang tenang untuk melatih komunikasi verbal, artikulasi, dan kemampuan bahasa reseptif-ekspresif.", Icon: "speech"},
		{Title: "Ruang Sensori Integrasi", Description: "Dilengkapi dengan ayunan, perosotan, dan jaring sensori untuk melatih koordinasi motorik dan kepekaan indra anak.", Icon: "puzzle"},
		{Title: "Ruang Fisioterapi", Description: "Ruang terapi fisik anak dengan peralatan latihan kekuatan otot, keseimbangan, serta stimulasi motorik kasar.", Icon: "hands"},
		{Title: "Ruang Ortopedagogik (Terapi Belajar)", Description: "Ruang pendampingan belajar khusus secara individual untuk melatih kesiapan sekolah, kognitif, dan akademis anak.", Icon: "book"},
		{Title: "Ruang Eksplorasi & Bermain", Description: "Area eksplorasi motorik kasar dan bermain interaktif untuk menumbuhkan rasa percaya diri dan sosialisasi anak.", Icon: "run"},
		{Title: "Ruang Tunggu Utama", Description: "Area tunggu yang nyaman dan luas di lantai bawah lengkap dengan sarana informasi layanan dan administrasi.", Icon: "heart"},
		{Title: "Ruang Tunggu Lantai Atas", Description: "Area tunggu lantai dua yang tenang untuk orang tua/pengantar selama mendampingi proses terapi anak.", Icon: "heart"},
		{Title: "Ruang Administrasi & Pendaftaran", Description: "Layanan pendaftaran, pengaturan jadwal terapi, administrasi, dan pusat informasi bagi calon orang tua.", Icon: "clipboard"},
		{Title: "Musholla", Description: "Fasilitas musholla yang bersih dan tenang untuk kenyamanan beribadah bagi keluarga pasien maupun staf.", Icon: "heart"},
		{Title: "Area Parkir", Description: "Area parkir yang luas, aman, dan mudah diakses untuk kenyamanan orang tua dan kendaraan selama kunjungan.", Icon: "heart"},
	}
	for i, f := range items {
		f.SortOrder = i
		saveFacility(f)
	}
	log.Println("seeded facilities (add photos in the dashboard — image files were not copied from the frontend repo)")
}

func seedGallery() {
	items := []GalleryPhoto{
		{Category: "aktivitas", Caption: "Terapi Wicara"},
		{Category: "aktivitas", Caption: "Terapi Oral Motor"},
		{Category: "aktivitas", Caption: "Konsultasi & Asesmen Tumbuh Kembang"},
		{Category: "aktivitas", Caption: "Terapi Okupasi"},
		{Category: "aktivitas", Caption: "Fisioterapi"},
		{Category: "aktivitas", Caption: "Stimulasi Sensori Motorik"},
		{Category: "aktivitas", Caption: "Latihan Motorik Kasar"},
		{Category: "aktivitas", Caption: "Stimulasi Anak"},
		{Category: "aktivitas", Caption: "Terapi Perilaku"},
		{Category: "fasilitas", Caption: "Ruang Asesmen & Psikolog Anak"},
		{Category: "fasilitas", Caption: "Ruang Sensori Integrasi"},
		{Category: "fasilitas", Caption: "Ruang Terapi Wicara"},
		{Category: "fasilitas", Caption: "Ruang Fisioterapi"},
		{Category: "fasilitas", Caption: "Ruang Ortopedagogik (Terapi Belajar)"},
		{Category: "fasilitas", Caption: "Ruang Konsultasi"},
		{Category: "fasilitas", Caption: "Ruang Eksplorasi & Bermain"},
		{Category: "fasilitas", Caption: "Ruang Tunggu Utama"},
		{Category: "fasilitas", Caption: "Musholla"},
		{Category: "fasilitas", Caption: "Area Parkir"},
	}
	for i, g := range items {
		g.SortOrder = i
		saveGalleryPhoto(g)
	}
	log.Println("seeded gallery captions (add photos in the dashboard — image files were not copied from the frontend repo)")
}

func seedTeam() {
	items := []TeamMember{
		{Name: "Octariana P., M.Psi., Psikolog", Role: "Psikolog Anak"},
		{Name: "Risfa, A.Md.TW", Role: "Terapis Wicara"},
		{Name: "Eka Nur Novitasari, S.Tr.OT", Role: "Terapis Okupasi"},
		{Name: "Haya, S.Tr.OT", Role: "Terapis Okupasi"},
		{Name: "Janu C. Tyas, A.Md.OT", Role: "Terapis Okupasi"},
		{Name: "Novianan L., S.Pd. PLB.", Role: "Guru Ortopedagogik"},
		{Name: "Dewi L. Fauzia, A.Md.Fis", Role: "Fisioterapis"},
	}
	for i, t := range items {
		t.SortOrder = i
		saveTeamMember(t)
	}
	log.Println("seeded team (add photos in the dashboard — image files were not copied from the frontend repo)")
}

func seedBranches() {
	saveBranch(Branch{
		Slug: "lamongan", Name: "GenSA Kidz Lamongan",
		Address:  "Ruko Tambakboyo Regency No. 01–02, Tikung, Lamongan, Jawa Timur 62281",
		WhatsApp: "6281311992012", Phone: "0322314966",
		MapsQuery: "GenSA+Kidz+Terapi+Tumbuh+Kembang+Anak+Lamongan",
		MapsURL:   "https://www.google.com/maps/place/GenSA+Kidz+Terapi+Tumbuh+Kembang+Anak+Lamongan/@-7.1305286,112.4225374,17z",
		Schedules: []ScheduleItem{
			{Days: "Senin – Jumat", Hours: "08:00 – 16:00 WIB"},
			{Days: "Sabtu", Hours: "08:00 – 15:00 WIB"},
		},
	})
	log.Println("seeded branches")
}

func seedTestimonials() {
	items := []Testimonial{
		{Name: "Fany R.", Role: "Orang tua pasien terapi wicara", Quote: "Alhamdulillah setelah terapi 1 bulan di GenSA Kidz, Esa banyak sekali perubahannya — lebih cepat menangkap, konsentrasinya mulai bagus, dan kosakatanya lebih banyak. Saya sebagai orang tua juga dapat banyak ilmu dari terapisnya yang baik dan sabar."},
		{Name: "Ika C.", Role: "Orang tua pasien terapi perilaku", Quote: "GenSA Kidz the best! Anak saya yang awalnya hiperaktif dan susah diajak komunikasi, sekarang hipernya mulai menurun dan mulai mengerti kalau dikasih tahu. Baru 3 kali pertemuan sudah kelihatan hasilnya."},
		{Name: "Mega K.", Role: "Orang tua pasien fisioterapi", Quote: "Awalnya putri kami sangat minim motorik, sensori, dan belum kuat secara postural. Setelah 2 bulan belajar di GenSA Kidz, banyak sekali kemajuan yang ditunjukkan. Terima kasih GenSA Kidz."},
		{Name: "Sinta D.", Role: "Orang tua pasien terapi wicara", Quote: "Anak saya sudah terapi selama sebulan dan perubahannya signifikan — dari yang tadinya dipanggil tidak merespons, sekarang sudah mengerti kalau dipanggil namanya."},
		{Name: "Maria U.", Role: "Orang tua pasien stimulasi anak", Quote: "Bersyukur sekali dipertemukan dengan GenSA Kidz. Anak kami yang awalnya kontak mata tidak ada dan konsentrasi minim, sekarang sudah cerewet, suka bercerita, dan konsentrasinya mulai terbentuk."},
		{Name: "Orang tua Sheryl", Role: "Orang tua pasien terapi okupasi", Quote: "Sheryl awalnya tipe anak yang tidak mau dipegang orang lain, tapi terapisnya super sabar sampai Sheryl mau diarahkan belajar motorik kasar dan halusnya. Sekarang gerakannya lebih kalem dan konsentrasinya lebih dapat."},
	}
	for i, t := range items {
		t.SortOrder = i
		saveTestimonial(t)
	}
	log.Println("seeded testimonials")
}

func seedSettings() {
	setSetting("hero_badge", "Lamongan · Sejak 2020")
	setSetting("hero_title", "Setiap anak punya jalur tumbuh kembangnya sendiri.")
	setSetting("hero_subtitle", "GenSA Kidz mendampingi anak usia 0–16 tahun — baik dengan perkembangan umum maupun kebutuhan khusus (ABK) — lewat asesmen menyeluruh dan program terapi yang dipersonalisasi, 1 atap di Lamongan.")
	setSetting("about_text", "GenSA Kidz adalah pusat layanan terapi dan stimulasi tumbuh kembang anak di Lamongan. Sejak 2020, kami mendampingi anak dengan perkembangan umum maupun anak berkebutuhan khusus (ABK) melalui deteksi dini, stimulasi, dan penanganan klinis yang terpadu — mulai dari konsultasi awal hingga sesi terapi rutin, semua dalam satu tempat.")
	setSetting("visi", "Menjadi pusat layanan tumbuh kembang anak secara inklusif yang terintegrasi dengan program keberlanjutan anak sesuai dengan potensi perkembangannya, serta memberikan pendampingan terbaik kepada keluarga dalam mewujudkan masa depan anak secara optimal.")
	setSetting("misi", "Memberikan layanan deteksi dini yang komprehensif melalui asesmen menyeluruh guna memetakan potensi dan profil perkembangan unik setiap anak secara akurat sejak dini.\nMenyusun program stimulasi dan intervensi yang dipersonalisasi guna memastikan setiap anak mendapatkan penanganan yang spesifik dan tepat sasaran sesuai dengan fase perkembangannya.\nMenghadirkan pusat layanan terpadu satu atap yang mengintegrasikan berbagai disiplin terapi dan dukungan psikologi anak secara sinergis, memudahkan akses bagi orang tua dalam satu lokasi.\nMenyelaraskan kolaborasi tim ahli lintas disiplin untuk memberikan solusi penanganan yang holistik, berkesinambungan, dan terkoordinasi bagi setiap tantangan tumbuh kembang anak.\nMemberdayakan peran orang tua sebagai mitra utama melalui edukasi berkelanjutan dan komunikasi transparan, demi keberhasilan perkembangan anak yang konsisten baik di layanan maupun di rumah.\nMenciptakan lingkungan yang inklusif, aman, dan penuh kasih yang menghargai keberagaman cara belajar setiap anak, sehingga mereka merasa nyaman dalam berekspresi dan berkembang.\nMenjaga standar keunggulan layanan secara berkelanjutan dengan senantiasa mengikuti perkembangan ilmu pengetahuan dan teknologi terkini di bidang kesehatan dan tumbuh kembang anak.")
	setSetting("tim_intro", "Setiap anak mendapatkan penanganan yang komprehensif melalui kolaborasi lintas divisi. Evaluasi tumbuh kembang dilakukan secara kolektif oleh tim terapis ahli dan divalidasi langsung oleh Psikolog Klinis untuk memastikan ketajaman program terapi yang dipersonalisasi.")
	setSetting("ig_url", "https://www.instagram.com/gensakidz/")
	setSetting("fb_url", "https://www.facebook.com/GenSAKidz/")
	setSetting("tiktok_url", "https://www.tiktok.com/@gensa.kidz")
	setSetting("linkedin_url", "https://www.linkedin.com/company/gensa-kidz")
	setSetting("apply_email", "resource@gensakidz.com")
	log.Println("seeded settings")
}
