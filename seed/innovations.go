package seed

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InnovationData represents the seed data for an innovation
type InnovationData struct {
	GroupSlug string
	Name      string
	Slug      string
	Division  string
}

// InnovationsData contains all innovation seed data
var InnovationsData = []InnovationData{
	{
		GroupSlug: "pemprov-jabar",
		Name:      "Jabar Digital Academy",
		Slug:      "jabar-digital-academy",
		Division:  "PEMERINTAH PROVINSI JAWA BARAT",
	},
	{
		GroupSlug: "pemprov-jabar",
		Name:      "Data Potensi Digital Desa ( TAPAL DESA )",
		Slug:      "data-potensi-digital-desa-tapal-desa",
		Division:  "PEMERINTAH PROVINSI JAWA BARAT",
	},
	{
		GroupSlug: "pemprov-jabar",
		Name:      "GISA PRIMA ADMINDUK JABAR",
		Slug:      "gisa-prima-adminduk-jabar",
		Division:  "PEMERINTAH PROVINSI JAWA BARAT",
	},
	{
		GroupSlug: "pemprov-jabar",
		Name:      "Delman Sarah (Model Pemeliharaan Sapi Perah) di Jawa Barat",
		Slug:      "delman-sarah-model-pemeliharaan-sapi-perah-di-jawa-barat",
		Division:  "PEMERINTAH PROVINSI JAWA BARAT",
	},
	{
		GroupSlug: "pemprov-jabar",
		Name:      "Jabar Form",
		Slug:      "jabar-form",
		Division:  "PEMERINTAH PROVINSI JAWA BARAT",
	},
	{
		GroupSlug: "bumn-bumd",
		Name:      "Alat Pemecah Ombak (APO) Desa Mayangan Subang",
		Slug:      "alat-pemecah-ombak-apo-desa-mayangan-subang",
		Division:  "BUMN/BUMD",
	},
	{
		GroupSlug: "bumn-bumd",
		Name:      "SIMOTIP",
		Slug:      "simotip",
		Division:  "BUMN/BUMD",
	},
	{
		GroupSlug: "bumn-bumd",
		Name:      "Aplikasi Pemilu Elektronik (e-Voting)",
		Slug:      "aplikasi-pemilu-elektronik-e-voting",
		Division:  "BUMN/BUMD",
	},
	{
		GroupSlug: "bumn-bumd",
		Name:      "THR Asyik",
		Slug:      "thr-asyik",
		Division:  "BUMN/BUMD",
	},
	{
		GroupSlug: "kementrian-lembaga-pt",
		Name:      "Instrumen Deteksi Risiko Stunting pada Remaja (Insting)",
		Slug:      "instrumen-deteksi-risiko-stunting-pada-remaja-insting",
		Division:  "KEMENTERIAN/LEMBAGA/PT",
	},
	{
		GroupSlug: "kementrian-lembaga-pt",
		Name:      "Teknologi Hybrid Taman Sanitasi (HTS) untuk Pencegahan Pencemaran Lingkungan dan Daur Ulang Air",
		Slug:      "teknologi-hybrid-taman-sanitasi-hts-untuk-pencegahan-pencemaran-lingkungan-dan-daur-ulang-air",
		Division:  "KEMENTERIAN/LEMBAGA/PT",
	},
	{
		GroupSlug: "kementrian-lembaga-pt",
		Name:      "MENTARI (Mental Health Remaja Indonesia) Assessment",
		Slug:      "mentari-mental-health-remaja-indonesia-assessment",
		Division:  "KEMENTERIAN/LEMBAGA/PT",
	},
	{
		GroupSlug: "kementrian-lembaga-pt",
		Name:      "ISOPA (Intelligent Solar Panel)",
		Slug:      "isopa-intelligent-solar-panel",
		Division:  "KEMENTERIAN/LEMBAGA/PT",
	},
	{
		GroupSlug: "kementrian-lembaga-pt",
		Name:      "INOVASI SASCHIEVERSITY",
		Slug:      "inovasi-saschieversity",
		Division:  "KEMENTERIAN/LEMBAGA/PT",
	},
	{
		GroupSlug: "smp-sma-sederajat",
		Name:      "TONNETAR (Tongkat Tunanetra Pintar)",
		Slug:      "tonnetar-tongkat-tunanetra-pintar",
		Division:  "SMP/SMA/SEDERAJAT",
	},
	{
		GroupSlug: "smp-sma-sederajat",
		Name:      "Motor Lstrik Dengan Teknologi Finger Print",
		Slug:      "motor-lstrik-dengan-teknologi-finger-print",
		Division:  "SMP/SMA/SEDERAJAT",
	},
	{
		GroupSlug: "smp-sma-sederajat",
		Name:      "SAMVING BLOCK (SAMPAH PLASTIK MENJADI PAVING BLOCK)",
		Slug:      "samving-block-sampah-plastik-menjadi-paving-block",
		Division:  "SMP/SMA/SEDERAJAT",
	},
	{
		GroupSlug: "smp-sma-sederajat",
		Name:      "Penguatan Kompetensi LITNUM Melalui Lesson Study",
		Slug:      "penguatan-kompetensi-litnum-melalui-lesson-study",
		Division:  "SMP/SMA/SEDERAJAT",
	},
	{
		GroupSlug: "smp-sma-sederajat",
		Name:      "INOVASI SABUN NANAS TSANAWIYAH SATU",
		Slug:      "inovasi-sabun-nanas-tsanawiyah-satu",
		Division:  "SMP/SMA/SEDERAJAT",
	},
	{
		GroupSlug: "pemda-kota",
		Name:      "Konservasi Mata Air menjadi Ruang Terbuka Hijau/Ruang Publik",
		Slug:      "konservasi-mata-air-menjadi-ruang-terbuka-hijau-ruang-publik",
		Division:  "PEMERINTAH DAERAH KOTA",
	},
	{
		GroupSlug: "pemda-kota",
		Name:      "Bogor Smart Health",
		Slug:      "bogor-smart-health",
		Division:  "PEMERINTAH DAERAH KOTA",
	},
	{
		GroupSlug: "pemda-kota",
		Name:      "Gerakan Orang Cimahi Pilah Sampah (Grak Ompimpah)",
		Slug:      "gerakan-orang-cimahi-pilah-sampah-grak-ompimpah",
		Division:  "PEMERINTAH DAERAH KOTA",
	},
	{
		GroupSlug: "pemda-kota",
		Name:      "BUNG SENJA (Tabungan Sedot Tinja )",
		Slug:      "bung-senja-tabungan-sedot-tinja",
		Division:  "PEMERINTAH DAERAH KOTA",
	},
	{
		GroupSlug: "pemda-kota",
		Name:      "SMART - K (Sistem Manajemen Akuakultur, Rekayasa Teknologi dan Kemitraan)",
		Slug:      "smart-k-sistem-manajemen-akuakultur-rekayasa-teknologi-dan-kemitraan",
		Division:  "PEMERINTAH DAERAH KOTA",
	},
	{
		GroupSlug: "pemda-kabupaten",
		Name:      "Si Pintar Online",
		Slug:      "si-pintar-online",
		Division:  "PEMERINTAH DAERAH KABUPATEN",
	},
	{
		GroupSlug: "pemda-kabupaten",
		Name:      "EKONOMI BANGKIT HARAPAN TERBIT: SI DARA PUBER BUKA JALAN SEJAHTERA UNTUK 5.260 ORANG MISKIN DI KABUPATEN SUMEDANG (SISTEM PEMBERDAYAAN MASYARAKAT MISKIN DENGAN PENGEMBANGAN EKONOMI PRODUKTIF MELALUI KELOMPOK USAHA BERSAMA)",
		Slug:      "ekonomi-bangit-harapan-terbit-si-dara-puber-buka-jalan-sejahtera-untuk-5-260-orang-miskin-di-kabupaten-sumedang-sistem-pemberdayaan-masyarakat-miskin-dengan-pengembangan-ekonomi-produktif-melalui-kelompok-usaha-bersama",
		Division:  "PEMERINTAH DAERAH KABUPATEN",
	},
	{
		GroupSlug: "pemda-kabupaten",
		Name:      "Sistem Informasi Manajemen Perlindungan Pertanian (SIMARLIN)",
		Slug:      "sistem-informasi-manajemen-perlindungan-pertanian-simarlin",
		Division:  "PEMERINTAH DAERAH KABUPATEN",
	},
	{
		GroupSlug: "pemda-kabupaten",
		Name:      "KETUPAT LEBARAN (Kegunaan Kartu Kepatuhan Minum Tablet Tambah darah)",
		Slug:      "ketupat-lebaran-kegunaan-kartu-kepatuhan-minum-tablet-tambah-darah",
		Division:  "PEMERINTAH DAERAH KABUPATEN",
	},
	{
		GroupSlug: "pemda-kabupaten",
		Name:      "NGUPAHAN (NGABAGI, NGUBAH, NGURAI, SAMPAH PANGAN) DINAS KETAHANAN PANGAN KAB.BOGOR",
		Slug:      "ngupahan-ngabagi-ngubah-ngurai-sampah-pangan-dinas-ketahanan-pangan-kab-bogor",
		Division:  "PEMERINTAH DAERAH KABUPATEN",
	},
	{
		GroupSlug: "pemda-kabupaten",
		Name:      "NYAI (INDRAMAYU ARTIFICIAL INTELLIGENCE)",
		Slug:      "nyai-indramayu-artificial-intelligence",
		Division:  "PEMERINTAH DAERAH KABUPATEN",
	},
}

// SeedInnovations inserts innovation seed data into the database
func SeedInnovations(ctx context.Context, pool *pgxpool.Pool) error {
	query := `
		INSERT INTO innovations (group_slug, slug, name, division, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (group_slug, slug) DO UPDATE
		SET name = EXCLUDED.name,
		    division = EXCLUDED.division,
		    updated_at = NOW()
	`

	log.Println("Starting innovation seeding...")

	for i, innovation := range InnovationsData {
		_, err := pool.Exec(ctx, query,
			innovation.GroupSlug,
			innovation.Slug,
			innovation.Name,
			innovation.Division,
		)
		if err != nil {
			return fmt.Errorf("failed to seed innovation %s: %w", innovation.Name, err)
		}
		log.Printf("Seeded %d/%d: %s", i+1, len(InnovationsData), innovation.Name)
	}

	log.Printf("Successfully seeded %d innovations", len(InnovationsData))
	return nil
}




