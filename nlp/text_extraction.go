package nlp

import (
	database "backend-crm/internal/database/mongodb"
	"backend-crm/pkg/enum"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"

	"time"
)

// returns ETA and status extracted from the master email content
// "subject": "UOG IOANNIS V-202402 / ST Shipping / 8 FEB 2024 / Order + Clean // 12 Hrs ETA Notice",
// "bodycontent": "Dear Sir, Good Day,\n\n\n\nETA Singapore (PEBGC): 0600 Lt/26th Feb'24 AGW WSNP\n\n\n\nAgents to kindly inform all
// concerned & advise berthing prospects.\n\n\n\nThanks and Best Regards,\n\n\n\nCapt. Devendra
// Joshi\n\n\n\nMaster\n\n\n\nMT UOG IOANNIS V\n\n\n\nMarshal Islan / V898793/9937938",

// "subject": "UOG IOANNIS V-202402 / ST Shipping / 8 FEB 2024 / Order + Clean // 12 Hrs ETA Notice",
// "body_content": "ETA Next Calling Dear Sir Speed: (PURONG TERMINAL): 21 Feb 2024/ 12:00 \n\n\n\nAgents to kindly inform all concerned & advise berthing prospects.\n\n\n\nThanks and Best
// Regards,\n\n\n\nCapt. Devendra Joshi\n\n\n\nMaster\n\n\n\nM",

// "subject": "UOG IOANNIS V-202402 / ST Shipping / 8 FEB 2024 / Order + Clean // NOR Tendered - Singapore"
// "bodycontent": "Dear Sir, \nGood Day.\nPls find attached NOR tendered.\nAgent to kindly ack receipt & inform all concerned.\nThanks and Best Regards,
// Capt. Devendra Joshi\n\nMaster\nMT UOG IOANNIS V\nMarshal Island / V898793/9937938

