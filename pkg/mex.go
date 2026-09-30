package pkg

import (
	"fmt"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Line numbers are 1-indexed as this is convention
type Region struct {
	Path    string
	Start   int
	End     int
	Content string
}

func FromGrep(content string) ([]Region, error) {
	lineRe := regexp.MustCompile(`(.+):(\d+):(.*)`)

	lines := strings.Lines(strings.TrimSpace(content))

	regions := []Region{}
	for line := range lines {
		current, err := parseGrepLine(lineRe, line)

		if err != nil {
			return regions, err
		}

		if len(regions) == 0 {
			regions = append(regions, current)
			continue
		}

		last := regions[len(regions)-1]

		if last.Path == current.Path && last.End == current.Start-1 {
			// Extend the current block since it's immediately after us
			last.Content += current.Content
			last.End++

			regions[len(regions)-1] = last
		} else {
			// Otherwise, append the new block
			regions = append(regions, current)
		}

	}

	return regions, nil
}

func parseGrepLine(lineRe *regexp.Regexp, line string) (Region, error) {
	parts := lineRe.FindStringSubmatch(line)

	if len(parts) < 4 {
		return Region{}, fmt.Errorf("Unexpected metadata line: %s\n\nMake sure you run grep with the -H and -n flags to include the file name and line number", line)
	}

	path := parts[1]
	num, err := strconv.Atoi(parts[2])
	content := parts[3] + "\n"

	region := Region{
		Path:    path,
		Start:   num,
		End:     num,
		Content: content,
	}

	if err != nil {
		return region, fmt.Errorf("Invalid line number %s at: %s", parts[2], line)
	}

	return region, nil

}

func FromMd(content string) ([]Region, error) {
	startRe := regexp.MustCompile("`````\\w+ ")
	endSep := "\n`````"

	blocks := startRe.Split(content, -1)[1:]

	regions := make([]Region, len(blocks))

	for b := range blocks {
		block := strings.Split(blocks[b], endSep)
		if len(block) < 1 {
			continue
		}

		region, err := parseMdBlock(block[0] + endSep)
		if err != nil {
			return regions, err
		}

		regions[b] = region

	}

	return regions, nil
}

func parseMdBlock(block string) (Region, error) {
	lines := slices.Collect(strings.Lines(block))

	meta := strings.TrimSpace(lines[0])
	content := strings.Join(lines[1:(len(lines)-1)], "")

	pathRange := strings.SplitN(meta, ":", 2)
	path := pathRange[0]

	ranges := strings.SplitN(pathRange[1], "-", 2)
	start, err := strconv.Atoi(ranges[0])

	if err != nil {
		return Region{}, fmt.Errorf("Invalid start line at: %s", meta)
	}

	end, err := strconv.Atoi(ranges[1])
	if err != nil {
		return Region{}, fmt.Errorf("Invalid end line at: %s", meta)
	}

	region := Region{
		Path:    path,
		Start:   start,
		End:     end,
		Content: content,
	}

	return region, nil
}

func ToMdBlocks(regions []Region) string {
	md := ""

	for r := range regions {
		region := regions[r]

		ext := strings.Replace(path.Ext(region.Path), ".", "", 1)
		md += fmt.Sprintf("`````%s %s:%d-%d\n%s`````\n\n", ext, region.Path, region.Start, region.End, region.Content)
	}

	return md
}

func ToMdFile(title string, note string, regions []Region) string {
	return fmt.Sprintf("# %s\n\n> %s\n\n%s", title, note, ToMdBlocks(regions))
}

func ApplyRegions(content string, regions []Region) string {
	grouped := map[int]Region{}

	for r := range regions {
		region := regions[r]
		grouped[region.Start] = region
	}

	lines := slices.Collect(strings.Lines(content))
	output := ""

	i := 0
	for i < len(lines) {
		n := i + 1

		region, exists := grouped[n]

		if exists {
			output += region.Content
			i = region.End
		} else {
			output += lines[i]
			i++
		}

	}

	return output
}
