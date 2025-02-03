package handlers

import (
	"errors"
	"fmt"
	"go-conference/internal/database"
	"go-conference/internal/mail"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type Ticket struct {
	Name    string `json:"name" form:"name"`
	Email   string `json:"email" form:"email"`
	Tickets uint   `json:"tickets" form:"tickets"`
}

func Home(c *gin.Context) {
	fmt.Println("Rendering index.html...")
	c.HTML(http.StatusOK, "index.html", nil)
}

func TicketForm(c *gin.Context) {
	fmt.Println("Tickets page requested")
	c.HTML(http.StatusOK, "tickets.html", nil)
}

func ValidateUserInput(name string, email string, tickets uint) error {
	// Validate name length
	if len(name) < 3 {
		return errors.New("Name must be at least 3 characters long")
	}

	// Validate email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("Invalid email format")
	}

	// Validate ticket count
	if tickets < 1 || tickets > 10 {
		return errors.New("Ticket count must be between 1 and 10")
	}

	return nil // Input is valid
}

func BuyTicket(c *gin.Context) {
	var ticket Ticket
	if err := c.ShouldBind(&ticket); err != nil {
		c.HTML(http.StatusBadRequest, "tickets.html", gin.H{
			"ticket": ticket,
			"formErrors": map[string]string{
				"general": "Invalid input format",
			},
		})
		return
	}
	fmt.Printf("Received ticket data: Name: %s, Email: %s, Tickets: %d\n", ticket.Name, ticket.Email, ticket.Tickets)

	// Validate user input
	formErrors := make(map[string]string)
	fmt.Printf("Validating: Name=%s (len=%d), Email=%s, Tickets=%d\n",
		ticket.Name, len(ticket.Name), ticket.Email, ticket.Tickets)

	if err := ValidateUserInput(ticket.Name, ticket.Email, ticket.Tickets); err != nil {
		// Map specific error messages to their respective fields
		switch err.Error() {
		case "Name must be at least 3 characters long":
			formErrors["Name"] = err.Error()
		case "Invalid email format":
			formErrors["Email"] = err.Error()
		case "Ticket count must be between 1 and 10":
			formErrors["Tickets"] = err.Error()
		default:
			formErrors["general"] = err.Error()
		}
	}

	if len(formErrors) > 0 {
		c.HTML(http.StatusBadRequest, "tickets.html", gin.H{
			"ticket":     ticket,
			"formErrors": formErrors,
		})
		return
	}

	// Insert ticket into database
	query := "INSERT INTO tickets (name, email, tickets) VALUES ($1, $2, $3)"
	_, err := database.DB.Exec(query, ticket.Name, ticket.Email, ticket.Tickets)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Could not process ticket"})
		return
	}

	// Send confirmation email
	err = mail.SendConfirmation(ticket.Email, ticket.Name, ticket.Tickets)
	if err != nil {
		fmt.Printf("Error sending confirmation email: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Could not send confirmation email"})
		return
	}

	// Return confirmation page
	c.HTML(http.StatusOK, "confirmation.html", gin.H{
		"name":    ticket.Name,
		"email":   ticket.Email,
		"tickets": ticket.Tickets,
	})
}
