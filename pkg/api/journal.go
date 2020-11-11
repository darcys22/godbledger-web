package api

import (
	m "github.com/darcys22/godbledger-web/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetJournals(c *gin.Context) {
	journalsModel := m.NewJournalsListing()
	err := journalsModel.SearchJournals()
	if err != nil {
		log.Errorf("Could not get journal listing (%v)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(200, journalsModel)
}

func PostJournal(c *gin.Context) {
	var journal m.PostJournalCommand

	if err := c.BindJSON(&journal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := journal.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(200, journal)
}

func DeleteJournal(c *gin.Context) {
	id := c.Params.ByName("id")

	if err := m.DeleteJournalCommand(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.String(http.StatusOK, "Success")
}
