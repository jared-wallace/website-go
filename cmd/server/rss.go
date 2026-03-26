package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// RSS structures for parsing BBC World News feed
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
}

// Enhanced news handler with actual RSS feed fetching
func (s *Server) newsHandler(w http.ResponseWriter, r *http.Request) {
	// BBC World News RSS feed
	feedURL := "http://feeds.bbci.co.uk/news/world/rss.xml"

	news, err := fetchRSSFeed(feedURL)
	if err != nil {
		// Fallback to sample data if RSS fetch fails
		news = []NewsItem{
			{
				Title:       "Unable to fetch latest news",
				Link:        "https://bbc.com/news",
				Description: "Please check BBC News directly for the latest updates",
				PubDate:     time.Now().Format(time.RFC3339),
			},
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300") // Cache for 5 minutes
	json.NewEncoder(w).Encode(news)
}

func fetchRSSFeed(feedURL string) ([]NewsItem, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(feedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RSS feed returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read RSS response: %w", err)
	}

	var rss RSS
	if err := xml.Unmarshal(body, &rss); err != nil {
		return nil, fmt.Errorf("failed to parse RSS XML: %w", err)
	}

	var newsItems []NewsItem
	for _, item := range rss.Channel.Items {
		// Clean up description (remove HTML tags for ticker)
		description := strings.ReplaceAll(item.Description, "<p>", "")
		description = strings.ReplaceAll(description, "</p>", "")
		description = strings.ReplaceAll(description, "<br>", " ")
		description = strings.ReplaceAll(description, "<br/>", " ")

		// Limit description length for ticker
		if len(description) > 150 {
			description = description[:147] + "..."
		}

		newsItem := NewsItem{
			Title:       item.Title,
			Link:        item.Link,
			Description: description,
			PubDate:     item.PubDate,
		}
		newsItems = append(newsItems, newsItem)

		// Limit to first 10 items for performance
		if len(newsItems) >= 10 {
			break
		}
	}

	return newsItems, nil
}
