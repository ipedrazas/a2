package security

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// NetworkCheck detects suspicious network operations that could indicate data exfiltration.
type NetworkCheck struct {
	Patterns    map[string][]*regexp.Regexp
	SafeDomains map[string]bool
}

// ID returns the unique identifier for this check.
func (c *NetworkCheck) ID() string {
	return "security:network"
}

// Name returns the human-readable name for this check.
func (c *NetworkCheck) Name() string {
	return "Network Exfiltration Detection"
}

// Run executes the network exfiltration detection check.
func (c *NetworkCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	// Initialize patterns and safe domains if not already done
	if c.Patterns == nil {
		c.Patterns = c.getPatterns()
	}
	if c.SafeDomains == nil {
		c.SafeDomains = c.getSafeDomains()
	}

	// Scan all source files
	findings := c.scanDirectory(path)

	if len(findings) == 0 {
		return rb.Pass("No suspicious network operations detected"), nil
	}

	// Format findings into message
	msg := c.formatFindings(findings)
	return rb.Warn(msg), nil // Network is warn by default (network might be legit)
}

// getSafeDomains returns a list of known safe/common domains.
func (c *NetworkCheck) getSafeDomains() map[string]bool {
	return map[string]bool{
		// Cloud providers
		"amazonaws.com":      true,
		"amazonaws.com.cn":   true,
		"azure.com":          true,
		"azure-devices.net":  true,
		"googleapis.com":     true,
		"google.com":         true,
		"gcp.gvt2.com":       true,
		"cloudfunctions.net": true,
		"alistandard.com":    true,
		// CDN and static assets
		"cloudflare.com":         true,
		"cloudflareinsights.com": true,
		"cloudinary.com":         true,
		"akamai.net":             true,
		"akamaihd.net":           true,
		// Package registries
		"npmjs.org":              true,
		"yarnpkg.com":            true,
		"registry.npmjs.org":     true,
		"pypi.org":               true,
		"files.pythonhosted.org": true,
		"crates.io":              true,
		"repo.maven.apache.org":  true,
		"repo1.maven.org":        true,
		// Monitoring/analytics
		"segment.io":           true,
		"segment.com":          true,
		"amplitude.com":        true,
		"mixpanel.com":         true,
		"google-analytics.com": true,
		"googletagmanager.com": true,
		"analytics.google.com": true,
		"datadoghq.com":        true,
		"datadoghq.eu":         true,
		"newrelic.com":         true,
		"rollbar.com":          true,
		"sentry.io":            true,
		"bugsnag.com":          true,
		"honeybadger.io":       true,
		"logentries.com":       true,
		"loggly.com":           true,
		"papertrailapp.com":    true,
		// Payment processing
		"stripe.com":            true,
		"paypal.com":            true,
		"braintreepayments.com": true,
		"authorize.net":         true,
		// Authentication
		"auth0.com":    true,
		"okta.com":     true,
		"onelogin.com": true,
		// Email services
		"sendgrid.com":    true,
		"mailchimp.com":   true,
		"mailgun.com":     true,
		"postmarkapp.com": true,
		// Common APIs
		"api.github.com":    true,
		"github.com":        true,
		"gitlab.com":        true,
		"bitbucket.org":     true,
		"api.twilio.com":    true,
		"slack.com":         true,
		"discord.com":       true,
		"api.openai.com":    true,
		"api.anthropic.com": true,
		// Local development
		"localhost": true,
		"127.0.0.1": true,
		"0.0.0.0":   true,
		"::1":       true,
	}
}

