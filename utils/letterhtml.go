package utils

import (
	"dontkeep/futureme-backend/models"
	"fmt"
)

// BuildLetterHTML returns a styled HTML email body for the letter
func BuildLetterHTML(letter models.Letter, body string) string {
	return fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; background: #f9fafb; padding: 32px;">
		  <div style="max-width:600px;margin:auto;background:#fff;border-radius:16px;box-shadow:0 2px 8px #0001;padding:32px;">
			<h1 style="color:#22223b;font-size:2.2rem;margin-bottom:0.5rem;letter-spacing:1px;">
			  W<span style="color:#a3c4f3;">hi</span>sper
			</h1>
			<p style="color:#555;font-size:1rem;margin-bottom:1.5rem;">
			  The following is a letter from <b>%s</b>, delivered from the past by <b>Whisper</b>
			</p>
			<hr style="margin:1.5rem 0;">
			<div style="color:#222;font-size:1.1rem;white-space:pre-line;">%s</div>
			<hr style="margin:1.5rem 0;">
			<footer style="color:#888;font-size:0.9rem;text-align:center;">
			  &copy; 2025 made with <span style="color:#a3c4f3;">&#10084;</span> by Loyalty
			</footer>
		  </div>
		</div>
	`,
		letter.CreatedAt.Format("January 02, 2006"), body)
}
