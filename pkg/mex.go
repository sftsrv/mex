package pkg

import (
	"fmt"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

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
	blockRe := regexp.MustCompile("`````\\w+ ")

	blocks := blockRe.Split(content, -1)[1:]

	regions := make([]Region, len(blocks))

	for b := range blocks {
		block := slices.Collect(strings.Lines(strings.TrimSpace(blocks[b])))

		region, err := parseMdBlock(block)
		if err != nil {
			return regions, err
		}

		regions[b] = region

	}

	return regions, nil
}

func parseMdBlock(lines []string) (Region, error) {
	meta := strings.TrimSpace(lines[0])
	content := strings.Join(lines[1:(len(lines)-1)], "")

	pathRange := strings.SplitN(meta, ":", 2)
	path := pathRange[0]

	fmt.Println(pathRange)
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

		ext := path.Ext(region.Path)
		md += fmt.Sprintf("`````%s %s:%d-%d\n%s\n`````", ext, region.Path, region.Start, region.End, region.Content)
	}

	return md
}
