package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageURLForNewImageBlock(t *testing.T) {
	imageText := NewTextBlockObject("plain_text", "Location", false, false)
	imageBlock := NewImageBlock("https://api.slack.com/img/blocks/bkb_template_images/tripAgentLocationMarker.png", "Marker", "test", imageText)

	assert.Equal(t, imageBlock.BlockType(), MBTImage)
	assert.Equal(t, string(imageBlock.Type), "image")
	assert.Equal(t, imageBlock.Title.Type, "plain_text")
	assert.Equal(t, imageBlock.ID(), "test")
	assert.Equal(t, imageBlock.BlockID, "test")
	assert.Contains(t, imageBlock.Title.Text, "Location")
	assert.Contains(t, imageBlock.ImageURL, "tripAgentLocationMarker.png")
}

func TestSlackFileForNewImageBlock(t *testing.T) {
	imageText := NewTextBlockObject("plain_text", "Location", false, false)
	slackFile := &SlackFileObject{URL: "https://api.slack.com/img/blocks/bkb_template_images/tripAgentLocationMarker.png"}
	imageBlock := NewImageBlockSlackFile(slackFile, "Marker", "test", imageText)

	assert.Equal(t, imageBlock.BlockType(), MBTImage)
	assert.Equal(t, string(imageBlock.Type), "image")
	assert.Equal(t, imageBlock.Title.Type, "plain_text")
	assert.Equal(t, imageBlock.ID(), "test")
	assert.Equal(t, imageBlock.BlockID, "test")
	assert.Contains(t, imageBlock.Title.Text, "Location")
	assert.Contains(t, imageBlock.SlackFile.URL, "tripAgentLocationMarker.png")
}
