package main

import "github.com/spf13/cobra"

var (
	generateCoverCmdPreset   string
	generateCoverCmdArticle  string
	generateCoverCmdTitle    string
	generateCoverCmdSummary  string
	generateCoverCmdKeywords string
	generateCoverCmdStyle    string
	generateCoverCmdAspect   string
	generateCoverCmdSize     string

	generateInfographicCmdPreset   string
	generateInfographicCmdArticle  string
	generateInfographicCmdTitle    string
	generateInfographicCmdSummary  string
	generateInfographicCmdKeywords string
	generateInfographicCmdStyle    string
	generateInfographicCmdAspect   string
	generateInfographicCmdSize     string
)

var generateCoverCmd = &cobra.Command{
	Use:   "generate_cover",
	Short: "Generate an article cover image from a preset",
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGeneratePresetImage("cover", "cover-default", generateImageInput{
			Preset:   generateCoverCmdPreset,
			Article:  generateCoverCmdArticle,
			Title:    generateCoverCmdTitle,
			Summary:  generateCoverCmdSummary,
			Keywords: generateCoverCmdKeywords,
			Style:    generateCoverCmdStyle,
			Aspect:   generateCoverCmdAspect,
			Size:     generateCoverCmdSize,
		})
	},
}

var generateInfographicCmd = &cobra.Command{
	Use:   "generate_infographic",
	Short: "Generate an infographic image from a preset",
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGeneratePresetImage("infographic", "infographic-default", generateImageInput{
			Preset:   generateInfographicCmdPreset,
			Article:  generateInfographicCmdArticle,
			Title:    generateInfographicCmdTitle,
			Summary:  generateInfographicCmdSummary,
			Keywords: generateInfographicCmdKeywords,
			Style:    generateInfographicCmdStyle,
			Aspect:   generateInfographicCmdAspect,
			Size:     generateInfographicCmdSize,
		})
	},
}

func init() {
	addPresetImageFlags(generateCoverCmd, &generateCoverCmdPreset, &generateCoverCmdArticle, &generateCoverCmdTitle, &generateCoverCmdSummary, &generateCoverCmdKeywords, &generateCoverCmdStyle, &generateCoverCmdAspect, &generateCoverCmdSize)
	addPresetImageFlags(generateInfographicCmd, &generateInfographicCmdPreset, &generateInfographicCmdArticle, &generateInfographicCmdTitle, &generateInfographicCmdSummary, &generateInfographicCmdKeywords, &generateInfographicCmdStyle, &generateInfographicCmdAspect, &generateInfographicCmdSize)
}

func addPresetImageFlags(cmd *cobra.Command, preset, article, title, summary, keywords, style, aspect, size *string) {
	cmd.Flags().StringVar(preset, "preset", "", "Prompt preset from the image prompt catalog")
	cmd.Flags().StringVarP(article, "article", "a", "", "Article markdown file used to render a preset prompt")
	cmd.Flags().StringVar(title, "title", "", "Article title used to render a preset prompt")
	cmd.Flags().StringVar(summary, "summary", "", "Article summary used to render a preset prompt")
	cmd.Flags().StringVar(keywords, "keywords", "", "Keywords used to render a preset prompt")
	cmd.Flags().StringVar(style, "style", "", "Visual style used to render a preset prompt")
	cmd.Flags().StringVar(aspect, "aspect", "", "Aspect ratio hint used to render a preset prompt, e.g. 16:9 or 3:4")
	cmd.Flags().StringVarP(size, "size", "s", "", "Image size (e.g., 2560x1440 for 16:9)")
}