// getPatterns returns language-specific patterns for detecting suspicious network operations.
func (c *NetworkCheck) getPatterns() map[string][]*regexp.Regexp {
	patterns := make(map[string][]*regexp.Regexp)

	// Go patterns
	patterns["go"] = []*regexp.Regexp{
		// HTTP requests with variable URLs
		regexp.MustCompile(`http\.(Get|Post|Put|Delete|Patch)\s*\(\s*[a-zA-Z_]\w*\s*[\+,]`),
		regexp.MustCompile(`http\.NewRequest\s*\(\s*["'][A-Z]+["'],\s*[a-zA-Z_]\w*\s*[\+,]`),
		// Client requests
		regexp.MustCompile(`&?http\.Client\s*`),
		regexp.MustCompile(`\.Do\s*\(`),
		// Suspicious encoding before send
		regexp.MustCompile(`base64\.(StdEncoding|RawStdEncoding)\.EncodeToString\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`hex\.EncodeToString\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`json\.Marshal\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		// External URLs
		regexp.MustCompile(`["']https?://[^"']+`),
	}

	// Python patterns
	patterns["python"] = []*regexp.Regexp{
		// HTTP requests with variables
		regexp.MustCompile(`requests\.(get|post|put|delete|patch|request)\s*\(\s*[a-zA-Z_]\w*\s*[\$,=]`),
		regexp.MustCompile(`urllib\.(request\.)?(.urlopen|Request)\s*\(\s*[a-zA-Z_]\w*\s*[\$,]`),
		regexp.MustCompile(`urllib2\.urlopen\s*\(\s*[a-zA-Z_]\w*\s*[\$,]`),
		regexp.MustCompile(`http\.client\s*\(`),
		// Suspicious encoding before send
		regexp.MustCompile(`base64\.(b64encode|standard_b64encode)\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`binascii\.hexlify\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`json\.dumps\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`pickle\.(dumps|dump)\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		// External URLs
		regexp.MustCompile(`["']https?://[^"']+`),
		// Socket connections
		regexp.MustCompile(`socket\.(connect|send|sendall)\s*\(\s*[a-zA-Z_]\w*\s*[\$,]`),
	}

	// Node/JavaScript patterns
	patterns["node"] = []*regexp.Regexp{
		// HTTP requests with variables
		regexp.MustCompile(`fetch\s*\(\s*[a-zA-Z_$]\w*\s*[\+,]`),
		regexp.MustCompile(`axios\.(get|post|put|delete|patch)\s*\(\s*[a-zA-Z_$]\w*\s*[\+,]`),
		regexp.MustCompile(`http\.(get|post|request)\s*\(\s*[a-zA-Z_$]\w*\s*[\+,]`),
		regexp.MustCompile(`https\.(get|post|request)\s*\(\s*[a-zA-Z_$]\w*\s*[\+,]`),
		// Suspicious encoding before send
		regexp.MustCompile(`btoa\s*\(\s*[a-zA-Z_$]\w*\s*\)`),
		regexp.MustCompile(`Buffer\.from\s*\(\s*[a-zA-Z_$]\w*\s*\)\.toString\s*\(\s*["']base64`),
		regexp.MustCompile(`JSON\.stringify\s*\(\s*[a-zA-Z_$]\w*\s*\)`),
		// External URLs
		regexp.MustCompile(`["']https?://[^"']+`),
		// WebSockets
		regexp.MustCompile(`new WebSocket\s*\(\s*[a-zA-Z_$]\w*\s*[\+,]`),
	}

	// TypeScript patterns (same as Node)
	patterns["typescript"] = patterns["node"]

	// Java patterns
	patterns["java"] = []*regexp.Regexp{
		// HTTP requests with variables
		regexp.MustCompile(`HttpURLConnection\s*.*\.connect\s*\(`),
		regexp.MustCompile(`HttpClient\s*(.*\.send\s*\()`),
		regexp.MustCompile(`RestTemplate\.(exchange|getFor|postFor)\s*\(\s*[a-zA-Z_]\w*\s*[\$,]`),
		regexp.MustCompile(`OkHttpClient\s*.*\.execute\s*\(`),
		// Suspicious encoding
		regexp.MustCompile(`Base64\.getEncoder\(\)\.encodeToString\s*\(`),
		regexp.MustCompile(` DatatypeConverter\.printBase64Binary\s*\(`),
		regexp.MustCompile(`JSON\.toJSONString\s*\(`),
		// External URLs
		regexp.MustCompile(`["']https?://[^"']+`),
	}

	// Ruby patterns
	patterns["ruby"] = []*regexp.Regexp{
		// HTTP requests with variables
		regexp.MustCompile(`Net::HTTP\.(get|post|put|delete|patch)\s*\(\s*[a-zA-Z_]\w*\s*`),
		regexp.MustCompile(`\.open\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`request\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		// Suspicious encoding
		regexp.MustCompile(`Base64\.encode64\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`\[.*\]\.pack\s*\(`),
		regexp.MustCompile(`JSON\.generate\s*\(\s*[a-zA-Z_]\w*\s*\)`),
		// External URLs
		regexp.MustCompile(`["']https?://[^"']+`),
		// Mechanize
		regexp.MustCompile(`Mechanize\s*new`),
	}

	// PHP patterns
	patterns["php"] = []*regexp.Regexp{
		// HTTP requests with variables
		regexp.MustCompile(`curl_(exec|init|_setopt)\s*\(\s*\$[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`file_get_contents\s*\(\s*\$[a-zA-Z_]\w*\s*[\$,]`),
		regexp.MustCompile(`fopen\s*\(\s*["']https?://.*\$[a-zA-Z_]\w*`),
		// Suspicious encoding
		regexp.MustCompile(`base64_encode\s*\(\s*\$[a-zA-Z_]\w*\s*\)`),
		regexp.MustCompile(`json_encode\s*\(\s*\$[a-zA-Z_]\w*\s*\)`),
		// External URLs
		regexp.MustCompile(`["']https?://[^"']+`),
	}

	// Rust patterns
	patterns["rust"] = []*regexp.Regexp{
		regexp.MustCompile(`reqwest::(get|post|Client::new).*\.send\s*\(\)`),
		regexp.MustCompile(`ureq::(get|post)\s*\(\s*&?[a-zA-Z_]\w*`),
		regexp.MustCompile(`attohttpc\s*.*\.send\s*\(\)`),
		regexp.MustCompile(`["']https?://[^"']+`),
	}

	// C/C++ patterns
	patterns["c"] = []*regexp.Regexp{
		regexp.MustCompile(`curl_easy_(setopt|perform)\s*\(`),
		regexp.MustCompile(`socket\s*\(\s*AF_INET`),
		regexp.MustCompile(`connect\s*\(`),
		regexp.MustCompile(`send\s*\(\s*[a-zA-Z_]\w*\s*,\s*buffer`),
	}
	patterns["cpp"] = patterns["c"]

	// Swift patterns
	patterns["swift"] = []*regexp.Regexp{
		regexp.MustCompile(`URLSession\.shared\.(data|download)\s*\(\s*.*url:`),
		regexp.MustCompile(`URLRequest\s*\(url:\s*[a-zA-Z_]\w*`),
		regexp.MustCompile(`["']https?://[^"']+`),
	}

	return patterns
}

// scanDirectory scans all files in directory for suspicious network patterns.
func (c *NetworkCheck) scanDirectory(path string) []Finding {
	var findings []Finding

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories
		if info.IsDir() {
			if skipDirectories[info.Name()] {
				return filepath.SkipDir
			}
			if strings.HasPrefix(info.Name(), ".") && info.Name() != ".github" && info.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if this is a source file we should scan
		if !isSourceFile(filePath) {
			return nil
		}

		// Skip test files, examples, templates
		if isTestFile(filepath.Base(filePath)) || shouldSkipFile(filepath.Base(filePath)) {
			return nil
		}

		// Detect language from file extension
		language := detectLanguageFromPath(filePath)

		// Get patterns for this language
		patterns, ok := c.Patterns[language]
		if !ok || len(patterns) == 0 {
			return nil
		}

		// Scan file
		fileFindings := c.scanFile(path, filePath, language, patterns)
		findings = append(findings, fileFindings...)

		// Limit findings to avoid excessive output
		if len(findings) >= 50 {
			return filepath.SkipAll
		}

		return nil
	})

	_ = err
	return findings
}

// scanFile scans a single file for suspicious network patterns.
func (c *NetworkCheck) scanFile(root, filePath string, language string, patterns []*regexp.Regexp) []Finding {
	var findings []Finding

	file, err := safepath.OpenPath(root, filePath)
	if err != nil {
		return findings
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	// Get relative path for reporting
	relPath, err := filepath.Rel(root, filePath)
	if err != nil {
		relPath = filepath.Base(filePath)
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip comment lines
		if isCommentLine(line, language) {
			continue
		}

		// Check each pattern
		for _, pattern := range patterns {
			if pattern.MatchString(line) {
				// For URL patterns, check if domain is safe
				if c.isSafeURL(line) {
					continue
				}

				// Extract matched pattern for better reporting
				match := pattern.FindString(line)
				findings = append(findings, Finding{
					Type:        "network",
					File:        relPath,
					Line:        lineNum,
					Description: fmt.Sprintf("suspicious network operation: %s", c.sanitizeMatch(match)),
					Severity:    "medium",
				})
				break // One finding per line
			}
		}
	}

	return findings
}

// isSafeURL checks if a URL in the line is from a known safe domain.
func (c *NetworkCheck) isSafeURL(line string) bool {
	// Extract URLs from the line
	urlRegex := regexp.MustCompile(`https?://([a-zA-Z0-9.-]+)`)
	matches := urlRegex.FindAllStringSubmatch(line, -1)

	for _, match := range matches {
		if len(match) > 1 {
			domain := match[1]
			// Check main domain and parent domain
			if c.SafeDomains[domain] {
				return true
			}
			// Check parent domain
			parts := strings.Split(domain, ".")
			if len(parts) >= 2 {
				parentDomain := strings.Join(parts[len(parts)-2:], ".")
				if c.SafeDomains[parentDomain] {
					return true
				}
			}
		}
	}

	// Also try parsing the URL properly
	urlRegex2 := regexp.MustCompile(`https?://[^\s"']+"`)
	urls := urlRegex2.FindAllString(line, -1)
	for _, urlStr := range urls {
		u, err := url.Parse(urlStr)
		if err == nil && u.Host != "" {
			if c.SafeDomains[u.Host] {
				return true
			}
			// Check parent domain
			parts := strings.Split(u.Host, ".")
			if len(parts) >= 2 {
				parentDomain := strings.Join(parts[len(parts)-2:], ".")
				if c.SafeDomains[parentDomain] {
					return true
				}
			}
		}
	}

	return false
}

// sanitizeMatch sanitizes a matched pattern for display.
func (c *NetworkCheck) sanitizeMatch(match string) string {
	// Truncate very long matches
	if len(match) > 100 {
		return match[:100] + "..."
	}
	return match
}

// formatFindings formats findings into a readable message.
func (c *NetworkCheck) formatFindings(findings []Finding) string {
	if len(findings) == 1 {
		return fmt.Sprintf("Suspicious network operation detected: %s", findings[0].String())
	}

	if len(findings) <= 3 {
		return fmt.Sprintf("Suspicious network operations detected: %s", c.joinFindings(findings))
	}

	return fmt.Sprintf("%d suspicious network operations detected (e.g., %s)", len(findings),
		c.joinFindings(findings[:3]))
}

// joinFindings joins findings for display.
func (c *NetworkCheck) joinFindings(findings []Finding) string {
	var strs []string
	for _, f := range findings {
		strs = append(strs, f.String())
	}
	return strings.Join(strs, ", ")
}
