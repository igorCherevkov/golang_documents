package utils

import "documents/internal/models"

func CanAccessDocument(doc models.Document, requesterUserID, requesterLogin string) bool {
	if doc.UserID == requesterUserID {
		return true
	}

	if doc.IsPublic {
		return true
	}

	for _, login := range doc.Grant {
		if login == requesterLogin {
			return true
		}
	}

	return false
}
