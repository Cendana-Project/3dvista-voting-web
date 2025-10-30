package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"voteweb/internal/domain"
	"voteweb/internal/http/middleware"
)

type PageHandler struct {
	service domain.VoteService
	logger  *slog.Logger
}

func NewPageHandler(service domain.VoteService, logger *slog.Logger) *PageHandler {
	return &PageHandler{
		service: service,
		logger:  logger,
	}
}

func (h *PageHandler) ShowInnovation(c *gin.Context) {
	groupSlug := c.Param("group")
	slug := c.Param("slug")

	// Get innovation
	innovation, err := h.service.GetInnovation(c.Request.Context(), groupSlug, slug)
	if err != nil {
		if err == domain.ErrInnovationNotFound {
			c.HTML(http.StatusNotFound, "error.tmpl.html", gin.H{
				"Title":   "Innovation Not Found",
				"Message": "The innovation you're looking for does not exist.",
			})
			return
		}
		h.logger.ErrorContext(c.Request.Context(), "failed to get innovation",
			"group_slug", groupSlug,
			"slug", slug,
			"error", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl.html", gin.H{
			"Title":   "Error",
			"Message": "An error occurred while loading the page.",
		})
		return
	}

	// Get current vote count
	voteCount, err := h.service.GetVoteCount(c.Request.Context(), innovation.ID)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "failed to get vote count",
			"innovation_id", innovation.ID,
			"error", err)
		voteCount = 0
	}

	// Check if user has already voted
	clientIP := c.GetString("client_ip")
	hasVoted := false
	if clientIP != "" {
		voted, err := h.service.CheckHasVoted(c.Request.Context(), innovation.ID, clientIP)
		if err != nil {
			h.logger.ErrorContext(c.Request.Context(), "failed to check vote status",
				"innovation_id", innovation.ID,
				"error", err)
		} else {
			hasVoted = voted
		}
	}

	// Get CSRF token for the page
	csrfToken := middleware.GetCSRFToken(c)

	// Resolve hero asset (webp) based on group + slug
	var hero string
	var heroMobile string
	if innovation.GroupSlug == "bumn-bumd" {
		switch innovation.Slug {
		case "alat-pemecah-ombak-apo-desa-mayangan-subang":
			hero = "/static/bumn-bumd/6.webp"
			heroMobile = "/static/bumn-bumd/6-mobile.webp"
		case "simotip":
			hero = "/static/bumn-bumd/simotip.webp"
			heroMobile = "/static/bumn-bumd/simotip-mobile.webp"
		case "aplikasi-pemilu-elektronik-e-voting":
			hero = "/static/bumn-bumd/evoting.webp"
			heroMobile = "/static/bumn-bumd/evoting-mobile.webp"
		case "thr-asyik":
			hero = "/static/bumn-bumd/thr-asyik.webp"
			heroMobile = "/static/bumn-bumd/thr-asyik-mobile.webp"
		}
	} else if innovation.GroupSlug == "kementrian-lembaga-pt" {
		switch innovation.Slug {
		case "isopa-intelligent-solar-panel":
			hero = "/static/kementrian-lembaga-pt/isopa.webp"
			heroMobile = "/static/kementrian-lembaga-pt/isopa-mobile.webp"
		case "instrumen-deteksi-risiko-stunting-pada-remaja-insting":
			hero = "/static/kementrian-lembaga-pt/insting.webp"
			heroMobile = "/static/kementrian-lembaga-pt/insting-mobile.webp"
		case "teknologi-hybrid-taman-sanitasi-hts-untuk-pencegahan-pencemaran-lingkungan-dan-daur-ulang-air":
			hero = "/static/kementrian-lembaga-pt/hts.webp"
			heroMobile = "/static/kementrian-lembaga-pt/hts-mobile-1.webp" // corrected mobile asset
		case "inovasi-saschieversity":
			hero = "/static/kementrian-lembaga-pt/saschieversity.webp"
			heroMobile = "/static/kementrian-lembaga-pt/saschieversity-mobile.webp"
		case "mentari-mental-health-remaja-indonesia-assessment":
			hero = "/static/kementrian-lembaga-pt/mentari-assesment.webp"
			heroMobile = "/static/kementrian-lembaga-pt/mentari-assesment.webp"
		}
	} else if innovation.GroupSlug == "pemprov-jabar" {
		switch innovation.Slug {
		case "jabar-digital-academy":
			hero = "/static/pemprov-jabar/jabar-istimewa-digital-academy.webp"
			heroMobile = "/static/pemprov-jabar/jabar-istimewa-digital-academy-mobile.webp"
		case "delman-sarah-model-pemeliharaan-sapi-perah-di-jawa-barat":
			hero = "/static/pemprov-jabar/new-normal-persusuan-jawa-barat.webp"
			heroMobile = "/static/pemprov-jabar/new-normal-persusuan-jawa-barat-mobile.webp"
		case "jabar-form":
			hero = "/static/pemprov-jabar/jabar-form.webp"
			heroMobile = "/static/pemprov-jabar/jabar-form-mobile.webp"
		case "gisa-prima-adminduk-jabar":
			hero = "/static/pemprov-jabar/gisa-prima.webp"
			heroMobile = "/static/pemprov-jabar/gisa-prima-mobile.webp"
		case "data-potensi-digital-desa-tapal-desa":
			hero = "/static/pemprov-jabar/tapal-desa.webp"
			heroMobile = "/static/pemprov-jabar/tapal-desa-mobile.webp"
		}
	} else if innovation.GroupSlug == "smp-sma-sederajat" {
		switch innovation.Slug {
		case "penguatan-kompetensi-litnum-melalui-lesson-study":
			hero = "/static/smp-sma/litnum.webp"
			heroMobile = "/static/smp-sma/litnum-mobile.webp"
		case "motor-lstrik-dengan-teknologi-finger-print":
			hero = "/static/smp-sma/motor-listrik-dengan-teknologi-finger-print.webp"
			heroMobile = "/static/smp-sma/motor-listrik-dengan-teknologi-finger-print-mobile.webp"
		case "samving-block-sampah-plastik-menjadi-paving-block":
			hero = "/static/smp-sma/samving-block.webp"
			heroMobile = "/static/smp-sma/samving-block-mobile.webp"
		case "inovasi-sabun-nanas-tsanawiyah-satu":
			hero = "/static/smp-sma/sanatsu.webp"
			heroMobile = "/static/smp-sma/sanatsu-mobile.webp"
		case "tonnetar-tongkat-tunanetra-pintar":
			hero = "/static/smp-sma/tonnetar.webp"
			heroMobile = "/static/smp-sma/tonnetar-mobile.webp"
		}
	} else if innovation.GroupSlug == "pemda-kabupaten" {
		switch innovation.Slug {
		case "si-pintar-online":
			hero = "/static/pemda-jabar/pintar-on-line.webp"
			heroMobile = "/static/pemda-jabar/pintar-on-line-mobile.webp"
		case "ekonomi-bangit-harapan-terbit-si-dara-puber-buka-jalan-sejahtera-untuk-5-260-orang-miskin-di-kabupaten-sumedang-sistem-pemberdayaan-masyarakat-miskin-dengan-pengembangan-ekonomi-produktif-melalui-kelompok-usaha-bersama":
			hero = "/static/pemda-jabar/sidara-puber.webp"
			heroMobile = "/static/pemda-jabar/sidara-puber-mobile.webp"
		case "sistem-informasi-manajemen-perlindungan-pertanian-simarlin":
			hero = "/static/pemda-jabar/simarlin.webp"
			heroMobile = "/static/pemda-jabar/simarlin-mobile.webp"
		case "nyai-indramayu-artificial-intelligence":
			hero = "/static/pemda-jabar/nyai.webp"
			heroMobile = "/static/pemda-jabar/nyai.webp"
		case "ngupahan-ngabagi-ngubah-ngurai-sampah-pangan-dinas-ketahanan-pangan-kab-bogor":
			hero = "/static/pemda-jabar/ngupahan.webp"
			heroMobile = "/static/pemda-jabar/ngupahan-mobile.webp"
		case "ketupat-lebaran-kegunaan-kartu-kepatuhan-minum-tablet-tambah-darah":
			hero = "/static/pemda-jabar/ketupat-lebaran.webp"
			heroMobile = "/static/pemda-jabar/ketupat-lebaran-mobile.webp"
		}
	} else if innovation.GroupSlug == "pemda-kota" {
		// City governments (pemkot) mappings
		switch innovation.Slug {
		case "smart-k-sistem-manajemen-akuakultur-rekayasa-teknologi-dan-kemitraan":
			hero = "/static/pemkot/smart-k.webp"
			heroMobile = "/static/pemkot/smart-k-mobile.webp"
		case "bung-senja-tabungan-sedot-tinja":
			hero = "/static/pemkot/buang-senja.webp"
			heroMobile = "/static/pemkot/buang-senja-mobile.webp"
		case "gerakan-orang-cimahi-pilah-sampah-grak-ompimpah":
			hero = "/static/pemkot/grak-ompimpah.webp"
			heroMobile = "/static/pemkot/grak-ompimpah-mobile.webp"
		case "bogor-smart-health":
			hero = "/static/pemkot/bogor-smart-health.webp"
			heroMobile = "/static/pemkot/bogor-smart-health-mobile.webp"
		case "konservasi-mata-air-menjadi-ruang-terbuka-hijau-ruang-publik":
			hero = "/static/pemkot/konversi-mata-air.webp"
			heroMobile = "/static/pemkot/konversi-mata-air-mobile.webp"
		}
	}

	c.HTML(http.StatusOK, "innovation.tmpl.html", gin.H{
		"Innovation": innovation,
		"VoteCount":  voteCount,
		"CSRFToken":  csrfToken,
		"HasVoted":   hasVoted,
		"Hero":       hero,
		"HeroMobile": heroMobile,
	})
}