func TranslateMasterEmailContentToStatusUpdate(subject string, content string, currentShipmentResponse database.ShipmentResponse) (time.Time, enum.ShipmentStatus) {
	var status enum.ShipmentStatus

	// Define regex patterns for extracting ETA and status
	// text := "ETA Singapore (PEBGC): 0600 Lt/26th Feb'24 AGW WSNP"
	// etaPattern := regexp.MustCompile(`ETA.*?: (\d{1,2}:\d{2}|\d{4}).*?(\d{1,2})(?:st|nd|rd|th)? (\w{3})'(\d{2})`)
	// etaPattern := regexp.MustCompile(`ETA.*?: (\d{1,2}:\d{2}|\d{4}).*?(\d{1,2})(?:st|nd|rd|th)? (\w{3})'(\d{2})|(\d{1,2} \w{3} \d{4})/ (\d{1,2}:\d{2})`)
	etaPattern := regexp.MustCompile(`ETA`)
	// etaPattern := regexp.MustCompile(`ETA.*?: (\d{1,2}:\d{2}|\d{4}) .*?(\d{1,2})(?:st|nd|rd|th)? (\w{3}) ?'?(20\d{2}|\d{2})?`)
	// etaPattern := regexp.MustCompile(`ETA.*?: (\d{1,2}:\d{2}|\d{4})[^\d]*(\d{1,2})(?:st|nd|rd|th)? (\w{3})[^\d]*(\d{2,4})?`)

	// etaPattern := regexp.MustCompile(`ETA.*: (\d{4}) .*?(\d{1,2}(?:st|nd|rd|th)? \w{3}'\d{2})`)
	eospPattern := regexp.MustCompile(`EOSP`)
	anchoragePattern := regexp.MustCompile(`ANCHOR\w*|ANCHORAGE\w*`)
	norTenderedPattern := regexp.MustCompile(`NOR TENDERED`)
	norReTenderedPattern := regexp.MustCompile(`NOR RE-TENDERED`)
	berthedPattern := regexp.MustCompile(`BERTH\w*|ALL FAST`)
	activityCommencedPattern := regexp.MustCompile(`COMMENCE\w*`)
	activityCompletePattern := regexp.MustCompile(`COMPLETED\w*`)
	cospPattern := regexp.MustCompile(`COSP\w*`)
	// Check if the subject or content matches specific status patterns
	// According to priority of statuses, meaning if there is "COMPLETED" and "NOR TENDERED",
	// obviously, the current status = "COMPLETED"
	// output = "COMPLETED"
	switch {
	case cospPattern.MatchString(content):
		status = enum.SHIPMENT_STATUS_COSP
	case activityCompletePattern.MatchString(content):
		status = enum.SHIPMENT_STATUS_ACTIVITY_COMPLETED
	case activityCommencedPattern.MatchString(content):
		status = enum.SHIPMENT_STATUS_ACTIVITY_COMMENCED
	case berthedPattern.MatchString(content):
		status = enum.SHIPMENT_STATUS_BERTHED
	case norReTenderedPattern.MatchString(content):
		status = enum.SHIPMENT_STATUS_NOR_RETENDERED
	case norTenderedPattern.MatchString(content):
		status = enum.SHIPMENT_STATUS_NOR_TENDERED
	case anchoragePattern.MatchString(content):
		status = enum.SHIPMENT_STATUS_AT_ANCHORAGE
	case eospPattern.MatchString(content):
		status = enum.SHIPMENT_STATUS_EOSP
	case etaPattern.MatchString(content):
		status = enum.SHIPMENT_STATUS_EN_ROUTE

	default:
		status = enum.SHIPMENT_STATUS_NOT_STARTED
	}
	log.Println(status, "status")
	if status == enum.SHIPMENT_STATUS_EN_ROUTE {
		// Define the regex to match exactly 4 digits
		timeRegex := regexp.MustCompile(`\b\d{4}\b`)

		// Define the regex to match date suffixes (e.g., "1st", "2nd", "3rd", "4th")
		dateSuffixRegex := regexp.MustCompile(`\b(\d{1,2})(st|nd|rd|th)\b`)

		// Remove the date suffixes by replacing them with just the digits
		contentWithoutDateSuffix := dateSuffixRegex.ReplaceAllString(content, `$1`)

		// Replace the matched time with the formatted time "HH:MM"
		updatedContent := timeRegex.ReplaceAllStringFunc(contentWithoutDateSuffix, func(match string) string {
			// Check if this match is already formatted with a colon (i.e., "HH:MM")
			if len(match) == 4 && !strings.Contains(match, ":") && isLikelyTime(match, contentWithoutDateSuffix) {
				// Insert colon to format the time
				log.Println(match)
				return match[:2] + ":" + match[2:]
			}
			return match
		})
		w := when.New(nil)
		w.Add(en.All...)
		w.Add(common.All...)

		systemTimeNow := time.Now()
		systemTimeNowLocation := systemTimeNow.Location()
		// Check if the location is not "Local"
		if systemTimeNowLocation.String() != "Local" {
			// Define a fixed time zone with +8 hours offset from UTC
			systemTimeNowLocation = time.FixedZone("UTC+8", 8*60*60)

		}
		fmt.Println(updatedContent, "updatedcontent")
		currentTime := time.Now().In(systemTimeNowLocation)
		r, err := w.Parse(updatedContent, currentTime)

		if err != nil {
			fmt.Println("Error parsing date:", err)
			return currentShipmentResponse.CurrentETA, currentShipmentResponse.CurrentStatus
		}

		if r == nil {
			fmt.Println("No date found")
			return currentShipmentResponse.CurrentETA, currentShipmentResponse.CurrentStatus
		}

		fmt.Println("Parsed date:", r.Time)
		eta := r.Time.In(systemTimeNowLocation)
		fmt.Println("Constructed ETA:", eta)
		return eta, status
	}

	return currentShipmentResponse.CurrentETA, status
}

// isLikelyTime checks if the matched string is more likely a time than a year or part of a date.
func isLikelyTime(match, originalText string) bool {
	// Regex to detect a date pattern without ordinal suffixes (e.g., "21 Feb 2024")
	dateRegex1 := regexp.MustCompile(`\d{1,2} \w{3} \d{4}`)

	// Regex to detect a date pattern with ordinal suffixes (e.g., "21st Feb 2024" or "2nd Jan 2024")
	// dateRegex2 := regexp.MustCompile(`\d{1,2}(st|nd|rd|th) \w{3} \d{4}`)

	// Regex to detect a date pattern with slashes (e.g., "2/12/2024" or "12/2/2024")
	dateRegex3 := regexp.MustCompile(`\b\d{1,2}/\d{1,2}/\d{4}\b`)

	// Find all date matches in the text
	dateMatches1 := dateRegex1.FindAllString(originalText, -1)
	// dateMatches2 := dateRegex2.FindAllString(originalText, -1)
	dateMatches3 := dateRegex3.FindAllString(originalText, -1)

	// Combine all types of matches into a single slice
	dateMatches := append(dateMatches1, dateMatches3...)
	// dateMatches = append(dateMatches, ...)

	for _, dateMatch := range dateMatches {
		// If the match is part of a date (e.g., "21 Feb 2024" or "2/12/2024"), treat it as a year, not a time
		if strings.Contains(dateMatch, match) {
			log.Println(match, "false")
			return false
		}
	}

	// If there's no obvious date pattern nearby, assume the four digits represent a time
	log.Println(match, "true")
	return true
}
