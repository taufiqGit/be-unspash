package repositories

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

type EmailRepository interface {
	SendVerificationEmail(ctx context.Context, to, name, verificationLink string, token string) error
	SendResetPasswordEmail(ctx context.Context, to, name, resetLink string) error
	SendOrderConfirmation(ctx context.Context, to, orderID string, items []string) error
}

type emailRepository struct {
	client *resend.Client
	from   string
}

func NewEmailRepository() EmailRepository {
	apiKey := os.Getenv("RESEND_API_KEY")
	from := os.Getenv("EMAIL_FROM") // misal: noreply@yourdomain.com

	client := resend.NewClient(apiKey)

	return &emailRepository{
		client: client,
		from:   from,
	}
}

func (r *emailRepository) SendVerificationEmail(ctx context.Context, to, name, verificationLink string, token string) error {
	params := &resend.SendEmailRequest{
		From: r.from,
		To:   []string{to},
		Template: &resend.EmailTemplate{
			Id: "email-verification",
			Variables: map[string]interface{}{
				"NAME": name,
				"LINK": fmt.Sprintf("%s?token=%s", verificationLink, token),
			},
		},
	}

	_, err := r.client.Emails.Send(params)
	return err
}

func (r *emailRepository) SendResetPasswordEmail(ctx context.Context, to, name, resetLink string) error {
	// Similar implementation
	return nil
}

func (r *emailRepository) SendOrderConfirmation(ctx context.Context, to, orderID string, items []string) error {
	// Similar implementationn
	return nil
}
