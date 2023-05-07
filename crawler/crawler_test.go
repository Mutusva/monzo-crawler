package crawler

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func TestFindLinks(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		filters []string
		exp     []string
		scheme  string
		host    string
	}{
		{
			name: "No links",
			html: `
             <!DOCTYPE html>
				<html>
				<body>
				<h1>HTML Links</h1>
				</body>
              </html>
             `,
			exp: []string{},
		},
		{
			name: "One links and no filter",
			html: `
             <!DOCTYPE html>
				<html>
				<body>
				<h1>HTML Links</h1>
				<p><a href="https://monzo.com">Visit monzo!</a></p>
				</body>
              </html>
             `,
			scheme: "https",
			host:   "monzo.com",
			exp:    []string{"https://monzo.com"},
		},
		{
			name: "One links and one filter",
			html: `
             <!DOCTYPE html>
				<html>
				<body>
				<h1>HTML Links</h1>
				<p><a href="https://monzo.com">Visit monzo!</a></p>
				</body>
              </html>
             `,
			scheme:  "https",
			host:    "monzo.com",
			filters: []string{"monzo.com"},
			exp:     []string{"https://monzo.com"},
		},
		{
			name: "One links and one filter",
			html: `
             <!DOCTYPE html>
				<html>
				<body>
				<h1>HTML Links</h1>
				<p><a href="https://community.monzo.com">Visit monzo</a></p>
				</body>
              </html>
             `,
			host:    "monzo.com",
			scheme:  "https",
			filters: []string{"monzo.com"},
			exp:     []string{},
		},
		{
			name: "relative links and one filter",
			html: `
             <!DOCTYPE html>
				<html>
				<body>
				<h1>HTML Links</h1>
				<p><a href="https://community.monzo.com">Visit community</a></p>
                <div>
                  <ul>
					<li><a href="/feedback/premium">feedback</a></li>
				</ul>
                </div>
				</body>
              </html>
             `,
			host:    "monzo.com",
			scheme:  "https",
			filters: []string{"monzo.com"},
			exp:     []string{"https://monzo.com/feedback/premium"},
		},
		{
			name: "Multiple links and filters",
			html: `
             <!DOCTYPE html>
				<html>
				<body>
				<h1>HTML Links</h1>
                 <div>
					<p><a href="https://community.monzo.com">Visit W3Schools.com!</a></p>
     				<div>
                	  <p><a href="https://monzo.com/about">Visit W3Schools.com!</a></p>
 					</div>
                 </div>
                <p><a href="https://monzo.com/contact">Visit W3Schools.com!</a></p>
                <p><a href="https://facebook.com">Visit W3Schools.com!</a></p>
				</body>
              </html>
             `,
			host:    "monzo.com",
			scheme:  "https",
			filters: []string{"monzo.com"},
			exp: []string{
				"https://monzo.com/about",
				"https://monzo.com/contact",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(*testing.T) {
			doc, err := html.Parse(strings.NewReader(tc.html))
			assert.Nil(t, err)
			got := findLinks(doc, tc.filters, tc.scheme, tc.host)
			assert.ElementsMatch(t, tc.exp, got)
		})
	}
}
