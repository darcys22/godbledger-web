package api

import (
	m "github.com/darcys22/godbledger-web/backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetJournals(ctx *gin.Context) {
	journalsModel := m.NewJournalsListing()
	err := journalsModel.SearchJournals()
	if err != nil {
		log.Errorf("Could not get journal listing (%v)", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, journalsModel)
}

func PostJournal(ctx *gin.Context) {
	var journal m.PostJournalCommand

	if err := ctx.BindJSON(&journal); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := journal.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, journal)
}

func DeleteJournal(ctx *gin.Context) {
	id := ctx.Params.ByName("id")

	if err := m.DeleteJournalCommand(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.String(http.StatusOK, "Success")
}

func GetJournal(ctx *gin.Context) {
	id := ctx.Params.ByName("id")

	journal, err := m.GetJournalCommand(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, journal)
}

func EditJournal(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	var journal m.PostJournalCommand

	if err := ctx.BindJSON(&journal); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := journal.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if err := m.DeleteJournalCommand(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, journal)
}
