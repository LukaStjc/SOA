package controllers

import (
	"go-tourm/initializers"
	"go-tourm/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateTour(c *gin.Context) {
	var body struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Type        uint              `json:"type"`
		Tags        string            `json:"tags"`
		Price       float64           `json:"price"`
		AvgRate     float64           `json:"avgRate"`
		UserID      uint              `json:"userId"`
		KeyPoints   []models.KeyPoint `json:"keyPoints"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	tour := models.Tour{
		Name:        body.Name,
		Description: body.Description,
		Type:        models.TourType(body.Type),
		Tags:        body.Tags,
		Price:       body.Price,
		AvgRate:     body.AvgRate,
		UserID:      body.UserID,
	}

	result := initializers.DB.Create(&tour)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create tour",
		})
		return
	}

	for _, kp := range body.KeyPoints {
		kp.TourID = tour.ID
		if kpResult := initializers.DB.Create(&kp); kpResult.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create keypoints",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"data": tour})
}

// Ovaj endpoint bi mogao da posluzi da se svakome ko se uloguje prikazu
// sve ture. I pored ture da stoji textbox sa dugmetom i mogucnoscu unosa ocene
// to je enabled samo ako se ovde pri vracanju tih tura kaze da na osnovu
// potrazivanja ulogovanog korisnika ima pravo da postavi ocenu
// ta logika ce se takodje ovde proveriti
// => prosiriti ovaj endpoint
func GetToursByUser(c *gin.Context) {
	// moram biti kupac te ture

	// mogu ostaviti ocenu samo ako u bazi imam da sam posetio vise od polovine tacaka

	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var tours []models.Tour
	result := initializers.DB.Preload("KeyPoints").Where("user_id = ?", userID).Find(&tours)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tours"})
		return
	}

	// kontam da je najlakse da se vracaju parovi tours i bool vrednost
	// da li sme ili ne sme da ostavi ocenu
	c.JSON(http.StatusOK, gin.H{"tours": tours})
}

// U ovom endpointu cemo samo ostavljati ocenu bez ikakvih provera uz azuriranje ocene
// A u onom endpointu gde prikazujemo ture cemo prikazivati sve ture
// bez obzira ko je vlasnik. I pored toga cemo imati textbox za unos ocene
// kao i dugme za potvrdu ukoliko korisnik koji je ulogovan zadovoljava
// uslove ocenjivanja te ture, i nakon toga se uradi refresh strane
// odnosno ponovo se dobave ture sa azuriranom ocenom.
func PostReview(c *gin.Context) {
	var body struct {
		TourID uint
		Rate   uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	} // mogu jos i neke validacije, a i ne moraju posto ce to pokriti front

	// dobavi korisnika
	authUser, _ := c.Get("user")
	user := authUser.(models.User)

	// poenta je da se izvuku kljucne tacke koje postoje i koje su jedinstvene
	// za konkretnu turu, i onda da se vidi u tour check ins da li je korisnik
	// otkacio vise od polovine tura tako sto se u polju date proveri da li je vece
	// od pocetnog datuma to se proverava sa isZero() nesto tako
	// kada se utvrdi to je to, bice malo vise forova, jer sam tako projektovao semu
	// ali moze da se dobavi, nije nista specijalno

	// ostavi ocenu za turu
	// userTourCheckIns := initializers.DB.Where("userid = ?", user.ID).Table("tour_check_ins")
	// filtriraj samo one koji su validni, koje si prethodno posetio

	// na osnovu ture imas i keypointse
	// userTour := initializers.DB.Where("tourid = ?", body.TourID).First()

	// podesim ocenu u orderu

	// automatski azuriram avgrate te ture, prelazeci preko svih ordera haha

}

// Nebitna
func Validate(c *gin.Context) {

	// ako hoces da dobavis neko polje usera
	// onda zapocinjes komandu na sledeci nacin
	// user.(models.User).

	userInterface, _ := c.Get("user")

	// NE SME *models.User!
	user, _ := userInterface.(models.User)

	if user.Role.String() == "Administrator" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Administrator Content.",
		})
	} else if user.Role.String() == "Guide" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Guide Content.",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Tourist Content.",
		})
	}
}
