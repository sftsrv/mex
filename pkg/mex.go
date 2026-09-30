package pkg

import (
	"fmt"
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

func FromGrep(content string) []Region {
	lineRe := regexp.MustCompile(`(.+):(\d+):(.*)`)

	lines := strings.Lines(strings.TrimSpace(content))

	regions := []Region{}
	for line := range lines {
		parts := lineRe.FindStringSubmatch(line)
		fmt.Println(parts)

		path := parts[1]
		num, err := strconv.Atoi(parts[2])

		if err != nil {
			panic("Invalid line number " + parts[2] + " at: " + line)

		}

		content := parts[3]

		newBlock := Region{
			Path:    path,
			Start:   num,
			End:     num,
			Content: content,
		}

		if len(regions) == 0 {
			regions = append(regions, newBlock)
			continue
		}

		last := regions[len(regions)-1]

		if last.Path == path && last.End == num-1 {
			// Extend the current block since it's immediately after us
			last.Content += "\\n" + content
			last.End++

			regions[len(regions)-1] = last
		} else {
			// Otherwise, append the new block
			regions = append(regions, newBlock)
		}

	}

	return regions
}

func FromMd(content string) []Region {
	blockRe := regexp.MustCompile("`````\\w+ ")

	blocks := blockRe.Split(content, -1)[1:]

	regions := make([]Region, len(blocks))

	for b := range blocks {
		block := slices.Collect(strings.Lines(blocks[b]))

		meta := strings.TrimSpace(block[0])
		content := strings.Join(block[1:(len(block)-2)], "\n")

		pathRange := strings.SplitN(meta, ":", 2)
		path := pathRange[0]

		fmt.Println(pathRange)
		ranges := strings.SplitN(pathRange[1], "-", 2)
		start, err := strconv.Atoi(ranges[0])

		fmt.Println(ranges)

		if err != nil {
			panic("Invalid start line at: " + meta)
		}

		end, err := strconv.Atoi(ranges[1])
		if err != nil {
			panic("Invalid end line at: " + meta)
		}

		regions[b] = Region{
			Path:    path,
			Start:   start,
			End:     end,
			Content: content,
		}

		fmt.Println(regions[b])

	}

	return regions
}
